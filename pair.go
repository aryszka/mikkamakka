package main

type tpair struct {
	car, cdr *val
}

var (
	vnil       = &val{mtype: mnil}
	mixedTypes = &val{merror, "mixed types"}
)

func cons(car, cdr *val) *val {
	// if car.mtype != cdr.mtype && cdr.mtype != mnil && cdr.mtype != pair {
	// 	panic(mixedTypes)
	// }

	return &val{pair, &tpair{car, cdr}}
}

func car(p *val) *val {
	checkType(p, pair)
	return p.value.(*tpair).car
}

func cdr(p *val) *val {
	checkType(p, pair)
	return p.value.(*tpair).cdr
}

func isPair(a *val) *val {
	if a.mtype == pair {
		return vtrue
	}

	return vfalse
}

func isNil(a *val) *val {
	if a == vnil {
		return vtrue
	}

	return vfalse
}

func list(a ...*val) *val {
	l := vnil
	for i := len(a) - 1; i >= 0; i-- {
		l = cons(a[i], l)
	}

	return l
}

func reverse(l *val) *val {
	checkType(l, pair, mnil)

	r := vnil
	for {
		if l == vnil {
			return r
		}

		r = cons(car(l), r)
		l = cdr(l)
	}
}

func reverseIrregular(l *val) *val {
	checkType(l, pair, mnil)

	r := cons(car(cdr(l)), car(l))
	l = cdr(cdr(l))
	for {
		if l == vnil {
			return r
		}

		r = cons(car(l), r)
		l = cdr(l)
	}
}
