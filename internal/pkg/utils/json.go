package utils

import "encoding/json"

func Json(x interface{}) string {
	marshal, _ := json.Marshal(x)
	return string(marshal)
}
