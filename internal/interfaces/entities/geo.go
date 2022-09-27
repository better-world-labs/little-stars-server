package entities

import "aed-api-server/internal/pkg/location"

type GeoValue struct {
	Name string
	location.Coordinate
}

type DistancedGeoValue struct {
	GeoValue

	Distance float64
}

func GeoValueFromDevice(device *BaseDevice) *GeoValue {
	return &GeoValue{
		Name:       device.Id,
		Coordinate: location.Coordinate{Longitude: device.Longitude, Latitude: device.Latitude},
	}
}

func GeoValuesFromDevices(devices []*BaseDevice) []*GeoValue {
	var values []*GeoValue

	for _, d := range devices {
		values = append(values, GeoValueFromDevice(d))
	}

	return values
}

func GeoValuesMapNames(values []*GeoValue) []string {
	var res []string

	for _, v := range values {
		res = append(res, v.Name)
	}

	return res
}

func DistancedGeoValuesMapNames(values []*DistancedGeoValue) []string {
	var res []string

	for _, v := range values {
		res = append(res, v.Name)
	}

	return res
}
