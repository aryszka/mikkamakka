package mikkamakka

type Function func([]*Val) *Val

type fn struct {
	compiled  Function
	composite *Val
}

var (
	InvalidArgs  = &Val{merror, "invalid arguments"}
	NotCompiled  = &Val{merror, "not a compiled function"}
	NotComposite = &Val{merror, "not a composite function"}
)

func NewComposite(v *Val) *Val {
	return &Val{function, &fn{composite: v}}
}

// needs the names
func NewCompiled(argCount int, variadic bool, f Function) *Val {
	return &Val{
		function,
		&fn{
			compiled: func(a []*Val) *Val {
				if len(a) < argCount {
					return fatal(InvalidArgs)
				}

				if !variadic && len(a) != argCount {
					return fatal(InvalidArgs)
				}

				return f(a)
			}}}
}

func IsCompiledFunction(e *Val) *Val {
	if e.mtype == function && e.value.(*fn).compiled != nil {
		return True
	}

	return False
}

func IsFunction(v *Val) *Val {
	if v.mtype == function {
		return True
	}

	return False
}

func FunctionToString(f *Val) *Val {
	checkType(f, function)
	return fromString("<function>")
}

func listToSlice(l *Val) []*Val {
	var s []*Val
	for {
		if IsNil(l) != False {
			break
		}

		s, l = append(s, Car(l)), Cdr(l)
	}

	return s
}

func ApplyCompiled(f, a *Val) *Val {
	checkType(f, function)
	checkType(a, pair, mnil)

	ft := f.value.(*fn)
	if ft.compiled == nil {
		return fatal(NotCompiled)
	}

	return ft.compiled(listToSlice(a))
}

func Composite(f *Val) *Val {
	checkType(f, function)

	ft := f.value.(*fn)
	if ft.composite == nil {
		return fatal(NotComposite)
	}

	return ft.composite
}