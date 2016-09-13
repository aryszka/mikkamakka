package mikkamakka

/*
when panic, when error?
*/

var (
	invalidQuote         = &Val{merror, "invalid quote"}
	invalidDef           = &Val{merror, "invalid definition"}
	invalidIf            = &Val{merror, "invalid if expression"}
	invalidAnd           = &Val{merror, "invalid and expression"}
	invalidOr            = &Val{merror, "invalid or expression"}
	invalidFn            = &Val{merror, "invalid function expression"}
	invalidSequence      = &Val{merror, "invalid sequence"}
	invalidVector        = &Val{merror, "invalid vector"}
	invalidCond          = &Val{merror, "invalid cond expression"}
	invalidLet           = &Val{merror, "invalid let expression"}
	invalidTest          = &Val{merror, "invalid test"}
	invalidApplication   = &Val{merror, "invalid application"}
	invalidExpression    = &Val{merror, "invalid expression"}
	definitionExpression = &Val{merror, "definition in expression position"}
	testFailed           = &Val{merror, "test failed"}
)

func isTaggedBy(v, s *Val) *Val {
	if isPair(v) != vfalse && isSymbol(car(v)) != vfalse && smeq(car(v), s) != vfalse {
		return vtrue
	}

	return vfalse
}

func isQuote(v *Val) *Val {
	return isTaggedBy(v, sfromString("quote"))
}

func evalQuote(e, v *Val) *Val {
	if isPair(cdr(v)) == vfalse || cdr(cdr(v)) != Nil {
		fatal(invalidQuote)
	}

	return car(cdr(v))
}

func makeFn(a, b *Val) *Val {
	return cons(sfromString("fn"), cons(a, b))
}

func isDef(v *Val) *Val {
	return isTaggedBy(v, sfromString("def"))
}

func isVectorForm(v *Val) *Val {
	return isTaggedBy(v, sfromString("vector:"))
}

func evalVector(e, v *Val) *Val {
	if isPair(v) == vfalse {
		return fatal(invalidVector)
	}

	return vectorFromList(valueList(e, cdr(v)))
}

func isStructForm(v *Val) *Val {
	return isTaggedBy(v, sfromString("struct:"))
}

func evalStructValues(e, v *Val) *Val {
	if isNil(v) != vfalse {
		return Nil
	}

	if isPair(v) == vfalse || isPair(cdr(v)) == vfalse {
		return fatal(invalidStruct)
	}

	return cons(
		car(v),
		cons(
			evalExp(e, car(cdr(v))),
			evalStructValues(e, cdr(cdr(v))),
		),
	)
}

func evalStruct(e, v *Val) *Val {
	if isPair(v) == vfalse {
		return fatal(invalidStruct)
	}

	return structFromList(evalStructValues(e, cdr(v)))
}

func nameOfDef(v *Val) *Val {
	return car(cdr(v))
}

func nameOfFunctionDef(v *Val) *Val {
	if isPair(car(cdr(v))) == vfalse || isSymbol(car(car(cdr(v)))) == vfalse {
		return fatal(invalidDef)
	}

	return car(car(cdr(v)))
}

func defName(v *Val) *Val {
	if isPair(cdr(v)) == vfalse {
		return fatal(invalidDef)
	}

	if isSymbol(car(cdr(v))) != vfalse {
		return nameOfDef(v)
	}

	return nameOfFunctionDef(v)
}

func valueOfDef(v *Val) *Val {
	if isPair(cdr(cdr(v))) == vfalse || cdr(cdr(cdr(v))) != Nil {
		return fatal(invalidDef)
	}

	return car(cdr(cdr(v)))
}

func defValue(v *Val) *Val {
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

func evalDef(e, v *Val) *Val {
	return define(e, defName(v), evalExp(e, defValue(v)))
}

func isIf(v *Val) *Val {
	return isTaggedBy(v, sfromString("if"))
}

func ifPredicate(v *Val) *Val {
	if isPair(cdr(v)) == vfalse {
		return fatal(invalidIf)
	}

	return car(cdr(v))
}

func ifConsequent(v *Val) *Val {
	if isPair(cdr(v)) == vfalse || isPair(cdr(cdr(v))) == vfalse {
		return fatal(invalidIf)
	}

	return car(cdr(cdr(v)))
}

func ifAlternative(v *Val) *Val {
	if isPair(cdr(v)) == vfalse ||
		isPair(cdr(cdr(v))) == vfalse ||
		isPair(cdr(cdr(cdr(v)))) == vfalse {
		return fatal(invalidIf)
	}

	return car(cdr(cdr(cdr(v))))
}

func evalIf(e, v *Val) *Val {
	if evalExp(e, ifPredicate(v)) != vfalse {
		return evalExp(e, ifConsequent(v))
	}

	return evalExp(e, ifAlternative(v))
}

func isAnd(v *Val) *Val {
	return isTaggedBy(v, sfromString("and"))
}

func evalAnd(e, v *Val) *Val {
	if isNil(v) != vfalse {
		return vtrue
	}

	if isPair(v) == vfalse {
		return invalidAnd
	}

	r := evalExp(e, car(v))
	if isNil(cdr(v)) != vfalse {
		return r
	}

	if r == vfalse {
		return vfalse
	}

	return evalAnd(e, cdr(v))
}

func isOr(v *Val) *Val {
	return isTaggedBy(v, sfromString("or"))
}

func evalOr(e, v *Val) *Val {
	if isNil(v) != vfalse {
		return vfalse
	}

	if isPair(v) == vfalse {
		return invalidOr
	}

	r := evalExp(e, car(v))
	if isNil(cdr(v)) != vfalse {
		return r
	}

	if r != vfalse {
		return r
	}

	return evalAnd(e, cdr(v))
}

func isFn(v *Val) *Val {
	return isTaggedBy(v, sfromString("fn"))
}

func fnParams(v *Val) *Val {
	if isPair(v) == vfalse || isPair(cdr(v)) == vfalse ||
		isSymbol(car(cdr(v))) == vfalse && isPair(car(cdr(v))) == vfalse && isNil(car(cdr(v))) == vfalse {
		return fatal(invalidFn)
	}

	return car(cdr(v))
}

func fnBody(v *Val) *Val {
	if isPair(v) == vfalse || isPair(cdr(v)) == vfalse || isPair(cdr(cdr(v))) == vfalse {
		return fatal(invalidFn)
	}

	return cdr(cdr(v))
}

func fnToProc(e, v *Val) *Val {
	return newProc(e, fnParams(v), fnBody(v))
}

func isBegin(v *Val) *Val {
	return isTaggedBy(v, sfromString("begin"))
}

func beginSeq(v *Val) *Val {
	if isPair(v) == vfalse {
		return fatal(invalidSequence)
	}

	return cdr(v)
}

func evalSeq(e, v *Val) *Val {
	if isPair(v) == vfalse {
		return fatal(invalidSequence)
	}

	if cdr(v) == Nil {
		return evalExp(e, car(v))
	}

	eval(e, car(v))
	return evalSeq(e, cdr(v))
}

func isCond(v *Val) *Val {
	return isTaggedBy(v, sfromString("cond"))
}

func seqToExp(v *Val) *Val {
	if isPair(v) == vfalse {
		return fatal(invalidCond)
	}

	if isNil(cdr(v)) != vfalse {
		return car(v)
	}

	return cons(sfromString("begin"), v)
}

func expandCond(v *Val) *Val {
	if isPair(v) == vfalse {
		return fatal(invalidCond)
	}

	cond := car(v)
	rest := cdr(v)

	if isPair(cond) == vfalse {
		return fatal(invalidCond)
	}

	pred := car(cond)

	if isSymbol(pred) != vfalse && smeq(pred, sfromString("else")) != vfalse {
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

func condToIf(v *Val) *Val {
	if isPair(v) == vfalse {
		return fatal(invalidCond)
	}

	return expandCond(cdr(v))
}

func isLet(v *Val) *Val {
	return isTaggedBy(v, sfromString("let"))
}

func letDefs(v *Val) *Val {
	if isNil(v) != vfalse {
		return Nil
	}

	if isPair(v) == vfalse || isPair(cdr(v)) == vfalse {
		return fatal(invalidLet)
	}

	return cons(
		list(sfromString("def"), car(v), car(cdr(v))),
		letDefs(cdr(cdr(v))),
	)
}

func letFnBody(v *Val) *Val {
	if isPair(v) == vfalse || isPair(cdr(v)) == vfalse || isNil(cdr(cdr(v))) != vfalse {
		return fatal(invalidLet)
	}

	return mappend(letDefs(car(cdr(v))), cdr(cdr(v)))
}

func expandLet(v *Val) *Val {
	return list(makeFn(Nil, letFnBody(v)))
}

func isTest(v *Val) *Val {
	return isTaggedBy(v, sfromString("test"))
}

// TODO: should be and
func evalTest(e, v *Val) *Val {
	if isPair(v) == vfalse {
		return fatal(invalidTest)
	}

	if isNil(cdr(v)) != vfalse {
		return vtrue
	}

	result := evalSeq(newEnv(e), cdr(v))
	if result == vfalse {
		return fatal(testFailed)
	}

	if isError(result) != vfalse {
		return fatal(result)
	}

	return sfromString("test-complete")
}

func isApplication(v *Val) *Val {
	return isPair(v)
}

func valueList(e, v *Val) *Val {
	if isNil(v) != vfalse {
		return Nil
	}

	if isPair(v) == vfalse {
		return fatal(invalidApplication)
	}

	return cons(evalExp(e, car(v)), valueList(e, cdr(v)))
}

func evalApply(e, v *Val) *Val {
	if isPair(v) == vfalse {
		return fatal(invalidApplication)
	}

	return apply(evalExp(e, car(v)), valueList(e, cdr(v)))
}

func evalExp(e, v *Val) *Val {
	switch {
	case isDef(v) != vfalse:
		return fatal(definitionExpression)
	default:
		return eval(e, v)
	}
}

func eval(e, v *Val) *Val {
	switch {
	case isNumber(v) != vfalse:
		return v
	case isString(v) != vfalse:
		return v
	case isBool(v) != vfalse:
		return v
	case isNil(v) != vfalse:
		return v
	case isQuote(v) != vfalse:
		return evalQuote(e, v)
	case isSymbol(v) != vfalse:
		return lookupDef(e, v)
	case isDef(v) != vfalse:
		return evalDef(e, v)
	case isVectorForm(v) != vfalse:
		return evalVector(e, v)
	case isStructForm(v) != vfalse:
		return evalStruct(e, v)
	case isIf(v) != vfalse:
		return evalIf(e, v)
	case isAnd(v) != vfalse:
		return evalAnd(e, cdr(v))
	case isOr(v) != vfalse:
		return evalOr(e, cdr(v))
	case isFn(v) != vfalse:
		return fnToProc(e, v)
	case isBegin(v) != vfalse:
		return evalSeq(e, beginSeq(v))
	case isCond(v) != vfalse:
		return eval(e, condToIf(v))
	case isLet(v) != vfalse:
		return eval(e, expandLet(v))
	case isTest(v) != vfalse:
		return evalTest(e, v)
	case isApplication(v) != vfalse:
		return evalApply(e, v)
	default:
		println(v.mtype)
		return fatal(invalidExpression)
	}
}

func Eval(e, v *Val) *Val {
	return eval(e, v)
}
