package dbconnect

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/viper"

	"database/sql"

	_ "github.com/lib/pq"
)

// db connection parameters
type CfgConnectToDb = struct {
	Host     string `json:"host"`
	Dbname   string `json:"dbname"`
	User     string `json:"user"`
	Password string `json:"password"`
}

var cfgConnect CfgConnectToDb
var PathCfg = "./configs/"
var PathPublic = ""
var PathStatic = ""
var PathTmp = ""

// init config to connect
func CfgInit() {

	viper.SetConfigFile(PathCfg + ".env")
	viper.ReadInConfig()

	// db connection parameters
	cfgConnect = CfgConnectToDb{}
	cfgConnect.Host = viper.Get("HOST").(string)
	cfgConnect.Dbname = viper.Get("DBNAME").(string)
	cfgConnect.User = viper.Get("USER").(string)
	cfgConnect.Password = viper.Get("PASSWORD").(string)

	// path
	PathPublic = viper.Get("PATHPUBLIC").(string)
	PathStatic = viper.Get("PATHSTATIC").(string)
	PathTmp = viper.Get("PATHTMP").(string)
}

// connect to db
func Open() (db *sql.DB, err error) {
	connStr := fmt.Sprintf("host=%s dbname=%s user=%s password=%s sslmode=disable",
		cfgConnect.Host, cfgConnect.Dbname, cfgConnect.User, cfgConnect.Password)
	db, err = sql.Open("postgres", connStr)
	return
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
