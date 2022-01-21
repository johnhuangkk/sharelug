package entity

import "time"

type TwIdVerifyData struct {
	Id               int       `xorm:"pk unique int autoincr comment('ID')"`
	UserId           string    `xorm:"varchar(50) notnull comment('會員ID')"`
	IdentityName     string    `xorm:"varchar(50) notnull comment('中文姓名')"`
	IdentityId       string    `xorm:"varchar(10) notnull comment('身份字號')"`
	IdentityType     string    `xorm:"varchar(2) notnull comment('領補換類別')"`
	IdentityDate     string    `xorm:"varchar(10) notnull comment('發證日期')"`
	IdentityCounties string    `xorm:"varchar(10) notnull comment('發證地點')"`
	HttpResponseCode string    `xorm:"varchar(50) comment('回應碼')"`
	HttpResponseMsg  string    `xorm:"varchar(50) comment('回應訊息')"`
	Response         string    `xorm:"text comment('驗證回傳內容')"`
	ResponseMsg      string    `xorm:"varchar(100) comment('驗證回應訊息')"`
	ResponseCode     string    `xorm:"varchar(50) comment('驗證回應碼')"`
	CreateTime       time.Time `xorm:"datetime notnull comment('建立時間')"`
	UpdateTime       time.Time `xorm:"datetime notnull comment('建立時間')"`
}
