package activity

import (
	"aed-api-server/internal/interfaces/entities"
)

const (
	CategorySceneReport       = "scene-report"
	CategorySceneArrived      = "scene-arrived"
	CategoryVolunteerActivity = "volunteer-activity"
	CategoryDeviceGot         = "device-got"
	CategorySystem            = "system"
)

var categoryOrder = []string{CategorySceneReport, CategorySceneArrived, CategoryDeviceGot, CategoryVolunteerActivity, CategorySystem}

var categoryForClasses = map[string]string{
	ClassSceneReport:       CategorySceneReport,
	ClassSceneArrived:      CategorySceneArrived,
	ClassDeviceGot:         CategoryDeviceGot,
	ClassVolunteerNotified: CategorySystem,
	ClassAidCalled:         CategorySystem,
}

func GetCategory(class string) string {
	c, exists := categoryForClasses[class]

	if !exists {
		return CategoryVolunteerActivity
	}

	return c
}

func PutCategory(activity *entities.Activity) {
	activity.Category = GetCategory(activity.Class)
}
