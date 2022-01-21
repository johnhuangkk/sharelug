package transfer

import (
	"api/services/Enum"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"time"
)

//新增轉帳資料
func InsertTransfer(engine *database.MysqlSession, data entity.TransferData) (entity.TransferData, error) {
	data.CreateTime = time.Now()
	data.UpdateTime = time.Now()
	_, err := engine.Session.Table(entity.TransferData{}).Insert(&data)
	if err != nil {
		log.Error("Database Error", err)
		return data, err
	}
	return data, nil
}
//取出訂單的轉帳資料
func GetTransferByOrderId(engine *database.MysqlSession, orderId string) (entity.TransferData,error) {
	var transferData entity.TransferData
	if _, err := engine.Engine.Table(entity.TransferData{}).Select("*").
		Where("order_id = ? AND transfer_status = ?", orderId, Enum.TransferInit).Get(&transferData); err != nil {
		log.Error("get transfer Database Error", err)
		return transferData, err
	}
	return transferData, nil
}

func GetTransferByOrderIds(engine *database.MysqlSession, orderId string) (entity.TransferData,error) {
	var transferData entity.TransferData
	if _, err := engine.Engine.Table(entity.TransferData{}).Select("*").
		Where("order_id = ?", orderId).Get(&transferData); err != nil {
		log.Error("get transfer Database Error", err)
		return transferData, err
	}
	return transferData, nil
}

//取出轉帳帳號資料
func GetTransferByAccount(engine *database.MysqlSession, Account string) (entity.TransferData,error) {
	var data entity.TransferData
	if _, err := engine.Engine.Table(entity.TransferData{}).Select("*").
		Where("bank_account = ?", Account).Get(&data); err != nil {
		log.Error("get transfer Database Error", err)
		return data, err
	}
	return data, nil
}

func GetTransferAndOrderByAccount(engine *database.MysqlSession, Account string) (entity.OrderAndTransfer,error) {
	var data entity.OrderAndTransfer
	if _, err := engine.Engine.Table(entity.TransferData{}).Select("*").
		Join("LEFT OUTER", entity.OrderData{}, "transfer_data.order_id = order_data.order_id").
		Where("bank_account = ?", Account).Get(&data); err != nil {
		log.Error("get transfer Database Error", err)
		return data, err
	}
	return data, nil
}

//更新Transfer Data
func UpdateTransferDate(engine *database.MysqlSession, Id int, Data entity.TransferData) error {
	Data.UpdateTime = time.Now()
	if _, err := engine.Session.Table(entity.TransferData{}).ID(Id).AllCols().Update(Data); err != nil {
		return err
	}
	return nil
}
//取出過期的轉帳資料
func GetTransferExpire(engine *database.MysqlSession, date string) ([]entity.TransferData,error) {
	var data []entity.TransferData
	if err := engine.Engine.Table(entity.TransferData{}).
		Where("transfer_status = ? AND expire_date < ?", Enum.TransferInit, date).
		Asc("create_time").Find(&data); err != nil {
		log.Error("get transfer Database Error", err)
		return data, err
	}
	return data, nil
}

func GetTransferQuery(engine *database.MysqlSession, date string) ([]entity.TransferData,error) {
	var data []entity.TransferData
	if err := engine.Engine.Table(entity.TransferData{}).
		Where("transfer_status = ? AND create_time < ?", Enum.TransferInit, date).
		Asc("create_time").Find(&data); err != nil {
		log.Error("get transfer Database Error", err)
		return data, err
	}
	return data, nil
}


//取所有的轉帳資料
func GetAllTransfer(engine *database.MysqlSession) ([]entity.TransferData,error) {
	var data []entity.TransferData
	if err := engine.Engine.Table(entity.TransferData{}).Select("*").Find(&data); err != nil {
		log.Error("get transfer Database Error", err)
		return data, err
	}
	return data, nil
}
