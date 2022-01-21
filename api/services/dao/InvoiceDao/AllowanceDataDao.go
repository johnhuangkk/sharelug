package InvoiceDao

import (
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"time"
)

func InsertAllowanceData(engine *database.MysqlSession, data entity.AllowanceData) error {
	data.CreateTime = time.Now()
	if _, err := engine.Session.Table(entity.AllowanceData{}).Insert(&data); err != nil {
		log.Error("Insert Invoice Database Error", err)
		return err
	}
	return nil
}

func GetAllowanceDataByAllowanceId(engine *database.MysqlSession, allowanceId string) (entity.AllowanceData, error) {
	var data entity.AllowanceData
	if _, err := engine.Engine.Table(entity.AllowanceData{}).
		Select("*").Where("allowance_id = ? ", allowanceId).Get(&data); err != nil {
		return data, err
	}
	return data, nil
}

func UpdateAllowanceData(engine *database.MysqlSession, data entity.AllowanceData) error {
	if _, err := engine.Session.Table(entity.AllowanceData{}).Where("allowance_id = ?", data.AllowanceId).Update(data); err != nil {
		log.Error("Update Database Error", err)
		return  err
	}
	return nil
}



