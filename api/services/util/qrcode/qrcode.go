package qrcode

import (
	"api/services/util/images"
	"api/services/util/log"
	"api/services/util/tools"
	"fmt"
	"github.com/spf13/viper"
	"github.com/yeqown/go-qrcode"
	"image"
)

func GeneratorQrCode(url string, productId string, img image.Image) error {
	var qrc *qrcode.QRCode
	var err error
	if img != nil {
		srcImg, _ := images.ResizeIcon(img)
		oo := qrcode.WithLogoImage(srcImg)
		qrc, err = qrcode.New(url, oo)

	} else {
		qrc, err = qrcode.New(url)
	}
	if err != nil {
		log.Error("could not generate QRCode: %v", err)
		return err
	}
	// save file
	path, _ := tools.GetImageFilePath(viper.GetString("qrcode.path"), "", 0)
	var saveToPath = fmt.Sprintf("%s%s.jpg", path, productId)
	log.Debug("Save image:", saveToPath)
	if err := qrc.Save(saveToPath); err != nil {
		log.Error("Save Qrcode image: %v", err)
		return err
	}
	return nil
}

func GetQrcodeImageLink(productId string) string {
	return fmt.Sprintf("/static/images/qrcode/%s.jpg", productId)
}

func GetTinyUrl(TinyUrl string) string {
	var tinyUrl = viper.GetString("TinyURL")
	return fmt.Sprintf("%s/%s", tinyUrl, TinyUrl)
}
