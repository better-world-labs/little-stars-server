package image_process_test

import (
	"gitlab.openviewtech.com/openview-pub/gopkg/log"
	"testing"
)

func TestCall(s *testing.T) {
	log.Init(log.LogConfig{
		Level:        5,
		Output:       "stdout",
		ReportCaller: false,
	})
	//TODO  依赖第三方服务 暂时注视掉
	//service := image_process.NewImageBotService("http://localhost:8888/call")
	//url, err := service.Call("user-charity-card", map[string]interface{}{
	//	"username":        "username",
	//	"userAvatar":      "https://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTLtZvbOWWopA7216libKCVabh9EcLLmh3UWYZAIQ6XMaxibIZpicRPB7lyibJ5d2zLricwS4wuYfgMfPCA/132",
	//	"medals":          []string{"https://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTLtZvbOWWopA7216libKCVabh9EcLLmh3UWYZAIQ6XMaxibIZpicRPB7lyibJ5d2zLricwS4wuYfgMfPCA/132"},
	//	"donationPoints":  "5000",
	//	"donationProject": "1",
	//	"addStarDays":     "21",
	//	"qrContent":       "qrContent",
	//}, "hahaha")
	//if err != nil {
	//	panic(err)
	//}
	//
	//fmt.Println(url)
}
