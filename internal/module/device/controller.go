package device

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/interfaces/entities"
	service2 "aed-api-server/internal/interfaces/service"
	"aed-api-server/internal/pkg"
	"aed-api-server/internal/pkg/feedback"
	"aed-api-server/internal/pkg/location"
	page "aed-api-server/internal/pkg/query"
	"aed-api-server/internal/pkg/response"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gitlab.openviewtech.com/openview-pub/gopkg/route"
)

type Controller struct {
	Service service2.DeviceService `inject:"-"`
}

func NewController() *Controller {
	return &Controller{}
}

func (con *Controller) MountAuthRouter(r *route.Router) {
	r.GET("/devices/:deviceId/gallery", con.DeviceGallery)
	r.GET("/aed/devices", con.ListDevices)
	r.POST("/devices", con.MarkDevice)
	r.POST("/aed/add/device", con.AddDevice)
	r.GET("/aed/device", con.InfoDevice)
	r.POST("/aed/device/add_guide", con.AddGuide)
	r.GET("/aed/device/guide", con.ListGuide)
	r.GET("/aed/device/guide_info", con.GetDeviceGuideInfoById)
}

func (con *Controller) MountAdminRouter(r *route.Router) {
	r.POST("/import-devices", con.ImportDevices)
	r.POST("/sync-devices", con.SyncDevices)
}

// 查询设备列表
func (con *Controller) ListDevices(c *gin.Context) (interface{}, error) {
	req := new(ListDevice)
	err := c.ShouldBindQuery(req)
	if err != nil {
		log.Errorf("ListDevices bind error: %v", err)
		return nil, err
	}

	list, err := con.Service.ListDevices(location.Coordinate{Longitude: req.Longitude, Latitude: req.Latitude}, req.Distance, req.Query)
	if err != nil {
		log.Errorf("ListDevices error: %v", err)
		return nil, err
	}

	log.Debugf("ListDevices lnt: %v,lat: %v, len:%v", req.Longitude, req.Latitude, len(list))

	return page.Result{List: list, Total: len(list)}, nil
}

// 标记设备
func (con *Controller) MarkDevice(c *gin.Context) (interface{}, error) {
	req := new(entities.AddDevice)
	err := c.ShouldBindJSON(req)
	if err != nil {
		return nil, err
	}
	userId := c.MustGet(pkg.AccountIDKey).(int64)
	pointRst, err := con.Service.AddDevice(userId, req)
	if err != nil {
		return nil, err
	}
	back := feedback.NewValuableFeedBack()
	err = back.AddPointsEventRsts(pointRst)
	if err != nil {
		return nil, err
	}

	return back, nil
}

// 添加aed设备
func (con *Controller) AddDevice(c *gin.Context) (interface{}, error) {
	req := new(entities.AddDevice)
	err := c.ShouldBindJSON(req)
	if err != nil {
		log.Errorf("AddDevice bind error: %v", err)
		return nil, err
	}
	accountID := c.MustGet(pkg.AccountIDKey).(int64)
	device, err := con.Service.AddDevice(accountID, req)
	if err != nil {
		log.Errorf("AddDevice error: %v", err)
		return nil, err
	}

	return device, nil
}

// 设备详情
func (con *Controller) InfoDevice(c *gin.Context) (interface{}, error) {
	req := new(InfoDevice)
	err := c.ShouldBindQuery(req)
	if err != nil {
		return nil, err
	}

	info, err := con.Service.InfoDevice(req.Longitude, req.Latitude, req.UdId)
	if err != nil {
		log.Errorf("InfoDevice error: %v", err)
		return nil, err
	}

	return info, nil
}

func (con *Controller) AddGuide(c *gin.Context) (interface{}, error) {
	accountID := c.MustGet(pkg.AccountIDKey).(int64)
	req := new(AddDeviceGuideDto)
	err := c.ShouldBindJSON(req)
	if err != nil {
		log.Errorf("AddGuide bind error: %v", err)
		return nil, err
	}

	var desc []string
	var remarks []string
	pics := [][]string{}
	for _, v := range req.Info {
		desc = append(desc, v.Desc)
		remarks = append(remarks, v.Remark)
		pics = append(pics, v.Pic)
	}

	pointRst, err := con.Service.AddGuideInfo(accountID, req.DeviceId, desc, remarks, pics)
	if err != nil {
		return nil, err
	}
	back := feedback.NewValuableFeedBack()
	err = back.AddPointsEventRsts(pointRst)
	if err != nil {
		return nil, err
	}

	return back, nil
}

func (con *Controller) ListGuide(c *gin.Context) (interface{}, error) {
	req := new(entities.DeviceGuideList)
	err := c.ShouldBindQuery(req)
	if err != nil {
		log.Errorf("ListGuide bind error: %v", err)
		return nil, err
	}

	res, err := con.Service.GetDeviceGuideInfo(req.DeviceId)
	if err != nil {
		log.Errorf("ListGuide error: %v", err)
		return nil, err
	}

	return page.Result{List: res.Info, Total: len(res.Info)}, nil
}

func (con *Controller) GetDeviceGuideInfoById(c *gin.Context) (interface{}, error) {
	req := new(DeviceGuideDto)
	err := c.ShouldBindQuery(req)
	if err != nil {
		log.Errorf("GetDeviceGuideInfoById bind error: %v", err)
		return nil, err
	}

	res, err := con.Service.GetGuideInfoById(req.Uid)
	if err != nil {
		log.Errorf("GetDeviceGuideInfoById error: %v", err)
		return nil, err
	}

	return res, nil
}

func (con *Controller) DeviceGallery(ctx *gin.Context) (interface{}, error) {
	deviceId := ctx.Param("deviceId")

	var query struct {
		Latest int `form:"latest"`
	}

	err := ctx.ShouldBindQuery(&query)
	if err != nil {
		return nil, response.NewIllegalArgumentError(err.Error())
	}

	gallery, err := con.Service.GetDeviceGallery(deviceId, query.Latest)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{"gallery": gallery}, nil
}

func (con *Controller) CountDevice(ctx *gin.Context) (interface{}, error) {
	count, err := con.Service.CountDeviceByCredibleState()
	if err != nil {
		return nil, response.NewIllegalArgumentError(err.Error())
	}

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

	return map[string]interface{}{
		"total":      total,
		"needPicket": needPicket,
	}, nil
}

func (con *Controller) SyncDevices(ctx *gin.Context) (interface{}, error) {
	err := con.Service.SyncDevices()
	return nil, err
}

func (con *Controller) ImportDevices(ctx *gin.Context) (interface{}, error) {
	form, err := ctx.MultipartForm()
	if err != nil {
		return nil, response.NewIllegalArgumentError(err.Error())
	}

	if file, ok := form.File["file"]; ok {
		header := file[0]
		f, err := header.Open()
		if err != nil {
			return nil, err
		}

		err = interfaces.S.Device.ImportDevices(f)
		return nil, err
	}

	return nil, response.NewIllegalArgumentError("invalid params")
}
