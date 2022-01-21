package Withdraw

import (
	"api/services/Enum"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"fmt"
	"time"
)

func InsertBankCode(data entity.BankCodeData) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()

	_, err := engine.Session.Table(entity.BankCodeData{}).Insert(&data)
	if err != nil {
		log.Error("Insert Bank Code Database Error", err)
		return err
	}
	return nil
}

func GetBankCode(engine *database.MysqlSession) ([]entity.BankCodeData, error) {
	var data []entity.BankCodeData
	sql := fmt.Sprintf("SELECT * FROM bank_code_data WHERE bank_status = ? ORDER BY id ASC")
	err := engine.Engine.SQL(sql, Enum.OrderSuccess).Find(&data)
	if err != nil {
		log.Error("get transfer Database Error", err)
		return data, err
	}
	return data, nil
}

func GetBankInfoByBranchCode(engine *database.MysqlSession, BranchCode string) (entity.BankCodeData, error) {
	var data entity.BankCodeData
	sql := fmt.Sprintf("SELECT * FROM bank_code_data WHERE branch_code = ?")
	_, err := engine.Engine.SQL(sql, BranchCode).Get(&data)
	if err != nil {
		log.Error("get transfer Database Error", err)
		return data, err
	}
	return data, nil
}

func GetBankInfoByBankCode(engine *database.MysqlSession, code string) (entity.BankCodeData, error) {
	var data entity.BankCodeData
	if _, err := engine.Engine.Table(entity.BankCodeData{}).Select("*").Where("bank_code = ?", code).Get(&data); err != nil {
		log.Error("get transfer Database Error", err)
		return data, err
	}
	return data, nil
}

func InsertIndustryData(data entity.IndustryData) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	data.CreateTime = time.Now()
	if _, err := engine.Session.Table(entity.IndustryData{}).Insert(&data); err != nil {
		log.Error("Insert Industry Database Error", err)
		return err
	}
	return nil
}

func GetIndustryData(engine *database.MysqlSession) ([]entity.IndustryData, error) {
	var data []entity.IndustryData
	if err := engine.Engine.Table(entity.IndustryData{}).Select("*").Asc("sort").Find(&data); err != nil {
		log.Error("get Industry Database Error", err)
		return data, err
	}
	return data, nil
}

func GetIndustryCode(engine *database.MysqlSession, code string) (entity.IndustryData, error) {
	var data entity.IndustryData
	if _, err := engine.Engine.Table(entity.IndustryData{}).Select("*").Where("industry_id = ?", code).Get(&data); err != nil {
		log.Error("get Industry Database Error", err)
		return data, err
	}
	return data, nil
}