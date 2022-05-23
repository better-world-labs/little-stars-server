package device

import (
	"aed-api-server/internal/interfaces/entities"
	service2 "aed-api-server/internal/interfaces/service"
	"aed-api-server/internal/pkg"
	"aed-api-server/internal/pkg/feedback"
	"aed-api-server/internal/pkg/global"
	"aed-api-server/internal/pkg/location"
	page "aed-api-server/internal/pkg/query"
	"aed-api-server/internal/pkg/response"
	"aed-api-server/internal/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.openviewtech.com/openview-pub/gopkg/log"
)

type Controller struct {
	service service2.DeviceService
}

func NewController(s service2.DeviceService) *Controller {
	return &Controller{service: s}
}

// 查询设备列表
func (con *Controller) ListDevices(c *gin.Context) {
	req := new(ListDevice)
	err := c.BindQuery(req)
	if err != nil {
		log.DefaultLogger().Errorf("ListDevices bind error: %v", err)
		response.ReplyError(c, err)
		return
	}

	list, err := con.service.ListDevices(location.Coordinate{Longitude: req.Longitude, Latitude: req.Latitude}, req.Distance, req.Query)
	if err != nil {
		log.DefaultLogger().Errorf("ListDevices error: %v", err)
		response.ReplyError(c,err)
		return
	}

	log.DefaultLogger().Debugf("ListDevices lnt: %v,lat: %v, len:%v", req.Longitude, req.Latitude, len(list))

	response.ReplyOK(c, page.Result{List: list, Total: len(list)})
}

// 标记设备
func (con *Controller) MarkDevice(c *gin.Context) {
	req := new(entities.AddDevice)
	err := c.BindJSON(req)
	if err != nil {
		log.DefaultLogger().Errorf("AddDevice bind error: %v", err)
		response.ReplyError(c, err)
		return
	}
	userId := c.MustGet(pkg.AccountIDKey).(int64)
	pointRst, err := con.service.AddDevice(userId, req)
	if err != nil {
		log.DefaultLogger().Errorf("AddDevice error: %v", err)
		response.ReplyError(c, err)
		return
	}
	back := feedback.NewValuableFeedBack()
	err = back.AddPointsEventRsts(pointRst)
	if err != nil {
		log.DefaultLogger().Errorf("AddDevice error: %v", err)
		response.ReplyError(c, err)
		return
	}

	response.ReplyOK(c, back)
}

// 添加aed设备
func (con *Controller) AddDevice(c *gin.Context) {
	req := new(entities.AddDevice)
	err := c.BindJSON(req)
	if err != nil {
		log.DefaultLogger().Errorf("AddDevice bind error: %v", err)
		response.ReplyError(c, err)
		return
	}
	accountID := c.MustGet(pkg.AccountIDKey).(int64)
	device, err := con.service.AddDevice(accountID, req)
	if err != nil {
		log.DefaultLogger().Errorf("AddDevice error: %v", err)
		response.ReplyError(c, err)
		return
	}

	response.ReplyOK(c, device)
}

// 设备详情
func (con *Controller) InfoDevice(c *gin.Context) {
	req := new(InfoDevice)
	err := c.BindQuery(req)
	if err != nil {
		log.DefaultLogger().Errorf("InfoDevice bind error: %v", err)
		response.ReplyError(c, err)
		return
	}

	info, err := con.service.InfoDevice(req.Longitude, req.Latitude, req.UdId)
	if err != nil {
		log.DefaultLogger().Errorf("InfoDevice error: %v", err)
		response.ReplyError(c, err)
		return
	}

	response.ReplyOK(c, info)
}

func (con *Controller) AddGuide(c *gin.Context) {
	accountID := c.MustGet(pkg.AccountIDKey).(int64)
	req := new(AddDeviceGuideDto)
	err := c.BindJSON(req)
	if err != nil {
		log.DefaultLogger().Errorf("AddGuide bind error: %v", err)
		response.ReplyError(c, err)
		return
	}

	var desc []string
	var remarks []string
	pics := [][]string{}
	for _, v := range req.Info {
		desc = append(desc, v.Desc)
		remarks = append(remarks, v.Remark)
		pics = append(pics, v.Pic)
	}

	pointRst, err := con.service.AddGuideInfo(accountID, req.DeviceId, desc, remarks, pics)
	if err != nil {
		log.DefaultLogger().Errorf("AddGuide error: %v", err)
		response.ReplyError(c, err)
		return
	}
	back := feedback.NewValuableFeedBack()
	err = back.AddPointsEventRsts(pointRst)
	if err != nil {
		log.DefaultLogger().Errorf("AddGuide error: %v", err)
		response.ReplyError(c, err)
		return
	}

	response.ReplyOK(c, back)
}

func (con *Controller) ListGuide(c *gin.Context) {
	req := new(entities.DeviceGuideList)
	err := c.BindQuery(req)
	if err != nil {
		log.DefaultLogger().Errorf("ListGuide bind error: %v", err)
		response.ReplyError(c, err)
		return
	}

	res, err := con.service.GetDeviceGuideInfo(req.DeviceId)
	if err != nil {
		log.DefaultLogger().Errorf("ListGuide error: %v", err)
		response.ReplyError(c, err)
		return
	}

	response.ReplyOK(c, page.Result{List: res.Info, Total: len(res.Info)})
}

func (con *Controller) GetDeviceGuideInfoById(c *gin.Context) {
	req := new(DeviceGuideDto)
	err := c.BindQuery(req)
	if err != nil {
		log.DefaultLogger().Errorf("GetDeviceGuideInfoById bind error: %v", err)
		response.ReplyError(c, err)
		return
	}

	res, err := con.service.GetGuideInfoById(req.Uid)
	if err != nil {
		log.DefaultLogger().Errorf("GetDeviceGuideInfoById error: %v", err)
		response.ReplyError(c, err)
		return
	}

	response.ReplyOK(c, res)
}

func (con *Controller) DeviceGallery(ctx *gin.Context) {
	deviceId := ctx.Param("deviceId")

	latest, exists := ctx.GetQuery("latest")
	l := 0
	if exists {
		la, err := strconv.Atoi(latest)
		l = la
		utils.MustNil(err, global.ErrorInvalidParam)
	}

	gallery, err := con.service.GetDeviceGallery(deviceId, l)
	utils.MustNil(err, global.ErrorInvalidParam)

	response.ReplyOK(ctx, map[string]interface{}{"gallery": gallery})
}

func (con *Controller) CountDevice(ctx *gin.Context) {
	count, err := con.service.CountDeviceByCredibleState()
	utils.MustNil(err, err)

	total := 0
	needPicket := 0
	for _, c := range count {
		total += c.Count
	}

	for _, c := range count {
		if c.CredibleState == CredibleStatusDeviceNotFound {
			needPicket += c.Count
		}
	}

	response.ReplyOK(ctx, map[string]interface{}{
		"total":      total,
		"needPicket": needPicket,
	})
}
