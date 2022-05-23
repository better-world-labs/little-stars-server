package task

import (
	"aed-api-server/internal/interfaces"
	"testing"
)

func Test_GenJobsByUserLocation(t *testing.T) {
	println("this is a test2")
	println("this is a test")
}

func Test_FindUserTaskByUserIdAndDeviceId(t *testing.T) {
	t.Cleanup(InitDbAndConfig())
	mockPicketService()
	interfaces.S.Task.FindUserTaskByUserIdAndDeviceId(53, "a4ee1809-a932-4c40-9029-1856b8a10dde")
}
