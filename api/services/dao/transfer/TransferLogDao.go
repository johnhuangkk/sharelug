package transfer

import (
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"time"
)

func InsertTransferLog(engine *database.MysqlSession, json string) error {
	var data entity.TransferLogData
	data.Response = json
	data.CreateTime = time.Now()
	if _, err := engine.Session.Table(entity.TransferLogData{}).Insert(&data); err != nil {
		log.Error("Database Error", err)
		return err
	}
	return nil
}
