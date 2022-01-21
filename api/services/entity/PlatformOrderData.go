package entity

import "time"

type PlatformOrderData struct {
	OrderId    string    `xorm:"varchar(50) unique comment('訂單編號')"`
	UserId     string    `xorm:"varchar(50) notnull comment('購買ID')"`
	Program    string    `xorm:"varchar(50) notnull comment('購買方案')"`
	Amount     int64     `xorm:"decimal(10,2) notnull comment('金額')"`
	CreateTime time.Time `xorm:"datetime notnull comment('建立時間')"`
}
