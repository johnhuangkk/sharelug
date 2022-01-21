package entity

import "time"

type DeviceData struct {
	DeviceUuid 		string 		`xorm:"pk varchar(100) unique"` //裝置ID
	Token			string		`xorm:"varchar(100)"`
	Platform		string		`xorm:"text"` 		//作業系統
	PlatformVersion	string		`xorm:"varchar(50)"`		//作業系統版本
	PlatformDevice	string		`xorm:"varchar(50)"`		//裝置廠牌和型號
	PlatformIP		string		`xorm:"varchar(30)"`		//裝置IP
	FcmToken		string		`xorm:"varchar(100)"`		//FcmToken
	CreateTime		time.Time	`xorm:"datetime notnull"`	//建立時間
	UpdateTime		time.Time	`xorm:"datetime notnull"`	//更新時間
}
