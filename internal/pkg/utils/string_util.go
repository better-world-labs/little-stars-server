package utils

// StringLimitHidden 字符窜多余部分隐藏并以 hiddenTail 代替
func StringLimitHidden(s string, limit int, hiddenTail string) string {
	runes := []rune(s)
	if len(runes) > limit {
		runes = runes[:limit]
		if hiddenTail != "" {
			tail := []rune(hiddenTail)
			runes = append(runes, tail...)
		}
	}

	return string(runes)
}
