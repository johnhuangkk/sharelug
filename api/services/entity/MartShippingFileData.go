package entity

type MartShippingFileData struct {
	Id          int    `xorm:"pk int(32) autoincr"`
	Vendor      string `xorm:"varchar(20) unique(oi)  comment('供應商，如萊爾富、全家等')"`
	ParentId    string `xorm:"varchar(5) unique(oi) comment('母特店代號')"`
	EshopId     string `xorm:"varchar(5) unique(oi) comment('子特店代號')"`
	FileId      string `xorm:"varchar(255) unique(oi)  comment('檔案Id')"`
	FileType 	string `xorm:"varchar(20) comment('購買訂單代號')"`
	FileDate    string `xorm:"varchar(20) comment('檔案日期 YYYY/MM/DD-HH')"`
	Content     string `xorm:"text comment('內容')"`
}
