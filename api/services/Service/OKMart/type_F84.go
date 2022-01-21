package OKMart

import (
	"encoding/xml"
)

//<F84DOC>
//	<DOCHEAD>
//		<DOCDATE>20180913</DOCDATEDOCDATE>
//		<FROM PARTNERCODE>CVS</FROM PARTNERCODE>
//		<TOPARTNERCODE >526</TOPARTNERCODE>
//	</DOCHEAD>
//	<F84CONTENT>
//		<ECNO>526</ECNO>
//		<STNO>K001234</STNO>
//		<ODNO>12345678901</ODNO>
//		<DCSTDT>20180913142848</DCSTDT>
//		<EASYECNO></EASYECNO>
//	</F84CONTENT>
//	<F84CONTENT>
//		<ECNO>526</ECNO>
//		<STNO>K001234</STNO>
//		<ODNO>12345678901</ODNO>
//		<DCSTDT>20180913142848</DCSTDT>
//		<EASYECNO></EASYECNO>
//	</F84CONTENT>
//</F84DOC>

type F84Doc struct {
	XMLName  xml.Name     `xml:"F84DOC"`
	Head     F84DocHead   `xml:"DOCHEAD"`
	Body     []F84Content `xml:"F84CONTENT"`
}

type F84DocHead struct {
	DocDate         string `xml:"DOCDATE"`
	FromPartnerCode string `xml:"FROMPARTNERCODE"`
	ToPartenerCode  string `xml:"TOPARTNERCODE"`
}

type F84Content struct {
	EcNo          string `xml:"ECNO"`
	RrStoreId     string `xml:"STNO"`
	OrderNo       string `xml:"ODNO"`
	UpDateTime    string `xml:"DCSTDT"` // 實際離店時間
	OtherCode     string `xml:"EASYECNO"`
}

func (receiver *F84Doc) DecodeXML(data []byte) (err error) {
	err = xml.Unmarshal(data, receiver)
	return
}