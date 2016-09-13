package mikkamakka

type builtin func([]*Val) *Val

type Function func(*Val) *Val

type fn struct {
	builtin  builtin
	argCount int
	varArgs  bool
	params   *Val
	body     *Val
	env      *Val
}

var invalidArguments = &Val{merror, "invalid arguments"}

func newBuiltin(p builtin, argCount int, varArgs bool) *Val {
	return &Val{
		function,
		&fn{
			builtin:  p,
			argCount: argCount,
			varArgs:  varArgs,
		},
	}
}

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
	args := make([]*Val, 0, p.argCount)

	for {
		if isNil(a) != False || !p.varArgs && len(args) == p.argCount {
			break
		}

		if isPair(a) == False {
			return fatal(invalidArguments)
		}

		args = append(args, car(a))
		a = cdr(a)
	}

	if isNil(a) == False || !p.varArgs && len(args) != p.argCount || p.varArgs && len(args) < p.argCount {
		return fatal(invalidArguments)
	}

	return p.builtin(args)
}

func applyStruct(s, a *Val) *Val {
	checkType(s, mstruct)
	checkType(a, pair)

	if isNil(cdr(a)) == False {
		return invalidArguments
	}

	return field(s, car(a))
}

func applyLang(p *fn, a *Val) *Val {
	return evalSeq(extendEnv(p.env, p.params, a), p.body)
}

func apply(p, a *Val) *Val {
	checkType(p, function, mstruct)
	checkType(a, pair, mnil)

	if isStruct(p) != False {
		return applyStruct(p, a)
	}

	pt := p.value.(*fn)

	if pt.builtin != nil {
		return applyBuiltin(pt, a)
	}

	return applyLang(pt, a)
}

func bapply(a []*Val) *Val {
	return apply(a[0], a[1])
}

func isFn(e *Val) *Val {
	if e.mtype == function {
		return True
	}

	return False
}

func toBuiltin(c Function) builtin {
	return func(a []*Val) *Val {
		al := Nil
		for i := len(a) - 1; i >= 0; i-- {
			al = cons(a[i], al)
		}

		return c(al)
	}
}

// needs the names
func NewCompiled(p Function, argCount int, varArgs bool) *Val {
	return newBuiltin(toBuiltin(p), argCount, varArgs)
}

func Apply(a *Val) *Val {
	p := car(a)
	a = cdr(a)

	av := Nil
	for {
		if isNil(a) != False {
			break
		}

		if isPair(a) == False {
			return fatal(invalidArguments)
		}

		av, a = cons(car(a), av), cdr(a)
	}

	return apply(p, reverse(av))
}
