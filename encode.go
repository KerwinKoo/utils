package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"log"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"

	"io/ioutil"

	"strings"

	"github.com/dchest/uniuri"
)

const encodeStd = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"

// pwdChars is a set of standard characters allowed in uniuri string.
var pwdChars = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789.+-_")

var coder = base64.NewEncoding(encodeStd)

// Base64Decode decode base64-code to src result string
func Base64Decode(base64Str string) string {
	data, err := coder.DecodeString(base64Str)
	if err != nil {
		panic(err)
	}

	result := string(data)
	return result
}

// Base64Encode encode src to BASE64
func Base64Encode(src string) string {
	result := coder.EncodeToString([]byte(src))

	return result
}

//Base64HmacSHA1Encode get BASE64(HmacSHA1(source, key)) encode result
func Base64HmacSHA1Encode(source, key string) string {
	src := []byte(source)
	//hmac_sha1
	mac := hmac.New(sha1.New, []byte(key))
	mac.Write(src)
	// code64 := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	code64 := coder.EncodeToString(mac.Sum(nil))

	return code64
}

// DESEncode 3DES encode check ok
// the password length must be 16 or 24
// return result have been base64ed
func DESEncode(key, src string) string {
	var result string
	if len(key) != 16 && len(key) != 24 {
		panic("key length error, must be 16 or 24")
	}

	tripleDESKey := make([]byte, 0, 24)
	if len(key) == 16 {
		tripleDESKey = append(tripleDESKey, key[:16]...)
		tripleDESKey = append(tripleDESKey, key[:8]...)
	} else {
		tripleDESKey = append(tripleDESKey, key[:]...)
	}

	td, err := des.NewTripleDESCipher(tripleDESKey)
	if err != nil {
		panic(err)
	}

	mod := len(src) % td.BlockSize()
	v := td.BlockSize() - mod

	data := []byte(src)
	for i := 0; i < v; i++ {
		data = append(data, byte(v))
	}

	n := len(data) / td.BlockSize()
	var rb []byte
	for i := 0; i < n; i++ {
		dst := make([]byte, td.BlockSize())
		td.Encrypt(dst, data[i*8:(i+1)*8])
		rb = append(rb, dst[:]...)
	}

	result = coder.EncodeToString(rb) // BASE64 endcoded the 3DES result
	return result
}

// DESDecode 3DES decode check ok
// the password length must be 16 or 24
// result what returned is clear text
// return uncode&nil result or "" and error
func DESDecode(key, src string) (string, error) {
	if len(key) != 16 && len(key) != 24 {
		panic("key length error, must be 16 or 24")
	}

	data, err := coder.DecodeString(src)
	if err != nil {
		return "", err
	}

	tripleDESKey := make([]byte, 0, 24)
	if len(key) == 16 {
		tripleDESKey = append(tripleDESKey, key[:16]...)
		tripleDESKey = append(tripleDESKey, key[:8]...)
	} else {
		tripleDESKey = append(tripleDESKey, key[:]...)
	}

	td, err := des.NewTripleDESCipher(tripleDESKey)
	if err != nil {
		return "", err
	}

	n := len(data) / td.BlockSize()
	var rb []byte
	for i := 0; i < n; i++ {
		dst := make([]byte, td.BlockSize())
		td.Decrypt(dst, data[i*8:(i+1)*8])
		rb = append(rb, dst[:]...)
	}

	lastValue := int(rb[len(rb)-1])
	if lastValue >= len(rb) {
		err := errors.New("3DES encode string format err")
		return "", err
	}
	return string(rb[0 : len(rb)-lastValue]), nil
}

// New3DESPassword create 24 length 3DES password and return it
func New3DESPassword() string {
	return uniuri.NewLenChars(24, pwdChars)
}

// NewHmacSHA1Password create 32 length SHA1-password and return it
func NewHmacSHA1Password() string {
	return uniuri.NewLenChars(32, pwdChars)
}

// NewConnectVerifyPassword create 16 length client connect password
func NewConnectVerifyPassword() string {
	return uniuri.NewLenChars(24, pwdChars)
}

// MD5_32Lower MD5 32bit encode lower letter
func MD5_32Lower(m string) string {
	h := md5.New()
	h.Write([]byte(m))
	cipherStr := h.Sum(nil)
	result := hex.EncodeToString(cipherStr)
	return result
}

// MD5_32Upper MD5 32bit encode upper letter
func MD5_32Upper(m string) string {
	h := md5.New()
	h.Write([]byte(m))
	cipherStr := h.Sum(nil)
	result := hex.EncodeToString(cipherStr)
	upperStr := strings.ToUpper(result)
	return upperStr
}

// AESCBCEncodeWithIV using pkcs7paddin and CBCmode
// this Crypt method using IV to realve comprex platform
func AESCBCEncodeWithIV(plantText, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	log.Println("block length:", block.BlockSize())
	plantText = PKCS7Padding(plantText, block.BlockSize())
	blockModel := cipher.NewCBCEncrypter(block, iv)
	ciphertext := make([]byte, len(plantText))
	blockModel.CryptBlocks(ciphertext, plantText)
	return ciphertext, nil
}

// PKCS7Padding PKCS#7 padding
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

//GBK2UTF8 trans from GBK to UTF-8
func GBK2UTF8(src []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(src), simplifiedchinese.GBK.NewDecoder())
	result, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// UTF82GBK trans from UTF-8 to GBK code
func UTF82GBK(src []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(src), simplifiedchinese.GBK.NewEncoder())
	d, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return d, nil
}

/*
//PKCS5Padding ECB PKCS5Padding
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

//PKCS5Unpadding ECB PKCS5Unpadding
func PKCS5Unpadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

//TripleEcbDesEncrypt [golang ECB 3DES Encrypt]
func TripleEcbDesEncrypt(origData, key []byte) ([]byte, error) {
	tkey := make([]byte, 24, 24)
	copy(tkey, key)
	k1 := tkey[:8]
	k2 := tkey[8:16]
	k3 := tkey[16:]

	block, err := des.NewCipher(k1)
	if err != nil {
		return nil, err
	}
	bs := block.BlockSize()
	origData = PKCS5Padding(origData, bs)

	buf1, err := encrypt(origData, k1)
	if err != nil {
		return nil, err
	}
	buf2, err := decrypt(buf1, k2)
	if err != nil {
		return nil, err
	}
	out, err := encrypt(buf2, k3)
	if err != nil {
		return nil, err
	}
	return out, nil
}

//TripleEcbDesDecrypt [golang ECB 3DES Decrypt]
func TripleEcbDesDecrypt(crypted, key []byte) ([]byte, error) {
	tkey := make([]byte, 24, 24)
	copy(tkey, key)
	k1 := tkey[:8]
	k2 := tkey[8:16]
	k3 := tkey[16:]
	buf1, err := decrypt(crypted, k3)
	if err != nil {
		return nil, err
	}
	buf2, err := encrypt(buf1, k2)
	if err != nil {
		return nil, err
	}
	out, err := decrypt(buf2, k1)
	if err != nil {
		return nil, err
	}
	out = PKCS5Unpadding(out)
	return out, nil
}

// TripleCBCDesEncrypt 3DES加密
func TripleCBCDesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, err
	}
	origData = PKCS5Padding(origData, block.BlockSize())
	// origData = ZeroPadding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key[:8])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

// TripleCBCDesDecrypt 3DES解密
func TripleCBCDesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, key[:8])
	origData := make([]byte, len(crypted))
	// origData := crypted
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5Unpadding(origData)
	// origData = ZeroUnPadding(origData)
	return origData, nil
}
*/
/* private funcs */

/*
//encrypt Des encrypt
func encrypt(origData, key []byte) ([]byte, error) {
	if len(origData) < 1 || len(key) < 1 {
		return nil, errors.New("wrong data or key")
	}
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	bs := block.BlockSize()
	if len(origData)%bs != 0 {
		return nil, errors.New("wrong padding")
	}
	out := make([]byte, len(origData))
	dst := out
	for len(origData) > 0 {
		block.Encrypt(dst, origData[:bs])
		origData = origData[bs:]
		dst = dst[bs:]
	}
	return out, nil
}

//decrypt Des decrypt
func decrypt(crypted, key []byte) ([]byte, error) {
	if len(crypted) < 1 || len(key) < 1 {
		return nil, errors.New("wrong data or key")
	}
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	out := make([]byte, len(crypted))
	dst := out
	bs := block.BlockSize()
	if len(crypted)%bs != 0 {
		return nil, errors.New("wrong crypted size")
	}

	for len(crypted) > 0 {
		block.Decrypt(dst, crypted[:bs])
		crypted = crypted[bs:]
		dst = dst[bs:]
	}

	return out, nil
}
*/
