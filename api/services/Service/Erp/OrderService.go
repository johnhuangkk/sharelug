package Erp

import (
	"api/services/Enum"
	"api/services/VO/Response"
	"api/services/dao/Balance"
	"api/services/dao/Orders"
	"api/services/dao/Store"
	"api/services/dao/Withdraw"
	"api/services/database"
	"api/services/util/log"
	"api/services/util/tools"
)

// func GetOrderByMultipleField() {
// 	orders.GetOrderByMultipleField()
// 	engine * database.MysqlSession
// }

//訂單明細
func GetErpOrderData(engine *database.MysqlSession, OrderId string) (Response.ErpOrderResponse, error) {
	var resp Response.ErpOrderResponse
	data, err := Orders.GetOrderByOrderId(engine, OrderId)
	if err != nil {
		return resp, err
	}
	resp.OrderId = data.OrderId
	resp.OrderDate = data.CreateTime.Format("2006/01/02 15:04")
	resp.ProductTotal = int64(data.SubTotal)
	resp.OrderTotal = int64(data.TotalAmount)
	resp.ShippingFee = int64(data.ShipFee)
	resp.PlatformShipFee = int64(data.PlatformShipFee)
	resp.PlatformInfoFee = int64(data.PlatformInfoFee)
	resp.PlatformPayFee = int64(data.PlatformPayFee)
	resp.PlatformFee = int64(data.PlatformTransFee + data.PlatformShipFee + data.PlatformInfoFee + data.PlatformPayFee)
	storeData, err := Store.GetStoreDataBySellerIdAndStoreId(engine, data.SellerId, data.StoreId)
	if err != nil {
		return resp, err
	}
	resp.StoreName = storeData.StoreName
	resp.SellerName = storeData.Username
	resp.SellerAcct = tools.MaskerPhoneLater(storeData.Mphone)

	detail, err := Orders.GetOrderDetailByOrderId(engine, data.OrderId)
	if err != nil {
		return resp, err
	}
	for _, v := range detail {
		var res Response.ErpOrderDetail
		res.ProductName = v.ProductName
		res.ProductSpec = v.ProductSpecName
		res.ProductQty = v.ProductQuantity
		res.ProductAmt = v.ProductPrice
		res.ShippingFee = v.ShipFee
		res.ShipMerge = false
		if v.ShipMerge == 1 {
			res.ShipMerge = true
		}
		resp.OrderDetail = append(resp.OrderDetail, res)
	}
	return resp, nil
}

//本次退款內容
func GetErpOrderRefundData(engine *database.MysqlSession, RefundId string) (Response.ErpOrderRefundResponse, error) {
	var resp Response.ErpOrderRefundResponse
	//取得退款資料
	data, err := Orders.GetOrderRefundByRefundId(engine, RefundId)
	if err != nil {
		return resp, err
	}
	resp.ApplyDate = data.CreateTime.Format("2006/01/02 15:04")
	resp.RefundId = data.RefundId
	resp.RefundStatus = Enum.RefundStatus[data.Status]
	resp.RefundAmount = int64(data.Amount)
	//取出訂單資料
	orderData, err := Orders.GetOrderByOrderId(engine, data.OrderId)
	if err != nil {
		return resp, err
	}
	resp.OrderAmount = int64(orderData.TotalAmount)
	resp.Reason = Enum.OrderCaptureStatus[orderData.CaptureStatus]
	//取賣家資料
	storeData, err := Store.GetStoreDataBySellerIdAndStoreId(engine, orderData.SellerId, orderData.StoreId)
	if err != nil {
		return resp, err
	}
	resp.SellerAccount = tools.MaskerPhoneLater(storeData.Mphone)
	//取得賣家餘額
	balance, err := Balance.GetBalanceAccountLastByUserId(engine, orderData.SellerId)
	if err != nil {
		return resp, err
	}
	resp.SellerBalance = int64(balance.Balance)
	return resp, nil
}

//取出訂單的退款記錄
func GetErpRefundDataList(engine *database.MysqlSession, orderId string) ([]Response.ErpOrderRefundListResponse, error) {
	var resp []Response.ErpOrderRefundListResponse
	data, err := Orders.GetOrderRefundAllByOrderId(engine, orderId)
	if err != nil {
		return resp, err
	}
	for _, v := range data {
		res := Response.ErpOrderRefundListResponse{
			RefundId:     v.RefundId,
			RefundTime:   v.RefundTime.Format("2006/01/02 15:04"),
			RefundReason: "",
			RefundAmount: int64(v.Amount),
			RefundStatus: Enum.RefundStatus[v.Status],
		}
		resp = append(resp, res)
	}
	return resp, nil
}

//取出訂單的退貨記錄
func GetErpReturnDataList(engine *database.MysqlSession, orderId string) ([]Response.ErpOrderReturnListResponse, error) {
	var resp []Response.ErpOrderReturnListResponse
	data, err := Orders.GetOrderReturnAllByOrderId(engine, orderId)
	if err != nil {
		return resp, err
	}
	for _, v := range data {
		res := Response.ErpOrderReturnListResponse{
			ReturnId:      v.RefundId,
			ReturnTime:    v.RefundTime.Format("2006/01/02 15:04"),
			ReturnStatus:  Enum.RefundStatus[v.Status],
			ReturnProduct: v.ProductName,
			ReturnSpec:    "",
			ReturnQty:     v.Qty,
			ReturnPrice:   0,
		}
		resp = append(resp, res)
	}
	return resp, nil
}

func GetShipTrader(shipType, ShipText string) string {
	switch shipType {
		case Enum.DELIVERY_POST_BAG1, Enum.DELIVERY_POST_BAG2, 
			 Enum.DELIVERY_POST_BAG3, Enum.DELIVERY_I_POST_BAG1, 
			 Enum.DELIVERY_I_POST_BAG2, Enum.DELIVERY_I_POST_BAG3, Enum.I_POST:
			return "中華郵政"
		case Enum.CVS_7_ELEVEN:
			return "統一超商"
		case Enum.CVS_FAMILY:
			return "全家超商"
		case Enum.CVS_OK_MART:
			return "OK超商"
		case Enum.CVS_HI_LIFE:
			return "萊爾富超商"
		case Enum.DELIVERY_T_CAT:
			return "黑貓宅急便"
		case Enum.DELIVERY_E_CAN:
			return "宅配通"
		case Enum.DELIVERY_OTHER:
			return ShipText
	}
	return ""
}

func AccountGetBankName(engine *database.MysqlSession, account string) string {
	if len(account) == 0 {
		return ""
	}
	code := account[0:3]
	log.Debug("bank name", code)
	data, err := Withdraw.GetBankInfoByBankCode(engine, code)
	if err != nil {
		log.Error("Get Bank Info Error")
	}
	return data.BankName
}

