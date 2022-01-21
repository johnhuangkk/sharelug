package CustomerVo

type CustomerResponse struct {
	QuestionsDate     string `json:"QuestionsDate"`     //來信日期
	QuestionsTitle    string `json:"QuestionsTitle"`    //來信主旨
	QuestionsContents string `json:"QuestionsContents"` //來信內容
	AnswerDate        string `json:"AnswerDate"`        //回覆日期
	AnswerTitle       string `json:"AnswerTitle"`       //回覆主旨
	AnswerContents    string `json:"AnswerContents"`    //回覆內容
}

type CustomerRemarkResponse struct {
	RemarkDate     string `json:"RemarkDate"`     //日期
	RemarkStaff    string `json:"RemarkStaff"`    //處理人員
	RemarkContents string `json:"RemarkContents"` //內容
}
