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

func sliceToList(s []*Val) *Val {
	l := Nil
	for _, si := range s {
		l = cons(si, l)
	}

	return reverse(l)
}

func cons(car, cdr *Val) *Val {
	// if car.mtype != cdr.mtype && cdr.mtype != mnil && cdr.mtype != pair {
	// 	panic(mixedTypes)
	// }

	return &Val{pair, &tpair{car, cdr}}
}

func bcons(a []*Val) *Val {
	return cons(a[0], a[1])
}

func car(p *Val) *Val {
	checkType(p, pair)
	return p.value.(*tpair).car
}

func bcar(a []*Val) *Val {
	return car(a[0])
}

func cdr(p *Val) *Val {
	checkType(p, pair)
	return p.value.(*tpair).cdr
}

func bcdr(a []*Val) *Val {
	return cdr(a[0])
}

func isPair(a *Val) *Val {
	if a.mtype == pair {
		return True
	}

	return False
}

func bisPair(a []*Val) *Val {
	return isPair(a[0])
}

func isNil(a *Val) *Val {
	if a == Nil {
		return True
	}

	return False
}

func bisNil(a []*Val) *Val {
	return isNil(a[0])
}

func list(a ...*Val) *Val {
	l := Nil
	for i := len(a) - 1; i >= 0; i-- {
		l = cons(a[i], l)
	}

	return l
}

func blist(a []*Val) *Val {
	return list(a...)
}

func reverse(l *Val) *Val {
	checkType(l, pair, mnil)

	r := Nil
	for {
		if l == Nil {
			return r
		}

		r = cons(car(l), r)
		l = cdr(l)
	}
}

func reverseIrregular(l *Val) *Val {
	checkType(l, pair, mnil)

	r := cons(car(cdr(l)), car(l))
	l = cdr(cdr(l))
	for {
		if l == Nil {
			return r
		}

		r = cons(car(l), r)
		l = cdr(l)
	}
}

func mappend(left, right *Val) *Val {
	checkType(left, pair, mnil)
	checkType(right, pair, mnil)

	if isNil(left) != False {
		return right
	}

	return cons(car(left), mappend(cdr(left), right))
}

func List(a ...*Val) *Val {
	av := make([]*Val, len(a))
	for i, ai := range a {
		av[i] = ai
	}

	return list(av...)
}

func Cons(car, cdr *Val) *Val {
	return cons(car, cdr)
}
