package OKMart

import "encoding/xml"

//<F01DOC>
//		<DOCHEAD>
//			<DOCDATE>20201110</DOCDATE><FROMPARTNERCODE>CVS</FROMPARTNERCODE><TOPARTENERCODE>462</TOPARTENERCODE>
//		</DOCHEAD>
//		<F01CONTENT>
//			<STNO>K000002</STNO><STNM>福林店</STNM><STTEL>02-66170259</STTEL>
//			<STCITY>台北市</STCITY><STCNTRY>士林區</STCNTRY>
//			<STADR>台北市士林區中山北路５段702號</STADR><ZIPCD>111</ZIPCD>
//			<DCRONO></DCRONO><SDATE>00000000</SDATE><EDATE>00000000</EDATE>
//		</F01CONTENT>
//		<F01CONTENT>
//			<STNO>K000004</STNO><STNM>港墘店</STNM><STTEL>02-66170386</STTEL>
//			<STCITY>台北市</STCITY><STCNTRY>內湖區</STCNTRY>
//			<STADR>台北市內湖區港墘路84號</STADR><ZIPCD>114</ZIPCD>
//			<DCRONO></DCRONO><SDATE>00000000</SDATE><EDATE>00000000</EDATE>
//		</F01CONTENT>
//</F01DOC>

type F01Doc struct {
	Head     F01DocHead   `xml:"DOCHEAD"`
	Contents []F01Content `xml:"F01CONTENT"`
}

type F01DocHead struct {
	DOCDATE         string `xml:"DOCDATE"`
	FROMPARTNERCODE string `xml:"FROMPARTNERCODE"`
	TOPARTENERCODE  string `xml:"TOPARTENERCODE"`
}

type F01Content struct {
	StoreId      string `xml:"STNO"`
	StoreName    string `xml:"STNM"`
	StoreTel     string `xml:"STTEL"`
	StoreCity    string `xml:"STCITY"`
	StoreDisct   string `xml:"STCNTRY"`
	StoreAddress string `xml:"STADR"`
	Zipcode      string `xml:"ZIPCD"`
	DCRONO       string `xml:"DCRONO"`
	SDATE        string `xml:"SDATE"`
	EDATE        string `xml:"EDATE"`
}

func (receiver *F01Doc) DecodeXML(data []byte) (err error) {
	err = xml.Unmarshal(data, receiver)
	return
}
