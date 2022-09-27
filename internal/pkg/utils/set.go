package utils

type Set[K comparable] map[K]struct{}

func NewSet[K comparable]() Set[K] {
	set := make(Set[K])
	return set
}

func (i Set[K]) Add(v K) {
	i[v] = struct{}{}
}

func (i Set[K]) ToSlice() (slice []K) {
	for k := range i {
		slice = append(slice, k)
	}
	return
}

func (i Set[K]) Contains(item K) bool {
	_, exists := i[item]
	return exists
}

func (i Set[K]) Remove(item K) {
	delete(i, item)
}

func (i Set[K]) RemoveAll(items []K) {
	for _, e := range items {
		i.Remove(e)
	}
}

func (i Set[K]) AddAll(items []K) {
	for _, e := range items {
		i.Add(e)
	}
}
