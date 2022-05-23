package entities

type Feedback struct {
	Type    int      `json:"type"`
	Content string   `json:"content"`
	Images  []string `json:"images"`
}
