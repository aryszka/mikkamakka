package main

type builtin func([]*val) *val

type proc struct {
	builtin  builtin
	argCount int
	params   *val
	body     *val
	env      *val
}

var invalidArguments = &val{merror, "invalid arguments"}

func newBuiltin(p builtin, argCount int) *val {
	return &val{
		procedure,
		&proc{
			builtin:  p,
			argCount: argCount,
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
		if isNil(a) != vfalse || len(args) == p.argCount {
			break
		}

		if isPair(a) == vfalse {
			return fatal(invalidArguments)
		}

		args = append(args, car(a))
	}

	if isNil(a) == vfalse || len(args) != p.argCount {
		return fatal(invalidArguments)
	}

	return p.builtin(args)
}

func applyLang(p *proc, a *val) *val {
	return evalSeq(extendEnv(p.env, p.params, a), p.body)
}

func apply(p, a *val) *val {
	checkType(p, procedure)
	checkType(a, pair, mnil)
	pt := p.value.(*proc)

	if pt.builtin != nil {
		return applyBuiltin(pt, a)
	}

	return applyLang(pt, a)
}

func isProc(e *val) *val {
	if e.mtype == procedure {
		return vtrue
	}

	return vfalse
}
