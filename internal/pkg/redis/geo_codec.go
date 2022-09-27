package redis

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/pkg/location"
	"errors"
	"strconv"
)

func DecodeDistancedGeoValue(reply interface{}) (*entities.DistancedGeoValue, error) {
	if i, ok := reply.([]interface{}); ok {
		name, err := decodeString(i[0])
		if err != nil {
			return nil, err
		}

		distance, err := decodeFloat(i[1])
		if err != nil {
			return nil, err
		}

		coordinate, err := decodeCoordinate(i[2])
		if err != nil {
			return nil, err
		}

		return &entities.DistancedGeoValue{
			GeoValue: entities.GeoValue{
				Name:       name,
				Coordinate: *coordinate,
			},
			Distance: distance,
		}, nil
	}

	return nil, errors.New("invalid type")
}

func decodeString(n interface{}) (string, error) {
	if name, ok := n.([]uint8); ok {
		return string(name), nil
	}

	return "", errors.New("invalid type")
}

func decodeCoordinate(n interface{}) (*location.Coordinate, error) {
	if coordinate, ok := n.([]interface{}); ok {
		longitude, err := decodeFloat(coordinate[0])
		if err != nil {
			return nil, err
		}

		latitude, err := decodeFloat(coordinate[1])
		if err != nil {
			return nil, err
		}

		return &location.Coordinate{
			Longitude: longitude,
			Latitude:  latitude,
		}, nil
	}

	return nil, errors.New("invalid type")
}

func decodeFloat(n interface{}) (float64, error) {
	if longitude, ok := n.([]uint8); ok {
		return strconv.ParseFloat(string(longitude), 32)
	}

	return 0, errors.New("invalid type")
}
