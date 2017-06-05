package utils

import (
	"bytes"
	"errors"

	"strings"

	govalidator "gopkg.in/asaskevich/govalidator.v4"
)

// IsPhoneNumber check value is phone number or not
func IsPhoneNumber(value string) bool {
	phlen := len(value)
	tmp := "12"
	head := tmp[0]

	if phlen < 11 {
		return false
	} else if value[0] != head || (govalidator.IsNumeric(value) == false) {
		return false
	}

	return true
}

// IsMacAddress check mac is valid or not
// Only allowed format:
// 		01:23:45:67:89:AB
//		01.23.45.67.89.AB
// 		0123456789AB
func IsMacAddress(mac string) bool {
	if govalidator.IsMAC(mac) == true {
		macLenMax := len("01:23:45:67:89:ab")
		if len(mac) > macLenMax {
			return false
		}

		dotFieldLen := 0
		for i := 0; i < len(mac); i++ {
			if mac[i] == '-' {
				return false
			}

			if mac[i] == '.' {
				if dotFieldLen > 2 {
					return false
				}
				dotFieldLen = 0
			}
			dotFieldLen++
		}
	} else {
		if len(mac) != len("0123456789AB") {
			return false
		}

		if govalidator.IsHexadecimal(mac) == false {
			return false
		}
	}

	return true
}

// IsBase64 return true if str is base64 encoded
func IsBase64(str string) bool {
	return govalidator.IsBase64(str)
}

// ToHexadecimalMac to Hexadecimal Mac address with upper character
func ToHexadecimalMac(macStr string) (string, error) {
	macRet := ""
	var err error
	if IsMacAddress(macStr) == false {
		err = errors.New("mac source is invalid")

		return macRet, err
	}

	var macBytes []byte
	for i := 0; i < len(macStr); i++ {
		if macStr[i] == '.' || macStr[i] == ':' {
			continue
		} else {
			macBytes = append(macBytes, byte(macStr[i]))
		}
	}

	macRet = string(macBytes)
	macRet = strings.ToUpper(macRet)
	return macRet, nil
}

// GetKeyValueSign create KuaiChon sign by source
// Md5-digitl sign (Lower case to upper case) format:
// eg:
// Md5("AAA=1&BBB=xxxx&CCC=xxxx&DDD=xxxxxxxxxxxx")
// return :
//		clear text and MD5-32bit upper result
func GetKeyValueSign(signKeySort []string, signData map[string]string) (string, string) {
	var buffer bytes.Buffer
	i := 0

	for i = 0; i < len(signKeySort); i++ {
		for key, value := range signData {

			if key == signKeySort[i] {
				if i != 0 {
					buffer.WriteString("&")
				}
				buffer.WriteString(key)
				buffer.WriteString("=")
				buffer.WriteString(value)
				break
			}
		}
	}

	resultMing := buffer.String()
	resultMd5 := MD5_32Upper(resultMing)

	return resultMing, resultMd5
}

// GetMd5Sign create KuaiChon sign by source
// Md5-digitl sign (Lower case to upper case) format:
// eg:
// Md5("value1value2....")
// return :
//		clear text and MD5-32bit lower result
func GetMd5Sign(signKeySort []string, signData map[string]string) (string, string) {
	var buffer bytes.Buffer
	i := 0

	for i = 0; i < len(signKeySort); i++ {
		for key, value := range signData {

			if key == signKeySort[i] {
				buffer.WriteString(value)
				break
			}
		}
	}

	resultMing := buffer.String()
	resultMd5 := MD5_32Lower(resultMing)

	return resultMing, resultMd5
}
