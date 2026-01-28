package xyOrm

import (
	"testing"
	"reflect"
)

func TestKModelCreateFields(t *testing.T) {
	newUser := User{}

	newUser.RealModel = true
	newUser.Name = reflect.ValueOf(newUser).Type().Name()
	AssignFieldNames(&newUser)

	testModel := newUser.KModel
	testModel.CreateFields()
	if len(testModel.Fields) != 1 {
		t.Errorf("TestKModelCreateFields expect Fields count: %d, but got: %d", 1, len(testModel.Fields))
	} else {
		if testModel.Fields[0].Name != "FID" || testModel.Fields[0].Type != "string" || 
		testModel.Fields[0].Length != 50 || !testModel.Fields[0].IsKey {
			t.Errorf("TestKModelCreateFields expect Field0 value: %s %s %d %t, but got: %s %s %d %t", 
			"FID", "string", 50, true,
			testModel.Fields[0].Name, testModel.Fields[0].Type, testModel.Fields[0].Length, 
		    testModel.Fields[0].IsKey)
		}
		if testModel.FID != "FID" {
			t.Errorf("TestKModelCreateFields expect testModel.FID value: %s, but got: %s", 
			"FID", testModel.FID)
		}
	}	
}

func TestKModelSql(t *testing.T) {
	setSqlTest()

	testModel := GetUser().KModel

	testModel.Insert(map[string]string{
		"FID":"001",
		"FUserName":"aaa",
		"FPassword":"bbb",
	})
	testModel.Insert(map[string]string{
		"FID":"002",
		"FUserName":"aaa2",
		"FPassword":"bbb2",
	})

	if record, err := testModel.SelectByPk("002"); err != nil {
		t.Errorf("TestKModelSql SelectByPk 002 error: %s", err)
	} else {
		if (*record)["FID"] != "002" || (*record)["FUserName"] != "aaa2" || 
		(*record)["FPassword"] != "bbb2" {
			t.Errorf("TestKModelSql SelectByPk 002 record expect values: %s %s %s, but got: %s %s %s,", 
				"002","aaa2","bbb2",
				(*record)["FID"], (*record)["FUserName"],
				(*record)["FPassword"])
		}
	}

	if _, err := testModel.SelectByPk("003"); err == nil {
		t.Errorf("TestKModelSql SelectAll 003 expect an error, but got nil")
	} else {
		if err.Error() != "No record found" {
			t.Errorf("TestKModelSql SelectAll 003 expect an error: %s, but got: %s",
				"No record found", err.Error())
		}
	}
	
	if err := testModel.UpdateByByPk(map[string]string{
		"FUserName":"aaa2u",
		"FPassword":"bbb2u",
	}, "002"); err != nil {
		t.Errorf("TestKModelSql UpdateByByPk 002 error: %s", err)
	}
	if record, err := testModel.SelectByPk("002"); err != nil {
		t.Errorf("TestKModelSql SelectByPk 002 error: %s", err)
	} else {
		if (*record)["FID"] != "002" || (*record)["FUserName"] != "aaa2u" || 
		(*record)["FPassword"] != "bbb2u" {
			t.Errorf("TestKModelSql SelectByPk 002 record expect values: %s %s %s, but got: %s %s %s,", 
				"002","aaa2u","bbb2u",
				(*record)["FID"], (*record)["FUserName"],
				(*record)["FPassword"])
		}
	}
	
	if err := testModel.DeleteByByPk("002"); err != nil {
		t.Errorf("TestKModelSql DeleteByByPk 002 error: %s", err)
	}
	if _, err := testModel.SelectByPk("002"); err == nil {
		t.Errorf("TestKModelSql SelectAll 002 expect an error, but got nil")
	} else {
		if err.Error() != "No record found" {
			t.Errorf("TestKModelSql SelectAll 002 expect an error: %s, but got: %s",
				"No record found", err.Error())
		}
	}

	cleanDb()
}