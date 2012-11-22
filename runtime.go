package main

import "fmt"
import "go/build"

func DumpBuildInfo() {
  fmt.Printf("Default Build Context: %#v", build.Default)
}

func init() {
  DumpBuildInfo()
}

type Any interface{}

type PyObject interface {
	PyAdd(other PyObject) PyObject
	PyStr() PyString
	PyInt() PyInteger
}

// PyString
type PyString string

func (s PyString) PyAdd(other PyObject) PyObject {
	return PyString( string(s) + string(other.(PyString)) )
}

func (s PyString) PyStr() PyString {
	return s
}

func (s PyString) PyInt() PyInteger {
	panic("Can't convert string to integer.")
}

// ~~~PyString

// PyInteger
type PyInteger int

func (i PyInteger) PyAdd(other PyObject) PyObject {
	return PyInteger( int(i) + int(other.(PyInteger)) )
}

func (i PyInteger) PyStr() PyString {
	return PyString(fmt.Sprintf("%d", int(i)))
}

func (i PyInteger) PyInt() PyInteger {
	return i
}

// ~~~PyInteger

// PyBoolean
type PyBoolean bool

func (b PyBoolean) PyAdd(other PyObject) PyObject {
	return PyInteger( b.PyInt() + other.(PyBoolean).PyInt())
}

func (b PyBoolean) PyStr() PyString {
	return PyString(fmt.Sprintf("%v", bool(b)))
}

func (b PyBoolean) PyInt() PyInteger {
	if bool(b) {
		return 1
	}

	return 0
}

// ~~~PyBoolean

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
