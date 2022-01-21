package History

import (
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
)

func InsertHistoryLog(engine *database.MysqlSession, data entity.StatusHistoryLog) (entity.StatusHistoryLog, error) {
	_, err := engine.Session.Table(entity.StatusHistoryLog{}).Insert(&data)
	if err != nil {
		log.Error("Insert Status history Log Database Error", err)
		return data, err
	}
	return data, nil
}

//取出歷史訊息內容
func GetStatusHistoryByOrderId(engine *database.MysqlSession, Table string, Field string, DataId string) (entity.StatusHistoryLog, error) {
	var historyLog entity.StatusHistoryLog
	_, err := engine.Engine.Table(entity.StatusHistoryLog{}).
		Select("*").Where("data_id = ? ", DataId).And("table = ?", Table).And("Field = ?", Field).Get(&historyLog)
	if err != nil {
		return historyLog, err
	}
	return historyLog, nil
}

func InsertProductLog(engine *database.MysqlSession, data entity.ProductHistoryLog) error {
	_, err := engine.Session.Table(entity.ProductHistoryLog{}).Insert(&data)
	if err != nil {
		log.Error("Insert Status history Log Database Error", err)
		return err
	}
	return nil
}
