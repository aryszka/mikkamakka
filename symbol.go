package mikkamakka

var EmptySymbol = &Val{merror, "empty symbol not allowed"}

func SymbolFromRawString(s string) *Val {
	if s == "" {
		return EmptySymbol
	}

	return &Val{symbol, s}
}

func SymbolFromString(a *Val) *Val {
	return SymbolFromRawString(RawString(a))
}

func RawSymbolString(s *Val) string {
	checkType(s, symbol)
	return s.value.(string)
}

func SymbolToString(s *Val) *Val {
	return StringFromRaw(RawSymbolString(s))
}

func IsSymbol(a *Val) *Val {
	if a.mtype == symbol {
		return True
	}

	return False
}

func symbolEq(left, right *Val) *Val {
	if RawSymbolString(left) == RawSymbolString(right) {
		return True
	}

	return False
}
