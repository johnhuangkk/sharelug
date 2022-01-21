package Orders

import (
	"api/services/Enum"
	"api/services/database"
	entity "api/services/entity"
	"api/services/util/log"
	"fmt"
	"time"
)

//取這訂單最後一筆退款
func GetOrderRefundLastByOrderId(engine *database.MysqlSession, OrderId string) (entity.OrderRefundData, error) {
	var resp entity.OrderRefundData
	sql := fmt.Sprintf("SELECT * FROM order_refund_data WHERE refund_type = ? AND order_id = ? ORDER BY create_time DESC LIMIT 0,1")
	_, err := engine.Engine.SQL(sql, Enum.TypeRefund, OrderId).Get(&resp)
	if err != nil {
		log.Error("Get Order refund data Database Error", err)
		return resp, err
	}
	return resp, nil
}

//退款列表
func GetOrderRefundsByOrderId(engine *database.MysqlSession, OrderId string) ([]entity.OrderRefundData, error) {
	var resp []entity.OrderRefundData
	sql := fmt.Sprintf("SELECT * FROM order_refund_data WHERE refund_type = ? AND order_id = ? ORDER BY create_time DESC")
	err := engine.Engine.SQL(sql, Enum.TypeRefund, OrderId).Find(&resp)
	if err != nil {
		log.Error("Get Order refund data Database Error", err)
		return resp, err
	}
	return resp, nil
}

func InsertOrderRefundData(engine *database.MysqlSession, data entity.OrderRefundData) error {
	data.RefundType = Enum.TypeRefund
	data.CreateTime = time.Now()
	_, err := engine.Session.Table(entity.OrderRefundData{}).Insert(&data)
	if err != nil {
		log.Error("Order Refund Data Database Insert Error", err)
		return err
	}
	return nil
}


func CountOrderRefundByOrderId(engine *database.MysqlSession, OrderId string) (int64, error) {
	sql := fmt.Sprintf("SELECT count(*) FROM order_refund_data WHERE refund_type = ? AND order_id = ?")
	result, err := engine.Engine.SQL(sql, Enum.TypeRefund, OrderId).Count()
	if err != nil {
		log.Error("Get Order refund data Database Error", err)
		return result, err
	}
	return result, nil
}

func GetOrderRefundByRefundId(engine *database.MysqlSession, refundId string) (entity.OrderRefundData, error) {
	var resp entity.OrderRefundData
	sql := fmt.Sprintf("SELECT * FROM order_refund_data WHERE refund_id = ?")
	_, err := engine.Engine.SQL(sql, refundId).Get(&resp)
	if err != nil {
		log.Error("Get Order refund data Database Error", err)
		return resp, err
	}
	return resp, nil
}

func GetOrderRefundAllByOrderId(engine *database.MysqlSession, refundId string) ([]entity.OrderRefundData, error) {
	var resp []entity.OrderRefundData
	sql := fmt.Sprintf("SELECT * FROM order_refund_data WHERE refund_type = ? AND order_id = ?")
	err := engine.Engine.SQL(sql, Enum.TypeRefund, refundId).Find(&resp)
	if err != nil {
		log.Error("Get Order refund data Database Error", err)
		return resp, err
	}
	return resp, nil
}

func GetOrderReturnAllByOrderId(engine *database.MysqlSession, refundId string) ([]entity.OrderRefundData, error) {
	var resp []entity.OrderRefundData
	sql := fmt.Sprintf("SELECT * FROM order_refund_data WHERE refund_type = ? AND order_id = ?")
	err := engine.Engine.SQL(sql, Enum.TypeReturn, refundId).Find(&resp)
	if err != nil {
		log.Error("Get Order refund data Database Error", err)
		return resp, err
	}
	return resp, nil
}

func GetReturnAndRefundByOrderId(engine *database.MysqlSession, orderId string) ([]entity.OrderRefundData, error) {
	var data []entity.OrderRefundData
	if err := engine.Engine.Table(entity.OrderRefundData{}).Select("*").Where("order_id = ?", orderId).Find(&data); err != nil {
		log.Error("Get Order refund Database Error", err)
		return data, err
	}
	return data, nil
}

func GetOrderRefundsByOrderIdAndStatus(engine *database.MysqlSession, OrderId string) (entity.OrderRefundData, error) {
	var data entity.OrderRefundData
	if _, err := engine.Engine.Table(entity.OrderRefundData{}).Select("*").
		Where("order_id = ? AND status = ?", OrderId, Enum.RefundStatusAudit).Asc("refund_time").Get(&data); err != nil {
		log.Error("Get Order refund data Database Error", err)
		return data, err
	}
	return data, nil
}


func UpdateOrderRefundData(engine *database.MysqlSession, data entity.OrderRefundData) error {
	data.UpdateTime = time.Now()
	_, err := engine.Session.Table(entity.OrderRefundData{}).ID(data.RefundId).AllCols().Update(&data)
	if err != nil {
		log.Error("Order Refund Data Database Insert Error", err)
		return err
	}
	return nil
}
