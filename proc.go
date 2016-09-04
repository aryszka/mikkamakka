package main

type proc struct {
	params *val
	body   *val
	env    *val
}

func newProc(e, p, b *val) *val {
	return &val{procedure, &proc{p, b, e}}
}

func procString(e *val) *val {
	checkType(e, procedure)
	return fromString("<procedure>")
}
