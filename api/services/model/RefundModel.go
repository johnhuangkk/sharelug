package model

import (
	"api/services/Enum"
	"api/services/Service/Balance"
	"api/services/Service/Notification"
	"api/services/VO/Request"
	"api/services/VO/Response"
	"api/services/dao/Orders"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"api/services/util/tools"
	"fmt"
	"github.com/spf13/viper"
	"time"
)

//搜尋退貨退款列表
func SearchRefundList(params Request.OrderRefundSearch) (Response.SearchRefundResponse, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var resp Response.SearchRefundResponse

	var Limit = tools.CheckIsZero(params.Limit, 50)
	var Start = tools.CheckIsZero(params.Start, 0)

	where, bind := setRefundParams(params)
	data, err := Orders.GetRefundList(engine, where, bind, Limit, Start)
	if err != nil {
		log.Error("Database Get Error", err)
		return resp, fmt.Errorf("系統錯誤！")
	}

	for _, v := range data {
		var res Response.SearchRefundList
		res.OrderId = v.Order.OrderId
		res.RefundId = v.Refund.RefundId
		res.RefundType = v.Refund.RefundType
		if v.Refund.RefundType == Enum.TypeRefund {
			res.RefundTime = v.Refund.CreateTime.Format("2006/01/02 15:04")
			res.ThisRefundAmt = int64(v.Refund.Amount)
			res.TotalRefundAmt = int64(v.Refund.Total)
			res.RefundStatus = v.Refund.Status
			res.RefundStatusText = Enum.RefundStatus[v.Refund.Status]
		} else {
			res.ReturnId = v.Refund.RefundId
			res.ReturnTime = v.Refund.CreateTime.Format("2006/01/02 15:04")
			res.ReturnCheckTime = ""
			if !v.Refund.RefundTime.IsZero() {
				res.ReturnCheckTime = v.Refund.RefundTime.Format("2006/01/02 15:04")
			}
			res.ReturnProductName = v.Refund.ProductName
			res.ThisReturnQty = v.Refund.Qty
			res.TotalReturnQty = v.Refund.Sum
			res.ReturnStatus = v.Refund.Status
			res.ReturnStatusText = Enum.ReturnStatus[v.Refund.Status]
		}
		resp.SearchRefundList = append(resp.SearchRefundList, res)
	}
	var ref Response.RefundTabsResponse
	ref.RefundAll = Orders.GetRefundCount(engine, params.StoreId, "", "")
	ref.RefundSuc = Orders.GetRefundCount(engine, params.StoreId, Enum.TypeRefund, Enum.RefundStatusSuccess)
	ref.RefundWait = Orders.GetRefundCount(engine, params.StoreId, Enum.TypeRefund, Enum.RefundStatusWait)
	ref.ReturnSuc = Orders.GetRefundCount(engine, params.StoreId, Enum.TypeReturn, Enum.RefundStatusSuccess)
	ref.ReturnWait = Orders.GetRefundCount(engine, params.StoreId, Enum.TypeReturn, Enum.RefundStatusWait)
	resp.Tabs = ref

	return resp, nil
}

// set Params
func setRefundParams(params Request.OrderRefundSearch) ([]string, []interface{}) {
	var sql 	[]string
	var bind 	[]interface{}

	if len(params.StoreId) != 0 {
		sql = append(sql, "store_id = ?")
		bind = append(bind, params.StoreId)
	}

	switch params.Tab {
		case "RefundWait":
			sql = append(sql, "refund_type = ?")
			bind = append(bind, Enum.TypeRefund)
			sql = append(sql, "status = ?")
			bind = append(bind, Enum.RefundStatusWait)
		case "RefundSuc":
			sql = append(sql, "refund_type = ?")
			bind = append(bind, Enum.TypeRefund)
			sql = append(sql, "status = ?")
			bind = append(bind, Enum.RefundStatusSuccess)
		case "ReturnWait":
			sql = append(sql, "refund_type = ?")
			bind = append(bind, Enum.TypeReturn)
			sql = append(sql, "status = ?")
			bind = append(bind, Enum.RefundStatusWait)
		case "ReturnSuc":
			sql = append(sql, "refund_type = ?")
			bind = append(bind, Enum.TypeReturn)
			sql = append(sql, "status = ?")
			bind = append(bind, Enum.RefundStatusSuccess)
	}
	return sql, bind
}

//取退款列表
func QueryRefundList(params Request.RefundQuery) ([]Response.QueryRefundListResponse, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var resps []Response.QueryRefundListResponse

	data, err := Orders.GetOrderRefundsByOrderId(engine, params.OrderId)
	if err != nil {
		log.Error("Database Get Error", err)
		return resps, fmt.Errorf("系統錯誤！")
	}

	for _, v := range data {
		var resp Response.QueryRefundListResponse
		resp.RefundId = v.RefundId
		resp.RefundStatus = v.Status
		resp.RefundStatusText = Enum.RefundStatus[v.Status]
		resp.ApplyTime = v.CreateTime.Format("2006/01/02 15:04:05")
		resp.CompleteTime = v.RefundTime.Format("2006/01/02 15:04:05")
		resp.ThisRefund = int64(v.Amount)
		resp.TotalRefund = int64(v.Total)
		resps = append(resps, resp)
	}
	return resps, nil
}

/**
 * 取退款資料
 */
func QueryRefund(params Request.RefundQuery) (Response.QueryRefundResponse, error) {

	engine := database.GetMysqlEngine()
	defer engine.Close()

	var resp Response.QueryRefundResponse

	orderData, err := Orders.GetOrderByOrderId(engine, params.OrderId)
	if err != nil {
		log.Error("Database Get Error", err)
		return resp, fmt.Errorf("系統錯誤！")
	}

	data, err := Orders.GetOrderRefundLastByOrderId(engine, params.OrderId)
	if err != nil {
		log.Error("Database Get Error", err)
		return resp, fmt.Errorf("系統錯誤！")
	}
	resp.OrderId = orderData.OrderId
	resp.OrderAmount = int64(orderData.TotalAmount)
	resp.ProductAmt = int64(orderData.SubTotal)
	resp.ShipFee = int64(orderData.ShipFee)
	resp.RefundedAmt = int64(data.Total)
	return resp, nil
}

/**
 * 退款API
 */
func HandleRefund(storeData entity.StoreDataResp, params *Request.RefundParams) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	if err := engine.Session.Begin(); err != nil {
		log.Error("Database Begin Error", err)
		return fmt.Errorf("系統錯誤！")
	}
	if storeData.SellerId == "" {
		return fmt.Errorf("系統錯誤！")
	}
	//判斷餘額是否足夠 主帳號的UID
	balance := Balance.GetBalanceByUid(engine, storeData.SellerId)
	//log.Debug("aaa", balance)
	if float64(params.Amount) > balance {
		return fmt.Errorf("帳戶餘額不足無法辦理退款。")
	}
	refundData, err := Orders.GetOrderRefundLastByOrderId(engine, params.OrderId)
	if err != nil {
		log.Error("Database Get Error", err)
		return fmt.Errorf("系統錯誤！")
	}
	//判斷還可退的金額是否足夠
	orderData, err := Orders.GetOrderByOrderId(engine, params.OrderId)
	if err != nil {
		log.Error("Database Get Error", err)
		return fmt.Errorf("系統錯誤！")
	}
	if float64(params.Amount) > orderData.TotalAmount - refundData.Amount {
		return fmt.Errorf("退款超過訂單總金額！")
	}
	//產生退款單
	data, err := GeneratorReturnData(engine, params, orderData, refundData)
	if err != nil {
		log.Error("Database Insert Error", err)
		err = engine.Session.Rollback()
		return fmt.Errorf("系統錯誤！")
	}
	comment := fmt.Sprintf("%s<br>%s", orderData.OrderId, storeData.StoreName) //收銀機名稱前10字＋訂單編號
	//賣家扣款
	if err := Balance.Withdrawal(storeData.SellerId, data.RefundId, data.Amount, Enum.BalanceTypeRefund, comment); err != nil {
		err = engine.Session.Rollback()
		log.Error("Database balance Error", err)
		return fmt.Errorf("系統錯誤！")
	}
	//更改訂單狀態
	if err := updateOrderRefundStatus(engine, Enum.OrderRefund, orderData); err != nil {
		err = engine.Session.Rollback()
		log.Error("Database Update Return Error", err)
		return fmt.Errorf("系統錯誤！")
	}
	//買家退款
	if err := buyerRefundPayment(engine, orderData.OrderId, orderData.BuyerId, orderData.PayWay, data.RefundId, comment, data.Amount); err != nil {
		err = engine.Session.Rollback()
		log.Error("BuyerId Refund Payment Error", err)
		return fmt.Errorf("系統錯誤！")
	}
	//發送訊息
	if err := Notification.SendRefundApplyMessage(engine, orderData, data); err != nil {
		return fmt.Errorf("系統錯誤！")
	}
	if err = engine.Session.Commit(); err != nil {
		err = engine.Session.Rollback()
		log.Error("Database Commit Error", err)
		return fmt.Errorf("系統錯誤！")
	}
	return nil
}

//買家退款處理
func buyerRefundPayment(engine *database.MysqlSession, OrderId, BuyerId, PayWay, RefundId, comment string, Amount float64) error {
	//分不同退款回方式 信用卡退回信用卡 其他都退回餘額
	switch PayWay {
		case Enum.Credit:
			//退款金額寫入金流 退款人 收款人 退款編號 付款方式 退款金額 狀態
			if err := Balance.RefundPaymentProcess(engine, OrderId, Amount); err != nil {
				log.Error("Refund Deposit Payment Error", err)
				return err
			}
			userId := viper.GetString("PLATFORM.USERID")
			if err := Balance.Deposit(userId, RefundId, Amount, Enum.BalanceTypeCreditWait, comment); err != nil {
				log.Error("Balance Deposit Error", err)
				return err
			}
		case Enum.Transfer, Enum.Balance, Enum.CvsPay:
			if err := Balance.Deposit(BuyerId, RefundId, Amount, Enum.BalanceTypeRefund, comment); err != nil {
				log.Error("Database balance Error", err)
				return err
			}
	}
	return nil
}

//更新訂單狀態
func updateOrderRefundStatus(engine *database.MysqlSession, newStatus string, OrderData entity.OrderData) error {
	OldStatus := OrderData.ShipStatus
	OrderData.UpdateTime = time.Now()
	OrderData.RefundStatus = newStatus
	if _, err := Orders.UpdateOrderData(engine, OrderData.OrderId, OrderData); err != nil {
		log.Error("Update Order Data Error", err)
		return err
	}
	if err := SetOrderStatusLog(Enum.StatusLogOrderDataOrderStatus,  OrderData.OrderId, OldStatus, OrderData.RefundStatus, ""); err != nil {
		log.Error("Update Status History Log Error", err)
		return err
	}
	return nil
}

//產生退款單
func GeneratorReturnData(engine *database.MysqlSession, params *Request.RefundParams, order entity.OrderData, refund entity.OrderRefundData) (entity.OrderRefundData, error) {
	var data entity.OrderRefundData
	data.RefundId = tools.GeneratorOrderRefundId()
	data.OrderId = params.OrderId
	data.Amount = float64(params.Amount)
	data.Total = refund.Amount + float64(params.Amount)
	data.Status = Enum.RefundStatusSuccess
	if order.PayWay == Enum.Credit {
		data.Status = Enum.RefundStatusAudit
	}
	data.RefundTime = time.Now()
	data.UpdateTime = time.Now()
	if err := Orders.InsertOrderRefundData(engine, data); err != nil {
		log.Error("Database Insert Error", err)
		return data, err
	}
	return data, nil
}

