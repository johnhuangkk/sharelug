package MartHiLife

//<?xml version="1.0" encoding="utf-8"?>
//<doc>
//	<HEADER>
//		<RDFMT>1</RDFMT>
//		<SNCD>HIL</SNCD>
//		<PRDT>2019-01-19</PRDT>
//	</HEADER>
//	<BODY>

//	</BODY>
//	<FOOTER>
//			<RDFMT>3</RDFMT>
//			<RDCNT>1</RDCNT>
//	</FOOTER>
//</doc>

type RHeader struct {
	RDFMT         string `xml:"RDFMT"`
	SNCD          string `xml:"SNCD"`
	PRDT          string `xml:"PRDT"`
}

type DecodeXML interface {
	DecodeXML(data []byte) (err error)
}
