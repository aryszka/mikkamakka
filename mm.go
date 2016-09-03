/*
experiment bootstrapping a scheme value space
*/
package main

import "strings"

var (
	voidError      = &val{merror, "void error"}
	invalidToken   = &val{merror, "invalid token"}
	notImplemented = &val{merror, "not implemented"}
	psempty        = &val{number, 0}
	pssymbol       = &val{number, 1}
	psstring       = &val{number, 2}
	pscomment = &val{number, 3}
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
		"in":         in,
		"state":      voidError,
		"parseState": psempty,
	})
}

func isNewline(s *val) *val {
	if stringVal(s) == "\n" {
		return vtrue
	}

	return vfalse
}

func isWhiteSpace(s *val) *val {
	if isNewline(s) != vfalse || stringVal(s) == " " {
		return vtrue
	}

	return vfalse
}

func isDigit(s *val) *val {
	if strings.Index("0123456789", stringVal(s)) >= 0 {
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

func read(r *val) *val {
	in := field(r, sfromString("in"))

	var loop func(*val) *val
	loop = func(token *val) *val {
		in = fread(in, fromInt(1))
		st := fstate(in)

		ps := field(r, sfromString("parseState"))
		if isError(st) != vfalse {
			return fromMap(map[string]*val{
				"in":         in,
				"state":      st,
				"parseState": ps,
			})
		}

		switch ps {
		case psempty:
			if isWhiteSpace(st) != vfalse {
				return loop(token)
			}

			if isStringDelimiter(st) != vfalse {
				st = fromString("")
				r = fromMap(map[string]*val{
					"in":         in,
					"state":      st,
					"parseState": psstring,
				})

				return loop(st)
			}

			if isComment(st) != vfalse {
				st = fromString("")
				r = fromMap(map[string]*val{
					"in":         in,
					"state":      st,
					"parseState": pscomment,
				})

				return loop(st)
			}

			r = fromMap(map[string]*val{
				"in":         in,
				"state":      st,
				"parseState": pssymbol,
			})

			return loop(st)
		case pscomment:
			if isNewline(st) != vfalse {
				st = fromString("")
				r = fromMap(map[string]*val{
					"in": in,
					"state": st,
					"parseState": psempty,
				})
			}

			return loop(st)
		case pssymbol:
			if isWhiteSpace(st) != vfalse || isComment(st) != vfalse {
				pst := psempty
				if isComment(st) != vfalse {
					pst = pscomment
				}

				st = nfromString(stringVal(token))
				if isError(st) == vfalse {
					return fromMap(map[string]*val{
						"in":         in,
						"state":      st,
						"parseState": pst,
					})
				}

				st = bfromString(stringVal(token))
				if isError(st) == vfalse {
					return fromMap(map[string]*val{
						"in":         in,
						"state":      st,
						"parseState": pst,
					})
				}

				return fromMap(map[string]*val{
					"in":         in,
					"state":      sfromString(stringVal(token)),
					"parseState": pst,
				})
			}

			if isStringDelimiter(st) != vfalse {
				r = fromMap(map[string]*val{
					"in":         in,
					"state":      nfromString(stringVal(token)),
					"parseState": psstring,
				})

				return loop(fromString(""))
			}

			return loop(appendString(token, st))
		case psstring:
			if isStringDelimiter(st) != vfalse {
				r = fromMap(map[string]*val{
					"in":         in,
					"state":      token,
					"parseState": psempty,
				})

				return r
			}

			return loop(appendString(token, st))
		}

		return fromMap(map[string]*val{
			"in":         in,
			"state":      invalidToken,
			"parseState": field(r, sfromString("parseState")),
		})
	}

	return loop(fromString(""))
}

func printer(out *val) *val {
	return fromMap(map[string]*val{
		"out":   out,
		"state": voidError,
	})
}

func mprint(p, v *val) *val {
	f := field(p, sfromString("out"))

	if isSymbol(v) != vfalse {
		v = symbolToString(v)
	} else if isNumber(v) != vfalse {
		v = numberToString(v)
	} else if isString(v) != vfalse {
		v = appendString(fromString(`"`), v, fromString(`"`))
	} else if isBool(v) != vfalse {
		v = boolToString(v)
	} else {
		return fromMap(map[string]*val{
			"out":   f,
			"state": notImplemented,
		})
	}

	f = writeln(f, v)
	if st := fstate(f); isError(st) != vfalse {
		return fromMap(map[string]*val{
			"out":   p,
			"state": st,
		})
	}

	return fromMap(map[string]*val{
		"out":   f,
		"state": v,
	})
}

func loop(in, out *val) {
	// TODO:
	// - need to drain input for OSX terminal
	// - fix ctl-d behavior

	in = read(in)
	v := field(in, sfromString("state"))
	if isError(v) != vfalse && v != voidError {
		if v == eof {
			return
		}

		fatal(v)
	}

	out = mprint(out, v)
	v = field(out, sfromString("state"))
	if isError(v) != vfalse && v != voidError {
		fatal(v)
	}

	loop(in, out)
}

func main() {
	in := reader(stdin())
	out := printer(stdout())
	loop(in, out)
}

func main0() {
	sin := stdin()
	sout := stdout()
	input := fromString("")
	escaped := vfalse

	var loop func()
	loop = func() {
		sin = fread(sin, fromInt(1))
		istate := fstate(sin)
		if istate == eof {
			if greater(stringLength(input), fromInt(0)) == vfalse {
				return
			} else {
				sout = writeln(sout, fromString(""))
				if ostate := fstate(sout); isError(ostate) != vfalse {
					fatal(ostate)
				}

				sout = writeln(sout, input)
				if ostate := fstate(sout); isError(ostate) != vfalse {
					fatal(ostate)
				}

				input = fromString("")
				loop()
				return
			}
		} else if isError(istate) != vfalse {
			fatal(istate)
		}

		appendEscaped := func(s *val) *val {
			if escaped == vfalse {
				return vfalse
			}

			input = appendString(input, s)
			escaped = vfalse
			return vtrue
		}

		switch stringVal(istate) {
		case "\\":
			if appendEscaped(istate) == vfalse {
				escaped = vtrue
			}
		case "\n":
			if appendEscaped(istate) == vfalse && greater(stringLength(input), fromInt(0)) != vfalse {
				sout = writeln(sout, input)
				if ostate := fstate(sout); isError(ostate) != vfalse {
					fatal(ostate)
				}

				input = fromString("")
			}
		default:
			input = appendString(input, istate)
		}

		loop()
	}

	loop()
}
