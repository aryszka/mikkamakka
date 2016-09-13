package mikkamakka

type sym struct {
	val string
}

var emptySymbol = &Val{merror, "empty symbol not allowed"}

func sfromString(s string) *Val {
	if s == "" {
		return emptySymbol
	}

	return &Val{symbol, &sym{s}}
}

func sstringVal(s *Val) string {
	return s.value.(*sym).val
}

func symbolToString(s *Val) *Val {
	return fromString(sstringVal(s))
}

func bsymbolToString(a []*Val) *Val {
	return symbolToString(a[0])
}

func stringToSymbol(a []*Val) *Val {
	checkType(a[0], mstring)
	return sfromString(stringVal(a[0]))
}

func isSymbol(a *Val) *Val {
	if a.mtype == symbol {
		return True
	}

	return False
}

func bisSymbol(a []*Val) *Val {
	return isSymbol(a[0])
}

func smeq(left, right *Val) *Val {
	checkType(left, symbol)
	checkType(right, symbol)
	if sstringVal(left) == sstringVal(right) {
		return True
	}

	return False
}

func SfromString(s string) *Val {
	return sfromString(s)
}
