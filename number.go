package main

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
