package main

import . "fmt"
import "os"
import "reflect"
import "strconv"

var _ = os.Open

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
	Xindex(a Pobj) Pobj
	Xslice(a Pobj, b Pobj) Pobj
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

func (me Pstr) Xindex(a Pobj) Pobj {
	s := string(me)
	Printf("Pstr::Xindex %#v %#v\n", s, a)
	i := a.Int64()
	if i < 0 {
	  i += int64(len(s))
	}
	return Pstr(s[i : i+1])  // One char string.
}
func (me Pstr) Xslice(a Pobj, b Pobj) Pobj {
	s := string(me)
	Printf("Pstr::Xslice %#v %#v %#v\n", s, a, b)
	if a != nil && b == nil {
	  return Pstr(s[a.Int64() : ])
	}
	if a == nil && b != nil {
	  return Pstr(s[ : b.Int64()])
	}
	if a != nil && b != nil {
	  return Pstr(s[a.Int64() : b.Int64()])
	}
	return me  // Both nil.
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
	return Sprintf("%d", int64(i))
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

func (me Pint) Xindex(a Pobj) Pobj {
	panic("Pstr cannot index")
}
func (s Pint) Xslice(a Pobj, b Pobj) Pobj {
	panic("Pint cannot slice")
}
// ~~~Pint

// Pbool
type Pbool bool

func (me Pbool) Bool() bool {
	return bool(me)
}
func (me Pbool) Int64() int64 {
	if bool(me) { return 1 }
	return 0
}

func (me Pbool) String() string {
	return Sprintf("%v", bool(me))
}

func (me Pbool) Xstr() Pstr {
	return Pstr(me.String())
}

func (me Pbool) Xint() Pint {
	if bool(me) {
		return 1
	}
	return 0
}

func (me Pbool) Xadd(o Pobj) Pobj {
	return Pint( me.Xint() + o.(Pbool).Xint())
}

func (me Pbool) Xsub(o Pobj) Pobj {
	return Pint( me.Xint() - o.(Pbool).Xint() )
}

func (me Pbool) Xmul(o Pobj) Pobj {
	return Pint( me.Xint() * o.(Pbool).Xint() )
}

func (me Pbool) Xmod(o Pobj) Pobj {
	return Pint( me.Xint() % o.(Pbool).Xint() )
}

func (me Pbool) Xdiv(o Pobj) Pobj {
	return Pint( me.Xint() / o.(Pbool).Xint() )
}

func (me Pbool) Xcmp(o Pobj) Pint {
	if me == o.(Pbool) {
		return 0
	}

	if me {
		return 1
	}

	return -1
}

func (me Pbool) Xlt(o Pobj) Pbool {
	if me.Xcmp(o) < 0 {
		return true
	}

	return false
}

func (me Pbool) Xle(o Pobj) Pbool {
	if me.Xcmp(o) <= 0 {
		return true
	}

	return false
}

func (me Pbool) Xeq(o Pobj) Pbool {
	if me.Xcmp(o) == 0 {
		return true
	}

	return false
}

func (me Pbool) Xne(o Pobj) Pbool {
	if me.Xcmp(o) != 0 {
		return true
	}

	return false
}

func (me Pbool) Xgt(o Pobj) Pbool {
	if me.Xcmp(o) > 0 {
		return true
	}

	return false
}

func (me Pbool) Xge(o Pobj) Pbool {
	if me.Xcmp(o) >= 0 {
		return true
	}

	return false
}

func (me Pbool) Xindex(a Pobj) Pobj {
	panic("Pbool cannot index")
}
func (s Pbool) Xslice(a Pobj, b Pobj) Pobj {
	panic("Pbool cannot slice")
}
// ~~~Pbool

// Pgo
type Pgo struct {
	P	interface{}
}

func NewPgo(a interface{}) Pgo {
	return Pgo{P: a}
}

func (me Pgo) Bool() bool {
	panic("Pgo cannot Bool")
}
func (me Pgo) Int64() int64 {
	panic("Pgo cannot Bool")
}

func (me Pgo) String() string {
	return Sprintf("%v", me)
}

func (me Pgo) Xstr() Pstr {
	return Pstr(me.String())
}

func (me Pgo) Xint() Pint {
	panic("Pgo cannot int")
}

func (me Pgo) Xadd(o Pobj) Pobj {
	panic("Pgo cannot add")
}

func (me Pgo) Xsub(o Pobj) Pobj {
	panic("Pgo cannot sub")
}

func (me Pgo) Xmul(o Pobj) Pobj {
	panic("Pgo cannot mul")
}

func (me Pgo) Xmod(o Pobj) Pobj {
	panic("Pgo cannot mod")
}

func (me Pgo) Xdiv(o Pobj) Pobj {
	panic("Pgo cannot div")
}

func (me Pgo) Xcmp(it Pobj) Pint {
	a := reflect.ValueOf(me.P).Addr().Uint()
	switch x := it.(type) {
	case Pgo:
	    	b := reflect.ValueOf(x.P).Addr().Uint()
		if a < b {
			return -1
		}
		if a > b {
			return 1
		}
	default:
		panic("Pgo cannot cmp non-Pgo")
	}
	return 0
}

func (me Pgo) Xlt(o Pobj) Pbool {
	return me.Xcmp(o) < 0
}

func (me Pgo) Xle(o Pobj) Pbool {
	if me.Xcmp(o) <= 0 {
		return true
	}
	return false
}

func (me Pgo) Xeq(o Pobj) Pbool {
	if me.Xcmp(o) == 0 {
		return true
	}
	return false
}

func (me Pgo) Xne(o Pobj) Pbool {
	if me.Xcmp(o) != 0 {
		return true
	}
	return false
}

func (me Pgo) Xgt(o Pobj) Pbool {
	if me.Xcmp(o) > 0 {
		return true
	}
	return false
}

func (me Pgo) Xge(o Pobj) Pbool {
	if me.Xcmp(o) >= 0 {
		return true
	}
	return false
}

func (me Pgo) Xindex(a Pobj) Pobj {
	panic("Pgo cannot index")
}
func (s Pgo) Xslice(a Pobj, b Pobj) Pobj {
	panic("Pgo cannot slice")
}
// ~~~Pbool
