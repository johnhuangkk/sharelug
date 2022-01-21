package entity

type Sequence struct {
	Name 			string `xorm:"pk varchar(50) notnull"`
	CurrentValue	int `xorm:"int(11)"`
	Increment		int	`xorm:"int(11)"`
}
