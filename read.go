package mikkamakka

import "unicode"

var (
	invalidToken    = &Val{merror, "invalid token"}
	notImplemented  = &Val{merror, "not implemented"}
	unexpectedClose = &Val{merror, "unexpected close"}
	irregularCons   = &Val{merror, "irregular cons"}
	VoidError       = &Val{merror, "void error"}
	ttnone          = &Val{number, 0}
	ttcomment       = &Val{number, 1}
	ttsymbol        = &Val{number, 2}
	ttstring        = &Val{number, 3}
	ttlist          = &Val{number, 4}
	ttquote         = &Val{number, 5}
	ttvector        = &Val{number, 7}
	ttstruct        = &Val{number, 8}
)

func reader(in *Val) *Val {
	return fromMap(map[string]*Val{
		"in":            in,
		"token-type":    ttnone,
		"value":         VoidError,
		"escaped":       vfalse,
		"last-char":     fromString(""),
		"current-token": fromString(""),
		"in-list":       fromString(""),
		"close-list":    vfalse,
		"cons":          vfalse,
		"cons-items":    fromInt(0),
	})
}

func charCheck(c string) func(*Val) *Val {
	return func(s *Val) *Val {
		if stringVal(s) == c {
			return vtrue
		}

		return vfalse
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
	if unicode.IsSpace(rune(stringVal(s)[0])) {
		return vtrue
	}

	return vfalse
}

func symbolToken(t *Val) *Val {
	v := nfromString(stringVal(t))
	if isError(v) == vfalse {
		return v
	}

	v = bfromString(stringVal(t))
	if isError(v) == vfalse {
		return v
	}

	return sfromString(stringVal(t))
}

func readChar(r *Val) *Val {
	in := fread(field(r, sfromString("in")), fromInt(1))
	st := fstate(in)

	if isError(st) != vfalse {
		return assign(r, fromMap(map[string]*Val{
			"in":    in,
			"value": st,
		}))
	}

	return assign(r, fromMap(map[string]*Val{
		"in":        in,
		"last-char": st,
	}))
}

func readError(r *Val) bool {
	v := field(r, sfromString("value"))
	return isError(v) != vfalse && v != VoidError

}

func lastChar(r *Val) *Val {
	return field(r, sfromString("last-char"))
}

func currentTokenType(r *Val) *Val {
	return field(r, sfromString("token-type"))
}

func setTokenType(r *Val, t *Val) *Val {
	return assign(r, fromMap(map[string]*Val{
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
	return assign(r, fromMap(map[string]*Val{
		"current-token": fromString(""),
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
	return assign(r, fromMap(map[string]*Val{"escaped": vtrue}))
}

func unsetEscaped(r *Val) *Val {
	return assign(r, fromMap(map[string]*Val{"escaped": vfalse}))
}

func isEscaped(r *Val) *Val {
	return field(r, sfromString("escaped"))
}

func unescapeSymbolChar(c *Val) *Val {
	return c
}

func unescapeStringChar(c *Val) *Val {
	switch stringVal(c) {
	case "b":
		return fromString("\b")
	case "f":
		return fromString("\f")
	case "n":
		return fromString("\n")
	case "r":
		return fromString("\r")
	case "t":
		return fromString("\t")
	case "v":
		return fromString("\v")
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
	if isEscaped(r) != vfalse {
		c = unescapeChar(field(r, sfromString("token-type")), c)
	}

	return assign(r, fromMap(map[string]*Val{
		"current-token": appendString(field(r, sfromString("current-token")), c),
	}))
}

func setInvalid(r *Val) *Val {
	return assign(r, fromMap(map[string]*Val{
		"err": invalidToken,
	}))
}

func processSymbol(r *Val) *Val {
	return assign(r, fromMap(map[string]*Val{
		"value": symbolToken(field(r, sfromString("current-token"))),
	}))
}

func processString(r *Val) *Val {
	return assign(r, fromMap(map[string]*Val{
		"value": field(r, sfromString("current-token")),
	}))
}

func closeChar(c *Val) *Val {
	switch stringVal(c) {
	case "(":
		return fromString(")")
	case "[":
		return fromString("]")
	case "{":
		return fromString("}")
	default:
		return fromString("")
	}
}

func setClose(r, c *Val) *Val {
	if seq(closeChar(field(r, sfromString("in-list"))), c) == vfalse {
		return setUnexpectedClose(r)
	}

	return assign(r, fromMap(map[string]*Val{
		"close-list": vtrue,
	}))
}

func hasCons(r *Val) bool {
	return greater(field(r, sfromString("cons-items")), fromInt(0)) != vfalse
}

func consSet(r *Val) bool {
	return field(r, sfromString("cons")) != vfalse
}

func setCons(r *Val) *Val {
	if hasCons(r) {
		return setIrregularCons(r)
	}

	return assign(r, fromMap(map[string]*Val{
		"cons": vtrue,
	}))
}

func setUnexpectedClose(r *Val) *Val {
	return assign(r, fromMap(map[string]*Val{
		"value": unexpectedClose,
	}))
}

func setIrregularCons(r *Val) *Val {
	return assign(r, fromMap(map[string]*Val{
		"value": irregularCons,
	}))
}

func readList(r, c *Val) *Val {
	lr := reader(field(r, sfromString("in")))
	lr = assign(lr, fromMap(map[string]*Val{
		"list-items": Nil,
		"in-list":    c,
	}))

	var loop func(*Val) *Val
	loop = func(lr *Val) *Val {
		lr = read(lr)
		if readError(lr) {
			return assign(r, fromMap(map[string]*Val{
				"in":    field(lr, sfromString("in")),
				"value": field(lr, sfromString("value")),
			}))
		}

		v := field(lr, sfromString("value"))
		if v != VoidError {
			lr = assign(lr, fromMap(map[string]*Val{
				"list-items": cons(
					v,
					field(lr, sfromString("list-items")),
				),
				"value": VoidError,
			}))

			if hasCons(lr) {
				lr = assign(lr, fromMap(map[string]*Val{
					"cons-items": add(field(lr, sfromString("cons-items")), fromInt(1)),
				}))
			}
		}

		if consSet(lr) {
			if field(lr, sfromString("list-items")) == Nil ||
				neq(field(lr, sfromString("cons-items")), fromInt(0)) == vfalse {
				return setIrregularCons(assign(r, fromMap(map[string]*Val{
					"in": field(lr, sfromString("in")),
				})))
			}

			lr = assign(lr, fromMap(map[string]*Val{
				"cons-items": fromInt(1),
				"cons":       vfalse,
			}))
		}

		if field(lr, sfromString("close-list")) != vfalse {
			if hasCons(lr) {
				if neq(field(lr, sfromString("cons-items")), fromInt(2)) == vfalse {
					return setIrregularCons(assign(r, fromMap(map[string]*Val{
						"in": field(lr, sfromString("in")),
					})))
				}

				return assign(r, fromMap(map[string]*Val{
					"in":    field(lr, sfromString("in")),
					"value": reverseIrregular(field(lr, sfromString("list-items"))),
				}))
			}

			return assign(r, fromMap(map[string]*Val{
				"in":    field(lr, sfromString("in")),
				"value": reverse(field(lr, sfromString("list-items"))),
			}))
		}

		return loop(lr)
	}

	return loop(lr)
}

func readQuote(r *Val) *Val {
	lr := reader(field(r, sfromString("in")))
	if seq(closeChar(field(r, sfromString("in-list"))), fromString("")) == vfalse {
		lr = assign(lr, fromMap(map[string]*Val{
			"in-list": field(r, sfromString("in-list")),
		}))
	}

	lr = read(lr)
	if readError(lr) {
		return assign(r, fromMap(map[string]*Val{
			"in":    field(lr, sfromString("in")),
			"value": field(lr, sfromString("value")),
		}))
	}

	return assign(r, fromMap(map[string]*Val{
		"in":         field(lr, sfromString("in")),
		"value":      list(sfromString("quote"), field(lr, sfromString("value"))),
		"close-list": field(lr, sfromString("close-list")),
	}))
}

func readVector(r *Val) *Val {
	r = readList(r, fromString("["))
	if readError(r) {
		return r
	}

	return assign(r, fromMap(map[string]*Val{
		"value": cons(sfromString("vector:"), field(r, sfromString("value"))),
	}))
}

func readStruct(r *Val) *Val {
	r = readList(r, fromString("{"))
	if readError(r) {
		return r
	}

	return assign(r, fromMap(map[string]*Val{
		"value": cons(sfromString("struct:"), field(r, sfromString("value"))),
	}))
}

func read(r *Val) *Val {
	t := currentTokenType(r)
	if isTList(t) {
		return setNone(readList(r, fromString("(")))
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
		case isWhitespace(c) != vfalse:
			return read(r)
		case isEscapeChar(c) != vfalse:
			return read(setEscaped(setSymbol(r)))
		case isStringDelimiter(c) != vfalse:
			return read(setString(r))
		case isComment(c) != vfalse:
			return read(setComment(r))
		case isListOpen(c) != vfalse:
			return read(setList(r))
		case isListClose(c) != vfalse:
			return setClose(r, c)
		case isCons(c) != vfalse:
			return setCons(r)
		case isQuoteChar(c) != vfalse:
			return read(setQuote(r))
		case isOpenVector(c) != vfalse:
			return read(setVector(r))
		case isCloseVector(c) != vfalse:
			return setClose(r, c)
		case isOpenStruct(c) != vfalse:
			return read(setStruct(r))
		case isCloseStruct(c) != vfalse:
			return setClose(r, c)
		default:
			return read(appendToken(setSymbol(r)))
		}

	case isTComment(t):
		switch {
		case isNewline(c) != vfalse:
			return read(closeComment(r))
		}

		return read(r)

	case isTSymbol(t):
		switch {
		case isEscaped(r) != vfalse:
			return read(unsetEscaped(appendToken(r)))
		case isWhitespace(c) != vfalse:
			return closeSymbol(r)
		case isEscapeChar(c) != vfalse:
			return read(setEscaped(r))
		case isComment(c) != vfalse:
			return setComment(closeSymbol(r))
		case isStringDelimiter(c) != vfalse:
			return setString(closeSymbol(r))
		case isListOpen(c) != vfalse:
			return setList(closeSymbol(r))
		case isListClose(c) != vfalse:
			return setClose(closeSymbol(r), c)
		case isCons(c) != vfalse:
			return setCons(closeSymbol(r))
		case isOpenVector(c) != vfalse:
			return setVector(closeSymbol(r))
		case isCloseVector(c) != vfalse:
			return setClose(closeSymbol(r), c)
		case isOpenStruct(c) != vfalse:
			return setStruct(closeSymbol(r))
		case isCloseStruct(c) != vfalse:
			return setClose(closeSymbol(r), c)
		default:
			return read(appendToken(r))
		}

	case isTString(t):
		switch {
		case isEscaped(r) != vfalse:
			return read(unsetEscaped(appendToken(r)))
		case isEscapeChar(c) != vfalse:
			return read(setEscaped(r))
		case isStringDelimiter(c) != vfalse:
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
