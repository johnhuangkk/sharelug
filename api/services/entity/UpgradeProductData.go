package entity

import "time"

type UpgradeProductData struct {
	ProductId    string    `xorm:"pk varchar(50) unique comment('商品編號')"`
	ProductName  string    `xorm:"varchar(50) notnull comment('商品名稱')"`
	Description  string    `xorm:"text notnull comment('商品說明')"`
	Note         string    `xorm:"varchar(50) comment('商品備註')"`
	Amount       int64     `xorm:"int(10) notnull comment('商品價格')"`
	UpgradeLevel int64     `xorm:"int(10) notnull comment('等級')"`
	Store        int64     `xorm:"int(10) notnull comment('店數')"`
	Manager      int64     `xorm:"int(10) notnull comment('管理者數')"`
	CreateTime   time.Time `xorm:"datetime notnull comment('建立時間')"`
	UpdateTime   time.Time `xorm:"datetime notnull comment('建立時間')"`
}
