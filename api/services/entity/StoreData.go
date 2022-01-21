package entity

import (
	"api/services/VO/Response"
	"time"
)

type StoreData struct {
	StoreId          string    `xorm:"pk varchar(50) unique comment('商店ID')"`
	SellerId         string    `xorm:"varchar(50) notnull comment('賣家ID')"`
	StoreName        string    `xorm:"varchar(50) notnull comment('商店名稱')"`
	StoreTax         string    `xorm:"varchar(10) comment('商家統一編號')"`
	StorePicture     string    `xorm:"varchar(255) comment('商店頭貼')"`
	StoreDefault     int       `xorm:"int(1) notnull default 0 comment('是否為預設店')"`
	StoreStatus      string    `xorm:"varchar(20) notnull default 'SUCCESS' comment('賣場狀態')"`
	SubMerchantId    string    `xorm:"varchar(50) comment('次特店代碼')"`
	VerifyIdentity   int64     `xorm:"tinyint(1) default 0 comment('是否身份認證')"`
	ExpireTime       string    `xorm:"varchar(30) comment('商店到期日')"`
	FreeShipKey      string    `xorm:"varchar(20) default 'NONE' comment('免運費方式')"`
	FreeShip         int64     `xorm:"int(10) default 0 comment('免運費')"`
	SelfDelivery     bool      `xorm:"tinyint(1) default 0 comment('是否啟用外送')"`
	SelfDeliveryKey  string    `xorm:"varchar(20) default 'NONE' comment('外送免運方式')"`
	FreeSelfDelivery int64     `xorm:"int(10) default 0 comment('外送免運費')"`
	EnablePromo      bool      `xorm:"tinyint(1) default 0 comment('啟用優惠活動')"`
	CreateTime       time.Time `xorm:"datetime notnull comment('建立時間')"`
	UpdateTime       time.Time `xorm:"datetime notnull comment('更新時間')"`
}

func (s *StoreData) GetStoreLoginInfo() Response.StoreInfo {
	var si Response.StoreInfo
	si.Sid = s.StoreId
	si.Name = s.StoreName
	si.Picture = s.StorePicture
	return si
}

func (s *StoreDataResp) GetStoreInfo() Response.StoreInfo {
	var si Response.StoreInfo
	si.Sid = s.StoreId
	si.Name = s.StoreName
	si.Picture = s.StorePicture
	si.Rank = s.Rank
	return si
}

type StoreDataResp struct {
	StoreId       string
	StoreData     `xorm:"extends"`
	StoreRankData `xorm:"extends"`
}

type StoreUserData struct {
	StoreData  `xorm:"extends"`
	MemberData `xorm:"extends"`
}

type StoreFreeShip struct {
	FreeShipKey string `xorm:"varchar(20) default 'none' comment('免運費方式')"`
	FreeShip    int64  `xorm:"int(10) default 0 comment('免運費')"`
}

type StoreSocialMediaData struct {
	StoreId       string    `xorm:"pk unique varchar(50) notnull comment('收銀機ID')"`
	FacebookLink  string    `xorm:"varchar(500) comment('FB帳號連結')"`
	FacebookShow  int       `xorm:"tinyint(1) default 0 comment('FB帳號連結是否顯示')"`
	InstagramLink string    `xorm:"varchar(500) comment('IG帳號連結')"`
	InstagramShow int       `xorm:"tinyint(1) default 0 comment('IG帳號連結是否顯示')"`
	LineLink      string    `xorm:"varchar(500) comment('Line帳號連結')"`
	LineShow      int       `xorm:"tinyint(1) default 0 comment('Line帳號連結是否顯示')"`
	TelegramLink  string    `xorm:"varchar(500) comment('Tg帳號連結')"`
	TelegramShow  int       `xorm:"tinyint(1) default 0 comment('Tg帳號連結是否顯示')"`
	CreateTime    time.Time `xorm:"datetime notnull comment('建立時間')"`
	UpdateTime    time.Time `xorm:"datetime notnull comment('更新時間')"`
}

func (s *StoreSocialMediaData) GetStoreSocialMediaInfo() Response.StoreSocialMediaResponse {
	var ssm Response.StoreSocialMediaResponse
	ssm.Facebook.Link = s.FacebookLink
	ssm.Facebook.Show = s.FacebookShow
	ssm.Instagram.Link = s.InstagramLink
	ssm.Instagram.Show = s.InstagramShow
	ssm.Line.Link = s.LineLink
	ssm.Line.Show = s.LineShow
	ssm.Telegram.Link = s.TelegramLink
	ssm.Telegram.Show = s.TelegramShow
	return ssm
}
