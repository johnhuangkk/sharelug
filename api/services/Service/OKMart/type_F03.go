package OKMart

import (
	"encoding/xml"
)

//<F03DOC>
//	<DOCHEAD>
//		<DOCDATE>20180915</DOCDATE><FROMPARTNERCODE>CVS</FROMPARTNERCODE><TOPARTNERCODE>417</TOPARTNERCODE>
//	</DOCHEAD>
//	<F03CONTENT>
//		<ECNO>417</ECNO><ODNO>12345678901</ODNO><CUNAME>簡大翔</CUNAME>
//		<PRODTYPE>0</PRODTYPE><PINCODE>12345678901</PINCODE>
//		<RTDT>20180915172249</RTDT>
//	</F03CONTENT>
//	<F03CONTENT>
//		<ECNO>417</ECNO><ODNO>12345678901</ODNO><CUNAME>簡大翔</CUNAME>
//		<PRODTYPE>0</PRODTYPE><PINCODE>12345678901</PINCODE>
//		<RTDT>20180915172249</RTDT>
//	</F03CONTENT>
//</F03DOC>


type F03Doc struct {
	XMLName  xml.Name     `xml:"F03DOC"`
	Head     F03DocHead   `xml:"DOCHEAD"`
	Body     []F03Content `xml:"F03CONTENT"`
}

type F03DocHead struct {
	DocDate         string `xml:"DOCDATE"`
	FromPartnerCode string `xml:"FROMPARTNERCODE"`
	ToPartenerCode  string `xml:"TOPARTNERCODE"`
}

type F03Content struct {
	EcNo            string `xml:"ECNO"`
	OrderNo         string `xml:"ODNO"`
	RrName          string `xml:"CUNAME"`
	ProductType     string `xml:"PRODTYPE"`
	UpDateTime      string `xml:"RTDT"`
	SendNo          string `xml:"PINCODE"`
}

func (receiver *F03Doc) DecodeXML(data []byte) (err error) {
	err = xml.Unmarshal(data, receiver)
	return
}

