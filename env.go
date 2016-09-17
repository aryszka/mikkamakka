package mikkamakka

type env struct {
	current Struct
	parent  *Val
}

var (
	undefined        = ErrorFromRawString("undefined")
	notAnEnvironment = ErrorFromRawString("not an environment")
	definitionExists = ErrorFromRawString("definition exists")
)

func newEnv(p *Val) *Val {
	return &Val{
		environment,
		&env{
			current: make(Struct),
			parent:  p,
		},
	}
}

func lookupDef(e, n *Val) *Val {
	checkType(e, environment)
	checkType(n, symbol)
	et, ok := e.value.(*env)
	if !ok {
		return Fatal(notAnEnvironment)
	}

	ns := RawSymbolString(n)
	if v, ok := et.current[ns]; ok {
		return v
	}

	if et.parent == nil {
		println("undefined reference", ns)
		return Fatal(undefined)
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
		Field(s, Car(names)),
	)

	return defineStruct(e, n, s, Cdr(names))
}

func define(e, n, v *Val) *Val {
	checkType(e, environment)
	checkType(n, symbol)
	et, ok := e.value.(*env)
	if !ok {
		return Fatal(notAnEnvironment)
	}

	ns := RawSymbolString(n)
	if _, has := et.current[ns]; has {
		return Fatal(definitionExists)
	}

	et.current[ns] = v
	if IsStruct(v) != False {
		return defineStruct(e, n, v, StructNames(v))
	}

	return v
}

func defineAll(e, n, a *Val) *Val {
	for {
		if IsNil(n) != False && IsNil(a) != False {
			break
		}

		if IsPair(a) == False && IsNil(a) == False {
			return Fatal(InvalidArgs)
		}

		if IsPair(n) == False {
			define(e, n, a)
			return e
		}

		if IsNil(a) != False {
			return Fatal(InvalidArgs)
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
	define(env, SymbolFromRawString("error?"), newBuiltin1(IsError))
	define(env, SymbolFromRawString("string->error"), newBuiltin1(StringToError))
	define(env, SymbolFromRawString("fatal"), newBuiltin1(Fatal))
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
	define(env, SymbolFromRawString("fopen"), newBuiltin1(Fopen))
	define(env, SymbolFromRawString("fclose"), newBuiltin1(Fclose))
	define(env, SymbolFromRawString("fread"), newBuiltin2(Fread))
	define(env, SymbolFromRawString("fwrite"), newBuiltin2(Fwrite))
	define(env, SymbolFromRawString("fstate"), newBuiltin1(Fstate))
	define(env, SymbolFromRawString("failing-io"), newBuiltin0(FailingIO))
	define(env, SymbolFromRawString("eof"), Eof)
	define(env, SymbolFromRawString("stdin"), newBuiltin0(Stdin))
	define(env, SymbolFromRawString("stderr"), newBuiltin0(Stderr))
	define(env, SymbolFromRawString("stdout"), newBuiltin0(Stdout))
	define(env, SymbolFromRawString("buffer"), newBuiltin0(Buffer))
	define(env, SymbolFromRawString("argv"), newBuiltin0(Argv))
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
