package entity

import "time"

type SmsLogData struct {
	Id         int64     `xorm:"pk int(10) unique autoincr comment('序號')"`
	MerchantId string    `xorm:"varchar(10) notnull comment('簡訊廠商')"`
	Phone      string    `xorm:"varchar(20) notnull comment('手機號碼')"`
	Content	   string    `xorm:"text notnull comment('簡訊內容')"`
	ResultCode string    `xorm:"varchar(20) comment('回應代碼')"`
	ResultMsg  string    `xorm:"text comment('回應內容')"`
	CreateTime time.Time `xorm:"datetime notnull comment('建立時間')"`
}
