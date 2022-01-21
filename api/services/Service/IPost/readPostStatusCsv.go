package IPost

import (
	"api/services/model"
	"api/services/util/log"
	"api/services/util/tools"
	"encoding/csv"
	"fmt"
	"github.com/jlaffaye/ftp"
	"github.com/spf13/viper"
	"io"
	"os"
	"regexp"
	"time"
)

// 下載
func download(conn *ftp.ServerConn, file string) error {
	res, err := conn.Retr(file)
	defer res.Close()
	if err != nil {
		log.Fatal(`ftp.Retr :`, err)
		return err
	}

	path := tools.GetFilePath(viper.GetString("IPOST.csvPath"), "", 0)
	outFile, err := os.Create(path + file)

	log.Debug("outFile =>", outFile, file)

	defer outFile.Close()
	//
	if err != nil {
		log.Fatal(`ftp.Create :`, err)
		return err
	}

	_, err = io.Copy(outFile, res)
	if err != nil {
		log.Fatal(`ftp.Copy :`, err)
		return err
	}

	return nil
}

// 搬移
func mv(conn *ftp.ServerConn, file string) error {
	err := conn.Rename(file, `success/` + file)
	if err != nil {
		log.Fatal("os.Rename :", err.Error())
		return err
	}

	return nil
}

// dev02 ftp 處理 貨態資料
func dev2IPostFtp() error {
	var filesName []string

	config := viper.GetStringMapString("Ftp")
	c, err := ftp.Dial(config["path"], ftp.DialWithTimeout(5*time.Second), ftp.DialWithDisabledEPSV(true))

	if err != nil {
		log.Fatal("ftp.Dial :", err)
		return err
	}

	err = c.Login(config["username"], config["password"])
	defer c.Quit()

	if err != nil {
		log.Fatal("ftp.Login :", err)
		return err
	}

	current, _ := c.CurrentDir()

	log.Info(`CurrentDir [%v]`, current)
	// 切換所在資料夾
	err = c.ChangeDir("/home/ipost159")
	if err != nil {
		log.Fatal("ftp.ChangeDir :", err)
		return err
	}

	current, _ = c.CurrentDir()
	log.Info(`ChangeDir CurrentDir [%v]`, current)

	nameList, err := c.NameList(".")
	log.Error(`NameList Error`, err)
	log.Info("NameList :", nameList)
	//
	for _, n := range nameList {
		// 年月日-大宗碼-000000-年月日-時分秒.CSV
		// "200822-101277-000000-201021-100120"
		if bl, _ := regexp.MatchString(".*-101277-.*\\.(?i)CSV$", n); bl {
			filesName = append(filesName, n)
			err = download(c, n)
			// 下載失敗寫log
			if err != nil {
				log.Error("ftp.download fail %s: %s", n ,err)
			} else {
				// 下載成功搬移到其他資料夾
				emv := mv(c, n)
				if emv != nil {
					log.Error("ftp.mv fail %s: %s", n ,emv)
				}
			}
		}
	}

	log.Info("dev2IPostFtp regexp.MatchString [%s]", filesName)

	return err
}

// 連線dev02 ftp 下載貨態資料 並 寫入檔案
func HandleIPostCVS()  {
	log.Info("HandleIPostCVS [start]")
	// 連線 dev2 ftp 處理下載 搬移動作
	_ = dev2IPostFtp()
	readPostShippingStatusCsv()
	log.Info("HandleIPostCVS [end]")
}


// 讀寫資料 - 找尋資料夾過濾檔案
func readPostShippingStatusCsv() {

	path := tools.GetFilePath(viper.GetString("IPOST.csvPath"), "", 0)
	_, files, _ := tools.GetDirList(path)

	log.Debug("files", files)

	for _, f := range files {
		boolean, _ := regexp.Match(".*-101277-.*\\.(?i)CSV$", []byte(f))
		log.Info("path + f", path + f)
		if boolean {
			log.Info("path + f", path + f)
			err := toReadAndWriteData(path + f)
			if err != nil {
				continue
			} else {
				err = os.Remove(path + f)
				if err != nil {
					log.Error(`Remove [%v]`,f )
					log.Error("Remove err: [%s]", err)
				}
			}
		}
	}
}

// 讀寫資料 - 開啟檔案寫入DB
func toReadAndWriteData(csvFile string) error {
	file, err := os.Open(csvFile)
	if err != nil {
		log.Error("toReadAndWriteData file Error:", err)
		return err
	}
	defer file.Close()
	reader := csv.NewReader(file)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("toReadAndWriteData reader Error:", err, record)
			return nil
		}
		model.InsertPostShippingStatus(record)
	}

	return nil
}