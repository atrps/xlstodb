// load data from Excel to table PostgreSQL
package xlstodb

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"

	"database/sql"

	_ "github.com/lib/pq"
)

// import settings Excel file
type CfgXlsToDb = struct {
	Filename  string      `json:"filename"`
	Sheet     string      `json:"sheet"`
	Rowfirst  int         `json:"rowfirst,omitempty"`
	Rowlast   int         `json:"rowlast,omitempty"`
	Tablename string      `json:"tablename"`
	Fields    []FieldItem `json:"fields"`
}

// db connection parameters
type CfgConnectToDb = struct {
	Host     string `json:"host"`
	Dbname   string `json:"dbname"`
	User     string `json:"user"`
	Password string `json:"password"`
}

// fields of table
type FieldItem = struct {
	Name       string `json:"name"`
	Place      int    `json:"place"`
	Formatdate string `json:"formatdate,omitempty"`
}

// reading import settings Excel file
func ReadCfgXls(fjsname string) (CfgXlsToDb, error) {
	var cfg CfgXlsToDb

	jsFile, err := os.Open(fjsname)
	if err != nil {
		return cfg, err
	}
	defer jsFile.Close()

	rd := bufio.NewReader(jsFile)
	data, err := ioutil.ReadAll(rd)
	if err != nil {
		fmt.Println(err)
		return cfg, err
	}

	if err := json.Unmarshal(data, &cfg); err != nil {
		fmt.Println(err)
		return cfg, err
	}
	return cfg, nil
}

// reading db connection parameters
func ReadCfgConnect(fjsname string) (CfgConnectToDb, error) {
	var cfg CfgConnectToDb

	jsFile, err := os.Open(fjsname)
	if err != nil {
		return cfg, err
	}
	defer jsFile.Close()

	rd := bufio.NewReader(jsFile)
	data, err := ioutil.ReadAll(rd)
	if err != nil {
		fmt.Println(err)
		return cfg, err
	}

	if err := json.Unmarshal(data, &cfg); err != nil {
		fmt.Println(err)
		return cfg, err
	}
	return cfg, nil
}

// parse string of date to format YYYY-MM-DD
func strToYMD(s, f string) (string, error) {
	var st string
	var ft string
	var dt time.Time

	// convert from format DD.MM.YYYY
	if strings.Contains(f, ".") && strings.Contains(s, "/") {
		st = strings.ReplaceAll(s, "/", ".")
	} else {
		st = s
	}
	ft = strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(f, "YYYY", "2006"), "MM", "01"), "DD", "02")
	dt, err := time.Parse(ft, st)
	if err != nil {
		return "", err
	}
	// to format YYYY-MM-DD
	st = dt.Format("2006-01-02")
	return st, nil
}

// append records to table on settings
func AppendRec(cfgXls CfgXlsToDb, cfgConnect CfgConnectToDb) (cntAppend int, cntSkip int, lastId int64, err error) {
	var r1, r2 int
	var listNames, listHolders string
	var cntFields int
	var flSkip bool

	// open Excel file
	r1, r2 = cfgXls.Rowfirst, cfgXls.Rowlast
	f, err := excelize.OpenFile(cfgXls.Filename)
	defer func() {
		// Close Excel file
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	if err != nil {
		return
	}

	// get list fields
	cntFields = len(cfgXls.Fields)
	for i, field := range cfgXls.Fields {
		listNames += field.Name
		listHolders += "$" + strconv.Itoa(i+1)
		if i < cntFields-1 {
			listNames += ", "
			listHolders += ", "
		}
	}
	// slice of fields value
	listV := make([]any, cntFields)

	// connect to db
	connStr := fmt.Sprintf("host=%s dbname=%s user=%s password=%s sslmode=disable",
		cfgConnect.Host, cfgConnect.Dbname, cfgConnect.User, cfgConnect.Password)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return
	}
	defer db.Close()

	// statement INSERT
	statement := "insert into " + cfgXls.Tablename + " (" + listNames + ") values (" + listHolders + ") returning id"
	stmt, err := db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	// get all the rows in the Sheet
	rows, err := f.GetRows(cfgXls.Sheet)
	if err != nil {
		fmt.Println(err)
		return
	}
	if r1 > 0 && r2 > 0 {
		for _, row := range rows[r1-1 : r2] {
			flSkip = false
			for i, field := range cfgXls.Fields {
				if field.Formatdate != "" {
					// рarse date
					listV[i], err = strToYMD(row[field.Place-1], field.Formatdate)
					if err != nil {
						flSkip = true
						break
					}
				} else {
					listV[i] = row[field.Place-1]
				}
				fmt.Print(listV[i], "\t")
			}

			if flSkip {
				fmt.Println(" -> skipped")
				cntSkip++
			} else {
				// use QueryRow to execute INSERT and scan the returned id
				err = stmt.QueryRow(listV...).Scan(&lastId)
				if err != nil {
					fmt.Println(" -> skipped")
					cntSkip++
				} else {
					fmt.Println(" -> appended with id=", lastId)
					cntAppend++
				}
			}
		}
	}
	return
}
