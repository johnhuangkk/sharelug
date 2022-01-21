package InvoiceXml

//分支機構配號檔
type BranchTrack struct {
	Main    BranchTrackMain    `xml:"Main"`
	Details BranchTrackDetails `xml:"Details"`
}

type BranchTrackMain struct {
	HeadBan        string `xml:"HeadBan"`        //總公司統一編號
	BranchBan      string `xml:"BranchBan"`      //分支機構統一編號
	InvoiceType    string `xml:"InvoiceType"`    //發票類別
	YearMonth      string `xml:"YearMonth"`      //發票期別
	InvoiceTrack   string `xml:"InvoiceTrack"`   //發票字軌
	InvoiceBeginNo string `xml:"InvoiceBeginNo"` //發票起號
	InvoiceEndNo   string `xml:"InvoiceEndNo"`   //發票迄號
}

type BranchTrackDetails struct {
	BranchTrackItem BranchTrackItem `xml:"BranchTrackItem"` //分支機構配號項目資料
}

type BranchTrackItem struct {
	InvoiceBeginNo string `xml:"InvoiceBeginNo"` //發票起號
	InvoiceEndNo   string `xml:"InvoiceEndNo"`   //發票迄號
	InvoiceBooklet string `xml:"InvoiceBooklet"` //本組數
}


