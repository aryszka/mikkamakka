package mikkamakka

func printer(out *Val) *Val {
	return fromMap(map[string]*Val{
		"output": out,
		"state":  Nil,
	})
}

func printState(p *Val) *Val {
	return field(p, sfromString("state"))
}

func printRaw(p *Val, r *Val) *Val {
	f := fwrite(field(p, sfromString("output")), r)
	return Assign(p, fromMap(map[string]*Val{
		"output": f,
		"state":  fstate(f),
	}))
}

func printQuoteSign(p *Val) *Val {
	return printRaw(p, fromString("'"))
}

func printSymbol(p, v, q *Val) *Val {
	if q == False {
		p = printQuoteSign(p)
	}

	return printRaw(p, symbolToString(v))
}

func printQuote(p, v *Val) *Val {
	p = printQuoteSign(p)
	return mprintq(p, car(cdr(v)), False)
}

func printPair(p, v, q *Val) *Val {
	if q == False {
		p = printQuoteSign(p)
	}

	p = printRaw(p, fromString("("))
	if st := printState(p); isError(st) != False {
		return p
	}

	var loop func(*Val, *Val, *Val) *Val
	loop = func(p *Val, v *Val, first *Val) *Val {
		if isNil(v) != False {
			return printRaw(p, fromString(")"))
		}

		if first == False {
			p = printRaw(p, fromString(" "))
			if st := printState(p); isError(st) != False {
				return p
			}
		}

		if isPair(cdr(v)) == False && isNil(cdr(v)) == False {
			p = mprintq(p, car(v), True)
			if st := field(p, sfromString("state")); isError(st) != False {
				return p
			}

			p = printRaw(p, fromString(" . "))
			if st := printState(p); isError(st) != False {
				return p
			}

			p = mprintq(p, cdr(v), True)
			if st := field(p, sfromString("state")); isError(st) != False {
				return p
			}

			return printRaw(p, fromString(")"))
		}

		p = mprintq(p, car(v), True)
		if st := field(p, sfromString("state")); isError(st) != False {
			return p
		}

		return loop(p, cdr(v), False)
	}

	return loop(p, v, True)
}

func printVector(p, v *Val) *Val {
	p = printRaw(p, fromString("["))
	if st := field(p, sfromString("state")); isError(st) != False {
		return p
	}

	var loop func(*Val, *Val, *Val) *Val
	loop = func(p, i, f *Val) *Val {
		if neq(i, vectorLength(v)) != False {
			return p
		}

		if f == False {
			p = printRaw(p, fromString(" "))
			if st := field(p, sfromString("state")); isError(st) != False {
				return p
			}
		}

		p = mprintq(p, vectorRef(v, i), True)
		if st := field(p, sfromString("state")); isError(st) != False {
			return p
		}

		return loop(p, add(i, fromInt(1)), False)
	}

	p = loop(p, fromInt(0), True)
	return printRaw(p, fromString("]"))
}

func printStruct(p, v *Val) *Val {
	p = printRaw(p, fromString("{"))
	if st := field(p, sfromString("state")); isError(st) != False {
		return p
	}

	var loop func(*Val, *Val, *Val) *Val
	loop = func(p, n, f *Val) *Val {
		if n == Nil {
			return p
		}

		if f == False {
			p = printRaw(p, fromString(" "))
			if st := field(p, sfromString("state")); isError(st) != False {
				return p
			}
		}

		p = mprintq(p, car(n), True)
		if st := field(p, sfromString("state")); isError(st) != False {
			return p
		}

		p = printRaw(p, fromString(" "))
		if st := field(p, sfromString("state")); isError(st) != False {
			return p
		}

		p = mprintq(p, field(v, car(n)), True)
		if st := field(p, sfromString("state")); isError(st) != False {
			return p
		}

		return loop(p, cdr(n), False)
	}

	p = loop(p, structNames(v), True)
	return printRaw(p, fromString("}"))
}

func mprintq(p, v, q *Val) *Val {
	if isSymbol(v) != False {
		return printSymbol(p, v, q)
	} else if isNumber(v) != False {
		v = numberToString(v)
	} else if isString(v) != False {
		v = appendString(fromString(`"`), v, fromString(`"`))
	} else if isBool(v) != False {
		v = boolToString(v)
	} else if isSys(v) != False {
		v = sstring(v)
	} else if isError(v) != False {
		v = estring(v)
	} else if isPair(v) != False && isSymbol(car(v)) != False && smeq(car(v), sfromString("quote")) != False {
		return printQuote(p, v)
	} else if isPair(v) != False || isNil(v) != False {
		return printPair(p, v, q)
	} else if isVector(v) != False {
		return printVector(p, v)
	} else if isStruct(v) != False {
		return printStruct(p, v)
	} else if isEnv(v) != False {
		v = envString(v)
	} else if IsFunction(v) != False {
		v = FunctionToString(v)
	} else {
		return Assign(p, fromMap(map[string]*Val{
			"state": notImplemented,
		}))
	}

	f := fwrite(field(p, sfromString("output")), v)
	if st := fstate(f); isError(st) != False {
		return Assign(p, fromMap(map[string]*Val{
			"output": f,
			"state":  st,
		}))
	}

	return Assign(p, fromMap(map[string]*Val{
		"output": f,
		"state":  v,
	}))
}

func mprint(p, v *Val) *Val {
	return mprintq(p, v, False)
}
