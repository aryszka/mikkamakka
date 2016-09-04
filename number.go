package main

import "strconv"

func intVal(n *val) int {
	return n.value.(int)
}

func fromInt(i int) *val {
	return &val{number, i}
}

func greater(left, right *val) *val {
	checkType(left, number)
	checkType(right, number)
	if left.value.(int) > right.value.(int) {
		return vtrue
	}

	return vfalse
}

func nfromString(s string) *val {
	n, err := strconv.Atoi(s)
	if err != nil {
		return invalidToken
	}

	return fromInt(n)
}

func numberToString(n *val) *val {
	checkType(n, number)
	return fromString(strconv.Itoa(n.value.(int)))
}

func isNumber(a *val) *val {
	return is(a, number)
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

func neq(left, right *val) bool {
	return left.value.(int) == right.value.(int)
}
