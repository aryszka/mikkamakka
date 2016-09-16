package mikkamakka

type Struct map[string]*Val

var InvalidStruct = &Val{merror, "invalid struct"}

func FromMap(m Struct) *Val {
	return &Val{mstruct, m}
}

func Assign(a ...*Val) *Val {
	m := make(Struct)
	for _, ai := range a {
		checkType(ai, mstruct)
		for k, v := range ai.value.(Struct) {
			m[k] = v
		}
	}

	return FromMap(m)
}

func StructFromList(l *Val) *Val {
	m := make(Struct)
	for {
		if l == Nil {
			break
		}

		if IsPair(l) == False || IsPair(Cdr(l)) == False || IsSymbol(Car(l)) == False {
			return fatal(InvalidStruct)
		}

		m[RawSymbolString(Car(l))], l = Car(Cdr(l)), Cdr(Cdr(l))
	}

	return FromMap(m)
}

func IsStruct(a *Val) *Val {
	return is(a, mstruct)
}

func Field(s, f *Val) *Val {
	checkType(s, mstruct)

	name := RawSymbolString(f)
	v, ok := s.value.(Struct)[name]
	if !ok {
		return fatal(&Val{merror, "undefined field name: " + name})
	}

	return v
}

func StructNames(s *Val) *Val {
	checkType(s, mstruct)

	n := Nil
	for k, _ := range s.value.(Struct) {
		n = Cons(SymbolFromRawString(k), n)
	}

	return n
}
