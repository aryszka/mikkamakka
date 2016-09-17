package mikkamakka

import "unicode"

var (
	invalidToken       = SysStringToError("invalid token")
	notImplemented     = SysStringToError("not implemented")
	unexpectedClose    = SysStringToError("unexpected close")
	irregularCons      = SysStringToError("irregular cons")
	UndefinedReadValue = SysStringToError("undefined read value")
	ttnone             = SysIntToNumber(0)
	ttcomment          = SysIntToNumber(1)
	ttsymbol           = SysIntToNumber(2)
	ttstring           = SysIntToNumber(3)
	ttlist             = SysIntToNumber(4)
	ttquote            = SysIntToNumber(5)
	ttvector           = SysIntToNumber(7)
	ttstruct           = SysIntToNumber(8)
)

func reader(in *Val) *Val {
	return SysMapToStruct(map[string]*Val{
		"in":            in,
		"token-type":    ttnone,
		"value":         UndefinedReadValue,
		"escaped":       False,
		"last-char":     SysStringToString(""),
		"current-token": SysStringToString(""),
		"in-list":       SysStringToString(""),
		"close-list":    False,
		"cons":          False,
		"cons-items":    SysIntToNumber(0),
	})
}

func charCheck(c string) func(*Val) *Val {
	return func(s *Val) *Val {
		if StringToSysString(s) == c {
			return True
		}

		return False
	}
}

var (
	isEscapeChar      = charCheck("\\")
	isNewline         = charCheck("\n")
	isStringDelimiter = charCheck(`"`)
	isComment         = charCheck(";")
	isListOpen        = charCheck("(")
	isListClose       = charCheck(")")
	isCons            = charCheck(".")
	isQuoteChar       = charCheck("'")
	isOpenVector      = charCheck("[")
	isCloseVector     = charCheck("]")
	isOpenStruct      = charCheck("{")
	isCloseStruct     = charCheck("}")
)

func isWhitespace(s *Val) *Val {
	if unicode.IsSpace(rune(StringToSysString(s)[0])) {
		return True
	}

	return False
}

func symbolToken(t *Val) *Val {
	v := StringToNumber(t)
	if IsError(v) == False {
		return v
	}

	v = SysStringToBool(StringToSysString(t))
	if IsError(v) == False {
		return v
	}

	return SymbolFromRawString(StringToSysString(t))
}

func readChar(r *Val) *Val {
	in := Fread(Field(r, SymbolFromRawString("in")), SysIntToNumber(1))
	st := Fstate(in)

	if IsError(st) != False {
		return Assign(r, SysMapToStruct(map[string]*Val{
			"in":    in,
			"value": st,
		}))
	}

	return Assign(r, SysMapToStruct(map[string]*Val{
		"in":        in,
		"last-char": st,
	}))
}

func readError(r *Val) bool {
	v := Field(r, SymbolFromRawString("value"))
	return IsError(v) != False && v != UndefinedReadValue

}

func lastChar(r *Val) *Val {
	return Field(r, SymbolFromRawString("last-char"))
}

func currentTokenType(r *Val) *Val {
	return Field(r, SymbolFromRawString("token-type"))
}

func setTokenType(r *Val, t *Val) *Val {
	return Assign(r, SysMapToStruct(map[string]*Val{
		"token-type": t,
	}))
}

func isTNone(t *Val) bool    { return t == ttnone }
func isTComment(t *Val) bool { return t == ttcomment }
func isTSymbol(t *Val) bool  { return t == ttsymbol }
func isTString(t *Val) bool  { return t == ttstring }
func isTList(t *Val) bool    { return t == ttlist }
func isTQuote(t *Val) bool   { return t == ttquote }
func isTVector(t *Val) bool  { return t == ttvector }
func isTStruct(t *Val) bool  { return t == ttstruct }

func setNone(r *Val) *Val    { return setTokenType(r, ttnone) }
func setString(r *Val) *Val  { return setTokenType(r, ttstring) }
func setComment(r *Val) *Val { return setTokenType(r, ttcomment) }
func setSymbol(r *Val) *Val  { return setTokenType(r, ttsymbol) }
func setList(r *Val) *Val    { return setTokenType(r, ttlist) }
func setQuote(r *Val) *Val   { return setTokenType(r, ttquote) }
func setVector(r *Val) *Val  { return setTokenType(r, ttvector) }
func setStruct(r *Val) *Val  { return setTokenType(r, ttstruct) }

func clearToken(r *Val) *Val {
	return Assign(r, SysMapToStruct(map[string]*Val{
		"current-token": SysStringToString(""),
	}))
}

func closeComment(r *Val) *Val {
	return clearToken(setTokenType(r, ttnone))
}

func closeSymbol(r *Val) *Val {
	return clearToken(processSymbol(setTokenType(r, ttnone)))
}

func closeString(r *Val) *Val {
	return clearToken(processString(setTokenType(r, ttnone)))
}

func setEscaped(r *Val) *Val {
	return Assign(r, SysMapToStruct(map[string]*Val{"escaped": True}))
}

func unsetEscaped(r *Val) *Val {
	return Assign(r, SysMapToStruct(map[string]*Val{"escaped": False}))
}

func isEscaped(r *Val) *Val {
	return Field(r, SymbolFromRawString("escaped"))
}

func unescapeSymbolChar(c *Val) *Val {
	return c
}

func unescapeStringChar(c *Val) *Val {
	switch StringToSysString(c) {
	case "b":
		return SysStringToString("\b")
	case "f":
		return SysStringToString("\f")
	case "n":
		return SysStringToString("\n")
	case "r":
		return SysStringToString("\r")
	case "t":
		return SysStringToString("\t")
	case "v":
		return SysStringToString("\v")
	default:
		return c
	}
}

func unescapeChar(tokenType, c *Val) *Val {
	switch tokenType {
	case ttsymbol:
		return unescapeSymbolChar(c)
	case ttstring:
		return unescapeStringChar(c)
	default:
		return invalidToken
	}
}

func appendToken(r *Val) *Val {
	c := lastChar(r)
	if isEscaped(r) != False {
		c = unescapeChar(Field(r, SymbolFromRawString("token-type")), c)
	}

	return Assign(r, SysMapToStruct(map[string]*Val{
		"current-token": AppendString(Field(r, SymbolFromRawString("current-token")), c),
	}))
}

func setInvalid(r *Val) *Val {
	return Assign(r, SysMapToStruct(map[string]*Val{
		"err": invalidToken,
	}))
}

func processSymbol(r *Val) *Val {
	return Assign(r, SysMapToStruct(map[string]*Val{
		"value": symbolToken(Field(r, SymbolFromRawString("current-token"))),
	}))
}

func processString(r *Val) *Val {
	return Assign(r, SysMapToStruct(map[string]*Val{
		"value": Field(r, SymbolFromRawString("current-token")),
	}))
}

func closeChar(c *Val) *Val {
	switch StringToSysString(c) {
	case "(":
		return SysStringToString(")")
	case "[":
		return SysStringToString("]")
	case "{":
		return SysStringToString("}")
	default:
		return SysStringToString("")
	}
}

func setClose(r, c *Val) *Val {
	if stringEq(closeChar(Field(r, SymbolFromRawString("in-list"))), c) == False {
		return setUnexpectedClose(r)
	}

	return Assign(r, SysMapToStruct(map[string]*Val{
		"close-list": True,
	}))
}

func hasCons(r *Val) bool {
	return Greater(Field(r, SymbolFromRawString("cons-items")), SysIntToNumber(0)) != False
}

func consSet(r *Val) bool {
	return Field(r, SymbolFromRawString("cons")) != False
}

func setCons(r *Val) *Val {
	if hasCons(r) {
		return setIrregularCons(r)
	}

	return Assign(r, SysMapToStruct(map[string]*Val{
		"cons": True,
	}))
}

func setUnexpectedClose(r *Val) *Val {
	return Assign(r, SysMapToStruct(map[string]*Val{
		"value": unexpectedClose,
	}))
}

func setIrregularCons(r *Val) *Val {
	return Assign(r, SysMapToStruct(map[string]*Val{
		"value": irregularCons,
	}))
}

func reverse(l *Val) *Val {
	checkType(l, Pair, Nil)

	r := NilVal
	for {
		if l == NilVal {
			return r
		}

		r = Cons(Car(l), r)
		l = Cdr(l)
	}
}

func reverseIrregular(l *Val) *Val {
	checkType(l, Pair, Nil)

	r := Cons(Car(Cdr(l)), Car(l))
	l = Cdr(Cdr(l))
	for {
		if l == NilVal {
			return r
		}

		r = Cons(Car(l), r)
		l = Cdr(l)
	}
}

func readList(r, c *Val) *Val {
	lr := reader(Field(r, SymbolFromRawString("in")))
	lr = Assign(lr, SysMapToStruct(map[string]*Val{
		"list-items": NilVal,
		"in-list":    c,
	}))

	var loop func(*Val) *Val
	loop = func(lr *Val) *Val {
		lr = read(lr)
		if readError(lr) {
			return Assign(r, SysMapToStruct(map[string]*Val{
				"in":    Field(lr, SymbolFromRawString("in")),
				"value": Field(lr, SymbolFromRawString("value")),
			}))
		}

		v := Field(lr, SymbolFromRawString("value"))
		if v != UndefinedReadValue {
			lr = Assign(lr, SysMapToStruct(map[string]*Val{
				"list-items": Cons(
					v,
					Field(lr, SymbolFromRawString("list-items")),
				),
				"value": UndefinedReadValue,
			}))

			if hasCons(lr) {
				lr = Assign(lr, SysMapToStruct(map[string]*Val{
					"cons-items": Add(Field(lr, SymbolFromRawString("cons-items")), SysIntToNumber(1)),
				}))
			}
		}

		if consSet(lr) {
			if Field(lr, SymbolFromRawString("list-items")) == NilVal ||
				numberEq(Field(lr, SymbolFromRawString("cons-items")), SysIntToNumber(0)) == False {
				return setIrregularCons(Assign(r, SysMapToStruct(map[string]*Val{
					"in": Field(lr, SymbolFromRawString("in")),
				})))
			}

			lr = Assign(lr, SysMapToStruct(map[string]*Val{
				"cons-items": SysIntToNumber(1),
				"cons":       False,
			}))
		}

		if Field(lr, SymbolFromRawString("close-list")) != False {
			if hasCons(lr) {
				if numberEq(Field(lr, SymbolFromRawString("cons-items")), SysIntToNumber(2)) == False {
					return setIrregularCons(Assign(r, SysMapToStruct(map[string]*Val{
						"in": Field(lr, SymbolFromRawString("in")),
					})))
				}

				return Assign(r, SysMapToStruct(map[string]*Val{
					"in":    Field(lr, SymbolFromRawString("in")),
					"value": reverseIrregular(Field(lr, SymbolFromRawString("list-items"))),
				}))
			}

			return Assign(r, SysMapToStruct(map[string]*Val{
				"in":    Field(lr, SymbolFromRawString("in")),
				"value": reverse(Field(lr, SymbolFromRawString("list-items"))),
			}))
		}

		return loop(lr)
	}

	return loop(lr)
}

func readQuote(r *Val) *Val {
	lr := reader(Field(r, SymbolFromRawString("in")))
	if stringEq(closeChar(Field(r, SymbolFromRawString("in-list"))), SysStringToString("")) == False {
		lr = Assign(lr, SysMapToStruct(map[string]*Val{
			"in-list": Field(r, SymbolFromRawString("in-list")),
		}))
	}

	lr = read(lr)
	if readError(lr) {
		return Assign(r, SysMapToStruct(map[string]*Val{
			"in":    Field(lr, SymbolFromRawString("in")),
			"value": Field(lr, SymbolFromRawString("value")),
		}))
	}

	return Assign(r, SysMapToStruct(map[string]*Val{
		"in":         Field(lr, SymbolFromRawString("in")),
		"value":      List(SymbolFromRawString("quote"), Field(lr, SymbolFromRawString("value"))),
		"close-list": Field(lr, SymbolFromRawString("close-list")),
	}))
}

func readVector(r *Val) *Val {
	r = readList(r, SysStringToString("["))
	if readError(r) {
		return r
	}

	return Assign(r, SysMapToStruct(map[string]*Val{
		"value": Cons(SymbolFromRawString("vector:"), Field(r, SymbolFromRawString("value"))),
	}))
}

func readStruct(r *Val) *Val {
	r = readList(r, SysStringToString("{"))
	if readError(r) {
		return r
	}

	return Assign(r, SysMapToStruct(map[string]*Val{
		"value": Cons(SymbolFromRawString("struct:"), Field(r, SymbolFromRawString("value"))),
	}))
}

func read(r *Val) *Val {
	t := currentTokenType(r)
	if isTList(t) {
		return setNone(readList(r, SysStringToString("(")))
	}

	if isTQuote(t) {
		return setNone(readQuote(r))
	}

	if isTVector(t) {
		return setNone(readVector(r))
	}

	if isTStruct(t) {
		return setNone(readStruct(r))
	}

	r = readChar(r)
	if readError(r) {
		return r
	}

	c := lastChar(r)

	switch {
	case isTNone(t):
		switch {
		case isWhitespace(c) != False:
			return read(r)
		case isEscapeChar(c) != False:
			return read(setEscaped(setSymbol(r)))
		case isStringDelimiter(c) != False:
			return read(setString(r))
		case isComment(c) != False:
			return read(setComment(r))
		case isListOpen(c) != False:
			return read(setList(r))
		case isListClose(c) != False:
			return setClose(r, c)
		case isCons(c) != False:
			return setCons(r)
		case isQuoteChar(c) != False:
			return read(setQuote(r))
		case isOpenVector(c) != False:
			return read(setVector(r))
		case isCloseVector(c) != False:
			return setClose(r, c)
		case isOpenStruct(c) != False:
			return read(setStruct(r))
		case isCloseStruct(c) != False:
			return setClose(r, c)
		default:
			return read(appendToken(setSymbol(r)))
		}

	case isTComment(t):
		switch {
		case isNewline(c) != False:
			return read(closeComment(r))
		}

		return read(r)

	case isTSymbol(t):
		switch {
		case isEscaped(r) != False:
			return read(unsetEscaped(appendToken(r)))
		case isWhitespace(c) != False:
			return closeSymbol(r)
		case isEscapeChar(c) != False:
			return read(setEscaped(r))
		case isComment(c) != False:
			return setComment(closeSymbol(r))
		case isStringDelimiter(c) != False:
			return setString(closeSymbol(r))
		case isListOpen(c) != False:
			return setList(closeSymbol(r))
		case isListClose(c) != False:
			return setClose(closeSymbol(r), c)
		case isCons(c) != False:
			return setCons(closeSymbol(r))
		case isOpenVector(c) != False:
			return setVector(closeSymbol(r))
		case isCloseVector(c) != False:
			return setClose(closeSymbol(r), c)
		case isOpenStruct(c) != False:
			return setStruct(closeSymbol(r))
		case isCloseStruct(c) != False:
			return setClose(closeSymbol(r), c)
		default:
			return read(appendToken(r))
		}

	case isTString(t):
		switch {
		case isEscaped(r) != False:
			return read(unsetEscaped(appendToken(r)))
		case isEscapeChar(c) != False:
			return read(setEscaped(r))
		case isStringDelimiter(c) != False:
			return closeString(r)
		default:
			return read(appendToken(r))
		}

	default:
		return setInvalid(r)
	}
}

func Reader(in *Val) *Val { return reader(in) }
func Read(r *Val) *Val    { return read(r) }
