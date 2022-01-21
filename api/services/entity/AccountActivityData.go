package entity

import "time"

type AccountActivityData struct {
	Id         int       `xorm:"pk int(10) unique autoincr comment('序號')"`
	UserId     string    `xorm:"varchar(50) notnull comment('使用者')"`
	Action     string    `xorm:"varchar(10) notnull comment('動作')"`
	Message    string    `xorm:"varchar(50) notnull comment('內容')"`
	CreateTime time.Time `xorm:"datetime notnull comment('建立時間')"`
}
