package InvoiceXml

import "encoding/xml"

//開立折讓證明單｜傳送折讓證明單通知
type Allowance struct {
	XMLName        xml.Name         `xml:"Allowance"`
	Text           string           `xml:",chardata"`
	Xmlns          string           `xml:"xmlns,attr"`
	Xsi            string           `xml:"xmlns:xsi,attr"`
	SchemaLocation string           `xml:"xsi:schemaLocation,attr"`
	Main           AllowanceMain    `xml:"Main"`
	Details        AllowanceDetails `xml:"Details"`
	Amount         AllowanceAmount  `xml:"Amount"`
}

type AllowanceMain struct {
	AllowanceNumber string        `xml:"AllowanceNumber"` //折讓證明單號碼
	AllowanceDate   string        `xml:"AllowanceDate"`   //折讓證明單開立日期
	Seller          InvoiceSeller `xml:"Seller"`          //賣方資訊
	Buyer           InvoiceBuyer  `xml:"Buyer"`           //買方資訊
	AllowanceType   string        `xml:"AllowanceType"`   //折讓種類
	Attachment      string        `xml:"Attachment"`      //證明文件
}

type AllowanceDetails struct {
	ProductItem []AllowanceProductItem `xml:"ProductItem"`
}

type AllowanceProductItem struct {
	OriginalInvoiceDate     string `xml:"OriginalInvoiceDate"`     //原發票日期
	OriginalInvoiceNumber   string `xml:"OriginalInvoiceNumber"`   //原發票號碼
	OriginalSequenceNumber  string `xml:"OriginalSequenceNumber"`  //原明細排列序號
	OriginalDescription     string `xml:"OriginalDescription"`     //原品名
	Quantity                int64  `xml:"Quantity"`                //數量
	Unit                    string `xml:"Unit"`                    //單位
	UnitPrice               int64  `xml:"UnitPrice"`               //單價
	Amount                  int64  `xml:"Amount"`                  //金額
	Tax                     int64  `xml:"Tax"`                     //營業稅額
	AllowanceSequenceNumber string `xml:"AllowanceSequenceNumber"` //折讓證明單明細排列序號
	TaxType                 string `xml:"TaxType"`                 //課稅別
}

type AllowanceAmount struct {
	TaxAmount   int64 `xml:"TaxAmount"`   //營業稅額合計
	TotalAmount int64 `xml:"TotalAmount"` //金額合計(不含稅之進貨額合計)
}
