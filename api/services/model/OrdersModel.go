package model

import (
	"api/services/Enum"
	"api/services/Service/Balance"
	"api/services/Service/Excel"
	"api/services/Service/History"
	"api/services/Service/MemberService"
	"api/services/Service/Notification"
	"api/services/Service/OrderService"
	"api/services/Service/StoreService"
	"api/services/Service/TransferService"
	"api/services/Service/UserAddressService"
	"api/services/VO/CartsVo"
	"api/services/VO/ExcelVo"
	"api/services/VO/OrderVo"
	"api/services/VO/Request"
	"api/services/VO/Response"
	"api/services/dao/Orders"
	"api/services/dao/Store"
	"api/services/dao/transfer"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"api/services/util/tools"
	"api/services/util/validate"
	"fmt"
	"strings"
	"time"
)

// 訂單明細
func setOrderDetail(engine *database.MysqlSession, orderData entity.OrderData) (Response.OrderDetailResponse, error) {
	var resp Response.OrderDetailResponse
	data, err := ModifyOrderDetail(engine, orderData)
	if err != nil {
		return resp, err
	}
	//不可合拼計算
	for _, v := range data.Merge {
		resp.MergeFee += int(v.ShipFee)
	}
	resp.Merge = append(resp.Merge, data.Merge...)
	//免運計算
	amount := float64(0)
	qty := int64(0)
	for _, v := range data.Free {
		amount += float64(v.ProductPrice * v.ProductQuantity)
		qty += v.ProductQuantity
	}
	//判斷是否達免運條件
	if orderData.FreeShipKey == Enum.FreeShipAmount && int64(amount) >= orderData.FreeShip {
		resp.Free = append(resp.Free, data.Free...)
	} else if orderData.FreeShipKey == Enum.FreeShipQuantity && qty >= orderData.FreeShip {
		resp.Free = append(resp.Free, data.Free...)
	} else {
		//未達免運丟進可合拼
		data.General = append(data.General, data.Free...)
	}
	//可合拼計算
	for _, v := range data.General {
		if v.ShipFee > int64(resp.GeneralFee) {
			resp.GeneralFee = int(v.ShipFee)
		}
	}
	resp.General = append(resp.General, data.General...)
	return resp, nil
}

func ModifyOrderDetail(engine *database.MysqlSession, orderData entity.OrderData) (Response.OrderDetailResponse, error) {
	var resp Response.OrderDetailResponse
	detail, err := Orders.GetOrderDetailByOrderId(engine, orderData.OrderId)
	if err != nil {
		return resp, fmt.Errorf("系統錯誤")
	}
	//修改為 不可合拼、可合拼、免運 三個部份
	for _, v := range detail {
		data, _ := Orders.GetOrderReturnLastByOrderId(engine, v.OrderId, v.ProductSpecId)
		id := strings.Split(v.ProductSpecId, "-")
		var vo Response.OrderDetails
		vo.ProductId = id[0]
		vo.ProductName = v.ProductName
		vo.ProductSpecId = v.ProductSpecId
		vo.ProductSpecName = v.ProductSpecName
		vo.ProductQuantity = v.ProductQuantity
		vo.ProductPrice = v.ProductPrice
		vo.ShipFee = v.ShipFee
		vo.ShipMerge = v.ShipMerge
		vo.ReturnQty = data.Sum
		if orderData.FreeShipKey != Enum.FreeShipNone && v.IsFreeShip == true && v.ShipMerge == 1 {
			resp.Free = append(resp.Free, vo)
		} else if v.ShipMerge == 0 {
			resp.Merge = append(resp.Merge, vo)
		} else {
			resp.General = append(resp.General, vo)
		}
	}
	return resp, nil
}

// 設定配送資訊
func setShipInfo(engine *database.MysqlSession, orderData entity.OrderData) Response.Shipping {
	var shipping Response.Shipping
	//05/28 修改讓面交也可以看到資訊
	//if orderData.ShipType == Enum.F2F {
	//	shipping.Type = orderData.ShipType
	//	shipping.Text = Enum.Shipping[orderData.ShipType]
	//	return shipping
	//}

	shipping.Type = orderData.ShipType
	shipping.Text = Enum.Shipping[orderData.ShipType]
	shipping.Receiver.ReceiverName = orderData.ReceiverName
	shipping.Receiver.ReceiverPhone = orderData.ReceiverPhone
	shipping.ReceiverMasker.ReceiverName = tools.MaskerName(orderData.ReceiverName)
	shipping.ReceiverMasker.ReceiverPhone = tools.MaskerPhone(orderData.ReceiverPhone)

	switch strings.ToUpper(orderData.ShipType) {
	case Enum.I_POST:
		iPost := GetPostBoxAddressById(engine, orderData.ReceiverAddress)
		shipping.Receiver.ReceiverAddress = iPost.Address
		shipping.ReceiverMasker.ReceiverAddress = iPost.Address
		shipping.ReceiverMasker.ReceiverAlias = iPost.Alias
	case Enum.CVS_FAMILY, Enum.CVS_HI_LIFE, Enum.CVS_OK_MART, Enum.CVS_7_ELEVEN:
		address := UserAddressService.HandleCVSAddress(engine, orderData.ShipType, orderData.ReceiverAddress)
		shipping.Receiver.ReceiverAddress = address.Address
		shipping.ReceiverMasker.ReceiverAddress = address.Address
		shipping.ReceiverMasker.ReceiverAlias = address.Alias
	default:
		log.Debug("address", orderData.ReceiverAddress)
		shipping.Receiver.ReceiverAddress = orderData.ReceiverAddress
		shipping.ReceiverMasker.ReceiverAddress = tools.MaskerAddress(orderData.ReceiverAddress)
	}
	return shipping
}

// 設定付款方式
func setPayWayInfo(engine *database.MysqlSession, orderData entity.OrderData) (interface{}, error) {
	var payWay interface{}

	switch strings.ToUpper(orderData.PayWay) {
	case Enum.Transfer:
		var payment entity.TransferPayWay
		transferData, err := transfer.GetTransferByOrderId(engine, orderData.OrderId)
		if err != nil {
			return payment, err
		}
		payment.Type = orderData.PayWay
		payment.Text = Enum.PayWay[orderData.PayWay]
		payment.BankAccount = transferData.BankAccount
		payment.BankName = transferData.BankCode + " " + transferData.BankName
		payment.BankExpireDate = transferData.ExpireDate.Format("2006/01/02") + " 23:59"
		payWay = payment
	case Enum.Balance:
		var payment entity.BalancePayWay
		payment.Type = orderData.PayWay
		payment.Text = Enum.PayWay[orderData.PayWay]
		payment.OrderAmount = orderData.TotalAmount
		balance := Balance.GetBalanceByUid(engine, orderData.BuyerId)
		payment.Balance = balance
		payWay = payment
	default:
		var payment entity.OtherPayWay
		payment.Type = orderData.PayWay
		payment.Text = Enum.PayWay[orderData.PayWay]
		payWay = payment
	}

	return payWay, nil
}

func CheckOrderPayment(OrderId string) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	orderData, err := Orders.GetOrderByOrderId(engine, OrderId)
	if err != nil || len(orderData.OrderId) == 0 {
		return err
	}
	if orderData.OrderStatus != Enum.OrderSuccess {
		return fmt.Errorf("尚未付款")
	}
	return nil
}

/**
 * 取訂單內容
 */
func GetOrderData(user entity.MemberData, store entity.StoreDataResp, OrderId string) (Response.OrderResponse, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var orderResp Response.OrderResponse
	orderData, err := Orders.GetOrderByOrderId(engine, OrderId)
	if err != nil || len(orderData.OrderId) == 0 {
		log.Error("GetOrderByOrderId Error", orderData)
		return orderResp, fmt.Errorf("系統錯誤")
	}
	if user.Uid != orderData.BuyerId && store.StoreId != orderData.StoreId {
		return orderResp, fmt.Errorf("系統錯誤")
	}
	storeData, err := Store.GetStoreDataByStoreId(engine, orderData.StoreId)
	if err != nil {
		return orderResp, fmt.Errorf("系統錯誤")
	}
	orderResp = orderData.GetOrderResponse(storeData)
	orderResp.OrderStatusText = OrderService.GetOrderStatusText(orderData)
	orderResp.Income = Balance.CalculateIncome(orderData)
	data, err := Orders.GetOrderRefundLastByOrderId(engine, orderData.OrderId)
	if err != nil {
		log.Error("Get Order Refund Last Error", err)
		return orderResp, fmt.Errorf("系統錯誤！")
	}
	orderResp.RefundAmount = int64(data.Total)
	orderResp.FormUrl = orderData.FormUrl
	//取退貨退款
	if err := OrderService.GetReturnAndRefundByOrderId(engine, &orderResp); err != nil {
		return orderResp, fmt.Errorf("系統錯誤")
	}
	// 訂單明細
	orderResp.Detail, err = setOrderDetail(engine, orderData)
	if err != nil {
		return orderResp, err
	}
	// 設定配送資訊
	orderResp.Shipping = setShipInfo(engine, orderData)
	// 設定付款資訊
	orderResp.Payment, err = setPayWayInfo(engine, orderData)
	if err != nil {
		return orderResp, err
	}

	return orderResp, nil
}

/**
 * 處理訂單付款
 */
func HandleCreateOrder(buyerData entity.MemberData, params *Request.PayParams, carts CartsVo.Carts) (Response.PayResponse, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var resp Response.PayResponse
	if err := engine.Session.Begin(); err != nil {
		log.Error("Set Begin Error", err)
		return resp, fmt.Errorf("1005002")
	}
	//取商品資料及檢查庫存數以及購買數不足時回應錯誤
	products, err := OrderService.CheckProductStock(engine, carts)
	if err != nil {
		return resp, err
	}
	storeData, _ := Store.GetStoreDataByStoreId(engine, carts.StoreId)
	if storeData.SellerId == "" {
		log.Error("get storeData Error", storeData)
		return resp, fmt.Errorf("1005002")
	}
	if carts.Shipping == "" {
		return resp, fmt.Errorf("1005011")
	}
	if params.PayWay == "" {
		return resp, fmt.Errorf("1005012")
	}
	//判斷是否為信用卡
	if params.PayWay == Enum.Credit && len(params.CardId) == 0 {
		//驗證信用卡卡號是否正確
		if !validate.IsValidCreditNumber(params.CardNumber) {
			return resp, fmt.Errorf("1005015")
		}
	}
	//撿查是否為自有商品
	if err := OrderService.CheckNonOwnedProduct(engine, buyerData.Uid, carts.StoreId); err != nil {
		return resp, fmt.Errorf("1005001")
	}
	//寫訂單
	orderData, err := createOrder(engine, storeData, buyerData.Uid, carts, params, products)
	if err != nil {
		engine.Session.Rollback()
		log.Error("Create Order Error", err)
		return resp, fmt.Errorf("1005002")
	}
	//變更真實姓名
	if err := ChangeMemberRealName(engine, buyerData.Uid, params.BuyerName); err != nil {
		engine.Session.Rollback()
		return resp, fmt.Errorf("1005002")
	}
	// 加入常用地址
	if len(params.ReceiverId) == 0 {
		//面交不要記錄地址
		if carts.Shipping != Enum.F2F {
			if err := UserAddressService.NewShipReceiverAddress(engine, buyerData, carts.Shipping, params.BuyerName, params.BuyerPhone, params.ReceiverName, params.ReceiverPhone, params.Address);
				err != nil {
				return resp, fmt.Errorf("1005002")
			}
		}
	}
	//產生訂單DETAIL
	if err := OrderService.ProcessOrderDetail(engine, carts, products, orderData); err != nil {
		engine.Session.Rollback()
		return resp, fmt.Errorf("1005002")
	}
	//處理優惠卷記錄
	if err := HandleCouponUsedRecord(engine, orderData); err != nil {
		engine.Session.Rollback()
		return resp, fmt.Errorf("1005002")
	}
	if err := engine.Session.Commit(); err != nil {
		return resp, fmt.Errorf("1005002")
	}
	//執行付款
	res, err := processPayment(engine, buyerData.Uid, params, orderData)
	if err != nil {
		//退還商品庫存
		log.Error("payment Error", err)
		if err := OrderService.PaymentFailReturnStock(engine, orderData.OrderId); err != nil {
			log.Error("Process Return Stock Error", err)
		}
		return resp, err
	}
	//付款失敗將庫存回寫及即時帳單打開
	if res.Status == Enum.OrderFail {
		if err := OrderService.PaymentFailReturnStock(engine, orderData.OrderId); err != nil {
			log.Error("Process Return Stock Error", err)
		}
	}
	resp.OrderId = orderData.OrderId
	log.Info("Order Return Status", res.Status, res.RtnHtml)
	if res.Status == Enum.OrderWait {
		resp.RtnURL = res.RtnHtml
		return resp, nil
	}
	return resp, nil
}

//C2C付款
func processPayment(engine *database.MysqlSession, buyerId string, params *Request.PayParams, orderData entity.OrderData) (Response.PaymentResponse, error) {
	var resp Response.PaymentResponse
	switch strings.ToUpper(params.PayWay) {
	case Enum.Transfer: //轉帳
		//建立轉帳資料
		Ent, err := TransferService.CreateTransfer(engine, orderData.OrderId, int64(orderData.TotalAmount), Enum.OrderTransC2c)
		if err != nil {
			log.Error("Create transfer", err)
			return resp, fmt.Errorf("1005004")
		}
		//變更訂單狀態為等待付款
		if err := OrderService.OrderWaitPayment(engine, orderData, Enum.OrderWait); err != nil {
			log.Error("Order change Status Error", err)
			return resp, fmt.Errorf("1005004")
		}
		//發送轉帳訊息
		if err := Notification.SendAtmSuccessfullyOrderedMessage(engine, orderData.OrderId); err != nil {
			log.Error("Send Atm Successfully Ordered Message Error", err)
			return resp, fmt.Errorf("1005004")
		}
		date := Ent.ExpireDate.Format("2006/01/02") + " 23:59"
		resp.Status = Enum.OrderSuccess
		resp.Message = "你已透過Check'Ne成功訂購商品，請在" + date + "前至ATM完成繳款，請登入網站追蹤訂單。"
		return resp, nil
	case Enum.Balance: //餘額付款
		//扣除買家餘額
		err := Balance.C2cBalanceTransaction(engine, orderData)
		if err != nil {
			log.Error("Balance transfer error", err)
			return resp, fmt.Errorf("1005016")
		}
		resp.Status = Enum.OrderSuccess
		resp.Message = "你已透過Check'Ne成功訂購商品，並已選擇餘額付款，請登入網站追蹤訂單。"
		return resp, nil
	case Enum.Credit: //信用卡
		//取信用卡資料 或新增信用卡資料
		cardData, err := MemberService.GetCreditData(engine, params.CardId, params.CardNumber, params.CardExpiration, buyerId)
		if err != nil {
			log.Error("Get Card data Error", err)
			return resp, fmt.Errorf("1005006")
		}
		vo := OrderVo.CreditPaymentVo{
			OrderId:     orderData.OrderId,
			TotalAmount: orderData.TotalAmount,
			OrderType:   Enum.OrderTransC2c,
			AuditStatus: Enum.OrderAuditRelease,
			BuyerId:     buyerId,
			SellerId:    orderData.SellerId,
		}
		//信用卡取授權
		res, err := Balance.HandleAuth(vo, params, cardData)
		if err != nil {
			log.Error("Auth handle Error", err)
			return resp, fmt.Errorf("1005006")
		}
		resp.Status = res.Status
		resp.Message = "你已透過Check'Ne成功訂購商品，並以信用卡完成付款，請登入網站追蹤訂單。"
		resp.RtnHtml = res.RtnURL
		return resp, nil
	case Enum.CvsPay: //超商貨到付款
		if err := Balance.OrderCheckout(engine, orderData.OrderId, Enum.OrderSuccess); err != nil {
			log.Error("超商取貨付款 Order Checkout Error", err)
			return resp, fmt.Errorf("1005007")
		}
		resp.Status = Enum.OrderSuccess
		resp.Message = "你已透過Check'Ne成功訂購商品，並已選擇超商取貨付款，請登入網站追蹤訂單。"
		return resp, nil
	}
	return resp, fmt.Errorf("1005002")
}

/**
 *建立訂單
 */
func createOrder(engine *database.MysqlSession, storeData entity.StoreData, buyerId string, carts CartsVo.Carts, params *Request.PayParams, products []entity.ProductsData) (entity.OrderData, error) {
	//取訂單編號
	orderId := OrderService.GetNewOrderId(carts.Realtime, params.PayType)
	if len(params.ReceiverId) != 0 {
		params.Address = UserAddressService.GetReceiverAddress(engine, params.ReceiverId)
	}
	ent := carts.GenerateOrder(orderId, buyerId, storeData, params)
	vo := ent.GetOrder()
	fee := Balance.CalculatePlatFormFee(vo)
	ent.PlatformShipFee = fee.PlatformShipFee
	ent.PlatformTransFee = fee.PlatformTransFee
	ent.PlatformPayFee = fee.PlatformPayFee
	ent.PlatformInfoFee = fee.PlatformInfoFee
	ent.CaptureAmount = fee.CaptureAmount
	ent.FormUrl = OrderService.GetDetailProductFormUrl(engine, products)
	// 宅配檢查
	if (strings.ToLower(carts.Shipping) == Enum.DELIVERY_POST || strings.ToLower(carts.Shipping) == Enum.DELIVERY_E_CAN ||
		strings.ToLower(carts.Shipping) == Enum.DELIVERY_T_CAT || strings.ToLower(carts.Shipping) == Enum.DELIVERY_OTHER) &&
		len(strings.Split(params.Address, ",")) != 4 {
		return ent, fmt.Errorf("地址格式錯誤")
	}
	setting, err := StoreService.GetStoreFreeShipping(engine, carts.StoreId)
	if err != nil {
		log.Debug("get free ship setting error", err)
	}
	if carts.Shipping == Enum.SELF_DELIVERY {
		ent.FreeShipKey = setting.SelfDeliveryKey
		ent.FreeShip = setting.SelfDeliveryFree
	} else {
		ent.FreeShipKey = setting.FreeShipKey
		ent.FreeShip = setting.FreeShip
	}
	orderData, err := Orders.InsertOrderData(engine, ent)
	if err != nil {
		log.Debug("Insert Order data Error!!")
		return orderData, err
	}
	//寫入LOG
	err = SetOrderStatusLog(Enum.StatusLogOrderDataOrderStatus, ent.OrderId, "", ent.OrderStatus, buyerId)
	if err != nil {
		log.Error("Update Status History Log Error", err)
		return orderData, err
	}
	return orderData, nil
}

//func OrderCvsPayCheckoutOfSession(engine *database.MysqlSession, orderId string, status string) error {
//	OrderData, err := Orders.GetSessionOrderByOrderId(engine, orderId)
//	if err != nil {
//		log.Error("Get Order Data Error", err)
//		return err
//	}
//	log.Debug("OrderData", OrderData)
//	OldStatus := Enum.OrderInit
//	if OrderData.OrderStatus != "" {
//		OldStatus = OrderData.OrderStatus
//	}
//	OrderData.CsvCheck = 0
//	OrderData.OrderStatus = status
//	err = OrderService.UpdateOrderDataByOrderStatus(engine, OldStatus, OrderData)
//	if err != nil {
//		log.Error("Update Order Data Error", err)
//		return err
//	}
//	return nil
//}

////付款完成
//func OrderCheckoutOfSession(engine *database.MysqlSession, orderId string, status string) error {
//	OrderData, err := Orders.GetSessionOrderByOrderId(engine, orderId)
//	if err != nil {
//		log.Error("Get Order Data Error", err)
//		return err
//	}
//	log.Debug("OrderData", OrderData)
//	OldStatus := Enum.OrderInit
//	if OrderData.OrderStatus != "" {
//		OldStatus = OrderData.OrderStatus
//	}
//	OrderData.OrderStatus = status
//	OrderData.PayWayTime = time.Now()
//	err = OrderService.UpdateOrderDataByOrderStatus(engine, OldStatus, OrderData)
//	if err != nil {
//		log.Error("Update Order Data Error", err)
//		return err
//	}
//	return nil
//}

// 寫入貨運編號
func WriteOrderDataShipNumber(engine *database.MysqlSession, shipNumber string, OrderData entity.OrderData) (entity.OrderData, error) {
	OrderData.ShipNumber = shipNumber
	OrderData.ShipStatus = Enum.OrderShipTake
	_, err := Orders.UpdateOrderData(engine, OrderData.OrderId, OrderData)
	if err != nil {
		return OrderData, err
	}
	return OrderData, nil
}

////更新訂單狀態
//func updateOrderDataByOrderStatus(engine *database.MysqlSession, OldStatus string, OrderData entity.OrderData) error {
//	OrderData.UpdateTime = time.Now()
//	log.Debug("updateOrderDataByOrderStatus", OrderData)
//	_, err := Orders.UpdateOrderData(engine, OrderData.OrderId, OrderData)
//	if err != nil {
//		log.Error("Update Order Data Error", err)
//		return err
//	}
//
//	err = SetOrderStatusLog(Enum.StatusLogOrderDataOrderStatus,  OrderData.OrderId, OldStatus, OrderData.OrderStatus, "")
//	if err != nil {
//		log.Error("Update Status History Log Error", err)
//		return err
//	}
//	return nil
//}

// 記錄狀態變更
func SetOrderStatusLog(Field string, OrderId string, OldStatus string, NewStatus string, UserName string) error {
	err := History.GenerateStatusLog(Enum.StatusLogOrderData, Field, OrderId, OldStatus, NewStatus, UserName)
	if err != nil {
		log.Error("Update Status History Log Error", err)
		return err
	}
	return nil
}

//fixme 搜尋銷售訂單出貨列表
func GetSearchShipData(storeData entity.StoreDataResp, params Request.OrderSearch) (Response.ShipSearchResponse, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var resp Response.ShipSearchResponse
	if storeData.StoreId == "" {
		return resp, fmt.Errorf("系統錯誤")
	}
	var page = tools.CheckIsZero(params.Page, 1)
	where, bind, _ := setOrderParams(storeData.StoreId, "", params)
	resp.OrderCount = Orders.CountSearchOrderData(engine, where, bind)
	var per = tools.CheckIsZero(params.Length, int(resp.OrderCount))
	OrderData, err := Orders.SearchOrderData(engine, where, bind, per, page)
	if err != nil {
		log.Error("get order list", err)
		return resp, err
	}
	OrderList, err := getOrders(engine, OrderData, storeData.StoreId)
	if err != nil {
		log.Error("get order list", err)
		return resp, err
	}

	resp.UnreadCount = Orders.CountUnreadOrderData(engine, storeData.StoreId)
	resp.OrderList = OrderList
	resp.Tabs.ShipWait = Orders.CountShipWaitOrderData(engine, storeData.StoreId)
	resp.Tabs.ShipOverdue = Orders.CountShipOverdueOrderData(engine, storeData.StoreId, Enum.OrderShipOverdue)
	resp.Tabs.Shipment = Orders.CountShippedOrderData(engine, storeData.StoreId, Enum.OrderShipment)
	resp.Tabs.ShipSucc = Orders.CountShipSuccessOrderData(engine, storeData.StoreId, Enum.OrderShipSuccess)
	return resp, nil
}

//買家搜尋銷售訂單列表
func GetSearchBuyerOrderData(userData entity.MemberData, request Request.OrderSearch) (Response.BuyerOrderSearchResponse, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var OrderResp Response.BuyerOrderSearchResponse
	if userData.Uid == "" {
		return OrderResp, fmt.Errorf("系統錯誤")
	}
	var per = tools.CheckIsZero(request.Length, 50)
	var page = tools.CheckIsZero(request.Page, 1)
	where, bind, _ := setOrderParams("", userData.Uid, request)
	OrderData, err := Orders.SearchOrderData(engine, where, bind, per, page)
	if err != nil {
		log.Error("get order list", err)
		return OrderResp, err
	}
	OrderList, err := getOrders(engine, OrderData, userData.Uid)
	if err != nil {
		log.Error("get order list", err)
		return OrderResp, err
	}
	OrderResp.OrderCount = Orders.CountSearchOrderData(engine, where, bind)
	OrderResp.OrderList = OrderList
	OrderResp.Tabs.All = Orders.CountBuyerAllOrderData(engine, userData.Uid)
	OrderResp.Tabs.Shipment = Orders.CountBuyerOrderByShipStatus(engine, userData.Uid)
	OrderResp.Tabs.OrderWait = Orders.CountBuyerOrderWaitData(engine, userData.Uid, Enum.OrderWait)
	OrderResp.Tabs.OrderCancel = Orders.CountBuyerOrderCancelData(engine, userData.Uid)
	OrderResp.Tabs.Refund = Orders.CountBuyerOrderRefundData(engine, userData.Uid)
	return OrderResp, nil
}

//搜尋銷售訂單列表
func GetSearchSellerOrderData(storeData entity.StoreDataResp, params Request.OrderSearch) (Response.OrderSearchResponse, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var OrderResp Response.OrderSearchResponse
	var per = tools.CheckIsZero(params.Length, 50)
	var page = tools.CheckIsZero(params.Page, 1)

	if storeData.StoreId == "" {
		return OrderResp, fmt.Errorf("系統錯誤")
	}
	where, bind, orderBy := setOrderParams(storeData.StoreId, "", params)
	data, err := Orders.SearchOrderAndDetailData(engine, where, bind, orderBy, per, page)
	if err != nil {
		log.Error("get order list", err)
		return OrderResp, err
	}
	OrderList, err := getOrders(engine, data, storeData.StoreId)
	if err != nil {
		log.Error("get order list", err)
		return OrderResp, err
	}
	OrderResp.OrderCount = Orders.CountSearchOrderData(engine, where, bind)
	OrderResp.UnreadCount = Orders.CountUnreadOrderData(engine, storeData.StoreId)
	OrderResp.OrderList = OrderList
	OrderResp.Tabs.All = Orders.CountAllOrderData(engine, storeData.StoreId)
	OrderResp.Tabs.OrderWait = Orders.CountWaitOrderData(engine, storeData.StoreId)
	OrderResp.Tabs.ShipWait = Orders.CountShipWaitOrderData(engine, storeData.StoreId)
	OrderResp.Tabs.Expire = Orders.CountOrderExpireData(engine, storeData.StoreId)
	OrderResp.Tabs.OrderCancel = Orders.CountOrderCancelData(engine, storeData.StoreId)
	return OrderResp, nil
}

func getOrders(engine *database.MysqlSession, orders []entity.OrderData, id string) ([]Response.OrderListResponse, error) {
	var OrderList []Response.OrderListResponse
	for _, v := range orders {
		order, err := composeOrderContents(engine, v, id)
		if err != nil {
			return nil, err
		}
		OrderList = append(OrderList, order)
	}
	return OrderList, nil
}

func composeOrderContents(engine *database.MysqlSession, data entity.OrderData, id string) (Response.OrderListResponse, error) {
	var order Response.OrderListResponse
	order.OrderId = data.OrderId
	order.StoreId = data.StoreId
	order.OrderStatusType = data.OrderStatus
	order.OrderStatusText = OrderService.GetOrderStatusText(data)
	order.RefundStatusType = data.RefundStatus
	order.RefundStatusText = Enum.OrderRefundStatus[data.RefundStatus]
	order.ShipStatusType = data.ShipStatus
	order.ShipStatusText = Enum.OrderShipStatus[data.ShipStatus]
	order.CaptureStatusType = data.CaptureStatus
	order.CaptureStatusText = Enum.OrderCaptureStatus[data.CaptureStatus]
	order.ShipNumber = data.ShipNumber
	order.Ship = Enum.Shipping[data.ShipType]
	order.ShipCompany = data.ShipText
	order.ShipType = data.ShipType
	order.PayWayType = data.PayWay
	order.PayWayText = Enum.PayWay[data.PayWay]
	order.ShipTime = ""
	order.BuyerNotes = data.BuyerNotes
	if !data.ShipTime.IsZero() {
		order.ShipTime = data.ShipTime.Format("2006/01/02 15:04:05")
	}
	order.PayWayTime = ""
	if !data.PayWayTime.IsZero() {
		order.PayWayTime = data.PayWayTime.Format("2006/01/02 15:04:05")
	}
	order.ShipExpire = ""
	if !data.ShipExpire.IsZero() {
		order.ShipExpire = data.ShipExpire.Format("2006/01/02 15:04:05")
	}
	order.CaptureTime = ""
	if !data.CaptureTime.IsZero() {
		order.CaptureTime = data.CaptureTime.Format("2006/01/02 15:04:05")
	}
	order.InvoiceNumber = data.InvoiceNumber
	order.BuyerName = tools.MaskerName(data.BuyerName)
	storeData, err := Store.GetStoreDataByStoreId(engine, data.StoreId)
	if err != nil {
		log.Error("Get Store Data Error", err)
	}
	order.StoreName = storeData.StoreName
	order.CreateTime = data.CreateTime.Format("2006/01/02 15:04:05")
	order.Price = int64(data.TotalAmount)
	order.FormUrl = data.FormUrl
	order.Unread = false
	if id[0:1] != "U" {
		if data.SellerUnread == 1 {
			order.Unread = true
		}
	} else {
		if data.BuyerUnread == 1 {
			order.Unread = true
		}
	}
	details, err := GetOrderDetail(engine, data.OrderId)
	if err != nil {
		log.Error("get order detail list", err)
		return order, err
	}
	order.Detail = details
	return order, err
}

func HandleOrderMemo(params Request.OrderMemoParams) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	orderData, err := Orders.GetOrderByOrderId(engine, params.OrderId)
	if err != nil {
		log.Error("get order Error", err)
		return err
	}
	if len(orderData.OrderId) == 0 {
		return fmt.Errorf("訂單不存在！")
	}
	orderData.OrderMemo = params.OrderMemo
	_, err = Orders.UpdateOrderData(engine, orderData.OrderId, orderData)
	if err != nil {
		log.Error("update order Error", err)
		return err
	}
	return nil
}

//已讀處理
func HandleUnread(UserData entity.MemberData, params Request.OrderReadParams) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	for _, v := range params.OrderId {
		orderData, err := Orders.GetOrderByOrderId(engine, v)
		if err != nil {
			return err
		}

		if orderData.BuyerId == UserData.Uid {
			orderData.BuyerUnread = 1
		} else {
			orderData.SellerUnread = 1
		}
		_, err = Orders.UpdateOrderData(engine, v, orderData)
		if err != nil {
			return err
		}
	}
	return nil
}

//取得訂單內容
func GetOrderDetail(engine *database.MysqlSession, orderId string) ([]Response.OrderDetail, error) {
	var orderDetail []Response.OrderDetail

	detail, err := Orders.GetOrderDetailByOrderId(engine, orderId)
	if err != nil {
		log.Error("get order detail list", err)
		return orderDetail, err
	}

	for _, v := range detail {
		var detail Response.OrderDetail
		detail.ProductName = v.ProductName
		orderDetail = append(orderDetail, detail)
	}
	return orderDetail, nil
}

//組MYSQL
func setOrderParams(StoreId string, UserId string, request Request.OrderSearch) ([]string, []interface{}, string) {
	var sql []string
	var bind []interface{}

	if len(StoreId) != 0 {
		sql = append(sql, "store_id = ?")
		bind = append(bind, StoreId)
	}

	if len(UserId) != 0 {
		sql = append(sql, "buyer_id = ?")
		bind = append(bind, UserId)
	}

	switch request.Tab {
	case "OrderWait": //尚未付款
		sql = append(sql, "order_status = ?")
		bind = append(bind, Enum.OrderWait)
	case "OrderCancel": //訂單取消
		sql = append(sql, "order_status = ?")
		bind = append(bind, Enum.OrderCancel)
	case "ShipWait": //待出貨
		sql = append(sql, "order_status = ?")
		sql = append(sql, "(ship_status = ? OR ship_status = ?)")
		bind = append(bind, Enum.OrderSuccess)
		bind = append(bind, Enum.OrderShipInit)
		bind = append(bind, Enum.OrderShipTake)
	case "Expire":
		sql = append(sql, "order_status = ?")
		bind = append(bind, Enum.OrderExpire)
	case "ShipOverdue": //逾期未寄
		sql = append(sql, "order_status = ?")
		bind = append(bind, Enum.OrderSuccess)
		sql = append(sql, "ship_status = ?")
		bind = append(bind, Enum.OrderShipOverdue)
	case "Shipment": //已出貨
		sql = append(sql, "order_status = ?")
		sql = append(sql, "(ship_status = ? OR ship_status = ? OR ship_status = ? OR ship_status = ?)")
		bind = append(bind, Enum.OrderSuccess)
		bind = append(bind, Enum.OrderShipment)
		bind = append(bind, Enum.OrderShipTransit)
		bind = append(bind, Enum.OrderShipShop)
		bind = append(bind, Enum.OrderShipNone)
	case "ShipSucc": //完成出貨
		sql = append(sql, "order_status = ?")
		sql = append(sql, "ship_status = ?")
		bind = append(bind, Enum.OrderSuccess)
		bind = append(bind, Enum.OrderShipSuccess)
	case "Refund": //退貨/退款
		sql = append(sql, "refund_status = ?")
		bind = append(bind, Enum.OrderRefund)
	case "All":
		sql = append(sql, "order_status NOT IN (?, ?)")
		bind = append(bind, Enum.OrderFail)
		bind = append(bind, Enum.OrderInit)
	}
	if len(request.ShipType) != 0 {
		sql = append(sql, "ship_type = ?")
		bind = append(bind, request.ShipType)
	}
	if len(request.OrderId) != 0 {
		sql = append(sql, "order_id = ?")
		bind = append(bind, request.OrderId)
	}
	if request.Duration != 0 {
		now := time.Now()
		sql = append(sql, "create_time BETWEEN ? AND ?")
		bind = append(bind, now.AddDate(0, 0, - request.Duration))
		bind = append(bind, now)
	}
	by := "order_data.create_time"
	switch request.OrderBy {
	case "name":
		by = "order_detail.product_spec_id"
	case "price":
		by = "order_detail.product_price"
	}
	return sql, bind, by
}

//處理取消訂單
func HandleCancelOrder(StoreData entity.StoreDataResp, params Request.SetPaymentParams) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	if err := engine.Session.Begin(); err != nil {
		log.Error("Set Begin Error", err)
		return fmt.Errorf("1001001")
	}
	orderData, err := Orders.GetOrderByOrderIdAndStoreId(engine, params.OrderId, StoreData.StoreId)
	if err != nil || len(orderData.OrderId) == 0 {
		log.Error("GetOrderByOrderId Error", orderData)
		return fmt.Errorf("1001001")
	}
	//訂單不是SUCCESS或AUDIT不能取消
	if orderData.OrderStatus == Enum.OrderSuccess || orderData.OrderStatus == Enum.OrderAudit {
		//安排出貨前才可以取消 ship_status == OrderShipInit
		if orderData.ShipStatus != Enum.OrderShipInit {
			return fmt.Errorf("1008001")
		}
		//撿查是否已完成付款 //（如果為 CVS_PAY 無需退款）
		if !tools.InArray([]string{Enum.OrderInit, Enum.OrderWait}, orderData.OrderStatus) {
			//撿查付款方式 (不同付款方式 不同退款方式)
			if err := Balance.OrderCancelPaymentRefund(engine, orderData); err != nil {
				engine.Session.Rollback()
				return fmt.Errorf("1001001")
			}
		}
		//餘額扣除平台費用
		if err := Balance.OrderPlatformDeduction(engine, &orderData); err != nil {
			engine.Session.Rollback()
			return fmt.Errorf("1001001")
		}
		//訂單改變裝態
		if err := OrderService.ChangeOrderStatus(engine, orderData, Enum.OrderCancel); err != nil {
			engine.Session.Rollback()
			return fmt.Errorf("1001001")
		}
		if err := OrderService.ProcessReturnStock(engine, orderData.OrderId); err != nil {
			engine.Session.Rollback()
			return fmt.Errorf("1001001")
		}
		if err := Notification.SendOrderCancelMessage(engine, orderData.OrderId); err != nil {
			engine.Session.Rollback()
			return fmt.Errorf("1001001")
		}
		if err := engine.Session.Commit(); err != nil {
			return fmt.Errorf("1005002")
		}
		return nil
	} else {
		if orderData.OrderStatus == Enum.OrderWait {
			return fmt.Errorf("1008002")
		} else {
			return fmt.Errorf("1008001")
		}
	}
}

//建立帳單轉成訂單
func CreateBillOrder(engine *database.MysqlSession, storeData entity.StoreDataResp, bill entity.BillOrderData) (entity.OrderData, error) {
	var Entity entity.OrderData
	Entity.OrderId = bill.BillId
	Entity.SellerId = storeData.SellerId
	Entity.StoreId = storeData.StoreId
	Entity.OrderStatus = Enum.OrderInit
	Entity.RefundStatus = Enum.OrderRefundInit
	Entity.ShipStatus = Enum.OrderShipInit
	Entity.CaptureStatus = Enum.OrderCaptureInit
	Entity.CaptureApply = Enum.OrderCaptureInit
	Entity.PayWay = bill.PayWayType
	//信用卡付款時間

	Entity.ShipType = bill.ShipType
	Entity.SubTotal = bill.SubTotal
	Entity.ShipFee = float64(bill.ShipFee)
	Entity.TotalAmount = bill.TotalAmount
	Entity.BuyerId = bill.BuyerId
	Entity.BuyerName = bill.BuyerName
	Entity.BuyerPhone = bill.BuyerPhone
	Entity.ReceiverName = bill.ReceiverName
	Entity.ReceiverPhone = bill.ReceiverPhone
	Entity.CsvCheck = 1
	if bill.PayWayType == Enum.CvsPay {
		Entity.CsvCheck = 0
	}
	Entity.PlatformShipFee = bill.PlatformShipFee
	Entity.PlatformTransFee = bill.PlatformTransFee
	Entity.PlatformPayFee = bill.PlatformPayFee
	Entity.PlatformInfoFee = bill.PlatformInfoFee
	Entity.CaptureAmount = bill.CaptureAmount
	Entity.ReceiverAddress = bill.ReceiverAddress
	orderData, err := Orders.InsertOrderData(engine, Entity)
	if err != nil {
		log.Debug("Insert Order data Error!!")
		return orderData, err
	}
	if err := OrderService.CreateBillOrderDetail(engine, bill.BillId, storeData.SellerId, bill); err != nil {
		log.Debug("Insert Order Detail data Error!!")
		return orderData, err
	}
	//寫入LOG
	err = SetOrderStatusLog(Enum.StatusLogOrderDataOrderStatus, Entity.OrderId, "", Entity.OrderStatus, storeData.StoreId)
	if err != nil {
		log.Error("Update Status History Log Error", err)
		return orderData, err
	}
	return orderData, nil
}

func GetShipAddress(engine *database.MysqlSession, ShipType, ReceiverAddress string) string {
	switch strings.ToUpper(ShipType) {
	case Enum.I_POST:
		iPost := GetPostBoxAddressById(engine, ReceiverAddress)
		return iPost.Alias
	case Enum.CVS_FAMILY, Enum.CVS_HI_LIFE, Enum.CVS_OK_MART, Enum.CVS_7_ELEVEN:
		address := UserAddressService.HandleCVSAddress(engine, ShipType, ReceiverAddress)
		return address.Alias
	default:
		return ReceiverAddress
	}
}

//匯出批次出貨單
func HandleBatchShipExport(storeData entity.StoreDataResp, params Request.ExportOrderShippingParams) (string, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var data []entity.OrderDetailData
	if params.SelectAll {
		data, _ = Orders.GetOrderAndDetailByStatus(engine, storeData.StoreId, params.OrderBy, Enum.OrderShipInit, Enum.DELIVERY_OTHER)
	} else {
		data, _ = Orders.GetOrderAndDetailByOrderIds(engine, storeData.StoreId, params.OrderBy, params.Orders)
	}
	order, err := modifyOrders(engine, storeData.StoreId, params.OrderBy, data)
	if err != nil {
		log.Error("Get order database Error", err)
		return "", err
	}
	log.Debug("ssss", order)
	var report []ExcelVo.ShipReportVo
	for k, v := range order {
		addr := strings.Split(v.Order.ReceiverAddress, ",")
		res := ExcelVo.ShipReportVo{
			Id:              int64(k + 1),
			OrderId:         v.Order.OrderId,
			ShipIdn:         "",
			ShipNumber:      "",
			ProductName:     v.Detail.ProductName,
			Price:           tools.IntToString(int(v.Detail.ProductPrice)),
			Pieces:          tools.IntToString(int(v.Detail.ProductQuantity)),
			ReceiverName:    v.Order.ReceiverName,
			ReceiverPhone:   v.Order.ReceiverPhone,
			ReceiverCode:    addr[0],
			ReceiverAddress: fmt.Sprintf("%s%s%s", addr[1], addr[2], addr[3]),
			OrderMemo:       v.Order.OrderMemo,
			BuyerNotes:      v.Order.BuyerNotes,
		}
		report = append(report, res)
	}
	filename, err := Excel.ShippingNew().ToShippingReportFile(report)
	if err != nil {
		return "", fmt.Errorf("系統錯誤！")
	}
	return filename, nil
}

func modifyOrders(engine *database.MysqlSession, storeId, orderBy string, data []entity.OrderDetailData) ([]entity.OrderDetailData, error) {
	var order []string
	for _, v := range data {
		if tools.InArray(order, v.Order.OrderId) {
			continue
		}
		order = append(order, v.Order.OrderId)
	}
	var resp []entity.OrderDetailData
	for _, v := range order {
		detail, err := Orders.GetOrderAndDetailByOrderId(engine, storeId, v, orderBy)
		if err != nil {
			log.Error("Get order database Error", err)
			return resp, err
		}
		resp = append(resp, detail...)
	}
	return resp, nil
}

func modifyOrderList(engine *database.MysqlSession, data []entity.OrderDetailData) ([]entity.OrderData, error) {
	var order []string
	for _, v := range data {
		if tools.InArray(order, v.Order.OrderId) {
			continue
		}
		order = append(order, v.Order.OrderId)
	}
	var resp []entity.OrderData
	for _, v := range order {
		data, err := Orders.GetOrderByOrderId(engine, v)
		if err != nil {
			log.Error("Get order database Error", err)
			return resp, err
		}
		resp = append(resp, data)
	}
	return resp, nil
}

//宅配匯出批次出貨單PDF
func HandleBatchDeliverySendExport(storeData entity.StoreDataResp, params Request.ExportOrderShippingParams) ([]Response.ShipReportPdfResponse, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var resp []Response.ShipReportPdfResponse
	var data []entity.OrderData
	if params.SelectAll {
		data, _ = Orders.GetOrderByStatus(engine, storeData.StoreId, Enum.OrderShipInit, Enum.DELIVERY_OTHER)
	} else {
		data, _ = Orders.GetOrderByOrderIds(engine, storeData.StoreId, params.Orders)
	}
	for k, v := range data {
		addr := strings.Split(v.ReceiverAddress, ",")
		res := Response.ShipReportPdfResponse{
			Id:              int64(k + 1),
			OrderId:         v.OrderId,
			ReceiverName:    v.ReceiverName,
			ReceiverPhone:   v.ReceiverPhone,
			ReceiverCode:    addr[0],
			ReceiverCity:    addr[1],
			ReceiverArea:    addr[2],
			ReceiverAddress: addr[3],
		}
		resp = append(resp, res)
	}
	return resp, nil
}

func HandleBatchDeliveryPdfExport(storeData entity.StoreDataResp, params Request.ExportOrderShippingParams) ([]Response.DeliveryReportResponse, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var resp []Response.DeliveryReportResponse
	var data []entity.OrderData
	if params.SelectAll {
		data, _ = Orders.GetOrderByStatus(engine, storeData.StoreId, Enum.OrderShipInit, Enum.DELIVERY_OTHER)
	} else {
		data, _ = Orders.GetOrderByOrderIds(engine, storeData.StoreId, params.Orders)
	}
	for _, v := range data {
		addr := strings.Split(v.ReceiverAddress, ",")
		var res Response.DeliveryReportResponse
		res.OrderId = v.OrderId
		res.OrderDate = v.CreateTime.Format("2006/01/02")
		res.StoreName = storeData.StoreName
		res.ReceiverName = v.ReceiverName
		res.ReceiverPhone = v.BuyerPhone
		res.ReceiverCode = addr[0]
		res.ReceiverCity = addr[1]
		res.ReceiverArea = addr[2]
		res.ReceiverAddress = addr[3]
		res.ShipFee = int64(v.ShipFee)
		res.TotalAmount = int64(v.TotalAmount)
		res.OrderMemo = v.OrderMemo
		res.BuyerNotes = v.BuyerNotes
		detail, _ := Orders.GetOrderDetailByOrderId(engine, v.OrderId)
		for k, row := range detail {
			var name []string
			name = append(name, row.ProductName)
			if len(row.ProductSpecName) != 0 {
				name = append(name, row.ProductSpecName)
			}
			var d Response.DeliveryDetail
			d.Id = int64(k + 1)
			d.ProductName = strings.Join(name, "-")
			d.Quantity = row.ProductQuantity
			d.Price = row.ProductPrice
			res.Details = append(res.Details, d)
		}
		resp = append(resp, res)
	}
	return resp, nil
}

func HandleBatchSelfDeliveryPdfExport(storeData entity.StoreDataResp, params Request.ExportOrderShippingParams) ([]Response.DeliveryReportResponse, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var resp []Response.DeliveryReportResponse
	var data []entity.OrderData
	if params.SelectAll {
		data, _ = Orders.GetOrderByStatus(engine, storeData.StoreId, Enum.OrderShipInit, Enum.SELF_DELIVERY)
	} else {
		data, _ = Orders.GetOrderByOrderIds(engine, storeData.StoreId, params.Orders)
	}
	for _, v := range data {
		addr := strings.Split(v.ReceiverAddress, ",")
		var res Response.DeliveryReportResponse
		res.OrderId = v.OrderId
		res.OrderDate = v.CreateTime.Format("2006/01/02")
		res.StoreName = storeData.StoreName
		res.ReceiverName = v.ReceiverName
		res.ReceiverPhone = v.BuyerPhone
		res.ReceiverCode = addr[0]
		res.ReceiverCity = addr[1]
		res.ReceiverArea = addr[2]
		res.ReceiverAddress = addr[3]
		res.ShipFee = int64(v.ShipFee)
		res.TotalAmount = int64(v.TotalAmount)
		res.OrderMemo = v.OrderMemo
		res.BuyerNotes = v.BuyerNotes
		detail, _ := Orders.GetOrderDetailByOrderId(engine, v.OrderId)
		for k, row := range detail {
			var name []string
			name = append(name, row.ProductName)
			if len(row.ProductSpecName) != 0 {
				name = append(name, row.ProductSpecName)
			}
			var d Response.DeliveryDetail
			d.Id = int64(k + 1)
			d.ProductName = strings.Join(name, "-")
			d.Quantity = row.ProductQuantity
			d.Price = row.ProductPrice
			res.Details = append(res.Details, d)
		}
		resp = append(resp, res)
	}
	return resp, nil
}

//宅配出貨單匯入
func HandleBatchShipFile(storeData entity.StoreDataResp, filename string) (Response.BatchOrderShippingResponse, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var resp Response.BatchOrderShippingResponse
	//讀取EXCEL檔
	excel, err := Excel.ReadOrderShippingExcel(filename)
	if err != nil {
		log.Debug("Read Order Ship Excel Error", err)
		return resp, fmt.Errorf("1001001")
	}
	var orders []Response.BatchOrders
	//取出所有訂單編號
	i := int64(1)
	for _, v := range excel {
		//檢查orderId是否正確
		if tools.IsOrderId(v.OrderId) {
			data, err := Orders.GetOrderByOrderIdAndStoreId(engine, v.OrderId, storeData.StoreId)
			if err != nil {
				log.Error("Get order database Error", err)
				return resp, fmt.Errorf("1001001")
			}
			if len(data.OrderId) == 0 || len(v.ShipNumber) == 0 || len(v.ShipIdn) == 0 || v.ShipNumber == "#N/A" {
				continue
			}
			var rsp Response.BatchOrders
			rsp.Id = i
			rsp.OrderId = data.OrderId
			rsp.Number = v.ShipNumber
			rsp.Trader = v.ShipIdn
			orders = append(orders, rsp)
			i++
		}
	}
	//寫入記錄
	before, _ := tools.JsonEncode(excel)
	after, _ := tools.JsonEncode(orders)
	data, err := Orders.InsertOrderBatchShipData(engine, after, before, storeData.StoreId)
	if err != nil {
		return resp, fmt.Errorf("1001001")
	}
	resp.BatchId = data.BatchId
	resp.BatchOrders = orders
	return resp, nil
}

//宅配出貨單匯入出貨
func HandleBatchShip(storeData entity.StoreDataResp, batchId string) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	err := engine.Session.Begin()
	if err != nil {
		log.Error("Begin Database Error", err)
		return fmt.Errorf("系統錯誤！")
	}
	data, err := Orders.GetOrderBatchShipData(engine, batchId, storeData.StoreId)
	if err != nil {
		engine.Session.Rollback()
		log.Error("Get Order Batch Ship Data Database Error", err)
		return fmt.Errorf("1001001")
	}
	if len(data.BatchId) == 0 {
		return fmt.Errorf("1001002")
	}
	var orders []Response.BatchOrders
	if err := tools.JsonDecode([]byte(data.AfterContent), &orders); err != nil {
		log.Error("json decode Error", err)
		engine.Session.Rollback()
		return fmt.Errorf("1001001")
	}
	for _, v := range orders {
		order, err := Orders.GetOrderByOrderId(engine, v.OrderId)
		if err != nil {
			log.Error("Get order database Error", err)
			engine.Session.Rollback()
			return fmt.Errorf("1001001")
		}
		if len(order.ShipNumber) != 0 {
			continue
		}
		if err := ChangeOrderShipStatus(engine, order, v.Trader, v.Number); err != nil {
			log.Error("Change Order Ship Status Error", err)
			engine.Session.Rollback()
			return fmt.Errorf("1001001")
		}
	}
	//更新批次時間和狀態
	if err := Orders.UpdateOrderBatchShipData(engine, data); err != nil {
		log.Error("Update Order Batch Ship Data Error", err)
		engine.Session.Rollback()
		return fmt.Errorf("1001001")
	}
	if err := engine.Session.Commit(); err != nil {
		log.Error("Commit Database Error", err)
		return fmt.Errorf("1001001")
	}
	return nil
}

//面交匯出批次出貨單
func HandleBatchF2fExport(storeData entity.StoreDataResp, params Request.ExportOrderShippingParams) (string, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var data []entity.OrderDetailData
	if params.SelectAll {
		data, _ = Orders.GetOrderAndDetailByStatus(engine, storeData.StoreId, params.OrderBy, Enum.OrderShipInit, Enum.F2F)
	} else {
		data, _ = Orders.GetOrderAndDetailByOrderIds(engine, storeData.StoreId, params.OrderBy, params.Orders)
	}
	order, err := modifyOrders(engine, storeData.StoreId, params.OrderBy, data)
	if err != nil {
		log.Error("Get order database Error", err)
		return "", err
	}
	var report []ExcelVo.ShipReportVo
	for k, v := range order {
		var name []string
		name = append(name, v.Detail.ProductName)
		if len(v.Detail.ProductSpecName) != 0 {
			name = append(name, v.Detail.ProductSpecName)
		}
		res := ExcelVo.ShipReportVo{
			Id:            int64(k + 1),
			OrderId:       v.Order.OrderId,
			ProductName:   strings.Join(name, "-"),
			Price:         tools.IntToString(int(v.Detail.ProductPrice)),
			Pieces:        tools.IntToString(int(v.Detail.ProductQuantity)),
			ReceiverName:  v.Order.ReceiverName,
			ReceiverPhone: v.Order.ReceiverPhone,
			OrderMemo:     v.Order.OrderMemo,
			BuyerNotes:    v.Order.BuyerNotes,
		}
		report = append(report, res)
	}
	filename, err := Excel.F2fNew().ToShippingReportFile(report)
	if err != nil {
		return "", fmt.Errorf("系統錯誤！")
	}
	return filename, nil
}

//面交匯出批次出貨單PDF
func HandleBatchF2fSndExport(storeData entity.StoreDataResp, params Request.ExportOrderShippingParams) ([]Response.F2fReportPdfResponse, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var resp []Response.F2fReportPdfResponse
	var data []entity.OrderData
	if params.SelectAll {
		data, _ = Orders.GetOrderByStatus(engine, storeData.StoreId, Enum.OrderShipInit, Enum.F2F)
	} else {
		data, _ = Orders.GetOrderByOrderIds(engine, storeData.StoreId, params.Orders)
	}
	for _, v := range data {
		var products []Response.F2fProduct
		detail, _ := Orders.GetOrderDetailByOrderId(engine, v.OrderId)
		for _, row := range detail {
			var name []string
			name = append(name, row.ProductName)
			if len(row.ProductSpecName) != 0 {
				name = append(name, row.ProductSpecName)
			}
			var rep Response.F2fProduct
			rep.ProductName = strings.Join(name, "-")
			rep.ProductQuantity = row.ProductQuantity
			products = append(products, rep)
		}
		res := Response.F2fReportPdfResponse{
			OrderId:       v.OrderId,
			ReceiverName:  v.ReceiverName,
			ReceiverPhone: v.ReceiverPhone,
			Products:      products,
		}
		resp = append(resp, res)
	}
	return resp, nil
}


func HandleBatchSelfDeliveryExport(storeData entity.StoreDataResp, params Request.ExportOrderShippingParams) (string, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var data []entity.OrderDetailData
	if params.SelectAll {
		data, _ = Orders.GetOrderAndDetailByStatus(engine, storeData.StoreId, params.OrderBy, Enum.OrderShipInit, Enum.SELF_DELIVERY)
	} else {
		data, _ = Orders.GetOrderAndDetailByOrderIds(engine, storeData.StoreId, params.OrderBy, params.Orders)
	}
	order, err := modifyOrders(engine, storeData.StoreId, params.OrderBy, data)
	if err != nil {
		log.Error("Get order database Error", err)
		return "", err
	}
	log.Debug("ssss", order)
	var report []ExcelVo.ShipReportVo
	for k, v := range order {
		addr := strings.Split(v.Order.ReceiverAddress, ",")
		res := ExcelVo.ShipReportVo{
			Id:              int64(k + 1),
			OrderId:         v.Order.OrderId,
			ShipIdn:         "",
			ShipNumber:      "",
			ProductName:     v.Detail.ProductName,
			Price:           tools.IntToString(int(v.Detail.ProductPrice)),
			Pieces:          tools.IntToString(int(v.Detail.ProductQuantity)),
			ReceiverName:    v.Order.ReceiverName,
			ReceiverPhone:   v.Order.ReceiverPhone,
			ReceiverCode:    addr[0],
			ReceiverAddress: fmt.Sprintf("%s%s%s", addr[1], addr[2], addr[3]),
			OrderMemo:       v.Order.OrderMemo,
			BuyerNotes:      v.Order.BuyerNotes,
		}
		report = append(report, res)
	}
	filename, err := Excel.ShippingNew().ToShippingReportFile(report)
	if err != nil {
		return "", fmt.Errorf("系統錯誤！")
	}
	return filename, nil
}

//外送匯出批次交寄單
func HandleBatchSelfDeliverySendExport(storeData entity.StoreDataResp, params Request.ExportOrderShippingParams) ([]Response.ShipReportPdfResponse, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var resp []Response.ShipReportPdfResponse
	var data []entity.OrderData
	if params.SelectAll {
		data, _ = Orders.GetOrderByStatus(engine, storeData.StoreId, Enum.OrderShipInit, Enum.SELF_DELIVERY)
	} else {
		data, _ = Orders.GetOrderByOrderIds(engine, storeData.StoreId, params.Orders)
	}
	for k, v := range data {
		addr := strings.Split(v.ReceiverAddress, ",")
		res := Response.ShipReportPdfResponse{
			Id:              int64(k + 1),
			OrderId:         v.OrderId,
			ReceiverName:    v.ReceiverName,
			ReceiverPhone:   v.ReceiverPhone,
			ReceiverCode:    addr[0],
			ReceiverCity:    addr[1],
			ReceiverArea:    addr[2],
			ReceiverAddress: addr[3],
		}
		resp = append(resp, res)
	}
	return resp, nil
}