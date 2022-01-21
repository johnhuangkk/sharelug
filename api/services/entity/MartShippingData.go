package entity

type MartHiLifeShippingData struct {
	Id         int64                  `xorm:"pk int(10) unique autoincr"`
	ShipNo     string                 `xorm:"varchar(20) unique comment('運送單號')"`
	EcOrderNo  string                 `xorm:"varchar(20) comment('購買訂單代號')"`
	ParentId   string                 `xorm:"varchar(5) comment('母特店代號')"`
	EshopId    string                 `xorm:"varchar(5) comment('子特店代號')"`
	RrStoreId  string                 `xorm:"varchar(10) comment('收件超商編號')"`
	RrName     string				  `xorm:"varchar(30) comment('收件人姓名')"`
	RrPhone    string				  `xorm:"varchar(16) comment('收件人電話')"`
	SrStoreId  string                 `xorm:"varchar(10) comment('送件超商編號')"`
	SrName     string				  `xorm:"varchar(30) comment('送件人姓名')"`
	SrPhone    string				  `xorm:"varchar(16) comment('送件人電話')"`
	State      string                 `xorm:"varchar(255) comment('前端顯示運送狀態')"`
	StateCode  string                 `xorm:"varchar(10)  comment('運送狀態碼')"`
	Record     string                 `xorm:"text comment('前端顯示運送日誌')"`
	Log        string                 `xorm:"text comment('系統存入Log')"`
	NeedPay    bool                   `xorm:"tinyint(1) NOT NULL comment('是否取貨付款')"`
	NeedChange bool                   `xorm:"tinyint(1) NOT NULL comment('是否需要改店')"`
	Amount     string                 `xorm:"int(11) comment('商品金額')"`
	ShipFee    string                 `xorm:"int(11) comment('運費')"`
	IsLose     bool                   `xorm:"tinyint(1) NOT NULL comment('商品遺失')"`
	OnReturn   bool					  `xorm:"tinyint(1) NOT NULL comment('退貨流程中')"`
	CreateDT   string			      `xorm:"varchar(30) comment('寄送建立時間')"`
	UpdateDT   string			      `xorm:"varchar(30) comment('狀態更新時間')"`
	SwitchDT   string			      `xorm:"varchar(30) comment('閉轉期限')"`
}

type MartOkShippingData struct {
	Id         int64                  `xorm:"pk int(10) unique autoincr"`
	ShipNo     string                 `xorm:"varchar(20) unique comment('運送單號')"`
	EcOrderNo  string                 `xorm:"varchar(20) comment('購買訂單代號')"`
	ParentId   string                 `xorm:"varchar(5) comment('母特店代號')"`
	EshopId    string                 `xorm:"varchar(5) comment('子特店代號')"`
	RrStoreId  string                 `xorm:"varchar(10) comment('收件超商編號')"`
	RrName     string				  `xorm:"varchar(30) comment('收件人姓名')"`
	RrPhone    string				  `xorm:"varchar(16) comment('收件人電話')"`
	SrStoreId  string                 `xorm:"varchar(10) comment('送件超商編號')"`
	SrName     string				  `xorm:"varchar(30) comment('送件人姓名')"`
	SrPhone    string				  `xorm:"varchar(16) comment('送件人電話')"`
	State      string                 `xorm:"varchar(255) comment('前端顯示運送狀態')"`
	StateCode  string                 `xorm:"varchar(10)  comment('運送狀態碼')"`
	Record     string                 `xorm:"text comment('前端顯示運送日誌')"`
	Log        string                 `xorm:"text comment('系統存入Log')"`
	NeedPay    bool                   `xorm:"tinyint(1) NOT NULL comment('是否取貨付款')"`
	NeedChange bool                   `xorm:"tinyint(1) NOT NULL comment('是否需要改店')"`
	Amount     string                 `xorm:"int(11) comment('商品金額')"`
	ShipFee    string                 `xorm:"int(11) comment('運費')"`
	IsLose     bool                   `xorm:"tinyint(1) NOT NULL comment('商品遺失')"`
	OnReturn   bool					  `xorm:"tinyint(1) NOT NULL comment('退貨流程中')"`
	CreateDT   string			      `xorm:"varchar(30) comment('寄送建立時間')"`
	UpdateDT   string			      `xorm:"varchar(30) comment('狀態更新時間')"`
	SwitchDT   string			      `xorm:"varchar(30) comment('閉轉期限')"`
}

type MartFamilyShippingData struct {
	Id         int64                  `xorm:"pk int(10) unique autoincr"`
	ShipNo     string                 `xorm:"varchar(20) unique comment('運送單號')"`
	EcOrderNo  string                 `xorm:"varchar(20) comment('購買訂單代號')"`
	ParentId   string                 `xorm:"varchar(5) comment('母特店代號')"`
	EshopId    string                 `xorm:"varchar(5) comment('子特店代號')"`
	RrStoreId  string                 `xorm:"varchar(10) comment('收件超商編號')"`
	RrName     string				  `xorm:"varchar(30) comment('收件人姓名')"`
	RrPhone    string				  `xorm:"varchar(16) comment('收件人電話')"`
	SrStoreId  string                 `xorm:"varchar(10) comment('送件超商編號')"`
	SrName     string				  `xorm:"varchar(30) comment('送件人姓名')"`
	SrPhone    string				  `xorm:"varchar(16) comment('送件人電話')"`
	State      string                 `xorm:"varchar(255) comment('前端顯示運送狀態')"`
	StateCode  string                 `xorm:"varchar(10)  comment('運送狀態碼')"`
	Record     string                 `xorm:"text comment('前端顯示運送日誌')"`
	Log        string                 `xorm:"text comment('系統存入Log')"`
	NeedPay    bool                   `xorm:"tinyint(1) NOT NULL comment('是否取貨付款')"`
	NeedChange bool                   `xorm:"tinyint(1) NOT NULL comment('是否需要改店')"`
	Amount     string                 `xorm:"int(11) comment('商品金額')"`
	ShipFee    string                 `xorm:"int(11) comment('運費')"`
	IsLose     bool                   `xorm:"tinyint(1) NOT NULL comment('商品遺失')"`
	OnReturn   bool					  `xorm:"tinyint(1) NOT NULL comment('退貨流程中')"`
	CreateDT   string			      `xorm:"varchar(30) comment('寄送建立時間')"`
	UpdateDT   string			      `xorm:"varchar(30) comment('狀態更新時間')"`
	SwitchDT   string			      `xorm:"varchar(30) comment('閉轉期限')"`
}