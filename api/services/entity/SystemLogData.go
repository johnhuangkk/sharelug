package entity

import "time"

//	- 應用程式亦需產生log，內容至少須包含:
//  - 使用者登出、登入紀錄（包含登入成功與失敗）
//  - 使用者存取信用卡資料的紀錄
//  - 系統管理者的所有行為
//  - log內容包含日期、時間、使用者帳號、事件行為、成功或失敗、受影響的資源等資訊

type SystemLog struct {
	Id           int64     `xorm:"pk int(10) unique autoincr comment('序號')"`
	UserId       string    `xorm:"varchar(50) notnull comment('使用者ID')"`
	LoginIp      string    `xorm:"varchar(30) notnull comment('登入IP')"`
	Action       string    `xorm:"varchar(20) notnull comment('行為')"`
	Content      string    `xorm:"varchar(200) notnull comment('行為內容')"`
	CreateTime   time.Time `xorm:"datetime notnull comment('建立時間')"`
}