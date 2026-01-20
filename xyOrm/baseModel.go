package xyOrm

import (
	"reflect"
	"github.com/doobcontrol/gDb/xyDb"
	"slices"
)

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

}

func (bm *baseModel) SelectAll() []map[string]string {
	xyDb.DService.SelectAll("")
	return nil
}


// AssignFieldNames dynamically sets all exported string fields of a struct to their own names.
var ExcludedFields = []string{"Name", "Fields"}
func AssignFieldNames(s interface{}) error {
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