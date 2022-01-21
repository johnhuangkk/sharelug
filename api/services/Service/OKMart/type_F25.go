package OKMart

import (
	"encoding/xml"
)

//<F25DOC>
//	<DOCHEAD>
//		<DOCDATE>20180912</DOCDATEDOCDATE> <FROMPARTNERCODE>CVS</FROM PARTNERCODE>
//		<TOPARTNERCODE>526</TOPARTNERCODE>
//	</DOCHEAD>
//	<F25CONTENT>
//		<ECNO>526</ECNO>
//		<ODNO>12345678901</ODNO>
//		<STNO>K001234</STNO>
//		<RTDT>20180912203045</RTDT>
//	</F25CONTENT>
//</F25DOC>

type F25Doc struct {
	XMLName  xml.Name     `xml:"F25DOC"`
	Head     F25DocHead   `xml:"DOCHEAD"`
	Body     []F25Content `xml:"F25CONTENT"`
}

type F25DocHead struct {
	DocDate         string `xml:"DOCDATE"`
	FromPartnerCode string `xml:"FROMPARTNERCODE"`
	ToPartenerCode  string `xml:"TOPARTNERCODE"`
}

type F25Content struct {
	BarCode1        string `xml:"BC1"`
	BarCode2        string `xml:"BC2"`
	RrStoreId       string `xml:"STNO"`
	UpDateTime      string `xml:"RTDT"`
	CheckDateTime   string `xml:"TKDT"`
	VendorNo        string `xml:"VENDORNO"`
}

func (receiver *F25Doc) DecodeXML(data []byte) (err error) {
	err = xml.Unmarshal(data, receiver)
	return
}