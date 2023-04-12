// load data from Excel to table PostgreSQL
package loader

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

	"xlstodb/pkg/dbconnect"
)

// import settings Excel file
type CfgXlsToDb = struct {
	Filename  string      `json:"filename"`
	Sheet     string      `json:"sheet"`
	Rowfirst  int         `json:"rowfirst,omitempty"`
	Rowlast   int         `json:"rowlast,omitempty"`
	BlockType string      `json:"blocktype,omitempty"`
	BlockInfo string      `json:"blockinfo,omitempty"`
	Tablename string      `json:"tablename"`
	Fields    []FieldItem `json:"fields"`
	BlockId   int64       `json:"blockid,omitempty"`
	CntAppend int         `json:"cntappend,omitempty"`
	CntSkip   int         `json:"cntskip,omitempty"`
	LastId    int64       `json:"lastid,omitempty"`
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

func CreateBlock(cfgXls *CfgXlsToDb) (err error) {
	var blockId int64
	var cntFields int

	// connect to db
	db, err := dbconnect.Open()
	if err != nil {
		return
	}
	defer db.Close()

	// get Tablename
	if cfgXls.Tablename == "" {
		cfgXls.Tablename = GetBlockTypeTableName(db, cfgXls.BlockType)
		// fmt.Println("Filled Tablename", cfgXls.Tablename)
		if cfgXls.Tablename == "" {
			return
		}
	}
	// check count fields
	cntFields = len(cfgXls.Fields)
	if cntFields == 0 {
		// from block_type_dtl
		cfgXls.Fields = GetBlockTypeDtlFields(db, cfgXls.BlockType)
		// fmt.Printf("Filled Fields in cfgXls %+v\n", cfgXls)
	}
	// block create
	block_statement := "select block_create($1, $2) as block_id"
	if err = db.QueryRow(block_statement, cfgXls.BlockType, cfgXls.BlockInfo).Scan(&blockId); err != nil {
		if err == sql.ErrNoRows {
			return
		}
	}
	cfgXls.BlockId = blockId
	fmt.Println("==create block with id=", blockId)

	return
}

// append records to table on settings
func AppendRec(cfgXls *CfgXlsToDb) (err error) {
	var blockId int64
	var cntAppend int
	var cntSkip int
	var lastId int64
	var r1, r2 int
	var listNames, listHolders string
	var cntFields int
	var flSkip bool

	blockId = cfgXls.BlockId

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

	// connect to db
	db, err := dbconnect.Open()
	if err != nil {
		return
	}
	defer db.Close()

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
	// statement INSERT
	statement := "insert into " + cfgXls.Tablename + " (block_id, " + listNames +
		") values (" + strconv.FormatInt(blockId, 10) + ", " + listHolders + ") returning id"

	fmt.Println("generate statement ", statement)

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
	// in block change status
	_, err = db.Exec("CALL block_change_status($1, $2, $3)", blockId, "P", cntAppend)
	// load
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
					if cntAppend%2 == 0 {
						// in block change count loaded record
						_, err = db.Exec("CALL block_change_status($1, $2, $3)", blockId, "P", cntAppend)
					}
				}
			}
		}
	}

	// block complete
	if err == nil {
		_, err = db.Exec("CALL block_complete($1, $2)", blockId, cntAppend)
	} else {
		_, err = db.Exec("CALL block_change_status($1, $2, $3)", blockId, "E", cntAppend)
	}
	cfgXls.CntAppend = cntAppend
	cfgXls.CntSkip = cntSkip
	cfgXls.LastId = lastId

	// fmt.Println(GetBlockStatus(db, blockId).Info)
	return
}
