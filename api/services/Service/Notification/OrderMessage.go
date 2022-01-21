package Notification

import (
	"api/services/Enum"
	"api/services/Service/Mail"
	"api/services/Service/Sms"
	"api/services/dao/Orders"
	"api/services/dao/Store"
	"api/services/dao/member"
	"api/services/dao/transfer"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"fmt"
)

//成功訂單訊息傳送
func SendSuccessfullyOrderedMessage(engine *database.MysqlSession, OrderId string) error {
	OrderData, err := Orders.GetOrderByOrderId(engine, OrderId)
	if err != nil {
		return err
	}
	//賣場資料,管理者資料
	StoreData, ManagerData, err := GetStoreManagerData(engine, OrderData.StoreId)
	if err != nil {
		return err
	}
	DetailData, err := Orders.GetOrderDetailByOrderId(engine, OrderData.OrderId)
	if err != nil {
		return err
	}
	for _, v := range ManagerData {
		UserData, err := member.GetMemberDataByUid(engine, v.UserId)
		if err != nil {
			return err
		}
		if len(UserData.Email) != 0 {
			if err := Mail.SendSuccessfullySellerOrderedMail(UserData, StoreData, OrderData); err != nil {
				return err
			}
		}
		if err := Sms.SendSuccessfullySellerOrderedSms(UserData, StoreData); err != nil {
				return err
			}
	}
	msg := fmt.Sprintf("%s 新增一筆訂單，訂單編號：%s，訂購商品：[%s]...。",
		OrderData.CreateTime.Format("2006/01/02 15:04"), OrderData.OrderId, DetailData[0].ProductName)
	if err := SendSystemNotify(engine, StoreData.StoreId, msg, Enum.NotifyMsgTypeOrderSeller, OrderData.OrderId); err != nil {
		return err
	}
	//買家資料
	BuyerData, err := member.GetMemberDataByUid(engine, OrderData.BuyerId)
	if err != nil {
		return err
	}
	if len(BuyerData.Email) != 0 {
		if err := Mail.SendSuccessfullyBuyerOrderedMail(BuyerData, StoreData, OrderData); err != nil {
			return err
		}
	}
	if err := Sms.SendSuccessfullyBuyerOrderSms(BuyerData, StoreData); err != nil {
			return err
		}
	msg = fmt.Sprintf("%s 新增一筆訂單，訂單編號：%s，訂購商品：[%s]...。",
		OrderData.CreateTime.Format("2006/01/02 15:04"), OrderData.OrderId, DetailData[0].ProductName)
	if err := SendSystemNotify(engine, BuyerData.Uid, msg, Enum.NotifyMsgTypeOrderBuyer, OrderData.OrderId); err != nil {
		return err
	}
	return nil
}

//ATM 成功訂單訊息傳送
func SendAtmSuccessfullyOrderedMessage(engine *database.MysqlSession, OrderId string) error {
	OrderData, err := Orders.GetOrderByOrderId(engine, OrderId)
	if err != nil {
		return err
	}
	//賣場資料,管理者資料
	StoreData, ManagerData, err := GetStoreManagerData(engine, OrderData.StoreId)
	if err != nil {
		return err
	}
	DetailData, err := Orders.GetOrderDetailByOrderId(engine, OrderData.OrderId)
	if err != nil {
		return err
	}
	for _, v := range ManagerData {
		UserData, err := member.GetMemberDataByUid(engine, v.UserId)
		if err != nil {
			return err
		}
		if len(UserData.Email) != 0 {
			err := Mail.SendSuccessfullySellerOrderedMail(UserData, StoreData, OrderData)
			if err != nil {
				return err
			}
		}
		if err := Sms.SendSuccessfullySellerOrderedSms(UserData, StoreData); err != nil {
				return err
			}
	}
	msg := fmt.Sprintf("%s 新增一筆訂單，訂單編號：%s，訂購商品：[%s]...。",
		OrderData.CreateTime.Format("2006/01/02 15:04"), OrderData.OrderId, DetailData[0].ProductName)
	if err := SendSystemNotify(engine, StoreData.StoreId, msg, Enum.NotifyMsgTypeOrderSeller, OrderData.OrderId); err != nil {
		return err
	}
	//買家資料
	BuyerData, err := member.GetMemberDataByUid(engine, OrderData.BuyerId)
	if err != nil {
		return err
	}
	transferData, err := transfer.GetTransferByOrderId(engine, OrderId)
	if err != nil {
		return err
	}
	if len(BuyerData.Email) != 0 {
		err := Mail.SendAtmSuccessfullyBuyerOrderedMail(BuyerData, StoreData, OrderData, transferData)
		if err != nil {
			return err
		}
	}
	if err := Sms.SendAtmSuccessfullyBuyerOrderedSms(BuyerData, StoreData, transferData); err != nil {
			return err
		}
	date := transferData.ExpireDate.Format("2006/01/02") + " 23:59"
	bank := fmt.Sprintf("%s%s", transferData.BankCode, transferData.BankName)
	msg = fmt.Sprintf("%s 新增訂單，訂單編號：%s，訂購商品：%s...，請於%s前繳款，金額%d元，\n轉帳銀行：%s，轉帳帳號%s。",
		OrderData.CreateTime.Format("2006/01/02 15:04"), OrderData.OrderId, DetailData[0].ProductName, date, transferData.Amount, bank, transferData.BankAccount)
	if err := SendSystemNotify(engine, BuyerData.Uid, msg, Enum.NotifyMsgTypeOrderBuyer, OrderData.OrderId); err != nil {
		return err
	}
	return nil
}

//訂單通知待出貨 （找不到那理要用）
func SendWaitShipMessage(engine *database.MysqlSession, OrderData entity.OrderData) error {
	StoreData, ManagerData, err := GetStoreManagerData(engine, OrderData.StoreId)
	if err != nil {
		return err
	}
	for _, v := range ManagerData {
		UserData, err := member.GetMemberDataByUid(engine, v.UserId)
		if err != nil {
			return err
		}
		if len(UserData.Email) != 0 {
			if err := Mail.SendWaitShipMail(UserData, StoreData, OrderData); err != nil {
				return err
			}
		}
		if err := Sms.SendWaitShipSms(UserData); err != nil {
				return err
			}
		var msg string
		msg += fmt.Sprintf("訂單編號：%s，訂購時間：%s，請儘速安排出貨。", OrderData.OrderId, OrderData.CreateTime.Format("2006/01/02 15:04"))
		if err = SendSystemNotify(engine, UserData.Uid, msg, Enum.NotifyMsgTypeOrder, OrderData.OrderId); err != nil {
			return err
		}
	}
	return nil
}

//訂單2天未出貨
func SendNotShippedMessage(engine *database.MysqlSession, storeId string, OrderData []entity.OrderData, day int64) error {
	StoreData, ManagerData, err := GetStoreManagerData(engine, storeId)
	if err != nil {
		return err
	}
	for _, v := range ManagerData {
		UserData, err := member.GetMemberDataByUid(engine, v.UserId)
		if err != nil {
			return err
		}
		if len(UserData.Email) != 0 {
			if err := Mail.SendNotShippedMail(UserData, len(OrderData)); err != nil {
				return err
			}
		}
		if err := Sms.SendNotShippedSms(UserData, len(OrderData)); err != nil {
				return err
			}
	}
	var msg string
	for _, v := range OrderData {
		msg += fmt.Sprintf("訂單編號：%s，訂購時間：%s 已超過%v天尚未出貨，請儘速安排出貨。", v.OrderId, v.PayWayTime.Format("2006/01/02"), day)
	}
	err = SendSystemNotify(engine, StoreData.StoreId, msg, Enum.NotifyMsgTypeOrderSeller, "")
	if err != nil {
		return err
	}
	return nil
}

//訂單已取消
func SendOrderCancelMessage(engine *database.MysqlSession, OrderId string) error {
	OrderData, err := Orders.GetOrderByOrderId(engine, OrderId)
	if err != nil {
		return err
	}
	StoreData, ManagerData, err := GetStoreManagerData(engine, OrderData.StoreId)
	if err != nil {
		return err
	}
	for _, v := range ManagerData {
		UserData, err := member.GetMemberDataByUid(engine, v.UserId)
		if err != nil {
			return err
		}
		if err := Sms.SendOrderCancelSellerSms(UserData, StoreData); err != nil {
			return err
		}
	}
	msg := fmt.Sprintf("訂單編號：%s已辦理訂單取消，此筆訂單若已繳款，將會在3個營業日內辦理退款。", OrderData.OrderId)
	if err = SendSystemNotify(engine, StoreData.StoreId, msg, Enum.NotifyMsgTypeOrderSeller, OrderData.OrderId); err != nil {
		return err
	}
	BuyerData, err := member.GetMemberDataByUid(engine, OrderData.BuyerId)
	if err != nil {
		return err
	}
	if err := Sms.SendOrderCancelBuyerSms(BuyerData, StoreData); err != nil {
		return err
	}
	msg = fmt.Sprintf("你在[%s]的訂單編號：%s已取消，若已繳款，將會在3個營業日內辦理退款。", StoreData.StoreName, OrderData.OrderId)
	if err := SendSystemNotify(engine, BuyerData.Uid, msg, Enum.NotifyMsgTypeOrderBuyer, OrderData.OrderId); err != nil {
		return err
	}
	return nil
}

//訂單逾期未寄 ---(有綁定支援物流)
func SendOrderShipOverdueMessage(engine *database.MysqlSession, OrderData entity.OrderData) error {
	StoreData, ManagerData, err := GetStoreManagerData(engine, OrderData.StoreId)
	if err != nil {
		return err
	}
	for _, v := range ManagerData {
		log.Debug("ManagerData", ManagerData)
		UserData, err := member.GetMemberDataByUid(engine, v.UserId)
		if err != nil {
			return err
		}
		if len(UserData.Email) != 0 {
			err := Mail.SendOrderShipOverdueMail(UserData, StoreData, OrderData)
			if err != nil {
				return err
			}
		}
		if err := Sms.SendOrderShipOverdueSms(UserData, OrderData); err != nil {
				return err
			}
	}
	msg := fmt.Sprintf("訂單編號：%s 逾期未寄，出貨單號：%s 已失效無法寄送，若此筆訂單已繳款，請儘早辦理退款。",
		OrderData.OrderId, OrderData.ShipNumber)
	err = SendSystemNotify(engine, StoreData.StoreId, msg, Enum.NotifyMsgTypeOrderSeller, OrderData.OrderId)
	if err != nil {
		return err
	}
	//BuyerData, err := member.GetMemberDataByUid(engine, OrderData.BuyerId)
	//if err != nil {
	//	return err
	//}
	//msg = fmt.Sprintf("[%s]訂單編號：%s 逾期未出貨，若已繳款，請與商家聯繫退款事宜。", StoreData.StoreName, OrderData.OrderId)
	//err = SendSystemNotify(engine, BuyerData.Uid, msg, Enum.NotifyMsgTypeOrderBuyer, OrderData.OrderId)
	//if err != nil {
	//	return err
	//}
	return nil
}

//訂單-已出貨
func SendShippedMessage(engine *database.MysqlSession, OrderData entity.OrderData) error {
	StoreData, _, err := GetStoreManagerData(engine, OrderData.StoreId)
	if err != nil {
		return err
	}
	buyerData, err := member.GetMemberDataByUid(engine, OrderData.BuyerId)
	if err != nil {
		return err
	}
	if err := Sms.SendShippedMessageSms(buyerData, StoreData, OrderData); err != nil {
		log.Error("交寄簡訊寄送失敗", err.Error(), OrderData.OrderId)
		return err
	}
	msg := fmt.Sprintf("checkne[%s]訂單編號：%s 已出貨，追蹤配送進度請到checkne.com查詢。", StoreData.StoreName, OrderData.OrderId)
	if err := SendSystemNotify(engine, buyerData.Uid, msg, Enum.NotifyMsgTypeOrderBuyer, OrderData.OrderId); err != nil {
		return err
	}
	return nil
}

//出貨買家取件閉轉店通知
func SendToBuyerShipClosedShopMessage(engine *database.MysqlSession, OrderData entity.OrderData) error {
	StoreData, ManagerData, err := GetStoreManagerData(engine, OrderData.StoreId)
	if err != nil {
		return err
	}
	for _, v := range ManagerData {
		UserData, err := member.GetMemberDataByUid(engine, v.UserId)
		if err != nil {
			return err
		}
		if len(UserData.Email) != 0 {
			if err := Mail.SendForSellerShipClosedShopMail(UserData, StoreData, OrderData); err != nil {
				return err
			}
		} else {
			if err := Sms.SendForSellerShipClosedShopSms(UserData, OrderData); err != nil {
				return err
			}
		}
	}
	msg := fmt.Sprintf("訂單編號：%s  指定收件超商因故無法提供取件服務，請通知買家在2日內重新選擇配送門市，以免配送時間過期。", OrderData.OrderId)
	if err := SendSystemNotify(engine, StoreData.StoreId, msg, Enum.NotifyMsgTypeOrderSeller, OrderData.OrderId); err != nil {
		return err
	}
	BuyerData, err := member.GetMemberDataByUid(engine, OrderData.BuyerId)
	if err != nil {
		return err
	}
	if len(BuyerData.Email) != 0 {
		if err := Mail.SendForBuyerShipClosedShopMail(BuyerData, StoreData, OrderData); err != nil {
			return err
		}
	}
	if err := Sms.SendForBuyerShipClosedShopSms(BuyerData, OrderData); err != nil {
			return err
		}
	msg = fmt.Sprintf("訂單編號：%s 指定收件超商因故無法提供取件服務，請重新選擇配送門市，以免配送時間過期。", OrderData.OrderId)
	if err := SendSystemNotify(engine, BuyerData.Uid, msg, Enum.NotifyMsgTypeOrderBuyer, OrderData.OrderId); err != nil {
		return err
	}
	return nil
}

//出貨貨到門市第一天
func SendShipToShopFirstDayMessage(engine *database.MysqlSession, OrderData entity.OrderData) error {
	StoreData, _, err := GetStoreManagerData(engine, OrderData.StoreId)
	if err != nil {
		return err
	}
	msg := fmt.Sprintf("訂單編號：%s 包裹已到指定門市。請聯繫買家儘速取件以免逾期退回。", OrderData.OrderId)
	if err = SendSystemNotify(engine, StoreData.StoreId, msg, Enum.NotifyMsgTypeOrderSeller, OrderData.OrderId); err != nil {
		return err
	}
	BuyerData, err := member.GetMemberDataByUid(engine, OrderData.BuyerId)
	if err != nil {
		return err
	}
	msg = fmt.Sprintf("訂單編號：%s  包裹已到指定門市，請儘速取件以免逾期退回。", OrderData.OrderId)
	if err := SendSystemNotify(engine, BuyerData.Uid, msg, Enum.NotifyMsgTypeOrderBuyer, OrderData.OrderId); err != nil {
		return err
	}
	return nil
}

//出貨貨到門市第四天
func SendShipToShopFourDayMessage(engine *database.MysqlSession, OrderData entity.OrderData) error {
	StoreData, _, err := GetStoreManagerData(engine, OrderData.StoreId)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("訂單編號：%s 包裹已到指定門市第四天，買家尚未取件。請聯繫買家儘速取件以免逾期退回。", OrderData.OrderId)
	if err := SendSystemNotify(engine, StoreData.StoreId, msg, Enum.NotifyMsgTypeOrderSeller, OrderData.OrderId); err != nil {
		return err
	}
	BuyerData, err := member.GetMemberDataByUid(engine, OrderData.BuyerId)
	if err != nil {
		return err
	}
	msg = fmt.Sprintf("訂單編號：%s 包裹已到指定門市第四天，請儘速取件以免逾期退回。", OrderData.OrderId)
	if err := SendSystemNotify(engine, BuyerData.Uid, msg, Enum.NotifyMsgTypeOrderBuyer, OrderData.OrderId); err != nil {
		return err
	}
	return nil
}

//訂單逾期未取
func SendOrderOverdueMessage(engine *database.MysqlSession, OrderData entity.OrderData) error {
	StoreData, _, err := GetStoreManagerData(engine, OrderData.StoreId)
	if err != nil {
		return err
	}
	msg := fmt.Sprintf("CheckNe訂單編號：%s 包裹送達指定門市逾期未領取，已退回超商物流中心，包裹將退回原交寄超商門市，請留意手機取件通知。\n如此筆訂單買家已繳款完成，請與買家確認後辦理退款。", OrderData.OrderId)
	if err := SendSystemNotify(engine, StoreData.StoreId, msg, Enum.NotifyMsgTypeOrderSeller, OrderData.OrderId); err != nil {
		return err
	}
	BuyerData, err := member.GetMemberDataByUid(engine, OrderData.BuyerId)
	if err != nil {
		return err
	}
	msg = fmt.Sprintf("Check'Ne訂單編號：%s  包裹送達指定門市逾期未領取，已啟動退回賣家，如需辦理退款，請與賣家聯絡。", OrderData.OrderId)
	if err := SendSystemNotify(engine, BuyerData.Uid, msg, Enum.NotifyMsgTypeOrderBuyer, OrderData.OrderId); err != nil {
		return err
	}
	return nil
}

//退貨-賣家取件-閉轉店通知
func SendReturnShipClosedShopMessage(engine *database.MysqlSession, OrderData entity.OrderData) error {
	StoreData, ManagerData, err := GetStoreManagerData(engine, OrderData.StoreId)
	if err != nil {
		return err
	}
	for _, v := range ManagerData {
		UserData, err := member.GetMemberDataByUid(engine, v.UserId)
		if err != nil {
			return err
		}
		if err := Sms.SendReturnShipClosedShopSms(UserData, OrderData); err != nil {
			return err
		}
	}
	msg := fmt.Sprintf("CheckNe訂單編號：%s 包裹逾期未取，原交寄門市因故無法提供退回商品之取件服務，請在2日內重新選擇門市位置，以免損失權益。", OrderData.OrderId)
	if err := SendSystemNotify(engine, StoreData.StoreId, msg, Enum.NotifyMsgTypeOrderSeller, OrderData.OrderId); err != nil {
		return err
	}
	return nil
}

//退貨-退貨申請
func SendReturnApplyMessage(engine *database.MysqlSession, OrderData entity.OrderData, returnData entity.OrderRefundData) error {
	StoreData, _, err := GetStoreManagerData(engine, OrderData.StoreId)
	if err != nil {
		return err
	}
	msg := fmt.Sprintf("Check'Ne訂單編號：%s，商品名稱：%s... 已完成退貨申請。",
		OrderData.OrderId, returnData.ProductName)
	if err := SendSystemNotify(engine, StoreData.StoreId, msg, Enum.NotifyMsgTypeOrderSeller, OrderData.OrderId); err != nil {
		return err
	}
	BuyerData, err := member.GetMemberDataByUid(engine, OrderData.BuyerId)
	if err != nil {
		return err
	}
	msg = fmt.Sprintf("Check'Ne訂單編號：%s，商品名稱：%s... 已完成退貨申請，請儘速將退貨商品寄回。", OrderData.OrderId, returnData.ProductName)
	if err := SendSystemNotify(engine, BuyerData.Uid, msg, Enum.NotifyMsgTypeOrderBuyer, OrderData.OrderId); err != nil {
		return err
	}
	return nil
}

//退貨-退貨完成
func SendReturnSuccessMessage(engine *database.MysqlSession, OrderData entity.OrderData, returnData entity.OrderRefundData) error {
	StoreData, _, err := GetStoreManagerData(engine, OrderData.StoreId)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("訂單編號：%s，商品名稱：%s... 已完成退貨。",
		OrderData.OrderId, returnData.ProductName)
	if err := SendSystemNotify(engine, StoreData.StoreId, msg, Enum.NotifyMsgTypeOrderSeller, OrderData.OrderId); err != nil {
		return err
	}

	BuyerData, err := member.GetMemberDataByUid(engine, OrderData.BuyerId)
	if err != nil {
		return err
	}
	msg = fmt.Sprintf("Check'Ne訂單編號：%s，商品名稱：%s... 已完成退貨。", OrderData.OrderId, returnData.ProductName)
	if err := SendSystemNotify(engine, BuyerData.Uid, msg, Enum.NotifyMsgTypeOrderBuyer, OrderData.OrderId); err != nil {
		return err
	}
	return nil
}

//退款申請
func SendRefundApplyMessage(engine *database.MysqlSession, OrderData entity.OrderData, refundData entity.OrderRefundData) error {
	StoreData, _, err := GetStoreManagerData(engine, OrderData.StoreId)
	if err != nil {
		return err
	}
	msg := fmt.Sprintf("訂單編號：%s 已完成退款申請，退款金額NT$%v元，將退到買家的check'Ne餘額中。",
		OrderData.OrderId, int64(refundData.Total))
	if err := SendSystemNotify(engine, StoreData.StoreId, msg, Enum.NotifyMsgTypeOrderSeller, OrderData.OrderId); err != nil {
		return err
	}
	BuyerData, err := member.GetMemberDataByUid(engine, OrderData.BuyerId)
	if err != nil {
		return err
	}
	if err := Sms.SendBuyerRefundApplySms(BuyerData, StoreData, refundData); err != nil {
		return err
	}
	msg = fmt.Sprintf("訂單編號：%s 退款金額NT$%v元，將退到check'Ne餘額中。", OrderData.OrderId, int64(refundData.Total))
	if err := SendSystemNotify(engine, BuyerData.Uid, msg, Enum.NotifyMsgTypeOrderBuyer, OrderData.OrderId); err != nil {
		return err
	}
	return nil
}
//訂購單賣家接受訂單
func SendBillOrderSuccessMessage(engine *database.MysqlSession, OrderData entity.OrderData) error {
	detail, err := Orders.GetOrderDetailByOrderId(engine, OrderData.OrderId)
	if err != nil {
		return err
	}
	StoreData, manager, err := GetStoreManagerData(engine, OrderData.StoreId)
	if err != nil {
		return err
	}
	msg := fmt.Sprintf("%s 新增一筆新接受的訂單，訂單編號：%s，訂購商品：%s...，請儘速安排出貨。",
		OrderData.CreateTime.Format("2006/01/02 15:04"), OrderData.OrderId, detail[0].ProductName)
	if err := SendSystemNotify(engine, StoreData.StoreId, msg, Enum.NotifyMsgTypeOrderSeller, OrderData.OrderId); err != nil {
		return err
	}
	for _, v := range manager {
		UserData, err := member.GetMemberDataByUid(engine, v.UserId)
		if err != nil {
			return err
		}
		if len(UserData.Email) != 0 {
			if err := Mail.SendSellerBillOrderSuccessMail(UserData, StoreData, OrderData); err != nil {
				return err
			}
		}
		if err := Sms.SendSellerBillOrderSuccessSms(UserData, StoreData); err != nil {
			return err
		}
	}
	buyer, err := member.GetMemberDataByUid(engine, OrderData.BuyerId)
	if err != nil {
		return err
	}
	if len(buyer.Email) != 0 {
		if err := Mail.SendBuyerBillOrderSuccessMail(buyer, StoreData, OrderData); err != nil {
			return err
		}
	}
	if err := Sms.SendBuyerBillOrderSuccessSms(buyer, StoreData); err != nil {
		return err
	}
	msg = fmt.Sprintf("%s 新增一筆賣家接受的訂單，訂單編號：%s，訂購商品：%s...，賣家會儘速安排出貨。",
		OrderData.CreateTime.Format("2006/01/02 15:04"), OrderData.OrderId, detail[0].ProductName)
	if err := SendSystemNotify(engine, buyer.Uid, msg, Enum.NotifyMsgTypeOrderBuyer, OrderData.OrderId); err != nil {
		return err
	}
	return nil
}

//ATM 成功訂單訊息傳送
func SendBillOrderAtmMessage(engine *database.MysqlSession, OrderData entity.OrderData) error {
	//賣場資料,管理者資料
	StoreData, ManagerData, err := GetStoreManagerData(engine, OrderData.StoreId)
	if err != nil {
		return err
	}
	DetailData, err := Orders.GetOrderDetailByOrderId(engine, OrderData.OrderId)
	if err != nil {
		return err
	}
	for _, v := range ManagerData {
		UserData, err := member.GetMemberDataByUid(engine, v.UserId)
		if err != nil {
			return err
		}
		if len(UserData.Email) != 0 {
			err := Mail.SendSellerAtmBillOrderMail(UserData, StoreData, OrderData)
			if err != nil {
				return err
			}
		}
		if err := Sms.SendSellerBillOrderSuccessSms(UserData, StoreData); err != nil {
			return err
		}
	}
	msg := fmt.Sprintf("%s 新增一筆新接受的訂單，訂單編號：%s，訂購商品：%s...，此為ATM繳款帳單，買家繳款成功後，請儘速安排出貨流程。",
		OrderData.CreateTime.Format("2006/01/02 15:04"), OrderData.OrderId, DetailData[0].ProductName)
	if err := SendSystemNotify(engine, StoreData.StoreId, msg, Enum.NotifyMsgTypeOrderSeller, OrderData.OrderId); err != nil {
		return err
	}
	//買家資料
	BuyerData, err := member.GetMemberDataByUid(engine, OrderData.BuyerId)
	if err != nil {
		return err
	}
	transferData, err := transfer.GetTransferByOrderId(engine, OrderData.OrderId)
	if err != nil {
		return err
	}
	if len(BuyerData.Email) != 0 {
		err := Mail.SendBuyerAtmBillOrderMail(BuyerData, StoreData, OrderData, transferData)
		if err != nil {
			return err
		}
	}
	if err := Sms.SendBuyerAtmBillOrderSms(BuyerData, StoreData, transferData); err != nil {
		return err
	}
	date := transferData.ExpireDate.Format("2006/01/02") + " 23:59"
	bank := fmt.Sprintf("%s%s", transferData.BankCode, transferData.BankName)
	msg = fmt.Sprintf("%s 新增一筆賣家接受的訂單，訂單編號：%s，訂購商品：%s...，請於%s前繳款，金額%d元，\n轉帳銀行：%s，轉帳帳號%s，賣家會儘速安排出貨，",
		OrderData.CreateTime.Format("2006/01/02 15:04"), OrderData.OrderId, DetailData[0].ProductName, date, transferData.Amount, bank, transferData.BankAccount)
	if err := SendSystemNotify(engine, BuyerData.Uid, msg, Enum.NotifyMsgTypeOrderBuyer, OrderData.OrderId); err != nil {
		return err
	}
	return nil
}

func GetStoreManagerData(engine *database.MysqlSession, StoreId string) (entity.StoreData, []entity.StoreRankData, error) {
	var StoreData entity.StoreData
	var ManagerData []entity.StoreRankData
	StoreData, err := Store.GetStoreDataByStoreId(engine, StoreId)
	if err != nil {
		return StoreData, ManagerData, err
	}
	ManagerData, err = Store.GetStoreByStoreId(engine, StoreId)
	if err != nil {
		return StoreData, ManagerData, err
	}
	return StoreData, ManagerData, nil
}

