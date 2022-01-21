package OKMart

import (
	"encoding/xml"
)

//<F67DOC>
//    <DOCHEAD>
//        <DOCDATE>20201116</DOCDATE>
//        <FROMPARTNERCODE>CVS</FROMPARTNERCODE>
//        <TOPARTENERCODE>462</TOPARTENERCODE>
//    </DOCHEAD>
//    <F67CONTENT>
//		  <RETM>T08</RETM>
//		  <ECNO>516</ECNO>
//		  <STNO>K001234</STNO>
//		  <ODNO>12345678901</ODNO>
//		  <RTDCDT>20180913144333</RTDCDT>
//		  <FRTDCDT>20180902144333</FRTDCDT>
//		  <VENDOR>EC 廠商</VENDOR>
//		  <VENDORNO>TW1234567890</VENDORNO>
//	  </F67CONTENT>
//    <F67CONTENT>
//		  <RETM>T08</RETM><ECNO>516</ECNO><STNO>K001234</STNO>
//		  <ODNO>12345678901</ODNO><RTDCDT>20180913144333</RTDCDT>
//		  <FRTDCDT>20180902144333</FRTDCDT>
//		  <VENDOR>EC 廠商</VENDOR><VENDORNO>TW1234567890</VENDORNO>
//	  </F67CONTENT>
//</F67DOC>


type F67Doc struct {
	XMLName  xml.Name     `xml:"F67DOC"`
	Head     F67DocHead   `xml:"DOCHEAD"`
	Body     []F67Content `xml:"F67CONTENT"`
}

type F67DocHead struct {
	DocDate         string `xml:"DOCDATE"`
	FromPartnerCode string `xml:"FROMPARTNERCODE"`
	ToPartenerCode  string `xml:"TOPARTENERCODE"`
}

type F67Content struct {
	ReturnCode    string `xml:"RET_M"`
	EcNo          string `xml:"ECNO"`
	RrStoreId     string `xml:"STNO"`
	OrderNo       string `xml:"ODNO"`
	UpDateTime    string `xml:"RTDCDT"`
	CheckDateTime string `xml:"FRTDCDT"`
	Vendor        string `xml:"VENDOR"`
	VendorNo      string `xml:"VENDORNO"`
}

func (receiver *F67Doc) DecodeXML(data []byte) (err error) {
	err = xml.Unmarshal(data, receiver)
	return
}

