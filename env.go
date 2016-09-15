package mikkamakka

type env struct {
	current map[string]*Val
	parent  *Val
}

var (
	undefined        = &Val{merror, "undefined"}
	notAnEnvironment = &Val{merror, "not an environment"}
	definitionExists = &Val{merror, "definition exists"}
)

func newEnv(p *Val) *Val {
	return &Val{
		environment,
		&env{
			current: make(map[string]*Val),
			parent:  p,
		},
	}
}

func lookupDef(e, n *Val) *Val {
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

func defineStruct(e, n, s, names *Val) *Val {
	checkType(e, environment)
	checkType(n, symbol)
	checkType(s, mstruct)
	checkType(names, pair, mnil)

	if IsNil(names) != False {
		return s
	}

	define(
		e,
		sfromString(
			stringVal(
				appendString(
					symbolToString(n),
					fromString(":"),
					symbolToString(Car(names)),
				),
			),
		),
		structVal(s, Car(names)),
	)

	return defineStruct(e, n, s, Cdr(names))
}

func define(e, n, v *Val) *Val {
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
	if isStruct(v) != False {
		return defineStruct(e, n, v, structNames(v))
	}

	return v
}

func defineAll(e, n, a *Val) *Val {
	for {
		if IsNil(n) != False && IsNil(a) != False {
			break
		}

		if IsPair(a) == False && IsNil(a) == False {
			return fatal(InvalidArgs)
		}

		if IsPair(n) == False {
			define(e, n, a)
			return e
		}

		if IsNil(a) != False {
			return fatal(InvalidArgs)
		}

		ni := Car(n)
		ai := Car(a)
		define(e, ni, ai)
		n, a = Cdr(n), Cdr(a)
	}

	return e
}

func extendEnv(e, n, a *Val) *Val {
	e = newEnv(e)
	return defineAll(e, n, a)
}

func envString(e *Val) *Val {
	checkType(e, environment)
	return fromString("<environment>")
}

func isEnv(e *Val) *Val {
	if e.mtype == environment {
		return True
	}

	return False
}

func newBuiltin0(f func() *Val) *Val {
	return NewCompiled(0, false, func(a []*Val) *Val {
		return f()
	})
}

func newBuiltin0V(f func(...*Val) *Val) *Val {
	return NewCompiled(0, true, func(a []*Val) *Val {
		return f(a...)
	})
}

func newBuiltin1(f func(*Val) *Val) *Val {
	return NewCompiled(1, false, func(a []*Val) *Val {
		return f(a[0])
	})
}

func newBuiltin2(f func(*Val, *Val) *Val) *Val {
	return NewCompiled(2, false, func(a []*Val) *Val {
		return f(a[0], a[1])
	})
}

func InitialEnv() *Val {
	env := newEnv(nil)

	define(env, sfromString("nil"), Nil)
	define(env, sfromString("nil?"), newBuiltin1(IsNil))
	define(env, sfromString("pair?"), newBuiltin1(IsPair))
	define(env, sfromString("cons"), newBuiltin2(Cons))
	define(env, sfromString("car"), newBuiltin1(Car))
	define(env, sfromString("cdr"), newBuiltin1(Cdr))
	define(env, sfromString("list"), newBuiltin0V(List))
	define(env, sfromString("apply"), newBuiltin2(Apply))
	define(env, sfromString("error?"), newBuiltin1(isError))
	define(env, sfromString("string->error"), newBuiltin1(stringToError))
	define(env, sfromString("fatal"), newBuiltin1(fatal))
	define(env, sfromString("not"), newBuiltin1(not))
	define(env, sfromString("="), newBuiltin0V(Eq))
	define(env, sfromString(">"), newBuiltin0V(greater))
	define(env, sfromString("+"), newBuiltin0V(add))
	define(env, sfromString("try-string->number"), newBuiltin1(tryNumberFromString))
	define(env, sfromString("try-string->bool"), newBuiltin1(tryBoolFromString))
	define(env, sfromString("symbol?"), newBuiltin1(isSymbol))
	define(env, sfromString("symbol->string"), newBuiltin1(symbolToString))
	define(env, sfromString("string->symbol"), newBuiltin1(stringToSymbol))
	define(env, sfromString("number?"), newBuiltin1(isNumber))
	define(env, sfromString("number->string"), newBuiltin1(numberToString))
	define(env, sfromString("bool?"), newBuiltin1(isBool))
	define(env, sfromString("bool->string"), newBuiltin1(boolToString))
	define(env, sfromString("string?"), newBuiltin1(isString))
	define(env, sfromString("assign"), newBuiltin0V(Assign))
	define(env, sfromString("fopen"), newBuiltin1(fopen))
	define(env, sfromString("fclose"), newBuiltin1(fclose))
	define(env, sfromString("fread"), newBuiltin2(fread))
	define(env, sfromString("fwrite"), newBuiltin2(fwrite))
	define(env, sfromString("fstate"), newBuiltin1(fstate))
	define(env, sfromString("derived-object?"), newBuiltin2(derivedObject))
	define(env, sfromString("failing-reader"), newBuiltin0(failingReader))
	define(env, sfromString("eof"), Eof)
	define(env, sfromString("stdin"), newBuiltin0(stdin))
	define(env, sfromString("stderr"), newBuiltin0(stderr))
	define(env, sfromString("stdout"), newBuiltin0(stdout))
	define(env, sfromString("buffer"), newBuiltin0(buffer))
	define(env, sfromString("argv"), newBuiltin0(argv))
	define(env, sfromString("invalid-token"), invalidToken)
	define(env, sfromString("string-append"), newBuiltin0V(appendString))
	define(env, sfromString("escape-compiled-string"), newBuiltin1(escapeCompiled))
	define(env, sfromString("printer"), newBuiltin1(printer))
	define(env, sfromString("print"), newBuiltin2(mprint))

	return (*Val)(env)
}

func Define(e, n, v *Val) *Val {
	return define(e, n, v)
}

func ExtendEnv(e, n, v *Val) *Val {
	return extendEnv(e, n, v)
}

func LookupDef(e, n *Val) *Val {
	return lookupDef(e, n)
}
