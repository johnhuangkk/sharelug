package Response

type PaySendOtpResponse struct {
	Otp       int    `json:"Otp"`
	Email     string `json:"Email"`
	OtpExpire string `json:"OtpExpire"`
}

type ValidateOtpResponse struct {
	Uuid   string     `json:"UUID"`
	Token  string     `json:"Token"`
	Member MemberInfo `json:"Member"`
	Store  StoreInfo  `json:"Store"`
}

type MemberInfo struct {
	Uid            string `json:"Uid"`
	Email          string `json:"Mail"`
	Mphone         string `json:"Mphone"`
	Category       string `json:"Category"`
	Username       string `json:"Username"`
	RealName       string `json:"RealName"`
	IdentityName   string `json:"IdentityName"`
	Picture        string `json:"Picture"`
	VerifyIdentity int64  `json:"VerifyIdentity"`
	VerifyBusiness int64  `json:"VerifyBusiness"`
	UpgradeLevel   int64  `json:"UpgradeLevel"`
	StoreMax       int64  `json:"StoreMax"`
	StoreCount     int64  `json:"StoreCount"`
}

type StoreInfo struct {
	Sid      string `json:"Sid"`
	Name     string `json:"Name"`
	Picture  string `json:"Picture"`
	Rank     string `json:"Rank"`
	IsExpire bool   `json:"IsExpire"`
	Count    int64  `json:"Count"`
}
