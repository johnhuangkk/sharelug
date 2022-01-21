package CvsShipping

import (
	"api/services/Enum"
	"api/services/Service/Balance"
	"api/services/Service/Notification"
	"api/services/Service/OrderService"
	"api/services/Service/Sms"
	"api/services/dao/Cvs"
	"api/services/dao/Orders"
	"api/services/database"
	"api/services/entity"
	"api/services/model"
	"api/services/util/log"
	"fmt"
	"time"
)

type UpdateCvsShipping struct {
	// 運送編號
	ShipNo string
	// 訂單編號
	OrderNo      string
	ShipType     string
	Type         string
	DateTime     string
	DetailStatus string
	// N正向 R逆向
	FlowType string
	// 原始資訊
	Log      string
	FileName string
}

// 確認訂單資料
func checkOrderData(engine *database.MysqlSession, u *UpdateCvsShipping) (orderData entity.OrderData, err error) {
	if len(u.OrderNo) == 0 {
		orderData, err = Orders.GetOrderDataByShip(engine, u.ShipNo, u.ShipType)
	} else {
		orderData, err = Orders.GetOrderDataByShipTypeOrderId(engine, u.OrderNo, u.ShipType)
	}

	if err != nil {
		log.Error("UpdateCvsShippingShipment Error [%s]", err.Error())
		return orderData, fmt.Errorf("資料庫異常")
	}
	if len(orderData.OrderId) == 0 {
		return orderData, fmt.Errorf("訂單無此寄件編號")
	}

	return orderData, nil
}

// 更新超商訂單配送狀態 寫入log
func updateCvsShippingShipment(
	engine *database.MysqlSession,
	orderData entity.OrderData,
	u *UpdateCvsShipping,
	shippingStatus string) (err error) {

	// 更新貨運狀態
	err = model.UpdateOrderDataShipStatus(engine, orderData, shippingStatus)
	if err != nil {
		log.Error("UpdateOrderDataShipStatus Error [%s]", err.Error())
		return fmt.Errorf("資料庫異常")
	}

	u.ShipNo = orderData.ShipNumber

	// 寫入超商log
	err = writeCvsLogData(engine, u, orderData, true)
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	return nil
}

// 寫入cvsLog
func writeCvsLogData(
	engine *database.MysqlSession,
	u *UpdateCvsShipping,
	orderData entity.OrderData,
	show bool) (err error) {

	var cvsDataLog entity.CvsShippingLogData

	cvsDataLog.IsShow = show
	cvsDataLog.CvsType = u.ShipType
	cvsDataLog.ShipNo = orderData.ShipNumber
	cvsDataLog.Type = u.Type
	cvsDataLog.DateTime = u.DateTime
	cvsDataLog.Log = u.Log
	cvsDataLog.SetLogDataText(u.DetailStatus)
	cvsDataLog.FileName = u.FileName

	log.Info("cvsDataLog [%s] [%s]", u.Type, cvsDataLog)

	// 寫入超商log
	err = Cvs.InsertCvsShippingLogData(engine, cvsDataLog)
	if err != nil {
		return err
	}

	return nil
}

// 成功交寄寫入交寄時間
func (u *UpdateCvsShipping) updateCvsDataSendTime(engine *database.MysqlSession, orderData entity.OrderData) (err error) {
	// 成功交寄寫入交寄時間
	shippingData := Cvs.GetCvsShippingData(engine, orderData)
	shippingData.SendTime = orderData.ShipTime

	return Cvs.UpdateCvsShippingData(engine, shippingData)

}

// 成功寫入取件時間
func (u *UpdateCvsShipping) updateCvsDataReceiveTime(engine *database.MysqlSession, orderData entity.OrderData) (err error) {
	// 成功寫入取件時間
	shippingData := Cvs.GetCvsShippingData(engine, orderData)
	shippingData.ReceiveTime, err = time.Parse(`2006-01-02 15:04:05`, u.DateTime)
	if err != nil {
		log.Error(`updateCvsDataSendTime Error`, err.Error())
		return err
	}
	return Cvs.UpdateCvsShippingData(engine, shippingData)

}

// 閉轉通知
func (u *UpdateCvsShipping) UpdateCvsShippingSwitch(engine *database.MysqlSession) error {
	orderData, err := checkOrderData(engine, u)
	if err != nil {
		return err
	}

	//todo 發送訊息 買｜賣家

	// 取得配送資訊
	shippingData := Cvs.GetCvsShippingData(engine, orderData)
	shippingData.FlowType = u.FlowType
	shippingData.Switch = `1`
	// 兩天的閉轉變更時間
	shippingData.SwitchDeadline = time.Now().Add(time.Hour * 48).Format(`2006-01-02`)

	u.ShipNo = orderData.ShipNumber
	// 更新配送資訊
	err = Cvs.UpdateCvsShippingData(engine, shippingData)
	if err != nil {
		return err
	}

	// 小鈴鐺通知
	if u.FlowType == `R` {
		err = Notification.SendReturnShipClosedShopMessage(engine, orderData)
	} else {
		err = Notification.SendToBuyerShipClosedShopMessage(engine, orderData)
	}

	if err != nil {
		return err
	}

	// 寫入log
	return writeCvsLogData(engine, u, orderData, true)
}

// 訂單貨運狀態為 未出貨 則修改狀態為 已出貨
func (u *UpdateCvsShipping) UpdateCvsShippingShipment(engine *database.MysqlSession) error {

	orderData, err := checkOrderData(engine, u)
	if err != nil {
		return err
	}

	if Enum.OrderShipTake == orderData.ShipStatus {
		// 非超取付時 交寄後變成可撥款
		if orderData.PayWay != Enum.CvsPay {
			OrderService.OrderCaptureRelease(&orderData, time.Time{})
		}
		_ = Balance.OrderShipDeduction(engine, &orderData)
		// 回寫訂單出貨時間
		orderData.ShipTime, err = time.ParseInLocation(`2006-01-02 15:04:05`, u.DateTime, time.Local)
		if err != nil {
			return err
		}
		err = updateCvsShippingShipment(engine, orderData, u, Enum.OrderShipment)
		if err != nil {
			return err
		}

		// 寫入交寄時間
		err = u.updateCvsDataSendTime(engine, orderData)
		if err != nil {
			return err
		}

		// 小鈴鐺通知
		return Notification.SendShippedMessage(engine, orderData)
	}
	return nil
}

// 訂單貨運狀態為 已出貨 則修改狀態為 配送中
func (u *UpdateCvsShipping) UpdateCvsShippingTransit(engine *database.MysqlSession) error {

	orderData, err := checkOrderData(engine, u)
	if err != nil {
		return err
	}

	if u.FlowType == `R` {
		return writeCvsLogData(engine, u, orderData, true)
	} else if Enum.OrderShipment == orderData.ShipStatus {
		return updateCvsShippingShipment(engine, orderData, u, Enum.OrderShipTransit)
	}

	return nil
}

// 訂單貨運狀態為 配送中 則修改狀態為 到店
func (u *UpdateCvsShipping) UpdateCvsShippingShop(engine *database.MysqlSession) error {

	orderData, err := checkOrderData(engine, u)
	if err != nil {
		return err
	}

	if u.FlowType == `R` {

		if orderData.ShipType == Enum.CVS_OK_MART {
			// OK逆向 只看日檔發送簡訊
			if u.ShipType == `F64` || u.ShipType == `F04` {
				Sms.CvsPushShipStatusShopMessageSms(engine, orderData)
			}
		}

		return writeCvsLogData(engine, u, orderData, true)
	} else if Enum.OrderShipTransit == orderData.ShipStatus {
		// 商品到店時間寫入訂單 用於算 1 4 7 天要發送提醒通知
		orderData.ArrivedTime = time.Now()
		err = updateCvsShippingShipment(engine, orderData, u, Enum.OrderShipShop)
		if err != nil {
			return err
		}
		err = Notification.SendShipToShopFirstDayMessage(engine, orderData)
		if err != nil {
			return err
		}

		// OK正向到店時 須發送簡訊
		if orderData.ShipType == Enum.CVS_OK_MART {
			Sms.CvsPushShipStatusShopMessageSms(engine, orderData)
		}
	}

	return nil
}

// 訂單貨運狀態為 到店 則修改狀態為 取貨
func (u *UpdateCvsShipping) UpdateCvsShippingSuccess(engine *database.MysqlSession) error {

	orderData, err := checkOrderData(engine, u)

	if err != nil {
		return err
	}

	if u.FlowType == `R` {
		return writeCvsLogData(engine, u, orderData, true)
	} else if Enum.OrderShipShop == orderData.ShipStatus {
		// 超取付取貨時 才可變成撥付
		if orderData.PayWay == Enum.CvsPay {
			orderData.PayWayTime = time.Now()
			OrderService.OrderCaptureRelease(&orderData, time.Time{})
			if err := Balance.RetainDeposit(orderData.SellerId, orderData.OrderId, orderData.TotalAmount, Enum.BalanceTypeDeposit, "訂單交易存入"); err != nil {
				return err
			}
		}
		err = u.updateCvsDataReceiveTime(engine, orderData)
		if err != nil {
			log.Error(`cvs receive time update`, err.Error())
			return err
		}
		return updateCvsShippingShipment(engine, orderData, u, Enum.OrderShipSuccess)
	}

	return nil
}

// 買家未取退回
func (u *UpdateCvsShipping) UpdateCvsShippingBuyerNotPickUp(engine *database.MysqlSession) error {
	orderData, err := checkOrderData(engine, u)
	if err != nil {
		return err
	}

	shippingData := Cvs.GetCvsShippingData(engine, orderData)

	log.Debug("shippingData", shippingData)

	shippingData.FlowType = `R`
	shippingData.StateCode = u.DetailStatus

	err = Cvs.UpdateCvsShippingData(engine, shippingData)

	if err != nil {
		return err
	}
	err = updateCvsShippingShipment(engine, orderData, u, Enum.OrderShipNotTaken)
	if err != nil {
		return err
	}

	return Notification.SendOrderOverdueMessage(engine, orderData)
}

// 訂單貨運狀態修改為 配送失敗
func (u *UpdateCvsShipping) UpdateCvsShippingFail(engine *database.MysqlSession) error {

	orderData, err := checkOrderData(engine, u)
	if err != nil {
		return err
	}

	shippingData := Cvs.GetCvsShippingData(engine, orderData)

	log.Debug("shippingData", shippingData)

	shippingData.FlowType = u.FlowType
	shippingData.StateCode = u.DetailStatus

	// 退貨逆向物流
	switch u.ShipType {
	case Enum.CVS_FAMILY, Enum.CVS_HI_LIFE:
		if u.Type == `R08` && u.FlowType == `N` {
			shippingData.FlowType = `R`
		}
	case Enum.CVS_OK_MART:
		if (u.Type == `F07` || u.Type == `F67`) && u.DetailStatus == `T00` {
			shippingData.FlowType = `R`
		}
	}

	err = Cvs.UpdateCvsShippingData(engine, shippingData)

	if err != nil {
		return err
	}
	if (u.ShipType == Enum.CVS_HI_LIFE || u.ShipType == Enum.CVS_FAMILY) && u.DetailStatus == `XXX` {
		return updateCvsShippingShipment(engine, orderData, u, Enum.OrderShipOverdue)
	}

	return updateCvsShippingShipment(engine, orderData, u, Enum.OrderShipFail)
}

// 只寫入運送狀態
func (u *UpdateCvsShipping) OnlyWriteShippingLog(engine *database.MysqlSession, show bool) error {
	orderData, err := checkOrderData(engine, u)
	if err != nil {
		return err
	}

	err = writeCvsLogData(engine, u, orderData, show)

	if err != nil {
		return err
	}

	return nil
}
