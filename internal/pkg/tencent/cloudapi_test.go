package tencent

import (
	"aed-api-server/internal/pkg/location"
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

func Test_DistanceCompute(t *testing.T) {
	from := location.Coordinate{Longitude: 30.577662, Latitude: 104.065795}
	to := location.Coordinate{Longitude: 30.877662, Latitude: 104.065795}
	tos := []location.Coordinate{to, to}
	DistanceFrom(from, tos)
}
