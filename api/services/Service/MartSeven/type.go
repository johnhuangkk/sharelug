package MartSeven

import "encoding/xml"

type Store struct {
	Id      string
	Name    string
	Address string
	Detail  string
}

type OrderAddRequest struct {
	XMLName			 xml.Name `xml:"C2C"`
	EshopId          string `xml:"eshopid"`
	EshopSonId       string `xml:"eshopsonid"`
	OrderNo          string `xml:"orderno"`
	ServiceType      string `xml:"service_type"`
	Account          string `xml:"account"`
	PaymentCpName    string `xml:"payment_cpname"`
	TradeDescription string `xml:"trade_description"`
	CpRemark01       string `xml:"cp_remark01"`
	CpRemark02       string `xml:"cp_remark02"`
	CpRemark03       string `xml:"cp_remark03"`
	DeadlineDate     string `xml:"deadlinedate"`
	DeadlineTime     string `xml:"deadlinetime"`
	ShowType         string `xml:"show_type"`
	DaishouAccount   string `xml:"daishou_account"`
	Sender           string `xml:"sender"`
	SenderPhone      string `xml:"sender_phone"`
	Receiver         string `xml:"receiver"`
	ReceiverPhone    string `xml:"receiver_phone"`
	ReceiverStoreId  string `xml:"receiver_storeid"`
	ReturnStoreId    string `xml:"return_storeid"`
}

func (receiver* OrderAddRequest) EncodeXML() ([]byte, error) {
	data,err := xml.Marshal(receiver)
	if err != nil {
		return nil, err
	}
	headerData := []byte(xml.Header)
	headerData = append(headerData,data...)
	return headerData, nil
}