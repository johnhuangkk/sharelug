package entity

import "time"

type TaiwanArea struct {
	Id       int64     `xorm:"pk int(10) autoincr"`
	CityCode string    `xorm:"varchar(15) NOT NULL comment('縣市代碼')"`
	CityName string    `xorm:"varchar(15) NOT NULL comment('縣市名稱')"`
	ZipCode  string    `xorm:"varchar(5) NOT NULL comment('地區郵遞區號')"`
	AreaName string    `xorm:"varchar(20) NOT NULL comment('地區名稱')"`
	Created  time.Time `xorm:"timestamp created"`
	Updated  time.Time `xorm:"timestamp updated"`
}

type TaiwanCity struct {
	Id      int64     `xorm:"pk int(10) autoincr"`
	Name    string    `xorm:"varchar(15) NOT NULL comment('縣市名稱')"`
	Code    string    `xorm:"varchar(15) NOT NULL comment('縣市代碼')"`
	NameEn  string    `xorm:"varchar(15) NOT NULL comment('縣市名稱英')"`
	Created time.Time `xorm:"timestamp created"`
}
