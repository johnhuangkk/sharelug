package UserAddress

// backend 接收 frontend 呼叫 地址資訊格式
type AddressInfo struct {
	Address string `form:"Address" validate:"required"`
	Ship    string `form:"Ship" validate:"required"`
	Type    string `form:"Type" validate:"required"`
	Name    string `form:"Name"`
	Phone    string `form:"Phone"`
}

type DeleteAddress struct {
	UaId string `form:"UaId" validate:"required"`
}

type AddressInfoResponse struct {
	Id      string
	Alias   string
	Ship    string
	Country string
	City    string
	SpecId  string
	Address string
	Status  string
	Name    string
	Phone    string
}
