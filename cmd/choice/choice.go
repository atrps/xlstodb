package main

import (
	"github.com/labstack/echo/v4"

	"xlstodb/pkg/dbconnect"
	"xlstodb/pkg/handlers"
)

func main() {

	// db connection
	dbconnect.CfgInit()
	db, err := dbconnect.Open()
	if err != nil {
		return
	}
	defer db.Close()

	// create echo
	e := echo.New()

	e.File("/", "web/public/index.html")
	e.Static("/static/", "web/public/")
	e.GET("/blocks", handlers.GetBlocks(db))
	e.GET("/block/:id", handlers.GetBlockOnId(db))
	e.POST("/load", handlers.LoadBlock(db))
	e.POST("/uploadfile", handlers.UpLoadFile())

	e.Logger.Fatal(e.Start(":8080"))
}
