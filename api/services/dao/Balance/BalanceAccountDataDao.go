package Balance

import (
	"api/services/Enum"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"fmt"
	"time"
)

func InsertBalanceAccountData(engine *database.MysqlSession, data entity.BalanceAccountData) error {
	data.CreateTime = time.Now()
	if _, err := engine.Session.Table(entity.BalanceAccountData{}).Insert(&data); err != nil {
		log.Error("Balance Account Database Insert Error", err)
		return err
	}
	return nil
}


func GetBalanceAccountLastByUserId(engine *database.MysqlSession, UserId string) (entity.BalanceAccountData, error) {
	var resp entity.BalanceAccountData
	if _, err := engine.Engine.Table(entity.BalanceAccountData{}).Select("*").
		Where("user_id = ?", UserId).Desc("id").Get(&resp); err != nil {
		log.Error("Get Balance Account Database Error", err)
		return resp, err
	}
	return resp, nil
}

func CountBalanceListByUserId(engine *database.MysqlSession, UserId, start, end string) (int64, error) {
	sql := fmt.Sprintf("SELECT count(*) FROM balance_account_data WHERE user_id = ? AND create_time BETWEEN ? AND ?")
	count, err := engine.Engine.SQL(sql, UserId, start, end).Count()
	if err != nil {
		log.Error("Count Balance Database Error", err)
		return count, err
	}
	return count, nil
}

func GetBalanceListNoLimitByUserId(engine *database.MysqlSession, UserId, start, end string) ([]entity.BalanceAccountData, error) {
	var data []entity.BalanceAccountData
	if err := engine.Engine.Table(entity.BalanceAccountData{}).Select("*").
		Where("user_id = ? AND create_time BETWEEN ? AND ? ", UserId, start, end).Desc("id").Find(&data); err != nil {
		log.Error("Get Balance Database Error", err)
		return data, err
	}
	return data, nil
}

func GetBalanceListByUserId(engine *database.MysqlSession, UserId, start, end string, limit, page int64 ) ([]entity.BalanceAccountData, error) {
	var data []entity.BalanceAccountData
	page = (page - 1) * limit
	sql := fmt.Sprintf("SELECT * FROM balance_account_data WHERE user_id = ? AND create_time BETWEEN ? AND ? ORDER BY id DESC LIMIT %v OFFSET %v", limit, page)
	err := engine.Engine.SQL(sql, UserId, start, end).Find(&data)
	if err != nil {
		log.Error("Get Balance Database Error", err)
		return data, err
	}
	return data, nil
}

func GetBalancesByDateAndUserId(engine *database.MysqlSession, UserId, start, end string) ([]entity.BalanceAccountData, error) {
	var data []entity.BalanceAccountData
	err := engine.Engine.Table(entity.BalanceAccountData{}).Select("*").
		Where("user_id = ? AND create_time BETWEEN ? AND ?", UserId, start, end).Desc("id").Find(&data)
	if err != nil {
		log.Error("Get Balance Database Error", err)
		return data, err
	}
	return data, nil
}



func CountBalanceByUserIdAndOrderId(engine *database.MysqlSession, UserId, OrderId string) (int64, error) {
	count, err := engine.Engine.Table(entity.BalanceAccountData{}).Select("*").
		Where("user_id = ?", UserId).And("data_id = ?", OrderId).And("trans_type = ?", Enum.BalanceTypePlatform).Count()
	if err != nil {
		log.Error("Get Balance Database Error", err)
		return count, err
	}
	return count, nil
}

func GetBalancesByUserId(engine *database.MysqlSession, UserId string) ([]entity.BalanceAccountData, error) {
	var data []entity.BalanceAccountData
	if err := engine.Engine.Table(entity.BalanceAccountData{}).Select("*").
	 	Where("user_id = ?", UserId).Asc("id").Find(&data); err != nil {
		log.Error("Get Balance Database Error", err)
		return data, err
	}
	return data, nil
}

func GetBalanceById(engine *database.MysqlSession, Id int64) (entity.BalanceAccountData, error) {
	var data entity.BalanceAccountData
	if _, err := engine.Engine.Table(entity.BalanceAccountData{}).Select("*").
		Where("id = ?", Id).Get(&data); err != nil {
		log.Error("Get Balance Database Error", err)
		return data, err
	}
	return data, nil
}


//更新 GwCreditData
func UpdateBalancesData(engine *database.MysqlSession, data entity.BalanceAccountData) error {
	if _, err := engine.Session.Table(entity.BalanceAccountData{}).Where("id = ?", data.Id).Update(data); err != nil {
		return err
	}
	return nil
}
