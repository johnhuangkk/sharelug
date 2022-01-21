package PdfService

import (
	"api/services/util/log"
	"api/services/util/tools"
	"fmt"
	wkhtml "github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"strings"
	"time"
)

func GeneratorPdfFile(html string) (string, error) {
	pdf, err :=  wkhtml.NewPDFGenerator()
	if err != nil{
		log.Debug("Error", err)
		return "", err
	}
	pdf.AddPage(wkhtml.NewPageReader(strings.NewReader(html)))
	// Create PDF document in internal buffer
	if err := pdf.Create(); err != nil {
		log.Debug("Error", err)
		return "", err
	}
	path := tools.GetFilePath("/temp/", "", 0)
	filename := fmt.Sprintf("%s.%s", time.Now().Format("pdf20060102150405"), "pdf")
	file := fmt.Sprintf("%s%s", path, filename)
	err = pdf.WriteFile(file)
	if err != nil {
		log.Debug("Error", err)
		return "", err
	}
	return filename, nil
}
