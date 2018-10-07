package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gumper23/sql/rs"
)

func main() {
	user := os.Getenv("MYSQL_USERNAME")
	pass := os.Getenv("MYSQL_PASSWORD")
	db, err := sql.Open("mysql", user+":"+pass+"@tcp(127.0.0.1:13306)/information_schema")
	if err != nil {
		log.Fatalf("%s\n", err.Error())
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("%s\n", err.Error())
	}

	var rs rs.Resultset
	err = rs.QueryRows(db, "show slave status")
	if err != nil {
		log.Fatalf("%s\n", err.Error())
	}
	rs.Vprint()

	err = rs.QueryRows(db, "select app_id, name, playtime_forever, created_at from steam.game order by playtime_forever desc limit 10")
	if err != nil {
		log.Fatalf("%s\n", err.Error())
	}
	rs.Hprint()
}
