package xyOrm

import (
	"errors"
	"os"
	"reflect"
	"testing"

	"github.com/doobcontrol/gDb/xyDb"
	"github.com/doobcontrol/gDbSqlite/xyDbSqlite"
)

func TestInitModel(t *testing.T) {
	cleanDb()

	setTest()

	if _, err := InitModel(
		testDFile,
		&map[string]string{xyDbSqlite.S_dbFile: testDFile}, 
		[]XyModel{GetUser()}); err != nil {
		t.Errorf("TestInitModel expect err is nil, but got: %s", err.Error())
	} else {
		if !fileExists(testDFile) {
			t.Errorf("TestInitModel expect database file: %s, but not exists", testDFile)
		}
	}

	cleanDb()
}
func TestConfigModel(t *testing.T) {
	cleanDb()

	setTest()

	connectString, _ := InitModel(
		testDFile,
		&map[string]string{xyDbSqlite.S_dbFile: testDFile}, 
		[]XyModel{GetUser()})
	xyDb.DService.DbAccess.Close()

	if err := ConfigModel(connectString); err != nil {
		t.Errorf("TestConfigModel expect err is nil, but got: %s", err.Error())
	} else {
		if _, err := GetUser().SelectAll(); err != nil {
			t.Errorf("TestConfigModel select test expect err is nil, but got: %s", err.Error())
		}
	}

	cleanDb()
}
func TestClean(t *testing.T) {
	cleanDb()

	setTest()

	connectString, _ := InitModel(
		testDFile,
		&map[string]string{xyDbSqlite.S_dbFile: testDFile}, 
		[]XyModel{GetUser()})
	xyDb.DService.DbAccess.Close()

	if err := ConfigModel(connectString); err != nil {
		t.Errorf("TestClean ConfigModel error: %s", err.Error())
	} else {
		Clean()
		_, err = GetUser().SelectAll()
		if _, err := GetUser().SelectAll(); err == nil {
			t.Errorf("TestClean after clean, select expect an error, but got nil")
		} else {
			if err.Error() != "sql: database is closed" {
				t.Errorf("TestClean after clean, select expect an error: %s, but got: %s",
				"sql: database is closed", err.Error())
			}
		}
	}

	cleanDb()
}

var testDFile = "./testDb"
func cleanDb(){
	os.Remove(testDFile)
}

func setTest()  {
	SetDbAccess(&xyDbSqlite.DbSqliteAccess{})
}

//Define a test kModel
type User struct{
	KModel
	FUserName string
	FPassword string
}

func (km *User) CreateFields() {
	//parent model fields
	km.KModel.CreateFields()

	km.Fields = append(
		km.Fields, 
		TbField{
			Name: km.FUserName,
			Type: "string",
			Length: 50,
			IsKey: false,
		})
	km.Fields = append(
		km.Fields, 
		TbField{
			Name: km.FPassword,
			Type: "string",
			Length: 50,
			IsKey: false,
		})
}

//Factory Function to implement a constructor
func NewUser() *User {
	newUser := User{}

	newUser.RealModel = true
	newUser.Name = reflect.ValueOf(newUser).Type().Name()
	newUser.AssignFieldNames(&newUser)
	newUser.CreateFields()

	return &newUser
}
//Singleton
var signleUser *User
func GetUser() *User {
	if signleUser == nil {
		signleUser = NewUser()
	}
	return signleUser
}

// fileExists checks if a file exists and is not a directory.
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if errors.Is(err, os.ErrNotExist) {
		return false // File does not exist
	}
    // If err is nil, check if it's a directory
	return err == nil && !info.IsDir()
}