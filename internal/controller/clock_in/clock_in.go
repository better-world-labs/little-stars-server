package clock_in

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/interfaces/service"
	"aed-api-server/internal/pkg"
	"aed-api-server/internal/pkg/feedback"
	"aed-api-server/internal/pkg/location"
	"errors"
	"github.com/gin-gonic/gin"
	"gitlab.openviewtech.com/openview-pub/gopkg/route"
)

type Controller struct {
	ClockIn service.ClockInService `inject:"-"`
}

func NewController() *Controller {
	return &Controller{}
}

func (con Controller) MountAuthRouter(r *route.Router) {
	clockInR := r.Group("/clock-ins")
	clockInR.POST("", con.PostClockIn)
	clockInR.GET("", con.GetDeviceClockInList)
	clockInR.GET("/last", con.GetDeviceLastClockIn)
}

func (con Controller) MountNoAuthRouter(r *route.Router) {
	clockInR := r.Group("/clock-ins")
	clockInR.GET("/stat", con.GetClockInStat)
}

func (con Controller) GetClockInStat(c *gin.Context) (interface{}, error) {
	stat, err := con.ClockIn.GetDeviceClockInStat()
	if err != nil {
		return nil, err
	}

	return stat, nil
}

func (con Controller) PostClockIn(c *gin.Context) (interface{}, error) {
	var clockInBaseInfo entities.ClockInBaseInfo
	if err := c.ShouldBind(&clockInBaseInfo); err != nil {
		return nil, err
	}
	rst, err := con.ClockIn.DoDeviceClockIn(&clockInBaseInfo, c.MustGet(pkg.AccountIDKey).(int64))
	if err != nil {
		return nil, err
	}
	back := feedback.NewValuableFeedBack()
	err = back.AddPointsEventRsts(rst)
	if err != nil {
		return nil, err
	}

	return back, nil
}

func (con Controller) GetDeviceClockInList(c *gin.Context) (interface{}, error) {
	deviceId, _ := c.GetQuery("deviceId")
	if len(deviceId) == 0 {
		return nil, errors.New("deviceId is required")
	}

	list, err := con.ClockIn.GetDeviceClockInList(deviceId)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (con *Controller) GetDeviceLastClockIn(c *gin.Context) (interface{}, error) {
	deviceId, _ := c.GetQuery("deviceId")
	if len(deviceId) == 0 {
		return nil, errors.New("deviceId is required")
	}

	var from location.Coordinate
	err := c.ShouldBindQuery(&from)
	if err != nil {
		return nil, err
	}

	rst, err := con.ClockIn.GetDeviceLastClockIn(from, deviceId)
	if err != nil {
		return nil, err
	}

	return rst, nil
}
