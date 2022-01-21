package Request

type WithdrawRequest struct {
	Email       string `form:"Email" json:"Email"`
	AccountName string `form:"AccountName" json:"AccountName"`
	BankCode    string `form:"BankCode" json:"BankCode"`
	BankAccount string `form:"BankAccount" json:"BankAccount"`
	WithdrawAmt string `form:"WithdrawAmt" json:"WithdrawAmt"`
}

type EditWithdrawRequest struct {
	AccountId string `form:"AccountId" json:"AccountId"`
}

type EditCreditRequest struct {
	CreditId string `form:"CreditId" json:"CreditId"`
}