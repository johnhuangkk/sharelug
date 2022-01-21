package InvoiceXml

import "encoding/xml"

//作廢折讓證明單訊息
type CancelAllowance struct {
	XMLName               xml.Name `xml:"CancelAllowance"`
	Text                  string   `xml:",chardata"`
	Xmlns                 string   `xml:"xmlns,attr"`
	Xsi                   string   `xml:"xmlns:xsi,attr"`
	SchemaLocation        string   `xml:"xsi:schemaLocation,attr"`
	CancelAllowanceNumber string   `xml:"CancelAllowanceNumber"` //作廢折讓證明單號碼
	AllowanceDate         string   `xml:"AllowanceDate"`         //折讓證明單日期
	BuyerId               string   `xml:"BuyerId"`               //買方統一編號
	SellerId              string   `xml:"SellerId"`              //賣方統一編號
	CancelDate            string   `xml:"CancelDate"`            //折讓證明單作廢日期
	CancelTime            string   `xml:"CancelTime"`            //折讓證明單作廢時間
	CancelReason          string   `xml:"CancelReason"`          //折讓證明單作廢原因
	Remark                string   `xml:"Remark"`                //備註
}
