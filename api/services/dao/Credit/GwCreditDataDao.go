package Credit

import (
	"api/services/Enum"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"fmt"
	"time"
)

//新增 GwCreditData
func InsertGwCreditData(engine *database.MysqlSession, data entity.GwCreditAuthData) (entity.GwCreditAuthData, error) {
	data.CreateTime = time.Now()
	data.UpdateTime = time.Now()
	if _, err := engine.Session.Table(entity.GwCreditAuthData{}).Insert(&data); err != nil {
		log.Error("Insert Gw Credit Database Error", err)
		return data, err
	}
	return data, nil
}
//更新 GwCreditData
func UpdateGwCreditData(engine *database.MysqlSession, data entity.GwCreditAuthData) error {
	if _, err := engine.Session.Table(entity.GwCreditAuthData{}).Where("auth_id = ?", data.AuthId).Update(data); err != nil {
		return err
	}
	return nil
}
//新增 GwCreditAuthLog
func InsertGwCreditAuthLog(engine *database.MysqlSession, params string) error {
	var data entity.GwCreditAuthLog
	data.Response = params
	data.CreateTime = time.Now()
	if _, err := engine.Session.Table(entity.GwCreditAuthLog{}).Insert(&data); err != nil {
		log.Error("Insert Gw Credit Auth Log Database Error", err)
		return err
	}
	return nil
}
//取GwCredit By OrderId And Tran Type
func GetGwCreditByOrderIdAndTranType(engine *database.MysqlSession, OrderId, TranType string) (entity.GwCreditAuthData, error) {
	var data entity.GwCreditAuthData
	if _, err := engine.Engine.Table(entity.GwCreditAuthData{}).
		Select("*").Where("order_id = ? ", OrderId).And("trans_type = ?", TranType).Get(&data); err != nil {
		log.Error("Get Gw Credit Database Error", err)
		return data, err
	}
	return data, nil
}
//取GwCredit By OrderId
func GetGwCreditByOrderId(engine *database.MysqlSession, OrderId string) (entity.GwCreditAuthData, error) {
	var data entity.GwCreditAuthData
	if _, err := engine.Engine.Table(entity.GwCreditAuthData{}).
		Select("*").Where("order_id = ? ", OrderId).Get(&data); err != nil {
		log.Error("Get Gw Credit Database Error", err)
		return data, err
	}
	return data, nil
}
//取GwCredit By Success
func GetGwCreditBySuccess(engine *database.MysqlSession, payType, creditType string) ([]entity.GwCreditAuthData, error) {
	var data []entity.GwCreditAuthData
	if err := engine.Engine.Table(entity.GwCreditAuthData{}).Select("*").
		Where("batch_id = ?", 0).And("pay_type = ?", payType).And("credit_type = ?", creditType).
		And("trans_type != ?", Enum.CreditTransTypeVoid).And("trans_status = ?",  Enum.OrderSuccess).
		And("audit_status = ?", Enum.CreditAuditRelease).Desc("trans_time").Find(&data); err != nil {
		log.Error("Get Gw Credit Database Error", err)
		return data, err
	}
	return data, nil
}
// 更新 Credit Gw
func UpdateGwCreditBatchId(engine *database.MysqlSession, OrderId, AuthCode string, BatchId string) error {
	sql := fmt.Sprintf("UPDATE gw_credit_auth_data SET batch_id = ?, batch_time = ? WHERE order_id = ? AND approve_code = ?")
	_, err := engine.Session.Exec(sql, BatchId, time.Now(), OrderId, AuthCode)
	if err != nil {
		return err
	}
	return nil
}
//取GwCredit By OrderId
func GetGwCreditByOrderIdAndTransType(engine *database.MysqlSession, OrderId, transType string) (entity.GwCreditAuthData, error) {
	var data entity.GwCreditAuthData
	if _, err := engine.Engine.Table(entity.GwCreditAuthData{}).Select("*").
		Where("order_id = ? AND trans_type = ?", OrderId, transType).Get(&data); err != nil {
		log.Error("Get Gw Credit Database Error", err)
		return data, err
	}
	return data, nil
}
//變更狀能
func ChangeGwCreditStatus(engine *database.MysqlSession, Data entity.GwCreditAuthData) error {
	Data.UpdateTime = time.Now()
	if err := UpdateGwCreditData(engine, Data); err != nil {
		return err
	}
	return nil
}
//取出信用卡刷卡資料
func GetCreditByOrderId(engine *database.MysqlSession, OrderId string) (entity.GwCreditAuthData, error) {
	var data entity.GwCreditAuthData
	if _, err := engine.Engine.Table(entity.GwCreditAuthData{}).Select("*").
		Where("trans_type =?", Enum.CreditTransTypeAuth).And("order_id = ?", OrderId).Get(&data); err != nil {
		log.Error("Get Gw Credit Database Error", err)
		return data, err
	}
	return data, nil
}
//取出信用卡刷卡審單資料
func GetAllGwCreditDataByAuditStatus(engine *database.MysqlSession, AuditStatus string) ([]entity.GwCreditAuthData, error) {
	var data []entity.GwCreditAuthData
	if err := engine.Engine.Table(entity.GwCreditAuthData{}).Select("*").
		Where("trans_type =?", Enum.CreditTransTypeAuth).And("audit_status = ?", AuditStatus).
		Desc("trans_time").Find(&data); err != nil {
		log.Error("Get Gw Credit Database Error", err)
		return data, err
	}
	return data, nil
}
//統計信用卡
func CountGwCreditByAuditStatus(engine *database.MysqlSession, transType, auditStatus string) (int64, error) {
	count, err := engine.Engine.Table(entity.GwCreditAuthData{}).Select("count(*)").
		Where("trans_type =?", transType).And("audit_status = ?", auditStatus).Count()
	if err != nil {
		log.Error("Count Gw Credit Database Error", err)
		return count, err
	}
	return count, nil
}
//新增銀行BinCode
func InsertBankBinCodeData(data entity.BankBinCode) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	data.CreateTime = time.Now()
	data.UpdateTime = time.Now()
	if _, err := engine.Session.Table(entity.BankBinCode{}).Insert(&data); err != nil {
		log.Error("Insert Bank Bin Code Database Error", err)
		return err
	}
	return nil
}
//信用卡銀行BIN CODE
func GetBankBinCodeByBinCode(engine *database.MysqlSession, code string) ([]entity.BankBinCode, error) {
	var data []entity.BankBinCode
	if err := engine.Session.Table(entity.BankBinCode{}).Where("bin_number LIKE ?", code + "%").Find(&data); err != nil {
		log.Error("Get Bank Bin Code Database Error", err)
		return data, err
	}
	return data, nil
}
//取出信用卡刷卡及使用者
func GetGwCreditAndMemberCreditByOrderId(engine *database.MysqlSession, orderId, transType string) (entity.MemberGwCreditAuth, error) {
	var data entity.MemberGwCreditAuth
	if _, err := engine.Engine.Table(entity.GwCreditAuthData{}).Select("*").
		Join("LEFT", entity.MemberCardData{}, "gw_credit_auth_data.card_id = member_card_data.card_id").
		Where("gw_credit_auth_data.order_id = ?", orderId).And("gw_credit_auth_data.trans_type = ?", transType).Get(&data); err != nil {
		log.Error("Get Gw Credit Auth And Member Card Database Error", err)
		return data, err
	}
	return data, nil
}

//取GwCredit By OrderId
func GetGwCreditAllByStatus(engine *database.MysqlSession, now string) ([]entity.GwCreditAuthData, error) {
	var data []entity.GwCreditAuthData
	if err := engine.Engine.Table(entity.GwCreditAuthData{}).
		Select("*").Where("trans_status = ? AND create_time <= ?", Enum.CreditTransStatusInit, now).Find(&data); err != nil {
		log.Error("Get Gw Credit Database Error", err)
		return data, err
	}
	return data, nil
}
//取出Seller的次特店代碼
func GetSellerMerchantIdBySellerId(engine *database.MysqlSession, sellerId string) (entity.KgiSpecialStore, error) {
	var data entity.KgiSpecialStore
	if _, err := engine.Engine.Table(entity.KgiSpecialStore{}).Select("*").
		Where("user_id = ?", sellerId).Get(&data); err != nil {
		log.Error("Get Kgi Special Store Database Error", err)
		return data, err
	}
	return data, nil
}
//取出未使用的次特店代碼
func GetSellerMerchantIdUnused(engine *database.MysqlSession) (entity.KgiSpecialStore, error) {
	var data entity.KgiSpecialStore
	if _, err := engine.Engine.Table(entity.KgiSpecialStore{}).Select("*").
		Where("is_used = ?", false).Asc("id").Get(&data); err != nil {
		log.Error("Get Kgi Special Store Database Error", err)
		return data, err
	}
	return data, nil
}
//更新次特店代碼
func UpdateSellerMerchantId(engine *database.MysqlSession, data entity.KgiSpecialStore) error {
	data.UpdateTime = time.Now()
	if _, err := engine.Engine.Table(entity.KgiSpecialStore{}).ID(data.Id).AllCols().Update(&data); err != nil {
		log.Error("Get Kgi Special Store Database Error", err)
		return err
	}
	return nil
}
//新增次特店代碼
func InsertSellerMerchantId(engine *database.MysqlSession, data entity.KgiSpecialStore) error {
	data.CreateTime = time.Now()
	data.UpdateTime = time.Now()
	if _, err := engine.Session.Table(entity.KgiSpecialStore{}).Insert(&data); err != nil {
		log.Error("Insert Kgi Special Store Database Error", err)
		return err
	}
	return nil
}

