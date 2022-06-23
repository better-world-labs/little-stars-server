package donation

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/pkg/db"
	"aed-api-server/internal/pkg/utils"
	"fmt"
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
)

func mkMsgContent(name, sex, region, mobile, job, community string) string {
	return fmt.Sprintf(`
### 用户申请积分捐赠项目
- 姓名：%s
- 性别：%s
- 地区：%s
- 手机号：%s
- 职业：%s
- 小区：%s
`, name, sex, region, mobile, job, community)
}

func sendMsg(apply entities.DonationApply, applyId int64) error {
	config := interfaces.GetConfig()
	var str = mkMsgContent(apply.Name, apply.Sex, apply.Region, apply.Mobile, apply.Job, apply.Community)

	msg := utils.DingTalkMsg{
		Msgtype: "actionCard",
		ActionCard: utils.ActionCard{
			Title:          "用户申请积分捐赠项目",
			Text:           str,
			BtnOrientation: "0",
			SingleTitle:    "<环境:" + config.Server.Env + ">点击查看json数据",
			SingleURL:      "https://" + config.Server.Host + "/api/donations/apply/explain?id=" + strconv.FormatInt(applyId, 10),
		},
	}
	return utils.SendDingTalkBot(config.DonationApplyNotify, &msg)
}

type DonationApplyDO struct {
	entities.DonationApply `xorm:"extends"`
	Id                     int64 `json:"id" xorm:"id pk autoincr"`
	UserId                 int64
	CreatedAt              time.Time
}

func (s *Service) Apply(apply entities.DonationApply, userId int64) error {
	var applyDo = DonationApplyDO{
		DonationApply: apply,
		UserId:        userId,
		CreatedAt:     time.Now(),
	}

	_, err := db.Table("donation_apply").Insert(&applyDo)
	if err != nil {
		return err
	}
	utils.Go(func() {
		err := sendMsg(apply, applyDo.Id)
		if err != nil {
			log.Error("sendMsg error", err)
		}
	})
	return nil
}
