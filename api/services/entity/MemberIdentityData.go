package entity

import "time"

type MemberIdentityData struct {
	Uuid         string     `xorm:"pk varchar(50) unique"`
	Twid         string     `xorm:"varchar(50) notnull"`
	Name         string     `xorm:"varchar(10) notnull"`
	Verify       string     `xorm:"varchar(10) notnull"`
	CreateTime   time.Time  `xorm:"datetime notnull"`
	UpdateTime   time.Time  `xorm:"datetime notnull"`
}

