package model

import (
	"api/services/Enum"
	"api/services/Service/Balance"
	"api/services/Service/OrderService"
	"api/services/VO/Request"
	"api/services/dao/Orders"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"api/services/util/tools"
	"fmt"
	"time"
)
//其他運送方式
func HandleSetShipNumber(storeData entity.StoreDataResp, params Request.SetShipNumberParams) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	err := engine.Session.Begin()
	if err != nil {
		log.Error("Begin Database Error", err)
		return fmt.Errorf("系統錯誤！")
	}
	for _, value := range params.ShipNumberList {
		data, err := Orders.GetOrderByOrderId(engine, value.OrderId)
		if err != nil {
			log.Error("Get Order Database Error", err)
			_ = engine.Session.Rollback()
			return fmt.Errorf("系統錯誤！")
		}
		if data.StoreId != storeData.StoreId {
			return fmt.Errorf("訂單賣家錯誤！")
		}
		if data.OrderStatus != Enum.OrderSuccess {
			_ = engine.Session.Rollback()
			return fmt.Errorf("訂單尚未購買完成！")
		}
		if data.ShipStatus != Enum.OrderShipInit {
			_ = engine.Session.Rollback()
			return fmt.Errorf("訂單已非未出貨狀態！")
		}
		//寫入運送單號  運送方式  物流名稱
		if err := ChangeOrderShipStatus(engine, data, value.ShipText, value.Number); err != nil {
			_ = engine.Session.Rollback()
			return fmt.Errorf("系統錯誤！")
		}
	}
	err = engine.Session.Commit()
	if err != nil {
		log.Error("Commit Database Error", err)
		return fmt.Errorf("系統錯誤！")
	}
	return nil
}	

//提前撥付
func HandleAdvancePayment(userData entity.MemberData, params Request.SetPaymentParams) error {
	//todo 出貨狀態 只要以出貨就可提前撥付(排除面交自取及 CVS_PAY 貨到付款)
	engine := database.GetMysqlEngine()
	defer engine.Close()
	err := engine.Session.Begin()
	if err != nil {
		log.Error("Begin Database Error", err)
		return fmt.Errorf("系統錯誤！")
	}
	data, err := Orders.GetOrderByOrderIdAndBuyerId(engine, params.OrderId, userData.Uid)
	if err != nil {
		log.Error("Get Order Database Error", err)
		return fmt.Errorf("系統錯誤！")
	}
	if data.BuyerId != userData.Uid {
		log.Error("not User Error", err)
		return fmt.Errorf("系統錯誤！")
	}
	if data.OrderStatus != Enum.OrderSuccess {
		return fmt.Errorf("訂單尚未購買完成！")
	}
	if data.CaptureStatus == Enum.OrderCaptureSuccess {
		return fmt.Errorf("訂單已撥款！")
	}
	if tools.InArray([]string{Enum.F2F ,Enum.NONE}, data.PayWay) || data.ShipType == Enum.CvsPay {
		return fmt.Errorf("此訂單不得提前付款！")
	}
	if !tools.InArray([]string{Enum.OrderShipInit, Enum.OrderShipTake}, data.ShipStatus) {
		OrderService.OrderCaptureRelease(&data, time.Now())
		log.Debug("Order", data)
		if _, err := Orders.UpdateOrderData(engine, data.OrderId, data); err != nil {
			_ = engine.Session.Rollback()
			log.Error("Update Order Database Error", err)
			return fmt.Errorf("系統錯誤！")
		}
		if err := engine.Session.Commit(); err != nil {
			_ = engine.Session.Rollback()
			log.Error("Commit Database Error", err)
			return fmt.Errorf("系統錯誤！")
		}
	} else {
		return fmt.Errorf("訂單已非出貨狀態！")
	}
	return nil
}

//延長撥付
func HandleExtensionPayment(userData entity.MemberData, params Request.SetPaymentParams) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	if err := engine.Session.Begin(); err != nil {
		log.Error("Begin Database Error", err)
		return fmt.Errorf("系統錯誤！")
	}
	data, err := Orders.GetOrderByOrderIdAndBuyerId(engine, params.OrderId, userData.Uid)
	if err != nil {
		log.Error("Get Order Database Error", err)
		return fmt.Errorf("系統錯誤！")
	}
	if data.OrderStatus != Enum.OrderSuccess {
		return fmt.Errorf("訂單尚未購買完成！")
	}
	if data.CaptureStatus == Enum.OrderCaptureSuccess {
		return fmt.Errorf("訂單已撥款！")
	}
	OrderService.OrderCaptureRelease(&data, time.Time{})
	_, err = Orders.UpdateOrderData(engine, data.OrderId, data)
	if err != nil {
		_ = engine.Session.Rollback()
		log.Error("Update Order Database Error", err)
		return fmt.Errorf("系統錯誤！")
	}
	if err := engine.Session.Commit(); err != nil {
		_ = engine.Session.Rollback()
		log.Error("Commit Database Error", err)
		return fmt.Errorf("系統錯誤！")
	}
	return nil
}

// fixme 面交及無需配送訂單完成交易 (信用卡部份的問題)
func HandleCompleteTransaction(storeData entity.StoreDataResp, params Request.SetPaymentParams) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	err := engine.Session.Begin()
	if err != nil {
		log.Error("Begin Database Error", err)
		return fmt.Errorf("系統錯誤！")
	}
	data, err := Orders.GetOrderByOrderIdAndStoreId(engine, params.OrderId, storeData.StoreId)
	if err != nil {
		log.Error("Get Order Database Error", data)
		return fmt.Errorf("系統錯誤！")
	}
	if data.StoreId != storeData.StoreId {
		return fmt.Errorf("訂單有誤！")
	}
	if !tools.InArray([]string {Enum.F2F, Enum.NONE}, data.ShipType) {
		return fmt.Errorf("訂單有誤！")
	}
	if data.OrderStatus != Enum.OrderSuccess {
		return fmt.Errorf("訂單有誤！")
	}
	data.ShipStatus = Enum.OrderShipNone
	data.ShipTime = time.Now()
	OrderService.OrderCaptureRelease(&data, time.Time{})
	if _, err := Orders.UpdateOrderData(engine, data.OrderId, data); err != nil {
		_ = engine.Session.Rollback()
		log.Error("Update Order Database Error", err)
		return fmt.Errorf("系統錯誤！")
	}
	//支付平台費用
	if err := Balance.OrderShipDeduction(engine, &data); err != nil {
		_ = engine.Session.Rollback()
		log.Error("Order Ship Deduction Error", err)
		return fmt.Errorf("系統錯誤！")
	}
	if err := engine.Session.Commit(); err != nil {
		_ = engine.Session.Rollback()
		log.Error("Commit Database Error", err)
		return fmt.Errorf("系統錯誤！")
	}
	return nil
}

// 更新貨運狀態
func UpdateOrderDataShipStatus(engine *database.MysqlSession, OrderData entity.OrderData, newStatus string) (err error) {
	// 寫入貨運變更歷程
	err = SetOrderStatusLog(Enum.StatusLogOrderDataOrderShipStatus,  OrderData.OrderId, OrderData.OrderStatus, newStatus, "")
	if err != nil {
		log.Error("StatusLogOrderDataOrderShipStatus Error", err)
	}

	// 更新訂單
	OrderData.ShipStatus = newStatus
	log.Debug(`UpdateOrderDataShipStatus`, OrderData)
	_, err = Orders.UpdateOrderData(engine, OrderData.OrderId, OrderData)
	if err != nil {
		log.Error("UpdateOrderDataShipStatus Orders.UpdateOrderData Error", err)
		return err
	}

	return nil
}
//出貨單號、出貨業者、出貨狀態 變更
func ChangeOrderShipStatus(engine *database.MysqlSession, data entity.OrderData, trader, number string) error {
	data.ShipText = trader
	data.ShipNumber = number
	data.ShipStatus = Enum.OrderShipment
	data.ShipTime = time.Now()
	OrderService.OrderCaptureRelease(&data, time.Time{})
	//支付平台費用
	if err := Balance.OrderShipDeduction(engine, &data); err != nil {
		log.Error("Order Ship Deduction Error", err)
		return err
	}
	//變更訂單出貨狀態
	if _, err := Orders.UpdateOrderData(engine, data.OrderId, data); err != nil {
		log.Error("Update Order Database Error", err)
		return err
	}
	return nil
}
