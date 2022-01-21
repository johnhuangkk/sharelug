package entity

import (
	"api/services/VO/Response"
	"time"
)

type OrderDetail struct {
	Id              int64     `xorm:"pk int(10) unique autoincr comment('序號')"`
	SellerId        string    `xorm:"varchar(50) notnull comment('賣家ID')"`
	OrderId         string    `xorm:"varchar(50) notnull comment('訂單編號')"`
	ProductSpecId   string    `xorm:"varchar(50) notnull comment('商品編號')"`
	ProductSpecName string    `xorm:"varchar(50) comment('商品規格')"`
	ProductName     string    `xorm:"varchar(100) notnull comment('商品名稱')"`
	ProductQuantity int64     `xorm:"int(10) notnull comment('商品數量')"`
	ProductPrice    int64     `xorm:"int(10) notnull comment('商品單價')"`
	ShipMerge       int64     `xorm:"tinyint(1) default 0 comment('是否合併運費')"`
	ShipFee         int64     `xorm:"int(10) notnull comment('運費')"`
	IsFreeShip      bool      `xorm:"tinyint(1) default 0 comment('是否免運費')"`
	Subtotal        int64     `xorm:"int(10) notnull comment('小計')"`
	CreateTime      time.Time `xorm:"datetime notnull comment('建立時間')"`
}

func (d *OrderDetail) GetOrderDetail() Response.SearchOrderDetail {
	var data Response.SearchOrderDetail
	data.ShipFee = d.ShipFee
	data.ProductName  = d.ProductName
	data.ProductSpec = d.ProductSpecName
	data.Amount = d.ProductPrice
	data.Quantity = d.ProductQuantity
	data.IsMerge = d.ShipMerge
	return data
}