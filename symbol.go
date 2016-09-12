package mikkamakka

type sym struct {
	val string
}

var emptySymbol = &val{merror, "empty symbol not allowed"}

func sfromString(s string) *val {
	if s == "" {
		return emptySymbol
	}

	return &val{symbol, &sym{s}}
}

func sstringVal(s *val) string {
	return s.value.(*sym).val
}

func symbolToString(s *val) *val {
	return fromString(sstringVal(s))
}

func bsymbolToString(a []*val) *val {
	return symbolToString(a[0])
}

func stringToSymbol(a []*val) *val {
	checkType(a[0], mstring)
	return sfromString(stringVal(a[0]))
}

func isSymbol(a *val) *val {
	if a.mtype == symbol {
		return vtrue
	}

	return vfalse
}

func bisSymbol(a []*val) *val {
	return isSymbol(a[0])
}

func smeq(left, right *val) *val {
	checkType(left, symbol)
	checkType(right, symbol)
	if sstringVal(left) == sstringVal(right) {
		return vtrue
	}

	return vfalse
}

func SfromString(s string) *Val {
	return (*Val)(sfromString(s))
}
