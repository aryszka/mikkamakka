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

func newEnv(p *val) *val {
	return &val{
		environment,
		&env{
			current: make(map[string]*val),
			parent:  p,
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
		println("undefined reference", ns)
		return fatal(undefined)
	}

	return lookupDef(et.parent, n)
}

func defineStruct(e, n, s, names *val) *val {
	checkType(e, environment)
	checkType(n, symbol)
	checkType(s, mstruct)
	checkType(names, pair, mnil)

	if isNil(names) != vfalse {
		return s
	}

	define(
		e,
		sfromString(
			stringVal(
				appendString(
					symbolToString(n),
					fromString(":"),
					symbolToString(car(names)),
				),
			),
		),
		structVal(s, car(names)),
	)

	return defineStruct(e, n, s, cdr(names))
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
	if isStruct(v) != vfalse {
		return defineStruct(e, n, v, structNames(v))
	}

	return v
}

func defineAll(e, n, a *val) *val {
	for {
		if isNil(n) != vfalse && isNil(a) != vfalse {
			break
		}

		if isPair(a) == vfalse && isNil(a) == vfalse {
			return fatal(invalidArguments)
		}

		if isPair(n) == vfalse {
			define(e, n, a)
			return e
		}

		if isNil(a) != vfalse {
			return fatal(invalidArguments)
		}

		ni := car(n)
		ai := car(a)
		define(e, ni, ai)
		n, a = cdr(n), cdr(a)
	}

	return e
}

func extendEnv(e, n, a *val) *val {
	e = newEnv(e)
	return defineAll(e, n, a)
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
