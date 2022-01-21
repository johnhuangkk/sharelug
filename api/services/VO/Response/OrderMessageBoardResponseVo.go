package Response

type OrderMessageBoardResponse struct {
	Tabs        Tabs          `json:"Tabs"`
	MessageList []MessageList `json:"MessageList"`
	BoardCount  int64         `json:"BoardCount"`
}

type MessageList struct {
	OrderId        string `json:"OrderId"`
	OrderTime      string `json:"OrderTime"`
	MessageTime    string `json:"MessageTime"`
	MessageContent string `json:"MessageContent"`
	MessageReply   string `json:"MessageReply"`
}

type Tabs struct {
	Inbox   int64 `json:"Inbox"`
	Replied int64 `json:"Replied"`
}
