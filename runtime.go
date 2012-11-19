package main

import "fmt"

type Any interface{}

func PrintableValue(a Any) string {
	switch t := a.(type) {
	case bool:
		return fmt.Sprintf("%v", t)
	case int:
		return fmt.Sprintf("%d", t)
	case string:
		return t
	}
	panic("Type not handled in PrintableValue")
}

func CompareLt(a Any, b Any) Any {
	return a.(int) < b.(int)
}

func CompareLtE(a Any, b Any) Any {
	return a.(int) <= b.(int)
}

func CompareGt(a Any, b Any) Any {
	return a.(int) > b.(int)
}

func CompareGtE(a Any, b Any) Any {
	return a.(int) >= b.(int)
}

func CompareEq(a Any, b Any) Any {
	return a.(int) == b.(int)
}

func CompareNotEq(a Any, b Any) Any {
	return a.(int) != b.(int)
}

func BinOpAdd(a Any, b Any) Any {
	return a.(int) + b.(int)
}

func BinOpSub(a Any, b Any) Any {
	return a.(int) - b.(int)
}

func BinOpMult(a Any, b Any) Any {
	return a.(int) * b.(int)
}

func BinOpDiv(a Any, b Any) Any {
	return a.(int) / b.(int)
}

func BinOpMod(a Any, b Any) Any {
	return a.(int) % b.(int)
}
