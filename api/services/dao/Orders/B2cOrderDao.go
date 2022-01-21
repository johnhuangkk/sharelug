package Orders

import (
	"api/services/Enum"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"fmt"
	"time"
)

func InsertB2cOrderData(engine *database.MysqlSession, data entity.B2cOrderData) error {
	data.CreateTime = time.Now()
	data.UpdateTime = time.Now()
	_, err := engine.Session.Table(entity.B2cOrderData{}).Insert(&data)
	if err != nil {
		log.Error("Database Error", err)
		return err
	}
	return nil
}

//取出訂單內容
func GetB2cOrderByOrderId(engine *database.MysqlSession, OrderId string) (entity.B2cOrderData, error) {
	var data entity.B2cOrderData
	 if _, err := engine.Engine.Table(entity.B2cOrderData{}).Select("*").
	 	Where("order_id = ?", OrderId).Get(&data); err != nil {
		return data, err
	}
	return data, nil
}

//取出最後一筆成功訂單
func GetB2cOrderLastSuccessByUserId(engine *database.MysqlSession, UserId string) (entity.B2cOrderData, error) {
	var data entity.B2cOrderData
	if _, err := engine.Engine.Table(entity.B2cOrderData{}).Select("*").
		Where("user_id = ? AND order_status = ?", UserId, Enum.OrderSuccess).Desc("create_time").Get(&data); err != nil {
		log.Error("Get B2c Order Database Error", err)
		return data, err
	}
	return data, nil
}
//統計未付款帳單
func CountB2cUnpaidBillsByUserId(engine *database.MysqlSession, UserId string) (int64, error) {
	count, err := engine.Engine.Table(entity.B2cBillingData{}).Where("user_id = ? AND billing_status = ?",  UserId, Enum.OrderWait).Count()
	if err != nil {
		log.Error("Get B2c Order Database Error", err)
		return 0, err
	}
	return count, err
}

//統計未付款完成
func CountB2cUnpaidOrdersByUserId(engine *database.MysqlSession, UserId string) (int64, error) {
	sql := fmt.Sprintf("SELECT count(*) FROM b2c_order_data WHERE user_id = ? AND order_sys = ? AND order_status = ?")
	count, err := engine.Engine.SQL(sql, UserId, 1, Enum.OrderWait).Count()
	if err != nil {
		log.Error("Get B2c Order Database Error", err)
		return 0, err
	}
	return count, err
}

//取未付款完成訂單
func GetB2cUnpaidOrdersByUserId(engine *database.MysqlSession, UserId string) ([]entity.B2cOrderData, error) {
	var data []entity.B2cOrderData
	sql := fmt.Sprintf("SELECT * FROM b2c_order_data WHERE user_id = ? AND order_sys = ? AND order_status = ? ORDER BY create_time ASC")
	err := engine.Engine.SQL(sql, UserId, 1, Enum.OrderWait).Find(&data)
	if err != nil {
		log.Error("Get B2c Order Database Error", err)
		return data, err
	}
	return data, err
}

//更新Order Data
func UpdateB2cOrderData(engine *database.MysqlSession, OrderData entity.B2cOrderData) error {
	_, err := engine.Session.Table("b2c_order_data").ID(OrderData.OrderId).AllCols().Update(OrderData)
	if err != nil {
		return err
	}
	return nil
}

//取出訂單內容
func GetB2cOrderTransferByUserId(engine *database.MysqlSession, userId string) (entity.B2cOrderData, error) {
	var OrderData entity.B2cOrderData
	sql := fmt.Sprintf("SELECT * FROM b2c_order_data WHERE user_id = ? AND order_status = ? AND payment = ? ORDER BY create_time DESC LIMIT 0,1")
	_, err := engine.Engine.SQL(sql, userId, Enum.OrderWait, Enum.Transfer).Get(&OrderData)
	if err != nil {
		return OrderData, err
	}
	return OrderData, nil
}
//取出待付款的帳單
func GetB2cOrdersByUserIdAndExpire(engine *database.MysqlSession, UserId, StoreId string) (entity.B2cOrderData, error) {
	var data entity.B2cOrderData
	if _, err := engine.Engine.Table(entity.B2cOrderData{}).Select("*").
		Where("user_id = ? AND store_id = ? AND order_sys = ? AND order_status = ?", UserId, StoreId, 1, Enum.OrderWait).
		Desc("create_time").Get(&data); err != nil {
		log.Error("Get B2c Order Database Error", err)
		return data, err
	}
	return data, nil
}
//取出未開發票的訂單
func GetB2cOrderByNotInvoice(engine *database.MysqlSession) ([]entity.B2cOrderData, error) {
	var data []entity.B2cOrderData
	if err := engine.Engine.Table(entity.B2cOrderData{}).Where("order_status = ?", Enum.OrderSuccess).And("ask_invoice = ?", 1).
		And("invoice_status = ?", Enum.InvoiceOpenStatusNot).And("payment != ?", " ") .Find(&data); err != nil {
		log.Error("Get B2c Order Database Error", err)
		return data, err
	}
	return data, nil
}
//取出最早的訂單
func GetUpgradeProductDataByOrderSys(engine *database.MysqlSession, userId, sys string) (entity.B2cOrderData, error) {
	var data entity.B2cOrderData
	if _, err := engine.Engine.Table(entity.B2cOrderData{}).Select("*").
		Where("user_id = ? AND order_sys = ? AND order_status = ?", userId, sys, Enum.OrderSuccess) .Desc("create_time").Get(&data); err != nil {
		log.Error("Get Upgrade Product Database Error", err)
		return data, err
	}
	return data, nil
}

func SumB2cWaitOrdersByUserId(engine *database.MysqlSession, UserId string) (int64, error) {
	var data entity.B2cOrderData
	sum, err := engine.Engine.Where("user_id = ? AND order_sys = ? AND order_status = ?", UserId, 1, Enum.OrderWait).Sum(data, "amount")
	if err != nil {
		log.Error("Get B2c Order Database Error", err)
		return int64(sum), err
	}
	return int64(sum), err
}
//取出最後的訂單
func GetUpgradeProductDataByOrder(engine *database.MysqlSession, userId string) (entity.B2cOrderData, error) {
	var data entity.B2cOrderData
	if _, err := engine.Engine.Table(entity.B2cOrderData{}).Select("*").
		Where("user_id = ? AND order_status = ?", userId, Enum.OrderSuccess) .Desc("create_time").Get(&data); err != nil {
		log.Error("Get Upgrade Product Database Error", err)
		return data, err
	}
	return data, nil
}
//取出B2C訂單中的帳單
func GetB2cBillOrders(engine *database.MysqlSession) ([]entity.B2cOrderData, error) {
	var data []entity.B2cOrderData
	sql := fmt.Sprintf("SELECT * FROM b2c_order_data WHERE order_sys = ? ORDER BY create_time ASC")
	err := engine.Engine.SQL(sql,1).Find(&data)
	if err != nil {
		log.Error("Get B2c Order Database Error", err)
		return data, err
	}
	return data, nil
}
//新增一筆B2C帳單
func InsertB2cBillData(engine *database.MysqlSession, data entity.B2cBillingData) error {
	data.CreateTime = time.Now()
	data.UpdateTime = time.Now()
	_, err := engine.Session.Table(entity.B2cBillingData{}).Insert(&data)
	if err != nil {
		log.Error("Database Error", err)
		return err
	}
	return nil
}
//取未付款帳單
func GetB2cUnpaidBillByUserId(engine *database.MysqlSession, UserId string) ([]entity.B2cBillingData, error) {
	var data []entity.B2cBillingData
	err := engine.Engine.Table(entity.B2cBillingData{}).Select("*").
		Where("user_id = ? AND billing_status = ?",  UserId, Enum.OrderWait).Asc("create_time").Find(&data)
	if err != nil {
		log.Error("Get B2c Bill Database Error", err)
		return data, err
	}
	return data, err
}
//取出待付款的帳單
func GetB2cBillingByUserIdAndExpire(engine *database.MysqlSession, UserId, StoreId string) (entity.B2cBillingData, error) {
	var data entity.B2cBillingData
	if _, err := engine.Engine.Table(entity.B2cBillingData{}).Select("*").
		Where("user_id = ? AND store_id = ? AND billing_status = ?", UserId, StoreId, Enum.OrderWait).
		Desc("create_time").Get(&data); err != nil {
		log.Error("Get B2c Billing Database Error", err)
		return data, err
	}
	return data, nil
}
//取出帳單內容
func GetB2cBillByOrderId(engine *database.MysqlSession, OrderId string) (entity.B2cBillingData, error) {
	var data entity.B2cBillingData
	if _, err := engine.Engine.Table(entity.B2cBillingData{}).Select("*").
		Where("billing_id = ?", OrderId).Get(&data); err != nil {
		return data, err
	}
	return data, nil
}
//取未付款帳單
func SumB2cUnpaidBillByUserId(engine *database.MysqlSession, UserId string) (int64, error) {
	var data entity.B2cBillingData
	sum, err := engine.Engine.Table(entity.B2cBillingData{}).
		Where("user_id = ? AND billing_status = ?",  UserId, Enum.OrderWait).SumInt(data,"amount")
	if err != nil {
		log.Error("Get B2c Bill Database Error", err)
		return 0, err
	}
	return sum, err
}
//更新帳單
func UpdateB2cBillData(engine *database.MysqlSession, data entity.B2cBillingData) error {
	data.UpdateTime = time.Now()
	_, err := engine.Session.Table(entity.B2cBillingData{}).ID(data.BillingId).AllCols().Update(data)
	if err != nil {
		return err
	}
	return nil
}
