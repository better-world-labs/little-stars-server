package healthy_game

import (
	"aed-api-server/internal/interfaces/entities"
	log "github.com/sirupsen/logrus"
)

type Operator string

const (
	OpAnd    Operator = "$and"
	OpOr     Operator = "$or"
	OpTotal  Operator = "$score"
	OpAnswer Operator = "$answer"
	OpLt     Operator = "$lt"
	OpEq     Operator = "$eq"
)

type ctx struct {
	Answers    []*entities.Answer  //答案
	AnswerMap  map[question]answer //答案
	TotalScore int                 //总分

	ScoreCached bool
}

type op struct {
	Operator Operator  //运算符号
	Params   []opParam //运算参数
}

type rule struct {
	Condition *op
	Result    *entities.Result
}

type opParam interface {
	Execute(*ctx) interface{}
}

type score int

func (s score) Execute(ctx *ctx) interface{} {
	return int(s)
}

type answer string

func (a answer) Execute(ctx *ctx) interface{} {
	return a
}

type question int

func (q question) Execute(ctx *ctx) interface{} {
	return q
}

func (o *op) Execute(ctx *ctx) interface{} {
	switch o.Operator {
	case OpEq:
		if len(o.Params) < 2 {
			log.Warn("healthy game: `OpAnswer` params less than 2")
			return false
		}

		left := o.Params[0].Execute(ctx)
		right := o.Params[1].Execute(ctx)
		leftInt, ok := left.(int)
		if !ok {
			log.Warn("healthy game: `OpLt` first param executed result must int")
		}
		rightInt, ok := right.(int)
		if !ok {
			log.Warn("healthy game: `OpLt` second param executed result must int")
		}
		return leftInt == rightInt
	case OpLt:
		if len(o.Params) < 2 {
			log.Warn("healthy game: `OpAnswer` params less than 2")
			return false
		}

		left := o.Params[0].Execute(ctx)
		right := o.Params[1].Execute(ctx)
		leftInt, ok := left.(int)
		if !ok {
			log.Warn("healthy game: `OpLt` first param executed result must int")
		}
		rightInt, ok := right.(int)
		if !ok {
			log.Warn("healthy game: `OpLt` second param executed result must int")
		}
		return leftInt < rightInt

	case OpOr:
		for _, op := range o.Params {
			execute := op.Execute(ctx)
			b, ok := execute.(bool)
			if !ok {
				log.Warn("healthy game: `OpOr` params executed result must be bool")
				return false
			}
			if b {
				return true
			}
		}
		return false
	case OpAnd:
		for _, op := range o.Params {
			execute := op.Execute(ctx)
			b, ok := execute.(bool)
			if !ok {
				log.Warn("healthy game: `OpAnd` params executed result must be bool")
				return false
			}
			if !b {
				return false
			}
		}
		return true
	case OpTotal:
		if !ctx.ScoreCached {
			o.calculateTotalScore(ctx)
		}
		return ctx.TotalScore
	case OpAnswer:
		if ctx.AnswerMap == nil {
			o.buildAnswerMap(ctx)
		}
		if len(o.Params) < 2 {
			log.Warn("healthy game: `OpAnswer` params less than 2")
			return false
		}
		q, ok := o.Params[0].(question)
		if !ok {
			log.Warn("healthy game: `OpAnswer` first param must be question id")
			return false
		}

		a, ok := o.Params[1].(answer)
		if !ok {
			log.Warn("healthy game: `OpAnswer` second param must be answer index")
			return false
		}
		return ctx.AnswerMap[q] == a
	default:
		return 0
	}
}

func (o *op) buildAnswerMap(ctx *ctx) {
	ctx.AnswerMap = make(map[question]answer)
	for _, a := range ctx.Answers {
		ctx.AnswerMap[question(a.QuestionId)] = answer(a.Select)
	}
}

func (o *op) calculateTotalScore(ctx *ctx) {
	if ctx.AnswerMap == nil {
		o.buildAnswerMap(ctx)
	}

	ctxScore := 0
	for _, q := range questions {
		a := ctx.AnswerMap[question(q.Id)]
		for _, option := range q.Options {
			if answer(option.Index) == a {
				ctxScore += option.Score
				break
			}
		}
	}
	ctx.TotalScore = ctxScore
	ctx.ScoreCached = true
}
