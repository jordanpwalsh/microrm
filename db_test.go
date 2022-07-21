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
	_, err := Open("./unit_test.db")
	if err != nil {
		t.Errorf("Failed to create database")
	}
}

func TestCreateTable(t *testing.T) {
	var microrm *Microrm
	microrm = setUpTest(microrm)
	defer tearDownTest(microrm)

	createResult, error := microrm.CreateTable(TestStructure{})
	if createResult != true || error != nil {
		t.Errorf("Failed to create table")
	}
}

//TODO: refactor this to fix warning and crappy design
func setUpTest(microrm *Microrm) *Microrm {
	golog.Debug("Setting up test")
	err := os.Remove("./unit_test.db")
	if err != nil {
		golog.Errorf("Cannot remove database file")
	}

	microrm, _ = Open("./unit_test.db")
	microrm.CreateTable(TestStructure{})
	return microrm
}

func tearDownTest(microrm *Microrm) {
	microrm.Close()
}

//insert fails if tablObj is pointer.
func TestInsertOne(t *testing.T) {
	var microrm *Microrm
	microrm = setUpTest(microrm)
	defer tearDownTest(microrm)

	testStruct := TestStructure{
		Name:      "testVarName",
		Byte_val:  22,
		Float_val: 3.14159,
	}
	err := microrm.InsertOne(&testStruct)
	if err != nil {
		t.Error("Error inserting row", err)
	}
}

func TestFindExpectEmpty(t *testing.T) {
	type TestEmptyStruct struct {
		Id   int `microrm:"pk"`
		name string
	}

	var microrm *Microrm
	microrm = setUpTest(microrm)
	defer tearDownTest(microrm)

	microrm.CreateTable(TestEmptyStruct{})
	var testEmptyStruct TestEmptyStruct

	res, err := microrm.Find(&testEmptyStruct, 1)
	if err != nil {
		t.Errorf("Error finding row")
	}

	if res == false {
		return
	}
}

func TestFindExpectRow(t *testing.T) {
	type TestStruct struct {
		Id   int `microrm:"pk"`
		Name string
	}

	var microrm *Microrm
	microrm = setUpTest(microrm)
	defer tearDownTest(microrm)

	microrm.CreateTable(TestStruct{})
	testStruct := TestStruct{Name: "Jordan"}
	microrm.InsertOne(&testStruct)

	//overwrite struct
	testStruct = TestStruct{}

	_, err := microrm.Find(&testStruct, 1)
	if err != nil {
		t.Errorf("Error finding row")
	}

	if testStruct.Name != "Jordan" {
		t.Error("Find row value not as expected")
	}
}

func TestDropTable(t *testing.T) {
	var microrm *Microrm

	microrm = setUpTest(microrm)
	defer tearDownTest(microrm)

	//refactor this
	dropResult, error := microrm.DropTable(TestStructure{})
	if dropResult != true || error != nil {
		t.Errorf("Failed to drop table")
	}
}

func TestClose(t *testing.T) {
	var microrm *Microrm
	microrm = setUpTest(microrm)
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
