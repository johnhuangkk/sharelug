package Withdraw

import (
	"api/services/Enum"
	"api/services/VO/Request"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"fmt"
	"strings"
	"time"
)

//建立提領資料
func InsertWithdrawData(engine *database.MysqlSession, data entity.WithdrawData) error {
	data.CreateTime = time.Now()
	data.UpdateTime = time.Now()
	_, err := engine.Session.Table(entity.WithdrawData{}).Insert(&data)
	if err != nil {
		log.Error("Insert Withdraw Database Error", err)
		return err
	}
	return nil
}

//取得提領資料
func GetEachWithdrawDataByStatus(engine *database.MysqlSession, Status string) ([]entity.WithdrawData, error) {
	var data []entity.WithdrawData
	bankCode := []string{"0040037", "0050418", "0060567", "0070937", "0081005", "0095185", "0110026", "0122009", "0130017", "0172015", "0480011", "0500108", "0540537", "1030019", "8030021", "8060219", "8070014", "8090267", "8100364", "8120012", "8150015", "8220901"}
	now := time.Now().Format("2006-01-02")
	sql := fmt.Sprintf("SELECT * FROM withdraw_data WHERE bank_code IN(%s) AND withdraw_status = ? AND create_time < ?", strings.Join(bankCode, ","))
	err := engine.Engine.SQL(sql, Status, now).Find(&data)
	if err != nil {
		return data, err
	}
	return data, nil
}

//取得提領資料
func GetAchWithdrawDataByStatus(engine *database.MysqlSession, Status string) ([]entity.WithdrawData, error) {
	var data []entity.WithdrawData
	bankCode := []string{"0040037", "0050418", "0060567", "0070937", "0081005", "0095185", "0110026", "0122009", "0130017", "0172015", "0480011", "0500108", "0540537", "1030019", "8030021", "8060219", "8070014", "8090267", "8100364", "8120012", "8150015", "8220901"}
	now := time.Now().Format("2006-01-02")
	sql := fmt.Sprintf("SELECT * FROM withdraw_data WHERE bank_code NOT IN(%s) AND withdraw_status = ? AND create_time < ?", strings.Join(bankCode, ","))
	err := engine.Engine.SQL(sql, Status, now).Find(&data)
	if err != nil {
		return data, err
	}
	return data, nil
}

//更新提領資料
func UpdateWithdrawDataByStatus(engine *database.MysqlSession, Data entity.WithdrawData, BatchId string) error {
	Data.WithdrawStatus = Enum.WithdrawStatusSuccess
	Data.BatchId = BatchId
	Data.ExportTime = time.Now()
	Data.UpdateTime = time.Now()
	_, err := engine.Session.Table(entity.WithdrawData{}).ID(Data.WithdrawId).Update(Data)
	if err != nil {
		return err
	}
	return nil
}

func GetAchWithdrawDataByTransId(engine *database.MysqlSession, transId string) (entity.WithdrawData, error) {
	var data entity.WithdrawData
	sql := fmt.Sprintf("SELECT * FROM withdraw_data WHERE trans_id = ?")
	if _, err := engine.Engine.SQL(sql, transId).Get(&data); err != nil {
		return data, err
	}
	return data, nil
}

//取出提領資料
func GetAchWithdrawDataByWithdrawId(engine *database.MysqlSession, WithdrawId string) (entity.WithdrawData, error) {
	var data entity.WithdrawData
	if _, err := engine.Engine.Table(entity.WithdrawData{}).Select("*").Where("withdraw_id = ?", WithdrawId).Get(&data); err != nil {
		return data, err
	}
	return data, nil
}

func GetWithdrawData(engine *database.MysqlSession, params Request.SearchWithdrawRequest) ([]entity.QueryWithdraw, error) {
	var data []entity.QueryWithdraw
	where, bind := ComposeSearchWithdrawParams(params.Search, params.Tabs)
	err := engine.Engine.Table(entity.WithdrawData{}).Select("*").
		Join("LEFT", entity.MemberData{}, "withdraw_data.user_id = member_data.uid").
		Where(strings.Join(where, " AND "), bind...).Find(&data)
	if err != nil {
		return data, err
	}
	return data, nil
}

func ComposeSearchWithdrawParams(params Request.SearchWithdraw, Tabs string) ([]string, []interface{}) {
	var where []string
	var bind []interface{}
	if len(Tabs) != 0 {
		where = append(where, "withdraw_data.withdraw_status = ?")
		bind = append(bind, Tabs)
	}
	if len(params.WithdrawDate) != 0 {
		date := strings.Split(params.WithdrawDate, "-")
		where = append(where, "withdraw_data.withdraw_time BETWEEN ? AND ?")
		bind = append(bind, date[0])
		bind = append(bind, date[1])
	}
	if len(params.Buyer) != 0 {
		where = append(where, "member_data.mphone = ?")
		bind = append(bind, params.Buyer)
	}
	if len(params.BuyerEmail) != 0 {
		where = append(where, "member_data.email = ?")
		bind = append(bind, params.BuyerEmail)
	}
	if len(params.WithdrawStatus) != 0 {
		where = append(where, "withdraw_data.response_status = ?")
		bind = append(bind, params.WithdrawStatus)
	}
	if len(params.WithdrawId) != 0 {
		where = append(where, "withdraw_data.withdraw_id = ?")
		bind = append(bind, params.WithdrawId)
	}
	return where, bind
}

func CountWithdrawDataByStatus(engine *database.MysqlSession, status string) int64 {
	count, err := engine.Engine.Table(entity.WithdrawData{}).Select("count(*)").Where("withdraw_status = ?", status).Count()
	if err != nil {
		log.Error("Count Withdraw Database Error", err)
		return count
	}
	return count
}

func CountWithdrawByUserId(engine *database.MysqlSession, userId string) int64 {
	count, err := engine.Engine.Table(entity.WithdrawData{}).Select("count(*)").
		Where("user_id = ?", userId).Count()
	if err != nil {
		log.Error("Count Withdraw Database Error", err)
		return 0
	}
	return count
}

func UpdateWithdrawDataByResponse(engine *database.MysqlSession, Data entity.WithdrawData) error {
	Data.ResponseTime = time.Now()
	_, err := engine.Session.Table(entity.WithdrawData{}).ID(Data.WithdrawId).Update(Data)
	if err != nil {
		return err
	}
	return nil
}

func UpdateWithdrawData(engine *database.MysqlSession, Data entity.WithdrawData) error {
	Data.UpdateTime = time.Now()
	_, err := engine.Session.Table(entity.WithdrawData{}).Where("withdraw_id = ?", Data.WithdrawId).Update(Data)
	if err != nil {
		return err
	}
	return nil
}

func CountWithdrawSuccessByUserId(engine *database.MysqlSession, userId string) int64 {
	count, err := engine.Engine.Table(entity.WithdrawData{}).Select("count(*)").
		Where("user_id = ? AND withdraw_status=?", userId, Enum.WithdrawStatusSuccess).Count()
	if err != nil {
		log.Error("Count Withdraw Database Error", err)
		return 0
	}
	return count
}
