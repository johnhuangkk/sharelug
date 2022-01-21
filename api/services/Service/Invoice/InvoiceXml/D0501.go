package InvoiceXml

type CancelAllowanceD0501 struct {
	CancelAllowanceNumber string `xml:"CancelAllowanceNumber"` //作廢折讓證明單號碼
	AllowanceDate         string `xml:"AllowanceDate"`			//折讓證明單日期
	BuyerId               string `xml:"BuyerId"`				//買方統一編號
	SellerId              string `xml:"SellerId"`				//賣方統一編號
	CancelDate            string `xml:"CancelDate"`				//折讓證明單作廢日期
	CancelTime            string `xml:"CancelTime"`				//折讓證明單作廢時間
	CancelReason          string `xml:"CancelReason"`			//折讓證明單作廢原因
	Remark                string `xml:"Remark"`					//備註
}
