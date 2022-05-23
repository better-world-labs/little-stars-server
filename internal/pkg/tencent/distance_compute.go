package tencent

import (
	"aed-api-server/internal/pkg/location"
	"aed-api-server/internal/pkg/utils"
	"encoding/json"
	"fmt"

	"gitlab.openviewtech.com/openview-pub/gopkg/log"
)

// 距离批量计算 from起始位置 to终点位置
func DistanceFrom(from location.Coordinate, to []location.Coordinate) (distances []int64) {
	defer utils.TimeStat("DistanceFrom")()
	if len(to) == 0 {
		return []int64{}
	}

	log.DefaultLogger().Infof("tencent_distance_from called from=%v, to=%v", from, to)
	if !inChina(from) {
		l := len(to)
		distances = make([]int64, l)
		for i := 0; i < l; i++ {
			distances[i] = -1
		}

		return
	}

	distances, err := DistanceFromCloud(from, to)
	if err != nil {
		log.DefaultLogger().Warnf("DistanceFromCloud failed: %v", err)

		// 腾讯API 请求失败则本地计算直线距离
		distances = DistanceFromLocal(from, to)
	}

	return
}

func DistanceFromCloud(from location.Coordinate, to []location.Coordinate) (distances []int64, err error) {
	distances = make([]int64, len(to))
	mode := "driving"
	APIKey := config.APIKey
	params := make(map[string]string)
	params["mode"] = mode
	params["key"] = APIKey
	params["from"] = from.ToTencentStr()
	toStr := ""
	for _, v := range to {
		toStr = toStr + v.ToTencentStr() + ";"
	}
	params["to"] = toStr[0 : len(toStr)-1]
	sign := getSign("/ws/distance/v1/matrix", params)
	url := fmt.Sprintf("%s?mode=%s&from=%s&to=%s&key=%s&sig=%s", "https://apis.map.qq.com/ws/distance/v1/matrix", mode, params["from"], params["to"], APIKey, sign)
	res, err := Get(url)
	if err != nil {
		return nil, err
	}
	resp := new(DistRsep)
	err = json.Unmarshal(res, resp)
	if err != nil {
		return nil, err
	}

	if resp.Status != 0 || resp.Message != "query ok" {
		log.DefaultLogger().Errorf("DistanceFrom error: url %v", url)
		return nil, fmt.Errorf("/ws/distance/v1/matrix status:%v message:%v", resp.Status, resp.Message)
	}

	for i, v := range resp.Result.Rows[0].Elements {
		distances[i] = v.Distance
	}

	return
}

func DistanceFromLocal(from location.Coordinate, to []location.Coordinate) (distances []int64) {
	distances = make([]int64, len(to))
	for i, t := range to {
		distances[i] = from.DistanceOf(t)
	}

	return
}

func inChina(from location.Coordinate) bool {
	return from.Longitude >= 73.33 && from.Longitude <= 135.05 &&
		from.Latitude >= 3.15 && from.Latitude <= 53.33
}
