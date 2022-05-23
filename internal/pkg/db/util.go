package db

import "bytes"

func ParamPlaceHolder(len int) string {
	b := bytes.Buffer{}
	b.WriteString("(")
	for i := 0; i < len; i++ {
		b.WriteString("?")
		if i < len-1 {
			b.WriteString(",")
		}
	}

	b.WriteString(")")
	return b.String()
}
