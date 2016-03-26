(setlocale LC_ALL "")


(define (log . args)
  (write args (current-error-port))
  (format (current-error-port) "\n"))


(define (panic message . args)
  (format (current-error-port) message)
  (format (current-error-port) " ")
  (write args (current-error-port))
  (format (current-error-port) "\n")
  (exit -1))


(define (mkreftable)
  (let ((symbols '())
        (current-ref 0))
    (define (lookup symbol-list symbol)
      (cond ((null? symbol-list) #f)
            ((eq? (caar symbol-list) symbol)
             (cadar symbol-list))
            (else
              (lookup (cdr symbol-list)
                      symbol))))
    (define (create-ref symbol)
      (cond ((lookup symbols symbol)
             (panic "symbol ref exists")))
      (let ((ref current-ref))
        (set! symbols
          (cons (list symbol current-ref)
                symbols))
        (set! current-ref (+ current-ref 1))
        ref))
    (define (symbol-ref symbol-list symbol)
      (let ((ref (lookup symbols symbol)))
        (if ref ref (panic "symbol not found"))))
    (lambda (message symbol)
      (cond ((eq? message 'lookup)
             (lookup symbols symbol))
            ((eq? message 'ref)
             (symbol-ref symbols symbol))
            ((eq? message 'create-ref)
             (create-ref symbol))))))


(define (lookup-ref reftable symbol) (reftable 'lookup symbol))
(define (symbol-ref reftable symbol) (reftable 'ref symbol))
(define (create-ref reftable symbol) (reftable 'create-ref symbol))


(define variable-ref
  (let ((variables (mkreftable)))
    (lambda (name)
      (let ((ref (lookup-ref variables name)))
        (if ref ref (create-ref variables name))))))


(define make-label
  (let ((label-count 0))
    (lambda ()
      (set! label-count (+ label-count 1))
      label-count)))


(define (list-union list1 list2)
  (cond ((null? list1) list2)
        ((memq (car list1) list2) (list-union (cdr list1) list2))
        (else (cons (car list1) (list-union (cdr list1) list2)))))


(define (list-difference list1 list2)
  (cond ((null? list1) '())
        ((memq (car list1) list2) (list-difference (cdr list1) list2))
        (else (cons (car list1) (list-difference (cdr list1) list2)))))


(define (make-instruction-sequence needs modifies statements)
  (list needs modifies statements))


(define (registers-needed instructions)
  (cond ((number? instructions) '())
        ((pair? instructions) (car instructions))
        (else
          (display instructions)
          (newline)
          (panic "something went wrong"))))


(define (registers-modified instructions)
  (cond ((number? instructions) '())
        ((pair? instructions) (cadr instructions))
        (else
          (display instructions)
          (newline)
          (panic "something went wrong"))))


(define (statements instructions)
  (cond ((number? instructions) (list instructions))
        ((pair? instructions) (caddr instructions))
        (else
          (display instructions)
          (newline)
          (panic "something went wrong"))))


(define (empty-instruction-sequence)
  (make-instruction-sequence '() '() '()))


(define (append-instruction-sequences . seqs)
  (define (append-2-sequences instructions1 instructions2)
    (make-instruction-sequence
      (list-union (registers-needed instructions1)
                  (list-difference (registers-needed instructions2)
                                   (registers-modified instructions1)))
      (list-union (registers-modified instructions1)
                  (registers-modified instructions2))
      (append (statements instructions1) (statements instructions2))))
  (define (append-seq-list seqs)
    (if (null? seqs)
      (empty-instruction-sequence)
      (append-2-sequences (car seqs)
                          (append-seq-list (cdr seqs)))))
  (append-seq-list seqs))


(define (needs-register? instructions register)
  (memq register (registers-needed instructions)))


(define (modifies-register? instructions register)
  (memq register (registers-modified instructions)))


(define (preserving regs instructions1 instructions2)
  (if (null? regs)
    (append-instruction-sequences instructions1 instructions2)
    (let ((first-reg (car regs)))
      (if (and (needs-register? instructions2 first-reg)
               (modifies-register? instructions1 first-reg))
        (preserving
          (cdr regs)
          (make-instruction-sequence
            (list-union (list first-reg)
                        (registers-needed instructions1))
            (list-difference (registers-modified instructions1)
                             (list first-reg))
            (append (list (if (eq? first-reg 'env)
                            (list 'saveenv)
                            (list 'save first-reg)))
                    (statements instructions1)
                    (list (if (eq? first-reg 'env)
                            (list 'restoreenv)
                            (list 'restore first-reg)))))
          instructions2)
        (preserving (cdr regs) instructions1 instructions2)))))


(define (compile-linkage linkage)
  (cond ((eq? linkage 'return)
         (make-instruction-sequence
           '(continue) '()
           '((goto (reg continue)))))
        ((eq? linkage 'next)
         (empty-instruction-sequence))
        (else
          (make-instruction-sequence
            '() '()
            (list (list 'goto (list 'label linkage)))))))


(define (end-with-linkage linkage instructions)
  (preserving
    '(continue)
    instructions
    (compile-linkage linkage)))


(define (tack-on-instruction-sequence instructions body-instructions)
  (make-instruction-sequence
    (registers-needed instructions)
    (registers-modified instructions)
    (append (statements instructions) (statements body-instructions))))


(define (self-evaluating? code)
  (or (number? code)
      (string? code)))


(define (compile-self-evaluating code target linkage)
  (end-with-linkage
    linkage
    (make-instruction-sequence
      '() (list target)
      (list (list 'initreg target (list 'const code))))))


(define (quoted? code) (tagged-list? code 'quote))


(define (text-of-quotation code) (cadr code))


(define (compile-quoted code target linkage)
  (end-with-linkage
    linkage
    (make-instruction-sequence
      '()
      (list target)
      (list (list 'initreg
                  target
                  (list 'const
                        (text-of-quotation code)))))))


(define (variable? code) (symbol? code))


(define (compile-variable code target linkage)
  (end-with-linkage
    linkage
    (make-instruction-sequence
      '(env) (list target)
      (list (list 'get-variable target code)))))


(define (assignment? code) (tagged-list? code 'set!))


(define (assignment-variable code) (cadr code))


(define (assignment-value code) (caddr code))


(define (compile-assignment code target linkage)
  (let ((var (assignment-variable code))
        (get-value-code
          (compile (assignment-value code) 'val 'next)))
    (end-with-linkage
      linkage
      (preserving
        '(env)
        get-value-code
        (make-instruction-sequence
          '(env val) (list target)
          (list (list 'set-variable-value target var)))))))


(define (definition? code)
  (tagged-list? code 'define))


(define (definition-variable code)
  (if (symbol? (cadr code))
    (cadr code)
    (caadr code)))


(define (definition-value code)
  (if (symbol? (cadr code))
    (caddr code)
    (make-lambda (cdadr code) (cddr code))))


(define (compile-definition code target linkage)
  (let ((var (definition-variable code))
        (value-code
          (compile (definition-value code) 'val 'next)))
    (end-with-linkage
      linkage
      (preserving
        '(env)
        value-code
        (make-instruction-sequence
          '(env val) (list target)
          (list (list 'define-variable var)))))))


(define (begin? code)
  (tagged-list? code 'begin))


(define (begin-actions code)
  (cdr code))


(define (last-exp? code) (null? (cdr code)))
(define (first-exp code) (car code))
(define (rest-exps code) (cdr code))


(define (compile-sequence code target linkage)
  (if (last-exp? code)
    (compile (first-exp code) target linkage)
    (preserving
      '(env continue)
      (compile (first-exp code) target 'next)
      (compile-sequence (rest-exps code) target linkage))))


(define (tagged-list? code tag)
  (and (pair? code)
       (eq? (car code) tag)))


(define (if? code) (tagged-list? code 'if))


(define (if-predicate code) (cadr code))


(define (if-consequent code) (caddr code))


(define (if-alternative code)
  (if (null? (cdddr code))
    'false
    (cadddr code)))


(define (compile-if code target linkage)
  (let ((t-branch (make-label))
        (f-branch (make-label))
        (after-if (make-label)))
    (let ((consequent-linkage
            (if (eq? linkage 'next) after-if linkage)))
      (let ((p-code (compile (if-predicate code) 'val 'next))
            (c-code (compile (if-consequent code) target consequent-linkage))
            (a-code (compile (if-alternative code) target linkage)))
        (preserving
          '(env continue)
          p-code
          (append-instruction-sequences
            (make-instruction-sequence
              '(val) '()
              (list (list 'branchval f-branch)))
            (parallel-instruction-sequences
              (append-instruction-sequences t-branch c-code)
              (append-instruction-sequences f-branch a-code))
            after-if))))))


(define (lambda? code) (tagged-list? code 'lambda))


(define (lambda-parameters code) (cadr code))
(define (lambda-body code) (cddr code))


(define (make-lambda parameters body)
  (cons 'lambda (cons parameters body)))


(define (compile-lambda-body code proc-label)
  (let ((names (lambda-parameters code)))
    (append-instruction-sequences
      (make-instruction-sequence
        '(env proc args) '(env)
        (list proc-label
              (list 'init-proc-env names)))
      (compile-sequence (lambda-body code) 'val 'return))))


(define (compile-lambda code target linkage)
  (let ((proc-label (make-label))
        (after-lambda (make-label)))
    (let ((lambda-linkage
            (if (eq? linkage 'next) after-lambda linkage)))
      (append-instruction-sequences
        (tack-on-instruction-sequence
          (end-with-linkage
            lambda-linkage
            (make-instruction-sequence
              '(env)
              (list target)
              (list (list 'make-compiled-procedure
                          target
                          proc-label))))
          (compile-lambda-body code proc-label))
        after-lambda))))


(define (cond? code) (tagged-list? code 'cond))


(define (cond-predicate clause) (car clause))


(define (cond-else-clause? clause)
  (eq? (cond-predicate clause) 'else))


(define (make-begin seq) (cons 'begin seq))


(define (sequence->exp seq)
  (cond ((null? seq) '())
        ((last-exp? seq) (first-exp seq))
        (else (make-begin seq))))


(define (cond-actions clause) (cdr clause))


(define (make-if predicate consequent alternative)
  (list 'if predicate consequent alternative))


(define (expand-clauses clauses)
  (if (null? clauses)
    'false
    (let ((first (car clauses))
          (rest (cdr clauses)))
      (if (cond-else-clause? first)
        (if (null? rest)
          (sequence->exp (cond-actions first))
          (panic "invalid cond clause" clauses))
        (make-if (cond-predicate first)
                 (sequence->exp (cond-actions first))
                 (expand-clauses rest))))))


(define (cond-clauses code) (cdr code))


(define (cond->if code)
  (expand-clauses (cond-clauses code)))


(define (application? code) (pair? code))
(define (operator code) (car code))
(define (operands code) (cdr code))


(define (code-to-get-rest-args operand-codes)
  (let ((code-for-next-arg
          (preserving
            '(args)
            (car operand-codes)
            (make-instruction-sequence
              '(val args) '(args)
              '((addarg))))))
    (if (null? (cdr operand-codes))
      code-for-next-arg
      (preserving
        '(env)
        code-for-next-arg
        (code-to-get-rest-args (cdr operand-codes))))))


(define (construct-arglist code)
  (let ((operand-codes (reverse code)))
    (if (null? operand-codes)
      (make-instruction-sequence
        '() '(args)
        '((initargs)))
      (let ((code-to-get-last-arg
              (append-instruction-sequences
                (car operand-codes)
                (make-instruction-sequence
                  '(val) '(args)
                  '((initargs) (addarg))))))
        (if (null? (cdr operand-codes))
          code-to-get-last-arg
          (preserving
            '(env)
            code-to-get-last-arg
            (code-to-get-rest-args
              (cdr operand-codes))))))))


(define (parallel-instruction-sequences instructions1 instructions2)
  (make-instruction-sequence
    (list-union (registers-needed instructions1)
                (registers-needed instructions2))
    (list-union (registers-modified instructions1)
                (registers-modified instructions2))
    (append (statements instructions1) (statements instructions2))))


(define all-regs
  '(val continue proc args env))


(define (compile-proc-appl target linkage)
  (cond ((and (eq? target 'val) (not (eq? linkage 'return)))
         (make-instruction-sequence
           '(proc) all-regs
           (list (list 'initreg 'continue (list 'label linkage))
                 '(takeproclabel)
                 '(goto (reg val)))))
        ((and (not (eq? target 'val))
              (not (eq? linkage 'return)))
         (let ((proc-return (make-label)))
           (make-instruction-sequence
             '(proc) all-regs
             (list (list 'initreg 'continue (list 'label proc-return))
                   '(takeproclabel)
                   '(goto (reg val))
                   proc-return
                   (list 'initreg target '(reg val))
                   (list 'goto (list 'label linkage))))))
        ((and (eq? target 'val) (eq? linkage 'return))
         (make-instruction-sequence
           '(proc continue) all-regs
           '((takeproclabel)
             (goto (reg val)))))
        (else
          (panic "invalid procedure application"))))


(define (compile-procedure-call target linkage)
  (let ((primitive-branch (make-label))
        (compiled-branch (make-label))
        (after-call (make-label)))
    (let ((compiled-linkage
            (if (eq? linkage 'next) after-call linkage)))
      (append-instruction-sequences
        (make-instruction-sequence
          '(proc) '()
          (list (list 'branchproc primitive-branch)))
        (parallel-instruction-sequences
          (append-instruction-sequences
            compiled-branch
            (compile-proc-appl target compiled-linkage))
          (append-instruction-sequences
            primitive-branch
            (end-with-linkage
              linkage
              (make-instruction-sequence
                '(proc args)
                (list target)
                (list (list 'apply-primitive-procedure
                            target))))))
        after-call))))


(define (compile-application code target linkage)
  (let ((proc-code (compile (operator code) 'proc 'next))
        (operand-codes
          (map (lambda (operand) (compile operand 'val 'next))
               (operands code))))
    (preserving
      '(env continue)
      proc-code
      (preserving
        '(proc continue)
        (construct-arglist operand-codes)
        (compile-procedure-call target linkage)))))


(define (or? code) (tagged-list? code 'or))


(define (or->if code)
  (define (or->if args)
    (cond ((null? args) 'false)
          ((null? (cdr args)) (car args))
          (else
            (list 'if
                  (car args)
                  (car args)
                  (or->if (cdr args))))))
  (or->if (cdr code)))


(define (and? code) (tagged-list? code 'and))


(define (and->if code)
  (define (and->if args)
    (cond ((null? args) 'true)
          ((null? (cdr args)) (car args))
          (else
            (list 'if
                  (car args)
                  (and->if (cdr args))
                  'false))))
  (and->if (cdr code)))


(define (let? code) (tagged-list? code 'let))


(define (let->lambda code)
  (append
    (list
      (append (list 'lambda)
              (list (map car (cadr code)))
              (cddr code)))
    (map cadr (cadr code))))


(define (compile code target linkage)
  (cond ((self-evaluating? code)
         (compile-self-evaluating code target linkage))
        ((quoted? code) (compile-quoted code target linkage))
        ((variable? code)
         (compile-variable code target linkage))
        ((assignment? code)
         (compile-assignment code target linkage))
        ((definition? code)
         (compile-definition code target linkage))
        ((if? code) (compile-if code target linkage))
        ((lambda? code) (compile-lambda code target linkage))
        ((begin? code)
         (compile-sequence (begin-actions code) target linkage))
        ((cond? code) (compile (cond->if code) target linkage))
        ((or? code) (compile (or->if code) target linkage))
        ((and? code) (compile (and->if code) target linkage))
        ((let? code) (compile (let->lambda code) target linkage))
        ((application? code)
         (compile-application code target linkage))
        (else
          (panic "invalid expression"))))


(define max-small 2000000)
(define min-small -2000000)


(define (register-name register)
  (define (req? rname) (eq? register rname))
  (cond ((req? 'continue) "cont")
        ((req? 'proc) "proc")
        ((req? 'args) "args")
        ((req? 'val) "val")
        ((req? 'env) "env")
        (else (panic "invalid register" register))))


(define (proclabel? code)
  (eq? code 'proclabel))


(define (substr string f len)
  (substring string f (+ f len)))


(define (assembly-escape string escaped)
  (define (char-eq? char1 char2) (eq? (string-ref char1 0) (string-ref char2 0)))
  (define (escape-char? char) (char-eq? char "\\"))
  (define (quotation-char? char) (char-eq? char "\""))
  (define (newline-char? char) (char-eq? char "\n"))
  (if (eq? (string-length string) 0)
    escaped
    (let ((char (substr string 0 1))
          (rest (substr string 1 (- (string-length string) 1))))
      (assembly-escape
        rest
        (string-append
          escaped
          (cond ((or (escape-char? char)
                     (quotation-char? char))
                 "\\")
                ((newline-char? char) "\\n")
                (else ""))
          (if (newline-char? char) "" char))))))


(define (assemble-chars string)
  (string-append "\""
                 (assembly-escape string "")
                 "\""))


(define (assemble-symbol-name symbol)
  (assemble-chars (symbol->string symbol)))


(define (assemble-symbol-value symbol)
  (string-append "symval("
                 (assemble-symbol-name symbol)
                 ")"))


(define (small-number? number)
  (and (integer? number)
       (>= number min-small)
       (<= number max-small)))


(define (assemble-small-number number)
  (string-append "numvali("
                 (number->string number)
                 ")"))


(define (assemble-large-number number)
  (let ((nums (number->string number)))
    (string-append "numvalc("
                   (assemble-chars nums)
                   ")")))


(define (assemble-number number)
  (cond ((small-number? number) (assemble-small-number number))
        (else (assemble-large-number number))))


(define (assemble-string string)
  (string-append "stringval("
                 (assemble-chars string)
                 ")"))


(define (assemble-self-evaluating const)
  (cond ((number? const) (assemble-number const))
        ((string? const) (assemble-string const))
        (else (panic "invalid const expression" const))))


(define (assemble-pair pair)
  (string-append "pairval("
                 (assemble-const (car pair))
                 ", "
                 (assemble-const (cdr pair))
                 ")"))


(define (assemble-null) "null")


(define (const? code)
  (and (tagged-list? code 'const)))


(define (assemble-const const)
  (cond ((symbol? const) (assemble-symbol-value const))
        ((self-evaluating? const) (assemble-self-evaluating const))
        ((pair? const) (assemble-pair const))
        ((null? const) (assemble-null))
        (else (panic "invalid const expression" const))))


(define (label? code)
  (and (tagged-list? code 'label)
       (integer? (cadr code))))


(define (assemble-label label)
  (number->string label))


(define (assemble-label-value label)
  (string-append "numvali("
                 (assemble-label label)
                 ")"))


(define (register? code)
  (and (tagged-list? code 'reg)
       (symbol? (cadr code))))


(define (assemble-register reg)
  (string-append "rm->"
                 (register-name reg)))


(define (assemble-expression expression)
  (cond ((const? expression) (assemble-const (cadr expression)))
        ((label? expression) (assemble-label-value (cadr expression)))
        ((register? expression) (assemble-register (cadr expression)))
        (else (panic "invalid expression" expression))))


(define (stack-operation? instruction operation)
  (and (tagged-list? instruction operation)
       (symbol? (cadr instruction))))


(define (save-env? instruction)
  (tagged-list? instruction 'saveenv))


(define (assemble-save-env instruction)
  (string-append "saveenv(rm);\n"))


(define (restore-env? instruction)
  (tagged-list? instruction 'restoreenv))


(define (assemble-restore-env instruction)
  (string-append "restoreenv(rm);\n"))


(define (save? instruction)
  (stack-operation? instruction 'save))


(define (assemble-save instruction)
  (string-append "savereg(rm, "
                 (assemble-register (cadr instruction))
                 ");\n"))


(define (restore? instruction)
  (stack-operation? instruction 'restore))


(define (assemble-restore instruction)
  (string-append "restorereg(rm, (void *)&"
                 (assemble-register (cadr instruction))
                 ");\n"))


(define (goto? instruction)
  (and (tagged-list? instruction 'goto)
       (or (label? (cadr instruction))
           (register? (cadr instruction))
           (proclabel? (cadr instruction)))))


(define (assemble-goto instruction)
  (cond ((label? (cadr instruction))
         (string-append "gotolabel(rm, "
                        (assemble-label (cadadr instruction))
                        "); break;\n"))
        ((register? (cadr instruction))
         (string-append "gotoreg(rm, "
                        (assemble-register (cadadr instruction))
                        "); break;\n"))
        ((proclabel? (cadr instruction)) "gotoproc(rm); break;\n")))


(define (takeproclabel? instruction)
  (and (pair? instruction)
       (eq? (car instruction) 'takeproclabel)))


(define (assemble-takeproclabel instruction)
  "takeproclabel(rm);\n")


(define (initreg? instruction)
  (and (tagged-list? instruction 'initreg)
       (symbol? (cadr instruction))
       (or (const? (caddr instruction))
           (label? (caddr instruction))
           (register? (caddr instruction)))))


(define (assemble-initreg instruction)
  (string-append "initreg((void *)&"
                 (assemble-register (cadr instruction))
                 ", "
                 (assemble-expression (caddr instruction))
                 ");\n"))


(define (get-variable? instruction)
  (and (tagged-list? instruction 'get-variable)
       (symbol? (cadr instruction))
       (symbol? (caddr instruction))))


(define (assemble-get-variable instruction)
  (string-append "getenvvar(rm, (void *)&"
                 (assemble-register (cadr instruction))
                 ", "
                 (assemble-symbol-name (caddr instruction))
                 ");\n"))


(define (set-variable? instruction)
  (and (tagged-list? instruction 'set-variable-value)
       (symbol? (cadr instruction))
       (symbol? (caddr instruction))))


(define (assemble-set-variable instruction)
  (string-append "setenvvar(rm, (void *)&"
                 (assemble-register (cadr instruction))
                 ", "
                 (assemble-symbol-name (caddr instruction))
                 ");\n"))


(define (define-variable? instruction)
  (and (tagged-list? instruction 'define-variable)
       (symbol? (cadr instruction))))


(define (assemble-define-variable instruction)
  (string-append "defenvvar(rm, "
                 (assemble-symbol-name (cadr instruction))
                 ");\n"))


(define (branchval? instruction)
  (and (tagged-list? instruction 'branchval)
       (integer? (cadr instruction))))


(define (assemble-branchval instruction)
  (string-append "if (branchval(rm, "
                 (assemble-label (cadr instruction))
                 ")) { break; }\n"))


(define (initprocenv? instruction)
  (and (tagged-list? instruction 'init-proc-env)
       (or (null? (cadr instruction))
           (pair? (cadr instruction))
           (symbol? (cadr instruction)))))


(define (assemble-initprocenv instruction)
  (string-append "initprocenv(rm, "
                 (assemble-const (cadr instruction))
                 ");\n"))


(define (makeproc? instruction)
  (and (tagged-list? instruction 'make-compiled-procedure)
       (symbol? (cadr instruction))
       (integer? (caddr instruction))))


(define (assemble-makeproc instruction)
  (string-append "mkcompiledprocreg(rm, (void *)&"
                 (assemble-register (cadr instruction))
                 ", "
                 (assemble-label (caddr instruction))
                 ");\n"))


(define (initargs? instruction)
  (and (pair? instruction)
       (eq? (car instruction) 'initargs)))


(define (assemble-initargs)
  "initargs(rm);\n")


(define (addarg? instruction)
  (and (pair? instruction)
       (eq? (car instruction) 'addarg)))


(define (assemble-addarg)
  "addarg(rm);\n")


(define (branchproc? instruction)
  (and (tagged-list? instruction 'branchproc)
       (integer? (cadr instruction))))


(define (assemble-branchproc instruction)
  (string-append "if (branchproc(rm, "
                 (assemble-label (cadr instruction))
                 ")) { break; }\n"))


(define (applyprimitive? instruction)
  (and (tagged-list? instruction 'apply-primitive-procedure)
       (symbol? (cadr instruction))))


(define (assemble-applyprimitive instruction)
  (string-append "applyprimitivereg(rm, (void *)&"
                 (assemble-register (cadr instruction))
                 ");\n"))


(define (label-def? instruction)
  (integer? instruction))


(define (assemble-label-def instruction)
  (string-append "case "
                 (assemble-label instruction)
                 ":\n"))


(define (assemble-instruction instruction)
  (cond ((save-env? instruction) (assemble-save-env instruction))
        ((restore-env? instruction) (assemble-restore-env instruction))
        ((save? instruction) (assemble-save instruction))
        ((restore? instruction) (assemble-restore instruction))
        ((goto? instruction) (assemble-goto instruction))
        ((takeproclabel? instruction) (assemble-takeproclabel instruction))
        ((initreg? instruction) (assemble-initreg instruction))
        ((get-variable? instruction) (assemble-get-variable instruction))
        ((set-variable? instruction) (assemble-set-variable instruction))
        ((define-variable? instruction) (assemble-define-variable instruction))
        ((branchval? instruction) (assemble-branchval instruction))
        ((initprocenv? instruction) (assemble-initprocenv instruction))
        ((makeproc? instruction) (assemble-makeproc instruction))
        ((initargs? instruction) (assemble-initargs))
        ((addarg? instruction) (assemble-addarg))
        ((branchproc? instruction) (assemble-branchproc instruction))
        ((applyprimitive? instruction) (assemble-applyprimitive instruction))
        ((label-def? instruction) (assemble-label-def instruction))
        (else (panic "invalid assembly instruction" instruction))))


(define (assemble assembly output)
  (if (null? assembly) output
    (assemble (cdr assembly)
              (string-append output
                             (assemble-instruction (car assembly))))))


(define output-head "
#include <stdlib.h>
#include <stdio.h>
#include \"../src-head/sys.h\"
#include \"../src-head/error.h\"
#include \"../src-head/sysio.h\"
#include \"../src-head/number.h\"
#include \"../src-head/string.h\"
#include \"../src-head/compound-types.h\"
#include \"../src-head/stack.h\"
#include \"../src-head/io.h\"
#include \"../src-head/value.h\"
#include \"../src-head/register-machine.h\"
#include \"../src-head/primitives.h\"
#include \"../src-head/registry.h\"

int main() {
    initsys();
    initmodule_sysio();
    initmodule_io();
    initmodule_number();
    initmodule_value();
    initmodule_primitives();
    initmodule_registry();

    regmachine rm = mkregmachine();
    for (;;) {
        long labelval = valrawint(rm->label);
        switch (labelval) {
        case 0:
        ")


(define output-tail "
        default:
            // printf(\"%s\\n\", sprintraw(rm->val));
            return 0;
        }
    }

    freeregmachine(rm);

    freemodule_registry();
    freemodule_primitives();
    freemodule_value();
    freemodule_number();
    freemodule_io();
    freemodule_sysio();

    return 0;
}
")


(define (write-head)
  (format (current-output-port) output-head))


(define (write-tail)
  (format (current-output-port) output-tail))


(define (write-code output-code)
  (format (current-output-port) output-code))


(define code
  '((lambda ()
      (define (debug . args)
        (cond ((null? args) 'ok)
              ((null? (cdr args))
               (write-file stderr (car args))
               (write-file stderr "\n")
               (car args))
              (else
                (write-file stderr (car args))
                (write-file stderr " ")
                (apply debug (cdr args)))))


      (debug "running")


      (define (read-file-utf8 f k)
        (define read-length (if (< k 6) 6 k))
        (define (fail) (error "file not utf8"))
        (define (utf8-or-fail bytes)
          (if (utf8-string? bytes) bytes (fail)))
        (define (cut-and-seek current)
          (define string (string-copy current 0 k))
          (seek-file f
                     (- (bytes-length string) (bytes-length current))
                     file-seek-mode-current)
          string)
        (define (read previous)
          (define bytes (read-file f read-length))
          (if (eof? bytes)
            (if (eq? (bytes-length previous) 0)
              eof
              (utf8-or-fail previous))
            (begin
              (define current (string-append previous bytes))
              (cond ((not (> (string-length current)
                             (string-length previous)))
                     (fail))
                    ((>= (string-length current) k)
                     (cut-and-seek current))
                    (else (read current))))))
        (read ""))


      (define tokenizer-expression
        '(";[^\\n]*\\n?|"                    ; comment
          "\\(|\\)|"                         ; list open, list/vector/struct close
          ; "#\\(|"                          ; vector open
          ; "#s\\(|"                         ; struct open
          "'|"                               ; quote
          "\"(\\\\\\\\|\\\\\"|[^\"])*\"?|"   ; string
          "(\\\\.|"                          ; symbol, single escape
          "\\|(\\\\\\\\|\\\\\\||[^|])*\\|?|" ; symbol, range escape
          "[^;()#'|\"\\s])+"))               ; symbol, no comment/list/type-escape/quote/string/whitespace


      ; this is kind of unnecessary to be called separately
      (define token-complete-expression
        '("^;[^\\n]*\\n|"                         ; comment
          "\\(|"                                  ; list open
          "\\)|"                                  ; list/vector/struct close
          "#\\(|"                                 ; vector open
          "#s\\(|"                                ; struct open
          "'|"                                    ; quote
          "\"(\\\\\"|\\\\[^\"]|[^\\\\\"])*\"|"    ; string
          "(\\\\.|"                               ; symbol, single escape
          "\\|(\\\\\\||\\\\[^\\|]|[^\\\\|])*\\||" ; symbol, range escape
          "[^;()#'|\"\\s\\\\])+$"))               ; symbol, no comment/list/type-escape/quote/string/whitespace/escape


      ; stricter syntax check is needed
      (define remove-whitespace
        ((lambda ()
           (define rx (make-regex "^\\s*" 0))
           (lambda (string)
             (define m (regex-match rx string))
             (if (null? m)
               string
               (begin
                 (define wl (car (cdr (car m))))
                 (string-copy string wl (- (string-length string) wl))))))))


      (define get-token
        ((lambda ()
           (define rx (make-regex
                        (apply string-append tokenizer-expression)
                        0))
           (lambda (string)
             (define m (regex-match rx string))
             (if (null? m)
               ""
               (string-copy string
                            (car (car m))
                            (car (cdr (car m)))))))))


      (define token-complete?
        ((lambda ()
           (define rx (make-regex
                        (apply string-append token-complete-expression)
                        0))
           (lambda (string) (not (null? (regex-match rx string)))))))


      (define f (open-file "some-file" file-mode-read))
      (define s (read-file-utf8 f 48))
      (close-file f)
      (define s (remove-whitespace "fűzfánfütyülő (some (list (to the file)))"))


      (define (make-token-port f)
        (define string "")
        (define (fail) (error "read error"))
        (define (return-token token)
          (set! string
            (remove-whitespace
              (string-copy
                string
                (string-length token)
                (- (string-length string)
                   (string-length token)))))
          token)
        (define (read-file)
          (define string (read-file-utf8 f 8192))
          (if (or (eof? string) (> (string-length string) 0))
            string
            (fail)))
        (define (read-file-continue on-eof)
          (define s (read-file))
          (if (eof? s)
            (on-eof)
            (begin
              (set! string (remove-whitespace (string-append string s)))
              (read))))
        (define (read-file-continue-or-eof)
          (read-file-continue (lambda () eof)))
        (define (read-file-continue-or-fail)
          (read-file-continue fail))
        (define (read-file-continue-or-token token)
          (read-file-continue (lambda () (return-token token))))
        (define (read)
          (if (== (string-length string) 0)
            (read-file-continue-or-eof)
            (begin
              (define token (get-token string))
              (define full? (== (string-length token)
                                (string-length string)))
              ; this could be a bit more optimized not to use a separate rx
              (define complete? (token-complete? token))
              (cond ((== (string-length token) 0)
                     (read-file-continue-or-fail))
                    ((and complete? full?)
                     (read-file-continue-or-token token))
                    ((and complete? (not full?))
                     (return-token token))
                    ((and (not complete?) full?)
                     (read-file-continue-or-fail))
                    (else (fail))))))
        (lambda () (read)))


      (define (read-token p) (p))


      (define f (open-file "some-file" file-mode-read))
      (define p (make-token-port f))
      (define (read-all-tokens p)
        (define token (read-token p))
        (if (eof? token)
          'done
          (read-all-tokens p)))
      ; (read-all-tokens p)
      (close-file f)


      (define (unescape string)
        (define (unescape-char char)
          (cond ; ((== char "b") "\b")
            ; ((== char "t") "\t")
            ((== char "n") "\n")
            ; ((== char "v") "\v")
            ; ((== char "f") "\f")
            ; ((== char "r") "\r")
            (else char)))
        (define (unescape string)
          (cond ((== (string-length string) 0) "")
                ((and (== (string-copy string 0 1) "\\")
                      (== (string-length string) 1))
                 (error "invalid escape sequence"))
                ((== (string-copy string 0 1) "\\")
                 (string-append (unescape-char (string-copy string 1 1))
                                (unescape
                                  (string-copy string
                                               2
                                               (- (string-length string) 2)))))
                (else
                  (string-append (string-copy string 0 1)
                                 (unescape
                                   (string-copy string
                                                1
                                                (- (string-length string) 1)))))))
        (unescape string))


      (define (unescape-symbol string)
        (define (unescape escaped? string)
          (cond ((and escaped? (== (string-length string) 0))
                 (error "invalid escape sequence"))
                ((== (string-length string) 0) "")
                ((== (string-copy string 0 1) "|")
                 (unescape (not escaped?)
                           (string-copy string
                                        1
                                        (- (string-length string) 1))))
                ((and (== (string-copy string 0 1) "\\")
                      (== (string-length string) 1))
                 (error "invalid escape sequence"))
                ((and escaped?
                      (== (string-copy string 0 1) "\\")
                      (or (== (string-copy string 1 1) "\\")
                          (== (string-copy string 1 1) "|")))
                 (string-append (string-copy string 1 1)
                                (unescape
                                  true
                                  (string-copy string
                                               2
                                               (- (string-length string) 2)))))
                ((and (not escaped?)
                      (== (string-copy string 0 1) "\\"))
                 (string-append (string-copy string 1 1)
                                (unescape
                                  false
                                  (string-copy string
                                               2
                                               (- (string-length string) 2)))))
                (else
                  (string-append (string-copy string 0 1)
                                 (unescape
                                   escaped?
                                   (string-copy string
                                                1
                                                (- (string-length string) 1)))))))
        (unescape false string))


      (define (list . x) x)


      (define (caar l) (car (car l)))
      (define (cadr l) (car (cdr l)))
      (define (cdar l) (cdr (car l)))
      (define (cddr l) (cdr (cdr l)))

      (define (caaar l) (car (car (car l))))
      (define (caadr l) (car (car (cdr l))))
      (define (cadar l) (car (cdr (car l))))
      (define (caddr l) (car (cdr (cdr l))))
      (define (cdaar l) (cdr (car (car l))))
      (define (cdadr l) (cdr (car (cdr l))))
      (define (cddar l) (cdr (cdr (car l))))
      (define (cdddr l) (cdr (cdr (cdr l))))

      (define (caaaar l) (car (car (car (car l)))))
      (define (caaadr l) (car (car (car (cdr l)))))
      (define (caadar l) (car (car (cdr (car l)))))
      (define (caaddr l) (car (car (cdr (cdr l)))))
      (define (cadaar l) (car (cdr (car (car l)))))
      (define (cadadr l) (car (cdr (car (cdr l)))))
      (define (caddar l) (car (cdr (cdr (car l)))))
      (define (cadddr l) (car (cdr (cdr (cdr l)))))
      (define (cdaaar l) (cdr (car (car (car l)))))
      (define (cdaadr l) (cdr (car (car (cdr l)))))
      (define (cdadar l) (cdr (car (cdr (car l)))))
      (define (cdaddr l) (cdr (car (cdr (cdr l)))))
      (define (cddaar l) (cdr (cdr (car (car l)))))
      (define (cddadr l) (cdr (cdr (car (cdr l)))))
      (define (cdddar l) (cdr (cdr (cdr (car l)))))
      (define (cddddr l) (cdr (cdr (cdr (cdr l)))))


      (define (map f l)
        (cond ((null? l) '())
              (else (cons (f (car l)) (map f (cdr l))))))


      (define (member m l c)
        (cond ((null? l) false)
              ((c m (car l)) l)
              (else (member m (cdr l) c))))


      (define (memq m l) (member m l ==))


      (define (append . l)
        (cond ((null? l) '())
              ((null? (car l)) (apply append (cdr l)))
              (else (cons (caar l) (apply append (cons (cdar l) (cdr l)))))))


      (define (reverse l)
        (cond ((null? l) '())
              (else (append (reverse (cdr l)) (list (car l))))))


      (define (read p)
        (define (comment? token) (== (string-copy token 0 1) ";"))
        (define (quote-mark? token) (== token "'"))
        (define (list-open? token) (== token "("))
        (define (list-close? token) (== token ")"))
        (define (pair-separator? token) (== token "."))
        (define (token->number token) (string->number-safe token))
        (define (token->string token)
          (and (== (string-copy token 0 1) "\"")
               (unescape
                 (string-copy token 1 (- (string-length token) 2)))))
        (define (token->symbol token) (string->symbol (unescape-symbol token)))
        (define (read-datum token)
          (or (token->number token)
              (token->string token)
              (token->symbol token)))
        (define (reverse-over l over)
          (cond ((null? l) over)
                (else (reverse-over (cdr l)
                                    (cons (car l) over)))))
        (define (finish-list l)
          (reverse-over l '()))
        (define (finish-pair l)
          (define last (read l))
          (if (not (== (read l) l))
            (error "invalid list")
            (reverse-over l last)))
        (define (read-list l)
          (define object (read l))
          (cond ((eof? object) (error "unclosed list"))
                ((== object l) (finish-list l))
                (object (read-list (cons object l)))
                (else (finish-pair l))))
        (define (read l)
          (define token (read-token p))
          (cond ((eof? token) eof)
                ((comment? token) (read l))
                ((quote-mark? token) (list 'quote (read false)))
                ((list-open? token) (read-list '()))
                ((and l (list-close? token)) l)
                ((and l (pair-separator? token)) false)
                ((list-close? token) (error "unexpected list close"))
                ((pair-separator? token) (error "unexpected cons"))
                (else (read-datum token))))
        (read false))


      (define (mkreftable)
        (let ((symbols '())
              (current-ref 0))
          (define (lookup symbol-list symbol)
            (cond ((null? symbol-list) false)
                  ((eq? (caar symbol-list) symbol)
                   (cadar symbol-list))
                  (else
                    (lookup (cdr symbol-list)
                            symbol))))
          (define (create-ref symbol)
            (cond ((lookup symbols symbol)
                   (error "symbol ref exists")))
            (let ((ref current-ref))
              (set! symbols
                (cons (list symbol current-ref)
                      symbols))
              (set! current-ref (+ current-ref 1))
              ref))
          (define (symbol-ref symbol-list symbol)
            (let ((ref (lookup symbols symbol)))
              (if ref ref (error "symbol not found"))))
          (lambda (message symbol)
            (cond ((eq? message 'lookup)
                   (lookup symbols symbol))
                  ((eq? message 'ref)
                   (symbol-ref symbols symbol))
                  ((eq? message 'create-ref)
                   (create-ref symbol))))))


      (define (lookup-ref reftable symbol) (reftable 'lookup symbol))
      (define (symbol-ref reftable symbol) (reftable 'ref symbol))
      (define (create-ref reftable symbol) (reftable 'create-ref symbol))


      (define variable-ref
        (let ((variables (mkreftable)))
          (lambda (name)
            (let ((ref (lookup-ref variables name)))
              (if ref ref (create-ref variables name))))))


      (define make-label
        (let ((label-count 0))
          (lambda ()
            (set! label-count (+ label-count 1))
            label-count)))


      (define (list-union list1 list2)
        (cond ((null? list1) list2)
              ((memq (car list1) list2) (list-union (cdr list1) list2))
              (else (cons (car list1) (list-union (cdr list1) list2)))))


      (define (list-difference list1 list2)
        (cond ((null? list1) '())
              ((memq (car list1) list2) (list-difference (cdr list1) list2))
              (else (cons (car list1) (list-difference (cdr list1) list2)))))


      (define (make-instruction-sequence needs modifies statements)
        (list needs modifies statements))


      (define (registers-needed instructions)
        (cond ((number? instructions) '())
              ((pair? instructions) (car instructions))
              (else
                (display instructions)
                (newline)
                (error "something went wrong"))))


      (define (registers-modified instructions)
        (cond ((number? instructions) '())
              ((pair? instructions) (cadr instructions))
              (else
                (display instructions)
                (newline)
                (error "something went wrong"))))


      (define (statements instructions)
        (cond ((number? instructions) (list instructions))
              ((pair? instructions) (caddr instructions))
              (else
                (display instructions)
                (newline)
                (error "something went wrong"))))


      (define (empty-instruction-sequence)
        (make-instruction-sequence '() '() '()))


      (define (append-instruction-sequences . seqs)
        (define (append-2-sequences instructions1 instructions2)
          (make-instruction-sequence
            (list-union (registers-needed instructions1)
                        (list-difference (registers-needed instructions2)
                                         (registers-modified instructions1)))
            (list-union (registers-modified instructions1)
                        (registers-modified instructions2))
            (append (statements instructions1) (statements instructions2))))
        (define (append-seq-list seqs)
          (if (null? seqs)
            (empty-instruction-sequence)
            (append-2-sequences (car seqs)
                                (append-seq-list (cdr seqs)))))
        (append-seq-list seqs))


      (define (needs-register? instructions register)
        (memq register (registers-needed instructions)))


      (define (modifies-register? instructions register)
        (memq register (registers-modified instructions)))


      (define (preserving regs instructions1 instructions2)
        (if (null? regs)
          (append-instruction-sequences instructions1 instructions2)
          (let ((first-reg (car regs)))
            (if (and (needs-register? instructions2 first-reg)
                     (modifies-register? instructions1 first-reg))
              (preserving
                (cdr regs)
                (make-instruction-sequence
                  (list-union (list first-reg)
                              (registers-needed instructions1))
                  (list-difference (registers-modified instructions1)
                                   (list first-reg))
                  (append (list (if (eq? first-reg 'env)
                                  (list 'saveenv)
                                  (list 'save first-reg)))
                          (statements instructions1)
                          (list (if (eq? first-reg 'env)
                                  (list 'restoreenv)
                                  (list 'restore first-reg)))))
                instructions2)
              (preserving (cdr regs) instructions1 instructions2)))))


      (define (compile-linkage linkage)
        (cond ((eq? linkage 'return)
               (make-instruction-sequence
                 '(continue) '()
                 '((goto (reg continue)))))
              ((eq? linkage 'next)
               (empty-instruction-sequence))
              (else
                (make-instruction-sequence
                  '() '()
                  (list (list 'goto (list 'label linkage)))))))


      (define (end-with-linkage linkage instructions)
        (preserving
          '(continue)
          instructions
          (compile-linkage linkage)))


      (define (tack-on-instruction-sequence instructions body-instructions)
        (make-instruction-sequence
          (registers-needed instructions)
          (registers-modified instructions)
          (append (statements instructions) (statements body-instructions))))


      (define (self-evaluating? code)
        (or (number? code)
            (string? code)
            (== code true)
            (== code false)))


      (define (compile-self-evaluating code target linkage)
        (end-with-linkage
          linkage
          (make-instruction-sequence
            '() (list target)
            (list (list 'initreg target (list 'const code))))))


      (define (quoted? code) (tagged-list? code 'quote))


      (define (text-of-quotation code) (cadr code))


      (define (compile-quoted code target linkage)
        (end-with-linkage
          linkage
          (make-instruction-sequence
            '()
            (list target)
            (list (list 'initreg
                        target
                        (list 'const
                              (text-of-quotation code)))))))


      (define (variable? code) (symbol? code))


      (define (compile-variable code target linkage)
        (end-with-linkage
          linkage
          (make-instruction-sequence
            '(env) (list target)
            (list (list 'get-variable target code)))))


      (define (assignment? code) (tagged-list? code 'set!))


      (define (assignment-variable code) (cadr code))


      (define (assignment-value code) (caddr code))


      (define (compile-assignment code target linkage)
        (let ((var (assignment-variable code))
              (get-value-code
                (compile (assignment-value code) 'val 'next)))
          (end-with-linkage
            linkage
            (preserving
              '(env)
              get-value-code
              (make-instruction-sequence
                '(env val) (list target)
                (list (list 'set-variable-value target var)))))))


      (define (definition? code)
        (tagged-list? code 'define))


      (define (definition-variable code)
        (if (symbol? (cadr code))
          (cadr code)
          (caadr code)))


      (define (definition-value code)
        (if (symbol? (cadr code))
          (caddr code)
          (make-lambda (cdadr code) (cddr code))))


      (define (compile-definition code target linkage)
        (let ((var (definition-variable code))
              (value-code
                (compile (definition-value code) 'val 'next)))
          (end-with-linkage
            linkage
            (preserving
              '(env)
              value-code
              (make-instruction-sequence
                '(env val) (list target)
                (list (list 'define-variable var)))))))


      (define (begin? code)
        (tagged-list? code 'begin))


      (define (begin-actions code)
        (cdr code))


      (define (last-exp? code) (null? (cdr code)))
      (define (first-exp code) (car code))
      (define (rest-exps code) (cdr code))


      (define (compile-sequence code target linkage)
        (if (last-exp? code)
          (compile (first-exp code) target linkage)
          (preserving
            '(env continue)
            (compile (first-exp code) target 'next)
            (compile-sequence (rest-exps code) target linkage))))


      (define (tagged-list? code tag)
        (and (pair? code)
             (eq? (car code) tag)))


      (define (if? code) (tagged-list? code 'if))


      (define (if-predicate code) (cadr code))


      (define (if-consequent code) (caddr code))


      (define (if-alternative code)
        (if (null? (cdddr code))
          'false
          (cadddr code)))


      (define (compile-if code target linkage)
        (let ((t-branch (make-label))
              (f-branch (make-label))
              (after-if (make-label)))
          (let ((consequent-linkage
                  (if (eq? linkage 'next) after-if linkage)))
            (let ((p-code (compile (if-predicate code) 'val 'next))
                  (c-code (compile (if-consequent code) target consequent-linkage))
                  (a-code (compile (if-alternative code) target linkage)))
              (preserving
                '(env continue)
                p-code
                (append-instruction-sequences
                  (make-instruction-sequence
                    '(val) '()
                    (list (list 'branchval f-branch)))
                  (parallel-instruction-sequences
                    (append-instruction-sequences t-branch c-code)
                    (append-instruction-sequences f-branch a-code))
                  after-if))))))


      (define (lambda? code) (tagged-list? code 'lambda))


      (define (lambda-parameters code) (cadr code))
      (define (lambda-body code) (cddr code))


      (define (make-lambda parameters body)
        (cons 'lambda (cons parameters body)))


      (define (compile-lambda-body code proc-label)
        (let ((names (lambda-parameters code)))
          (append-instruction-sequences
            (make-instruction-sequence
              '(env proc args) '(env)
              (list proc-label
                    (list 'init-proc-env names)))
            (compile-sequence (lambda-body code) 'val 'return))))


      (define (compile-lambda code target linkage)
        (let ((proc-label (make-label))
              (after-lambda (make-label)))
          (let ((lambda-linkage
                  (if (eq? linkage 'next) after-lambda linkage)))
            (append-instruction-sequences
              (tack-on-instruction-sequence
                (end-with-linkage
                  lambda-linkage
                  (make-instruction-sequence
                    '(env)
                    (list target)
                    (list (list 'make-compiled-procedure
                                target
                                proc-label))))
                (compile-lambda-body code proc-label))
              after-lambda))))


      (define (cond? code) (tagged-list? code 'cond))


      (define (cond-predicate clause) (car clause))


      (define (cond-else-clause? clause)
        (eq? (cond-predicate clause) 'else))


      (define (make-begin seq) (cons 'begin seq))


      (define (sequence->exp seq)
        (cond ((null? seq) '())
              ((last-exp? seq) (first-exp seq))
              (else (make-begin seq))))


      (define (cond-actions clause) (cdr clause))


      (define (make-if predicate consequent alternative)
        (list 'if predicate consequent alternative))


      (define (expand-clauses clauses)
        (if (null? clauses)
          'false
          (let ((first (car clauses))
                (rest (cdr clauses)))
            (if (cond-else-clause? first)
              (if (null? rest)
                (sequence->exp (cond-actions first))
                (error "invalid cond clause"))
              (make-if (cond-predicate first)
                       (sequence->exp (cond-actions first))
                       (expand-clauses rest))))))


      (define (cond-clauses code) (cdr code))


      (define (cond->if code)
        (expand-clauses (cond-clauses code)))


      (define (application? code) (pair? code))
      (define (operator code) (car code))
      (define (operands code) (cdr code))


      (define (code-to-get-rest-args operand-codes)
        (let ((code-for-next-arg
                (preserving
                  '(args)
                  (car operand-codes)
                  (make-instruction-sequence
                    '(val args) '(args)
                    '((addarg))))))
          (if (null? (cdr operand-codes))
            code-for-next-arg
            (preserving
              '(env)
              code-for-next-arg
              (code-to-get-rest-args (cdr operand-codes))))))


      (define (construct-arglist code)
        (let ((operand-codes (reverse code)))
          (if (null? operand-codes)
            (make-instruction-sequence
              '() '(args)
              '((initargs)))
            (let ((code-to-get-last-arg
                    (append-instruction-sequences
                      (car operand-codes)
                      (make-instruction-sequence
                        '(val) '(args)
                        '((initargs) (addarg))))))
              (if (null? (cdr operand-codes))
                code-to-get-last-arg
                (preserving
                  '(env)
                  code-to-get-last-arg
                  (code-to-get-rest-args
                    (cdr operand-codes))))))))


      (define (parallel-instruction-sequences instructions1 instructions2)
        (make-instruction-sequence
          (list-union (registers-needed instructions1)
                      (registers-needed instructions2))
          (list-union (registers-modified instructions1)
                      (registers-modified instructions2))
          (append (statements instructions1) (statements instructions2))))


      (define all-regs
        '(val continue proc args env))


      (define (compile-proc-appl target linkage)
        (cond ((and (eq? target 'val) (not (eq? linkage 'return)))
               (make-instruction-sequence
                 '(proc) all-regs
                 (list (list 'initreg 'continue (list 'label linkage))
                       '(takeproclabel)
                       '(goto (reg val)))))
              ((and (not (eq? target 'val))
                    (not (eq? linkage 'return)))
               (let ((proc-return (make-label)))
                 (make-instruction-sequence
                   '(proc) all-regs
                   (list (list 'initreg 'continue (list 'label proc-return))
                         '(takeproclabel)
                         '(goto (reg val))
                         proc-return
                         (list 'initreg target '(reg val))
                         (list 'goto (list 'label linkage))))))
              ((and (eq? target 'val) (eq? linkage 'return))
               (make-instruction-sequence
                 '(proc continue) all-regs
                 '((takeproclabel)
                   (goto (reg val)))))
              (else
                (error "invalid procedure application"))))


      (define (compile-procedure-call target linkage)
        (let ((primitive-branch (make-label))
              (compiled-branch (make-label))
              (after-call (make-label)))
          (let ((compiled-linkage
                  (if (eq? linkage 'next) after-call linkage)))
            (append-instruction-sequences
              (make-instruction-sequence
                '(proc) '()
                (list (list 'branchproc primitive-branch)))
              (parallel-instruction-sequences
                (append-instruction-sequences
                  compiled-branch
                  (compile-proc-appl target compiled-linkage))
                (append-instruction-sequences
                  primitive-branch
                  (end-with-linkage
                    linkage
                    (make-instruction-sequence
                      '(proc args)
                      (list target)
                      (list (list 'apply-primitive-procedure
                                  target))))))
              after-call))))


      (define (compile-application code target linkage)
        (let ((proc-code (compile (operator code) 'proc 'next))
              (operand-codes
                (map (lambda (operand) (compile operand 'val 'next))
                     (operands code))))
          (preserving
            '(env continue)
            proc-code
            (preserving
              '(proc continue)
              (construct-arglist operand-codes)
              (compile-procedure-call target linkage)))))


      (define (or? code) (tagged-list? code 'or))


      (define (or->if code)
        (define (or->if args)
          (cond ((null? args) 'false)
                ((null? (cdr args)) (car args))
                (else
                  (list 'if
                        (car args)
                        (car args)
                        (or->if (cdr args))))))
        (or->if (cdr code)))


      (define (and? code) (tagged-list? code 'and))


      (define (and->if code)
        (define (and->if args)
          (cond ((null? args) 'true)
                ((null? (cdr args)) (car args))
                (else
                  (list 'if
                        (car args)
                        (and->if (cdr args))
                        'false))))
        (and->if (cdr code)))


      (define (let? code) (tagged-list? code 'let))


      (define (let->lambda code)
        (append
          (list
            (append (list 'lambda)
                    (list (map car (cadr code)))
                    (cddr code)))
          (map cadr (cadr code))))


      (define (compile code target linkage)
        (debug code)
        (cond ((self-evaluating? code)
               (compile-self-evaluating code target linkage))
              ((quoted? code)
               (compile-quoted code target linkage))
              ((variable? code)
               (compile-variable code target linkage))
              ((assignment? code)
               (compile-assignment code target linkage))
              ((definition? code)
               (compile-definition code target linkage))
              ((if? code)
               (compile-if code target linkage))
              ((lambda? code)
               (compile-lambda code target linkage))
              ((begin? code)
               (compile-sequence (begin-actions code) target linkage))
              ((cond? code)
               (compile (cond->if code) target linkage))
              ((or? code)
               (compile (or->if code) target linkage))
              ((and? code)
               (compile (and->if code) target linkage))
              ((let? code)
               (compile (let->lambda code) target linkage))
              ((application? code)
               (compile-application code target linkage))
              (else
                (error "compile: invalid expression"))))


      (define max-small 2000000)
      (define min-small -2000000)


      (define (register-name register)
        (define (req? rname) (eq? register rname))
        (cond ((req? 'continue) "cont")
              ((req? 'proc) "proc")
              ((req? 'args) "args")
              ((req? 'val) "val")
              ((req? 'env) "env")
              (else (error "invalid register"))))


      (define (proclabel? code)
        (eq? code 'proclabel))


      (define (assembly-escape string)
        (define (escape-char? char) (== char "\\"))
        (define (quotation-char? char) (== char "\""))
        (define (newline-char? char) (== char "\n"))
        (if (eq? (string-length string) 0)
          ""
          (let ((char (string-copy string 0 1))
                (rest (string-copy string 1 (- (string-length string) 1))))
            (string-append
              (cond ((or (escape-char? char)
                         (quotation-char? char))
                     "\\")
                    ((newline-char? char) "\\n")
                    (else ""))
              (if (newline-char? char) "" char)
              (assembly-escape rest)))))


      (define (assemble-chars string)
        (string-append "\""
                       (assembly-escape string)
                       "\""))


      (define (assemble-symbol-name symbol)
        (assemble-chars (symbol->string symbol)))


      (define (assemble-symbol-value symbol)
        (string-append "symval("
                       (assemble-symbol-name symbol)
                       ")"))


      (define (small-number? number)
        (and (integer? number)
             (>= number min-small)
             (<= number max-small)))


      (define (assemble-small-number number)
        (string-append "numvali("
                       (number->string number)
                       ")"))


      (define (assemble-large-number number)
        (let ((nums (number->string number)))
          (string-append "numvalc("
                         (assemble-chars nums)
                         ")")))


      (define (assemble-number number)
        (cond ((small-number? number) (assemble-small-number number))
              (else (assemble-large-number number))))


      (define (assemble-string string)
        (string-append "stringval("
                       (assemble-chars string)
                       ")"))


      (define (assemble-self-evaluating const)
        (cond ((number? const) (assemble-number const))
              ((string? const) (assemble-string const))
              (else (error "invalid const expression"))))


      (define (assemble-pair pair)
        (string-append "pairval("
                       (assemble-const (car pair))
                       ", "
                       (assemble-const (cdr pair))
                       ")"))


      (define (assemble-null) "null")


      (define (const? code)
        (and (tagged-list? code 'const)))


      (define (assemble-const const)
        (cond ((symbol? const) (assemble-symbol-value const))
              ((self-evaluating? const) (assemble-self-evaluating const))
              ((pair? const) (assemble-pair const))
              ((null? const) (assemble-null))
              (else (error "invalid const expression"))))


      (define (label? code)
        (and (tagged-list? code 'label)
             (integer? (cadr code))))


      (define (assemble-label label)
        (number->string label))


      (define (assemble-label-value label)
        (string-append "numvali("
                       (assemble-label label)
                       ")"))


      (define (register? code)
        (and (tagged-list? code 'reg)
             (symbol? (cadr code))))


      (define (assemble-register reg)
        (string-append "rm->"
                       (register-name reg)))


      (define (assemble-expression expression)
        (cond ((const? expression) (assemble-const (cadr expression)))
              ((label? expression) (assemble-label-value (cadr expression)))
              ((register? expression) (assemble-register (cadr expression)))
              (else (error "assemble: invalid expression"))))


      (define (stack-operation? instruction operation)
        (and (tagged-list? instruction operation)
             (symbol? (cadr instruction))))


      (define (save-env? instruction)
        (tagged-list? instruction 'saveenv))


      (define (assemble-save-env instruction)
        (string-append "saveenv(rm);\n"))


      (define (restore-env? instruction)
        (tagged-list? instruction 'restoreenv))


      (define (assemble-restore-env instruction)
        (string-append "restoreenv(rm);\n"))


      (define (save? instruction)
        (stack-operation? instruction 'save))


      (define (assemble-save instruction)
        (string-append "savereg(rm, "
                       (assemble-register (cadr instruction))
                       ");\n"))


      (define (restore? instruction)
        (stack-operation? instruction 'restore))


      (define (assemble-restore instruction)
        (string-append "restorereg(rm, (void *)&"
                       (assemble-register (cadr instruction))
                       ");\n"))


      (define (goto? instruction)
        (and (tagged-list? instruction 'goto)
             (or (label? (cadr instruction))
                 (register? (cadr instruction))
                 (proclabel? (cadr instruction)))))


      (define (assemble-goto instruction)
        (cond ((label? (cadr instruction))
               (string-append "gotolabel(rm, "
                              (assemble-label (cadadr instruction))
                              "); break;\n"))
              ((register? (cadr instruction))
               (string-append "gotoreg(rm, "
                              (assemble-register (cadadr instruction))
                              "); break;\n"))
              ((proclabel? (cadr instruction)) "gotoproc(rm); break;\n")))


      (define (takeproclabel? instruction)
        (and (pair? instruction)
             (eq? (car instruction) 'takeproclabel)))


      (define (assemble-takeproclabel instruction)
        "takeproclabel(rm);\n")


      (define (initreg? instruction)
        (and (tagged-list? instruction 'initreg)
             (symbol? (cadr instruction))
             (or (const? (caddr instruction))
                 (label? (caddr instruction))
                 (reg? (caddr instruction)))))


      (define (assemble-initreg instruction)
        (string-append "initreg((void *)&"
                       (assemble-register (cadr instruction))
                       ", "
                       (assemble-expression (caddr instruction))
                       ");\n"))


      (define (get-variable? instruction)
        (and (tagged-list? instruction 'get-variable)
             (symbol? (cadr instruction))
             (symbol? (caddr instruction))))


      (define (assemble-get-variable instruction)
        (string-append "getenvvar(rm, (void *)&"
                       (assemble-register (cadr instruction))
                       ", "
                       (assemble-symbol-name (caddr instruction))
                       ");\n"))


      (define (set-variable? instruction)
        (and (tagged-list? instruction 'set-variable-value)
             (symbol? (cadr instruction))
             (symbol? (caddr instruction))))


      (define (assemble-set-variable instruction)
        (string-append "setenvvar(rm, (void *)&"
                       (assemble-register (cadr instruction))
                       ", "
                       (assemble-symbol-name (caddr instruction))
                       ");\n"))


      (define (define-variable? instruction)
        (and (tagged-list? instruction 'define-variable)
             (symbol? (cadr instruction))))


      (define (assemble-define-variable instruction)
        (string-append "defenvvar(rm, "
                       (assemble-symbol-name (cadr instruction))
                       ");\n"))


      (define (branchval? instruction)
        (and (tagged-list? instruction 'branchval)
             (integer? (cadr instruction))))


      (define (assemble-branchval instruction)
        (string-append "if (branchval(rm, "
                       (assemble-label (cadr instruction))
                       ")) { break; }\n"))


      (define (initprocenv? instruction)
        (and (tagged-list? instruction 'init-proc-env)
             (or (null? (cadr instruction))
                 (pair? (cadr instruction))
                 (symbol? (cadr instruction)))))


      (define (assemble-initprocenv instruction)
        (string-append "initprocenv(rm, "
                       (assemble-const (cadr instruction))
                       ");\n"))


      (define (makeproc? instruction)
        (and (tagged-list? instruction 'make-compiled-procedure)
             (symbol? (cadr instruction))
             (integer? (caddr instruction))))


      (define (assemble-makeproc instruction)
        (string-append "mkcompiledprocreg(rm, (void *)&"
                       (assemble-register (cadr instruction))
                       ", "
                       (assemble-label (caddr instruction))
                       ");\n"))


      (define (initargs? instruction)
        (and (pair? instruction)
             (eq? (car instruction) 'initargs)))


      (define (assemble-initargs)
        "initargs(rm);\n")


      (define (addarg? instruction)
        (and (pair? instruction)
             (eq? (car instruction) 'addarg)))


      (define (assemble-addarg)
        "addarg(rm);\n")


      (define (branchproc? instruction)
        (and (tagged-list? instruction 'branchproc)
             (integer? (cadr instruction))))


      (define (assemble-branchproc instruction)
        (string-append "if (branchproc(rm, "
                       (assemble-label (cadr instruction))
                       ")) { break; }\n"))


      (define (applyprimitive? instruction)
        (and (tagged-list? instruction 'apply-primitive-procedure)
             (symbol? (cadr instruction))))


      (define (assemble-applyprimitive instruction)
        (string-append "applyprimitivereg(rm, (void *)&"
                       (assemble-register (cadr instruction))
                       ");\n"))


      (define (label-def? instruction)
        (integer? instruction))


      (define (assemble-label-def instruction)
        (string-append "case "
                       (assemble-label instruction)
                       ":\n"))


      (define (assemble-instruction instruction)
        (debug instruction)
        (cond ((save-env? instruction) (assemble-save-env instruction))
              ((restore-env? instruction) (assemble-restore-env instruction))
              ((save? instruction) (assemble-save instruction))
              ((restore? instruction) (assemble-restore instruction))
              ((goto? instruction) (assemble-goto instruction))
              ((takeproclabel? instruction) (assemble-takeproclabel instruction))
              ((initreg? instruction) (assemble-initreg instruction))
              ((get-variable? instruction) (assemble-get-variable instruction))
              ((set-variable? instruction) (assemble-set-variable instruction))
              ((define-variable? instruction) (assemble-define-variable instruction))
              ((branchval? instruction) (assemble-branchval instruction))
              ((initprocenv? instruction) (assemble-initprocenv instruction))
              ((makeproc? instruction) (assemble-makeproc instruction))
              ((initargs? instruction) (assemble-initargs))
              ((addarg? instruction) (assemble-addarg))
              ((branchproc? instruction) (assemble-branchproc instruction))
              ((applyprimitive? instruction) (assemble-applyprimitive instruction))
              ((label-def? instruction) (assemble-label-def instruction))
              (else (error "invalid assembly instruction"))))


      (define (assemble assembly)
        (if (null? assembly) ""
          (string-append (assemble-instruction (car assembly))
                         (assemble (cdr assembly)))))


        (define output-head "
        #include <stdlib.h>
        #include <stdio.h>
        #include \"../src-head/sys.h\"
        #include \"../src-head/error.h\"
        #include \"../src-head/sysio.h\"
        #include \"../src-head/number.h\"
        #include \"../src-head/string.h\"
        #include \"../src-head/compound-types.h\"
        #include \"../src-head/stack.h\"
        #include \"../src-head/io.h\"
        #include \"../src-head/value.h\"
        #include \"../src-head/register-machine.h\"
        #include \"../src-head/primitives.h\"
        #include \"../src-head/registry.h\"

        int main() {
            initsys();
            initmodule_sysio();
            initmodule_io();
            initmodule_number();
            initmodule_value();
            initmodule_primitives();
            initmodule_registry();

            regmachine rm = mkregmachine();
            for (;;) {
                long labelval = valrawint(rm->label);
                switch (labelval) {
                case 0:
                ")


        (define output-tail "
                default:
                    // printf(\"%s\\n\", sprintraw(rm->val));
                    return 0;
                }
            }

            freeregmachine(rm);

            freemodule_primitives();
            freemodule_value();
            freemodule_number();
            freemodule_io();
            freemodule_sysio();

            return 0;
        }
        ")


      (define (write-head)
        (write-file stdout output-head))


      (define (write-tail)
        (write-file stdout output-tail))


      (define (write-code output-code)
        (write-file stdout output-code))


      (define f (open-file "src-mm/mm.scm" file-mode-read))
      (define p (make-token-port f))
      (debug "read started")
      (define code (read p))
      (debug "read done")


      (debug "compiling")
      (define assembly (statements (compile code 'val 'return)))
      (debug "compile done")


      (define (print-assembly assembly)
        (if (null? assembly)
          (newline)
          (begin
            (write (car assembly))
            (newline)
            (print-assembly (cdr assembly)))))


      (define (write-output-code code)
        (write-head)
        (write-code code)
        (write-tail))


      ; (print-assembly assembly)
      (debug "assembling" assembly)
      (define output-code (assemble assembly))
      (debug "done")
      (write-output-code output-code)


      'ok)))


(define assembly (statements (compile code 'val 'return)))


(define (print-assembly assembly)
  (if (null? assembly)
    (newline)
    (begin
      (write (car assembly))
      (newline)
      (print-assembly (cdr assembly)))))


(define (write-output-code code)
  (write-head)
  (write-code code)
  (write-tail))


; (print-assembly assembly)
(define output-code (assemble assembly ""))
(write-output-code output-code)
