(def definition-expression (string->error "definition in expression position"))
(def invalid-expression (string->error "invalid expression"))
(def inalid-token (string->error "invalid-token"))
(def circular-import (string->error "circular-import"))
(def not-implemented (string->error "not implemented"))
(def invalid-literal (string->error "invalid literal"))
(def invalid-application (string->error "invalid application"))
(def invalid-value-list (string->error "invalid value list"))
(def invalid-import (string->error "invalid import expression"))

(def (trace message . values)
  (def (trace out values)
    (let (out (print out (car values)))
      (cond ((nil? (cdr values))
             (fwrite out:output "\n")
             (car values))
            (else
              (trace
                (assign out {output (fwrite out:output " ")})
                (cdr values))))))
  (trace (printer (stderr)) (cons message values)))


(def (id x) x)


(def (list . x) x)


(def (apply f a)
  (cond ((vector? f) (vector-ref f (car a)))
        ((struct? f) (field f (car a)))
        ((compiled-function? f) (apply-compiled f a))
        (else (let (c (composite f))
                (eval-seq (extend-env (car c) (car (cdr c)) a) (cdr (cdr c)))))))


(def (call f . a) (apply f a))


(def (fold f i l)
  (cond ((nil? l) i)
        (else (fold f (f (car l) i) (cdr l)))))


(def (foldr f i l)
  (cond ((nil? l) i)
        (else (f (car l) (foldr f i (cdr l))))))


(def (map f . l)
  (cond ((nil? (car l)) '())
        ((nil? (cdr l))
         (cons (f (car (car l)))
               (map f (cdr (car l)))))
        (else
          (cons (apply f (map car l))
                (apply map (cons f (map cdr l)))))))


(def (append . l)
  (cond ((nil? l) nil)
        ((nil? (cdr l)) (car l))
        (else (foldr cons
                     (apply append (cdr l))
                     (car l)))))


(def (part f . a) (fn b (apply f (append a b))))


(def (partr f . a) (fn b (apply f (append b a))))


(def reverse (part fold cons nil))


(def (reverse-irregular l)
  (cond ((or (nil? l) (nil? (cdr l))) irregular-cons)
        (else (fold cons (cons (car (cdr l)) (car l)) (cdr (cdr l))))))


(def (inc n) (+ n 1))


(def (dec n) (- n 1))


(def (>= . n)
  (cond ((nil? n) false)
        ((nil? (cdr n)) true)
        ((and (not (> (car n) (car (cdr n))))
              (not (= (car n) (car (cdr n)))))
         false)
        (else (apply >= (cdr n)))))


(def list-len (part fold (fn (_ c) (inc c)) 0))


(def (len v)
  (cond ((vector? v) (vector-len v))
        ((struct? v) (len (struct-names s)))
        (else (list-len v))))


(def (mem f l)
  (cond ((nil? l) false)
        ((f (car l)) l)
        (else (mem f (cdr l)))))


(def (memq v l) (mem (part = v) l))


(def (notf f) (fn (i) (not (f i))))


(def (every? f . l) (not (mem (notf f) l)))


(def (any? . v) (and (mem id v) true))


(def (take n l)
  (cond ((= n 0) nil)
        (else (cons (car l)
                    (take (- n 1) (cdr l))))))


(def (drop n l)
  (cond ((= n 0) l)
        (else (drop (dec n) (cdr l)))))


(def (pad f n) (fn a (apply f (drop n a))))


(def (padr f n) (fn a (apply f (take (- (len a) n) a))))


(def (flip f) (fn a (apply f (reverse a))))


(def (check-types v . types)
  (apply any? (map (fn (t?) (t? v)) types)))


(def (compose . f) (partr (part fold call) f))


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
            (read-compile (compile c r:state) r)))))


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
         compiler-append "m := mm.LoadedModule(env, "
         compile module-name
         compiler-append ");"
         compiler-append "if m == mm.UndefinedModule {"
         compiler-append "mm.LoadCompiledModule(env, "
         compile module-name
         compiler-append "); m = mm.LoadedModule(env, "
         compile module-name
         compiler-append ")}; "
         (partr compile-module-define import-name) module-name))
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
         env := mm.InitialEnv()
     ")


(def compiled-tail "}")


(def irregular-cons (string->error "irregular cons expression"))
(def unexpected-close (string->error "unexpected close token"))
(def invalid-statement (string->error "invalid expression in statement position"))
(def invalid-cond (string->error "invalid cond expression"))


(def token-type
  {none    0
   symbol  1
   string  2
   comment 3
   list    4
   quote   5
   vector  6
   struct  7})


(def list-type
  {none   0
   lisp   1
   vector 2
   struct 3})


; something that doesn't equal anything else
; only itself
(def undefined {})


(def (reader input)
  {input         input 
   token-type    token-type:none
   state         undefined
   escaped?       false
   char          ""
   token         ""
   list-type     list-type:none
   close-list?   false
   list-items    nil
   list-cons?    false
   cons-items    0})


(def (read-char r)
  (let (next-input (fread r:input 1)
       state       (fstate next-input))
    (assign r
      {input next-input}
      (if (error? state)
        {state state}
        {char state}))))


; (test "read-char"
;   (test "returns error"
;     (let (r (reader (failing-io))
;           next (read-char r))
;       (error? next:state)))
; 
;   (test "returns eof"
;     (let (r (reader (buffer))
;           next (read-char r))
;       (= next:state eof)))
; 
;   (test "returns single char"
;     (let (r (reader (fwrite (buffer) "abc"))
;           next (read-char r))
;       (= next:char "a"))))


(def (newline? c) (= c "\n"))


(def (char-check c . cc)
  (cond ((nil? cc) false)
        ((= (car cc) c) true)
        (else (apply char-check (cons c (cdr cc))))))


(def (make-char-check . cc) (fn (c) (apply char-check (cons c cc))))
(def whitespace? (make-char-check " " "\b" "\f" "\n" "\r" "\t" "\v"))
(def escape-char? (make-char-check "\\"))
(def string-delimiter? (make-char-check "\""))
(def comment-char? (make-char-check ";"))
(def list-open? (make-char-check "("))
(def list-close? (make-char-check ")"))
(def cons-char? (make-char-check "."))
(def quote-char? (make-char-check "'"))
(def vector-open? (make-char-check "["))
(def vector-close? (make-char-check "]"))
(def struct-open? (make-char-check "{"))
(def struct-close? (make-char-check "}"))


(def (set-escaped r)
  (assign r {escaped? true}))


(def (append-token r)
  (assign r {token (string-append r:token r:char)}))


; (test "append-token"
;   (let (r    (assign (reader (buffer)) {token "ab" char "c"})
;         next (append-token r))
;     (= next:token "abc")))


(def (append-token-escaped r)
  (assign r {token (string-append r:token (cond ((= r:char "b") "\b")
                                                ((= r:char "f") "\f")
                                                ((= r:char "n") "\n")
                                                ((= r:char "r") "\r")
                                                ((= r:char "t") "\t")
                                                ((= r:char "v") "\v")
                                                (else r:char)))}))


(def (clear-token r) (assign r {token-type token-type:none token ""}))


(def (make-set-token-type tt) (fn (r) (assign r {token-type tt})))
(def set-symbol (make-set-token-type token-type:symbol))
(def set-string (make-set-token-type token-type:string))
(def set-comment (make-set-token-type token-type:comment))
(def set-list (make-set-token-type token-type:list))
(def set-quote (make-set-token-type token-type:quote))
(def set-vector (make-set-token-type token-type:vector))
(def set-struct (make-set-token-type token-type:struct))


(def (set-close r list-type)
  (if (= list-type r:list-type)
    (assign r {close-list? true})
    (assign r {state unexpected-close})))


(def (set-cons r)
  (assign r (if (= r:list-type list-type:lisp)
              {list-cons? true}
              {state irregular-cons})))


(def (symbol-token r)
  (def (try-parse . parsers)
    (if (nil? parsers)
      (string->symbol r:token)
      (let (v ((car parsers) r:token))
        (if (error? v) (apply try-parse (cdr parsers)) v))))
  (assign r {state (try-parse string->number
                              string->bool)}))


; (test "symbol-token"
;   (test "no empty token allowed"
;     (let (r (symbol-token (reader (buffer))))
;       (error? r:state)))
; 
;   (test "takes number"
;     (let (r (symbol-token (assign (reader (buffer)) {token "123"})))
;       (= r:state 123)))
; 
;   (test "takes bool"
;     (let (r (symbol-token (assign (reader (buffer)) {token "false"})))
;       (not r:state)))
;   
;   (test "takes symbol"
;     (let (r (symbol-token (assign (reader (buffer)) {token "a-symbol"})))
;       (= r:state 'a-symbol))))


; TODO: clean this up
(def (finalize-token r)
  (cond ((= r:token-type token-type:symbol)
         (clear-token (symbol-token r)))
        ((= r:token-type token-type:string)
         (clear-token (assign r {state r:token})))
        ((= r:token-type token-type:none) (clear-token r))
        (else
          (assign r {state (if (= r:token "") undefined invalid-token)}))))


; (test "finalize-token"
;   (test "invalid token type"
;     (let (r (finalize-token
;               (assign (reader (buffer))
;                 {token-type token-type:vector
;                  token "a"})))
;       (error? r:state)))
; 
;   (test "empty token"
;     (let (r (finalize-token
;               (assign (reader (buffer))
;                       {token-type token-type:none
;                        token ""})))
;       (= r:state undefined)))
; 
;   (test "number"
;     (let (r (finalize-token
;               (assign (reader (buffer))
;                 {token-type token-type:symbol
;                  token "123"})))
;       (= r:state 123)))
; 
;   (test "bool"
;     (let (r (finalize-token
;               (assign (reader (buffer))
;                 {token-type token-type:symbol
;                  token "false"})))
;       (not r:state)))
; 
;   (test "symbol"
;     (let (r (finalize-token
;               (assign (reader (buffer))
;                 {token-type token-type:symbol
;                  token "a-symbol"})))
;       (= r:state 'a-symbol))))


; todo: unimperativize
(def (read-list r list-type)
  (def (read-item lr)
    (let (next (read lr))
      (cond ((error? next:state) next)
            ((= next:state undefined) next)
            (else
              (assign next
                {list-items (cons next:state next:list-items)
                 cons-items (if (> next:cons-items 0) (inc lr:cons-items) 0)
                 state      undefined})))))

  (def (check-cons lr)
    (if lr:list-cons?
      (assign lr
        (if (or (nil? lr:list-items) (> lr:cons-items 0))
          {state irregular-cons}
          {cons-items 1 list-cons? false}))
      lr))

  (def (complete-list lr)
    (assign r
      {input lr:input
       token-type token-type:none
       state (cond ((and (> lr:cons-items 0) (not (= lr:cons-items 2)))
                    irregular-cons)
                   ((> lr:cons-items 0)
                    (reverse-irregular lr:list-items))
                   (else (reverse lr:list-items)))}))

  (def (read-items lr)
    (let (next (check-cons (read-item lr)))
      (cond ((error? next:state) (assign r {input next:input state next:state}))
            (next:close-list? (complete-list next))
            (else (read-items next)))))
  
  (read-items (assign (reader r:input) {list-type list-type})))


(def (read-quote r)
  (let (lr (assign (reader r:input) {list-type r:list-type})
        next (read lr))
    (assign r {input next:input
               token-type token-type:none
               state (if (error? next:state) next:state (list 'quote next:state))
               close-list? next:close-list?})))


(def (read-vector r)
  (let (next (read-list r list-type:vector))
    (assign next
      {state (if (error? next:state)
               next:state
               (cons 'vector: next:state))})))


(def (read-struct r)
  (let (next (read-list r list-type:struct))
    (assign next
      {state (if (error? next:state)
               next:state
               (cons 'struct: next:state))})))


(def (read r)
  (cond
    ((= r:token-type token-type:list)
     (assign (read-list r list-type:lisp)
             {token-type: token-type:none}))

    ((= r:token-type token-type:quote)
     (assign (read-quote r)
             {token-type: token-type:none}))

    ((= r:token-type token-type:vector)
     (assign (read-vector r)
             {token-type: token-type:none}))

    ((= r:token-type token-type:struct)
     (assign (read-struct r)
             {token-type: token-type:none}))

    (else
      (let (next (read-char r))
        (cond ((= next:state eof) (finalize-token next))

              ((error? next:state) next)

              ((= next:token-type token-type:none)
               (cond ((whitespace? next:char) (read next))
                     ((escape-char? next:char)
                      (read (set-symbol (set-escaped next))))
                     ((string-delimiter? next:char)
                      (read (set-string next)))
                     ((comment-char? next:char)
                      (read (set-comment next)))
                     ((list-open? next:char)
                      (read (set-list next)))
                     ((list-close? next:char)
                      (set-close next list-type:lisp))
                     ((cons-char? next:char) (set-cons next))
                     ((quote-char? next:char) (read (set-quote next)))
                     ((vector-open? next:char) (read (set-vector next)))
                     ((vector-close? next:char) (set-close next list-type:vector))
                     ((struct-open? next:char) (read (set-struct next)))
                     ((struct-close? next:char) (set-close next list-type:struct))
                     (else (read (append-token (set-symbol next))))))

              ((= next:token-type token-type:comment)
               (cond ((newline? next:char) (read (clear-token next)))
                     (else (read next))))

              ((= next:token-type token-type:symbol)
               (cond (next:escaped?
                      (read (assign (append-token-escaped next) {escaped? false})))
                     ((whitespace? next:char) (finalize-token next))
                     ((escape-char? next:char) (read (set-escaped next)))
                     ((string-delimiter? next:char)
                      (set-string (finalize-token next)))
                     ((comment-char? next:char)
                      (set-comment (finalize-token next)))
                     ((list-open? next:char)
                      (set-list (finalize-token next)))
                     ((list-close? next:char)
                      (set-close (finalize-token next) list-type:lisp))
                     ((vector-open? next:char) (set-vector (finalize-token next)))
                     ((vector-close? next:char)
                      (set-close (finalize-token next) list-type:vector))
                     ((struct-open? next:char) (set-struct (finalize-token next)))
                     ((struct-close? next:char)
                      (set-close (finalize-token next) list-type:struct))
                     (else (read (append-token next)))))

              ((= next:token-type token-type:string)
               (cond (next:escaped? (read (assign (append-token-escaped next) {escaped? false})))
                     ((escape-char? next:char) (read (set-escaped next)))
                     ((string-delimiter? next:char) (finalize-token next))
                     (else (read (append-token next)))))

              (else (assign next {state invalid-token})))))))


; (test "read"
;   (def (assert-read result test-value)
;     (cond ((= test-value undefined) (= result undefined))
;           ((number? test-value) (= result test-value))
;           ((bool? test-value) (= result test-value))
;           ((symbol? test-value) (= result test-value))
;           ((string? test-value) (= result test-value))
;           ((nil? test-value) (nil? result))
;           ((pair? test-value) (and (pair? result)
;                                    (assert-read (car result) (car test-value))
;                                    (assert-read (cdr result) (cdr test-value))))
;           ((= test-value eof) (= result eof))
;           (else (test-value result))))
; 
;   (def (read-string s) (read (reader (fwrite (buffer) s))))
; 
;   (def (assert-read-string string expected)
;     (let (r (read-string string))
;       (assert-read r:state expected)))
; 
;   (test "returns error"
;     (let (r (read (reader (failing-io))))
;       (assert-read r:state error?)))
; 
;   (test "returns on eof" (assert-read-string "a" 'a))
; 
;   (test "ignores whitespace" (assert-read-string " " eof))
; 
;   (test "reads number" (assert-read-string "123" 123))
; 
;   (test "reads bool" (assert-read-string "false" not))
; 
;   (test "reads symbol" (assert-read-string "a-symbol" 'a-symbol))
; 
;   (test "reads symbol with escape" (assert-read-string "a\\(b" 'a\(b))
; 
;   (test "reads symbol with escaped escape char" (assert-read-string "a\\\\b" 'a\\b))
; 
;   (test "reads symbol with escaped whitespace only" (assert-read-string "\\ " '\ ))
; 
;   (test "reads symbol closed by whitespace" (assert-read-string "a b" 'a))
; 
;   (test "reads symbol closed by string" (assert-read-string "a\"b\"" 'a))
; 
;   (test "reads symbol closed by comment" (assert-read-string "a; a comment" 'a))
; 
;   (test "reads symbol closed by list" (assert-read-string "a(b)" 'a))
; 
;   (test "reads symbol closed by vector" (assert-read-string "a[b]" 'a))
; 
;   (test "reads symbol closed by struct" (assert-read-string "a{b 2}" 'a))
; 
;   (test "reads string" (assert-read-string "\"a string\"" "a string"))
; 
;   (test "reads string with escape" (assert-read-string "\"a \\\"string\\\"\"" "a \"string\""))
; 
;   (test "reads string with escaped escape char" (assert-read-string "\"a \\\\\"" "a \\"))
; 
;   (test "ignores comments" (assert-read-string "; a comment" undefined))
; 
;   (test "new line closes comments" (assert-read-string "; a comment\na" 'a))
; 
;   (test "reads list" (assert-read-string "(a b c)" '(a b c)))
; 
;   (test "reads irregular list" (assert-read-string "(a b . c)" '(a b . c)))
; 
;   (test "reads quote" (assert-read-string "'a" ''a))
; 
;   (test "reads quoted list" (assert-read-string "'(a b c)" ''(a b c)))
; 
;   (test "reads vector" (assert-read-string "[a b c]" '[a b c]))
; 
;   (test "reads struct" (assert-read-string "{a 1 b 2 c 3}" '{a 1 b 2 c 3}))
;   
;   (test "yields ready objects"
;     (def (test-yield r expected)
;       (let (next (read r))
;         (cond ((and (nil? expected) (= next:state eof)) 'ok)
;               ((and (not (nil? expected)) (= next:state eof))
;                (string->error "more values expected"))
;               ((nil? expected) (string->error "unexpected state"))
;               (else (and (assert-read next:state (car expected))
;                          (test-yield next (cdr expected)))))))
;     (test-yield (reader (fwrite (buffer) "a b c")) '(a b c))
;     (test-yield (reader (fwrite (buffer) "a {b 2} [c]")) '(a {b 2} [c]))))


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


; (read-compile-write)
(read-eval-print)
