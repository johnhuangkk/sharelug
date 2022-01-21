package entity

import (
	"time"
)

// 超商配送資料
type CvsShippingData struct {
	Id            int64     `xorm:"pk int(10) unique autoincr comment('序號')"`
	ParentId      string    `xorm:"varchar(5) comment('母特店代號')"`
	ShipNo        string    `xorm:"varchar(20) unique comment('運送單號')"`
	CvsType       string    `xorm:"varchar(20) comment('超商')"`
	EcOrderNo     string    `xorm:"varchar(20) unique comment('訂單編號')"`
	ServiceType   string    `xorm:"char(1) default '0' NOT NULL comment('是否取貨付款： 1 取貨付款 / 0 取貨不付款')"`
	SenderName    string    `xorm:"varchar(30) comment('送件人姓名')"`
	SenderPhone   string    `xorm:"varchar(16) comment('送件人電話')"`
	SenderStoreId string    `xorm:"varchar(16) comment('寄件店號為空值：若有值代表退回賣家時發生閉轉店，需求無法退回原寄件店')"`
	SendTime      time.Time `xorm:"timestamp comment('到店交寄時間')"`

	OriReceiverAddress    string    `xorm:"varchar(16) NOT NULL comment('原收件地址')"`
	SwitchReceiverAddress string    `xorm:"varchar(16) comment('閉轉收件地址')"`
	ReceiveTime           time.Time `xorm:"timestamp comment('到店取件時間')"`
	// 參照 超商異常狀態.md
	StateCode      string    `xorm:"varchar(6) default '1' comment('運送狀態碼')"`
	Switch         string    `xorm:"char(1) default '0' comment('是否需要閉轉店與StateCode一起看：  1 需要 / 0 不需要')"`
	SwitchDeadline string    `xorm:"varchar(20) comment('閉轉店期限')"`
	FlowType       string    `xorm:"char(1) default 'N' comment('進退貨狀態： Ｎ進貨 / R退貨')"`
	CreateTime     time.Time `xorm:"datetime comment('建立時間')"`
	UpdateTime     time.Time `xorm:"datetime comment('更新時間')"`
	CheckSend      bool      `xorm:"char(1) default '0' comment('寄件核帳：  1 核帳 / 0 未核帳')"`
}
type CvsShippingWithAmount struct {
	CvsShippingData `xorm:"extends"`
	Amount          float64
}

func (data *CvsShippingData) InitInsert(cvsType string) {
	data.SenderStoreId = ``
	data.StateCode = `1`
	data.Switch = `0`
	data.FlowType = `N`
	data.CvsType = cvsType
	data.SwitchDeadline = ``
	data.CreateTime = time.Now()
	data.UpdateTime = time.Now()
}
