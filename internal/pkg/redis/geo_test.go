package redis

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/pkg/cache"
	"aed-api-server/internal/pkg/location"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRadius(t *testing.T) {
	cache.InitPool(cache.RedisConfig{
		Server:    "redis-star-dev.openviewtech.com:6379",
		Password:  "",
		MaxIdle:   3,
		MaxActive: 10,
	})

	i, err := GeoRadius("device", location.Coordinate{Longitude: 104.065064, Latitude: 30.568875}, 100000)

	require.Nil(t, err)
	fmt.Println(i)

}
func TestCreate(t *testing.T) {
	cache.InitPool(cache.RedisConfig{
		Server:    "redis-star-dev.openviewtech.com:6379",
		Password:  "",
		MaxIdle:   3,
		MaxActive: 10,
	})

	items := []*entities.GeoValue{
		{
			Name: "1s",
			Coordinate: location.Coordinate{
				Longitude: 104.084002,
				Latitude:  30.420189,
			},
		},
		{
			Name: "2s",
			Coordinate: location.Coordinate{
				Longitude: 104.084002,
				Latitude:  30.420189,
			},
		},
		{
			Name: "3s",
			Coordinate: location.Coordinate{
				Longitude: 104.084002,
				Latitude:  30.420189,
			},
		},
		{
			Name: "4s",
			Coordinate: location.Coordinate{
				Longitude: 104.084002,
				Latitude:  30.420189,
			},
		},
	}
	i, err := GeoAdd("device", items)

	require.Nil(t, err)
	fmt.Println(i)
}

func TestGeoRemove(t *testing.T) {
	cache.InitPool(cache.RedisConfig{
		Server:    "redis-star-dev.openviewtech.com:6379",
		Password:  "",
		MaxIdle:   3,
		MaxActive: 10,
	})

	deleted, err := GeoRemove("device", []string{"1", "2", "3", "4"})
	require.Nil(t, err)
	fmt.Println(deleted)
}
