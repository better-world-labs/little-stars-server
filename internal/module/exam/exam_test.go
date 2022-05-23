package exam

import (
	"aed-api-server/internal/domains"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestExam(t *testing.T) {
	var questions = []*domains.Question{
		{
			ID:        1,
			ProjectID: 1,
			Title:     "问题1",
			Options:   []string{"A", "B", "C", "D"},
			Answers:   []int{0},
			CreatedAt: time.Now(),
		},
		{
			ID:        2,
			ProjectID: 1,
			Title:     "问题2",
			Options:   []string{"A", "B", "C", "D"},
			Answers:   []int{0},
			CreatedAt: time.Now(),
		},
		{
			ID:        3,
			ProjectID: 1,
			Title:     "问题3",
			Options:   []string{"A", "B", "C", "D"},
			Answers:   []int{1},
			CreatedAt: time.Now(),
		},
		{
			ID:        4,
			ProjectID: 1,
			Title:     "问题4",
			Options:   []string{"A", "B", "C", "D"},
			Answers:   []int{3},
			CreatedAt: time.Now(),
		},
		{
			ID:        5,
			ProjectID: 1,
			Title:     "问题5",
			Options:   []string{"A", "B", "C", "D"},
			Answers:   []int{0, 3},
			CreatedAt: time.Now(),
		},
	}
	exam := domains.NewExam(1, 1, 49, questions)
	paper := map[int64][]int{
		1: {0},
		2: {0},
		3: {1},
		4: {3},
		5: {0, 0},
	}
	err := exam.SaveExam(paper)
	require.Nil(t, err)

	paper[4] = []int{2}
	err = exam.SubmitExam(paper)
	require.Nil(t, err)

	fmt.Printf("%+v", exam)
}
