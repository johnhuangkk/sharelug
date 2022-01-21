package InvoiceXml

//發票接收確認訊息
type InvoiceConfirm struct {
	InvoiceNumber     string `xml:"InvoiceNumber"`  	//發票號碼
	InvoiceDate       string `xml:"InvoiceDate"`		//發票開立日期
	BuyerId           string `xml:"BuyerId"`			//買方統一編號
	SellerId          string `xml:"SellerId"`			//賣方統一編號
	ReceiveDate       string `xml:"ReceiveDate"`		//發票接收日期
	ReceiveTime       string `xml:"ReceiveTime"`		//發票接收時間
	BuyerRemark       string `xml:"BuyerRemark"`		//買受人註記欄
	Remark            string `xml:"Remark"`				//備註
	BondedAreaConfirm string `xml:"BondedAreaConfirm"`	//買受人簽署適用零稅率註記
}
