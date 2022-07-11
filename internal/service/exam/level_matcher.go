package exam

func MatchLevel(last7Exams []*ExamDo) int {
	var count, avgLatest3, avgLatest5, avgTotal, sumLatest3, sumLatest5, sumTotal, level int
	for _, e := range last7Exams {
		count++
		sumTotal += e.Score
		avgTotal = sumTotal / count
		if count <= 3 {
			sumLatest3 += e.Score
			avgLatest3 = sumLatest3 / 3
		}

		if count <= 5 {
			sumLatest5 += e.Score
			avgLatest5 = sumLatest5 / 5
		}

		// 白银
		if count >= 1 && avgTotal >= 80 {
			level = 1
		}

		// 黄金
		if count >= 3 && avgLatest3 >= 85 {
			level = 2
		}

		// 铂金
		if count >= 5 && avgLatest3 >= 85 {
			level = 3
		}

		// 钻石
		if count >= 7 && avgLatest5 >= 85 {
			level = 4
		}

		// 星耀
		if count >= 7 && avgLatest5 >= 90 {
			level = 5
		}
	}

	return level
}
