package InvoiceXml

import "encoding/xml"

//開立發票訊息規格
type InvoiceA0101 struct {
	XMLName        xml.Name       `xml:"Invoice"`
	Text           string         `xml:",chardata"`
	Xmlns          string         `xml:"xmlns,attr"`
	Xsi            string         `xml:"xsi,attr"`
	SchemaLocation string         `xml:"schemaLocation,attr"`
	Main           InvoiceMain    `xml:"Main"`
	Details        InvoiceDetails `xml:"Details"`
	Amount         InvoiceAmount  `xml:"Amount"`
}

type InvoiceMain struct {
	InvoiceNumber        string        `xml:"InvoiceNumber"`        //發票號碼
	InvoiceDate          string        `xml:"InvoiceDate"`          //發票開立日期
	InvoiceTime          string        `xml:"InvoiceTime"`          //發票開立時間
	Seller               InvoiceSeller `xml:"Seller"`               //賣方資訊
	Buyer                InvoiceBuyer  `xml:"Buyer"`                //買方資訊
	CheckNumber          string        `xml:"CheckNumber"`          //發票檢查碼
	BuyerRemark          string        `xml:"BuyerRemark"`          //買受人註記欄
	MainRemark           string        `xml:"MainRemark"`           //總備註
	CustomsClearanceMark string        `xml:"CustomsClearanceMark"` //通關方式註記
	Category             string        `xml:"Category"`             //沖帳別
	RelateNumber         string        `xml:"RelateNumber"`         //相關號碼
	InvoiceType          string        `xml:"InvoiceType"`          //發票類別
	//GroupMark            string        `xml:"GroupMark"`            //彙開註記
	DonateMark           string        `xml:"DonateMark"`           //捐贈註記
	Attachment           string        `xml:"Attachment"`           //證明附件
}

type InvoiceSeller struct {
	Identifier      string `xml:"Identifier"`      //識別碼
	Name            string `xml:"Name"`            //名稱
	Address         string `xml:"Address"`         //地址
	PersonInCharge  string `xml:"PersonInCharge"`  //負責人姓名
	TelephoneNumber string `xml:"TelephoneNumber"` //電話號碼
	FacsimileNumber string `xml:"FacsimileNumber"` //傳真號碼
	EmailAddress    string `xml:"EmailAddress"`    //電子郵件地址
	CustomerNumber  string `xml:"CustomerNumber"`  //客戶編號
	RoleRemark      string `xml:"RoleRemark"`      //營業人角色註記
}

type InvoiceBuyer struct {
	Identifier      string `xml:"Identifier"`      //識別碼
	Name            string `xml:"Name"`            //名稱
	//Address         string `xml:"Address"`         //地址
	//PersonInCharge  string `xml:"PersonInCharge"`  //負責人姓名
	//TelephoneNumber string `xml:"TelephoneNumber"` //電話號碼
	//FacsimileNumber string `xml:"FacsimileNumber"` //傳真號碼
	//EmailAddress    string `xml:"EmailAddress"`    //電子郵件地址
	//CustomerNumber  string `xml:"CustomerNumber"`  //客戶編號
	//RoleRemark      string `xml:"RoleRemark"`      //營業人角色註記
}

type InvoiceDetails struct {
	ProductItem []ProductItem `xml:"ProductItem"` //商品項目資料
}

type ProductItem struct {
	Description    string `xml:"Description"`    //品名
	Quantity       int64  `xml:"Quantity"`       //數量
	Unit           string `xml:"Unit"`           //單位
	UnitPrice      int64  `xml:"UnitPrice"`      //單價
	Amount         int64  `xml:"Amount"`         //金額
	SequenceNumber int64  `xml:"SequenceNumber"` //明細排列序號
	Remark         string `xml:"Remark"`         //單一欄位備註
	RelateNumber   string `xml:"RelateNumber"`   //相關號碼
}

type InvoiceAmount struct {
	SalesAmount            int64  `xml:"SalesAmount"`            //銷售額合計
	TaxType                string `xml:"TaxType"`                //課稅別
	TaxRate                int64  `xml:"TaxRate"`                //稅率
	TaxAmount              int64  `xml:"TaxAmount"`              //營業稅額
	TotalAmount            int64  `xml:"TotalAmount"`            //總計
	DiscountAmount         int64  `xml:"DiscountAmount"`         //折扣金額
	OriginalCurrencyAmount int64  `xml:"OriginalCurrencyAmount"` //原幣金額
	ExchangeRate           int64  `xml:"ExchangeRate"`           //匯率
	Currency               string `xml:"Currency"`               //幣別
}
