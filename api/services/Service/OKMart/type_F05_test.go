package OKMart

import (
	"errors"
	"testing"
)

func TestF05Doc_DecodeXML(t *testing.T) {
	dataStr := `<F05DOC>
		<DOCHEAD>
			<DOCDATE>20180915</DOCDATE>
			<FROMPARTNERCODE>CVS</FROMPARTNERCODE>
			<TOPARTNERCODE>416</TOPARTNERCODE>
		</DOCHEAD>
		<F05CONTENT>
			<BC1>416375993</BC1>
			<BC2>3849859930000015</BC2>
			<STNO>K001234</STNO>
			<RTDT>20180915230037</RTDT>
			<TKDT>20180915230037</TKDT>
			<VENDORNO>TW1234567890</VENDORNO>
		</F05CONTENT>
		<F05CONTENT>
			<BC1>416375993</BC1>
			<BC2>3849859930000015</BC2>
			<STNO>K001234</STNO>
			<RTDT>20180915230037</RTDT>
			<TKDT>20180915230037</TKDT>
			<VENDORNO>TW1234567890</VENDORNO>
		</F05CONTENT>
	</F05DOC>`
	data := []byte(dataStr)

	doc := F05Doc{}
	if 	err := doc.DecodeXML(data); err != nil {
		t.Fatalf("DecodeXML() error = %v", err)
	}

	head := doc.Head
	if head.DocDate != "20180915" || head.ToPartenerCode != "416" || head.FromPartnerCode != "CVS" {
		t.Fatalf("DecodeXML() error = %v", errors.New("Head not match"))
	}

	if len(doc.Body) != 2 {
		t.Fatalf("DecodeXML() error = %v", errors.New("Content Length not match"))
	}
	content := doc.Body[0]
	if  content.BarCode1 != "416375993" || content.BarCode2 != "3849859930000015" {
		t.Fatalf("DecodeXML() error = %v", errors.New("Body not match"))
	}

	if  content.RrStoreId != "K001234" || content.UpDateTime != "20180915230037" {
		t.Fatalf("DecodeXML() error = %v", errors.New("Body not match"))
	}

	if  content.CheckDateTime != "20180915230037" || content.VendorNo != "TW1234567890" {
		t.Fatalf("DecodeXML() error = %v", errors.New("Body not match"))
	}
}
