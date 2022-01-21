package MartSeven

import "encoding/xml"

//<SPDoc>
//	<DocHead>
//		<DocNo>2009010909520125</DocNo>
//		<DocDate>2009-01-09</DocDate>
//		<From>
//			<FromPartnerCode>000</FromPartnerCode>
//		</From>
//		<To>
//			<ToPartnerCode>M16</ToPartnerCode>
//		</To>
//		<DocCount>1</DocCount>
//	</DocHead>
//	<DocContent>
//		<SP>
//			<ParentId>M16</ParentId>
//			<TotalCount>1</TotalCount>
//			<TotalAmount>1999</TotalAmount>
//			<SPDetail>
//				<ParentId>M16</ParentId><EshopId>A01</EshopId><EC3GParentId>935</EC3GParentId><EC3GEshopId>A01</EC3GEshopId>
//				<PaymentNo>C0000001</PaymentNo><DCStoreStatus>+</DCStoreStatus><StoreId>886750</StoreId><SPDate>2009-01-07</SPDate>
//				<SPAmount>1999</SPAmount><ServiceType>6</ServiceType><SPNo>00000001</SPNo>
//			</SPDetail>
//          <SPDetail>
//				<ParentId>M16</ParentId><EshopId>A01</EshopId><EC3GParentId>935</EC3GParentId><EC3GEshopId>A01</EC3GEshopId>
//				<PaymentNo>C0000001</PaymentNo>
//				<DCStoreStatus>+</DCStoreStatus>
//				<StoreId>886750</StoreId><SPDate>2009-01-07</SPDate>
//				<SPAmount>1999</SPAmount><ServiceType>6</ServiceType><SPNo>00000001</SPNo>
//			</SPDetail>
//		</SP>
//	</DocContent>
//</SPDoc>

type CESP struct {
	XMLName xml.Name `xml:"SPDoc"`
	Body    struct {
		Contents SP `xml:"SP"`
	} `xml:"DocContent"`
}

type SP struct {
	ParentId    string     `xml:"ParentId"`
	TotalCount  string     `xml:"TotalCount"`
	TotalAmount string     `xml:"TotalAmount"`
	Details     []SPDetail `xml:"SPDetail"`
}

type SPDetail struct {
	ParentId      string `xml:"ParentId"`
	EshopId       string `xml:"EshopId"`
	EC3GParentId  string `xml:"EC3GParentId"`
	EC3GEshopId   string `xml:"EC3GEshopId"`
	PaymentNo     string `xml:"PaymentNo"`
	DCStoreStatus string `xml:"DCStoreStatus"`
	StoreId       string `xml:"StoreId"`
	SPDate        string `xml:"SPDate"`
	SPAmount      string `xml:"SPAmount"`
	ServiceType   string `xml:"ServiceType"`
	SPNo          string `xml:"SPNo"`
}

func (receiver *CESP) DecodeXML(data []byte) error {
	return xml.Unmarshal(data, &receiver)
}
