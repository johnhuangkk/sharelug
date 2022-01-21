package Erp

import (
	"api/services/Enum"
	"api/services/VO/Response"
	"api/services/dao/Credit"
	"api/services/dao/Orders"
	"api/services/dao/member"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"api/services/util/tools"
	"fmt"
)

func GetAuditList(engine *database.MysqlSession, data []entity.GwCreditAuthData, resp *Response.AuditListResponse) error {
	for _, v := range data {
		if v.PayType == Enum.OrderTransC2c {
			res, err := GetC2cOrder(engine, v)
			if err != nil {
				return err
			}
			resp.AuditListData = append(resp.AuditListData, res)
		} else {
			res, err := GetB2cOrder(engine, v)
			if err != nil {
				return err
			}
			resp.AuditListData = append(resp.AuditListData, res)
		}
	}
	return nil
}

func GetC2cOrder(engine *database.MysqlSession, data entity.GwCreditAuthData) (Response.AuditListData, error) {
	var res Response.AuditListData
	OrderData, err := Orders.GetOrderByOrderId(engine, data.OrderId)
	if err != nil {
		log.Error("Get C2c Order Error", err)
		return res, err
	}
	userData, err := member.GetMemberDataByUid(engine, OrderData.BuyerId)
	if err != nil {
		log.Error("Get Member Error", err)
		return res, err
	}
	cardData, err := member.GetMemberCreditDataByCardIdAndUserId(engine, userData.Uid, data.CardId)
	if err != nil {
		log.Error("Get Member Card Error", err)
		return res, err
	}
	res.OrderId = data.OrderId
	res.PaymentTime = data.CreateTime.Format("2006/01/02 15:04")
	res.BuyerName = OrderData.BuyerName
	res.OrderStatus = Enum.ErpOrderStatus[OrderData.OrderStatus]
	res.OrderAmount = int64(OrderData.TotalAmount)
	res.BuyerAccount = tools.MaskerPhoneLater(userData.Mphone)
	res.CardAccount = fmt.Sprintf("**** **** **** %s", cardData.Last4Digits)
	res.CardVerify = "N"
	if data.CreditType == Enum.OrderTrans3D {
		res.CardVerify = "Y"
	}
	return res, nil
}

func GetB2cOrder(engine *database.MysqlSession, data entity.GwCreditAuthData) (Response.AuditListData, error) {
	var res Response.AuditListData
	OrderData, err := Orders.GetB2cOrderByOrderId(engine, data.OrderId)
	if err != nil {
		log.Error("Get C2c Order Error", err)
		return res, err
	}
	userData, err := member.GetMemberDataByUid(engine, OrderData.UserId)
	if err != nil {
		log.Error("Get Member Error", err)
		return res, err
	}
	cardData, err := member.GetMemberCreditDataByCardIdAndUserId(engine, userData.Uid, data.CardId)
	if err != nil {
		log.Error("Get Member Card Error", err)
		return res, err
	}
	res.OrderId = data.OrderId
	res.PaymentTime = data.CreateTime.Format("2006/01/02 15:04")
	res.BuyerName = userData.RealName
	res.OrderStatus = Enum.ErpOrderStatus[OrderData.OrderStatus]
	res.OrderAmount = OrderData.Amount
	res.BuyerAccount = tools.MaskerPhoneLater(userData.Mphone)
	res.CardAccount = fmt.Sprintf("**** **** **** %s", cardData.Last4Digits)
	res.CardVerify = "N"
	if data.CreditType == Enum.OrderTrans3D {
		res.CardVerify = "Y"
	}
	return res, nil
}

func CountGwCreditByAuditStatus(engine *database.MysqlSession, transType, auditStatus string) int64 {
	count, err := Credit.CountGwCreditByAuditStatus(engine, transType, auditStatus)
	if err != nil {
		return 0
	}
	return count
}

//取出刷卡資料
func GetCreditAuthData(engine *database.MysqlSession, OrderData entity.OrderData) (Response.GetCreditAuthResponse, error) {
	var resp Response.GetCreditAuthResponse
	data, err := Credit.GetCreditByOrderId(engine, OrderData.OrderId)
	if err != nil {
		return resp, err
	}
	card, err := member.GetMemberCreditDataByCardIdAndUserId(engine, OrderData.BuyerId, data.CardId)
	if err != nil {
		return resp, err
	}
	resp.BankName = card.BankName
	resp.CardType = card.CardType
	resp.ResponseCode = data.ResponseCode
	resp.ResponseMsg = data.ResponseMsg
	resp.ApproveCode = data.ApproveCode
	resp.FirstTrans = "Y"
	if card.Frequency > 0 {
		resp.FirstTrans = "N"
	}
	resp.Foreign = "N"
	if card.IsForeign == 1 {
		resp.Foreign = "Y"
	}
	//fixme 這兩個沒資料
	resp.RiskNote = "--"
	resp.RiskList = "--"
	resp.IP = data.AuthIp
	resp.AuditStatus = Enum.OrderAuditStatus[data.AuditStatus]
	resp.PendingDate = "--"
	if !data.PendingTime.IsZero() {
		resp.PendingDate = data.PendingTime.Format("2006/01/02 15:04")
	}
	resp.NoteDate = "--"
	if !data.NoteTime.IsZero() {
		resp.NoteDate = data.NoteTime.Format("2006/01/02 15:04")
	}
	resp.RefusedDate = "--"
	if !data.RefusedTime.IsZero() {
		resp.RefusedDate = data.RefusedTime.Format("2006/01/02 15:04")
	}
	resp.ReleaseDate = "--"
	if !data.ReleaseTime.IsZero() {
		resp.ReleaseDate = data.ReleaseTime.Format("2006/01/02 15:04")
	}
	resp.CaptureStatus = Enum.CreditCaptureStatus[data.CaptureStatus]
	resp.CaptureDate = "--"
	if !data.CaptureTime.IsZero() {
		resp.CaptureDate = data.CaptureTime.Format("2006/01/02 15:04")
	}
	resp.BatchDate = "--"
	if !data.BatchTime.IsZero() {
		resp.BatchDate = data.BatchTime.Format("2006/01/02 15:04")
	}
	resp.VoidDate = "--" //fixme 這兩個沒資料
	resp.RetreatDate = "--"
	return resp, nil
}


