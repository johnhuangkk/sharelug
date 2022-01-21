package export

import (
	"api/services/util/log"
	"encoding/json"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
	"os"
)

// 寫出 i郵箱JSON給前端頁面用
func JsonToData(posts interface{}, fileName string) error {

	var config = viper.GetStringMapString("Data")

	dir := config["gopath"]
	// 開資料夾
	_ = os.MkdirAll(dir, os.ModePerm)

	source := config["gopath"] + fileName
	destination := config["wwwpath"] + fileName


	if _, err := os.Stat(source); err == nil {
		_ = os.Remove(source)
	}

	if _, err := os.Stat(destination); err == nil {
		_ = os.Remove(destination)
	}

	// 新增|開啟 檔案
	f, err := os.OpenFile(source, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		log.Debug("OpenFile writePostBoxData err", err)
		return err
	}
	defer f.Close()

	data, _ := json.Marshal(posts)
	_, err = f.WriteString(string(data))

	if err != nil {
		log.Debug("WriteString writePostBoxData err", err)
		return err
	}

	_, err = CopyFile(destination, source)

	if err != nil {
		log.Debug("CopyFile  err", err)
		return err
	}

	return nil
}

func JsonToData2(posts interface{}, fileName string) error {

	var config = viper.GetStringMapString("Data")

	dir := config["gopath"]
	// 開資料夾
	_ = os.MkdirAll(dir, os.ModePerm)

	source := dir + fileName
	destination := config["wwwpath"] + fileName

	if _, err := os.Stat(source); err == nil {
		log.Info(`JsonToData2 Remove source`, source)
		err = os.Remove(source)
		if err != nil {
			log.Info(`JsonToData2 Remove source err`, err)
		}
	}

	if _, err := os.Stat(destination); err == nil {
		log.Info(`JsonToData2 Remove destination`, destination)
		err = os.Remove(destination)
		if err != nil {
			log.Info(`JsonToData2 Remove destination err`, err)
		}
	}


	data, _ := json.Marshal(posts)

	err := ioutil.WriteFile(source, data, 0677)

	if err != nil {
		log.Debug("JsonToData2 WriteFile err", err)
		return err
	}
	log.Debug("JsonToData2", destination, source)
	_, err = CopyFile(destination, source)
	if err != nil {
		log.Debug("CopyFile  err", err)
		return err
	}

	return nil
}

// 複製到前端data資料夾下
func CopyFile(destination, source string) (written int64, err error) {
	src, err := os.Open(source)
	if err != nil {
		return
	}
	defer src.Close()

	dst, err := os.OpenFile(destination, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	defer dst.Close()

	return io.Copy(dst, src)
}
