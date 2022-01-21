package tools

import (
	"api/services/util/log"
	"bufio"
	"fmt"
	"io"
	"mime/multipart"
	"os"
)

//寫入檔案
func CreateFile(path, fileContent, filename string) (string, error) {
	fileName := path + filename
	err := os.Remove(fileName)
	if err != nil {
		log.Error("delete file error", err)
	}
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0755)
	// 是用 OpenFile,不是只用Open,因為還要設定模式. 建立檔案 只有寫入  UNIX檔案權限
	if err != nil {
		log.Debug("開檔錯誤!")
		return fileName, err
	}
	defer file.Close()
	output := bufio.NewWriter(file)
	_, err = output.WriteString(fileContent)
	if err != nil {
		return "", err
	}
	err = output.Flush()
	if err != nil {
		return "", err
	}
	return fileName, nil
}

func UploadFile(file multipart.File, header *multipart.FileHeader) (string, error) {
	filename := header.Filename
	log.Debug("Content-Type", header.Header.Get("Content-Type"))
	out, err := os.Create("./data/temp/" + filename)
	if err != nil {
		log.Debug("file Error", err)
		return "", err
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		log.Debug("file Error", err)
		return "", err
	}
	return filename, nil
}

func WriteFile(file io.Reader, filename string) error {
	path := GetFilePath("/invoice/E0501/", "", 0)
	outFile, err := os.Create(path + filename)
	if err != nil {
		log.Debug("file Error", err)
		return err
	}
	defer outFile.Close()
	_, err = io.Copy(outFile, file)
	if err != nil {
		log.Debug("file Error", err)
		return err
	}
	return nil
}

func UploadAwardedFile(file multipart.File, header *multipart.FileHeader) (string, error) {
	filename := header.Filename
	log.Debug("Content-Type", header.Header.Get("Content-Type"))
	path := GetFilePath("/Invoice/awarded/src/", "", 0)
	out, err := os.Create(path + filename)
	if err != nil {
		log.Debug("file Error", err)
		return "", err
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		log.Debug("file Error", err)
		return "", err
	}
	return filename, nil
}

func MoveAwardedFile(filename string) error {
	path := GetFilePath("/Invoice/awarded/bak/", "", 0)
	if err := os.Rename(fmt.Sprintf("./data/Invoice/awarded/src/%s", filename), path + filename); err!= nil {
		log.Error("Move file err", err)
		return err
	}
	return nil
}
