package tw_address

import (
	"api/services/Enum"
	"bytes"
	"strings"
	"unicode/utf8"
)

func SeperateCityAndDistrict(addr string) (city, dist string) {
	city = privateFinCity(addr)
	dist = privateFinDistrict(addr, city)
	city = privateCharatarCheck(city)
	return
}

func privateCharatarCheck(old string) (new string) {
	old = strings.ReplaceAll(old, "桃園縣", "桃園市")
	return strings.ReplaceAll(old, "台", "臺")
}

func privateFinCity(input string) (city string) {
	return string([]rune(input)[0:3])
}

func privateFinDistrict(input, city string) (District string) {
	address := strings.TrimPrefix(input, city)
	for _, c := range []rune{'區', '鄉', '鎮', '市'} {
		if index := strings.IndexRune(address, c); index != -1 {
			return address[0:index] + string(c)
		}
	}
	return
}

func AddressSplit(address string) (string, string, string, bool) {

	if country, ok := FindCountry(address); ok {
		subAddress := strings.TrimPrefix(address, country)
		if exchangeName, ok := Enum.CountryExchange[country]; ok {
			country = exchangeName
		}

		if district, ok := FindDistrictByCountry(subAddress, country); ok {

			return country, district, strings.TrimPrefix(subAddress, district), true
		}
		return "", "", "", false

	}
	return "", "", "", false

}
func FindCountry(address string) (string, bool) {

	for _, country := range Enum.CountryList {
		clen := utf8.RuneCountInString(country)
		subAddress := SubStr(address, 0, clen-1)
		if subAddress == country {
			return country, true
		}
	}
	return "", false
}
func FindDistrictByCountry(address string, country string) (string, bool) {

	if _, ok := Enum.CountryCityList[country]; ok {

		for _, district := range Enum.CountryCityList[country] {
			clen := utf8.RuneCountInString(district)
			subAddress := SubStr(address, 0, clen-1)
			if subAddress == district {
				return district, true
			}
		}
		return "", false
	}
	return "", false
}

func SubStr(source interface{}, start int, end int) string {
	str := source.(string)
	var r = []rune(str)
	length := len(r)
	if length == 0 {
		return str
	}
	subLen := end - start

	for {
		if start < 0 {
			break
		}
		if start == 0 && subLen == length {
			break
		}
		if end > length {
			subLen = length - start
		}
		if end < 0 {
			subLen = length - start + end
		}
		var substring bytes.Buffer
		if end > 0 {
			subLen = subLen + 1
		}
		for i := start; i < subLen; i++ {
			substring.WriteString(string(r[i]))
		}
		str = substring.String()

		break
	}

	return str
}
