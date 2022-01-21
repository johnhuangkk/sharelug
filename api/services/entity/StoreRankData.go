package entity

import "time"

type StoreRankData struct {
	RankId     int       `xorm:"pk int(10) unique autoincr comment('序號')"`
	StoreId    string    `xorm:"varchar(50) notnull comment('StoreId')"`
	UserId     string    `xorm:"varchar(50) notnull comment('UserId')"`
	Rank       string    `xorm:"varchar(10) notnull comment('身份')"`
	RankStatus string    `xorm:"varchar(10) notnull comment('狀態')"`
	Email      string    `xorm:"varchar(100) comment('信箱')"`
	CreateTime time.Time `xorm:"datetime notnull comment('建立時間')"`
	UpdateTime time.Time `xorm:"datetime notnull comment('更新時間')"`
}

type StoreList struct {
	StoreId       string
	StoreRankData `xorm:"extends"`
	StoreData     `xorm:"extends"`
}

type StoreRankResp struct {
	StoreRank StoreRankData `xorm:"extends"`
	Member    MemberData    `xorm:"extends"`
}

type StoreAndStoreRankEnt struct {
	StoreData StoreData     `xorm:"extends"`
	StoreRank StoreRankData `xorm:"extends"`
}


type StoreWithRank struct {
	SellerId  string
	Mphone    string
	UserName  string
	StoreName string
	StoreRank StoreRankData `xorm:"extends"`
}
