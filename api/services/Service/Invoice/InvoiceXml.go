package Invoice

import (
	"api/services/Service/Invoice/InvoiceXml"
	"api/services/util/log"
	"api/services/util/tools"
	"api/services/util/xml"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)


func BranchTrackDecoder(str string) (InvoiceXml.BranchTrack, error) {
	value := xml.InvoiceXmlDecoder(str, "BranchTrack")
	var data InvoiceXml.BranchTrack
	if len(value.String()) != 0 {
		err := json.Unmarshal([]byte(value.String()), &data)
		if err != nil {
			return data, err
		}
	}
	return data, nil
}

func GenerateCancelAllowanceXml(data InvoiceXml.CancelAllowance) error {
	data.Xmlns = "urn:GEINV:eInvoiceMessage:D0501:3.1"
	data.Xsi = "http://www.w3.org/2001/XMLSchema-instance"
	data.SchemaLocation = "urn:GEINV:eInvoiceMessage:D0501:3.1 D0501.xsd"
	content, err := xml.InvoiceXmlEncoder(data)
	if err != nil {
		log.Error("Generate Invoice Error", err)
		return err
	}
	path := tools.GetFilePath("/invoice/allowance/cancel/", "", 0)
	filename := fmt.Sprintf("%s.xml", time.Now().Format("20060102150405"))
	file, err := tools.CreateFile(path, string(content), filename)
	if err != nil {
		log.Debug("Create File Error", err)
		return err
	}
	if err := Connect().UploadFolder("D0501", file, filename); err != nil {
		log.Error("Upload Ftp Error", err)
		return err
	}
	return nil
}

func GenerateAllowanceXml(data InvoiceXml.Allowance) error {
	data.Xmlns = "urn:GEINV:eInvoiceMessage:D0401:3.1"
	data.Xsi = "http://www.w3.org/2001/XMLSchema-instance"
	data.SchemaLocation = "urn:GEINV:eInvoiceMessage:D0401:3.1 D0401.xsd"
	content, err := xml.InvoiceXmlEncoder(data)
	if err != nil {
		log.Error("Generate Invoice Error", err)
		return err
	}
	path := tools.GetFilePath("/invoice/allowance/", "", 0)
	filename := fmt.Sprintf("%s.xml", data.Main.AllowanceNumber)
	file, err := tools.CreateFile(path, string(content), filename)
	if err != nil {
		log.Debug("Create File Error", err)
		return err
	}
	if err := Connect().UploadFolder("D0401", file, filename); err != nil {
		log.Error("Upload Ftp Error", err)
		return err
	}
	return nil
}

func GenerateVoidInvoiceXml(data InvoiceXml.VoidInvoiceC0701) error {
	data.Xmlns = "urn:GEINV:eInvoiceMessage:C0701:3.1"
	data.Xsi = "http://www.w3.org/2001/XMLSchema-instance"
	data.SchemaLocation = "urn:GEINV:eInvoiceMessage:C0701:3.1 C0701.xsd"
	//產生Xml Schema
	content, err := xml.InvoiceXmlEncoder(data)
	if err != nil {
		log.Error("Generate Invoice Error", err)
		return err
	}
	path := tools.GetFilePath("/invoice/void/", "", 0)
	filename := fmt.Sprintf("%s.xml", data.VoidInvoiceNumber)
	file, err := tools.CreateFile(path, string(content), filename)
	if err != nil {
		log.Debug("Create File Error", err)
		return err
	}
	if err := Connect().UploadFolder("C0701", file, filename); err != nil {
		log.Error("Upload Ftp Error", err)
		return err
	}
	return nil
}

func GenerateCancelInvoiceXml(data InvoiceXml.CancelInvoiceC0501) error {
	data.Xmlns = "urn:GEINV:eInvoiceMessage:C0501:3.1"
	data.Xsi = "http://www.w3.org/2001/XMLSchema-instance"
	data.SchemaLocation = "urn:GEINV:eInvoiceMessage:C0501:3.1 C0501.xsd"
	//產生Xml Schema
	content, err := xml.InvoiceXmlEncoder(data)
	if err != nil {
		log.Error("Generate Invoice Error", err)
		return err
	}
	path := tools.GetFilePath("/invoice/cancel/", "", 0)
	filename := fmt.Sprintf("%s.xml", data.CancelInvoiceNumber)
	file, err := tools.CreateFile(path, string(content), filename)
	if err != nil {
		log.Debug("Create File Error", err)
		return err
	}
	if err := Connect().UploadFolder("C0501", file, filename); err != nil {
		log.Error("Upload Ftp Error", err)
		return err
	}
	return nil
}

func GenerateInvoiceXml(data InvoiceXml.InvoiceC0401) error {
	data.Xmlns = "urn:GEINV:eInvoiceMessage:C0401:3.1"
	data.Xsi = "http://www.w3.org/2001/XMLSchema-instance"
	data.SchemaLocation = "urn:GEINV:eInvoiceMessage:C0401:3.1 C0401.xsd"
	//產生Xml Schema
	content, err := xml.InvoiceXmlEncoder(data)
	if err != nil {
		log.Error("Generate Invoice Error", err)
		return err
	}
	path := tools.GetFilePath("/invoice/invoice/", "", 0)
	filename := fmt.Sprintf("%s.xml", data.Main.InvoiceNumber)
	file, err := tools.CreateFile(path, string(content), filename)
	if err != nil {
		log.Debug("Create File Error", err)
		return err
	}
	if err := Connect().UploadFolder("C0401", file, filename); err != nil {
		log.Error("Upload Ftp Error", err)
		return err
	}
	return nil
}




//取得發票字軌
func GetInvoiceAssignNumberFile(path string) ([]os.FileInfo, error) {
	if err := Connect().DownLoadFolder(); err != nil {
		log.Error("Upload Ftp Error", err)
		return nil, err
	}
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Error("Open File Error", err)
		return nil, err
	}
	return files, nil
}