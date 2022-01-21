package InvoiceXml

import "encoding/xml"

type VoidInvoiceC0701 struct {
	XMLName           xml.Name `xml:"VoidInvoice"`
	Text              string   `xml:",chardata"`
	Xmlns             string   `xml:"xmlns,attr"`
	Xsi               string   `xml:"xmlns:xsi,attr"`
	SchemaLocation    string   `xml:"xsi:schemaLocation,attr"`
	VoidInvoiceNumber string   `xml:"VoidInvoiceNumber"` //註銷發票號碼
	InvoiceDate       string   `xml:"InvoiceDate"`       //發票日期
	BuyerId           string   `xml:"BuyerId"`           //買方統一編號
	SellerId          string   `xml:"SellerId"`          //賣方統一編號
	VoidDate          string   `xml:"VoidDate"`          //註銷日期
	VoidTime          string   `xml:"VoidTime"`          //註銷時間
	VoidReason        string   `xml:"VoidReason"`        //註銷原因
	Remark            string   `xml:"Remark"`            //備註
}
