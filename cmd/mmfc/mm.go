/*
experiment bootstrapping a scheme value space
*/
package main

import mm "github.com/aryszka/mikkamakka"

func loop(env, in *mm.Val) {
	// TODO:
	// - need to drain input for OSX terminal
	// - fix ctl-d behavior
	// - display errors

	in = mm.Read(in)
	v := mm.Field(in, mm.SfromString("value"))
	if mm.IsError(v) != mm.False {
		if v == mm.Eof {
			return
		}

		if v == mm.UndefinedReadValue {
			mm.Fatal(mm.FromString("failed to read value"))
			return
		}

		mm.Fatal(v)
		return
	}

	v = mm.Eval(env, v)

	loop(env, in)
}

func main() {
	env := mm.InitialEnv()
	in := mm.Reader(mm.Stdin())
	loop(env, in)
}
