package entity

import "time"

// 使用者寄送地址資訊
type UserAddressData struct {
	UaId       string    `xorm:"pk varchar(50) unique "`                              // timestamp md5
	Uid        string    `xorm:"varchar(50) notnull"`                                 //
	Ship       string    `xorm:"varchar(30) notnull comment('運送方式')"`                 //運送方式
	Address    string    `xorm:"varchar(200) notnull"`                                //地址
	Status     string    `xorm:"char(1) notnull default 'Y' comment('狀態 啟用 Y 停用 N')"` //狀態 啟用 1 刪除 0
	Type       string    `xorm:"char(1) notnull default 'R' comment('類型 寄件 S 收件 R')"` //類型 寄件 S 收件 R
	RealName   string    `xorm:"varchar(20)  comment('真實姓名')"`                        //真實姓名
	Phone      string    `xorm:"varchar(20)  comment('收件聯絡電話')"`                      //真實姓名
	UpdateTime time.Time `xorm:"datetime notnull"`                                    //建立時間
}
