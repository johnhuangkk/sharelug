package TransferService

import (
	"api/services/Enum"
	"api/services/entity"
	"github.com/spf13/viper"
	"time"
)

func GenerateTransferXml(Account string, StartDate, EndDate, Temp, Type string) (entity.SMX, error) {
	header := generateHeader(Type)
	body := generateBody(Account, StartDate, EndDate, Temp, Type)
	smx := entity.SMX{}
	smx.Header = header
	smx.SvcRq = body
	return smx, nil
}

func GenerateTransferAccountQueryXml(Account string, StartDate, EndDate, Type string) (entity.SMX, error) {
	header := generateHeader(Type)
	body := generateAccountQueryBody(Account, StartDate, EndDate, Type)
	smx := entity.SMX{}
	smx.Header = header
	smx.SvcRq = body
	return smx, nil
}


func generateHeader(Type string) entity.Header {
	t := time.Now()
	header := entity.Header{}
	header.Password = viper.GetString("KgiBank.C2C.Password")
	header.SenderID = viper.GetString("KgiBank.C2C.SenderID")
	if Type == Enum.OrderTransB2c {
		header.Password = viper.GetString("KgiBank.B2C.Password")
		header.SenderID = viper.GetString("KgiBank.B2C.SenderID")
	}
	header.ReceID = viper.GetString("KgiBank.BankCode")
	header.Date = t.Format("20060102")
	header.Time = t.Format("150405")
	header.TxnId = "V522"
	return header
}

func generateBody(Account string, StartDate, EndDate, Temp string, Type string) entity.SvcRq {
	body := entity.SvcRq{}
	body.IDNO = viper.GetString("KgiBank.SenderID")
	body.INQOPTNO = viper.GetString("KgiBank.C2C.AtmCode")
	body.ACNO = viper.GetString("KgiBank.C2C.Account")
	if Type == Enum.OrderTransB2c {
		body.INQOPTNO = viper.GetString("KgiBank.B2C.AtmCode")
		body.ACNO = viper.GetString("KgiBank.B2C.Account")
	}
	body.BDATE = ""
	body.EDATE = ""
	body.SBDATE = StartDate
	body.SEDATE = EndDate
	body.SBTIME = "0000"
	body.SETIME = "2359"
	body.VACNO = Account
	body.TEMPD = Temp
	return body
}

func generateAccountQueryBody(Account string, StartDate, EndDate string, Type string) entity.SvcRq {
	body := entity.SvcRq{}
	body.IDNO = viper.GetString("KgiBank.SenderID")
	body.INQOPTNO = viper.GetString("KgiBank.C2C.AtmCode")
	body.ACNO = viper.GetString("KgiBank.C2C.Account")
	if Type == Enum.OrderTransB2c {
		body.INQOPTNO = viper.GetString("KgiBank.B2C.AtmCode")
		body.ACNO = viper.GetString("KgiBank.B2C.Account")
	}
	body.BDATE = StartDate
	body.EDATE = EndDate
	body.SBDATE = ""
	body.SEDATE = ""
	body.SBTIME = ""
	body.SETIME = ""
	body.VACNO = Account
	body.TEMPD = ""
	return body
}