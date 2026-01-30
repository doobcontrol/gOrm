package xyOrm

import (
	"testing"
	"github.com/doobcontrol/gDb/xyDb"
	"github.com/doobcontrol/gDbSqlite/xyDbSqlite"
)

func TestMakeInsertStr(t *testing.T) {
	recordMap := map[string]string{
		"a":"1",
		"b":"2",
	}
	
	fields, values := MakeInsertStr(recordMap)
	if fields != "a,b" {
		t.Errorf("TestMakeInsertStr expect fields: %s, but got: %s", "a,b", fields)
	}
	if values != "'1','2'" {
		t.Errorf("TestMakeInsertStr expect values: %s, but got: %s", "'1','2'", values)
	}
}

func TestMakeUpdateStr(t *testing.T) {
	recordMap := map[string]string{
		"a":"1",
		"b":"2",
	}
	
	updateString := MakeUpdateStr(recordMap)
	if updateString != "a='1',b='2'" {
		t.Errorf("TestMakeUpdateStr expect updateString: %s, but got: %s", "a='1',b='2'", updateString)
	}
}
func TestBaseModelAssignFieldNames(t *testing.T) {
	type testModel struct{
		baseModel
		FID string
	}

	tModel := &testModel{}

	if err := tModel.AssignFieldNames(tModel); err != nil {
		t.Errorf("TestBaseModelAssignFieldNames error: %s", err)
	} else {
		if tModel.FID != "FID" {
			t.Errorf("TestBaseModelAssignFieldNames expect tModel.FID value as: %s, but got: %s", 
			"F1", tModel.FID)
		}
		if tModel.Name == "Name" {
			t.Errorf("TestBaseModelAssignFieldNames expect tModel.Name excluded, but still got: %s", 
			tModel.Name)
		}
	}
}

func TestBaseModelGetDbTable(t *testing.T) {
	bModel := &baseModel{}
	bModel.Name = "testTable"
	bModel.Fields = []TbField{
		{
			Name: "F1",
			Type: "string",
			Length: 50,
			IsKey: true,
		},
		{
			Name: "F2",
			Type: "string",
			Length: 50,
			IsKey: false,
		},
	}
	dTable := bModel.GetDbTable()
	if dTable.TableName != "testTable" {
		t.Errorf("TestBaseModelGetDbTable expect dTable.TableName: %s, but got: %s", "testTable", dTable.TableName)
	}
	if len(dTable.Fields) != 2 {
		t.Errorf("TestBaseModelGetDbTable expect dTable.Fields: %d, but got: %d", 2, len(dTable.Fields))
	}
	if dTable.Fields[0].FieldName != "F1" || dTable.Fields[0].DataType != "string" || 
		dTable.Fields[0].Length != 50 || !dTable.Fields[0].IsKey {
		t.Errorf("TestBaseModelGetDbTable expect dTable.Fields[0] values: %s %s %d %t, but got: %s %s %d %t,", 
		"F1","string",50,true,
		dTable.Fields[0].FieldName,dTable.Fields[0].DataType,
		dTable.Fields[0].Length,dTable.Fields[0].IsKey)
	}
}

func TestBaseModelSql(t *testing.T) {
	setSqlTest()

	bModel := GetUser().baseModel

	if err := bModel.Insert(map[string]string{
		"FID":"001",
		"FUserName":"aaa",
		"FPassword":"bbb",
	}); err != nil {
		t.Errorf("TestBaseModelSql Insert 001 error: %s", err)
	}
	
	if err := bModel.Insert(map[string]string{
		"FID":"002",
		"FUserName":"aaa2",
		"FPassword":"bbb2",
	}); err != nil {
		t.Errorf("TestBaseModelSql Insert 002 error: %s", err)
	}

	if record, err := bModel.SelectAll(); err != nil {
		t.Errorf("TestBaseModelSql SelectAll error: %s", err)
	} else {
		if len(*record) != 2 {
			t.Errorf("TestBaseModelSql SelectAll expect count: %d, but got: %d", 2, len(*record))
		} else {
			if (*record)[1]["FID"] != "002" || (*record)[1]["FUserName"] != "aaa2" || 
			(*record)[1]["FPassword"] != "bbb2" {
				t.Errorf("TestBaseModelSql SelectAll record2 expect values: %s %s %s, but got: %s %s %s,", 
					"002","aaa2","bbb2",
					(*record)[1]["FID"], (*record)[1]["FUserName"],
					(*record)[1]["FPassword"])
			}
		}
	}
	
	if err := bModel.UpdateByField(map[string]string{
		"FUserName":"aaa2u",
		"FPassword":"bbb2u",
	}, "FID", "002"); err != nil {
		t.Errorf("TestBaseModelSql update 002 error: %s", err)
	}

	if record, err := bModel.SelectAll(); err != nil {
		t.Errorf("TestBaseModelSql update_SelectAll error: %s", err)
	} else {
		if len(*record) != 2 {
			t.Errorf("TestBaseModelSql update_SelectAll expect count: %d, but got: %d", 2, len(*record))
		} else {
			if (*record)[1]["FID"] != "002" || (*record)[1]["FUserName"] != "aaa2u" || 
			(*record)[1]["FPassword"] != "bbb2u" {
				t.Errorf("TestBaseModelSql update_SelectAll record2 expect values: %s %s %s, but got: %s %s %s,", 
					"002","aaa2u","bbb2u",
					(*record)[1]["FID"], (*record)[1]["FUserName"],
					(*record)[1]["FPassword"])
			}
		}
	}
	
	if err := bModel.DeleteByField("FID", "002"); err != nil {
		t.Errorf("TestBaseModelSql DeleteByField 002 error: %s", err)
	}

	if record, err := bModel.SelectAll(); err != nil {
		t.Errorf("TestBaseModelSql delete_SelectAll error: %s", err)
	} else {
		if len(*record) != 1 {
			t.Errorf("TestBaseModelSql delete_SelectAll expect count: %d, but got: %d", 1, len(*record))
		} else {
			if (*record)[0]["FID"] != "001" || (*record)[0]["FUserName"] != "aaa" || 
			(*record)[0]["FPassword"] != "bbb" {
				t.Errorf("TestBaseModelSql delete_SelectAll record1 expect values: %s %s %s, but got: %s %s %s,", 
					"002","aaa2u","bbb2u",
					(*record)[0]["FID"], (*record)[0]["FUserName"],
					(*record)[0]["FPassword"])
			}
		}
	}

	cleanDb()
}

func setSqlTest(){
	cleanDb()
	setTest()

	connectString, _ := InitModel(
		testDFile,
		&map[string]string{xyDbSqlite.S_dbFile: testDFile}, 
		[]XyModel{GetUser()})
	xyDb.DService.DbAccess.Close()

	ConfigModel(connectString)
}
