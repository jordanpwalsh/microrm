package microrm

import (
	"os"
	"testing"

	"github.com/kataras/golog"
)

type TestStructure struct {
	Id        int `microrm:"pk"`
	Name      string
	Byte_val  byte
	Float_val float64
}

var microrm *Microrm
var err error

func TestInferTableName(t *testing.T) {
	tableName := inferTableName(TestStructure{})
	if tableName != "TestStructure" {
		golog.Info("TableName:", tableName)
		t.Error("Table name not as expected")
	}

	testStruct := TestStructure{
		Id:        1,
		Name:      "testVarName",
		Byte_val:  22,
		Float_val: 3.14159,
	}
	tableName = inferTableName(testStruct)
	if tableName != "TestStructure" {
		golog.Info("TableName:", tableName)
		t.Error("Table name not as expected")
	}

}

func TestMapField(t *testing.T) {
	testStruct := TestStructure{
		Id:        1,
		Name:      "testVarName",
		Byte_val:  22,
		Float_val: 3.14159,
	}

	fieldMappings := mapRecordFields(testStruct)
	golog.Debug("field mappings length:", len(fieldMappings))
	golog.Debug("name:", fieldMappings[0].name)
	golog.Debug("type:", fieldMappings[0].dataType)
	golog.Debug("sql type:", fieldMappings[0].sqlType)
	golog.Debug("tag:", fieldMappings[0].tag)
	//TODO: write some test conditions

}
func TestOpen(t *testing.T) {
	microrm, err = Open("./unit_test.db")
	if err != nil {
		t.Errorf("Failed to create database")
	}
}

func TestCreateTable(t *testing.T) {
	createResult, error := microrm.CreateTable(TestStructure{})
	if createResult != true || error != nil {
		t.Errorf("Failed to create table")
	}
}

func TestInsertOne(t *testing.T) {
	testStruct := TestStructure{
		Name:      "testVarName",
		Byte_val:  22,
		Float_val: 3.14159,
	}
	err := microrm.InsertOne(testStruct)
	if err != nil {
		t.Error("Error inserting row", err)
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

func init() {
	if os.Getenv("GOLOG_LEVEL") == "debug" {
		golog.SetLevel("debug")
	}
}
