package geo

import "github.com/mmcloughlin/geohash"

const Precision uint = 31

var EncodeIntWithPrecision = geohash.EncodeIntWithPrecision

func Hash(lat, lng float64) uint64 {
	return geohash.EncodeIntWithPrecision(lat, lng, Precision)
}
