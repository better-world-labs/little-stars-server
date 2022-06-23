package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_goFileParse(t *testing.T) {
	rst, err := goFileParse("./testdata/activity_service.go")
	assert.Nil(t, err)

	assert.Equal(t, rst.PkgName, "service")
	assert.Equal(t, rst.InjectNames[0], "NewActivityService")
	assert.Equal(t, rst.InjectNames[1], "XxTest22321_x")
}
