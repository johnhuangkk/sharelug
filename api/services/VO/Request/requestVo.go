package Request

type ProductList struct {
	Status    string `form:"status"`
	Name      string `form:"name"`
	Ship      string `form:"ship"`
	PayWay    string `form:"payWay"`
	Length    int    `form:"Limit"`
	Page      int    `form:"Start"`
	Price     string `form:"price"`
	ShipMerge string `form:"shipMerge"`
}

type BatchProductShipping struct {
	ProductId    []string      `form:"ProductId" json:"ProductId" validate:"required"`
	ShippingList []NewShipping `form:"ShippingList" json:"ShippingList" validate:"required"` //運送方式加運費
	ShipMerge    int           `form:"ShipMerge" json:"ShipMerge"`                           //是否合併運費
}

type BatchProductPayWay struct {
	ProductId  []string `form:"ProductId" json:"ProductId" validate:"required"`
	PayWayList []string `form:"PayWayList" json:"PayWayList" validate:"required"`
}

type BatchProductStatus struct {
	ProductId   []string `form:"ProductId" json:"ProductId" validate:"required"`
	ProductDown bool     `form:"ProductDown" json:"ProductDown"`
}

type BatchProductFreeShip struct {
	ProductId  []string `form:"ProductId" json:"ProductId" validate:"required"`
	IsFreeShip bool     `form:"IsFreeShip" json:"IsFreeShip"`
}

type OrderSearch struct {
	Page     int    `form:"page" json:"page"`
	Length   int    `form:"per" json:"per"`
	Status   string `form:"status" json:"status"`
	Tab      string `form:"tab" json:"tab"`
	OrderId  string `form:"orderId" json:"orderId"`
	ShipType string `form:"shipType" json:"shipType"`
	Duration int    `form:"duration" json:"duration"`
	OrderBy  string `form:"orderBy" json:"orderBy"`
}

type NotifySearch struct {
	Tab string `form:"tab"`
}

type RealTimesRequest struct {
	Tab   string `form:"Tab" json:"Tab"`
	Limit int    `form:"Limit" json:"Limit"`
	Start int    `form:"Start" json:"Start"`
}

type NotifyRequest struct {
	Email string `form:"email" json:"email"`
	Tel   string `form:"tel" json:"tel"`
}

type MemberEmailVerifyRequest struct {
	VerifyCode string `form:"VerifyCode" json:"VerifyCode"`
}
type MemberCompanyVerifyRequest struct {
	Mphone           string `form:"Mphone" json:"Mphone"`
	CompanyName      string `form:"CompanyName" json:"CompanyName"`
	CompanyAddr      string `form:"CompanyAddr" json:"CompanyAddr"`
	Representative   string `form:"Representative" json:"Representative"`
	RepresentativeId string `form:"RepresentativeId" json:"RepresentativeId"`
	Identity         string `form:"Identity" json:"Identity"`
}
type ContactRequest struct {
	Email     string `form:"Email" json:"Email"`
	UserName  string `form:"UserName" json:"UserName"`
	Company   string `form:"Company" json:"Company"`
	Telephone string `form:"Telephone" json:"Telephone"`
	Contents  string `form:"Contents" json:"Contents"`
}

type NotificationRequest struct {
	Tab   string `form:"Tab" json:"Tab"`
	Limit int64  `form:"Limit" json:"Limit"`
	Start int64  `form:"Start" json:"Start"`
}

type NotificationReadRequest struct {
	MessageId int64 `form:"MessageId" json:"MessageId"`
}

type GetB2CPayRequest struct {
	ProductId string `form:"ProductId" json:"ProductId" validate:"required"`
}

type B2CPayRequest struct {
	OrderId        string         `form:"OrderId" json:"OrderId"`
	Amount         int64          `form:"Amount" json:"Amount"`
	Payment        string         `form:"Payment" json:"Payment"`
	CardId         string         `form:"CardId" json:"CardId"`
	CardNumber     string         `form:"CardNumber" json:"CardNumber"`
	CardExpiration string         `form:"CardExpiration" json:"CardExpiration"`
	CardSecurity   string         `form:"CardSecurity" json:"CardSecurity"`
	Carrier        CarrierRequest `form:"Carrier" json:"Carrier"`
}

func (p *B2CPayRequest) GetB2cCreditPayment() PayParams {
	var resp PayParams
	resp.CardId = p.CardId
	resp.CardNumber = p.CardNumber
	resp.CardExpiration = p.CardExpiration
	resp.CardSecurity = p.CardSecurity
	return resp
}

type B2COrderRequest struct {
	OrderId string `form:"OrderId" json:"OrderId"`
}

type CreationStoreRequest struct {
	StorePicture string `form:"StorePicture" json:"StorePicture"`
	StoreName    string `form:"StoreName" json:"StoreName"`
}

type CreationManagerRequest struct {
	ManagerEmail string `form:"ManagerEmail" json:"ManagerEmail"`
	ManagerPhone string `form:"ManagerPhone" json:"ManagerPhone"`
}

type InviteManagerRequest struct {
	ManagerId int `json:"ManagerId"`
}

type DeleteManagerRequest struct {
	ManagerId int `json:"ManagerId"`
}

type LogsRequest struct {
	Logs string `form:"Logs"`
}

type CarrierRequest struct {
	InvoiceType string `form:"InvoiceType" validate:"required"`
	CompanyBan  string `form:"CompanyBan"`
	CompanyName string `form:"CompanyName"`
	DonateBan   string `form:"DonateBan"`
	CarrierType string `form:"CarrierType" validate:"required"`
	CarrierId   string `form:"CarrierId"`
}

type BindPlatformRequest struct {
	Token string `form:"token"`
	Ban   string `form:"ban"`
}

type BindCarrierRequest struct {
	Token string `form:"Token" json:"Token"`
}

type VerifyDonateCodeRequest struct {
	Version string `json:"version"`
	Action  string `json:"action"`
	PCode   string `json:"pCode"`
	TxID    string `json:"TxId"`
	AppId   string `json:"appId"`
}

type InvoiceListRequest struct {
	Limit int64 `form:"Limit" json:"Limit"`
	Start int64 `form:"Start" json:"Start"`
}

type StoreInfoRequest struct {
	Address  string `json:"Address"`
	Industry string `json:"Industry"`
}

type MemberCompanyRequest struct {
	CompanyName      string  `form:"CompanyName" json:"CompanyName"`
	CompanyAddr      string  `form:"CompanyAddr" json:"CompanyAddr"`
	CompanyAddrEn    string  `form:"CompanyAddrEn" json:"CompanyAddrEn"`
	Representative   string  `form:"Representative" json:"Representative"`
	RepresentativeId string  `form:"RepresentativeId" json:"RepresentativeId"`
	RepresentFirst   string  `form:"RepresentFirst" json:"RepresentFirst"`
	RepresentLast    string  `form:"RepresentLast" json:"RepresentLast"`
	Capital          float64 `form:"Capital" json:"Capital"`
	Establish        string  `form:"Establist" json:"Establish"`
	ZipCode          string  `form:"ZipCode" json:"ZipCode"`
	Identity         string  `form:"Identity" json:"Identity"`
	Contact          string  `form:"Contact" json:"Contact"`
	ContactPhone     string  `form:"ContactPhone" json:"ContactPhone"`
}

type MemberPersonalRequest struct {
	RepresentFirst string `form:"RepresentFirst" json:"RepresentFirst"`
	RepresentLast  string `form:"RepresentLast" json:"RepresentLast"`
	AddrEn         string `form:"CompanyAddrEn" json:"AddrEn"`
}

type MemberPersonalSpecialStoreRequest struct {
	MemberPhone string `form:"MemberPhone" json:"MemberPhone"`
	MerchantId  string `form:"MerchantId" json:"MerchantId"`
}

type KgiSpecialStoreExcel struct {
	Filename string `json:"Filename"`
}

type TranslateAddrEn struct {
	Addr string `json:"Addr"`
}
