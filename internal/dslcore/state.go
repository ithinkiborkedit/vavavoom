package dslcore

import "sync"

type ExecutionState struct {
	CurrentRepo   string
	CurrentBranch string
	Vars          map[string]interface{}
	mu            sync.Mutex
}

func NewState() *ExecutionState {
	return &ExecutionState{Vars: make(map[string]interface{})}
}

func (s *ExecutionState) SetVar(name string, value interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Vars[name] = value

}

func (s *ExecutionState) GetVar(name string, value interface{}) (interface{}, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	val, ok := s.Vars[name]
	return val, ok
}

func (s *ExecutionState) SetRepo(repo string) {
	s.CurrentRepo = repo
}

func (s *ExecutionState) SetBranch(branch string) {
	s.CurrentBranch = branch
}

type CommandFunc func(state *ExecutionState, opts map[string]string) error
