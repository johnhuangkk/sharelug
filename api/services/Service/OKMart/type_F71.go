package OKMart

import (
	"encoding/xml"
)

//<F71DOC>
//		<DOCHEAD>
//			<DOCDATE >20180913</DOCDATEDOCDATE> <FROMPARTNERCODE >CVS</FROMPARTNERCODE> <TOPARTNERCODE>526</TOPARTNERCODE>
//      </DOCHEAD>
//     	<F71CONTENT>
//			<ECNO>526</ECNO>
//			<ODNO>12345678901</ODNOODNO>
//			<STNO>K001234</STNO>
//			<RTDT>20180913195004</RTDT>
//		</F71CONTENT>
//		<F71CONTENT>
//			<ECNO>526</ECNO>
//			<ODNO>12345678901</ODNOODNO>
//			<STNO>K001234</STNO>
//			<RTDT>20180913195004</RTDT>
//		</F71CONTENT>
//</F71DOC>

type F71Doc struct {
	XMLName  xml.Name     `xml:"F71DOC"`
	Head     F71DocHead   `xml:"DOCHEAD"`
	Body     []F71Content `xml:"F71CONTENT"`
}

type F71DocHead struct {
	DocDate         string `xml:"DOCDATE"`
	FromPartnerCode string `xml:"FROMPARTNERCODE"`
	ToPartenerCode  string `xml:"TOPARTNERCODE"`
}

type F71Content struct {
	EcNo          string `xml:"ECNO"`
	OrderNo       string `xml:"ODNO"`
	RrStoreId     string `xml:"STNO"`
	UpDateTime    string `xml:"RTDT"`
}

func (receiver *F71Doc) DecodeXML(data []byte) (err error) {
	err = xml.Unmarshal(data, receiver)
	return
}
