package tools

import (
	"fmt"
	"strings"
	"unicode"
)

func validateDigit(s string) bool {
	dotCount := 0 //There are several decimal points in statistics, only one decimal point can appear
	for _, v := range s {
		if v == '.' {
			dotCount++
			if dotCount > 1 { //Only one decimal point is allowed
				return false
			}
			continue
		}
		if !unicode.IsDigit(v) {
			return false
		}
	}
	return true
}

//Function to remove extra characters
func cleanUnneededChar(price string) string {
	//Statistics whether there are extra positive/negative signs in the head, only one positive/negative sign is allowed in the head
	symbolI, symbolJ := 0, 1
	for {
		if symbolJ == len(price)-1 {
			break
		}

		if price[symbolI] == '-' || price[symbolI] == '+' {
			if price[symbolJ] == '-' || price[symbolJ] == '+' {
				symbolI++
				symbolJ++
			} else {
				break
			}
		} else {
			break
		}
	}

	price = price[symbolI:] //Cut off the extra positive/negative signs at the top

	//Check if this string has a plus/minus sign. If there is a symbol, extract the symbol separately first
	symbolString := ""
	if price[0] == '-' || price[0] == '+' {
		symbolString = string(price[0])
		price = price[1:]
	}

	if len(price) == 1 {
		return symbolString + price
	}

	//Count how many extras are in the head 0
	zeroI := 0
	for {
		if price[zeroI] != '0' {
			break
		} else {
			zeroI++
		}
	}

	//If the value is not a floating point number, then the following 0 is not allowed to move
	dotIndex := strings.Index(price, ".")
	if dotIndex == -1 {
		if zeroI > 0 { //zeroI is to count whether there is an extra 0 in the head, if this value is greater than 0, it means that there is an extra 0 in the head
			return symbolString + price[zeroI:] //Cut off the extra 0 at the top, then put the symbol and return
		}

		return symbolString + price //There is no extra 0 at the top, no need to crop
	}

	//The code can go here to show that this must be a floating point number, at least with a decimal point
	//Find out if there are any characters after the decimal point. If there is no character, it means that there is only a solitary decimal point character, then remove this extra decimal point
	if price[dotIndex:] == "." {
		price = price[zeroI:dotIndex]
	}

	//Only when it is a floating point value, will the extra 0 at the end be removed
	end := len(price) //The last subscript, push forward
	for {
		if price[end-1] != '0' {
			if price[end-1] == '.' { //Judge if there is an extra'.' at the end
				end--
			}
			break
		} else {
			end--
		}
	}

	return symbolString + price[zeroI:end]
}

//Core function: from right to left, every 3 digits, add a single-character comma','
func comma(s string) string {
	if len(s) <= 3 {
		return s
	}

	return comma(s[:len(s)-3]) + "," + comma(s[len(s)-3:])
}

//Output a numeric string as financial (add a comma every 3 digits)
/*
 Idea: The way of violent traversal
 1. Clean up the extra positive/negative signs and extra 0s in the string.
 2. Check whether the string has a symbol, if there is, extract the symbol and use it for splicing at the end.
 3. The extra characters have been cleaned, and the verification is performed to determine whether the character string is a legal value. There can be at most one decimal point `.` in the value, and no non-numeric characters can appear in the value (not considering scientific notation and complex type).
 4. If the string is a floating point value, extract all the decimal part and the decimal point and use it for splicing at the end.
 5. To perform core functions, add a single-character comma `,` every 3 digits (code source: "Goc Programming Language" Chapter 3.5.4, P54, gopl.io/ch3/comma).
 6. In the final splicing, in order, the positive/negative sign, the separated character string, and the decimal part are spliced ​​into a final complete character string.
*/
func FormatFinancialString(price string) string {
	//Clean up extra characters. For example: 0 at the end of a floating point number, 0 at the beginning, extra positive/negative signs
	price = cleanUnneededChar(price)

	//Check whether this string has a plus/minus sign. If there is a symbol, extract the symbol separately first
	symbolString := ""
	if price[0] == '-' || price[0] == '+' {
		symbolString = string(price[0])
		price = price[1:]
	}

	//After cleaning the extra characters, start to verify this numeric string
	if !validateDigit(price) {
		return "Illegal value! Please check whether the value you provided is correct! The value is allowed to be a floating point number, and a positive/negative sign is allowed in the front of the number!"
	}

	//If there is no 0 before the decimal point, just add a 0 to it to make the number string look better
	if price[0] == '.' {
		return "0" + price
	}

	// Determine if this number is a floating point value
	dotIndex, decimalString := strings.Index(price, "."), ""
	if dotIndex != -1 {
		decimalString = price[dotIndex:]
		price = price[:dotIndex]
	} else if dotIndex == -1 {
		dotIndex = len(price)
	}

	return fmt.Sprintf("%s%s%s", symbolString, comma(price[:dotIndex]), decimalString)
}

