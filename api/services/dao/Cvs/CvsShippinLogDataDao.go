package Cvs

import (
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"api/services/util/tools"
	"fmt"
)

// 寫入超商配送log
func InsertCvsShippingLogData(engine *database.MysqlSession, logData entity.CvsShippingLogData) error {
	logData.CreateTime = tools.Now(`YmdHis`)
	_, err := engine.Session.Table(entity.CvsShippingLogData{}).Insert(logData)
	if err != nil {
		log.Error("InsertShippingLogData Data Error: [%s]", logData)
		log.Error("InsertShippingLogData Error: [%s]", err.Error())
		return fmt.Errorf("資料庫異常")
	}
	return  nil
}

// 取得超商配送log
func GetCvsShippingLogData(engine *database.MysqlSession, orderData entity.OrderData) (data []entity.CvsShippingLogData) {
	var query = map[string]interface{}{}
	query[`ship_no`] = orderData.ShipNumber
	query[`cvs_type`] = orderData.ShipType
	query[`is_show`] = `1`

	err := engine.Engine.Table(entity.CvsShippingLogData{}).Where(query).Desc(`date_time`).Find(&data)
	if err != nil {
		log.Error("GetCvsShippingLogData Error: [%s]", err.Error())
		return data
	}
	return data
}

// 取得一筆紀錄
func GetOneCvsShippingLogData(engine *database.MysqlSession, shipNo, typeX string) (data entity.CvsShippingLogData) {
	var query = map[string]interface{}{}
	query[`ship_no`] = shipNo
	query[`type`] = typeX

	_, err := engine.Engine.Table(entity.CvsShippingLogData{}).Where(query).Desc(`date_time`).Limit(1, 0).Get(&data)

	if err != nil {
		log.Error("GetOneCvsShippingLogData Error: [%s]", err.Error())
		return data
	}
	return data
}