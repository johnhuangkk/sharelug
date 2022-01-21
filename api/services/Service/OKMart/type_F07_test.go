package OKMart

import (
	"errors"
	"testing"
)

func TestF07Doc_DecodeXML(t *testing.T) {
	dataStr := `<F07DOC>
		<DOCHEAD>
			<DOCDATE>20180915</DOCDATE>
			<FROMPARTNERCODE>CVS</FROMPARTNERCODE>
			<TOPARTNERCODE>416</TOPARTNERCODE>
		</DOCHEAD>
		<F07CONTENT>
			<RET_M>T08</RET_M><ECNO>416</ECNO>
			<STNO>K001234</STNO><ODNO>12345678901</ODNO>
			<RTDCDT>20180915200006</RTDCDT><FRTDCDT>20180915200006</FRTDCDT>
		</F07CONTENT>
		<F07CONTENT>
			<RET_M>T08</RET_M><ECNO>416</ECNO>
			<STNO>K001234</STNO><ODNO>12345678901</ODNO>
			<RTDCDT>20180915200006</RTDCDT><FRTDCDT>20180915200006</FRTDCDT>
		</F07CONTENT>
	</F07DOC>`

	doc := F07Doc{}
	err := doc.DecodeXML([]byte(dataStr))
	if err != nil {
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
	if  content.EcNo != "416" || content.OrderNo != "12345678901" {
		t.Fatalf("DecodeXML() error = %v", errors.New("Body not match"))
	}

	if  content.ReturnCode != "T08" || content.RrStoreId != "K001234" {
		t.Fatalf("DecodeXML() error = %v", errors.New("Body not match"))
	}
	if  content.UpDateTime != "20180915200006" || content.CheckDateTime != "20180915200006" {
		t.Fatalf("DecodeXML() error = %v", errors.New("Body not match"))
	}
}
