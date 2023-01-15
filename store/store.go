package store

type Store struct {
	Variables map[string]any
}

func New(variables map[string]any) *Store {
	return &Store{
		Variables: variables,
	}
}

func (s *Store) Set(k string, v any) {
	s.Variables[k] = v
}

func (s *Store) Get(k string) (any, bool) {
	v, ok := s.Variables[k]
	return v, ok
}

func (s *Store) List() map[string]any {
	return s.Variables
}
