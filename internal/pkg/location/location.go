package location

import (
	"fmt"
	"math"
)

// Coordinate 经纬度坐标
type Coordinate struct {
	Longitude float64 `xorm:"longitude" json:"longitude" form:"longitude"`
	Latitude  float64 `xorm:"latitude" json:"latitude" form:"latitude"`
}

// DistanceOf 坐标之间的距离计算
// @param other 另一个坐标
// @return 距离 (m)
func (c *Coordinate) DistanceOf(other Coordinate) int64 {
	distance := computeDistance(c.Longitude, c.Latitude, other.Longitude, other.Latitude)
	return int64(math.Floor(distance + 0/5))
}

func (c *Coordinate) ToTencentStr() string {
	return fmt.Sprintf("%v,%v", c.Latitude, c.Longitude)
}

func rad(d float64) (r float64) {
	r = d * math.Pi / 180.0
	return
}

func computeDistance(lon1, lat1, lon2, lat2 float64) (distance float64) {
	//赤道半径(单位m)
	const EARTH_RADIUS = 6378137
	rad_lat1 := rad(lat1)
	rad_lon1 := rad(lon1)
	rad_lat2 := rad(lat2)
	rad_lon2 := rad(lon2)
	if rad_lat1 < 0 {
		rad_lat1 = math.Pi/2 + math.Abs(rad_lat1)
	}
	if rad_lat1 > 0 {
		rad_lat1 = math.Pi/2 - math.Abs(rad_lat1)
	}
	if rad_lon1 < 0 {
		rad_lon1 = math.Pi*2 - math.Abs(rad_lon1)
	}
	if rad_lat2 < 0 {
		rad_lat2 = math.Pi/2 + math.Abs(rad_lat2)
	}
	if rad_lat2 > 0 {
		rad_lat2 = math.Pi/2 - math.Abs(rad_lat2)
	}
	if rad_lon2 < 0 {
		rad_lon2 = math.Pi*2 - math.Abs(rad_lon2)
	}
	x1 := EARTH_RADIUS * math.Cos(rad_lon1) * math.Sin(rad_lat1)
	y1 := EARTH_RADIUS * math.Sin(rad_lon1) * math.Sin(rad_lat1)
	z1 := EARTH_RADIUS * math.Cos(rad_lat1)

	x2 := EARTH_RADIUS * math.Cos(rad_lon2) * math.Sin(rad_lat2)
	y2 := EARTH_RADIUS * math.Sin(rad_lon2) * math.Sin(rad_lat2)
	z2 := EARTH_RADIUS * math.Cos(rad_lat2)
	d := math.Sqrt((x1-x2)*(x1-x2) + (y1-y2)*(y1-y2) + (z1-z2)*(z1-z2))
	theta := math.Acos((EARTH_RADIUS*EARTH_RADIUS + EARTH_RADIUS*EARTH_RADIUS - d*d) / (2 * EARTH_RADIUS * EARTH_RADIUS))
	distance = theta * EARTH_RADIUS
	return
}
