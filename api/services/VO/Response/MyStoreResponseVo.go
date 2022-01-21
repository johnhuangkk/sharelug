package Response

type MyStoreResponse struct {
	StoreName           string                `json:"StoreName"`
	MemberPhone         string                `json:"MemberPhone"`
	Rank                string                `json:"Rank"`
	RankText            string                `json:"RankText"`
	BalanceAmount       string                `json:"BalanceAmount"`
	OrderWarn           OrderWarn             `json:"OrderWarn"`
	AccountActivityList []AccountActivityList `json:"AccountActivityList"`
}

type OrderWarn struct {
	BillMessage     int64 `json:"BillMessage"`
	ShipMessage     int64 `json:"ShipMessage"`
	RefundMessage   int64 `json:"RefundMessage"`
	CustomerMessage int64 `json:"CustomerMessage"`
}

type AccountActivityList struct {
	Message string `json:"Message"`
	Time    string `json:"Time"`
}

type SalesReportResponse struct {
	StartTime        string `json:"StartTime"`
	EndTime          string `json:"EndTime"`
	SalesCount       int64  `json:"SalesCount"`
	SalesAmount      int64  `json:"SalesAmount"`
	CaptureAmount    int64  `json:"CaptureAmount"`
	RecCaptureAmount int64  `json:"RecCaptureAmount"`
}

type SettingStoreResponse struct {
	StoreStatus      string          `json:"StoreStatus"`
	StoreStatusText  string          `json:"StoreStatusText"`
	StoreName        string          `json:"StoreName"`
	StorePicture     string          `json:"StorePicture"`
	FreeShipKey      string          `json:"FreeShipKey"`
	FreeShip         int64           `json:"FreeShip"`
	SelfDeliveryKey  string          `json:"SelfDeliveryKey"`
	SelfDeliveryFree int64           `json:"SelfDeliveryFree"`
	SelfDelivery     bool            `json:"SelfDelivery"`
	ManagerCount     int64           `json:"ManagerCount"`
	ManagerMax       int64           `json:"ManagerMax"`
	MyStoreList      []StoreDataResp `json:"MyStoreList"`
}

type StoreDataResp struct {
	StoreId      string `json:"StoreId"`
	SellerId     string `json:"SellerId"`
	StoreName    string `json:"StoreName"`
	StoreTax     string `json:"StoreTax"`
	StorePicture string `json:"StorePicture"`
	StoreDefault int    `json:"StoreDefault"`
	StoreStatus  string `json:"StoreStatus"`
	ExpireTime   string `json:"ExpireTime"`
	RankId       int    `json:"RankId"`
	UserId       string `json:"UserId"`
	Rank         string `json:"Rank"`
	RankStatus   string `json:"RankStatus"`
	RankExpire   string `json:"RankExpire"`
	ManagerCount int64  `json:"ManagerCount"`
}

type SettingUserResponse struct {
	UserPhone   string        `json:"UserPhone"`
	Nickname    string        `json:"Nickname"`
	Email       string        `json:"Mail"`
	VerifyEmail string        `json:"VerifyEmail"`
	Picture     string        `json:"Picture"`
	BankAccount []BankAccount `json:"BankAccount"`
	CreditCard  []CreditCard  `json:"CreditCard"`
}

type BankAccount struct {
	AccountId     string `json:"AccountId"`
	AccountNumber string `json:"AccountNumber"`
	IsDefault     bool   `json:"IsDefault"`
}

type CreditCard struct {
	CardId     string `json:"CardId"`
	CardNumber string `json:"CardNumber"`
	IsDefault  bool   `json:"IsDefault"`
}

type BalanceResponse struct {
	StartTime           string               `json:"StartTime"`
	EndTime             string               `json:"EndTime"`
	BalanceAccountCount int64                `json:"BalanceAccountCount"`
	BalanceAccountList  []BalanceAccountList `json:"BalanceAccountList"`
}

type BalanceAccountList struct {
	Date      string `json:"Date"`
	TransText string `json:"TransText"`
	In        int64  `json:"In"`
	Out       int64  `json:"Out"`
	Balance   int64  `json:"Balance"`
	Comment   string `json:"Comment"`
}

type SellerBalanceResponse struct {
	Account         string `json:"Account"`
	SellerId        string `json:"SellerId"`
	Balance         int64  `json:"Balance"`
	RetainBalance   int64  `json:"RetainBalance"`
	DetainBalance   int64  `json:"DetainBalance"`
	WithholdBalance int64  `json:"WithholdBalance"`
}

type MyAccountResponse struct {
	Balance       int64           `json:"Balance"`
	RetentionList []RetentionList `json:"RetentionList"`
}

type RetentionList struct {
	Date    string `json:"Date"`
	Store   string `json:"Store"`
	OrderId string `json:"OrderId"`
	Amount  int64  `json:"Amount"`
}

type RetainAccountResponse struct {
	RetainAccountCount int64               `json:"RetainAccountCount"`
	RetainAccountList  []RetainAccountList `json:"RetainAccountList"`
}

type RetainAccountList struct {
	Date      string `json:"Date"`
	StoreName string `json:"StoreName"`
	OrderId   string `json:"OrderId"`
	Amount    int64  `json:"Amount"`
}
