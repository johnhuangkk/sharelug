package InvoiceXml

type RejectInvoiceConfirm struct {
	RejectInvoiceNumber string `xml:"RejectInvoiceNumber"`  //退回(拒收)發票號碼
	InvoiceDate         string `xml:"InvoiceDate"`			//發票日期
	BuyerId             string `xml:"BuyerId"`				//買方統一編號
	SellerId            string `xml:"SellerId"`				//賣方統一編號
	RejectDate          string `xml:"RejectDate"`			//退回(拒收)日期
	RejectTime          string `xml:"RejectTime"`			//退回(拒收)時間
	Remark              string `mxl:"Remark"`				//備註
}
