package models

import (
	"database/sql"

	_ "github.com/lib/pq"
)

// Block is a struct containing Block_load data
type Block struct {
	Id           int64  `json:"id"`
	Info         string `json:"info"`
	Type_code    string `json:"type_code"`
	Created_at   string `json:"created_at"`
	Status       string `json:"status"`
	Status_name  string `json:"status_name"`
	Status_date  string `json:"status_date"`
	Count_rec    int    `json:"count_rec"`
	Completed_at string `json:"completed_at"`
	Filename     string `json:"filename"`
	Sheet        string `json:"sheet"`
	Rowfirst     int    `json:"rowfirst"`
	Rowlast      int    `json:"rowlast"`
}

// BlockList is list of Block
type BlockList struct {
	Blocks []Block `json:"items"`
}

// query for get block
func getSql(oneRec bool) string {
	sql := `
	SELECT
		b.Id,
		b.Info,
		t.Code as Type_code,
		to_char(b.Created_at, 'DD.MM.YYYY HH24:MI:SS') as Created_at,
		b.Status,
		block_status_name(b.Status) as Status_name,
		to_char(coalesce(b.Completed_at, b.Created_at), 'DD.MM.YYYY HH24:MI:SS') as Status_date,
		coalesce(Count_rec, 0) as Count_rec,
		coalesce(to_char(b.Completed_at, 'DD.MM.YYYY HH24:MI:SS'),'') as Completed_at
	FROM
		block_load b
		INNER JOIN block_type_ref t ON b.type_id = t.id
	`
	if oneRec {
		sql = sql + `
		WHERE b.id = $1
		`
	} else {
		sql = sql + `
		ORDER BY b.id DESC LIMIT 10
		`
	}
	return sql
}

// GetBlock on Id
func GetBlockOnId(db *sql.DB, id int64) Block {
	bk := Block{}
	sql := getSql(true)

	rows, err := db.Query(sql, id)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		err2 := rows.Scan(
			&bk.Id,
			&bk.Info,
			&bk.Type_code,
			&bk.Created_at,
			&bk.Status,
			&bk.Status_name,
			&bk.Status_date,
			&bk.Count_rec,
			&bk.Completed_at)
		if err2 != nil {
			panic(err2)
		}
	}

	return bk
}

// GetBlocks from the DB
func GetBlocks(db *sql.DB) BlockList {
	result := BlockList{}
	sql := getSql(false)

	rows, err := db.Query(sql)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		bk := Block{}
		err2 := rows.Scan(
			&bk.Id,
			&bk.Info,
			&bk.Type_code,
			&bk.Created_at,
			&bk.Status,
			&bk.Status_name,
			&bk.Status_date,
			&bk.Count_rec,
			&bk.Completed_at)
		// Exit if we get an error
		if err2 != nil {
			panic(err2)
		}
		result.Blocks = append(result.Blocks, bk)
	}

	return result
}
