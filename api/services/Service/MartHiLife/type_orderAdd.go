package MartHiLife

import (
	"api/services/Enum"
	"api/services/entity"
	"api/services/util/tools"
	"encoding/xml"
	"github.com/spf13/viper"
	"strconv"
)

//<?xml version="1.0" encoding="utf-8"?>
//<Doc>
//	<ParentId>124</ParentId>
//	<EshopId>901</EshopId>
//	<EcDcNo>D11</EcDcNo>
//	<EcCvs>HILIFEC2C</EcCvs>
//	<ReceiverStoreId>3750</ReceiverStoreId>
//	<ReturnStoreId></ReturnStoreId>
//	<ServiceType>3</ServiceType>
//	<VdrOrderNo>000012345678</VdrOrderNo>
//	<OrderDate>2020-09-09</OrderDate>
//	<OrderAmount>500</OrderAmount>
//	<AgencyFee>0</AgencyFee>
//	<Sercode>865</Sercode>
//	<ShipDate>2020-09-09</ShipDate>
//	<SenderName>程又青</SenderName>
//	<SenderPhone>0922123456</SenderPhone>
//	<ReceiverName>李大仁</ReceiverName>
//	<ReceiverPhone>0955123456</ReceiverPhone>
//	<Remarks></Remarks>
//	<ChkMac>4B1223532ACF6FF1E50B52D23E03A558</ChkMac>
//</Doc>`

type OrderAddRequest struct {
	XMLName         xml.Name `xml:"Doc"`
	ParentId        string `xml:"ParentId"`
	EshopId         string `xml:"EshopId"`
	EcDcNo          string `xml:"EcDcNo"`
	EcCvs           string `xml:"EcCvs"`
	ReceiverStoreId string `xml:"ReceiverStoreId"`
	ReturnStoreId   string `xml:"ReturnStoreId"`
	ServiceType     string `xml:"ServiceType"`
	VdrOrderNo      string `xml:"VdrOrderNo"`
	OrderDate       string `xml:"OrderDate"`
	OrderAmount     string `xml:"OrderAmount"`
	AgencyFee       string `xml:"AgencyFee"`
	Sercode         string `xml:"Sercode"`
	ShipDate        string `xml:"ShipDate"`
	SenderName      string `xml:"SenderName"`
	SenderPhone     string `xml:"SenderPhone"`
	ReceiverName    string `xml:"ReceiverName"`
	ReceiverPhone   string `xml:"ReceiverPhone"`
	Remarks         string `xml:"Remarks"`
	ChkMac          string `xml:"ChkMac"`
}

func (receiver *OrderAddRequest) SetAddShipNoParams(orderData entity.OrderData, sellerData entity.MemberData) {
	config := viper.GetStringMapString(`MartHiLife`)

	var ServiceType = `3`
	var AgencyFee = `0`

	// 付款方式為貨到付款時 需填入代收
	if orderData.PayWay == Enum.CvsPay {
		ServiceType = `1`
		AgencyFee = strconv.Itoa(int(orderData.TotalAmount))
	}

	receiver.ParentId = config[`parentid`]
	receiver.EshopId = config[`eshopid`]
	receiver.EcCvs = config[`eccvs`]
	receiver.EcDcNo = config[`ecdcno`]
	receiver.Sercode = config[`sercode`]
	receiver.ReceiverStoreId = orderData.ReceiverAddress
	receiver.ReturnStoreId = ``
	receiver.ServiceType = ServiceType
	receiver.VdrOrderNo = orderData.OrderId
	receiver.OrderDate = tools.Now(`Ymd`)
	receiver.OrderAmount = strconv.Itoa(int(orderData.TotalAmount))
	receiver.AgencyFee = AgencyFee
	receiver.ShipDate = tools.Now(`Ymd`)
	receiver.SenderName = sellerData.SendName
	receiver.SenderPhone = sellerData.Mphone
	receiver.ReceiverName = orderData.ReceiverName
	receiver.ReceiverPhone = orderData.ReceiverPhone

	receiver.ChkMac = GenerateCheckSum1(receiver.ParentId,
		receiver.EshopId,
		receiver.EcDcNo,
		receiver.EcCvs,
		receiver.VdrOrderNo,
		config[`hashkey`],
		config[`hashiv`])
}

// 1-取貨付款 3-取貨不付款
func (receiver *OrderAddRequest) GetServiceType() string {
	if receiver.ServiceType == `1` {
		return receiver.ServiceType
	}
	return `0`
}

func (receiver*OrderAddRequest) EncodeXML() []byte {
	d,_ := xml.Marshal(receiver)
	hd := []byte(xml.Header)
	hd = append(hd,d...)
	return hd
}

func (receiver*OrderAddRequest) CheckMac()  {

	config := viper.GetStringMapString(`MartHiLife`)

	receiver.ChkMac = GenerateCheckSum1(receiver.ParentId,
		receiver.EshopId,
		receiver.EcDcNo,
		receiver.EcCvs,
		receiver.VdrOrderNo,
		config[`hashkey`],
		config[`hashiv`])
}

//<?xml version="1.0" encoding="utf-8"?>
//<Doc>
//	<VdrOrderNo>20201120001H</VdrOrderNo>
//	<OrderNo>20RKB40416476</OrderNo>
//	<ErrorCode>000</ErrorCode>
//	<ErrorMessage>成功</ErrorMessage>
//</Doc>

type OrderAddResp struct {
	XMLName         xml.Name `xml:"Doc"`
	VdrOrderNo      string   `xml:"VdrOrderNo"`
	OrderNo         string   `xml:"OrderNo"`
	ErrorCode       string   `xml:"ErrorCode"`
	ErrorMessage    string   `xml:"ErrorMessage"`
}

func (me *OrderAddResp)	DecodeXML(data []byte) (err error) {
	return xml.Unmarshal(data,me)
}

