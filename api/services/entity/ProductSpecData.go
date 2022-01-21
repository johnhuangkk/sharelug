package entity

import "time"

type ProductSpecData struct {
	SpecId       string `xorm:"pk varchar(50) notnull comment('商品規格編號')"`
	ProductId    string `xorm:"varchar(50) notnull comment('商品編號')"`
	SpecName     string `xorm:"varchar(100) comment('商品規格名稱')"`
	Quantity 	 int64  `xorm:"int(10) notnull comment('商品規格數量')"`
	SpecPrice    int64  `xorm:"int(10) notnull comment('商品規格價格')"`
	SpecStatus   string `xorm:"varchar(10) notnull comment('商品規格狀態')"`
	CreateTime   time.Time `xorm:"datetime notnull comment('建立時間')"`
}
