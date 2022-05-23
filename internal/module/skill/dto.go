package skill

import "aed-api-server/internal/pkg/global"

type QueryProjects struct {
	Account string `json:"account"`
}

// 用户认证项目返回
type UserProjectDto struct {
	ProjectID   string   `json:"projectId"`
	Status      string   `json:"status"`
	Description string   `json:"description"`
	Title       string   `json:"title"`
	Pic         string   `json:"pic"`
	Points      float64  `json:"points"`
	Images      []string `json:"images"`
	Exams       []struct {
		Id      string   `json:"id"`
		Title   string   `json:"title"`
		Options []string `json:"options"`
		Type    int      `json:"type"`
		Sort    int      `json:"sort"`
		Answers []string `json:"answers"`
	} `json:"exams"`
}

type UserProjectsDto []UserProjectDto

func (arr *UserProjectsDto) Push(e ...UserProjectDto) {
	*arr = append(*arr, e...)
}

// 用户保存答题
type UserExamSaveDto struct {
	ProjectID string `json:"projectId"`
	List      []struct {
		ID      string   `json:"examId"`
		Answers []string `json:"answers"`
	} `json:"list"`
}

// 用户考试提交
type UserExamSumbitDto struct {
	ProjectID string `json:"projectId"`
	List      []struct {
		ID      string   `json:"examId"`
		Answers []string `json:"answers"`
	} `json:"list"`
}

type ExamDto struct {
	Id            int      `json:"id,omitempty"`
	ExamId        string   `json:"examId,omitempty"`
	ProjectId     int      `json:"projectId,omitempty"`
	Title         string   `json:"title,omitempty"`
	Options       []string `json:"options,omitempty"`
	CorrectAnswer []string `json:"correctAnswer,omitempty"`
	Sort          int      `json:"sort,omitempty"`
	Answer        []string `json:"answers,omitempty"`
	Type          int
}

// 学习时长3分钟
type UseStudyDuration struct {
	ProjectID string `json:"id"`
}

type QueryProject struct {
	ProjectId string `form:"projectId,omitempty"`
}
type DtoCert struct {
	Uid         string               `json:"certId,omitempty"`
	ProjectId   int                  `json:"projectId"`
	ProjectName string               `json:"projectName"`
	Origin      string               `json:"origin"`
	Thumbnail   string               `json:"thumbnail"`
	Created     global.FormattedTime `json:"time,omitempty"`
}
