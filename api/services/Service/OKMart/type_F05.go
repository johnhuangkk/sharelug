package OKMart

import "encoding/xml"

//<F05DOC>
//		<DOCHEAD>
//			<DOCDATE>20180915</DOCDATE>
//			<FROMPARTNERCODE>CVS</FROMPARTNERCODE>
//			<TOPARTNERCODE>416</TOPARTNERCODE>
//		</DOCHEAD>
//		<F05CONTENT>
//			<BC1>416375993</BC1>
//			<BC2>3849859930000015</BC2>
//			<STNO>K001234</STNO>
//			<RTDT>20180915230037</RTDT>
//			<TKDT>20180915230037</TKDT>
//			<VENDORNO>TW1234567890</VENDORNO>
//		</F05CONTENT>
//		<F05CONTENT>
//			<BC1>416375993</BC1>
//			<BC2>3849859930000015</BC2>
//			<STNO>K001234</STNO>
//			<RTDT>20180915230037</RTDT>
//			<TKDT>20180915230037</TKDT>
//			<VENDORNO>TW1234567890</VENDORNO>
//		</F05CONTENT>
//</F05DOC>


type F05Doc struct {
	XMLName  xml.Name     `xml:"F05DOC"`
	Head     F05DocHead   `xml:"DOCHEAD"`
	Body     []F05Content `xml:"F05CONTENT"`
}

type F05DocHead struct {
	DocDate         string `xml:"DOCDATE"`
	FromPartnerCode string `xml:"FROMPARTNERCODE"`
	ToPartenerCode  string `xml:"TOPARTNERCODE"`
}

type F05Content struct {
	BarCode1        string `xml:"BC1"`
	BarCode2        string `xml:"BC2"`
	RrStoreId       string `xml:"STNO"`
	UpDateTime      string `xml:"RTDT"`
	CheckDateTime   string `xml:"TKDT"`
	VendorNo        string `xml:"VENDORNO"`
}

func (receiver *F05Doc) DecodeXML(data []byte) (err error) {
	err = xml.Unmarshal(data, receiver)
	return
}

