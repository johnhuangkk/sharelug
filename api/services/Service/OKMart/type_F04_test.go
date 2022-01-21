package OKMart

import (
	"errors"
	"testing"
)

func TestF04Doc_DecodeXML(t *testing.T) {
	dataStr := `<F04DOC>
		<DOCHEAD>
			<DOCDATE>20180915</DOCDATE>
			<FROMPARTNERCODE>CVS</FROMPARTNERCODE>
			<TOPARTNERCODE>417</TOPARTNERCODE>
		</DOCHEAD>
		<F04CONTENT>
			<ECNO>417</ECNO>
			<STNO>K001234</STNO>
			<ODNO>123456778901</ODNO>
			<DCSTDT>20180915025928</DCSTDT>
		</F04CONTENT>
      <F04CONTENT>
 			<ECNO>417</ECNO>
 			<STNO>K001234</STNO>
 			<ODNO>123456778901</ODNO>
 			<DCSTDT>20180915025928</DCSTDT>
 		</F04CONTENT>
	</F04DOC>`
	data := []byte(dataStr)

	doc := F04Doc{}
	err := doc.DecodeXML(data)
	if err != nil {
		t.Fatalf("DecodeXML() error = %v", err)
	}

	head := doc.Head
	if head.DocDate != "20180915" || head.ToPartenerCode != "417" || head.FromPartnerCode != "CVS" {
		t.Fatalf("DecodeXML() error = %v", errors.New("Head not match"))
	}

	if len(doc.Body) != 2 {
		t.Fatalf("DecodeXML() error = %v", errors.New("Content Length not match"))
	}
	content := doc.Body[0]
	if  content.EcNo != "417" || content.OrderNo != "123456778901" {
		t.Fatalf("DecodeXML() error = %v", errors.New("Body not match"))
	}

	if  content.RrStoreId != "K001234" || content.UpDateTime != "20180915025928" {
		t.Fatalf("DecodeXML() error = %v", errors.New("Body not match"))
	}
}
