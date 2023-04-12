// status loading
package loader

import (
	"fmt"

	"database/sql"

	"github.com/lib/pq"
)

// block status
type BlockStatus = struct {
	Status string `json:"status"`
	Count  string `json:"count"`
	Sdate  string `json:"sdate"`
	Info   string `json:"info"`
}

func GetBlockStatus(db *sql.DB, block_id int) BlockStatus {
	var ablock_status []string
	var statement string

	result := BlockStatus{}
	statement = "select block_current_status($1) as block_status"
	if err := db.QueryRow(statement, block_id).Scan(pq.Array(&ablock_status)); err != nil {
		if err == sql.ErrNoRows {
			return result
		}
	}
	fmt.Println("array block status ", ablock_status)
	result.Status = ablock_status[0]
	result.Count = ablock_status[1]
	result.Sdate = ablock_status[2]
	if result.Status == "L" {
		result.Info = fmt.Sprintf("%s load is completed, processed %s records ", ablock_status[2], ablock_status[1])
	} else {
		result.Info = fmt.Sprintf("%s loading %s records ", ablock_status[2], ablock_status[1])
	}
	return result
}
