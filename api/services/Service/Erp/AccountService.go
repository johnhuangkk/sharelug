package Erp

import (
	"api/services/dao/Balance"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"time"
)

func HandleAccountAmount(UserId string, dataId string, amount float64, transType string, comment string) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	data, err := Balance.GetBalanceAccountLastByUserId(engine, UserId)
	if err != nil {
		log.Error("Get Balance Account Database Error")
		return err
	}
	var ent entity.BalanceAccountData
	ent.UserId = UserId
	ent.DataId = dataId
	ent.TransType = transType
	ent.In = 0
	ent.Out = amount
	ent.Balance = data.Balance - amount
	ent.Comment = comment
	ent.CreateTime = time.Now()
	err = Balance.InsertBalanceAccountData(engine, ent)
	if err != nil {
		log.Error("Insert Balance Account Database Error")
		return err
	}
	return nil
}
