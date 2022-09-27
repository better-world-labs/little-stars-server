package entities

import "encoding/json"

type AnswerOptionIndex string

//AnswerOption 答案选项
type AnswerOption struct {
	Index AnswerOptionIndex `json:"index"` //标号
	Desc  string            `json:"desc"`  //描述
	Score int               `json:"-"`     //分数，放大100倍记录整数
}

type Question struct {
	Id       int             `json:"id"`      //题目ID
	OriginNo string          `json:"-"`       //原问卷ID
	Type     string          `json:"type"`    //题目类型
	SubType  string          `json:"-"`       //题目子类
	Desc     string          `json:"desc"`    //题目描述
	Options  []*AnswerOption `json:"options"` //题目选项
}

type Answer struct {
	QuestionId int                    `json:"questionId"` //题目ID
	Select     AnswerOptionIndex      `json:"select"`     //选项
	Extra      map[string]interface{} `json:"extra"`      //额外输入，json
}

func (a *Answer) FromDB(b []byte) error {
	return json.Unmarshal(b, a)
}

func (a *Answer) ToDB() ([]byte, error) {
	return json.Marshal(a)
}

type Result struct {
	LevelId int    `json:"levelId"`
	Level   string `json:"level"`
	Explain string `json:"explain"`
}

func (r *Result) FromDB(b []byte) error {
	return json.Unmarshal(b, r)
}

func (r *Result) ToDB() ([]byte, error) {
	return json.Marshal(r)
}

type ShareUser struct {
	NickName string `json:"nickName"`
	Avatar   string `json:"avatar"`
}

type HealthyGameInfo struct {
	HadResult      bool        `json:"hadResult"`
	TodayDeadCount int         `json:"todayDeadCount"`
	ShareUser      *ShareUser  `json:"shareUser"`
	Questions      []*Question `json:"questions"`
}
