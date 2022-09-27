package controller

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/interfaces/service"
	"aed-api-server/internal/pkg"
	"aed-api-server/internal/pkg/feedback"
	"aed-api-server/internal/pkg/location"
	page "aed-api-server/internal/pkg/query"
	"aed-api-server/internal/pkg/response"
	"aed-api-server/internal/pkg/utils"
	"aed-api-server/internal/service/device"
	"errors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gitlab.openviewtech.com/openview-pub/gopkg/route"
)

type DeviceController struct {
	Service service.DeviceService `inject:"-"`
}

//go:inject-component
func NewDeviceController() *DeviceController {
	return &DeviceController{}
}

func (con *DeviceController) MountAuthRouter(r *route.Router) {
	r.GET("/devices/:deviceId/gallery", con.DeviceGallery)
	r.GET("/aed/devices", con.ListDevices)
	r.GET("/encrypted/aed/devices", con.ListDeviceEncrypted)
	r.GET("/aed/latest/:latest/devices", con.ListLatestDevices)
	r.POST("/aed/devices/query-by-ids", con.queryByIds)
	r.POST("/devices", con.MarkDevice)
	r.POST("/aed/add/device", con.AddDevice)

	r.GET("/aed/device", con.InfoDevice)
	r.GET("/aed/risk-area", con.RiskArea)
	r.POST("/aed/device/add_guide", con.AddGuide)
	r.GET("/aed/device/guide", con.ListGuide)
	r.GET("/aed/device/guide_info", con.GetDeviceGuideInfoById)
}

func (con *DeviceController) MountAdminRouter(r *route.Router) {
	v1 := r.Group("/v1")
	v1.POST("/import-devices", con.ImportDevices)
	v1.POST("/sync-devices", con.SyncDevices)
	v1.GET("/devices", con.AdminListDevices)
	v1.POST("/devices", con.AdminCreateDevice)
	v1.GET("/devices/:id", con.AdminGetDevice)
	v1.PUT("/devices/:id", con.AdminUpdateDevice)
	v1.DELETE("/devices", con.AdminDeleteDevices)
}

// ListDevices 查询设备列表
func (con *DeviceController) ListDevices(c *gin.Context) (interface{}, error) {
	type ListDevice struct {
		location.Coordinate

		Distance float64 `form:"distance,omitempty"`
	}

	req := new(ListDevice)
	err := c.ShouldBindQuery(req)
	if err != nil {
		log.Errorf("ListDevices bind error: %v", err)
		return nil, err
	}

	list, err := con.Service.ListDevices(req.Coordinate, req.Distance)
	if err != nil {
		log.Errorf("ListDevices error: %v", err)
		return nil, err
	}

	log.Debugf("ListDevices lnt: %v,lat: %v, len:%v", req.Longitude, req.Latitude, len(list))

	return page.NewResult[*entities.Device](list, len(list)), nil
}

// 标记设备
func (con *DeviceController) MarkDevice(c *gin.Context) (interface{}, error) {
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
func (con *DeviceController) AddDevice(c *gin.Context) (interface{}, error) {
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

// InfoDevice 设备详情
func (con *DeviceController) InfoDevice(c *gin.Context) (interface{}, error) {
	type InfoDevice struct {
		UdId      string  `form:"id,omitempty"`
		Longitude float64 `form:"longitude,omitempty"`
		Latitude  float64 `form:"latitude,omitempty"`
	}

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

func (con *DeviceController) AddGuide(c *gin.Context) (interface{}, error) {
	type AddDeviceGuideDto struct {
		DeviceId string               `json:"deviceId"`
		Info     []entities.GuideInfo `json:"info"`
	}

	accountID := c.MustGet(pkg.AccountIDKey).(int64)
	req := new(AddDeviceGuideDto)
	err := c.ShouldBindJSON(req)
	if err != nil {
		log.Errorf("AddGuide bind error: %v", err)
		return nil, err
	}

	var desc []string
	var remarks []string
	var pics [][]string
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

func (con *DeviceController) ListGuide(c *gin.Context) (interface{}, error) {
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

	return page.NewResult[entities.DeviceGuideListItem](res.Info, len(res.Info)), nil
}

func (con *DeviceController) GetDeviceGuideInfoById(c *gin.Context) (interface{}, error) {
	type DeviceGuideDto struct {
		Uid string `json:"uid" form:"uid"`
	}

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

func (con *DeviceController) DeviceGallery(ctx *gin.Context) (interface{}, error) {
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

func (con *DeviceController) CountDevice(ctx *gin.Context) (interface{}, error) {
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
		if c.CredibleState == device.CredibleStatusDeviceNotFound {
			needPicket += c.Count
		}
	}

	return map[string]interface{}{
		"total":      total,
		"needPicket": needPicket,
	}, nil
}

func (con *DeviceController) SyncDevices(ctx *gin.Context) (interface{}, error) {
	err := con.Service.SyncDevices()
	return nil, err
}

func (con *DeviceController) ImportDevices(ctx *gin.Context) (interface{}, error) {
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

func (con *DeviceController) queryByIds(context *gin.Context) (interface{}, error) {
	type Req struct {
		location.Coordinate
		DeviceIds []string `json:"deviceIds"`
	}
	var req Req
	if err := context.ShouldBindJSON(&req); err != nil {
		return nil, err
	}

	if req.DeviceIds == nil || len(req.DeviceIds) == 0 {
		return nil, errors.New("设备Id不能为空")
	}
	return con.Service.ListDevicesByIDs(req.Coordinate, req.DeviceIds)
}

func (con *DeviceController) ListLatestDevices(ctx *gin.Context) (interface{}, error) {
	latest, err := utils.GetContextPathParamInt64(ctx, "latest")
	if err != nil {
		return nil, err
	}

	devices, err := con.Service.ListLatestDevices(latest)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"devices": devices,
	}, nil

}

func (con *DeviceController) RiskArea(ctx *gin.Context) (interface{}, error) {
	var center struct {
		Longitude float64 `form:"longitude" binding:"required"`
		Latitude  float64 `form:"latitude" binding:"required"`
	}

	err := ctx.ShouldBindQuery(&center)
	if err != nil {
		return nil, err
	}

	area, err := con.Service.RiskArea(location.Coordinate{Longitude: center.Longitude, Latitude: center.Latitude})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"areas": area,
	}, nil
}

func (con *DeviceController) ListDeviceEncrypted(c *gin.Context) (interface{}, error) {
	userId := utils.GetContextUserId(c)
	type ListDevice struct {
		location.Coordinate

		Distance   float64 `form:"distance,omitempty"`
		KeyVersion int     `form:"keyVersion,omitempty"`
	}

	req := new(ListDevice)
	err := c.ShouldBindQuery(req)
	if err != nil {
		log.Errorf("ListDevices bind error: %v", err)
		return nil, err
	}

	encrypted, err := con.Service.ListDevicesEncrypted(userId, req.Coordinate, req.Distance, req.KeyVersion)
	if err != nil {
		log.Errorf("ListDevices error: %v", err)
		return nil, err
	}

	log.Debugf("ListDevices lnt: %v,lat: %v, len:%v", req.Longitude, req.Latitude, len(encrypted))

	return map[string]interface{}{
		"encryptedData": encrypted,
	}, nil
}

func (con *DeviceController) AdminListDevices(ctx *gin.Context) (interface{}, error) {
	page, err := page.BindPageQuery(ctx)
	if err != nil {
		return nil, err
	}

	keyword := ctx.Query("keyword")
	devices, err := con.Service.PageDevices(page, keyword)
	if err != nil {
		return nil, err
	}

	return devices, nil
}

func (con *DeviceController) AdminGetDevice(ctx *gin.Context) (interface{}, error) {
	id := ctx.Param("id")
	d, err := con.Service.GetById(id)
	if err != nil {
		return nil, err
	}

	return d, nil
}

func (con *DeviceController) AdminCreateDevice(ctx *gin.Context) (interface{}, error) {
	userId := utils.GetContextUserId(ctx)
	var d entities.BaseDevice

	err := ctx.ShouldBindJSON(&d)
	if err != nil {
		return nil, err
	}

	d.CreateBy = userId
	return nil, con.Service.CreateDevice(&d)
}

func (con *DeviceController) AdminUpdateDevice(ctx *gin.Context) (interface{}, error) {
	id := ctx.Param("id")

	var d entities.BaseDevice
	err := ctx.ShouldBindJSON(&d)

	if err != nil {
		return nil, err
	}

	d.Id = id

	return nil, con.Service.UpdateDevice(&d)
}

func (con *DeviceController) AdminDeleteDevices(ctx *gin.Context) (interface{}, error) {
	var param struct {
		Ids []string `form:"ids" binding:"required"`
	}

	err := ctx.ShouldBindJSON(&param)
	if err != nil {
		return nil, err
	}

	return nil, con.Service.DeleteDevices(param.Ids)
}
