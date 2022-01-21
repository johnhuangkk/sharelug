package FamilyMartLogistics

import (
	"encoding/xml"
	"time"
)

// 商品寄件檔/多批次(R22)
/*
<?xml version="1.0" encoding="utf-8"?>
<doc>
	<HEADER>
		<RDFMT><![CDATA[1]]></RDFMT>
		<SNCD><![CDATA[DFM]]></SNCD>
		<PRDT><![CDATA[2020-06-01]]></PRDT>
	</HEADER>
	<BODY>
		<R22>
			<RDFMT><![CDATA[2]]></RDFMT>
			<ParentId><![CDATA[861]]></ParentId>
			<EshopId><![CDATA[0001]]></EshopId>
			<OrderNo><![CDATA[03918823876]]></OrderNo>
			<OrderDate><![CDATA[2020-06-01]]></OrderDate>
			<OrderTime><![CDATA[201911]]></OrderTime>
			<OPMode><![CDATA[1]]></OPMode>
			<StoreId><![CDATA[010459]]></StoreId>
		</R22>
	</BODY>
	<FOOTER>
		<RDFMT><![CDATA[3]]></RDFMT>
		<RDCNT><![CDATA[1]]></RDCNT>
	</FOOTER>
</doc>
*/

type FMLR22 struct {
	Header FMLDataHeader `xml:"HEADER"`
	Body   struct {
		Data []FMLR22Data `xml:"R22"`
	} `xml:"BODY"`
	Footer FMLDataFooter `xml:"FOOTER"`
}

type FMLR22Data struct {
	RDFMT     string   `xml:"RDFMT"`
	ParentId  string   `xml:"ParentId"`
	EshopId   string   `xml:"EshopId"`
	OrderNo   string   `xml:"OrderNo"`
	OrderDate string   `xml:"OrderDate"`
	OrderTime string   `xml:"OrderTime"`
	OPMode    string   `xml:"OPMode"`
	StoreId   string   `xml:"StoreId"`
}

func (s *FMLR22Data) GetDateTime () string {
	t, _ := time.Parse(`2006-01-02150405`, s.OrderDate + s.OrderTime)
	return t.Format(`2006-01-02 15:04:05`)
}

func (receiver *FMLR22) DecodeXML(data []byte) error {
	return xml.Unmarshal(data, receiver)
}
