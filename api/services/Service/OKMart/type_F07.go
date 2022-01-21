package OKMart

import "encoding/xml"

//<F07DOC>
//	<DOCHEAD>
//		<DOCDATE>20180915</DOCDATE>
//		<FROMPARTNERCODE>CVS</FROMPARTNERCODE>
//		<TOPARTNERCODE>416</TOPARTNERCODE>
//  </DOCHEAD>
// 	<F07CONTENT>
//		<RET_M>T08</RET_M>
//		<ECNO>416</ECNO>
//		<STNO>K001234</STNO>
//		<ODNO>12345678901</ODNO>
//		<RTDCDT>20180915200006</RTDCDT>
//		<FRTDCDT>20180915200006</FRTDCDT>
//	</F07CONTENT>
// 	<F07CONTENT>
//		<RET_M>T08</RET_M>
//		<ECNO>416</ECNO>
//		<STNO>K001234</STNO>
//		<ODNO>12345678901</ODNO>
//		<RTDCDT>20180915200006</RTDCDT>
//		<FRTDCDT>20180915200006</FRTDCDT>
//	</F07CONTENT>
//</F07DOC>


type F07Doc struct {
	XMLName  xml.Name     `xml:"F07DOC"`
	Head     F07DocHead   `xml:"DOCHEAD"`
	Body     []F07Content `xml:"F07CONTENT"`
}

type F07DocHead struct {
	DocDate         string `xml:"DOCDATE"`
	FromPartnerCode string `xml:"FROMPARTNERCODE"`
	ToPartenerCode  string `xml:"TOPARTNERCODE"`
}

type F07Content struct {
	ReturnCode      string `xml:"RET_M"`
	EcNo            string `xml:"ECNO"`
	OrderNo         string `xml:"ODNO"`
	RrStoreId       string `xml:"STNO"`
	UpDateTime      string `xml:"RTDCDT"`
	CheckDateTime   string `xml:"FRTDCDT"`
}

func (receiver *F07Doc) DecodeXML(data []byte) (err error) {
	err = xml.Unmarshal(data, receiver)
	return
}

