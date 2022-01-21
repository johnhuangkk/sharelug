package model

import (
	"api/services/Enum"
	"api/services/Service/Notification"
	"api/services/VO/Request"
	"api/services/VO/Response"
	"api/services/dao/Orders"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"api/services/util/tools"
	"fmt"
	"time"
)

func QueryReturnList(params *Request.ReturnListQuery) ([]Response.QueryReturnListResponse, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var res []Response.QueryReturnListResponse

	data, err := Orders.GetOrderReturnByOrderIdAndSpecId(engine, params.OrderId, params.ProductSpecId)
	if err != nil {
		log.Error("Database Get Error", err)
		return res, fmt.Errorf("系統錯誤！")
	}

	for _, v := range data {
		var resp Response.QueryReturnListResponse
		resp.ReturnId = v.RefundId
		resp.ReturnStatus = v.Status
		resp.ReturnStatusText = Enum.ReturnStatus[v.Status]
		resp.ReturnTime = v.CreateTime.Format("2006/01/02 15:04:05")
		resp.CompleteTime = ""
		if !v.RefundTime.IsZero() {
			resp.CompleteTime = v.RefundTime.Format("2006/01/02 15:04:05")
		}
		resp.ProductName = v.ProductName
		resp.ThisReturn = v.Qty
		resp.TotalReturn = v.Sum
		res = append(res, resp)
	}
	return res, nil
}

func QueryReturn(params *Request.ReturnQuery) (Response.QueryReturnResponse, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()

	var resp Response.QueryReturnResponse
	orderData, err := Orders.GetOrderByOrderId(engine, params.OrderId)
	if err != nil {
		log.Error("Database Get Error", err)
		return resp, fmt.Errorf("系統錯誤！")
	}

	orderDetail, err := Orders.GetOrderDetailByOrderId(engine, orderData.OrderId)
	if err != nil {
		log.Error("Database Get Error", err)
		return resp, fmt.Errorf("系統錯誤！")
	}
	var ReturnLists []Request.ReturnList
	for _, v := range orderDetail {
		var ReturnList Request.ReturnList
		ReturnList.ProductSpecId = v.ProductSpecId
		ReturnLists = append(ReturnLists, ReturnList)
	}

	data, err := GetOrderReturnLast(engine, orderData.OrderId, ReturnLists)
	if err != nil {
		log.Error("Database Get Error", err)
		return resp, fmt.Errorf("系統錯誤！")
	}

	resp.OrderId = orderData.OrderId
	for _, v := range orderDetail {
		var detail Response.ReturnProductList
		detail.ProductName = v.ProductName
		detail.ProductSpecId = v.ProductSpecId
		detail.ProductSpecName = v.ProductSpecName
		detail.ProductPrice = v.ProductPrice
		detail.ProductQty = v.ProductQuantity
		detail.Refundable = v.ProductQuantity - data[v.ProductSpecId].Sum
		resp.ProductList = append(resp.ProductList, detail)
	}
	return resp, nil
}

/**
 * 退貨API
 */
func HandleReturn(params *Request.ReturnParams) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	if err := engine.Session.Begin(); err != nil {
		log.Error("Database Begin Error", err)
		return fmt.Errorf("系統錯誤！")
	}
	orderData, err := Orders.GetOrderByOrderId(engine, params.OrderId)
	if err != nil {
		log.Error("Database Get Error", err)
		return fmt.Errorf("系統錯誤！")
	}
	//取出商品
	orderDetail, err := getOrderDetailByOrderId(engine, orderData.OrderId)
	if err != nil {
		log.Error("Database Get Error", err)
		return fmt.Errorf("系統錯誤！")
	}
	//取出退貨單
	data, err := GetOrderReturnLast(engine, params.OrderId, params.ReturnList)
	if err != nil {
		log.Error("Database Get Error", err)
		return fmt.Errorf("系統錯誤！")
	}
	//判斷退貨是否超過
	for _, value := range params.ReturnList {
		for _, v := range orderDetail {
			if value.ProductSpecId == v.ProductSpecId {
				if value.Qty > v.ProductQuantity - data[v.ProductSpecId].Sum {
					return fmt.Errorf("退貨超過訂單總數量！")
				}
			}
		}
	}

	for _, value := range params.ReturnList {
		if value.ProductSpecId == orderDetail[value.ProductSpecId].ProductSpecId {
			var Return entity.OrderRefundData
			Return.RefundId = tools.GeneratorOrderReturnId()
			Return.OrderId = params.OrderId
			Return.ProductSpecId = value.ProductSpecId
			Return.ProductName = orderDetail[value.ProductSpecId].ProductName
			Return.Qty = value.Qty
			Return.Sum = data[value.ProductSpecId].Sum + value.Qty
			IsReturn := Enum.RefundStatusWait
			if params.IsReturn == 1 {
				IsReturn = Enum.ReturnStatusSuccess
				Return.RefundTime = time.Now()
			}
			Return.Status = IsReturn
			Return.UpdateTime = time.Now()
			err = Orders.InsertOrderReturnData(engine, Return)
			if err != nil {
				log.Error("Database Insert Return Detail Error", err)
				err = engine.Session.Rollback()
				return fmt.Errorf("系統錯誤！")
			}
			if err := Notification.SendReturnApplyMessage(engine, orderData, Return); err != nil {
				err = engine.Session.Rollback()
				return fmt.Errorf("系統錯誤！")
			}
		}
	}
	err = updateOrderRefundStatus(engine, Enum.OrderRefund, orderData)
	if err != nil {
		err = engine.Session.Rollback()
		log.Error("Database Update Return Error", err)
		return fmt.Errorf("系統錯誤！")
	}

	err = engine.Session.Commit()
	return nil
}

//退貨完成確認
func HandleReturnConfirm(params *Request.ReturnConfirmParams) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	if err := engine.Session.Begin(); err != nil {
		log.Error("Database Begin Error", err)
		return fmt.Errorf("系統錯誤！")
	}
	data, err := Orders.GetOrderReturnByReturnId(engine, params.ReturnId)
	if err != nil {
		log.Error("Get Return Database Error", err)
		return fmt.Errorf("系統錯誤！")
	}
	orderData, err := Orders.GetOrderByOrderId(engine, data.OrderId)
	if err != nil {
		log.Error("Database Get Error", err)
		return fmt.Errorf("系統錯誤！")
	}
	if len(data.RefundId) == 0 {
		return fmt.Errorf("無此退貨單！")
	}
	data.Status = Enum.ReturnStatusSuccess
	_, err = Orders.UpdateReturnData(engine, data)
	if err != nil {
		log.Error("Database Update Return Error", err)
		err = engine.Session.Rollback()
		return fmt.Errorf("系統錯誤！")
	}
	if err := Notification.SendReturnSuccessMessage(engine, orderData, data); err != nil {
		err = engine.Session.Rollback()
		return fmt.Errorf("系統錯誤！")
	}
	err = engine.Session.Commit()
	return nil
}

func getOrderDetailByOrderId(engine *database.MysqlSession, OrderId string) (map[string]entity.OrderDetail, error) {
	data := make(map[string]entity.OrderDetail)
	//取出訂單
	orderDetail, err := Orders.GetOrderDetailByOrderId(engine, OrderId)
	if err != nil {
		log.Error("Database Get Error", err)
		return data, fmt.Errorf("系統錯誤！")
	}
	for _, value := range orderDetail {
		data[value.ProductSpecId] = value
	}
	return data, nil
}

//取出退貨單
func GetOrderReturnLast(engine *database.MysqlSession, OrderId string, ReturnList []Request.ReturnList) (map[string]entity.OrderRefundData, error) {
	detail := make(map[string]entity.OrderRefundData)
	for _, value := range ReturnList {
		data, err := Orders.GetOrderReturnLastByOrderId(engine, OrderId, value.ProductSpecId)
		if err != nil {
			log.Error("Database Get Error", err)
			return detail, fmt.Errorf("系統錯誤！")
		}
		detail[value.ProductSpecId] = data
	}
	return detail, nil
}

