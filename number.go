package mikkamakka

import "strconv"

var InvalidNumberString = ErrorFromRawString("invalid number string")

func NumberFromRawInt(i int) *Val {
	return &Val{number, i}
}

func NumberFromRawString(s string) *Val {
	n, err := strconv.Atoi(s)
	if err != nil {
		return InvalidNumberString
	}

	return NumberFromRawInt(n)
}

func NumberFromString(s *Val) *Val {
	return NumberFromRawString(RawString(s))
}

func RawInt(n *Val) int {
	checkType(n, number)
	return n.value.(int)
}

func NumberToString(n *Val) *Val {
	checkType(n, number)
	return StringFromRaw(strconv.Itoa(n.value.(int)))
}

func IsNumber(a *Val) *Val {
	return is(a, number)
}

func numberEq(left, right *Val) *Val {
	if RawInt(left) == RawInt(right) {
		return True
	}

	return False
}

func op(n int, a []*Val, f func(int, int) int) *Val {
	for {
		if len(a) == 0 {
			return NumberFromRawInt(n)
		}

		n, a = f(n, RawInt(a[0])), a[1:]
	}
}

func Sub(a0 *Val, a ...*Val) *Val {
	if len(a) == 0 {
		return NumberFromRawInt(0 - RawInt(a0))
	}

	return op(RawInt(a0), a, func(prev, next int) int {
		return prev - next
	})
}

func Add(a ...*Val) *Val {
	return op(0, a, func(prev, next int) int {
		return prev + next
	})
}

func Greater(a ...*Val) *Val {
	for {
		if len(a) == 0 {
			return False
		}

		checkType(a[0], number)
		if len(a) == 1 {
			return True
		}

		if RawInt(a[0]) <= RawInt(a[1]) {
			return False
		}

		a = a[1:]
	}
}
