package utils

import (
	"strconv"
)

func PointsString(points int) string {
	b := []byte(strconv.Itoa(points))
	var res []byte

	separation := []byte(",")
	count := 1
	for i := len(b) - 1; i >= 0; i-- {
		res = append([]byte{b[i]}, res...)
		if (count)%3 == 0 && i > 0 {
			res = append(separation, res...)
		}
		count++
	}

	return string(res)
}
