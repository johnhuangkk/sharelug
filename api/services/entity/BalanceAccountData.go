package entity

import "time"

type BalanceAccountData struct {
	Id         int64     `xorm:"pk int(10) unique autoincr comment('序號')"`
	UserId     string    `xorm:"varchar(50) notnull comment('帳戶使用者')"`
	DataId     string    `xorm:"varchar(50) notnull comment('交易單號')"`
	TransType  string    `xorm:"varchar(10) notnull comment('交易類別')"`
	In         float64   `xorm:"decimal(10,2) notnull comment('進項金額')"`
	Out        float64   `xorm:"decimal(10,2) notnull comment('出項金額')"`
	Balance    float64   `xorm:"decimal(10,2) notnull comment('餘額金額')"`
	Comment    string    `xorm:"varchar(100) comment('備註')"`
	CreateTime time.Time `xorm:"datetime notnull comment('建立時間')"`
}

type BalanceRetainAccountData struct {
	Id         int64     `xorm:"pk int(10) unique autoincr comment('序號')"`
	UserId     string    `xorm:"varchar(50) notnull comment('帳戶使用者')"`
	DataId     string    `xorm:"varchar(50) notnull comment('交易單號')"`
	TransType  string    `xorm:"varchar(10) notnull comment('交易類別')"`
	In         float64   `xorm:"decimal(10,2) notnull comment('進項金額')"`
	Out        float64   `xorm:"decimal(10,2) notnull comment('出項金額')"`
	Balance    float64   `xorm:"decimal(10,2) notnull comment('餘額金額')"`
	Comment    string    `xorm:"varchar(100) comment('備註')"`
	CreateTime time.Time `xorm:"datetime notnull comment('建立時間')"`
}

type BalanceRetainByOrderData struct {
	Retain BalanceRetainAccountData `xorm:"extends"`
	Order  OrderData                `xorm:"extends"`
}
