package FamilyMartLogistics

import (
	"encoding/xml"
	"time"
)

/*
<?xml version="1.0" encoding="UTF-8"?> <doc>
<HEADER>
	<RDFMT>1</RDFMT>
	<SNCD>DFM</SNCD>
	<PRDT>2020-06-01</PRDT>
</HEADER>
<BODY>
	<R25>
		<RDFMT>2</RDFMT>
		<ParentId>861</ParentId>
		<EshopId>0001</EshopId>
		<ShipmentNo>0001100340169</ShipmentNo>
		<DCReceiveDate>2020-06-01</DCReceiveDate>
		<DCReceiveTime>141300</DCReceiveTime>
		<DCReceiveStatus>1</DCReceiveStatus>
		<FlowType>N</FlowType>
	</R25>
</BODY>
<FOOTER>
	<RDFMT>3</RDFMT>
	<RDCNT>1</RDCNT>
</FOOTER>
</doc>
 */


type FMLR25 struct {
	Header FMLDataHeader `xml:"HEADER"`
	Body  struct {
		Data []FMLR25Data `xml:"R25"`
	} `xml:"BODY"`
	Footer FMLDataFooter `xml:"FOOTER"`
}

type FMLR25Data struct {
	RDFMT           string   `xml:"RDFMT"`
	ParentId        string   `xml:"ParentId"`
	EshopId         string   `xml:"EshopId"`
	ShipmentNo      string   `xml:"ShipmentNo"`
	TotalAmount     string   `xml:"TotalAmount"`
	DCReceiveDate   string   `xml:"DCReceiveDate"`
	DCReceiveTime   string   `xml:"DCReceiveTime"`
	DCReceiveStatus string   `xml:"DCReceiveStatus"`
	FlowType        string   `xml:"FlowType"`
}




func (s *FMLR25Data) GetDateTime () string {
	t, _ := time.Parse(`2006-01-02150405`, s.DCReceiveDate + s.DCReceiveTime)
	return t.Format(`2006-01-02 15:04:05`)
}

func (receiver *FMLR25) DecodeXML(data []byte) error {
	return xml.Unmarshal(data, receiver)
}