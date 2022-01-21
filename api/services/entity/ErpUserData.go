package entity

type ErpUser struct {
	Uid      string `xorm:" pk unique varchar(50) comment('使用者ID')"`
	Name     string `xorm:" varchar(50) comment('使用者名稱')"`
	Email    string `xorm:" unique varchar(50) comment('使用者信箱')"`
	Password string `xorm:"varchar(100)"`
	Created  string `xorm:"timestamp created notnull comment('建立時間')"`
	Enable   bool   `xorm:"tinyint(1) comment('是否啟用')"`
}
