package Orders

import (
	"api/services/Enum"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"fmt"
	"strings"
	"time"
)

//取出最後一筆退貨單
func GetOrderReturnLastByOrderId(engine *database.MysqlSession, OrderId string, ProductSpecId string) (entity.OrderRefundData, error) {
	var data entity.OrderRefundData
	if _, err := engine.Engine.Table(entity.OrderRefundData{}).Select("*").
		Where("refund_type = ? AND order_id = ? AND product_spec_id = ?", Enum.TypeReturn, OrderId, ProductSpecId).
		Desc("create_time").Get(&data); err != nil {
		log.Error("Get Order refund data Database Error", err)
		return data, err
	}
	return data, nil
}
//新增一筆退貨單
func InsertOrderReturnData(engine *database.MysqlSession, data entity.OrderRefundData) error {
	data.RefundType = Enum.TypeReturn
	data.CreateTime = time.Now()
	_, err := engine.Session.Table(entity.OrderRefundData{}).Insert(&data)
	if err != nil {
		log.Error("Order Refund Data Database Insert Error", err)
		return err
	}
	return nil
}
//取出退貨單 BY OrderId And SpecId
func GetOrderReturnByOrderIdAndSpecId(engine *database.MysqlSession, OrderId string, ProductSpecId string) ([]entity.OrderRefundData, error) {
	var resp []entity.OrderRefundData
	if err := engine.Engine.Table(entity.OrderRefundData{}).Select("*").
		Where("refund_type = ? AND order_id = ? AND product_spec_id = ?", Enum.TypeReturn, OrderId, ProductSpecId).
		Desc("create_time").Find(&resp); err != nil {
		log.Error("Get Order refund data Database Error", err)
		return resp, err
	}
	return resp, nil
}
//取出退貨單 BY ReturnId
func GetOrderReturnByReturnId(engine *database.MysqlSession, ReturnId string) (entity.OrderRefundData, error) {
	var resp entity.OrderRefundData
	if _, err := engine.Engine.Table(entity.OrderRefundData{}).Select("*").
		Where("refund_type = ? AND refund_id = ?", Enum.TypeReturn, ReturnId).Get(&resp); err != nil {
		log.Error("Get Order refund data Database Error", err)
		return resp, err
	}
	return resp, nil
}
//更新Order Data
func UpdateReturnData(engine *database.MysqlSession, data entity.OrderRefundData) (int64, error) {
	data.UpdateTime = time.Now()
	affected, err := engine.Session.Table(entity.OrderRefundData{}).ID(data.RefundId).Update(data)
	if err != nil {
		return affected, err
	}
	return affected, nil
}
//取出退貨單列表
func GetRefundList(engine *database.MysqlSession, where []string, bind []interface{}, limit int, start int) ([]entity.OrderResp, error) {
	var data []entity.OrderResp
	if err := engine.Engine.Table(entity.OrderData{}).
		Join("RIGHT", entity.OrderRefundData{}, "order_data.order_id = order_refund_data.order_id").
		Select("*").Where(strings.Join(where, " AND "), bind...).Desc("order_data.create_time").
		Limit(limit, start).Find(&data); err != nil {
		log.Error("Get Order Database Error", err)
		return data, err
	}
	return data, nil
}
//計算退貨退款數
func GetRefundCount(engine *database.MysqlSession, StoreId, RefundType, RefundStatus string) int64 {
	var where 	[]string
	var bind 	[]interface{}
	where = append(where, "store_id = ?")
	bind = append(bind, StoreId)

	if RefundType != "" {
		where = append(where, "refund_type = ?")
		bind = append(bind, RefundType)
		where = append(where, "status = ?")
		bind = append(bind, RefundStatus)
	}
	sql := fmt.Sprintf("SELECT count(*) FROM order_data orders RIGHT JOIN order_refund_data refund ON orders.order_id = refund.order_id " +
		"WHERE %s", strings.Join(where, " AND "))
	result, err := engine.Engine.SQL(sql, bind...).Count()
	if err != nil {
		log.Error("count refund Order Database Error", err)
		return 0
	}
	return result
}

func CountOrderReturnByOrderId(engine *database.MysqlSession, OrderId string) (int64, error) {
	result, err := engine.Engine.Table(entity.OrderRefundData{}).Select("count(*)").
		Where("refund_type = ? AND order_id = ?", Enum.TypeReturn, OrderId).Count()
	if err != nil {
		log.Error("Get Order refund data Database Error", err)
		return result, err
	}
	return result, err
}
