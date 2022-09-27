package redis

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/pkg/cache"
	"aed-api-server/internal/pkg/location"
	"errors"
	"github.com/gomodule/redigo/redis"
)

var (
	MaxCreateOnce         = 2000
	MaxRadius     float64 = 30000
)

func GeoAdd(key string, values []*entities.GeoValue) (int64, error) {
	conn := cache.GetConn()
	defer conn.Close()

	i, err := geoAdd(conn, key, values)
	return i, err
}

func GeoRadius(key string, from location.Coordinate, radius float64) ([]*entities.DistancedGeoValue, error) {
	conn := cache.GetConn()
	defer conn.Close()

	if radius > MaxRadius {
		radius = MaxRadius
	}

	reply, err := conn.Do("GEORADIUS", key, from.Longitude, from.Latitude, radius, "m", "WITHDIST", "WITHCOORD", "ASC")
	if err != nil {
		return nil, err
	}

	if arr, ok := reply.([]interface{}); ok {
		var values []*entities.DistancedGeoValue
		for _, i := range arr {
			value, err := DecodeDistancedGeoValue(i)
			if err != nil {
				return nil, err
			}

			values = append(values, value)
		}

		return values, nil
	}

	return nil, errors.New("invalid type")
}

func GeoRadiusCount(key string, from location.Coordinate, radius float64) (int, error) {
	conn := cache.GetConn()
	defer conn.Close()

	if radius > MaxRadius {
		radius = MaxRadius
	}

	reply, err := conn.Do("GEORADIUS", key, from.Longitude, from.Latitude, radius, "m")
	if err != nil {
		return 0, err
	}

	return len(reply.([]interface{})), nil
}

func geoAdd(conn redis.Conn, key string, points []*entities.GeoValue) (int64, error) {
	if len(points) > MaxCreateOnce {
		row1, err := geoAdd(conn, key, points[:MaxCreateOnce])
		if err != nil {
			return row1, err
		}

		row2, err := geoAdd(conn, key, points[MaxCreateOnce:])
		if err != nil {
			return row2, err
		}

		return row1 + row2, nil
	}

	var args []interface{}
	args = append(args, key)

	for _, p := range points {
		args = append(args, p.Longitude, p.Latitude, p.Name)
	}

	reply, err := conn.Do("GEOADD", args...)

	if err != nil {
		return 0, err
	}

	return reply.(int64), nil
}

func GeoUpdate(key string, points []*entities.GeoValue) (int64, error) {
	conn := cache.GetConn()
	defer conn.Close()

	_, err := geoRemove(conn, key, entities.GeoValuesMapNames(points))
	if err != nil {
		return 0, err
	}

	return geoAdd(conn, key, points)
}

func GeoRemove(key string, names []string) (int64, error) {
	conn := cache.GetConn()
	defer conn.Close()

	return geoRemove(conn, key, names)
}

func GeoListNames(key string) ([]string, error) {
	conn := cache.GetConn()
	defer conn.Close()

	reply, err := conn.Do("ZRANGE", key, 0, -1)
	var res []string
	rep := reply.([]interface{})

	for _, r := range rep {
		res = append(res, string(r.([]uint8)))
	}

	return res, err
}

func geoRemove(conn redis.Conn, key string, names []string) (int64, error) {
	if len(names) > 2000 {
		deleted1, err := geoRemove(conn, key, names[:2000])
		if err != nil {
			return 0, err
		}

		deleted2, err := geoRemove(conn, key, names[2000:])
		if err != nil {
			return 0, err
		}

		return deleted1 + deleted2, nil
	}

	var args []interface{}
	args = append(args, key)

	for _, n := range names {
		args = append(args, n)
	}

	reply, err := conn.Do("ZREM", args...)
	if err != nil {
		return 0, err
	}

	return reply.(int64), nil
}
