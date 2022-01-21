package Enum

const (
	MemberStatusSuccess = "SUCC"
	MemberBuyer         = "Buyer"
	MemberSeller        = "Seller"

	CategoryMember  = "MEMBER"
	CategoryCompany = "COMPANY"

	VerifyIdentitySuccess = 1

	InvoiceTypePersonal = "PERSONAL" //個人發票
	InvoiceTypeCompany  = "COMPANY"  //公司發票
	InvoiceTypeDonate   = "DONATE"   //捐贈發票

	InvoiceCarrierTypeMember = "MEMBER" //會員載具
	InvoiceCarrierTypeMobile = "MOBILE" //手機載具
	InvoiceCarrierTypeCert   = "CERT"   //憑證載具
	InvoiceCarrierTypePrint  = "PRINT"  //已索取紙本電子發票
)

var MemberStatus = map[string]string{
	MemberStatusSuccess: "正常",
}

var InvoiceCarrierType = map[string]string{
	InvoiceCarrierTypeMember: "會員載具",
	InvoiceCarrierTypeMobile: "手機載具",
	InvoiceCarrierTypeCert:   "憑證載具",
	InvoiceCarrierTypePrint:  "已索取紙本電子發票",
}

var CarrierType = map[string]string{
	InvoiceCarrierTypeMember: "EJ1515",
	InvoiceCarrierTypeMobile: "3J0002",
	InvoiceCarrierTypeCert:   "CQ0001",
}
