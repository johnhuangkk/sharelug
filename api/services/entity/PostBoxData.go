package entity

import "time"

/**
i郵箱櫃體資料
"ADMId": "2179",
"ADMName": "布袋過溝郵局ｉ郵箱",
"ADMAlias": "布袋過溝郵局ｉ郵箱",
"ADMLocation": "大門左側",
"country": "嘉義縣",
"zip": "625",
"city": "布袋鎮",
"address": "中安里頂厝11號",
"Longitude": "120.182545",
"Latitude": "23.420347",
"POSTGOV_No": "600043"
*/

type IPOSTBoxData struct {
	ADMId       string `json:"ADMId"`
	ADMName     string `json:"ADMName"`
	ADMAlias    string `json:"ADMAlias"`
	ADMLocation string `json:"ADMLocation"`
	Country     string `json:"Country"`
	Zip         string `json:"Zip"`
	City        string `json:"City"`
	Address     string `json:"address"`
	Longitude   string `json:"Longitude"`
	Latitude    string `json:"Latitude"`
	GovNo       string `json:"POSTGOV_No"`
}

type IPostBoxArray []IPOSTBoxData

type PostBoxData struct {
	AdmId       string    `xorm:"pk varchar(8) comment('編號') "`
	AdmName     string    `xorm:"varchar(50) notnull comment('名稱')"`
	AdmAlias    string    `xorm:"varchar(50) notnull "`
	AdmLocation string    `xorm:"varchar(50) notnull comment('位置描述')"`
	Country     string    `xorm:"varchar(10) notnull comment('縣市')"`
	Zip         string    `xorm:"varchar(8)  notnull comment('郵遞區號')"`
	City        string    `xorm:"varchar(10) notnull comment('鄉鎮區')"`
	Address     string    `xorm:"varchar(50) notnull comment('地址')"`
	Longitude   string    `xorm:"varchar(20) notnull comment('經度')"`
	Latitude    string    `xorm:"varchar(20) notnull comment('緯度')"`
	GovNo       string    `xorm:"varchar(10) notnull comment('')"`
	BoxStatus   string    `xorm:"char(1) default 'Y' comment('啟用裝態 Y/N')"`
	UpdateTime  time.Time `xorm:"datetime notnull comment('更新時間')"`
}
