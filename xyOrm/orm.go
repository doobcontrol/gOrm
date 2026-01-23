package xyOrm

import (
	"github.com/doobcontrol/gDb/xyDb"
)

type XyModel interface {
	GetDbTable() xyDb.DbTable
}

func SetDbAccess(dbAccess xyDb.IDbAccess){
	xyDb.DService.DbAccess = dbAccess
}

func InitModel(
	dbName string, 
	initPars *map[string]string,
	models []XyModel,
	) (string, error) {
	dbStructure := xyDb.DbStructure{
		DbName: dbName,
		Tables: []xyDb.DbTable{},
	}

	for _, model := range models {
		tableStruc := model.GetDbTable()
		dbStructure.Tables = append(dbStructure.Tables, tableStruc)
	}

	connectString, err := xyDb.DService.Init(initPars, dbStructure)
	if err != nil {
		return "", err
	}

	return connectString, nil
}

func ConfigModel(
	dbConnectString string,
	) error {
		xyDb.DService.Set(dbConnectString)
		return nil
}