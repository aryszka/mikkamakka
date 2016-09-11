/*
experiment bootstrapping a scheme value space
*/
package main

func initialEnv() *val {
	env := newEnv(nil)

	define(env, sfromString("nil"), vnil)
	define(env, sfromString("nil?"), newBuiltin(bisNil, 1, false))
	define(env, sfromString("pair?"), newBuiltin(bisPair, 1, false))
	define(env, sfromString("cons"), newBuiltin(bcons, 2, false))
	define(env, sfromString("car"), newBuiltin(bcar, 1, false))
	define(env, sfromString("cdr"), newBuiltin(bcdr, 1, false))
	define(env, sfromString("list"), newBuiltin(blist, 0, true))
	define(env, sfromString("apply"), newBuiltin(bapply, 2, false))
	define(env, sfromString("error?"), newBuiltin(bisError, 1, false))
	define(env, sfromString("string->error"), newBuiltin(stringToError, 1, false))
	define(env, sfromString("not"), newBuiltin(not, 1, false))
	define(env, sfromString("="), newBuiltin(beq, 0, true))
	define(env, sfromString(">"), newBuiltin(bgreater, 2, false))
	define(env, sfromString("+"), newBuiltin(badd, 0, true))
	define(env, sfromString("try-string->number"), newBuiltin(btryNumberFromString, 1, false))
	define(env, sfromString("try-string->bool"), newBuiltin(btryBoolFromString, 1, false))
	define(env, sfromString("symbol?"), newBuiltin(bisSymbol, 1, false))
	define(env, sfromString("string->symbol"), newBuiltin(stringToSymbol, 1, false))
	define(env, sfromString("number?"), newBuiltin(bisNumber, 1, false))
	define(env, sfromString("number->string"), newBuiltin(bnumberToString, 1, false))
	define(env, sfromString("bool?"), newBuiltin(bisBool, 1, false))
	define(env, sfromString("string?"), newBuiltin(bisString, 1, false))
	define(env, sfromString("assign"), newBuiltin(bassign, 1, true))
	define(env, sfromString("fopen"), newBuiltin(bfopen, 1, false))
	define(env, sfromString("fclose"), newBuiltin(bfclose, 1, false))
	define(env, sfromString("fread"), newBuiltin(bfread, 2, false))
	define(env, sfromString("fwrite"), newBuiltin(bfwrite, 2, false))
	define(env, sfromString("fstate"), newBuiltin(bfstate, 1, false))
	define(env, sfromString("derived-object?"), newBuiltin(derivedObject, 2, false))
	define(env, sfromString("failing-reader"), newBuiltin(failingReader, 0, false))
	define(env, sfromString("eof"), eof)
	define(env, sfromString("stdin"), newBuiltin(bstdin, 0, false))
	define(env, sfromString("stderr"), newBuiltin(bstderr, 0, false))
	define(env, sfromString("stdout"), newBuiltin(bstdout, 0, false))
	define(env, sfromString("buffer"), newBuiltin(bbuffer, 0, false))
	define(env, sfromString("argv"), newBuiltin(argv, 0, false))
	define(env, sfromString("invalid-token"), invalidToken)
	define(env, sfromString("string-append"), newBuiltin(bappendString, 0, true))
	define(env, sfromString("escape-compiled-string"), newBuiltin(escapeCompiled, 1, false))
	define(env, sfromString("printer"), newBuiltin(bprinter, 1, false))
	define(env, sfromString("print"), newBuiltin(bprint, 2, false))

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

	// out = mprint(out, v)
	// v = field(out, sfromString("state"))
	// if isError(v) != vfalse {
	// 	fatal(v)
	// }

	f := field(out, sfromString("output"))
	// f = fwrite(f, fromString("\n"))
	// if isError(fstate(f)) != vfalse {
	// 	fatal(fstate(f))
	// }

	loop(env, in, assign(out, fromMap(map[string]*val{
		"output": f,
	})))
}

func main() {
	env := initialEnv()
	in := reader(stdin())
	out := printer(stdout())
	loop(env, in, out)
}
