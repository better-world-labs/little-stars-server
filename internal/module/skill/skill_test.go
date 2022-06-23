package skill

import (
	_ "aed-api-server/internal/pkg/db"
	"aed-api-server/internal/service"
	"testing"
)

func initData() {
	//db.InitEngine("mysql", "root:1qaz.2wsx@tcp(116.62.220.222:3306)/aed?charset=utf8mb4")
	initExamData()
}

var s = service.NewService(nil, nil)

func Test_SaveExamToDb(t *testing.T) {
	// initData()
}

func Test_checkIfCorrect(t *testing.T) {
	answer := []string{"a", "b"}
	correctAnswer := []string{"a", "b"}
	actual := s.checkCorrect(answer, correctAnswer)
	expected := true
	if actual != expected {
		t.Fail()
	}

	answer1 := []string{"a", "b"}
	correctAnswer1 := []string{"a"}
	actual = s.checkCorrect(answer1, correctAnswer1)
	expected = false
	if actual != expected {
		t.Fail()
	}
}
