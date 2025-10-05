package main

type Store struct {
	Values map[string]interface{}
}

type StoreInterface interface {
	GetValue(key string) interface{}
	SetValue(key string, value interface{})
	DeleteValue(key string) bool
}

func (s *Store) GetValue(key string) interface{} {
	return s.Values[key]
}

func (s *Store) SetValue(key string, value interface{}) {
	if s.Values == nil {
		s.Values = make(map[string]interface{})
	}
	s.Values[key] = value
}

func (s *Store) DeleteValue(key string) bool {
	if _, exists := s.Values[key]; exists {
		delete(s.Values, key)
		return true
	}
	return false
}
