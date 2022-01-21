package OKMart

import (
	"errors"
	"testing"
)

func TestF27Doc_DecodeXML(t *testing.T) {
	dataStr := `
	<F27DOC>
		<DOCHEAD>
			<DOCDATE>20180912</DOCDATE>
			<FROMPARTNERCODE>CVS</FROMPARTNERCODE>
			<TOPARTNERCODE>526</TOPARTNERCODE>
		</DOCHEAD>
		<F27CONTENT>
			<ECNO>526</ECNO>
			<ODNO>12345678901</ODNO>
			<STNO>K001234</STNO>
			<RTDT>20180912203045</RTDT>
		</F27CONTENT>
		<F27CONTENT>
			<ECNO>526</ECNO>
			<ODNO>12345678901</ODNO>
			<STNO>K001234</STNO>
			<RTDT>20180912203045</RTDT>
		</F27CONTENT>
	</F27DOC>`
	data := []byte(dataStr)

	doc := F27Doc{}
	err := doc.DecodeXML(data)
	if err != nil {
		t.Errorf("DecodeXML() error = %v", err)
	}

	head := doc.Head
	if head.DocDate != "20180912" || head.ToPartenerCode != "526" || head.FromPartnerCode != "CVS" {
		t.Fatalf("DecodeXML() error = %v", errors.New("Head not match"))
	}

	if len(doc.Body) != 2 {
		t.Fatalf("DecodeXML() error = %v", errors.New("Content Length not match"))
	}
	content := doc.Body[0]
	if  content.EcNo != "526" || content.OrderNo != "12345678901" {
		t.Fatalf("DecodeXML() error = %v", errors.New("Body not match"))
	}

	if  content.RrStoreId != "K001234" || content.UpDateTime != "20180912203045" {
		t.Fatalf("DecodeXML() error = %v", errors.New("Body not match"))
	}
}
