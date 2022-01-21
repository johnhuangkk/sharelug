package OKMart

import "encoding/xml"

//<F65DOC>
//    <DOCHEAD>
//        <DOCDATE>20201116</DOCDATE>
//        <FROMPARTNERCODE>CVS</FROMPARTNERCODE>
//        <TOPARTENERCODE>462</TOPARTENERCODE>
//    </DOCHEAD>
//    <F65CONTENT>
//		<BC1>516725992</BC1>
//		<BC2>3916921910002084</BC2>
//		<STNO>>K001234</STNO>
//		<RTDT>20180914140036</RTDT>
//		<TKDT>20180914140036</TKDT>
//		<VENDOR>EC 廠商</VENDOR>
//		<VENDORNO>TW1234567890</VENDORNO>
//	  </F65CONTENT>
//    <F65CONTENT>
//		<BC1>516725992</BC1>
//		<BC2>3916921910002084</BC2>
//		<STNO>>K001234</STNO>
//		<RTDT>20180914140036</RTDT>
//		<TKDT>20180914140036</TKDT>
//		<VENDOR>EC 廠商</VENDOR>
//		<VENDORNO>TW1234567890</VENDORNO>
//	  </F65CONTENT>
//</F65DOC>


type F65Doc struct {
	XMLName  xml.Name     `xml:"F65DOC"`
	Head     F65DocHead   `xml:"DOCHEAD"`
	Body     []F65Content `xml:"F65CONTENT"`
}

type F65DocHead struct {
	DocDate         string `xml:"DOCDATE"`
	FromPartnerCode string `xml:"FROMPARTNERCODE"`
	ToPartenerCode  string `xml:"TOPARTENERCODE"`
}

type F65Content struct {
	BarCode1        string `xml:"BC1"`
	BarCode2        string `xml:"BC2"`
	RrStoreId       string `xml:"STNO"`
	RrPickDateTime  string `xml:"RTDT"`
	CheckDateTime   string `xml:"TKDT"`
	Vendor          string `xml:"VENDOR"`
	VendorNo        string `xml:"VENDORNO"`
}

func (receiver *F65Doc) DecodeXML(data []byte) (err error) {
	err = xml.Unmarshal(data, receiver)
	return
}

