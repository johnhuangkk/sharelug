package Notification

import (
	"api/services/Enum"
	"api/services/Service/Mail"
	"api/services/Service/Sms"
	"api/services/dao/Orders"
	"api/services/dao/member"
	"api/services/database"
	"api/services/entity"
	"fmt"
)

//訂單客服-留言
func SendOrderCustomerMessage(engine *database.MysqlSession, OrderId string, count int64, data entity.OrderMessageBoardData) error {
	OrderData, err := Orders.GetOrderByOrderId(engine, OrderId)
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
			if err := Mail.SendBuyerOrderCustomerMail(UserData, data); err != nil {
				return err
			}
		} else {
			if err := Sms.SendOrderCustomerSms(UserData, StoreData, count); err != nil {
				return err
			}
		}

	}
	msg := fmt.Sprintf("訂單編號：%s，提出一筆客服問題，請儘速安排回覆。", OrderData.OrderId)
	if err := SendSystemNotify(engine, StoreData.StoreId, msg, Enum.NotifyMsgTypeCustomer, OrderData.OrderId); err != nil {
		return err
	}
	msg = fmt.Sprintf("你在 訂單編號：%s 中提出一筆客服問題，可以至訂單留言追蹤處理進度。", OrderData.OrderId)
	err = SendSystemNotify(engine, OrderData.BuyerId, msg, Enum.NotifyMsgTypeCustomer, OrderData.OrderId)
	if err != nil {
		return err
	}
	return nil
}

//訂單客服-回覆留言
func SendReplyOrderCustomerMessage(engine *database.MysqlSession, OrderData entity.OrderData) error {
	StoreData, _, err := GetStoreManagerData(engine, OrderData.StoreId)
	if err != nil {
		return err
	}
	BuyerData, err := member.GetMemberDataByUid(engine, OrderData.BuyerId)
	if err != nil {
		return err
	}
	err = Sms.SendReplyOrderCustomerSms(BuyerData, StoreData, OrderData)
	if err != nil {
		return err
	}
	msg := fmt.Sprintf("%s已針對訂單編號：%s 中提出客服問題，請至訂單留言 查看。", StoreData.StoreName, OrderData.OrderId)
	err = SendSystemNotify(engine, OrderData.BuyerId, msg, Enum.NotifyMsgTypeCustomer, OrderData.OrderId)
	if err != nil {
		return err
	}
	return nil
}

func SendUserCustomerMessage(engine *database.MysqlSession, userId, orderId, category string, questionId string) error {
	userData, err := member.GetMemberDataByUid(engine, userId)
	if err != nil {
		return err
	}
	var msg string
	if orderId != "" {
		msg = fmt.Sprintf("你提出的客服問題:%s。訂單編號:%s", category, orderId)
	} else {
		msg = fmt.Sprintf("你提出的客服問題:%s。", category)
	}
	err = SendPlatformNotify(engine, userData.Uid, msg, Enum.NotifyTypePlatformUser, questionId)
	return nil
}


//平台客服回覆
func SendCustomerReplyMessage(engine *database.MysqlSession, userId, orderId, category string, questionId string) error {
	userData, err := member.GetMemberDataByUid(engine, userId)
	if err != nil {
		return err
	}
	var msg string
	if orderId != "" {
		msg = fmt.Sprintf("來自 Check'Ne 的客服回覆:%s。訂單編號:%s", category, orderId)
	} else {
		msg = fmt.Sprintf("來自 Check'Ne 的客服回覆:%s。", category)
	}
	err = SendPlatformNotify(engine, userData.Uid, msg, Enum.NotifyTypePlatformService, questionId)
	return nil
}
