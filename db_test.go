package microrm

import (
	"database/sql"
	"os"
	"testing"
)

type TestStructure struct {
	id        int `microrm:"pk"`
	name      string
	byte_val  byte
	float_val float64

	//add all the go primitives
}

var db *sql.DB

func TestCreateTable(t *testing.T) {
	createResult, error := CreateTable(db, "test_table", TestStructure{})
	if createResult != true || error != nil {
		t.Errorf("Failed to create table")
	}
}

func TestDropTable(t *testing.T) {
	dropResult, error := DropTable(db, "test_table")
	if dropResult != true || error != nil {
		t.Errorf("Failed to drop table")
	}
}

func TestMain(m *testing.M) {
	var err error
	os.Remove("./unit_test.db")
	db, err = sql.Open("sqlite3", "./unit_test.db")

	if err != nil {
		panic("Failed to create database")
	}
	defer db.Close()

	retCode := m.Run()

	os.Exit(retCode)

}
