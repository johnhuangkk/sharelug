package main

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	//"os"
	"regexp"
)


type Param struct {
	SysClientId string	`json:"sysClientId"`
	ExpNum		string	`json:"expNum"`
	ClientPhone string	`json:"ClientPhone"`
	ReturnUrl 	string	`json:"Return_Url"`
}

type IPostBox struct {
	ADMId			string 						`json:"ADM_Id"`
	ADMAddress		string 						`json:"ADM_Address"`
	ADMName			string 						`json:"ADM_Name"`
	ADMLocation		string 						`json:"ADM_Location"`
	ADMZip			string 						`json:"ADM_Zip"`
	Longitude		string 						`json:"Longitude"`
	Latitude		string 						`json:"Latitude"`
	Deep			string 						`json:"Deep"`
	Width			string 						`json:"Width"`
	High			string 						`json:"High"`
	Param			Param					`json:"param"`
}


type IPostBoxes []IPostBox


type IPOSTData struct {
	ADMId         	string		`json:"ADMId"` 		// i郵箱櫃體編號
	ADMName        	string		`json:"ADMName"`    // i郵箱櫃體名稱
	ADMAlias      	string		`json:"ADMAlias"`    //
	ADMLocation  	string		`json:"ADMLocation"`    // i郵箱櫃體所在位置描述
	Country 		string		`json:"country"`    //
	Zip 			string		`json:"zip"`		//
	City 			string		`json:"city"`    //
	Address 		string		`json:"address"`    // i郵箱櫃體地址
	Longitude 		string		`json:"Longitude"`    //
	Latitude 		string		`json:"Latitude"`    //
	POSTGOV_No 		string		`json:"POSTGOV_No"`    //
}

type IPost []IPOSTData

func regReplace(reg string, str string, rpl string) string {
	re := regexp.MustCompile(reg)
	return re.ReplaceAllString(str, rpl)
}

func ttt()  {

	//var urlencode = `%5B%7B%22ADM_Id%22:%22324%22,%22ADM_Address%22:%22%E9%AB%98%E9%9B%84%E5%B8%82%E4%B8%89%E6%B0%91%E5%8D%80%E5%BB%BA%E5%9C%8B%E4%BA%8C%E8%B7%AF318%E8%99%9FB1%22,%22ADM_Name%22:%22%E9%AB%98%E9%9B%84%E7%81%AB%E8%BB%8A%E7%AB%99%EF%BD%89%E9%83%B5%E7%AE%B1%22,%22ADM_Location%22:%22%E9%AB%98%E9%9B%84%E5%B8%82%E4%B8%89%E6%B0%91%E5%8D%80%E5%BB%BA%E5%9C%8B%E4%BA%8C%E8%B7%AF318%E8%99%9FB1%22,%22Longitude%22:%22120.301731%22,%22Latitude%22:%2222.638367%22,%22Deep%22:%22470%22,%22Width%22:%22200%22,%22High%22:%22480%22,%22ADM_zip%22:%22807%22,%22param%22:%22%7B%5C%22sysClientId%5C%22:%5C%223%5C%22,%5C%22expNum%5C%22:%5C%224%5C%22,%5C%22ClientPhone%5C%22:%5C%225%5C%22,%5C%22Return_Url%5C%22:%5C%22https://local.api.sharelug.com/v1/iPostBox/map%5C%22%7D%22%7D%5D`;

	//m, _ := url.QueryUnescape(urlencode)

	//fmt.Println(`urldecode`)
	//fmt.Println(m)

	//m2 := regReplace(`\\`, m, ``)
	//m2 = regReplace(`}"`, m2, `}`)
	//m2 = regReplace(`"{`, m2, `{`)
	//fmt.Println(`m2`)
	//fmt.Println(m2)

	//var re = regexp.MustCompile(`(?m),"param".*"}"`)
	//var str = `[{"ADM_Id":"414","ADM_Address":"屏東縣枋寮鄉枋寮村儲運路18號","ADM_Name":"枋寮火車站ｉ郵箱","ADM_Location":"屏東縣枋寮鄉枋寮村儲運路18號","Longitude":"120.595101","Latitude":"22.368142","Deep":"470","Width":"200","High":"480","ADM_zip":"940","param":"{"sysClientId":"4","expNum":"5","ClientPhone":"6","Return_Url":"https://local.api.sharelug.com/v1/iPostBox/map"}"}]`
	//
	//regxData := re.ReplaceAllString(str, "")

	//regxData := regReplace(`}"`, str, `}`)
	//regxData = regReplace(`"{`, regxData, `{`)


	//var iPostBoxes IPostBoxes
	//
	//json.Unmarshal([]byte(m2), &iPostBoxes)
	//fmt.Printf("%T", iPostBoxes)
	//fmt.Println(iPostBoxes)
	//jsonData, _ := json.Marshal(iPostBoxes)
	//fmt.Println(string(jsonData))
}

func IPxx() {
	iPost := IPost{}


	client := &http.Client{}
	postValues := url.Values{}

	resp, err := client.PostForm("https://ipickup.post.gov.tw/Api/ITRI/Query_POSTBOX_Address", postValues)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	if resp.StatusCode == 200 {
		body, _ := ioutil.ReadAll(resp.Body)

		_ = json.Unmarshal(body, &iPost)
		//fmt.Println(iPost)
		i := 0
		//IJSON := IJSON{}
		var ctys []City
		ii  := IPostJson{}
		for _, v := range iPost {
			i++
			fmt.Println(v)
			//不等於同區 新增區域
			//fmt.Println(len(cty.Name))
			c := false
			for _, v1 := range ctys {
				if v1.Name == v.City {
					c = true
					continue
				}
			}
			if c == false {
				//ii.City = append(ii.City, cty)
				cty := City{}
				cty.Name = v.City
				cty.Zip = v.Zip
				//fmt.Println(cty)
				ctys = append(ctys, cty)
			}

			fmt.Println(ctys)
			//
			//adr := Address{}
			//adr.Name = v.Address
			//adr.Location = v.ADMLocation
			//cty.Address = append(cty.Address, adr)
			////fmt.Println(cty)
			////不等於同縣市 新增縣市
			//if ii.Country != v.Country {
			//	IJSON = append(IJSON, ii)
			//	ii  := IPostJson{}
			//	ii.Country = v.Country
			//}

		}

		fmt.Println(ii)

		fmt.Println(i)
		//
		////fmt.Printf("%T", iPost)
		//fmt.Println(iPost)
		//jsonData, _ := json.Marshal(IJSON)
		////jsonData, _ := json.Marshal(iPost)
		//fmt.Println(string(jsonData))
		//fmt.Println(string(ii))
		//each(jsonData)

		//write(jsonData)
		//fmt.Printf("%T", (jsonData))
	}
}

type IJSON []IPostJson



type IPostJson struct{
	Country string `json:"country"`
	City []City
}

type City struct {
	Name string `json:"city"`
	Zip string `json:"zip"`
	Address []Address
}

type Address struct {
	Name string `json:"address"`
	Location string `json:"Location"`
}

func write(iPost []byte)  {

	//os.MkdirAll("dir1/dir2/dir3", os.ModePerm)   //创建多级目录
	//
	//f, _ := os.OpenFile("./dir1/dir2/dir3/4.json", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	//defer f.Close()
	//f.WriteString(`{`)

	var ii  []IPostJson
	_ = json.Unmarshal(iPost, &ii)

	fmt.Print(ii)

	//for _, value := range iPost {
	//	a, _ := json.Marshal(value)
	//	//_, _ = f.Write(a)
	//}
	//f.WriteString(`}`)
}

func hash384()  {
	str := `330450000001900023209201811190930500900000001090000000220190830`
	timestamp := `20190903171345`
	key := `iBoxEC.UP_MAILINFO`

	w := md5.New()
	io.WriteString(w, str)   //将str写入到w中
	md5str2 := fmt.Sprintf("%x", w.Sum(nil))  //w.Sum(nil)将w的hash转成[]byte格式


	message := timestamp + `,` + md5str2
	fmt.Println(message)

	mac := hmac.New(sha512.New384, []byte(key))
	io.WriteString(mac, message)
	c := fmt.Sprintf("%x", mac.Sum(nil))
	b64 := base64.StdEncoding.EncodeToString([]byte(c))
	fmt.Println("aaaa",b64)

}

func main() {
	IPxx()

}



