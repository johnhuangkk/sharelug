package tools

import (
	"api/services/util/log"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"strings"
	"time"
)

func GetFilePath(category string, subFolderPath string, subLength int) string {

	var subFolder []string
	var basePath = viper.GetString("FILE.UPLOAD_FOLDER")

	if subLength != 0 {
		if subFolderPath == "" {
			subFolderPath = time.Now().Format("0601")
		}
		for i := 0; i < len(subFolderPath); i = i + subLength {
			subFolder = append(subFolder, subFolderPath[i: i + subLength])
		}
	}
	dir := strings.Split(basePath, "/")
	cat := strings.Split(category, "/")
	dir = append(dir, cat...)
	dir = append(dir, subFolder...)
	dir = deleteEmpty(dir)
	fp := strings.Join(dir, "/")

	_, err := os.Stat(fp)
	if err != nil {
		err := os.MkdirAll(fp, os.ModePerm)
		if err != nil {
			log.Error("mkdir Error", err)
		}
	}
	return fmt.Sprintf("%s/", fp)
}

func GetImageFilePath(category string, subFolderPath string, subLength int) (string, string) {

	var subFolder []string
	var basePath = viper.GetString("FILE.UPLOAD_IMAGES_FOLDER")

	if subLength != 0 {
		if subFolderPath == "" {
			subFolderPath = time.Now().Format("0601")
		}
		for i := 0; i < len(subFolderPath); i = i + subLength {
			subFolder = append(subFolder, subFolderPath[i: i + subLength])
		}
	}
	dir := strings.Split(basePath, "/")
	cat := strings.Split(category, "/")
	dir = append(dir, cat...)
	dir = append(dir, subFolder...)
	dir = deleteEmpty(dir)
	fp := strings.Join(dir, "/")
	dir1 := delete(dir)
	path := strings.Join(dir1, "/")

	_, err := os.Stat(fp)
	if err != nil {
		err := os.MkdirAll(fp, os.ModePerm)
		if err != nil {
			log.Error("mkdir Error", err)
		}
	}
	return fmt.Sprintf("%s/", fp), fmt.Sprintf("%s/", path)
}


func deleteEmpty (s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

func delete (s []string) []string {
	var r []string
	for _, str := range s {
		if str != "." && str != "www" {
			r = append(r, str)
		}
	}
	return r
}