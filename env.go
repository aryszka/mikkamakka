package main

type env struct {
	current map[string]*val
	parent  *val
}

var (
	undefined        = &val{merror, "undefined"}
	notAnEnvironment = &val{merror, "not an environment"}
	definitionExists = &val{merror, "definition exists"}
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

func lookupDef(e, n *val) *val {
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
		return fatal(undefined)
	}

	return lookupDef(et.parent, n)
}

func define(e, n, v *val) *val {
	checkType(e, environment)
	checkType(n, symbol)
	et, ok := e.value.(*env)
	if !ok {
		return fatal(notAnEnvironment)
	}

	ns := sstringVal(n)
	if _, has := et.current[ns]; has {
		return fatal(definitionExists)
	}

	et.current[ns] = v
	return v
}

func extendEnv(e, n, a *val) *val {
	m := make(map[string]*val)
	for {
		if isNil(n) != vfalse && isNil(a) != vfalse {
			break
		}

		if isPair(n) == vfalse || isPair(a) == vfalse {
			return fatal(invalidArguments)
		}

		m[sstringVal(car(n))] = car(a)
		n, a = cdr(n), cdr(a)
	}

	return &val{
		environment,
		&env{
			current: m,
			parent:  e,
		},
	}
}

func envString(e *val) *val {
	checkType(e, environment)
	return fromString("<environment>")
}

func isEnv(e *val) *val {
	if e.mtype == environment {
		return vtrue
	}

	return vfalse
}
