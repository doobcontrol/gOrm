package xyOrm

import (
	"github.com/doobcontrol/gDb/xyDb"
)

type XyModel interface {
	GetDbTable() xyDb.DbTable
}

func InitModel(
	dbAccess xyDb.IDbAccess, 
	dbName string, 
	initPars map[string]string,
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

	xyDb.DService.DbAccess = dbAccess
	connectString, err := xyDb.DService.Init(initPars, dbStructure)
	if err != nil {
		return "", err
	}

	return connectString, nil
}

func ConfigModel(
	dbAccess xyDb.IDbAccess,
	dbConnectString string,
	) error {
		xyDb.DService.DbAccess = dbAccess
		xyDb.DService.Set(dbConnectString)
		return nil
}