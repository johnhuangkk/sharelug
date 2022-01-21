package InvoiceXml

//開立折讓證明／通知單接收確認訊息
type AllowanceConfirm struct {
	AllowanceNumber string `xml:"AllowanceNumber"` 	//折讓證明單號碼
	AllowanceDate   string `xml:"AllowanceDate"`	//折讓證明單日期
	BuyerId 		string `xml:"BuyerId"`			//買方統一編號
	SellerId 		string `xml:"SellerId"`			//賣方統一編號
	ReceiveDate		string `xml:"ReceiveDate"`		//折讓證明單接收日
	ReceiveTime 	string `xml:"ReceiveTime"`		//折讓證明單接收時間
	Remark			string `xml:"Remark"`			//折讓種類
	AllowanceType	string `xml:"AllowanceType"`	//備註
}
