package interpreter

import (
	"fmt"
	"time"

	"golox/pkg/ast"
)

type callable interface {
	arity() int
	call(i *interpreter, args []interface{}) interface{}
}

type clock struct{}

func (c *clock) arity() int { return 0 }

func (c *clock) call(i *interpreter, args []interface{}) interface{} {
	return float64(time.Now().UnixMilli() / 1000)
}

func (c clock) String() string {
	return "<native function clock>"
}

type function struct {
	declaration *ast.FunStmt
}

func (f *function) arity() int { return len(f.declaration.Params) }

func (f *function) call(i *interpreter, args []interface{}) (value interface{}) {
	env := &environment{i.global, make(map[string]interface{})}
	for i := 0; i < f.arity(); i++ {
		env.define(f.declaration.Params[i].Lexeme, args[i])
	}

	prev := i.env
	defer func() {
		i.env = prev
		r := recover()
		if err, ok := r.(error); ok {
			panic(err)
		} else {
			value = r
		}
	}()

	i.env = env
	for _, stmt := range f.declaration.Body.Statements {
		stmt.Accept(i)
	}

	return value
}

func (f function) String() string {
	return fmt.Sprintf("<function %s >", f.declaration.Name.Lexeme)
}
