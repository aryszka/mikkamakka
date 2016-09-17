package mikkamakka

import "strconv"

var InvalidNumberString = SysStringToError("invalid number string")

func SysIntToNumber(i int) *Val {
	return newVal(Number, i)
}

func SysStringToNumber(s string) *Val {
	n, err := strconv.Atoi(s)
	if err != nil {
		return InvalidNumberString
	}

	return SysIntToNumber(n)
}

func StringToNumber(s *Val) *Val {
	return SysStringToNumber(StringToSysString(s))
}

func NumberToSysInt(n *Val) int {
	checkType(n, Number)
	return n.value.(int)
}

func NumberToString(n *Val) *Val {
	checkType(n, Number)
	return SysStringToString(strconv.Itoa(n.value.(int)))
}

func IsNumber(a *Val) *Val {
	return is(a, Number)
}

func numberEq(left, right *Val) *Val {
	if NumberToSysInt(left) == NumberToSysInt(right) {
		return True
	}

	return False
}

func op(n int, a []*Val, f func(int, int) int) *Val {
	for {
		if len(a) == 0 {
			return SysIntToNumber(n)
		}

		n, a = f(n, NumberToSysInt(a[0])), a[1:]
	}
}

func Sub(a0 *Val, a ...*Val) *Val {
	if len(a) == 0 {
		return SysIntToNumber(0 - NumberToSysInt(a0))
	}

	return op(NumberToSysInt(a0), a, func(prev, next int) int {
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

		checkType(a[0], Number)
		if len(a) == 1 {
			return True
		}

		if NumberToSysInt(a[0]) <= NumberToSysInt(a[1]) {
			return False
		}

		a = a[1:]
	}
}
