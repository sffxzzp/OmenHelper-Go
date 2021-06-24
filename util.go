package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func hex2int(hexStr string) uint64 {
	n, err := strconv.ParseUint(hexStr, 16, 64)
	if err != nil {
		panic(err)
	}
	return n
}

func uuid2bytes(uuidStr string) []byte {
	uuidStr = strings.ReplaceAll(uuidStr, "-", "")
	uuidLength := len(uuidStr) / 2
	uuidBytes := make([]byte, uuidLength)
	for i := 0; i < uuidLength; i++ {
		uuidBytes[i] = uint8(hex2int(uuidStr[i*2 : i*2+2]))
	}
	return uuidBytes
}

func genSign(message []byte, secret []byte) string {
	h := hmac.New(sha256.New, secret)
	h.Write(message)
	sha := h.Sum(nil)
	return base64.StdEncoding.EncodeToString([]byte(sha))
}

func getSign(applicationId string, sessionId string, eventName string, strTime []string) string {
	arr := uuid2bytes(applicationId)
	arr2 := uuid2bytes(sessionId)
	bLen := 16
	arr3 := make([]byte, bLen)
	for i := 0; i < bLen; i++ {
		if i < 8 {
			arr3[i] = arr[i*2+1]
		} else {
			arr3[i] = arr2[(i-8)*2]
		}
	}
	var text bytes.Buffer
	text.WriteString(fmt.Sprintf("%s%s%s%d", eventName, strTime[0], strTime[1], 1))
	sign := genSign(text.Bytes(), arr3)
	return sign
}

func formatTime(time time.Time) string {
	strTime := time.Format("2006-01-02T15:04:05.999999Z")
	strTime = strTime[:len(strTime)-4] + "000Z"
	return strTime
}

func genTime(minutes int) []string {
	cTime := time.Now().UTC()
	start := cTime.Add(time.Duration(int64(-minutes) * int64(time.Minute)))
	end := cTime
	return []string{formatTime(start), formatTime(end)}
}
