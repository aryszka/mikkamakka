package mikkamakka

/*
when panic, when error?
*/

var (
	invalidQuote         = ErrorFromRawString("invalid quote")
	invalidDef           = ErrorFromRawString("invalid definition")
	invalidIf            = ErrorFromRawString("invalid if expression")
	invalidAnd           = ErrorFromRawString("invalid and expression")
	invalidOr            = ErrorFromRawString("invalid or expression")
	invalidFn            = ErrorFromRawString("invalid function expression")
	invalidSequence      = ErrorFromRawString("invalid sequence")
	invalidVector        = ErrorFromRawString("invalid vector")
	invalidCond          = ErrorFromRawString("invalid cond expression")
	invalidLet           = ErrorFromRawString("invalid let expression")
	invalidTest          = ErrorFromRawString("invalid test")
	invalidApplication   = ErrorFromRawString("invalid application")
	invalidExpression    = ErrorFromRawString("invalid expression")
	definitionExpression = ErrorFromRawString("definition in expression position")
	notFunction          = ErrorFromRawString("not a function")
	testFailed           = ErrorFromRawString("test failed")
)

func List(a ...*Val) *Val {
	l := Nil
	for i := len(a) - 1; i >= 0; i-- {
		l = Cons(a[i], l)
	}

	return l
}

func mappend(left, right *Val) *Val {
	checkType(left, pair, mnil)
	checkType(right, pair, mnil)

	if IsNil(left) != False {
		return right
	}

	return Cons(Car(left), mappend(Cdr(left), right))
}

func isTaggedBy(v, s *Val) *Val {
	if IsPair(v) != False && IsSymbol(Car(v)) != False && Eq(Car(v), s) != False {
		return True
	}

	return False
}

func isQuote(v *Val) *Val {
	return isTaggedBy(v, SymbolFromRawString("quote"))
}

func evalQuote(e, v *Val) *Val {
	if IsPair(Cdr(v)) == False || Cdr(Cdr(v)) != Nil {
		Fatal(invalidQuote)
	}

	return Car(Cdr(v))
}

func makeFn(a, b *Val) *Val {
	return Cons(SymbolFromRawString("fn"), Cons(a, b))
}

func isDef(v *Val) *Val {
	return isTaggedBy(v, SymbolFromRawString("def"))
}

func isVectorForm(v *Val) *Val {
	return isTaggedBy(v, SymbolFromRawString("vector:"))
}

func evalVector(e, v *Val) *Val {
	if IsPair(v) == False {
		return Fatal(invalidVector)
	}

	return VectorFromList(valueList(e, Cdr(v)))
}

func isStructForm(v *Val) *Val {
	return isTaggedBy(v, SymbolFromRawString("struct:"))
}

func evalStructValues(e, v *Val) *Val {
	if IsNil(v) != False {
		return Nil
	}

	if IsPair(v) == False || IsPair(Cdr(v)) == False {
		return Fatal(InvalidStruct)
	}

	return Cons(
		Car(v),
		Cons(
			evalExp(e, Car(Cdr(v))),
			evalStructValues(e, Cdr(Cdr(v))),
		),
	)
}

func evalStruct(e, v *Val) *Val {
	if IsPair(v) == False {
		return Fatal(InvalidStruct)
	}

	return StructFromList(evalStructValues(e, Cdr(v)))
}

func nameOfDef(v *Val) *Val {
	return Car(Cdr(v))
}

func nameOfFunctionDef(v *Val) *Val {
	if IsPair(Car(Cdr(v))) == False || IsSymbol(Car(Car(Cdr(v)))) == False {
		return Fatal(invalidDef)
	}

	return Car(Car(Cdr(v)))
}

func defName(v *Val) *Val {
	if IsPair(Cdr(v)) == False {
		return Fatal(invalidDef)
	}

	if IsSymbol(Car(Cdr(v))) != False {
		return nameOfDef(v)
	}

	return nameOfFunctionDef(v)
}

func valueOfDef(v *Val) *Val {
	if IsPair(Cdr(Cdr(v))) == False || Cdr(Cdr(Cdr(v))) != Nil {
		return Fatal(invalidDef)
	}

	return Car(Cdr(Cdr(v)))
}

func defValue(v *Val) *Val {
	if IsPair(Cdr(v)) == False {
		return Fatal(invalidDef)
	}

	if IsSymbol(Car(Cdr(v))) != False {
		return valueOfDef(v)
	}

	if IsPair(Car(Cdr(v))) == False {
		return Fatal(invalidDef)
	}

	return makeFn(Cdr(Car(Cdr(v))), Cdr(Cdr(v)))
}

func evalDef(e, v *Val) *Val {
	return define(e, defName(v), evalExp(e, defValue(v)))
}

func isIf(v *Val) *Val {
	return isTaggedBy(v, SymbolFromRawString("if"))
}

func ifPredicate(v *Val) *Val {
	if IsPair(Cdr(v)) == False {
		return Fatal(invalidIf)
	}

	return Car(Cdr(v))
}

func ifConsequent(v *Val) *Val {
	if IsPair(Cdr(v)) == False || IsPair(Cdr(Cdr(v))) == False {
		return Fatal(invalidIf)
	}

	return Car(Cdr(Cdr(v)))
}

func ifAlternative(v *Val) *Val {
	if IsPair(Cdr(v)) == False ||
		IsPair(Cdr(Cdr(v))) == False ||
		IsPair(Cdr(Cdr(Cdr(v)))) == False {
		return Fatal(invalidIf)
	}

	return Car(Cdr(Cdr(Cdr(v))))
}

func evalIf(e, v *Val) *Val {
	if evalExp(e, ifPredicate(v)) != False {
		return evalExp(e, ifConsequent(v))
	}

	return evalExp(e, ifAlternative(v))
}

func isAnd(v *Val) *Val {
	return isTaggedBy(v, SymbolFromRawString("and"))
}

func evalAnd(e, v *Val) *Val {
	if IsNil(v) != False {
		return True
	}

	if IsPair(v) == False {
		return invalidAnd
	}

	r := evalExp(e, Car(v))
	if IsNil(Cdr(v)) != False {
		return r
	}

	if r == False {
		return r
	}

	return evalAnd(e, Cdr(v))
}

func isOr(v *Val) *Val {
	return isTaggedBy(v, SymbolFromRawString("or"))
}

func evalOr(e, v *Val) *Val {
	if IsNil(v) != False {
		return False
	}

	if IsPair(v) == False {
		return invalidOr
	}

	r := evalExp(e, Car(v))
	if IsNil(Cdr(v)) != False {
		return r
	}

	if r != False {
		return r
	}

	return evalAnd(e, Cdr(v))
}

func isFunctionLiteral(v *Val) *Val {
	return isTaggedBy(v, SymbolFromRawString("fn"))
}

func fnParams(v *Val) *Val {
	if IsPair(v) == False || IsPair(Cdr(v)) == False ||
		IsSymbol(Car(Cdr(v))) == False && IsPair(Car(Cdr(v))) == False && IsNil(Car(Cdr(v))) == False {
		return Fatal(invalidFn)
	}

	return Car(Cdr(v))
}

func fnBody(v *Val) *Val {
	if IsPair(v) == False || IsPair(Cdr(v)) == False || IsPair(Cdr(Cdr(v))) == False {
		return Fatal(invalidFn)
	}

	return Cdr(Cdr(v))
}

func fnToFunc(e, v *Val) *Val {
	return NewComposite(Cons(e, Cons(fnParams(v), fnBody(v))))
}

func isBegin(v *Val) *Val {
	return isTaggedBy(v, SymbolFromRawString("begin"))
}

func beginSeq(v *Val) *Val {
	if IsPair(v) == False {
		return Fatal(invalidSequence)
	}

	return Cdr(v)
}

func evalSeq(e, v *Val) *Val {
	if IsPair(v) == False {
		return Fatal(invalidSequence)
	}

	if Cdr(v) == Nil {
		return evalExp(e, Car(v))
	}

	eval(e, Car(v))
	return evalSeq(e, Cdr(v))
}

func isCond(v *Val) *Val {
	return isTaggedBy(v, SymbolFromRawString("cond"))
}

func seqToExp(v *Val) *Val {
	if IsPair(v) == False {
		return Fatal(invalidCond)
	}

	if IsNil(Cdr(v)) != False {
		return Car(v)
	}

	return Cons(SymbolFromRawString("begin"), v)
}

func expandCond(v *Val) *Val {
	if IsPair(v) == False {
		return Fatal(invalidCond)
	}

	cond := Car(v)
	rest := Cdr(v)

	if IsPair(cond) == False {
		return Fatal(invalidCond)
	}

	pred := Car(cond)

	if IsSymbol(pred) != False && Eq(pred, SymbolFromRawString("else")) != False {
		if IsNil(rest) == False {
			return Fatal(invalidCond)
		}

		return seqToExp(Cdr(cond))
	}

	return List(
		SymbolFromRawString("if"),
		pred,
		seqToExp(Cdr(cond)),
		expandCond(rest),
	)
}

func condToIf(v *Val) *Val {
	if IsPair(v) == False {
		return Fatal(invalidCond)
	}

	return expandCond(Cdr(v))
}

func isLet(v *Val) *Val {
	return isTaggedBy(v, SymbolFromRawString("let"))
}

func letDefs(v *Val) *Val {
	if IsNil(v) != False {
		return Nil
	}

	if IsPair(v) == False || IsPair(Cdr(v)) == False {
		return Fatal(invalidLet)
	}

	return Cons(
		List(SymbolFromRawString("def"), Car(v), Car(Cdr(v))),
		letDefs(Cdr(Cdr(v))),
	)
}

func letFnBody(v *Val) *Val {
	if IsPair(v) == False || IsPair(Cdr(v)) == False || IsNil(Cdr(Cdr(v))) != False {
		return Fatal(invalidLet)
	}

	return mappend(letDefs(Car(Cdr(v))), Cdr(Cdr(v)))
}

func expandLet(v *Val) *Val {
	return List(makeFn(Nil, letFnBody(v)))
}

func isTest(v *Val) *Val {
	return isTaggedBy(v, SymbolFromRawString("test"))
}

// TODO: should be and
func evalTest(e, v *Val) *Val {
	if IsPair(v) == False {
		return Fatal(invalidTest)
	}

	if IsNil(Cdr(v)) != False {
		return True
	}

	result := evalSeq(newEnv(e), Cdr(v))
	if result == False {
		return Fatal(testFailed)
	}

	if IsError(result) != False {
		return Fatal(result)
	}

	return SymbolFromRawString("test-complete")
}

func isApplication(v *Val) *Val {
	return IsPair(v)
}

func valueList(e, v *Val) *Val {
	if IsNil(v) != False {
		return Nil
	}

	if IsPair(v) == False {
		return Fatal(invalidApplication)
	}

	return Cons(evalExp(e, Car(v)), valueList(e, Cdr(v)))
}

func evalApply(e, v *Val) *Val {
	if IsPair(v) == False {
		return Fatal(invalidApplication)
	}

	return Apply(evalExp(e, Car(v)), valueList(e, Cdr(v)))
}

func Apply(f, a *Val) *Val {
	if IsStruct(f) != False {
		return Field(f, Car(a))
	}

	if IsCompiledFunction(f) != False {
		return ApplyCompiled(f, a)
	}

	if IsFunction(f) == False {
		return notFunction
	}

	cf := Composite(f)
	return evalSeq(extendEnv(Car(cf), Car(Cdr(cf)), a), Cdr(Cdr(cf)))
}

func evalExp(e, v *Val) *Val {
	switch {
	case isDef(v) != False:
		return Fatal(definitionExpression)
	default:
		return eval(e, v)
	}
}

func eval(e, v *Val) *Val {
	switch {
	case IsNumber(v) != False:
		return v
	case IsString(v) != False:
		return v
	case IsBool(v) != False:
		return v
	case IsNil(v) != False:
		return v
	case isQuote(v) != False:
		return evalQuote(e, v)
	case IsSymbol(v) != False:
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
		return evalAnd(e, Cdr(v))
	case isOr(v) != False:
		return evalOr(e, Cdr(v))
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
		return Fatal(invalidExpression)
	}
}

func Eval(e, v *Val) *Val {
	return eval(e, v)
}
