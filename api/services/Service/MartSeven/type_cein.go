package MartSeven

import "encoding/xml"

//<DCReceiveDoc>
//	<DocHead>
//		<DocNo>2009010200000123</DocNo>
//		<DocDate>2009-01-02</DocDate>
//		<From>
//			<FromPartnerCode>000</FromPartnerCode>
//		</From>
//		<To>
//			<ToPartnerCode>M16</ToPartnerCode>
//		</To>
//		<DocCount>2</DocCount>
//	</DocHead>
//	<DocContent>
//		<DCReceive>
//			<ParentId>M16</ParentId>
//			<EshopId>A01</EshopId>
//			<EC3GParentId>935</EC3GParentId>
//			<EC3GEshopId>A01</EC3GEshopId>
//			<PaymentNo>C0000001</PaymentNo>
//			<ShipmentNo>00000001</ShipmentNo>
//			<DCReceiveDate>2009-01-02</DCReceiveDate>
//			<DCStoreStatus>+</DCStoreStatus>
//			<DCReceiveStatus>00</DCReceiveStatus>
//			<DCRecName>進驗成功</DCRecName>
//			<DCStoreDate>2009-01-03</DCStoreDate>
//		</DCReceive>
//		<DCReceive>
//			<ParentId>M16</ParentId>
//			<EshopId>A01</EshopId>
//			<EC3GParentId>935</EC3GParentId>
//			<EC3GEshopId>A01</EC3GEshopId>
//			<PaymentNo>C0000002</PaymentNo>
//			<ShipmentNo>00000002</ShipmentNo>
//			<DCReceiveDate>2009-01-02</DCReceiveDate>
//			<DCStoreStatus>+</DCStoreStatus>
//			<DCReceiveStatus>36</DCReceiveStatus>
//			<DCRecName>門市關轉</DCRecName>
//			<DCStoreDate></DCStoreDate>
//		</DCReceive>
//	</DocContent>
//</DCReceiveDoc>


type CEIN struct {
	XMLName xml.Name `xml:"DCReceiveDoc"`
	Body  struct{
		Contents []CEINContent `xml:"DCReceive"`
	} `xml:"DocContent"`
}

type CEINContent struct {
	ParentId        string `xml:"ParentId"`
	EshopId         string `xml:"EshopId"`
	EC3GParentId    string `xml:"EC3GParentId"`
	EC3GEshopId     string `xml:"EC3GEshopId"`
	PaymentNo       string `xml:"PaymentNo"`
	ShipmentNo      string `xml:"ShipmentNo"`
	DCReceiveDate   string `xml:"DCReceiveDate"`
	DCStoreStatus   string `xml:"DCStoreStatus"`
	DCReceiveStatus string `xml:"DCReceiveStatus"`
	DCRecName       string `xml:"DCRecName"`
	DCStoreDate     string `xml:"DCStoreDate"`
}

func (receiver *CEIN) DecodeXML(data []byte) error {
	return xml.Unmarshal(data, &receiver)
}
