package main

type builtin func([]*val) *val

type proc struct {
	builtin  builtin
	argCount int
	varArgs  bool
	params   *val
	body     *val
	env      *val
}

var invalidArguments = &val{merror, "invalid arguments"}

func newBuiltin(p builtin, argCount int, varArgs bool) *val {
	return &val{
		procedure,
		&proc{
			builtin:  p,
			argCount: argCount,
			varArgs:  varArgs,
		},
	}
}

func newProc(e, p, b *val) *val {
	return &val{
		procedure,
		&proc{params: p, body: b, env: e},
	}
}

func procString(p *val) *val {
	checkType(p, procedure)
	return fromString("<procedure>")
}

func applyBuiltin(p *proc, a *val) *val {
	args := make([]*val, 0, p.argCount)

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

func applyStruct(s, a *val) *val {
	checkType(s, mstruct)
	checkType(a, pair)

	if isNil(cdr(a)) == vfalse {
		return invalidArguments
	}

	return field(s, car(a))
}

func applyLang(p *proc, a *val) *val {
	return evalSeq(extendEnv(p.env, p.params, a), p.body)
}

func apply(p, a *val) *val {
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

func bapply(a []*val) *val {
	return apply(a[0], a[1])
}

func isProc(e *val) *val {
	if e.mtype == procedure {
		return vtrue
	}

	return vfalse
}
