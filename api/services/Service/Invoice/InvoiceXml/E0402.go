package InvoiceXml

//空白未使用字軌檔
type BranchTrackBlank struct {
	Main    BranchTrackBlankMain    `xml:"Main"`
	Details BranchTrackBlankDetails `xml:"Details"`
}

type BranchTrackBlankMain struct {
	HeadBan      string `xml:"HeadBan"`      //總公司統一編號
	BranchBan    string `xml:"BranchBan"`    //分支機構統一編號
	InvoiceType  string `xml:"InvoiceType"`  //發票類別
	YearMonth    string `xml:"YearMonth"`    //發票期別
	InvoiceTrack string `xml:"InvoiceTrack"` //空白字軌
}

type BranchTrackBlankDetails struct {
	BranchTrackBlankItem []BranchTrackBlankItem `xml:"BranchTrackBlankItem"`
}

type BranchTrackBlankItem struct {
	InvoiceBeginNo string `xml:"InvoiceBeginNo"` //空白發票起號
	InvoiceEndNo   string `xml:"InvoiceEndNo"`   //空白發票迄號
}
