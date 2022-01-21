package InvoiceXml

import "encoding/xml"

type CancelInvoiceC0501 struct {
	XMLName                 xml.Name `xml:"CancelInvoice"`
	Text                    string   `xml:",chardata"`
	Xmlns                   string   `xml:"xmlns,attr"`
	Xsi                     string   `xml:"xmlns:xsi,attr"`
	SchemaLocation          string   `xml:"xsi:schemaLocation,attr"`
	CancelInvoiceNumber     string   `xml:"CancelInvoiceNumber"`     //作廢發票號碼
	InvoiceDate             string   `xml:"InvoiceDate"`             //發票日期
	BuyerId                 string   `xml:"BuyerId"`                 //買方統一編號
	SellerId                string   `xml:"SellerId"`                //賣方統一編號
	CancelDate              string   `xml:"CancelDate"`              //作廢時間
	CancelTime              string   `xml:"CancelTime"`              //作廢時間
	CancelReason            string   `xml:"CancelReason"`            //作廢原因
	//ReturnTaxDocumentNumber string   `xml:"ReturnTaxDocumentNumber"` //專案作廢核淮文號
	Remark                  string   `xml:"Remark"`                  //備註
}
