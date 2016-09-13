package mikkamakka

type tstruct struct {
	sys map[string]*Val
}

var invalidStruct = &Val{merror, "invalid struct"}

func fromMap(m map[string]*Val) *Val {
	return &Val{mstruct, &tstruct{m}}
}

func field(s, f *Val) *Val {
	checkType(s, mstruct)
	checkType(f, symbol)
	name := sstringVal(f)
	v, ok := s.value.(*tstruct).sys[name]
	if !ok {
		panic("undefined field name: " + name)
	}

	return v
}

func Assign(a *Val) *Val {
	next := make(map[string]*Val)
	s := car(a)
	a = cdr(a)
	for k, v := range s.value.(*tstruct).sys {
		next[k] = v
	}

	for {
		if isNil(a) != False {
			break
		}

		for k, v := range car(a).value.(*tstruct).sys {
			next[k] = v
		}

		a = cdr(a)
	}

	return fromMap(next)
}

func bassign(a []*Val) *Val {
	return Assign(cons(a[0], sliceToList(a[1:])))
}

func structFromList(l *Val) *Val {
	sys := make(map[string]*Val)
	for {
		if l == Nil {
			break
		}

		if isPair(l) == False || isPair(cdr(l)) == False || isSymbol(car(l)) == False {
			return fatal(invalidStruct)
		}

		sys[sstringVal(car(l))], l = car(cdr(l)), cdr(cdr(l))
	}

	return fromMap(sys)
}

func isStruct(a *Val) *Val {
	if a.mtype == mstruct {
		return True
	}

	return False
}

func structNames(s *Val) *Val {
	checkType(s, mstruct)

	n := Nil
	for k, _ := range s.value.(*tstruct).sys {
		n = cons(sfromString(k), n)
	}

	return reverse(n)
}

func structVal(s, n *Val) *Val {
	checkType(s, mstruct)
	checkType(n, symbol)
	ns := sstringVal(n)

	if v, ok := s.value.(*tstruct).sys[ns]; !ok {
		return fatal(undefined)
	} else {
		return v
	}
}

func Field(s, f *Val) *Val {
	return field(s, f)
}

func FromMap(m map[string]*Val) *Val {
	mv := make(map[string]*Val)
	for k, v := range m {
		mv[k] = v
	}

	return fromMap(mv)
}

func StructFromList(l *Val) *Val {
	return structFromList(l)
}
