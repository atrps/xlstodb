package main

import (
	"fmt"
	"log"

	"github.com/atrps/xlstodb"
)

func main() {

	// db connection parameters
	dtCfgConnect, err := xlstodb.ReadCfgConnect("connectprm.json")
	if err != nil {
		log.Fatal(err)
	}

	// example load svc
	dtCfgXls1, err := xlstodb.ReadCfgXls("svc1.json")
	if err != nil {
		log.Fatal(err)
	}
	// example load pay
	dtCfgXls2, err := xlstodb.ReadCfgXls("pay2.json")
	if err != nil {
		log.Fatal(err)
	}

	// db
	// append svc
	cntAppend, cntSkip, id, err := xlstodb.AppendRec(dtCfgXls1, dtCfgConnect)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("*** append svc records ", cntAppend, ", skip records ", cntSkip)
	fmt.Println("*** last appended rec with id=", id)

	// append pay
	cntAppend, cntSkip, id, err = xlstodb.AppendRec(dtCfgXls2, dtCfgConnect)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("*** append pay records ", cntAppend, ", skip records ", cntSkip)
	fmt.Println("*** last appended rec with id=", id)

}
