package gotype

import "sync"

type Set struct {
	m map[interface{}]struct{}
	sync.RWMutex
}

func NewSet() *Set {
	return &Set{
		m: map[interface{}]struct{}{},
	}
}

func (s *Set) Add(item interface{}) {
	s.Lock()
	defer s.Unlock()
	s.m[item] = struct{}{}
}
func (s *Set) Remove(item interface{}) {
	s.Lock()
	defer s.Unlock()
	delete(s.m, item)
}
func (s *Set) Contains(item interface{}) bool {
	s.RLock()
	defer s.RUnlock()
	_, ok := s.m[item]
	return ok
}
func (s *Set) Len() int {
	return len(s.List())
}
func (s *Set) Clear() {
	s.RLock()
	defer s.Unlock()
	s.m = map[interface{}]struct{}{}
}
func (s *Set) IsEmpty() bool {
	return s.Len() == 0
}
func (s *Set) List() []interface{} {
	s.RLock()
	defer s.RUnlock()
	var list []interface{}
	for item := range s.m {
		list = append(list, item)
	}
	return list
}
