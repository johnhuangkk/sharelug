package OKMart

import (
	"encoding/xml"
)

//<F44DOC>
//    <DOCHEAD>
//        <DOCDATE>20201116</DOCDATE>
//        <FROMPARTNERCODE>CVS</FROMPARTNERCODE>
//        <TOPARTENERCODE>462</TOPARTENERCODE>
//    </DOCHEAD>
//    </F44CONTENT>
//		<ECNO>516</ECNO>
//		<ODNO>12345678901</ODNO>
//		<STNO>K001234</STNO>
//		<DCSTDT>20180913030035</DCSTDT>
//		<VENDOR>EC 廠商</VENDOR>
//		<VENDORNO>TW1234567890</VENDORNO>
//	  </F44CONTENT>
//</F44DOC>


type F44Doc struct {
	XMLName  xml.Name     `xml:"F44DOC"`
	Head     F44DocHead   `xml:"DOCHEAD"`
	Body     []F44Content `xml:"F44CONTENT"`
}

type F44DocHead struct {
	DocDate         string `xml:"DOCDATE"`
	FromPartnerCode string `xml:"FROMPARTNERCODE"`
	ToPartenerCode  string `xml:"TOPARTENERCODE"`
}

type F44Content struct {
	EcNo          string `xml:"ECNO"`
	OrderNo       string `xml:"ODNO"`
	RrStoreId     string `xml:"STNO"`
	RrInDateTime  string `xml:"DCSTDT"`
	Vendor        string `xml:"VENDOR"`
	VendorNo      string `xml:"VENDORNO"`
}

func (receiver *F44Doc) DecodeXML(data []byte) (err error) {
	err = xml.Unmarshal(data, receiver)
	return
}
