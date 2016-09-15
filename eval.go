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
	notFunction          = &Val{merror, "not a function"}
	testFailed           = &Val{merror, "test failed"}
)

func isTaggedBy(v, s *Val) *Val {
	if isPair(v) != False && isSymbol(car(v)) != False && smeq(car(v), s) != False {
		return True
	}

	return False
}

func isQuote(v *Val) *Val {
	return isTaggedBy(v, sfromString("quote"))
}

func evalQuote(e, v *Val) *Val {
	if isPair(cdr(v)) == False || cdr(cdr(v)) != Nil {
		fatal(invalidQuote)
	}

	return car(cdr(v))
}

func makeFn(a, b *Val) *Val {
	return Cons(sfromString("fn"), Cons(a, b))
}

func isDef(v *Val) *Val {
	return isTaggedBy(v, sfromString("def"))
}

func isVectorForm(v *Val) *Val {
	return isTaggedBy(v, sfromString("vector:"))
}

func evalVector(e, v *Val) *Val {
	if isPair(v) == False {
		return fatal(invalidVector)
	}

	return vectorFromList(valueList(e, cdr(v)))
}

func isStructForm(v *Val) *Val {
	return isTaggedBy(v, sfromString("struct:"))
}

func evalStructValues(e, v *Val) *Val {
	if isNil(v) != False {
		return Nil
	}

	if isPair(v) == False || isPair(cdr(v)) == False {
		return fatal(invalidStruct)
	}

	return Cons(
		car(v),
		Cons(
			evalExp(e, car(cdr(v))),
			evalStructValues(e, cdr(cdr(v))),
		),
	)
}

func evalStruct(e, v *Val) *Val {
	if isPair(v) == False {
		return fatal(invalidStruct)
	}

	return structFromList(evalStructValues(e, cdr(v)))
}

func nameOfDef(v *Val) *Val {
	return car(cdr(v))
}

func nameOfFunctionDef(v *Val) *Val {
	if isPair(car(cdr(v))) == False || isSymbol(car(car(cdr(v)))) == False {
		return fatal(invalidDef)
	}

	return car(car(cdr(v)))
}

func defName(v *Val) *Val {
	if isPair(cdr(v)) == False {
		return fatal(invalidDef)
	}

	if isSymbol(car(cdr(v))) != False {
		return nameOfDef(v)
	}

	return nameOfFunctionDef(v)
}

func valueOfDef(v *Val) *Val {
	if isPair(cdr(cdr(v))) == False || cdr(cdr(cdr(v))) != Nil {
		return fatal(invalidDef)
	}

	return car(cdr(cdr(v)))
}

func defValue(v *Val) *Val {
	if isPair(cdr(v)) == False {
		return fatal(invalidDef)
	}

	if isSymbol(car(cdr(v))) != False {
		return valueOfDef(v)
	}

	if isPair(car(cdr(v))) == False {
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
	if isPair(cdr(v)) == False {
		return fatal(invalidIf)
	}

	return car(cdr(v))
}

func ifConsequent(v *Val) *Val {
	if isPair(cdr(v)) == False || isPair(cdr(cdr(v))) == False {
		return fatal(invalidIf)
	}

	return car(cdr(cdr(v)))
}

func ifAlternative(v *Val) *Val {
	if isPair(cdr(v)) == False ||
		isPair(cdr(cdr(v))) == False ||
		isPair(cdr(cdr(cdr(v)))) == False {
		return fatal(invalidIf)
	}

	return car(cdr(cdr(cdr(v))))
}

func evalIf(e, v *Val) *Val {
	if evalExp(e, ifPredicate(v)) != False {
		return evalExp(e, ifConsequent(v))
	}

	return evalExp(e, ifAlternative(v))
}

func isAnd(v *Val) *Val {
	return isTaggedBy(v, sfromString("and"))
}

func evalAnd(e, v *Val) *Val {
	if isNil(v) != False {
		return True
	}

	if isPair(v) == False {
		return invalidAnd
	}

	r := evalExp(e, car(v))
	if isNil(cdr(v)) != False {
		return r
	}

	if r == False {
		return False
	}

	return evalAnd(e, cdr(v))
}

func isOr(v *Val) *Val {
	return isTaggedBy(v, sfromString("or"))
}

func evalOr(e, v *Val) *Val {
	if isNil(v) != False {
		return False
	}

	if isPair(v) == False {
		return invalidOr
	}

	r := evalExp(e, car(v))
	if isNil(cdr(v)) != False {
		return r
	}

	if r != False {
		return r
	}

	return evalAnd(e, cdr(v))
}

func isFunctionLiteral(v *Val) *Val {
	return isTaggedBy(v, sfromString("fn"))
}

func fnParams(v *Val) *Val {
	if isPair(v) == False || isPair(cdr(v)) == False ||
		isSymbol(car(cdr(v))) == False && isPair(car(cdr(v))) == False && isNil(car(cdr(v))) == False {
		return fatal(invalidFn)
	}

	return car(cdr(v))
}

func fnBody(v *Val) *Val {
	if isPair(v) == False || isPair(cdr(v)) == False || isPair(cdr(cdr(v))) == False {
		return fatal(invalidFn)
	}

	return cdr(cdr(v))
}

func fnToFunc(e, v *Val) *Val {
	return NewComposite(Cons(e, Cons(fnParams(v), fnBody(v))))
}

func isBegin(v *Val) *Val {
	return isTaggedBy(v, sfromString("begin"))
}

func beginSeq(v *Val) *Val {
	if isPair(v) == False {
		return fatal(invalidSequence)
	}

	return cdr(v)
}

func evalSeq(e, v *Val) *Val {
	if isPair(v) == False {
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
	if isPair(v) == False {
		return fatal(invalidCond)
	}

	if isNil(cdr(v)) != False {
		return car(v)
	}

	return Cons(sfromString("begin"), v)
}

func expandCond(v *Val) *Val {
	if isPair(v) == False {
		return fatal(invalidCond)
	}

	cond := car(v)
	rest := cdr(v)

	if isPair(cond) == False {
		return fatal(invalidCond)
	}

	pred := car(cond)

	if isSymbol(pred) != False && smeq(pred, sfromString("else")) != False {
		if isNil(rest) == False {
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
	if isPair(v) == False {
		return fatal(invalidCond)
	}

	return expandCond(cdr(v))
}

func isLet(v *Val) *Val {
	return isTaggedBy(v, sfromString("let"))
}

func letDefs(v *Val) *Val {
	if isNil(v) != False {
		return Nil
	}

	if isPair(v) == False || isPair(cdr(v)) == False {
		return fatal(invalidLet)
	}

	return Cons(
		list(sfromString("def"), car(v), car(cdr(v))),
		letDefs(cdr(cdr(v))),
	)
}

func letFnBody(v *Val) *Val {
	if isPair(v) == False || isPair(cdr(v)) == False || isNil(cdr(cdr(v))) != False {
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
	if isPair(v) == False {
		return fatal(invalidTest)
	}

	if isNil(cdr(v)) != False {
		return True
	}

	result := evalSeq(newEnv(e), cdr(v))
	if result == False {
		return fatal(testFailed)
	}

	if isError(result) != False {
		return fatal(result)
	}

	return sfromString("test-complete")
}

func isApplication(v *Val) *Val {
	return isPair(v)
}

func valueList(e, v *Val) *Val {
	if isNil(v) != False {
		return Nil
	}

	if isPair(v) == False {
		return fatal(invalidApplication)
	}

	return Cons(evalExp(e, car(v)), valueList(e, cdr(v)))
}

func evalApply(e, v *Val) *Val {
	if isPair(v) == False {
		return fatal(invalidApplication)
	}

	return Apply(evalExp(e, car(v)), valueList(e, cdr(v)))
}

func Apply(f, a *Val) *Val {
	if isStruct(f) != False {
		return field(f, car(a))
	}

	if IsCompiledFunction(f) != False {
		return ApplyCompiled(f, a)
	}

	if IsFunction(f) == False {
		return notFunction
	}

	cf := Composite(f)
	return evalSeq(extendEnv(car(cf), car(cdr(cf)), a), cdr(cdr(cf)))
}

func evalExp(e, v *Val) *Val {
	switch {
	case isDef(v) != False:
		return fatal(definitionExpression)
	default:
		return eval(e, v)
	}
}

func eval(e, v *Val) *Val {
	switch {
	case isNumber(v) != False:
		return v
	case isString(v) != False:
		return v
	case isBool(v) != False:
		return v
	case isNil(v) != False:
		return v
	case isQuote(v) != False:
		return evalQuote(e, v)
	case isSymbol(v) != False:
		return lookupDef(e, v)
	case isDef(v) != False:
		return evalDef(e, v)
	case isVectorForm(v) != False:
		return evalVector(e, v)
	case isStructForm(v) != False:
		return evalStruct(e, v)
	case isIf(v) != False:
		return evalIf(e, v)
	case isAnd(v) != False:
		return evalAnd(e, cdr(v))
	case isOr(v) != False:
		return evalOr(e, cdr(v))
	case isFunctionLiteral(v) != False:
		return fnToFunc(e, v)
	case isBegin(v) != False:
		return evalSeq(e, beginSeq(v))
	case isCond(v) != False:
		return eval(e, condToIf(v))
	case isLet(v) != False:
		return eval(e, expandLet(v))
	case isTest(v) != False:
		return evalTest(e, v)
	case isApplication(v) != False:
		return evalApply(e, v)
	default:
		println(v.mtype)
		return fatal(invalidExpression)
	}
}

func Eval(e, v *Val) *Val {
	return eval(e, v)
}
