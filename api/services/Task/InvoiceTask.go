package Task

import (
	"api/services/Enum"
	"api/services/Service/Invoice"
	"api/services/dao/Orders"
	"api/services/database"
	"api/services/model"
	"api/services/util/log"
	"time"
)

func HandleInvoiceTask() {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	if err := openB2cOrderInvoice(engine); err != nil {
		log.Error("open B2c Order Error")
	}
	if err := openC2cOrderInvoice(engine); err != nil {
		log.Error("open c2c Order Error")
	}
}

func openB2cOrderInvoice(engine *database.MysqlSession) error {
	//取出應開發票的訂單
	data, err := Orders.GetB2cOrderByNotInvoice(engine)
	if err != nil {
		log.Error("Get B2c Order Error", err)
		return err
	}
	for _, v := range data {
		//開發票
		if err := Invoice.ProcessCreateB2cInvoice(engine, v.OrderId); err != nil {
			log.Error("Create Invoice Error!!", err)
			return err
		}
		v.InvoiceStatus = Enum.InvoiceOpenStatusOpen
		if err := Orders.UpdateB2cOrderData(engine, v); err != nil {
			log.Error("Update B2c Order Data Error", err)
			return err
		}
	}
	return nil
}

func HandleDayStatementTask() {
	start := time.Now().Add(- 24 * time.Hour).Format("2006/01/02")
	end := time.Now().Format("2006/01/02")
	if err := model.GeneratorDayStatement(start, end); err != nil {
		log.Error("Day Statement Error", err)
	}
}

func openC2cOrderInvoice(engine *database.MysqlSession) error {
	data, err := Orders.GetOrderByNotInvoice(engine)
	if err != nil {
		log.Error("Get B2c Order Error")
		return err
	}
	for _, v := range data {
		//開發票
		if err := Invoice.ProcessCreateServiceInvoice(engine, v.OrderId); err != nil {
			log.Error("Create Invoice Error!!")
			return err
		}
		v.InvoiceStatus = Enum.InvoiceOpenStatusOpen
		_, err := Orders.UpdateOrderData(engine, v.OrderId, v)
		if err != nil {
			log.Error("Update B2c Order Data Error", err)
			return err
		}
	}
	return nil
}

func HandleInvoiceAssignNumber() {
	if err:= Invoice.GetInvoiceAssignNumber(); err != nil {
		log.Error("Get Invoice Assign Error", err)
	}
}
