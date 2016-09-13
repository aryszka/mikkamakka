package mikkamakka

import "strconv"

func intVal(n *Val) int {
	return n.value.(int)
}

func fromInt(i int) *Val {
	return &Val{number, i}
}

func greater(a ...*Val) *Val {
	return bgreater(a)
}

func bgreater(a []*Val) *Val {
	for {
		if len(a) == 0 {
			return vfalse
		}

		checkType(a[0], number)

		if len(a) == 1 {
			return vtrue
		}

		checkType(a[1], number)

		if a[0].value.(int) <= a[1].value.(int) {
			return vfalse
		}

		a = a[1:]
	}
}

func nfromString(s string) *Val {
	n, err := strconv.Atoi(s)
	if err != nil {
		return invalidToken
	}

	return fromInt(n)
}

func tryNumberFromString(s *Val) *Val {
	checkType(s, mstring)
	return nfromString(stringVal(s))
}

func btryNumberFromString(a []*Val) *Val {
	return tryNumberFromString(a[0])
}

func numberToString(n *Val) *Val {
	checkType(n, number)
	return fromString(strconv.Itoa(n.value.(int)))
}

func bnumberToString(a []*Val) *Val {
	return numberToString(a[0])
}

func isNumber(a *Val) *Val {
	return is(a, number)
}

func bisNumber(a []*Val) *Val {
	return isNumber(a[0])
}

func sub(left, right *Val) *Val {
	checkType(left, number)
	checkType(right, number)
	return fromInt(left.value.(int) - right.value.(int))
}

func add(left, right *Val) *Val {
	checkType(left, number)
	checkType(right, number)
	return fromInt(left.value.(int) + right.value.(int))
}

func badd(a []*Val) *Val {
	s := 0
	for {
		if len(a) == 0 {
			return fromInt(s)
		}

		checkType(a[0], number)
		s += intVal(a[0])
		a = a[1:]
	}
}

func neq(left, right *Val) *Val {
	if left.value.(int) == right.value.(int) {
		return vtrue
	}

	return vfalse
}

func FromInt(i int) *Val {
	return (*Val)(fromInt(i))
}
