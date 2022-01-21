package Task

import (
	"api/services/Enum"
	"api/services/Service/Balance"
	"api/services/Service/Notification"
	"api/services/dao/Orders"
	"api/services/database"
	"api/services/util/log"
	"time"
)

//撥款處理
func HandleAppropriationTask() {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	now := time.Now().Format("2006-01-02 15:04")
	data, err := Orders.GetAppropriationOrder(engine, now)
	if err != nil {
		log.Error("Get Appropriation Order Error", err)
	}
	log.Debug("get appropriation order", data)
	for _, value := range data {
		//付款狀態為 Enum.CvsPay 時要檢查 csv_check 是否等於1
		if value.ShipType == Enum.CvsPay && value.CsvCheck == 1 || value.ShipType != Enum.CvsPay {
			//撥款
			err := Balance.OrderAppropriation(engine, value)
			if err != nil {
				log.Error("Order Appropriation Error", err)
			}
		}
	}
}

//逾期未寄排程 抓ORDER ship_status 及 ship_expire
func HandleShipExpireTask() {
	log.Info("HandleShipExpireTask [start]")
	engine := database.GetMysqlEngine()
	defer engine.Close()
	now := time.Now().Format("2006-01-02 15:04")
	data, err := Orders.GetOrderShipExpire(engine, now)
	if err != nil {
		log.Error("Get Ship Expire Order Error", err)
	}
	for _, v := range data {
		v.ShipStatus = Enum.OrderShipOverdue
		_, err := Orders.UpdateOrderData(engine, v.OrderId, v)
		if err != nil {
			log.Error("Update Order Ship Status Error", err)
		}
		if err := Notification.SendOrderShipOverdueMessage(engine, v); err != nil {
			log.Error("Send Order Ship Expire Message Error", err)
		}
	}
	log.Info("HandleShipExpireTask [end]")
}
