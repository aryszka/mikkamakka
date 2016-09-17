package mikkamakka

type pair struct {
	car, cdr *Val
}

var NilVal = newVal(Nil, nil)

func Cons(car, cdr *Val) *Val {
	return newVal(Pair, &pair{car, cdr})
}

func Car(p *Val) *Val {
	checkType(p, Pair)
	return p.value.(*pair).car
}

func Cdr(p *Val) *Val {
	checkType(p, Pair)
	return p.value.(*pair).cdr
}

func IsPair(a *Val) *Val {
	return is(a, Pair)
}

func IsNil(a *Val) *Val {
	return is(a, Nil)
}
