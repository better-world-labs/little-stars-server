package utils

import (
	"fmt"
	"testing"
)

type User struct {
	Id   int64
	Name string
}

func TestToMap(t *testing.T) {
	userSlice := []User{
		{
			1, "张三",
		},
		{
			1, "张三",
		},
		{
			2, "张三2",
		},
	}

	userMap := ToMap(userSlice, func(in User) (int64, User) {
		return in.Id, in
	})

	fmt.Println(userMap)
}

func TestDistinct(t *testing.T) {
	userSlice := []User{
		{
			1, "张三",
		},
		{
			1, "张三",
		},
		{
			2, "张三2",
		},
	}

	userDistincted := Distinct(userSlice)
	fmt.Println(userDistincted)
}

func TestMap(t *testing.T) {
	userSlice := []User{
		{
			1, "张三",
		},
		{
			1, "张三",
		},
		{
			2, "张三2",
		},
	}

	fmt.Println(Map(userSlice, func(in User) int64 {
		return in.Id
	}))
}
