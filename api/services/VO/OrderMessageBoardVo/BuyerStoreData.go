package OrderMessageBoardVo

type Picture struct {
	Id, Name, Picture string
}

type BuyerStorePictureData struct {
	OrderId   string    `json:"OrderId"`
	IsBuyer   bool      `json:"IsBuyer"`
	Buyer     Picture   `json:"Buyer"`
	Store     Picture   `json:"Store"`
	OrderInfo OrderInfo `json:"OrderInfo"`
}

type OrderInfo struct {
	Products   []string `json:"Products"`
	BuyerName  string   `json:"BuyerName"`
	BuyerPhone string   `json:"BuyerPhone"`
}
