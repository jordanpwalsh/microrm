package microrm

import (
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

func TestCreateTable(t *testing.T) {
	microrm, err := Open("./unit_test.db")
	if err != nil {
		t.Errorf("Failed to create database")
	}

	createResult, error := microrm.CreateTable("test_table", TestStructure{})
	if createResult != true || error != nil {
		t.Errorf("Failed to create table")
	}
}

// func TestDropTable(t *testing.T) {
// 	dropResult, error := microrm.DropTable("test_table")
// 	if dropResult != true || error != nil {
// 		t.Errorf("Failed to drop table")
// 	}
// }

func TestMain(m *testing.M) {
	var err error

	microrm, err := Open("./unit_test.db")
	if err != nil {
		panic("Failed to create database")
	}
	defer microrm.Close()

	retCode := m.Run()

	os.Exit(retCode)

}
