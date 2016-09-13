package mikkamakka

func printer(out *Val) *Val {
	return fromMap(map[string]*Val{
		"output": out,
		"state":  Nil,
	})
}

func bprinter(a []*Val) *Val {
	return printer(a[0])
}

func printState(p *Val) *Val {
	return field(p, sfromString("state"))
}

func printRaw(p *Val, r *Val) *Val {
	f := fwrite(field(p, sfromString("output")), r)
	return assign(p, fromMap(map[string]*Val{
		"output": f,
		"state":  fstate(f),
	}))
}

func printQuoteSign(p *Val) *Val {
	return printRaw(p, fromString("'"))
}

func printSymbol(p, v, q *Val) *Val {
	if q == vfalse {
		p = printQuoteSign(p)
	}

	return printRaw(p, symbolToString(v))
}

func printQuote(p, v *Val) *Val {
	p = printQuoteSign(p)
	return mprintq(p, car(cdr(v)), vfalse)
}

func printPair(p, v, q *Val) *Val {
	if q == vfalse {
		p = printQuoteSign(p)
	}

	p = printRaw(p, fromString("("))
	if st := printState(p); isError(st) != vfalse {
		return p
	}

	var loop func(*Val, *Val, *Val) *Val
	loop = func(p *Val, v *Val, first *Val) *Val {
		if isNil(v) != vfalse {
			return printRaw(p, fromString(")"))
		}

		if first == vfalse {
			p = printRaw(p, fromString(" "))
			if st := printState(p); isError(st) != vfalse {
				return p
			}
		}

		if isPair(cdr(v)) == vfalse && isNil(cdr(v)) == vfalse {
			p = mprintq(p, car(v), vtrue)
			if st := field(p, sfromString("state")); isError(st) != vfalse {
				return p
			}

			p = printRaw(p, fromString(" . "))
			if st := printState(p); isError(st) != vfalse {
				return p
			}

			p = mprintq(p, cdr(v), vtrue)
			if st := field(p, sfromString("state")); isError(st) != vfalse {
				return p
			}

			return printRaw(p, fromString(")"))
		}

		p = mprintq(p, car(v), vtrue)
		if st := field(p, sfromString("state")); isError(st) != vfalse {
			return p
		}

		return loop(p, cdr(v), vfalse)
	}

	return loop(p, v, vtrue)
}

func printVector(p, v *Val) *Val {
	p = printRaw(p, fromString("["))
	if st := field(p, sfromString("state")); isError(st) != vfalse {
		return p
	}

	var loop func(*Val, *Val, *Val) *Val
	loop = func(p, i, f *Val) *Val {
		if neq(i, vectorLength(v)) != vfalse {
			return p
		}

		if f == vfalse {
			p = printRaw(p, fromString(" "))
			if st := field(p, sfromString("state")); isError(st) != vfalse {
				return p
			}
		}

		p = mprintq(p, vectorRef(v, i), vtrue)
		if st := field(p, sfromString("state")); isError(st) != vfalse {
			return p
		}

		return loop(p, add(i, fromInt(1)), vfalse)
	}

	p = loop(p, fromInt(0), vtrue)
	return printRaw(p, fromString("]"))
}

func printStruct(p, v *Val) *Val {
	p = printRaw(p, fromString("{"))
	if st := field(p, sfromString("state")); isError(st) != vfalse {
		return p
	}

	var loop func(*Val, *Val, *Val) *Val
	loop = func(p, n, f *Val) *Val {
		if n == Nil {
			return p
		}

		if f == vfalse {
			p = printRaw(p, fromString(" "))
			if st := field(p, sfromString("state")); isError(st) != vfalse {
				return p
			}
		}

		p = mprintq(p, car(n), vtrue)
		if st := field(p, sfromString("state")); isError(st) != vfalse {
			return p
		}

		p = printRaw(p, fromString(" "))
		if st := field(p, sfromString("state")); isError(st) != vfalse {
			return p
		}

		p = mprintq(p, field(v, car(n)), vtrue)
		if st := field(p, sfromString("state")); isError(st) != vfalse {
			return p
		}

		return loop(p, cdr(n), vfalse)
	}

	p = loop(p, structNames(v), vtrue)
	return printRaw(p, fromString("}"))
}

func mprintq(p, v, q *Val) *Val {
	if isSymbol(v) != vfalse {
		return printSymbol(p, v, q)
	} else if isNumber(v) != vfalse {
		v = numberToString(v)
	} else if isString(v) != vfalse {
		v = appendString(fromString(`"`), v, fromString(`"`))
	} else if isBool(v) != vfalse {
		v = boolToString(v)
	} else if isSys(v) != vfalse {
		v = sstring(v)
	} else if isError(v) != vfalse {
		v = estring(v)
	} else if isPair(v) != vfalse && isSymbol(car(v)) != vfalse && smeq(car(v), sfromString("quote")) != vfalse {
		return printQuote(p, v)
	} else if isPair(v) != vfalse || isNil(v) != vfalse {
		return printPair(p, v, q)
	} else if isVector(v) != vfalse {
		return printVector(p, v)
	} else if isStruct(v) != vfalse {
		return printStruct(p, v)
	} else if isEnv(v) != vfalse {
		v = envString(v)
	} else if isProc(v) != vfalse {
		v = procString(v)
	} else {
		return assign(p, fromMap(map[string]*Val{
			"state": notImplemented,
		}))
	}

	f := fwrite(field(p, sfromString("output")), v)
	if st := fstate(f); isError(st) != vfalse {
		return assign(p, fromMap(map[string]*Val{
			"output": f,
			"state":  st,
		}))
	}

	return assign(p, fromMap(map[string]*Val{
		"output": f,
		"state":  v,
	}))
}

func mprint(p, v *Val) *Val {
	return mprintq(p, v, vfalse)
}

func bprint(a []*Val) *Val {
	return mprint(a[0], a[1])
}
