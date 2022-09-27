package healthy_game

import "aed-api-server/internal/interfaces/entities"

func and(ops ...opParam) *op {
	return &op{
		Operator: OpAnd,
		Params:   ops,
	}
}

func or(ops ...opParam) *op {
	return &op{
		Operator: OpOr,
		Params:   ops,
	}
}

func lt(left, right opParam) *op {
	return &op{
		Operator: OpLt,
		Params: []opParam{
			left, right,
		},
	}
}

func eq(left, right opParam) *op {
	return &op{
		Operator: OpEq,
		Params: []opParam{
			left, right,
		},
	}
}

func totalScore() *op {
	return &op{
		Operator: OpTotal,
	}
}

func is(question question, answer answer) *op {
	return &op{
		Operator: OpAnswer,
		Params: []opParam{
			question, answer,
		},
	}
}

var rules = []*rule{
	//高风险
	{
		Condition: or(
			is(1, answer("A")),
			is(2, answer("A")),
		),
		Result: &entities.Result{
			LevelId: 5,
			Level:   "高风险",
			Explain: "已经确诊心源性猝死病因的相关疾病，需要谨遵医嘱严格以药物&治疗形式控制病情，避免导致猝死结果。",
		},
	},

	//中风险
	{
		Condition: and(
			is(1, "B"),
			is(2, "B"),
			is(3, "A"),
		),
		Result: &entities.Result{
			LevelId: 4,
			Level:   "中风险",
			Explain: "未确诊，但是有直接的家族史。猝死相关病因都有一定比例的遗传因素，需要全面体检，早筛查早确诊早治疗。",
		},
	},

	//潜在高风险
	{
		Condition: and(
			is(1, "B"),
			is(2, "B"),
			or(
				is(3, "B"),
				is(3, "C"),
			),

			//900 <= 总分
			or(
				lt(score(900), totalScore()),
				eq(score(900), totalScore()),
			),
		),

		Result: &entities.Result{
			LevelId: 3,
			Level:   "潜在高风险",
			Explain: "目前的生活习惯&生理指标中具有单个严重或多个猝死相关疾病的危险因素，具有强烈的猝死相关疾病导向性，需要即刻采取就医药物&治疗的形式、同时改变不良的生活习惯，控制异常指标，避免发展为猝死高风险。",
		},
	},

	//潜在中风险
	{
		Condition: and(
			is(1, "B"),
			is(2, "B"),
			or(
				is(3, "B"),
				is(3, "C"),
			),

			// 250 < 总分 < 900
			lt(totalScore(), score(900)),
			lt(score(250), totalScore()),
		),
		Result: &entities.Result{
			LevelId: 2,
			Level:   "潜在中风险",
			Explain: "目前的生活习惯&生理指标中具有1-2个猝死相关疾病的危险因素，需要即刻关注并采取相应的医疗或者非医疗手段去除此危险因素，避免进一步发展为潜在高风险。",
		},
	},

	//低风险
	{
		Condition: and(
			is(1, "B"),
			is(2, "B"),
			or(
				is(3, "B"),
				is(3, "C"),
			),

			// 总分 <= 250
			or(
				lt(totalScore(), score(250)),
				eq(totalScore(), score(250)),
			),
		),
		Result: &entities.Result{
			LevelId: 1,
			Level:   "低风险",
			Explain: "生理指标&生活习惯正常且健康，请坚持。",
		},
	},
}
