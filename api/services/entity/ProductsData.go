package entity

import (
	"api/services/VO/Response"
	"api/services/util/qrcode"
	"api/services/util/tools"
	"time"
)

type ProductData struct {
	ProductId     string    `xorm:"pk varchar(50) unique comment('商品編號')"`
	ProductName   string    `xorm:"varchar(100) notnull comment('商品名稱')"`
	StoreId       string    `xorm:"varchar(50) notnull comment('店家ID')"`
	IsRealtime    int       `xorm:"tinyint(1) notnull default 1 comment('是否為即時帳單')"`
	IsSpec        int       `xorm:"tinyint(1) notnull default 1 comment('是否有規格')"`
	Price         int64     `xorm:"int(10) notnull comment('商品價格')"`
	Shipping      string    `xorm:"text comment('運送方式')"`
	ShipMerge     int       `xorm:"tinyint(1) default 0 comment('是否合併運費')"`
	PayWay        string    `xorm:"text comment('付款方式')"`
	QrCode        string    `xorm:"varchar(100) comment('Qrcode')"`
	TinyUrl       string    `xorm:"varchar(100) comment('短網址')"`
	ProductStatus string    `xorm:"varchar(10) notnull comment('商品狀態')"`
	Stock         int64     `xorm:"int(10) default 0 comment('商品總庫存')"`
	FormUrl       string    `xorm:"varchar(250) comment('表單網址')"`
	LimitKey      string    `xorm:"varchar(20) default 'NONE' notnull comment('限制')"`
	LimitQty      int64     `xorm:"int(10) default 0 notnull comment('限制數量')"`
	IsFreeShip    bool      `xorm:"tinyint(1) default 1 comment('是否免運費')"`
	ExpireDate    time.Time `xorm:"datetime comment('即時帳單到期日期')"`
	CreateTime    time.Time `xorm:"datetime notnull comment('建立時間')"`
	UpdateTime    time.Time `xorm:"datetime notnull comment('更新時間')"`
}

type ProductsData struct {
	ProductId       string
	ProductSpecData `xorm:"extends"`
	ProductData     `xorm:"extends"`
}

type NewShippingMode struct {
	Type   string `json:"Type"`
	Text   string `json:"Text"`
	Price  int    `json:"Price"`
	Remark string `json:"Remark"`
}

type NewProductResponse struct {
	ProductId string `json:"ProductId"`
}

func (p *ProductData) GetSearchProduct(userData MemberData, store StoreData) Response.SearchProductResponse {
	var data Response.SearchProductResponse
	data.CreateDate = p.CreateTime.Format("2006/01/02 15:04")
	data.UpdateDate = p.UpdateTime.Format("2006/01/02 15:04")
	data.UserAccount = userData.Mphone //拿掉隱碼
	data.StoreName = store.StoreName
	data.ProductId = p.ProductId
	data.ProductName = p.ProductName
	data.Amount = p.Price
	data.Quantity = p.Stock
	data.ShipMerge = "N"
	if p.ShipMerge == 1 {
		data.ShipMerge = "Y"
	}
	data.GoogleFormLink = ""
	data.ProductLink = qrcode.GetTinyUrl(p.TinyUrl)
	return data
}

func (p *ProductData) GetRealtimeProduct(userData MemberData, store StoreData) Response.SearchProductResponse {
	var data Response.SearchProductResponse
	data.CreateDate = p.CreateTime.Format("2006/01/02 15:04")
	data.UserAccount = tools.MaskerPhoneLater(userData.Mphone)
	data.StoreName = store.StoreName
	data.ProductId = p.ProductId
	data.ProductName = p.ProductName
	data.Amount = p.Price
	data.GoogleFormLink = ""
	data.ProductLink = qrcode.GetTinyUrl(p.TinyUrl)
	return data
}
