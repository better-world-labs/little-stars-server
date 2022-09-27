package audit

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestAuditTextScan(t *testing.T) {
	a := &Audit{
		"LTAI5tN8fcJYPbrQGetxFL7c",
		"4o40jlAWKPboMFbk7MajzMOBX7fNbj",
		"green.cn-shanghai.aliyuncs.com",
	}

	pass, err := a.ScanText("日你妈，退钱")
	require.Nil(t, err)
	fmt.Printf("%v", pass)

}

func TestAuditImageScan(t *testing.T) {
	a := &Audit{
		"LTAI5tN8fcJYPbrQGetxFL7c",
		"4o40jlAWKPboMFbk7MajzMOBX7fNbj",
		"green.cn-shanghai.aliyuncs.com",
	}

	pass, err := a.ScanImage("https://openview-oss.oss-cn-chengdu.aliyuncs.com/aed-/111/111_1657607308716020608.jpg")
	require.Nil(t, err)
	fmt.Printf("%v", pass)

}

func TestGmt(t *testing.T) {
	gmt := FormatGMT(time.Now())
	fmt.Println(gmt)
}
