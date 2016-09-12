package mikkamakka

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

func InitialEnv() *Val {
	env := newEnv(nil)

	define(env, sfromString("nil"), vnil)
	define(env, sfromString("nil?"), newBuiltin(bisNil, 1, false))
	define(env, sfromString("pair?"), newBuiltin(bisPair, 1, false))
	define(env, sfromString("cons"), newBuiltin(bcons, 2, false))
	define(env, sfromString("car"), newBuiltin(bcar, 1, false))
	define(env, sfromString("cdr"), newBuiltin(bcdr, 1, false))
	define(env, sfromString("list"), newBuiltin(blist, 0, true))
	define(env, sfromString("apply"), newBuiltin(bapply, 2, false))
	define(env, sfromString("error?"), newBuiltin(bisError, 1, false))
	define(env, sfromString("string->error"), newBuiltin(stringToError, 1, false))
	define(env, sfromString("fatal"), newBuiltin(bfatal, 1, false))
	define(env, sfromString("not"), newBuiltin(not, 1, false))
	define(env, sfromString("="), newBuiltin(beq, 0, true))
	define(env, sfromString(">"), newBuiltin(bgreater, 2, false))
	define(env, sfromString("+"), newBuiltin(badd, 0, true))
	define(env, sfromString("try-string->number"), newBuiltin(btryNumberFromString, 1, false))
	define(env, sfromString("try-string->bool"), newBuiltin(btryBoolFromString, 1, false))
	define(env, sfromString("symbol?"), newBuiltin(bisSymbol, 1, false))
	define(env, sfromString("symbol->string"), newBuiltin(bsymbolToString, 1, false))
	define(env, sfromString("string->symbol"), newBuiltin(stringToSymbol, 1, false))
	define(env, sfromString("number?"), newBuiltin(bisNumber, 1, false))
	define(env, sfromString("number->string"), newBuiltin(bnumberToString, 1, false))
	define(env, sfromString("bool?"), newBuiltin(bisBool, 1, false))
	define(env, sfromString("bool->string"), newBuiltin(bboolToString, 1, false))
	define(env, sfromString("string?"), newBuiltin(bisString, 1, false))
	define(env, sfromString("assign"), newBuiltin(bassign, 1, true))
	define(env, sfromString("fopen"), newBuiltin(bfopen, 1, false))
	define(env, sfromString("fclose"), newBuiltin(bfclose, 1, false))
	define(env, sfromString("fread"), newBuiltin(bfread, 2, false))
	define(env, sfromString("fwrite"), newBuiltin(bfwrite, 2, false))
	define(env, sfromString("fstate"), newBuiltin(bfstate, 1, false))
	define(env, sfromString("derived-object?"), newBuiltin(derivedObject, 2, false))
	define(env, sfromString("failing-reader"), newBuiltin(failingReader, 0, false))
	define(env, sfromString("eof"), eof)
	define(env, sfromString("stdin"), newBuiltin(bstdin, 0, false))
	define(env, sfromString("stderr"), newBuiltin(bstderr, 0, false))
	define(env, sfromString("stdout"), newBuiltin(bstdout, 0, false))
	define(env, sfromString("buffer"), newBuiltin(bbuffer, 0, false))
	define(env, sfromString("argv"), newBuiltin(argv, 0, false))
	define(env, sfromString("invalid-token"), invalidToken)
	define(env, sfromString("string-append"), newBuiltin(bappendString, 0, true))
	define(env, sfromString("escape-compiled-string"), newBuiltin(escapeCompiled, 1, false))
	define(env, sfromString("printer"), newBuiltin(bprinter, 1, false))
	define(env, sfromString("print"), newBuiltin(bprint, 2, false))

	return (*Val)(env)
}

func Define(e, n, v *Val) *Val {
	return (*Val)(define((*val)(e), (*val)(n), (*val)(v)))
}

func ExtendEnv(e, n, v *Val) *Val {
	return (*Val)(extendEnv((*val)(e), (*val)(n), (*val)(v)))
}

func LookupDef(e, n *Val) *Val {
	return (*Val)(lookupDef((*val)(e), (*val)(n)))
}
