package mikkamakka

var InvalidStruct = SysStringToError("invalid struct")

func SysMapToStruct(m map[string]*Val) *Val {
	return newVal(Struct, m)
}

func Assign(a ...*Val) *Val {
	m := make(map[string]*Val)
	for _, ai := range a {
		checkType(ai, Struct)
		for k, v := range ai.value.(map[string]*Val) {
			m[k] = v
		}
	}

	return SysMapToStruct(m)
}

func ListToStruct(l *Val) *Val {
	m := make(map[string]*Val)
	for {
		if l == NilVal {
			break
		}

		if IsPair(l) == False || IsPair(Cdr(l)) == False || IsSymbol(Car(l)) == False {
			return Fatal(InvalidStruct)
		}

		m[SymbolToSysString(Car(l))], l = Car(Cdr(l)), Cdr(Cdr(l))
	}

	return SysMapToStruct(m)
}

func IsStruct(a *Val) *Val {
	return is(a, Struct)
}

func Field(s, f *Val) *Val {
	checkType(s, Struct)

	name := SymbolToSysString(f)
	v, ok := s.value.(map[string]*Val)[name]
	if !ok {
		return Fatal(SysStringToError("undefined field name: " + name))
	}

	return v
}

func StructNames(s *Val) *Val {
	checkType(s, Struct)

	n := NilVal
	for k, _ := range s.value.(map[string]*Val) {
		n = Cons(SymbolFromRawString(k), n)
	}

	return n
}
