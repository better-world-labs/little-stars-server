package tencent

import (
	"aed-api-server/internal/pkg/location"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	APIKey    = "VKMBZ-RGA66-6NNSX-ECTFS-42I73-5FFBP"
	SecretKey = "JctfIJRbnRe3Gu9DIG4W14yDZ3p"

	TblDevice = "j4d9MlkKBM7ro1" // aed设备云端表
	TblTest   = "aed"            // 测试表
)

func init() {
	c := new(Config)
	c.APIKey = APIKey
	c.SecretKey = SecretKey
	c.TblDevice = TblDevice
	config = c
}

func Test_QueryDeviceTbl(t *testing.T) {
	params := make(map[string]string)
	params["table_id"] = TblDevice
	params["key"] = APIKey
	sign := getSign("/place_cloud/table/list", params)
	url := fmt.Sprintf(`https://apis.map.qq.com/place_cloud/table/list?key=%v&table_id=%v&sig=%v`, APIKey, TblDevice, sign)

	res, _ := Get(url)
	fmt.Println(string(res))
}

func Test_AddDevice(t *testing.T) {
	//udid, err := AddDevice(30.657789, 104.065795, "测试设备3", "detailAddress", []string{"1.png", "2.png", "3.png", "4.png"})
	//_ = udid
	//if err != nil {
	//	t.Error("Test_AddDevice failed")
	//}
}

func Test_SearchMapDevice(t *testing.T) {
	// "30.657789,104.065795"
	lng := 104.069969
	lat := 30.577662

	Init(&Config{APIKey: APIKey, SecretKey: SecretKey, TblDevice: TblDevice})
	ids, err := SearchMapDevice(lng, lat, 4000, 1, 10)
	//ids, err := ListRangeDeviceIDs(location.Coordinate{lng, lat}, 4000, page.Query{Size: 200})
	require.Nil(t, err)
	fmt.Println(ids, err)
}

func Test_UpdateDevice(t *testing.T) {
	//UpdateDevice("95d58ad6-d5d0-4ae8-b60e-da8454123431", 10, 10, "test", "detailAddress", []string{"1.png", "2.png", "3.png", "4.png"})
}

func Test_DelDevice(t *testing.T) {
	//DelDevice("1")
}

func Test_ListDevice(t *testing.T) {
	device, err := ListDevice(1, 10)
	require.Nil(t, err)
	fmt.Println(device)
}

func Test_DistanceCompute(t *testing.T) {
	from := location.Coordinate{Longitude: 30.577662, Latitude: 104.065795}
	to := location.Coordinate{Longitude: 30.877662, Latitude: 104.065795}
	tos := []location.Coordinate{to, to}
	DistanceFrom(from, tos)
}
