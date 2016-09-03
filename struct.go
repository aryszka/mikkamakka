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
	return s.value.(*tstruct).sys[sstringVal(f)]
}
