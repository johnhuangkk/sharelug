package FamilyMartLogistics

import (
	"encoding/xml"
)

/*
<?xml version="1.0" encoding="utf-8"?>
<doc>
    <HEADER><RDFMT>1</RDFMT><SNCD>DFM</SNCD><PRDT>2020-06-01</PRDT></HEADER>
    <BODY>
        <RS9>
            <RDFMT>2</RDFMT>
            <ParentId>861</ParentId>
            <EshopId>0001</EshopId>
            <ShipmentNo>800006429</ShipmentNo>
            <ReceiveStoreId>010307</ReceiveStoreId>
            <DCReceiveDate>2020-06-01</DCReceiveDate>
            <DCReceiveStatus>2</DCReceiveStatus>
            <StatusDetails>N05</StatusDetails>
            <StatusRemark>門市遺失</StatusRemark>
            <FlowType>N</FlowType>
        </RS9>
        <RS9><RDFMT>2</RDFMT><ParentId>861</ParentId><EshopId>0001</EshopId><ShipmentNo>800006430</ShipmentNo><ReceiveStoreId>003822</ReceiveStoreId><DCReceiveDate>2020-06-01</DCReceiveDate><DCReceiveStatus>3</DCReceiveStatus><StatusDetails>D04</StatusDetails><StatusRemark>包裝廠包裝不良(滲漏) </StatusRemark><FlowType>N</FlowType></RS9>
        <RS9><RDFMT>2</RDFMT><ParentId>861</ParentId><EshopId>0001</EshopId><ShipmentNo>800006411</ShipmentNo><ReceiveStoreId>010458</ReceiveStoreId><DCReceiveDate>2020-06-01</DCReceiveDate><DCReceiveStatus>5</DCReceiveStatus><StatusDetails>999</StatusDetails><StatusRemark>店家未到貨</StatusRemark><FlowType>N</FlowType></RS9>
    </BODY>
    <FOOTER><RDFMT>3</RDFMT><RDCNT>3</RDCNT></FOOTER>
</doc>
*/

type FMLRS9 struct {
	Header FMLDataHeader `xml:"HEADER"`
	Body   struct {
		Data []FMLRS9Data `xml:"RS9"`
	} `xml:"BODY"`
	Footer FMLDataFooter `xml:"FOOTER"`
}

type FMLRS9Data struct {
	RDFMT           string `xml:"RDFMT"`
	ParentId        string `xml:"ParentId"`
	EshopId         string `xml:"EshopId"`
	ShipmentNo      string `xml:"ShipmentNo"`
	ReceiveStoreId  string `xml:"ReceiveStoreId"`
	DCReceiveDate   string `xml:"DCReceiveDate"`
	DCReceiveStatus string `xml:"DCReceiveStatus"`
	StatusDetails   string `xml:"StatusDetails"`
	StatusRemark    string `xml:"StatusRemark"`
	FlowType        string `xml:"FlowType"`
}

func (receiver *FMLRS9) DecodeXML(data []byte) error {
	return xml.Unmarshal(data, receiver)
}
