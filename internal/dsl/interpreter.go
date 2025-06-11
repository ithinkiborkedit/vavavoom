package dsl

import (
	"fmt"
	"nelly/internal/dslcore"
	"strings"
)

// func (p *Program) Execute(state *dslcore.ExecutionState) error {
// 	if p.Repo != nil {
// 		if err := p.Rep.Execute(state); err != nil {
// 			return fmt.Errorf("header execution failed: %w", err)
// 		}
// 	}
// 	for _, stmt := range p.Statements {
// 		if err := stmt.Execute(state); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

func (p *Program) Execute(state *dslcore.ExecutionState) error {
	if p.Repo == nil {
		return fmt.Errorf("no repo specified")
	}
	if p.Branch == nil {
		return fmt.Errorf("no branch specified")
	}
	repo := strings.Trim(p.Repo.URL, `"`)
	branch := strings.Trim(p.Branch.Name, `"`)

	run := func(name string, opts map[string]string) error {
		cmd, ok := Commands[name]
		if !ok {
			return fmt.Errorf("command not found: %s", name)
		}
		return cmd(state, opts)
	}

	if strings.HasPrefix(repo, "http://") ||
		strings.HasPrefix(repo, "https://") ||
		strings.HasPrefix(repo, "git@") {
		if err := run("clone", map[string]string{"url": repo}); err != nil {
			return fmt.Errorf("git clone failed: %w", err)
		}
		if err := run("checkout", map[string]string{"name": branch}); err != nil {
			return fmt.Errorf("git checkout failed: %w", err)
		}
	} else {
		if err := run("init", map[string]string{"directory": repo}); err != nil {
			return fmt.Errorf("git init failed: %w", err)
		}
		if err := run("createBranch", map[string]string{"name": branch}); err != nil {
			return fmt.Errorf("git create branch failed: %w", err)
		}
		if err := run("track", map[string]string{"name": branch}); err != nil {
			return fmt.Errorf("git track branch failed: %w", err)
		}
		if err := run("push", map[string]string{}); err != nil {
			return fmt.Errorf("git push failed: %w", err)
		}
	}

	for _, stmt := range p.Statements {
		if err := stmt.Execute(state); err != nil {
			return err
		}
	}
	return nil
}

// func isRemoteURL(s string) bool {
// 	return len(s) > 0 &&
// 		(strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://") || strings.HasPrefix(s, "git@"))
// }

// func (h *Header) Execute(state *dslcore.ExecutionState) error {
// 	if isRemoteURL(h.Repo) {
// 		cmds := []*CommandStmt{
// 			{Name: "clone", Options: []*CommandOption{
// 				{Dot: ".", Name: "url", Value: &Expr{String: &h.Repo}},
// 			}},
// 			{Name: "checkout", Options: []*CommandOption{
// 				{Dot: ".", Name: "name", Value: &Expr{String: &h.Branch}},
// 			}},
// 		}
// 		for _, c := range cmds {

// 			if err := c.Execute(state); err != nil {
// 				return err
// 			}
// 		}
// 	} else {
// 		cmds := []*CommandStmt{
// 			{Name: "init", Options: []*CommandOption{
// 				{Dot: ".", Name: "directory", Value: &Expr{String: &h.Repo}},
// 			}},
// 			{Name: "createBranch", Options: []*CommandOption{
// 				{Dot: ".", Name: "name", Value: &Expr{String: &h.Branch}},
// 			}},
// 			{Name: "push", Options: nil},
// 			{Name: "track", Options: []*CommandOption{
// 				{Dot: ".", Name: "name", Value: &Expr{String: &h.Branch}},
// 			}},
// 		}
// 		for _, c := range cmds {
// 			if err := c.Execute(state); err != nil {
// 				return err
// 			}
// 		}
// 	}
// 	return nil
// }

func (s *Statement) Execute(state *dslcore.ExecutionState) error {
	switch {
	case s.Let != nil:
		return s.Let.Execute(state)
	case s.Command != nil:
		return s.Command.Execute(state)
	case s.For != nil:
		return s.For.Execute(state)
	case s.If != nil:
		return s.If.Execute(state)
	case s.Expr != nil:
		return s.Expr.Execute(state)
	default:
		return nil
	}
}

func (l *LetStmt) Execute(state *dslcore.ExecutionState) error {
	val, err := l.Value.Eval(state)
	if err != nil {
		return err
	}
	if state.Vars == nil {
		state.Vars = map[string]interface{}{}
	}
	state.Vars[l.Name] = val
	return nil
}

func (c *CommandStmt) Execute(state *dslcore.ExecutionState) error {
	handler, ok := Commands[c.Name]
	if !ok {
		return fmt.Errorf("unkown command: %s", c.Name)
	}
	opts := map[string]string{}
	for _, opt := range c.Options {
		if opt.Value != nil {
			val, err := opt.Value.Eval(state)
			if err != nil {
				return fmt.Errorf("invalid value for option %s: %v", opt.Name, err)
			}
			opts[opt.Name] = fmt.Sprintf("%v", val)
		} else {
			opts[opt.Name] = "true"
		}
	}
	return handler(state, opts)
}

func (f *ForStmt) Execute(state *dslcore.ExecutionState) error {
	val, err := f.Range.Eval(state)
	if err != nil {
		return err
	}
	arr, ok := val.([]interface{})
	if !ok {
		return fmt.Errorf("for: range is not iterable")
	}
	for _, item := range arr {
		state.Vars[f.Var] = item
		for _, stmt := range f.Body {
			if err := stmt.Execute(state); err != nil {
				return err
			}
		}
	}
	return nil
}

func (i *IfStmt) Execute(state *dslcore.ExecutionState) error {
	cond, err := i.Cond.Eval(state)
	if err != nil {
		return err
	}
	isTrue, ok := cond.(bool)
	if !ok {
		return fmt.Errorf("if: condition is not boolean")
	}
	if isTrue {
		for _, stmt := range i.Then {
			if err := stmt.Execute(state); err != nil {
				return err
			}
		}
	} else if i.Else != nil {
		for _, stmt := range i.Else.Body {
			if err := stmt.Execute(state); err != nil {
				return err
			}
		}
	}
	return nil
}

func (e *ExprStmt) Execute(state *dslcore.ExecutionState) error {
	_, err := e.Expr.Eval(state)
	return err
}
