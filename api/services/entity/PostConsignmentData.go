package entity

import "time"

// Consignment
type PostConsignmentData struct {
	OrderId     string `xorm:"pk unique varchar(50) notnull comment('訂單編號')"`
	MerchantId  string `xorm:"varchar(30) notnull comment('合作代號')"`
	ShipNumber  string `xorm:"varchar(20) notnull comment('郵寄編號')"`
	SellerId    string `xorm:"varchar(50) notnull comment('賣家Id')"`
	SellerName  string `xorm:"varchar(100) notnull comment('賣家名稱')"`
	SellerPhone string `xorm:"varchar(20) notnull comment('賣家電話')"`
	SellerZip   string `xorm:"varchar(10) notnull comment('賣家郵遞區號')"`
	SellerAddr  string `xorm:"varchar(100) notnull comment('賣家地址')"`
	CreateTime  string `xorm:"datetime notnull comment('建立時間')"`
}
type ConsignmentNote struct {
	OrderId, MerchantId, ShipNumber, ShipType      string
	SellerName, SellerPhone, SellerZip, SellerAddr string
	ReceiverName, ReceiverAddress, ReceiverPhone   string
}

type PostBagConsignmentData struct {
	OrderId        string    `xorm:"pk unique varchar(50) notnull comment('訂單編號')"`
	MerchantId     string    `xorm:"varchar(30) notnull comment('合作代號')"`
	ShipNumber     string    `xorm:"varchar(20) notnull comment('郵寄編號')"`
	SellerId       string    `xorm:"varchar(50) notnull comment('賣家Id')"`
	SellerName     string    `xorm:"varchar(50) notnull comment('賣家名稱')"`
	SellerPhone    string    `xorm:"varchar(20) notnull comment('賣家電話')"`
	SellerZip      string    `xorm:"varchar(10) notnull comment('賣家郵遞區號')"`
	SellerAddr     string    `xorm:"varchar(100) notnull comment('賣家地址')"`
	FileName       string    `xorm:"varchar(100) comment('郵件號驗證檔案名')"`
	VerifyStatus   bool      `xorm:"tinyint(1) comment('驗證狀態')"`
	VerifyTime     time.Time `xorm:"timestamp comment('郵件編號驗證時間')"`
	VerifyFileName string    `xorm:"varchar(100) comment('已驗證之檔名')"`
	CreateTime     time.Time `xorm:"datetime created notnull comment('建立時間')"`
}

type PostBagFileData struct {
	Type       string    `xorm:"varchar(30) notnull comment('')"`
	Date       string    `xorm:"varchar(15) notnull comment('')"`
	List       string    `xorm:"text notnull comment('')"`
	CreateTime time.Time `xorm:"datetime created notnull comment('建立時間')"`
}
