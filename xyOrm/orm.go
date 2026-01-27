package xyOrm

import (
	"errors"

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
	if 	dbName == "" {
		return "", errors.New("dbName is null")
	}
	if 	initPars == nil {
		return "", errors.New("initPars must not be nil")
	}
	if 	models == nil {
		return "", errors.New("models must not be nil")
	} else {
		if 	len(models) == 0 {
			return "", errors.New("no models")
		}
	}

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