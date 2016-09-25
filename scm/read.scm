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

(export reader reader read read)
