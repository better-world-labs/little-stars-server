package tencent

type CreateReq struct {
	Key     string      `json:"key"`
	TableID string      `json:"table_id"`
	Data    interface{} `json:"data"`
}

// 腾讯地点云表结构
type TblData struct {
	ID       string   `json:"id,omitempty"`
	UdID     string   `json:"ud_id,omitempty"`
	Title    string   `json:"title,omitempty"`
	Location Location `json:"location,omitempty"`
}

// 位置
type Location struct {
	Lat float64 `json:"lat"` // 纬度
	Lng float64 `json:"lng"` // 经度
}

// 地图查询结果
type SearchMapResp struct {
	Status    int64  `json:"status"`
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
	Result    Result `json:"result"`
}

type Result struct {
	Count   int64         `json:"count"`
	Data    []Device      `json:"data"`
	Failure []interface{} `json:"failure"`
}

type Device struct {
	Adcode       int64           `json:"adcode"`
	Address      string          `json:"address"`
	City         string          `json:"city"`
	CreateTime   int64           `json:"create_time"`
	District     string          `json:"district"`
	GeometryType int64           `json:"geometry_type"`
	ID           string          `json:"id"`
	Location     Location        `json:"location"`
	Polygon      string          `json:"polygon"`
	Province     string          `json:"province"`
	Tel          string          `json:"tel"`
	Title        string          `json:"title"` // 详细地址
	UdID         string          `json:"ud_id"`
	UpdateTime   int64           `json:"update_time"`
	X            DeviceTblCustom `json:"x"`
}

// 添加地图aed设备数据返回
type CreateResp struct {
	Status    int64  `json:"status"`
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
	Result    Result `json:"result"`
}

// aed设备自定义字段
type DeviceTblCustom struct {
	DeviceName   string `json:"device_name,omitempty"`
	Origin       string `json:"origin,omitempty"`
	Thumbnail    string `json:"thumbnail,omitempty"`
	EnvOrigin    string `json:"env_origin,omitempty"`
	EnvThumbnail string `json:"env_thumbnail,omitempty"`
	State        int    `json:"state,omitempty"`
}

// 更新地图aed设备数据请求
type UpdateReq struct {
	Key     string      `json:"key"`
	TableID string      `json:"table_id"`
	Data    interface{} `json:"data"`
	Filter  string      `json:"filter"`
}

// 更新地图aed设备数据返回
type UpdateResp struct {
	Status    int64  `json:"status"`
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
	Result    Result `json:"result"`
}

// 删除地图aed设备返回
type DelResp struct {
	Status    int64  `json:"status"`
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
	Result    Result `json:"result"`
}

// 查询地图aed设备列表返回
type ListResp struct {
	Status    int64  `json:"status"`
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
	Result    Result `json:"result"`
}

// 距离计算返回
type DistRsep struct {
	Status    int64      `json:"status"`
	Message   string     `json:"message"`
	RequestID string     `json:"request_id"`
	Result    ResultDist `json:"result"`
}

type ResultDist struct {
	Rows []Row `json:"rows"`
}

type Row struct {
	Elements []Element `json:"elements"`
}

type Element struct {
	Distance int64 `json:"distance"`
	Duration int64 `json:"duration"`
}
