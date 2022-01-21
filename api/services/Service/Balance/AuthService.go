package Balance

import (
	"api/services/Enum"
	"api/services/Service/CreditService"
	"api/services/VO/OrderVo"
	"api/services/VO/Request"
	"api/services/dao/Credit"
	"api/services/dao/Orders"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"api/services/util/tools"
	"fmt"
	"github.com/spf13/viper"
	"time"
)

//刷卡處理
func HandleAuth(vo OrderVo.CreditPaymentVo, params *Request.PayParams, cardData entity.MemberCardData) (entity.AuthResponse, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	AuthParams := cardData.GenerateAuthRequest(vo, params)
	request, _ := tools.JsonEncode(AuthParams)
	log.Info("Request Credit Auth", request)
	auth, err := AuthProcess(engine, AuthParams)
	if err != nil {
		log.Error("Auth process Error!!")
		return auth, err
	}
	response, _ := tools.JsonEncode(auth)
	log.Info("Response Credit Auth", response)
	if auth.Status == Enum.OrderWait {
		return auth, nil
	} else if auth.Status == Enum.OrderFail {
		if err := AuthCheckout(engine, vo.OrderId, auth.Status, vo.OrderType); err != nil {
			log.Error("Order check out Error!!")
			return auth, err
		}
		return auth, fmt.Errorf("1005006")
	} else {
		if err := AuthCheckout(engine, vo.OrderId, Enum.OrderSuccess, vo.OrderType); err != nil {
			log.Error("Order check out Error!!")
			return auth, err
		}
	}
	return auth, nil
}
//信用卡非3D認證回應處理
func AuthCheckout(engine *database.MysqlSession, orderId, Status, Type string) error {
	switch Type {
		case Enum.OrderTransBill:
			if err := OrderBillCheckout(engine, orderId, Status); err != nil {
				return err
			}
		case Enum.OrderTransC2c:
			if err := OrderCheckout(engine, orderId, Status); err != nil {
				return err
			}
		case Enum.OrderTransB2c:
			if err := B2cOrderCheckout(engine, orderId, Status, ""); err != nil {
				return err
			}
	}
	return nil
}
//信用卡3D認證回應處理
func Auth3DCheckout(engine *database.MysqlSession, PayType string, params *Request.Credit3dCheckParams) error {
	var orderId string
	switch PayType {
		case Enum.OrderTransBill:
			data, err := Orders.GetBillOrderByOrderId(engine, params.OrderID)
			if err != nil {
				log.Error("Get B2c Order Data Error", err)
				return err
			}
			Status := Enum.OrderFail
			if params.ResponseCode == "00" {
				Status = Enum.OrderSuccess
			}
			if len(data.BillId) != 0 {
				orderId = data.BillId
				if err := OrderBillCheckout(engine, orderId, Status); err != nil {
					return err
				}
			}
		case Enum.OrderTransC2c:
			data, err := Orders.GetOrderByOrderId(engine, params.OrderID)
			if err != nil {
				log.Error("Get Order data Error", err)
				return err
			}
			Status := Enum.OrderFail
			if params.ResponseCode == "00" {
				Status = Enum.OrderSuccess
			}
			if len(data.OrderId) != 0 {
				orderId = data.OrderId
				if err := OrderCheckout(engine, orderId, Status); err != nil {
					return err
				}
			}
		case Enum.OrderTransB2c:
			data, err := Orders.GetB2cOrderByOrderId(engine, params.OrderID)
			if err != nil {
				log.Error("Get B2c Order Data Error", err)
				return err
			}

			Status := Enum.OrderFail
			if params.ResponseCode == "00" {
				Status = Enum.OrderSuccess
			}
			if len(data.OrderId) != 0 {
				orderId = data.OrderId
				if err := B2cOrderCheckout(engine, orderId, Status, ""); err != nil {
					return err
				}
			}
	}
	if len(orderId) == 0 {
		orderId = params.OrderID
	}
	GwData, err := Credit.GetGwCreditByOrderId(engine, orderId)
	if err != nil {
		log.Error("Get Gw Credit data Error", err)
		return err
	}
	if err := updateAuthGwData(engine, GwData, params); err != nil {
		log.Error("Update Gw Credit data Error", err)
		return err
	}
	return nil
}
//刷卡取授權
func AuthProcess(engine *database.MysqlSession, params entity.AuthRequest) (entity.AuthResponse, error) {
	var resp entity.AuthResponse
	//建立刷卡記錄
	var gwData entity.GwCreditAuthData
	var response entity.AuthResult
	//送刷卡資料到銀行
	gwData, err := CreditService.NewAuthGwData(engine, params)
	if err != nil {
		log.Error("Create gw data Error", err)
		return resp, err
	}
	params.MerchantId = gwData.MerchantId
	params.TerminalId = gwData.TerminalId
	response, err = CreditService.New(params).DoAuth()
	if err != nil {
		log.Error("Get Auth config Error", err)
		return resp, err
	}

	log.Info("Auth Credit response", response)
	if response.RtnCode == "1" {
		resp.RtnURL = response.RtnHtml
		resp.Status = Enum.OrderWait
		return resp, nil
	} else {
		//更新刷卡資料資料庫
		err := updateGwData(engine, gwData, response)
		if err != nil {
			return resp, err
		}
		if response.ResponseCode != "00" {
			resp.Status = Enum.OrderFail
		} else {
			resp.Status = Enum.OrderSuccess
		}
	}
	return resp, nil
}
//處理取消訂單
func VoidProcess(engine *database.MysqlSession, params entity.CancelRequest) error {
	//取出刷卡資料
	gwData, err := Credit.GetGwCreditByOrderId(engine, params.OrderId)
	if err != nil {
		log.Error("Get Credit Gw Error", err)
		return err
	}
	//判斷是否請款
	if gwData.BatchId == 0 {
		data, err := createVoidGwData(engine, gwData)
		if err != nil {
			log.Error("Creat credit auth error", err)
			return fmt.Errorf("系統錯誤")
		}
		//未請款 執行 取消受權
		params.MerchantId = data.MerchantId
		params.TerminalId = data.TerminalId
		response, err := sendC2cVoid(params)
		if err != nil {
			log.Error("Send credit auth error", err)
			return fmt.Errorf("系統錯誤")
		}
		data.ApproveCode = response.ApproveCode
		data.ResponseCode = response.ResponseCode
		data.ResponseMsg = response.ResponseMsg
		if response.ResponseCode == "00" {

			//改變此筆狀態
			data.TransStatus = Enum.CreditTransStatusSuccess
			if err = Credit.ChangeGwCreditStatus(engine, data); err != nil {
				log.Error("Change Gw Credit Status Error", err)
				return fmt.Errorf("系統錯誤")
			}
		} else {
			data.TransStatus = Enum.CreditTransStatusFail
			if err := Credit.ChangeGwCreditStatus(engine, data); err != nil {
				log.Error("Change Gw Credit Status Error", err)
				return fmt.Errorf("系統錯誤")
			}
			log.Error("Send credit void error", response)
			return fmt.Errorf("退款失敗")
		}
	} else {
		//已請款 執行 退款
		err := createRefundGwData(engine, gwData)
		if err != nil {
			log.Error("Create Credit Void Data Error", err)
			return fmt.Errorf("系統錯誤")
		}
	}
	return nil
}
//信用卡退款
func RefundPaymentProcess(engine *database.MysqlSession, OrderId string, Amount float64) error {
	gwData, err := Credit.GetGwCreditByOrderId(engine, OrderId)
	if err != nil {
		log.Error("Get Gw Credit Database Error", err)
		return err
	}
	gwData.OrderId = OrderId
	gwData.TramsAmount = int64(Amount)
	err = createRefundGwData(engine, gwData)
	if err != nil {
		log.Error("Create Credit Void Data Error", err)
		return fmt.Errorf("系統錯誤")
	}
	return nil
}
//送C2C取消訂單
func sendC2cVoid(params entity.CancelRequest) (entity.AuthResult, error) {
	var response entity.AuthResult
	//MerchantId := viper.GetString("KgiCredit.C2C.N3D.MerchantID")
	//TerminalId := viper.GetString("KgiCredit.C2C.N3D.TerminalID")
	AuthConfig := CreditService.SetConnectionConfig{
		ResponseLink: fmt.Sprintf("https://www.checkne.com/v1/pay/credit/confirm/c2c"),
		HostName: viper.GetString("KgiCredit.hostname"),
		Path: viper.GetString("KgiCredit.AuthPath"),
		MerchantId: params.MerchantId,
		TerminalId: params.TerminalId,
	}
	response, err := AuthConfig.DoVoid(params)
	if err != nil {
		log.Error("Send credit auth error", err)
		return response, err
	}
	return response, nil
}
//信用卡回傳結果回寫至GW刷卡資料
func updateGwData(engine *database.MysqlSession, gw entity.GwCreditAuthData, data entity.AuthResult) error {
	gw.ApproveCode = data.ApproveCode
	gw.ResponseCode = data.ResponseCode
	gw.ResponseMsg = data.ResponseMsg
	if data.ResponseCode != "00" {
		gw.TransStatus = Enum.CreditTransStatusFail
	} else {
		gw.TransStatus = Enum.CreditTransStatusSuccess
	}
	if err := Credit.UpdateGwCreditData(engine, gw); err != nil {
		return err
	}
	return nil
}
//取消受權
func createVoidGwData(engine *database.MysqlSession, GwData entity.GwCreditAuthData) (entity.GwCreditAuthData, error) {
	result := GwData.GenerateGwCreditVoidData()
	data, err := Credit.InsertGwCreditData(engine, result)
	if err != nil {
		return data, err
	}
	return data, nil
}
//建立信用卡退款資料
func createRefundGwData(engine *database.MysqlSession, GwData entity.GwCreditAuthData) error {
	result := GwData.GenerateGwCreditRefundData()
	_, err := Credit.InsertGwCreditData(engine, result)
	if err != nil {
		return err
	}
	return nil
}
//信用卡3D認證回應
func Credit3DConfirm(PayType string, params *Request.Credit3dCheckParams) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	jsonData, _ := tools.JsonEncode(params)
	if err := Credit.InsertGwCreditAuthLog(engine, jsonData); err != nil {
		log.Error("insert gw credit log error", err)
		return err
	}
	if params.TransCode != "00" {
		return fmt.Errorf("credit error code: %s", params.TransCode)
	}
	//取出訂單
	if err := Auth3DCheckout(engine, PayType, params); err != nil {
		log.Error("Auth Check Out Error", err)
		return err
	}
	return nil
}
//變更信用卡記錄
func updateAuthGwData(engine *database.MysqlSession, gw entity.GwCreditAuthData, data *Request.Credit3dCheckParams) error {
	trans := fmt.Sprintf("%s %s", data.TransDate, data.TransTime)
	day, _ := time.ParseInLocation("20060102 150405", trans, time.Local)
	gw.ApproveCode = data.ApproveCode
	gw.ResponseCode = data.ResponseCode
	gw.ResponseMsg = data.ResponseMsg
	gw.TransTime = day
	if data.ResponseCode != "00" {
		gw.TransStatus = Enum.CreditTransStatusFail
	} else {
		gw.TransStatus = Enum.CreditTransStatusSuccess
	}
	if err := Credit.UpdateGwCreditData(engine, gw); err != nil {
		return err
	}
	return nil
}
//信用卡回傳結果 寫入LOG
func AuthResponse(params entity.AuthParams) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	jsonData, _ := tools.JsonEncode(params)
	if err := Credit.InsertGwCreditAuthLog(engine, jsonData); err != nil {
		return err
	}
	return nil
}

//送C2C訂單查詢
func SendC2cQuery(params entity.QueryRequest) (entity.QueryResponse, error) {
	var response entity.QueryResponse
	//MerchantId := viper.GetString("KgiCredit.C2C.N3D.MerchantID")
	//TerminalId := viper.GetString("KgiCredit.C2C.N3D.TerminalID")
	AuthConfig := CreditService.SetConnectionConfig{
		ResponseLink: fmt.Sprintf("https://www.checkne.com/v1/pay/credit/confirm/c2c"),
		HostName: viper.GetString("KgiCredit.hostname"),
		Path: viper.GetString("KgiCredit.AuthPath"),
		MerchantId: params.MerchantId,
		TerminalId: params.TerminalId,
	}
	response, err := AuthConfig.DoQuery(params)
	if err != nil {
		log.Error("Send credit auth error", err)
		return response, err
	}
	return response, nil
}



