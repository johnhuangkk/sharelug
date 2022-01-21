package entity

import "time"

type NotifyMessageData struct {
	Id         int64     `xorm:"pk int(10) unique autoincr comment('序號')"`
	Email      string    `xorm:"varchar(50) comment('Mail')"`
	Phone      string    `xorm:"varchar(20) comment('Phone')"`
	CreateTime time.Time `xorm:"datetime notnull comment('建立時間')"`
}

type ContactData struct {
	Id         int64     `xorm:"pk int(10) unique autoincr comment('序號')"`
	Email      string    `xorm:"varchar(100) notnull comment('信箱')"`
	UserName   string    `xorm:"varchar(50) notnull comment('名稱')"`
	Company    string    `xorm:"varchar(100) notnull comment('公司名稱')"`
	Telephone  string    `xorm:"varchar(20) notnull comment('聯絡電話')"`
	Contents   string    `xorm:"text notnull comment('留言內容')"`
	CreateTime time.Time `xorm:"datetime notnull comment('建立時間')"`
}

type OnlineNotifyData struct {
	Id         int64     `xorm:"pk int(10) unique autoincr comment('序號')"`
	Type       string    `xorm:"varchar(20) notnull comment('類別')"`
	UserId     string    `xorm:"varchar(50) notnull comment('使用者')"`
	Messages   string    `xorm:"text notnull comment('訊息內容')"`
	Unread     int64     `xorm:"tinyint(1) default 0 comment('已讀')"`
	MsgType    string    `xorm:"varchar(20) notnull comment('訊息類別')"`
	OrderId    string	 `xorm:"varchar(50) comment('訂單編號')"`
	CreateTime time.Time `xorm:"datetime notnull comment('建立時間')"`
	UpdateTime time.Time `xorm:"datetime notnull comment('更新時間')"`
}
