package mikkamakka

type module struct {
	all  map[string]*Val
	path *Val
}

type env struct {
	current         map[string]*Val
	parent          *Val
	export          map[string]*Val
	module          *module
	compiledModules map[string]func(*Val)
}

var (
	Undefined        = SysStringToError("undefined")
	NotAnEnvironment = SysStringToError("not an environment")
	DefinitionExists = SysStringToError("definition exists")
	UndefinedModule  = SysStringToError("undefined module")
)

func newEnv(p *Val, m *module) *Val {
	if m == nil {
		m = &module{make(map[string]*Val), NilVal}
	}

	return newVal(
		Environment,
		&env{
			current:         make(map[string]*Val),
			parent:          p,
			module:          m,
			export:          make(map[string]*Val),
			compiledModules: make(map[string]func(*Val)),
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

func DefineAll(e, s *Val) *Val {
	n := StructNames(s)
	for {
		if IsNil(n) != False {
			break
		}

		Define(e, Car(n), Field(s, Car(n)))
		n = Cdr(n)
	}

	return s
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
	checkType(e, Environment)
	env := e.value.(*env)
	e = newEnv(e, env.module)
	return defineAll(e, n, a)
}

func Export(e, s *Val) *Val {
	checkType(e, Environment)
	env := e.value.(*env)

	n := StructNames(s)
	for {
		if IsNil(n) != False {
			return s
		}

		env.export[SymbolToSysString(Car(n))], n = Field(s, Car(n)), Cdr(n)
	}
}

func Exports(e *Val) *Val {
	checkType(e, Environment)
	env := e.value.(*env)
	return SysMapToStruct(env.export)
}

func ModulePath(e *Val) *Val {
	checkType(e, Environment)
	env := e.value.(*env)
	return env.module.path
}

func ModuleEnv(e, n *Val) *Val {
	checkType(e, Environment)
	env := e.value.(*env)
	return newEnv(nil, &module{env.module.all, Cons(n, env.module.path)})
}

func LoadedModule(e, n *Val) *Val {
	checkType(e, Environment)
	env := e.value.(*env)
	if exp, ok := env.module.all[StringToSysString(n)]; ok {
		return exp
	}

	return UndefinedModule
}

func StoreModule(e, n, m *Val) *Val {
	checkType(e, Environment)
	env := e.value.(*env)
	env.module.all[StringToSysString(n)] = m
	return m
}

func ModuleLoader(e, n *Val, m func(*Val)) {
	checkType(e, Environment)
	env := e.value.(*env)
	env.compiledModules[StringToSysString(n)] = m
}

func LoadCompiledModule(e, n *Val) {
	checkType(e, Environment)
	env := e.value.(*env)
	env.compiledModules[StringToSysString(n)](e)
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
	env := newEnv(nil, nil)

	defs := map[string]*Val{
		"nil":                    NilVal,
		"nil?":                   newBuiltin1(IsNil),
		"pair?":                  newBuiltin1(IsPair),
		"cons":                   newBuiltin2(Cons),
		"car":                    newBuiltin1(Car),
		"cdr":                    newBuiltin1(Cdr),
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
		"module-export":          newBuiltin2(Export),
		"module-path":            newBuiltin1(ModulePath),
		"loaded-module":          newBuiltin2(LoadedModule),
		"undefined-module":       UndefinedModule,
		"module-env":             newBuiltin2(ModuleEnv),
		"exports":                newBuiltin1(Exports),
		"store-module":           newBuiltin3(StoreModule),
		"error->string":          newBuiltin1(ErrorToString),
		"sys->string":            newBuiltin1(SysToString),
	}

	for k, v := range defs {
		Define(env, SymbolFromRawString(k), v)
	}

	return env
}
