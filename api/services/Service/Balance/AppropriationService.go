package Balance

import (
	"api/services/Enum"
	"api/services/Service/OrderService"
	"api/services/Service/SysLog"
	"api/services/dao/Store"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"fmt"
	"github.com/spf13/viper"
)

//訂單撥付處理
func OrderAppropriation(engine *database.MysqlSession, OrderData entity.OrderData) error {
	if err := engine.Session.Begin(); err != nil {
		log.Error("Set Begin Error", err)
		return err
	}
	//取收銀機名稱
	storeData, err := Store.GetStoreDataByStoreId(engine, OrderData.StoreId)
	if err != nil {
		engine.Session.Rollback()
		log.Error("Get Store Error", err)
		return err
	}
	//剩餘金額存入賣家餘額
	comment := fmt.Sprintf("%s<br>%s", OrderData.OrderId, storeData.StoreName) //收銀機名稱 + 訂單號碼
	if err := Deposit(OrderData.SellerId, OrderData.OrderId, OrderData.TotalAmount, Enum.BalanceTypePayment, comment); err != nil {
		engine.Session.Rollback()
		log.Error("insert Balance Error", err)
		return err
	}
	//將保留款項扣除撥款金額
	if err := RetainWithdrawal(OrderData.SellerId, OrderData.OrderId, OrderData.TotalAmount, Enum.BalanceTypeWithdrawal, "撥款扣除"); err != nil {
		engine.Session.Rollback()
		log.Error("insert Balance Error", err)
		return err
	}
	//變更訂單撥款狀態
	if err := OrderService.ChangeOrderCaptureStatus(engine, OrderData, Enum.OrderCaptureSuccess); err != nil {
		engine.Session.Rollback()
		log.Error("Update order Error", err)
		return err
	}
	if err := SysLog.AppropriationSystemLog(OrderData.SellerId, int64(OrderData.TotalAmount)); err != nil {
		log.Error("Appropriation System Log Error", err)
	}
	if err := engine.Session.Commit(); err != nil {
		return err
	}
	return nil
}

//訂單提前扣除平台費用及物流費 fixme 無需配送未加
func OrderShipDeduction(engine *database.MysqlSession, OrderData *entity.OrderData) error {
	storeData, err := Store.GetStoreDataByStoreId(engine, OrderData.StoreId)
	if err != nil {
		log.Error("Get Store Error", err)
		return err
	}
	comment := fmt.Sprintf("%s<br>%s", OrderData.OrderId, storeData.StoreName) //收銀機名稱 + 訂單號碼
	fee := OrderData.PlatformShipFee + OrderData.PlatformTransFee + OrderData.PlatformInfoFee + OrderData.PlatformPayFee
	//賣家扣除 手續費及運費
	err = Withdrawal(OrderData.SellerId, OrderData.OrderId, fee, Enum.BalanceTypePlatform, comment)
	if err != nil {
		log.Error("Balance Withdrawal Error", err)
		return err
	}
	//手續費及運費 存入平台帳戶 fixme
	userId := viper.GetString("PLATFORM.USERID")
	err = Deposit(userId, OrderData.OrderId, fee, Enum.BalanceTypePlatform, comment)
	if err != nil {
		log.Error("Balance Deposit Error", err)
		return err
	}
	if err := OrderService.OrderOpenInvoice(engine, OrderData); err != nil {
		log.Error("Create Service Invoice Error", err)
		return err
	}
	return nil
}

//訂單提前扣除平台費用
func OrderPlatformDeduction(engine *database.MysqlSession, OrderData *entity.OrderData) error {
	storeData, err := Store.GetStoreDataByStoreId(engine, OrderData.StoreId)
	if err != nil {
		log.Error("Get Store Error", err)
		return err
	}
	comment := fmt.Sprintf("%s<br>%s", OrderData.OrderId, storeData.StoreName) //收銀機名稱 + 訂單號碼
	fee := OrderData.PlatformTransFee + OrderData.PlatformInfoFee + OrderData.PlatformPayFee
	//賣家扣除 手續費及運費
	err = Withdrawal(OrderData.SellerId, OrderData.OrderId, fee, Enum.BalanceTypePlatform, comment)
	if err != nil {
		log.Error("Balance Withdrawal Error", err)
		return err
	}
	//手續費及運費 存入平台帳戶
	userId := viper.GetString("PLATFORM.USERID")
	err = Deposit(userId, OrderData.OrderId, fee, Enum.BalanceTypePlatform, comment)
	if err != nil {
		log.Error("Balance Deposit Error", err)
		return err
	}
	if err := OrderService.OrderOpenInvoice(engine, OrderData); err != nil {
		log.Error("Create Service Invoice Error", err)
		return err
	}
	return nil
}
//訂單取消交易
func OrderCancelPaymentRefund(engine *database.MysqlSession, OrderData entity.OrderData) error {
	if OrderData.PayWay == Enum.Credit {
		//信用卡 請款前 取消受權 請款後 退款
		var data entity.CancelRequest
		data.OrderId = OrderData.OrderId
		if err := VoidProcess(engine, data); err != nil {
			return fmt.Errorf("退款失敗！")
		}
	}
	//轉帳及餘額 都退回至買家餘額
	if OrderData.PayWay == Enum.Transfer  || OrderData.PayWay == Enum.Balance {
		err := Deposit(OrderData.BuyerId, OrderData.OrderId, OrderData.TotalAmount, Enum.BalanceTypeRefund, OrderData.OrderId)
		if err != nil {
			log.Error("insert Balance Error", err)
		}
	}
	//將保留款項扣除撥款金額
	if err := RetainWithdrawal(OrderData.SellerId, OrderData.OrderId, OrderData.TotalAmount, Enum.BalanceTypeWithdrawal, "取消訂單"); err != nil {
		log.Error("insert Balance Error", err)
	}
	return nil
}