package model

import (
	"api/services/Enum"
	"api/services/Service/Balance"
	"api/services/Service/OrderService"
	"api/services/VO/IPOSTVO"
	"api/services/dao/Orders"
	"api/services/dao/iPost"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"api/services/util/tools"
	"fmt"
	"strings"
	"time"
)

// 新增郵箱貨態
func InsertPostShippingStatus(record []string)  {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var postShippingStatus = &entity.PostShippingStatus{}
	var orderData entity.OrderData

	/**
	[0]
	[1]96783010005518
	[2]2018-06-04
	[3]10:18:26
	[4]運輸途中
	[5]臺中郵局特投股
	*/
	MailNo := strings.Replace(record[1], " ", "", -1)
	MailNo = strings.Replace(MailNo, "\t", "", -1)
	postShippingStatus.MailNo = tools.Big5ToUtf8(MailNo)
	postShippingStatus.HandleTime = tools.Big5ToUtf8(record[2] + " " + record[3])
	postShippingStatus.ShippingStatus = tools.Big5ToUtf8(record[4])
	postShippingStatus.Branch = tools.Big5ToUtf8(record[5])
	postShippingStatus.CreateTime = tools.Now(`YmdHis`)
	postShippingStatus.Detail = tools.Big5ToUtf8(strings.Join(record, ","))

	// 排除頭尾字
	if len(postShippingStatus.MailNo) == 14 {

		if postShippingStatus.ShippingStatus == `投遞成功` {
			// i郵箱 箱到箱
			orderData, _ = Orders.GetOrderDataByShip(engine, postShippingStatus.MailNo, Enum.I_POST)
			// i郵箱 箱到箱找不到 找 箱到宅
			if len(orderData.OrderId) == 0 {
				orderData, _ = Orders.GetOrderDataByShip(engine, postShippingStatus.MailNo, Enum.DELIVERY_I_POST_BAG1)
			}

			// 有訂單編號才回寫成功
			if len(orderData.OrderId) > 0 {
				if orderData.ShipStatus != Enum.OrderShipSuccess {
					_ = UpdateOrderDataShipStatus(engine, orderData, Enum.OrderShipSuccess)
				}
			}
		}

		err := iPost.InsertPostShippingStatus(engine, *postShippingStatus)
		if err != nil {
			log.Error("InsertPostShippingStatus fail %v", postShippingStatus)
		}
	}
}

// 郵局即時 Api 通知
func NotificationInsertPostShippingStatus(params IPOSTVO.ShipStatusNotify) error {

	engine := database.GetMysqlEngine()
	defer engine.Close()

	var postShippingStatus = entity.PostShippingStatus{}
	postShippingStatus.SetData(params)

	// 撈出訂單
	orderData, err := Orders.GetOrderDataByShip(engine, postShippingStatus.MailNo, Enum.I_POST)

	if err != nil {
		log.Error("GetOrderDataByShip Error [%v] ", err.Error())
		return fmt.Errorf("系統錯誤")
	}

	if len(orderData.OrderId) == 0 {
		orderData, err = Orders.GetOrderDataByShip(engine, postShippingStatus.MailNo, Enum.DELIVERY_I_POST_BAG1)
		if err != nil {
			log.Error("GetOrderDataByShip Error [%v] ", err.Error())
			return fmt.Errorf("系統錯誤")
		}
	}

	var updateStatus string

	switch params.ShipStatus {
	case `30`: // i 郵箱收寄(30)
		updateStatus = Enum.OrderShipment
		OrderService.OrderCaptureRelease(&orderData, time.Time{})
		_ = Balance.OrderShipDeduction(engine, &orderData)
		// 回寫訂單出貨時間
		orderData.ShipTime, _ = time.Parse(`2006-01-02 15:04:05`, postShippingStatus.HandleTime)
	case `20`: // 到達買家 i 郵箱(20)
		//寫入抵達時間
		orderData.ArrivedTime, _ = time.Parse(`2006-01-02 15:04:05`, postShippingStatus.HandleTime)
		updateStatus = Enum.OrderShipShop
	case `80`: // 買家 i 郵箱取件成功(80)
		updateStatus = Enum.OrderShipSuccess
	case `120`:
		//到達賣家 i 郵箱(120) todo 小鈴鐺
	default:
		updateStatus = Enum.OrderShipFail
	}


	log.Info(`orderData`, orderData)
	log.Info(`updateStatus`, updateStatus)
	if orderData.ShipStatus != updateStatus {
		if UpdateOrderDataShipStatus(engine, orderData, updateStatus) != nil {
			log.Error("UpdateOrderDataShipStatus postShippingStatus Error [%v] ", postShippingStatus)
			return fmt.Errorf("訂單更新貨態失敗")
		}
	}

	err = iPost.InsertPostShippingStatus(engine, postShippingStatus)
	if err != nil {
		log.Error("InsertPostShippingStatus Error [%v] ", postShippingStatus, err)
		return fmt.Errorf("新增郵箱貨態錯誤")
	}

	return nil
}
