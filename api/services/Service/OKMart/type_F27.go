package OKMart

import (
	"encoding/xml"
)

//<F27DOC>
//	<DOCHEAD>
//		<DOCDATE>20180912</DOCDATEDOCDATE> <FROMPARTNERCODE>CVS</FROM PARTNERCODE>
//		<TOPARTNERCODE>526</TOPARTNERCODE>
//	</DOCHEAD>
//	<F27CONTENT>
//		<ECNO>526</ECNO>
//		<ODNO>12345678901</ODNO>
//		<STNO>K001234</STNO>
//		<RTDT>20180912203045</RTDT>
//	</F27CONTENT>
//</F27DOC>

type F27Doc struct {
	XMLName  xml.Name     `xml:"F27DOC"`
	Head     F27DocHead   `xml:"DOCHEAD"`
	Body     []F27Content `xml:"F27CONTENT"`
}

type F27DocHead struct {
	DocDate         string `xml:"DOCDATE"`
	FromPartnerCode string `xml:"FROMPARTNERCODE"`
	ToPartenerCode  string `xml:"TOPARTNERCODE"`
}

type F27Content struct {
	EcNo          string `xml:"ECNO"`
	OrderNo       string `xml:"ODNO"`
	RrStoreId     string `xml:"STNO"`
	UpDateTime    string `xml:"RTDT"` // 實際寄件代收日期
}

func (receiver *F27Doc) DecodeXML(data []byte) (err error) {
	return xml.Unmarshal(data, receiver)
}