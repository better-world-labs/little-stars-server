package utils

// StringSet string set
type StringSet map[string]struct{}

func (i StringSet) Set(v string) {
	i[v] = struct{}{}
}

func (i StringSet) ToSlice() (slice []string) {
	for k := range i {
		slice = append(slice, k)
	}
	return
}
