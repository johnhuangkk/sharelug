package Response

type WithdrawResponse struct {
	Balance         int64             `json:"Balance"`
	IsVerifyEmail   bool              `json:"IsVerifyEmail"`
	WithdrawAccount []WithdrawAccount `json:"WithdrawAccount"`
	BackCodeList    []BackCodeList    `json:"BackCodeList"`
}

type WithdrawAccount struct {
	AccountName string `json:"AccountName"`
	BankCode    string `json:"BankCode"`
	BankAccount string `json:"BankAccount"`
	Default     bool   `json:"Default"`
}

type BackCodeList struct {
	BackCode string `json:"BackCode"`
	BankName string `json:"BankName"`
}

type StoreInfoResponse struct {
	Username string           `json:"Username"`
	Industry []IndustryCategory `json:"Industry"`
}

type IndustryCategory struct {
	Category string
	Industry []IndustryVo
}

type IndustryVo struct {
	Industry string
	Mcc      string
}
