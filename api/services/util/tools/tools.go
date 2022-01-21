package tools

import (
	"api/services/util/log"
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"math/rand"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const couponCodeCharset = "abcdefghijkmnpqrstuvwxyz"

//產生亂數字申
func RandString(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

//產生亂數數字串
func RangeNumber(max int, length int) string {
	rand.Seed(time.Now().UnixNano())
	randNum := rand.Intn(max)
	num := strconv.Itoa(randNum)
	return fmt.Sprintf("%0*s", length, num)
}

//左邊補0
func StringPadLeft(input string, length int) string {
	return fmt.Sprintf("%0*s", length, input)
}

//目前時間
func Now(format string) string {
	switch format {
	case "YmdHis":
		return time.Now().Format("2006-01-02 15:04:05")
	case "Ymd":
		return time.Now().Format("2006-01-02")
	case "Hms":
		return time.Now().Format("15:04:05")
	case "TwDate":
		year, _ := strconv.Atoi(time.Now().Format("2006"))
		s := StringPadLeft(strconv.Itoa(year-1911), 4)
		return fmt.Sprintf("%s%s", s, time.Now().Format("0102"))
	case "TwYear":
		year, _ := strconv.Atoi(time.Now().Format("2006"))
		s := StringPadLeft(strconv.Itoa(year-1911), 3)
		return fmt.Sprintf("%s", s)
	default:
		return time.Now().String()
	}
}

func GetInvoiceYearMonth() string {
	month := time.Now().Format("01")
	switch month {
	case "01", "02":
		return fmt.Sprintf("%s%s", Now("TwYear"), "02")
	case "03", "04":
		return fmt.Sprintf("%s%s", Now("TwYear"), "04")
	case "05", "06":
		return fmt.Sprintf("%s%s", Now("TwYear"), "06")
	case "07", "08":
		return fmt.Sprintf("%s%s", Now("TwYear"), "08")
	case "09", "10":
		return fmt.Sprintf("%s%s", Now("TwYear"), "10")
	case "11", "12":
		return fmt.Sprintf("%s%s", Now("TwYear"), "12")
	}
	return ""
}

func GetYearMonth(month string) []string {
	switch month {
	case "02":
		return []string{"1", "2"}
	case "04":
		return []string{"3", "4"}
	case "06":
		return []string{"5", "6"}
	case "08":
		return []string{"7", "8"}
	case "10":
		return []string{"9", "10"}
	case "12":
		return []string{"11", "12"}
	}
	return []string{}
}

func NowHHmmss() string {
	return Now("Hms")
}

func NowYYYYMMDDHHmmss(afterDay int) string {
	now := time.Now().Add(time.Hour * 24 * time.Duration(afterDay))
	return now.Format("2006-01-02 15:04:05")
}

// 取得資料夾列表
func GetDirList(path string) ([]string, []string, error) {
	var dir, files []string

	if len(path) == 0 {
		path = "."
	}

	fs, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal("GetDirList", err)
		return files, nil, err
	}

	for _, f := range fs {
		name := f.Name()
		if f.IsDir() {
			dir = append(dir, name)
		} else {
			files = append(files, name)
		}
	}

	return dir, files, nil
}

// 寫檔案
func WriteFileByByte(filePath string, fileName string, context []byte) {
	var dest = filePath + fileName

	err := os.MkdirAll(filePath, 0755)
	if err != nil {
		log.Fatal(`WriteFile Mkdir Error`, dest)
		log.Fatal(`WriteFile Mkdir Error`, err.Error())
	}
	outFile, err := os.Create(dest)
	if err != nil {
		log.Fatal(`WriteFile Create Error`, err.Error())
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, bytes.NewReader(context))
	if err != nil {
		log.Fatal(`WriteFile Copy Error`, err.Error())
	}
}

// php in Array
func InArray(needArr []string, need string) bool {
	for _, v := range needArr {
		if need == v {
			return true
		}
	}
	return false
}

//字串轉數字
func StringToInt64(str string) int64 {
	int, _ := strconv.Atoi(str)
	return int64(int)
}

//字串轉數字
func ParseInt(str string) int {
	int, _ := strconv.Atoi(str)
	return int
}

func IntToString(i int) string {
	return strconv.Itoa(i)
}

// 檢查是否值為0
func CheckIsZero(param int, def int) int {
	if param == 0 {
		return def
	}
	return param
}

func Sign(i float64) string {
	if i >= 0 {
		return "+"
	} else {
		return "-"
	}
}

/**
轉換 字串陣 -> []interface
*/
func StringArrayToInterface(strAry []string) []interface{} {
	var i []interface{}
	for _, s := range strAry {
		i = append(i, s)
	}
	return i
}

//
func XmlToString(v interface{}) string {
	l, _ := xml.Marshal(v)
	return string(l)
}

func Trim(s string) string {
	return strings.Replace(s, " ", "", -1)
}

func Nl2br(s string) string {
	s1 := strings.Replace(s, " ", "", -1)
	return strings.Replace(s1, "\n", "", -1)
}

func JsonDecode(data []byte, inf interface{}) error {
	log.Debug("Decode Json Source", string(data))
	if err := json.Unmarshal(data, &inf); err != nil {
		log.Error("Decode Json Error", err)
		return err
	}
	return nil
}

func JsonEncode(vo interface{}) (string, error) {
	result, err := json.Marshal(vo)
	if err != nil {
		log.Error("Encode Json Error", err)
		return "", err
	}
	return string(result), nil
}

//四捨五入
func Round(x float64) float64 {
	return math.Floor(x + 0.5)
}

// 是否為正式環境
func EnvIsProduction() bool {
	return viper.GetString(`ENV`) == `prod`
}

func GetWithdrawType(bankCode string) string {
	if InArray([]string{"0040037", "0050418", "0060567", "0070937", "0081005", "0095185", "0110026", "0122009", "0130017", "0172015", "0480011", "0500108", "0540537", "1030019", "8030021", "8060219", "8070014", "8090267", "8100364", "8120012", "8150015", "8220901"}, bankCode) {
		return "eACH"
	} else {
		return "ACH"
	}
}

func IsOrderId(orderId string) bool {
	match, err := regexp.MatchString("^[A-Z]\\d{12}$", orderId)
	if err != nil {
		log.Error("Check Order Id Error", err)
	}
	return match
}

func IsAlphaNumeric(content string) bool {
	match, err := regexp.MatchString("^.[A-Za-z0-9]+$", content)
	if err != nil {
		log.Error("Is Alpha Numeric Error", err)
	}
	return match
}

//解析URL
func UrlParse(rawurl string) *url.URL {
	u, err := url.Parse(rawurl)
	if err != nil {
		log.Error("Url Parse Error")
	}
	return u
}

//信用卡卡號解析
func ParseCredit(params string) (string, []string) {
	old := Trim(params)
	card := strings.Split(old, "")
	var number []string
	number = append(number, strings.Join(card[:4], ""))
	number = append(number, strings.Join(card[4:8], ""))
	number = append(number, strings.Join(card[8:12], ""))
	number = append(number, strings.Join(card[12:], ""))
	return old, number
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
func RandLowerString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = couponCodeCharset[rand.Intn(len(couponCodeCharset))]
	}
	return string(b)
}

func ChangeCityName(s string) string {
	switch s {
		case "臺北市":
			return "台北市"
		case "臺中市":
			return "台中市"
		case "臺南市":
			return "台南市"
		case "臺東市":
			return "台東市"
		default:
			return s
	}
}