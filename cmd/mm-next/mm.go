package main

import mm "github.com/aryszka/mikkamakka"

func main() {
	initialEnv := mm.InitialEnv()
	env := initialEnv
	mm.ModuleLoader(env, mm.SysStringToString("scm/lang.scm"), func(eenv *mm.Val) {
		env := mm.ModuleEnv(eenv, mm.SysStringToString("scm/lang.scm"))
		mm.Define(env, mm.SymbolFromRawString("trace"), mm.NewCompiled(1, true, func(a []*mm.Val) *mm.Val {
			env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("message"), mm.SymbolFromRawString("values")), mm.SliceToList(a))
			env = env
			mm.Define(env, mm.SymbolFromRawString("trace"), mm.NewCompiled(2, false, func(a []*mm.Val) *mm.Val {
				env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("out"), mm.Cons(mm.SymbolFromRawString("values"), mm.NilVal)), mm.SliceToList(a))
				env = env
				return mm.ApplySys(mm.NewCompiled(0, false, func(a []*mm.Val) *mm.Val {
					env := mm.ExtendEnv(env, mm.NilVal, mm.SliceToList(a))
					env = env
					mm.Define(env, mm.SymbolFromRawString("out"), mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("print")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("out")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("values")), mm.NilVal)), mm.NilVal))))
					return func() *mm.Val {
						if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("nil?")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("values")), mm.NilVal)), mm.NilVal)) != mm.False {
							return func() *mm.Val {
								mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("fwrite")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("out:output")), mm.Cons(mm.SysStringToString("\n"), mm.NilVal)))
								return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("values")), mm.NilVal))
							}()
						} else {
							return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("trace")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("assign")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("out")), mm.Cons(mm.ListToStruct(mm.Cons(mm.SymbolFromRawString("output"), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("fwrite")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("out:output")), mm.Cons(mm.SysStringToString(" "), mm.NilVal))), mm.NilVal))), mm.NilVal))), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("values")), mm.NilVal)), mm.NilVal)))
						}
					}()
				}), mm.NilVal)
			}))
			return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("trace")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("printer")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("stderr")), mm.NilVal), mm.NilVal)), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cons")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("message")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("values")), mm.NilVal))), mm.NilVal)))
		}))
		mm.Define(env, mm.SymbolFromRawString("id"), mm.NewCompiled(1, false, func(a []*mm.Val) *mm.Val {
			env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("x"), mm.NilVal), mm.SliceToList(a))
			env = env
			return mm.LookupDef(env, mm.SymbolFromRawString("x"))
		}))
		mm.Define(env, mm.SymbolFromRawString("list"), mm.NewCompiled(0, true, func(a []*mm.Val) *mm.Val {
			env := mm.ExtendEnv(env, mm.SymbolFromRawString("x"), mm.SliceToList(a))
			env = env
			return mm.LookupDef(env, mm.SymbolFromRawString("x"))
		}))
		mm.Define(env, mm.SymbolFromRawString("apply"), mm.NewCompiled(2, false, func(a []*mm.Val) *mm.Val {
			env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("f"), mm.Cons(mm.SymbolFromRawString("a"), mm.NilVal)), mm.SliceToList(a))
			env = env
			return func() *mm.Val {
				if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("vector?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("f")), mm.NilVal)) != mm.False {
					return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("vector-ref")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("f")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("a")), mm.NilVal)), mm.NilVal)))
				} else {
					return func() *mm.Val {
						if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("struct?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("f")), mm.NilVal)) != mm.False {
							return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("field")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("f")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("a")), mm.NilVal)), mm.NilVal)))
						} else {
							return func() *mm.Val {
								if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compiled-function?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("f")), mm.NilVal)) != mm.False {
									return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("apply-compiled")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("f")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("a")), mm.NilVal)))
								} else {
									return mm.ApplySys(mm.NewCompiled(0, false, func(a []*mm.Val) *mm.Val {
										env := mm.ExtendEnv(env, mm.NilVal, mm.SliceToList(a))
										env = env
										mm.Define(env, mm.SymbolFromRawString("c"), mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("composite")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("f")), mm.NilVal)))
										return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("eval-seq")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("extend-env")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.NilVal)), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.NilVal)), mm.NilVal)), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("a")), mm.NilVal)))), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.NilVal)), mm.NilVal)), mm.NilVal)))
									}), mm.NilVal)
								}
							}()
						}
					}()
				}
			}()
		}))
		mm.Define(env, mm.SymbolFromRawString("call"), mm.NewCompiled(1, true, func(a []*mm.Val) *mm.Val {
			env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("f"), mm.SymbolFromRawString("a")), mm.SliceToList(a))
			env = env
			return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("apply")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("f")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("a")), mm.NilVal)))
		}))
		mm.Define(env, mm.SymbolFromRawString("fold"), mm.NewCompiled(3, false, func(a []*mm.Val) *mm.Val {
			env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("f"), mm.Cons(mm.SymbolFromRawString("i"), mm.Cons(mm.SymbolFromRawString("l"), mm.NilVal))), mm.SliceToList(a))
			env = env
			return func() *mm.Val {
				if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("nil?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("l")), mm.NilVal)) != mm.False {
					return mm.LookupDef(env, mm.SymbolFromRawString("i"))
				} else {
					return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("fold")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("f")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("f")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("l")), mm.NilVal)), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("i")), mm.NilVal))), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("l")), mm.NilVal)), mm.NilVal))))
				}
			}()
		}))
		mm.Define(env, mm.SymbolFromRawString("foldr"), mm.NewCompiled(3, false, func(a []*mm.Val) *mm.Val {
			env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("f"), mm.Cons(mm.SymbolFromRawString("i"), mm.Cons(mm.SymbolFromRawString("l"), mm.NilVal))), mm.SliceToList(a))
			env = env
			return func() *mm.Val {
				if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("nil?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("l")), mm.NilVal)) != mm.False {
					return mm.LookupDef(env, mm.SymbolFromRawString("i"))
				} else {
					return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("f")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("l")), mm.NilVal)), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("foldr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("f")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("i")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("l")), mm.NilVal)), mm.NilVal)))), mm.NilVal)))
				}
			}()
		}))
		mm.Define(env, mm.SymbolFromRawString("map"), mm.NewCompiled(1, true, func(a []*mm.Val) *mm.Val {
			env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("f"), mm.SymbolFromRawString("l")), mm.SliceToList(a))
			env = env
			return func() *mm.Val {
				if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("nil?")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("l")), mm.NilVal)), mm.NilVal)) != mm.False {
					return mm.NilVal
				} else {
					return func() *mm.Val {
						if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("nil?")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("l")), mm.NilVal)), mm.NilVal)) != mm.False {
							return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cons")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("f")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("l")), mm.NilVal)), mm.NilVal)), mm.NilVal)), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("map")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("f")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("l")), mm.NilVal)), mm.NilVal)), mm.NilVal))), mm.NilVal)))
						} else {
							return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cons")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("apply")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("f")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("map")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("l")), mm.NilVal))), mm.NilVal))), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("apply")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("map")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cons")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("f")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("map")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("l")), mm.NilVal))), mm.NilVal))), mm.NilVal))), mm.NilVal)))
						}
					}()
				}
			}()
		}))
		mm.Define(env, mm.SymbolFromRawString("append"), mm.NewCompiled(0, true, func(a []*mm.Val) *mm.Val {
			env := mm.ExtendEnv(env, mm.SymbolFromRawString("l"), mm.SliceToList(a))
			env = env
			return func() *mm.Val {
				if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("nil?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("l")), mm.NilVal)) != mm.False {
					return mm.LookupDef(env, mm.SymbolFromRawString("nil"))
				} else {
					return func() *mm.Val {
						if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("nil?")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("l")), mm.NilVal)), mm.NilVal)) != mm.False {
							return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("l")), mm.NilVal))
						} else {
							return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("foldr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("cons")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("apply")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("append")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("l")), mm.NilVal)), mm.NilVal))), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("l")), mm.NilVal)), mm.NilVal))))
						}
					}()
				}
			}()
		}))
		mm.Define(env, mm.SymbolFromRawString("part"), mm.NewCompiled(1, true, func(a []*mm.Val) *mm.Val {
			env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("f"), mm.SymbolFromRawString("a")), mm.SliceToList(a))
			env = env
			return mm.NewCompiled(0, true, func(a []*mm.Val) *mm.Val {
				env := mm.ExtendEnv(env, mm.SymbolFromRawString("b"), mm.SliceToList(a))
				env = env
				return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("apply")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("f")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("append")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("a")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("b")), mm.NilVal))), mm.NilVal)))
			})
		}))
		mm.Define(env, mm.SymbolFromRawString("partr"), mm.NewCompiled(1, true, func(a []*mm.Val) *mm.Val {
			env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("f"), mm.SymbolFromRawString("a")), mm.SliceToList(a))
			env = env
			return mm.NewCompiled(0, true, func(a []*mm.Val) *mm.Val {
				env := mm.ExtendEnv(env, mm.SymbolFromRawString("b"), mm.SliceToList(a))
				env = env
				return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("apply")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("f")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("append")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("b")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("a")), mm.NilVal))), mm.NilVal)))
			})
		}))
		mm.Define(env, mm.SymbolFromRawString("reverse"), mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("part")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("fold")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("cons")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("nil")), mm.NilVal)))))
		mm.Define(env, mm.SymbolFromRawString("reverse-irregular"), mm.NewCompiled(1, false, func(a []*mm.Val) *mm.Val {
			env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("l"), mm.NilVal), mm.SliceToList(a))
			env = env
			return func() *mm.Val {
				if func() *mm.Val {
					if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("nil?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("l")), mm.NilVal)) != mm.False {
						return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("nil?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("l")), mm.NilVal))
					} else {
						return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("nil?")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("l")), mm.NilVal)), mm.NilVal))
					}
				}() != mm.False {
					return mm.LookupDef(env, mm.SymbolFromRawString("irregular-cons"))
				} else {
					return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("fold")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("cons")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cons")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("l")), mm.NilVal)), mm.NilVal)), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("l")), mm.NilVal)), mm.NilVal))), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("l")), mm.NilVal)), mm.NilVal)), mm.NilVal))))
				}
			}()
		}))
		mm.Define(env, mm.SymbolFromRawString("inc"), mm.NewCompiled(1, false, func(a []*mm.Val) *mm.Val {
			env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("n"), mm.NilVal), mm.SliceToList(a))
			env = env
			return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("+")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("n")), mm.Cons(mm.SysIntToNumber(1), mm.NilVal)))
		}))
		mm.Define(env, mm.SymbolFromRawString("dec"), mm.NewCompiled(1, false, func(a []*mm.Val) *mm.Val {
			env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("n"), mm.NilVal), mm.SliceToList(a))
			env = env
			return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("-")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("n")), mm.Cons(mm.SysIntToNumber(1), mm.NilVal)))
		}))
		mm.Define(env, mm.SymbolFromRawString(">="), mm.NewCompiled(0, true, func(a []*mm.Val) *mm.Val {
			env := mm.ExtendEnv(env, mm.SymbolFromRawString("n"), mm.SliceToList(a))
			env = env
			return func() *mm.Val {
				if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("nil?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("n")), mm.NilVal)) != mm.False {
					return mm.False
				} else {
					return func() *mm.Val {
						if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("nil?")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("n")), mm.NilVal)), mm.NilVal)) != mm.False {
							return mm.True
						} else {
							return func() *mm.Val {
								if func() *mm.Val {
									if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("not")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString(">")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("n")), mm.NilVal)), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("n")), mm.NilVal)), mm.NilVal)), mm.NilVal))), mm.NilVal)) != mm.False {
										return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("not")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("=")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("n")), mm.NilVal)), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("n")), mm.NilVal)), mm.NilVal)), mm.NilVal))), mm.NilVal))
									} else {
										return mm.False
									}
								}() != mm.False {
									return mm.False
								} else {
									return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("apply")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString(">=")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("n")), mm.NilVal)), mm.NilVal)))
								}
							}()
						}
					}()
				}
			}()
		}))
		mm.Define(env, mm.SymbolFromRawString("list-len"), mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("part")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("fold")), mm.Cons(mm.NewCompiled(2, false, func(a []*mm.Val) *mm.Val {
			env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("_"), mm.Cons(mm.SymbolFromRawString("c"), mm.NilVal)), mm.SliceToList(a))
			env = env
			return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("inc")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.NilVal))
		}), mm.Cons(mm.SysIntToNumber(0), mm.NilVal)))))
		mm.Define(env, mm.SymbolFromRawString("len"), mm.NewCompiled(1, false, func(a []*mm.Val) *mm.Val {
			env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("v"), mm.NilVal), mm.SliceToList(a))
			env = env
			return func() *mm.Val {
				if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("vector?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)) != mm.False {
					return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("vector-len")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal))
				} else {
					return func() *mm.Val {
						if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("struct?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)) != mm.False {
							return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("len")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("struct-names")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("s")), mm.NilVal)), mm.NilVal))
						} else {
							return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("list-len")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal))
						}
					}()
				}
			}()
		}))
		mm.Define(env, mm.SymbolFromRawString("mem"), mm.NewCompiled(2, false, func(a []*mm.Val) *mm.Val {
			env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("f"), mm.Cons(mm.SymbolFromRawString("l"), mm.NilVal)), mm.SliceToList(a))
			env = env
			return func() *mm.Val {
				if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("nil?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("l")), mm.NilVal)) != mm.False {
					return mm.False
				} else {
					return func() *mm.Val {
						if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("f")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("l")), mm.NilVal)), mm.NilVal)) != mm.False {
							return mm.LookupDef(env, mm.SymbolFromRawString("l"))
						} else {
							return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("mem")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("f")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("l")), mm.NilVal)), mm.NilVal)))
						}
					}()
				}
			}()
		}))
		mm.Define(env, mm.SymbolFromRawString("memq"), mm.NewCompiled(2, false, func(a []*mm.Val) *mm.Val {
			env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("v"), mm.Cons(mm.SymbolFromRawString("l"), mm.NilVal)), mm.SliceToList(a))
			env = env
			return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("mem")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("part")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("=")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal))), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("l")), mm.NilVal)))
		}))
		mm.Define(env, mm.SymbolFromRawString("notf"), mm.NewCompiled(1, false, func(a []*mm.Val) *mm.Val {
			env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("f"), mm.NilVal), mm.SliceToList(a))
			env = env
			return mm.NewCompiled(1, false, func(a []*mm.Val) *mm.Val {
				env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("i"), mm.NilVal), mm.SliceToList(a))
				env = env
				return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("not")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("f")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("i")), mm.NilVal)), mm.NilVal))
			})
		}))
		mm.Define(env, mm.SymbolFromRawString("every?"), mm.NewCompiled(1, true, func(a []*mm.Val) *mm.Val {
			env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("f"), mm.SymbolFromRawString("l")), mm.SliceToList(a))
			env = env
			return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("not")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("mem")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("notf")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("f")), mm.NilVal)), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("l")), mm.NilVal))), mm.NilVal))
		}))
		mm.Define(env, mm.SymbolFromRawString("any?"), mm.NewCompiled(0, true, func(a []*mm.Val) *mm.Val {
			env := mm.ExtendEnv(env, mm.SymbolFromRawString("v"), mm.SliceToList(a))
			env = env
			return func() *mm.Val {
				if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("mem")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("id")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal))) != mm.False {
					return mm.True
				} else {
					return mm.False
				}
			}()
		}))
		mm.Define(env, mm.SymbolFromRawString("take"), mm.NewCompiled(2, false, func(a []*mm.Val) *mm.Val {
			env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("n"), mm.Cons(mm.SymbolFromRawString("l"), mm.NilVal)), mm.SliceToList(a))
			env = env
			return func() *mm.Val {
				if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("=")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("n")), mm.Cons(mm.SysIntToNumber(0), mm.NilVal))) != mm.False {
					return mm.LookupDef(env, mm.SymbolFromRawString("nil"))
				} else {
					return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cons")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("l")), mm.NilVal)), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("take")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("-")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("n")), mm.Cons(mm.SysIntToNumber(1), mm.NilVal))), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("l")), mm.NilVal)), mm.NilVal))), mm.NilVal)))
				}
			}()
		}))
		mm.Define(env, mm.SymbolFromRawString("drop"), mm.NewCompiled(2, false, func(a []*mm.Val) *mm.Val {
			env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("n"), mm.Cons(mm.SymbolFromRawString("l"), mm.NilVal)), mm.SliceToList(a))
			env = env
			return func() *mm.Val {
				if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("=")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("n")), mm.Cons(mm.SysIntToNumber(0), mm.NilVal))) != mm.False {
					return mm.LookupDef(env, mm.SymbolFromRawString("l"))
				} else {
					return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("drop")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("dec")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("n")), mm.NilVal)), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("l")), mm.NilVal)), mm.NilVal)))
				}
			}()
		}))
		mm.Define(env, mm.SymbolFromRawString("pad"), mm.NewCompiled(2, false, func(a []*mm.Val) *mm.Val {
			env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("f"), mm.Cons(mm.SymbolFromRawString("n"), mm.NilVal)), mm.SliceToList(a))
			env = env
			return mm.NewCompiled(0, true, func(a []*mm.Val) *mm.Val {
				env := mm.ExtendEnv(env, mm.SymbolFromRawString("a"), mm.SliceToList(a))
				env = env
				return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("apply")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("f")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("drop")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("n")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("a")), mm.NilVal))), mm.NilVal)))
			})
		}))
		mm.Define(env, mm.SymbolFromRawString("padr"), mm.NewCompiled(2, false, func(a []*mm.Val) *mm.Val {
			env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("f"), mm.Cons(mm.SymbolFromRawString("n"), mm.NilVal)), mm.SliceToList(a))
			env = env
			return mm.NewCompiled(0, true, func(a []*mm.Val) *mm.Val {
				env := mm.ExtendEnv(env, mm.SymbolFromRawString("a"), mm.SliceToList(a))
				env = env
				return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("apply")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("f")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("take")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("-")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("len")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("a")), mm.NilVal)), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("n")), mm.NilVal))), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("a")), mm.NilVal))), mm.NilVal)))
			})
		}))
		mm.Define(env, mm.SymbolFromRawString("flip"), mm.NewCompiled(1, false, func(a []*mm.Val) *mm.Val {
			env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("f"), mm.NilVal), mm.SliceToList(a))
			env = env
			return mm.NewCompiled(0, true, func(a []*mm.Val) *mm.Val {
				env := mm.ExtendEnv(env, mm.SymbolFromRawString("a"), mm.SliceToList(a))
				env = env
				return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("apply")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("f")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("reverse")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("a")), mm.NilVal)), mm.NilVal)))
			})
		}))
		mm.Define(env, mm.SymbolFromRawString("check-types"), mm.NewCompiled(1, true, func(a []*mm.Val) *mm.Val {
			env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("v"), mm.SymbolFromRawString("types")), mm.SliceToList(a))
			env = env
			return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("apply")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("any?")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("map")), mm.Cons(mm.NewCompiled(1, false, func(a []*mm.Val) *mm.Val {
				env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("t?"), mm.NilVal), mm.SliceToList(a))
				env = env
				return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("t?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal))
			}), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("types")), mm.NilVal))), mm.NilVal)))
		}))
		mm.Define(env, mm.SymbolFromRawString("compose"), mm.NewCompiled(0, true, func(a []*mm.Val) *mm.Val {
			env := mm.ExtendEnv(env, mm.SymbolFromRawString("f"), mm.SliceToList(a))
			env = env
			return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("partr")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("part")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("fold")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("call")), mm.NilVal))), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("f")), mm.NilVal)))
		}))
		mm.Export(env, mm.ListToStruct(mm.Cons(mm.SymbolFromRawString("trace"), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("trace")), mm.Cons(mm.SymbolFromRawString("id"), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("id")), mm.Cons(mm.SymbolFromRawString("list"), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("list")), mm.Cons(mm.SymbolFromRawString("apply"), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("apply")), mm.Cons(mm.SymbolFromRawString("call"), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("call")), mm.Cons(mm.SymbolFromRawString("fold"), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("fold")), mm.Cons(mm.SymbolFromRawString("foldr"), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("foldr")), mm.Cons(mm.SymbolFromRawString("map"), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("map")), mm.Cons(mm.SymbolFromRawString("append"), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("append")), mm.Cons(mm.SymbolFromRawString("part"), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("part")), mm.Cons(mm.SymbolFromRawString("partr"), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("partr")), mm.Cons(mm.SymbolFromRawString("reverse"), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("reverse")), mm.Cons(mm.SymbolFromRawString("reverse-irregular"), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("reverse-irregular")), mm.Cons(mm.SymbolFromRawString("inc"), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("inc")), mm.Cons(mm.SymbolFromRawString("dec"), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("dec")), mm.Cons(mm.SymbolFromRawString(">="), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString(">=")), mm.Cons(mm.SymbolFromRawString("len"), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("len")), mm.Cons(mm.SymbolFromRawString("mem"), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("mem")), mm.Cons(mm.SymbolFromRawString("memq"), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("memq")), mm.Cons(mm.SymbolFromRawString("notf"), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("notf")), mm.Cons(mm.SymbolFromRawString("every?"), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("every?")), mm.Cons(mm.SymbolFromRawString("any?"), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("any?")), mm.Cons(mm.SymbolFromRawString("take"), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("take")), mm.Cons(mm.SymbolFromRawString("drop"), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("drop")), mm.Cons(mm.SymbolFromRawString("pad"), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("pad")), mm.Cons(mm.SymbolFromRawString("padr"), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("padr")), mm.Cons(mm.SymbolFromRawString("flip"), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("flip")), mm.Cons(mm.SymbolFromRawString("check-types"), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("check-types")), mm.Cons(mm.SymbolFromRawString("compose"), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compose")), mm.NilVal))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))
		mm.StoreModule(eenv, mm.SysStringToString("scm/lang.scm"), mm.Exports(env))
	})
	func() {
		m := mm.LoadedModule(env, mm.SysStringToString("scm/lang.scm"))
		if m == mm.UndefinedModule {
			mm.LoadCompiledModule(initialEnv, mm.SysStringToString("scm/lang.scm"))
			m = mm.LoadedModule(env, mm.SysStringToString("scm/lang.scm"))
		}
		mm.DefineAll(env, m)
	}()
	mm.Define(env, mm.SymbolFromRawString("definition-expression"), mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("string->error")), mm.Cons(mm.SysStringToString("definition in expression position"), mm.NilVal)))
	mm.Define(env, mm.SymbolFromRawString("invalid-expression"), mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("string->error")), mm.Cons(mm.SysStringToString("invalid expression"), mm.NilVal)))
	mm.Define(env, mm.SymbolFromRawString("inalid-token"), mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("string->error")), mm.Cons(mm.SysStringToString("invalid-token"), mm.NilVal)))
	mm.Define(env, mm.SymbolFromRawString("circular-import"), mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("string->error")), mm.Cons(mm.SysStringToString("circular-import"), mm.NilVal)))
	mm.Define(env, mm.SymbolFromRawString("not-implemented"), mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("string->error")), mm.Cons(mm.SysStringToString("not implemented"), mm.NilVal)))
	mm.Define(env, mm.SymbolFromRawString("invalid-literal"), mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("string->error")), mm.Cons(mm.SysStringToString("invalid literal"), mm.NilVal)))
	mm.Define(env, mm.SymbolFromRawString("invalid-application"), mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("string->error")), mm.Cons(mm.SysStringToString("invalid application"), mm.NilVal)))
	mm.Define(env, mm.SymbolFromRawString("invalid-value-list"), mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("string->error")), mm.Cons(mm.SysStringToString("invalid value list"), mm.NilVal)))
	mm.Define(env, mm.SymbolFromRawString("invalid-import"), mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("string->error")), mm.Cons(mm.SysStringToString("invalid import expression"), mm.NilVal)))
	mm.Define(env, mm.SymbolFromRawString("compiler"), mm.NewCompiled(0, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.NilVal, mm.SliceToList(a))
		env = env
		return mm.ListToStruct(mm.Cons(mm.SymbolFromRawString("output"), mm.Cons(mm.SysStringToString(""), mm.Cons(mm.SymbolFromRawString("error"), mm.Cons(mm.False, mm.Cons(mm.SymbolFromRawString("compiled-modules"), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("nil")), mm.Cons(mm.SymbolFromRawString("current-import-path"), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("nil")), mm.NilVal)))))))))
	}))
	mm.Define(env, mm.SymbolFromRawString("compiler-append"), mm.NewCompiled(1, true, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("c"), mm.SymbolFromRawString("a")), mm.SliceToList(a))
		env = env
		return func() *mm.Val {
			if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("nil?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("a")), mm.NilVal)) != mm.False {
				return mm.LookupDef(env, mm.SymbolFromRawString("c"))
			} else {
				return func() *mm.Val {
					if mm.LookupDef(env, mm.SymbolFromRawString("c:error")) != mm.False {
						return mm.LookupDef(env, mm.SymbolFromRawString("c"))
					} else {
						return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("apply")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compiler-append")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cons")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("assign")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.ListToStruct(mm.Cons(mm.SymbolFromRawString("output"), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("string-append")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c:output")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("a")), mm.NilVal)), mm.NilVal))), mm.NilVal))), mm.NilVal))), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("a")), mm.NilVal)), mm.NilVal))), mm.NilVal)))
					}
				}()
			}
		}()
	}))
	mm.Define(env, mm.SymbolFromRawString("compiler-error"), mm.NewCompiled(2, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("c"), mm.Cons(mm.SymbolFromRawString("e"), mm.NilVal)), mm.SliceToList(a))
		env = env
		return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("assign")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.ListToStruct(mm.Cons(mm.SymbolFromRawString("error"), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("e")), mm.NilVal))), mm.NilVal)))
	}))
	mm.Define(env, mm.SymbolFromRawString("compiler-compose"), mm.NewCompiled(1, true, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("c"), mm.SymbolFromRawString("i")), mm.SliceToList(a))
		env = env
		return mm.ApplySys(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("apply")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compose")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("foldr")), mm.Cons(mm.NewCompiled(2, false, func(a []*mm.Val) *mm.Val {
			env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("i"), mm.Cons(mm.SymbolFromRawString("p"), mm.NilVal)), mm.SliceToList(a))
			env = env
			return func() *mm.Val {
				if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("function?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("i")), mm.NilVal)) != mm.False {
					return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cons")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("partr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("i")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("p")), mm.NilVal)), mm.NilVal))), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("p")), mm.NilVal)), mm.NilVal)))
				} else {
					return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cons")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("i")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("p")), mm.NilVal)))
				}
			}()
		}), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("nil")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("i")), mm.NilVal)))), mm.NilVal))), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.NilVal))
	}))
	mm.Define(env, mm.SymbolFromRawString("compile-number"), mm.NewCompiled(2, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("c"), mm.Cons(mm.SymbolFromRawString("v"), mm.NilVal)), mm.SliceToList(a))
		env = env
		return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compiler-append")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.SysStringToString("mm.SysIntToNumber("), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("number->string")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.Cons(mm.SysStringToString(")"), mm.NilVal)))))
	}))
	mm.Define(env, mm.SymbolFromRawString("compile-string"), mm.NewCompiled(2, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("c"), mm.Cons(mm.SymbolFromRawString("v"), mm.NilVal)), mm.SliceToList(a))
		env = env
		return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compiler-append")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.SysStringToString("mm.SysStringToString("), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("escape-compiled-string")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.Cons(mm.SysStringToString(")"), mm.NilVal)))))
	}))
	mm.Define(env, mm.SymbolFromRawString("compile-bool"), mm.NewCompiled(2, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("c"), mm.Cons(mm.SymbolFromRawString("v"), mm.NilVal)), mm.SliceToList(a))
		env = env
		return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compiler-append")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(func() *mm.Val {
			if mm.LookupDef(env, mm.SymbolFromRawString("v")) != mm.False {
				return mm.SysStringToString("mm.True")
			} else {
				return mm.SysStringToString("mm.False")
			}
		}(), mm.NilVal)))
	}))
	mm.Define(env, mm.SymbolFromRawString("compile-nil"), mm.NewCompiled(2, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("c"), mm.Cons(mm.SymbolFromRawString("v"), mm.NilVal)), mm.SliceToList(a))
		env = env
		return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compiler-append")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.SysStringToString("mm.NilVal"), mm.NilVal)))
	}))
	mm.Define(env, mm.SymbolFromRawString("compile-symbol"), mm.NewCompiled(2, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("c"), mm.Cons(mm.SymbolFromRawString("v"), mm.NilVal)), mm.SliceToList(a))
		env = env
		return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compiler-append")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.SysStringToString("mm.SymbolFromRawString("), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("escape-compiled-string")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("symbol->string")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.NilVal)), mm.Cons(mm.SysStringToString(")"), mm.NilVal)))))
	}))
	mm.Define(env, mm.SymbolFromRawString("compile-quote-literal"), mm.NewCompiled(2, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("c"), mm.Cons(mm.SymbolFromRawString("v"), mm.NilVal)), mm.SliceToList(a))
		env = env
		return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compiler-compose")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compiler-append")), mm.Cons(mm.SysStringToString("mm.Cons("), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compile-symbol")), mm.Cons(mm.SymbolFromRawString("quote"), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compiler-append")), mm.Cons(mm.SysStringToString(", mm.Cons("), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compile-literal")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.NilVal)), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compiler-append")), mm.Cons(mm.SysStringToString(", mm.NilVal))"), mm.NilVal))))))))))))
	}))
	mm.Define(env, mm.SymbolFromRawString("compile-quote"), mm.NewCompiled(2, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("c"), mm.Cons(mm.SymbolFromRawString("v"), mm.NilVal)), mm.SliceToList(a))
		env = env
		return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compile-literal")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.NilVal)), mm.NilVal)))
	}))
	mm.Define(env, mm.SymbolFromRawString("compile-pair-literal"), mm.NewCompiled(2, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("c"), mm.Cons(mm.SymbolFromRawString("v"), mm.NilVal)), mm.SliceToList(a))
		env = env
		return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compiler-compose")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compiler-append")), mm.Cons(mm.SysStringToString("mm.Cons("), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compile-literal")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compiler-append")), mm.Cons(mm.SysStringToString(", "), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compile-literal")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compiler-append")), mm.Cons(mm.SysStringToString(")"), mm.NilVal))))))))))))
	}))
	mm.Define(env, mm.SymbolFromRawString("make-fn"), mm.NewCompiled(2, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("args"), mm.Cons(mm.SymbolFromRawString("body"), mm.NilVal)), mm.SliceToList(a))
		env = env
		return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cons")), mm.Cons(mm.SymbolFromRawString("fn"), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cons")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("args")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("body")), mm.NilVal))), mm.NilVal)))
	}))
	mm.Define(env, mm.SymbolFromRawString("value-def?"), mm.NewCompiled(1, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("v"), mm.NilVal), mm.SliceToList(a))
		env = env
		return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("symbol?")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.NilVal)), mm.NilVal))
	}))
	mm.Define(env, mm.SymbolFromRawString("valid-def?"), mm.NewCompiled(1, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("v"), mm.NilVal), mm.SliceToList(a))
		env = env
		return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("every?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("pair?")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("take")), mm.Cons(mm.SysIntToNumber(3), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal))), mm.NilVal)))
	}))
	mm.Define(env, mm.SymbolFromRawString("valid-value-def?"), mm.NewCompiled(1, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("v"), mm.NilVal), mm.SliceToList(a))
		env = env
		return func() *mm.Val {
			if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("valid-def?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)) != mm.False {
				return func() *mm.Val {
					if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("symbol?")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.NilVal)), mm.NilVal)) != mm.False {
						return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("nil?")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("drop")), mm.Cons(mm.SysIntToNumber(3), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal))), mm.NilVal))
					} else {
						return mm.False
					}
				}()
			} else {
				return mm.False
			}
		}()
	}))
	mm.Define(env, mm.SymbolFromRawString("valid-function-def?"), mm.NewCompiled(1, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("v"), mm.NilVal), mm.SliceToList(a))
		env = env
		return func() *mm.Val {
			if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("valid-def?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)) != mm.False {
				return func() *mm.Val {
					if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("pair?")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.NilVal)), mm.NilVal)) != mm.False {
						return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("symbol?")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.NilVal)), mm.NilVal)), mm.NilVal))
					} else {
						return mm.False
					}
				}()
			} else {
				return mm.False
			}
		}()
	}))
	mm.Define(env, mm.SymbolFromRawString("def-name"), mm.NewCompiled(1, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("v"), mm.NilVal), mm.SliceToList(a))
		env = env
		return func() *mm.Val {
			if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("value-def?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)) != mm.False {
				return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.NilVal))
			} else {
				return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.NilVal)), mm.NilVal))
			}
		}()
	}))
	mm.Define(env, mm.SymbolFromRawString("def-value"), mm.NewCompiled(1, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("v"), mm.NilVal), mm.SliceToList(a))
		env = env
		return func() *mm.Val {
			if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("value-def?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)) != mm.False {
				return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.NilVal)), mm.NilVal))
			} else {
				return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("make-fn")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.NilVal)), mm.NilVal)), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.NilVal)), mm.NilVal)))
			}
		}()
	}))
	mm.Define(env, mm.SymbolFromRawString("compile-def"), mm.NewCompiled(2, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("c"), mm.Cons(mm.SymbolFromRawString("v"), mm.NilVal)), mm.SliceToList(a))
		env = env
		return func() *mm.Val {
			if func() *mm.Val {
				if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("value-def?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)) != mm.False {
					return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("not")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("valid-value-def?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.NilVal))
				} else {
					return mm.False
				}
			}() != mm.False {
				return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compiler-error")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("invalid-def")), mm.NilVal)))
			} else {
				return func() *mm.Val {
					if func() *mm.Val {
						if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("not")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("value-def?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.NilVal)) != mm.False {
							return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("not")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("valid-function-def?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.NilVal))
						} else {
							return mm.False
						}
					}() != mm.False {
						return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compiler-error")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("invalid-def")), mm.NilVal)))
					} else {
						return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compiler-compose")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compiler-append")), mm.Cons(mm.SysStringToString("mm.Define(env, "), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compile-literal")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("def-name")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compiler-append")), mm.Cons(mm.SysStringToString(", "), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compile-exp")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("def-value")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compiler-append")), mm.Cons(mm.SysStringToString(")"), mm.NilVal))))))))))))
					}
				}()
			}
		}()
	}))
	mm.Define(env, mm.SymbolFromRawString("compile-seq"), mm.NewCompiled(2, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("c"), mm.Cons(mm.SymbolFromRawString("v"), mm.NilVal)), mm.SliceToList(a))
		env = env
		return func() *mm.Val {
			if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("not")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("pair?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.NilVal)) != mm.False {
				return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compiler-error")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("invalid-seq")), mm.NilVal)))
			} else {
				return func() *mm.Val {
					if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("nil?")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.NilVal)) != mm.False {
						return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compiler-compose")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compiler-append")), mm.Cons(mm.SysStringToString("return "), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compile-exp")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.NilVal))))))
					} else {
						return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compiler-compose")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compile-statement")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compiler-append")), mm.Cons(mm.SysStringToString(";\n"), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compile-seq")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.NilVal))))))))
					}
				}()
			}
		}()
	}))
	mm.Define(env, mm.SymbolFromRawString("fn-signature"), mm.NewCompiled(1, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("v"), mm.NilVal), mm.SliceToList(a))
		env = env
		return func() *mm.Val {
			if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("nil?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)) != mm.False {
				return mm.ListToStruct(mm.Cons(mm.SymbolFromRawString("count"), mm.Cons(mm.SysIntToNumber(0), mm.Cons(mm.SymbolFromRawString("var?"), mm.Cons(mm.False, mm.Cons(mm.SymbolFromRawString("names"), mm.Cons(mm.NilVal, mm.NilVal)))))))
			} else {
				return func() *mm.Val {
					if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("symbol?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)) != mm.False {
						return mm.ListToStruct(mm.Cons(mm.SymbolFromRawString("count"), mm.Cons(mm.SysIntToNumber(0), mm.Cons(mm.SymbolFromRawString("var?"), mm.Cons(mm.True, mm.Cons(mm.SymbolFromRawString("names"), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)))))))
					} else {
						return func() *mm.Val {
							if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("pair?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)) != mm.False {
								return mm.ApplySys(mm.NewCompiled(0, false, func(a []*mm.Val) *mm.Val {
									env := mm.ExtendEnv(env, mm.NilVal, mm.SliceToList(a))
									env = env
									mm.Define(env, mm.SymbolFromRawString("signature"), mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("fn-signature")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.NilVal)))
									return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("assign")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("signature")), mm.Cons(mm.ListToStruct(mm.Cons(mm.SymbolFromRawString("count"), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("inc")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("signature:count")), mm.NilVal)), mm.Cons(mm.SymbolFromRawString("names"), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cons")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("signature:names")), mm.NilVal))), mm.NilVal))))), mm.NilVal)))
								}), mm.NilVal)
							} else {
								return mm.LookupDef(env, mm.SymbolFromRawString("invalid-fn"))
							}
						}()
					}
				}()
			}
		}()
	}))
	mm.Define(env, mm.SymbolFromRawString("compile-fn"), mm.NewCompiled(2, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("c"), mm.Cons(mm.SymbolFromRawString("v"), mm.NilVal)), mm.SliceToList(a))
		env = env
		return func() *mm.Val {
			if func() *mm.Val {
				if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("not")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("pair?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.NilVal)) != mm.False {
					return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("not")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("pair?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.NilVal))
				} else {
					return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("not")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("pair?")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.NilVal)), mm.NilVal))
				}
			}() != mm.False {
				return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compiler-error")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("invalid-fn")), mm.NilVal)))
			} else {
				return mm.ApplySys(mm.NewCompiled(0, false, func(a []*mm.Val) *mm.Val {
					env := mm.ExtendEnv(env, mm.NilVal, mm.SliceToList(a))
					env = env
					mm.Define(env, mm.SymbolFromRawString("signature"), mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("fn-signature")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.NilVal)), mm.NilVal)))
					mm.Define(env, mm.SymbolFromRawString("body"), mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.NilVal)))
					return func() *mm.Val {
						if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("error?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("signature")), mm.NilVal)) != mm.False {
							return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compiler-error")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("signature")), mm.NilVal)))
						} else {
							return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compiler-compose")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compiler-append")), mm.Cons(mm.SysStringToString("mm.NewCompiled("), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compiler-append")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("number->string")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("signature:count")), mm.NilVal)), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compiler-append")), mm.Cons(mm.SysStringToString(", "), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compiler-append")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("bool->string")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("signature:var?")), mm.NilVal)), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compiler-append")), mm.Cons(mm.SysStringToString(", func(a []*mm.Val) *mm.Val { env := mm.ExtendEnv(env, "), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compile-literal")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("signature:names")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compiler-append")), mm.Cons(mm.SysStringToString(", mm.SliceToList(a)); env = env; "), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compile-seq")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("body")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compiler-append")), mm.Cons(mm.SysStringToString("})"), mm.NilVal))))))))))))))))))))
						}
					}()
				}), mm.NilVal)
			}
		}()
	}))
	mm.Define(env, mm.SymbolFromRawString("compile-if"), mm.NewCompiled(2, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("c"), mm.Cons(mm.SymbolFromRawString("v"), mm.NilVal)), mm.SliceToList(a))
		env = env
		return func() *mm.Val {
			if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("not")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("=")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("len")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.Cons(mm.SysIntToNumber(4), mm.NilVal))), mm.NilVal)) != mm.False {
				return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compiler-error")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("invalid-if")), mm.NilVal)))
			} else {
				return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compiler-compose")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compiler-append")), mm.Cons(mm.SysStringToString(" func() *mm.Val { if "), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compile-exp")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.NilVal)), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compiler-append")), mm.Cons(mm.SysStringToString(" != mm.False { return "), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compile-exp")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.NilVal)), mm.NilVal)), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compiler-append")), mm.Cons(mm.SysStringToString(" } else { return "), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compile-exp")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.NilVal)), mm.NilVal)), mm.NilVal)), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compiler-append")), mm.Cons(mm.SysStringToString(" }}() "), mm.NilVal))))))))))))))))
			}
		}()
	}))
	mm.Define(env, mm.SymbolFromRawString("and->if"), mm.NewCompiled(1, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("v"), mm.NilVal), mm.SliceToList(a))
		env = env
		return func() *mm.Val {
			if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("nil?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)) != mm.False {
				return mm.True
			} else {
				return func() *mm.Val {
					if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("nil?")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.NilVal)) != mm.False {
						return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal))
					} else {
						return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("list")), mm.Cons(mm.SymbolFromRawString("if"), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("and->if")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.NilVal)), mm.Cons(mm.False, mm.NilVal)))))
					}
				}()
			}
		}()
	}))
	mm.Define(env, mm.SymbolFromRawString("compile-and"), mm.NewCompiled(2, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("c"), mm.Cons(mm.SymbolFromRawString("v"), mm.NilVal)), mm.SliceToList(a))
		env = env
		return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compile")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("and->if")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.NilVal)), mm.NilVal)))
	}))
	mm.Define(env, mm.SymbolFromRawString("or->if"), mm.NewCompiled(1, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("v"), mm.NilVal), mm.SliceToList(a))
		env = env
		return func() *mm.Val {
			if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("nil?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)) != mm.False {
				return mm.False
			} else {
				return func() *mm.Val {
					if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("nil?")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.NilVal)) != mm.False {
						return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal))
					} else {
						return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("list")), mm.Cons(mm.SymbolFromRawString("if"), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("or->if")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.NilVal)), mm.NilVal)))))
					}
				}()
			}
		}()
	}))
	mm.Define(env, mm.SymbolFromRawString("compile-or"), mm.NewCompiled(2, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("c"), mm.Cons(mm.SymbolFromRawString("v"), mm.NilVal)), mm.SliceToList(a))
		env = env
		return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compile")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("or->if")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.NilVal)), mm.NilVal)))
	}))
	mm.Define(env, mm.SymbolFromRawString("compile-begin"), mm.NewCompiled(2, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("c"), mm.Cons(mm.SymbolFromRawString("v"), mm.NilVal)), mm.SliceToList(a))
		env = env
		return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compiler-compose")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compiler-append")), mm.Cons(mm.SysStringToString("func() *mm.Val {"), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compile-seq")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compiler-append")), mm.Cons(mm.SysStringToString("}()"), mm.NilVal))))))))
	}))
	mm.Define(env, mm.SymbolFromRawString("cond->if"), mm.NewCompiled(1, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("v"), mm.NilVal), mm.SliceToList(a))
		env = env
		mm.Define(env, mm.SymbolFromRawString("seq->exp"), mm.NewCompiled(1, false, func(a []*mm.Val) *mm.Val {
			env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("v"), mm.NilVal), mm.SliceToList(a))
			env = env
			return func() *mm.Val {
				if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("nil?")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.NilVal)) != mm.False {
					return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal))
				} else {
					return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cons")), mm.Cons(mm.SymbolFromRawString("begin"), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)))
				}
			}()
		}))
		mm.Define(env, mm.SymbolFromRawString("expand"), mm.NewCompiled(1, false, func(a []*mm.Val) *mm.Val {
			env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("v"), mm.NilVal), mm.SliceToList(a))
			env = env
			return func() *mm.Val {
				if func() *mm.Val {
					if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("not")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("pair?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.NilVal)) != mm.False {
						return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("not")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("pair?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.NilVal))
					} else {
						return func() *mm.Val {
							if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("not")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("pair?")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.NilVal)), mm.NilVal)) != mm.False {
								return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("not")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("pair?")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.NilVal)), mm.NilVal))
							} else {
								return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("not")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("pair?")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.NilVal)), mm.NilVal)), mm.NilVal))
							}
						}()
					}
				}() != mm.False {
					return mm.LookupDef(env, mm.SymbolFromRawString("invalid-cond"))
				} else {
					return func() *mm.Val {
						if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("=")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.NilVal)), mm.Cons(mm.SymbolFromRawString("else"), mm.NilVal))) != mm.False {
							return func() *mm.Val {
								if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("not")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("nil?")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.NilVal)), mm.NilVal)) != mm.False {
									return mm.LookupDef(env, mm.SymbolFromRawString("invalid-cond"))
								} else {
									return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("seq->exp")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.NilVal)), mm.NilVal))
								}
							}()
						} else {
							return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("list")), mm.Cons(mm.SymbolFromRawString("if"), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.NilVal)), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("seq->exp")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.NilVal)), mm.NilVal)), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("expand")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.NilVal)), mm.NilVal)))))
						}
					}()
				}
			}()
		}))
		return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("expand")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.NilVal))
	}))
	mm.Define(env, mm.SymbolFromRawString("compile-cond"), mm.NewCompiled(2, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("c"), mm.Cons(mm.SymbolFromRawString("v"), mm.NilVal)), mm.SliceToList(a))
		env = env
		return mm.ApplySys(mm.NewCompiled(0, false, func(a []*mm.Val) *mm.Val {
			env := mm.ExtendEnv(env, mm.NilVal, mm.SliceToList(a))
			env = env
			mm.Define(env, mm.SymbolFromRawString("vi"), mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cond->if")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)))
			return func() *mm.Val {
				if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("error?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("vi")), mm.NilVal)) != mm.False {
					return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compiler-error")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("vi")), mm.NilVal)))
				} else {
					return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compile")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("vi")), mm.NilVal)))
				}
			}()
		}), mm.NilVal)
	}))
	mm.Define(env, mm.SymbolFromRawString("let-body"), mm.NewCompiled(1, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("v"), mm.NilVal), mm.SliceToList(a))
		env = env
		mm.Define(env, mm.SymbolFromRawString("let-defs"), mm.NewCompiled(1, false, func(a []*mm.Val) *mm.Val {
			env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("v"), mm.NilVal), mm.SliceToList(a))
			env = env
			return func() *mm.Val {
				if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("nil?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)) != mm.False {
					return mm.LookupDef(env, mm.SymbolFromRawString("nil"))
				} else {
					return func() *mm.Val {
						if func() *mm.Val {
							if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("not")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("pair?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.NilVal)) != mm.False {
								return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("not")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("pair?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.NilVal))
							} else {
								return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("not")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("pair?")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.NilVal)), mm.NilVal))
							}
						}() != mm.False {
							return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("fatal")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("invalid-let")), mm.NilVal))
						} else {
							return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cons")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("list")), mm.Cons(mm.SymbolFromRawString("def"), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.NilVal)), mm.NilVal)))), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("let-defs")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.NilVal)), mm.NilVal)), mm.NilVal)))
						}
					}()
				}
			}()
		}))
		return func() *mm.Val {
			if func() *mm.Val {
				if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("not")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("pair?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.NilVal)) != mm.False {
					return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("not")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("pair?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.NilVal))
				} else {
					return func() *mm.Val {
						if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("not")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("pair?")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.NilVal)), mm.NilVal)) != mm.False {
							return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("not")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("pair?")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.NilVal)), mm.NilVal))
						} else {
							return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("nil?")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.NilVal)), mm.NilVal))
						}
					}()
				}
			}() != mm.False {
				return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("fatal")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("invalid-let")), mm.NilVal))
			} else {
				return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("append")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("let-defs")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.NilVal)), mm.NilVal)), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.NilVal)), mm.NilVal)))
			}
		}()
	}))
	mm.Define(env, mm.SymbolFromRawString("compile-let"), mm.NewCompiled(2, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("c"), mm.Cons(mm.SymbolFromRawString("v"), mm.NilVal)), mm.SliceToList(a))
		env = env
		return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compile")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("list")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("make-fn")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("nil")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("let-body")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.NilVal))), mm.NilVal)), mm.NilVal)))
	}))
	mm.Define(env, mm.SymbolFromRawString("compile-test"), mm.NewCompiled(2, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("c"), mm.Cons(mm.SymbolFromRawString("v"), mm.NilVal)), mm.SliceToList(a))
		env = env
		mm.Define(env, mm.SymbolFromRawString("compile-test-seq"), mm.NewCompiled(2, false, func(a []*mm.Val) *mm.Val {
			env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("c"), mm.Cons(mm.SymbolFromRawString("v"), mm.NilVal)), mm.SliceToList(a))
			env = env
			return func() *mm.Val {
				if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("nil?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)) != mm.False {
					return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compiler-compose")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compiler-append")), mm.Cons(mm.SysStringToString("return "), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compile-literal")), mm.Cons(mm.SymbolFromRawString("test-complete"), mm.NilVal))))))
				} else {
					return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compiler-compose")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compiler-append")), mm.Cons(mm.SysStringToString("if result := func() *mm.Val { return "), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compile")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compiler-append")), mm.Cons(mm.SysStringToString("}(); result == mm.False { return mm.Fatal("), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compile-exp")), mm.Cons(mm.SysStringToString("test failed"), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compiler-append")), mm.Cons(mm.SysStringToString(") } else if mm.IsError(result) != mm.False "), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compiler-append")), mm.Cons(mm.SysStringToString(" { return mm.Fatal(result) }; "), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compile-test-seq")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.NilVal))))))))))))))))
				}
			}()
		}))
		return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compiler-compose")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compiler-append")), mm.Cons(mm.SysStringToString("func() *mm.Val {"), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compiler-append")), mm.Cons(mm.SysStringToString("env := mm.ExtendEnv(env, mm.NilVal, mm.NilVal); env = env;"), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compile-test-seq")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compiler-append")), mm.Cons(mm.SysStringToString("}()"), mm.NilVal))))))))))
	}))
	mm.Define(env, mm.SymbolFromRawString("compile-export"), mm.NewCompiled(2, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("c"), mm.Cons(mm.SymbolFromRawString("v"), mm.NilVal)), mm.SliceToList(a))
		env = env
		return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compiler-compose")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compiler-append")), mm.Cons(mm.SysStringToString("mm.Export(env, "), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compile")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cons")), mm.Cons(mm.SymbolFromRawString("struct:"), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.NilVal))), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compiler-append")), mm.Cons(mm.SysStringToString(")"), mm.NilVal))))))))
	}))
	mm.Define(env, mm.SymbolFromRawString("read-compile"), mm.NewCompiled(2, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("c"), mm.Cons(mm.SymbolFromRawString("r"), mm.NilVal)), mm.SliceToList(a))
		env = env
		return mm.ApplySys(mm.NewCompiled(0, false, func(a []*mm.Val) *mm.Val {
			env := mm.ExtendEnv(env, mm.NilVal, mm.SliceToList(a))
			env = env
			mm.Define(env, mm.SymbolFromRawString("r"), mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("read")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("r")), mm.NilVal)))
			return func() *mm.Val {
				if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("=")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("r:state")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("eof")), mm.NilVal))) != mm.False {
					return mm.LookupDef(env, mm.SymbolFromRawString("c"))
				} else {
					return func() *mm.Val {
						if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("error?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("r:state")), mm.NilVal)) != mm.False {
							return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compiler-error")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("r:state")), mm.NilVal)))
						} else {
							return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("read-compile")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compiler-compose")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compile")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("r:state")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compiler-append")), mm.Cons(mm.SysStringToString("\n"), mm.NilVal)))))), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("r")), mm.NilVal)))
						}
					}()
				}
			}()
		}), mm.NilVal)
	}))
	mm.Define(env, mm.SymbolFromRawString("compile-module"), mm.NewCompiled(2, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("c"), mm.Cons(mm.SymbolFromRawString("module-name"), mm.NilVal)), mm.SliceToList(a))
		env = env
		return mm.ApplySys(mm.NewCompiled(0, false, func(a []*mm.Val) *mm.Val {
			env := mm.ExtendEnv(env, mm.NilVal, mm.SliceToList(a))
			env = env
			mm.Define(env, mm.SymbolFromRawString("f"), mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("fopen")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("module-name")), mm.NilVal)))
			return func() *mm.Val {
				if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("error?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("f")), mm.NilVal)) != mm.False {
					return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compiler-error")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("f")), mm.NilVal)))
				} else {
					return mm.ApplySys(mm.NewCompiled(0, false, func(a []*mm.Val) *mm.Val {
						env := mm.ExtendEnv(env, mm.NilVal, mm.SliceToList(a))
						env = env
						mm.Define(env, mm.SymbolFromRawString("c"), mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compiler-compose")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compiler-append")), mm.Cons(mm.SysStringToString("mm.ModuleLoader(env, "), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compile")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("module-name")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compiler-append")), mm.Cons(mm.SysStringToString(", func(eenv *mm.Val) { env := mm.ModuleEnv(eenv, "), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compile")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("module-name")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compiler-append")), mm.Cons(mm.SysStringToString(");"), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("read-compile")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("reader")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("f")), mm.NilVal)), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compiler-append")), mm.Cons(mm.SysStringToString("; mm.StoreModule(eenv, "), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compile")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("module-name")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compiler-append")), mm.Cons(mm.SysStringToString(", mm.Exports(env)) });"), mm.NilVal)))))))))))))))))))))
						return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("assign")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.ListToStruct(mm.Cons(mm.SymbolFromRawString("compiled-modules"), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cons")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("module-name")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c:compiled-modules")), mm.NilVal))), mm.NilVal))), mm.NilVal)))
					}), mm.NilVal)
				}
			}()
		}), mm.NilVal)
	}))
	mm.Define(env, mm.SymbolFromRawString("compile-import"), mm.NewCompiled(2, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("c"), mm.Cons(mm.SymbolFromRawString("exp"), mm.NilVal)), mm.SliceToList(a))
		env = env
		mm.Define(env, mm.SymbolFromRawString("compile-module-define"), mm.NewCompiled(3, false, func(a []*mm.Val) *mm.Val {
			env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("c"), mm.Cons(mm.SymbolFromRawString("module-name"), mm.Cons(mm.SymbolFromRawString("import-name"), mm.NilVal))), mm.SliceToList(a))
			env = env
			return func() *mm.Val {
				if mm.LookupDef(env, mm.SymbolFromRawString("import-name")) != mm.False {
					return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compiler-compose")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compiler-append")), mm.Cons(mm.SysStringToString("mm.Define(env, "), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compile-literal")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("import-name")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compiler-append")), mm.Cons(mm.SysStringToString(", m)"), mm.NilVal))))))))
				} else {
					return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compiler-append")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.SysStringToString("mm.DefineAll(env, m)"), mm.NilVal)))
				}
			}()
		}))
		mm.Define(env, mm.SymbolFromRawString("compile-load"), mm.NewCompiled(3, false, func(a []*mm.Val) *mm.Val {
			env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("c"), mm.Cons(mm.SymbolFromRawString("module-name"), mm.Cons(mm.SymbolFromRawString("import-name"), mm.NilVal))), mm.SliceToList(a))
			env = env
			return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compiler-compose")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compiler-append")), mm.Cons(mm.SysStringToString("func() { m := mm.LoadedModule(env, "), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compile")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("module-name")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compiler-append")), mm.Cons(mm.SysStringToString(");"), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compiler-append")), mm.Cons(mm.SysStringToString("if m == mm.UndefinedModule {"), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compiler-append")), mm.Cons(mm.SysStringToString("mm.LoadCompiledModule(initialEnv, "), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compile")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("module-name")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compiler-append")), mm.Cons(mm.SysStringToString("); m = mm.LoadedModule(env, "), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compile")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("module-name")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compiler-append")), mm.Cons(mm.SysStringToString(")};"), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("partr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compile-module-define")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("import-name")), mm.NilVal))), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("module-name")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compiler-append")), mm.Cons(mm.SysStringToString("}();"), mm.NilVal))))))))))))))))))))))))
		}))
		return mm.ApplySys(mm.NewCompiled(0, false, func(a []*mm.Val) *mm.Val {
			env := mm.ExtendEnv(env, mm.NilVal, mm.SliceToList(a))
			env = env
			mm.Define(env, mm.SymbolFromRawString("i"), mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("import-def")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("exp")), mm.NilVal)))
			return func() *mm.Val {
				if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("error?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("i")), mm.NilVal)) != mm.False {
					return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compiler-error")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("i")), mm.NilVal)))
				} else {
					return func() *mm.Val {
						if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("memq")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("i:module-name")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c:current-import-path")), mm.NilVal))) != mm.False {
							return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compiler-error")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("circular-import")), mm.NilVal)))
						} else {
							return func() *mm.Val {
								if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("memq")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("i:module-name")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c:compiled-modules")), mm.NilVal))) != mm.False {
									return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compile-load")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("i:module-name")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("i:import-name")), mm.NilVal))))
								} else {
									return mm.ApplySys(mm.NewCompiled(0, false, func(a []*mm.Val) *mm.Val {
										env := mm.ExtendEnv(env, mm.NilVal, mm.SliceToList(a))
										env = env
										mm.Define(env, mm.SymbolFromRawString("cc"), mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("assign")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compiler")), mm.NilVal), mm.Cons(mm.ListToStruct(mm.Cons(mm.SymbolFromRawString("current-import-path"), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cons")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("i:module-name")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c:current-import-path")), mm.NilVal))), mm.Cons(mm.SymbolFromRawString("compiled-modules"), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c:compiled-modules")), mm.NilVal))))), mm.NilVal))))
										mm.Define(env, mm.SymbolFromRawString("cr"), mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compile-module")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("cc")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("i:module-name")), mm.NilVal))))
										mm.Define(env, mm.SymbolFromRawString("c"), mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("assign")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.ListToStruct(mm.Cons(mm.SymbolFromRawString("error"), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("cr:error")), mm.Cons(mm.SymbolFromRawString("compiled-modules"), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cons")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("i:module-name")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c:compiled-modules")), mm.NilVal))), mm.NilVal))))), mm.NilVal))))
										return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compiler-compose")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compiler-append")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("cr:output")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("partr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compile-load")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("i:import-name")), mm.NilVal))), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("i:module-name")), mm.NilVal))))))
									}), mm.NilVal)
								}
							}()
						}
					}()
				}
			}()
		}), mm.NilVal)
	}))
	mm.Define(env, mm.SymbolFromRawString("compile-lookup"), mm.NewCompiled(2, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("c"), mm.Cons(mm.SymbolFromRawString("v"), mm.NilVal)), mm.SliceToList(a))
		env = env
		return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compiler-compose")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compiler-append")), mm.Cons(mm.SysStringToString("mm.LookupDef(env, "), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compile-symbol")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compiler-append")), mm.Cons(mm.SysStringToString(")"), mm.NilVal))))))))
	}))
	mm.Define(env, mm.SymbolFromRawString("compile-vector"), mm.NewCompiled(2, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("c"), mm.Cons(mm.SymbolFromRawString("v"), mm.NilVal)), mm.SliceToList(a))
		env = env
		return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compiler-compose")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compiler-append")), mm.Cons(mm.SysStringToString("mm.ListToVector("), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compile-literal")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compiler-append")), mm.Cons(mm.SysStringToString(")"), mm.NilVal))))))))
	}))
	mm.Define(env, mm.SymbolFromRawString("compile-struct"), mm.NewCompiled(2, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("c"), mm.Cons(mm.SymbolFromRawString("v"), mm.NilVal)), mm.SliceToList(a))
		env = env
		mm.Define(env, mm.SymbolFromRawString("compile-struct-values"), mm.NewCompiled(2, false, func(a []*mm.Val) *mm.Val {
			env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("c"), mm.Cons(mm.SymbolFromRawString("v"), mm.NilVal)), mm.SliceToList(a))
			env = env
			return func() *mm.Val {
				if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("nil?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)) != mm.False {
					return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compiler-append")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.SysStringToString("mm.NilVal"), mm.NilVal)))
				} else {
					return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compiler-compose")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compiler-append")), mm.Cons(mm.SysStringToString("mm.Cons("), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compile-literal")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compiler-append")), mm.Cons(mm.SysStringToString(", mm.Cons("), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compile-exp")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.NilVal)), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compiler-append")), mm.Cons(mm.SysStringToString(", "), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compile-struct-values")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.NilVal)), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compiler-append")), mm.Cons(mm.SysStringToString("))"), mm.NilVal))))))))))))))))
				}
			}()
		}))
		return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compiler-compose")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compiler-append")), mm.Cons(mm.SysStringToString("mm.ListToStruct("), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compile-struct-values")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compiler-append")), mm.Cons(mm.SysStringToString(")"), mm.NilVal))))))))
	}))
	mm.Define(env, mm.SymbolFromRawString("compile-value-list"), mm.NewCompiled(2, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("c"), mm.Cons(mm.SymbolFromRawString("v"), mm.NilVal)), mm.SliceToList(a))
		env = env
		return func() *mm.Val {
			if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("nil?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)) != mm.False {
				return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compiler-append")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.SysStringToString("mm.NilVal"), mm.NilVal)))
			} else {
				return func() *mm.Val {
					if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("not")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("pair?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.NilVal)) != mm.False {
						return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compiler-error")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("invalid-value-list")), mm.NilVal)))
					} else {
						return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compiler-compose")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compiler-append")), mm.Cons(mm.SysStringToString("mm.Cons("), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compile-exp")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compiler-append")), mm.Cons(mm.SysStringToString(", "), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compile-value-list")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compiler-append")), mm.Cons(mm.SysStringToString(")"), mm.NilVal))))))))))))
					}
				}()
			}
		}()
	}))
	mm.Define(env, mm.SymbolFromRawString("compile-application"), mm.NewCompiled(2, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("c"), mm.Cons(mm.SymbolFromRawString("v"), mm.NilVal)), mm.SliceToList(a))
		env = env
		return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compiler-compose")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compiler-append")), mm.Cons(mm.SysStringToString("mm.ApplySys("), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compile-exp")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compiler-append")), mm.Cons(mm.SysStringToString(", "), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compile-value-list")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compiler-append")), mm.Cons(mm.SysStringToString(")"), mm.NilVal))))))))))))
	}))
	mm.Define(env, mm.SymbolFromRawString("compile-current-env"), mm.NewCompiled(1, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("c"), mm.NilVal), mm.SliceToList(a))
		env = env
		return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compiler-append")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.SysStringToString("func() *mm.Val { return env }()"), mm.NilVal)))
	}))
	mm.Define(env, mm.SymbolFromRawString("compile-literal"), mm.NewCompiled(2, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("c"), mm.Cons(mm.SymbolFromRawString("v"), mm.NilVal)), mm.SliceToList(a))
		env = env
		return func() *mm.Val {
			if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("check-types")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("number?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("string?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("bool?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("nil?")), mm.NilVal)))))) != mm.False {
				return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compile")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)))
			} else {
				return func() *mm.Val {
					if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("quote?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)) != mm.False {
						return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compile-quote-literal")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)))
					} else {
						return func() *mm.Val {
							if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("symbol?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)) != mm.False {
								return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compile-symbol")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)))
							} else {
								return func() *mm.Val {
									if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("pair?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)) != mm.False {
										return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compile-pair-literal")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)))
									} else {
										return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compiler-error")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("invalid-literal")), mm.NilVal)))
									}
								}()
							}
						}()
					}
				}()
			}
		}()
	}))
	mm.Define(env, mm.SymbolFromRawString("compile-exp"), mm.NewCompiled(2, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("c"), mm.Cons(mm.SymbolFromRawString("v"), mm.NilVal)), mm.SliceToList(a))
		env = env
		return func() *mm.Val {
			if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("check-types")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("number?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("string?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("bool?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("nil?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("quote?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("symbol?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("vector-form?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("struct-form?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("if?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("and?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("or?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("fn?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("begin?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("cond?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("let?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("current-env?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("application?")), mm.NilVal))))))))))))))))))) != mm.False {
				return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compile")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)))
			} else {
				return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compiler-error")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("invalid-expression")), mm.NilVal)))
			}
		}()
	}))
	mm.Define(env, mm.SymbolFromRawString("compile-statement"), mm.NewCompiled(2, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("c"), mm.Cons(mm.SymbolFromRawString("v"), mm.NilVal)), mm.SliceToList(a))
		env = env
		return func() *mm.Val {
			if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("check-types")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("def?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("if?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("and?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("or?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("begin?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("cond?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("let?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("import?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("export?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("test?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("application?")), mm.NilVal))))))))))))) != mm.False {
				return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compile")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)))
			} else {
				return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compiler-error")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("invalid-statement")), mm.NilVal)))
			}
		}()
	}))
	mm.Define(env, mm.SymbolFromRawString("compile"), mm.NewCompiled(2, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("c"), mm.Cons(mm.SymbolFromRawString("v"), mm.NilVal)), mm.SliceToList(a))
		env = env
		return func() *mm.Val {
			if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("number?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)) != mm.False {
				return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compile-number")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)))
			} else {
				return func() *mm.Val {
					if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("string?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)) != mm.False {
						return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compile-string")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)))
					} else {
						return func() *mm.Val {
							if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("bool?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)) != mm.False {
								return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compile-bool")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)))
							} else {
								return func() *mm.Val {
									if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("nil?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)) != mm.False {
										return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compile-nil")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)))
									} else {
										return func() *mm.Val {
											if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("quote?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)) != mm.False {
												return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compile-quote")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)))
											} else {
												return func() *mm.Val {
													if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("symbol?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)) != mm.False {
														return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compile-lookup")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)))
													} else {
														return func() *mm.Val {
															if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("vector-form?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)) != mm.False {
																return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compile-vector")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)))
															} else {
																return func() *mm.Val {
																	if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("struct-form?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)) != mm.False {
																		return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compile-struct")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)))
																	} else {
																		return func() *mm.Val {
																			if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("def?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)) != mm.False {
																				return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compile-def")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)))
																			} else {
																				return func() *mm.Val {
																					if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("if?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)) != mm.False {
																						return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compile-if")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)))
																					} else {
																						return func() *mm.Val {
																							if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("and?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)) != mm.False {
																								return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compile-and")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)))
																							} else {
																								return func() *mm.Val {
																									if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("or?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)) != mm.False {
																										return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compile-or")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)))
																									} else {
																										return func() *mm.Val {
																											if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("fn?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)) != mm.False {
																												return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compile-fn")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)))
																											} else {
																												return func() *mm.Val {
																													if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("begin?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)) != mm.False {
																														return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compile-begin")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)))
																													} else {
																														return func() *mm.Val {
																															if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cond?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)) != mm.False {
																																return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compile-cond")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)))
																															} else {
																																return func() *mm.Val {
																																	if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("let?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)) != mm.False {
																																		return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compile-let")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)))
																																	} else {
																																		return func() *mm.Val {
																																			if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("export?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)) != mm.False {
																																				return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compile-export")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)))
																																			} else {
																																				return func() *mm.Val {
																																					if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("import?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)) != mm.False {
																																						return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compile-import")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)))
																																					} else {
																																						return func() *mm.Val {
																																							if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("test?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)) != mm.False {
																																								return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compile-test")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)))
																																							} else {
																																								return func() *mm.Val {
																																									if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("current-env?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)) != mm.False {
																																										return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compile-current-env")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.NilVal))
																																									} else {
																																										return func() *mm.Val {
																																											if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("application?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)) != mm.False {
																																												return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compile-application")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)))
																																											} else {
																																												return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compiler-error")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("not-implemented")), mm.NilVal)))
																																											}
																																										}()
																																									}
																																								}()
																																							}
																																						}()
																																					}
																																				}()
																																			}
																																		}()
																																	}
																																}()
																															}
																														}()
																													}
																												}()
																											}
																										}()
																									}
																								}()
																							}
																						}()
																					}
																				}()
																			}
																		}()
																	}
																}()
															}
														}()
													}
												}()
											}
										}()
									}
								}()
							}
						}()
					}
				}()
			}
		}()
	}))
	mm.Define(env, mm.SymbolFromRawString("compiled-head"), mm.SysStringToString("package main\n\n     import mm \"github.com/aryszka/mikkamakka\"\n\n     func main() {\n          initialEnv := mm.InitialEnv()\n         env := initialEnv\n     "))
	mm.Define(env, mm.SymbolFromRawString("compiled-tail"), mm.SysStringToString("}"))
	mm.ModuleLoader(env, mm.SysStringToString("scm/read.scm"), func(eenv *mm.Val) {
		env := mm.ModuleEnv(eenv, mm.SysStringToString("scm/read.scm"))
		func() {
			m := mm.LoadedModule(env, mm.SysStringToString("scm/lang.scm"))
			if m == mm.UndefinedModule {
				mm.LoadCompiledModule(initialEnv, mm.SysStringToString("scm/lang.scm"))
				m = mm.LoadedModule(env, mm.SysStringToString("scm/lang.scm"))
			}
			mm.DefineAll(env, m)
		}()
		mm.Define(env, mm.SymbolFromRawString("irregular-cons"), mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("string->error")), mm.Cons(mm.SysStringToString("irregular cons expression"), mm.NilVal)))
		mm.Define(env, mm.SymbolFromRawString("unexpected-close"), mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("string->error")), mm.Cons(mm.SysStringToString("unexpected close token"), mm.NilVal)))
		mm.Define(env, mm.SymbolFromRawString("invalid-statement"), mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("string->error")), mm.Cons(mm.SysStringToString("invalid expression in statement position"), mm.NilVal)))
		mm.Define(env, mm.SymbolFromRawString("invalid-cond"), mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("string->error")), mm.Cons(mm.SysStringToString("invalid cond expression"), mm.NilVal)))
		mm.Define(env, mm.SymbolFromRawString("token-type"), mm.ListToStruct(mm.Cons(mm.SymbolFromRawString("none"), mm.Cons(mm.SysIntToNumber(0), mm.Cons(mm.SymbolFromRawString("symbol"), mm.Cons(mm.SysIntToNumber(1), mm.Cons(mm.SymbolFromRawString("string"), mm.Cons(mm.SysIntToNumber(2), mm.Cons(mm.SymbolFromRawString("comment"), mm.Cons(mm.SysIntToNumber(3), mm.Cons(mm.SymbolFromRawString("list"), mm.Cons(mm.SysIntToNumber(4), mm.Cons(mm.SymbolFromRawString("quote"), mm.Cons(mm.SysIntToNumber(5), mm.Cons(mm.SymbolFromRawString("vector"), mm.Cons(mm.SysIntToNumber(6), mm.Cons(mm.SymbolFromRawString("struct"), mm.Cons(mm.SysIntToNumber(7), mm.NilVal))))))))))))))))))
		mm.Define(env, mm.SymbolFromRawString("list-type"), mm.ListToStruct(mm.Cons(mm.SymbolFromRawString("none"), mm.Cons(mm.SysIntToNumber(0), mm.Cons(mm.SymbolFromRawString("lisp"), mm.Cons(mm.SysIntToNumber(1), mm.Cons(mm.SymbolFromRawString("vector"), mm.Cons(mm.SysIntToNumber(2), mm.Cons(mm.SymbolFromRawString("struct"), mm.Cons(mm.SysIntToNumber(3), mm.NilVal))))))))))
		mm.Define(env, mm.SymbolFromRawString("undefined"), mm.ListToStruct(mm.NilVal))
		mm.Define(env, mm.SymbolFromRawString("reader"), mm.NewCompiled(1, false, func(a []*mm.Val) *mm.Val {
			env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("input"), mm.NilVal), mm.SliceToList(a))
			env = env
			return mm.ListToStruct(mm.Cons(mm.SymbolFromRawString("input"), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("input")), mm.Cons(mm.SymbolFromRawString("token-type"), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("token-type:none")), mm.Cons(mm.SymbolFromRawString("state"), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("undefined")), mm.Cons(mm.SymbolFromRawString("escaped?"), mm.Cons(mm.False, mm.Cons(mm.SymbolFromRawString("char"), mm.Cons(mm.SysStringToString(""), mm.Cons(mm.SymbolFromRawString("token"), mm.Cons(mm.SysStringToString(""), mm.Cons(mm.SymbolFromRawString("list-type"), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("list-type:none")), mm.Cons(mm.SymbolFromRawString("close-list?"), mm.Cons(mm.False, mm.Cons(mm.SymbolFromRawString("list-items"), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("nil")), mm.Cons(mm.SymbolFromRawString("list-cons?"), mm.Cons(mm.False, mm.Cons(mm.SymbolFromRawString("cons-items"), mm.Cons(mm.SysIntToNumber(0), mm.NilVal)))))))))))))))))))))))
		}))
		mm.Define(env, mm.SymbolFromRawString("read-char"), mm.NewCompiled(1, false, func(a []*mm.Val) *mm.Val {
			env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("r"), mm.NilVal), mm.SliceToList(a))
			env = env
			return mm.ApplySys(mm.NewCompiled(0, false, func(a []*mm.Val) *mm.Val {
				env := mm.ExtendEnv(env, mm.NilVal, mm.SliceToList(a))
				env = env
				mm.Define(env, mm.SymbolFromRawString("next-input"), mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("fread")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("r:input")), mm.Cons(mm.SysIntToNumber(1), mm.NilVal))))
				mm.Define(env, mm.SymbolFromRawString("state"), mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("fstate")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next-input")), mm.NilVal)))
				return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("assign")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("r")), mm.Cons(mm.ListToStruct(mm.Cons(mm.SymbolFromRawString("input"), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next-input")), mm.NilVal))), mm.Cons(func() *mm.Val {
					if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("error?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("state")), mm.NilVal)) != mm.False {
						return mm.ListToStruct(mm.Cons(mm.SymbolFromRawString("state"), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("state")), mm.NilVal)))
					} else {
						return mm.ListToStruct(mm.Cons(mm.SymbolFromRawString("char"), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("state")), mm.NilVal)))
					}
				}(), mm.NilVal))))
			}), mm.NilVal)
		}))
		mm.Define(env, mm.SymbolFromRawString("newline?"), mm.NewCompiled(1, false, func(a []*mm.Val) *mm.Val {
			env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("c"), mm.NilVal), mm.SliceToList(a))
			env = env
			return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("=")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.SysStringToString("\n"), mm.NilVal)))
		}))
		mm.Define(env, mm.SymbolFromRawString("char-check"), mm.NewCompiled(1, true, func(a []*mm.Val) *mm.Val {
			env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("c"), mm.SymbolFromRawString("cc")), mm.SliceToList(a))
			env = env
			return func() *mm.Val {
				if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("nil?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("cc")), mm.NilVal)) != mm.False {
					return mm.False
				} else {
					return func() *mm.Val {
						if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("=")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("cc")), mm.NilVal)), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.NilVal))) != mm.False {
							return mm.True
						} else {
							return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("apply")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("char-check")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cons")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("cc")), mm.NilVal)), mm.NilVal))), mm.NilVal)))
						}
					}()
				}
			}()
		}))
		mm.Define(env, mm.SymbolFromRawString("make-char-check"), mm.NewCompiled(0, true, func(a []*mm.Val) *mm.Val {
			env := mm.ExtendEnv(env, mm.SymbolFromRawString("cc"), mm.SliceToList(a))
			env = env
			return mm.NewCompiled(1, false, func(a []*mm.Val) *mm.Val {
				env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("c"), mm.NilVal), mm.SliceToList(a))
				env = env
				return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("apply")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("char-check")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cons")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("cc")), mm.NilVal))), mm.NilVal)))
			})
		}))
		mm.Define(env, mm.SymbolFromRawString("whitespace?"), mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("make-char-check")), mm.Cons(mm.SysStringToString(" "), mm.Cons(mm.SysStringToString("\b"), mm.Cons(mm.SysStringToString("\f"), mm.Cons(mm.SysStringToString("\n"), mm.Cons(mm.SysStringToString("\r"), mm.Cons(mm.SysStringToString("\t"), mm.Cons(mm.SysStringToString("\v"), mm.NilVal)))))))))
		mm.Define(env, mm.SymbolFromRawString("escape-char?"), mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("make-char-check")), mm.Cons(mm.SysStringToString("\\"), mm.NilVal)))
		mm.Define(env, mm.SymbolFromRawString("string-delimiter?"), mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("make-char-check")), mm.Cons(mm.SysStringToString("\""), mm.NilVal)))
		mm.Define(env, mm.SymbolFromRawString("comment-char?"), mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("make-char-check")), mm.Cons(mm.SysStringToString(";"), mm.NilVal)))
		mm.Define(env, mm.SymbolFromRawString("list-open?"), mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("make-char-check")), mm.Cons(mm.SysStringToString("("), mm.NilVal)))
		mm.Define(env, mm.SymbolFromRawString("list-close?"), mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("make-char-check")), mm.Cons(mm.SysStringToString(")"), mm.NilVal)))
		mm.Define(env, mm.SymbolFromRawString("cons-char?"), mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("make-char-check")), mm.Cons(mm.SysStringToString("."), mm.NilVal)))
		mm.Define(env, mm.SymbolFromRawString("quote-char?"), mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("make-char-check")), mm.Cons(mm.SysStringToString("'"), mm.NilVal)))
		mm.Define(env, mm.SymbolFromRawString("vector-open?"), mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("make-char-check")), mm.Cons(mm.SysStringToString("["), mm.NilVal)))
		mm.Define(env, mm.SymbolFromRawString("vector-close?"), mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("make-char-check")), mm.Cons(mm.SysStringToString("]"), mm.NilVal)))
		mm.Define(env, mm.SymbolFromRawString("struct-open?"), mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("make-char-check")), mm.Cons(mm.SysStringToString("{"), mm.NilVal)))
		mm.Define(env, mm.SymbolFromRawString("struct-close?"), mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("make-char-check")), mm.Cons(mm.SysStringToString("}"), mm.NilVal)))
		mm.Define(env, mm.SymbolFromRawString("set-escaped"), mm.NewCompiled(1, false, func(a []*mm.Val) *mm.Val {
			env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("r"), mm.NilVal), mm.SliceToList(a))
			env = env
			return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("assign")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("r")), mm.Cons(mm.ListToStruct(mm.Cons(mm.SymbolFromRawString("escaped?"), mm.Cons(mm.True, mm.NilVal))), mm.NilVal)))
		}))
		mm.Define(env, mm.SymbolFromRawString("append-token"), mm.NewCompiled(1, false, func(a []*mm.Val) *mm.Val {
			env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("r"), mm.NilVal), mm.SliceToList(a))
			env = env
			return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("assign")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("r")), mm.Cons(mm.ListToStruct(mm.Cons(mm.SymbolFromRawString("token"), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("string-append")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("r:token")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("r:char")), mm.NilVal))), mm.NilVal))), mm.NilVal)))
		}))
		mm.Define(env, mm.SymbolFromRawString("append-token-escaped"), mm.NewCompiled(1, false, func(a []*mm.Val) *mm.Val {
			env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("r"), mm.NilVal), mm.SliceToList(a))
			env = env
			return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("assign")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("r")), mm.Cons(mm.ListToStruct(mm.Cons(mm.SymbolFromRawString("token"), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("string-append")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("r:token")), mm.Cons(func() *mm.Val {
				if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("=")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("r:char")), mm.Cons(mm.SysStringToString("b"), mm.NilVal))) != mm.False {
					return mm.SysStringToString("\b")
				} else {
					return func() *mm.Val {
						if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("=")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("r:char")), mm.Cons(mm.SysStringToString("f"), mm.NilVal))) != mm.False {
							return mm.SysStringToString("\f")
						} else {
							return func() *mm.Val {
								if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("=")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("r:char")), mm.Cons(mm.SysStringToString("n"), mm.NilVal))) != mm.False {
									return mm.SysStringToString("\n")
								} else {
									return func() *mm.Val {
										if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("=")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("r:char")), mm.Cons(mm.SysStringToString("r"), mm.NilVal))) != mm.False {
											return mm.SysStringToString("\r")
										} else {
											return func() *mm.Val {
												if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("=")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("r:char")), mm.Cons(mm.SysStringToString("t"), mm.NilVal))) != mm.False {
													return mm.SysStringToString("\t")
												} else {
													return func() *mm.Val {
														if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("=")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("r:char")), mm.Cons(mm.SysStringToString("v"), mm.NilVal))) != mm.False {
															return mm.SysStringToString("\v")
														} else {
															return mm.LookupDef(env, mm.SymbolFromRawString("r:char"))
														}
													}()
												}
											}()
										}
									}()
								}
							}()
						}
					}()
				}
			}(), mm.NilVal))), mm.NilVal))), mm.NilVal)))
		}))
		mm.Define(env, mm.SymbolFromRawString("clear-token"), mm.NewCompiled(1, false, func(a []*mm.Val) *mm.Val {
			env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("r"), mm.NilVal), mm.SliceToList(a))
			env = env
			return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("assign")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("r")), mm.Cons(mm.ListToStruct(mm.Cons(mm.SymbolFromRawString("token-type"), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("token-type:none")), mm.Cons(mm.SymbolFromRawString("token"), mm.Cons(mm.SysStringToString(""), mm.NilVal))))), mm.NilVal)))
		}))
		mm.Define(env, mm.SymbolFromRawString("make-set-token-type"), mm.NewCompiled(1, false, func(a []*mm.Val) *mm.Val {
			env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("tt"), mm.NilVal), mm.SliceToList(a))
			env = env
			return mm.NewCompiled(1, false, func(a []*mm.Val) *mm.Val {
				env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("r"), mm.NilVal), mm.SliceToList(a))
				env = env
				return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("assign")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("r")), mm.Cons(mm.ListToStruct(mm.Cons(mm.SymbolFromRawString("token-type"), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("tt")), mm.NilVal))), mm.NilVal)))
			})
		}))
		mm.Define(env, mm.SymbolFromRawString("set-symbol"), mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("make-set-token-type")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("token-type:symbol")), mm.NilVal)))
		mm.Define(env, mm.SymbolFromRawString("set-string"), mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("make-set-token-type")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("token-type:string")), mm.NilVal)))
		mm.Define(env, mm.SymbolFromRawString("set-comment"), mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("make-set-token-type")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("token-type:comment")), mm.NilVal)))
		mm.Define(env, mm.SymbolFromRawString("set-list"), mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("make-set-token-type")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("token-type:list")), mm.NilVal)))
		mm.Define(env, mm.SymbolFromRawString("set-quote"), mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("make-set-token-type")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("token-type:quote")), mm.NilVal)))
		mm.Define(env, mm.SymbolFromRawString("set-vector"), mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("make-set-token-type")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("token-type:vector")), mm.NilVal)))
		mm.Define(env, mm.SymbolFromRawString("set-struct"), mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("make-set-token-type")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("token-type:struct")), mm.NilVal)))
		mm.Define(env, mm.SymbolFromRawString("set-close"), mm.NewCompiled(2, false, func(a []*mm.Val) *mm.Val {
			env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("r"), mm.Cons(mm.SymbolFromRawString("list-type"), mm.NilVal)), mm.SliceToList(a))
			env = env
			return func() *mm.Val {
				if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("=")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("list-type")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("r:list-type")), mm.NilVal))) != mm.False {
					return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("assign")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("r")), mm.Cons(mm.ListToStruct(mm.Cons(mm.SymbolFromRawString("close-list?"), mm.Cons(mm.True, mm.NilVal))), mm.NilVal)))
				} else {
					return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("assign")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("r")), mm.Cons(mm.ListToStruct(mm.Cons(mm.SymbolFromRawString("state"), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("unexpected-close")), mm.NilVal))), mm.NilVal)))
				}
			}()
		}))
		mm.Define(env, mm.SymbolFromRawString("set-cons"), mm.NewCompiled(1, false, func(a []*mm.Val) *mm.Val {
			env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("r"), mm.NilVal), mm.SliceToList(a))
			env = env
			return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("assign")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("r")), mm.Cons(func() *mm.Val {
				if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("=")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("r:list-type")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("list-type:lisp")), mm.NilVal))) != mm.False {
					return mm.ListToStruct(mm.Cons(mm.SymbolFromRawString("list-cons?"), mm.Cons(mm.True, mm.NilVal)))
				} else {
					return mm.ListToStruct(mm.Cons(mm.SymbolFromRawString("state"), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("irregular-cons")), mm.NilVal)))
				}
			}(), mm.NilVal)))
		}))
		mm.Define(env, mm.SymbolFromRawString("symbol-token"), mm.NewCompiled(1, false, func(a []*mm.Val) *mm.Val {
			env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("r"), mm.NilVal), mm.SliceToList(a))
			env = env
			mm.Define(env, mm.SymbolFromRawString("try-parse"), mm.NewCompiled(0, true, func(a []*mm.Val) *mm.Val {
				env := mm.ExtendEnv(env, mm.SymbolFromRawString("parsers"), mm.SliceToList(a))
				env = env
				return func() *mm.Val {
					if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("nil?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("parsers")), mm.NilVal)) != mm.False {
						return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("string->symbol")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("r:token")), mm.NilVal))
					} else {
						return mm.ApplySys(mm.NewCompiled(0, false, func(a []*mm.Val) *mm.Val {
							env := mm.ExtendEnv(env, mm.NilVal, mm.SliceToList(a))
							env = env
							mm.Define(env, mm.SymbolFromRawString("v"), mm.ApplySys(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("parsers")), mm.NilVal)), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("r:token")), mm.NilVal)))
							return func() *mm.Val {
								if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("error?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)) != mm.False {
									return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("apply")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("try-parse")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("parsers")), mm.NilVal)), mm.NilVal)))
								} else {
									return mm.LookupDef(env, mm.SymbolFromRawString("v"))
								}
							}()
						}), mm.NilVal)
					}
				}()
			}))
			return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("assign")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("r")), mm.Cons(mm.ListToStruct(mm.Cons(mm.SymbolFromRawString("state"), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("try-parse")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("string->number")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("string->bool")), mm.NilVal))), mm.NilVal))), mm.NilVal)))
		}))
		mm.Define(env, mm.SymbolFromRawString("finalize-token"), mm.NewCompiled(1, false, func(a []*mm.Val) *mm.Val {
			env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("r"), mm.NilVal), mm.SliceToList(a))
			env = env
			return func() *mm.Val {
				if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("=")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("r:token-type")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("token-type:symbol")), mm.NilVal))) != mm.False {
					return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("clear-token")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("symbol-token")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("r")), mm.NilVal)), mm.NilVal))
				} else {
					return func() *mm.Val {
						if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("=")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("r:token-type")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("token-type:string")), mm.NilVal))) != mm.False {
							return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("clear-token")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("assign")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("r")), mm.Cons(mm.ListToStruct(mm.Cons(mm.SymbolFromRawString("state"), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("r:token")), mm.NilVal))), mm.NilVal))), mm.NilVal))
						} else {
							return func() *mm.Val {
								if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("=")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("r:token-type")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("token-type:none")), mm.NilVal))) != mm.False {
									return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("clear-token")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("r")), mm.NilVal))
								} else {
									return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("assign")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("r")), mm.Cons(mm.ListToStruct(mm.Cons(mm.SymbolFromRawString("state"), mm.Cons(func() *mm.Val {
										if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("=")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("r:token")), mm.Cons(mm.SysStringToString(""), mm.NilVal))) != mm.False {
											return mm.LookupDef(env, mm.SymbolFromRawString("undefined"))
										} else {
											return mm.LookupDef(env, mm.SymbolFromRawString("invalid-token"))
										}
									}(), mm.NilVal))), mm.NilVal)))
								}
							}()
						}
					}()
				}
			}()
		}))
		mm.Define(env, mm.SymbolFromRawString("read-list"), mm.NewCompiled(2, false, func(a []*mm.Val) *mm.Val {
			env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("r"), mm.Cons(mm.SymbolFromRawString("list-type"), mm.NilVal)), mm.SliceToList(a))
			env = env
			mm.Define(env, mm.SymbolFromRawString("read-item"), mm.NewCompiled(1, false, func(a []*mm.Val) *mm.Val {
				env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("lr"), mm.NilVal), mm.SliceToList(a))
				env = env
				return mm.ApplySys(mm.NewCompiled(0, false, func(a []*mm.Val) *mm.Val {
					env := mm.ExtendEnv(env, mm.NilVal, mm.SliceToList(a))
					env = env
					mm.Define(env, mm.SymbolFromRawString("next"), mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("read")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("lr")), mm.NilVal)))
					return func() *mm.Val {
						if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("error?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next:state")), mm.NilVal)) != mm.False {
							return mm.LookupDef(env, mm.SymbolFromRawString("next"))
						} else {
							return func() *mm.Val {
								if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("=")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next:state")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("undefined")), mm.NilVal))) != mm.False {
									return mm.LookupDef(env, mm.SymbolFromRawString("next"))
								} else {
									return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("assign")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next")), mm.Cons(mm.ListToStruct(mm.Cons(mm.SymbolFromRawString("list-items"), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cons")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next:state")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next:list-items")), mm.NilVal))), mm.Cons(mm.SymbolFromRawString("cons-items"), mm.Cons(func() *mm.Val {
										if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString(">")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next:cons-items")), mm.Cons(mm.SysIntToNumber(0), mm.NilVal))) != mm.False {
											return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("inc")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("lr:cons-items")), mm.NilVal))
										} else {
											return mm.SysIntToNumber(0)
										}
									}(), mm.Cons(mm.SymbolFromRawString("state"), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("undefined")), mm.NilVal))))))), mm.NilVal)))
								}
							}()
						}
					}()
				}), mm.NilVal)
			}))
			mm.Define(env, mm.SymbolFromRawString("check-cons"), mm.NewCompiled(1, false, func(a []*mm.Val) *mm.Val {
				env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("lr"), mm.NilVal), mm.SliceToList(a))
				env = env
				return func() *mm.Val {
					if mm.LookupDef(env, mm.SymbolFromRawString("lr:list-cons?")) != mm.False {
						return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("assign")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("lr")), mm.Cons(func() *mm.Val {
							if func() *mm.Val {
								if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("nil?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("lr:list-items")), mm.NilVal)) != mm.False {
									return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("nil?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("lr:list-items")), mm.NilVal))
								} else {
									return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString(">")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("lr:cons-items")), mm.Cons(mm.SysIntToNumber(0), mm.NilVal)))
								}
							}() != mm.False {
								return mm.ListToStruct(mm.Cons(mm.SymbolFromRawString("state"), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("irregular-cons")), mm.NilVal)))
							} else {
								return mm.ListToStruct(mm.Cons(mm.SymbolFromRawString("cons-items"), mm.Cons(mm.SysIntToNumber(1), mm.Cons(mm.SymbolFromRawString("list-cons?"), mm.Cons(mm.False, mm.NilVal)))))
							}
						}(), mm.NilVal)))
					} else {
						return mm.LookupDef(env, mm.SymbolFromRawString("lr"))
					}
				}()
			}))
			mm.Define(env, mm.SymbolFromRawString("complete-list"), mm.NewCompiled(1, false, func(a []*mm.Val) *mm.Val {
				env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("lr"), mm.NilVal), mm.SliceToList(a))
				env = env
				return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("assign")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("r")), mm.Cons(mm.ListToStruct(mm.Cons(mm.SymbolFromRawString("input"), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("lr:input")), mm.Cons(mm.SymbolFromRawString("token-type"), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("token-type:none")), mm.Cons(mm.SymbolFromRawString("state"), mm.Cons(func() *mm.Val {
					if func() *mm.Val {
						if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString(">")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("lr:cons-items")), mm.Cons(mm.SysIntToNumber(0), mm.NilVal))) != mm.False {
							return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("not")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("=")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("lr:cons-items")), mm.Cons(mm.SysIntToNumber(2), mm.NilVal))), mm.NilVal))
						} else {
							return mm.False
						}
					}() != mm.False {
						return mm.LookupDef(env, mm.SymbolFromRawString("irregular-cons"))
					} else {
						return func() *mm.Val {
							if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString(">")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("lr:cons-items")), mm.Cons(mm.SysIntToNumber(0), mm.NilVal))) != mm.False {
								return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("reverse-irregular")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("lr:list-items")), mm.NilVal))
							} else {
								return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("reverse")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("lr:list-items")), mm.NilVal))
							}
						}()
					}
				}(), mm.NilVal))))))), mm.NilVal)))
			}))
			mm.Define(env, mm.SymbolFromRawString("read-items"), mm.NewCompiled(1, false, func(a []*mm.Val) *mm.Val {
				env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("lr"), mm.NilVal), mm.SliceToList(a))
				env = env
				return mm.ApplySys(mm.NewCompiled(0, false, func(a []*mm.Val) *mm.Val {
					env := mm.ExtendEnv(env, mm.NilVal, mm.SliceToList(a))
					env = env
					mm.Define(env, mm.SymbolFromRawString("next"), mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("check-cons")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("read-item")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("lr")), mm.NilVal)), mm.NilVal)))
					return func() *mm.Val {
						if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("error?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next:state")), mm.NilVal)) != mm.False {
							return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("assign")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("r")), mm.Cons(mm.ListToStruct(mm.Cons(mm.SymbolFromRawString("input"), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next:input")), mm.Cons(mm.SymbolFromRawString("state"), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next:state")), mm.NilVal))))), mm.NilVal)))
						} else {
							return func() *mm.Val {
								if mm.LookupDef(env, mm.SymbolFromRawString("next:close-list?")) != mm.False {
									return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("complete-list")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next")), mm.NilVal))
								} else {
									return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("read-items")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next")), mm.NilVal))
								}
							}()
						}
					}()
				}), mm.NilVal)
			}))
			return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("read-items")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("assign")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("reader")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("r:input")), mm.NilVal)), mm.Cons(mm.ListToStruct(mm.Cons(mm.SymbolFromRawString("list-type"), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("list-type")), mm.NilVal))), mm.NilVal))), mm.NilVal))
		}))
		mm.Define(env, mm.SymbolFromRawString("read-quote"), mm.NewCompiled(1, false, func(a []*mm.Val) *mm.Val {
			env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("r"), mm.NilVal), mm.SliceToList(a))
			env = env
			return mm.ApplySys(mm.NewCompiled(0, false, func(a []*mm.Val) *mm.Val {
				env := mm.ExtendEnv(env, mm.NilVal, mm.SliceToList(a))
				env = env
				mm.Define(env, mm.SymbolFromRawString("lr"), mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("assign")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("reader")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("r:input")), mm.NilVal)), mm.Cons(mm.ListToStruct(mm.Cons(mm.SymbolFromRawString("list-type"), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("r:list-type")), mm.NilVal))), mm.NilVal))))
				mm.Define(env, mm.SymbolFromRawString("next"), mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("read")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("lr")), mm.NilVal)))
				return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("assign")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("r")), mm.Cons(mm.ListToStruct(mm.Cons(mm.SymbolFromRawString("input"), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next:input")), mm.Cons(mm.SymbolFromRawString("token-type"), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("token-type:none")), mm.Cons(mm.SymbolFromRawString("state"), mm.Cons(func() *mm.Val {
					if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("error?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next:state")), mm.NilVal)) != mm.False {
						return mm.LookupDef(env, mm.SymbolFromRawString("next:state"))
					} else {
						return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("list")), mm.Cons(mm.SymbolFromRawString("quote"), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next:state")), mm.NilVal)))
					}
				}(), mm.Cons(mm.SymbolFromRawString("close-list?"), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next:close-list?")), mm.NilVal))))))))), mm.NilVal)))
			}), mm.NilVal)
		}))
		mm.Define(env, mm.SymbolFromRawString("read-vector"), mm.NewCompiled(1, false, func(a []*mm.Val) *mm.Val {
			env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("r"), mm.NilVal), mm.SliceToList(a))
			env = env
			return mm.ApplySys(mm.NewCompiled(0, false, func(a []*mm.Val) *mm.Val {
				env := mm.ExtendEnv(env, mm.NilVal, mm.SliceToList(a))
				env = env
				mm.Define(env, mm.SymbolFromRawString("next"), mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("read-list")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("r")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("list-type:vector")), mm.NilVal))))
				return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("assign")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next")), mm.Cons(mm.ListToStruct(mm.Cons(mm.SymbolFromRawString("state"), mm.Cons(func() *mm.Val {
					if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("error?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next:state")), mm.NilVal)) != mm.False {
						return mm.LookupDef(env, mm.SymbolFromRawString("next:state"))
					} else {
						return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cons")), mm.Cons(mm.SymbolFromRawString("vector:"), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next:state")), mm.NilVal)))
					}
				}(), mm.NilVal))), mm.NilVal)))
			}), mm.NilVal)
		}))
		mm.Define(env, mm.SymbolFromRawString("read-struct"), mm.NewCompiled(1, false, func(a []*mm.Val) *mm.Val {
			env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("r"), mm.NilVal), mm.SliceToList(a))
			env = env
			return mm.ApplySys(mm.NewCompiled(0, false, func(a []*mm.Val) *mm.Val {
				env := mm.ExtendEnv(env, mm.NilVal, mm.SliceToList(a))
				env = env
				mm.Define(env, mm.SymbolFromRawString("next"), mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("read-list")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("r")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("list-type:struct")), mm.NilVal))))
				return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("assign")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next")), mm.Cons(mm.ListToStruct(mm.Cons(mm.SymbolFromRawString("state"), mm.Cons(func() *mm.Val {
					if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("error?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next:state")), mm.NilVal)) != mm.False {
						return mm.LookupDef(env, mm.SymbolFromRawString("next:state"))
					} else {
						return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cons")), mm.Cons(mm.SymbolFromRawString("struct:"), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next:state")), mm.NilVal)))
					}
				}(), mm.NilVal))), mm.NilVal)))
			}), mm.NilVal)
		}))
		mm.Define(env, mm.SymbolFromRawString("read"), mm.NewCompiled(1, false, func(a []*mm.Val) *mm.Val {
			env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("r"), mm.NilVal), mm.SliceToList(a))
			env = env
			return func() *mm.Val {
				if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("=")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("r:token-type")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("token-type:list")), mm.NilVal))) != mm.False {
					return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("assign")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("read-list")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("r")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("list-type:lisp")), mm.NilVal))), mm.Cons(mm.ListToStruct(mm.Cons(mm.SymbolFromRawString("token-type:"), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("token-type:none")), mm.NilVal))), mm.NilVal)))
				} else {
					return func() *mm.Val {
						if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("=")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("r:token-type")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("token-type:quote")), mm.NilVal))) != mm.False {
							return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("assign")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("read-quote")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("r")), mm.NilVal)), mm.Cons(mm.ListToStruct(mm.Cons(mm.SymbolFromRawString("token-type:"), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("token-type:none")), mm.NilVal))), mm.NilVal)))
						} else {
							return func() *mm.Val {
								if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("=")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("r:token-type")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("token-type:vector")), mm.NilVal))) != mm.False {
									return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("assign")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("read-vector")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("r")), mm.NilVal)), mm.Cons(mm.ListToStruct(mm.Cons(mm.SymbolFromRawString("token-type:"), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("token-type:none")), mm.NilVal))), mm.NilVal)))
								} else {
									return func() *mm.Val {
										if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("=")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("r:token-type")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("token-type:struct")), mm.NilVal))) != mm.False {
											return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("assign")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("read-struct")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("r")), mm.NilVal)), mm.Cons(mm.ListToStruct(mm.Cons(mm.SymbolFromRawString("token-type:"), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("token-type:none")), mm.NilVal))), mm.NilVal)))
										} else {
											return mm.ApplySys(mm.NewCompiled(0, false, func(a []*mm.Val) *mm.Val {
												env := mm.ExtendEnv(env, mm.NilVal, mm.SliceToList(a))
												env = env
												mm.Define(env, mm.SymbolFromRawString("next"), mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("read-char")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("r")), mm.NilVal)))
												return func() *mm.Val {
													if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("=")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next:state")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("eof")), mm.NilVal))) != mm.False {
														return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("finalize-token")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next")), mm.NilVal))
													} else {
														return func() *mm.Val {
															if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("error?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next:state")), mm.NilVal)) != mm.False {
																return mm.LookupDef(env, mm.SymbolFromRawString("next"))
															} else {
																return func() *mm.Val {
																	if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("=")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next:token-type")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("token-type:none")), mm.NilVal))) != mm.False {
																		return func() *mm.Val {
																			if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("whitespace?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next:char")), mm.NilVal)) != mm.False {
																				return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("read")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next")), mm.NilVal))
																			} else {
																				return func() *mm.Val {
																					if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("escape-char?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next:char")), mm.NilVal)) != mm.False {
																						return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("read")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("set-symbol")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("set-escaped")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next")), mm.NilVal)), mm.NilVal)), mm.NilVal))
																					} else {
																						return func() *mm.Val {
																							if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("string-delimiter?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next:char")), mm.NilVal)) != mm.False {
																								return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("read")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("set-string")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next")), mm.NilVal)), mm.NilVal))
																							} else {
																								return func() *mm.Val {
																									if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("comment-char?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next:char")), mm.NilVal)) != mm.False {
																										return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("read")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("set-comment")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next")), mm.NilVal)), mm.NilVal))
																									} else {
																										return func() *mm.Val {
																											if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("list-open?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next:char")), mm.NilVal)) != mm.False {
																												return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("read")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("set-list")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next")), mm.NilVal)), mm.NilVal))
																											} else {
																												return func() *mm.Val {
																													if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("list-close?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next:char")), mm.NilVal)) != mm.False {
																														return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("set-close")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("list-type:lisp")), mm.NilVal)))
																													} else {
																														return func() *mm.Val {
																															if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cons-char?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next:char")), mm.NilVal)) != mm.False {
																																return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("set-cons")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next")), mm.NilVal))
																															} else {
																																return func() *mm.Val {
																																	if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("quote-char?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next:char")), mm.NilVal)) != mm.False {
																																		return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("read")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("set-quote")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next")), mm.NilVal)), mm.NilVal))
																																	} else {
																																		return func() *mm.Val {
																																			if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("vector-open?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next:char")), mm.NilVal)) != mm.False {
																																				return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("read")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("set-vector")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next")), mm.NilVal)), mm.NilVal))
																																			} else {
																																				return func() *mm.Val {
																																					if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("vector-close?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next:char")), mm.NilVal)) != mm.False {
																																						return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("set-close")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("list-type:vector")), mm.NilVal)))
																																					} else {
																																						return func() *mm.Val {
																																							if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("struct-open?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next:char")), mm.NilVal)) != mm.False {
																																								return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("read")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("set-struct")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next")), mm.NilVal)), mm.NilVal))
																																							} else {
																																								return func() *mm.Val {
																																									if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("struct-close?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next:char")), mm.NilVal)) != mm.False {
																																										return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("set-close")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("list-type:struct")), mm.NilVal)))
																																									} else {
																																										return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("read")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("append-token")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("set-symbol")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next")), mm.NilVal)), mm.NilVal)), mm.NilVal))
																																									}
																																								}()
																																							}
																																						}()
																																					}
																																				}()
																																			}
																																		}()
																																	}
																																}()
																															}
																														}()
																													}
																												}()
																											}
																										}()
																									}
																								}()
																							}
																						}()
																					}
																				}()
																			}
																		}()
																	} else {
																		return func() *mm.Val {
																			if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("=")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next:token-type")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("token-type:comment")), mm.NilVal))) != mm.False {
																				return func() *mm.Val {
																					if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("newline?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next:char")), mm.NilVal)) != mm.False {
																						return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("read")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("clear-token")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next")), mm.NilVal)), mm.NilVal))
																					} else {
																						return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("read")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next")), mm.NilVal))
																					}
																				}()
																			} else {
																				return func() *mm.Val {
																					if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("=")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next:token-type")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("token-type:symbol")), mm.NilVal))) != mm.False {
																						return func() *mm.Val {
																							if mm.LookupDef(env, mm.SymbolFromRawString("next:escaped?")) != mm.False {
																								return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("read")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("assign")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("append-token-escaped")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next")), mm.NilVal)), mm.Cons(mm.ListToStruct(mm.Cons(mm.SymbolFromRawString("escaped?"), mm.Cons(mm.False, mm.NilVal))), mm.NilVal))), mm.NilVal))
																							} else {
																								return func() *mm.Val {
																									if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("whitespace?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next:char")), mm.NilVal)) != mm.False {
																										return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("finalize-token")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next")), mm.NilVal))
																									} else {
																										return func() *mm.Val {
																											if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("escape-char?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next:char")), mm.NilVal)) != mm.False {
																												return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("read")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("set-escaped")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next")), mm.NilVal)), mm.NilVal))
																											} else {
																												return func() *mm.Val {
																													if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("string-delimiter?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next:char")), mm.NilVal)) != mm.False {
																														return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("set-string")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("finalize-token")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next")), mm.NilVal)), mm.NilVal))
																													} else {
																														return func() *mm.Val {
																															if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("comment-char?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next:char")), mm.NilVal)) != mm.False {
																																return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("set-comment")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("finalize-token")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next")), mm.NilVal)), mm.NilVal))
																															} else {
																																return func() *mm.Val {
																																	if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("list-open?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next:char")), mm.NilVal)) != mm.False {
																																		return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("set-list")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("finalize-token")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next")), mm.NilVal)), mm.NilVal))
																																	} else {
																																		return func() *mm.Val {
																																			if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("list-close?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next:char")), mm.NilVal)) != mm.False {
																																				return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("set-close")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("finalize-token")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next")), mm.NilVal)), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("list-type:lisp")), mm.NilVal)))
																																			} else {
																																				return func() *mm.Val {
																																					if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("vector-open?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next:char")), mm.NilVal)) != mm.False {
																																						return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("set-vector")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("finalize-token")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next")), mm.NilVal)), mm.NilVal))
																																					} else {
																																						return func() *mm.Val {
																																							if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("vector-close?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next:char")), mm.NilVal)) != mm.False {
																																								return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("set-close")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("finalize-token")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next")), mm.NilVal)), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("list-type:vector")), mm.NilVal)))
																																							} else {
																																								return func() *mm.Val {
																																									if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("struct-open?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next:char")), mm.NilVal)) != mm.False {
																																										return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("set-struct")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("finalize-token")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next")), mm.NilVal)), mm.NilVal))
																																									} else {
																																										return func() *mm.Val {
																																											if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("struct-close?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next:char")), mm.NilVal)) != mm.False {
																																												return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("set-close")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("finalize-token")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next")), mm.NilVal)), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("list-type:struct")), mm.NilVal)))
																																											} else {
																																												return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("read")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("append-token")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next")), mm.NilVal)), mm.NilVal))
																																											}
																																										}()
																																									}
																																								}()
																																							}
																																						}()
																																					}
																																				}()
																																			}
																																		}()
																																	}
																																}()
																															}
																														}()
																													}
																												}()
																											}
																										}()
																									}
																								}()
																							}
																						}()
																					} else {
																						return func() *mm.Val {
																							if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("=")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next:token-type")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("token-type:string")), mm.NilVal))) != mm.False {
																								return func() *mm.Val {
																									if mm.LookupDef(env, mm.SymbolFromRawString("next:escaped?")) != mm.False {
																										return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("read")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("assign")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("append-token-escaped")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next")), mm.NilVal)), mm.Cons(mm.ListToStruct(mm.Cons(mm.SymbolFromRawString("escaped?"), mm.Cons(mm.False, mm.NilVal))), mm.NilVal))), mm.NilVal))
																									} else {
																										return func() *mm.Val {
																											if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("escape-char?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next:char")), mm.NilVal)) != mm.False {
																												return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("read")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("set-escaped")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next")), mm.NilVal)), mm.NilVal))
																											} else {
																												return func() *mm.Val {
																													if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("string-delimiter?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next:char")), mm.NilVal)) != mm.False {
																														return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("finalize-token")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next")), mm.NilVal))
																													} else {
																														return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("read")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("append-token")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next")), mm.NilVal)), mm.NilVal))
																													}
																												}()
																											}
																										}()
																									}
																								}()
																							} else {
																								return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("assign")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("next")), mm.Cons(mm.ListToStruct(mm.Cons(mm.SymbolFromRawString("state"), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("invalid-token")), mm.NilVal))), mm.NilVal)))
																							}
																						}()
																					}
																				}()
																			}
																		}()
																	}
																}()
															}
														}()
													}
												}()
											}), mm.NilVal)
										}
									}()
								}
							}()
						}
					}()
				}
			}()
		}))
		mm.Export(env, mm.ListToStruct(mm.Cons(mm.SymbolFromRawString("reader"), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("reader")), mm.Cons(mm.SymbolFromRawString("read"), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("read")), mm.NilVal))))))
		mm.StoreModule(eenv, mm.SysStringToString("scm/read.scm"), mm.Exports(env))
	})
	func() {
		m := mm.LoadedModule(env, mm.SysStringToString("scm/read.scm"))
		if m == mm.UndefinedModule {
			mm.LoadCompiledModule(initialEnv, mm.SysStringToString("scm/read.scm"))
			m = mm.LoadedModule(env, mm.SysStringToString("scm/read.scm"))
		}
		mm.DefineAll(env, m)
	}()
	mm.Define(env, mm.SymbolFromRawString("tagged?"), mm.NewCompiled(2, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("v"), mm.Cons(mm.SymbolFromRawString("t"), mm.NilVal)), mm.SliceToList(a))
		env = env
		return func() *mm.Val {
			if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("pair?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)) != mm.False {
				return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("=")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("t")), mm.NilVal)))
			} else {
				return mm.False
			}
		}()
	}))
	mm.Define(env, mm.SymbolFromRawString("quote?"), mm.NewCompiled(1, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("v"), mm.NilVal), mm.SliceToList(a))
		env = env
		return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("tagged?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.Cons(mm.SymbolFromRawString("quote"), mm.NilVal)))
	}))
	mm.Define(env, mm.SymbolFromRawString("def?"), mm.NewCompiled(1, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("v"), mm.NilVal), mm.SliceToList(a))
		env = env
		return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("tagged?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.Cons(mm.SymbolFromRawString("def"), mm.NilVal)))
	}))
	mm.Define(env, mm.SymbolFromRawString("vector-form?"), mm.NewCompiled(1, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("v"), mm.NilVal), mm.SliceToList(a))
		env = env
		return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("tagged?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.Cons(mm.SymbolFromRawString("vector:"), mm.NilVal)))
	}))
	mm.Define(env, mm.SymbolFromRawString("struct-form?"), mm.NewCompiled(1, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("v"), mm.NilVal), mm.SliceToList(a))
		env = env
		return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("tagged?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.Cons(mm.SymbolFromRawString("struct:"), mm.NilVal)))
	}))
	mm.Define(env, mm.SymbolFromRawString("if?"), mm.NewCompiled(1, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("v"), mm.NilVal), mm.SliceToList(a))
		env = env
		return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("tagged?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.Cons(mm.SymbolFromRawString("if"), mm.NilVal)))
	}))
	mm.Define(env, mm.SymbolFromRawString("and?"), mm.NewCompiled(1, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("v"), mm.NilVal), mm.SliceToList(a))
		env = env
		return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("tagged?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.Cons(mm.SymbolFromRawString("and"), mm.NilVal)))
	}))
	mm.Define(env, mm.SymbolFromRawString("or?"), mm.NewCompiled(1, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("v"), mm.NilVal), mm.SliceToList(a))
		env = env
		return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("tagged?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.Cons(mm.SymbolFromRawString("or"), mm.NilVal)))
	}))
	mm.Define(env, mm.SymbolFromRawString("fn?"), mm.NewCompiled(1, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("v"), mm.NilVal), mm.SliceToList(a))
		env = env
		return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("tagged?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.Cons(mm.SymbolFromRawString("fn"), mm.NilVal)))
	}))
	mm.Define(env, mm.SymbolFromRawString("begin?"), mm.NewCompiled(1, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("v"), mm.NilVal), mm.SliceToList(a))
		env = env
		return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("tagged?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.Cons(mm.SymbolFromRawString("begin"), mm.NilVal)))
	}))
	mm.Define(env, mm.SymbolFromRawString("cond?"), mm.NewCompiled(1, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("v"), mm.NilVal), mm.SliceToList(a))
		env = env
		return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("tagged?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.Cons(mm.SymbolFromRawString("cond"), mm.NilVal)))
	}))
	mm.Define(env, mm.SymbolFromRawString("let?"), mm.NewCompiled(1, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("v"), mm.NilVal), mm.SliceToList(a))
		env = env
		return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("tagged?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.Cons(mm.SymbolFromRawString("let"), mm.NilVal)))
	}))
	mm.Define(env, mm.SymbolFromRawString("test?"), mm.NewCompiled(1, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("v"), mm.NilVal), mm.SliceToList(a))
		env = env
		return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("tagged?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.Cons(mm.SymbolFromRawString("test"), mm.NilVal)))
	}))
	mm.Define(env, mm.SymbolFromRawString("export?"), mm.NewCompiled(1, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("v"), mm.NilVal), mm.SliceToList(a))
		env = env
		return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("tagged?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.Cons(mm.SymbolFromRawString("export"), mm.NilVal)))
	}))
	mm.Define(env, mm.SymbolFromRawString("import?"), mm.NewCompiled(1, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("v"), mm.NilVal), mm.SliceToList(a))
		env = env
		return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("tagged?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.Cons(mm.SymbolFromRawString("import"), mm.NilVal)))
	}))
	mm.Define(env, mm.SymbolFromRawString("application?"), mm.NewCompiled(1, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("v"), mm.NilVal), mm.SliceToList(a))
		env = env
		return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("pair?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal))
	}))
	mm.Define(env, mm.SymbolFromRawString("current-env?"), mm.NewCompiled(1, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("v"), mm.NilVal), mm.SliceToList(a))
		env = env
		return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("tagged?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.Cons(mm.SymbolFromRawString("current-env"), mm.NilVal)))
	}))
	mm.Define(env, mm.SymbolFromRawString("eval-quote"), mm.NewCompiled(1, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("exp"), mm.NilVal), mm.SliceToList(a))
		env = env
		return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("exp")), mm.NilVal)), mm.NilVal))
	}))
	mm.Define(env, mm.SymbolFromRawString("eval-def"), mm.NewCompiled(2, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("env"), mm.Cons(mm.SymbolFromRawString("exp"), mm.NilVal)), mm.SliceToList(a))
		env = env
		return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("define")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("env")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("def-name")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("exp")), mm.NilVal)), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("eval-exp")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("env")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("def-value")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("exp")), mm.NilVal)), mm.NilVal))), mm.NilVal))))
	}))
	mm.Define(env, mm.SymbolFromRawString("value-list"), mm.NewCompiled(2, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("env"), mm.Cons(mm.SymbolFromRawString("exp"), mm.NilVal)), mm.SliceToList(a))
		env = env
		return func() *mm.Val {
			if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("nil?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("exp")), mm.NilVal)) != mm.False {
				return mm.LookupDef(env, mm.SymbolFromRawString("nil"))
			} else {
				return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cons")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("eval-exp")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("env")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("exp")), mm.NilVal)), mm.NilVal))), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("value-list")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("env")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("exp")), mm.NilVal)), mm.NilVal))), mm.NilVal)))
			}
		}()
	}))
	mm.Define(env, mm.SymbolFromRawString("eval-vector"), mm.NewCompiled(2, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("env"), mm.Cons(mm.SymbolFromRawString("exp"), mm.NilVal)), mm.SliceToList(a))
		env = env
		return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("list->vector")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("value-list")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("env")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("exp")), mm.NilVal)), mm.NilVal))), mm.NilVal))
	}))
	mm.Define(env, mm.SymbolFromRawString("struct-values"), mm.NewCompiled(2, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("env"), mm.Cons(mm.SymbolFromRawString("exp"), mm.NilVal)), mm.SliceToList(a))
		env = env
		return func() *mm.Val {
			if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("nil?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("exp")), mm.NilVal)) != mm.False {
				return mm.LookupDef(env, mm.SymbolFromRawString("nil"))
			} else {
				return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cons")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("exp")), mm.NilVal)), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cons")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("eval-exp")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("env")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("exp")), mm.NilVal)), mm.NilVal)), mm.NilVal))), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("struct-values")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("env")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("exp")), mm.NilVal)), mm.NilVal)), mm.NilVal))), mm.NilVal))), mm.NilVal)))
			}
		}()
	}))
	mm.Define(env, mm.SymbolFromRawString("eval-struct"), mm.NewCompiled(2, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("env"), mm.Cons(mm.SymbolFromRawString("exp"), mm.NilVal)), mm.SliceToList(a))
		env = env
		return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("list->struct")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("struct-values")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("env")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("exp")), mm.NilVal)), mm.NilVal))), mm.NilVal))
	}))
	mm.Define(env, mm.SymbolFromRawString("eval-exp"), mm.NewCompiled(2, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("env"), mm.Cons(mm.SymbolFromRawString("exp"), mm.NilVal)), mm.SliceToList(a))
		env = env
		return func() *mm.Val {
			if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("def?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("exp")), mm.NilVal)) != mm.False {
				return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("fatal")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("definition-expression")), mm.NilVal))
			} else {
				return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("eval-env")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("env")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("exp")), mm.NilVal)))
			}
		}()
	}))
	mm.Define(env, mm.SymbolFromRawString("eval-if"), mm.NewCompiled(2, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("env"), mm.Cons(mm.SymbolFromRawString("exp"), mm.NilVal)), mm.SliceToList(a))
		env = env
		return func() *mm.Val {
			if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("eval-exp")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("env")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("exp")), mm.NilVal)), mm.NilVal)), mm.NilVal))) != mm.False {
				return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("eval-exp")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("env")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("exp")), mm.NilVal)), mm.NilVal)), mm.NilVal)), mm.NilVal)))
			} else {
				return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("eval-exp")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("env")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("exp")), mm.NilVal)), mm.NilVal)), mm.NilVal)), mm.NilVal)), mm.NilVal)))
			}
		}()
	}))
	mm.Define(env, mm.SymbolFromRawString("eval-and"), mm.NewCompiled(2, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("env"), mm.Cons(mm.SymbolFromRawString("exp"), mm.NilVal)), mm.SliceToList(a))
		env = env
		return func() *mm.Val {
			if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("nil?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("exp")), mm.NilVal)) != mm.False {
				return mm.True
			} else {
				return func() *mm.Val {
					if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("nil?")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("exp")), mm.NilVal)), mm.NilVal)) != mm.False {
						return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("eval-exp")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("env")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("exp")), mm.NilVal)), mm.NilVal)))
					} else {
						return func() *mm.Val {
							if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("not")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("eval-exp")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("env")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("exp")), mm.NilVal)), mm.NilVal))), mm.NilVal)) != mm.False {
								return mm.False
							} else {
								return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("eval-and")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("env")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("exp")), mm.NilVal)), mm.NilVal)))
							}
						}()
					}
				}()
			}
		}()
	}))
	mm.Define(env, mm.SymbolFromRawString("eval-or"), mm.NewCompiled(2, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("env"), mm.Cons(mm.SymbolFromRawString("exp"), mm.NilVal)), mm.SliceToList(a))
		env = env
		return func() *mm.Val {
			if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("nil?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("exp")), mm.NilVal)) != mm.False {
				return mm.False
			} else {
				return func() *mm.Val {
					if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("nil?")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("exp")), mm.NilVal)), mm.NilVal)) != mm.False {
						return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("eval-exp")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("env")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("exp")), mm.NilVal)), mm.NilVal)))
					} else {
						return mm.ApplySys(mm.NewCompiled(0, false, func(a []*mm.Val) *mm.Val {
							env := mm.ExtendEnv(env, mm.NilVal, mm.SliceToList(a))
							env = env
							mm.Define(env, mm.SymbolFromRawString("v"), mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("eval-exp")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("env")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("exp")), mm.NilVal)), mm.NilVal))))
							return func() *mm.Val {
								if mm.LookupDef(env, mm.SymbolFromRawString("v")) != mm.False {
									return mm.LookupDef(env, mm.SymbolFromRawString("v"))
								} else {
									return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("eval-or")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("env")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("exp")), mm.NilVal)), mm.NilVal)))
								}
							}()
						}), mm.NilVal)
					}
				}()
			}
		}()
	}))
	mm.Define(env, mm.SymbolFromRawString("eval-fn"), mm.NewCompiled(2, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("env"), mm.Cons(mm.SymbolFromRawString("exp"), mm.NilVal)), mm.SliceToList(a))
		env = env
		return mm.ApplySys(mm.NewCompiled(0, false, func(a []*mm.Val) *mm.Val {
			env := mm.ExtendEnv(env, mm.NilVal, mm.SliceToList(a))
			env = env
			mm.Define(env, mm.SymbolFromRawString("signature"), mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("fn-signature")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("exp")), mm.NilVal)), mm.NilVal)), mm.NilVal)))
			mm.Define(env, mm.SymbolFromRawString("body"), mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("exp")), mm.NilVal)), mm.NilVal)))
			return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("make-composite")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cons")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("env")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cons")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("signature:names")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("body")), mm.NilVal))), mm.NilVal))), mm.NilVal))
		}), mm.NilVal)
	}))
	mm.Define(env, mm.SymbolFromRawString("eval-seq"), mm.NewCompiled(2, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("env"), mm.Cons(mm.SymbolFromRawString("exp"), mm.NilVal)), mm.SliceToList(a))
		env = env
		return func() *mm.Val {
			if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("not")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("pair?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("exp")), mm.NilVal)), mm.NilVal)) != mm.False {
				return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("fatal")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("invalid-sequence")), mm.NilVal))
			} else {
				return func() *mm.Val {
					if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("nil?")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("exp")), mm.NilVal)), mm.NilVal)) != mm.False {
						return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("eval-exp")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("env")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("exp")), mm.NilVal)), mm.NilVal)))
					} else {
						return func() *mm.Val {
							mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("eval-env")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("env")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("exp")), mm.NilVal)), mm.NilVal)))
							return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("eval-seq")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("env")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("exp")), mm.NilVal)), mm.NilVal)))
						}()
					}
				}()
			}
		}()
	}))
	mm.Define(env, mm.SymbolFromRawString("eval-test"), mm.NewCompiled(2, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("env"), mm.Cons(mm.SymbolFromRawString("exp"), mm.NilVal)), mm.SliceToList(a))
		env = env
		return func() *mm.Val {
			if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("nil?")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("exp")), mm.NilVal)), mm.NilVal)) != mm.False {
				return mm.True
			} else {
				return mm.ApplySys(mm.NewCompiled(0, false, func(a []*mm.Val) *mm.Val {
					env := mm.ExtendEnv(env, mm.NilVal, mm.SliceToList(a))
					env = env
					mm.Define(env, mm.SymbolFromRawString("result"), mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("eval-seq")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("extend-env")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("env")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("nil")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("nil")), mm.NilVal)))), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("exp")), mm.NilVal)), mm.NilVal))))
					return func() *mm.Val {
						if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("error?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("result")), mm.NilVal)) != mm.False {
							return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("fatal")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("result")), mm.NilVal))
						} else {
							return func() *mm.Val {
								if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("not")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("result")), mm.NilVal)) != mm.False {
									return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("fatal")), mm.Cons(mm.SysStringToString("test failed"), mm.NilVal))
								} else {
									return mm.SymbolFromRawString("test-complete")
								}
							}()
						}
					}()
				}), mm.NilVal)
			}
		}()
	}))
	mm.Define(env, mm.SymbolFromRawString("eval-export"), mm.NewCompiled(2, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("env"), mm.Cons(mm.SymbolFromRawString("exp"), mm.NilVal)), mm.SliceToList(a))
		env = env
		return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("module-export")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("env")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("list->struct")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("struct-values")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("env")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("exp")), mm.NilVal)), mm.NilVal))), mm.NilVal)), mm.NilVal)))
	}))
	mm.Define(env, mm.SymbolFromRawString("read-eval"), mm.NewCompiled(2, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("env"), mm.Cons(mm.SymbolFromRawString("r"), mm.NilVal)), mm.SliceToList(a))
		env = env
		return mm.ApplySys(mm.NewCompiled(0, false, func(a []*mm.Val) *mm.Val {
			env := mm.ExtendEnv(env, mm.NilVal, mm.SliceToList(a))
			env = env
			mm.Define(env, mm.SymbolFromRawString("r"), mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("read")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("r")), mm.NilVal)))
			return func() *mm.Val {
				if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("=")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("r:state")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("eof")), mm.NilVal))) != mm.False {
					return mm.LookupDef(env, mm.SymbolFromRawString("env"))
				} else {
					return func() *mm.Val {
						if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("error?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("r:state")), mm.NilVal)) != mm.False {
							return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("fatal")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("r:state")), mm.NilVal))
						} else {
							return func() *mm.Val {
								mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("eval-env")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("env")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("r:state")), mm.NilVal)))
								return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("read-eval")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("env")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("r")), mm.NilVal)))
							}()
						}
					}()
				}
			}()
		}), mm.NilVal)
	}))
	mm.Define(env, mm.SymbolFromRawString("load-module"), mm.NewCompiled(2, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("env"), mm.Cons(mm.SymbolFromRawString("module-name"), mm.NilVal)), mm.SliceToList(a))
		env = env
		return mm.ApplySys(mm.NewCompiled(0, false, func(a []*mm.Val) *mm.Val {
			env := mm.ExtendEnv(env, mm.NilVal, mm.SliceToList(a))
			env = env
			mm.Define(env, mm.SymbolFromRawString("f"), mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("fopen")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("module-name")), mm.NilVal)))
			return func() *mm.Val {
				if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("error?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("f")), mm.NilVal)) != mm.False {
					return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("fatal")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("f")), mm.NilVal))
				} else {
					return mm.ApplySys(mm.NewCompiled(0, false, func(a []*mm.Val) *mm.Val {
						env := mm.ExtendEnv(env, mm.NilVal, mm.SliceToList(a))
						env = env
						mm.Define(env, mm.SymbolFromRawString("r"), mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("reader")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("f")), mm.NilVal)))
						mm.Define(env, mm.SymbolFromRawString("menv"), mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("read-eval")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("module-env")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("env")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("module-name")), mm.NilVal))), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("r")), mm.NilVal))))
						mm.Define(env, mm.SymbolFromRawString("exp"), mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("exports")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("menv")), mm.NilVal)))
						return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("store-module")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("env")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("module-name")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("exp")), mm.NilVal))))
					}), mm.NilVal)
				}
			}()
		}), mm.NilVal)
	}))
	mm.Define(env, mm.SymbolFromRawString("import-def"), mm.NewCompiled(1, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("exp"), mm.NilVal), mm.SliceToList(a))
		env = env
		return func() *mm.Val {
			if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("=")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("len")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("exp")), mm.NilVal)), mm.Cons(mm.SysIntToNumber(2), mm.NilVal))) != mm.False {
				return mm.ListToStruct(mm.Cons(mm.SymbolFromRawString("import-name"), mm.Cons(mm.False, mm.Cons(mm.SymbolFromRawString("module-name"), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("exp")), mm.NilVal)), mm.NilVal)), mm.NilVal)))))
			} else {
				return func() *mm.Val {
					if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("=")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("len")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("exp")), mm.NilVal)), mm.Cons(mm.SysIntToNumber(3), mm.NilVal))) != mm.False {
						return mm.ListToStruct(mm.Cons(mm.SymbolFromRawString("import-name"), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("exp")), mm.NilVal)), mm.NilVal)), mm.Cons(mm.SymbolFromRawString("module-name"), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("exp")), mm.NilVal)), mm.NilVal)), mm.NilVal)), mm.NilVal)))))
					} else {
						return mm.LookupDef(env, mm.SymbolFromRawString("invalid-import"))
					}
				}()
			}
		}()
	}))
	mm.Define(env, mm.SymbolFromRawString("eval-import"), mm.NewCompiled(2, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("env"), mm.Cons(mm.SymbolFromRawString("exp"), mm.NilVal)), mm.SliceToList(a))
		env = env
		mm.Define(env, mm.SymbolFromRawString("define-import"), mm.NewCompiled(2, false, func(a []*mm.Val) *mm.Val {
			env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("n"), mm.Cons(mm.SymbolFromRawString("m"), mm.NilVal)), mm.SliceToList(a))
			env = env
			return func() *mm.Val {
				if mm.LookupDef(env, mm.SymbolFromRawString("n")) != mm.False {
					return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("define")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("env")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("n")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("m")), mm.NilVal))))
				} else {
					return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("fold")), mm.Cons(mm.NewCompiled(2, false, func(a []*mm.Val) *mm.Val {
						env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("n"), mm.Cons(mm.SymbolFromRawString("m"), mm.NilVal)), mm.SliceToList(a))
						env = env
						mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("define")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("env")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("n")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("m")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("n")), mm.NilVal)), mm.NilVal))))
						return mm.LookupDef(env, mm.SymbolFromRawString("m"))
					}), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("m")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("struct-names")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("m")), mm.NilVal)), mm.NilVal))))
				}
			}()
		}))
		return mm.ApplySys(mm.NewCompiled(0, false, func(a []*mm.Val) *mm.Val {
			env := mm.ExtendEnv(env, mm.NilVal, mm.SliceToList(a))
			env = env
			mm.Define(env, mm.SymbolFromRawString("i"), mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("import-def")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("exp")), mm.NilVal)))
			mm.Define(env, mm.SymbolFromRawString("current-import-path"), mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("module-path")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("env")), mm.NilVal)))
			return func() *mm.Val {
				if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("error?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("i")), mm.NilVal)) != mm.False {
					return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("fatal")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("i")), mm.NilVal))
				} else {
					return func() *mm.Val {
						if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("memq")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("i:module-name")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("current-import-path")), mm.NilVal))) != mm.False {
							return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("fatal")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("circular-import")), mm.NilVal))
						} else {
							return mm.ApplySys(mm.NewCompiled(0, false, func(a []*mm.Val) *mm.Val {
								env := mm.ExtendEnv(env, mm.NilVal, mm.SliceToList(a))
								env = env
								mm.Define(env, mm.SymbolFromRawString("module"), mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("loaded-module")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("env")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("i:module-name")), mm.NilVal))))
								return func() *mm.Val {
									if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("=")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("module")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("undefined-module")), mm.NilVal))) != mm.False {
										return mm.ApplySys(mm.NewCompiled(0, false, func(a []*mm.Val) *mm.Val {
											env := mm.ExtendEnv(env, mm.NilVal, mm.SliceToList(a))
											env = env
											mm.Define(env, mm.SymbolFromRawString("module"), mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("load-module")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("env")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("i:module-name")), mm.NilVal))))
											return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("define-import")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("i:import-name")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("module")), mm.NilVal)))
										}), mm.NilVal)
									} else {
										return func() *mm.Val {
											if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("error?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("module")), mm.NilVal)) != mm.False {
												return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("fatal")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("module")), mm.NilVal))
											} else {
												return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("define-import")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("i:import-name")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("module")), mm.NilVal)))
											}
										}()
									}
								}()
							}), mm.NilVal)
						}
					}()
				}
			}()
		}), mm.NilVal)
	}))
	mm.Define(env, mm.SymbolFromRawString("eval-apply"), mm.NewCompiled(2, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("env"), mm.Cons(mm.SymbolFromRawString("exp"), mm.NilVal)), mm.SliceToList(a))
		env = env
		return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("apply")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("eval-exp")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("env")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("exp")), mm.NilVal)), mm.NilVal))), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("value-list")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("env")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("exp")), mm.NilVal)), mm.NilVal))), mm.NilVal)))
	}))
	mm.Define(env, mm.SymbolFromRawString("eval-env"), mm.NewCompiled(2, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("env"), mm.Cons(mm.SymbolFromRawString("exp"), mm.NilVal)), mm.SliceToList(a))
		env = env
		return func() *mm.Val {
			if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("number?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("exp")), mm.NilVal)) != mm.False {
				return mm.LookupDef(env, mm.SymbolFromRawString("exp"))
			} else {
				return func() *mm.Val {
					if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("string?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("exp")), mm.NilVal)) != mm.False {
						return mm.LookupDef(env, mm.SymbolFromRawString("exp"))
					} else {
						return func() *mm.Val {
							if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("bool?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("exp")), mm.NilVal)) != mm.False {
								return mm.LookupDef(env, mm.SymbolFromRawString("exp"))
							} else {
								return func() *mm.Val {
									if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("nil?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("exp")), mm.NilVal)) != mm.False {
										return mm.LookupDef(env, mm.SymbolFromRawString("exp"))
									} else {
										return func() *mm.Val {
											if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("vector-form?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("exp")), mm.NilVal)) != mm.False {
												return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("eval-vector")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("env")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("exp")), mm.NilVal)))
											} else {
												return func() *mm.Val {
													if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("struct-form?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("exp")), mm.NilVal)) != mm.False {
														return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("eval-struct")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("env")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("exp")), mm.NilVal)))
													} else {
														return func() *mm.Val {
															if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("quote?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("exp")), mm.NilVal)) != mm.False {
																return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("eval-quote")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("exp")), mm.NilVal))
															} else {
																return func() *mm.Val {
																	if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("symbol?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("exp")), mm.NilVal)) != mm.False {
																		return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("lookup-def")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("env")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("exp")), mm.NilVal)))
																	} else {
																		return func() *mm.Val {
																			if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("def?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("exp")), mm.NilVal)) != mm.False {
																				return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("eval-def")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("env")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("exp")), mm.NilVal)))
																			} else {
																				return func() *mm.Val {
																					if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("if?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("exp")), mm.NilVal)) != mm.False {
																						return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("eval-if")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("env")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("exp")), mm.NilVal)))
																					} else {
																						return func() *mm.Val {
																							if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("and?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("exp")), mm.NilVal)) != mm.False {
																								return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("eval-and")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("env")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("exp")), mm.NilVal)), mm.NilVal)))
																							} else {
																								return func() *mm.Val {
																									if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("or?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("exp")), mm.NilVal)) != mm.False {
																										return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("eval-or")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("env")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("exp")), mm.NilVal)), mm.NilVal)))
																									} else {
																										return func() *mm.Val {
																											if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("fn?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("exp")), mm.NilVal)) != mm.False {
																												return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("eval-fn")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("env")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("exp")), mm.NilVal)))
																											} else {
																												return func() *mm.Val {
																													if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("begin?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("exp")), mm.NilVal)) != mm.False {
																														return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("eval-seq")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("env")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("exp")), mm.NilVal)), mm.NilVal)))
																													} else {
																														return func() *mm.Val {
																															if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cond?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("exp")), mm.NilVal)) != mm.False {
																																return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("eval-env")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("env")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cond->if")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("exp")), mm.NilVal)), mm.NilVal)))
																															} else {
																																return func() *mm.Val {
																																	if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("let?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("exp")), mm.NilVal)) != mm.False {
																																		return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("eval-env")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("env")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("list")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("make-fn")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("nil")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("let-body")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("exp")), mm.NilVal)), mm.NilVal))), mm.NilVal)), mm.NilVal)))
																																	} else {
																																		return func() *mm.Val {
																																			if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("export?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("exp")), mm.NilVal)) != mm.False {
																																				return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("eval-export")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("env")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("exp")), mm.NilVal)))
																																			} else {
																																				return func() *mm.Val {
																																					if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("import?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("exp")), mm.NilVal)) != mm.False {
																																						return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("eval-import")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("env")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("exp")), mm.NilVal)))
																																					} else {
																																						return func() *mm.Val {
																																							if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("test?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("exp")), mm.NilVal)) != mm.False {
																																								return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("eval-test")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("env")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("exp")), mm.NilVal)))
																																							} else {
																																								return func() *mm.Val {
																																									if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("application?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("exp")), mm.NilVal)) != mm.False {
																																										return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("eval-apply")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("env")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("exp")), mm.NilVal)))
																																									} else {
																																										return mm.LookupDef(env, mm.SymbolFromRawString("invalid-expression"))
																																									}
																																								}()
																																							}
																																						}()
																																					}
																																				}()
																																			}
																																		}()
																																	}
																																}()
																															}
																														}()
																													}
																												}()
																											}
																										}()
																									}
																								}()
																							}
																						}()
																					}
																				}()
																			}
																		}()
																	}
																}()
															}
														}()
													}
												}()
											}
										}()
									}
								}()
							}
						}()
					}
				}()
			}
		}()
	}))
	mm.Define(env, mm.SymbolFromRawString("printer"), mm.NewCompiled(1, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("output"), mm.NilVal), mm.SliceToList(a))
		env = env
		return mm.ListToStruct(mm.Cons(mm.SymbolFromRawString("output"), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("output")), mm.Cons(mm.SymbolFromRawString("state"), mm.Cons(mm.SymbolFromRawString("ok"), mm.NilVal)))))
	}))
	mm.Define(env, mm.SymbolFromRawString("print-raw"), mm.NewCompiled(2, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("p"), mm.Cons(mm.SymbolFromRawString("s"), mm.NilVal)), mm.SliceToList(a))
		env = env
		return func() *mm.Val {
			if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("error?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("p:state")), mm.NilVal)) != mm.False {
				return mm.LookupDef(env, mm.SymbolFromRawString("p"))
			} else {
				return mm.ApplySys(mm.NewCompiled(0, false, func(a []*mm.Val) *mm.Val {
					env := mm.ExtendEnv(env, mm.NilVal, mm.SliceToList(a))
					env = env
					mm.Define(env, mm.SymbolFromRawString("output"), mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("fwrite")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("p:output")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("s")), mm.NilVal))))
					return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("assign")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("p")), mm.Cons(mm.ListToStruct(mm.Cons(mm.SymbolFromRawString("output"), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("output")), mm.Cons(mm.SymbolFromRawString("state"), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("fstate")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("output")), mm.NilVal)), mm.NilVal))))), mm.NilVal)))
				}), mm.NilVal)
			}
		}()
	}))
	mm.Define(env, mm.SymbolFromRawString("print-quote-sign"), mm.NewCompiled(1, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("p"), mm.NilVal), mm.SliceToList(a))
		env = env
		return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("print-raw")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("p")), mm.Cons(mm.SysStringToString("'"), mm.NilVal)))
	}))
	mm.Define(env, mm.SymbolFromRawString("print-symbol"), mm.NewCompiled(3, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("p"), mm.Cons(mm.SymbolFromRawString("v"), mm.Cons(mm.SymbolFromRawString("quoted?"), mm.NilVal))), mm.SliceToList(a))
		env = env
		return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("print-raw")), mm.Cons(func() *mm.Val {
			if mm.LookupDef(env, mm.SymbolFromRawString("quoted?")) != mm.False {
				return mm.LookupDef(env, mm.SymbolFromRawString("p"))
			} else {
				return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("print-quote-sign")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("p")), mm.NilVal))
			}
		}(), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("symbol->string")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.NilVal)))
	}))
	mm.Define(env, mm.SymbolFromRawString("print-quote"), mm.NewCompiled(2, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("p"), mm.Cons(mm.SymbolFromRawString("v"), mm.NilVal)), mm.SliceToList(a))
		env = env
		return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("printq")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("print-quote-sign")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("p")), mm.NilVal)), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.NilVal)), mm.Cons(mm.False, mm.NilVal))))
	}))
	mm.Define(env, mm.SymbolFromRawString("print-pair"), mm.NewCompiled(3, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("p"), mm.Cons(mm.SymbolFromRawString("v"), mm.Cons(mm.SymbolFromRawString("quoted?"), mm.NilVal))), mm.SliceToList(a))
		env = env
		mm.Define(env, mm.SymbolFromRawString("print-space"), mm.NewCompiled(2, false, func(a []*mm.Val) *mm.Val {
			env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("p"), mm.Cons(mm.SymbolFromRawString("v"), mm.NilVal)), mm.SliceToList(a))
			env = env
			return func() *mm.Val {
				if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("nil?")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.NilVal)) != mm.False {
					return mm.LookupDef(env, mm.SymbolFromRawString("p"))
				} else {
					return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("print-raw")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("p")), mm.Cons(mm.SysStringToString(" "), mm.NilVal)))
				}
			}()
		}))
		mm.Define(env, mm.SymbolFromRawString("print-items"), mm.NewCompiled(2, false, func(a []*mm.Val) *mm.Val {
			env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("p"), mm.Cons(mm.SymbolFromRawString("v"), mm.NilVal)), mm.SliceToList(a))
			env = env
			return func() *mm.Val {
				if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("nil?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)) != mm.False {
					return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("print-raw")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("p")), mm.Cons(mm.SysStringToString(")"), mm.NilVal)))
				} else {
					return func() *mm.Val {
						if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("not")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("pair?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.NilVal)) != mm.False {
							return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("print-raw")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("printq")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("print-raw")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("p")), mm.Cons(mm.SysStringToString(". "), mm.NilVal))), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.Cons(mm.True, mm.NilVal)))), mm.Cons(mm.SysStringToString(")"), mm.NilVal)))
						} else {
							return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("print-items")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("print-space")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("printq")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("p")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.Cons(mm.True, mm.NilVal)))), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal))), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.NilVal)))
						}
					}()
				}
			}()
		}))
		return mm.ApplySys(mm.NewCompiled(0, false, func(a []*mm.Val) *mm.Val {
			env := mm.ExtendEnv(env, mm.NilVal, mm.SliceToList(a))
			env = env
			mm.Define(env, mm.SymbolFromRawString("p"), mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("print-raw")), mm.Cons(func() *mm.Val {
				if mm.LookupDef(env, mm.SymbolFromRawString("quoted?")) != mm.False {
					return mm.LookupDef(env, mm.SymbolFromRawString("p"))
				} else {
					return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("print-quote-sign")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("p")), mm.NilVal))
				}
			}(), mm.Cons(mm.SysStringToString("("), mm.NilVal))))
			return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("print-items")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("p")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)))
		}), mm.NilVal)
	}))
	mm.Define(env, mm.SymbolFromRawString("print-vector"), mm.NewCompiled(2, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("p"), mm.Cons(mm.SymbolFromRawString("v"), mm.NilVal)), mm.SliceToList(a))
		env = env
		mm.Define(env, mm.SymbolFromRawString("print-space"), mm.NewCompiled(2, false, func(a []*mm.Val) *mm.Val {
			env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("p"), mm.Cons(mm.SymbolFromRawString("i"), mm.NilVal)), mm.SliceToList(a))
			env = env
			return func() *mm.Val {
				if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString(">=")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("i")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("-")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("len")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.Cons(mm.SysIntToNumber(1), mm.NilVal))), mm.NilVal))) != mm.False {
					return mm.LookupDef(env, mm.SymbolFromRawString("p"))
				} else {
					return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("print-raw")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("p")), mm.Cons(mm.SysStringToString(" "), mm.NilVal)))
				}
			}()
		}))
		mm.Define(env, mm.SymbolFromRawString("print-items"), mm.NewCompiled(2, false, func(a []*mm.Val) *mm.Val {
			env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("p"), mm.Cons(mm.SymbolFromRawString("i"), mm.NilVal)), mm.SliceToList(a))
			env = env
			return func() *mm.Val {
				if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("=")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("i")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("len")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.NilVal))) != mm.False {
					return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("print-raw")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("p")), mm.Cons(mm.SysStringToString("]"), mm.NilVal)))
				} else {
					return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("print-items")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("print-space")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("printq")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("p")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("vector-ref")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("i")), mm.NilVal))), mm.Cons(mm.True, mm.NilVal)))), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("i")), mm.NilVal))), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("inc")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("i")), mm.NilVal)), mm.NilVal)))
				}
			}()
		}))
		return mm.ApplySys(mm.NewCompiled(0, false, func(a []*mm.Val) *mm.Val {
			env := mm.ExtendEnv(env, mm.NilVal, mm.SliceToList(a))
			env = env
			mm.Define(env, mm.SymbolFromRawString("p"), mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("print-raw")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("p")), mm.Cons(mm.SysStringToString("["), mm.NilVal))))
			return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("print-items")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("p")), mm.Cons(mm.SysIntToNumber(0), mm.NilVal)))
		}), mm.NilVal)
	}))
	mm.Define(env, mm.SymbolFromRawString("print-struct"), mm.NewCompiled(2, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("p"), mm.Cons(mm.SymbolFromRawString("s"), mm.NilVal)), mm.SliceToList(a))
		env = env
		mm.Define(env, mm.SymbolFromRawString("print-space"), mm.NewCompiled(2, false, func(a []*mm.Val) *mm.Val {
			env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("p"), mm.Cons(mm.SymbolFromRawString("n"), mm.NilVal)), mm.SliceToList(a))
			env = env
			return func() *mm.Val {
				if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("nil?")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("n")), mm.NilVal)), mm.NilVal)) != mm.False {
					return mm.LookupDef(env, mm.SymbolFromRawString("p"))
				} else {
					return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("print-raw")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("p")), mm.Cons(mm.SysStringToString(" "), mm.NilVal)))
				}
			}()
		}))
		mm.Define(env, mm.SymbolFromRawString("print-items"), mm.NewCompiled(2, false, func(a []*mm.Val) *mm.Val {
			env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("p"), mm.Cons(mm.SymbolFromRawString("n"), mm.NilVal)), mm.SliceToList(a))
			env = env
			return func() *mm.Val {
				if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("nil?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("n")), mm.NilVal)) != mm.False {
					return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("print-raw")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("p")), mm.Cons(mm.SysStringToString("}"), mm.NilVal)))
				} else {
					return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("print-items")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("print-space")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("printq")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("print-raw")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("printq")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("p")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("n")), mm.NilVal)), mm.Cons(mm.True, mm.NilVal)))), mm.Cons(mm.SysStringToString(" "), mm.NilVal))), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("field")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("s")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("car")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("n")), mm.NilVal)), mm.NilVal))), mm.Cons(mm.True, mm.NilVal)))), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("n")), mm.NilVal))), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("cdr")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("n")), mm.NilVal)), mm.NilVal)))
				}
			}()
		}))
		return mm.ApplySys(mm.NewCompiled(0, false, func(a []*mm.Val) *mm.Val {
			env := mm.ExtendEnv(env, mm.NilVal, mm.SliceToList(a))
			env = env
			mm.Define(env, mm.SymbolFromRawString("p"), mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("print-raw")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("p")), mm.Cons(mm.SysStringToString("{"), mm.NilVal))))
			return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("print-items")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("p")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("struct-names")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("s")), mm.NilVal)), mm.NilVal)))
		}), mm.NilVal)
	}))
	mm.Define(env, mm.SymbolFromRawString("printq"), mm.NewCompiled(3, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("p"), mm.Cons(mm.SymbolFromRawString("v"), mm.Cons(mm.SymbolFromRawString("quoted?"), mm.NilVal))), mm.SliceToList(a))
		env = env
		return func() *mm.Val {
			if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("number?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)) != mm.False {
				return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("print-raw")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("p")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("number->string")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.NilVal)))
			} else {
				return func() *mm.Val {
					if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("string?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)) != mm.False {
						return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("print-raw")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("p")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)))
					} else {
						return func() *mm.Val {
							if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("bool?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)) != mm.False {
								return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("print-raw")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("p")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("bool->string")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.NilVal)))
							} else {
								return func() *mm.Val {
									if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("symbol?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)) != mm.False {
										return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("print-symbol")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("p")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("quoted?")), mm.NilVal))))
									} else {
										return func() *mm.Val {
											if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("sys?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)) != mm.False {
												return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("print-raw")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("p")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("sys->string")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.NilVal)))
											} else {
												return func() *mm.Val {
													if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("error?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)) != mm.False {
														return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("print-raw")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("p")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("error->string")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.NilVal)))
													} else {
														return func() *mm.Val {
															if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("quote?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)) != mm.False {
																return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("print-quote")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("p")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)))
															} else {
																return func() *mm.Val {
																	if func() *mm.Val {
																		if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("pair?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)) != mm.False {
																			return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("pair?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal))
																		} else {
																			return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("nil?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal))
																		}
																	}() != mm.False {
																		return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("print-pair")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("p")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("quoted?")), mm.NilVal))))
																	} else {
																		return func() *mm.Val {
																			if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("vector?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)) != mm.False {
																				return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("print-vector")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("p")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)))
																			} else {
																				return func() *mm.Val {
																					if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("struct?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)) != mm.False {
																						return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("print-struct")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("p")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)))
																					} else {
																						return func() *mm.Val {
																							if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("env?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)) != mm.False {
																								return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("print-raw")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("p")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("env->string")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.NilVal)))
																							} else {
																								return func() *mm.Val {
																									if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("function?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)) != mm.False {
																										return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("print-raw")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("p")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("function->string")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)), mm.NilVal)))
																									} else {
																										return mm.LookupDef(env, mm.SymbolFromRawString("not-implemented"))
																									}
																								}()
																							}
																						}()
																					}
																				}()
																			}
																		}()
																	}
																}()
															}
														}()
													}
												}()
											}
										}()
									}
								}()
							}
						}()
					}
				}()
			}
		}()
	}))
	mm.Define(env, mm.SymbolFromRawString("print"), mm.NewCompiled(2, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("p"), mm.Cons(mm.SymbolFromRawString("v"), mm.NilVal)), mm.SliceToList(a))
		env = env
		return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("printq")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("p")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.Cons(mm.False, mm.NilVal))))
	}))
	mm.Define(env, mm.SymbolFromRawString("read-eval-print"), mm.NewCompiled(0, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.NilVal, mm.SliceToList(a))
		env = env
		mm.Define(env, mm.SymbolFromRawString("loop"), mm.NewCompiled(3, false, func(a []*mm.Val) *mm.Val {
			env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("env"), mm.Cons(mm.SymbolFromRawString("r"), mm.Cons(mm.SymbolFromRawString("p"), mm.NilVal))), mm.SliceToList(a))
			env = env
			return mm.ApplySys(mm.NewCompiled(0, false, func(a []*mm.Val) *mm.Val {
				env := mm.ExtendEnv(env, mm.NilVal, mm.SliceToList(a))
				env = env
				mm.Define(env, mm.SymbolFromRawString("r"), mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("read")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("r")), mm.NilVal)))
				return func() *mm.Val {
					if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("=")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("r:state")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("eof")), mm.NilVal))) != mm.False {
						return mm.SymbolFromRawString("ok")
					} else {
						return func() *mm.Val {
							if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("error?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("r:state")), mm.NilVal)) != mm.False {
								return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("fatal")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("r:state")), mm.NilVal))
							} else {
								return mm.ApplySys(mm.NewCompiled(0, false, func(a []*mm.Val) *mm.Val {
									env := mm.ExtendEnv(env, mm.NilVal, mm.SliceToList(a))
									env = env
									mm.Define(env, mm.SymbolFromRawString("v"), mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("eval-env")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("env")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("r:state")), mm.NilVal))))
									return func() *mm.Val {
										if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("error?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal)) != mm.False {
											return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("fatal")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal))
										} else {
											return mm.ApplySys(mm.NewCompiled(0, false, func(a []*mm.Val) *mm.Val {
												env := mm.ExtendEnv(env, mm.NilVal, mm.SliceToList(a))
												env = env
												mm.Define(env, mm.SymbolFromRawString("p"), mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("print")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("p")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("v")), mm.NilVal))))
												return func() *mm.Val {
													if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("error?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("p:state")), mm.NilVal)) != mm.False {
														return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("fatal")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("p:state")), mm.NilVal))
													} else {
														return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("loop")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("env")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("r")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("print-raw")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("p")), mm.Cons(mm.SysStringToString("\n"), mm.NilVal))), mm.NilVal))))
													}
												}()
											}), mm.NilVal)
										}
									}()
								}), mm.NilVal)
							}
						}()
					}
				}()
			}), mm.NilVal)
		}))
		return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("loop")), mm.Cons(func() *mm.Val { return env }(), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("reader")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("stdin")), mm.NilVal), mm.NilVal)), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("printer")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("stdout")), mm.NilVal), mm.NilVal)), mm.NilVal))))
	}))
	mm.Define(env, mm.SymbolFromRawString("read-compile-write"), mm.NewCompiled(0, false, func(a []*mm.Val) *mm.Val {
		env := mm.ExtendEnv(env, mm.NilVal, mm.SliceToList(a))
		env = env
		mm.Define(env, mm.SymbolFromRawString("compile-top"), mm.NewCompiled(2, false, func(a []*mm.Val) *mm.Val {
			env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("c"), mm.Cons(mm.SymbolFromRawString("exp"), mm.NilVal)), mm.SliceToList(a))
			env = env
			return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compiler-compose")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compile")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("exp")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compiler-append")), mm.Cons(mm.SysStringToString(";"), mm.NilVal))))))
		}))
		mm.Define(env, mm.SymbolFromRawString("loop"), mm.NewCompiled(2, false, func(a []*mm.Val) *mm.Val {
			env := mm.ExtendEnv(env, mm.Cons(mm.SymbolFromRawString("r"), mm.Cons(mm.SymbolFromRawString("c"), mm.NilVal)), mm.SliceToList(a))
			env = env
			return mm.ApplySys(mm.NewCompiled(0, false, func(a []*mm.Val) *mm.Val {
				env := mm.ExtendEnv(env, mm.NilVal, mm.SliceToList(a))
				env = env
				mm.Define(env, mm.SymbolFromRawString("r"), mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("read")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("r")), mm.NilVal)))
				return func() *mm.Val {
					if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("error?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("r:state")), mm.NilVal)) != mm.False {
						return mm.ListToStruct(mm.Cons(mm.SymbolFromRawString("input"), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("r:input")), mm.Cons(mm.SymbolFromRawString("output"), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c:output")), mm.Cons(mm.SymbolFromRawString("read-error"), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("r:state")), mm.Cons(mm.SymbolFromRawString("compile-error"), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c:error")), mm.NilVal)))))))))
					} else {
						return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("loop")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("r")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compile-top")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("c")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("r:state")), mm.NilVal))), mm.NilVal)))
					}
				}()
			}), mm.NilVal)
		}))
		return mm.ApplySys(mm.NewCompiled(0, false, func(a []*mm.Val) *mm.Val {
			env := mm.ExtendEnv(env, mm.NilVal, mm.SliceToList(a))
			env = env
			mm.Define(env, mm.SymbolFromRawString("r"), mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("loop")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("reader")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("stdin")), mm.NilVal), mm.NilVal)), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("compiler")), mm.NilVal), mm.NilVal))))
			return func() *mm.Val {
				if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("error?")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("r:compile-error")), mm.NilVal)) != mm.False {
					return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("fatal")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("r:compile-error")), mm.NilVal))
				} else {
					return func() *mm.Val {
						if mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("=")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("r:read-error")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("eof")), mm.NilVal))) != mm.False {
							return func() *mm.Val {
								mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("fclose")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("r:input")), mm.NilVal))
								mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("fwrite")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("stdout")), mm.NilVal), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compiled-head")), mm.NilVal)))
								mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("fwrite")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("stdout")), mm.NilVal), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("r:output")), mm.NilVal)))
								mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("fwrite")), mm.Cons(mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("stdout")), mm.NilVal), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("compiled-tail")), mm.NilVal)))
								return mm.SymbolFromRawString("ok")
							}()
						} else {
							return mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("fatal")), mm.Cons(mm.LookupDef(env, mm.SymbolFromRawString("r:read-error")), mm.NilVal))
						}
					}()
				}
			}()
		}), mm.NilVal)
	}))
	mm.ApplySys(mm.LookupDef(env, mm.SymbolFromRawString("read-compile-write")), mm.NilVal)
}
