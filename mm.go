/*
experiment bootstrapping a scheme value space
*/
package main

import ()

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
)

func writeln(f *val, s *val) *val {
	if greater(stringLength(s), fromInt(0)) != vfalse {
		f = fwrite(f, s)
		if isError(fstate(f)) != vfalse {
			return f
		}
	}

	return fwrite(f, fromString("\n"))
}

func reader(in *val) *val {
	return fromMap(map[string]*val{
		"in":            in,
		"token-type":    ttnone,
		"value":         voidError,
		"current-token": fromString(""),
		"in-list":       vfalse,
		"close-list":    vfalse,
		"cons":          vfalse,
		"cons-items":    fromInt(0),
	})
}

func isNewline(s *val) *val {
	if stringVal(s) == "\n" {
		return vtrue
	}

	return vfalse
}

func isWhitespace(s *val) *val {
	if isNewline(s) != vfalse || stringVal(s) == " " {
		return vtrue
	}

	return vfalse
}

func isStringDelimiter(s *val) *val {
	if stringVal(s) == `"` {
		return vtrue
	}

	return vfalse
}

func isComment(s *val) *val {
	if stringVal(s) == ";" {
		return vtrue
	}

	return vfalse
}

func isListOpen(s *val) *val {
	if stringVal(s) == "(" {
		return vtrue
	}

	return vfalse
}

func isListClose(s *val) *val {
	if stringVal(s) == ")" {
		return vtrue
	}

	return vfalse
}

func isCons(s *val) *val {
	if stringVal(s) == "." {
		return vtrue
	}

	return vfalse
}

func isQuote(s *val) *val {
	if stringVal(s) == "'" {
		return vtrue
	}

	return vfalse
}

func isOpenVector(s *val) *val {
	if stringVal(s) == "[" {
		return vtrue
	}

	return vfalse
}

func isCloseVector(s *val) *val {
	if stringVal(s) == "]" {
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

func setNone(r *val) *val    { return setTokenType(r, ttnone) }
func setString(r *val) *val  { return setTokenType(r, ttstring) }
func setComment(r *val) *val { return setTokenType(r, ttcomment) }
func setSymbol(r *val) *val  { return setTokenType(r, ttsymbol) }
func setList(r *val) *val    { return setTokenType(r, ttlist) }
func setQuote(r *val) *val   { return setTokenType(r, ttquote) }
func setVector(r *val) *val  { return setTokenType(r, ttvector) }

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

func appendToken(r *val) *val {
	return assign(r, fromMap(map[string]*val{
		"current-token": appendString(field(r, sfromString("current-token")), lastChar(r)),
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

func inList(r *val) bool {
	return field(r, sfromString("in-list")) != vfalse
}

func setClose(r *val) *val {
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

func readList(r *val) *val {
	lr := reader(field(r, sfromString("in")))
	lr = assign(lr, fromMap(map[string]*val{
		"list-items": vnil,
		"in-list":    vtrue,
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
				!neq(field(lr, sfromString("cons-items")), fromInt(0)) {
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
				if !neq(field(lr, sfromString("cons-items")), fromInt(2)) {
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
	lr = read(lr)
	if readError(lr) {
		return assign(r, fromMap(map[string]*val{
			"in":    field(lr, sfromString("in")),
			"value": field(lr, sfromString("value")),
		}))
	}

	return assign(r, fromMap(map[string]*val{
		"in":    field(lr, sfromString("in")),
		"value": list(sfromString("quote"), field(lr, sfromString("value"))),
	}))
}

func readVector(r *val) *val {
	r = readList(r)
	return assign(r, fromMap(map[string]*val{
		"value": vectorFromList(field(r, sfromString("value"))),
	}))
}

func read(r *val) *val {
	t := currentTokenType(r)
	if isTList(t) {
		return setNone(readList(r))
	}

	if isTQuote(t) {
		return setNone(readQuote(r))
	}

	if isTVector(t) {
		return setNone(readVector(r))
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
		case isStringDelimiter(c) != vfalse:
			return read(setString(r))
		case isComment(c) != vfalse:
			return read(setComment(r))
		case isListOpen(c) != vfalse:
			return read(setList(r))
		case isListClose(c) != vfalse:
			if !inList(r) {
				return setUnexpectedClose(r)
			}

			return setClose(r)
		case isCons(c) != vfalse:
			if !inList(r) {
				return setIrregularCons(r)
			}

			return setCons(r)
		case isQuote(c) != vfalse:
			return read(setQuote(r))
		case isOpenVector(c) != vfalse:
			return read(setVector(r))
		case isCloseVector(c) != vfalse:
			if !inList(r) {
				return setUnexpectedClose(r)
			}

			return setClose(r)
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
		case isWhitespace(c) != vfalse:
			return closeSymbol(r)
		case isComment(c) != vfalse:
			return setComment(closeSymbol(r))
		case isStringDelimiter(c) != vfalse:
			return setString(closeSymbol(r))
		case isListOpen(c) != vfalse:
			return setList(closeSymbol(r))
		case isListClose(c) != vfalse:
			if !inList(r) {
				return setUnexpectedClose(r)
			}

			return setClose(closeSymbol(r))
		case isCons(c) != vfalse:
			if !inList(r) {
				return setIrregularCons(r)
			}

			return setCons(closeSymbol(r))
		case isOpenVector(c) != vfalse:
			return setVector(closeSymbol(r))
		case isCloseVector(c) != vfalse:
			if !inList(r) {
				return setUnexpectedClose(r)
			}

			return setClose(closeSymbol(r))
		default:
			return read(appendToken(r))
		}
	case isTString(t):
		switch {
		case isStringDelimiter(c) != vfalse:
			return closeString(r)
		default:
			return read(appendToken(r))
		}
	default:
		return setInvalid(r)
	}
}

func printer(out *val) *val {
	return fromMap(map[string]*val{
		"out":   out,
		"state": vnil,
	})
}

func printState(p *val) *val {
	return field(p, sfromString("state"))
}

func printRaw(p *val, r *val) *val {
	f := fwrite(field(p, sfromString("out")), r)
	return assign(p, fromMap(map[string]*val{
		"out":   f,
		"state": fstate(f),
	}))
}

func printQuoteSign(p *val) *val {
	return printRaw(p, fromString("'"))
}

func printSymbol(p, v, q *val) *val {
	if q == vfalse {
		p = printQuoteSign(p)
	}

	return printRaw(p, symbolToString(v))
}

func printQuote(p, v *val) *val {
	p = printQuoteSign(p)
	return mprintq(p, car(cdr(v)), vfalse)
}

func printPair(p, v, q *val) *val {
	if q == vfalse {
		p = printQuoteSign(p)
	}

	p = printRaw(p, fromString("("))
	if st := printState(p); isError(st) != vfalse {
		return p
	}

	var loop func(*val, *val, *val) *val
	loop = func(p *val, v *val, first *val) *val {
		if isNil(v) != vfalse {
			return printRaw(p, fromString(")"))
		}

		if first == vfalse {
			p = printRaw(p, fromString(" "))
			if st := printState(p); isError(st) != vfalse {
				return p
			}
		}

		if isPair(cdr(v)) == vfalse && isNil(cdr(v)) == vfalse {
			p = mprintq(p, car(v), vtrue)
			if st := field(p, sfromString("state")); isError(st) != vfalse {
				return p
			}

			p = printRaw(p, fromString(" . "))
			if st := printState(p); isError(st) != vfalse {
				return p
			}

			p = mprintq(p, cdr(v), vtrue)
			if st := field(p, sfromString("state")); isError(st) != vfalse {
				return p
			}

			return printRaw(p, fromString(")"))
		}

		p = mprintq(p, car(v), vtrue)
		if st := field(p, sfromString("state")); isError(st) != vfalse {
			return p
		}

		return loop(p, cdr(v), vfalse)
	}

	return loop(p, v, vtrue)
}

func printVector(p, v *val) *val {
	p = printRaw(p, fromString("["))
	if st := field(p, sfromString("state")); isError(st) != vfalse {
		return p
	}

	var loop func(*val, *val, *val) *val
	loop = func(p, i, f *val) *val {
		if neq(i, vectorLength(v)) {
			return p
		}

		if f == vfalse {
			p = printRaw(p, fromString(" "))
			if st := field(p, sfromString("state")); isError(st) != vfalse {
				return p
			}
		}

		p = mprintq(p, vectorRef(v, i), vtrue)
		if st := field(p, sfromString("state")); isError(st) != vfalse {
			return p
		}

		return loop(p, add(i, fromInt(1)), vfalse)
	}

	p = loop(p, fromInt(0), vtrue)
	return printRaw(p, fromString("]"))
}

func mprintq(p, v, q *val) *val {
	if isSymbol(v) != vfalse {
		return printSymbol(p, v, q)
	} else if isNumber(v) != vfalse {
		v = numberToString(v)
	} else if isString(v) != vfalse {
		v = appendString(fromString(`"`), v, fromString(`"`))
	} else if isBool(v) != vfalse {
		v = boolToString(v)
	} else if isSys(v) != vfalse {
		v = sstring(v)
	} else if isError(v) != vfalse {
		v = estring(v)
	} else if isPair(v) != vfalse && isSymbol(car(v)) != vfalse && smeq(car(v), sfromString("quote")) != vfalse {
		return printQuote(p, v)
	} else if isPair(v) != vfalse || isNil(v) != vfalse {
		return printPair(p, v, q)
	} else if isVector(v) != vfalse {
		return printVector(p, v)
	} else {
		return assign(p, fromMap(map[string]*val{
			"state": notImplemented,
		}))
	}

	f := fwrite(field(p, sfromString("out")), v)
	if st := fstate(f); isError(st) != vfalse {
		return assign(p, fromMap(map[string]*val{
			"out":   f,
			"state": st,
		}))
	}

	return assign(p, fromMap(map[string]*val{
		"out":   f,
		"state": v,
	}))
}

func mprint(p, v *val) *val {
	return mprintq(p, v, vfalse)
}

func loop(in, out *val) {
	// TODO:
	// - need to drain input for OSX terminal
	// - fix ctl-d behavior
	// - display errors

	in = read(in)
	v := field(in, sfromString("value"))
	if isError(v) != vfalse {
		if v == eof {
			return
		}

		if v == voidError {
			fatal(&val{merror, "failed to read value"})
		}

		fatal(v)
	}

	out = mprint(out, v)
	v = field(out, sfromString("state"))
	if isError(v) != vfalse {
		fatal(v)
	}

	f := field(out, sfromString("out"))
	f = fwrite(f, fromString("\n"))
	if isError(fstate(f)) != vfalse {
		fatal(fstate(f))
	}

	loop(in, assign(out, fromMap(map[string]*val{
		"out": f,
	})))
}

func main() {
	in := reader(stdin())
	out := printer(stdout())
	loop(in, out)
}
