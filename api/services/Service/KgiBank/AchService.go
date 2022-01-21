package KgiBank

import (
	"api/services/Enum"
	"api/services/Service/MemberService"
	"api/services/VO/KgiBank"
	"api/services/dao/Withdraw"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"api/services/util/tools"
	"fmt"
	"github.com/spf13/viper"
	"strconv"
	"time"
)

func ExporterKgiEachFile() (string, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	data, err := Withdraw.GetEachWithdrawDataByStatus(engine, Enum.WithdrawStatusWait)
	if err != nil {
		log.Error("Get Withdraw Error", err)
		return "", err
	}
	BatchId := fmt.Sprintf("EACH%s", time.Now().Format("20060102"))
	if len(data) == 0 {
		return "", nil
	}
	var config = viper.GetStringMapString("KgiAch")
	content := SetAchHeader()
	count := 0
	amt := 0
	for _, v := range data {
		identityId := MemberService.TakeMemberIdentity(engine, v.UserId)
		if len(identityId) == 0 {
			log.Error("Get identity Error", err)
			return "", err
		}
		content += SetAchBody(v, identityId, v.TransId, config)
		count ++
		amt += int(v.WithdrawAmt)
		err = Withdraw.UpdateWithdrawDataByStatus(engine, v, BatchId)
		if err != nil {
			log.Error("Update Withdraw Error", err)
			return "", err
		}
	}
	content += SetAchFooter(amt, count)
	path := tools.GetFilePath(config["path"], "", 0)
	filename, err := tools.CreateFile(path, content, BatchId + ".txt")
	if err != nil {
		log.Error("Create File Error", err)
		return "", err
	}
	return filename, nil
}

//產生提領ACH檔案
func ExporterKgiAchFile() (string, error) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	data, err := Withdraw.GetAchWithdrawDataByStatus(engine, Enum.WithdrawStatusWait)
	if err != nil {
		log.Error("Get Withdraw Error", err)
		return "", err
	}
	BatchId := fmt.Sprintf("ACH%s", time.Now().Format("20060102"))
	if len(data) == 0 {
		log.Error("Not Withdraw Data")
		return "", nil
	}
	var config = viper.GetStringMapString("KgiAch")
	content := SetAchHeader()
	count := 0
	amt := 0
	for _, v := range data {
		identityId := MemberService.TakeMemberIdentity(engine, v.UserId)
		if len(identityId) == 0 {
			log.Error("Get identity Error", err)
			return "", err
		}
		content += SetAchBody(v, identityId, v.TransId, config)
		count ++
		amt += int(v.WithdrawAmt)
		err = Withdraw.UpdateWithdrawDataByStatus(engine, v, BatchId)
		if err != nil {
			log.Error("Update Withdraw Error", err)
			return "", err
		}
	}
	content += SetAchFooter(amt, count)
	path := tools.GetFilePath(config["path"], "", 0)
	filename, err := tools.CreateFile(path, content, BatchId + ".txt")
	if err != nil {
		log.Error("Create File Error", err)
		return "", err
	}
	return filename, nil
}

//產生ACH表頭
func SetAchHeader() string {
	var header = KgiBank.AchHeader{
		BOF:    "BOF",
		CDATA:  "ACHP01",
		TDATE:  tools.Now("TwDate"),
		TTIME:  time.Now().Format("150405"),
		SORG:   "8090072",
		RORG:   "9990250",
		VERNO:  "V10",
		FILLER: "",
	}
	return NewKgiAchHeaderRule().ToString(header)
}

//產生ACH內容
func SetAchBody(data entity.WithdrawData, identity string, Seq string, config map[string]string) string {
	body := KgiBank.AchBody{
		TYPE: "N",
		TXTYPE: "SC",
		TXID: "405",
		SEQ: Seq,
		PBANK: config["bankcode"],
		PCLNO: config["bankaccount"],
		RBANK: data.BankCode,
		RCLNO: data.BankAccount,
		AMT: strconv.Itoa(int(data.WithdrawAmt)),
		RCODE: "",
		SCHD: "B",
		CID: config["senderid"],
		PID: identity,
		SID: "",
		PDATE: "",
		PSEQ: "",
		PSCHD: "",
		CNO: "",
		NOTE: "",
		MEMO: "CheckNe",
		CFEE: "",
		NOTEB: "",
		FILLER: "",
	}
	return NewKgiAchBodyRule().ToString(body)
}

//產生ACH表尾
func SetAchFooter(amt, count int) string {
	footer := KgiBank.AchFooter{
		EOF: "EOF",
		CDATA: "ACHP01",
		TDATE:  tools.Now("TwDate"),
		SORG: "8090072",
		RORG: "9990250",
		TCOUNT: strconv.Itoa(count),
		TAMT: strconv.Itoa(amt),
		YDATE: "",
		FILLER: "",
	}
	return NewKgiAchFooterRule().ToString(footer)
}
