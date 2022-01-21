package entity

type ScsBankAccountData struct {
	Id 		    int		   `xorm:"pk bigint unique autoincr comment('序號')"`
	TxKey       string     `xorm:"varchar(40) unique"`  // 收款帳號 + 索引日 + 序號 組合而成的
	RAccount    string     `xorm:"varchar(20)"`  // 收款帳號
	BtxDate     string     `xorm:"varchar(8) "`  // 索引日
	SeqNo       string     `xorm:"varchar(8) "`  // 序號
	SAccount    string     `xorm:"varchar(20)"`  // 轉出帳號
	BusDate     string     `xorm:"varchar(8)"`   // 營業日
	TxDate      string     `xorm:"varchar(8)"`   // 交易日期yyyymmdd
	TxTime      string     `xorm:"varchar(8)"`   // 交易時間hhmmss
	VAccount    string     `xorm:"varchar(20)"`  // 虛擬帳號
	Amount      int64      `xorm:"int(10)"`      // 繳款金額
	DC          string     `xorm:"varchar(5)"`   // 借貸別
	Notes       string     `xorm:"varchar(255)"` // 附註(交易類別)
	CheckSum    string     `xorm:"varchar(10)"`  // 檢核碼
	IsValid     int        `xorm:"tinyint(1)"`   // 檢核碼檢核是否有效
	RawData     string     `xorm:"text"`         // API提供的原始資料
}
