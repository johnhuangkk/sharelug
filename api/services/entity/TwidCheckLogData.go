package entity

import "time"

type TwidCheckLogData struct {
	Id 	              int	     `xorm:"pk int(10) unique autoincr"`
	Uuid              string     `xorm:"varchar(50) notnull"`
	ApplyDate         string     `xorm:"varchar(10) notnull"`
	ApplyCode         string     `xorm:"varchar(10) notnull"`
	ApplyDist         string     `xorm:"varchar(10) notnull"`
	HttpCode          string     `xorm:"varchar(50)"`
	HttpMessage       string     `xorm:"varchar(50)"`
	RdCode            string     `xorm:"varchar(50)"`
	RdMessage         string     `xorm:"varchar(50)"`
	CheckIdCardApply  string     `xorm:"varchar(10)"`
	Type              string     `xorm:"varchar(50) notnull"`
	UpdateTime        time.Time  `xorm:"datetime notnull"`
}


