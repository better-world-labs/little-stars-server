package user

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/pkg/config"
	"github.com/stretchr/testify/assert"
	"gitlab.openviewtech.com/openview-pub/gopkg/log"
	"testing"
)

func Test_genImgCommand(t *testing.T) {
	c, err := config.LoadConfig("../../../config-local.yaml")
	if err != nil {
		panic("get config error")
	}
	interfaces.InitConfig(c)
	log.Init(c.Log)

	command, err := genImgCommand(
		"dapeng",
		"https://thirdwx.qlogo.cn/mmopen/vi_32/HWxbOaAUg2RWXR34LjZAHdL2efLAFKtfePZtRrKZjJmfgJuCwKs4l6pT9xOdJxXuvV2ma41gegsWrmTXtWPRjw/132",
		[]string{
			"https://openview-oss.oss-cn-chengdu.aliyuncs.com/star-static/%E5%AE%8C%E6%88%90%E6%8D%90%E7%8C%AE_small.png",
			"https://openview-oss.oss-cn-chengdu.aliyuncs.com/star-static/%E8%8D%A3%E8%AA%89%E5%B7%A1%E6%A3%80%E5%91%98.png",
		},
		10,
		2,
		20,
		"https://openview-oss.oss-cn-chengdu.aliyuncs.com/star-static/%E8%8D%A3%E8%AA%89%E5%B7%A1%E6%A3%80%E5%91%98.png",
		"card-117.jpeg",
	)

	assert.Nil(t, err)
	log.Info(command)
}
