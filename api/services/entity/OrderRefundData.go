package entity

import (
	"api/services/Enum"
	"api/services/VO/Response"
	"time"
)

type OrderRefundData struct {
	RefundId      string    `xorm:"pk varchar(50) unique comment('退款單號')"`
	OrderId       string    `xorm:"varchar(50) notnull comment('訂單編號')"`
	RefundType    string    `xorm:"varchar(10) notnull comment('退款類別')"`
	Amount        float64   `xorm:"decimal(10,2) default 0.00 comment('退款金額')"`
	Total         float64   `xorm:"decimal(10,2) default 0.00 comment('合計退款金額')"`
	ProductSpecId string    `xorm:"varchar(50) comment('商品編號')"`
	ProductName   string    `xorm:"varchar(50) comment('商品名稱')"`
	Qty           int64     `xorm:"int(10) comment('數量')"`
	Sum           int64     `xorm:"int(10) comment('合計數量')"`
	Status        string    `xorm:"varchar(10) notnull comment('退款狀態')"`
	AuditStatus   string    `xorm:"varchar(10) default 'INIT' comment('審單狀態')"`
	RefundTime    time.Time `xorm:"datetime comment('退款退貨時間')"`
	CreateTime    time.Time `xorm:"datetime notnull comment('建立時間')"`
	UpdateTime    time.Time `xorm:"datetime notnull comment('更新時間')"`
}

func (r *OrderRefundData) GetRefund() Response.RefundResponse {
	var data Response.RefundResponse
	data.RefundDate = r.RefundTime.Format("2006/01/02 15:04")
	data.RefundId = r.RefundId
	data.Reason = ""
	data.Amount = int64(r.Amount)
	data.RefundStatus = r.Status
	data.RefundStatusText = Enum.RefundStatus[r.Status]
	return data
}

func (r *OrderRefundData) GetReturn() Response.ReturnResponse {
	var data Response.ReturnResponse
	data.ReturnDate = r.RefundTime.Format("2006/01/02 15:04")
	data.ReturnId = r.RefundId
	data.ReturnStatus = r.Status
	data.ReturnStatusText = Enum.ReturnStatus[r.Status]
	data.ProductName = r.ProductName
	data.ProductSpec = ""
	data.Amount = int64(r.Amount)
	data.Quantity = r.Sum
	return data
}