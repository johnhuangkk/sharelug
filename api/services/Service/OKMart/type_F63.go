package OKMart

import (
	"encoding/xml"
)

//<F63DOC>
//	<DOCHEAD>
//		<DOCDATE>20180913</DOCDATE>
//		<FROMPARTNERCODE>CVS</FROMPARTNERCODE>
//		<TOPARTNERCODE>516</TOPARTNERCODE>
//	</DOCHEAD>
//	<F63CONTENT>
//		<ECNO>516</ECNO> <ODNO>12345678901</ODNO>
//		<CNNO>TOK</CNNO> <CUNAME>簡大翔</CUNAME> <PRODTYPE>0</PRODTYPE>
//		<PINCODE>12345678901</PINCODE>
//		<RTDT>20180913163005</RTDT>
//	</F63CONTENT>
//  <F63CONTENT>
//		<ECNO>516</ECNO>
//		<ODNO>12345678901</ODNO>
//		<CNNO>TOK</CNNO> <CUNAME>簡大翔</CUNAME> <PRODTYPE>0</PRODTYPE>
//		<PINCODE>12345678901</PINCODE>
//		<RTDT>20180913163005</RTDT>
//	</F63CONTENT>
//</F63DOC>

type F63Doc struct {
	XMLName  xml.Name     `xml:"F63DOC"`
	Head     F63DocHead   `xml:"DOCHEAD"`
	Body     []F63Content `xml:"F63CONTENT"`
}

type F63DocHead struct {
	DocDate         string `xml:"DOCDATE"`
	FromPartnerCode string `xml:"FROMPARTNERCODE"`
	ToPartenerCode  string `xml:"TOPARTNERCODE"`
}

type F63Content struct {
	EcNo          string `xml:"ECNO"`
	OrderNo       string `xml:"ODNO"`
	CnNo          string `xml:"CNNO"`
	RrName        string `xml:"CUNAME"`
	ProductType   string `xml:"PRODTYPE"`
	SendNo        string `xml:"PINCODE"`
	UpDateTime    string `xml:"RTDT"`
}

func (receiver *F63Doc) DecodeXML(data []byte) (err error) {
	err = xml.Unmarshal(data, receiver)
	return
}
