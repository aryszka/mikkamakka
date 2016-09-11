package main

import "unicode"

var (
	invalidToken    = &val{merror, "invalid token"}
	notImplemented  = &val{merror, "not implemented"}
	unexpectedClose = &val{merror, "unexpected close"}
	irregularCons   = &val{merror, "irregular cons"}
	voidError       = &val{merror, "void error"}
	ttnone          = &val{number, 0}
	ttcomment       = &val{number, 1}
	ttsymbol        = &val{number, 2}
	ttstring        = &val{number, 3}
	ttlist          = &val{number, 4}
	ttquote         = &val{number, 5}
	ttvector        = &val{number, 7}
	ttstruct        = &val{number, 8}
)

func reader(in *val) *val {
	return fromMap(map[string]*val{
		"in":            in,
		"token-type":    ttnone,
		"value":         voidError,
		"escaped":       vfalse,
		"last-char":     fromString(""),
		"current-token": fromString(""),
		"in-list":       fromString(""),
		"close-list":    vfalse,
		"cons":          vfalse,
		"cons-items":    fromInt(0),
	})
}

func charCheck(c string) func(*val) *val {
	return func(s *val) *val {
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

func isWhitespace(s *val) *val {
	if unicode.IsSpace(rune(stringVal(s)[0])) {
		return vtrue
	}

	return vfalse
}

func symbolToken(t *val) *val {
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

func readChar(r *val) *val {
	in := fread(field(r, sfromString("in")), fromInt(1))
	st := fstate(in)

	if isError(st) != vfalse {
		return assign(r, fromMap(map[string]*val{
			"in":    in,
			"value": st,
		}))
	}

	return assign(r, fromMap(map[string]*val{
		"in":        in,
		"last-char": st,
	}))
}

func readError(r *val) bool {
	v := field(r, sfromString("value"))
	return isError(v) != vfalse && v != voidError

}

func lastChar(r *val) *val {
	return field(r, sfromString("last-char"))
}

func currentTokenType(r *val) *val {
	return field(r, sfromString("token-type"))
}

func setTokenType(r *val, t *val) *val {
	return assign(r, fromMap(map[string]*val{
		"token-type": t,
	}))
}

func isTNone(t *val) bool    { return t == ttnone }
func isTComment(t *val) bool { return t == ttcomment }
func isTSymbol(t *val) bool  { return t == ttsymbol }
func isTString(t *val) bool  { return t == ttstring }
func isTList(t *val) bool    { return t == ttlist }
func isTQuote(t *val) bool   { return t == ttquote }
func isTVector(t *val) bool  { return t == ttvector }
func isTStruct(t *val) bool  { return t == ttstruct }

func setNone(r *val) *val    { return setTokenType(r, ttnone) }
func setString(r *val) *val  { return setTokenType(r, ttstring) }
func setComment(r *val) *val { return setTokenType(r, ttcomment) }
func setSymbol(r *val) *val  { return setTokenType(r, ttsymbol) }
func setList(r *val) *val    { return setTokenType(r, ttlist) }
func setQuote(r *val) *val   { return setTokenType(r, ttquote) }
func setVector(r *val) *val  { return setTokenType(r, ttvector) }
func setStruct(r *val) *val  { return setTokenType(r, ttstruct) }

func clearToken(r *val) *val {
	return assign(r, fromMap(map[string]*val{
		"current-token": fromString(""),
	}))
}

func closeComment(r *val) *val {
	return clearToken(setTokenType(r, ttnone))
}

func closeSymbol(r *val) *val {
	return clearToken(processSymbol(setTokenType(r, ttnone)))
}

func closeString(r *val) *val {
	return clearToken(processString(setTokenType(r, ttnone)))
}

func setEscaped(r *val) *val {
	return assign(r, fromMap(map[string]*val{"escaped": vtrue}))
}

func unsetEscaped(r *val) *val {
	return assign(r, fromMap(map[string]*val{"escaped": vfalse}))
}

func isEscaped(r *val) *val {
	return field(r, sfromString("escaped"))
}

func unescapeSymbolChar(c *val) *val {
	return c
}

func unescapeStringChar(c *val) *val {
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

func unescapeChar(tokenType, c *val) *val {
	switch tokenType {
	case ttsymbol:
		return unescapeSymbolChar(c)
	case ttstring:
		return unescapeStringChar(c)
	default:
		return invalidToken
	}
}

func appendToken(r *val) *val {
	c := lastChar(r)
	if isEscaped(r) != vfalse {
		c = unescapeChar(field(r, sfromString("token-type")), c)
	}

	return assign(r, fromMap(map[string]*val{
		"current-token": appendString(field(r, sfromString("current-token")), c),
	}))
}

func setInvalid(r *val) *val {
	return assign(r, fromMap(map[string]*val{
		"err": invalidToken,
	}))
}

func processSymbol(r *val) *val {
	return assign(r, fromMap(map[string]*val{
		"value": symbolToken(field(r, sfromString("current-token"))),
	}))
}

func processString(r *val) *val {
	return assign(r, fromMap(map[string]*val{
		"value": field(r, sfromString("current-token")),
	}))
}

func closeChar(c *val) *val {
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

func setClose(r, c *val) *val {
	if seq(closeChar(field(r, sfromString("in-list"))), c) == vfalse {
		return setUnexpectedClose(r)
	}

	return assign(r, fromMap(map[string]*val{
		"close-list": vtrue,
	}))
}

func hasCons(r *val) bool {
	return greater(field(r, sfromString("cons-items")), fromInt(0)) != vfalse
}

func consSet(r *val) bool {
	return field(r, sfromString("cons")) != vfalse
}

func setCons(r *val) *val {
	if hasCons(r) {
		return setIrregularCons(r)
	}

	return assign(r, fromMap(map[string]*val{
		"cons": vtrue,
	}))
}

func setUnexpectedClose(r *val) *val {
	return assign(r, fromMap(map[string]*val{
		"value": unexpectedClose,
	}))
}

func setIrregularCons(r *val) *val {
	return assign(r, fromMap(map[string]*val{
		"value": irregularCons,
	}))
}

func readList(r, c *val) *val {
	lr := reader(field(r, sfromString("in")))
	lr = assign(lr, fromMap(map[string]*val{
		"list-items": vnil,
		"in-list":    c,
	}))

	var loop func(*val) *val
	loop = func(lr *val) *val {
		lr = read(lr)
		if readError(lr) {
			return assign(r, fromMap(map[string]*val{
				"in":    field(lr, sfromString("in")),
				"value": field(lr, sfromString("value")),
			}))
		}

		v := field(lr, sfromString("value"))
		if v != voidError {
			lr = assign(lr, fromMap(map[string]*val{
				"list-items": cons(
					v,
					field(lr, sfromString("list-items")),
				),
				"value": voidError,
			}))

			if hasCons(lr) {
				lr = assign(lr, fromMap(map[string]*val{
					"cons-items": add(field(lr, sfromString("cons-items")), fromInt(1)),
				}))
			}
		}

		if consSet(lr) {
			if field(lr, sfromString("list-items")) == vnil ||
				neq(field(lr, sfromString("cons-items")), fromInt(0)) == vfalse {
				return setIrregularCons(assign(r, fromMap(map[string]*val{
					"in": field(lr, sfromString("in")),
				})))
			}

			lr = assign(lr, fromMap(map[string]*val{
				"cons-items": fromInt(1),
				"cons":       vfalse,
			}))
		}

		if field(lr, sfromString("close-list")) != vfalse {
			if hasCons(lr) {
				if neq(field(lr, sfromString("cons-items")), fromInt(2)) == vfalse {
					return setIrregularCons(assign(r, fromMap(map[string]*val{
						"in": field(lr, sfromString("in")),
					})))
				}

				return assign(r, fromMap(map[string]*val{
					"in":    field(lr, sfromString("in")),
					"value": reverseIrregular(field(lr, sfromString("list-items"))),
				}))
			}

			return assign(r, fromMap(map[string]*val{
				"in":    field(lr, sfromString("in")),
				"value": reverse(field(lr, sfromString("list-items"))),
			}))
		}

		return loop(lr)
	}

	return loop(lr)
}

func readQuote(r *val) *val {
	lr := reader(field(r, sfromString("in")))
	if seq(closeChar(field(r, sfromString("in-list"))), fromString("")) == vfalse {
		lr = assign(lr, fromMap(map[string]*val{
			"in-list": field(r, sfromString("in-list")),
		}))
	}

	lr = read(lr)
	if readError(lr) {
		return assign(r, fromMap(map[string]*val{
			"in":    field(lr, sfromString("in")),
			"value": field(lr, sfromString("value")),
		}))
	}

	return assign(r, fromMap(map[string]*val{
		"in":         field(lr, sfromString("in")),
		"value":      list(sfromString("quote"), field(lr, sfromString("value"))),
		"close-list": field(lr, sfromString("close-list")),
	}))
}

func readVector(r *val) *val {
	r = readList(r, fromString("["))
	if readError(r) {
		return r
	}

	return assign(r, fromMap(map[string]*val{
		"value": cons(sfromString("vector:"), field(r, sfromString("value"))),
	}))
}

func readStruct(r *val) *val {
	r = readList(r, fromString("{"))
	if readError(r) {
		return r
	}

	return assign(r, fromMap(map[string]*val{
		"value": cons(sfromString("struct:"), field(r, sfromString("value"))),
	}))
}

func read(r *val) *val {
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
