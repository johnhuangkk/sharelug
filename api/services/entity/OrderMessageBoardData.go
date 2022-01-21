package entity

import "time"

// 訂單留言板
type OrderMessageBoardData struct {
	Id          int       `xorm:"pk int(11) unique autoincr comment('序號')"`
	BuyerId     string    `xorm:"varchar(50) comment('買家ID')"`
	BuyerName   string    `xorm:"varchar(50) comment('購買人名稱')"`
	StoreId     string    `xorm:"varchar(50) comment('商店ID')"`
	StoreName   string    `xorm:"varchar(50) comment('商店名稱')"`
	OrderId     string    `xorm:"varchar(50) index notnull comment('訂單編號')"`
	MessageRole string    `xorm:"varchar(10) notnull comment('留言角色')"`
	Message     string    `xorm:"text comment('訊息')"`
	Reply       int       `xorm:"tinyint(1) default 0 comment('是否回覆')"`
	CreateTime  time.Time `xorm:"datetime notnull comment('建立時間')"`
	UpdateTime  time.Time `xorm:"datetime comment('更新時間')"`
}

func (om *OrderMessageBoardData) GetBuyerMessageData(member MemberData) {
	om.BuyerId = member.Uid
	om.BuyerName = member.Username
}

func (om *OrderMessageBoardData) GetStoreMessageData(store StoreDataResp) {
	om.StoreId = store.StoreId
	om.StoreName = store.StoreName
}

type OrderMessageBoardByOrderData struct {
	MessageBoard OrderMessageBoardData `xorm:"extends"`
	Order        OrderData             `xorm:"extends"`
}
