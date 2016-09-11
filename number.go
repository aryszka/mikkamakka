package mikkamakka

import "strconv"

func intVal(n *val) int {
	return n.value.(int)
}

func fromInt(i int) *val {
	return &val{number, i}
}

func greater(a ...*val) *val {
	return bgreater(a)
}

func bgreater(a []*val) *val {
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

func nfromString(s string) *val {
	n, err := strconv.Atoi(s)
	if err != nil {
		return invalidToken
	}

	return fromInt(n)
}

func tryNumberFromString(s *val) *val {
	checkType(s, mstring)
	return nfromString(stringVal(s))
}

func btryNumberFromString(a []*val) *val {
	return tryNumberFromString(a[0])
}

func numberToString(n *val) *val {
	checkType(n, number)
	return fromString(strconv.Itoa(n.value.(int)))
}

func bnumberToString(a []*val) *val {
	return numberToString(a[0])
}

func isNumber(a *val) *val {
	return is(a, number)
}

func bisNumber(a []*val) *val {
	return isNumber(a[0])
}

func sub(left, right *val) *val {
	checkType(left, number)
	checkType(right, number)
	return fromInt(left.value.(int) - right.value.(int))
}

func add(left, right *val) *val {
	checkType(left, number)
	checkType(right, number)
	return fromInt(left.value.(int) + right.value.(int))
}

func badd(a []*val) *val {
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

func neq(left, right *val) *val {
	if left.value.(int) == right.value.(int) {
		return vtrue
	}

	return vfalse
}

func FromInt(i int) *Val {
	return (*Val)(fromInt(i))
}
