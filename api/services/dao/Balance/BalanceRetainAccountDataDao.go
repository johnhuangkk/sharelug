package Balance

import (
	"api/services/Enum"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"fmt"
	"time"
)

func InsertBalanceRetainAccountData(engine *database.MysqlSession, data entity.BalanceRetainAccountData) error {
	data.CreateTime = time.Now()
	_, err := engine.Session.Table(entity.BalanceRetainAccountData{}).Insert(&data)
	if err != nil {
		log.Error("Balance Retain Account Database Insert Error", err)
		return err
	}
	return nil
}

func CountBalanceRetainListByUserId(engine *database.MysqlSession, UserId string) (int64, error) {
	sql := fmt.Sprintf("SELECT count(*) FROM balance_retain_account_data a LEFT JOIN order_data b ON a.data_id = b.order_id WHERE a.user_id = ? AND a.trans_type = ?")
	count, err := engine.Engine.SQL(sql, UserId, Enum.BalanceTypeRetain).Count()
	if err != nil {
		log.Error("Count Balance Retain Database Error", err)
		return count, err
	}
	return count, nil
}



func GetBalanceRetainListByUserId(engine *database.MysqlSession, UserId string, limit, page int64 ) ([]entity.BalanceRetainByOrderData, error) {
	var data []entity.BalanceRetainByOrderData
	page = (page - 1) * limit
	sql := fmt.Sprintf("SELECT * FROM balance_retain_account_data a LEFT JOIN order_data b ON a.data_id = b.order_id WHERE a.user_id = ? AND a.trans_type = ? ORDER BY a.create_time DESC LIMIT %v OFFSET %v", limit, page)
	err := engine.Engine.SQL(sql, UserId, Enum.BalanceTypeRetain).Find(&data)
	if err != nil {
		log.Error("Get Balance Retain Database Error", err)
		return data, err
	}
	return data, nil
}

func GetBalanceRetainAccountLastByUserId(engine *database.MysqlSession, SellerId string) (entity.BalanceRetainAccountData, error) {
	var resp entity.BalanceRetainAccountData
	if _, err := engine.Engine.Table(entity.BalanceRetainAccountData{}).Select("*").
		Where("user_id = ?", SellerId).Desc("id").Get(&resp); err != nil {
		log.Error("Get Balance Account Database Error", err)
		return resp, err
	}
	return resp, nil
}

func GetBalanceRetainsByUserId(engine *database.MysqlSession, UserId string) ([]entity.BalanceRetainAccountData, error) {
	var data []entity.BalanceRetainAccountData
	if err := engine.Engine.Table(entity.BalanceRetainAccountData{}).Select("*").
		Where("user_id = ?", UserId).Asc("id").Find(&data); err != nil {
		log.Error("Get Balance Database Error", err)
		return data, err
	}
	return data, nil
}

func GetBalanceRetainById(engine *database.MysqlSession, Id int64) (entity.BalanceRetainAccountData, error) {
	var data entity.BalanceRetainAccountData
	if _, err := engine.Engine.Table(entity.BalanceRetainAccountData{}).Select("*").
		Where("id = ?", Id).Get(&data); err != nil {
		log.Error("Get Balance Database Error", err)
		return data, err
	}
	return data, nil
}

//更新 GwCreditData
func UpdateBalanceRetainData(engine *database.MysqlSession, data entity.BalanceRetainAccountData) error {
	if _, err := engine.Session.Table(entity.BalanceRetainAccountData{}).Where("id = ?", data.Id).AllCols().Update(data); err != nil {
		log.Error("Update Balance Retain Error", err)
		return err
	}
	return nil
}
