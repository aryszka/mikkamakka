package main

type env struct {
	current map[string]*val
	parent  *val
}

var (
	undefinedVariable = &val{merror, "undefined variable"}
	notAnEnvironment  = &val{merror, "not an environment"}
	variableExists    = &val{merror, "variable exists"}
)

func newEnv() *val {
	return &val{
		environment,
		&env{
			current: make(map[string]*val),
			parent:  nil,
		},
	}
}

func lookupVar(e, n *val) *val {
	checkType(e, environment)
	checkType(n, symbol)
	et, ok := e.value.(*env)
	if !ok {
		return fatal(notAnEnvironment)
	}

	ns := sstringVal(n)
	if v, ok := et.current[ns]; ok {
		return v
	}

	if et.parent == nil {
		return fatal(undefinedVariable)
	}

	return lookupVar(et.parent, n)
}

func defVar(e, n, v *val) *val {
	checkType(e, environment)
	checkType(n, symbol)
	et, ok := e.value.(*env)
	if !ok {
		return fatal(notAnEnvironment)
	}

	ns := sstringVal(n)
	if _, has := et.current[ns]; has {
		return fatal(variableExists)
	}

	et.current[ns] = v
	return v
}

func envString(e *val) *val {
	checkType(e, environment)
	return fromString("<environment>")
}
