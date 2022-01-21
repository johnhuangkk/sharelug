package Cvs

import (
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"api/services/util/tools"
	"fmt"
	"time"
)

// 新增超商配送資訊
func InsertCvsShippingData(engine *database.MysqlSession, data entity.CvsShippingData) error {
	_, err := engine.Session.Table(entity.CvsShippingData{}).Insert(&data)
	if err != nil {
		log.Error("InsertCvsShippingData Data [%s]", &data)
		log.Error("InsertCvsShippingData Error[%s]", err.Error())
		return fmt.Errorf("寫入失敗")
	}

	return nil
}

// 更新超商配送資訊
func UpdateCvsShippingData(engine *database.MysqlSession, data entity.CvsShippingData) error {
	data.UpdateTime = time.Now()
	_, err := engine.Session.Table(entity.CvsShippingData{}).ID(data.Id).Update(&data)
	if err != nil {
		log.Error("UpdateCvsShippingData Data [%s]", &data)
		log.Error("UpdateCvsShippingData Error[%s]", err.Error())
		return fmt.Errorf("更新失敗")
	}

	return nil
}

// 更新超商配送資訊 可選擇欄位
func UpdateCvsShippingDataForFields(engine *database.MysqlSession, data entity.CvsShippingData, clos []string) error {
	data.UpdateTime = time.Now()
	clos = append(clos, `update_time`)

	_, err := engine.Session.Table(entity.CvsShippingData{}).Cols(clos...).ID(data.Id).Update(&data)
	if err != nil {
		log.Error("UpdateCvsShippingDataForFields Data [%s]", &data)
		log.Error("UpdateCvsShippingDataForFields Error[%s]", err.Error())
		return fmt.Errorf("更新失敗")
	}

	return nil
}

// 取得超商配送資訊
func GetCvsShippingData(engine *database.MysqlSession, orderData entity.OrderData) entity.CvsShippingData {
	var data entity.CvsShippingData
	query := map[string]interface{}{}
	query["ec_order_no"] = orderData.OrderId
	query["cvs_type"] = orderData.ShipType

	log.Debug("query", query)

	_, _ = engine.Engine.Table(entity.CvsShippingData{}).Select("*").Where(query).Get(&data)
	return data
}

/**
取得貨運單號
*/
func GetCvsShippingDataByShipNo(engine *database.MysqlSession, shipNo []string) ([]entity.CvsShippingData, error)  {

	var data []entity.CvsShippingData

	shipNoInterface := tools.StringArrayToInterface(shipNo)

	err := engine.Engine.Table(entity.CvsShippingData{}).Select("*").In("ship_no", shipNoInterface...).Find(&data)
	if err != nil {
		log.Error("GetCvsShippingDataByShipNo Error", err)
		return data, err
	}

	return data, nil
}
