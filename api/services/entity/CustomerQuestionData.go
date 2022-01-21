package entity

import (
	"api/services/VO/CustomerVo"
	"time"
)

type CustomerQuestionData struct {
	Id           int64     `xorm:"pk int(10) unique autoincr comment('序號')"`
	Question     string    `xorm:"varchar(50) notnull comment('問題類別')"`
	QuestionType int64     `xorm:"tinyint(1) default 0 comment('是否有特殊欄位')"`
	Sort         int64     `xorm:"tinyint(3) notnull comment('掛序')"`
	CreateTime   time.Time `xorm:"datetime notnull comment('建立時間')"`
}

type CustomerData struct {
	Id          int64     `xorm:"pk int(10) unique autoincr comment('序號')"`
	QuestionId  string    `xorm:"varchar(20) unique comment('問題編號')"`
	RelatedId   string    `xorm:"varchar(20) comment('相關問題編號')"`
	Status      string    `xorm:"varchar(50) default 'none' comment('回復狀態')"`
	UserId      string    `xorm:"varchar(50) notnull comment('使用者ID')"`
	Question    string    `xorm:"varchar(50) notnull comment('問題')"`
	OrderId     string    `xorm:"varchar(50) notnull comment('訂單編號')"`
	Contents    string    `xorm:"text notnull comment('問題內容')"`
	ReplyTitle  string    `xorm:"text comment('回覆主題')"`
	Reply       string    `xorm:"text comment('回覆內容')"`
	ReplyTime   time.Time `xorm:"datetime comment('回覆時間')"`
	Remark      string    `xorm:"text comment('客服備註')"`
	RemarkTime  time.Time `xorm:"datetime comment('備註時間')"`
	RemarkStaff string    `xorm:"varchar(20) comment('處理人員')"`
	CreateTime  time.Time `xorm:"datetime notnull comment('建立時間')"`
	UpdateTime  time.Time `xorm:"datetime comment('更新時間')"`
}
type CustomerMemo struct {
	Id             int64     `xorm:"pk int(10) unique autoincr comment('序號')"`
	CustomerDataId int64     `xorm:"int(10) comment('問題ID')"`
	Content        string    `xorm:"text comment('備註')"`
	Staff          string    `xorm:"varchar(50) comment('備註者')"`
	CreateTime     time.Time `xorm:"timestamp created"`
}

func (c *CustomerData) GetCustomerData() CustomerVo.CustomerResponse {
	var data CustomerVo.CustomerResponse
	data.QuestionsDate = c.CreateTime.Format("2006/01/02 15:04")
	data.QuestionsTitle = c.Question
	data.QuestionsContents = c.Contents
	data.AnswerDate = c.ReplyTime.Format("2006/01/02 15:04")
	data.AnswerTitle = c.ReplyTitle
	data.AnswerContents = c.Reply
	return data
}

func (c *CustomerData) GetCustomerRemarkData() CustomerVo.CustomerRemarkResponse {
	var data CustomerVo.CustomerRemarkResponse
	data.RemarkDate = c.RemarkTime.Format("2006/01/02 15:04")
	data.RemarkStaff = c.RemarkStaff
	data.RemarkContents = c.Remark
	return data
}
