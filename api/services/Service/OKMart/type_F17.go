package OKMart

import (
	"encoding/xml"
	"time"
)

//<F17DOC>
//    <DOCHEAD>
//        <DOCDATE>20201116</DOCDATE>
//        <FROMPARTNERCODE>CVS</FROMPARTNERCODE>
//        <TOPARTENERCODE>462</TOPARTENERCODE>
//    </DOCHEAD>
//	  <F17CONTENT>
//	  	<BC1> 516725992 </BC1>
//	  	<BC2> 3916921910002084 </BC2>
//	  	<STNO>K001234</STNO>
//	  	<RTDT> 20180914140036 </RTDT>
//	  	<PINCODE></PINCODE>
//	  </F17CONTENT>
//	  <F17CONTENT>
//	  <BC1> 516725992 </BC1>
//	  <BC2> 3916921910002084 </BC2> <STNO>K001234</STNO> <RTDT> 20180914140036 </RTDT> <PINCODE></PINCODE>
//	  </F17CONTENT>
//</F17DOC>


type F17Doc struct {
	XMLName  xml.Name     `xml:"F17DOC"`
	Head     F17DocHead   `xml:"DOCHEAD"`
	Body     []F17Content `xml:"F17CONTENT"`
}

type F17DocHead struct {
	DocDate         string `xml:"DOCDATE"`
	FromPartnerCode string `xml:"FROMPARTNERCODE"`
	ToPartenerCode  string `xml:"TOPARTENERCODE"`
}

type F17Content struct {
	BarCode1        string `xml:"BC1"`
	BarCode2        string `xml:"BC2"`
	RrStoreId       string `xml:"STNO"`
	RrPickDateTime  string `xml:"RTDT"`
	SendNo          string `xml:"PINCODE"`
}

func (s *F17Content) GetDateTime () string {
	t, _ := time.Parse(`20060102150405`, s.RrPickDateTime)
	return t.Format(`2006-01-02 15:04:05`)
}

func (receiver *F17Doc) DecodeXML(data []byte) (err error) {
	err = xml.Unmarshal(data, receiver)
	return
}

