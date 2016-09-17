package mikkamakka

type env struct {
	current Struct
	parent  *Val
}

var (
	Undefined        = ErrorFromRawString("undefined")
	NotAnEnvironment = ErrorFromRawString("not an environment")
	DefinitionExists = ErrorFromRawString("definition exists")
)

func NewEnv() *Val {
	return &Val{
		environment,
		&env{
			current: make(Struct),
			parent:  nil,
		},
	}
}

func IsEnv(e *Val) *Val {
	return is(e, environment)
}

func LookupDef(e, n *Val) *Val {
	checkType(e, environment)

	et := e.value.(*env)
	ns := RawSymbolString(n)
	if v, ok := et.current[ns]; ok {
		return v
	}

	if et.parent == nil {
		println("undefined reference", ns)
		return Fatal(Undefined)
	}

	return LookupDef(et.parent, n)
}

func defineDerived(e, n, k, v *Val) *Val {
	return Define(
		e,
		SymbolFromRawString(
			RawString(
				AppendString(
					SymbolToString(n),
					StringFromRaw(":"),
					k))),
		v,
	)
}

func defineStruct(e, n, s, names *Val) *Val {
	if IsNil(names) != False {
		return s
	}

	k := Car(names)
	defineDerived(e, n, SymbolToString(k), Field(s, k))
	return defineStruct(e, n, s, Cdr(names))
}

func defineVector(e, n, v, l *Val) *Val {
	if Eq(l, NumberFromRawInt(0)) != False {
		return v
	}

	i := Sub(l, NumberFromRawInt(1))
	defineDerived(e, n, i, VectorRef(v, i))
	return defineVector(e, n, v, i)
}

func Define(e, n, v *Val) *Val {
	checkType(e, environment)
	et, ok := e.value.(*env)
	if !ok {
		return Fatal(NotAnEnvironment)
	}

	ns := RawSymbolString(n)
	if _, has := et.current[ns]; has {
		return Fatal(DefinitionExists)
	}

	et.current[ns] = v

	if IsStruct(v) != False {
		return defineStruct(e, n, v, StructNames(v))
	}

	if IsVector(v) != False {
		return defineVector(e, n, v, VectorLen(v))
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
			Define(e, n, a)
			return e
		}

		if IsNil(a) != False {
			return Fatal(InvalidArgs)
		}

		ni := Car(n)
		ai := Car(a)
		Define(e, ni, ai)
		n, a = Cdr(n), Cdr(a)
	}

	return e
}

func ExtendEnv(e, n, a *Val) *Val {
	e = &Val{
		environment,
		&env{
			current: make(Struct),
			parent:  e,
		},
	}
	return defineAll(e, n, a)
}

func envString(e *Val) *Val {
	checkType(e, environment)
	return StringFromRaw("<environment>")
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
	env := NewEnv()

	defs := map[string]*Val{
		"nil":                    Nil,
		"nil?":                   newBuiltin1(IsNil),
		"pair?":                  newBuiltin1(IsPair),
		"cons":                   newBuiltin2(Cons),
		"car":                    newBuiltin1(Car),
		"cdr":                    newBuiltin1(Cdr),
		"list":                   newBuiltin0V(List),
		"apply":                  newBuiltin2(Apply),
		"error?":                 newBuiltin1(IsError),
		"string->error":          newBuiltin1(StringToError),
		"fatal":                  newBuiltin1(Fatal),
		"not":                    newBuiltin1(Not),
		"=":                      newBuiltin0V(Eq),
		">":                      newBuiltin0V(Greater),
		"+":                      newBuiltin0V(Add),
		"string->number":         newBuiltin1(NumberFromString),
		"string->bool":           newBuiltin1(BoolFromString),
		"symbol?":                newBuiltin1(IsSymbol),
		"symbol->string":         newBuiltin1(SymbolToString),
		"string->symbol":         newBuiltin1(SymbolFromString),
		"number?":                newBuiltin1(IsNumber),
		"number->string":         newBuiltin1(NumberToString),
		"bool?":                  newBuiltin1(IsBool),
		"bool->string":           newBuiltin1(BoolToString),
		"string?":                newBuiltin1(IsString),
		"assign":                 newBuiltin0V(Assign),
		"fopen":                  newBuiltin1(Fopen),
		"fclose":                 newBuiltin1(Fclose),
		"fread":                  newBuiltin2(Fread),
		"fwrite":                 newBuiltin2(Fwrite),
		"fstate":                 newBuiltin1(Fstate),
		"failing-io":             newBuiltin0(FailingIO),
		"eof":                    Eof,
		"stdin":                  newBuiltin0(Stdin),
		"stderr":                 newBuiltin0(Stderr),
		"stdout":                 newBuiltin0(Stdout),
		"buffer":                 newBuiltin0(Buffer),
		"argv":                   newBuiltin0(Argv),
		"invalid-token":          invalidToken,
		"string-append":          newBuiltin0V(AppendString),
		"escape-compiled-string": newBuiltin1(EscapeCompiled),
		"printer":                newBuiltin1(printer),
		"print":                  newBuiltin2(mprint),
	}

	for k, v := range defs {
		Define(env, SymbolFromRawString(k), v)
	}

	return env
}
