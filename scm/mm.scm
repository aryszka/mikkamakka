(import "scm/lang.scm")


(def definition-expression (string->error "definition in expression position"))
(def invalid-expression (string->error "invalid expression"))
(def inalid-token (string->error "invalid-token"))
(def circular-import (string->error "circular-import"))
(def not-implemented (string->error "not implemented"))
(def invalid-literal (string->error "invalid literal"))
(def invalid-application (string->error "invalid application"))
(def invalid-value-list (string->error "invalid value list"))
(def invalid-import (string->error "invalid import expression"))


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


(def (make-fn args body) (cons 'fn (cons args body)))


(def (value-def? v) (symbol? (car (cdr v))))


; makes no sense, fails in take, use list? and len
(def (valid-def? v) (every? pair? (take 3 v)))


(def (valid-value-def? v)
  (and (valid-def? v)
       (symbol? (car (cdr v)))
       (nil? (drop 3 v))))


(def (valid-function-def? v)
  (and (valid-def? v)
       (pair? (car (cdr v)))
       (symbol? (car (car (cdr v))))))


(def (def-name v)
  (cond ((value-def? v) (car (cdr v)))
        (else (car (car (cdr v))))))


(def (def-value v)
  (cond ((value-def? v) (car (cdr (cdr v))))
        (else (make-fn (cdr (car (cdr v))) (cdr (cdr v))))))


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


(def (fn-signature v)
  (cond ((nil? v) {count 0 var? false names '()})
        ((symbol? v) {count 0 var? true names v})
        ((pair? v) (let (signature (fn-signature (cdr v)))
                     (assign signature {count (inc signature:count)
                                         names (cons (car v) signature:names)})))
        (else invalid-fn)))


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


(def (and->if v)
  (cond ((nil? v) true)
        ((nil? (cdr v)) (car v))
        (else (list 'if (car v) (and->if (cdr v)) false))))


(def (compile-and c v) (compile c (and->if (cdr v))))


(def (or->if v)
  (cond ((nil? v) false)
        ((nil? (cdr v)) (car v))
        (else (list 'if (car v) (car v) (or->if (cdr v))))))


(def (compile-or c v) (compile c (or->if (cdr v))))


(def (compile-begin c v)
  (compiler-compose
    c
    compiler-append "func() *mm.Val {"
    compile-seq (cdr v)
    compiler-append "}()"))


(def (cond->if v)
  (def (seq->exp v)
    (cond ((nil? (cdr v)) (car v))
          (else (cons 'begin v))))
  (def (expand v)
    (cond ((or (not (pair? v))
               (not (pair? (car v)))
               (not (pair? (cdr (car v)))))
           invalid-cond)
          ((= (car (car v)) 'else)
           (cond ((not (nil? (cdr v)))
                  invalid-cond)
                 (else (seq->exp (cdr (car v))))))
          (else
            (list 'if (car (car v))
                  (seq->exp (cdr (car v)))
                  (expand (cdr v))))))
  (expand (cdr v)))


(def (compile-cond c v)
  (let (vi (cond->if v))
    (if (error? vi) (compiler-error c vi) (compile c vi))))


(def (let-body v)
  (def (let-defs v)
    (cond ((nil? v) nil)
          ((or (not (pair? v)) (not (pair? (cdr v))))
           (fatal invalid-let))
          (else (cons (list 'def (car v) (car (cdr v)))
                      (let-defs (cdr (cdr v)))))))
  (cond ((or (not (pair? v))
             (not (pair? (cdr v)))
             (nil? (cdr (cdr v))))
         (fatal invalid-let))
        (else (append (let-defs (car (cdr v))) (cdr (cdr v))))))


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
         env := initialEnv
     ")


(def compiled-tail "}")


(import "scm/read.scm")


(def (tagged? v t) (and (pair? v) (= (car v) t)))


(def (quote? v) (tagged? v 'quote))
(def (def? v) (tagged? v 'def))
(def (vector-form? v) (tagged? v 'vector:))
(def (struct-form? v) (tagged? v 'struct:))
(def (if? v) (tagged? v 'if))
(def (and? v) (tagged? v 'and))
(def (or? v) (tagged? v 'or))
(def (fn? v) (tagged? v 'fn))
(def (begin? v) (tagged? v 'begin))
(def (cond? v) (tagged? v 'cond))
(def (let? v) (tagged? v 'let))
(def (test? v) (tagged? v 'test))
(def (export? v) (tagged? v 'export))
(def (import? v) (tagged? v 'import))
(def (application? v) (pair? v))
(def (current-env? v) (tagged? v 'current-env))


(def (eval-quote exp) (car (cdr exp)))


(def (eval-def env exp)
  (define env (def-name exp) (eval-exp env (def-value exp))))


(def (value-list env exp)
  (cond ((nil? exp) nil)
        (else (cons (eval-exp env (car exp)) (value-list env (cdr exp))))))


(def (eval-vector env exp)
  (list->vector (value-list env (cdr exp))))


(def (struct-values env exp)
  (cond ((nil? exp) nil)
        (else (cons (car exp)
                    (cons (eval-exp env (car (cdr exp)))
                          (struct-values env (cdr (cdr exp))))))))


(def (eval-struct env exp)
  (list->struct (struct-values env (cdr exp))))


(def (eval-exp env exp)
  (cond ((def? exp) (fatal definition-expression))
        (else (eval-env env exp))))


(def (eval-if env exp)
  (if (eval-exp env (car (cdr exp)))
    (eval-exp env (car (cdr (cdr exp))))
    (eval-exp env (car (cdr (cdr (cdr exp)))))))


(def (eval-and env exp)
  (cond ((nil? exp) true)
        ((nil? (cdr exp)) (eval-exp env (car exp)))
        (else (if (not (eval-exp env (car exp)))
                false
                (eval-and env (cdr exp))))))


(def (eval-or env exp)
  (cond ((nil? exp) false)
        ((nil? (cdr exp)) (eval-exp env (car exp)))
        (else (let (v (eval-exp env (car exp)))
                (if v v (eval-or env (cdr exp)))))))


(def (eval-fn env exp)
  (let (signature (fn-signature (car (cdr exp)))
        body      (cdr (cdr exp)))
    (make-composite (cons env (cons signature:names body)))))


(def (eval-seq env exp)
  (cond ((not (pair? exp)) (fatal invalid-sequence))
        ((nil? (cdr exp)) (eval-exp env (car exp)))
        (else (eval-env env (car exp))
              (eval-seq env (cdr exp)))))


(def (eval-test env exp)
  (cond ((nil? (cdr exp)) true)
        (else (let (result (eval-seq (extend-env env nil nil) (cdr exp)))
                (cond ((error? result) (fatal result))
                      ((not result) (fatal "test failed"))
                      (else 'test-complete))))))


(def (eval-export env exp)
  (module-export env (list->struct (struct-values env (cdr exp)))))


(def (read-eval env r)
  (let (r (read r))
    (cond ((= r:state eof) env)
          ((error? r:state) (fatal r:state))
          (else
            (eval-env env r:state)
            (read-eval env r)))))


(def (load-module env module-name)
  (let (f (fopen module-name))
    (if (error? f)
      (fatal f)
      (let (r (reader f)
            menv (read-eval (module-env env module-name) r)
            exp  (exports menv))
        (store-module env module-name exp)))))


(def (import-def exp)
  (cond ((= (len exp) 2) {import-name false module-name (car (cdr exp))})
        ((= (len exp) 3) {import-name (car (cdr exp)) module-name (car (cdr (cdr exp)))})
        (else invalid-import)))


(def (eval-import env exp)
  (def (define-import n m)
    (if n
      (define env n m)
      (fold (fn (n m) (define env n (m n)) m)
            m
            (struct-names m))))
  (let (i (import-def exp)
        current-import-path (module-path env))
    (cond ((error? i) (fatal i))
          ((memq i:module-name current-import-path)
           (fatal circular-import))
          (else (let (module (loaded-module env i:module-name))
                  (cond ((= module undefined-module)
                         (let (module (load-module env i:module-name))
                           (define-import i:import-name module)))
                        ((error? module) (fatal module))
                        (else (define-import i:import-name module))))))))


(def (eval-apply env exp)
  (apply (eval-exp env (car exp)) (value-list env (cdr exp))))


(def (eval-env env exp)
  (cond ((number? exp) exp)
        ((string? exp) exp)
        ((bool? exp) exp)
        ((nil? exp) exp)
        ((vector-form? exp) (eval-vector env exp))
        ((struct-form? exp) (eval-struct env exp))
        ((quote? exp) (eval-quote exp))
        ((symbol? exp) (lookup-def env exp))
        ((def? exp) (eval-def env exp))
        ((if? exp) (eval-if env exp))
        ((and? exp) (eval-and env (cdr exp)))
        ((or? exp) (eval-or env (cdr exp)))
        ((fn? exp) (eval-fn env exp))
        ((begin? exp) (eval-seq env (cdr exp)))
        ((cond? exp) (eval-env env (cond->if exp)))
        ((let? exp) (eval-env env (list (make-fn nil (let-body exp)))))
        ((export? exp) (eval-export env exp))
        ((import? exp) (eval-import env exp))
        ((test? exp) (eval-test env exp))
        ((application? exp) (eval-apply env exp))
        (else invalid-expression)))


(def (printer output)
  {output output state 'ok})


(def (print-raw p s)
  (if (error? p:state) p
    (let (output (fwrite p:output s))
      (assign p {output output state (fstate output)}))))


(def (print-quote-sign p) (print-raw p "'"))


(def (print-symbol p v quoted?)
  (print-raw
    (if quoted? p (print-quote-sign p))
    (symbol->string v)))


(def (print-quote p v)
  (printq (print-quote-sign p) (car (cdr v)) false))


(def (print-pair p v quoted?)
  (def (print-space p v)
    (if (nil? (cdr v)) p
      (print-raw p " ")))
  (def (print-items p v)
    (cond ((nil? v) (print-raw p ")"))
          ((not (pair? v)) (print-raw (printq (print-raw p ". ") v true) ")"))
          (else (print-items
                  (print-space
                    (printq p (car v) true)
                    v)
                  (cdr v)))))
  (let (p (print-raw
            (if quoted? p (print-quote-sign p)) "("))
    (print-items p v)))


(def (print-vector p v)
  (def (print-space p i)
    (if (>= i (- (len v) 1)) p
      (print-raw p " ")))
  (def (print-items p i)
    (cond ((= i (len v)) (print-raw p "]"))
          (else (print-items
                  (print-space
                    (printq p (vector-ref v i) true)
                    i)
                  (inc i)))))
  (let (p (print-raw p "["))
    (print-items p 0)))


(def (print-struct p s)
  (def (print-space p n)
    (if (nil? (cdr n)) p
      (print-raw p " ")))
  (def (print-items p n)
    (cond ((nil? n) (print-raw p "}"))
          (else (print-items
                  (print-space
                    (printq
                      (print-raw (printq p (car n) true) " ")
                      (field s (car n))
                      true)
                    n)
                  (cdr n)))))
  (let (p (print-raw p "{"))
    (print-items p (struct-names s))))


(def (printq p v quoted?)
  (cond ((number? v) (print-raw p (number->string v)))
        ((string? v) (print-raw p v))
        ((bool? v) (print-raw p (bool->string v)))
        ((symbol? v) (print-symbol p v quoted?))
        ((sys? v) (print-raw p (sys->string v)))
        ((error? v) (print-raw p (error->string v)))
        ((quote? v) (print-quote p v))
        ((or (pair? v) (nil? v)) (print-pair p v quoted?))
        ((vector? v) (print-vector p v))
        ((struct? v) (print-struct p v))
        ((env? v) (print-raw p (env->string v)))
        ((function? v) (print-raw p (function->string v)))
        (else not-implemented)))


(def (print p v) (printq p v false))


(def (read-eval-print)
  (def (loop env r p)
    (let (r (read r))
      (cond
        ((= r:state eof) 'ok)
        ((error? r:state) (fatal r:state))
        (else
          (let (v (eval-env env r:state))
            (cond
              ((error? v) (fatal v))
              (else
                (let (p (print p v))
                  (cond
                    ((error? p:state) (fatal p:state))
                    (else (loop env r (print-raw p "\n"))))))))))))
  (loop (current-env) (reader (stdin)) (printer (stdout))))


(def (read-compile-write)
  (def (compile-top c exp)
    (compiler-compose
      c
      compile exp
      compiler-append ";"))
  (def (loop r c)
    (let (r (read r))
      (cond ((error? r:state)
             {input r:input
              output c:output
              read-error r:state
              compile-error c:error})
            (else (loop r (compile-top c r:state))))))
  (let (r  (loop (reader (stdin)) (compiler)))
    (cond ((error? r:compile-error) (fatal r:compile-error))
          ((= r:read-error eof)
           (fclose r:input)
           (fwrite (stdout) compiled-head)
           (fwrite (stdout) r:output)
           (fwrite (stdout) compiled-tail)
           'ok)
          (else (fatal r:read-error)))))


(read-compile-write)
; (read-eval-print)
