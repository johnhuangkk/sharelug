package entity

import (
	"api/services/Enum"
	"api/services/VO/Request"
	"api/services/util/tools"
	"time"
)

type BillOrderData struct {
	BillId           string    `xorm:"pk unique varchar(50) comment('訂單編號')"`
	BuyerId          string    `xorm:"varchar(50) notnull comment('買家ID')"`
	ProductImage     string    `xorm:"varchar(100) comment('商品圖')"`
	ProductName      string    `xorm:"varchar(100) notnull comment('商品名稱')"`
	ProductLink      string    `xorm:"varchar(200) notnull comment('商品連結')"`
	ProductSpec      string    `xorm:"varchar(100) comment('商品規格')"`
	ProductPrice     int64     `xorm:"int(10) notnull comment('商品單價')"`
	ProductQty       int64     `xorm:"int(10) notnull comment('商品數量')"`
	SubTotal         float64   `xorm:"decimal(10,2) notnull default 0.00 comment('小計金額')"`
	TotalAmount      float64   `xorm:"decimal(10,2) notnull default 0.00 comment('總金額')"`
	ShipType         string    `xorm:"varchar(20) null comment('運送方式')"`
	ShipFee          int64     `xorm:"varchar(20) null comment('運費')"`
	BuyerName        string    `xorm:"varchar(50) notnull comment('購買人名稱')"`
	BuyerPhone       string    `xorm:"varchar(20) notnull comment('購買人電話')"`
	ReceiverName     string    `xorm:"varchar(50) notnull comment('收件人姓名')"`
	ReceiverAddress  string    `xorm:"varchar(100) notnull comment('收件人地址')"`
	ReceiverPhone    string    `xorm:"varchar(20) notnull comment('收件人電話')"`
	PayWayType       string    `xorm:"varchar(10) notnull comment('付款方式')"`
	PayWayStatus     string    `xorm:"varchar(20) notnull comment('付款狀態')"`
	TinyUrl          string    `xorm:"varchar(100) comment('短網址')"`
	PlatformShipFee  float64   `xorm:"decimal(10,2) notnull default 0.00 comment('平台運費')"`
	PlatformTransFee float64   `xorm:"decimal(10,2) notnull default 0.00 comment('交易手續費')"`
	PlatformInfoFee  float64   `xorm:"decimal(10,2) notnull default 0.00 comment('資料手續費')"`
	PlatformPayFee   float64   `xorm:"decimal(10,2) notnull default 0.00 comment('金流手續費')"`
	CaptureAmount    float64   `xorm:"decimal(10,2) notnull default 0.00 comment('請款金額')"`
	BillStatus       string    `xorm:"varchar(20) notnull comment('帳單狀態')"`
	BillExpire       time.Time `xorm:"datetime comment('帳單到期日')"`
	IsExtension      bool      `xorm:"tinyint(1) default 0 comment('是否已延期')"`
	PayWayTime       time.Time `xorm:"datetime comment('付款時間')"`
	CreateTime       time.Time `xorm:"datetime notnull comment('建立時間')"`
	UpdateTime       time.Time `xorm:"datetime notnull comment('更新時間')"`
}

type BillOrderParams struct {
	ProductImage   string `form:"ProductImage" json:"ProductImage"`
	ProductName    string `form:"ProductName" json:"ProductName"`
	ProductLink    string `form:"ProductLink" json:"ProductLink"`
	ProductSpec    string `form:"ProductSpec" json:"ProductSpec"`
	ProductPrice   int64  `form:"ProductPrice" json:"ProductPrice"`
	ProductQty     int64  `form:"ProductQty" json:"ProductQty"`
	ShipType       string `form:"ShipType" json:"ShipType"`
	ShipFee        int64  `form:"ShipFee" json:"ShipFee"`
	BuyerPhone     string `form:"BuyerPhone" json:"BuyerPhone"`
	BuyerName      string `form:"BuyerName" json:"BuyerName"`
	ReceiverId     string `form:"ReceiverId" json:"ReceiverId"`
	ReceiverName   string `form:"ReceiverName" json:"ReceiverName"`
	ReceiverPhone  string `form:"ReceiverPhone" json:"ReceiverPhone"`
	Address        string `form:"Address" json:"Address"`
	PayWayType     string `form:"PayWayType" json:"PayWayType"`
	CardId         string `form:"CreditId" json:"CreditId"`
	CardNumber     string `form:"CreditNumber" json:"CreditNumber"`
	CardExpiration string `form:"CreditExpiration" json:"CreditExpiration"`
	CardSecurity   string `form:"CreditSecurity" json:"CreditSecurity"`
}

func (params *BillOrderParams) GeneratorBillOrderData(userData MemberData) BillOrderData {
	var data BillOrderData
	SubTotal := float64(params.ProductPrice * params.ProductQty)
	data.BillId = tools.GeneratorBillOrderId()
	data.BuyerId = userData.Uid
	data.ProductName = params.ProductName
	data.ProductLink = params.ProductLink
	data.ProductSpec = params.ProductSpec
	data.ProductPrice = params.ProductPrice
	data.ProductQty = params.ProductQty
	data.SubTotal = SubTotal
	data.TotalAmount = SubTotal + float64(params.ShipFee)
	data.ShipType = params.ShipType
	data.ShipFee = params.ShipFee
	data.BuyerName = params.BuyerName
	data.BuyerPhone = params.BuyerPhone
	if len(params.ReceiverName) == 0 {
		data.ReceiverName = params.BuyerName
		data.ReceiverPhone = params.BuyerPhone
	} else {
		data.ReceiverName = params.ReceiverName
		data.ReceiverPhone = params.ReceiverPhone
	}
	data.ReceiverAddress = params.Address
	data.PayWayType = params.PayWayType
	data.PayWayStatus = Enum.OrderWait
	data.PayWayTime = time.Time{}
	data.BillStatus = Enum.BillStatusInit
	data.BillExpire = time.Now().Add(48 * time.Hour)
	return data
}

func (params *BillOrderParams) GetCreditPayment() Request.PayParams {
	var resp Request.PayParams
	resp.BuyerPhone = params.BuyerPhone
	resp.BuyerName = params.BuyerName
	resp.ReceiverName = params.ReceiverName
	resp.ReceiverPhone = params.ReceiverPhone
	resp.Address = params.Address
	resp.CardId = params.CardId
	resp.CardNumber = params.CardNumber
	resp.CardExpiration = params.CardExpiration
	resp.CardSecurity = params.CardSecurity
	resp.PayWay = params.PayWayType
	return resp
}
