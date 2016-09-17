package mikkamakka

var EmptySymbol = SysStringToError("empty symbol not allowed")

func SymbolFromRawString(s string) *Val {
	if s == "" {
		return EmptySymbol
	}

	return newVal(Symbol, s)
}

func SymbolFromString(a *Val) *Val {
	return SymbolFromRawString(StringToSysString(a))
}

func SymbolToSysString(s *Val) string {
	checkType(s, Symbol)
	return s.value.(string)
}

func SymbolToString(s *Val) *Val {
	return SysStringToString(SymbolToSysString(s))
}

func IsSymbol(a *Val) *Val {
	return is(a, Symbol)
}

func symbolEq(left, right *Val) *Val {
	if SymbolToSysString(left) == SymbolToSysString(right) {
		return True
	}

	return False
}
