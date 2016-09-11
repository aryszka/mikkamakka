package mikkamakka

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

func bcons(a []*val) *val {
	return cons(a[0], a[1])
}

func car(p *val) *val {
	checkType(p, pair)
	return p.value.(*tpair).car
}

func bcar(a []*val) *val {
	return car(a[0])
}

func cdr(p *val) *val {
	checkType(p, pair)
	return p.value.(*tpair).cdr
}

func bcdr(a []*val) *val {
	return cdr(a[0])
}

func isPair(a *val) *val {
	if a.mtype == pair {
		return vtrue
	}

	return vfalse
}

func bisPair(a []*val) *val {
	return isPair(a[0])
}

func isNil(a *val) *val {
	if a == vnil {
		return vtrue
	}

	return vfalse
}

func bisNil(a []*val) *val {
	return isNil(a[0])
}

func list(a ...*val) *val {
	l := vnil
	for i := len(a) - 1; i >= 0; i-- {
		l = cons(a[i], l)
	}

	return l
}

func blist(a []*val) *val {
	return list(a...)
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

func mappend(left, right *val) *val {
	checkType(left, pair, mnil)
	checkType(right, pair, mnil)

	if isNil(left) != vfalse {
		return right
	}

	return cons(car(left), mappend(cdr(left), right))
}
