/*
experiment bootstrapping a scheme value space
*/
package main

func initialEnv() *val {
	env := newEnv()
	define(env, sfromString("nil"), vnil)
	return env
}

func loop(env, in, out *val) {
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

	v = eval(env, v)

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

	loop(env, in, assign(out, fromMap(map[string]*val{
		"out": f,
	})))
}

func main() {
	env := initialEnv()
	in := reader(stdin())
	out := printer(stdout())
	loop(env, in, out)
}
