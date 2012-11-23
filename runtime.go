package main

import "fmt"
import "go/build"
import "strconv"

func DumpBuildInfo() {
  fmt.Printf("Default Build Context: %#v\n", build.Default)
}

func init() {
  DumpBuildInfo()
}

type Any interface{}

type Pobj interface {
	Bool() bool
	Int64() int64
	String() string
	Xstr() Pstr
	Xint() Pint
	Xadd(o Pobj) Pobj
	Xsub(o Pobj) Pobj
	Xmul(o Pobj) Pobj
	Xmod(o Pobj) Pobj
	Xdiv(o Pobj) Pobj
	Xcmp(o Pobj) Pint
	Xlt(o Pobj) Pbool
	Xle(o Pobj) Pbool
	Xeq(o Pobj) Pbool
	Xne(o Pobj) Pbool
	Xgt(o Pobj) Pbool
	Xge(o Pobj) Pbool
}

// Pstr
type Pstr string

func (s Pstr) Bool() bool {
	return len(string(s)) > 0
}
func (s Pstr) Int64() int64 {
	z, e := strconv.ParseInt(string(s), 10, 64)
	if e != nil { panic("string not a valid int: <" + s + ">")}
	return z
}
func (s Pstr) String() string {
	return string(s)
}

func (s Pstr) Xstr() Pstr {
	return s
}

func (s Pstr) Xint() Pint {
	return Pint(s.Int64())
}

func (s Pstr) Xadd(o Pobj) Pobj {
	return Pstr( string(s) + o.String())
}

func (s Pstr) Xsub(o Pobj) Pobj {
	panic("Subtraction unsupported for strings.")
}

func (s Pstr) Xmul(o Pobj) Pobj {
	panic("Multiplication unsupported for strings.")
}

func (s Pstr) Xmod(o Pobj) Pobj {
	panic("Mod unsupported for strings.")
}

func (s Pstr) Xdiv(o Pobj) Pobj {
	panic("Division unsupported for strings.")
}

func (s Pstr) Xcmp(o Pobj) Pint {
	if string(s) > string(o.(Pstr)) {
		return 1
	}

	if string(s) < string(o.(Pstr)) {
		return -1
	}

	return 0
}

func (s Pstr) Xlt(o Pobj) Pbool {
	if s.Xcmp(o) < 0 {
		return true
	}

	return false
}

func (s Pstr) Xle(o Pobj) Pbool {
	if s.Xcmp(o) <= 0 {
		return true
	}

	return false
}

func (s Pstr) Xeq(o Pobj) Pbool {
	if s.Xcmp(o) == 0 {
		return true
	}

	return false
}

func (s Pstr) Xne(o Pobj) Pbool {
	if s.Xcmp(o) != 0 {
		return true
	}

	return false
}

func (s Pstr) Xgt(o Pobj) Pbool {
	if s.Xcmp(o) > 0 {
		return true
	}

	return false
}

func (s Pstr) Xge(o Pobj) Pbool {
	if s.Xcmp(o) >= 0 {
		return true
	}

	return false
}

// ~~~Pstr

// Pint
type Pint int64

func (i Pint) Bool() bool {
	return int(i) != 0
}
func (i Pint) Int64() int64 {
	return int64(i)
}

func (i Pint) String() string {
	return fmt.Sprintf("%d", int64(i))
}

func (i Pint) Xstr() Pstr {
	return Pstr(i.String())
}

func (i Pint) Xint() Pint {
	return i
}

func (i Pint) Xadd(o Pobj) Pobj {
	return Pint( i + o.(Pint) )
}

func (i Pint) Xsub(o Pobj) Pobj {
	return Pint( i - o.(Pint) )
}

func (i Pint) Xmul(o Pobj) Pobj {
	return Pint( i * o.(Pint) )
}

func (i Pint) Xmod(o Pobj) Pobj {
	return Pint( i % o.(Pint) )
}

func (i Pint) Xdiv(o Pobj) Pobj {
	return Pint( i / o.(Pint) )
}

func (i Pint) Xcmp(o Pobj) Pint {
	if i > o.(Pint) {
		return 1
	}

	if i < o.(Pint) {
		return -1
	}

	return 0
}

func (i Pint) Xlt(o Pobj) Pbool {
	if i.Xcmp(o) < 0 {
		return true
	}

	return false
}

func (i Pint) Xle(o Pobj) Pbool {
	if i.Xcmp(o) <= 0 {
		return true
	}

	return false
}

func (i Pint) Xeq(o Pobj) Pbool {
	if i.Xcmp(o) == 0 {
		return true
	}

	return false
}

func (i Pint) Xne(o Pobj) Pbool {
	if i.Xcmp(o) != 0 {
		return true
	}

	return false
}

func (i Pint) Xgt(o Pobj) Pbool {
	if i.Xcmp(o) > 0 {
		return true
	}

	return false
}

func (i Pint) Xge(o Pobj) Pbool {
	if i.Xcmp(o) >= 0 {
		return true
	}

	return false
}

// ~~~Pint

// Pbool
type Pbool bool

func (b Pbool) Bool() bool {
	return bool(b)
}
func (b Pbool) Int64() int64 {
	if bool(b) { return 1 }
	return 0
}

func (b Pbool) String() string {
	return fmt.Sprintf("%v", bool(b))
}

func (b Pbool) Xstr() Pstr {
	return Pstr(b.String())
}

func (b Pbool) Xint() Pint {
	if bool(b) {
		return 1
	}

	return 0
}

func (b Pbool) Xadd(o Pobj) Pobj {
	return Pint( b.Xint() + o.(Pbool).Xint())
}

func (b Pbool) Xsub(o Pobj) Pobj {
	return Pint( b.Xint() - o.(Pbool).Xint() )
}

func (b Pbool) Xmul(o Pobj) Pobj {
	return Pint( b.Xint() * o.(Pbool).Xint() )
}

func (b Pbool) Xmod(o Pobj) Pobj {
	return Pint( b.Xint() % o.(Pbool).Xint() )
}

func (b Pbool) Xdiv(o Pobj) Pobj {
	return Pint( b.Xint() / o.(Pbool).Xint() )
}

func (b Pbool) Xcmp(o Pobj) Pint {
	if b == o.(Pbool) {
		return 0
	}

	if b {
		return 1
	}

	return -1
}

func (b Pbool) Xlt(o Pobj) Pbool {
	if b.Xcmp(o) < 0 {
		return true
	}

	return false
}

func (b Pbool) Xle(o Pobj) Pbool {
	if b.Xcmp(o) <= 0 {
		return true
	}

	return false
}

func (b Pbool) Xeq(o Pobj) Pbool {
	if b.Xcmp(o) == 0 {
		return true
	}

	return false
}

func (b Pbool) Xne(o Pobj) Pbool {
	if b.Xcmp(o) != 0 {
		return true
	}

	return false
}

func (b Pbool) Xgt(o Pobj) Pbool {
	if b.Xcmp(o) > 0 {
		return true
	}

	return false
}

func (b Pbool) Xge(o Pobj) Pbool {
	if b.Xcmp(o) >= 0 {
		return true
	}

	return false
}

// ~~~Pbool
