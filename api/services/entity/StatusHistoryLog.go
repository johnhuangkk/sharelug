package entity

import "time"

type StatusHistoryLog struct {
	Id            int       `xorm:"pk int(10) unique autoincr comment('ID')"`
	Table         string    `xorm:"varchar(20) notnull comment('表格名稱')"`
	Field         string    `xorm:"varchar(20) notnull comment('欄位名稱')"`
	DataId        string    `xorm:"varchar(50) notnull comment('資料ID')"`
	OldValue      string    `xorm:"varchar(10) notnull comment('原始狀態')"`
	NewValue      string    `xorm:"varchar(10) notnull comment('新的狀態')"`
	OperateUserId string    `xorm:"varchar(50) notnull comment('變更者')"`
	CreateTime    time.Time `xorm:"datetime notnull comment('建立時間')"`
}

type ProductHistoryLog struct {
	Id            int         `xorm:"pk int(10) unique autoincr comment('ID')"`
	ProductId     string      `xorm:"varchar(50) notnull comment('資料ID')"`
	Action        string      `xorm:"varchar(20) notnull comment('動作')"`
	OldValue      ProductData `xorm:"json notnull comment('原始狀態')"`
	NewValue      ProductData `xorm:"json notnull comment('新的狀態')"`
	OperateUserId string      `xorm:"varchar(50) notnull comment('變更者')"`
	CreateTime    time.Time   `xorm:"timestamp created notnull comment('建立時間')"`
}
