package CvsAccounting

import (
	"api/services/Enum"
	"api/services/dao/Cvs"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"fmt"
	"time"
)

// 寄件核帳
func Send(engine *database.MysqlSession, cvsType string, limitTime string) {
	var Data []struct {
		A entity.CvsAccountingData `xorm:"extends"`
		B entity.CvsShippingData   `xorm:"extends"`
	}

	type_ := `S` // 寄件
	condition := `cvs_shipping_data.ec_order_no = cvs_accounting_data.data_id`

	// 萊爾富 dataId 對應 shipNo
	if cvsType == Enum.CVS_HI_LIFE {
		condition = `cvs_shipping_data.ship_no = cvs_accounting_data.data_id`
	}

	query := map[string]interface{}{`check_send`: 0, `checked`: 0, `cvs_accounting_data.cvs_type`: cvsType, `type`: type_}

	if err := engine.Engine.Table(entity.CvsAccountingData{}).
		Select("*").
		Join(`LEFT`, entity.CvsShippingData{}, condition).
		Where(query).
		And(`send_time <> ?`, ``).
		And(`cvs_accounting_data.create_time > ?`, limitTime).Find(&Data); err != nil {
		log.Error(`err`, err.Error())
	}

	for _, a := range Data {
		if err := engine.Session.Begin(); err != nil {
			log.Error(`Begin err`, err.Error())
		}

		queryA := map[string]interface{}{`type`: type_, `data_id`: a.A.DataId}
		updateAccountingData := entity.CvsAccountingData{}
		updateAccountingData.Checked = true

		if err := Cvs.UpdateAccountingDataForClos(engine, updateAccountingData, queryA, []string{`checked`}); err != nil {
			log.Error(`UpdateAccountingData err`, err.Error())
			if err := engine.Session.Rollback(); err != nil {
				log.Error(`Rollback err`, err.Error())
				continue
			}
		}

		cvsShippingData := entity.CvsShippingData{}
		cvsShippingData.Id = a.B.Id
		cvsShippingData.CheckSend = true

		if err := Cvs.UpdateCvsShippingDataForFields(engine, cvsShippingData, []string{`check_send`}); err != nil {
			log.Error(`UpdateCvsShippingData err`, err.Error())
			if err := engine.Session.Rollback(); err != nil {
				log.Error(`Rollback err`, err.Error())
				continue
			}
		}

		if err := engine.Session.Commit(); err != nil {
			log.Error(`Commit err`, err.Error())
		}
	}
}

// 收件核帳
func Receive(engine *database.MysqlSession, cvsType string, limitTime string) {
	var Data []struct {
		A entity.CvsAccountingData `xorm:"extends"`
		B entity.OrderData         `xorm:"extends"`
	}

	type_ := `P` // 收件
	query := map[string]interface{}{
		`csv_check`: 0,
		`ship_type`: cvsType, `pay_way`: Enum.CvsPay,
		`order_status`: Enum.OrderSuccess, `ship_status`: Enum.OrderShipSuccess,
		`service_type`: 1, `checked`: 0, `cvs_accounting_data.cvs_type`: cvsType, `type`: type_,
	}

	condition := `order_data.order_id = cvs_accounting_data.data_id`
	// 萊爾富 dataId 對應 shipNo
	if cvsType == Enum.CVS_HI_LIFE {
		condition = `order_data.ship_number = cvs_accounting_data.data_id`
	}

	if err := engine.Engine.Table(entity.CvsAccountingData{}).
		Select("*").
		Join(`LEFT`, entity.OrderData{}, condition).
		Where(query).And(`cvs_accounting_data.create_time > ?`, limitTime).Find(&Data); err != nil {
		log.Error(`err`, err.Error())
	}

	for _, a := range Data {
		if cvsType == Enum.CVS_7_ELEVEN {
			if (a.B.TotalAmount - a.B.ShipFee) != a.A.Amount {
				errMsg := fmt.Sprintf(`orderId : %s [%f, %f]`, a.B.OrderId, a.A.Amount, a.B.TotalAmount)
				log.Error(`Receive Error`, errMsg)
				continue
			}
		} else {
			if a.B.TotalAmount != a.A.Amount {
				errMsg := fmt.Sprintf(`orderId : %s [%f, %f]`, a.B.OrderId, a.A.Amount, a.B.TotalAmount)
				log.Error(`Receive Error`, errMsg)
				continue
			}
		}

		if err := engine.Session.Begin(); err != nil {
			log.Error(`Begin err`, err.Error())
		}

		updateAccountingData := entity.CvsAccountingData{}
		updateAccountingData.Checked = true

		queryA := map[string]interface{}{`type`: type_, `data_id`: a.A.DataId}
		if err := Cvs.UpdateAccountingDataForClos(engine, updateAccountingData, queryA, []string{`checked`}); err != nil {
			log.Error(`UpdateAccountingData err`, err.Error())
			if err := engine.Session.Rollback(); err != nil {
				log.Error(`Rollback err`, err.Error())
				continue
			}
		}

		orderData := entity.OrderData{}
		orderData.OrderId = a.B.OrderId
		orderData.CsvCheck = 1
		orderData.CaptureTime = time.Now()

		if _, err := engine.Session.Table(entity.OrderData{}).ID(a.B.OrderId).Update(orderData); err != nil {
			log.Error(`UpdateOrderData err`, err.Error())
			if err := engine.Session.Rollback(); err != nil {
				log.Error(`Rollback err`, err.Error())
				continue
			}
		}

		if err := engine.Session.Commit(); err != nil {
			log.Error(`Commit err`, err.Error())
		}
	}
}
