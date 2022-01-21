package OKApi

import (
	"api/services/Enum"
	om "api/services/Service/OKMart"
	"api/services/dao/Cvs"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"fmt"
	"strconv"
)

// 取號
func GetShippingOrderNo(engine *database.MysqlSession, orderData entity.OrderData, sellerData entity.MemberData) (orderNo string, err error) {
	var client om.Client
	client.GetClient()

	result, err := client.OrderAdd(orderData, sellerData)
	if err != nil {
		return "", err
	}

	if result.Body.Order.ErrCode != "000" {
		return "", fmt.Errorf("[" + result.Body.Order.ErrCode + "]" + result.Body.Order.ErrDesc)
	}

	shipNo := result.Body.Order.OdNo

	// 建立托運資訊
	var data entity.CvsShippingData
	data.InitInsert(Enum.CVS_OK_MART)

	data.ParentId = `OK`
	data.EcOrderNo = orderData.OrderId
	data.ShipNo = shipNo
	data.ServiceType = `0`
	if orderData.PayWay == Enum.CvsPay {
		data.ServiceType = `1`
	}
	data.SenderName = sellerData.SendName
	data.SenderPhone = sellerData.Mphone
	data.OriReceiverAddress = orderData.ReceiverAddress

	err = Cvs.InsertCvsShippingData(engine, data)

	if err != nil {
		log.Error("OK InsertCvsShippingData data Error: [%v]", data)
		log.Error("OK InsertCvsShippingData Error: [%v]", err.Error())
		return "", err
	}

	return shipNo, nil
}

// OK 印出面單
func PrintShippingOrderX(cvsData []entity.CvsShippingData) (data []byte, err error) {

	log.Info(`PrintShippingOrderX cvsData`, cvsData)

	var client om.Client
	client.GetClient()

	if len(cvsData) == 0 {
		return data, fmt.Errorf("無寄件資料")
	}

	var formData  []string

	for _, c := range cvsData {
		formData = append(formData, c.ShipNo + ":TOK:" + c.SenderPhone[len(c.SenderPhone)-3:])
	}

	return client.OrderPrintX(formData)
}


func PrintShippingOrder(shipNo []string) (data []byte, err error) {

	var client om.Client
	client.GetClient()

	var objs []entity.MartOkShippingData
	for _, v := range shipNo {
		obj := entity.MartOkShippingData{ShipNo: v}

		has, err := database.GetMysqlEngineGroup().Get(&obj)
		if err != nil {
			return nil, err
		}
		if has {
			objs = append(objs, obj)
		}
	}

	if len(shipNo) != len(objs) {
		log.Error("PrintShippingOrder:多筆單號列印失敗")
		return nil, fmt.Errorf("多筆單號列印失敗")
	}

	var shipNos []string
	var tels []string

	for _, obj := range objs {
		leng := len(obj.SrPhone)
		if leng < 3 {
			return nil, fmt.Errorf("Shipment order not correct")
		}
		// 寄件人手機末三碼
		phoneSubfix := obj.SrPhone[leng-3:]
		shipNos = append(shipNos, obj.ShipNo)
		tels = append(tels, phoneSubfix)
	}

	return client.OrderPrint(shipNos, tels)
}

// 閉轉店
func SwitchStore(orderData entity.OrderData, cvsShippingData entity.CvsShippingData, newStoreId string) (err error) {

	var client om.Client
	client.GetClient()

	obj := entity.MartOkShippingData{EcOrderNo: cvsShippingData.EcOrderNo, ShipNo: cvsShippingData.ShipNo}
	result, err := database.Mysql().Get(&obj)
	if err != nil {
		return err
	}

	if !result {
		return fmt.Errorf("%s", "Shipment not found")
	}

	// 預設正向
	ecNo := client.EcNo4
	CUTKNM := orderData.ReceiverName
	CUTKTL := orderData.ReceiverPhone[len(orderData.ReceiverPhone)-3 : len(orderData.ReceiverPhone)]

	// 逆向
	if cvsShippingData.FlowType == `R` {
		ecNo = client.EcNo5
		CUTKNM = cvsShippingData.SenderName
		CUTKTL = cvsShippingData.SenderPhone[len(cvsShippingData.SenderPhone)-3 : len(obj.RrPhone)]
	}

	newOlder, err := client.OrderReSend(
		ecNo,
		cvsShippingData.ShipNo,
		newStoreId,
		CUTKNM, CUTKTL,
		strconv.Itoa(int(orderData.TotalAmount)),
		cvsShippingData.ServiceType)

	if err != nil {
		log.Error(`OK OrderReSend Error [%v]`, err.Error())
		return fmt.Errorf("閉轉店變更失敗")
	}

	log.Info(`OK OrderReSend`, newOlder)

	return nil
}
