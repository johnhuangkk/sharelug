package Request

type CustomerRequest struct {
	Question   string `form:"Question" json:"Question"`
	OrderId    string `form:"OrderId" json:"OrderId"`
	Contents   string `form:"Contents" json:"Contents"`
	RelationId string `from:"RelationId" json:"RelationId"`
}
