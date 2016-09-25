(import "scm/lang.scm")
(import "scm/code.scm")
(import "scm/read.scm")


(def (compiler)
  {output ""
   error false
   compiled-modules nil
   current-import-path nil})


(def (compiler-append c . a)
  (cond ((nil? a) c)
        (c:error c)
        (else
          (apply
            compiler-append
            (cons
              (assign c {output (string-append c:output (car a))})
              (cdr a))))))


(def (compiler-error c e) (assign c {error e}))


(def (compiler-compose c . i)
  ((apply
     compose
     (foldr
       (fn (i p)
         (cond ((function? i)
                (cons (partr i (car p)) (cdr p)))
               (else (cons i p))))
       nil i))
   c))


(def (compile-number c v) (compiler-append c "mm.SysIntToNumber(" (number->string v) ")"))


(def (compile-string c v) (compiler-append c "mm.SysStringToString(" (escape-compiled-string v) ")"))


(def (compile-bool c v) (compiler-append c (if v "mm.True" "mm.False")))


(def (compile-nil c v) (compiler-append c "mm.NilVal"))


(def (compile-symbol c v)
  (compiler-append c
                   "mm.SymbolFromRawString("
                   (escape-compiled-string (symbol->string v))
                   ")"))


(def (compile-quote-literal c v)
  (compiler-compose
    c
    compiler-append "mm.Cons("
    compile-symbol 'quote
    compiler-append ", mm.Cons("
    compile-literal (car (cdr v))
    compiler-append ", mm.NilVal))"))


(def (compile-quote c v) (compile-literal c (car (cdr v))))


(def (compile-pair-literal c v)
  (compiler-compose
    c
    compiler-append "mm.Cons("
    compile-literal (car v)
    compiler-append ", "
    compile-literal (cdr v)
    compiler-append ")"))


(def (compile-def c v)
  (cond ((and (value-def? v) (not (valid-value-def? v)))
         (compiler-error c invalid-def))
        ((and (not (value-def? v)) (not (valid-function-def? v)))
         (compiler-error c invalid-def))
        (else
          (compiler-compose
            c
            compiler-append "mm.Define(env, "
            compile-literal (def-name v)
            compiler-append ", "
            compile-exp (def-value v)
            compiler-append ")"))))


(def (compile-seq c v)
  (cond ((not (pair? v))
         (compiler-error c invalid-seq))
        ((nil? (cdr v))
         (compiler-compose
           c
           compiler-append "return "
           compile-exp (car v)))
        (else (compiler-compose
                c
                compile-statement (car v)
                compiler-append ";\n"
                compile-seq (cdr v)))))


(def (compile-fn c v)
  (cond ((or (not (pair? v))
             (not (pair? (cdr v))))
         (compiler-error c invalid-fn))
        (else
          (let (signature (fn-signature (car (cdr v)))
                body      (cdr (cdr v)))
            (if (error? signature)
              (compiler-error c signature)
              (compiler-compose
                c
                compiler-append "mm.NewCompiled("
                compiler-append (number->string signature:count)
                compiler-append ", "
                compiler-append (bool->string signature:var?)
                compiler-append ", func(a []*mm.Val) *mm.Val { env := mm.ExtendEnv(env, "
                compile-literal signature:names
                compiler-append ", mm.SliceToList(a)); env = env; "
                compile-seq body
                compiler-append "})"))))))


(def (compile-if c v)
  (cond ((not (= (len v) 4)) (compiler-error c invalid-if))
        (else (compiler-compose
                c
                compiler-append " func() *mm.Val { if "
                compile-exp (car (cdr v))
                compiler-append " != mm.False { return "
                compile-exp (car (cdr (cdr v)))
                compiler-append " } else { return "
                compile-exp (car (cdr (cdr (cdr v))))
                compiler-append " }}() "))))


(def (compile-and c v) (compile c (and->if (cdr v))))


(def (compile-or c v) (compile c (or->if (cdr v))))


(def (compile-begin c v)
  (compiler-compose
    c
    compiler-append "func() *mm.Val {"
    compile-seq (cdr v)
    compiler-append "}()"))


(def (compile-cond c v)
  (let (vi (cond->if v))
    (if (error? vi) (compiler-error c vi) (compile c vi))))

(def (compile-let c v)
  (compile c (list (make-fn nil (let-body v)))))


(def (compile-test c v)
  (def (compile-test-seq c v)
     (if (nil? v)
      (compiler-compose
        c
        compiler-append "return "
        compile-literal 'test-complete)
      (compiler-compose
        c
        compiler-append "if result := func() *mm.Val { return "
        compile (car v)
        compiler-append "}(); result == mm.False { return mm.Fatal("
        compile-exp "test failed"
        compiler-append ") } else if mm.IsError(result) != mm.False "
        compiler-append " { return mm.Fatal(result) }; "
        compile-test-seq (cdr v))))
  (compiler-compose
    c
    compiler-append "func() *mm.Val {"
    compiler-append "env := mm.ExtendEnv(env, mm.NilVal, mm.NilVal); env = env;"
    compile-test-seq (cdr v)
    compiler-append "}()"))


(def (compile-export c v)
  (compiler-compose
    c
    compiler-append "mm.Export(env, "
    compile (cons 'struct: (cdr v))
    compiler-append ")"))


(def (read-compile c r)
  (let (r (read r))
    (cond ((= r:state eof) c)
          ((error? r:state) (compiler-error c r:state))
          (else
            (read-compile
              (compiler-compose
                c
                compile r:state
                compiler-append "\n")
              r)))))


(def (compile-module c module-name)
  (let (f (fopen module-name))
    (if (error? f)
      (compiler-error c f)
      (let (c (compiler-compose
                c
                compiler-append "mm.ModuleLoader(env, "
                compile module-name
                compiler-append ", func(eenv *mm.Val) { env := mm.ModuleEnv(eenv, "
                compile module-name
                compiler-append ");"
                read-compile (reader f)
                compiler-append "; mm.StoreModule(eenv, "
                compile module-name
                compiler-append ", mm.Exports(env)) });"))
        (assign
          c
          {compiled-modules
           (cons module-name c:compiled-modules)})))))


(def (compile-import c exp)
  (def (compile-module-define c module-name import-name)
    (if import-name
      (compiler-compose
        c
        compiler-append "mm.Define(env, "
        compile-literal import-name
        compiler-append ", m)")
      (compiler-append c "mm.DefineAll(env, m)")))
  (def (compile-load c module-name import-name)
       (compiler-compose
         c
         compiler-append "func() { m := mm.LoadedModule(env, "
         compile module-name
         compiler-append ");"
         compiler-append "if m == mm.UndefinedModule {"
         compiler-append "mm.LoadCompiledModule(initialEnv, "
         compile module-name
         compiler-append "); m = mm.LoadedModule(env, "
         compile module-name
         compiler-append ")};"
         (partr compile-module-define import-name) module-name
         compiler-append "}();"))
  (let (i (import-def exp))
    (cond ((error? i) (compiler-error c i))
          ((memq i:module-name c:current-import-path)
           (compiler-error c circular-import))
          ((memq i:module-name c:compiled-modules)
           (compile-load c i:module-name i:import-name))
          (else
            (let (cc (assign (compiler) {current-import-path (cons i:module-name c:current-import-path)
                                         compiled-modules c:compiled-modules})
                  cr (compile-module cc i:module-name)
                  c (assign c {error cr:error compiled-modules (cons i:module-name c:compiled-modules)}))
              (compiler-compose
                c
                compiler-append cr:output
                (partr compile-load i:import-name) i:module-name))))))


(def (compile-lookup c v)
  (compiler-compose
    c
    compiler-append "mm.LookupDef(env, "
    compile-symbol v
    compiler-append ")"))


(def (compile-vector c v)
  (compiler-compose
    c
    compiler-append "mm.ListToVector("
    compile-literal (cdr v)
    compiler-append ")"))


(def (compile-struct c v)
  (def (compile-struct-values c v)
    (cond ((nil? v)
           (compiler-append c "mm.NilVal"))
          (else
            (compiler-compose
              c
              compiler-append "mm.Cons("
              compile-literal (car v)
              compiler-append ", mm.Cons("
              compile-exp (car (cdr v))
              compiler-append ", "
              compile-struct-values (cdr (cdr v))
              compiler-append "))"))))
  (compiler-compose
    c
    compiler-append "mm.ListToStruct("
    compile-struct-values (cdr v)
    compiler-append ")"))


(def (compile-value-list c v)
  (cond ((nil? v) (compiler-append c "mm.NilVal"))
        ((not (pair? v)) (compiler-error c invalid-value-list))
        (else
          (compiler-compose
            c
            compiler-append "mm.Cons("
            compile-exp (car v)
            compiler-append ", "
            compile-value-list (cdr v)
            compiler-append ")"))))


(def (compile-application c v)
  (compiler-compose
    c
    compiler-append "mm.ApplySys("
    compile-exp (car v)
    compiler-append ", "
    compile-value-list (cdr v)
    compiler-append ")"))


(def (compile-current-env c)
  (compiler-append c "func() *mm.Val { return env }()"))


(def (compile-literal c v)
  (cond ((check-types
           v
           number? string? bool? nil?)
         (compile c v))
        ((quote? v) (compile-quote-literal c v))
        ((symbol? v) (compile-symbol c v))
        ((pair? v) (compile-pair-literal c v))
        (else (compiler-error c invalid-literal))))


(def (compile-exp c v)
  ; TODO: this is not valid, because application takes all pairs
  (if (check-types
        v
        number? string? bool? nil? quote? symbol?
        vector-form? struct-form? if? and? or? fn?
        begin? cond? let? current-env? application?)
    (compile c v)
    (compiler-error c invalid-expression)))


(def (compile-statement c v)
  (if (check-types
        v
        def? if? and? or? begin? cond? let?
        import? export? test? application?)
    (compile c v)
    (compiler-error c invalid-statement)))


(def (compile c v)
  (cond ((number? v) (compile-number c v))
        ((string? v) (compile-string c v))
        ((bool? v) (compile-bool c v))
        ((nil? v) (compile-nil c v))
        ((quote? v) (compile-quote c v))
        ((symbol? v) (compile-lookup c v))
        ((vector-form? v) (compile-vector c v))
        ((struct-form? v) (compile-struct c v))
        ((def? v) (compile-def c v))
        ((if? v) (compile-if c v))
        ((and? v) (compile-and c v))
        ((or? v) (compile-or c v))
        ((fn? v) (compile-fn c v))
        ((begin? v) (compile-begin c v))
        ((cond? v) (compile-cond c v))
        ((let? v) (compile-let c v))
        ((export? v) (compile-export c v))
        ((import? v) (compile-import c v))
        ((test? v) (compile-test c v))
        ((current-env? v) (compile-current-env c))
        ((application? v) (compile-application c v))
        (else (compiler-error c not-implemented))))


(def compiled-head
     "package main

     import mm \"github.com/aryszka/mikkamakka\"

     func main() {
          initialEnv := mm.InitialEnv()
         env := mm.ExtendEnv(initialEnv, mm.NilVal, mm.NilVal)
     ")


(def compiled-tail "}")


(def (program-compiler)
  (compiler-append (compiler) compiled-head))


(def (program-compile c exp)
  (compiler-compose
	c
	compile-statement exp
	compiler-append ";\n"))


(def (close-program-compiler c)
  (compiler-append c compiled-tail))


(export compiler program-compiler
        compile program-compile
		close-compiler close-program-compiler)
