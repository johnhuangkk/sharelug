package InvoiceXml

import "encoding/xml"

type InvoiceA0401 struct {
	XMLName        xml.Name       `xml:"Invoice"`
	Text           string         `xml:",chardata"`
	Xmlns          string         `xml:"xmlns,attr"`
	Xsi            string         `xml:"xsi,attr"`
	SchemaLocation string         `xml:"schemaLocation,attr"`
	Main           InvoiceMain    `xml:"Main"`
	Details        InvoiceDetails `xml:"Details"`
	Amount         InvoiceAmount  `xml:"Amount"`
}
