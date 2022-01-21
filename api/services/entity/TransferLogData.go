package entity

import "time"

type TransferLogData struct {
	Id         int       `xorm:"pk int(10) unique autoincr"`
	Response   string    `xorm:"text notnull"`
	CreateTime time.Time `xorm:"datetime notnull"`
}
