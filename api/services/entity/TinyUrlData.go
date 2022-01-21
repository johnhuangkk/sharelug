package entity

type TinyUrlData struct {
	Number		string 	`xorm:"pk varchar(20) notnull"`
	ProductId	string	`xorm:"varchar(50) notnull"`
	Url 		string	`xorm:"varchar(100) notnull"`
}
