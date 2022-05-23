package utils

import (
	"fmt"
	"strconv"
)

func ToInt(str string) int {
	n, _ := strconv.Atoi(str)
	return int(n)
}

func ToFloat(str string) float64 {
	v, _ := strconv.ParseFloat(str, 10)
	return v
}

func ToStr(n interface{}) string {
	return fmt.Sprintf("%v", n)
}
