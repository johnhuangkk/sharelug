package OKMart

import "encoding/xml"

//<F04DOC>
//		<DOCHEAD>
//			<DOCDATE>20180915</DOCDATE>
//			<FROMPARTNERCODE>CVS</FROMPARTNERCODE>
//			<TOPARTNERCODE>417</TOPARTNERCODE>
//		</DOCHEAD>
//		<F04CONTENT>
//			<ECNO>417</ECNO>
//			<STNO>K001234</STNO>
//			<ODNO>123456778901</ODNO>
//			<DCSTDT>20180915025928</DCSTDT>
//		</F04CONTENT>
//      <F04CONTENT>
// 			<ECNO>417</ECNO>
// 			<STNO>K001234</STNO>
// 			<ODNO>123456778901</ODNO>
// 			<DCSTDT>20180915025928</DCSTDT>
// 		</F04CONTENT>
//</F04DOC>


type F04Doc struct {
	XMLName  xml.Name     `xml:"F04DOC"`
	Head     F04DocHead   `xml:"DOCHEAD"`
	Body     []F04Content `xml:"F04CONTENT"`
}

type F04DocHead struct {
	DocDate         string `xml:"DOCDATE"`
	FromPartnerCode string `xml:"FROMPARTNERCODE"`
	ToPartenerCode  string `xml:"TOPARTNERCODE"`
}

type F04Content struct {
	EcNo            string `xml:"ECNO"`
	OrderNo         string `xml:"ODNO"`
	RrStoreId       string `xml:"STNO"`
	UpDateTime      string `xml:"DCSTDT"`
}

func (receiver *F04Doc) DecodeXML(data []byte) (err error) {
	err = xml.Unmarshal(data, receiver)
	return
}

