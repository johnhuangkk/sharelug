package model

import (
	"api/services/Enum"
	"api/services/Service/Balance"
	"api/services/Service/Mail"
	"api/services/Service/Notification"
	"api/services/Service/SysLog"
	"api/services/VO/Request"
	"api/services/VO/Response"
	"api/services/dao/Withdraw"
	"api/services/dao/member"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"api/services/util/tools"
	"fmt"
	"time"
)
//取銀行代碼
func GetBankCode(userData entity.MemberData) (Response.WithdrawResponse, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var resp Response.WithdrawResponse
	data, err := Withdraw.GetBankCode(engine)
	if err != nil {
		log.Error("Get Bank Code Error", err)
		return resp, err
	}
	for _, v := range data {
		var res Response.BackCodeList
		res.BackCode = v.BranchCode
		res.BankName = fmt.Sprintf("%s %s", v.BankCode, v.BankName)
		resp.BackCodeList = append(resp.BackCodeList, res)
	}
	account, err := member.GetMemberWithdrawListByUserId(engine, userData.Uid)
	if err != nil {
		log.Error("Get Bank Code Error", err)
		return resp, err
	}
	for _, v := range account {
		var rep Response.WithdrawAccount
		rep.BankAccount = v.Account
		rep.BankCode = v.BankCode
		rep.AccountName = v.BankName
		rep.Default = false
		if v.ActDefault == 1 {
			rep.Default = true
		}
		resp.WithdrawAccount = append(resp.WithdrawAccount, rep)
	}
	//增加EMAIL
	resp.IsVerifyEmail = false
	if len(userData.Email) != 0 {
		resp.IsVerifyEmail = true
	}
	balance := Balance.GetBalanceByUid(engine, userData.Uid)
	resp.Balance = int64(balance)
	return resp, nil
}
//處理刪除提領帳號
func HandleDeleteWithdraw(userData entity.MemberData, params Request.EditWithdrawRequest) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	data, err := member.GetMemberWithdrawByIdAndUserId(engine, userData.Uid, tools.StringToInt64(params.AccountId))
	if err != nil {
		log.Error("Get Withdraw Account Error", err)
		return err
	}
	data.Status = Enum.WithdrawStatusDelete
	if err := member.UpdateMemberWithdraw(engine, data.Id, data); err != nil {
		log.Error("Update Withdraw Account Error", err)
		return err
	}
	return nil
}
//變更預設提領帳號
func HandleChangeDefaultWithdraw(userData entity.MemberData, params Request.EditWithdrawRequest) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	if err := member.UpdateMemberWithdrawDeleteDefault(engine, userData.Uid); err != nil {
		return fmt.Errorf("系統錯誤")
	}
	data, err := member.GetMemberWithdrawByIdAndUserId(engine, userData.Uid, tools.StringToInt64(params.AccountId))
	if err != nil {
		log.Error("Get Withdraw Account Error", err)
		return err
	}
	data.ActDefault = 1
	if err := member.UpdateMemberWithdraw(engine, data.Id, data); err != nil {
		log.Error("Update Withdraw Account Error", err)
		return err
	}
	return nil
}
//提款處理
func HandleWithdraw(params Request.WithdrawRequest, userData entity.MemberData) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	if err := engine.Session.Begin(); err != nil {
		log.Error("engine Begin Error", err)
		return fmt.Errorf("1001001")
	}
	if userData.VerifyIdentity != 1 {
		return fmt.Errorf("1003012")
	}
	if userData.VerifyBusiness != 1 {
		return fmt.Errorf("1003013")
	}
	//取銀行資料
	bank, err := Withdraw.GetBankInfoByBranchCode(engine, params.BankCode)
	if err != nil {
		log.Error("Get Bank BranchCode Error", err)
		return fmt.Errorf("1001001")
	}
	//產生提領單
	data := generatorWithdraw(params, userData, bank)
	//取目前的餘額
	balance := Balance.GetBalanceByUid(engine, userData.Uid)
	//檢查餘額 金額外加10元
	if int64(balance) < data.WithdrawAmt {
		return fmt.Errorf("1012001")
	}
	//扣除餘額
	if err := Balance.WithdrawDeduction(data); err != nil {
		log.Error("Withdraw Deduction Error", err)
		engine.Session.Rollback()
		return fmt.Errorf("1001001")
	}
	//寫入提領單
	if err := Withdraw.InsertWithdrawData(engine, data); err != nil {
		log.Error("Insert Withdraw Error", err)
		engine.Session.Rollback()
		return fmt.Errorf("1001001")
	}
	//將帳戶記錄在常用名單中
	if err := processMemberWithdraw(engine, data); err != nil {
		log.Error("Process Member Used Withdraw Error", err)
		engine.Session.Rollback()
		return fmt.Errorf("1001001")
	}
	if err := Notification.SendApplyWithdrawMessage(engine, data.UserId, data.WithdrawAmt); err != nil {
		engine.Session.Rollback()
		log.Error("Send apply Withdraw Error", err)
		return fmt.Errorf("1001001")
	}
	//發信
	if len(params.Email) != 0 {
		if err := Mail.EmailVerify(engine, userData.Uid, params.Email); err != nil {
			log.Error("Send Mail Verify Error", err)
			engine.Session.Rollback()
			return fmt.Errorf("1001001")
		}
	}
	if err := engine.Session.Commit(); err != nil {
		log.Error("Commit error", err)
		return fmt.Errorf("1001001")
	}
	if err := SysLog.WithdrawSystemLog(data.UserId, data.BankName, data.WithdrawAmt); err != nil {
		log.Error("System Log Error", err)
	}
	return nil
}
//產生提領單
func generatorWithdraw(params Request.WithdrawRequest, userData entity.MemberData, bank entity.BankCodeData) entity.WithdrawData {
	var data entity.WithdrawData
	data.WithdrawId = tools.GeneratorWithdrawId()
	data.UserId = userData.Uid
	data.BankCode = bank.BranchCode
	data.BankName = bank.BankName
	data.BankAccount = params.BankAccount
	data.TransId = tools.GeneratorTransId()
	data.WithdrawAmt = tools.StringToInt64(params.WithdrawAmt)
	data.WithdrawFee = 0
	data.WithdrawTime = time.Now()
	data.WithdrawStatus = Enum.WithdrawStatusWait
	data.ResponseStatus = Enum.WithdrawStatusInit
	return data
}
//將提領帳戶記錄在常用名單中
func processMemberWithdraw(engine *database.MysqlSession, withdraw entity.WithdrawData) error {
	//取出銀行帳號
	data, err := member.GetMemberWithdrawByAccount(engine, withdraw.BankAccount)
	if err != nil {
		return fmt.Errorf("系統錯誤")
	}
	if len(data.Account) == 0 {
		if err := createMemberWithdraw(engine, withdraw); err != nil {
			return fmt.Errorf("系統錯誤")
		}
		if err := SysLog.SetBankSystemLog(withdraw.UserId, withdraw.BankName, withdraw.BankAccount); err != nil {
			log.Error("System Log Error", err)
		}
	}
	if err := member.UpdateMemberWithdrawDeleteDefault(engine, withdraw.UserId); err != nil {
		return fmt.Errorf("系統錯誤")
	}
	if err := member.UpdateMemberWithdrawSetDefault(engine, withdraw.UserId, withdraw.BankAccount); err != nil {
		return fmt.Errorf("系統錯誤")
	}
	return nil
}
//寫入提領帳戶常用名單
func createMemberWithdraw(engine *database.MysqlSession, withdraw entity.WithdrawData) error {
	var data entity.MemberWithdrawData
	data.UserId = withdraw.UserId
	data.BankName = withdraw.BankName
	data.BankCode = withdraw.BankCode
	data.Account = withdraw.BankAccount
	data.Status = Enum.WithdrawStatusSuccess
	err := member.InsertMemberWithdraw(engine, data)
	if err != nil {
		return fmt.Errorf("系統錯誤")
	}
	return nil
}
//提領退回
func WithdrawFailedRefund(engine *database.MysqlSession, data entity.WithdrawData) error {
	comment := fmt.Sprintf("提領日期:%s", data.WithdrawTime.Format("2006/01/02")) //收銀機名稱 + 訂單號碼
	amount := float64(data.WithdrawAmt + data.WithdrawFee)
	if err := Balance.Deposit(data.UserId, data.WithdrawId, amount, Enum.BalanceTypeWdFailed, comment); err != nil {
		log.Error("insert Balance Error", err)
		return err
	}
	if err := Notification.SendWithdrawFailedRefund(engine, data.UserId, data.WithdrawTime.Format("2006/01/02")); err != nil {
		log.Error("Send Withdraw Failed Refund Error", err)
		return err
	}
	return nil
}
