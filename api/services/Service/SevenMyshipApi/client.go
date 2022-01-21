package SevenMyshipApi

import (
	"api/services/Enum"
	"api/services/dao/Cvs"
	sevenmyshipdao "api/services/dao/SevenMyshipDao"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"bytes"
	"encoding/xml"
	"fmt"
	"html/template"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/viper"
)

type ChargeOrderXml struct {
	XMLName     xml.Name            `xml:"OLTP"`
	Header      ChargeOrderHeader   `xml:"HEADER"`
	TotalCount  string              `xml:"AP>TotalCount"`
	TotalAmount string              `xml:"AP>TotalAmount"`
	Detail      []ChargeOrderDetail `xml:"AP>Detail"`
}
type ChargeOrderHeader struct {
	VER      string
	FROM     string
	TERMINO  string
	TO       string
	BUSINESS string
	DATE     string
	TIME     string
	STATCODE string
	STATDESC string
}
type ChargeOrderDetail struct {
	SequenceNo  string
	OL_OI_NO    string
	OL_Code_1   string
	OL_Code_2   string
	OL_Code_3   string
	OL_Amount   string
	Status      string
	Description string
	OL_Print    string
}
type DivTableString struct {
	PayDiv    template.HTML
	NonPayDiv template.HTML
}
type shipDataXml struct {
	eshopid      string `xml:"eshopid"`
	eshopsonid   string `xml:"eshopsonid"`
	orderno      string `xml:"orderno"`
	serviceType  string `xml:"service_type"`
	account      string `xml:"account"`
	paymentno    string `xml:"paymentno"`
	validationno string `xml:"validationno"`
	status       string `xml:"status"`
}

const (
	showType              = "21"
	defaultReturnStoreID  = "NNNNNN" // 退回店面使用預設原寄件店面
	defaultDeadlineTime   = "2359"
	paymentCpname         = "CheckNe"
	defaultDaishouAccount = "0"
	StatusFail            = "F"
)

type ShipInfoResult struct {
	Xml            xml.Name `xml:"C2C"`
	EndPoint       string   `xml:"xmlns,attr"`
	Description    string   `xml:"description"`
	Status         string   `xml:"status"`
	Paymentno      string   `xml:"paymentno"`
	Validationno   string   `xml:"validationno"`
	EshopID        string   `xml:"eshopid"`
	EshopsonID     string   `xml:"eshopsonid"`
	OrderNo        string   `xml:"orderno"`
	PayAmount      string   `xml:"payamount"`
	DaishouAccount string   `xml:"daishou_account"`
	Sender         string   `xml:"sender"`
	SenderPhone    string   `xml:"sender_phone"`
}

var eshopID string
var serviceType string

func FetchShipment() {
	FetchPackageSendByOl()
	FetchCPPSStatus()
	FetchCEINStatus()
	FetchCEDRStatus()
	FetchCESPStatus()
	FetchCERTStatus()
}
func CreateShipOrder(orderData entity.OrderData, sellerData entity.MemberData) (string, error) {

	deadlineDate := orderData.ShipExpire.Format("20060102")
	endpoint := viper.GetString(`MyShip.postPackageUrl`)
	if orderData.PayWay != Enum.CvsPay {
		eshopID = viper.GetString(`MyShip.eshopid`)
		serviceType = "7"
	} else {
		eshopID = viper.GetString(`MyShip.eshopidbyPay`)
		serviceType = "6"
	}
	var sName string
	var rName string
	if len([]rune(sellerData.SendName)) > 5 {
		sName = string([]rune(sellerData.SendName)[:4])
	} else {
		sName = sellerData.SendName
	}
	if len([]rune(orderData.ReceiverName)) > 5 {
		rName = string([]rune(orderData.ReceiverName)[:4])
	} else {
		rName = orderData.ReceiverName
	}
	data := url.Values{}
	data.Add("eshopid", eshopID)
	data.Add("eshopsonid", viper.GetString(`MyShip.eshopsonid`))
	data.Add("orderno", orderData.OrderId)
	data.Add("service_type", serviceType)
	data.Add("account", strconv.Itoa(int(math.Floor(orderData.SubTotal))))
	data.Add("payment_cpname", paymentCpname)
	data.Add("trade_ description", "")
	data.Add("cp_remark01", "")
	data.Add("cp_remark02", "")
	data.Add("cp_remark03", "")
	data.Add("deadlinedate", deadlineDate)
	data.Add("deadlinetime", defaultDeadlineTime)
	data.Add("show_type", showType)
	data.Add("daishou_account", defaultDaishouAccount)
	data.Add("sender", sName)
	data.Add("sender_phone", sellerData.Mphone)
	data.Add("receiver", rName)
	data.Add("receiver_phone", orderData.ReceiverPhone)
	data.Add("return_storeid", defaultReturnStoreID)
	data.Add("receiver_storeid", orderData.ReceiverAddress)

	client := &http.Client{}
	r, err := http.NewRequest("POST", endpoint, strings.NewReader(data.Encode())) // URL-encoded payload

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(r)
	if err != nil {
		log.Fatal(err.Error())
		return "", err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	var shipDatas ShipInfoResult

	err = xml.Unmarshal([]byte(body), &shipDatas)

	if err != nil || shipDatas == (ShipInfoResult{}) {
		log.Error("Seven CreateShipOrder Error: [%s]", err.Error())
		return "", err
	}

	if shipDatas.Status == StatusFail {
		return orderData.OrderId, fmt.Errorf(shipDatas.Description)
	}

	shipOrderNo := shipDatas.Paymentno + shipDatas.Validationno
	var cvsData entity.CvsShippingData
	cvsData.InitInsert(Enum.CVS_7_ELEVEN)
	// 建立托運資訊

	cvsData.ParentId = ""
	cvsData.EcOrderNo = shipDatas.OrderNo
	cvsData.ShipNo = shipOrderNo
	cvsData.SenderName = shipDatas.Sender
	cvsData.SenderPhone = shipDatas.SenderPhone
	cvsData.OriReceiverAddress = orderData.ReceiverAddress
	if serviceType == "6" {
		cvsData.ServiceType = "1"
	} else {
		cvsData.ServiceType = "0"
	}
	engine := database.GetMysqlEngine()
	defer engine.Close()
	err = Cvs.InsertCvsShippingData(engine, cvsData)
	if err != nil {
		log.Error("OK InsertCvsShippingData data Error: [%v]", data)
		log.Error("OK InsertCvsShippingData Error: [%v]", err.Error())
		return "", err
	}
	var shipMap entity.SevenShipMapData
	shipMap.OrderId = orderData.OrderId
	shipMap.PaymentNo = shipDatas.Paymentno
	shipMap.VerifyCode = shipDatas.Validationno
	shipMap.PayWay = orderData.PayWay
	shipMap.PaymentNoWithCode = shipDatas.Paymentno + shipDatas.Validationno

	sevenmyshipdao.InsertSevenShipMap(shipMap)
	return shipOrderNo, nil
}
func PrintShippingOrder(orders map[string][]string) ([]byte, error) {

	DivTableString := DivTableString{}

	if len(orders["CvsPay"]) > 0 {

		DivTableString.PayDiv = callPrintClient(orders["CvsPay"], viper.GetString(`MyShip.eshopidbyPay`), viper.GetString(`MyShip.printPayOrderPwd`))
	}
	if len(orders["NonCvsPay"]) > 0 {

		DivTableString.NonPayDiv = callPrintClient(orders["NonCvsPay"], viper.GetString(`MyShip.eshopid`), viper.GetString(`MyShip.printNonPayOrderPwd`))
	}

	var tmpl = template.Must(template.ParseFiles("views/order/sevenShipOrder.html"))

	var tpl bytes.Buffer
	if err := tmpl.Execute(&tpl, DivTableString); err != nil {
		log.Error("Execute Template Error", err)

	}
	return tpl.Bytes(), nil
}
func callPrintClient(shipNumbers []string, eid string, epwd string) template.HTML {
	endpoint := viper.GetString(`MyShip.printOrderUrl`)
	data := url.Values{}
	data.Add("eshopid", eid)
	data.Add("PinCodes", strings.Join(shipNumbers, ","))
	data.Add("BackTag", "")
	data.Add("tempvar", "")
	data.Add("member_pwd", epwd)
	r, err := http.NewRequest("POST", endpoint, strings.NewReader(data.Encode())) // URL-encoded payload

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer res.Body.Close()

	var orderTableDiv string
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err == nil {
		doc.Find("div").Each(func(i int, s *goquery.Selection) {

			id, _ := s.Attr("id")
			if id == "Panel1" {
				orderTableDiv, _ = s.Html()
				return
			}
		})
	}

	return template.HTML(replaceSevenDomain(orderTableDiv))
}
func replaceSevenDomain(s string) string {
	fmt.Println(s)
	s = strings.Replace(s, "QRCode.ashx?", viper.GetString(`MyShip.sevenCodeDomain`)+"QRCode.ashx?", -1)
	s = strings.Replace(s, "BarCode.ashx?", viper.GetString(`MyShip.sevenCodeDomain`)+"BarCode.ashx?", -1)
	return s
}

//updatePackageStore by post xml form data
func updatePackageStore() {

}

//CreateChargeOrderRecordByPos from seven request when buyer fishish order proccess at store
func CreateChargeOrderRecordByPos(s []byte) ([]byte, []entity.SevenChargeOrderData) {
	v := ChargeOrderXml{}

	err := xml.Unmarshal(s, &v)
	if err != nil {
		log.Error("Seven Charge Order xml decode error")
	}
	datas := chargeOrderXmltoRecordData(v)

	_, err = sevenmyshipdao.InsertChargeOrderRecords(datas)

	if err != nil {
		log.Error("Seven Charge Order xml decode error")
	}
	v.Header.FROM, v.Header.TO = v.Header.TO, v.Header.FROM

	xml, err := xml.Marshal(v)
	if err != nil {
		log.Error("Seven Charge Order xml encode error")
	}
	return xml, datas
}
func chargeOrderXmltoRecordData(xml ChargeOrderXml) []entity.SevenChargeOrderData {
	datas := []entity.SevenChargeOrderData{}

	for _, detail := range xml.Detail {
		entity := entity.SevenChargeOrderData{
			From:        xml.Header.FROM,
			To:          xml.Header.TO,
			TermiNo:     xml.Header.TERMINO,
			Date:        xml.Header.DATE,
			Time:        xml.Header.TIME,
			StatCode:    xml.Header.STATCODE,
			StatDesc:    xml.Header.STATDESC,
			SequenceNo:  detail.SequenceNo,
			OLOiNo:      detail.OL_OI_NO,
			OLCode1:     detail.OL_Code_1,
			OLCode2:     detail.OL_Code_2,
			OLCode3:     detail.OL_Code_3,
			OLAmount:    detail.OL_Amount,
			Status:      detail.Status,
			Description: detail.Description,
			OLPrint:     detail.OL_Print}
		datas = append(datas, entity)
	}

	return datas
}

//retrieveStroreChangeRecord by ftp CCS file record
//If store close or storeId change
func retrieveStoreChangeRecord() {

}

//retrievePackageDeliveredToStoreRecord by ftp CEIN file record
//If package at sender store or receiver store
func retrievePackageDeliveredToStoreRecord() {

}

//retrievePackageReturnPackageRecord by ftp CERT file record
//The notification of package need to be return , receiver did not take the package in deadline
func retrieveReturnPackageRecord() {

}

//retrieiveReturnPackageToCenterRecord by ftp CEDR file record
//The return package has been delivered to center
func retrieveReturnPackageToCenterRecord() {

}

//retrieveReturnDailyPayCollectionRecord by ftp ACC file record
func retrieveDailyPayCollectionRecord() {

}

//retrieveCloseCaseRecord by ftp OL file record
func retrieveCloseCaseRecord() {

}

//retrievePackageAtStoreRecord by ftp CESP file record
//If sender delivered package to store or receiver return package to store
func retrievePackageAtStoreRecord() {

}

//batchCreateOrder by post xml data
//The order qauntity max <= 200
func batchCreateOrderForm() {

}
