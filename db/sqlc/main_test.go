package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"os"
	"testing"
)

var queryTest *Queries
var GlobalDB *sql.DB

const (
	driver = "postgres"
	path   = "postgresql://root:123456@127.0.0.1:5432/bank?sslmode=disable"
)

func TestMain(m *testing.M) {
	conn, err := sql.Open(driver, path)
	if err != nil {
		log.Fatal("db can not connection")
	}
	GlobalDB = conn
	queryTest = New(conn)
	os.Exit(m.Run())
}
