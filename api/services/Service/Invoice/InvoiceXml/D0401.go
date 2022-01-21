package InvoiceXml

type AllowanceD0401 struct {

	Main    AllowanceD0401Main `xml:"Main"`
	Details Details            `xml:"Details"`
	Amount  Amount             `xml:"Amount"`
}

type AllowanceD0401Main struct {
	AllowanceNumber string        `xml:"AllowanceNumber"` //折讓證明單號碼
	AllowanceDate   string        `xml:"AllowanceDate"`   //折讓證明單日期
	Seller          InvoiceSeller `xml:"Seller"`          //賣方資訊
	Buyer           InvoiceBuyer  `xml:"Buyer"`           //買方資訊
	AllowanceType   string        `xml:"AllowanceType"`   //折讓種類
	Attachment      string        `xml:"Attachment"`      //證明附件
}

type Details struct {
	ProductItem ProductItem `xml:"ProductItem"` //商品項目資料
}

type Amount struct {
	TaxAmount   string `xml:"TaxAmount"`   //營業稅額合計
	TotalAmount string `xml:"TotalAmount"` //金額合計(不含稅之進貨額合計)
}
