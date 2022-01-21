package Cvs

import (
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"fmt"
	"strings"
	"time"
)

// 新增對帳資訊
func InsertAccountingData(engine *database.MysqlSession, data entity.CvsAccountingData) error {
	data.UpdateTime = time.Now()
	data.CreateTime = time.Now()
	_, err := engine.Session.Table(entity.CvsAccountingData{}).Insert(&data)
	if err != nil {
		log.Error("InsertAccountingData Data [%s]", &data)
		log.Error("InsertAccountingData Error[%s]", err.Error())
		return fmt.Errorf("寫入失敗")
	}

	return nil
}

// 變更對帳資訊
func UpdateAccountingDataForClos(engine *database.MysqlSession, data entity.CvsAccountingData, condition map[string]interface{}, clos []string) error {
	data.UpdateTime = time.Now()
	clos = append(clos, `update_time`)
	_, err := engine.Session.Table(entity.CvsAccountingData{}).Cols(clos...).Where(condition).Update(data)

	if err != nil {
		log.Error("UpdateAccountingData Data [%s]", &data)
		log.Error("UpdateAccountingData Error[%s]", err.Error())
		return fmt.Errorf("寫入失敗")
	}

	return nil
}

// 超商ERP核帳
func GetCvsAccountingChecked(engine *database.MysqlSession, joinCondition string, whereCondition []string, whereBindCondition []interface{}) ([]entity.CvsShippingWithAmount, error) {
	var data []entity.CvsShippingWithAmount

	err := engine.Engine.Table(entity.CvsAccountingData{}).
		Select("*").
		Join("LEFT",
			"cvs_shipping_data",
			joinCondition).
		Where(strings.Join(whereCondition, " AND "), whereBindCondition...).
		Desc("cvs_accounting_data.create_time").
		Find(&data)

	if err != nil {
		log.Error("GetCvsAccountingChecked", err)
		return data, err
	}

	return data, nil
}
