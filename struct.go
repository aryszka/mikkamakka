package main

type tstruct struct {
	sys map[string]*val
}

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
