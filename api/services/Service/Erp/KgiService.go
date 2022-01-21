package Erp

import (
	"api/services/VO/Response"
	"api/services/util/log"
	"io"
	"os"
	"time"

	"github.com/jlaffaye/ftp"
	"github.com/spf13/viper"
)

func GetFolderList() ([]Response.KgiSpecialStoreFile, error) {
	var datas []Response.KgiSpecialStoreFile
	c, err := getClient()
	if err != nil {
		log.Error("KgiReport Ftp read error", err.Error())
		return datas, err
	}
	defer c.Quit()

	files, err := c.List("./")
	if err != nil {
		log.Error("KgiReport Ftp read error", err.Error())
		return datas, err
	}
	for _, file := range files {
		if file.Type == ftp.EntryTypeFile {
			data := Response.KgiSpecialStoreFile{
				Name:    file.Name,
				Created: file.Time.Format("2006-01-02 15:04:05"),
			}
			datas = append(datas, data)
		}

	}
	return datas, nil

}
func GetSpecialStoreFile(filename string) (*os.File, error) {
	var outFile *os.File
	c, err := getClient()
	if err != nil {
		log.Error("KgiReport Ftp read error", err.Error())
		return outFile, err
	}
	defer c.Quit()
	res, err := c.Retr(filename)
	if err != nil {
		log.Error("KgiReport Ftp read error", err.Error())
		return outFile, err
	}
	defer res.Close()

	outFile, err = os.Create(filename)
	if err != nil {
		log.Error("KgiReport Ftp read error", err.Error())
		return outFile, err
	}

	defer outFile.Close()

	_, err = io.Copy(outFile, res)
	if err != nil {
		log.Error("KgiReport Ftp read error", err.Error())
		return outFile, err

	}
	return outFile, nil

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
