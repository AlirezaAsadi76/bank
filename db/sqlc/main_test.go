package db

import (
	"database/sql"
	"firstprj/envi"
	_ "github.com/lib/pq"
	"log"
	"os"
	"testing"
)

var queryTest *Queries
var GlobalDB *sql.DB

func TestMain(m *testing.M) {
	config, err := envi.LoadConfig("../..")
	if err != nil {
		log.Fatal("envi can not load :", err)
	}
	conn, err := sql.Open(config.DbDriver, config.DbPath)
	if err != nil {
		log.Fatal("db can not connection")
	}
	GlobalDB = conn
	queryTest = New(conn)
	os.Exit(m.Run())
}
