; optimize:
; - reader again: pure reading, full tokens, complex construction
; - benchmarks and profiling
; - measure next biggest loss
; - create a checked mode and an unchecked mode, and apply the checks only there
; - complete syntax check during compilation
; debug utilities: stack
; sprint escaping: compound and compiled procedures
; check expression check. e.g. don't allow lambda parameters other then symbols
; tail recursion optimization:
; - tail context definition in r6rs
; don't override members in the global environment (currently in compile)
; remove code duplication between eval and compile, and others
; missing syntax:
; - let*, letrec
; arbitrary number of variables
; - check what needs to be a variadic
; introduce macros
; create quasiquotation
; generate ca/dr shortcuts
; delete functionality
; rename sbuilder, find similar lisp name
; introduce ports
; modules, tests
; load file into repl
; vim repl: http://www.vim.org/scripts/script.php?script_id=4336
; rewrite compilation with quasiquotation
; interop
; - reimplement io functions with interop
; verify that no foreign values can get in the mikkamakka env
; create a hashmap
; utilize standard input
; replace env with name table
; replace symbol-name with symbol->string
; check if eval shared needs to make a new environment every time
; don't compile precompiled if no expression defined
; make it more concise, maybe accept alternatives for lambda, like fn
; fix sprint: escape mikkamakka characters (non-printing and escaped characters)
; fix sprinting: '|\||
; unbind non-shared dependencies
; read/print struct
; review head to not allow emitting illegal types to scheme namespace
; convert string builder to function

(define (error where what arg)
  (cerror (sprint-quoted where)
          what
          (sprint arg)))

(define (symbol-eq? left right)
  (and (symbol? left)
       (symbol? right)
       (peq? (symbol-name left)
             (symbol-name right))))

(define (print s)
  (out (sprint s)))

(define (log exp)
  (clog (sprint exp)))

(define (newline)
  (out "\n"))

(define (print-stack error)
  (out (sprint-stack error))
  (newline))

(define (display s)
  (out (sunescape (sprint s))))

(define (self-evaluating? exp)
  (or (peq? exp true)
      (peq? exp false)
      (peq? exp no-print)
      (number? exp)
      (string? exp)))

(define (quote-text q) (cadr q))

(define (quote-eq? left right)
  (and (quote? left)
       (quote? right)
       (or (and (null? (cdr left))
                (null? (cdr right)))
           (eq? (quote-text left)
                (quote-text right)))))

(define (eq? left right)
  (or (and (null? left) (null? right))
      (symbol-eq? left right)
      (quote-eq? left right)
      (peq? left right)))

(define (not exp)
  (if (peq? exp false) true false))

(define (sescape-char char)
  (cond ((peq? char "\b") "\\b")
        ((peq? char "\t") "\\t")
        ((peq? char "\n") "\\n")
        ((peq? char "\v") "\\v")
        ((peq? char "\f") "\\f")
        ((peq? char "\r") "\\r")
        ((peq? char "\"") "\\\"")
        ((peq? char "\\") "\\\\")
        (else char)))

(define (escape-char? char)
  (not (peq? (sescape-char char) char)))

(define (sescaped-char char-symbol)
  (cond ((peq? char-symbol "b") "\b")
        ((peq? char-symbol "t") "\t")
        ((peq? char-symbol "n") "\n")
        ((peq? char-symbol "v") "\v")
        ((peq? char-symbol "f") "\f")
        ((peq? char-symbol "r") "\r")
        (else char-symbol)))

(define (sescape s)
  (define (sescapes ss ps)
    (if (peq? (slen ss) 0)
      ps
      (sescapes (subs ss 1 -1)
                (cats ps (sescape-char (subs ss 0 1))))))
  (sescapes s ""))

(define (sescape-symbol-name sn)
  (define (find-escape-char ss)
    (cond ((peq? (slen ss) 0)
           sn)
          ((escape-char? (subs ss 0 1))
           (cats "|" sn "|"))
          (else
            (find-escape-char (subs ss 1 -1)))))
  (find-escape-char sn))

(define (sunescape s)
  (define (sunescapes ss ps escaped)
    (cond ((peq? (slen ss) 0)
           (if escaped
             (error 'sunescape "invalid escape sequence" s)
             ps))
          (escaped
            (sunescapes (subs ss 1 -1)
                        (cats ps (sescaped-char (subs ss 0 1)))
                        false))
          ((peq? (subs ss 0 1) "\\")
           (sunescapes (subs ss 1 -1)
                       ps
                       true))
          (else
            (sunescapes (subs ss 1 -1)
                        (cats ps (subs ss 0 1))
                        false))))
  (sunescapes s "" false))

(define (sprintq-vector vector)
  (let ((builder (make-string-builder)))
    (set! builder (sbappend builder "#("))
    (let ((ref 0)
          (in-list? false))
      (for (lambda ()
             (if (eq? ref (vlen vector))
               true
               (begin
                 (cond ((> ref 0)
                        (set! builder
                          (sbappend builder " "))))
                 (set! builder
                   (sbappend builder
                             (sprintq (vref vector ref)
                                      true
                                      false)))
                 (set! in-list? true)
                 (set! ref (+ ref 1))
                 false)))))
    (set! builder (sbappend builder ")"))
    (builder->string builder)))

(define (sprintq exp quoted? in-list?)
  (cond ((peq? exp false) "false")
        ((peq? exp true) "true")
        ((peq? exp no-print) "")
        ((number? exp)
         (number->string exp))
        ((string? exp)
         (cats "\"" (sescape exp) "\""))
        ((symbol? exp)
         (if quoted?
           (sescape-symbol-name (symbol-name exp))
           (cats "'" (sescape-symbol-name (symbol-name exp)))))
        ((null? exp)
         (if quoted?
           "()"
           "'()"))
        ((quote? exp)
         (let ((text (try (lambda () (quote-text exp))
                          (lambda (_) "???"))))
           (cats "'" (sprintq text true false))))
        ((vector? exp)
         (sprintq-vector exp))
        ((pair? exp)
         (cats (cond ((not quoted?) "'(")
                     ((not in-list?) "(")
                     (else ""))
               (cond (in-list? " ")
                     (else ""))
               (sprintq (car exp) true false)
               (cond ((null? (cdr exp)) ")")
                     ((pair? (cdr exp)) (sprintq (cdr exp) true true))
                     (else (cats " . "
                                 (sprintq (cdr exp) true true)
                                 ")")))))
        ((table? exp) "[table]")
        ((error? exp)
         (sprint-error exp))
        (else
          "unknown type")))

(define (sprint exp)
  (sprintq exp false false))

(define (sprint-quoted exp)
  (sprintq exp true false))

(define (tagged-list? exp tag)
  (and (pair? exp) (eq? (car exp) tag)))

(define (len l)
  (if (null? l)
    0
    (+ 1 (len (cdr l)))))

(define (append left right)
  (if (null? left)
    right
    (cons (car left) (append (cdr left) right))))

(define (reverse l)
  (if (null? l)
    '()
    (append (reverse (cdr l)) (list (car l)))))

(define (map p l)
  (if (null? l)
    '()
    (cons (p (car l)) (map p (cdr l)))))

(define (deep-reverse l)
  (map (lambda (i)
         (if (pair? i)
           (deep-reverse i)
           i))
       (reverse l)))

(define (copy-table-filtered from to f)
  (define (iterate names copied-names)
    (if (null? names)
      copied-names
      (let ((name (car names)))
        (cond ((f name)
               (table-define to name (table-lookup from name))
               (set! copied-names (cons name copied-names))))
        (iterate (cdr names) copied-names))))
  (iterate (table-names from) '()))

(define (to-table l map-name map-value)
  (let ((table (make-name-table)))
    (map (lambda (i) (table-define table (map-name i) (map-value i))) l)
    table))

; quote
(define (quote? exp) (tagged-list? exp 'quote))

(define (check-quote exp)
  (cond ((or (not (tagged-list? exp 'quote))
             (not (peq? (len exp) 2)))
         (error 'check-quote "invalid quotation" '???))))

; assignment
(define (assignment-variable exp) (cadr exp))

(define (assignment? exp) (tagged-list? exp 'set!))

(define (check-assignment exp)
  (cond ((not (peq? (len exp) 3))
         (error 'check-assignment "invalid arity" exp))
        ((not (symbol? (assignment-variable exp)))
         (error 'check-assignment "invalid variable name" exp))))

(define (assignment-variable exp) (cadr exp))

(define (assignment-value exp) (caddr exp))

(define (eval-assignment env exp)
  (env (assignment-variable exp)
       false
       (eval-env env (assignment-value exp))))

; definition
(define (definition? exp) (tagged-list? exp 'define))

(define (definition-variable exp)
  (if (symbol? (cadr exp))
    (cadr exp)
    (caadr exp)))

(define (definition-value exp)
  (if (symbol? (cadr exp))
    (caddr exp)
    (make-lambda (cdadr exp) (cddr exp))))

(define (check-definition exp)
  (cond ((not (or (and (symbol? (cadr exp))
                       (peq? (len exp) 3))
                  (and (pair? (cadr exp))
                       (> (len exp) 2))))
         (error 'check-definition "invalid format" exp))
        ((not (symbol? (definition-variable exp)))
         (error 'check-definition "invalid variable name" exp))))

(define (eval-definition env exp)
  (env (definition-variable exp)
       true
       (eval-env env (definition-value exp))))

; if
(define (if? exp) (tagged-list? exp 'if))

(define (check-if exp)
  (cond ((not (peq? (len exp) 4))
         (error 'check-if "invalid arity" exp))))

(define (if-predicate exp) (cadr exp))

(define (if-consequent exp) (caddr exp))

(define (if-alternative exp) (cadddr exp))

(define (true? exp) (not (peq? exp false)))

(define (eval-if env exp)
  (if (true? (eval-env env (if-predicate exp)))
    (eval-env env (if-consequent exp))
    (eval-env env (if-alternative exp))))

(define (make-if predicate consequent alternative)
  (list 'if predicate consequent alternative))

; lambda
(define (make-lambda parameters body)
  (cons 'lambda (cons parameters body)))

(define (lambda? exp) (tagged-list? exp 'lambda))

(define (check-lambda exp)
  (cond ((< (len exp) 3)
         (error 'check-lambda "invalid arity" exp))
        ((and (not (null? (cadr exp)))
              (not (pair? (cadr exp))))
         (error 'check-lambda "invalid argument list" exp))))

(define (lambda-parameters exp) (cadr exp))

(define (lambda-body exp) (cddr exp))

(define (make-procedure parameters body env)
  (list 'compound parameters body env))

; begin
(define (begin? exp) (tagged-list? exp 'begin))

(define (check-begin exp)
  (cond ((< (len exp) 2)
         (error 'check-begin "invalid arity" exp))))

(define (begin-actions exp) (cdr exp))

(define (eval-sequence env seq)
  (cond ((null? seq)
         (error 'eval-sequence "unspecified sequence" seq))
        ((null? (cdr seq))
         (eval-env env (car seq)))
        (else
          (eval-env env (car seq))
          (eval-sequence env (cdr seq)))))

(define (make-begin seq) (cons 'begin seq))

(define (sequence->exp seq)
  (cond ((null? seq) '())
        ((null? (cdr seq)) (car seq))
        (else (make-begin seq))))

; cond
(define (cond? exp) (tagged-list? exp 'cond))

(define (cond-clauses exp) (cdr exp))

(define (cond-else-clause? clause) (eq? (cond-predicate clause) 'else))

(define (cond-actions clause) (cdr clause))

(define (cond-predicate clause) (car clause))

(define (expand-clauses clauses)
  (if (null? clauses)
    'true
    (let ((first (car clauses))
          (the-rest (cdr clauses)))
      (cond ((not (pair? first))
             (error 'expand-clauses "invalid syntax" first))
            ((cond-else-clause? first)
             (if (null? the-rest)
               (sequence->exp (cond-actions first))
               (error 'expand-clauses "else clause isn't last" clauses)))
            (else (make-if (cond-predicate first)
                           (sequence->exp (cond-actions first))
                           (expand-clauses the-rest)))))))

(define (cond->if exp) (expand-clauses (cond-clauses exp)))

; and
(define (and? exp) (tagged-list? exp 'and))

(define (and-expressions exp) (cdr exp))

(define (expand-and exps)
  (if (null? exps)
    'true
    (let ((first (car exps))
          (the-rest (cdr exps)))
      (if (null? the-rest)
        first
        (list (list 'lambda
                    '(predicate)
                    (list 'if
                          'predicate
                          (expand-and the-rest)
                          'predicate))
              first)))))

(define (and->if exp) (expand-and (and-expressions exp)))

; or
(define (or? exp) (tagged-list? exp 'or))

(define (or-expressions exp) (cdr exp))

(define (expand-or exps)
  (if (null? exps)
    'false
    (let ((first (car exps))
          (the-rest (cdr exps)))
      (if (null? the-rest)
        first
        (list (list 'lambda
                    '(predicate)
                    (list 'if
                          'predicate
                          'predicate
                          (expand-or the-rest)))
              first)))))

(define (or->if exp) (expand-or (or-expressions exp)))

; let
(define (let? exp) (tagged-list? exp 'let))

(define (check-let exp)
  (cond ((< (len exp) 3)
         (error 'check-let "invalid arity" exp))
        ((not (pair? (cadr exp)))
         (error 'check-let "invalid syntax" exp))))

(define (let-defs exp) (cadr exp))

(define (let-variables defs) (map car defs))

(define (let-values defs) (map cadr defs))

(define (let-body exp) (cddr exp))

(define (let->procedure exp)
  (cons (make-lambda (let-variables (let-defs exp))
                     (let-body exp))
        (let-values (let-defs exp))))

; application
(define (application? exp) (pair? exp))

(define (operator exp) (car exp))

(define (operands exp) (cdr exp))

(define (no-operands? ops) (null? ops))

(define (first-operand ops) (car ops))

(define (rest-operands ops) (cdr ops))

(define (compound-procedure? p) (tagged-list? p 'compound))

(define (procedure-parameters p) (cadr p))

(define (procedure-body p) (caddr p))

(define (procedure-environment p) (cadddr p))

(define (list-of-values env exps)
  (if (no-operands? exps)
    '()
    (cons (eval-env env (first-operand exps))
          (list-of-values env (rest-operands exps)))))

(define (extend-procedure-environment env names values)
  (define (define-args env names values)
    (cond ((null? names)
           (cond ((not (null? values))
                  (error 'procedure-environment
                         "too many values"
                         (list names values)))))
          ((null? values)
           (error 'procedure-environment
                  "too many names"
                  (list names values)))
          (else
            (env (car names) true (car values))
            (define-args env
                         (cdr names)
                         (cdr values)))))
  (let ((env (extend-env env)))
    (define-args env names values)
    env))

; eval/apply
(define (apply proc args)
  (cond ((compiled-procedure? proc)
         (capply proc args))
        ((compound-procedure? proc)
         (eval-sequence
           (extend-procedure-environment (procedure-environment proc)
                                         (procedure-parameters proc)
                                         args)
           (procedure-body proc)))
        (else (error 'apply "unknown procedure type" proc))))

(define (eval-env env exp)
  (cond ((self-evaluating? exp) exp)
        ((symbol? exp) (env exp))
        ((vector? exp) exp)
        ((quote? exp) (check-quote exp)
                      (quote-text exp))
        ((assignment? exp) (check-assignment exp)
                           (eval-assignment env exp))
        ((definition? exp) (check-definition exp)
                           (eval-definition env exp))
        ((if? exp) (check-if exp)
                   (eval-if env exp))
        ((lambda? exp) (check-lambda exp)
                       (make-procedure
                         (lambda-parameters exp)
                         (lambda-body exp)
                         env))
        ((begin? exp) (check-begin exp)
                      (eval-sequence env (begin-actions exp)))
        ((cond? exp) (eval-env env (cond->if exp)))
        ((and? exp) (eval-env env (and->if exp)))
        ((or? exp) (eval-env env (or->if exp)))
        ((let? exp) (check-let exp)
                    (eval-env env (let->procedure exp)))
        ((application? exp)
         (apply (eval-env env (operator exp))
                (list-of-values env (operands exp))))
        (else
          (error 'eval-env "unknown expression type" exp))))

(define (eval exp)
  (eval-env mikkamakka exp))

; compile
(define (compile-pair context p)
  (compile-write context "[")
  (compile-literal context (car p))
  (compile-write context ",")
  (compile-literal context (cdr p))
  (compile-write context "]"))

(define (compile-vector context v)
  (compile-write context "{vector:[")
  (let ((ref 0))
    (for (lambda ()
           (if (eq? ref (vlen v))
             true
             (begin
               (cond ((> ref 0) (compile-write context ",")))
               (compile-exp context (vref v ref))
               (set! ref (+ ref 1))
               false)))))
  (compile-write context "]}"))

(define (compile-literal context exp)
  (cond ((self-evaluating? exp) (compile-write context (sprint exp)))
        ((symbol? exp) (compile-write context (cats "[" (sprint (symbol-name exp)) "]")))
        ((null? exp) (compile-write context "[]"))
        ((pair? exp) (compile-pair context exp))
        ((vector? exp) (compile-vector context exp))
        (else (error 'compile-literal "unknown type" exp))))

(define (compile-variable context exp)
  (compile-write context "env(")
  (compile-write context (sprint (symbol-name exp)))
  (compile-write context ")")
  (add-compile-dependency context exp))

(define (compile-assignment context exp)
  (compile-write context "env(")
  (compile-write context (sprint (symbol-name (assignment-variable exp))))
  (compile-write context ",false,")
  (compile-exp context (assignment-value exp))
  (compile-write context ")"))

(define (compile-definition context exp)
  (add-compile-definition context (definition-variable exp) exp)
  (push-compile-scope context (definition-variable exp))
  (compile-write context "env(")
  (compile-write context (sprint (symbol-name (definition-variable exp))))
  (compile-write context ",true,")
  (compile-exp context (definition-value exp))
  (compile-write context ")")
  (pop-compile-scope context))

(define (compile-if context exp)
  (compile-write context "(false!==(")
  (compile-exp context (if-predicate exp))
  (compile-write context ")?")
  (compile-exp context (if-consequent exp))
  (compile-write context ":")
  (compile-exp context (if-alternative exp))
  (compile-write context ")"))

(define (compile-sequence context seq)
  (cond ((null? seq)
         (error 'compile-sequence "unspecified sequence" seq))
        ((null? (cdr seq))
         (compile-write context "return ")
         (compile-statement context (car seq)))
        (else
          (compile-statement context (car seq))
          (compile-sequence context (cdr seq)))))

(define (compile-procedure context parameters body)
  (define (compile-parameters params params-out count)
    (if (null? params)
      (begin
        (compile-write context "(arguments.length===")
        (compile-write context (sprint count))
        (compile-write context ")?false:env(\"cerror\")(\"arguments\",\"invalid arity\",")
        (compile-write context (sprint count))
        (compile-write context ");")
        (compile-write context (builder->string params-out)))
      (begin
        (add-compile-definition context (car params) 'argument)
        (compile-parameters (cdr params)
                            (sbappend
                              (sbappend
                                (sbappend
                                  (sbappend
                                    (sbappend
                                      (sbappend params-out "env(")
                                      (sprint (symbol-name (car params))))
                                    ",true,")
                                  "arguments[")
                                (sprint count))
                              "]);")
                            (+ count 1)))))
  (compile-write context "(function(penv){")
  (compile-write context "return function(){")
  (compile-write context "var env=penv(\"extend-env\")(penv);")
  (compile-parameters parameters (make-string-builder) 0)
  (compile-sequence context body)
  (compile-write context "};})(env)"))

(define (compile-begin context actions)
  (compile-write context "(function(){")
  (compile-sequence context actions)
  (compile-write context "})()"))

(define (compile-application context proc args)
  (define (compile-args args first)
    (cond ((not (null? args))
           (cond ((not first) (compile-write context ",")))
           (compile-exp context (car args))
           (compile-args (cdr args) false))))
  (compile-write context "(")
  (compile-exp context proc)
  (compile-write context ")(")
  (compile-args args true)
  (compile-write context ")"))

(define (compile-exp context exp)
  (cond ((self-evaluating? exp) (compile-literal context exp))
        ((symbol? exp) (compile-variable context exp))
        ((vector? exp) (compile-literal context exp))
        ((quote? exp) (compile-literal context (quote-text exp)))
        ((assignment? exp) (check-assignment exp)
                           (compile-assignment context exp))
        ((definition? exp) (check-definition exp)
                           (compile-definition context exp))
        ((if? exp) (check-if exp)
                   (compile-if context exp))
        ((lambda? exp) (check-lambda exp)
                       (compile-procedure
                         context
                         (lambda-parameters exp)
                         (lambda-body exp)))
        ((begin? exp) (check-begin exp)
                      (compile-begin
                        context
                        (begin-actions exp)))
        ((cond? exp) (compile-exp context (cond->if exp)))
        ((and? exp) (compile-exp context (and->if exp)))
        ((or? exp) (compile-exp context (or->if exp)))
        ((let? exp) (check-let exp)
                    (compile-exp context (let->procedure exp)))
        ((application? exp)
         (compile-application context
                              (operator exp)
                              (operands exp)))
        (else
          (error 'compile-exp "invalid expression" exp))))

(define (compile-statement context exp)
  (compile-exp context exp)
  (compile-write context ";"))

(define (make-compile-scope name)
  (let ((scope (make-name-table)))
    (table-define scope 'name name)
    (table-define scope 'definitions (make-name-table))
    (table-define scope 'dependencies (make-name-table))
    scope))

(define (make-name-gen)
  (define (inc base val)
    (if (null? val)
      '(0)
      (let ((lowest (+ (car val) 1)))
        (if (< lowest base)
          (cons lowest (cdr val))
          (cons 0 (inc base (cdr val)))))))
  (define (print digits val)
    (if (null? val)
      (make-string-builder)
      (sbappend (print digits (cdr val)) (char-at digits (car val)))))
  (let ((digits "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
        (ref '()))
    (lambda ()
      (set! ref (inc (slen digits) ref))
      (builder->string (print digits ref)))))

(define (make-name-map)
  (let ((name-gen (make-name-gen))
        (names (make-name-table)))
    (lambda (name)
      (if (table-has-name? names name)
        (table-lookup names name)
        (let ((mapped (name-gen)))
          (table-define names name mapped)
          mapped)))))

(define (make-compile-context names)
  (let ((context (make-name-table)))
    (table-define context 'names names)
    (table-define context 'expressions (make-name-table))
    (table-define context 'stack (list (make-compile-scope 'top)))
    (table-define context 'output (make-string-builder))
    context))

(define (context-names context)
  (table-lookup context 'names))

(define (map-name context name)
  ((context-names context) name))

(define (scope-name scope)
  (table-lookup scope 'name))

(define (compile-expressions context)
  (table-lookup context 'expressions))

(define (compile-stack context)
  (table-lookup context 'stack))

(define (compile-scope context)
  (car (compile-stack context)))

(define (scope-dependencies scope)
  (table-lookup scope 'dependencies))

(define (compile-dependencies context)
  (scope-dependencies (compile-scope context)))

(define (scope-definitions scope)
  (table-lookup scope 'definitions))

(define (compile-definitions context)
  (scope-definitions (compile-scope context)))

(define (add-compile-dependency context name)
  (table-define (compile-dependencies context) name true))

(define (add-compile-definition context name exp)
  (cond ((eq? (scope-name (compile-scope context)) 'top)
         (table-define (compile-expressions context) name exp)))
  (table-define (compile-definitions context)
                name
                (make-name-table)))

(define (push-compile-scope context name)
  (table-set! context
              'stack
              (cons (make-compile-scope name)
                    (compile-stack context))))

(define (pop-compile-scope context)
  (let ((current-scope (compile-scope context)))
    (let ((parent-definition
            (table-lookup (scope-definitions
                            (cadr (compile-stack context)))
                          (scope-name current-scope)))
          (current-definitions (scope-definitions current-scope)))
      (define (undefined? dep) (not (table-has-name? current-definitions dep)))
      (define (copy-def-deps def-names)
        (cond ((not (null? def-names))
               (copy-table-filtered (table-lookup current-definitions
                                                  (car def-names))
                                    parent-definition
                                    undefined?)
               (copy-def-deps (cdr def-names)))))
      (copy-def-deps (table-names current-definitions))
      (copy-table-filtered (scope-dependencies current-scope)
                           parent-definition
                           undefined?)
      (table-set! context 'stack (cdr (compile-stack context))))))

(define (compile-buffer context)
  (table-lookup context 'output))

(define (compile-write context s)
  (table-set! context
              'output
              (sbappend (compile-buffer context) s)))

(define (precompile-exp exp)
  (list 'set!
        'shared-precompiled
        (list 'cons
              (list 'quote exp)
              'shared-precompiled)))

(define (compile-head clb)
  (read-file "mm-head.js"
             (lambda (data)
               (out data)
               (clb))))

(define (compile-shared context)
  (define (iterate context shared)
    (cond ((not (null? shared))
           (compile-statement context (car shared))
           (iterate context (cdr shared)))))
  (cond ((defined? mikkamakka 'shared-precompiled)
         (iterate context shared-precompiled))))

; read
(define tokenizer-expression
  '(";[^\\n]*\\n?|"                    ; comment
    "\\(|\\)|"                         ; list open, list/vector close
    "#\\(|"                            ; vector open
    "'|"                               ; quote
    "\"(\\\\\\\\|\\\\\"|[^\"])*\"?|"   ; string
    "(\\\\.|"                          ; symbol: single escape
    "\\|(\\\\\\\\|\\\\\\||[^|])*\\|?|" ; symbol: range escape
    "[^;()'|\"\\s])+"))                ; symbol: no comment/list/type-escape/quote/string/whitespace

(define tokenizer-rx
  (make-regexp (apply cats tokenizer-expression) "g"))

(define token-complete-expression
  '("^(;[^\\n]*\\n|"                        ; comment
    "\\(|"                                  ; list open
    "\\)|"                                  ; list/vector close
    "#\\(|"                                 ; vector open
    "'|"                                    ; quote
    "\"(\\\\\"|\\\\[^\"]|[^\\\\\"])*\"|"    ; string
    "(\\\\.|"                               ; no single escape
    "\\|(\\\\\\||\\\\[^\\|]|[^\\\\|])*\\||" ; no range escape
    "[^;()#'|\"\\s\\\\])+)$"))              ; symbol: no comment/list/type-escape/quote/string/whitespace/escape

(define token-complete-rx
  (make-regexp (apply cats token-complete-expression) ""))

(define (make-read-context)
  (let ((context (make-name-table))
        (tokens (make-name-table)))
    (table-define tokens 'vectors '())
    (table-define tokens 'ref 0)
    (table-define context 'tokens tokens)
    (table-define context 'stack '())
    (table-define context 'expressions '())
    context))

(define (read-tokens context)
  (table-lookup context 'tokens))

(define (read-tokens-list context)
  (table-lookup (read-tokens context) 'vectors))

(define (read-token-ref context)
  (table-lookup (read-tokens context) 'ref))

(define (read-stack context)
  (table-lookup context 'stack))

(define (read-stack-empty? context)
  (null? (read-stack context)))

(define (current-stack-frame context)
  (car (read-stack context)))

(define (current-read-list-type context)
  (table-lookup (current-stack-frame context) 'type))

(define (current-read-list context)
  (table-lookup (current-stack-frame context) 'list))

(define (read-expressions context)
  (table-lookup context 'expressions))

(define (token-complete? token)
  (> (vlen (token-complete-rx token)) 0))

(define (unescape-char char-symbol)
  (cond ((peq? char-symbol "b") "\b")
        ((peq? char-symbol "t") "\t")
        ((peq? char-symbol "n") "\n")
        ((peq? char-symbol "v") "\v")
        ((peq? char-symbol "f") "\f")
        ((peq? char-symbol "r") "\r")
        (else char-symbol)))

(define (unescape string)
  (define (unescape builder s)
    (for (lambda ()
           (let ((eref (sidx s "\\\\")))
             (cond ((< eref 0)
                    (set! builder (sbappend builder s))
                    true)
                   ((eq? (slen s) (+ eref 1))
                    (error 'unescape "invalid escape sequence" string))
                   (else
                     (set! builder
                       (sbappend builder (subs s 0 eref)))
                     (set! builder
                       (sbappend builder (unescape-char (char-at s (+ eref 1)))))
                     (set! s (subs s (+ eref 2) -1))
                     false)))))
    (builder->string builder))
  (unescape (make-string-builder) string))

(define (unescape-symbol name)
  (define (unescape builder s)
    (let ((escaped false))
      (for (lambda ()
             (let ((eref (sidx s "\\\\|\\|")))
               (cond ((< eref 0)
                      (set! builder (sbappend builder s))
                      true)
                     (else
                       (set! builder
                         (sbappend builder (subs s 0 eref)))
                       (let ((echar (char-at s eref)))
                         (cond ((peq? echar "|")
                                (set! s (subs s (+ eref 1) -1))
                                (set! escaped (not escaped))
                                false)
                               ((eq? (slen s) (+ eref 1))
                                (set! builder (sbappend builder "\\"))
                                true)
                               ((eq? (char-at s (+ eref 1)) "|")
                                (set! builder (sbappend builder "|"))
                                (set! s (subs s (+ eref 2) -1))
                                false)
                               (escaped
                                 (set! builder (sbappend builder (subs s eref 2)))
                                 (set! s (subs s (+ eref 2) -1))
                                 false)
                               (else
                                 (set! builder (sbappend builder (char-at s (+ eref 1))))
                                 (set! s (subs s (+ eref 2) -1))
                                 false)))))))))
    (builder->string builder))
  (unescape (make-string-builder) name))

(define (push-list-type context type)
  (let ((frame (make-name-table)))
    (table-define frame 'type type)
    (table-define frame 'list '())
    (table-set! context
                'stack
                (cons frame (read-stack context)))))

(define (push-list context)
  (push-list-type context 'list))

(define (push-vector context)
  (push-list-type context 'vector))

(define (replace-read-list context list)
  (table-set! (current-stack-frame context) 'list list))

(define (cons-read-expression context tag)
  (table-set! context
              'expressions
              (cons tag (read-expressions context))))

(define (cons-read-tag context tag)
  ((if (read-stack-empty? context)
     cons-read-expression
     cons-stack-tag)
   context
   tag))

(define (close-current-list context)
  (if (read-stack-empty? context)
    (error 'close-current-list "unexpected closing paren" context)
    (let ((closed (if (eq? (current-read-list-type context) 'vector)
                    (apply vector (reverse (current-read-list context)))
                    (reverse (current-read-list context)))))
      (table-set! context 'stack (cdr (read-stack context)))
      (cons-read-tag context closed))))

(define (push-quote context)
  (push-list-type context 'quote))

(define (quoted-tag? context)
  (eq? (current-read-list-type context) 'quote))

(define (cons-stack-tag context tag)
  (let ((list (cons tag
                    (if (quoted-tag? context)
                      (cons 'quote (current-read-list context))
                      (current-read-list context)))))
    (replace-read-list context list)
    (cond ((quoted-tag? context)
           (close-current-list context)))))

(define (cons-string context token)
  (cons-read-tag context
                 (unescape (subs token
                                 1
                                 (- (slen token)
                                    2)))))

(define (cons-token context token)
  (let ((number (parse-number token)))
    (cons-read-tag
      context
      (if (number? number)
        number
        (string->symbol (unescape-symbol token))))))

(define (parse-token context token)
  (let ((char0 (char-at token 0)))
    (cond ((peq? token "(") (push-list context))
          ((peq? token "#(") (push-vector context))
          ((peq? token ")") (close-current-list context))
          ((peq? token "'") (push-quote context))
          ((peq? (char-at token 0) "\"") (cons-string context token))
          ((peq? (char-at token 0) "#")
           (error 'parse-token "unrecognized type escape" token))
          ((not (peq? (char-at token 0) ";")) (cons-token context token)))))

(define (shift-tokens context)
  (table-set! (read-tokens context) 'vectors '())
  (table-set! (read-tokens context) 'ref 0))

(define (parse-tokens context)
  (define (next-tokens tokens)
    (cond ((null? (cdr tokens))
           '())
          ((eq? (vlen (cadr tokens)) 0)
           (next-tokens (cdr tokens)))
          (else (cdr tokens))))
  (let ((tokens (read-tokens-list context))
        (ref (read-token-ref context)))
    (define (parse-incomplete token)
      (let ((next (next-tokens tokens))
            (exit? false))
        (for (lambda ()
               (cond ((null? next)
                      (set! exit? true)
                      true)
                     (else
                       (set! token
                         (sbappend token
                                   (if (eq? (vref (car next) 0) 'eof)
                                     "\n"
                                     (vref (car next) 0))))
                       (let ((token-string (builder->string token)))
                         (cond ((token-complete? token-string)
                                (parse-token context token-string)
                                (set! tokens next)
                                (set! ref 1)
                                true)
                               (else
                                 (set! next (next-tokens next))
                                 false)))))))
        exit?))
    (for (lambda ()
           (cond ((null? tokens)
                  (shift-tokens context)
                  true)
                 ((eq? (vlen (car tokens)) ref)
                  (set! tokens (next-tokens tokens))
                  (set! ref 0)
                  false)
                 (else
                   (let ((token (vref (car tokens) ref)))
                     (cond ((eq? token 'eof)
                            (shift-tokens context)
                            true)
                           ((token-complete? token)
                            (parse-token context token)
                            (set! ref (+ ref 1))
                            false)
                           (else
                             (parse-incomplete (sbappend (make-string-builder)
                                                         token)))))))))))

(define (append-tokens context tokens)
  (table-set! (read-tokens context)
              'vectors
              (append (read-tokens-list context)
                      (list tokens))))

(define (read-string context string)
  (let ((tokens (tokenizer-rx string)))
    (append-tokens context tokens)
    (parse-tokens context)))

(define (shift-expressions context)
  (let ((expressions (read-expressions context)))
    (table-set! context 'expressions '())
    (reverse expressions)))

(define (append-token context token)
  (append-tokens context (vector token)))

(define (read-complete? context)
  (and (null? (read-tokens-list context))
       (null? (read-stack context))))

(define (complete-read context)
  (append-token context 'eof)
  (parse-tokens context)
  (cond ((not (read-complete? context))
         (error 'complete-read
                "incomplete expression"
                (list (read-tokens-list context)
                      (read-stack context))))
        (else (shift-expressions context))))

(define (complete-read context)
  (append-token context 'eof)
  (parse-tokens context)
  (cond ((not (null? (read-tokens-list context)))
         (error 'complete-read
                "incomplete expression"
                (read-tokens-list context)))
        ((not (null? (read-stack context)))
         (error 'complete-read
                "incomplete expression"
                (read-stack context)))
        (else (shift-expressions context))))

(define (sread string)
  (let ((context (make-read-context)))
    (read-string context string)
    (complete-read context)))

; run
(define (command? arg)
  (or (peq? arg "repl")
      (peq? arg "run")
      (peq? arg "test")
      (peq? arg "compile")))

(define (explicit-command? argv)
  (command? (caddr argv)))

(define (compiled-run?)
  (compiled-procedure? get-command))

(define (get-command-from argv)
  (cond ((null? (cddr argv)) "repl")
        ((command? (caddr argv)) (caddr argv))
        (else "run")))

(define (get-arg-from argv)
  (cond ((null? (cddr argv))
         (error 'get-arg-from "missing argument" argv))
        ((explicit-command? argv)
         (if (null? (cdddr argv))
           (error 'get-arg-from "missing argument" argv)
           (cadddr argv)))
        (else (caddr argv))))

(define (argv-member op)
  (let ((argv (proc-argv)))
    (cond ((compiled-run?) (op argv))
          ((explicit-command? argv) (op (cddr argv)))
          (else (op (cdr argv))))))

(define (get-command)
  (argv-member get-command-from))

(define (get-arg)
  (argv-member get-arg-from))

(define (repl)
  (define env (create-env))
  (define (try-eval-sequence exps)
    (cond ((not (null? exps))
           (try
             (lambda ()
               (print (eval-env env (car exps))))
             (lambda (err)
               (display err)))
           (newline)
           (try-eval-sequence (cdr exps)))))
  (define read-context (make-read-context))
  (define (process-line rl line)
    (read-string read-context line)
    (try-eval-sequence (shift-expressions read-context))
    (set-prompt rl
                (if (read-complete? read-context)
                  "> "
                  ". "))
    (prompt rl))
  (define rl (read-line process-line "> "))
  (set-prompt rl "> ")
  (prompt rl))

(define (run)
  (define (read-and-eval scm)
    (let ((ctx (make-read-context)))
      (read-string ctx scm)
      (let ((exps (shift-expressions ctx)))
        (complete-read ctx)
        (eval-env (create-env)
                  (list (cons 'lambda (cons '() exps)))))))
  (read-file (get-arg) read-and-eval))

(define (compile)
  (define (precompile head done ccontext pcontext shared)
    (let ((shared-names (table-names shared)))
      (cond ((not (null? shared-names))
             (let ((name (car shared-names)))
               (cond ((and (not (table-has-name? done name))
                           (table-has-name?
                             (compile-expressions ccontext)
                             name))
                      (compile-statement
                        pcontext
                        (precompile-exp (table-lookup
                                          (compile-expressions ccontext)
                                          name)))
                      (copy-table-filtered
                        (table-lookup
                          (compile-definitions ccontext)
                          name)
                        shared
                        (lambda (name)
                          (not (or (table-has-name? done name)
                                   (defined? head name)))))))
               (table-define done name true)
               (table-delete! shared name)
               (precompile head done ccontext pcontext shared))))))
  (define (read-and-compile names scm)
    (let ((read-context (make-read-context))
          (ccontext (make-compile-context names)))
      (read-string read-context scm)
      (let ((exps (shift-expressions read-context)))
        (complete-read read-context)
        (compile-sequence ccontext exps)
        (let ((pcontext (make-compile-context names)))
          (compile-statement pcontext '(define shared-precompiled '()))
          (precompile (extend-env head)
                      (make-name-table)
                      ccontext
                      pcontext
                      (to-table (shared)
                                (lambda (i) (car i))
                                identity))
          (out (builder->string (compile-buffer pcontext)))
          (out (builder->string (compile-buffer ccontext)))))))
  (compile-head
    (lambda ()
      (let ((names (make-name-map)))
        (let ((ccontext (make-compile-context names)))
          (compile-shared ccontext)
          (out (builder->string (compile-buffer ccontext)))
          (read-file (get-arg) (lambda (scm)
                                 (read-and-compile names scm))))))))

(define (caar l) (car (car l)))
(define (cadr l) (car (cdr l)))
(define (cdar l) (cdr (car l)))
(define (cddr l) (cdr (cdr l)))
(define (caaar l) (car (car (car l))))
(define (caadr l) (car (car (cdr l))))
(define (cadar l) (car (cdr (car l))))
(define (cdaar l) (cdr (car (car l))))
(define (caddr l) (car (cdr (cdr l))))
(define (cdadr l) (cdr (car (cdr l))))
(define (cddar l) (cdr (cdr (car l))))
(define (cdddr l) (cdr (cdr (cdr l))))
(define (caaaar l) (car (car (car (car l)))))
(define (caaadr l) (car (car (car (cdr l)))))
(define (caadar l) (car (car (cdr (car l)))))
(define (cadaar l) (car (cdr (car (car l)))))
(define (cdaaar l) (cdr (car (car (car l)))))
(define (caaddr l) (car (car (cdr (cdr l)))))
(define (cadadr l) (car (cdr (car (cdr l)))))
(define (cdaadr l) (cdr (car (car (cdr l)))))
(define (caddar l) (car (cdr (cdr (car l)))))
(define (cdadar l) (cdr (car (cdr (car l)))))
(define (cddaar l) (cdr (cdr (car (car l)))))
(define (cadddr l) (car (cdr (cdr (cdr l)))))
(define (cdaddr l) (cdr (car (cdr (cdr l)))))
(define (cddadr l) (cdr (cdr (car (cdr l)))))
(define (cdddar l) (cdr (cdr (cdr (car l)))))
(define (cddddr l) (cdr (cdr (cdr (cdr l)))))

(define (shared)
  (list (list 'eval eval)
        (list 'apply apply)
        (list 'log log)
        (list 'print print)
        (list 'eq? eq?)
        (list 'newline newline)
        (list 'error error)
        (list 'not not)
        (list 'make-name-table make-name-table)
        (list 'table-define table-define)
        (list 'table-lookup table-lookup)
        (list 'table-set! table-set!)
        (list 'append append)
        (list 'sprint sprint)
        (list 'reverse reverse)
        (list 'display display)
        (list 'sescape sescape)
        (list 'caar caar)
        (list 'cadr cadr)
        (list 'cdar cdar)
        (list 'cddr cddr)
        (list 'caaar caaar)
        (list 'caadr caadr)
        (list 'cadar cadar)
        (list 'cdaar cdaar)
        (list 'caddr caddr)
        (list 'cdadr cdadr)
        (list 'cddar cddar)
        (list 'cdddr cdddr)
        (list 'caaaar caaaar)
        (list 'caaadr caaadr)
        (list 'caadar caadar)
        (list 'cadaar cadaar)
        (list 'cdaaar cdaaar)
        (list 'caaddr caaddr)
        (list 'cadadr cadadr)
        (list 'cdaadr cdaadr)
        (list 'caddar caddar)
        (list 'cdadar cdadar)
        (list 'cddaar cddaar)
        (list 'cadddr cadddr)
        (list 'cdaddr cdaddr)
        (list 'cddadr cddadr)
        (list 'cdddar cdddar)
        (list 'cddddr cddddr)))

(define (create-env)
  (define (list->table t l)
    (cond ((null? l) t)
          (else
            (table-define t (caar l) (cadar l))
            (list->table t (cdr l)))))
  (extend-env mikkamakka
              (list->table
                (make-name-table)
                (shared))))

(let ((cmd (get-command)))
  (cond ((peq? cmd "repl")
         (repl))
        ((peq? cmd "run")
         (run))
        ((peq? cmd "compile")
         (compile))))

; test
(define last-test-name false)

(define (string-contains? val pattern)
  (>= (sidx (sprint val) pattern) 0))

(define (list-eq? left right)
  (cond ((or (not (pair? left))
             (not (pair? right)))
         (eq? left right))
        ((not (list-eq? (car left) (car right))) false)
        (else (list-eq? (cdr left) (cdr right)))))

(define (assert val msg)
  (if val
    val
    (begin (print (cats last-test-name ": failed: " msg))
           (exit -1))))

(define (fail msg f)
  (try (lambda () (f) (assert false msg))
       (lambda (error)
         (cond ((not (string-contains? error msg))
                (print-stack error)
                (assert false msg))))))

(define (test name f)
  (set! last-test-name name)
  (let ((env (create-env)))
    (try (lambda () (f (lambda (exp) (eval-env env exp))))
         (lambda (error)
           (print (cats name ": error during test:"))
           (print-stack error)
           (exit)))))

(define (run-tests)
  (test "symbol?"
        (lambda (_)
          (assert (symbol? 'a) "true")
          (assert (not (symbol? "a")) "false")))

  (test "len"
        (lambda (_)
          (assert (eq? (len '(1 2 3)) 3) "3")
          (assert (eq? (len '()) 0) "0")))

  (test "quote-eq?"
        (lambda (_)
          (assert (quote-eq? ''a ''a) "eq")
          (assert (not (quote-eq? ''a ''b)) "not eq")))

  (test "eq"
        (lambda (_)
          (assert (eq? '() '()) "null")
          (assert (eq? 'a 'a) "symbol")
          (assert (eq? ''a ''a) "quote")
          (assert (eq? 1 1) "primitive")
          (let ((a (lambda () 1)))
            (let ((b a))
              (assert (eq? a b) "reference")))))

  (test "tagged-list?"
        (lambda (_)
          (assert (tagged-list? (list 'a) 'a) "true")
          (assert (not (tagged-list? (list 'a) 'b)) "wrong tag")
          (assert (not (tagged-list? 1 'a)) "not list")))

  (test "check-quote"
        (lambda (_)
          (assert (eq? (check-quote (list 'quote 'a)) true) "ok")))
  ; (fail "invalid quotation" (lambda () (check-quote (list 'quote))))))

  (test "escape"
        (lambda (_)
          (assert (eq? (sescape "\b\t\v\f\n\r\\\"")
                       "\\b\\t\\v\\f\\n\\r\\\\\\\"")
                  "escape")
          (assert (list-eq? (sread "\"\b\t\v\f\n\r\\\\\\\"\"")
                            '("\b\t\v\f\n\r\\\""))
                  "read")
          (assert (eq? (sunescape "\\b\\t\\v\\f\\n\\r\\\\\\\"")
                       "\b\t\v\f\n\r\\\"")
                  "unescape")
          (assert (list-eq? (sread "\"\\b\\t\\v\\f\\n\\r\\\\\\\"\"")
                            '("\b\t\v\f\n\r\\\""))
                  "read escaped")
          (assert (eq? (sescape-symbol-name "a")
                       "a")
                  "symbol name normal")
          (assert (eq? (sescape-symbol-name "a\na")
                       "|a\na|")
                  "symbol name escaped")))

  (test "sprint"
        (lambda (_)
          (assert (eq? (sprint 1) "1") "number")
          (assert (eq? (sprint "some string") "\"some string\"") "string")
          (assert (eq? (sprint "some string with \"apostrophs\" in it")
                       "\"some string with \\\"apostrophs\\\" in it\"")
                  "string, apostrophs")
          (assert (eq? (sprint 'a) "'a") "quote")
          (assert (eq? (sprint '(a b c)) "'(a b c)") "list")
          (assert (eq? (sprint '(a b (c 'd))) "'(a b (c 'd))") "list, embedded list")
          (assert (eq? (sprint '(a b '(c 'd))) "'(a b '(c 'd))") "list, embedded quoted list")
          (assert (eq? (sprint '(a b '(c 'd ''d))) "'(a b '(c 'd ''d))") "list, double quote")))

  (test "self-evaluating?"
        (lambda (_)
          (assert (self-evaluating? false) "false")
          (assert (self-evaluating? true) "true")
          (assert (self-evaluating? 1) "number")
          (assert (self-evaluating? "some string") "string")
          (assert (not (self-evaluating? (lambda () false))) "lambda")))

  (test "check-assignment"
        (lambda (_)
          (assert (eq? (check-assignment (list 'set! 'a 1)) true) "valid")
          (fail "invalid arity" (lambda () (check-assignment (list 'set! 'a))))
          (fail "invalid variable name"
                (lambda () (check-assignment (list 'set! 1 1))))))

  (test "definition-variable"
        (lambda (_)
          (assert
            (eq? (definition-variable (list 'define 'a 1))
                 'a)
            "variable")
          (assert
            (eq? (definition-variable (list 'define (list 'a) 1))
                 'a)
            "procedure")))

  (test "definition-value"
        (lambda (_)
          (assert (eq? (definition-value (list 'define 'a 1)) 1) "variable")
          (assert
            (list-eq?
              (definition-value (list 'define (list 'a) 1 2))
              (list 'lambda '() 1 2))
            "procedure")))

  (test "check-definition"
        (lambda (_)
          (assert (eq? (check-definition (list 'define 'a 1))
                       true)
                  "variable")
          (assert (eq? (check-definition (list 'define (list 'a) 1))
                       true)
                  "procedure")
          (fail "invalid format"
                (lambda () (check-definition (list 'define 'a))))
          (fail "invalid variable name"
                (lambda () (check-definition (list 'define (list 1) 1))))))

  (test "validate-if"
        (lambda (_)
          (assert (eq? (check-if (list 'if 1 2 3)) true) "valid")
          (fail "invalid arity" (lambda () (check-if (list 'if 1 2 3 4))))))

  (test "check-lambda"
        (lambda (_)
          (assert (eq? (check-lambda (list 'lambda '() 1)) true) "valid")
          (fail "invalid arity"
                (lambda () (check-lambda (list 'lambda '()))))
          (fail "invalid argument list"
                (lambda () (check-lambda (list 'lambda 1 2))))))

  (test "validate-begin"
        (lambda (_)
          (assert (check-begin (list 'begin 1)) "valid")
          (fail "invalid arity"
                (lambda () (check-begin (list 'begin))))))

  (test "sequence->exp"
        (lambda (_)
          (assert (eq? (sequence->exp '()) '()) "null")
          (assert (eq? (sequence->exp '(1)) 1) "single item")
          (assert (list-eq? (sequence->exp '(1 2)) '(begin 1 2)) "sequence")))

  (test "expand-clauses"
        (lambda (_)
          (assert (list-eq?
                    (expand-clauses '((1 1) (2 2) (else 3)))
                    '(if 1 1 (if 2 2 3))) "ok")
          (fail "invalid syntax"
                (lambda () (expand-clauses '(1 (2 2) (else 3)))))
          (fail "else clause isn't last"
                (lambda () (expand-clauses '((1 1) (else 2) (3 3)))))))

  (test "expand-and"
        (lambda (_)
          (assert (list-eq? (expand-and '()) 'true) "null")
          (assert (list-eq? (expand-and '(1 2 3))
                            '((lambda (predicate)
                                (if predicate
                                  ((lambda (predicate)
                                     (if predicate
                                       3
                                       predicate))
                                   2)
                                  predicate))
                              1))
                  "list")))

  (test "expand-or"
        (lambda (_)
          (assert (list-eq? (expand-or '()) 'false) "null")
          (assert (list-eq? (expand-or '(1 2 3))
                            '((lambda (predicate)
                                (if predicate
                                  predicate
                                  ((lambda (predicate)
                                     (if predicate
                                       predicate
                                       3))
                                   2)))
                              1))
                  "list")))

  (test "validate-let"
        (lambda (_)
          (assert (eq? (check-let '(let ((a 1)) a)) true) "ok")
          (fail "invalid arity" (lambda () (check-let '(let ((a 1))))))
          (fail "invalid syntax" (lambda () (check-let '(let 1 1))))))

  (test "let->procedure"
        (lambda (_)
          (assert (list-eq?
                    (let->procedure '(let ((a 1)) a))
                    '((lambda (a) a) 1)) "transform")))

  (test "false, true"
        (lambda (eval)
          (assert (eq? (eval 'false) false) "false")
          (assert (eq? (eval 'true) true) "true")))

  (test "self evaluation"
        (lambda (eval)
          (assert (eq? (eval 1) 1) "number")
          (assert (eq? (eval "some string") "some string") "string")))

  (test "variable"
        (lambda (eval)
          (fail ".*" (lambda () (eval 'a)))
          (assert (eq? (eval '(define a 1)) 1) "definition")
          (assert (eq? (eval 'a) 1) "lookup")))

  (test "quote"
        (lambda (eval)
          (fail "invalid quotation" (lambda () (eval '(quote))))
          (assert (eq? (eval '1) 1) "number")
          (assert (eq? (eval '"some string") "some string") "string")
          (assert (eq? (eval ''a) 'a) "quote")))

  (test "assignment"
        (lambda (eval)
          (assert (eq? (eval '(define a 1)) 1) "define")
          (assert (eq? (eval '(set! a 2)) 2) "assign")
          (assert (eq? (eval 'a) 2) "lookup")
          (assert (eq? (eval '((lambda ()
                                 (define a 3)
                                 (set! a 4)
                                 a))) 4)
                  "in lambda")
          (assert (eq? (eval '(apply (lambda ()
                                       (define a 3)
                                       (set! a 4)
                                       a)
                                     '()))
                       4)
                  "in apply")))

  (test "define"
        (lambda (eval)
          (assert (eq? (eval '(define a 1)) 1) "define variable")
          (assert (eq? (eval 'a) 1) "verify variable")
          (eval '(define (a) 2))))
  ; (assert (eq? (eval '(a)) 2) "call procedure")))

  (test "if"
        (lambda (eval)
          (assert (eq? (eval '(if false 1 0)) 0) "false")
          (assert (eq? (eval '(if true 1 0)) 1) "true")
          (assert (eq? (eval '(if "anything else" 1 0)) 1) "anything else")))

  (test "lambda"
        (lambda (eval)
          (assert (eq? (eval '((lambda () 1))) 1) "create and call lambda")))

  (test "begin"
        (lambda (eval)
          (assert (eq? (eval '(begin (define a 1) 2)) 2) "execute sequence")
          (assert (eq? (eval 'a) 1) "verify side effect")
          (fail "arity" (lambda () (eval '(begin))))))

  (test "cond"
        (lambda (eval)
          (assert (eq? (eval '(cond (true 0) (false 1) (else 2))) 0) "first")
          (assert (eq? (eval '(cond (false 0) (true 1) (else 2))) 1) "second")
          (assert (eq? (eval '(cond (false 0) (false 1) (else 2))) 2) "else")
          (assert (eq? (eval '(cond)) true) "no clauses")
          (assert (eq? (eval '(cond (false))) true) "implicit else clause")
          (fail "invalid syntax" (lambda () (eval '(cond 1))))))

  (test "and"
        (lambda (eval)
          (assert (eq? (eval '(and)) true) "no expression")
          (assert (eq? (eval '(and false)) false) "single, false")
          (assert (eq? (eval '(and true)) true) "single, true")
          (assert (eq? (eval '(and "something")) "something") "single, string")
          (assert (eq? (eval '(and 1 2)) 2) "two")
          (assert (eq? (eval '(and 1 2 3)) 3) "three")
          (assert (eq? (eval '(and 1 2 false)) false) "three, false")
          (assert (eq? (eval '(and 1 false 2)) false) "two, false")
          (assert (eq? (eval '(and false 1 2)) false) "one, false")))

  (test "or"
        (lambda (eval)
          (assert (eq? (eval '(or)) false) "no expression")
          (assert (eq? (eval '(or false)) false) "single, false")
          (assert (eq? (eval '(or true)) true) "single, true")
          (assert (eq? (eval '(or "something")) "something") "something")
          (assert (eq? (eval '(or 1)) 1) "number")
          (assert (eq? (eval '(or 1 2)) 1) "two")
          (assert (eq? (eval '(or 1 2 3)) 1) "three")
          (assert (eq? (eval '(or 1 false)) 1) "first")
          (assert (eq? (eval '(or false 1)) 1) "second")
          (assert (eq? (eval '(or false false 1)) 1) "third")))

  (test "let"
        (lambda (eval)
          (assert (eq? (eval '(let ((a 1) (b 2)) b a)) 1) "sequence")
          (assert (eq? (eval '(let ((a 1) (b 2)) (and a b))) 2) "operation")
          (fail "invalid arity" (lambda () (eval '(let))))
          (fail "invalid syntax" (lambda () (eval '(let 1 2 3))))
          (fail ".*" (lambda () (eval '(let (1) 2 3))))))

  (test "application"
        (lambda (eval)
          (assert (eq? (eval '((lambda () 1))) 1) "lambda")
          (eval '(define (a) 1))
          (assert (eq? (eval '(a)) 1) "procedure from shortcut")
          (eval '(define a (lambda () 2)))
          (assert (eq? (eval '(a)) 2) "procedure from variable")
          (eval '(define (a x) (if x "it's true" "it's false")))
          (assert (eq? (eval '(a true)) "it's true") "with args, true")
          (assert (eq? (eval '(a false)) "it's false") "with args, false")
          (assert (eq? (eval '(a "just anything")) "it's true") "with args, anything")))

  (test "primitive procedures"
        (lambda (eval)
          (assert (eq? (eval '(+ 1 2)) 3) "add")))

  (test "try"
        (lambda (eval)
          (assert (eq? (try (lambda () 1)
                            (lambda (err) 2))
                       1)
                  "try")
          (assert (eq? (try (lambda () (error 'try "err" "test"))
                            (lambda (err) 1))
                       1)
                  "catch")
          (assert (eq? (eval
                         '(try (lambda () 1)
                               (lambda (err) 2)))
                       1)
                  "eval, try")
          (assert (eq? (eval
                         '(try (lambda () (error 'try "err" "test"))
                               (lambda (err) 1)))
                       1)
                  "eval, catch")
          (assert (eq? (eval '(try (lambda () (car '(1 2)))
                                   (lambda (err) 2)))
                       1)
                  "primitive call, try")
          (assert (eq? (eval
                         '(try (lambda () (error 'try "err" "test"))
                               (lambda (error) (car '(1 2)))))
                       1)
                  "primitive call, catch")))

  (test "for"
        (lambda (eval)
          (assert (eq? (for (lambda () true)) true) "exit immediately")
          (assert (eq? (let ((n 0))
                         (for (lambda ()
                                (set! n (+ n 1))
                                (and (> n 2) n))))
                       3)
                  "count to 3")
          (assert (eq? (eval '(for (lambda () true))) true) "eval, exit immediately")
          (assert (eq? (eval '(let ((n 0))
                                (for (lambda ()
                                       (set! n (+ n 1))
                                       (and (> n 2) n)))))
                       3)
                  "eval, count to 3")))

  (test "eval"
        (lambda (eval)
          (assert (eq? (eval '(eval '(car '(1 2)))) 1) "inner eval")))

  (test "apply"
        (lambda (eval)
          (assert (eq? (eval '(apply car '((1 2)))) 1) "inner apply")))

  (test "vector"
        (lambda (eval)
          (assert (vector? #(1 2 3)) "vector literal")
          (assert (vector? (car (sread "#(1 2 3)"))) "read vector")
          (assert (vector? (eval (car (sread "#(1 2 3)")))) "eval vector")
          (assert (eq? (sprint (car (sread "#(1 2 3)"))) "#(1 2 3)") "print vector")
          (assert (eq? (sprint (car (sread "'(a b #(1 2 (1 2 3) c #(1 2 3 4)))")))
                       "'(a b #(1 2 (1 2 3) c #(1 2 3 4)))")
                  "compose vector")))

  'end-tests)

(cond ((peq? (get-command) "test") (run-tests)))
