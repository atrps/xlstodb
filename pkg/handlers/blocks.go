package handlers

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"xlstodb/pkg/dbconnect"
	"xlstodb/pkg/loader"
	"xlstodb/pkg/models"

	"github.com/labstack/echo/v4"
)

type H map[string]interface{}

func GetBlocks(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, models.GetBlocks(db))
	}
}

func GetBlockOnId(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
		return c.JSON(http.StatusOK, models.GetBlockOnId(db, id))
	}
}

func LoadBlock(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var bk models.Block
		var cfgXls = loader.CfgXlsToDb{}
		var err error

		// Map imcoming JSON body
		c.Bind(&bk)
		fmt.Printf("%+v\n", bk)
		// fill JSON Cfg
		cfgXls.Filename = dbconnect.PathTmp + bk.Filename
		cfgXls.Sheet = bk.Sheet
		cfgXls.Rowfirst = bk.Rowfirst
		cfgXls.Rowlast = bk.Rowlast
		cfgXls.BlockInfo = bk.Info
		cfgXls.BlockType = bk.Type_code
		err = loader.CreateBlock(&cfgXls)
		if err == nil {
			go func() {
				err = loader.AppendRec(&cfgXls)
				if err == nil {
					fmt.Println("*** load block ", cfgXls.BlockId)
					fmt.Println("*** append records ", cfgXls.CntAppend)
				}
			}()
		}
		bk.Id = cfgXls.BlockId
		// Return a JSON response if successful
		if err == nil {
			return c.JSON(http.StatusCreated, bk)
			// Handle any errors
		} else {
			return err
		}
	}
}

func UpLoadFile() echo.HandlerFunc {
	return func(c echo.Context) error {
		r := c.Request()

		// Get handler for filename, size and headers
		file, handler, err := r.FormFile("xlsfile")
		if err != nil {
			fmt.Println("Error Retrieving the File")
			fmt.Println(err)
			return err
		}

		defer file.Close()
		fmt.Printf("Uploaded File: %+v\n", handler.Filename)
		fmt.Printf("File Size: %+v\n", handler.Size)
		fmt.Printf("MIME Header: %+v\n", handler.Header)
		filename := handler.Filename

		// Create file
		dst, err := os.Create(dbconnect.PathTmp + handler.Filename)
		if err != nil {
			fmt.Println("Error create the File on server")
			return err
		}
		defer dst.Close()

		// Copy the uploaded file to the created file on the filesystem
		if _, err := io.Copy(dst, file); err != nil {
			fmt.Println("Error copy the File to server")
			return err
		}

		return c.JSON(http.StatusCreated, H{
			"filename": filename,
		})

	}
}
