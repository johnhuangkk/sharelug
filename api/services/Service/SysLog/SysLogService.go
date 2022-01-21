package SysLog

import (
	"api/config/middleware"
	"api/services/Enum"
	"api/services/dao/SysLogDao"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"api/services/util/tools"
	"fmt"
	"time"
)

//開設賣場：(賣場名稱)
func NewStoreSystemLog(userId, storeName string) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	content := fmt.Sprintf("開設賣場：%s", storeName)
	if err := generateSystemLog(engine, userId, Enum.ActivityStore, content); err != nil {
		log.Error("Generate System Log Error", err)
		return err
	}
	return nil
}
//指派管理帳號：(賣場名稱 / 隱碼之管理帳號之email)
func AssignManagerSystemLog(userId, storeName, email string) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	content := fmt.Sprintf("指派管理帳號：%s/%s", storeName, tools.MaskerEMail(email))
	if err := generateSystemLog(engine, userId, Enum.ActivityMgt, content); err != nil {
		log.Error("Generate System Log Error", err)
		return err
	}
	return nil
}
//款項撥付：NT$(金額)
func AppropriationSystemLog(userId string, amount int64) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	content := fmt.Sprintf("款項撥付：NT$ %v", tools.FormatFinancialString(tools.IntToString(int(amount))))
	if err := generateSystemLog(engine, userId, Enum.ActivityAppn, content); err != nil {
		log.Error("Generate System Log Error", err)
		return err
	}
	return nil
}
//提領：NT$(金額 / 銀行名稱)
func WithdrawSystemLog(userId, bankName string, amount int64) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	content := fmt.Sprintf("提領：NT$ %v/%s", tools.FormatFinancialString(tools.IntToString(int(amount))), bankName)
	if err := generateSystemLog(engine, userId, Enum.ActivityWithdraw, content); err != nil {
		log.Error("Generate System Log Error", err)
		return err
	}
	return nil
}
//設定銀行帳號：(銀行 / 隱碼之帳號)
func SetBankSystemLog(userId, bankName, account string) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	content := fmt.Sprintf("設定銀行帳號：%s/%s", bankName, tools.MaskerBankAccount(account))
	if err := generateSystemLog(engine, userId, Enum.ActivityWithdraw, content); err != nil {
		log.Error("Generate System Log Error", err)
		return err
	}
	return nil
}
//設定信用卡：(銀行 / 隱碼之信用卡卡號)
func SetCreditSystemLog(userId, bankName, Number string) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	content := fmt.Sprintf("設定信用卡：%s/************%s", bankName, Number)
	if err := generateSystemLog(engine, userId, Enum.ActivitySetCredit, content); err != nil {
		log.Error("Generate System Log Error", err)
		return err
	}
	return nil
}
//關閉賣場：(賣場名稱)
func ShutdownStoreSystemLog(userId, storeName string) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	content := fmt.Sprintf("關閉賣場：%s", storeName)
	if err := generateSystemLog(engine, userId, Enum.ActivityShutdown, content); err != nil {
		log.Error("Generate System Log Error", err)
		return err
	}
	return nil
}
//取消管理帳號：(賣場名稱 / 隱碼之管理帳號email)
func CancelManagerSystemLog(userId, storeName, email string) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	content := fmt.Sprintf("取消管理帳號：%s/%s", storeName, email)
	if err := generateSystemLog(engine, userId, Enum.ActivityCancelManager, content); err != nil {
		log.Error("Generate System Log Error", err)
		return err
	}
	return nil
}
//變更手機號碼：(隱碼之手機號碼)
func ChangePhoneSystemLog(userId, phone string) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	content := fmt.Sprintf("變更手機號碼：%s", phone)
	if err := generateSystemLog(engine, userId, Enum.ActivityChangePhone, content); err != nil {
		log.Error("Generate System Log Error", err)
		return err
	}
	return nil
}
//變更email：(隱碼之email)
func ChangeEmailSystemLog(userId, email string) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	log.Debug("ssss", email)
	content := fmt.Sprintf("變更EMail：%s", tools.MaskerEMail(email))
	if err := generateSystemLog(engine, userId, Enum.ActivityCancelEmail, content); err != nil {
		log.Error("Generate System Log Error", err)
		return err
	}
	return nil
}
//產生使用者操做記錄
func generateSystemLog(engine *database.MysqlSession, UserId, Action, content string) error {
	var data entity.SystemLog
	data.UserId = UserId
	data.Action = Action
	data.Content = content
	data.LoginIp = middleware.GetClientIP()
	data.CreateTime = time.Now()
	if err := SysLogDao.InsertSysLog(engine, data); err != nil {
		log.Error("Insert SysLog Database Error", err)
		return err
	}
	return nil
}