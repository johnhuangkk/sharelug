package Orders

import (
	"api/services/Enum"
	"api/services/VO/Request"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"api/services/util/tools"
	"strings"
	"time"
)

//新增 Bill Order
func InsertBillOrderData(engine *database.MysqlSession, data entity.BillOrderData) error {
	data.UpdateTime = time.Now()
	data.CreateTime = time.Now()
	if _, err := engine.Session.Table(entity.BillOrderData{}).Insert(&data); err != nil {
		log.Error("Insert Bill Order Database Error", err)
		return err
	}
	return nil
}

//取出訂單內容
func GetBillOrderByOrderId(engine *database.MysqlSession, OrderId string) (entity.BillOrderData, error) {
	var data entity.BillOrderData
	if _, err := engine.Engine.Table(entity.BillOrderData{}).Where("bill_id = ?", OrderId).Get(&data); err != nil {
		log.Error("Get Bill Order Database Error", err)
		return data, err
	}
	return data, nil
}

//取出訂單內容
func GetBillOrderByOrderIdAndUserId(engine *database.MysqlSession, OrderId, userId string) (entity.BillOrderData, error) {
	var data entity.BillOrderData
	if _, err := engine.Engine.Table(entity.BillOrderData{}).Where("bill_id = ?", OrderId).
		And("buyer_id = ?", userId).Get(&data); err != nil {
		log.Error("Get Bill Order Database Error", err)
		return data, err
	}
	return data, nil
}

//更新Order Data
func UpdateBillOrderData(engine *database.MysqlSession, data entity.BillOrderData) error {
	data.UpdateTime = time.Now()
	if _, err := engine.Session.Table(entity.BillOrderData{}).ID(data.BillId).AllCols().Update(data); err != nil {
		log.Error("Update Bill Order Database Error", err)
		return err
	}
	return nil
}

//取出訂單內容
func CountBillOrderByUserId(engine *database.MysqlSession, userId, tab string) (int64, error) {
	count, err := engine.Engine.Table(entity.BillOrderData{}).Where("buyer_id = ?", userId).
		And("bill_status = ?", tab).Count()
	if err != nil {
		log.Error("Get Bill Order Database Error", err)
		return count, err
	}
	return count, nil
}

func GetBillOrderListByUserId(engine *database.MysqlSession, BuyerId, tab string, limit, start int) ([]entity.BillOrderData, error) {
	limit = tools.CheckIsZero(limit, 10)
	start = tools.CheckIsZero(start, 1)
	start = (start - 1) * limit
	var data []entity.BillOrderData
	if err := engine.Engine.Table(entity.BillOrderData{}).Select("*").Where("buyer_id = ?", BuyerId).
		And("bill_status = ?", tab).Desc("create_time").Limit(limit, start).Find(&data); err != nil {
		return data, err
	}
	return data, nil
}

func CountAllBillOrderByUserId(engine *database.MysqlSession, userId string) (int64, error) {
	count, err := engine.Engine.Table(entity.BillOrderData{}).Where("buyer_id = ?", userId).Count()
	if err != nil {
		log.Error("Get Bill Order Database Error", err)
		return count, err
	}
	return count, nil
}

func GetAllBillListByUserId(engine *database.MysqlSession, BuyerId string, limit, start int) ([]entity.BillOrderData, error) {
	limit = tools.CheckIsZero(limit, 10)
	start = tools.CheckIsZero(start, 1)
	start = (start - 1) * limit
	var data []entity.BillOrderData
	if err := engine.Engine.Table(entity.BillOrderData{}).Select("*").Where("buyer_id = ?", BuyerId).
		Desc("create_time").Limit(limit, start).Find(&data); err != nil {
		return data, err
	}
	return data, nil
}

//取出過期買家帳單
func GetBillExpire(engine *database.MysqlSession, date string) ([]entity.BillOrderData, error) {
	var data []entity.BillOrderData
	err := engine.Engine.Table(entity.BillOrderData{}).Select("*").Where("bill_expire <= ?", date).
		And("bill_status = ?", Enum.BillStatusInit).Find(&data)
	if err != nil {
		log.Error("Select Appropriation Order Database Error", err)
		return data, err
	}
	return data, nil
}

func SearchBills(engine *database.MysqlSession, params Request.ErpSearchProductRequest) ([]entity.BillOrderData, error) {
	where, bind := ComposeSearchBillsParams(params)
	var data []entity.BillOrderData
	if err := engine.Engine.Table(entity.BillOrderData{}).Select("*").Where(strings.Join(where, " AND "), bind...).Find(&data); err != nil {
		log.Error("get products Database Error", err)
		return data, err
	}
	return data, nil
}

func ComposeSearchBillsParams(params Request.ErpSearchProductRequest) ([]string, []interface{}) {
	var where 	[]string
	var bind 	[]interface{}

	if len(params.ProductStatus) != 0 {
		where = append(where, "bill_status = ?")
		bind = append(bind, params.ProductStatus)
	}
	if len(params.CreateDate) != 0 {
		date := strings.Split(params.CreateDate, "-")
		where = append(where, "create_time BETWEEN ? AND ?")
		bind = append(bind, date[0])
		bind = append(bind, date[1])
	}
	if len(params.Amount) != 0 {
		price := strings.Split(params.Amount, "-")
		where = append(where, "total_amount BETWEEN ? AND ?")
		bind = append(bind, price[0])
		bind = append(bind, price[1])
	}
	if len(params.ProductName) != 0 {
		where = append(where, "product_name = ?")
		bind = append(bind, params.ProductName)
	}
	if len(params.ProductId) != 0 {
		where = append(where, "bill_id = ?")
		bind = append(bind, params.ProductId)
	}
	if len(params.ShipMode) != 0 {
		where = append(where, "ship_type = ?")
		bind = append(bind, params.ShipMode)
	}
	if len(params.PaymentMode) != 0 {
		where = append(where, "pay_way_type = ?")
		bind = append(bind, params.PaymentMode)
	}
	return where, bind
}