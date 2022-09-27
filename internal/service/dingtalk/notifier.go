package dingtalk

import (
	"aed-api-server/internal/interfaces/events"
	"aed-api-server/internal/interfaces/facility"
	"aed-api-server/internal/interfaces/service"
	"aed-api-server/internal/pkg/domain/emitter"
	"aed-api-server/internal/pkg/utils"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
)

var (
	ErrorInvalidEventType = errors.New("event assert failed, invalid event type")
)

type Notifier struct {
	DingTalkUrl string              `conf:"donation-apply-notify"`
	User        service.UserService `inject:"-"`
}

//go:inject-component
func NewNotifier() *Notifier {
	return &Notifier{}
}

func (n Notifier) Listen(on facility.OnEvent) {
	on(&events.DeviceMarkedEvent{}, n.handleDeviceMarked)
	on(&events.ClockInEvent{}, n.handleDeviceClockIns)
}

func (n Notifier) handleDeviceClockIns(e emitter.DomainEvent) error {
	logrus.Info("[dingtalk.Notifier] handleDeviceClockIns")

	if evt, ok := e.(*events.ClockInEvent); ok {
		var imgStr string
		for i := range evt.Images {
			imgStr += fmt.Sprintf("![](%s)\n", evt.Images[i])
		}

		u, exists, err := n.User.GetUserById(evt.CreatedBy)
		if err != nil {
			return err
		}

		if !exists {
			return errors.New("user not found")
		}

		clockInResult := "存在"
		if !evt.IsDeviceExisted {
			clockInResult = "不存在"
		}

		return utils.SendDingTalkBot(n.DingTalkUrl, &utils.DingTalkMsg{
			Msgtype: "markdown",
			Markdown: utils.Markdown{
				Title: "设备打卡",
				Text: fmt.Sprintf(`### 设备打卡
- 用户: %s(%d)[%s]
- 设备ID: %s
- 打卡结果: %s
- 打卡说明: %s
- 打卡图片: 
  %s
                `, u.Nickname, u.ID, u.Uid, evt.DeviceId, clockInResult, evt.Description, imgStr),
			},
		})
	}

	return ErrorInvalidEventType
}

func (n Notifier) handleDeviceMarked(e emitter.DomainEvent) error {
	logrus.Info("[dingtalk.Notifier] handleDeviceMarked")

	if evt, ok := e.(*events.DeviceMarkedEvent); ok {
		u, exists, err := n.User.GetUserById(evt.CreateBy)
		if err != nil {
			return err
		}

		if !exists {
			return errors.New("user not found")
		}

		return utils.SendDingTalkBot(n.DingTalkUrl, &utils.DingTalkMsg{
			Msgtype: "markdown",
			Markdown: utils.Markdown{
				Title: "设备新增",
				Text: fmt.Sprintf(`### 设备新增
- 用户: %s(%d)
- 设备ID: %s
- 设备标题: %s
- 设备地址 %s
- 设备坐标: (%f,%f)
- 设备图;
  ![](%s)
- 环境图:
  ![](%s)
                `, u.Nickname, u.ID, evt.Id, evt.Title, evt.Address, evt.Longitude, evt.Latitude, evt.DeviceImage, evt.EnvironmentImage),
			},
		})
	}

	return ErrorInvalidEventType
}
