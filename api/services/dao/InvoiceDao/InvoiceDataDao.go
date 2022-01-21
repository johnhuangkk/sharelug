package InvoiceDao

import (
	"api/services/Enum"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"api/services/util/tools"
	"time"
)

func InsertInvoiceData(engine *database.MysqlSession, data entity.InvoiceData) error {
	if _, err := engine.Session.Table(entity.InvoiceData{}).Insert(&data); err != nil {
		log.Error("Insert Invoice Database Error", err)
		return err
	}
	return nil
}
//取出發票資料
func GetInvoiceByOrderId(engine *database.MysqlSession, orderId string) (entity.InvoiceData, error) {
	var data entity.InvoiceData
	_, err := engine.Engine.Table(entity.InvoiceData{}).
		Select("*").Where("order_id = ? ", orderId).Get(&data)
	if err != nil {
		return data, err
	}
	return data, nil
}

/**
 * 計算搜尋結果數
 */
func CountInvoiceListByBuyerId(engine *database.MysqlSession, BuyerId string) (int64, error) {
	result, err := engine.Engine.Table(entity.InvoiceData{}).Where("invoice_status != ?", Enum.InvoiceStatusCancel).
		And("buyer_id = ? ", BuyerId).Count()
	if err != nil {
		log.Error("count search order Database Error", err)
		return 0, err
	}
	return result, nil
}

//取出發票資料
func GetInvoiceListByBuyerId(engine *database.MysqlSession, BuyerId string, limit, start int) ([]entity.InvoiceData, error) {
	limit = tools.CheckIsZero(limit, 10)
	start = tools.CheckIsZero(start, 1)
	start = (start - 1) * limit
	var data []entity.InvoiceData
	if err := engine.Engine.Table(entity.InvoiceData{}).Select("*").Where("invoice_status != ?", Enum.InvoiceStatusCancel).
		And("buyer_id = ? ", BuyerId).Desc("create_time").Limit(limit, start).Find(&data); err != nil {
		return data, err
	}
	return data, nil
}

func UpdateInvoiceAllData(engine *database.MysqlSession, year, month string) error {
	sql := "UPDATE invoice_data SET invoice_status = ? WHERE year = ? AND month = ? AND invoice_status = ?"
	if _, err := engine.Session.Exec(sql, Enum.InvoiceStatusLose, year, month, Enum.InvoiceStatusNot); err != nil {
		log.Error("Update Database Error", err)
		return  err
	}
	return nil
}

func UpdateInvoiceData(engine *database.MysqlSession, data entity.InvoiceData) error {
	if _, err := engine.Session.Table(entity.InvoiceData{}).Where("invoice_id = ?", data.InvoiceId).Update(data); err != nil {
		log.Error("Update Database Error", err)
		return  err
	}
	return nil
}

//取出
func GetInvoiceAssignNoByYearMonth(engine *database.MysqlSession, yearMonth string) (entity.InvoiceAssignNoData, error) {
	var data entity.InvoiceAssignNoData
	_, err := engine.Engine.Table(entity.InvoiceAssignNoData{}).Select("*").
		Where("month_year = ? AND invoice_status = ?", yearMonth, Enum.InvoiceAssignStatusEnable).
		Asc("create_time").Get(&data)
	if err != nil {
		return data, err
	}
	return data, nil
}

func GetInvoiceAssignNoByYearMonthAndBooklet(engine *database.MysqlSession, yearMonth, begin, end, track string) (entity.InvoiceAssignNoData, error) {
	var data entity.InvoiceAssignNoData
	_, err := engine.Engine.Table(entity.InvoiceAssignNoData{}).
		Select("*").Where("month_year = ?", yearMonth).And("invoice_track = ?", track).
		And("invoice_begin_no = ?", begin).And("invoice_end_no = ?", end).Get(&data)
	if err != nil {
		return data, err
	}
	return data, nil
}
//新增發票配號
func InsertInvoiceAssignNoData(engine *database.MysqlSession, data entity.InvoiceAssignNoData) error {
	data.CreateTime = time.Now()
	data.UpdateTime = time.Now()
	if _, err := engine.Session.Table(entity.InvoiceAssignNoData{}).Insert(&data); err != nil {
		log.Error("Insert Invoice Assign No Database Error", err)
		return err
	}
	return nil
}
//更新發票配號
func UpdateInvoiceAssignNoData(engine *database.MysqlSession, data entity.InvoiceAssignNoData) error {
	data.UpdateTime = time.Now()
	if _, err := engine.Session.Table(entity.InvoiceAssignNoData{}).Where("assign_id = ?", data.AssignId).Update(data); err != nil {
		log.Error("Update Database Error", err)
		return  err
	}
	return nil
}
//取出中獎發票資料
func GetInvoiceByTrackAndNumber(engine *database.MysqlSession, invoiceTrack, invoiceNumber string) (entity.InvoiceData, error) {
	var data entity.InvoiceData
	if _, err := engine.Engine.Table(entity.InvoiceData{}).Select("*").Where("invoice_track = ?", invoiceTrack).
		And("invoice_number = ?", invoiceNumber).Get(&data); err != nil {
		return data, err
	}
	return data, nil
}

//取出發票資料
func GetInvoiceByYearMonth(engine *database.MysqlSession, Year, Month string) ([]entity.InvoiceData, error) {
	var data []entity.InvoiceData
	if err := engine.Engine.Table(entity.InvoiceData{}).
		Select("*").Where("year = ? ", Year).And("month = ?", Month).Find(&data); err != nil {
		return data, err
	}
	return data, nil
}


func GetInvoiceAndOrderByUserid(engine *database.MysqlSession, Userid, date string) ([]entity.InvoiceData, error) {
	var data []entity.InvoiceData
	if err := engine.Engine.Table(entity.InvoiceData{}).Select("*").
		Where("buyer_id = ? ", Userid).And("create_time < ?", date).Find(&data); err != nil {
		return data, err
	}
	return data, nil
}
