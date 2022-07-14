package microrm

import (
	"fmt"
	"os"
	"testing"
)

type TestStructure struct {
	id        int `microrm:"pk"`
	name      string
	byte_val  byte
	float_val float64
}

var microrm *Microrm
var err error

func TestMapField(t *testing.T) {
	testStruct := TestStructure{
		id:        1,
		name:      "testVarName",
		byte_val:  22,
		float_val: 3.14159,
	}

	fieldMappings := mapField(testStruct)
	fmt.Println("field mappings length:", len(fieldMappings))
	fmt.Println("name:", fieldMappings[0].name)
	fmt.Println("type:", fieldMappings[0].dataType)
	fmt.Println("sql type:", fieldMappings[0].sqlType)
	fmt.Println("tag:", fieldMappings[0].tag)
	//TODO: write some test conditions

}
func TestOpen(t *testing.T) {
	microrm, err = Open("./unit_test.db")
	if err != nil {
		t.Errorf("Failed to create database")
	}
}

func TestCreateTable(t *testing.T) {
	fmt.Println("from testcreatetable:", microrm.path)
	createResult, error := microrm.CreateTable("test_table", TestStructure{})
	if createResult != true || error != nil {
		t.Errorf("Failed to create table")
	}
}

func TestDropTable(t *testing.T) {
	dropResult, error := microrm.DropTable("test_table")
	if dropResult != true || error != nil {
		t.Errorf("Failed to drop table")
	}
}

func TestCloseTable(t *testing.T) {
	microrm.Close()
	if err := os.Remove("./unit_test.db"); err != nil {
		t.Errorf("Failed to remove database file")
	}
}
