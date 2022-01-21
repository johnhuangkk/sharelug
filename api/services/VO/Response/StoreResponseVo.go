package Response

type StoreSocialMediaResponse struct {
	Facebook  Facebook  `json:"Facebook"`
	Instagram Instagram `json:"Instagram"`
	Line      Line      `json:"Line"`
	Telegram  Telegram  `json:"Telegram"`
}

type Facebook struct {
	Link string `form:"Link" json:"Link"`
	Show int    `form:"Show" json:"Show"`
}

type Instagram struct {
	Link string `form:"Link" json:"Link"`
	Show int    `form:"Show" json:"Show"`
}

type Line struct {
	Link string `form:"Link" json:"Link"`
	Show int    `form:"Show" json:"Show"`
}

type Telegram struct {
	Link string `form:"Link" json:"Link"`
	Show int    `form:"Show" json:"Show"`
}

type UserOperateRecordResponse struct {
	RecordTime    string `json:"RecordTime"`
	RecordContent string `json:"RecordContent"`
}

type StoreFreeShipResponse struct {
	FreeShipKey      string `json:"FreeShipKey"`
	FreeShip         int64  `json:"FreeShip"`
	FreeShipText     string `json:"FreeShipText"`
	SelfDelivery     bool   `json:"SelfDelivery"`
	SelfDeliveryKey  string `json:"SelfDeliveryKey"`
	SelfDeliveryFree int64  `json:"SelfDeliveryFree"`
	IsCoupon         bool   `json:"IsCoupon"`
}
