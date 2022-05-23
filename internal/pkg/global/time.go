package global

import (
	"fmt"
	"strings"
	"time"
)

type FormattedTime time.Time

func (t *FormattedTime) UnmarshalJSON(b []byte) error {
	s := string(b)
	fmt.Printf("%s", s)

	s = strings.ReplaceAll(s, "\"", "")
	time, err := time.Parse("2006-01-02 15:04:05", s)
	if err != nil {
		return err
	}

	*t = FormattedTime(time)
	return nil
}

//TODO 这里使用指针接收会丢失值
func (t FormattedTime) MarshalJSON() ([]byte, error) {
	var stamp = fmt.Sprintf("\"%s\"", time.Time(t).Format("2006-01-02 15:04:05"))
	return []byte(stamp), nil
}

func (t *FormattedTime) Time() time.Time {
	return time.Time(*t)
}
