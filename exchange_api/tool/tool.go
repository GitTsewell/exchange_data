package tool

import (
	"bytes"
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// AES 加密
func AesEncrypt(plaintext string) (string, error) {
	key := []byte("abcdefgh12345678")
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}
	cipher.NewCFBEncrypter(block, iv).XORKeyStream(ciphertext[aes.BlockSize:],
		[]byte(plaintext))
	return hex.EncodeToString(ciphertext), nil
}

// AES 解密
func AesDecrypt(d string) (string, error) {
	key := []byte("abcdefgh12345678")
	ciphertext, err := hex.DecodeString(d)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	cipher.NewCFBDecrypter(block, iv).XORKeyStream(ciphertext, ciphertext)
	return string(ciphertext), nil
}

// interface -> float64
func ArrInterfaceToFloat64(arrString [][]interface{}) [][]float64 {

	var arrFloat [][]float64
	for _,i := range arrString {

		var tmpArr []float64
		for k,j := range i{
			if k <= 1 {
				st := j.(string)
				tmp,_ := strconv.ParseFloat(st,64) // string 转 float64
				tmpArr = append(tmpArr,tmp)
			}
		}
		arrFloat = append(arrFloat,tmpArr)
	}

	return  arrFloat
}

// 解压火币WS数据
func GzipDecodeHuobi(in []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(in))
	if err != nil {
		var out []byte
		return out, err
	}
	defer reader.Close()

	return ioutil.ReadAll(reader)
}

func MsToTime(ms string) (time.Time, error) {
	msInt, err := strconv.ParseInt(ms, 10, 64)
	if err != nil {
		return time.Time{}, err
	}

	tm := time.Unix(0, msInt*int64(time.Millisecond))

	return tm, nil
}

func IsSpotOfOkex(symbol string) bool {
	pattern := "\\d+" //反斜杠要转义
	result,_ := regexp.MatchString(pattern,symbol)
	return !result
}

func IsSpotOfHuobi(symbol string) bool {
	return !strings.Contains(symbol,"_")
}

func IsSpotOfBitmex(symbol string) bool {
	return false
}

func IsSpotOfBinance(symbol string) bool {
	return true
}


func GetCurrentPath() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}
