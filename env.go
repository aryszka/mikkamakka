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

	ns := RawSymbolString(n)
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
		SymbolFromRawString(
			RawString(
				AppendString(
					SymbolToString(n),
					StringFromRaw(":"),
					SymbolToString(Car(names)),
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

	ns := RawSymbolString(n)
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
	return StringFromRaw("<environment>")
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

	define(env, SymbolFromRawString("nil"), Nil)
	define(env, SymbolFromRawString("nil?"), newBuiltin1(IsNil))
	define(env, SymbolFromRawString("pair?"), newBuiltin1(IsPair))
	define(env, SymbolFromRawString("cons"), newBuiltin2(Cons))
	define(env, SymbolFromRawString("car"), newBuiltin1(Car))
	define(env, SymbolFromRawString("cdr"), newBuiltin1(Cdr))
	define(env, SymbolFromRawString("list"), newBuiltin0V(List))
	define(env, SymbolFromRawString("apply"), newBuiltin2(Apply))
	define(env, SymbolFromRawString("error?"), newBuiltin1(isError))
	define(env, SymbolFromRawString("string->error"), newBuiltin1(stringToError))
	define(env, SymbolFromRawString("fatal"), newBuiltin1(fatal))
	define(env, SymbolFromRawString("not"), newBuiltin1(Not))
	define(env, SymbolFromRawString("="), newBuiltin0V(Eq))
	define(env, SymbolFromRawString(">"), newBuiltin0V(Greater))
	define(env, SymbolFromRawString("+"), newBuiltin0V(Add))
	define(env, SymbolFromRawString("string->number"), newBuiltin1(NumberFromString))
	define(env, SymbolFromRawString("string->bool"), newBuiltin1(BoolFromString))
	define(env, SymbolFromRawString("symbol?"), newBuiltin1(IsSymbol))
	define(env, SymbolFromRawString("symbol->string"), newBuiltin1(SymbolToString))
	define(env, SymbolFromRawString("string->symbol"), newBuiltin1(SymbolFromString))
	define(env, SymbolFromRawString("number?"), newBuiltin1(IsNumber))
	define(env, SymbolFromRawString("number->string"), newBuiltin1(NumberToString))
	define(env, SymbolFromRawString("bool?"), newBuiltin1(IsBool))
	define(env, SymbolFromRawString("bool->string"), newBuiltin1(BoolToString))
	define(env, SymbolFromRawString("string?"), newBuiltin1(IsString))
	define(env, SymbolFromRawString("assign"), newBuiltin0V(Assign))
	define(env, SymbolFromRawString("fopen"), newBuiltin1(fopen))
	define(env, SymbolFromRawString("fclose"), newBuiltin1(fclose))
	define(env, SymbolFromRawString("fread"), newBuiltin2(fread))
	define(env, SymbolFromRawString("fwrite"), newBuiltin2(fwrite))
	define(env, SymbolFromRawString("fstate"), newBuiltin1(fstate))
	define(env, SymbolFromRawString("derived-object?"), newBuiltin2(derivedObject))
	define(env, SymbolFromRawString("failing-reader"), newBuiltin0(failingReader))
	define(env, SymbolFromRawString("eof"), Eof)
	define(env, SymbolFromRawString("stdin"), newBuiltin0(stdin))
	define(env, SymbolFromRawString("stderr"), newBuiltin0(stderr))
	define(env, SymbolFromRawString("stdout"), newBuiltin0(stdout))
	define(env, SymbolFromRawString("buffer"), newBuiltin0(buffer))
	define(env, SymbolFromRawString("argv"), newBuiltin0(argv))
	define(env, SymbolFromRawString("invalid-token"), invalidToken)
	define(env, SymbolFromRawString("string-append"), newBuiltin0V(AppendString))
	define(env, SymbolFromRawString("escape-compiled-string"), newBuiltin1(EscapeCompiled))
	define(env, SymbolFromRawString("printer"), newBuiltin1(printer))
	define(env, SymbolFromRawString("print"), newBuiltin2(mprint))

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
