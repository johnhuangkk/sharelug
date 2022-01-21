package postBag

import (
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"time"
)

type OrderPostBagConsignment struct {
	Order   entity.OrderData              `xorm:"extends"`
	PostBag entity.PostBagConsignmentData `xorm:"extends"`
	Store   entity.StoreData              `xorm:"extends"`
}
type PostBagOrder struct {
	Order   entity.OrderData              `xorm:"extends"`
	PostBag entity.PostBagConsignmentData `xorm:"extends"`
}

func InsertPostBagConsignmentData(engine *database.MysqlSession, data entity.PostBagConsignmentData) (entity.PostBagConsignmentData, error) {
	_, err := engine.Session.Table(entity.PostBagConsignmentData{}).Insert(&data)
	if err != nil {
		log.Error("Database Error", err)
		return data, err
	}
	return data, nil
}
func UpdatePostBagConsignmentData(numbers []string, fileName string, nTime time.Time) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	if err := engine.Session.Begin(); err != nil {
		log.Error("Session Begin Error", err)
		return err
	}
	_, err := engine.Session.Table(entity.PostBagConsignmentData{}).In(`ship_number`, numbers).Update(&entity.PostBagConsignmentData{
		FileName:   fileName,
		VerifyTime: nTime,
	})
	if err != nil {
		if err := engine.Session.Rollback(); err != nil {
			log.Error("Rollback Error")
		}
		return err
	}
	if err := engine.Session.Commit(); err != nil {
		log.Error("Session Commit Error", err)
		return err
	}
	return nil
}
func FindRecentNonVerifyFilePostBagConsignmentData() ([]PostBagOrder, error) {
	var data []PostBagOrder
	engine := database.GetMysqlEngineGroup().Engine
	defer engine.Close()
	err := engine.Table(entity.PostBagConsignmentData{}).
		Select("post_bag_consignment_data.order_id,post_bag_consignment_data.ship_number,post_bag_consignment_data.verify_file_name,post_bag_consignment_data.file_name,post_bag_consignment_data.verify_time").
		Where(`post_bag_consignment_data.verify_file_name = ?`, "").
		Where(`post_bag_consignment_data.verify_time != ?`, "").
		Find(&data)
	if err != nil {
		log.Error("post bag Database Error", err)
		return data, err
	}

	return data, nil
}
func FindRecentNonVerifyPostBagConsignmentData() ([]PostBagOrder, error) {
	var data []PostBagOrder
	engine := database.GetMysqlEngineGroup().Engine
	defer engine.Close()
	err := engine.Table(entity.PostBagConsignmentData{}).
		Select("post_bag_consignment_data.order_id,post_bag_consignment_data.ship_number,order_data.ship_expire,order_data.order_id,order_data.ship_number").
		Join(`INNER`, entity.OrderData{}, "order_data.ship_number = post_bag_consignment_data.ship_number").
		Where(`post_bag_consignment_data.file_name = ?`, "").
		Where(`order_data.ship_status = ?`, "TAKE").
		Find(&data)
	if err != nil {
		log.Error("Database Error", err)
		return data, err
	}

	return data, nil
}

func GetConsignmentData(orderIds []string) ([]OrderPostBagConsignment, error) {
	engine := database.GetMysqlEngineGroup().Engine
	defer engine.Close()
	var datas []OrderPostBagConsignment
	err := engine.Table(entity.OrderData{}).
		Select("order_data.order_id,order_data.receiver_name,order_data.receiver_address,order_data.receiver_phone,order_data.ship_type,order_data.store_id,post_bag_consignment_data.order_id,post_bag_consignment_data.merchant_id,post_bag_consignment_data.ship_number,post_bag_consignment_data.seller_name,post_bag_consignment_data.seller_phone,post_bag_consignment_data.seller_zip,post_bag_consignment_data.seller_addr,store_data.store_id,store_data.store_name").
		Join(`INNER`, entity.PostBagConsignmentData{}, "post_bag_consignment_data.order_id = order_data.order_id").
		Join(`INNER`, entity.StoreData{}, "store_data.store_id = order_data.store_id").
		In("order_data.order_id", orderIds).Find(&datas)
	if err != nil {
		log.Error("Database PostBag Find ConsignmentData Error", err.Error())
		return datas, err
	}

	return datas, nil
}

func GetConsignmentByOrderId(engine *database.MysqlSession, orderId string) (entity.PostBagConsignmentData, error) {
	var data entity.PostBagConsignmentData
	if _, err := engine.Engine.Table(entity.PostBagConsignmentData{}).Select("*").Where("order_id = ?", orderId).Get(&data); err != nil {
		log.Error("Get PostBag Find Consignment Database Error", err)
		return data, err
	}
	return data, nil
}
