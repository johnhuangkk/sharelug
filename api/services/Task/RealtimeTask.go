package Task

import (
	"api/services/Enum"
	"api/services/Service/Balance"
	"api/services/dao/Orders"
	"api/services/dao/product"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"fmt"
	"time"
)

//帳單過期狀態變更為過期
func HandleRealtimeExpireTask() {
	log.Info("HandleRealtimeExpireTask [start]")
	engine := database.GetMysqlEngine()
	defer engine.Close()
	now := time.Now().Format("2006-01-02 15:04")
	data, err := product.GetRealtimeProductExpire(engine, now)
	if err != nil {
		log.Error("Get Product Realtime Expire Error", err)
	}
	for _, value := range data {
		err := product.UpdateProductStatusByProductId(engine, value.ProductId, Enum.ProductStatusOverdue)
		if err != nil {
			log.Error("Update Product Realtime Expire Error", err)
		}
	}
	log.Info("HandleRealtimeExpireTask [End]")
}

//帳單過期狀態變更為過期
func HandleBillExpireTask() {
	log.Info("HandleBillExpireTask [start]")
	engine := database.GetMysqlEngine()
	defer engine.Close()
	now := time.Now().Format("2006-01-02 15:04")
	data, err := Orders.GetBillExpire(engine, now)
	if err != nil {
		log.Error("Get Product Realtime Expire Error", err)
	}
	for _, value := range data {
		//判斷付款方式
		if value.PayWayType == Enum.Balance {
			comment := fmt.Sprintf("%s", value.BillId) //訂購單編號
			err = Balance.Deposit(value.BuyerId, value.BillId, value.TotalAmount, Enum.BalanceTypeBillFail, comment)
			if err != nil {
				log.Error("Balance Deposit Error", err)
			}
		}
		if value.PayWayType == Enum.Credit {
			var vo entity.CancelRequest
			vo.OrderId = value.BillId
			if err := Balance.VoidProcess(engine, vo); err != nil {
				log.Error("credit void Error", err)
			}
		}
		value.BillStatus = Enum.BillStatusOverdue
		err := Orders.UpdateBillOrderData(engine, value)
		if err != nil {
			log.Error("Update Product Realtime Expire Error", err)
		}
	}
	log.Info("HandleBillExpireTask [End]")
}

func HandleShippingExpireTask()  {
	log.Info("HandleShippingExpireTask [start]")
	engine := database.GetMysqlEngine()
	defer engine.Close()

	log.Info("HandleRealtimeExpireTask [End]")
}

func HandleCloseAllProductTask() {
	log.Info("HandleRealtimeExpireTask [start]")
	engine := database.GetMysqlEngine()
	defer engine.Close()
	day, _ := time.ParseInLocation("20060102 150405", "20210930 170000", time.Local)
	now := time.Now()
	if !now.Before(day) {
		if err:= product.CloseAllProduct(engine); err != nil {
			log.Error("Update Product status Error", err)
		}
	}
	log.Info("HandleRealtimeExpireTask [End]")
}
