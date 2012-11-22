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

type PObj interface {
	Xstr() PStr
	Xint() PInt
	Xadd(o PObj) PObj
	Xsub(o PObj) PObj
	Xmul(o PObj) PObj
	Xmod(o PObj) PObj
	Xdiv(o PObj) PObj
	Xcmp(o PObj) PInt
	Xlt(o PObj) PBool
	Xle(o PObj) PBool
	Xeq(o PObj) PBool
	Xne(o PObj) PBool
	Xgt(o PObj) PBool
	Xge(o PObj) PBool
}

// PStr
type PStr string

func (s PStr) Xstr() PStr {
	return s
}

func (s PStr) Xint() PInt {
	panic("Can't convert string to integer.")
}

func (s PStr) Xadd(o PObj) PObj {
	return PStr( string(s) + string(o.(PStr)) )
}

func (s PStr) Xsub(o PObj) PObj {
	panic("Subtraction unsupported for strings.")
}

func (s PStr) Xmul(o PObj) PObj {
	panic("Multiplication unsupported for strings.")
}

func (s PStr) Xmod(o PObj) PObj {
	panic("Mod unsupported for strings.")
}

func (s PStr) Xdiv(o PObj) PObj {
	panic("Division unsupported for strings.")
}

func (s PStr) Xcmp(o PObj) PInt {
	if string(s) > string(o.(PStr)) {
		return 1
	}

	if string(s) < string(o.(PStr)) {
		return -1
	}

	return 0
}

func (s PStr) Xlt(o PObj) PBool {
	if s.Xcmp(o) < -1 {
		return true
	}

	return false
}

func (s PStr) Xle(o PObj) PBool {
	if s.Xcmp(o) <= 0 {
		return true
	}

	return false
}

func (s PStr) Xeq(o PObj) PBool {
	if s.Xcmp(o) == 0 {
		return true
	}

	return false
}

func (s PStr) Xne(o PObj) PBool {
	if s.Xcmp(o) != 0 {
		return true
	}

	return false
}

func (s PStr) Xgt(o PObj) PBool {
	if s.Xcmp(o) > 0 {
		return true
	}

	return false
}

func (s PStr) Xge(o PObj) PBool {
	if s.Xcmp(o) >= 0 {
		return true
	}

	return false
}

// ~~~PStr

// PInt
type PInt int

func (i PInt) Xstr() PStr {
	return PStr(fmt.Sprintf("%d", int(i)))
}

func (i PInt) Xint() PInt {
	return i
}

func (i PInt) Xadd(o PObj) PObj {
	return PInt( i + o.(PInt) )
}

func (i PInt) Xsub(o PObj) PObj {
	return PInt( i - o.(PInt) )
}

func (i PInt) Xmul(o PObj) PObj {
	return PInt( i * o.(PInt) )
}

func (i PInt) Xmod(o PObj) PObj {
	return PInt( i % o.(PInt) )
}

func (i PInt) Xdiv(o PObj) PObj {
	return PInt( i / o.(PInt) )
}

func (i PInt) Xcmp(o PObj) PInt {
	if i > o.(PInt) {
		return 1
	}

	if i < o.(PInt) {
		return -1
	}

	return 0
}

func (i PInt) Xlt(o PObj) PBool {
	if i.Xcmp(o) < -1 {
		return true
	}

	return false
}

func (i PInt) Xle(o PObj) PBool {
	if i.Xcmp(o) <= 0 {
		return true
	}

	return false
}

func (i PInt) Xeq(o PObj) PBool {
	if i.Xcmp(o) == 0 {
		return true
	}

	return false
}

func (i PInt) Xne(o PObj) PBool {
	if i.Xcmp(o) != 0 {
		return true
	}

	return false
}

func (i PInt) Xgt(o PObj) PBool {
	if i.Xcmp(o) > 0 {
		return true
	}

	return false
}

func (i PInt) Xge(o PObj) PBool {
	if i.Xcmp(o) >= 0 {
		return true
	}

	return false
}

// ~~~PInt

// PBool
type PBool bool

func (b PBool) Xstr() PStr {
	return PStr(fmt.Sprintf("%v", bool(b)))
}

func (b PBool) Xint() PInt {
	if bool(b) {
		return 1
	}

	return 0
}

func (b PBool) Xadd(o PObj) PObj {
	return PInt( b.Xint() + o.(PBool).Xint())
}

func (b PBool) Xsub(o PObj) PObj {
	return PInt( b.Xint() - o.(PBool).Xint() )
}

func (b PBool) Xmul(o PObj) PObj {
	return PInt( b.Xint() * o.(PBool).Xint() )
}

func (b PBool) Xmod(o PObj) PObj {
	return PInt( b.Xint() % o.(PBool).Xint() )
}

func (b PBool) Xdiv(o PObj) PObj {
	return PInt( b.Xint() / o.(PBool).Xint() )
}

func (b PBool) Xcmp(o PObj) PInt {
	if b == o.(PBool) {
		return 0
	}

	if b {
		return 1
	}

	return -1
}

func (b PBool) Xlt(o PObj) PBool {
	if b.Xcmp(o) < -1 {
		return true
	}

	return false
}

func (b PBool) Xle(o PObj) PBool {
	if b.Xcmp(o) <= 0 {
		return true
	}

	return false
}

func (b PBool) Xeq(o PObj) PBool {
	if b.Xcmp(o) == 0 {
		return true
	}

	return false
}

func (b PBool) Xne(o PObj) PBool {
	if b.Xcmp(o) != 0 {
		return true
	}

	return false
}

func (b PBool) Xgt(o PObj) PBool {
	if b.Xcmp(o) > 0 {
		return true
	}

	return false
}

func (b PBool) Xge(o PObj) PBool {
	if b.Xcmp(o) >= 0 {
		return true
	}

	return false
}

// ~~~PBool



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
