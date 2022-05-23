package entities

import "aed-api-server/internal/pkg/global"

type (
	Activity struct {
		ID         int64                  `xorm:"id pk autoincr" json:"id,string"`
		Uuid       string                 `xorm:"uuid" json:"-"`
		HelpInfoID int64                  `xorm:"help_info_id" json:"aid,string"`
		Points     int                    `xorm:"points" json:"point"`
		Class      string                 `xorm:"class" json:"class"`
		Category   string                 `xorm:"-" json:"category"`
		UserID     *int64                 `xorm:"user_id" json:"userId,string,omitempty"`
		Record     map[string]interface{} `xorm:"record" json:"record,omitempty"`
		Created    global.FormattedTime   `xorm:"created" json:"created"`
	}
)
