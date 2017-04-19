package utils

import (
	"bytes"

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
		if len(mac) != len("0123456789AB") || govalidator.IsHexadecimal(mac) == false {
			return false
		}
	}

	return true
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
