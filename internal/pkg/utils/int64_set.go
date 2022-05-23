package utils

// Int64Set int64 set
type Int64Set map[int64]struct{}

func NewInt64Set() Int64Set {
	set := make(Int64Set)
	return set
}

func (i Int64Set) Set(v int64) {
	i[v] = struct{}{}
}

func (i Int64Set) ToSlice() (slice []int64) {
	for k := range i {
		slice = append(slice, k)
	}
	return
}

func (i Int64Set) Contains(item int64) bool {
	_, exists := i[item]
	return exists
}

func (i Int64Set) Remove(item int64) {
	delete(i, item)
}

func (i Int64Set) RemoveAll(items []int64) {
	for _, e := range items {
		i.Remove(e)
	}
}

func (i Int64Set) Add(item int64) {
	i[item] = struct{}{}
}

func (i Int64Set) AddAll(items []int64) {
	for _, e := range items {
		i.Add(e)
	}
}
