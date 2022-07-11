package geo

import (
	"aed-api-server/internal/pkg/location"
	"github.com/mmcloughlin/geohash"
	"testing"
)

func TestEncodeIntWithPrecision(t *testing.T) {
	Lat := 30.5801
	Lon := 104.0671
	//Lat := 30.5932
	//Lon := 104.1016
	//Lat := 31.3265
	//Lon := 104.3035

	inc := 0.0001
	coordinate := location.Coordinate{
		Longitude: Lon,
		Latitude:  Lat,
	}

	var w, h int64

	var Precision uint = 31

	//EncodeWithPrecision := geohash.EncodeWithPrecision
	EncodeWithPrecision := EncodeIntWithPrecision
	hash := EncodeWithPrecision(Lat, Lon, Precision)
	n := 0.0

	println("→")
	for {
		lat2 := n*inc + Lat
		hash2 := EncodeWithPrecision(lat2, Lon, Precision)
		println(lat2, Lon, hash2, coordinate.DistanceOf(location.Coordinate{Latitude: lat2, Longitude: Lon}), geohash.EncodeInt(lat2, Lon))

		hashX := geohash.Encode(lat2, Lon)
		lat, lng := geohash.DecodeCenter(hashX)
		println("=====>", hashX, lat, lng)

		if hash2 == hash {
			n++
		} else {
			w += coordinate.DistanceOf(location.Coordinate{Latitude: lat2, Longitude: Lon})
			break
		}
	}

	println("←")
	n = 0.0
	for {
		lat2 := n*inc + Lat
		hash2 := EncodeWithPrecision(lat2, Lon, Precision)
		println(lat2, Lon, hash2, coordinate.DistanceOf(location.Coordinate{Latitude: lat2, Longitude: Lon}), geohash.EncodeInt(lat2, Lon))
		if hash2 == hash {
			n--
		} else {
			w += coordinate.DistanceOf(location.Coordinate{Latitude: lat2, Longitude: Lon})
			break
		}
	}

	println("↑")
	n = 0.0
	for {
		lon := n*inc + Lon
		hash2 := EncodeWithPrecision(Lat, lon, Precision)
		println(Lat, lon, hash2, coordinate.DistanceOf(location.Coordinate{Latitude: Lat, Longitude: lon}), geohash.EncodeInt(Lat, lon))
		if hash2 == hash {
			n++
		} else {
			h += coordinate.DistanceOf(location.Coordinate{Latitude: Lat, Longitude: lon})
			break
		}
	}

	println("↓")
	n = 0.0
	for {
		lon := n*inc + Lon
		hash2 := EncodeWithPrecision(Lat, lon, Precision)
		println(Lat, lon, hash2, coordinate.DistanceOf(location.Coordinate{Latitude: Lat, Longitude: lon}), geohash.EncodeInt(Lat, lon))
		if hash2 == hash {
			n--
		} else {
			h += coordinate.DistanceOf(location.Coordinate{Latitude: Lat, Longitude: lon})
			break
		}
	}

	println("w=", w, "h=", h)
}
