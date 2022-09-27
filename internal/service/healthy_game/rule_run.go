package healthy_game

import "aed-api-server/internal/interfaces/entities"

func runRule(answers []*entities.Answer) (result *entities.Result, score int) {
	c := ctx{Answers: answers}
	for _, r := range rules {
		if rst, ok := r.Condition.Execute(&c).(bool); rst && ok {
			o := op{Operator: OpTotal}
			score = o.Execute(&c).(int)
			return r.Result, score
		}
	}
	return nil, 0
}
