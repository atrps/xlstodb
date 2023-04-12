// Block Type
package loader

import (
	"encoding/json"

	"database/sql"
)

func GetBlockTypeTableName(db *sql.DB, blockType string) string {
	var tbName string
	statement := "select table_name from block_type_ref where code = $1"
	if err := db.QueryRow(statement, blockType).Scan(&tbName); err != nil {
		if err == sql.ErrNoRows {
			return ""
		}
	}
	return tbName
}

// Fields get from BlockTypeDtl
func GetBlockTypeDtlFields(db *sql.DB, blockType string) []FieldItem {
	var statement string
	var dtlFields []FieldItem
	var jsonText []byte

	statement = "select block_type_dtl_fields($1) as dtlFields"
	if err := db.QueryRow(statement, blockType).Scan(&jsonText); err != nil {
		if err == sql.ErrNoRows {
			return []FieldItem{}
		}
	}
	if err := json.Unmarshal(jsonText, &dtlFields); err != nil {
		return []FieldItem{}
	}

	return dtlFields
}
