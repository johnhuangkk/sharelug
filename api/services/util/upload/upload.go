package upload

import (
	"api/services/util/images"
	"api/services/util/log"
	"api/services/util/tools"
	"fmt"
	"github.com/spf13/viper"
	"golang.org/x/exp/rand"
	"strings"
	"time"
)

func init() {
	rand.Seed(uint64(time.Now().UnixNano()))
}

func ProductImage(base64img string) (string, error) {
	base64img = strings.Replace(base64img, " ", "", -1)
	// 去除换行符
	//base64img = strings.Replace(base64img, "\n", "", -1)
	imgData, err := images.FormatBase64Images(base64img)
	if err != nil {
		log.Error("format images error", err)
		return "", err
	}
	extension := getExtension(imgData)
	path := viper.GetString("images.productPath")
	saveToPath, filename := getNewFilePath(extension, path)

	data, err := images.Resize(imgData)
	if err != nil {
		log.Error("Resize Error", err)
		return "", err
	}

	err = images.CreateImageFile(data, saveToPath)
	if err != nil {
		log.Error("Create Image File Error", err)
		return "", err
	}
	return filename, nil
}

func getExtension(image images.Image) string {
	switch image.MediaType {
	case "image/png":
		return "jpg"
	case "image/webp":
		return "webp"
	default:
		return "jpg"
	}

}

func getNewFilePath(extension string, path string) (string, string) {
	var filename = time.Now().Format("20060102150405") +"_"+ tools.RandString(16) + "." + extension
	imagePath, _ := tools.GetImageFilePath(path, "", 0)
	return fmt.Sprintf("%s%s", imagePath, filename), filename
}

func StorePicture(base64img string) (string, error) {
	base64img = strings.Replace(base64img, " ", "", -1)
	imgData, err := images.FormatBase64Images(base64img)
	if err != nil {
		log.Error("format images error", err)
		return "", err
	}
	extension := getExtension(imgData)
	path := viper.GetString("images.storePath")
	saveToPath, filename := getNewFilePath(extension, path)

	data, err := images.Resize(imgData)
	if err != nil {
		log.Error("Resize Error", err)
		return "", err
	}
	err = images.CreateImageFile(data, saveToPath)
	if err != nil {
		log.Error("Create Image File Error", err)
		return "", err
	}
	return filename, nil
}

func UserPicture(base64img string) (string, error) {
	base64img = strings.Replace(base64img, " ", "", -1)
	imgData, err := images.FormatBase64Images(base64img)
	if err != nil {
		log.Error("format images error", err)
		return "", err
	}
	extension := getExtension(imgData)
	path := viper.GetString("images.userPath")
	saveToPath, filename := getNewFilePath(extension, path)

	data, err := images.Resize(imgData)
	if err != nil {
		log.Error("Resize Error", err)
		return "", err
	}
	err = images.CreateImageFile(data, saveToPath)
	if err != nil {
		log.Error("Create Image File Error", err)
		return "", err
	}
	return filename, nil
}