package entities

type (
	AuditResult struct {
		Code    int    `json:"code"`
		Content string `json:"content"`
		DataId  string `json:"dataId"`
		TaskId  string `json:"taskId"`
		Results []struct {
			Label      string  `json:"label"`
			Suggestion string  `json:"suggestion"`
			Rate       float32 `json:"rate"`
		}
	}
)

// Suggestion 建议
// return Suggestion, label
func (d AuditResult) Suggestion() (string, string) {
	for _, r := range d.Results {
		// 广告不违规
		if r.Label == "ad" {
			continue
		}

		if r.Suggestion == "block" || r.Suggestion == "review" {
			return r.Suggestion, r.Label
		}
	}

	return "pass", ""
}

func (d AuditResult) CheckPass() bool {
	suggestion, _ := d.Suggestion()
	return suggestion == "pass"
}
