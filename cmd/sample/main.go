package main

import (
	"fmt"
	"log"

	// "github.com/atrps/xlstodb"

	"xlstodb/pkg/dbconnect"
	"xlstodb/pkg/loader"
)

func main() {

	dbconnect.CfgInit()

	// example load svc
	dtCfgXls1, err := loader.ReadCfgXls(dbconnect.PathCfg + "svc1.json")
	if err != nil {
		log.Fatal(err)
	}
	// example load pay
	dtCfgXls2, err := loader.ReadCfgXls(dbconnect.PathCfg + "pay2.json")
	if err != nil {
		log.Fatal(err)
	}

	// db
	// append svc
	err = loader.CreateBlock(&dtCfgXls1)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = loader.AppendRec(&dtCfgXls1)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("*** load block ", dtCfgXls1.BlockId)
	fmt.Println("*** append svc records ", dtCfgXls1.CntAppend, ", skip records ", dtCfgXls1.CntSkip)
	fmt.Println("*** last appended rec with id=", dtCfgXls1.LastId)

	// append pay
	err = loader.CreateBlock(&dtCfgXls2)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = loader.AppendRec(&dtCfgXls2)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("*** load block ", dtCfgXls2.BlockId)
	fmt.Println("*** append pay records ", dtCfgXls2.CntAppend, ", skip records ", dtCfgXls2.CntSkip)
	fmt.Println("*** last appended rec with id=", dtCfgXls2.LastId)
}
