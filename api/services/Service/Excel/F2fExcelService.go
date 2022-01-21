package Excel

func F2fNew() *excelFormats {
	fs := &excelFormats{}
	tag := []string{"A", "B", "C", "D", "E", "F", "G", "H"}
	title := []string{"序號", "訂單編號", "商品名稱", "單價", "件數", "收件人-姓名", "收件人-手機", "備註", "買家備註"}
	column := []string{"Id", "OrderId", "ProductName", "Price", "Pieces", "ReceiverName", "ReceiverPhone", "OrderMemo", "BuyerNotes"}
	for k, v := range tag {
		SetCell(v, column[k], title[k], fs)
	}
	return fs
}