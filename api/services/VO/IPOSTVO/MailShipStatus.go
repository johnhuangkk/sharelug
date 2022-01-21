package IPOSTVO

type MailShip struct {
	MailNo     string `form:"MailNo" json:"MailNo" validate:"required"`
	CreateTime string `form:"CreateTime" json:"CreateTime" validate:"required"`
}
