package InvoiceXml

//營業人電子發票配號檔
type InvoiceAssignNo struct {
	Ban            string `xml:"Ban"`            //公司統一編號
	InvoiceType    string `xml:"InvoiceType"`    //發票類別
	YearMonth      string `xml:"YearMonth"`      //發票期別
	InvoiceTrack   string `xml:"InvoiceTrack"`   //字軌
	InvoiceBeginNo string `xml:"InvoiceBeginNo"` //發票起號
	InvoiceEndNo   string `xml:"InvoiceEndNo"`   //發票迄號
	InvoiceBooklet string `xml:"InvoiceBooklet"` //本組數
}
