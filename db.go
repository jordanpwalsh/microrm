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

	_ "github.com/mattn/go-sqlite3"
)

type Microrm struct {
	sqlDb *sql.DB
	path  string
}

func Open(path string) (*Microrm, error) {
	db := new(Microrm) //investigate why fields are not set - nil
	var err error

	db.sqlDb, err = sql.Open("sqlite3", path)
	db.path = path
	if err != nil {
		return nil, err
	}
	fmt.Println(db)
	return db, nil
}

func (db *Microrm) Close() {
	db.sqlDb.Close()
}

//this will be replaced once migrations are a thing
func (microrm *Microrm) CreateTable(tableName string, tableStruct interface{}) (bool, error) {
	//Support struct tags for modifiers like not null - we'll allow nulls for the moment

	createQuery := fmt.Sprintf("CREATE TABLE %s (", tableName)

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
	fmt.Println(createQuery)

	if _, err := (microrm.sqlDb.Exec(createQuery)); err != nil {
		log.Fatal(err)
		return false, err
	}
	return true, nil
}

//this will be replaced once migrations are a thing
func (microrm *Microrm) DropTable(tableName string) (bool, error) {
	dropQuery := fmt.Sprintf("DROP TABLE IF EXISTS %s", tableName)

	_, err := microrm.sqlDb.Exec(dropQuery)
	if err != nil {
		log.Fatal(err)
		return false, err
	}
	return true, nil
}

//to make this chainable, just return a reciever
func FindOne(db *sql.DB, tableObj interface{}) (bool, error) {
	tableName := strings.ToLower(reflect.TypeOf(tableObj).Elem().Name())

	selectQuery := fmt.Sprintf("SELECT * FROM %s", tableName)

	rows, err := db.Query(selectQuery)
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
		reflect.ValueOf(tableObj).Elem().FieldByName(fieldNameCased).
			Set(reflect.ValueOf(value))
	}

	//make sure query rows.next again and error out cause we're only returning one here
	return true, nil
}
