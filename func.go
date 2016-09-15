package mikkamakka

type Function func(...*Val) *Val

type fn struct {
	compiled Function
	params   *Val
	body     *Val
	env      *Val
}

var invalidArgs = &Val{merror, "invalid arguments"}

func newFn(e, p, b *Val) *Val {
	return &Val{
		function,
		&fn{params: p, body: b, env: e},
	}
}

func fnString(p *Val) *Val {
	checkType(p, function)
	return fromString("<function>")
}

func applyBuiltin(p *fn, a *Val) *Val {
	return p.compiled(listToSlice(a)...)
}

func applyStruct(s, a *Val) *Val {
	checkType(s, mstruct)
	checkType(a, pair)

	if isNil(cdr(a)) == False {
		return invalidArgs
	}

	return field(s, car(a))
}

func applyLang(p *fn, a *Val) *Val {
	return evalSeq(extendEnv(p.env, p.params, a), p.body)
}

func Apply(p, a *Val) *Val {
	checkType(p, function, mstruct)
	checkType(a, pair, mnil)

	if isStruct(p) != False {
		return applyStruct(p, a)
	}

	pt := p.value.(*fn)

	if pt.compiled != nil {
		return applyBuiltin(pt, a)
	}

	return applyLang(pt, a)
}

func isFn(e *Val) *Val {
	if e.mtype == function {
		return True
	}

	return False
}

// needs the names
func NewCompiled(argCount int, variadic bool, f func([]*Val) *Val) *Val {
	return &Val{function, &fn{compiled: func(a ...*Val) *Val {
		if len(a) < argCount {
			return fatal(invalidArgs)
		}

		if !variadic && len(a) != argCount {
			return fatal(invalidArgs)
		}

		return f(a)
	}}}
}
