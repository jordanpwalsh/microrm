package microrm

//todo notes
//need to make type and implement create database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/kataras/golog"

	_ "github.com/mattn/go-sqlite3"
)

type Microrm struct {
	sqlDb *sql.DB
	path  string
}

type RecordMapping struct {
	name     string
	dataType reflect.Type
	sqlType  string
	tag      string
}

func inferTableName(record interface{}) string {
	//might need to check for empty here and return Elem.Name if not an empty struct
	tableName := reflect.TypeOf(record).Name()

	return tableName
}

func mapRecordFields(record interface{}) []RecordMapping {
	structFields := reflect.VisibleFields(reflect.TypeOf(record))
	recordMappings := make([]RecordMapping, 0)

	for _, field := range structFields {
		var recordMapping RecordMapping
		recordMapping.name = field.Name
		recordMapping.dataType = field.Type
		recordMapping.sqlType = field.Type.Kind().String() //expand on this later
		recordMapping.tag = field.Tag.Get("microrm")
		recordMappings = append(recordMappings, recordMapping)
	}
	return recordMappings
}

func Open(path string) (*Microrm, error) {
	db := Microrm{path: path}
	//db := new(Microrm) //investigate why fields are not set - nil
	var err error

	db.sqlDb, err = sql.Open("sqlite3", path)
	db.path = path
	if err != nil {
		return nil, err
	}
	return &db, nil
}

func (db *Microrm) Close() {
	db.sqlDb.Close()
}

//this will be replaced once migrations are a thing
func (microrm *Microrm) CreateTable(tableStruct interface{}) (bool, error) {
	tableName := inferTableName(tableStruct)
	golog.Info("createTable:inferTableName:", tableName)
	//Support struct tags for modifiers like not null - we'll allow nulls for the moment

	createQuery := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (", tableName)

	structFields := reflect.VisibleFields(reflect.TypeOf(tableStruct))
	for _, field := range structFields {
		//ok we have types and field named we can build a string from
		//handle int, string, and bool for the moment
		//pk todo: use rowid
		if field.Tag.Get("microrm") == "pk" {
			if field.Type.Kind() != reflect.Int {
				return false, errors.New("primary key must be type int")
			}
			createQuery += fmt.Sprintf("%s INTEGER PRIMARY KEY AUTOINCREMENT,", field.Name)
		} else {
			createQuery += fmt.Sprintf("%s %s,", field.Name, field.Type)
		}
	}
	createQuery = strings.TrimSuffix(createQuery, ",")
	createQuery += ")"

	golog.Debug(createQuery)

	if _, err := (microrm.sqlDb.Exec(createQuery)); err != nil {
		log.Fatal(err)
		return false, err
	}
	return true, nil
}

//this will be replaced once migrations are a thing
func (microrm *Microrm) DropTable(o interface{}) (bool, error) {
	var tableName string = inferTableName(o)
	dropQuery := fmt.Sprintf("DROP TABLE IF EXISTS %s", tableName)

	_, err := microrm.sqlDb.Exec(dropQuery)
	if err != nil {
		log.Fatal(err)
		return false, err
	}
	return true, nil
}

//refactor this to use field mapper function
func (microrm *Microrm) Find(tableObj interface{}, id int) (bool, error) {
	tableName := strings.ToLower(reflect.TypeOf(tableObj).Elem().Name())

	selectQuery := fmt.Sprintf("SELECT * FROM %s WHERE id=%d", tableName, id)

	rows, err := microrm.sqlDb.Query(selectQuery)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	rc, err := rows.ColumnTypes()
	if err != nil {
		return false, err
	}

	//make a new slice to scan the values
	obj := make([]interface{}, len(rc))
	for i := range rc {
		var iface interface{}
		obj[i] = &iface
	}

	//now we have a []interface{} filled with ready types and we're ready to scan
	if !rows.Next() {
		return false, nil
	}

	rows.Scan(obj...)

	//map to fields without calling them: todo fields <-> obj elems
	for i, field := range rc {
		var value = *(obj[i].(*interface{}))
		//get the type of the struct field
		fieldNameCased := strings.ToUpper(field.Name()[:1]) + field.Name()[1:]
		fieldType := reflect.ValueOf(tableObj).Elem().FieldByName(fieldNameCased).Type()
		//in the case of bools we need to check because sqlite does not have bools so
		//so they are stores as integers. The set() below will fail with cannot cast int64 to bool
		//TODO any other types to handle edge case? check the ones handled by scan? write a unit test to find out
		if fieldType.Kind() == reflect.Bool {
			//the problem was type assertion - research this
			if value.(int64) == 0 {
				value = false
			} else {
				value = true
			}
		}
		//reflect.ValueOf(tableObj).Elem().FieldByName(fieldNameCased).Set(reflect.ValueOf(value))
	}

	//make sure query rows.next again and error out cause we're only returning one here
	return true, nil
}

func (microrm *Microrm) InsertOne(tableObj interface{}) error {
	tableName := reflect.TypeOf(tableObj).Name()
	fieldMappings := mapRecordFields(tableObj)

	//get the field values from each object.. should the mapper do that?

	//build two strings: typeString: name, type and valueString: values
	//first from field mapping, second from fields by index if possible.
	var typeString string
	for _, field := range fieldMappings {
		typeString += fmt.Sprintf("%s,", field.name)
	}
	typeString = strings.TrimSuffix(typeString, ",")

	refTableObj := reflect.ValueOf(tableObj)
	var valueString string
	for i := 0; i < refTableObj.NumField(); i++ {
		valueString += fmt.Sprintf("\"%v\",", refTableObj.Field(i).Interface())
	}
	valueString = strings.TrimSuffix(valueString, ",")
	golog.Info("typeString:", typeString)
	golog.Info("valueString:", valueString)

	queryString := fmt.Sprintf("INSERT INTO %s(%s) VALUES(%s)", tableName, typeString, valueString)
	golog.Info("queryString:", queryString)
	result, err := microrm.sqlDb.Exec(queryString)
	if err != nil {
		golog.Error(result)
		return err
	}
	return nil

}
