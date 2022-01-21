package entity

import "time"

type SevenMyshipShopData struct {
	StoreID   string    `xorm:"pk varchar(10) notnull 'store_id' comment('店鋪編號')" json:"storeId"`
	StoreName string    `xorm:"varchar(255) notnull  'store_name' comment('店鋪名')" json:"storeName"`
	Country   string    `xorm:"'country' varchar(10) notnull  comment('縣市名')"  json:"-"`
	District  string    `xorm:"'district' varchar(15) notnull  comment('行政區名')"  json:"-"`
	Address   string    `xorm:"'address' text  notnull comment('店鋪地址')" json:"address"`
	Opened    bool      `xorm:"'opened' tinyint notnull index(CountryDistrictOpened) comment('是否營業')" json:"-"`
	Created   time.Time `xorm:"timestamp created" json:"-"`
	Updated   time.Time `xorm:"timestamp updated" json:"-"`
}

type SevenShops struct {
	Shops []SevenMyshipShopData
}
type SevenChargeOrderData struct {
	From        string    `xorm:"not null varchar(30) 'from' comment('來源號')"`
	To          string    `xorm:"not null varchar(30) 'to' comment('目的號')"`
	TermiNo     string    `xorm:"not null varchar(30) 'termi_no' comment('交易序號')"`
	Date        string    `xorm:"not null varchar(10) 'record_date' comment('電文產生日')"`
	Time        string    `xorm:"not null varchar(10) 'record_time' comment('電文產生時間')"`
	StatCode    string    `xorm:"not null varchar(5)  'record_status' comment('電文產生狀態')"`
	StatDesc    string    `xorm:"varchar(30)  'status_desc' comment('電文狀態敘述')"`
	SequenceNo  string    `xorm:"not null varchar(2)  'sequence_no' comment('交易次序')"`
	OLOiNo      string    `xorm:"not null varchar(5)  'ol_oi_no' comment('代收碼')"`
	OLCode1     string    `xorm:"not null varchar(10) 'ol_code_1' comment('第一段條碼')"`
	OLCode2     string    `xorm:"not null varchar(16) 'ol_code_2' comment('第二段條碼')"`
	OLCode3     string    `xorm:"not null varchar(16) 'ol_code_3' comment('第三段條碼')"`
	OLAmount    string    `xorm:"not null varchar(16) 'ol_amount' comment('交易金額')"`
	Status      string    `xorm:"not null varchar(1)  'ol_status' comment('處理結果')"`
	Description string    `xorm:" varchar(30) 'ol_description' comment('處理結果')"`
	OLPrint     string    `xorm:"varchar(2) 'ol_print' comment('列印結果')"`
	Created     time.Time `xorm:"created" json:"-"`
}

type SevenShipMapData struct {
	OrderId           string    `xorm:"pk varchar(15) notnull unique comment('訂單編號') 'order_id'"`
	PaymentNoWithCode string    `xorm:"pk varchar(16) not null unique comment('交貨便編號+驗證碼') 'paymentno_with_code'"`
	ShipNo            string    `xorm:"varchar(16) comment('運送編號') 'ship_no'"`
	VerifyCode        string    `xorm:"varchar(4) notnull comment('驗證碼') 'verify_code'"`
	PaymentNo         string    `xorm:"varchar(12) comment('交貨便編號')"`
	PayWay            string    `xorm:"varchar(15) notnull comment('付款方式') 'pay_way'"`
	Created           time.Time `xorm:"created" json:"-"`
	Updated           time.Time `xorm:"updated" json:"-"`
}
