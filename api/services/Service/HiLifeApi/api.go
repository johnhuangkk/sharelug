package HiLifeApi

import (
	"api/services/Enum"
	"api/services/Service/MartHiLife"
	"api/services/dao/Cvs"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"encoding/xml"
	"fmt"
	"strings"
)

// 萊爾富取號
func GetShippingOrderNo(engine *database.MysqlSession, orderData entity.OrderData, sellerData entity.MemberData) (string, error) {
	var order = MartHiLife.OrderAddRequest{}
	order.SetAddShipNoParams(orderData, sellerData)
	log.Info("GetShippingOrderNo Info Data: [%s]", order)

	var client MartHiLife.Client
	resp, err := client.GetClient().OrderAdd(order)

	if err != nil {
		log.Error("OrderAddRequest:", err)
		return ``, fmt.Errorf(`取號發生錯誤`)
	}
	if resp.ErrorCode != "000" {
		log.Error("GetShippingOrderNo Error: [%s]", resp.ErrorMessage)
		return ``, fmt.Errorf(`取號發生錯誤`)
	}

	var data entity.CvsShippingData
	data.InitInsert(Enum.CVS_HI_LIFE)

	data.ParentId = order.ParentId
	data.EcOrderNo = order.VdrOrderNo
	data.ShipNo = resp.OrderNo
	data.ServiceType = order.GetServiceType()
	data.SenderName = order.SenderName
	data.SenderPhone = order.SenderPhone
	data.OriReceiverAddress = orderData.ReceiverAddress

	// 建立託運單
	err = Cvs.InsertCvsShippingData(engine, data)

	return resp.OrderNo, nil
}

// 列印託運單
func PrintShippingOrder(shipNos []string) (data []byte, err error) {
	var client MartHiLife.Client

	shipNosString := strings.Join(shipNos, ";")
	data, err = client.GetClient().OrderPrint(shipNosString)

	return
}


// 修改閉轉店
func SwitchStore(cvsShippingData entity.CvsShippingData, switchLog entity.CvsShippingLogData, newStoreId string) (err error)  {
	var client MartHiLife.Client
	var req MartHiLife.OrderSwitchRequest

	var switchLogXml = struct {
		ParentId      string `xml:"ParentId" json:"ParentId"` // EC 客戶代號 length 3
		EshopId       string `xml:"EshopId" json:"EshopId"` // EC 網站代號 length 3
		OrderNo       string `xml:"OrderNo" json:"OrderNo"` // 寄件代碼 length 13
		EcOrderNo       string `xml:"EcOrderNo" json:"EcOrderNo"` // 訂單單號 length 11
		OriginStoreId       string `xml:"OriginStoreId" json:"OriginStoreId"` // 門市店號 length 4
		StoreType     string `xml:"StoreType" json:"StoreType"` // 店鋪類型 1:寄件店 2:取件店
		ChkMac     string `xml:"ChkMac" json:"ChkMac"` // 檢查碼*
	}{}


	err = xml.Unmarshal([]byte(switchLog.Log), &switchLogXml)

	if err != nil {
		log.Error("switchLogXml Error [%v]", err.Error())
		return fmt.Errorf("解析轉店失敗")
	}

	storeType := `2`
	if cvsShippingData.FlowType == `R` {
		storeType = `1`
	}

	req.GetRequest()
	req.ShipNo = cvsShippingData.ShipNo
	req.EcOrderNo = switchLogXml.EcOrderNo
	req.RcvStoreId = newStoreId
	req.StoreType = storeType

	resp, err := client.GetClient().OrderSwitchStore(req)

	if err != nil {
		log.Error("MartHiLifeSwitchStore", req, err)
		return err
	}

	if resp.Doc.ErrorCode != "000" {
		return fmt.Errorf(resp.Doc.ErrorMessage)
	}

	return nil
}


// 閉轉店
func MartHiLifeSwitchStore(shipNo, ecOrderNo, newStoreId string, isReceiveStore bool) error {
	var client MartHiLife.Client

	storeType := "1"
	if isReceiveStore {
		storeType = "2"
	}

	var req MartHiLife.OrderSwitchRequest
	req.GetRequest()
	req.ShipNo = shipNo
	req.EcOrderNo = ecOrderNo
	req.RcvStoreId = newStoreId
	req.StoreType = storeType

	resp, err := client.GetClient().OrderSwitchStore(req)
	if err != nil || resp.Doc.ErrorCode != "000" {
		log.Error("MartHiLifeSwitchStore", req, err)
		return err
	}

	return nil
}

// 閉轉通知
//func SwitchNotification(noti MartHiLife.ChangeNotification, raw string) error {
//
//	if !MartHiLife.SwitchCheckSum(noti) {
//		return fmt.Errorf("檢查失敗")
//	}

	//var UpdateCvsShipping CvsShipping.UpdateCvsShipping
	//UpdateCvsShipping.Type = `Switch`
	//UpdateCvsShipping.ShipType = Enum.CVS_HI_LIFE


	//for
	//for _,d := range noti.Doc {
	//
	//}
	//if noti.Doc.StoreType == `1` {
	//
	//} else {
	//
	//}

	//d := tools.NowYYYYMMDD()
	//dt := d + tools.NowHHmmss()
	//err := Mart.InsertFileHiLife(noti.Doc.ParentId, noti.Doc.EshopId, "關轉通知", dt, d, raw)
	//if err != nil {
	//	log.Error("SwitchNotification.Mart.InsertFileHiLife:", noti, raw)
	//}

	//state := "關轉店，須改店"
	//dirc := Enum.OrderShipSend
	//
	//dStr := tools.NowYYYYMMDD()
	//tStr := tools.NowHHmmss()
	//shipStatus :=  Enum.OrderShipReceiverStoreSwitch
	//if noti.Doc.StoreType == "2" {
	//	shipStatus = Enum.OrderShipSenderStoreSwitch
	//	dirc = Enum.OrderShipReturn
	//}

	//detail := generateRecord(dStr,tStr,dirc,state)
	//lg := generateLog(dStr,tStr,state, noti)
	// 更新關轉等待時間為2天
	//Mart.UpdateHiLifeSwitchTime(orderNo,tools.NowYYYYMMDDHHmmss(2))
	//Mart.UpdateOrderShippingStatus(ecOrderNo,shipStatus)

	//_ = Mart.WriteHiLifeShipStateByShipNo(orderNo, state,"1", detail,lg)
	//return nil
//}
