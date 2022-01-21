package entity

import "time"

type InvoiceData struct {
	InvoiceId     int64     `xorm:"pk int(11) unique autoincr comment('序號')"`
	Year          string    `xorm:"varchar(4) notnull comment('發票年份')"`
	Month         string    `xorm:"varchar(4) notnull comment('發票月份')"`
	InvoiceTrack  string    `xorm:"varchar(2) notnull comment('發票字軌')"`
	InvoiceNumber string    `xorm:"varchar(8) notnull comment('發票號碼')"`
	InvoiceType   string    `xorm:"varchar(2) notnull comment('發票類別')"`
	OrderId       string    `xorm:"varchar(50) notnull comment('訂單編號')"`
	BuyerId       string    `xorm:"varchar(50) notnull index(buyer_status) comment('會員編號')"`
	Detail        string    `xorm:"varchar(500) notnull comment('開立內容')"`
	Sales         int64     `xorm:"int(10) comment('銷售額')"`
	Tax           int64     `xorm:"int(10) comment('稅額')"`
	Amount        int64     `xorm:"int(10) notnull comment('開立金額')"`
	CheckNumber   string    `xorm:"varchar(10) notnull comment('發票檢查碼')"`
	RandomNumber  string    `xorm:"varchar(4) notnull comment('隨機碼')"`
	PrintMark     string    `xorm:"varchar(50) default 'N' comment('已列印')"`
	Buyer         string    `xorm:"varchar(50) notnull comment('買受人')"`
	Identifier    string    `xorm:"varchar(50) notnull comment('買受人統編')"`
	DonateMark    int64     `xorm:"int(1) default '0' comment('捐贈註記')"`
	DonateBan     string    `xorm:"varchar(50) comment('捐贈對象')"`
	CarrierType   string    `xorm:"varchar(10) comment('發票載具')"`
	Carrier       string    `xorm:"varchar(20) comment('載具')"`
	InvoiceStatus string    `xorm:"varchar(20) notnull index(buyer_status) comment('發票狀態')"`
	VoidReason    string    `xorm:"varchar(100) comment('註銷原因')"`
	CancelReason  string    `xorm:"varchar(100) comment('作廢原因')"`
	AwardModel    string    `xorm:"varchar(1) comment('載具類型')"`
	AwardTime     time.Time `xorm:"datetime comment('中獎時間')"`
	VoidTime      time.Time `xorm:"datetime comment('註銷時間')"`
	CancelTime    time.Time `xorm:"datetime comment('作廢時間')"`
	CreateTime    time.Time `xorm:"datetime notnull comment('開立時間')"`
}

type InvoiceResp struct {
	Invoice InvoiceData `xorm:"extends"`
	Order   OrderData   `xorm:"extends"`
}

//空白字軌
type InvoiceAssignNoData struct {
	AssignId       int64     `xorm:"pk int(11) unique autoincr comment('序號')"`
	InvoiceBan     string    `xorm:"varchar(10) notnull comment('公司統一編號')"`
	InvoiceType    string    `xorm:"varchar(2) notnull comment('發票類別')"`
	MonthYear      string    `xorm:"varchar(5) notnull comment('發票期別')"`
	InvoiceTrack   string    `xorm:"varchar(2) notnull comment('發票字軌')"`
	InvoiceBeginNo string    `xorm:"varchar(8) notnull comment('發票起號')"`
	InvoiceEndNo   string    `xorm:"varchar(8) notnull comment('發票迄號')"`
	InvoiceNowNo   int64     `xorm:"int(8) notnull comment('發票目前號碼')"`
	InvoiceStatus  string    `xorm:"varchar(20) notnull comment('發票狀態')"`
	InvoiceBooklet int64     `xorm:"int(2) notnull comment('本組數')"`
	CreateTime     time.Time `xorm:"datetime notnull comment('建立時間')"`
	UpdateTime     time.Time `xorm:"datetime notnull comment('更新時間')"`
}

type AllowanceData struct {
	AllowanceId     string    `xorm:"pk varchar(16) unique comment('折讓證明單號碼')"` //不可重複
	AllowanceDate   string    `xorm:"varchar(8) notnull comment('折讓證明單日期')"`
	AllowanceType   string    `xorm:"varchar(1) notnull comment('折讓種類')"`
	Identifier      string    `xorm:"varchar(50) notnull comment('買受人統編')"`
	Buyer           string    `xorm:"varchar(50) notnull comment('買受人')"`
	Details         string    `xorm:"text notnull comment('商品項目')"`
	TaxAmount       int64     `xorm:"int(10) notnull comment('營業稅額合計')"`
	TotalAmount     int64     `xorm:"int(10) notnull comment('金額合計')"`
	AllowanceStatus string    `xorm:"varchar(20) notnull comment('折讓單狀態')"`
	CancelReason    string    `xorm:"varchar(100) comment('作廢原因')"`
	CancelTime      time.Time `xorm:"datetime comment('作廢時間')"`
	CreateTime      time.Time `xorm:"datetime notnull comment('建立時間')"`
}

type CancelAllowanceData struct {
	CancelAllowanceId int64     `xorm:"pk int(16) unique comment('作廢折讓證明單號碼')"` //不可重複
	AllowanceDate     string    `xorm:"varchar(8) notnull comment('折讓證明單日期')"`
	BuyerId           string    `xorm:"varchar(50) notnull comment('買方統一編號')"`
	CancelDate        time.Time `xorm:"varchar(8) notnull comment('折讓證明單作廢日期時間')"`
	CancelReason      string    `xorm:"varchar(20) comment('折讓證明單作廢原因')"`
	Remark            string    `xorm:"text comment('備註')"`
	CreateTime        time.Time `xorm:"datetime notnull comment('建立時間')"`
}
