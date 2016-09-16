(def (trace message . values)
  (def (trace out values)
    (cond ((nil? (cdr values))
           (def out-complete (print out (car values)))
           (fwrite out-complete:output "\n")
           (car values))
          (else
            (def out-next (print out (car values)))
            (trace (assign out-next {output (fwrite out-next:output " ")})
                   (cdr values)))))
  (trace (printer (stderr)) (cons message values)))


(def (len l)
  (cond ((nil? l) 0)
        (else (inc (len (cdr l))))))


(def (append . l)
  (def (append-two left right)
    (cond ((nil? left) right)
          (else (cons (car left) (append-two (cdr left) right)))))
  (cond ((nil? l) nil)
        (else (append-two (car l) (apply append (cdr l))))))


(def (reverse l)
  (cond ((nil? l) nil)
        (else (append (reverse (cdr l)) (list (car l))))))


(def (reverse-irregular l)
  (def (reverse from to)
    (cond ((nil? from) to)
          (else (reverse (cdr from) (cons (car from) to)))))
  (cond ((or (nil? l) (nil? (cdr l))) irregular-cons)
        (else (reverse (cdr (cdr l)) (cons (car (cdr l)) (car l))))))


(def (inc n) (+ n 1))


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
   value         undefined
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
        {value state}
        {char state}))))


(test "read-char"
  (test "returns derived input"
    (let (b (buffer)
          r (read-char (reader b)))
      (derived-object? r:input b)))

  (test "returns error"
    (let (r (reader (failing-reader))
          next (read-char r))
      (error? next:value)))

  (test "returns eof"
    (let (r (reader (buffer))
          next (read-char r))
      (= next:value eof)))

  (test "returns single char"
    (let (r (reader (fwrite (buffer) "abc"))
          next (read-char r))
      (= next:char "a"))))


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


(test "append-token"
  (let (r    (assign (reader (buffer)) {token "ab" char "c"})
        next (append-token r))
    (= next:token "abc")))


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
    (assign r {value unexpected-close})))


(def (set-cons r)
  (assign r (if (= r:list-type list-type:lisp)
              {list-cons? true}
              {value irregular-cons})))


(def (symbol-token r)
  (def (try-parse . parsers)
    (if (nil? parsers)
      (string->symbol r:token)
      (let (v ((car parsers) r:token))
        (if (error? v) (apply try-parse (cdr parsers)) v))))
  (assign r {value (try-parse string->number
                              string->bool)}))


(test "symbol-token"
  (test "no empty token allowed"
    (let (r (symbol-token (reader (buffer))))
      (error? r:value)))

  (test "takes number"
    (let (r (symbol-token (assign (reader (buffer)) {token "123"})))
      (= r:value 123)))

  (test "takes bool"
    (let (r (symbol-token (assign (reader (buffer)) {token "false"})))
      (not r:value)))
  
  (test "takes symbol"
    (let (r (symbol-token (assign (reader (buffer)) {token "a-symbol"})))
      (= r:value 'a-symbol))))


; TODO: clean this up
(def (finalize-token r)
  (cond ((= r:token-type token-type:symbol)
         (clear-token (symbol-token r)))
        ((= r:token-type token-type:string)
         (clear-token (assign r {value r:token})))
        ((= r:token-type token-type:none) (clear-token r))
        (else
          (assign r {value (if (= r:token "") undefined invalid-token)}))))


(test "finalize-token"
  (test "invalid token type"
    (let (r (finalize-token
              (assign (reader (buffer))
                {token-type token-type:vector
                 token "a"})))
      (error? r:value)))

  (test "empty token"
    (let (r (finalize-token
              (assign (reader (buffer))
                      {token-type token-type:none
                       token ""})))
      (= r:value undefined)))

  (test "number"
    (let (r (finalize-token
              (assign (reader (buffer))
                {token-type token-type:symbol
                 token "123"})))
      (= r:value 123)))

  (test "bool"
    (let (r (finalize-token
              (assign (reader (buffer))
                {token-type token-type:symbol
                 token "false"})))
      (not r:value)))

  (test "symbol"
    (let (r (finalize-token
              (assign (reader (buffer))
                {token-type token-type:symbol
                 token "a-symbol"})))
      (= r:value 'a-symbol))))


; todo: unimperativize
(def (read-list r list-type)
  (def (read-item lr)
    (let (next (read lr))
      (cond ((error? next:value) next)
            ((= next:value undefined) next)
            (else
              (assign next
                {list-items (cons next:value next:list-items)
                 cons-items (if (> next:cons-items 0) (inc lr:cons-items) 0)
                 value      undefined})))))

  (def (check-cons lr)
    (if lr:list-cons?
      (assign lr
        (if (or (nil? lr:list-items) (> lr:cons-items 0))
          {value irregular-cons}
          {cons-items 1 list-cons? false}))
      lr))

  (def (complete-list lr)
    (assign r
      {input lr:input
       token-type token-type:none
       value (cond ((and (> lr:cons-items 0) (not (= lr:cons-items 2)))
                    irregular-cons)
                   ((> lr:cons-items 0)
                    (reverse-irregular lr:list-items))
                   (else (reverse lr:list-items)))}))

  (def (read-items lr)
    (let (next (check-cons (read-item lr)))
      (cond ((error? next:value) (assign r {input next:input value next:value}))
            (next:close-list? (complete-list next))
            (else (read-items next)))))
  
  (read-items (assign (reader r:input) {list-type list-type})))


(def (read-quote r)
  (let (lr (assign (reader r:input) {list-type r:list-type})
        next (read lr))
    (assign r {input next:input
               token-type token-type:none
               value (if (error? next:value) next:value (list 'quote next:value))
               close-list? next:close-list?})))


(def (read-vector r)
  (let (next (read-list r list-type:vector))
    (assign next
      {value (if (error? next:value)
               next:value
               (cons 'vector: next:value))})))


(def (read-struct r)
  (let (next (read-list r list-type:struct))
    (assign next
      {value (if (error? next:value)
               next:value
               (cons 'struct: next:value))})))


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
        (cond ((= next:value eof) (finalize-token next))

              ((error? next:value) next)

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

              (else (assign next {value invalid-token})))))))


(test "read"
  (def (assert-read result test-value)
    (cond ((= test-value undefined) (= result undefined))
          ((number? test-value) (= result test-value))
          ((bool? test-value) (= result test-value))
          ((symbol? test-value) (= result test-value))
          ((string? test-value) (= result test-value))
          ((nil? test-value) (nil? result))
          ((pair? test-value) (and (pair? result)
                                   (assert-read (car result) (car test-value))
                                   (assert-read (cdr result) (cdr test-value))))
          ((= test-value eof) (= result eof))
          (else (test-value result))))

  (def (read-string s) (read (reader (fwrite (buffer) s))))

  (def (assert-read-string string expected)
    (let (r (read-string string))
      (assert-read r:value expected)))

  (test "returns error"
    (let (r (read (reader (failing-reader))))
      (assert-read r:value error?)))

  (test "returns on eof" (assert-read-string "a" 'a))

  (test "ignores whitespace" (assert-read-string " " eof))

  (test "reads number" (assert-read-string "123" 123))

  (test "reads bool" (assert-read-string "false" not))

  (test "reads symbol" (assert-read-string "a-symbol" 'a-symbol))

  (test "reads symbol with escape" (assert-read-string "a\\(b" 'a\(b))

  (test "reads symbol with escaped escape char" (assert-read-string "a\\\\b" 'a\\b))

  (test "reads symbol with escaped whitespace only" (assert-read-string "\\ " '\ ))

  (test "reads symbol closed by whitespace" (assert-read-string "a b" 'a))

  (test "reads symbol closed by string" (assert-read-string "a\"b\"" 'a))

  (test "reads symbol closed by comment" (assert-read-string "a; a comment" 'a))

  (test "reads symbol closed by list" (assert-read-string "a(b)" 'a))

  (test "reads symbol closed by vector" (assert-read-string "a[b]" 'a))

  (test "reads symbol closed by struct" (assert-read-string "a{b 2}" 'a))

  (test "reads string" (assert-read-string "\"a string\"" "a string"))

  (test "reads string with escape" (assert-read-string "\"a \\\"string\\\"\"" "a \"string\""))

  (test "reads string with escaped escape char" (assert-read-string "\"a \\\\\"" "a \\"))

  (test "ignores comments" (assert-read-string "; a comment" undefined))

  (test "new line closes comments" (assert-read-string "; a comment\na" 'a))

  (test "reads list" (assert-read-string "(a b c)" '(a b c)))

  (test "reads irregular list" (assert-read-string "(a b . c)" '(a b . c)))

  (test "reads quote" (assert-read-string "'a" ''a))

  (test "reads quoted list" (assert-read-string "'(a b c)" ''(a b c)))

  (test "reads vector" (assert-read-string "[a b c]" '[a b c]))

  (test "reads struct" (assert-read-string "{a 1 b 2 c 3}" '{a 1 b 2 c 3}))
  
  (test "yields ready objects"
    (def (test-yield r expected)
      (let (next (read r))
        (cond ((and (nil? expected) (= next:value eof)) 'ok)
              ((and (not (nil? expected)) (= next:value eof))
               (string->error "more values expected"))
              ((nil? expected) (string->error "unexpected value"))
              (else (and (assert-read next:value (car expected))
                         (test-yield next (cdr expected)))))))
    (test-yield (reader (fwrite (buffer) "a b c")) '(a b c))
    (test-yield (reader (fwrite (buffer) "a {b 2} [c]")) '(a {b 2} [c]))))


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
(def (application? v) (pair? v))


(def (compile-number v) (string-append " mm.NumberFromRawInt(" (number->string v) ") "))
(def (compile-string v) (string-append " mm.StringFromRaw(" (escape-compiled-string v) ") "))
(def (compile-bool v) (if v " mm.True " " mm.False "))
(def (compile-nil v) " mm.Nil ")
(def (compile-quote-literal v) (string-append " mm.List(mm.SfromString("
                                      (escape-compiled-string (symbol->string 'quote))
                                      "), "
                                      (compile-literal (car (cdr v)))
                                      ") "))
(def (compile-quote v) (compile-literal (car (cdr v))))
(def (compile-symbol v) (string-append " mm.SfromString("
                                       (escape-compiled-string (symbol->string v))
                                       ") "))


(def (compile-pair v)
  (string-append " mm.Cons("
                 (compile-literal (car v))
                 ", "
                 (compile-literal (cdr v))
                 ") "))


(def (make-fn args body) (cons 'fn (cons args body)))


(def (def-name v)
  (cond ((or (not (pair? v)) (not (pair? (cdr v)))) (fatal invalid-def))
        ((symbol? (car (cdr v))) (car (cdr v)))
        (else
          (cond ((or (not (pair? (car (cdr v)))) (not (symbol? (car (car (cdr v))))))
                 (fatal invalid-def))
                (else (car (car (cdr v))))))))


(def (def-value v)
  (cond ((or (not (pair? v)) (not (pair? (cdr v)))) (fatal invalid-def))
        ((symbol? (car (cdr v))) 
         (cond ((or (not (pair? (cdr (cdr v)))) (not (nil? (cdr (cdr (cdr v))))))
                (fatal invalid-def))
               (else (car (cdr (cdr v))))))
        (else
          (cond ((not (pair? (car (cdr v)))) (fatal invalid-def))
                (else (make-fn (cdr (car (cdr v))) (cdr (cdr v))))))))


(def (compile-def v) (string-append " mm.Define(env, "
                                    (compile-literal (def-name v))
                                    ", "
                                    (compile-exp (def-value v))
                                    ") "))


(def (fn-signature v)
  (cond ((nil? v) {count 0 var? false names '()})
        ((symbol? v) {count 0 var? true names v})
        ((pair? v) (let (signature (fn-signature (cdr v)))
                     (assign signature {count (inc signature:count)
                                         names (cons (car v) signature:names)})))
        (else (fatal invalid-fn))))


(def (compile-seq v)
  (cond ((not (pair? v)) (fatal invalid-seq))
        ((nil? (cdr v)) (string-append "return " (compile-exp (car v))))
        (else (string-append (compile (car v)) ";\n" (compile-seq (cdr v))))))


(def (compile-fn v)
  (cond ((or (not (pair? v))
             (not (pair? (cdr v))))
         (fatal invalid-fn))
        (else
          (let (signature (fn-signature (car (cdr v)))
                body      (cdr (cdr v)))
            (string-append
              " mm.NewCompiled("
              (number->string signature:count)
              ", "
              (bool->string signature:var?)
              ", func(a []*mm.Val) *mm.Val { env := mm.ExtendEnv(env, "
              (compile-literal signature:names)
              ", mm.List(a...)); env = env; "
              (compile-seq body)
              "})")))))


(def (compile-if v)
  (cond ((not (= (len v) 4)) (fatal invalid-if))
        (else (string-append " func() *mm.Val { if "
                             (compile-exp (car (cdr v)))
                             " != mm.False { return "
                             (compile-exp (car (cdr (cdr v))))
                             " } else { return "
                             (compile-exp (car (cdr (cdr (cdr v)))))
                             " }}() "))))


(def (compile-and v)
  (def (compile-and s)
    (cond ((nil? s) " return mm.True ")
          ((nil? (cdr s))
           (string-append " return " (compile-exp (car s))))
          (else
            (string-append " if "
                           (compile-exp (car s))
                           " == mm.False { return mm.False }; "
                           (compile-and (cdr s))))))
  (string-append " func() *mm.Val { "
                 (compile-and (cdr v))
                 " }() "))


(def (compile-or v)
  (def (compile-or s)
    (cond ((nil? s) " return mm.False ")
          ((nil? (cdr s))
           (string-append " return " (compile-exp (car s))))
          (else
            (string-append "if v := "
                           (compile-exp (car s))
                           "; v != mm.False { return v }; "
                           (compile-or (cdr s))))))
  (string-append " func() *mm.Val { "
                 (compile-or (cdr v))
                 " }() "))


(def (compile-begin v)
  (string-append " func() *mm.Val { "
                 (compile-seq (cdr v))
                 " }() "))


(def (cond->if v)
  (def (seq->exp v)
    (cond ((nil? (cdr v)) (car v))
          (else (cons 'begin v))))
  (def (expand v)
    (cond ((or (not (pair? v))
               (not (pair? (car v)))
               (not (pair? (cdr (car v)))))
           (fatal invalid-cond))
          ((= (car (car v)) 'else)
           (cond ((not (nil? (cdr v)))
                  (fatal invalid-cond))
                 (else (seq->exp (cdr (car v))))))
          (else
            (list 'if (car (car v))
                  (seq->exp (cdr (car v)))
                  (expand (cdr v))))))
  (expand (cdr v)))


(def (compile-cond v) (compile (cond->if v)))


(def (compile-let v)
  (def (let-defs v)
    (cond ((nil? v) nil)
          ((or (not (pair? v)) (not (pair? (cdr v))))
           (fatal invalid-let))
          (else (cons (list 'def (car v) (car (cdr v)))
                      (let-defs (cdr (cdr v)))))))
  (def (let-body v)
    (cond ((or (not (pair? v))
               (not (pair? (cdr v)))
               (nil? (cdr (cdr v))))
           (fatal invalid-let))
          (else (append (let-defs (car (cdr v))) (cdr (cdr v))))))
  (compile (list (make-fn nil (let-body v)))))


(def (compile-test v)
  (def (compile-test-seq v)
    (cond ((nil? v) (string-append " return " (compile-literal 'test-complete)))
          (else (string-append " if result := func() *mm.Val { return "
                               (compile-exp-statement (car v))
                               " }(); result == mm.False { return mm.Fatal("
                               (compile-exp "test failed")
                               ") } else if mm.IsError(result) != mm.False "
                               " { return mm.Fatal(result) }; "
                               (compile-test-seq (cdr v))))))
  (string-append " func() *mm.Val { "
                 " env := mm.ExtendEnv(env, mm.Nil, mm.Nil); env = env; "
                 (compile-test-seq (cdr v))
                 " }() "))


(def (compile-lookup v)
  (string-append " mm.LookupDef(env, "
                 (compile-symbol v)
                 ") "))


(def (compile-vector v)
  (string-append " mm.VectorFromList("
                 (compile-literal (cdr v))
                 " )"))


(def (compile-struct v)
  (def (compile-struct-values v)
    (cond ((nil? v) " mm.Nil ")
          (else (string-append " mm.Cons( "
                               (compile-literal (car v))
                               ", mm.Cons( "
                               (compile-exp (car (cdr v)))
                               ", "
                               (compile-struct-values (cdr (cdr v)))
                               ")) "))))
  (string-append " mm.StructFromList("
                 (compile-struct-values (cdr v))
                 ") "))


(def (compile-value-list v)
  (cond ((nil? v) " mm.Nil ")
        ((not (pair? v)) (fatal invalid-application))
        (else (string-append " mm.Cons( "
                             (compile-exp (car v))
                             ", "
                             (compile-value-list (cdr v))
                             ") "))))


(def (compile-application v)
  (string-append " mm.Apply("
                 (compile-exp (car v))
                 ", "
                 (compile-value-list (cdr v))
                 ")"))


(def (compile-literal v)
  (cond ((number? v) (compile-number v))
        ((string? v) (compile-string v))
        ((bool? v) (compile-bool v))
        ((nil? v) (compile-nil v))
        ((quote? v) (compile-quote-literal v))
        ((symbol? v) (compile-symbol v))
        ((pair? v) (compile-pair v))
        (else (fatal invalid-literal))))


(def (compile-exp v)
  (cond ((number? v) (compile-number v))
        ((string? v) (compile-string v))
        ((bool? v) (compile-bool v))
        ((nil? v) (compile-nil v))
        ((quote? v) (compile-quote v))
        ((symbol? v) (compile-lookup v))
        ((vector-form? v) (compile-vector v))
        ((struct-form? v) (compile-struct v))
        ((if? v) (compile-if v))
        ((and? v) (compile-and v))
        ((or? v) (compile-or v))
        ((fn? v) (compile-fn v))
        ((begin? v) (compile-begin v))
        ((cond? v) (compile-cond v))
        ((let? v) (compile-let v))
        ((test? v) (compile-test v))
        ((application? v) (compile-application v))
        (else (fatal invalid-expression))))


(def (compile v)
  (cond ((def? v) (compile-def v))
        ((if? v) (compile-if v))
        ((and? v) (compile-and v))
        ((or? v) (compile-or v))
        ((begin? v) (compile-begin v))
        ((cond? v) (compile-cond v))
        ((let? v) (compile-let v))
        ((test? v) (compile-test v))
        ((application? v) (compile-application v))
        (else (fatal invalid-statement))))


(def (compile-exp-statement v)
  (cond ((number? v) (compile-number v))
        ((string? v) (compile-string v))
        ((bool? v) (compile-bool v))
        ((nil? v) (compile-nil v))
        ((quote? v) (compile-quote v))
        ((symbol? v) (compile-lookup v))
        ((vector-form? v) (compile-vector v))
        ((struct-form? v) (compile-struct v))
        ((def? v) (compile-def v))
        ((if? v) (compile-if v))
        ((and? v) (compile-and v))
        ((or? v) (compile-or v))
        ((fn? v) (compile-fn v))
        ((begin? v) (compile-begin v))
        ((cond? v) (compile-cond v))
        ((let? v) (compile-let v))
        ((test? v) (compile-test v))
        ((application? v) (compile-application v))
        (else (fatal invalid-expression))))


(def compiled-head
     "package main

     import mm \"github.com/aryszka/mikkamakka\"

     func main() {
         env := mm.InitialEnv()
     ")


(def compiled-tail "}")


(def (compile-file fin fout)
  (def (write-head fout) (fwrite fout compiled-head))
  (def (write-tail fout) (fwrite fout compiled-tail))

  (def (compile-reader r fout)
    (let (next-in (read r))
      (cond ((= next-in:value eof) (list next-in fout))
            ((error? next-in:value) next-in:value)
            (else
              (let (code (compile next-in:value))
                (cond ((error? code) code)
                      (else
                        (let (next-out (fwrite fout (string-append code ";\n")))
                          (cond ((error? (fstate next-out)) (fstate next-out))
                                (else (compile-reader next-in next-out)))))))))))

  (let (result (compile-reader (reader fin) (write-head fout)))
    (cond ((error? result) result)
          (else (list ((car result) 'input) (write-tail (car (cdr result))))))))


(let (fin  (fopen (car (cdr (argv))))
      fout (stdout))
  (cond ((error? (fstate fin)) (fstate fin))
        ((error? (fstate fout)) (fstate fout))
        (else (let (result (compile-file fin fout))
                (cond ((error? result) (fatal result))
                      (else
                        (fwrite (car (cdr result)) "\n")
                        (fclose (car result))))))))
