package OKMart

import (
	"encoding/xml"
	"time"
)

//<F64DOC>
//    <DOCHEAD>
//        <DOCDATE>20201116</DOCDATE>
//        <FROMPARTNERCODE>CVS</FROMPARTNERCODE>
//        <TOPARTENERCODE>462</TOPARTENERCODE>
//    </DOCHEAD>
//    <F64CONTENT>
//		<ECNO>516</ECNO>
//		<STNO>K001234</STNO>
//		<ODNO>12345678901</ODNO>
//		<DCSTDT>20180913030035</DCSTDT>
//		<VENDOR>EC 廠商</VENDOR>
//		<VENDORNO>TW1234567890</VENDORNO>
//	  </F64CONTENT>
//    <F64CONTENT>
//		<ECNO>516</ECNO>
//		<STNO>K001234</STNO>
//		<ODNO>12345678901</ODNO>
//		<DCSTDT>20180913030035</DCSTDT>
//		<VENDOR>EC 廠商</VENDOR>
//		<VENDORNO>TW1234567890</VENDORNO>
//	  </F64CONTENT>
//</F64DOC>


type F64Doc struct {
	XMLName  xml.Name     `xml:"F64DOC"`
	Head     F64DocHead   `xml:"DOCHEAD"`
	Body     []F64Content   `xml:"F64CONTENT"`
}

type F64DocHead struct {
	DocDate         string `xml:"DOCDATE"`
	FromPartnerCode string `xml:"FROMPARTNERCODE"`
	ToPartenerCode  string `xml:"TOPARTENERCODE"`
}

type F64Content struct {
	EcNo            string `xml:"ECNO"`
	RrStoreId       string `xml:"STNO"`
	OrderNo         string `xml:"ODNO"`
	RrInDateTime    string `xml:"DCSTDT"`
	Vendor          string `xml:"VENDOR"`
	VendorNo        string `xml:"VENDORNO"`
}

func (s *F64Content) GetDateTime () string {
	t, _ := time.Parse(`20060102150405`, s.RrInDateTime)
	return t.Format(`2006-01-02 15:04:05`)
}

func (receiver *F64Doc) DecodeXML(data []byte) (err error) {
	err = xml.Unmarshal(data, receiver)
	return
}

