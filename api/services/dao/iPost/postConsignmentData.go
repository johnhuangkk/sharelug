package iPost

import (
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"api/services/util/tools"
)

func InsertPostConsignmentData(engine *database.MysqlSession, data entity.PostConsignmentData) (entity.PostConsignmentData, error) {
	data.CreateTime = tools.Now("YmdHis")
	_, err := engine.Session.Table(entity.PostConsignmentData{}).Insert(&data)
	if err != nil {
		log.Error("Database Error", err)
		return data, err
	}
	return data, nil
}

func QueryPostConsignmentData(engine *database.MysqlSession, orderId []string) ([]entity.ConsignmentNote, error) {
	var consignments []entity.ConsignmentNote

	err := engine.Engine.Table(entity.PostConsignmentData{}).
		Select("*").
		Join("Left", "order_data", "order_data.order_id = post_consignment_data.order_id And order_data.ship_number = post_consignment_data.ship_number").
		In("order_data.order_id", orderId).
		Desc("order_data.order_id").
		Find(&consignments)

	log.Info("consignments:", consignments)

	if err != nil {
		log.Error("Database Error", err)
		return consignments, err
	}
	return consignments, nil

}