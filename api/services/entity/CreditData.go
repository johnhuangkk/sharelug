package entity

import "time"

type KgiCredit struct {
	MerchantID  string `form:"MerchantID"`  //商店代號
	TerminalID  string `form:"TerminalID"`  //端末機代號
	OrderID     string `form:"OrderID"`     //訂單編號
	TransCode   string `form:"TransCode"`   //交易代碼
	PAN         int    `form:"PAN"`         //交易卡號
	ExpireDate  int    `form:"ExpireDate"`  //卡片到期日
	ExtenNo     int    `form:"EXTENNO"`     //CVV2
	TransMode   string `form:"TransMode"`   //交易類別
	Install     int    `form:"Install"`     //分期期數
	TransAmt    int    `form:"TransAmt"`    //交易金額
	NotifyURL   string `form:"NotifyURL"`   //3D 回應網址
	PrivateData string `form:"PrivateData"` //自訂資料
}

type KgiSpecialStore struct {
	Id            int       `xorm:"pk int(10) unique autoincr comment('序號')"`
	MerchantId    string    `xorm:"varchar(12) notnull comment('次特店代碼')"`
	Terminal3dId  string    `xorm:"varchar(12) notnull comment('3D端末機代號')"`
	Terminaln3dId string    `xorm:"varchar(12) notnull comment('N3D端末機代號')"`
	UserId        string    `xorm:"varchar(50) comment('會員ID')"`
	IsUsed        bool      `xorm:"tinyint(1) default 0 comment('是否已使用')"`
	ChStoreName   string    `xorm:"varchar(50) comment('次特店中文名稱')"`
	EnStoreName   string    `xorm:"varchar(50) comment('次特店英文名稱')"`
	CreateTime    time.Time `xorm:"datetime notnull comment('建立時間')"`
	UpdateTime    time.Time `xorm:"datetime notnull comment('更新時間')"`
}
