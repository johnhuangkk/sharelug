package InvoiceXml

//作廢發票接收確認訊息
type CancelInvoiceConfirm struct {
	CancelInvoiceNumber string `xml:"CancelInvoiceNumber"` 	//作廢發票號碼
	InvoiceDate         string `xml:"InvoiceDate"`			//發票開立日期
	BuyerId             string `xml:"BuyerId"`				//買方統一編號
	SellerId            string `xml:"SellerId"`				//賣方統一編號
	CancelDate          string `xml:"CancelDate"`			//發票作廢日期
	CancelTime          string `xml:"CancelTime"`			//發票作廢時間
	Remark              string `xml:"Remark"`				//備註
}
