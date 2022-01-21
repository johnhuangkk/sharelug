package Notification

import (
	"api/services/Enum"
	"api/services/Service/Mail"
	"api/services/Service/Sms"
	"api/services/Service/Upgrade"
	"api/services/dao/Orders"
	"api/services/dao/Store"
	"api/services/dao/member"
	"api/services/dao/product"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"fmt"
	"time"
)

//升級服務方案
func SendUpgradeApplyMessage(engine *database.MysqlSession, StoreId, SellerId, OrderId string) error {
	//取store資料
	_, UserData, orderData, err := GetOrderInfo(engine, SellerId, StoreId, OrderId)
	if err != nil {
		return err
	}
	now := time.Now().Format("2006/01/02")
	//發送簡訊
	err = Sms.SendUpgradeApplySms(UserData, orderData)
	if err != nil {
		return err
	}
	//發送MAIL
	if len(UserData.Email) != 0 {
		err := Mail.SendUpgradeApplyMail(UserData, orderData)
		if err != nil {
			return err
		}
	}
	//發送站內訊息
	data, _ := product.GetUpgradeProductDataByLevel(engine, orderData.UpgradeLevel)
	log.Debug("data", data)
	msg := fmt.Sprintf("收銀機帳號升級：已成功購買「%s」，方案費用為NT：%v元，將於%s日生效。",  data.ProductName, orderData.Amount, now)
	err = SendSystemNotify(engine, UserData.Uid, msg, Enum.NotifyMsgTypePlaPlatform, "")
	if err != nil {
		return err
	}
	return nil
}

//服務方案中止
func SendUpgradeStopMessage(engine *database.MysqlSession, sellerData entity.MemberData) error {
	//取出賣場資料
	managerData, err := Store.GetStoreAllManagerListBySellerId(engine, sellerData.Uid)
	if err != nil {
		return err
	}
	//發送簡訊
	if err := Sms.SendUpgradeStopSms(sellerData); err != nil {
		return err
	}
	//發送MAIL
	if len(sellerData.Email) != 0 {
		err := Mail.SendUpgradeStopMail(sellerData)
		if err != nil {
			return err
		}
	}
	//發送站內訊息
	msg := fmt.Sprintf("申請提前終止「%s」，現有收銀機方案將於 %s日到期，在此之前你可以繼續使用本收銀機服務。", Upgrade.GetChangeProduct(sellerData.UpgradeLevel), sellerData.UpgradeExpire.Format("2006/01/02"))
	if err := SendSystemNotify(engine, sellerData.Uid, msg, Enum.NotifyMsgTypePlaPlatform, ""); err != nil {
		return err
	}
	for _, v := range managerData {
		UserData, err := member.GetMemberDataByUid(engine, v.StoreRank.UserId)
		if err != nil {
			return err
		}
		msg := fmt.Sprintf("申請終止「%s」，現有賣場方案將於 %s日到期，在此之前你可以繼續使用本賣場服務。", Upgrade.GetChangeProduct(sellerData.UpgradeLevel), sellerData.UpgradeExpire.Format("2006/01/02"))
		if err := SendSystemNotify(engine, UserData.Uid, msg, Enum.NotifyMsgTypePlaPlatform, ""); err != nil {
			return err
		}
	}
	return nil
}

//新增管理員權限
func SendAddManagerMessage(engine *database.MysqlSession, StoreId, SellerId string) error {
	StoreData, err := Store.GetStoreDataByStoreId(engine, StoreId)
	if err != nil {
		return err
	}
	UserData, err := member.GetMemberDataByUid(engine, SellerId)
	if err != nil {
		return err
	}
	err = Sms.SendAddManagerSms(UserData)
	if err != nil {
		return err
	}
	if len(UserData.Email) != 0 {
		err := Mail.SendAddManagerMail(UserData, StoreData)
		if err != nil {
			return err
		}
	}
	//發送站內訊息
	msg := fmt.Sprintf("%s 新增一名管理員權限，請聯繫管理員儘速完成登入驗證流程。", time.Now().Format("2006/01/02"))
	err = SendSystemNotify(engine, UserData.Uid, msg, Enum.NotifyMsgTypePlaPlatform, "")
	if err != nil {
		return err
	}
	return nil
}

//管理員完成驗證
func SendManagerVerifyCompleteMessage(engine *database.MysqlSession, StoreId, ManagerId string) error {
	StoreData, err := Store.GetStoreDataByStoreId(engine, StoreId)
	if err != nil {
		return err
	}
	//取得UserData
	masterData, err := member.GetMemberDataByUid(engine, StoreData.SellerId)
	if err != nil {
		return err
	}
	slaveData, err := member.GetMemberDataByUid(engine, ManagerId)
	if err != nil {
		return err
	}
	//發送站內訊息
	msg := fmt.Sprintf("%s 收銀機：%s 已完成管理員權限驗證。", StoreData.StoreName, slaveData.Username)
	err = SendSystemNotify(engine, masterData.Uid, msg, Enum.NotifyMsgTypePlaPlatform, "")
	if err != nil {
		return err
	}
	msg = fmt.Sprintf("你已加入[%s]，可以共同協作管理收銀機。", StoreData.StoreName)
	err = SendSystemNotify(engine, slaveData.Uid, msg, Enum.NotifyMsgTypePlaPlatform, "")
	if err != nil {
		return err
	}
	return nil
}

//刪除管理員權限
func SendDeleteManagerMessage(engine *database.MysqlSession, StoreId, ManagerId string) error {
	now := time.Now().Format("2006/01/02 15:04")
	StoreData, err := Store.GetStoreDataByStoreId(engine, StoreId)
	if err != nil {
		return err
	}
	masterData, err := member.GetMemberDataByUid(engine, StoreData.SellerId)
	if err != nil {
		return err
	}
	slaveData, err := member.GetMemberDataByUid(engine, ManagerId)
	if err != nil {
		return err
	}
	err = Sms.SendDeleteManagerSms(masterData)
	if err != nil {
		return err
	}
	if len(masterData.Email) != 0 {
		err := Mail.SendDeleteManagerMail(masterData, StoreData)
		if err != nil {
			return err
		}
	}
	//發送站內訊息
	msg := fmt.Sprintf("%s [ %s ]已刪除[ %s ]的管理權限。", now, StoreData.StoreName, slaveData.Username)
	err = SendSystemNotify(engine, masterData.Uid, msg, Enum.NotifyMsgTypePlaPlatform, "")
	if err != nil {
		return err
	}
	msg = fmt.Sprintf("%s 已移除Check'Ne[%s]管理員權限。", now, StoreData.StoreName)
	err = SendSystemNotify(engine, slaveData.Uid, msg, Enum.NotifyMsgTypePlaPlatform, "")
	if err != nil {
		return err
	}
	return nil
}

//繳費通知
func SendUpgradePaymentMessage(engine *database.MysqlSession, StoreId, SellerId, OrderId string) error {
	StoreData, UserData, OrderData, err := GetBillInfo(engine, SellerId, StoreId, OrderId)
	if err != nil {
		return err
	}
	err = Sms.SendUpgradePaymentSms(UserData)
	if err != nil {
		return err
	}
	if len(UserData.Email) != 0 {
		err := Mail.SendUpgradePaymentMail(UserData, StoreData, OrderData)
		if err != nil {
			return err
		}
	}
	//發送站內訊息
	data, _ := product.GetUpgradeProductDataByLevel(engine, OrderData.BillingLevel)
	sum, _ := Orders.SumB2cUnpaidBillByUserId(engine, OrderData.UserId)
	msg := fmt.Sprintf("方案為：[%s]，本期累積帳單總金額：NT＄%v，收銀機有效時間：%s。", data.ProductName, sum, OrderData.Expiration.Format("2006/01/02"))
	err = SendSystemNotify(engine, UserData.Uid, msg, Enum.NotifyMsgTypePlaPlatform, "")
	if err != nil {
		return err
	}
	return nil
}

//申請提領
func SendApplyWithdrawMessage(engine *database.MysqlSession, SellerId string, amount int64) error {
	UserData, err := member.GetMemberDataByUid(engine, SellerId)
	if err != nil {
		return err
	}
	err = Sms.SendApplyWithdrawSms(UserData)
	if err != nil {
		return err
	}
	if len(UserData.Email) != 0 {
		err := Mail.SendApplyWithdrawMail(UserData, amount)
		if err != nil {
			return err
		}
	}
	msg := fmt.Sprintf("新增一筆提領款項，提領金額NT$ %v 元。", amount)
	err = SendSystemNotify(engine, UserData.Uid, msg, Enum.NotifyMsgTypePlaPlatform, "")
	if err != nil {
		return err
	}
	return nil
}

//暫停收銀機
func SendSuspendStoreMessage(engine *database.MysqlSession, StoreId string) error {
	StoreData, err := Store.GetStoreDataByStoreId(engine, StoreId)
	if err != nil {
		return err
	}
	masterData, err := member.GetMemberDataByUid(engine, StoreData.SellerId)
	if err != nil {
		return err
	}
	err = Sms.SendSuspendStoreSms(masterData, StoreData)
	if err != nil {
		return err
	}
	if len(masterData.Email) != 0 {
		err := Mail.SendSuspendStoreMail(masterData, StoreData)
		if err != nil {
			return err
		}
	}
	msg := fmt.Sprintf("[%s]已申請暫停對外營業，商品下架、結帳連結無法使用。\n如要申請恢復收銀機經營，請至「我的收銀機」>開啟「收銀機狀態」；並更新商品列表中商品狀態：上架商品。", StoreData.StoreName)
	err = SendSystemNotify(engine, masterData.Uid, msg, Enum.NotifyMsgTypePlaPlatform, "")
	if err != nil {
		return err
	}
	manager, err := Store.GetStoreManagerByStoreId(engine, StoreData.StoreId)
	if err != nil {
		return err
	}
	for _, v := range manager {
		UserData, err := member.GetMemberDataByUid(engine, v.UserId)
		if err != nil {
			return err
		}
		msg := fmt.Sprintf("[%s]已申請暫停對外營業，商品下架、結帳連結無法使用。", StoreData.StoreName)
		err = SendSystemNotify(engine, UserData.Uid, msg, Enum.NotifyMsgTypePlaPlatform, "")
		if err != nil {
			return err
		}
	}
	return nil
}

func SendOrderInvoiceMessage(engine *database.MysqlSession, user entity.MemberData, data entity.InvoiceData) error {
	if len(user.Email) != 0 {
		if err := Mail.SendOrderInvoiceEmail(user, data); err != nil {
			return err
		}
	} else {
		if err := Sms.SendOrderInvoiceSms(user); err != nil {
			return err
		}
	}
	msg := fmt.Sprintf("平台服務費電子發票已開立如下：\n發票號碼：%s\n發票金額：NT$ %v元\n請至發票資訊查詢發票明細，也可以到財政部電子發票整合服務平台查詢發票明細資料。",
		fmt.Sprintf("%s%s", data.InvoiceTrack, data.InvoiceNumber), data.Amount)
	if err := SendSystemNotify(engine, user.Uid, msg, Enum.NotifyMsgTypePlaPlatform, ""); err != nil {
		return err
	}
	return nil
}

func SendServiceInvoiceMessage(engine *database.MysqlSession, user entity.MemberData, data entity.InvoiceData) error {
	if len(user.Email) != 0 {
		if err := Mail.SendServiceInvoiceEmail(user, data); err != nil {
			return err
		}
	} else {
		if err := Sms.SendServiceInvoiceSms(user); err != nil {
			return err
		}
	}
	msg := fmt.Sprintf("加值服務費 電子發票已開立如下：\n發票號碼：%s\n發票金額：NT$ %v元\n請至發票資訊查詢發票明細，也可以到財政部電子發票整合服務平台查詢發票明細資料。",
		fmt.Sprintf("%s%s", data.InvoiceTrack, data.InvoiceNumber), data.Amount)
	if err := SendSystemNotify(engine, user.Uid, msg, Enum.NotifyMsgTypePlaPlatform, ""); err != nil {
		return err
	}
	return nil
}


//平台暫停收銀機使用
func SendPlatformSuspendStoreMessage(engine *database.MysqlSession, StoreId string) error {
	StoreData, err := Store.GetStoreDataByStoreId(engine, StoreId)
	if err != nil {
		return err
	}
	masterData, err := member.GetMemberDataByUid(engine, StoreData.SellerId)
	if err != nil {
		return err
	}
	err = Sms.SendPlatformSuspendStoreSms(masterData, StoreData)
	if err != nil {
		return err
	}
	if len(masterData.Email) != 0 {
		err := Mail.SendPlatformSuspendStoreMail(masterData, StoreData)
		if err != nil {
			return err
		}
	}
	msg := fmt.Sprintf("[%s]因違反平台管理規則，目前強制暫停中，無法進行商品上架及結帳成交交易。", StoreData.StoreName)
	err = SendSystemNotify(engine, masterData.Uid, msg, Enum.NotifyMsgTypePlaPlatform, "")
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	msg = fmt.Sprintf("[%s]因違反平台管理規則，目前強制暫停中，無法進行商品上架及結帳成交交易。", StoreData.StoreName)
	err = SendSystemNotify(engine, StoreData.StoreId, msg, Enum.NotifyMsgTypePlaPlatform, "")
	if err != nil {
		return err
	}
	return nil
}

//發送提領失敗訊息
func SendWithdrawFailedRefund(engine *database.MysqlSession, UserId, WithdrawTime string) error {
	msg := fmt.Sprintf("你於 %s 辦理之提領，無法完成匯款，請確認指定銀行帳戶是否與身份資料相符後，再次辦理提領。", WithdrawTime)
	if err := SendSystemNotify(engine, UserId, msg, Enum.NotifyMsgTypePlaPlatform, ""); err != nil {
		return err
	}
	return nil
}

//發送電子郵件驗證完成
func SendEmailVerifySuccess(engine *database.MysqlSession, UserId string) error {
	msg := fmt.Sprintf("你已完成電子郵件驗證，如未有更新紀錄，請儘速與我們聯繫，謝謝。")
	if err := SendSystemNotify(engine, UserId, msg, Enum.NotifyMsgTypePlaPlatform, ""); err != nil {
		return err
	}
	return nil
}


func GetOrderInfo(engine *database.MysqlSession, SellerId, StoreId, OrderId string) (entity.StoreDataResp, entity.MemberData, entity.B2cOrderData, error) {
	var StoreData entity.StoreDataResp
	var UserData entity.MemberData
	var orderData entity.B2cOrderData

	StoreData, err := Store.GetStoreDataByUserIdAndStoreId(engine, SellerId, StoreId)
	if err != nil {
		return StoreData, UserData, orderData, err
	}
	//取得UserData
	UserData, err = member.GetMemberDataByUid(engine, StoreData.SellerId)
	if err != nil {
		return StoreData, UserData, orderData, err
	}
	//取訂單內容
	orderData, err = Orders.GetB2cOrderByOrderId(engine, OrderId)
	if err != nil {
		return StoreData, UserData, orderData, err
	}
	return StoreData, UserData, orderData, nil
}

func GetBillInfo(engine *database.MysqlSession, SellerId, StoreId, OrderId string) (entity.StoreDataResp, entity.MemberData, entity.B2cBillingData, error) {
	var StoreData entity.StoreDataResp
	var UserData entity.MemberData
	var billData entity.B2cBillingData

	StoreData, err := Store.GetStoreDataByUserIdAndStoreId(engine, SellerId, StoreId)
	if err != nil {
		return StoreData, UserData, billData, err
	}
	//取得UserData
	UserData, err = member.GetMemberDataByUid(engine, StoreData.SellerId)
	if err != nil {
		return StoreData, UserData, billData, err
	}
	//取訂單內容
	billData, err = Orders.GetB2cBillByOrderId(engine, OrderId)
	if err != nil {
		return StoreData, UserData, billData, err
	}
	return StoreData, UserData, billData, nil
}
