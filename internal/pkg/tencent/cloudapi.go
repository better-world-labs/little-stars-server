package tencent

import (
	"aed-api-server/internal/pkg/location"
	page "aed-api-server/internal/pkg/query"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

// 腾讯地点云 https://lbs.qq.com/service/placeCloud/placeCloudGuide/cloudOverview

type Config struct {
	APIKey    string `yaml:"apikey"`
	SecretKey string `yaml:"apiksecret"`
	TblDevice string `yaml:"tbl_device"` // aed设备云端表
}

var config *Config

func Init(c *Config) {
	config = c
}

func QueryDeviceTblStruct(pageNow, pageSize int) {
	url := fmt.Sprintf("https://apis.map.qq.com/place_cloud/data/list?table_id=%v&orderby=id&page_index=1&page_size=20&key=%v", config.TblDevice, config.APIKey)
	fmt.Println(Get(url))
}

func ListRangeDeviceIDs(center location.Coordinate, radius float64, query page.Query) ([]string, error) {
	if radius > 10000 {
		radius = 10000
	}

	if query.Page == 0 {
		query.Page = 1
	}

	if query.Size == 0 {
		query.Size = 200
	}

	params := make(map[string]string)
	sprintf := fmt.Sprintf("%v,%v", center.Latitude, center.Longitude)
	params["location"] = sprintf
	params["radius"] = fmt.Sprintf("%v", radius)
	//TODO 之前这里有个这玩意
	//params["filter"] = "x.device_name=aed"
	params["table_id"] = config.TblDevice
	params["key"] = config.APIKey
	//params["fields"] = fmt.Sprintf("ud_id,distance(%f,%f) distance", center.Latitude, center.Longitude)
	params["page_size"] = fmt.Sprintf("%v", query.Size)

	var dataParams string
	for k, _ := range params {
		dataParams = dataParams + k + "=" + params[k] + "&"
	}
	paramStr := dataParams[0 : len(dataParams)-1]
	sign := getSign("/place_cloud/search/nearby", params)
	url := fmt.Sprintf(`https://apis.map.qq.com/place_cloud/search/nearby?%v&sig=%v`, paramStr, sign)
	data, err := Get(url)
	if err != nil {
		return nil, err
	}

	resp := new(SearchMapResp)
	//resp := make(map[string]interface{}, 0)
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, err
	}

	if resp.Status != 0 {
		return nil, fmt.Errorf("place_cloud/search/nearby return status %d, message=%s", resp.Status, resp.Message)
	}

	var deviceIds []string
	for _, d := range resp.Result.Data {
		deviceIds = append(deviceIds, d.UdID)
	}

	return deviceIds, nil
}

func SearchMapDevice(lng, lat, distance float64, page, size int) ([]Device, error) {
	if distance > 10000 {
		distance = 10000
	}
	params := make(map[string]string)
	params["location"] = fmt.Sprintf("%v,%v", lat, lng)
	params["radius"] = fmt.Sprintf("%v", distance)
	params["filter"] = "x.device_name=aed"
	params["table_id"] = config.TblDevice
	params["key"] = config.APIKey
	params["page_size"] = fmt.Sprintf("%v", size)

	var dataParams string
	for k, _ := range params {
		dataParams = dataParams + k + "=" + params[k] + "&"
	}
	paramStr := dataParams[0 : len(dataParams)-1]
	sign := getSign("/place_cloud/search/nearby", params)
	url := fmt.Sprintf(`https://apis.map.qq.com/place_cloud/search/nearby?%v&sig=%v`, paramStr, sign)

	data, err := Get(url)
	if err != nil {
		return nil, err
	}

	resp := new(SearchMapResp)
	err = json.Unmarshal(data, resp)
	if err != nil {
		return nil, err
	}

	if resp.Status != 0 {
		return nil, fmt.Errorf("place_cloud/search/nearby return status %v", resp.Status)
	}

	return resp.Result.Data, nil
}

func AddDevice(lng, lat float64, title string) (string, error) {
	data := make([]*TblData, 0)
	device := new(TblData)
	device.UdID = uuid.New().String()
	device.Title = title
	device.Location = Location{Lng: lng, Lat: lat}
	device.Polygon = ""
	data = append(data, device)
	j, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	strData := string(j)
	params := make(map[string]string)
	params["key"] = config.APIKey
	params["table_id"] = config.TblDevice
	params["data"] = strData
	sign := getSign("/place_cloud/data/create", params)
	url := fmt.Sprintf("%v?sig=%v", "https://apis.map.qq.com/place_cloud/data/create", sign)

	req := new(CreateReq)
	req.Key = config.APIKey
	req.TableID = config.TblDevice
	req.Data = data
	b, _ := json.Marshal(req)
	res, err := Post(string(b), url)
	if err != nil {
		return "", err
	}

	resp := new(CreateResp)
	err = json.Unmarshal(res, resp)
	if err != nil {
		return "", err
	}
	if resp.Status != 0 {
		return "", fmt.Errorf("place_cloud/data/create return status %v", string(res))
	}
	if len(resp.Result.Failure) > 0 {
		return "", fmt.Errorf("place_cloud/data/update return failure %v", string(res))
	}

	return device.UdID, nil
}

func DelDevice(udid []string) error {
	del := ""
	for _, v := range udid {
		del += fmt.Sprintf("\"%v\"", v) + ","
	}

	params := make(map[string]string)
	params["key"] = config.APIKey
	params["table_id"] = config.TblDevice
	params["filter"] = fmt.Sprintf("ud_id in(%v)", del[0:len(del)-1])

	sign := getSign("/place_cloud/data/delete", params)
	url := fmt.Sprintf("%v?sig=%v", "https://apis.map.qq.com/place_cloud/data/delete", sign)

	j, _ := json.Marshal(params)
	// fmt.Println("post ", string(j))
	res, err := Post(string(j), url)
	if err != nil {
		return err
	}

	fmt.Println(string(res))
	resp := new(DelResp)
	err = json.Unmarshal(res, resp)
	if err != nil {
		return err
	}
	if resp.Status != 0 {
		return fmt.Errorf("place_cloud/data/delete return status %v", string(res))
	}
	return nil
}

func ListDevice(pageIndex, pageSize int) ([]Device, error) {
	params := make(map[string]string)
	params["key"] = config.APIKey
	params["table_id"] = config.TblDevice
	params["page_index"] = fmt.Sprintf("%v", pageIndex)
	params["page_size"] = fmt.Sprintf("%v", pageSize)

	var dataParams string
	for k, _ := range params {
		dataParams = dataParams + k + "=" + params[k] + "&"
	}
	paramStr := dataParams[0 : len(dataParams)-1]
	sign := getSign("/place_cloud/data/list", params)
	url := fmt.Sprintf(`https://apis.map.qq.com/place_cloud/data/list?%v&sig=%v`, paramStr, sign)
	res, err := Get(url)
	if err != nil {
		return nil, err
	}
	resp := new(ListResp)
	err = json.Unmarshal(res, resp)
	if err != nil {
		return nil, err
	}
	if resp.Status != 0 {
		return nil, fmt.Errorf("place_cloud/data/list return status %v", string(res))
	}
	return resp.Result.Data, nil
}
