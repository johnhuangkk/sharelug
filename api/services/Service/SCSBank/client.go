package SCSBank

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Client struct {
	host string
	client http.Client
}

func NewClient(host string) Client {
	return Client{host: host, client: http.Client{}}
}

// 取得帳務紀錄
func (me *Client) GetQueryAccountingData() (records []IncomeRecord,result error) {
	resp,err := me.client.Get(me.host + "/scsbstore?Func=query&Otype=5")
	if err != nil {
		return nil,err
	}

	data ,err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil,err
	}

	contentRaw := string(data)
	contents := strings.Split(contentRaw,"\n")
	for _,v := range contents {
		if v == "" {
			continue
		}
		records = append(records,NewIncomeRecord(v))
	}

	return records,nil
}

// 更新帳務紀錄的索引，更新過的索引資料就不再查詢到。
func (me *Client) GetAckAccountingData(baccno string, btxDate string, seqno string) (result bool, err error) {
	query := fmt.Sprintf("?Func=%s&Baccno=%s&Btxdate=%s&Seqno=%s","Ack",baccno,btxDate,seqno)
	resp,err := me.client.Get( me.host + "/scsbstore" + query)
	if err == nil {
		data ,err := ioutil.ReadAll(resp.Body)
		if err != nil {
			contentRaw := string(data)
			result = strings.Contains(contentRaw,"0")
			return true,nil
		}
	}
	return false, err
}


