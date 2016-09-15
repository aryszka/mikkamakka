package mikkamakka

import "strconv"

func intVal(n *Val) int {
	return n.value.(int)
}

func fromInt(i int) *Val {
	return &Val{number, i}
}

func greater(a ...*Val) *Val {
	for {
		if len(a) == 0 {
			return False
		}

		checkType(a[0], number)

		if len(a) == 1 {
			return True
		}

		checkType(a[1], number)

		if a[0].value.(int) <= a[1].value.(int) {
			return False
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

func numberToString(n *Val) *Val {
	checkType(n, number)
	return fromString(strconv.Itoa(n.value.(int)))
}

func isNumber(a *Val) *Val {
	return is(a, number)
}

func sub(left, right *Val) *Val {
	checkType(left, number)
	checkType(right, number)
	return fromInt(left.value.(int) - right.value.(int))
}

func add(a ...*Val) *Val {
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
		return True
	}

	return False
}

func FromInt(i int) *Val {
	return (*Val)(fromInt(i))
}
