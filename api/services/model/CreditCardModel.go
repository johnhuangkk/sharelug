package model

import (
	"api/services/Enum"
	"api/services/Service/Balance"
	"api/services/Service/OrderService"
	"api/services/Service/UserAddressService"
	"api/services/VO/Request"
	"api/services/VO/Response"
	"api/services/VO/UserAddress"
	"api/services/dao/Credit"
	"api/services/dao/Orders"
	"api/services/dao/UserAddressData"
	"api/services/dao/member"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"api/services/util/tools"
	"fmt"
	"strings"
	"time"
)

// 變更預設信用卡
func HandleChangeDefaultCredit(userData entity.MemberData, params Request.EditCreditRequest) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	err := member.UpdateMemberCardDataDefault(engine, userData.Uid)
	if err != nil {
		return fmt.Errorf("系統錯誤")
	}
	data, err := member.GetMemberCreditDataByCardIdAndUserId(engine, userData.Uid, params.CreditId)
	if err != nil {
		log.Error("Get Member Deposit  Error", err)
		return err
	}
	data.DefaultCard = "1"
	err = member.UpdateMemberCardData(engine, data)
	if err != nil {
		log.Error("Update Member Deposit Error", err)
		return err
	}
	return nil
}

//處理刪除提領帳號
func HandleDeleteCredit(userData entity.MemberData, params Request.EditCreditRequest) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	data, err := member.GetMemberCreditDataByCardIdAndUserId(engine, userData.Uid, params.CreditId)
	if err != nil {
		log.Error("Get Member Deposit Error", err)
		return err
	}
	data.Status = Enum.CreditStatusDelete
	err = member.UpdateMemberCardData(engine, data)
	if err != nil {
		log.Error("Update Member Deposit Error", err)
		return err
	}
	return nil
}

//取出卡片資料
func GetCreditData(userData entity.MemberData) (Response.MemberInfoResponse, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var resp Response.MemberInfoResponse
	var cardData []Response.MemberCardResponse
	data, err := member.GetMemberCreditDataByUserId(engine, userData.Uid)
	if err != nil {
		return resp, err
	}
	for _, v := range data {
		if len(v.CardId) != 0 {
			date := strings.Split(v.ExpiryDate, "")
			defaultCard := false
			if v.DefaultCard == "1" {
				defaultCard = true
			}
			rep := Response.MemberCardResponse{
				CardId: v.CardId,
				ExpiryDate: fmt.Sprintf("%s%s/%s%s", date[2], date[3], date[0], date[1]),
				CardNumber: fmt.Sprintf("**** **** **** %s", v.Last4Digits),
				DefaultCard: defaultCard,
			}
			cardData = append(cardData, rep)
		}
	}
	resp.MemberCard = cardData
	balance := Balance.GetBalanceByUid(engine, userData.Uid)
	resp.Balance = int64(balance)
	return resp, nil
}

func GetDeliveryLastAddress(userData entity.MemberData, params Request.GetAddressParams)  (UserAddress.AddressInfoResponse, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var resp UserAddress.AddressInfoResponse
	data, err := UserAddressData.GetDeliveryLastAddress(engine, userData.Uid, params.ShipType)
	if err != nil {
		return resp, fmt.Errorf("系統錯誤")
	}
	resp = UserAddressService.GetShipAddress(engine, data)
	return resp, nil
}

func QueryC2cCreditTrans(orderId string) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	data, err := Credit.GetGwCreditByOrderId(engine, orderId)
	if err != nil {
		return fmt.Errorf("1001001")
	}
	if len(data.OrderId) == 0 {
		return fmt.Errorf("1001010")
	}
	var query entity.QueryRequest
	query.OrderId = data.OrderId
	query.MerchantId = data.MerchantId
	query.TerminalId = data.TerminalId
	response, err := Balance.SendC2cQuery(query)
	if err != nil {
		return fmt.Errorf("1001001")
	}
	jsonData, _ := tools.JsonEncode(response)
	if err := Credit.InsertGwCreditAuthLog(engine, jsonData); err != nil {
		log.Error("insert gw credit log error", err)
		return err
	}
	params := response.GenerateCreditCheckParams()
	if err := Balance.Auth3DCheckout(engine, data.PayType, &params); err != nil {
		log.Error("Auth Check Out Error", err)
		return err
	}
	log.Debug("response", response)
	return nil
}
//信用卡受權檢查
func HandleCheckCredit() error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	//改為檢查30分鐘前的
	now := time.Now().Add(-time.Minute * 30 ).Format("2006-01-02 15:04")
	data, err := Credit.GetGwCreditAllByStatus(engine, now)
	if err != nil {
		return fmt.Errorf("1001001")
	}

	for _, v := range data {
		var query entity.QueryRequest
		query.OrderId = v.OrderId
		query.MerchantId = v.MerchantId
		query.TerminalId = v.TerminalId
		response, err := Balance.SendC2cQuery(query)
		if err != nil {
			return fmt.Errorf("1001001")
		}
		jsonData, _ := tools.JsonEncode(response)
		if err := Credit.InsertGwCreditAuthLog(engine, jsonData); err != nil {
			log.Error("insert gw credit log error", err)
			return err
		}
		params := response.GenerateCreditCheckParams()
		if err := Balance.Auth3DCheckout(engine, v.PayType, &params); err != nil {
			log.Error("Auth Check Out Error", err)
			return err
		}
		log.Debug("response", response)
	}
	return nil
}

func HandleCreditCancelOrder(orderId string) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	OrderData, err := Orders.GetOrderByOrderId(engine, orderId)
	if err != nil {
		log.Error("Get Order Error", err)
		return fmt.Errorf("1001001")
	}
	if len(OrderData.OrderId) == 0 {
		return fmt.Errorf("1001010")
	}
	if OrderData.OrderStatus != Enum.OrderSuccess {
		return fmt.Errorf("1001010")
	}
	if err := Balance.OrderCancelPaymentRefund(engine, OrderData); err != nil {
		return fmt.Errorf("1001001")
	}
	if err := OrderService.ChangeOrderStatus(engine, OrderData, Enum.OrderCancel); err != nil {
		return fmt.Errorf("1001001")
	}
	return nil
}

func HandleOrderShip(orderId string) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	OrderData, err := Orders.GetOrderByOrderId(engine, orderId)
	if err != nil {
		log.Error("Get Order Error", err)
		return fmt.Errorf("1001001")
	}
	if len(OrderData.OrderId) == 0 {
		return fmt.Errorf("1001010")
	}
	if OrderData.OrderStatus != Enum.OrderSuccess {
		return fmt.Errorf("1001010")
	}
	OrderData.ShipStatus = Enum.OrderShipSuccess
	OrderData.ShipTime = time.Now()
	OrderService.OrderCaptureRelease(&OrderData, time.Time{})
	if _, err := Orders.UpdateOrderData(engine, OrderData.OrderId, OrderData); err != nil {
		log.Error("Update Order Database Error", err)
		return fmt.Errorf("系統錯誤！")
	}
	return nil
}