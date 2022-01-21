package OKMart

import (
	"errors"
	"testing"
)

func TestF71Doc_DecodeXML(t *testing.T) {
	dataStr := `<F71DOC>
		<DOCHEAD>
			<DOCDATE>20180913</DOCDATE>
			<FROMPARTNERCODE>CVS</FROMPARTNERCODE>
			<TOPARTNERCODE>526</TOPARTNERCODE>
      	</DOCHEAD>
     	<F71CONTENT>
			<ECNO>526</ECNO>
			<ODNO>12345678901</ODNO>
			<STNO>K001234</STNO>
			<RTDT>20180913195004</RTDT>
		</F71CONTENT>
		<F71CONTENT>
			<ECNO>526</ECNO>
			<ODNO>12345678901</ODNO>
			<STNO>K001234</STNO>
			<RTDT>20180913195004</RTDT>
		</F71CONTENT>
</F71DOC>`
	data := []byte(dataStr)

	doc := F71Doc{}
	err := doc.DecodeXML(data)
	if err != nil {
		t.Errorf("DecodeXML() error = %v", err)
	}

	head := doc.Head
	if head.DocDate != "20180913" || head.ToPartenerCode != "526" || head.FromPartnerCode != "CVS" {
		t.Fatalf("DecodeXML() error = %v", errors.New("Head not match"))
	}

	if len(doc.Body) != 2 {
		t.Fatalf("DecodeXML() error = %v", errors.New("Content Length not match"))
	}
	content := doc.Body[0]
	if  content.EcNo != "526" || content.OrderNo != "12345678901" {
		t.Fatalf("DecodeXML() error = %v", errors.New("Barcode not match"))
	}

	if  content.RrStoreId != "K001234" || content.UpDateTime != "20180913195004" {
		t.Fatalf("DecodeXML() error = %v", errors.New("FromPartnerCode not match"))
	}
}
