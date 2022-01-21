package model

import (
	"api/services/Enum"
	"api/services/Service/Balance"
	"api/services/Service/Carts"
	"api/services/Service/MemberService"
	"api/services/Service/Notification"
	"api/services/Service/Upgrade"
	"api/services/VO/Request"
	"api/services/VO/Response"
	"api/services/dao/Orders"
	"api/services/dao/Store"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"fmt"
	"time"
)

//取出帳單及購買資料
func HandleGetB2CPay(cookie string, userData entity.MemberData, storeData entity.StoreDataResp, params Request.GetB2CPayRequest) (Response.GetB2cPayResponse, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var resp Response.GetB2cPayResponse
	data, err := Upgrade.GeneratorNewOrder(engine, userData, storeData, params, cookie, &resp)
	if err != nil {
		return resp, err
	}
	var unpaidOrder []Response.OrderList
	var upgradeList []Response.UpgradeList
	for _, v := range data.Detail {
		if v.ProductType == Enum.B2cOrderTypeUpgrade {
			sign := true
			if len(v.ProductId) == 0 {
				sign = false
			}
			OrderList := Response.UpgradeList{
				UpgradeText:  v.ProductName,
				UpgradePrice: v.ProductAmount,
				SignType: sign,
			}
			upgradeList = append(upgradeList, OrderList)
		}
		if v.ProductType == Enum.B2cOrderTypeBilling {
			OrderList := Response.OrderList{
				OrderText:  v.ProductName,
				OrderPrice: v.ProductAmount,
			}
			unpaidOrder = append(unpaidOrder, OrderList)
		}
	}
	resp.UpgradeList = upgradeList
	resp.OrderList = unpaidOrder
	return resp, nil
}
//升級方案 結帳
func HandleB2CPay(cookie string, params Request.B2CPayRequest) (Response.B2cPayResponse, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var resp Response.B2cPayResponse
	//從購物車入取出訂單內容
	data, err := Carts.GetB2cRedisCarts(cookie, Enum.StyleB2C)
	if err != nil {
		log.Error("Get Carts Error", err)
		return resp, fmt.Errorf("1001001")
	}
	data.Payment = params.Payment
	data.OrderStatus = Enum.OrderWait
	//寫入個人發票資訊
	if err := setOrderInvoiceInfo(params, &data); err != nil {
		log.Error("Set Order Invoice Data Error", err)
		return resp, fmt.Errorf("1001001")
	}
	//更新B2C訂單
	data.CarrierType = params.Carrier.CarrierType
	data.CarrierId = params.Carrier.CarrierId
	data.InvoiceType = params.Carrier.InvoiceType
	data.CompanyName = params.Carrier.CompanyName
	data.CompanyBan = params.Carrier.CompanyBan
	data.DonateBan = params.Carrier.DonateBan
	data.OrderSys = 2

	if err := Orders.InsertB2cOrderData(engine, data); err != nil {
		log.Error("Update B2c Order Data Error", err)
		return resp, fmt.Errorf("1001001")
	}
	//付款結帳
	resp, err = Balance.ProcessPayment(engine, data, params)
	if err != nil {
		return resp, err
	}
	return resp, nil
}
//發票資訊
func setOrderInvoiceInfo(params Request.B2CPayRequest, data *entity.B2cOrderData) error {
	//判斷是否有送值進來
	if len(params.Carrier.InvoiceType) == 0 {
		return fmt.Errorf("錯誤")
	} else {
		data.InvoiceType = params.Carrier.InvoiceType
		data.CompanyBan = params.Carrier.CompanyBan
		data.CompanyName = params.Carrier.CompanyName
		data.DonateBan = params.Carrier.DonateBan
		data.CarrierType = params.Carrier.CarrierType
		data.CarrierId = params.Carrier.CarrierId
		if  params.Carrier.InvoiceType == Enum.InvoiceTypeDonate {
			//驗證捐贈碼
			if err := MemberService.VerifyDonateCode(params.Carrier.DonateBan); err != nil {
				log.Error("Verify Donate Code Error", err)
				return err
			}
		}
		if params.Carrier.CarrierType == Enum.InvoiceCarrierTypeMobile {
			//驗證手機碼
			if err := MemberService.VerifyMobileCode(params.Carrier.CarrierId); err != nil {
				log.Error("Verify Mobile Code Error", err)
				return err
			}
		}
	}
	return nil
}

func GeneratorWaitPaymentOrder(engine *database.MysqlSession, UserData entity.MemberData) error {
	storeData, err := Store.GetStoreDefaultDataByUid(engine, UserData.Uid)
	if err != nil {
		return err
	}
	//取出待付款的帳單
	billData, err := Orders.GetB2cBillingByUserIdAndExpire(engine, UserData.Uid, storeData.StoreId)
	if err != nil {
		return err
	}
	if len(billData.BillingId) == 0 {
		//沒有就新增一筆 再扣款
		data, err := Upgrade.GeneratorB2cBill(engine, UserData, storeData, UserData.UpgradeExpire)
		if err != nil {
			return err
		}
		//餘額足夠使用餘額扣款
		if err := Balance.UpgradeAutoPayment(engine, data); err != nil {
			return err
		}
	} else {
		//判斷待付款的帳單是否已過期
		if time.Now().Before(billData.Expiration) {
			return nil
		}
		//過期就再新增一筆
		data, err := Upgrade.GeneratorB2cBill(engine, UserData, storeData, billData.Expiration)
		if err != nil {
			return err
		}
		if err := Notification.SendUpgradePaymentMessage(engine, storeData.StoreId, storeData.SellerId, data.BillingId); err != nil {
			return err
		}

	}
	return nil
}