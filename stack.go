package vigo

import "errors"

type stack struct {
	data [64]uint16
	size int
}

func (s *stack) push(v uint16) error {
	if s.size >= 64 {
		return errors.New("full stack")
	}
	s.data[s.size] = v
	s.size += 1
	return nil
}

func (s *stack) pop() (uint16, error) {
	if s.size <= 0 {
		return 0, errors.New("empty stack")
	}

	s.size -= 1
	v := s.data[s.size]

	return v, nil
}
