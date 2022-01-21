package InvoiceXml

type AllowanceB0401 struct {
	Main    AllowanceMain    `xml:"Main"`
	Details AllowanceDetails `xml:"Details"`
	Amount  AllowanceAmount  `xml:"Amount"`
}

