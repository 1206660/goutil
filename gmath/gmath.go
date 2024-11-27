package gmath

import (
	"math/rand"
)

func Clamp(val, min, max int) int {
	if val <= min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

func ConvertBoolToInt(val bool) int {
	switch val {
	case false:
		return 0
	case true:
		return 1
	}
	return 0
}

// 为数组添加唯一数据
func AppendDistinct(srcArr []int, tryAddVal int) (bool, []int) {
	for i := 0; i < len(srcArr); i++ {
		if srcArr[i] == tryAddVal {
			return false, srcArr
		}
	}
	return true, append(srcArr, tryAddVal)
}

//生成随机字符串
var __randomg_string_list = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func GetRandomString(length int) string {
	bytes := []byte(__randomg_string_list)
	result := []byte{}
	r := rand.New(rand.NewSource(rand.Int63()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}
