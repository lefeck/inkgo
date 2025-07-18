package set

import "sort"

type Empty struct {
}

type String map[string]Empty

func NewString(items ...string) String {
	ss := String{}
	return ss.Insert(items...)
}

func (s String) Insert(items ...string) String {
	for _, item := range items {
		s[item] = Empty{}
	}
	return s
}

func (s String) Delete(items ...string) String {
	for _, item := range items {
		delete(s, item)
	}
	return s
}

func (s String) Has(item string) bool {
	_, contained := s[item]
	return contained
}

func (s String) HasAll(items ...string) bool {
	for _, item := range items {
		if !s.Has(item) {
			return false
		}
	}
	return true
}

func (s String) HasAny(items ...string) bool {
	for _, item := range items {
		if s.Has(item) {
			return true
		}
	}
	return false
}

func (s String) Slice() []string {
	slice := make([]string, len(s))
	for item := range s {
		slice = append(slice, item)
	}
	sort.Strings(slice)
	return slice
}
