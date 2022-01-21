package OKMart

import (
	"errors"
	"testing"
)

func TestF84Doc_DecodeXML(t *testing.T) {
	dataStr := `
	<F84DOC>
		<DOCHEAD>
			<DOCDATE>20180913</DOCDATE>
			<FROMPARTNERCODE>CVS</FROMPARTNERCODE>
			<TOPARTNERCODE >526</TOPARTNERCODE>
		</DOCHEAD>
		<F84CONTENT>
			<ECNO>526</ECNO>
			<STNO>K001234</STNO>
			<ODNO>12345678901</ODNO>
			<DCSTDT>20180913142848</DCSTDT>
			<EASYECNO></EASYECNO>
		</F84CONTENT>
		<F84CONTENT>
			<ECNO>526</ECNO><STNO>K001234</STNO><ODNO>12345678901</ODNO>
			<DCSTDT>20180913142848</DCSTDT><EASYECNO></EASYECNO>
		</F84CONTENT>
	</F84DOC>`
	data := []byte(dataStr)

	doc := F84Doc{}
	err := doc.DecodeXML(data)
	if err != nil {
		t.Errorf("DecodeXML() error = %v", err)
	}

	head := doc.Head
	if head.DocDate != "20180913" || head.ToPartenerCode != "526" || head.FromPartnerCode != "CVS" {
		t.Fatalf("DecodeXML() error = %v", errors.New("DocDate not match"))
	}

	if len(doc.Body) != 2 {
		t.Fatalf("DecodeXML() error = %v", errors.New("Content Length not match"))
	}
	content := doc.Body[0]
	if  content.EcNo != "526" || content.OrderNo != "12345678901" {
		t.Fatalf("DecodeXML() error = %v", errors.New("Body not match"))
	}

	if  content.RrStoreId != "K001234" || content.UpDateTime != "20180913142848" {
		t.Fatalf("DecodeXML() error = %v", errors.New("Body not match"))
	}

	if  content.OtherCode != "" {
		t.Fatalf("DecodeXML() error = %v", errors.New("Body not match"))
	}
}
