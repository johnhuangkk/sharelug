package entity

type MartOkStoreData struct {
	StoreId 		string		`xorm:"pk varchar(10) notnull unique"`
	StoreName 		string		`xorm:"varchar(255) notnull"`
	StoreAddress 	string		`xorm:"text notnull"`
	StoreCloseDate  string		`xorm:"varchar(10)"`
	MdcStareDate	string		`xorm:"varchar(10)"`
	MdcEndDate		string		`xorm:"varchar(10)"`
	Route		    string		`xorm:"varchar(10)"`
	Step 		    string  	`xorm:"varchar(10)"`
	TelNo 		    string  	`xorm:"varchar(20)"`
	OldStore 		string  	`xorm:"varchar(10)"`
	Area 		    string  	`xorm:"varchar(10)"`
	EquipmentId     string  	`xorm:"varchar(255)"`
	City            string      `xorm:"varchar(255)"`
	District        string      `xorm:"varchar(255)"`
}
