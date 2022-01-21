package entity

import "time"

type ProductImagesData struct {
	Id          int       `xorm:"pk int(10) unique autoincr comment('序號')"`
	ProductId   string    `xorm:"varchar(50) notnull comment('商品編號')"`
	Image       string    `xorm:"varchar(100) notnull comment('商品圖')"`
	ImageSeq    int       `xorm:"int(10) notnull comment('排序')"`
	ImageStatus string    `xorm:"varchar(20) notnull comment('狀態')"`
	CreateTime  time.Time `xorm:"datetime notnull comment('建立時間')"`
}
