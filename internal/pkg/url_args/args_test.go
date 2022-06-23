package url_args

import (
	"github.com/sirupsen/logrus"
	"testing"
)

func Test_getLinkArgs(t *testing.T) {
	url := "/subcontract/article/detail?url=https%3A%2F%2Fmp.weixin.qq.com%2Fs%3F__biz%3DMzkxNjMyNDE4OA%3D%3D%26mid%3D2247484232%26idx%3D1%26sn%3Dc2daa014f118324ef6caec97018d2daa%26chksm%3Dc150e97bf627606de3a87c01f9c0506e4f86c24c56146608350e18c8b827c5f09c8f35df63ad%23rd"

	args := GetArgs(url)

	logrus.Infof("%v", args)
}
