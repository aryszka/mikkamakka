package mikkamakka

func printer(out *Val) *Val {
	return FromMap(Struct{
		"output": out,
		"state":  Nil,
	})
}

func printState(p *Val) *Val {
	return Field(p, SymbolFromRawString("state"))
}

func printRaw(p *Val, r *Val) *Val {
	f := fwrite(Field(p, SymbolFromRawString("output")), r)
	return Assign(p, FromMap(Struct{
		"output": f,
		"state":  fstate(f),
	}))
}

func printQuoteSign(p *Val) *Val {
	return printRaw(p, StringFromRaw("'"))
}

func printSymbol(p, v, q *Val) *Val {
	if q == False {
		p = printQuoteSign(p)
	}

	return printRaw(p, SymbolToString(v))
}

func printQuote(p, v *Val) *Val {
	p = printQuoteSign(p)
	return mprintq(p, Car(Cdr(v)), False)
}

func printPair(p, v, q *Val) *Val {
	if q == False {
		p = printQuoteSign(p)
	}

	p = printRaw(p, StringFromRaw("("))
	if st := printState(p); isError(st) != False {
		return p
	}

	var loop func(*Val, *Val, *Val) *Val
	loop = func(p *Val, v *Val, first *Val) *Val {
		if IsNil(v) != False {
			return printRaw(p, StringFromRaw(")"))
		}

		if first == False {
			p = printRaw(p, StringFromRaw(" "))
			if st := printState(p); isError(st) != False {
				return p
			}
		}

		if IsPair(Cdr(v)) == False && IsNil(Cdr(v)) == False {
			p = mprintq(p, Car(v), True)
			if st := Field(p, SymbolFromRawString("state")); isError(st) != False {
				return p
			}

			p = printRaw(p, StringFromRaw(" . "))
			if st := printState(p); isError(st) != False {
				return p
			}

			p = mprintq(p, Cdr(v), True)
			if st := Field(p, SymbolFromRawString("state")); isError(st) != False {
				return p
			}

			return printRaw(p, StringFromRaw(")"))
		}

		p = mprintq(p, Car(v), True)
		if st := Field(p, SymbolFromRawString("state")); isError(st) != False {
			return p
		}

		return loop(p, Cdr(v), False)
	}

	return loop(p, v, True)
}

func printVector(p, v *Val) *Val {
	p = printRaw(p, StringFromRaw("["))
	if st := Field(p, SymbolFromRawString("state")); isError(st) != False {
		return p
	}

	var loop func(*Val, *Val, *Val) *Val
	loop = func(p, i, f *Val) *Val {
		if numberEq(i, VectorLength(v)) != False {
			return p
		}

		if f == False {
			p = printRaw(p, StringFromRaw(" "))
			if st := Field(p, SymbolFromRawString("state")); isError(st) != False {
				return p
			}
		}

		p = mprintq(p, VectorRef(v, i), True)
		if st := Field(p, SymbolFromRawString("state")); isError(st) != False {
			return p
		}

		return loop(p, Add(i, NumberFromRawInt(1)), False)
	}

	p = loop(p, NumberFromRawInt(0), True)
	return printRaw(p, StringFromRaw("]"))
}

func printStruct(p, v *Val) *Val {
	p = printRaw(p, StringFromRaw("{"))
	if st := Field(p, SymbolFromRawString("state")); isError(st) != False {
		return p
	}

	var loop func(*Val, *Val, *Val) *Val
	loop = func(p, n, f *Val) *Val {
		if n == Nil {
			return p
		}

		if f == False {
			p = printRaw(p, StringFromRaw(" "))
			if st := Field(p, SymbolFromRawString("state")); isError(st) != False {
				return p
			}
		}

		p = mprintq(p, Car(n), True)
		if st := Field(p, SymbolFromRawString("state")); isError(st) != False {
			return p
		}

		p = printRaw(p, StringFromRaw(" "))
		if st := Field(p, SymbolFromRawString("state")); isError(st) != False {
			return p
		}

		p = mprintq(p, Field(v, Car(n)), True)
		if st := Field(p, SymbolFromRawString("state")); isError(st) != False {
			return p
		}

		return loop(p, Cdr(n), False)
	}

	p = loop(p, StructNames(v), True)
	return printRaw(p, StringFromRaw("}"))
}

func mprintq(p, v, q *Val) *Val {
	if IsSymbol(v) != False {
		return printSymbol(p, v, q)
	} else if IsNumber(v) != False {
		v = NumberToString(v)
	} else if IsString(v) != False {
		v = AppendString(StringFromRaw(`"`), v, StringFromRaw(`"`))
	} else if IsBool(v) != False {
		v = BoolToString(v)
	} else if isSys(v) != False {
		v = sstring(v)
	} else if isError(v) != False {
		v = estring(v)
	} else if IsPair(v) != False && IsSymbol(Car(v)) != False && Eq(Car(v), SymbolFromRawString("quote")) != False {
		return printQuote(p, v)
	} else if IsPair(v) != False || IsNil(v) != False {
		return printPair(p, v, q)
	} else if IsVector(v) != False {
		return printVector(p, v)
	} else if IsStruct(v) != False {
		return printStruct(p, v)
	} else if isEnv(v) != False {
		v = envString(v)
	} else if IsFunction(v) != False {
		v = FunctionToString(v)
	} else {
		return Assign(p, FromMap(Struct{
			"state": notImplemented,
		}))
	}

	f := fwrite(Field(p, SymbolFromRawString("output")), v)
	if st := fstate(f); isError(st) != False {
		return Assign(p, FromMap(Struct{
			"output": f,
			"state":  st,
		}))
	}

	return Assign(p, FromMap(Struct{
		"output": f,
		"state":  v,
	}))
}

func mprint(p, v *Val) *Val {
	return mprintq(p, v, False)
}
