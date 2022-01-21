package FamilyMartLogistics

import (
	"encoding/xml"
	"time"
)

// 寄件離店檔/多批次(R23)
/*
<?xml version="1.0" encoding="utf-8"?>
<doc>
	<HEADER>
		<RDFMT><![CDATA[1]]></RDFMT>
		<SNCD><![CDATA[DFM]]></SNCD>
		<PRDT><![CDATA[2020-06-01]]></PRDT>
	</HEADER>
	<BODY>
		<R23>
			<RDFMT><![CDATA[2]]></RDFMT>
			<ParentId><![CDATA[861]]></ParentId>
			<EshopId><![CDATA[0001]]></EshopId>
			<OrderNo><![CDATA[03918823876]]></OrderNo>
			<OrderDate><![CDATA[2020-06-01]]></OrderDate>
			<OrderTime><![CDATA[141300]]></OrderTime>
			<OPMode><![CDATA[1]]></OPMode>
			<StoreId><![CDATA[005989]]></StoreId>
		</R23>
		<R23>
			<RDFMT><![CDATA[2]]></RDFMT>
			<ParentId><![CDATA[861]]></ParentId>
			<EshopId><![CDATA[0001]]></EshopId>
			<OrderNo><![CDATA[03918823877]]></OrderNo>
			<OrderDate><![CDATA[2020-06-01]]></OrderDate>
			<OrderTime><![CDATA[141300]]></OrderTime>
			<OPMode><![CDATA[1]]></OPMode>
			<StoreId><![CDATA[005989]]></StoreId>
		</R23>
	</BODY>
	<FOOTER>
		<RDFMT><![CDATA[3]]></RDFMT>
		<RDCNT><![CDATA[2]]></RDCNT>
	</FOOTER>
</doc>
*/

type FMLR23 struct {
	Header FMLDataHeader `xml:"HEADER"`
	Body   struct {
		Data []FMLR23Data `xml:"R23"`
	} `xml:"BODY"`
	Footer FMLDataFooter `xml:"FOOTER"`
}

type FMLR23Data struct {
	RDFMT     string `xml:"RDFMT"`
	ParentId  string `xml:"ParentId"`
	EshopId   string `xml:"EshopId"`
	OrderNo   string `xml:"OrderNo"`
	OrderDate string `xml:"OrderDate"`
	OrderTime string `xml:"OrderTime"`
	OPMode    string `xml:"OPMode"`
	StoreId   string `xml:"StoreId"`
}

func (s *FMLR23Data) GetDateTime () string {
	t, _ := time.Parse(`2006-01-02150405`, s.OrderDate + s.OrderTime)
	return t.Format(`2006-01-02 15:04:05`)
}

func (receiver *FMLR23) DecodeXML(data []byte) error {
	return xml.Unmarshal(data, receiver)
}
