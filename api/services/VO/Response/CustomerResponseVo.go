package Response

type CustomerResponse struct {
	Question string `json:"Question"`
	Type     int64  `json:"Type"`
}

type CustomerQuestionResponse struct {
	QuestionId   string `json:"QuestionId"`
	RelatedId    string `json:"RelatedId"`
	Question     string `json:"Question"`
	OrderId      string `json:"OrderId"`
	Content      string `json:"Content"`
	CreateTime   string `json:"CreateTime"`
	ReplyContent string `json:"ReplyContent,omitempty"`
	ReplyTime    string `json:"ReplyTime,omitempty"`
}

type CustomerHistoryQuestionResponse struct {
	QuestionId       string                     `json:"QuestionId"`
	RelatedId        string                     `json:"RelatedId"`
	Question         string                     `json:"Question"`
	OrderId          string                     `json:"OrderId"`
	Content          string                     `json:"Content"`
	CreateTime       string                     `json:"CreateTime"`
	ReplyContent     string                     `json:"ReplyContent,omitempty"`
	ReplyTime        string                     `json:"ReplyTime,omitempty"`
	RelatedQuestions []CustomerQuestionResponse `json:"RelatedQuestions,omitempty"`
}
