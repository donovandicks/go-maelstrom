package utils

import "golang.org/x/exp/maps"

type Set[T comparable] struct {
	items map[T]struct{}
}

func NewSet[T comparable]() Set[T] {
	return Set[T]{
		items: make(map[T]struct{}),
	}
}

func (s *Set[T]) Add(item T) {
	s.items[item] = struct{}{}
}

func (s *Set[T]) Items() []T {
	return maps.Keys(s.items)
}

func (s *Set[T]) Contains(item T) bool {
	_, ok := s.items[item]
	return ok
}
