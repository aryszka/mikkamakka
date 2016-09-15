package mikkamakka

type tpair struct {
	car, cdr *Val
}

var (
	Nil = &Val{mtype: mnil}

	mixedTypes = &Val{merror, "mixed types"}
)

func listToSlice(l *Val) []*Val {
	var s []*Val
	for {
		if isNil(l) != False {
			break
		}

		s, l = append(s, car(l)), cdr(l)
	}

	return s
}

func Cons(car, cdr *Val) *Val {
	// if car.mtype != cdr.mtype && cdr.mtype != mnil && cdr.mtype != pair {
	// 	panic(mixedTypes)
	// }

	return &Val{pair, &tpair{car, cdr}}
}

func car(p *Val) *Val {
	checkType(p, pair)
	return p.value.(*tpair).car
}

func cdr(p *Val) *Val {
	checkType(p, pair)
	return p.value.(*tpair).cdr
}

func isPair(a *Val) *Val {
	if a.mtype == pair {
		return True
	}

	return False
}

func isNil(a *Val) *Val {
	if a == Nil {
		return True
	}

	return False
}

func list(a ...*Val) *Val {
	l := Nil
	for i := len(a) - 1; i >= 0; i-- {
		l = Cons(a[i], l)
	}

	return l
}

func reverse(l *Val) *Val {
	checkType(l, pair, mnil)

	r := Nil
	for {
		if l == Nil {
			return r
		}

		r = Cons(car(l), r)
		l = cdr(l)
	}
}

func reverseIrregular(l *Val) *Val {
	checkType(l, pair, mnil)

	r := Cons(car(cdr(l)), car(l))
	l = cdr(cdr(l))
	for {
		if l == Nil {
			return r
		}

		r = Cons(car(l), r)
		l = cdr(l)
	}
}

func mappend(left, right *Val) *Val {
	checkType(left, pair, mnil)
	checkType(right, pair, mnil)

	if isNil(left) != False {
		return right
	}

	return Cons(car(left), mappend(cdr(left), right))
}

func List(a ...*Val) *Val {
	av := make([]*Val, len(a))
	for i, ai := range a {
		av[i] = ai
	}

	return list(av...)
}
