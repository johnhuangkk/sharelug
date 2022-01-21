package InvoiceXml

import "encoding/xml"

type InvoiceC0401 struct {
	XMLName        xml.Name            `xml:"Invoice"`
	Text           string              `xml:",chardata"`
	Xmlns          string              `xml:"xmlns,attr"`
	Xsi            string              `xml:"xmlns:xsi,attr"`
	SchemaLocation string              `xml:"xsi:schemaLocation,attr"`
	Main           InvoiceC0401Main    `xml:"Main"`
	Details        InvoiceC0401Details `xml:"Details"`
	Amount         InvoiceC0401Amount  `xml:"Amount"`
}

type InvoiceC0401Main struct {
	InvoiceNumber        string        `xml:"InvoiceNumber"`        //發票號碼
	InvoiceDate          string        `xml:"InvoiceDate"`          //發票日期
	InvoiceTime          string        `xml:"InvoiceTime"`          //發票時間
	Seller               InvoiceSeller `xml:"Seller"`               //賣方資訊
	Buyer                InvoiceBuyer  `xml:"Buyer"`                //買方資訊
	CheckNumber          string        `xml:"CheckNumber"`          //發票檢查碼
	BuyerRemark          int64         `xml:"BuyerRemark"`          //買受人註記欄
	MainRemark           string        `xml:"MainRemark"`           //總備註
	CustomsClearanceMark int64         `xml:"CustomsClearanceMark"` //通關方式註記
	Category             string        `xml:"Category"`             //沖帳別
	RelateNumber         string        `xml:"RelateNumber"`         //相關號碼
	//GroupMark            string        `xml:"GroupMark"`            //發票類別
	InvoiceType          string        `xml:"InvoiceType"`          //彙開註記
	DonateMark           int64         `xml:"DonateMark"`           //捐贈註記
	CarrierType          string        `xml:"CarrierType"`          //載具類別號碼
	CarrierId1           string        `xml:"CarrierId1"`           //載具顯碼ID
	CarrierId2           string        `xml:"CarrierId2"`           //載具顯碼ID
	PrintMark            string        `xml:"PrintMark"`            //電子發票證明聯已列印註記
	RandomNumber         string        `xml:"RandomNumber"`         //發票防偽隨機碼
	NPOBAN               string        `xml:"NPOBAN"`               //發票捐贈對象
}

type InvoiceC0401Details struct {
	ProductItem []ProductItem `xml:"ProductItem"`
}

type InvoiceC0401Amount struct {
	SalesAmount            float64 `xml:"SalesAmount"`            //應稅銷售額合計
	FreeTaxSalesAmount     int64   `xml:"FreeTaxSalesAmount"`     //免稅銷售額合計
	ZeroTaxSalesAmount     int64   `xml:"ZeroTaxSalesAmount"`     //零稅率銷售額合計
	TaxType                string  `xml:"TaxType"`                //課稅別
	TaxRate                float64 `xml:"TaxRate"`                //稅率
	TaxAmount              float64 `xml:"TaxAmount"`              //營業稅額
	TotalAmount            int64   `xml:"TotalAmount"`            //總計
	DiscountAmount         int64   `xml:"DiscountAmount"`         //折扣金額
	OriginalCurrencyAmount int64   `xml:"OriginalCurrencyAmount"` //原幣金額
	ExchangeRate           int64   `xml:"ExchangeRate"`           //匯率
	//Currency               string  `xml:"Currency"`               //幣別
}
