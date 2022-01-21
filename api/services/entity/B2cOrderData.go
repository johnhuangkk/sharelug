package entity

import "time"

type B2cOrderData struct {
	OrderId       string    `xorm:"pk varchar(50) unique comment('訂單編號')"`
	UserId        string    `xorm:"varchar(50) notnull comment('使用者ID')"`
	StoreId       string    `xorm:"varchar(50) notnull comment('收銀機ID')"`
	ProductId     string    `xorm:"varchar(30) notnull comment('商品ID')"`
	ProductName   string    `xorm:"varchar(30) notnull comment('商品名稱')"`
	ProductDetail string    `xorm:"varchar(300) comment('商品內容')"`
	OrderDetail   string    `xorm:"text comment('訂單內容')"`
	BillingTime   string    `xorm:"varchar(300) comment('計費時間')"`
	UpgradeOrder  string    `xorm:"text comment('合併帳單')"`
	UpgradeLevel  int64     `xorm:"int(10) notnull comment('升級等級')"`
	Amount        int64     `xorm:"int(10) notnull comment('金額')"`
	Payment       string    `xorm:"varchar(20) comment('付款方式')"`
	OrderStatus   string    `xorm:"varchar(20) notnull comment('訂單狀態')"`
	OrderSys      int64     `xorm:"tinyint(1) default 0 comment('是否為系統訂單')"`
	InvoiceType   string    `xorm:"varchar(10) notnull default 'PERSONAL' comment('發票設定')"`
	CompanyBan    string    `xorm:"varchar(10) comment('統一編號')"`
	CompanyName   string    `xorm:"varchar(50) comment('公司名稱')"`
	DonateBan     string    `xorm:"varchar(10) comment('捐贈碼')"`
	CarrierType   string    `xorm:"varchar(10) notnull default 'MEMBER' comment('發票載具')"`
	CarrierId     string    `xorm:"varchar(64) comment('載具ID')"`
	AskInvoice    bool      `xorm:"tinyint(1) default 0 comment('是否需開立發票')"`
	InvoiceStatus string 	`xorm:"varchar(20) default 'INIT' comment('是否開立')"`
	Expiration    time.Time `xorm:"datetime comment('到期時間')"`
	CreateTime    time.Time `xorm:"datetime notnull comment('建立時間')"`
	UpdateTime    time.Time `xorm:"datetime notnull comment('更新時間')"`
}

type B2cBillingData struct {
	BillingId     string    `xorm:"pk varchar(50) unique comment('帳單編號')"`
	UserId        string    `xorm:"varchar(50) notnull comment('使用者ID')"`
	StoreId       string    `xorm:"varchar(50) notnull comment('收銀機ID')"`
	BillName      string    `xorm:"varchar(30) notnull comment('帳單名稱')"`
	ProductId     string    `xorm:"varchar(30) notnull comment('商品ID')"`
	ProductName   string    `xorm:"varchar(30) notnull comment('商品名稱')"`
	ProductDesc   string    `xorm:"varchar(50) notnull comment('商品說明')"`
	BillingTime   string    `xorm:"varchar(300) comment('計費時間')"`
	BillingLevel  int64     `xorm:"int(10) notnull comment('升級等級')"`
	Amount        int64     `xorm:"int(10) notnull comment('金額')"`
	OrderId       string    `xorm:"varchar(50) comment('訂單編號')"`
	BillingStatus string    `xorm:"varchar(20) notnull comment('計費狀態')"`
	ServiceType   string    `xorm:"varchar(20) notnull comment('狀態')"`
	Expiration    time.Time `xorm:"datetime comment('到期時間')"`
	CreateTime    time.Time `xorm:"datetime notnull comment('建立時間')"`
	UpdateTime    time.Time `xorm:"datetime notnull comment('更新時間')"`
}

type B2cOrder struct {
	Detail []B2cOrderDetail
}

//商品名稱、單價、類別
type B2cOrderDetail struct {
	ProductId     string `json:"ProductId"`
	ProductName   string `json:"ProductName"`
	ProductDetail string `json:"ProductDetail"`
	ProductAmount int64  `json:"ProductAmount"`
	ProductType   string `json:"ProductType"`
	BillingTime   string `json:"BillingTime"`
}

type B2cOrderVo struct {
	ProductId    string
	ProductName  string
	UserId       string
	StoreId      string
	OrderDetail  string
	BillingTime  string
	UpgradeLevel int64
	Amount       int64
	Expire       time.Time
	InvoiceType  string
	CompanyBan   string
	CompanyName  string
	DonateBan    string
	CarrierType  string
	CarrierId    string
}