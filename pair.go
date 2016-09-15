package mikkamakka

type tpair struct {
	car, cdr *Val
}

var Nil = &Val{mtype: mnil}

func Cons(car, cdr *Val) *Val {
	return &Val{pair, &tpair{car, cdr}}
}

func Car(p *Val) *Val {
	checkType(p, pair)
	return p.value.(*tpair).car
}

func Cdr(p *Val) *Val {
	checkType(p, pair)
	return p.value.(*tpair).cdr
}

func IsPair(a *Val) *Val {
	if a.mtype == pair {
		return True
	}

	return False
}

func IsNil(a *Val) *Val {
	if a == Nil {
		return True
	}

	return False
}
