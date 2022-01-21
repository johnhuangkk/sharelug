package Task

import (
	"api/services/model"
	"api/services/util/log"
	"os"
	"time"

	"github.com/jlaffaye/ftp"
	"github.com/spf13/viper"
)

// 上傳檔案FTP
func UploadKgiSpecialStoreRecord() {
	log.Info(`UploadKgiSpecialStoreRecord Start`)
	filename, path, ids, err := model.GetMemberSpecialStoreRecordExcelWithIds()
	c, err := getClient()
	if err != nil {
		log.Error("KgiReport Ftp read error", err.Error())
		return
	}
	defer c.Quit()
	file, err := os.Open(path + filename)

	if err != nil {
		log.Error(err.Error(), "Kgi Special Store FTP File Error")
	}
	defer file.Close()
	err = c.Stor(filename, file)
	if err != nil {
		log.Error("Kgi Special Store FTP FTP Upload Fail", err)
		return
	}

	err = model.UpdateMemberSpecialStoreRecordIsSendWithFilename(ids, filename)
	if err != nil {
		log.Error("Kgi Special Store record update fail", err)
		return
	}
	log.Info(`UploadKgiSpecialStoreRecord End`)
}

func getClient() (*ftp.ServerConn, error) {
	var host string
	ENV := viper.GetString("ENV")
	if ENV != "prod" {
		host = viper.GetString(`KgiReport.localhost`)
	} else {
		host = viper.GetString(`KgiReport.host`)
	}
	c, err := ftp.Dial(host, ftp.DialWithTimeout(5*time.Second))
	if err != nil {

		return c, err
	}
	err = c.Login(viper.GetString((`KgiReport.username`)), viper.GetString(`KgiReport.password`))
	if err != nil {

		return c, err
	}
	return c, nil
}
