package MartSeven

import (
	"api/services/util/SetupSFtp"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Client struct {
	ApiKey string
	PriKey string
	Host string
	FtpHost string
	FtpUsername string
	FtpPassword string
	client http.Client
}


func (receiver *Client) GetCSTD() (storeData []Store, err error) {
	t := "CSTD"
	data, _, err := receiver.privateDownloadAndUnzip(t)
	if err != nil {
		return nil, err
	}

	// 刪除前綴
	data = data[3:]
	rowData := [][]byte{}
	shift := 0
	for i,v := range data {
		if v == 0x0D && data[i+1] == 0x0A {
			rowData = append(rowData,data[shift:i+1])
			shift = i
		}
	}
	fmt.Println()
	for _,v := range rowData {
		tempRow := v
		if tempRow[0] == 0x0D && tempRow[1] == 0x0A {
			tempRow = tempRow[2:]
		}
		storeId := string(tempRow[:6])
		tempRow = tempRow[6:]
		content := string(tempRow)

		for strings.Contains(content,"   ") {
			content = strings.ReplaceAll(content,"   ","  ")
		}

		infos := strings.Split(content,"  ")


		store := Store{
			Id:      storeId,
			Name:    infos[0],
			Address: infos[2],
			Detail: content,
		}
		fmt.Println("Store:",store)
		storeData = append(storeData,store)

	}

	return storeData,nil
}

func (receiver *Client) OrderAdd(req OrderAddRequest) error {

	path := "/c2c/PaymentBack.ashx"

	reqData := `eshopid=`+ req.EshopId +`&
	eshopsonid=`+ req.EshopSonId +`&
	orderno=`+ req.OrderNo +`&
	service_type=`+ req.ServiceType +`&
	account=`+ req.Account +`&
	payment_cpname=`+ req.PaymentCpName +`&
	trade_description=`+ req.TradeDescription +`&
	cp_remark01=`+ req.CpRemark01 +`&
	cp_remark02=`+ req.CpRemark02 +`&
	cp_remark03=`+ req.CpRemark03 +`&
	deadlinedate=`+ req.DeadlineDate +`&
	deadlinetime=`+ req.DeadlineTime +`&
	daishou_account=`+ req.DeadlineDate +`&
	sender=`+ req.Sender +`&
	sender_phone=`+ req.SenderPhone +`&
	receiver=`+ req.Receiver +`&
	receiver_phone=`+ req.ReceiverPhone +`&
	receiver_storeid=`+ req.ReceiverStoreId +`&
	return_storeid=NNNNNN`

	client := http.Client{}
	buf := bytes.NewBuffer([]byte(reqData))
	resp,err := client.Post(receiver.Host + path,"application/x-www-form-urlencoded",buf)
	if err != nil {
		return err
	}

	data,err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(data))
	return nil
}

func (receiver *Client) GetCEIN() (doc CEIN, fileName string, err error) {
	raw, fName, err := receiver.privateDownloadAndUnzip("CEIN")
	if err != nil {
		return doc, "", err
	}

	err = doc.DecodeXML(raw)
	if err != nil {
		return doc, "", err
	}

	return doc, fName, nil
}

func (receiver* Client) privateDownloadAndUnzip(docType string) (raw []byte,fileName string, err error) {
	client,err := SetupSFtp.NewSFTPClientAndLogin(receiver.FtpHost,receiver.FtpUsername,receiver.FtpPassword)
	if err != nil {
		return
	}

	workF := "/" +docType + "/"

	rawData, fileName, err := client.GetAEarlierFileInDir(workF,"",docType)
	if err != nil {
		return nil,"",err
	}

	client.LogoutAndClose()
	return rawData, fileName, nil
}
