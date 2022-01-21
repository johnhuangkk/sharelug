package Request

type OrderMessageBoardRequest struct {
	Tabs    string `form:"Tabs" json:"Tabs"`
	OrderBy string `form:"OrderBy" json:"OrderBy"` //排序
	Limit   int    `form:"Limit" json:"Limit"`
	Start   int    `form:"Start" json:"Start"`
}
