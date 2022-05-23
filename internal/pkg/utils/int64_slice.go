package utils

func DistinctInt(s []int64) []int64 {
	l := len(s)
	set := make(map[int64]struct{}, l)
	res := make([]int64, 0, l)

	for _, e := range s {
		if _, exists := set[e]; !exists {
			res = append(res, e)
			set[e] = struct{}{}
		}
	}

	return res
}

func FlatmapInt(s [][]int64) []int64 {
	res := make([]int64, 0, len(s))

	for i := 0; i < len(s); i++ {
		for j := 0; j < len(s[i]); j++ {
			res = append(res, s[i][j])
		}
	}

	return res
}
