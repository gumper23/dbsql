package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "rsmith:***REMOVED***@tcp(127.0.0.1:13306)/information_schema")
	if err != nil {
		log.Fatalf("%s\n", err.Error())
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("%s\n", err.Error())
	}
	rs, err := GetResultset(db, "show slave status")
	if err != nil {
		log.Fatalf("%s\n", err.Error())
	}
	fmt.Printf("%+v\n", rs)
	fmt.Printf("%s:%s\n", rs[0]["Master_Host"], rs[0]["Master_Port"])

	rs, err = GetResultset(db, "select * from steam.game")
	if err != nil {
		log.Fatalf("%s\n", err.Error())
	}
	for _, row := range rs {
		fmt.Printf("%v\n", row)
	}

}

// GetResultset returns a slice of string->interfaces.
func GetResultset(db *sql.DB, query string) (resultset []map[string]interface{}, err error) {
	// Execute the query.
	rows, err := db.Query(query)
	if err != nil {
		return
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return
	}

	resultset = make([]map[string]interface{}, 0)
	for rows.Next() {
		// Create a slice of interface{}'s to represent each column,
		// and a second slice to contain pointers to each item in the columns slice.
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))

		for i := 0; i < len(columns); i++ {
			columnPointers[i] = &columns[i]
		}

		// Scan the result into the column pointers...
		if err := rows.Scan(columnPointers...); err != nil {
			return nil, err
		}

		// Create our map, and retrieve the value for each column from the pointers slice,
		// storing it in the map with the name of the column as the key.
		m := make(map[string]interface{})
		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			m[colName] = *val
		}

		// Outputs: map[columnName:value columnName2:value2 columnName3:value3 ...]
		fmt.Print(m)

		resultset = append(resultset, m)
	}
	err = rows.Err()
	return
}
