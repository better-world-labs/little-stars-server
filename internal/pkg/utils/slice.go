package utils

func Distinct[T comparable](s []T) []T {
	l := len(s)
	set := make(map[T]struct{}, l)
	res := make([]T, 0, l)

	for _, e := range s {
		if _, exists := set[e]; !exists {
			res = append(res, e)
			set[e] = struct{}{}
		}
	}

	return res
}

func Map[IN, OUT any](in []IN, mapper func(i IN) OUT) []OUT {
	out := make([]OUT, 0)

	for _, in := range in {
		out = append(out, mapper(in))
	}

	return out
}

func ToMap[IN, OutKey, OutValue comparable](in []IN, mapper func(in IN) (OutKey, OutValue)) map[OutKey]OutValue {
	out := make(map[OutKey]OutValue, 0)

	for _, in := range in {
		key, value := mapper(in)
		out[key] = value
	}

	return out
}

func Flatmap[T any](s [][]T) []T {
	res := make([]T, 0, len(s))

	for i := 0; i < len(s); i++ {
		for j := 0; j < len(s[i]); j++ {
			res = append(res, s[i][j])
		}
	}

	return res
}
