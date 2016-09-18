package mikkamakka

type env struct {
	current map[string]*Val
	parent  *Val
}

var (
	Undefined        = SysStringToError("undefined")
	NotAnEnvironment = SysStringToError("not an environment")
	DefinitionExists = SysStringToError("definition exists")
)

func NewEnv() *Val {
	return newVal(
		Environment,
		&env{
			current: make(map[string]*Val),
			parent:  nil,
		},
	)
}

func IsEnv(e *Val) *Val {
	return is(e, Environment)
}

func LookupDef(e, n *Val) *Val {
	checkType(e, Environment)

	et := e.value.(*env)
	ns := SymbolToSysString(n)
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
			StringToSysString(
				AppendString(
					SymbolToString(n),
					SysStringToString(":"),
					k))),
		v,
	)
}

func defineVector(e, n, v, l *Val) *Val {
	if Eq(l, SysIntToNumber(0)) != False {
		return v
	}

	i := Sub(l, SysIntToNumber(1))
	defineDerived(e, n, NumberToString(i), VectorRef(v, i))
	return defineVector(e, n, v, i)
}

func defineStruct(e, n, s, names *Val) *Val {
	if IsNil(names) != False {
		return s
	}

	k := Car(names)
	defineDerived(e, n, SymbolToString(k), Field(s, k))
	return defineStruct(e, n, s, Cdr(names))
}

func Define(e, n, v *Val) *Val {
	checkType(e, Environment)

	et := e.value.(*env)
	ns := SymbolToSysString(n)
	if _, has := et.current[ns]; has {
		println(SymbolToSysString(n))
		return Fatal(DefinitionExists)
	}

	et.current[ns] = v

	if IsVector(v) != False {
		return defineVector(e, n, v, VectorLen(v))
	}

	if IsStruct(v) != False {
		return defineStruct(e, n, v, StructNames(v))
	}

	return v
}

// TODO: clean this up
func defineAll(e, n, a *Val) *Val {
	for {
		if IsNil(n) != False && IsNil(a) != False {
			break
		}

		if IsPair(a) == False && IsNil(a) == False {
			println("invalid args 3")
			return Fatal(InvalidArgs)
		}

		if IsPair(n) == False {
			Define(e, n, a)
			return e
		}

		if IsNil(a) != False {
			println("invalid args 4", SymbolToSysString(Car(n)))
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
	e = newVal(
		Environment,
		&env{
			current: make(map[string]*Val),
			parent:  e,
		},
	)
	return defineAll(e, n, a)
}

func envString(e *Val) *Val {
	checkType(e, Environment)
	return SysStringToString("<environment>")
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

func newBuiltin1V(f func(*Val, ...*Val) *Val) *Val {
	return NewCompiled(1, true, func(a []*Val) *Val {
		return f(a[0], a[1:]...)
	})
}

func newBuiltin2(f func(*Val, *Val) *Val) *Val {
	return NewCompiled(2, false, func(a []*Val) *Val {
		return f(a[0], a[1])
	})
}

func newBuiltin3(f func(*Val, *Val, *Val) *Val) *Val {
	return NewCompiled(3, false, func(a []*Val) *Val {
		return f(a[0], a[1], a[2])
	})
}

func InitialEnv() *Val {
	env := NewEnv()

	defs := map[string]*Val{
		"nil":                    NilVal,
		"nil?":                   newBuiltin1(IsNil),
		"pair?":                  newBuiltin1(IsPair),
		"cons":                   newBuiltin2(Cons),
		"car":                    newBuiltin1(Car),
		"cdr":                    newBuiltin1(Cdr),
		"list":                   newBuiltin0V(List),
		"error?":                 newBuiltin1(IsError),
		"string->error":          newBuiltin1(StringToError),
		"fatal":                  newBuiltin1(Fatal),
		"yes":                    newBuiltin1(Yes),
		"not":                    newBuiltin1(Not),
		"=":                      newBuiltin0V(Eq),
		">":                      newBuiltin0V(Greater),
		"+":                      newBuiltin0V(Add),
		"-":                      newBuiltin1V(Sub),
		"string->number":         newBuiltin1(StringToNumber),
		"string->bool":           newBuiltin1(StringToBool),
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
		"vector?":                newBuiltin1(IsVector),
		"vector-len":             newBuiltin1(VectorLen),
		"struct?":                newBuiltin1(IsStruct),
		"exit":                   newBuiltin1(Exit),
		"lookup-def":             newBuiltin2(LookupDef),
		"define":                 newBuiltin3(Define),
		"list->vector":           newBuiltin1(ListToVector),
		"list->struct":           newBuiltin1(ListToStruct),
		"make-composite":         newBuiltin1(NewComposite),
		"extend-env":             newBuiltin3(ExtendEnv),
		"vector-ref":             newBuiltin2(VectorRef),
		"field":                  newBuiltin2(Field),
		"struct-names":           newBuiltin1(StructNames),
		"compiled-function?":     newBuiltin1(IsCompiledFunction),
		"composite-function?":    newBuiltin1(IsCompositeFunction),
		"apply-compiled":         newBuiltin2(ApplyCompiled),
		"composite":              newBuiltin1(Composite),
		"sys?":                   newBuiltin1(IsSys),
		"env?":                   newBuiltin1(IsEnv),
		"function?":              newBuiltin1(IsFunction),
		"function->string":       newBuiltin1(FunctionToString),
	}

	for k, v := range defs {
		Define(env, SymbolFromRawString(k), v)
	}

	return env
}
