package activity

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/service/activity"
)

type RecordSceneReportDTO struct {
	AidID       int64    `json:"aidId,string"`
	Description string   `json:"description"`
	Images      []string `json:"images"`
}

type BorrowDevice struct {
	AidId int64 `json:"aidId,omitempty,string"`
}

type GoingToDevice struct {
	AidId int64 `json:"aidId,omitempty,string"`
}

func resolveActivityImages(activities []*entities.Activity) {
	for _, a := range activities {
		resolveActivityImage(a)
	}
}

func resolveActivityImage(a *entities.Activity) {
	if a.Class == activity.ClassSceneReport {
		if img, ok := a.Record["images"]; ok {
			switch img.(type) {
			case []interface{}:
				imgStrings := parseImagesToStrings(img.([]interface{}))
				a.Record["images"] = imgStrings
			}
		}
	}
}

func parseImagesToStrings(r []interface{}) []string {
	var images []string
	for _, img := range r {
		switch img.(type) {
		case map[string]interface{}:
			if origin, ok := img.(map[string]interface{})["origin"]; ok {
				images = append(images, origin.(string))
			}
		default:
			images = append(images, img.(string))
		}
	}

	return images
}
