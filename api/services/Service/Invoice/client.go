package Invoice

import (
	"api/services/util/log"
	"api/services/util/tools"
	"github.com/jlaffaye/ftp"
	"github.com/spf13/viper"
	"os"
	"time"
)

type ConnectFtp struct {
	conn ftp.ServerConn
}

func Connect() *ConnectFtp {
	config := viper.GetStringMapString("INVOICE")
	conn, err := ftp.Dial(config["host"], ftp.DialWithTimeout(5*time.Second), ftp.DialWithDisabledEPSV(true))
	if err != nil {
		log.Error("ftp.Dial conn:", err)
	}
	if err := conn.Login(config["account"], config["password"]); err != nil {
		log.Error("ftp.Login :", err)
	}
	return &ConnectFtp {
		conn: *conn,
	}
}
//上傳檔案
func (conn *ConnectFtp) UploadFolder(folderName, file, filename string) error {
	defer conn.conn.Quit()
	//LIST 目錄下的檔案
	//current, _ := conn.conn.CurrentDir()
	//log.Debug("current", current)
	if viper.GetString("ENV") == "prod" {
		//	變更目錄
		if err := conn.conn.ChangeDir("/home/invoice/TurnkeyData/UpCast/B2CSTORAGE/" + folderName + "/SRC"); err != nil {
			log.Error("ftp Change Dir :", err)
			return err
		}
		open, err := os.Open(file)
		if err != nil {
			log.Error("Open file Error", err)
		}
		defer open.Close()
		if err := conn.conn.Stor(filename, open); err != nil {
			log.Error("Update file Error", err)
		}
	}
	return nil
}
//下載檔案
func (conn *ConnectFtp) DownLoadFolder() error {
	defer conn.conn.Quit()
	//LIST 目錄下的檔案
	folder1, err := conn.conn.List("/home/invoice/TurnkeyData/Unpack/BAK/E0501/")
	if err != nil {
		return err
	}
	for _, v1 := range folder1 {
		folder2, err := conn.conn.List("/home/invoice/TurnkeyData/Unpack/BAK/E0501/" + v1.Name + "/")
		if err != nil {
			return err
		}
		for _, v2 := range folder2 {
			folder3, err := conn.conn.List("/home/invoice/TurnkeyData/Unpack/BAK/E0501/" + v1.Name + "/" + v2.Name + "/")
			if err != nil {
				return err
			}
			for _, v3 := range folder3 {
				r, err := conn.conn.Retr("/home/invoice/TurnkeyData/Unpack/BAK/E0501/" + v1.Name + "/" + v2.Name + "/" + v3.Name)
				if err != nil {
					log.Error("Read Files Error", err)
					return err
				}
				if err := tools.WriteFile(r, v3.Name);err != nil {
					log.Error("Read Files Error", err)
					return err
				}
				r.Close()
			}
		}
	}
	return nil
}