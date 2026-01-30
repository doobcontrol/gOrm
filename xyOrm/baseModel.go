package xyOrm

import (
	"fmt"
	"reflect"
	"github.com/doobcontrol/gDb/xyDb"
	"slices"
	"strings"
)

var xydb = xyDb.DService

type TbField struct {
	Name string
	Type string
	Length int
	IsKey bool
}

//Define baseModel
type baseModel struct{
	Name string
	Fields []TbField
	RealModel bool //must be true to let the model map to database table
	FieldNamesSetted bool
}

func (bm *baseModel) GetDbTable() xyDb.DbTable {
	dbTable := xyDb.DbTable{
		TableName: bm.Name,
		Fields: []xyDb.DbField{},
	}

	for _, field := range bm.Fields {
		dbTable.Fields = append(
			dbTable.Fields,
			xyDb.DbField{
				FieldName: field.Name,
				DataType: field.Type,
				Length: field.Length,
				IsKey: field.IsKey,
			})
	}

	return dbTable
}

func (bm *baseModel) CreateFields() {
	bm.FieldNamesSetted = true
}

//query data
func (bm *baseModel) SelectAll() (*[]map[string]interface{}, error) {
	sql := fmt.Sprintf(
		"select * from %s", 
		bm.Name)
	return xydb.Query(sql)
}
func (bm *baseModel) SelectByField(field string, value string) (*[]map[string]interface{}, error) {
	sql := fmt.Sprintf(
		"select * from %s where %s='%s'", 
		bm.Name, field, value)
	return xydb.Query(sql)
}


//insert data
func (bm *baseModel) Insert(recordMap map[string]string) error {
	fieldsString,valuesString := MakeInsertStr(recordMap)
	sql := fmt.Sprintf(
		"insert into %s(%s) values(%s)", 
		bm.Name, 
		fieldsString, 
		valuesString)

	return xydb.ExSql(sql)
}

//update data
func (bm *baseModel) UpdateByField(recordMap map[string]string, field string, value string) error {
	setString  := MakeUpdateStr(recordMap)
	sql := fmt.Sprintf(
		"update %s set %s where %s='%s'", 
		bm.Name, 
		setString,
		field, 
		value)

	return xydb.ExSql(sql)
}

//delete data
func (bm *baseModel) DeleteByField(field string, value string) error {
	sql := fmt.Sprintf(
		"delete from %s where %s='%s'", 
		bm.Name, 
		field, 
		value)

	return xydb.ExSql(sql)
}

// AssignFieldNames dynamically sets all exported string fields of a struct to their own names.
var ExcludedFields = []string{"Name", "Fields"}
func (bm *baseModel) AssignFieldNames(s interface{}) error {
	if bm.FieldNamesSetted {
		return nil
	}
	
	// Get the reflect.Value of the interface.
	// We need a pointer to the struct to modify it, so use reflect.ValueOf(&s).Elem()
	// or ensure the input 's' is already a pointer and use reflect.ValueOf(s).Elem().
	val := reflect.ValueOf(s)

	// Check if it's a pointer and dereference if necessary
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// Ensure the value is a struct and is settable
	if val.Kind() != reflect.Struct {
		//return fmt.Errorf("input is not a struct or a pointer to a struct")
	}

	if !val.CanSet() {
		//return fmt.Errorf("struct is not settable (e.g., passed by value instead of pointer)")
	}

	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldName := typ.Field(i).Name

		// Check if the field is a string and is exported/settable
		found := slices.Contains(ExcludedFields, fieldName)
		fKind := field.Kind()
		canSet := field.CanSet()
		if (!found) && fKind == reflect.String && canSet {
			field.SetString(fieldName)
		}
	}

	return nil
}

//sql tool
func MakeInsertStr(recordMap map[string]string) (string, string) {
	var fieldsBuilder strings.Builder
	var valuesBuilder strings.Builder

	for key, value := range recordMap {
		if fieldsBuilder.Len() != 0{
			fieldsBuilder.WriteString(",")
			valuesBuilder.WriteString(",")
		}
		fieldsBuilder.WriteString(key)
		valuesBuilder.WriteString("'")
		valuesBuilder.WriteString(value)
		valuesBuilder.WriteString("'")
	}
	return fieldsBuilder.String(), valuesBuilder.String()
}
func MakeUpdateStr(recordMap map[string]string) string {
	var setsBuilder strings.Builder

	for key, value := range recordMap {
		if setsBuilder.Len() != 0{
			setsBuilder.WriteString(",")
		}
		setsBuilder.WriteString(fmt.Sprintf("%s='%s'", key, value))
	}
	return setsBuilder.String()
}