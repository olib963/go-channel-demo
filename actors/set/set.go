package set

type Set[T comparable] map[T]struct{}

func Of[T comparable](items ...T) Set[T] {
	set := make(Set[T])
	for _, item := range items {
		set.Add(item)
	}
	return set
}

func (set Set[T]) Add(item T) {
	set[item] = struct{}{}
}

func (set Set[T]) Remove(item T) {
	delete(set, item)
}

func (set Set[T]) Contains(item T) bool {
	_, exists := set[item]
	return exists
}
