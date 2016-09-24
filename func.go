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
					println("invalid args 1")
					return Fatal(InvalidArgs)
				}

				if !variadic && len(a) != argCount {
					println("invalid args 2", variadic, len(a), argCount)
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

func IsCompositeFunction(e *Val) *Val {
	if e.typ == Function && e.value.(*fn).composite != nil {
		return True
	}

	return False
}

func IsFunction(v *Val) *Val {
	return is(v, Function)
}

func FunctionToString(f *Val) *Val {
	checkType(f, Function)
	return SysStringToString("<function>")
}

func ListToSlice(l *Val) []*Val {
	var s []*Val
	for {
		if IsNil(l) != False {
			break
		}

		s, l = append(s, Car(l)), Cdr(l)
	}

	return s
}

func SliceToList(s []*Val) *Val {
	l := NilVal
	for i := len(s) - 1; i >= 0; i-- {
		l = Cons(s[i], l)
	}

	return l
}

func ApplyCompiled(f, a *Val) *Val {
	checkType(f, Function)

	ft := f.value.(*fn)
	if ft.compiled == nil {
		return Fatal(NotCompiled)
	}

	return ft.compiled(ListToSlice(a))
}

func ApplySys(f, a *Val) *Val {
	if IsVector(f) != False {
		return VectorRef(f, Car(a))
	}

	if IsStruct(f) != False {
		return Field(f, Car(a))
	}

	return ApplyCompiled(f, a)
}

func Composite(f *Val) *Val {
	checkType(f, Function)

	ft := f.value.(*fn)
	if ft.composite == nil {
		return Fatal(NotComposite)
	}

	return ft.composite
}
