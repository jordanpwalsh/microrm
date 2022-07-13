package microrm

import (
	"fmt"
	"testing"
)

type TestStructure struct {
	id        int `microrm:"pk"`
	name      string
	byte_val  byte
	float_val float64

	//add all the go primitives
}

var microrm *Microrm
var err error

func TestCreateTable(t *testing.T) {
	fmt.Println("from testcreatetable:", microrm.path)
	//microrm, err = Open("./unit_test.db")
	if err != nil {
		t.Errorf("Failed to create database")
	}

	createResult, error := microrm.CreateTable("test_table", TestStructure{})
	if createResult != true || error != nil {
		t.Errorf("Failed to create table")
	}
}

func TestDropTable(t *testing.T) {

	// //figure out how to reuse some of these so prevent a bunch of duplication.
	// microrm, err = Open("./unit_test.db")
	// if err != nil {
	// 	t.Errorf("Failed to create database")
	// }

	dropResult, error := microrm.DropTable("test_table")
	if dropResult != true || error != nil {
		t.Errorf("Failed to drop table")
	}
}

func init() {
	fmt.Println("init running")
	var err error
	microrm = &Microrm{path: "test"}

	//play with microrm instantiation
	microrm, err = Open("./unit_test.db")
	if err != nil {
		panic("Failed to create database")
	}
	//defer microrm.Close()
}
