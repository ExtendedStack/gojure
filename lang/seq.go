package lang

import (
	"fmt"

	pers "github.com/ExtendedStack/gojure/persistent"
)

type Seq interface {
	First() interface{}
	Rest() Seq
	Cons(x interface{}) Seq
}

func Count(s Seq) int {
	if s == nil {
		return 0
	}
	i := 1
	for s = s.Rest(); s != nil; s = s.Rest() {
		i++
	}
	return i
}

func Format(seq Seq, start string, end string) string {
	s := start
	if seq != nil {
		s += fmt.Sprint(seq.First())
		for seq = seq.Rest(); seq != nil; seq = seq.Rest() {
			s += " " + fmt.Sprint(seq.First())
		}
	}
	s += end
	return s
}

type LazySeq struct {
	val *struct {
		first interface{}
		rest  Seq
	}
	force func() (interface{}, Seq)
}

func (l *LazySeq) f() {
	if l.val == nil {
		l.val = new(struct {
			first interface{}
			rest  Seq
		})
		l.val.first, l.val.rest = l.force()
	}
}

func Lazy(from func() (interface{}, Seq)) Seq {
	return &LazySeq{force: from}
}

func (l *LazySeq) First() interface{} {
	l.f()
	return l.val.first
}

func (l *LazySeq) Rest() Seq {
	l.f()
	return l.val.rest
}

func (l *LazySeq) Cons(x interface{}) Seq {
	return Lazy(func() (interface{}, Seq) {
		return x, l
	})
}

func (l *LazySeq) String() string {
	return Format(l, "(", ")")
}

func Take(n int, seq Seq) Seq {
	if n == 0 || seq == nil {
		return nil
	}
	return Lazy(func() (interface{}, Seq) {
		return seq.First(), Take(n-1, seq.Rest())
	})
}

func Map(f func(interface{}) interface{}, seq Seq) Seq {
	if seq == nil {
		return nil
	}
	return Lazy(func() (interface{}, Seq) {
		return f(seq.First()), Map(f, seq.Rest())
	})
}

type List pers.List

func NewList(items ...interface{}) Seq {
	l := (*List)(pers.NewList(items...))
	if l == nil {
		return nil
	}
	return l
}

func (l *List) First() interface{} {
	return (*pers.List)(l).First()
}

func (l *List) Rest() Seq {
	next := (*pers.List)(l).Rest()
	if next == nil {
		return nil
	}
	seq := List(*next)
	return &seq
}

func (l *List) Cons(x interface{}) Seq {
	seq := List(*(*pers.List)(l).Cons(x))
	return &seq
}

func (l *List) String() string {
	return Format(l, "(", ")")
}

type Vector struct {
	*pers.Vector
	from int
}

func NewVector(items ...interface{}) Seq {
	v := &Vector{pers.NewVector(items...), 0}
	if v.Count() == 0 {
		return nil
	}
	return v
}

func (v *Vector) Count() int {
	return v.Vector.Count() - v.from
}

func (v *Vector) Nth(i int) interface{} {
	return v.Vector.Nth(i + v.from)
}

func (v *Vector) First() interface{} {
	return v.Nth(0)
}

func (v *Vector) Rest() Seq {
	if v.Count() <= 1 {
		return nil
	}
	return &Vector{v.Vector, v.from + 1}
}

func (v *Vector) Cons(x interface{}) Seq {
	newV := &Vector{pers.NewVector(x), 0}
	for i := 0; i < v.Count(); i++ {
		newV.Vector = newV.Vector.Conj(v.Nth(i))
	}
	return newV
}

func (v *Vector) String() string {
	return Format(v, "[", "]")
}
