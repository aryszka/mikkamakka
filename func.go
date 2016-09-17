package mikkamakka

type function func([]*Val) *Val

type fn struct {
	compiled  function
	composite *Val
}

var (
	InvalidArgs  = SysStringToError("invalid arguments")
	NotCompiled  = SysStringToError("not a compiled function")
	NotComposite = SysStringToError("not a composite function")
)

func NewComposite(v *Val) *Val {
	return newVal(Function, &fn{composite: v})
}

// needs the names
func NewCompiled(argCount int, variadic bool, f func([]*Val) *Val) *Val {
	return newVal(
		Function,
		&fn{
			compiled: func(a []*Val) *Val {
				if len(a) < argCount {
					return Fatal(InvalidArgs)
				}

				if !variadic && len(a) != argCount {
					return Fatal(InvalidArgs)
				}

				return f(a)
			}})
}

func IsCompiledFunction(e *Val) *Val {
	if e.typ == Function && e.value.(*fn).compiled != nil {
		return True
	}

	return False
}

func IsFunction(v *Val) *Val {
	if v.typ == Function {
		return True
	}

	return False
}

func FunctionToString(f *Val) *Val {
	checkType(f, Function)
	return SysStringToString("<function>")
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
	checkType(f, Function)

	ft := f.value.(*fn)
	if ft.compiled == nil {
		return Fatal(NotCompiled)
	}

	return ft.compiled(listToSlice(a))
}

func Composite(f *Val) *Val {
	checkType(f, Function)

	ft := f.value.(*fn)
	if ft.composite == nil {
		return Fatal(NotComposite)
	}

	return ft.composite
}
