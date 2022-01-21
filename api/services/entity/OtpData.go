package entity

import (
	"time"
)

type OtpData struct {
	Id         int       `xorm:"pk int(10) unique autoincr comment('序號')"`
	Uid        string    `xorm:"varchar(50) notnull comment('會員UID')"`
	Phone      string    `xorm:"varchar(10) notnull comment('手機號碼')"`
	Email      string    `xorm:"varchar(50) comment('電子信箱')"`
	OtpNumber  string    `xorm:"varchar(6) notnull comment('OTP號碼')"`
	OtpUse     int       `xorm:"tinyint(1) default 0 comment('是否使用過')"`
	SendFreq   int       `xorm:"int(2) default 1 comment('發送次數')"`
	ExpireTime time.Time `xorm:"datetime notnull comment('到期時間')"`
	CreateTime time.Time `xorm:"datetime notnull comment('建立時間')"`
}
