package Sms

import (
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"time"
)

func InsertSmsData(phone, content, merchantId string) (entity.SmsLogData, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var data entity.SmsLogData
	data.Phone = phone
	data.MerchantId = merchantId
	data.Content = content
	data.CreateTime = time.Now()
	if _, err := engine.Session.Table(entity.SmsLogData{}).Insert(&data); err != nil {
		log.Error("Insert Sms Logs Database", err)
		return data, err
	}
	return data, nil
}

func UpdateSmsData(data entity.SmsLogData, code, msg string) error {
	data.ResultCode = code
	data.ResultMsg = msg
	engine := database.GetMysqlEngine()
	defer engine.Close()
	if _, err := engine.Session.Table(entity.SmsLogData{}).ID(data.Id).AllCols().Update(&data); err != nil {
		log.Error("Update Sms Logs Database", err)
		return err
	}
	return nil
}