package mikkamakka

import "unicode"

var (
	invalidToken       = ErrorFromRawString("invalid token")
	notImplemented     = ErrorFromRawString("not implemented")
	unexpectedClose    = ErrorFromRawString("unexpected close")
	irregularCons      = ErrorFromRawString("irregular cons")
	UndefinedReadValue = &Val{symbol, "undefined read value"}
	ttnone             = &Val{number, 0}
	ttcomment          = &Val{number, 1}
	ttsymbol           = &Val{number, 2}
	ttstring           = &Val{number, 3}
	ttlist             = &Val{number, 4}
	ttquote            = &Val{number, 5}
	ttvector           = &Val{number, 7}
	ttstruct           = &Val{number, 8}
)

func reader(in *Val) *Val {
	return FromMap(Struct{
		"in":            in,
		"token-type":    ttnone,
		"value":         UndefinedReadValue,
		"escaped":       False,
		"last-char":     StringFromRaw(""),
		"current-token": StringFromRaw(""),
		"in-list":       StringFromRaw(""),
		"close-list":    False,
		"cons":          False,
		"cons-items":    NumberFromRawInt(0),
	})
}

func charCheck(c string) func(*Val) *Val {
	return func(s *Val) *Val {
		if RawString(s) == c {
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
	if unicode.IsSpace(rune(RawString(s)[0])) {
		return True
	}

	return False
}

func symbolToken(t *Val) *Val {
	v := NumberFromString(t)
	if IsError(v) == False {
		return v
	}

	v = BoolFromRawString(RawString(t))
	if IsError(v) == False {
		return v
	}

	return SymbolFromRawString(RawString(t))
}

func readChar(r *Val) *Val {
	in := Fread(Field(r, SymbolFromRawString("in")), NumberFromRawInt(1))
	st := Fstate(in)

	if IsError(st) != False {
		return Assign(r, FromMap(Struct{
			"in":    in,
			"value": st,
		}))
	}

	return Assign(r, FromMap(Struct{
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
	return Assign(r, FromMap(Struct{
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
	return Assign(r, FromMap(Struct{
		"current-token": StringFromRaw(""),
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
	return Assign(r, FromMap(Struct{"escaped": True}))
}

func unsetEscaped(r *Val) *Val {
	return Assign(r, FromMap(Struct{"escaped": False}))
}

func isEscaped(r *Val) *Val {
	return Field(r, SymbolFromRawString("escaped"))
}

func unescapeSymbolChar(c *Val) *Val {
	return c
}

func unescapeStringChar(c *Val) *Val {
	switch RawString(c) {
	case "b":
		return StringFromRaw("\b")
	case "f":
		return StringFromRaw("\f")
	case "n":
		return StringFromRaw("\n")
	case "r":
		return StringFromRaw("\r")
	case "t":
		return StringFromRaw("\t")
	case "v":
		return StringFromRaw("\v")
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

	return Assign(r, FromMap(Struct{
		"current-token": AppendString(Field(r, SymbolFromRawString("current-token")), c),
	}))
}

func setInvalid(r *Val) *Val {
	return Assign(r, FromMap(Struct{
		"err": invalidToken,
	}))
}

func processSymbol(r *Val) *Val {
	return Assign(r, FromMap(Struct{
		"value": symbolToken(Field(r, SymbolFromRawString("current-token"))),
	}))
}

func processString(r *Val) *Val {
	return Assign(r, FromMap(Struct{
		"value": Field(r, SymbolFromRawString("current-token")),
	}))
}

func closeChar(c *Val) *Val {
	switch RawString(c) {
	case "(":
		return StringFromRaw(")")
	case "[":
		return StringFromRaw("]")
	case "{":
		return StringFromRaw("}")
	default:
		return StringFromRaw("")
	}
}

func setClose(r, c *Val) *Val {
	if stringEq(closeChar(Field(r, SymbolFromRawString("in-list"))), c) == False {
		return setUnexpectedClose(r)
	}

	return Assign(r, FromMap(Struct{
		"close-list": True,
	}))
}

func hasCons(r *Val) bool {
	return Greater(Field(r, SymbolFromRawString("cons-items")), NumberFromRawInt(0)) != False
}

func consSet(r *Val) bool {
	return Field(r, SymbolFromRawString("cons")) != False
}

func setCons(r *Val) *Val {
	if hasCons(r) {
		return setIrregularCons(r)
	}

	return Assign(r, FromMap(Struct{
		"cons": True,
	}))
}

func setUnexpectedClose(r *Val) *Val {
	return Assign(r, FromMap(Struct{
		"value": unexpectedClose,
	}))
}

func setIrregularCons(r *Val) *Val {
	return Assign(r, FromMap(Struct{
		"value": irregularCons,
	}))
}

func reverse(l *Val) *Val {
	checkType(l, pair, mnil)

	r := Nil
	for {
		if l == Nil {
			return r
		}

		r = Cons(Car(l), r)
		l = Cdr(l)
	}
}

func reverseIrregular(l *Val) *Val {
	checkType(l, pair, mnil)

	r := Cons(Car(Cdr(l)), Car(l))
	l = Cdr(Cdr(l))
	for {
		if l == Nil {
			return r
		}

		r = Cons(Car(l), r)
		l = Cdr(l)
	}
}

func readList(r, c *Val) *Val {
	lr := reader(Field(r, SymbolFromRawString("in")))
	lr = Assign(lr, FromMap(Struct{
		"list-items": Nil,
		"in-list":    c,
	}))

	var loop func(*Val) *Val
	loop = func(lr *Val) *Val {
		lr = read(lr)
		if readError(lr) {
			return Assign(r, FromMap(Struct{
				"in":    Field(lr, SymbolFromRawString("in")),
				"value": Field(lr, SymbolFromRawString("value")),
			}))
		}

		v := Field(lr, SymbolFromRawString("value"))
		if v != UndefinedReadValue {
			lr = Assign(lr, FromMap(Struct{
				"list-items": Cons(
					v,
					Field(lr, SymbolFromRawString("list-items")),
				),
				"value": UndefinedReadValue,
			}))

			if hasCons(lr) {
				lr = Assign(lr, FromMap(Struct{
					"cons-items": Add(Field(lr, SymbolFromRawString("cons-items")), NumberFromRawInt(1)),
				}))
			}
		}

		if consSet(lr) {
			if Field(lr, SymbolFromRawString("list-items")) == Nil ||
				numberEq(Field(lr, SymbolFromRawString("cons-items")), NumberFromRawInt(0)) == False {
				return setIrregularCons(Assign(r, FromMap(Struct{
					"in": Field(lr, SymbolFromRawString("in")),
				})))
			}

			lr = Assign(lr, FromMap(Struct{
				"cons-items": NumberFromRawInt(1),
				"cons":       False,
			}))
		}

		if Field(lr, SymbolFromRawString("close-list")) != False {
			if hasCons(lr) {
				if numberEq(Field(lr, SymbolFromRawString("cons-items")), NumberFromRawInt(2)) == False {
					return setIrregularCons(Assign(r, FromMap(Struct{
						"in": Field(lr, SymbolFromRawString("in")),
					})))
				}

				return Assign(r, FromMap(Struct{
					"in":    Field(lr, SymbolFromRawString("in")),
					"value": reverseIrregular(Field(lr, SymbolFromRawString("list-items"))),
				}))
			}

			return Assign(r, FromMap(Struct{
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
	if stringEq(closeChar(Field(r, SymbolFromRawString("in-list"))), StringFromRaw("")) == False {
		lr = Assign(lr, FromMap(Struct{
			"in-list": Field(r, SymbolFromRawString("in-list")),
		}))
	}

	lr = read(lr)
	if readError(lr) {
		return Assign(r, FromMap(Struct{
			"in":    Field(lr, SymbolFromRawString("in")),
			"value": Field(lr, SymbolFromRawString("value")),
		}))
	}

	return Assign(r, FromMap(Struct{
		"in":         Field(lr, SymbolFromRawString("in")),
		"value":      List(SymbolFromRawString("quote"), Field(lr, SymbolFromRawString("value"))),
		"close-list": Field(lr, SymbolFromRawString("close-list")),
	}))
}

func readVector(r *Val) *Val {
	r = readList(r, StringFromRaw("["))
	if readError(r) {
		return r
	}

	return Assign(r, FromMap(Struct{
		"value": Cons(SymbolFromRawString("vector:"), Field(r, SymbolFromRawString("value"))),
	}))
}

func readStruct(r *Val) *Val {
	r = readList(r, StringFromRaw("{"))
	if readError(r) {
		return r
	}

	return Assign(r, FromMap(Struct{
		"value": Cons(SymbolFromRawString("struct:"), Field(r, SymbolFromRawString("value"))),
	}))
}

func read(r *Val) *Val {
	t := currentTokenType(r)
	if isTList(t) {
		return setNone(readList(r, StringFromRaw("(")))
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
