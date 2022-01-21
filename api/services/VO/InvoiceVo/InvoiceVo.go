package InvoiceVo

import (
	"api/services/VO/Request"
	"fmt"
	"github.com/spf13/viper"
	"strconv"
	"strings"
	"time"
)

type BindCarrierVo struct {
	Token string `json:"Token"`
	Ban   string `json:"Ban"`
}

type InvoiceResponse struct {
	Track  string `json:"Track"`
	Number string `json:"Number"`
	Type   string `json:"Type"`
	Year   string `json:"Year"`
	Month  string `json:"Month"`
}

type InvoiceVo struct {
	OrderId string
	Details []Details
	UserId  string
	Amount  int64

	Carrier Request.CarrierRequest
}

type Buyer struct {
	BuyerName  string
	Identifier string
}

type Details struct {
	ProductName   string `json:"ProductName"`
	Quantity      int64  `json:"Quantity"`
	ProductPrice  int64  `json:"ProductPrice"`
	ProductAmount int64  `json:"ProductAmount"`
	Sequence      int64  `json:"Sequence"`
}

type AllowanceVo struct {
	Identifier string          `json:"Identifier"`
	Buyer      string          `json:"Buyer"`
	Products   []ProductItemVo `json:"Products"`
}

type ProductItemVo struct {
	OriginalInvoiceDate     string `json:"OriginalInvoiceDate"`     //原發票日期
	OriginalInvoiceNumber   string `json:"OriginalInvoiceNumber"`   //原發票號碼
	OriginalSequenceNumber  string `json:"OriginalSequenceNumber"`  //原明細排列序號
	OriginalDescription     string `json:"OriginalDescription"`     //原品名
	Quantity                int64  `json:"Quantity"`                //數量
	Unit                    string `json:"Unit"`                    //單位
	UnitPrice               int64  `json:"UnitPrice"`               //單價
	Amount                  int64  `json:"Amount"`                  //金額
	Tax                     int64  `json:"Tax"`                     //營業稅額
	AllowanceSequenceNumber string `json:"AllowanceSequenceNumber"` //折讓證明單明細排列序號
	TaxType                 string `json:"TaxType"`                 //課稅別
}

type VerifyBindCarrierRequest struct {
	Token string `json:"token"`
	Nonce string `json:"nonce"`
}

type VerifyBindCarrierResponse struct {
	TokenFlag string `json:"token_flag"`
	Nonce     string `json:"nonce"`
	ErrMsg    string `json:"err_msg"`
}

type QRCodeInvVo struct {
	InvoiceNumber       string     //發票字軌號碼共10碼
	InvoiceDate         string     //發票開立年月日(中華民國年份月份日期)共7碼
	InvoiceTime         string     //發票開立時間 (24小時制)共6碼 時分秒
	RandomNumber        string     //4碼隨機碼
	SalesAmount         string     //以整數方式載入銷售額 (未稅)，若無法分離稅項則記載為0
	TaxAmount           string     //以整數方式載入稅額，若無法分離稅項則記載為0
	TotalAmount         string     //整數方式載入總計金額(含稅)
	BuyerIdentifier     string     //買受人統一編號，若買受人為一般消費者，請填入 00000000 8位字串
	RepresentIdentifier string     //代表店統一編號，電子發票證明聯二維條碼規格已不使用代表店，請填入00000000 8位字串
	SellerIdentifier    string     //銷售店統一編號
	BusinessIdentifier  string     //總機構統一編號，如無總機構請填入銷售店統一編號
	ProductArray        []ItemList `json:"ItemList"`
	AESKey              string     //加解密金鑰(QR種子密碼)
}

func (m *InvoiceDetailVo) GetQRCodeInvVo() QRCodeInvVo {
	t, _ := time.Parse("2006/01/02 15:04:05", m.CreateTime)
	var resp QRCodeInvVo
	resp.InvoiceNumber = strings.Replace(m.InvoiceNumber, "-", "", -1)
	resp.InvoiceDate = fmt.Sprintf("%s%s", m.InvoiceYear, t.Format("0102"))
	resp.InvoiceTime = t.Format("1504")
	resp.RandomNumber = m.InvoiceRandom
	resp.SalesAmount = fmt.Sprintf("%0*s", 8, strconv.FormatInt(m.Amount, 16))
	resp.TaxAmount = fmt.Sprintf("%0*s", 8, strconv.FormatInt(0, 16))
	resp.TotalAmount = fmt.Sprintf("%0*s", 8, strconv.FormatInt(m.Amount, 16))
	resp.BuyerIdentifier = "00000000"
	resp.RepresentIdentifier = "00000000"
	resp.SellerIdentifier = m.SellerBan
	resp.BusinessIdentifier = m.SellerBan
	resp.ProductArray = m.ItemList
	config := viper.GetStringMapString("INVOICE")
	resp.AESKey = config["aeskey"]
	return resp
}

type InvoiceDetailVo struct {
	InvoiceStatus     string     `json:"InvoiceStatus"`     //已捐贈
	InvoiceStatusText string     `json:"InvoiceStatusText"` //已捐贈
	InvoiceYear       string     `json:"InvoiceYear"`       //發票年份  109年03-04月
	InvoiceMonth      string     `json:"InvoiceMonth"`      //發票月份  109年03-04月
	CreateTime        string     `json:"InvoiceDate"`       //發票開立時間 2020/03/14 20:16:43
	InvoiceNumber     string     `json:"InvoiceNumber"`     //發票號碼 AA-00000000
	InvoiceRandom     string     `json:"InvoiceRandom"`     //隨機碼
	SellerBan         string     `json:"SellerBan"`         //賣方統編
	BuyerBan          string     `json:"BuyerBan"`          //買受人統編
	BuyerName         string     `json:"BuyerName"`         //買受人名稱
	ItemList          []ItemList `json:"ItemList"`
	Sales             int64      `json:"Sales"`  //銷售額
	Tax               int64      `json:"Tax"`    //稅額
	Amount            int64      `json:"Amount"` //總計
	QRCode1           string     `json:"QRCode1"`
	QRCode2           string     `json:"QRCode2"`
}

type ItemList struct {
	ProductName   string `json:"ProductName"`
	Quantity      int64  `json:"Quantity"`
	ProductPrice  int64  `json:"ProductPrice"`
	ProductAmount int64  `json:"ProductAmount"`
}

type InvoiceRequest struct {
	OrderId string `form:"OrderId" json:"OrderId"`
}

type Awarded struct {
	TCompanyBan    string //賣方-營業人總機構統一編號
	InvoiceYm      string //發票所屬年月
	InvoiceAxle    string //發票號碼-字軌
	InvoiceNumber  string //發票號碼-號碼
	SellerName     string //賣方-營業人名稱
	SellerBan      string //賣方-營業人統一編號
	Address        string //賣方-營業人地址
	InvoiceDate    string //發票日期
	InvoiceTime    string //發票時間
	TotalAmount    string //總計
	CarrierType    string //載具類別號碼
	CarrierName    string //載具類別名稱
	CarrierIdClear string //載具顯碼id
	CarrierIdHide  string //載具隱碼id
	RandomNumber   string //四位隨機碼
	PrizeType      string //中獎獎別
	PrizeAmt       string //中獎獎金
	BuyerBan       string //買受人-扣繳單位統一編號
	DepositMK      string //整合服務平台已匯款註記
	DataType       string //資料類別
	Code           string //例外代碼
	Print          string //列印格式
	Identifier     string //唯一識別碼
}
