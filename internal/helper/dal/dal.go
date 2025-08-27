package table

import (
	"cdp-admin-service/internal/config"
	"reflect"

	"gitlab.vrviu.com/cloud_esport_backend/dbwrap-public/dbwrap"
)

type TableInfo struct {
	DBWrap    *dbwrap.DBWrap
	TableName string
	Tpy       reflect.Type
}

var _TableMap map[string]*TableInfo = make(map[string]*TableInfo)

func InitDBWrap(c config.Config) {
	for _, tableInfo := range _TableMap {
		tableInfo.DBWrap = dbwrap.CreateDBWrap(c.EsportDal.EsportDalHost, 0, tableInfo.TableName, "", tableInfo.Tpy)

	}
}
