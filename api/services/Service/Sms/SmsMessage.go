package Sms

import (
	"api/services/Service/Upgrade"
	"api/services/entity"
	"fmt"
)

//升級服務方案
func SendUpgradeApplySms(UserData entity.MemberData, OrderData entity.B2cOrderData) error {
	content := fmt.Sprintf("感謝你在Check'Ne網站上完成收銀機管理NT:%v升級服務。如你未曾購買升級方案請儘速與我們聯繫。", OrderData.Amount)
	if err := PushMessageSms(UserData.Mphone, content); err != nil {
		return err
	}
	return nil
}
//服務方案中止
func SendUpgradeStopSms(UserData entity.MemberData) error {
	content := fmt.Sprintf("你好：已收到你申請提前終止「%s」。如你未曾申請提前終止請儘速與我們聯繫。", Upgrade.GetChangeProduct(UserData.UpgradeLevel))
	if err := PushMessageSms(UserData.Mphone, content); err != nil {
		return err
	}
	return nil
}
//新增管理員權限
func SendAddManagerSms(UserData entity.MemberData) error {
	content := fmt.Sprintf("你好：\n已收到Check'Ne收銀機：新增管理員權限，如你未曾申請新增管理員權限請儘速與我們聯繫。")
	if err := PushMessageSms(UserData.Mphone, content); err != nil {
		return err
	}
	return nil
}
//刪除管理員權限
func SendDeleteManagerSms(UserData entity.MemberData) error {
	content := fmt.Sprintf("你好：\n已收到Check'Ne收銀機：刪除管理員權限，如你未曾申請刪除管理員權限請儘速與我們聯繫。")
	if err := PushMessageSms(UserData.Mphone, content); err != nil {
		return err
	}
	return nil
}
//繳費通知
func SendUpgradePaymentSms(UserData entity.MemberData) error {
	content := fmt.Sprintf("Check'Ne收銀機新增一筆收銀機帳單，請登入checkne進行繳款，逾期繳費將會影響收銀機權益。如你未曾申請收銀機服務，請儘速與我們聯繫。")
	if err := PushMessageSms(UserData.Mphone, content); err != nil {
		return err
	}
	return nil
}
//申請提領
func SendApplyWithdrawSms(UserData entity.MemberData) error {
	content := fmt.Sprintf("Check'Ne新增一筆款項提領，如你未曾申請款項提領，請儘速與我們聯繫。")
	if err := PushMessageSms(UserData.Mphone, content); err != nil {
		return err
	}
	return nil
}
//暫停收銀機
func SendSuspendStoreSms(UserData entity.MemberData, StoreData entity.StoreData) error {
	content := fmt.Sprintf("Check'Ne %s已暫停對外營業，商品下架、結帳連結無法使用。如你未曾申請「暫停收銀機」，請儘速與我們聯繫。", StoreData.StoreName)
	if err := PushMessageSms(UserData.Mphone, content); err != nil {
		return err
	}
	return nil
}
//平台暫停收銀機使用
func SendPlatformSuspendStoreSms(UserData entity.MemberData, StoreData entity.StoreData) error {
	content := fmt.Sprintf("Check'Ne %s因違反平台管理規則，目前強制暫停中，無法進行商品上架及結帳成交交易，如需重新開啟收銀機管理功能，請儘速與我們聯繫。", StoreData.StoreName)
	if err := PushMessageSms(UserData.Mphone, content); err != nil {
		return err
	}
	return nil
}

func SendSuccessfullySellerOrderedSms(UserData entity.MemberData, StoreData entity.StoreData) error {
	content := fmt.Sprintf("Check'Ne %s有一筆新訂單，請登入Check'Ne安排出貨。", StoreData.StoreName)
	if err := PushMessageSms(UserData.Mphone, content); err != nil {
		return err
	}
	return nil
}

func SendSuccessfullyBuyerOrderSms(UserData entity.MemberData, StoreData entity.StoreData) error {
	content := fmt.Sprintf("已收到你在 %s完成的訂單，登入Check'Ne追蹤出貨進度。若你未曾訂購，請與我們聯絡。", StoreData.StoreName)
	if err := PushMessageSms(UserData.Mphone, content); err != nil {
		return err
	}
	return nil
}

func SendAtmSuccessfullyBuyerOrderedSms(UserData entity.MemberData, StoreData entity.StoreData, data entity.TransferData) error {
	date := data.ExpireDate.Format("01/02") + " 23:59"
	bank := fmt.Sprintf("%s%s", data.BankCode, data.BankName)
	content := fmt.Sprintf("有Check'Ne訂購請於%s前繳款%s帳號%s金額%v，未訂購勿轉帳", date, bank, data.BankAccount, data.Amount)
	//content := fmt.Sprintf("已收到你在%s的訂購，請於%s前繳款：%s，轉帳帳號%s，金額%d元，如未訂購請勿轉帳並通知CheckNe。", StoreData.StoreName, date, bank, data.BankAccount, data.Amount)
	if err := PushMessageSms(UserData.Mphone, content); err != nil {
		return err
	}
	return nil
}


func SendWaitShipSms(UserData entity.MemberData) error {
	content := fmt.Sprintf("你的Check'Ne收銀機有新進一筆訂單待出貨，請儘速安排出貨。")
	if err := PushMessageSms(UserData.Mphone, content); err != nil {
		return err
	}
	return nil
}

func SendNotShippedSms(UserData entity.MemberData, count int) error {
	content := fmt.Sprintf("你的Check'Ne收銀機有%v筆訂單待出貨已超過2天，請儘速安排出貨。", count)
	if err := PushMessageSms(UserData.Mphone, content); err != nil {
		return err
	}
	return nil
}

func SendOrderCancelSellerSms(UserData entity.MemberData, StoreData entity.StoreData) error {
	content := fmt.Sprintf("%s的一筆訂單已取消，你可至Check'Ne查詢訂單狀態。", StoreData.StoreName)
	if err := PushMessageSms(UserData.Mphone, content); err != nil {
		return err
	}
	return nil
}

func SendOrderShipOverdueSms(UserData entity.MemberData, OrderData entity.OrderData) error {
	content := fmt.Sprintf("Check'Ne訂單編號：%s  逾期未寄，出貨單號：%s 已失效無法寄送，請儘早辦理退款。", OrderData.OrderId, OrderData.ShipNumber)
	if err := PushMessageSms(UserData.Mphone, content); err != nil {
		return err
	}
	return nil
}

func SendOrderCancelBuyerSms(UserData entity.MemberData, StoreData entity.StoreData) error {
	content := fmt.Sprintf("你在%s的訂單已取消，你可至Check'Ne查詢訂單狀態並與賣家聯絡。", StoreData.StoreName)
	if err := PushMessageSms(UserData.Mphone, content); err != nil {
		return err
	}
	return nil
}

func SendForSellerShipClosedShopSms(UserData entity.MemberData, OrderData entity.OrderData) error {
	content := fmt.Sprintf("訂單編號：%s 配送超商因故無法提供取件服務，請通知買家在2日內重新選擇配送門市，以免配送時間過期。", OrderData.OrderId)
	if err := PushMessageSms(UserData.Mphone, content); err != nil {
		return err
	}
	return nil
}

func SendForBuyerShipClosedShopSms(UserData entity.MemberData, OrderData entity.OrderData) error {
	content := fmt.Sprintf("訂單編號：%s 配送超商因故無法提供取件服務，請在2日內至Check'Ne重新選擇配送門市，以免配送時間過期。", OrderData.OrderId)
	if err := PushMessageSms(UserData.Mphone, content); err != nil {
		return err
	}
	return nil
}

func SendReturnShipClosedShopSms(UserData entity.MemberData, OrderData entity.OrderData) error {
	content := fmt.Sprintf("Check'Ne訂編：%s 包裹逾期未取，原交寄門市因故無法提供退回商品之取件服務，請在2日內到Check'Ne訂單管理中重新選擇門市。", OrderData.OrderId)
	if err := PushMessageSms(UserData.Mphone, content); err != nil {
		return err
	}
	return nil
}

func SendBuyerRefundApplySms(UserData entity.MemberData, StoreData entity.StoreData, refundData entity.OrderRefundData) error {
	content := fmt.Sprintf("%s有一筆訂單退款，金額NT$%v元將退到Check'Ne餘額中，Check'Ne查看退款進度。",
		StoreData.StoreName, int64(refundData.Total))
	if err := PushMessageSms(UserData.Mphone, content); err != nil {
		return err
	}
	return nil
}

func SendOrderCustomerSms(UserData entity.MemberData, StoreData entity.StoreData, count int64) error {
	content := fmt.Sprintf("%s新增%v筆訂單客服，請儘速登入checkne.com回覆。", StoreData.StoreName, count)
	if err := PushMessageSms(UserData.Mphone, content); err != nil {
		return err
	}
	return nil
}

func SendReplyOrderCustomerSms(UserData entity.MemberData, StoreData entity.StoreData, OrderData entity.OrderData) error {
	content := fmt.Sprintf("%s已針對訂編：%s 中提出問題，請登入Check'Ne查看。",
		StoreData.StoreName, OrderData.OrderId)
	if err := PushMessageSms(UserData.Mphone, content); err != nil {
		return err
	}
	return nil
}

func SendShippedMessageSms(UserData entity.MemberData, StoreData entity.StoreData, OrderData entity.OrderData) error {
	content := fmt.Sprintf("Check'Ne %s 訂單編號：%s 已出貨，追蹤配送進度請到checkne.com查詢。", StoreData.StoreName, OrderData.ShipNumber)
	if err := PushMessageSms(UserData.Mphone, content); err != nil {
		return err
	}
	return nil
}

func SendBuyerBillOrderSuccessSms(UserData entity.MemberData, StoreData entity.StoreData) error {
	content := fmt.Sprintf("%s已接受你的訂購單，登入Check'Ne追蹤出貨進度。若你未曾訂購，請與我們聯絡。",
		StoreData.StoreName)
	if err := PushMessageSms(UserData.Mphone, content); err != nil {
		return err
	}
	return nil
}

func SendSellerBillOrderSuccessSms(UserData entity.MemberData, StoreData entity.StoreData) error {
	content := fmt.Sprintf("%s已接受一筆Check'Ne新訂單，請登入Check'Ne安排出貨。",
		StoreData.StoreName)
	if err := PushMessageSms(UserData.Mphone, content); err != nil {
		return err
	}
	return nil
}

func SendBuyerAtmBillOrderSms(UserData entity.MemberData, StoreData entity.StoreData, data entity.TransferData) error {
	date := data.ExpireDate.Format("01/02") + " 23:59"
	bank := fmt.Sprintf("%s%s", data.BankCode, data.BankName)
	content := fmt.Sprintf("Check'Ne收到訂單請於%s前繳款：%s，帳號%s，NT$%v，未訂購請勿轉帳", date, bank, data.BankAccount, data.Amount)
	//content := fmt.Sprintf("%s已接受你的訂購單，請於%s前繳款：%s，轉帳帳號%s，金額%v元，如未訂購請勿轉帳並通知CheckNe。", StoreData.StoreName, date, bank, data.BankAccount, data.Amount)
	if err := PushMessageSms(UserData.Mphone, content); err != nil {
		return err
	}
	return nil
}

func SendOrderInvoiceSms(UserData entity.MemberData) error {
	content := fmt.Sprintf("Check’Ne平台服務費電子發票已經開立，並已寄送通知到會員中心，請登入後查閱發票資訊。")
	if err := PushMessageSms(UserData.Mphone, content); err != nil {
		return err
	}
	return nil
}

func SendServiceInvoiceSms(UserData entity.MemberData) error {
	content := fmt.Sprintf("Check’Ne加值服務費電子發票已經開立，並已寄送通知到會員中心，請登入後查閱發票資訊。")
	if err := PushMessageSms(UserData.Mphone, content); err != nil {
		return err
	}
	return nil
}

func SendSystemSms(UserData entity.MemberData) error {
	content := fmt.Sprintf("有一則Check’Ne的重要系統通知，請登入Check’Ne系統查看。")
	if err := PushMessageSms(UserData.Mphone, content); err != nil {
		return err
	}
	return nil
}

