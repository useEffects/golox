package lox

import (
	"fmt"
	"strconv"
	"strings"
)

type Interpreter[T string] struct{}

func (i *Interpreter[T]) VisitLiteral(expr Literal[T]) T {
	return expr.Value.(T)
}

func (i *Interpreter[T]) VisitGrouping(expr Grouping[T]) T {
	return i.evaluate(expr.Expression)
}

func (i *Interpreter[T]) evaluate(expr Expr[T]) T {
	return expr.accept(i)
}

func (i *Interpreter[T]) VisitUnary(expr Unary[T]) T {
	right := i.evaluate(expr.Right)
	switch expr.Operator.Type {
	case BANG:
		return i.toAny(!i.isTruthy(right)).(T)
	case MINUS:
		i.checkNumberOperand(expr.Operator, right)
		val := i.toAny(right).(float64)
		return i.toAny(-val).(T)
	}
	return i.zeroValue()
}

func (i *Interpreter[T]) isTruthy(value any) bool {
	if value == nil {
		return false
	}
	if v, ok := value.(bool); ok {
		return v
	}
	return true
}

func (i *Interpreter[T]) VisitBinary(expr Binary[T]) T {
	left := i.toAny(i.evaluate(expr.Left))
	right := i.toAny(i.evaluate(expr.Right))

	switch expr.Operator.Type {
	case PLUS:
		if l, ok := left.(float64); ok {
			if r, ok := right.(float64); ok {
				return i.toAny(l + r).(T)
			}
		}
		if l, ok := left.(string); ok {
			if r, ok := right.(string); ok {
				return i.toAny(l + r).(T)
			}
		}
		err := RuntimeError{"Operands must be two numbers or two strings", expr.Operator}
		err.Error()
	case GREATER:
		i.checkNumberOperand(expr.Operator, left, right)
		return i.toAny(left.(float64) > right.(float64)).(T)
	case GREATER_EQUAL:
		i.checkNumberOperand(expr.Operator, left, right)
		return i.toAny(left.(float64) >= right.(float64)).(T)
	case LESS:
		i.checkNumberOperand(expr.Operator, left, right)
		return i.toAny(left.(float64) < right.(float64)).(T)
	case LESS_EQUAL:
		i.checkNumberOperand(expr.Operator, left, right)
		return i.toAny(left.(float64) <= right.(float64)).(T)
	case MINUS:
		i.checkNumberOperand(expr.Operator, left, right)
		leftVal, _ := strconv.ParseFloat(fmt.Sprintf("%v", left), 64)
		rightVal, _ := strconv.ParseFloat(fmt.Sprintf("%v", right), 64)
		return i.toAny(fmt.Sprintf("%f", leftVal-rightVal)).(T)
	case SLASH:
		i.checkNumberOperand(expr.Operator, left, right)
		return i.toAny(left.(float64) / right.(float64)).(T)
	case STAR:
		i.checkNumberOperand(expr.Operator, left, right)
		return i.toAny(left.(float64) * right.(float64)).(T)
	case BANG_EQUAL:
		return i.toAny(!i.isEqual(left, right)).(T)
	case EQUAL_EQUAL:
		return i.toAny(i.isEqual(left, right)).(T)
	}

	return i.zeroValue()
}

func (i *Interpreter[T]) isEqual(a, b any) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return false
	}
	return a == b
}

func (i *Interpreter[T]) checkNumberOperand(operator Token, operands ...any) {
	for _, o := range operands {
		switch v := o.(type) {
		case float64:
		case int:
			_ = float64(v)
		case string:
			if _, err := strconv.ParseFloat(v, 64); err != nil {
				err := RuntimeError{"Operand must be a number", operator}
				err.Error()
			}
		default:
			err := RuntimeError{"Operand must be a number", operator}
			err.Error()
		}
	}
}

func (i *Interpreter[T]) stringify(value any) (string, *RuntimeError) {
	switch v := value.(type) {
	case nil:
		return "nil", nil
	case float64:
		text := fmt.Sprintf("%v", int(v))
		if strings.HasSuffix(text, ".0") {
			return text[:len(text)-2], nil
		}
	case string:
		return v, nil
	default:
		return "", &RuntimeError{"Invalid value type", Token{}}
	}
	return "", nil
}

func (i *Interpreter[T]) Interpret(expression Expr[T]) {
	val, err := i.stringify(i.evaluate(expression))
	if err != nil {
		lox.RuntimeError(err)
	} else {
		fmt.Println(val)
	}
}

func (i *Interpreter[T]) zeroValue() T {
	var zero T
	return zero
}

func (i *Interpreter[T]) toAny(value any) any {
	return value
}
