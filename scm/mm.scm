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


(def token-type
  {none    0
   symbol  1
   string  2
   comment 3
   list    4})


(def list-type
  {none 0
   lisp 1})


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
(def whitespace? (make-char-check " " "\b" "\f" "\r" "\t" "\v"))
(def escape-char? (make-char-check "\\"))
(def string-delimiter? (make-char-check "\""))
(def comment-char? (make-char-check ";"))
(def list-open? (make-char-check "("))
(def list-close? (make-char-check ")"))
(def cons-char? (make-char-check "."))


(def (set-escaped r)
  (assign r {escaped? true}))


(def (append-token r)
  (assign r {token (string-append r:token r:char)}))


(test "append-token"
  (let (r    (assign (reader (buffer)) {token "ab" char "c"})
        next (append-token r))
    (= next:token "abc")))


(def (clear-token r) (assign r {token-type token-type:none token ""}))


(def (make-set-token-type tt) (fn (r) (assign r {token-type tt})))
(def set-symbol (make-set-token-type token-type:symbol))
(def set-string (make-set-token-type token-type:string))
(def set-comment (make-set-token-type token-type:comment))
(def set-list (make-set-token-type token-type:list))


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
  (assign r {value (try-parse try-string->number
                              try-string->bool)}))


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


(def (finalize-token r)
  (cond ((= r:token-type token-type:symbol)
         (clear-token (symbol-token r)))
        ((= r:token-type token-type:string)
         (clear-token (assign r {value r:token})))
        (else
          (assign r
            {value
             (if (= r:token "")
               undefined
               invalid-token)}))))


(test "finalize-token"
  (test "invalid token type"
    (let (r (finalize-token
              (assign (reader (buffer))
                {token-type token-type:none
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


(def (read r)
  (cond
    ((= r:token-type token-type:list)
     (assign (read-list r list-type:lisp)
             {token-type: token-type:none}))
    (else
      (let (next (read-char r))
        (cond ((= next:value eof)
               (finalize-token next))

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
                     (else (read (append-token (set-symbol next))))))

              ((= next:token-type token-type:comment)
               (cond ((newline? next:char) (clear-token next))
                     (else (read next))))

              ((= next:token-type token-type:symbol)
               (cond (next:escaped? (read (append-token next)))
                     ((whitespace? next:char) (finalize-token next))
                     ((list-close? next:char)
                      (set-close (finalize-token next) list-type:lisp))
                     (else (read (append-token next)))))

              ((= next:token-type token-type:string)
               (cond ((string-delimiter? next:char) (finalize-token next))
                     (else (read (append-token next)))))

              (else (assign next {value invalid-token})))))))


(test "read"
  (def (assert-value result test-value)
    (cond ((symbol? test-value) (= result test-value))
          ((= test-value undefined) (= result undefined))
          ((number? test-value) (= result test-value))
          ((bool? test-value) (= result test-value))
          ((symbol? test-value) (= result test-value))
          ((string? test-value) (= result test-value))
          ((nil? test-value) (nil? result))
          ((pair? test-value) (and (pair? result)
                                   (assert-value (car result) (car test-value))
                                   (assert-value (cdr result) (cdr test-value))))
          (else (test-value result))))

  (test "returns on eof"
    (let (r (read (reader (fwrite (buffer) "a"))))
      (assert-value r:value 'a)))

  (test "returns error"
    (let (r (read (reader (failing-reader))))
      (assert-value r:value error?)))

  (test "ignores whitespace"
    (let (r (read (reader (fwrite (buffer) " "))))
      (assert-value r:value undefined)))

  (test "reads number"
    (let (r (read (reader (fwrite (buffer) "123"))))
      (assert-value r:value 123)))

  (test "reads bool"
    (let (r (read (reader (fwrite (buffer) "false"))))
      (assert-value r:value not)))

  (test "reads symbol"
    (let (r (read (reader (fwrite (buffer) "a-symbol"))))
      (assert-value r:value 'a-symbol)))

  (test "reads string"
    (let (r (read (reader (fwrite (buffer) "\"a string\""))))
      (assert-value r:value "a string")))

  (test "ignores comments"
    (let (r (read (reader (fwrite (buffer) "; a comment"))))
      (assert-value r:value undefined)))

  (test "reads list"
    (let (r (read (reader (fwrite (buffer) "(a b c)"))))
      (assert-value r:value '(a b c))))

  (trace "test started")
  (test "reads irregular list"
    (let (r (read (reader (fwrite (buffer) "(a b . c)"))))
	  (assert-value (trace r:value) '(a b . c))))

  (test "escapes characters"
    (let (r (read (reader (fwrite (buffer) "\\ "))))
      (assert-value r:value '\ ))))
