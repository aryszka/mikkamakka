package main

type sym struct {
	val string
}

func sfromString(s string) *val {
	return &val{symbol, &sym{s}}
}

func sstringVal(s *val) string {
	return s.value.(*sym).val
}

func symbolToString(s *val) *val {
	return fromString(sstringVal(s))
}

func isSymbol(a *val) *val {
	if a.mtype == symbol {
		return vtrue
	}

	return vfalse
}
