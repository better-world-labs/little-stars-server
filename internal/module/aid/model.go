package aid

import (
	"aed-api-server/internal/interfaces/entities"
	"time"
)

type HelpImage struct {
	ID         int64     `xorm:"id pk autoincr"`
	HelpInfoID int64     `xorm:"help_info_id index"`
	Origin     string    `xorm:"origin"`
	Thumbnail  string    `xorm:"thumbnail"`
	Created    time.Time `xorm:"created"`
}

type HelpInfoImaged struct {
	*entities.HelpInfo

	images []*HelpImage
}
