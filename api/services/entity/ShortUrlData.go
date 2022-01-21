package entity

import "time"

type ShortUrlData struct {
	Short      string    `xorm:"pk varchar(20) unique comment('短網址')"`
	Url        string    `xorm:"varchar(50) notnull comment('網址')"`
	CreateTime time.Time `xorm:"datetime notnull comment('建立時間')"`
}
