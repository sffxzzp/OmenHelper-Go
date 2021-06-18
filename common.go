package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
)

func Input(str string) string {
	var inputStr string
	fmt.Print(str)
	fmt.Scan(&inputStr)
	return inputStr
}

func Read(path string) []byte {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil
	}
	return data
}

func Write(path, data string) {
	err := ioutil.WriteFile(path, []byte(data), 0777)
	if err != nil {
		panic(err)
	}
}

func Int2Str(i int) string {
	return strconv.Itoa(i)
}

func Str2Int(s string) int {
	num, _ := strconv.Atoi(s)
	return num
}
