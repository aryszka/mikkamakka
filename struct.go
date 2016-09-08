package main

type tstruct struct {
	sys map[string]*val
}

var invalidStruct = &val{merror, "invalid struct"}

func fromMap(m map[string]*val) *val {
	return &val{mstruct, &tstruct{m}}
}

func field(s, f *val) *val {
	checkType(s, mstruct)
	checkType(f, symbol)
	name := sstringVal(f)
	v, ok := s.value.(*tstruct).sys[name]
	if !ok {
		panic("undefined field name: " + name)
	}

	return v
}

func assign(s *val, a ...*val) *val {
	checkType(s, mstruct)

	next := make(map[string]*val)
	for k, v := range s.value.(*tstruct).sys {
		next[k] = v
	}

	for _, ai := range a {
		checkType(ai, mstruct)
		for k, v := range ai.value.(*tstruct).sys {
			next[k] = v
		}
	}

	return fromMap(next)
}

func structFromList(l *val) *val {
	sys := make(map[string]*val)
	for {
		if l == vnil {
			break
		}

		if isPair(l) == vfalse || isPair(cdr(l)) == vfalse || isSymbol(car(l)) == vfalse {
			return fatal(invalidStruct)
		}

		sys[sstringVal(car(l))], l = car(cdr(l)), cdr(cdr(l))
	}

	return fromMap(sys)
}

func isStruct(a *val) *val {
	if a.mtype == mstruct {
		return vtrue
	}

	return vfalse
}

func structNames(s *val) *val {
	checkType(s, mstruct)

	n := vnil
	for k, _ := range s.value.(*tstruct).sys {
		n = cons(sfromString(k), n)
	}

	return reverse(n)
}
