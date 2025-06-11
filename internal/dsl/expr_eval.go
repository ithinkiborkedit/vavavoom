package dsl

import (
	"errors"
	"fmt"
	"nelly/internal/dslcore"
	"strconv"
	"strings"
)

func (e *Expr) Eval(state *dslcore.ExecutionState) (interface{}, error) {
	switch {
	case e.Ident != nil:
		if val, ok := state.Vars[*e.Ident]; ok {
			return val, nil
		}
		switch *e.Ident {
		case "true":
			return true, nil
		case "false":
			return false, nil
		default:
			return nil, fmt.Errorf("undefined variable: %s", *e.Ident)
		}
	case e.String != nil:
		return strings.Trim(*e.String, `"`), nil
	case e.Number != nil:
		return *e.Number, nil
	case e.Bool != nil:
		return e.Bool.Value == "true", nil
	case e.Semver != nil:
		return e.Semver.String(), nil
	case e.Array != nil:
		var result []interface{}
		for _, item := range e.Array {
			v, err := item.Eval(state)
			if err != nil {
				return nil, err
			}
			result = append(result, v)
		}
		return result, nil
	case e.Call != nil:
		return e.Binary.Left.Eval(state)
	default:
		return nil, errors.New("invalid or unsupported expression")
	}
}

func (c *CallExpr) Eval(state *dslcore.ExecutionState) (interface{}, error) {
	switch c.Name {
	case "len":
		if len(c.Args) != 1 {
			return nil, errors.New("len expects one argument")
		}
		val, err := c.Args[0].Eval(state)
		if err != nil {
			return nil, err
		}
		switch v := val.(type) {
		case string:
			return float64(len(v)), nil
		case []interface{}:
			return float64(len(v)), nil
		default:
			return nil, fmt.Errorf("len: unsupported type %T", v)
		}
	case "print":
		var out []string
		for _, arg := range c.Args {
			val, err := arg.Eval(state)
			if err != nil {
				return nil, err
			}
			out = append(out, fmt.Sprint(val))
		}
		fmt.Println(strings.Join(out, " "))
		return nil, nil
	default:
		return nil, fmt.Errorf("unsuppported function: %s", c.Name)
	}
}

func asFloat(x interface{}) (float64, bool) {
	switch v := x.(type) {
	case float64:
		return v, true
	case int:
		return float64(v), true
	case string:
		f, err := strconv.ParseFloat(v, 64)
		return f, err == nil
	default:
		return 0, false
	}
}

func equals(a, b interface{}) bool {
	switch av := a.(type) {
	case float64:
		bv, ok := asFloat(b)
		return ok && av == bv
	case string:
		bs, ok := b.(string)
		return ok && av == bs
	case bool:
		bb, ok := b.(bool)
		return ok && av == bb
	default:
		return false
	}
}

func (b *BinaryExpr) Eval(state *dslcore.ExecutionState) (interface{}, error) {
	left, err := b.Left.Eval(state)
	if err != nil {
		return nil, err
	}
	right, err := b.Right.Eval(state)
	if err != nil {
		return nil, err
	}

	switch b.Operator {
	case "==":
		return equals(left, right), nil
	case "!=":
		return !equals(left, right), nil
	case "<", "<=", ">", ">=":
		lf, lok := asFloat(left)
		rf, rok := asFloat(right)
		if !lok || !rok {
			return nil, fmt.Errorf("comparison operands must be numbers: %v %v", left, right)
		}
		switch b.Operator {
		case "<":
			return lf < rf, nil
		case "<=":
			return lf <= rf, nil
		case ">":
			return lf > rf, nil
		case ">=":
			return lf >= rf, nil
		}
	case "+", "-", "*", "/":
		lf, lok := asFloat(left)
		rf, rok := asFloat(right)
		if !lok || !rok {
			return nil, fmt.Errorf("arithamtic operands must be numbers: %v %v", left, right)
		}
		switch b.Operator {
		case "+":
			return lf + rf, nil
		case "-":
			return lf - rf, nil
		case "*":
			return lf * rf, nil
		case "/":
			if rf == 0 {
				return nil, errors.New("divide by zero")
			}
			return lf / rf, nil
		}
	case "&&", "||":
		lb, lok := left.(bool)
		rb, rok := right.(bool)
		if !lok || !rok {
			return nil, fmt.Errorf("logical operands must be booleans: %v %v", left, right)
		}
		switch b.Operator {
		case "&&":
			return lb && rb, nil
		case "||":
			return lb && rb, nil
		}

	}
	return nil, fmt.Errorf("unsupported binary operator: %s", b.Operator)
}
