package entity

import (
	"api/services/Enum"
	"api/services/VO/OrderVo"
	"api/services/VO/Response"
	"api/services/util/images"
	"api/services/util/qrcode"
	"api/services/util/tools"
	"time"
)

type OrderData struct {
	OrderId          string    `xorm:"pk unique varchar(50) notnull comment('訂單編號')"`
	BuyerId          string    `xorm:"varchar(50) notnull comment('買家ID')"`
	SellerId         string    `xorm:"varchar(50) notnull comment('賣家ID')"`
	StoreId          string    `xorm:"varchar(50) notnull comment('店家ID')"`
	SubTotal         float64   `xorm:"decimal(10,2) notnull default 0.00 comment('小計金額')"`
	ShipFee          float64   `xorm:"decimal(10,2) notnull default 0.00 comment('運費')"`
	TotalAmount      float64   `xorm:"decimal(10,2) notnull default 0.00 comment('總金額')"`
	OrderStatus      string    `xorm:"varchar(10) notnull default 'INIT' comment('訂單狀態')"`
	RefundStatus     string    `xorm:"varchar(10) notnull default 'INIT' comment('退貨狀態')"`
	PayWay           string    `xorm:"varchar(10) notnull comment('付款方式')"`
	PayWayTime       time.Time `xorm:"datetime comment('付款時間')"`
	BuyerName        string    `xorm:"varchar(50) notnull comment('購買人名稱')"`
	BuyerPhone       string    `xorm:"varchar(20) notnull comment('購買人電話')"`
	ReceiverName     string    `xorm:"varchar(50) notnull comment('收件人姓名')"`
	ReceiverAddress  string    `xorm:"varchar(100) notnull comment('收件人地址')"`
	ReceiverPhone    string    `xorm:"varchar(20) notnull comment('收件人電話')"`
	ShipStatus       string    `xorm:"varchar(10) notnull comment('出貨狀態')"`
	ShipType         string    `xorm:"varchar(30) notnull comment('運送方式')"`
	ShipText         string    `xorm:"varchar(50) comment('貨運單位名稱')"`
	ShipTime         time.Time `xorm:"datetime comment('出貨時間＆交易完成時間')"`
	ShipNumber       string    `xorm:"varchar(50) comment('出貨單號')"`
	ShipExpire       time.Time `xorm:"datetime comment('出貨到期日')"`
	PlatformShipFee  float64   `xorm:"decimal(10,2) notnull default 0.00 comment('平台運費')"`
	PlatformTransFee float64   `xorm:"decimal(10,2) notnull default 0.00 comment('交易手續費')"`
	PlatformInfoFee  float64   `xorm:"decimal(10,2) notnull default 0.00 comment('資料手續費')"`
	PlatformPayFee   float64   `xorm:"decimal(10,2) notnull default 0.00 comment('金流手續費')"`
	InvoiceNumber    string    `xorm:"varchar(12) comment('發票號碼')"`
	InvoiceYearMonth string    `xorm:"varchar(12) comment('發票期別')"`
	SellerUnread     int       `xorm:"tinyint(1) default 0 comment('賣家是否讀取')"`
	BuyerUnread      int       `xorm:"tinyint(1) default 0 comment('買家是否讀取')"`
	CaptureStatus    string    `xorm:"varchar(10) notnull comment('請款狀態')"`
	CaptureApply     string    `xorm:"varchar(10) notnull default 'INIT' comment('請款申請狀態')"`
	CaptureTime      time.Time `xorm:"datetime comment('可請款時間')"`
	CaptureAudit     int       `xorm:"tinyint(1) default 1 comment('審核撥款')"`
	CaptureAmount    float64   `xorm:"decimal(10,2) notnull default 0.00 comment('請款金額')"`
	AskInvoice       bool      `xorm:"tinyint(1) default 0 comment('是否需開立發票')"`
	InvoiceStatus    string    `xorm:"varchar(20) default 'INIT' comment('是否開立')"`
	OrderMemo        string    `xorm:"text comment('備註')"`
	BuyerNotes       string    `xorm:"text comment('買家備註')"`
	CreateTime       time.Time `xorm:"datetime notnull comment('建立時間')"`
	UpdateTime       time.Time `xorm:"datetime notnull comment('更新時間')"`
	ArrivedTime      time.Time `xorm:"datetime comment('貨物到店時間')"`
	CsvCheck         int64     `xorm:"tinyint(1) default 1 comment('超取付入帳')"`
	FormUrl          string    `xorm:"varchar(250) comment('表單網址')"`
	FreeShipKey      string    `xorm:"varchar(20) default 'NONE' comment('免運費方式')"`
	FreeShip         int64     `xorm:"int(10) default 0 comment('免運費')"`
	BeforeTotal      float64   `xorm:"decimal(10,2) notnull default 0.00 comment('折扣前金額')"`
	Coupon           float64   `xorm:"decimal(10,2) notnull default 0.00 comment('優惠金額')"`
	CouponNumber     string    `xorm:"varchar(10) comment('優惠號碼')"`
}

func (o *OrderData) GetOrder() OrderVo.CalculatePlatFormFeeVo {
	var resp OrderVo.CalculatePlatFormFeeVo
	resp.PayWayType = o.PayWay
	resp.ShipType = o.ShipType
	resp.Amount = o.TotalAmount
	return resp
}

func (o *OrderData) GetOrderResponse(store StoreData) Response.OrderResponse {
	var resp Response.OrderResponse
	resp.OrderId = o.OrderId
	resp.Buyer.BuyerName = o.BuyerName
	resp.Buyer.BuyerPhone = o.BuyerPhone
	resp.Buyer.BuyerUid = o.BuyerId
	resp.BuyerMasker.BuyerName = tools.MaskerName(o.BuyerName)
	resp.BuyerMasker.BuyerPhone = tools.MaskerPhone(o.BuyerPhone)
	resp.BuyerMasker.BuyerUid = tools.MaskerPhone(o.BuyerId)
	resp.StoreName = store.StoreName
	resp.StoreId = store.StoreId
	resp.OrderTime = o.CreateTime.Format("2006/01/02 15:04")

	resp.OrderStatusType = o.OrderStatus
	resp.ShipStatusType = o.ShipStatus
	resp.ShipStatusText = Enum.OrderShipStatus[o.ShipStatus]
	resp.ShipCompany = o.ShipText

	resp.CaptureStatusType = o.CaptureStatus
	resp.CaptureStatusText = Enum.OrderCaptureStatus[o.CaptureStatus]
	resp.CaptureApplyType = o.CaptureApply
	resp.CaptureApplyText = Enum.OrderCaptureStatus[o.CaptureApply]
	resp.ShipExpire = ""
	if !o.ShipExpire.IsZero() {
		resp.ShipExpire = o.ShipExpire.Format("2006/01/02 15:04:05")
	}
	if !o.ShipTime.IsZero() {
		resp.ShipTime = o.ShipTime.Format("2006/01/02 15:04:05")
	}
	resp.PayWayTime = ""
	if !o.PayWayTime.IsZero() {
		resp.PayWayTime = o.PayWayTime.Format("2006/01/02 15:04:05")
	}
	resp.CaptureTime = ""
	if !o.CaptureTime.IsZero() {
		resp.CaptureTime = o.CaptureTime.Format("2006/01/02 15:04")
	}
	resp.InvoiceNumber = o.InvoiceNumber
	resp.TotalShipFee = int64(o.ShipFee)
	resp.Coupon = int64(o.Coupon)
	resp.CouponNumber = o.CouponNumber
	resp.TotalAmount = int64(o.TotalAmount)
	resp.SubTotal = int64(o.SubTotal)
	resp.ShipType = o.ShipType
	resp.ShipNumber = o.ShipNumber
	resp.PlatformShipFee = int64(o.PlatformShipFee)
	resp.PlatformTransFee = int64(o.PlatformTransFee)
	resp.PlatformInfoFee = int64(o.PlatformInfoFee)
	resp.PlatformPayFee = int64(o.PlatformPayFee)
	resp.OrderMemo = o.OrderMemo
	resp.BuyerNotes = o.BuyerNotes
	track := o.OrderId[0:2]
	resp.PayType = Enum.PayTypeDefault
	if track == "BM" {
		resp.PayType = Enum.PayTypeMarket
	}
	return resp
}

type PayParams struct {
	BuyerPhone     string `form:"BuyerPhone" validate:"required"`
	BuyerName      string `form:"BuyerName"`
	ReceiverName   string `form:"ReceiverName"`
	ReceiverPhone  string `form:"ReceiverPhone"`
	Address        string `form:"Address" validate:"required"`
	PayWay         string `form:"PayWay"  validate:"required"`
	CardId         string `form:"CreditId" json:"CreditId"`
	CardNumber     string `form:"CreditNumber" json:"CreditNumber"`
	CardExpiration string `form:"CreditExpiration" json:"CreditExpiration"`
	CardSecurity   string `form:"CreditSecurity" json:"CreditSecurity"`
}

type PayResponse struct {
	OrderId string `json:"OrderId"`
	RtnURL  string `json:"RtnURL"`
}

type OrderResponse struct {
	OrderId    string      `json:"OrderId"`
	StoreName  string      `json:"StoreName"`
	OrderTime  string      `json:"OrderTime"`
	BuyerName  string      `json:"BuyerName"`
	BuyerPhone string      `json:"BuyerPhone"`
	Payment    interface{} `json:"Payment"`
	Shipping   interface{} `json:"Shipping"`
}

type Shipping struct {
	OtherShipping
	ReceiverName    string `json:"ReceiverName"`
	ReceiverPhone   string `json:"ReceiverPhone"`
	ReceiverAddress string `json:"ReceiverAddress"`
}

type OtherShipping struct {
	Type string `json:"Type"`
	Text string `json:"Text"`
}

type BalancePayWay struct {
	Type        string  `json:"Type"`
	Text        string  `json:"Text"`
	Balance     float64 `json:"Balance"`
	OrderAmount float64 `json:"OrderAmount"`
}

type TransferPayWay struct {
	Type           string `json:"Type"`
	Text           string `json:"Text"`
	BankName       string `json:"BankName"`
	BankAccount    string `json:"BankAccount"`
	BankExpireDate string `json:"BankExpireDate"`
}

type OtherPayWay struct {
	Type string `json:"Type"`
	Text string `json:"Text"`
}

type OrderDetailData struct {
	Order  OrderData   `xorm:"extends"`
	Detail OrderDetail `xorm:"extends"`
}

type OrderResp struct {
	Order  OrderData       `xorm:"extends"`
	Refund OrderRefundData `xorm:"extends"`
}

type ErpSearchOrders struct {
	Order  OrderData  `xorm:"extends"`
	Store  StoreData  `xorm:"extends"`
	Member MemberData `xorm:"extends"`
}

type ErpSearchBuyerOrders struct {
	Order  OrderData  `xorm:"extends"`
	Member MemberData `xorm:"extends"`
}

type ErpOrder struct {
	Order   OrderData   `xorm:"extends"`
	Member  MemberData  `xorm:"extends"`
	Product ProductData `xorm:"extends"`
	Uid     string
}

func (o *BillOrderData) GetSearchBill(userData MemberData) Response.SearchProductResponse {
	var data Response.SearchProductResponse
	data.CreateDate = o.CreateTime.Format("2006/01/02 15:04")
	data.UserAccount = tools.MaskerPhoneLater(userData.Mphone)
	data.StoreName = userData.Username
	data.ProductId = o.BillId
	data.ProductName = o.ProductName
	data.Amount = o.ProductPrice
	data.Quantity = o.ProductQty
	data.ShipFee = o.ShipFee
	data.Total = int64(o.TotalAmount)
	data.PayWayMode = Enum.PayWay[o.PayWayType]
	data.ShipMode = Enum.Shipping[o.ShipType]
	data.ProductSpec = o.ProductSpec
	data.ProductLink = qrcode.GetTinyUrl(o.TinyUrl)
	data.BuyerName = o.BuyerName
	data.ReceiverName = o.ReceiverName
	data.Url = o.ProductLink
	return data
}

func (o *BillOrderData) GetBillOrder() OrderVo.CalculatePlatFormFeeVo {
	var resp OrderVo.CalculatePlatFormFeeVo
	resp.PayWayType = o.PayWayType
	resp.ShipType = o.ShipType
	resp.Amount = o.TotalAmount
	return resp
}

func (o *BillOrderData) GetBillToResponse() Response.BillResponse {
	var resp Response.BillResponse
	resp.BillId = o.BillId
	if len(o.ProductImage) != 0 {
		resp.ProductImage = images.GetImageUrl(o.ProductImage)
	} else {
		resp.ProductImage = ""
	}
	resp.ProductName = o.ProductName
	resp.ProductLink = o.ProductLink
	resp.ProductSpec = o.ProductSpec
	resp.ProductPrice = o.ProductPrice
	resp.ProductQty = o.ProductQty
	resp.BuyerName = o.BuyerName
	resp.ReceiverName = o.ReceiverName
	resp.ShipType = o.ShipType
	resp.ShipFee = o.ShipFee
	resp.PayWayType = o.PayWayType
	resp.TotalAmount = int64(o.TotalAmount)
	resp.PlatformInfoFee = int64(o.PlatformInfoFee)
	resp.PlatformPayFee = int64(o.PlatformPayFee)
	resp.PlatformTransFee = int64(o.PlatformTransFee)
	resp.PlatformShipFee = int64(o.PlatformShipFee)
	resp.CaptureAmount = int64(o.CaptureAmount)
	resp.TinyUrl = qrcode.GetTinyUrl(o.TinyUrl)
	resp.Expire = o.BillExpire.Format("2006/01/02 15:04:05")
	resp.IsExtension = o.IsExtension
	resp.BillStatus = o.BillStatus
	resp.PayWayStatus = o.PayWayStatus
	return resp
}

func (o *ErpSearchOrders) GetSearchOrders() Response.OrdersResponse {
	var data Response.OrdersResponse
	data.OrderDate = o.Order.CreateTime.Format("2006/01/02 15:04")
	data.Seller = o.Member.Mphone
	data.SellerId = o.Member.TerminalId
	data.StoreName = o.Store.StoreName
	data.OrderId = o.Order.OrderId
	data.OrderStatus = o.Order.OrderStatus
	data.OrderAmount = int64(o.Order.TotalAmount)
	data.PaymentType = o.Order.PayWay
	data.PaymentTypeText = Enum.PayWay[o.Order.PayWay]
	data.PaymentTime = o.Order.PayWayTime.Format("2006/01/02 15:04")
	data.Coupon = o.Order.CouponNumber
	data.CouponAmount = int64(o.Order.Coupon)
	data.ShipType = o.Order.ShipType
	data.ShipTypeText = Enum.Shipping[o.Order.ShipType]
	data.CaptureStatus = o.Order.CaptureStatus
	data.CaptureStatusText = Enum.OrderCaptureStatus[o.Order.CaptureStatus]
	data.ProductAmount = int64(o.Order.SubTotal)
	data.ShipFee = int64(o.Order.ShipFee)
	data.PlatformFee = int64(o.Order.PlatformTransFee + o.Order.PlatformShipFee + o.Order.PlatformInfoFee + o.Order.PlatformPayFee)
	return data
}

func (o *ErpSearchOrders) GetSearchOrder() Response.SearchOrderDetailResponse {
	var data Response.SearchOrderDetailResponse
	data.SellerName = o.Member.Username
	data.Seller = tools.MaskerPhone(o.Member.Mphone)
	data.StoreName = o.Store.StoreName
	data.OrderDate = o.Order.CreateTime.Format("2006/01/02 15:04")
	data.OrderId = o.Order.OrderId
	data.OrderStatus = o.Order.OrderStatus
	data.OrderAmount = int64(o.Order.TotalAmount)
	data.ProductAmount = int64(o.Order.SubTotal)
	data.ShipFee = int64(o.Order.ShipFee)
	data.PlatformFee = int64(o.Order.PlatformTransFee + o.Order.PlatformShipFee + o.Order.PlatformInfoFee + o.Order.PlatformPayFee)
	data.PlatformPayFee = int64(o.Order.PlatformPayFee)
	data.PlatformShipFee = int64(o.Order.PlatformShipFee)
	data.PlatformInfoFee = int64(o.Order.PlatformInfoFee)
	data.Coupon = o.Order.CouponNumber
	data.CouponAmount = int64(o.Order.Coupon)
	data.BuyerNotes = o.Order.BuyerNotes
	data.SellerNotes = o.Order.OrderMemo
	return data
}

type BatchShipExcelImport struct {
	BatchId       string    `xorm:"pk varchar(20) unique comment('序號')"`
	StoreId       string    `xorm:"varchar(50) notnull comment('賣場ID')"`
	BeforeContent string    `xorm:"LONGTEXT comment('處理前的內容')"`
	AfterContent  string    `xorm:"LONGTEXT comment('處理後的內容')"`
	ProcessStatus string    `xorm:"varchar(10) notnull comment('處理狀態')"`
	CreateTime    time.Time `xorm:"datetime notnull comment('建立時間')"`
	UpdateTime    time.Time `xorm:"datetime notnull comment('更新時間')"`
}

func (o OrderData) GeneratorCouponUsedRecordData(promotionId int64) CouponUsedRecord {
	var data CouponUsedRecord
	data.StoreId = o.StoreId
	data.OrderId = o.OrderId
	data.PromotionId = promotionId
	data.BuyerId = o.BuyerId
	data.BuyerPhone = o.BuyerPhone
	data.Code = o.CouponNumber
	data.Amount = o.BeforeTotal
	data.DiscountAmount = o.TotalAmount
	data.TransTime = o.CreateTime
	data.RecordStatus = Enum.RecordStatusSuccess
	data.Created = time.Now()
	return data
}
