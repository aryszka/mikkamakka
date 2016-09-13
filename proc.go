package mikkamakka

type builtin func([]*Val) *Val

type Compiled func(*Val) *Val

type proc struct {
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
		procedure,
		&proc{
			builtin:  p,
			argCount: argCount,
			varArgs:  varArgs,
		},
	}
}

func newProc(e, p, b *Val) *Val {
	return &Val{
		procedure,
		&proc{params: p, body: b, env: e},
	}
}

func procString(p *Val) *Val {
	checkType(p, procedure)
	return fromString("<procedure>")
}

func applyBuiltin(p *proc, a *Val) *Val {
	args := make([]*Val, 0, p.argCount)

	for {
		if isNil(a) != vfalse || !p.varArgs && len(args) == p.argCount {
			break
		}

		if isPair(a) == vfalse {
			return fatal(invalidArguments)
		}

		args = append(args, car(a))
		a = cdr(a)
	}

	if isNil(a) == vfalse || !p.varArgs && len(args) != p.argCount || p.varArgs && len(args) < p.argCount {
		return fatal(invalidArguments)
	}

	return p.builtin(args)
}

func applyStruct(s, a *Val) *Val {
	checkType(s, mstruct)
	checkType(a, pair)

	if isNil(cdr(a)) == vfalse {
		return invalidArguments
	}

	return field(s, car(a))
}

func applyLang(p *proc, a *Val) *Val {
	return evalSeq(extendEnv(p.env, p.params, a), p.body)
}

func apply(p, a *Val) *Val {
	checkType(p, procedure, mstruct)
	checkType(a, pair, mnil)

	if isStruct(p) != vfalse {
		return applyStruct(p, a)
	}

	pt := p.value.(*proc)

	if pt.builtin != nil {
		return applyBuiltin(pt, a)
	}

	return applyLang(pt, a)
}

func bapply(a []*Val) *Val {
	return apply(a[0], a[1])
}

func isProc(e *Val) *Val {
	if e.mtype == procedure {
		return vtrue
	}

	return vfalse
}

func toBuiltin(c Compiled) builtin {
	return func(a []*Val) *Val {
		al := Nil
		for i := len(a) - 1; i >= 0; i-- {
			al = cons(a[i], al)
		}

		return c(al)
	}
}

// needs the names
func NewCompiled(p Compiled, argCount int, varArgs bool) *Val {
	return newBuiltin(toBuiltin(p), argCount, varArgs)
}

func Apply(p, a *Val) *Val {
	av := Nil
	for {
		if isNil(a) != vfalse {
			break
		}

		if isPair(a) == vfalse {
			return fatal(invalidArguments)
		}

		av, a = cons(car(a), av), cdr(a)
	}

	return apply(p, reverse(av))
}
