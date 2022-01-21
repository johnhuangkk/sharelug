package Orders

import (
	"api/services/database"
	"api/services/entity"
)

//取出訂單明細內容
func GetOrderDetailByOrderId(engine *database.MysqlSession, orderId string) ([]entity.OrderDetail, error) {
	var OrderDetail []entity.OrderDetail
	err := engine.Engine.Table(entity.OrderDetail{}).
		Select("*").Where("order_id = ? ", orderId).Find(&OrderDetail)
	if err != nil {
		return OrderDetail, err
	}
	return OrderDetail, nil
}

func GetOrderDetailSingleByOrderId(engine *database.MysqlSession, orderId string) (entity.OrderDetail, error) {
	var OrderDetail entity.OrderDetail
	_, err := engine.Engine.Table(entity.OrderDetail{}).
		Select("*").Where("order_id = ? ", orderId).Get(&OrderDetail)
	if err != nil {
		return OrderDetail, err
	}
	return OrderDetail, nil
}

func GetOrderDetailListByOrderId(engine *database.MysqlSession, orderId string) ([]entity.OrderDetail, error) {
	var data []entity.OrderDetail
	if err := engine.Engine.Table(entity.OrderDetail{}).
		Select("*").Where("order_id = ? ", orderId).Desc("ship_merge").Find(&data); err != nil {
		return data, err
	}
	return data, nil
}