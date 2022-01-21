package Orders

import (
	"api/services/Enum"
	"api/services/VO/Request"
	"api/services/dao/member"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"api/services/util/tools"
	"fmt"
	"strings"
	"time"
)

//新增 Order Data
func InsertOrderData(engine *database.MysqlSession, data entity.OrderData) (entity.OrderData, error) {
	data.CreateTime = time.Now()
	data.UpdateTime = time.Now()
	if _, err := engine.Session.Table(entity.OrderData{}).Insert(&data); err != nil {
		log.Error("Database Error", err)
		return data, err
	}
	return data, nil
}

//更新Order Data
func UpdateOrderData(engine *database.MysqlSession, OrderId string, OrderData entity.OrderData) (int64, error) {
	OrderData.UpdateTime = time.Now()
	affected, err := engine.Session.Table(entity.OrderData{}).ID(OrderId).AllCols().Update(OrderData)
	if err != nil {
		return affected, err
	}
	return affected, nil
}

func UpdateOrdersData(engine *database.MysqlSession, OrderData *entity.OrderData) error {
	OrderData.UpdateTime = time.Now()
	if _, err := engine.Session.Table(entity.OrderData{}).ID(OrderData.OrderId).AllCols().Update(OrderData); err != nil {
		return err
	}
	return nil
}

//新增 Order Data Detail
func InsertOrderDetail(engine *database.MysqlSession, data entity.OrderDetail) error {
	data.CreateTime = time.Now()
	_, err := engine.Session.Table("order_detail").Insert(&data)
	if err != nil {
		log.Error("Database Error", err)
		return err
	}
	return nil
}

//賣家取出訂單內容
func GetOrderByOrderIdAndStoreId(engine *database.MysqlSession, OrderId, StoreId string) (entity.OrderData, error) {
	var OrderData entity.OrderData
	if _, err := engine.Engine.Table(entity.OrderData{}).Select("*").
		Where("order_id = ? AND store_id = ?", OrderId, StoreId).Get(&OrderData); err != nil {
		return OrderData, err
	}
	return OrderData, nil
}

//買家取出訂單內容
func GetOrderByOrderIdAndBuyerId(engine *database.MysqlSession, OrderId, BuyerId string) (entity.OrderData, error) {
	var OrderData entity.OrderData
	if _, err := engine.Engine.Table(entity.OrderData{}).
		Where("order_id = ? AND buyer_id = ?", OrderId, BuyerId).Get(&OrderData); err != nil {
		return OrderData, err
	}
	return OrderData, nil
}

//賣家取出指定訂單
func GetOrderByOrderIds(engine *database.MysqlSession, storeId string, OrderIds []string) ([]entity.OrderData, error) {
	var data []entity.OrderData
	if err := engine.Engine.Table(entity.OrderData{}).Select("*").Where("store_id = ?", storeId).
		In("order_id", OrderIds).Find(&data); err != nil {
		log.Error("Get Order Database Error", err)
		return data, err
	}
	return data, nil
}

func GetOrderByStatus(engine *database.MysqlSession, storeId, Status, PayWay string) ([]entity.OrderData, error) {
	var data []entity.OrderData
	if err := engine.Engine.Table(entity.OrderData{}).Select("*").
		Where("store_id = ? AND order_status = ? AND ship_type = ? AND ship_status = ?", storeId, Enum.OrderSuccess, PayWay, Status).
		Find(&data); err != nil {
		log.Error("Get Order Database Error", err)
		return data, err
	}
	return data, nil
}


//賣家取出指定訂單
func GetOrderAndDetailByOrderIds(engine *database.MysqlSession, storeId, orderBy string, OrderIds []string) ([]entity.OrderDetailData, error) {
	by := "order_data.create_time"
	switch orderBy {
	case "name":
		by = "order_detail.product_spec_id"
	case "price":
		by = "order_detail.product_price"
	case "code":
		by = "order_data.receiver_address"
	}
	var data []entity.OrderDetailData
	if err := engine.Engine.Table(entity.OrderData{}).Select("*").
		Join("LEFT", entity.OrderDetail{}, "order_data.order_id = order_detail.order_id").
		Where("order_data.store_id = ?", storeId).
		In("order_data.order_id", OrderIds).Desc(by).Find(&data); err != nil {
		return data, err
	}
	return data, nil
}

func GetOrderAndDetailByStatus(engine *database.MysqlSession, storeId, orderBy, Status, PayWay string) ([]entity.OrderDetailData, error) {
	by := "order_data.create_time"
	switch orderBy {
	case "name":
		by = "order_detail.product_spec_id"
	case "price":
		by = "order_detail.product_price"
	case "code":
		by = "order_data.receiver_address"
	}
	var data []entity.OrderDetailData
	if err := engine.Engine.Table(entity.OrderData{}).Select("*").
		Join("LEFT", entity.OrderDetail{}, "order_data.order_id = order_detail.order_id").
		Where("order_data.store_id = ? AND order_data.order_status = ? AND order_data.ship_type = ? AND order_data.ship_status = ?", storeId, Enum.OrderSuccess, PayWay, Status).
		Desc(by).Find(&data); err != nil {
		return data, err
	}
	return data, nil
}



func GetOrderAndDetailByOrderId(engine *database.MysqlSession, storeId, OrderId, orderBy string) ([]entity.OrderDetailData, error) {
	var data []entity.OrderDetailData
	by := "order_data.create_time"
	switch orderBy {
	case "name":
		by = "order_detail.product_spec_id"
	case "price":
		by = "order_detail.product_price"
	case "code":
		by = "order_data.receiver_address"
	}
	if err := engine.Engine.Table(entity.OrderData{}).Select("*").
		Join("LEFT", entity.OrderDetail{}, "order_detail.order_id = order_data.order_id").
		Where("order_data.store_id = ? AND order_data.order_id = ?", storeId, OrderId).
		Desc(by).Find(&data); err != nil {
		return data, err
	}
	return data, nil
}

//取出訂單內容
func GetOrderByOrderId(engine *database.MysqlSession, OrderId string) (entity.OrderData, error) {
	var OrderData entity.OrderData
	if _, err := engine.Engine.Table(entity.OrderData{}).Select("*").
		Where("order_id = ?", OrderId).Get(&OrderData); err != nil {
		log.Error("Get Order Database Error", err)
		return OrderData, err
	}
	return OrderData, nil
}

// 物流確認訂單狀態 by 物流編號
func GetOrderDataByShip(engine *database.MysqlSession, shipNumber, shipType string) (entity.OrderData, error) {
	var data entity.OrderData
	query := map[string]interface{}{}
	query["ship_number"] = shipNumber
	query["ship_type"] = shipType
	_, err := engine.Engine.Table(entity.OrderData{}).Select("*").Where(query).Get(&data)
	if err != nil {
		log.Error("GetOrderDataByShipTypeOrderId data Error [shipNumber: %s, ship_type: %s]", shipNumber, shipType)
		log.Error("GetOrderDataByShipTypeOrderId Error [%s]", err.Error())
		return data, fmt.Errorf("資料庫異常")
	}
	return data, nil
}

func GetOrderDataByPostBag(engine *database.MysqlSession, shipNumber string) (entity.OrderData, error) {
	var data entity.OrderData
	query := map[string]interface{}{}
	query["ship_number"] = shipNumber
	_, err := engine.Engine.Table(entity.OrderData{}).Select("*").Where(query).Get(&data)
	if err != nil {
		log.Error("GetOrderDataByShipTypeOrderId data Error [shipNumber: %s, ship_type: %s]", shipNumber)
		log.Error("GetOrderDataByShipTypeOrderId Error [%s]", err.Error())
		return data, fmt.Errorf("資料庫異常")
	}
	return data, nil
}

// 物流確認訂單狀態 by 訂單編號
func GetOrderDataByShipTypeOrderId(engine *database.MysqlSession, orderId, shipType string) (entity.OrderData, error) {
	var data entity.OrderData
	query := map[string]interface{}{}
	query["order_id"] = orderId
	query["ship_type"] = shipType
	_, err := engine.Engine.Table(entity.OrderData{}).Select("*").Where(query).Get(&data)
	if err != nil {
		log.Error("GetOrderDataByShipTypeOrderId data Error [orderId: %s, ship_type: %s]", orderId, shipType)
		log.Error("GetOrderDataByShipTypeOrderId Error [%s]", err.Error())
		return data, fmt.Errorf("資料庫異常")
	}

	return data, nil
}

//取出Session訂單內容
func GetSessionOrderByOrderId(engine *database.MysqlSession, OrderId string) (entity.OrderData, error) {
	var OrderData entity.OrderData
	sql := fmt.Sprintf("SELECT * FROM order_data WHERE order_id = ?")
	_, err := engine.Session.SQL(sql, OrderId).Get(&OrderData)
	if err != nil {
		return OrderData, err
	}
	return OrderData, nil
}

func SearchOrderAndDetailData(engine *database.MysqlSession, where []string, bind []interface{}, by string, limit int, start int) ([]entity.OrderData, error) {
	var data []entity.OrderData
	start = (start - 1) * limit
	if err := engine.Engine.Table(entity.OrderData{}).Select("*").
		Join("LEFT", entity.OrderDetail{}, "order_data.order_id = order_detail.order_id").
		Where(strings.Join(where, " AND "), bind...).
		GroupBy("order_data.order_id").Desc(by).Limit(limit, start).Find(&data); err != nil {
		return data, err
	}
	return data, nil
}

func SearchOrderData(engine *database.MysqlSession, where []string, bind []interface{}, limit int, start int) ([]entity.OrderData, error) {
	var OrdersData []entity.OrderData
	start = (start - 1) * limit
	if err := engine.Engine.Table(entity.OrderData{}).Select("*").
		Where(strings.Join(where, " AND "), bind...).Desc("create_time").
		Limit(limit, start).
		Find(&OrdersData); err != nil {
		log.Error("Get Order Database Error", err)
		return nil, err
	}
	return OrdersData, nil
}

//計算搜尋結果數
func CountSearchOrderData(engine *database.MysqlSession, where []string, bind []interface{}) int64 {
	result, err := engine.Engine.Table(entity.OrderData{}).Where(strings.Join(where, " AND "), bind...).Count()
	if err != nil {
		log.Error("count search order Database Error", err)
		return 0
	}
	return result
}

//計算未讀數
func CountUnreadOrderData(engine *database.MysqlSession, sellerId string) int64 {
	result, err := engine.Engine.Table(entity.OrderData{}).Select("count(*)").
		Where("seller_unread = ? AND seller_id = ?", 0, sellerId).Count()
	if err != nil {
		log.Error("count wait Order Database Error", err)
		return 0
	}
	return result
}

/**
 * 計算買家總訂單數
 */
func CountBuyerAllOrderData(engine *database.MysqlSession, UserId string) int64 {
	result, err := engine.Engine.Table(entity.OrderData{}).Select("count(*)").
		Where("buyer_id = ? AND order_status NOT IN (?, ?)", UserId, Enum.OrderFail, Enum.OrderInit).Count()
	if err != nil {
		log.Error("count refund Order Database Error", err)
		return 0
	}
	return result
}

/**
 * 買家計算運送狀態數
 */
func CountBuyerOrderByShipStatus(engine *database.MysqlSession, UserId string) int64 {
	result := GetOrderCount(engine, "", UserId, Enum.OrderShipment, "")
	return result
}

/**
 * 買家計算訂單狀態數
 */
func CountOrderByOrderStatus(engine *database.MysqlSession, UserId string, status string) (int64, error) {
	result := GetOrderCount(engine, "", UserId, status, status)
	return result, nil
}

/**
 * 買家計算訂單狀態數
 */
func CountBuyerOrderWaitData(engine *database.MysqlSession, UserId string, status string) int64 {
	result := GetOrderCount(engine, "", UserId, "ORDER", status)
	return result
}

/**
 * 買家計算退貨退款訂單
 */
func CountBuyerOrderRefundData(engine *database.MysqlSession, UserId string) int64 {
	result := GetOrderCount(engine, "", UserId, "RETURN", Enum.OrderRefund)
	return result
}

/**
 * 買家計算取消訂單
 */
func CountBuyerOrderCancelData(engine *database.MysqlSession, UserId string) int64 {
	result := GetOrderCount(engine, "", UserId, "ORDER", Enum.OrderCancel)
	return result
}

//計算總訂單數
func CountAllOrderData(engine *database.MysqlSession, StoreId string) int64 {
	result, err := engine.Engine.Table(entity.OrderData{}).Select("count(*)").
		Where("store_id = ? AND order_status NOT IN (?, ?)", StoreId, Enum.OrderFail, Enum.OrderInit).Count()
	if err != nil {
		log.Error("count refund Order Database Error", err)
		return 0
	}
	return result
}

/**
 * 計算待付款數
 */
func CountWaitOrderData(engine *database.MysqlSession, StoreId string) int64 {
	result, err := engine.Engine.Table(entity.OrderData{}).Select("count(*)").
		Where("store_id = ? AND order_status = ?", StoreId, Enum.OrderWait).Count()
	if err != nil {
		log.Error("count refund Order Database Error", err)
		return 0
	}
	return result
}

func CountShipOverdueOrderData(engine *database.MysqlSession, StoreId string, status string) int64 {
	result := GetOrderCount(engine, StoreId, "", status, status)
	return result
}

//計算轉帳付款過期
func CountOrderExpireData(engine *database.MysqlSession, StoreId string) int64 {
	result, err := engine.Engine.Table(entity.OrderData{}).Select("count(*)").
		Where("store_id = ? AND order_status = ?", StoreId, Enum.OrderExpire).Count()
	if err != nil {
		log.Error("count refund Order Database Error", err)
		return 0
	}
	return result
}

func CountOrderCancelData(engine *database.MysqlSession, StoreId string) int64 {
	result, err := engine.Engine.Table(entity.OrderData{}).Select("count(*)").
		Where("store_id = ? AND order_status = ?", StoreId, Enum.OrderCancel).Count()
	if err != nil {
		log.Error("count refund Order Database Error", err)
		return 0
	}
	return result
}



//計算待出貨數
func CountShipWaitOrderData(engine *database.MysqlSession, StoreId string) int64 {
	result, err := engine.Engine.Table(entity.OrderData{}).Select("count(*)").
		Where("store_id = ? AND order_status = ? AND (ship_status = ? OR ship_status = ?)", StoreId, Enum.OrderSuccess, Enum.OrderShipInit, Enum.OrderShipTake).Count()
	if err != nil {
		log.Error("count refund Order Database Error", err)
		return 0
	}
	return result
}

/**
 * 計算待出貨數
 */
func CountShipSuccessOrderData(engine *database.MysqlSession, StoreId string, status string) int64 {
	result := GetOrderCount(engine, StoreId, "", "SHIP", status)
	return result
}

/**
 * 計算已出貨數
 */
func CountShippedOrderData(engine *database.MysqlSession, StoreId string, status string) int64 {
	result := GetOrderCount(engine, StoreId, "", Enum.OrderShipment, status)
	return result
}

/**
 * 計算退貨數
 */
func GetOrderCount(engine *database.MysqlSession, StoreId, BuyerId, Type, Status string) int64 {
	var where []string
	var bind []interface{}
	if StoreId != "" {
		where = append(where, "store_id = ?")
		bind = append(bind, StoreId)
	} else {
		where = append(where, "buyer_id = ?")
		bind = append(bind, BuyerId)
	}
	switch strings.ToUpper(Type) {
	case Enum.OrderShipment:
		where = append(where, "order_status = ?")
		where = append(where, "(ship_status = ? OR ship_status = ? OR ship_status = ? OR ship_status = ?)")
		bind = append(bind, Enum.OrderSuccess)
		bind = append(bind, Enum.OrderShipment)
		bind = append(bind, Enum.OrderShipTransit)
		bind = append(bind, Enum.OrderShipShop)
		bind = append(bind, Enum.OrderShipNone)
	case "SHIP":
		where = append(where, "ship_status = ?")
		bind = append(bind, Status)
		where = append(where, "order_status = ?")
		bind = append(bind, Enum.OrderSuccess)
	case "ORDER":
		where = append(where, "order_status = ?")
		bind = append(bind, Status)
	case "RETURN":
		where = append(where, "refund_status = ?")
		bind = append(bind, Status)
	case "SHIPWAIT": //待出貨SHIPWAIT
		where = append(where, "order_status = ?")
		where = append(where, "(ship_status = ? OR ship_status = ?)")
		bind = append(bind, Enum.OrderSuccess)
		bind = append(bind, Enum.OrderShipInit)
		bind = append(bind, Enum.OrderShipTake)
	case Enum.OrderShipOverdue:
		where = append(where, "order_status = ?")
		where = append(where, "ship_status = ?")
		bind = append(bind, Enum.OrderSuccess)
		bind = append(bind, Status)
	case Enum.OrderExpire:
		where = append(where, "order_status = ?")
		bind = append(bind, Enum.OrderExpire)
	case "ALL":
		where = append(where, "order_status NOT IN (?, ?)")
		bind = append(bind, Enum.OrderFail)
		bind = append(bind, Enum.OrderInit)
	}
	result, err := engine.Engine.Table(entity.OrderData{}).Select("count(*)").
		Where(strings.Join(where, " AND "), bind...).Count()
	if err != nil {
		log.Error("count refund Order Database Error", err)
		return 0
	}
	return result
}

/**
確認多張單的配送方式是否一致
*/
func DistinctShipOrders(engine *database.MysqlSession, orders []interface{}) ([]string, error) {
	var d []string
	err := engine.Engine.Table("order_data").Distinct("ship_type").In("order_id", orders...).Find(&d)
	if err != nil {
		log.Error("count refund Order Database Error", err)
		return d, err
	}

	return d, nil
}

func DistinctSellerOrdersOwner(engine *database.MysqlSession, orders []interface{}) ([]string, error) {
	var d []string
	err := engine.Engine.Table("order_data").Distinct("store_id").In("order_id", orders...).Find(&d)
	if err != nil {
		log.Error("count refund Order Database Error", err)
		return d, err
	}

	return d, nil
}

/**
取得訂單貨運編號
*/
func GetShipNumber(engine *database.MysqlSession, orders []interface{}) ([]string, error) {
	var d []string
	err := engine.Engine.Table("order_data").Select("ship_number").In("order_id", orders...).Find(&d)
	if err != nil {
		log.Error("count refund Order Database Error", err)
		return d, err
	}

	return d, nil
}

//GetShipNumberGroupByPayWay

func GetShipNumberGroupByPayWay(orders []interface{}) (map[string][]string, error) {

	var datas []entity.SevenShipMapData

	engine := database.GetMysqlEngine()
	defer engine.Close()
	shipNumbers := make(map[string][]string)

	err := engine.Engine.Table("seven_ship_map_data").In("order_id", orders...).Find(&datas)
	if err != nil {
		log.Error("count refund Order Database Error", err)
		return shipNumbers, err
	}

	for _, data := range datas {
		shipNumber := data.PaymentNo + data.VerifyCode
		if data.PayWay != Enum.CvsPay {
			shipNumbers["NonCvsPay"] = append(shipNumbers["NonCvsPay"], shipNumber)
		} else {
			shipNumbers["CvsPay"] = append(shipNumbers["CvsPay"], shipNumber)
		}
	}
	return shipNumbers, nil
}

//銷售訂單報表 fixme 要加訂單狀態 order_status ＆ ship_status
func GetOrderDataByStoreIdAndDay(engine *database.MysqlSession, storeId, StartDate, EndDate string) ([]entity.OrderData, error) {
	start := fmt.Sprintf("%s %s", StartDate, "00:00:00")
	end := fmt.Sprintf("%s %s", EndDate, "23:59:59")
	var data []entity.OrderData
	sql := fmt.Sprintf("SELECT * FROM order_data WHERE store_id = ? AND (order_status = ? OR order_status = ? OR order_status = ? OR order_status = ?) AND create_time BETWEEN ? AND ?")
	err := engine.Engine.SQL(sql, storeId, Enum.OrderWait, Enum.OrderAudit, Enum.OrderSuccess, Enum.OrderCancel, start, end).Find(&data)
	if err != nil {
		log.Error("count Order Database Error", err)
		return data, err
	}
	return data, nil
}

//統計銷售數量
func CountOrderDataByStoreIdAndDay(engine *database.MysqlSession, storeId, StartDate, EndDate string) (int64, error) {
	start := fmt.Sprintf("%s %s", StartDate, "00:00:00")
	end := fmt.Sprintf("%s %s", EndDate, "23:59:59")
	sql := fmt.Sprintf("SELECT count(*) FROM order_data WHERE store_id = ? AND create_time BETWEEN ? AND ? AND (order_status = ? OR order_status = ? OR order_status = ? OR order_status = ?)")
	result, err := engine.Engine.SQL(sql, storeId, start, end, Enum.OrderWait, Enum.OrderAudit, Enum.OrderSuccess, Enum.OrderCancel).Count()
	if err != nil {
		log.Error("count Order Database Error", err)
		return 0, err
	}
	return result, nil
}

//統計銷售金額
func SumOrderDataByStoreIdAndDay(engine *database.MysqlSession, storeId, StartDate, EndDate string) (int64, error) {
	start := fmt.Sprintf("%s %s", StartDate, "00:00:00")
	end := fmt.Sprintf("%s %s", EndDate, "23:59:59")
	var OrderData entity.OrderData
	sql := fmt.Sprintf("SELECT COALESCE(SUM(total_amount), 0) as total FROM order_data WHERE store_id = ? AND create_time BETWEEN ? AND ? AND (order_status = ? OR order_status = ? OR order_status = ? OR order_status = ?)")
	result, err := engine.Engine.SQL(sql, storeId, start, end, Enum.OrderWait, Enum.OrderAudit, Enum.OrderSuccess, Enum.OrderCancel).Sum(OrderData, "total_amount")
	if err != nil {
		log.Error("count Order Database Error", err)
		return 0, err
	}
	return int64(result), nil
}

//統計銷售已撥款金額
func SumOrderAppropriationByStoreIdAndDay(engine *database.MysqlSession, storeId, StartDate, EndDate string) (int64, error) {
	start := fmt.Sprintf("%s %s", StartDate, "00:00:00")
	end := fmt.Sprintf("%s %s", EndDate, "23:59:59")
	var OrderData entity.OrderData
	sql := fmt.Sprintf("SELECT COALESCE(SUM(total_amount), 0) as total FROM order_data WHERE store_id = ? AND capture_status = ? AND capture_time BETWEEN ? AND ?")
	result, err := engine.Engine.SQL(sql, storeId, Enum.OrderCaptureSuccess, start, end).Sum(OrderData, "total_amount")
	if err != nil {
		log.Error("count Order Database Error", err)
		return 0, err
	}
	return int64(result), nil
}

//統計銷售未撥款金額
func SumOrderRecAppropriationByStoreIdAndDay(engine *database.MysqlSession, storeId, StartDate, EndDate string) (int64, error) {
	start := fmt.Sprintf("%s %s", StartDate, "00:00:00")
	end := fmt.Sprintf("%s %s", EndDate, "23:59:59")
	var OrderData entity.OrderData
	sql := fmt.Sprintf("SELECT COALESCE(SUM(capture_amount), 0) as total FROM order_data WHERE store_id = ? AND create_time BETWEEN ? AND ? AND (order_status = ? OR order_status = ? OR order_status = ? OR order_status = ?)")
	result, err := engine.Engine.SQL(sql, storeId, start, end, Enum.OrderWait, Enum.OrderAudit, Enum.OrderSuccess, Enum.OrderCancel).Sum(OrderData, "total_amount")
	if err != nil {
		log.Error("count Order Database Error", err)
		return 0, err
	}
	return int64(result), nil
}

//取出撥款訂單
func GetAppropriationOrder(engine *database.MysqlSession, date string) ([]entity.OrderData, error) {
	var data []entity.OrderData
	sql := fmt.Sprintf("SELECT * FROM order_data WHERE capture_time <= ? AND order_status = ? AND (capture_status = ? OR capture_status = ? OR capture_status = ?) ORDER BY create_time ASC")
	err := engine.Engine.SQL(sql, date, Enum.OrderSuccess, Enum.OrderCaptureAdvance, Enum.OrderCapturePostpone, Enum.OrderCaptureProgress).Find(&data)
	if err != nil {
		log.Error("Select Appropriation Order Database Error", err)
		return data, err
	}
	return data, nil
}

//訂單逾期未寄
func GetOrderShipExpire(engine *database.MysqlSession, date string) ([]entity.OrderData, error) {
	var data []entity.OrderData
	sql := fmt.Sprintf("SELECT * FROM order_data WHERE ship_status = ? AND ship_expire <= ? ORDER BY create_time ASC")
	err := engine.Engine.SQL(sql, Enum.OrderShipTake, date).Find(&data)
	if err != nil {
		log.Error("Select Appropriation Order Database Error", err)
		return data, err
	}
	return data, nil
}

//賣家取出訂單內容
func GetOrderByOrderStatus(engine *database.MysqlSession, Status string) ([]entity.OrderData, error) {
	var OrderData []entity.OrderData
	if err := engine.Engine.Table(entity.OrderData{}).Where("order_status != ?", Status).Find(&OrderData); err != nil {
		return OrderData, err
	}
	return OrderData, nil
}

func GetOrderById(engine *database.MysqlSession, OrderId []string) ([]entity.OrderData, error) {
	var OrderData []entity.OrderData
	if err := engine.Engine.Table(entity.OrderData{}).In("order_id", OrderId).Find(&OrderData); err != nil {
		return OrderData, err
	}
	return OrderData, nil
}

func GetOrderByNotInvoice(engine *database.MysqlSession) ([]entity.OrderData, error) {
	var data []entity.OrderData
	if err := engine.Engine.Table(entity.OrderData{}).Where("ask_invoice = ?", 1).
		And("invoice_status = ?", Enum.InvoiceOpenStatusNot).Find(&data); err != nil {
		log.Error("Get Order Database Error", err)
		return data, err
	}
	return data, nil
}

//取出訂單內容
func GetOrderAndSellerByOrderId(engine *database.MysqlSession, OrderId string) (entity.ErpSearchOrders, error) {
	var data entity.ErpSearchOrders
	_, err := engine.Engine.Table(entity.OrderData{}).Select("*").
		Join("LEFT", entity.StoreData{}, "order_data.store_id = store_data.store_id").
		Join("LEFT", entity.MemberData{}, "order_data.seller_id = member_data.uid").
		Where("order_data.order_id = ?", OrderId).Get(&data)
	if err != nil {
		return data, err
	}
	return data, nil
}

func GetOrderAndBuyerByOrderId(engine *database.MysqlSession, OrderId string) (entity.ErpSearchBuyerOrders, error) {
	var data entity.ErpSearchBuyerOrders
	_, err := engine.Engine.Table(entity.OrderData{}).Select("*").
		Join("LEFT", entity.MemberData{}, "order_data.buyer_id = member_data.uid").
		Where("order_data.order_id = ?", OrderId).Get(&data)
	if err != nil {
		return data, err
	}
	return data, nil
}

func CountOrderJoinStoreJoinRefund(engine *database.MysqlSession, params Request.SearchOrderRequest) (int64, error) {
	where, bind := ComposeSearchOrdersParams(engine, params.Search)
	count, err := engine.Engine.Table(entity.OrderData{}).Select("count(*)").
		//Join("LEFT", entity.OrderRefundData{}, "order_data.order_id = order_refund_data.order_id").
		Join("LEFT", entity.StoreData{}, "order_data.store_id = store_data.store_id").
		Join("LEFT", entity.MemberData{}, "order_data.seller_id = member_data.uid").
		Where(strings.Join(where, " AND "), bind...).Count()
	if err != nil {
		log.Error("Count Order Database Error", err)
		return count, err
	}
	return count, nil
}

func GetOrderJoinStoreJoinRefund(engine *database.MysqlSession, params Request.SearchOrderRequest) ([]entity.ErpSearchOrders, error) {
	where, bind := ComposeSearchOrdersParams(engine, params.Search)
	var data []entity.ErpSearchOrders
	if err := engine.Engine.Table(entity.OrderData{}).Select("*").
		//Join("LEFT", entity.OrderRefundData{}, "order_data.order_id = order_refund_data.order_id").
		Join("LEFT", entity.StoreData{}, "order_data.store_id = store_data.store_id").
		Join("LEFT", entity.MemberData{}, "order_data.seller_id = member_data.uid").
		Where(strings.Join(where, " AND "), bind...).Desc("order_data.create_time").Find(&data); err != nil {
		log.Error("Get Order Database Error", err)
		return data, err
	}
	return data, nil
}

func ComposeSearchOrdersParams(engine *database.MysqlSession, params Request.OrderRequest) ([]string, []interface{}) {
	var where []string
	var bind []interface{}
	if len(params.OrderId) != 0 {
		where = append(where, "order_data.order_id = ?")
		bind = append(bind, params.OrderId)
	}
	if len(params.Buyer) != 0 {
		data, _ := member.GetMemberDataByPhone(engine, params.Buyer)
		where = append(where, "order_data.buyer_id = ?")
		bind = append(bind, data.Uid)
	}
	if len(params.Seller) != 0 {
		data, _ := member.GetMemberDataByPhone(engine, params.Seller)
		where = append(where, "order_data.seller_id = ?")
		bind = append(bind, data.Uid)
	}
	if len(params.Buyer) != 0 {
		data, _ := member.GetMemberDataByPhone(engine, params.Buyer)
		where = append(where, "order_data.buyer_id = ?")
		bind = append(bind, data.Uid)
	}
	if len(params.SellerId) != 0 {
		data, _ := member.GetMemberDataByTerminalId(engine, params.SellerId)
		where = append(where, "order_data.seller_id = ?")
		bind = append(bind, data.Uid)
	}
	if len(params.OrderDate) != 0 {
		date := strings.Split(params.OrderDate, "-")
		where = append(where, "order_data.create_time BETWEEN ? AND ?")
		bind = append(bind, date[0])
		bind = append(bind, date[1])
	}
	if len(params.PaymentDate) != 0 {
		date := strings.Split(params.PaymentDate, "-")
		where = append(where, "order_data.pay_way_time BETWEEN ? AND ?")
		bind = append(bind, date[0])
		bind = append(bind, date[1])
	}
	if len(params.ShipDate) != 0 {
		date := strings.Split(params.ShipDate, "-")
		where = append(where, "order_data.ship_time BETWEEN ? AND ?")
		bind = append(bind, date[0])
		bind = append(bind, date[1])
	}
	if len(params.Receiver) != 0 {
		where = append(where, "order_data.receiver_name = ?")
		bind = append(bind, params.Receiver)
	}
	if len(params.ReceiverAddr1) != 0 {
		where = append(where, "receiver_name = ?")
		bind = append(bind, params.ReceiverAddr1)
	}
	if len(params.ReceiverAddr2) != 0 {
		where = append(where, "receiver_name = ?")
		bind = append(bind, params.ReceiverAddr2)
	}
	if len(params.PaymentType) != 0 {
		where = append(where, "order_data.pay_way = ?")
		bind = append(bind, params.PaymentType)
	}
	if len(params.ShipType) != 0 {
		where = append(where, "order_data.ship_type = ?")
		bind = append(bind, params.ShipType)
	}
	//Seller        string `json:"Seller"`        //會員帳號
	//Buyer         string `json:"Buyer"`         //訂購帳號

	if len(params.ReturnId) != 0 {
		where = append(where, "refund_type = ?")
		bind = append(bind, Enum.TypeReturn)
		where = append(where, "ship_type = ?")
		bind = append(bind, params.ReturnId)
	}

	log.Debug("bind", bind)
	//
	//
	//
	//SellerId      string `json:"SellerId"`      //會員代碼
	//BuyerName     string `json:"BuyerName"`     //訂購人
	//ReceiverPhone string `json:"ReceiverPhone"` //收件電話
	//OrderAmount   int64  `json:"OrderAmount"`   //訂單金額
	//ShipNumber    string `json:"ShipNumber"`    //出貨單號
	//ShipStatus    string `json:"ShipStatus"`    //退貨狀態
	//StoreName     string `json:"StoreName"`     //賣場名稱
	//OrderIp       string `json:"OrderIp"`       //訂購IP
	//
	//ProductAmount int64  `json:"ProductAmount"` //商品金額
	//ShipMode      string `json:"ShipMode"`      //物流業者
	//RefundId      string `json:"RefundId"`      //退款編號
	//OrderId       string `json:"OrderId"`       //訂單編號
	//OrderStatus   string `json:"OrderStatus"`   //訂單狀態
	//ProductId     string `json:"ProductId"`     //商品編號
	//ShipFee       int64  `json:"ShipFee"`       //運費金額
	//ReceiverAddr2 string `json:"ReceiverAddr2"` //收件地址
	//RefundStatus  string `json:"RefundStatus"`  //退款狀態
	return where, bind
}

//取出信用卡請款完成的訂單
func GetOrderAndCredit(engine *database.MysqlSession, start, end string) ([]entity.OrderGwCreditAuth, error) {
	var data []entity.OrderGwCreditAuth
	if err := engine.Engine.Table(entity.GwCreditAuthData{}).Select("*").
		Join("LEFT", entity.OrderData{}, "gw_credit_auth_data.order_id = order_data.order_id").
		Where("gw_credit_auth_data.capture_status = ? AND gw_credit_auth_data.pay_type = ?", Enum.CreditCaptureSuccess, Enum.OrderTransC2c).
		And("gw_credit_auth_data.capture_time BETWEEN ? AND ?", start, end).Find(&data); err != nil {
		return data, err
	}
	return data, nil
}

//日結表
func GetOrderAndTransfer(engine *database.MysqlSession, start, end string) ([]entity.OrderAndTransfer, error) {
	var data []entity.OrderAndTransfer
	if err := engine.Engine.Table(entity.TransferData{}).Select("*").
		Join("LEFT", entity.OrderData{}, "transfer_data.order_id = order_data.order_id").
		Where("transfer_data.transfer_status = ? AND transfer_data.trans_type = ?", Enum.TransferSuccess, Enum.OrderTransC2c).
		And("transfer_data.recd_date BETWEEN ? AND ?", start, end).Find(&data); err != nil {
		return data, err
	}
	return data, nil
}

func InsertOrderBatchShipData(engine *database.MysqlSession, after, before, storeId string) (entity.BatchShipExcelImport, error) {
	var data entity.BatchShipExcelImport
	data.BatchId = tools.GeneratorBatchShipId()
	data.AfterContent = after
	data.BeforeContent = before
	data.ProcessStatus = Enum.OrderInit
	data.StoreId = storeId
	data.CreateTime = time.Now()
	data.UpdateTime = time.Now()
	if _, err := engine.Session.Table(entity.BatchShipExcelImport{}).Insert(&data); err != nil {
		log.Error("Insert Order Batch Ship Database Error", err)
		return data, err
	}
	return data, nil
}

//更新Order Data
func UpdateOrderBatchShipData(engine *database.MysqlSession, data entity.BatchShipExcelImport) error {
	data.ProcessStatus = Enum.OrderSuccess
	data.UpdateTime = time.Now()
	if _, err := engine.Session.Table(entity.BatchShipExcelImport{}).ID(data.BatchId).AllCols().Update(data); err != nil {
		return err
	}
	return nil
}

func GetOrderBatchShipData(engine *database.MysqlSession, batchId, storeId string) (entity.BatchShipExcelImport, error) {
	var data entity.BatchShipExcelImport
	if _, err := engine.Engine.Table(entity.BatchShipExcelImport{}).Select("*").
		Where("batch_id = ? AND store_id = ? AND process_status = ?", batchId, storeId, Enum.OrderInit).Get(&data); err != nil {
		return data, err
	}
	return data, nil
}

//是否存在訂單
func IsExistOrderDataByBuyerAndSeller(engine *database.MysqlSession, BuyerId, SellerId string) int64 {
	result, err := engine.Engine.Table(entity.OrderData{}).Select("count(*)").
		Where("seller_id = ? AND buyer_id = ?", SellerId, BuyerId).Count()
	if err != nil {
		log.Error("count Order Database Error", err)
		return 0
	}
	return result
}

func CountOrderDataBySeller(engine *database.MysqlSession, SellerId string) int64 {
	result, err := engine.Engine.Table(entity.OrderData{}).Select("count(*)").
		Where("seller_id = ?", SellerId).Count()
	if err != nil {
		log.Error("count Order Database Error", err)
		return 0
	}
	return result
}

