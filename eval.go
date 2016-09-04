package main

/*
when panic, when error?
*/

var (
	invalidQuote       = &val{merror, "invalid quote"}
	invalidDef         = &val{merror, "invalid definition"}
	invalidIf          = &val{merror, "invalid if expression"}
	invalidFn          = &val{merror, "invalid function expression"}
	invalidSequence    = &val{merror, "invalid sequence"}
	invalidCond        = &val{merror, "invalid cond expression"}
	invalidApplication = &val{merror, "invalid application"}
	invalidExpression  = &val{merror, "invalid expression"}
)

func isTaggedBy(v, s *val) *val {
	if isPair(v) != vfalse && isSymbol(car(v)) != vfalse && smeq(car(v), s) != vfalse {
		return vtrue
	}

	return vfalse
}

func isQuote(v *val) *val {
	return isTaggedBy(v, sfromString("quote"))
}

func evalQuote(e, v *val) *val {
	if isPair(cdr(v)) == vfalse || cdr(cdr(v)) != vnil {
		fatal(invalidQuote)
	}

	return car(cdr(v))
}

func makeFn(a, b *val) *val {
	return list(sfromString("fn"), a, b)
}

func isDef(v *val) *val {
	return isTaggedBy(v, sfromString("def"))
}

func nameOfDef(v *val) *val {
	return car(cdr(v))
}

func nameOfFunctionDef(v *val) *val {
	if isPair(car(cdr(v))) == vfalse || isSymbol(car(car(cdr(v)))) == vfalse {
		return fatal(invalidDef)
	}

	return car(car(cdr(v)))
}

func defName(v *val) *val {
	if isPair(cdr(v)) == vfalse {
		return fatal(invalidDef)
	}

	if isSymbol(car(cdr(v))) != vfalse {
		return nameOfDef(v)
	}

	return nameOfFunctionDef(v)
}

func valueOfDef(v *val) *val {
	if isPair(cdr(cdr(v))) == vfalse || cdr(cdr(cdr(v))) != vnil {
		return fatal(invalidDef)
	}

	return car(cdr(cdr(v)))
}

func defValue(v *val) *val {
	if isPair(cdr(v)) == vfalse {
		return fatal(invalidDef)
	}

	if isSymbol(car(cdr(v))) != vfalse {
		return valueOfDef(v)
	}

	if isPair(car(cdr(v))) == vfalse {
		return fatal(invalidDef)
	}

	return makeFn(cdr(car(cdr(v))), cdr(cdr(v)))
}

func evalDef(e, v *val) *val {
	return define(e, defName(v), eval(e, defValue(v)))
}

func isIf(v *val) *val {
	return isTaggedBy(v, sfromString("if"))
}

func ifPredicate(v *val) *val {
	if isPair(cdr(v)) == vfalse {
		return fatal(invalidIf)
	}

	return car(cdr(v))
}

func ifConsequent(v *val) *val {
	if isPair(cdr(v)) == vfalse || isPair(cdr(cdr(v))) == vfalse {
		return fatal(invalidIf)
	}

	return car(cdr(cdr(v)))
}

func ifAlternative(v *val) *val {
	if isPair(cdr(v)) == vfalse ||
		isPair(cdr(cdr(v))) == vfalse ||
		isPair(cdr(cdr(cdr(v)))) == vfalse {
		return fatal(invalidIf)
	}

	return car(cdr(cdr(cdr(v))))
}

func evalIf(e, v *val) *val {
	if eval(e, ifPredicate(v)) != vfalse {
		return eval(e, ifConsequent(v))
	}

	return eval(e, ifAlternative(v))
}

func isFn(v *val) *val {
	return isTaggedBy(v, sfromString("fn"))
}

func procParams(v *val) *val {
	if isPair(v) == vfalse || isPair(cdr(v)) == vfalse ||
		isSymbol(car(cdr(v))) == vfalse && isPair(car(cdr(v))) == vfalse {
		return fatal(invalidFn)
	}

	return car(cdr(v))
}

func procBody(v *val) *val {
	if isPair(v) == vfalse || isPair(cdr(v)) == vfalse || isPair(cdr(cdr(v))) == vfalse {
		return fatal(invalidFn)
	}

	return car(cdr(cdr(v)))
}

func fnToProc(e, v *val) *val {
	return newProc(e, procParams(v), procBody(v))
}

func isBegin(v *val) *val {
	return isTaggedBy(v, sfromString("begin"))
}

func beginSeq(v *val) *val {
	if isPair(v) == vfalse {
		return fatal(invalidSequence)
	}

	return cdr(v)
}

func evalSeq(e, v *val) *val {
	if isPair(v) == vfalse {
		return fatal(invalidSequence)
	}

	if cdr(v) == vnil {
		return eval(e, car(v))
	}

	eval(e, car(v))
	return evalSeq(e, cdr(v))
}

func isCond(v *val) *val {
	return isTaggedBy(v, sfromString("cond"))
}

func seqToExp(v *val) *val {
	if isPair(v) == vfalse {
		return fatal(invalidCond)
	}

	if isNil(cdr(v)) == vfalse {
		return car(v)
	}

	return cons(sfromString("begin"), v)
}

func expandCond(v *val) *val {
	if isPair(v) == vfalse {
		return fatal(invalidCond)
	}

	cond := car(v)
	if isPair(v) == vfalse {
		return fatal(invalidCond)
	}

	rest := cdr(v)
	pred := car(cond)

	if smeq(pred, sfromString("else")) != vfalse {
		if isNil(rest) == vfalse {
			return fatal(invalidCond)
		}

		return seqToExp(cdr(cond))
	}

	return list(
		sfromString("if"),
		pred,
		seqToExp(cdr(cond)),
		expandCond(rest),
	)
}

func condToIf(v *val) *val {
	if isPair(v) == vfalse {
		return fatal(invalidCond)
	}

	return expandCond(cdr(v))
}

func isApplication(v *val) *val {
	return isPair(v)
}

func valueList(e, v *val) *val {
	if isNil(v) != vfalse {
		return vnil
	}

	if isPair(v) == vfalse {
		return fatal(invalidApplication)
	}

	return cons(eval(e, car(v)), valueList(e, cdr(v)))
}

func evalApply(e, v *val) *val {
	if isPair(v) == vfalse {
		return fatal(invalidApplication)
	}

	return apply(eval(e, car(v)), valueList(e, cdr(v)))
}

func eval(e, v *val) *val {
	switch {
	case isNumber(v) != vfalse:
		return v
	case isString(v) != vfalse:
		return v
	case isBool(v) != vfalse:
		return v
	case isQuote(v) != vfalse:
		return evalQuote(e, v)
	case isSymbol(v) != vfalse:
		return lookupDef(e, v)
	case isDef(v) != vfalse:
		return evalDef(e, v)
	case isIf(v) != vfalse:
		return evalIf(e, v)
	case isFn(v) != vfalse:
		return fnToProc(e, v)
	case isBegin(v) != vfalse:
		return evalSeq(e, beginSeq(v))
	case isCond(v) != vfalse:
		return eval(e, condToIf(v))
	case isApplication(v) != vfalse:
		return evalApply(e, v)
	default:
		return fatal(invalidExpression)
	}
}
