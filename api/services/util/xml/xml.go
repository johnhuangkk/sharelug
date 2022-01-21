package xml

import (
	"api/services/entity"
	"api/services/entity/Response"
	"api/services/util/log"
	"bytes"
	"encoding/json"
	"encoding/xml"
	"strings"

	xj "github.com/basgys/goxml2json"
	"github.com/tidwall/gjson"
)

func SmsXmlEncoder(v interface{}) ([]byte, error) {
	b := &bytes.Buffer{}
	b.Write([]byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n"))
	enc := xml.NewEncoder(b)
	enc.Indent("", "")
	if err := enc.Encode(v); err != nil {
		log.Error("error: %v\n", err)
		return nil, err
	}
	return b.Bytes(), nil
}

func SmsXmlDecoder(data string) (entity.SmsResult, error) {
	result := entity.SmsResult{}
	err := xml.Unmarshal([]byte(data), &result)
	if err != nil {
		log.Error("xml Decoder error !!", err)
		return result, err
	}
	return result, nil
}

func TransferXmlEncoder(v interface{}) ([]byte, error) {
	b := &bytes.Buffer{}
	b.Write([]byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?><S:Envelope xmlns:S=\"http://schemas.xmlsoap.org/soap/envelope/\">" +
		"<S:Body><getData xmlns=\"http://ws.kgibank.thesys.com\"><xstring>"))
	enc := xml.NewEncoder(b)
	enc.Indent("", "")
	if err := enc.Encode(v); err != nil {
		log.Error("error: %v\n", err)
		return nil, err
	}
	b.Write([]byte("</xstring></getData></S:Body></S:Envelope>\n"))
	return b.Bytes(), nil
}
//解析XML
func TransferXmlDecoder(xml string) (Response.SMX, []Response.DETAIL, error) {
	log.Debug("xml to json 0", xml)
	s := strings.NewReader(xml)
	content, _ := xj.Convert(s)
	body := gjson.Get(content.String(), "Envelope.Body.getDataResponse.getDataReturn.SMX")
	var smx Response.SMX
	if len(body.String()) != 0 {
		if err := json.Unmarshal([]byte(body.String()), &smx); err != nil {
			log.Error("transfer xml Decoder Error", err)
			return smx, nil,  err
		}
	}
	//判斷array
	detail := gjson.Get(content.String(), "Envelope.Body.getDataResponse.getDataReturn.SMX.SvcRs.DETAIL")
	var result []Response.DETAIL
	if len(detail.String()) != 0 {
		if detail.IsArray() {
			if err := json.Unmarshal([]byte(detail.String()), &result); err != nil {
				log.Error("transfer xml Decoder Error", err)
				return smx, result, err
			}
		} else {
			var data Response.DETAIL
			if err := json.Unmarshal([]byte(detail.String()), &data); err != nil {
				log.Error("transfer xml Decoder Error", err)
				return smx, result, err
			}
			result = append(result, data)
		}
	}
	log.Info("transfer xml Decoder", smx, result)
	return smx, result, nil
}

func InvoiceXmlDecoder(xml string, path string) gjson.Result {
	reader := strings.NewReader(xml)
	resp, _ := xj.Convert(reader)
	//log.Debug("ssssss", resp)
	value := gjson.Get(resp.String(), path)
	return value
}

func InvoiceXmlEncoder(v interface{}) ([]byte, error) {
	b := &bytes.Buffer{}
	b.Write([]byte("<?xml version=\"1.0\" encoding=\"utf-8\"?>"))
	enc := xml.NewEncoder(b)
	enc.Indent("", "")
	if err := enc.Encode(v); err != nil {
		log.Error("error: %v\n", err)
		return nil, err
	}
	return b.Bytes(), nil
}

func PostBagXmlEncoder(v interface{}) ([]byte, error) {
	b := &bytes.Buffer{}
	b.Write([]byte("<?xml version=\"1.0\" encoding=\"utf-8\"?>"))
	enc := xml.NewEncoder(b)
	enc.Indent("", "    ")
	if err := enc.Encode(v); err != nil {
		log.Error("error: %v\n", err)
		return nil, err
	}
	return b.Bytes(), nil
}
