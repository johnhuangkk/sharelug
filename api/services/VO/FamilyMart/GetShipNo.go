package FamilyMart

import (
	"api/services/Enum"
	"api/services/entity"
	"api/services/util/tools"
	"encoding/xml"
	"github.com/spf13/viper"
	"log"
	"strconv"
)

/*
<?xml version="1.0" encoding="UTF-8"?>
<Doc>
    <ParentId>888</ParentId>
    <EshopId>0001</EshopId>
    <OPMode>1</OPMode>
    <ReceiverStoreId>001850</ReceiverStoreId>
    <ReturnStoreId></ReturnStoreId>
    <ServiceType>1</ServiceType>
    <EcOrderNo>1234567890</EcOrderNo>
    <OrderDate>2020-06-01</OrderDate>
    <OrderAmount>400</OrderAmount>
    <AgencyFee>400</AgencyFee>
    <ShipDate>2017-12-08</ShipDate>
    <SenderName>tests</SenderName>
    <SenderPhone>0000000830</SenderPhone>
    <ReceiverName>王小明</ReceiverName>
    <ReceiverPhone>0912345678</ReceiverPhone>
    <Products>
		<Product>
			<ProductId>A001</ProductId>
			<ProductName>AAAA</ProductName>
			<Quantity>10</Quantity>
			<Price/>
		</Product>
	</Products>
    <Remarks/>
</Doc>
*/

type OrderAddRequest struct {
	XMLName         xml.Name `xml:"Doc"`
	ParentId        string
	EshopId         string
	OPMode          string
	ReceiverStoreId string
	ReturnStoreId   string
	ServiceType     string // 1-取貨付款 3-取貨不付款
	EcOrderNo       string
	OrderDate       string // 取號日期
	OrderAmount     string
	AgencyFee       string
	ShipDate        string
	SenderName      string
	SenderPhone     string
	ReceiverName    string
	ReceiverPhone   string
	Products        struct {
		Objs []OrderAddRequestProduct
	}
	Remarks string
}
type OrderAddRequestProduct struct {
	XMLName     xml.Name `xml:"Product"`
	ProductId   string
	ProductName string
	Quantity    int
	Price       int
}

// 設定初始值
func (oA *OrderAddRequest) SetParams(orderData entity.OrderData, sellerData entity.MemberData) {

	var config = viper.GetStringMapString(`MartFamily901`)
	var ServiceType = `3`
	var AgencyFee = `0`

	// 付款方式為貨到付款時 需填入代收
	if orderData.PayWay == Enum.CvsPay {
		ServiceType = `1`
		AgencyFee = strconv.Itoa(int(orderData.TotalAmount))
	}

	oA.ParentId = config[`parentid`]
	oA.EshopId = config[`eshopid`]
	oA.OPMode = config[`opmode`]
	oA.ReceiverStoreId = orderData.ReceiverAddress
	oA.ReturnStoreId = ""
	oA.ServiceType = ServiceType
	oA.EcOrderNo = orderData.OrderId
	oA.OrderDate = tools.Now(`Ymd`)
	oA.OrderAmount = strconv.Itoa(int(orderData.TotalAmount))
	oA.AgencyFee = AgencyFee
	oA.ShipDate = tools.Now(`Ymd`)
	oA.SenderName = sellerData.SendName
	oA.SenderPhone = sellerData.Mphone
	oA.ReceiverName = orderData.ReceiverName
	oA.ReceiverPhone = orderData.ReceiverPhone
}

// 1-取貨付款 3-取貨不付款
func (oA *OrderAddRequest) GetServiceType() string {
	if oA.ServiceType == `1` {
		return oA.ServiceType
	}
	return `0`
}

type OrderAddResponse struct {
	OrderNo      string
	EcOrderNo    string
	ErrorCode    string
	ErrorMessage string
}

func (oA *OrderAddRequest) EncodeXML() (string, bool) {
	data, err := xml.Marshal(oA)
	if err != nil {
		log.Println("EncodeXML:", err)
		return "", false
	}

	return xml.Header + string(data), true
}
