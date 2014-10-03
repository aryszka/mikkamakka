(define false #f)

(define true #t)

(define (memq item list)
  (cond ((null? list) false)
        ((eq? (car list) item) list)
        (else (memq item (cdr list)))))

(define (out exp)
  (display (cond ((eq? exp false) "false")
                 ((eq? exp true) "true")
                 (else exp))))

(define (log exp)
  (out exp)
  (newline))

(define (sprint exp)
  (with-output-to-string (lambda () (write exp))))

(define (entry? exp) (symbol? exp))

(define (entry-name exp) (symbol->string exp))

(define (make-instruction-sequence needs modifies statements)
  (list needs modifies statements))

(define (registers-needed s) (if (entry? s) '() (car s)))

(define (registers-modified s) (if (entry? s) '() (cadr s)))

(define (needs-register? seq reg)
  (memq reg (registers-needed seq)))

(define (modifies-register? seq reg)
  (memq reg (registers-modified seq)))

(define (statements s)
  (if (entry? s) (list s) (caddr s)))

(define (empty-instruction-sequence)
  (make-instruction-sequence '() '() '()))

(define (compile-linkage linkage)
  (cond ((eq? linkage 'return)
         (make-instruction-sequence
           '(cont) '()
           '((goto (reg cont)))))
        ((eq? linkage 'next)
         (empty-instruction-sequence))
        (else
          (make-instruction-sequence
            '() '()
            (list (list 'goto (list 'label linkage)))))))

(define (list-union s1 s2)
  (cond ((null? s1) s2)
        ((memq (car s1) s2) (list-union (cdr s1) s2))
        (else (cons (car s1) (list-union (cdr s1) s2)))))

(define (list-difference s1 s2)
  (cond ((null? s1) '())
        ((memq (car s1) s2) (list-difference (cdr s1) s2))
        (else (cons (car s1) (list-difference (cdr s1) s2)))))

(define (append-instruction-sequences . seqs)
  (define (append-2-sequences seq1 seq2)
    (make-instruction-sequence
      (list-union (registers-needed seq1)
                  (list-difference (registers-needed seq2)
                                   (registers-modified seq1)))
      (list-union (registers-modified seq1)
                  (registers-modified seq2))
      (append (statements seq1) (statements seq2))))
  (define (append-seq-list seqs)
    (if (null? seqs)
      (empty-instruction-sequence)
      (append-2-sequences (car seqs)
                          (append-seq-list (cdr seqs)))))
  (append-seq-list seqs))

(define (append-preserving regs seq1 seq2)
  (if (null? regs)
    (append-instruction-sequences seq1 seq2)
    (if (and (needs-register? seq2 (car regs))
             (modifies-register? seq1 (car regs)))
      (append-preserving (cdr regs)
                         (make-instruction-sequence
                           (list-union (list (car regs))
                                       (registers-needed seq1))
                           (list-difference (registers-modified seq1)
                                            (list (car regs)))
                           (append (list (list 'save (car regs)))
                                   (append (statements seq1)
                                           (list (list 'restore (car regs))))))
                         seq2)
      (append-preserving (cdr regs) seq1 seq2))))

(define (parallel-instruction-sequences seq1 seq2)
  (make-instruction-sequence
    (list-union (registers-needed seq1)
                (registers-needed seq2))
    (list-union (registers-modified seq1)
                (registers-modified seq2))
    (append (statements seq1) (statements seq2))))

(define (tack-on-instruction-sequence seq body-seq)
  (make-instruction-sequence
    (registers-needed seq)
    (registers-modified seq)
    (append (statements seq) (statements body-seq))))

(define (end-with-linkage linkage instruction-sequence)
  (append-preserving '(cont)
                     instruction-sequence
                     (compile-linkage linkage)))

(define (make-begin seq) (cons 'begin seq))

(define (compile-sequence seq target linkage)
  (if (null? (cdr seq))
    (compile (car seq) target linkage)
    (append-preserving '(env cont)
                       (compile (car seq) target 'next)
                       (compile-sequence (cdr seq) target linkage))))

(define label-counter -1)

(define (new-label-number)
  (set! label-counter (+ label-counter 1))
  label-counter)

(define (make-label name)
  (string->symbol
    (string-append (symbol->string name)
                   (number->string (new-label-number)))))

(define all-regs '(env proc val args cont))

(define (compile-self-evaluating exp target linkage)
  (end-with-linkage
    linkage
    (make-instruction-sequence
      '()
      (list target)
      (list (list 'assign target (list 'const exp))))))

(define (tagged-list? exp tag)
  (and (pair? exp) (eq? (car exp) tag)))

(define (quoted? exp) (tagged-list? exp 'quote))

(define (text-of-quotation exp) (cadr exp))

(define (compile-quoted exp target linkage)
  (end-with-linkage
    linkage
    (make-instruction-sequence
      '()
      (list target)
      (list (list
              'assign
              target
              (list 'const
                    (text-of-quotation exp)))))))

(define (variable? exp) (symbol? exp))

(define (compile-variable exp target linkage)
  (end-with-linkage
    linkage
    (make-instruction-sequence
      '(env)
      (list target)
      (list (list
              'assign
              target
              '(op lookup-variable-value)
              (list 'const exp)
              '(reg env))))))

(define (assignment? exp) (tagged-list? exp 'set!))

(define (assignment-variable exp) (cadr exp))

(define (assignment-value exp) (caddr exp))

(define (compile-assignment exp target linkage)
  (let ((var (assignment-variable exp))
        (value (compile (assignment-value exp) 'val 'next)))
    (end-with-linkage
      linkage
      (append-preserving '(env)
                         value
                         (make-instruction-sequence
                           '(env val)
                           (list target)
                           (list (list 'perform
                                       '(op set-variable-value!)
                                       (list 'const var)
                                       '(reg val)
                                       '(reg env))))))))

(define (definition? exp) (tagged-list? exp 'define))

(define (definition-variable exp)
  (if (pair? (cadr exp)) (caadr exp) (cadr exp)))

(define (definition-value exp)
  (if (pair? (cadr exp))
    (cons 'lambda (cons (cdadr exp) (cddr exp)))
    (caddr exp)))

(define (compile-definition exp target linkage)
  (let ((var (definition-variable exp))
        (value (compile (definition-value exp) 'val 'next)))
    (end-with-linkage
      linkage
      (append-preserving '(env)
                         value
                         (make-instruction-sequence
                           '(env val)
                           (list target)
                           (list (list 'perform
                                       '(op define-variable!)
                                       (list 'const var)
                                       '(reg val)
                                       '(reg env))))))))

(define (if? exp) (tagged-list? exp 'if))

(define (if-predicate exp) (cadr exp))

(define (if-consequent exp) (caddr exp))

(define (if-alternative exp) (cadddr exp))

(define (compile-if exp target linkage)
  (let ((after-if (make-label 'afterIf))
        (true-branch (make-label 'trueBranch))
        (false-branch (make-label 'falseBranch)))
    (let ((consequent-linkage
            (if (eq? linkage 'next) after-if linkage)))
      (let ((predicate (compile (if-predicate exp) 'val 'next))
            (consequent (compile (if-consequent exp) target consequent-linkage))
            (alternative (compile (if-alternative exp) target linkage)))
        (append-preserving '(env cont)
                           predicate
                           (append-instruction-sequences
                             (make-instruction-sequence
                               '(val) '()
                               (list '(test (op false?) (reg val))
                                     (list 'branch (list 'label false-branch))))
                             (parallel-instruction-sequences
                               (append-instruction-sequences true-branch consequent)
                               (append-instruction-sequences false-branch alternative))
                             after-if))))))

(define (lambda? exp) (tagged-list? exp 'lambda))

(define (lambda-parameters exp) (cadr exp))

(define (lambda-body exp) (cddr exp))

(define (compile-lambda-body exp proc-entry)
  (append-instruction-sequences
    (make-instruction-sequence
      '(env proc args) '(env)
      (list proc-entry
            '(assign env (op compiled-procedure-env) (reg proc))
            (list 'assign
                  'env
                  '(op extend-environment)
                  (list 'const (lambda-parameters exp))
                  '(reg args)
                  '(reg env))))
    (compile-sequence (lambda-body exp) 'val 'return)))

(define (compile-lambda exp target linkage)
  (let ((after-lambda (make-label 'afterLambda))
        (entry (make-label 'entry)))
    (let ((lambda-linkage
            (if (eq? linkage 'next) after-lambda linkage)))
      (append-instruction-sequences
        (tack-on-instruction-sequence
          (end-with-linkage
            lambda-linkage
            (make-instruction-sequence
              '(env)
              (list target)
              (list (list 'assign
                          target
                          '(op make-compiled-procedure)
                          (list 'label entry)
                          '(reg env)))))
          (compile-lambda-body exp entry))
        after-lambda))))

(define (cond? exp) (tagged-list? exp 'cond))

(define (cond-predicate clause) (car clause))

(define (cond-else-clause? clause) (eq? (cond-predicate clause) 'else))

(define (make-if predicate consequent alternative)
  (list 'if predicate consequent alternative))

(define (sequence->exp seq)
  (cond ((null? seq) '())
        ((null? (cdr seq)) (car seq))
        (else (make-begin seq))))

(define (expand-clauses clauses)
  (if (null? clauses)
    'true
    (let ((first (car clauses))
          (the-rest (cdr clauses)))
      (cond ((not (pair? first))
             (error 'expand-clauses "invalid syntax" first))
            ((cond-else-clause? first)
             (if (null? the-rest)
               (sequence->exp (cdr first))
               (error 'expand-clauses "else clause isn't last" clauses)))
            (else (make-if (car first)
                           (sequence->exp (cdr first))
                           (expand-clauses the-rest)))))))

(define (cond->if exp) (expand-clauses (cdr exp)))

(define (and? exp) (tagged-list? exp 'and))

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

(define (and->if exp) (expand-and (cdr exp)))

(define (or? exp) (tagged-list? exp 'or))

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

(define (or->if exp) (expand-or (cdr exp)))

(define (let? exp) (tagged-list? exp 'let))

(define (check-let exp)
  (cond ((< (length exp) 3)
         (error 'check-let "invalid arity" exp))
        ((not (pair? (cadr exp)))
         (error 'check-let "invalid syntax" exp))))

(define (let-variables defs) (map car defs))

(define (let-values defs) (map cadr defs))

(define (make-lambda parameters body)
  (cons 'lambda (cons parameters body)))

(define (let->procedure exp)
  (cons (make-lambda (let-variables (cadr exp))
                     (cddr exp))
        (let-values (cadr exp))))

(define (begin? exp) (tagged-list? exp 'begin))

(define (begin-actions exp) (cdr exp))

(define (application? exp) (pair? exp))

(define (operands exp) (cdr exp))

(define (code-to-get-rest-args operands)
  (let ((code-for-next-arg
          (append-preserving
            '(args)
            (car operands)
            (make-instruction-sequence
              '(val args) '(args)
              '((assign args
                        (op cons)
                        (reg val)
                        (reg args)))))))
    (if (null? (cdr operands))
      code-for-next-arg
      (append-preserving
        '(env)
        code-for-next-arg
        (code-to-get-rest-args (cdr operands))))))

(define (construct-arglist operands)
  (let ((operands (reverse operands)))
    (if (null? operands)
      (make-instruction-sequence
        '() '(args)
        '((assign args (const ()))))
      (let ((code-to-get-last-arg
              (append-instruction-sequences
                (car operands)
                (make-instruction-sequence
                  '(val) '(args)
                  '((assign args (op list) (reg val)))))))
        (if (null? (cdr operands))
          code-to-get-last-arg
          (append-preserving
            '(env)
            code-to-get-last-arg
            (code-to-get-rest-args (cdr operands))))))))

; (define (compile-proc-appl target linkage)
;   (cond ((and (eq? target 'val) (not (eq? linkage 'return)))
;          (make-instruction-sequence
;            '(proc) all-regs
;            (list (list 'assign
;                        'cont
;                        (list 'label linkage))
;                  '(assign val
;                           (op compiled-procedure-entry)
;                           (reg proc))
;                  '(goto (reg val)))))
;         ((and (not (eq? target 'val))
;               (not (eq? linkage 'return)))
;          (let ((proc-return (make-label 'procReturn)))
;            (make-instruction-sequence
;              '(proc) all-regs
;              (list (list 'assign
;                          'cont
;                          (list 'label proc-return))
;                    '(assign val
;                             (op compiled-procedure-entry)
;                             (reg proc))
;                    '(goto (reg val))
;                    proc-return
;                    (list 'assign target '(reg val))
;                    (list 'goto (list 'label linkage))))))
;         ((and (eq? target 'val) (eq? linkage 'return))
;          (make-instruction-sequence
;            '(proc cont)
;            all-regs
;            '((assign val
;                      (op compiled-procedure-entry)
;                      (reg proc))
;              (goto (reg val)))))
;         (else
;           (error 'compile-proc-appl
;                  "return linkage, target not val"
;                  target))))
; 
; (define (compile-procedure-call target linkage)
;   (let ((after-call (make-label 'afterCall))
;         (primitive-branch (make-label 'primitiveBranch))
;         (compiled-branch (make-label 'compiledBranch)))
;     (let ((compiled-linkage
;             (if (eq? linkage 'next) after-call linkage)))
;       (append-instruction-sequences
;         (make-instruction-sequence
;           '(proc) '()
;           (list '(test (op primitive-procedure?) (reg proc))
;                 (list 'branch (list 'label primitive-branch))))
;         (parallel-instruction-sequences
;           (append-instruction-sequences
;             compiled-branch
;             (compile-proc-appl target compiled-linkage))
;           (append-instruction-sequences
;             primitive-branch
;             (end-with-linkage
;               linkage
;               (make-instruction-sequence
;                 '(proc args) (list target)
;                 (list (list 'assign
;                             target
;                             '(op apply-primitive-procedure)
;                             '(reg proc)
;                             '(reg args)))))))
;         after-call))))

; proc args and cont if return
; all-regs
; (perform (op procedure-call) (const regs) compiled-linkage)
; if not return and not val:
; (perform (op procedure-call) (const regs) proc-return)
; proc-return
; target = val
; compiled-linkage

(define (compile-procedure-call target linkage)
  (let ((after-call (make-label 'afterCall))
        (uses (if (eq? linkage 'return)
                '(proc args cont)
                '(proc args)))
        (not-return-not-val? (and (not (eq? linkage 'return))
                                  (not (eq? target 'val)))))
    (let ((calling-linkage
            (cond (not-return-not-val? (make-label 'procReturn))
                  ((eq? linkage 'next) after-call)
                  (else 'false))))
      (let ((statements
              (if not-return-not-val?
                (list (list 'perform-continue
                            '(op procedure-call)
                            'regs
                            (list 'label calling-linkage))
                      calling-linkage
                      (list 'assign
                            target
                            '(reg val)))
                (list (list 'perform-continue
                            '(op procedure-call)
                            'regs
                            (list 'label calling-linkage))))))
        (make-instruction-sequence
          uses
          all-regs
          (if (eq? calling-linkage 'false)
            statements
            (append statements (list after-call))))))))

(define (operator exp) (car exp))

(define (compile-application exp target linkage)
  (let ((procedure (compile (operator exp) 'proc 'next))
        (operands
          (map (lambda (operand) (compile operand 'val 'next))
               (operands exp))))
    (append-preserving
      '(env cont)
      procedure
      (append-preserving
        '(proc cont)
        (construct-arglist operands)
        (compile-procedure-call target linkage)))))

(define (compile exp target linkage)
  (cond ((self-evaluating? exp) (compile-self-evaluating exp target linkage))
        ((quoted? exp)
         (compile-quoted exp target linkage))
        ((variable? exp)
         (compile-variable exp target linkage))
        ((assignment? exp)
         (compile-assignment exp target linkage))
        ((definition? exp)
         (compile-definition exp target linkage))
        ((if? exp) (compile-if exp target linkage))
        ((lambda? exp) (compile-lambda exp target linkage))
        ((cond? exp) (compile (cond->if exp) target linkage))
        ((and? exp) (compile (and->if exp) target linkage))
        ((or? exp) (compile (or->if exp) target linkage))
        ((let? exp) (check-let exp)
                    (compile (let->procedure exp) target linkage))
        ((begin? exp)
         (compile-sequence (begin-actions exp)
                           target
                           linkage))
        ((application? exp)
         (compile-application exp target linkage))
        (else (error 'compile "invalid expression type" exp))))

(define (optimize-asm asm) asm)

(define (register? exp) (tagged-list? exp 'reg))

(define (register-name exp)
  (string-append "regs." (symbol->string (cadr exp))))

(define (label? exp) (tagged-list? exp 'label))

(define (label-name exp) (symbol->string (cadr exp)))

(define (make-assembly-builder)
  (let ((builder '()))
    (lambda (mutate get)
      (set! builder (mutate builder))
      (get builder))))

(define (assembly-append builder string)
  (builder (lambda (b) (cons string b))
           identity))

(define (assembly->string builder)
  (builder (lambda (b) b)
           (lambda (b) (apply string-append (reverse b)))))

(define (assemble-pair builder exp)
  (assembly-append builder "[")
  (assemble-const builder (car exp))
  (assembly-append builder ",")
  (assemble-const builder (cdr exp))
  (assembly-append builder "]"))

(define (const? exp) (tagged-list? exp 'const))

(define (assemble-const builder exp)
  (cond ((self-evaluating? exp)
         (assembly-append builder (sprint exp)))
        ((symbol? exp)
         (assembly-append builder "[")
         (assembly-append builder (sprint (symbol->string exp)))
         (assembly-append builder "]"))
        ((null? exp)
         (assembly-append builder "[]"))
        ((pair? exp)
         (assemble-pair builder exp))
        (else (error 'assemble-const "invalid const" exp))))

(define (assemble-op-call builder op-name . args)
  (define (assemble-op-args first? args)
    (cond ((not (null? args))
           (cond ((not first?)
                  (assembly-append builder ",")))
           (cond ((register? (car args))
                  (assembly-append builder (register-name (car args))))
                 ((label? (car args))
                  (assembly-append builder (label-name (car args))))
                 ((const? (car args))
                  (assemble-const builder (cadar args)))
                 ((eq? (car args) 'regs)
                  (assembly-append builder "regs"))
                 ((eq? (car args) 'false)
                  (assembly-append builder "false"))
                 (else (error 'assemble-op-call
                              "invalid assembly argument"
                              (car args))))
           (assemble-op-args false (cdr args)))))
  (assembly-append builder op-name)
  (assembly-append builder "(")
  (assemble-op-args true args)
  (assembly-append builder ")"))

(define (assemble-lookup-variable builder exp)
  (assemble-op-call builder
                    "ops.lookupVar"
                    (caddr exp)
                    (cadr exp)))

(define (assemble-set-variable builder exp)
  (assemble-op-call builder
                    "ops.setVar"
                    (cadddr exp)
                    (cadr exp)
                    (caddr exp)))

(define (assemble-define-variable builder exp)
  (assemble-op-call builder
                    "ops.defineVar"
                    (cadddr exp)
                    (cadr exp)
                    (caddr exp)))

(define (assemble-false-check builder exp)
  (assemble-op-call builder
                    "ops.isFalse"
                    (cadr exp)))

(define (assemble-compiled-procedure-env builder exp)
  (assemble-op-call builder
                    "ops.compiledProcedureEnv"
                    (cadr exp)))

(define (assemble-extend-environment builder exp)
  (assemble-op-call builder
                    "ops.extendEnv"
                    (cadddr exp)
                    (cadr exp)
                    (caddr exp)))

(define (assemble-make-procedure builder exp)
  (assemble-op-call builder
                    "ops.makeProcedure"
                    (cadr exp)
                    (caddr exp)))

(define (assemble-cons-op builder exp)
  (assemble-op-call builder
                    "ops.cons"
                    (cadr exp)
                    (caddr exp)))

(define (assemble-list-op builder exp)
  (apply assemble-op-call
         (cons builder (cons "ops.list" (cdr exp)))))

(define (assemble-compiled-entry builder exp)
  (assemble-op-call builder
                    "ops.compiledEntry"
                    (cadr exp)))

(define (assemble-apply builder exp)
  (assemble-op-call builder
                    "ops.call"
                    (cadr exp)
                    (caddr exp)))

(define (assembly-op? exp) (tagged-list? exp 'op))

(define (op-name exp) (cadar exp))

(define (assemble-op builder exp)
  (cond ((eq? (op-name exp) 'lookup-variable-value)
         (assemble-lookup-variable builder exp))
        ((eq? (op-name exp) 'set-variable-value!)
         (assemble-set-variable builder exp))
        ((eq? (op-name exp) 'define-variable!)
         (assemble-define-variable builder exp))
        ((eq? (op-name exp) 'false?)
         (assemble-false-check builder exp))
        ((eq? (op-name exp) 'compiled-procedure-env)
         (assemble-compiled-procedure-env builder exp))
        ((eq? (op-name exp) 'extend-environment)
         (assemble-extend-environment builder exp))
        ((eq? (op-name exp) 'make-compiled-procedure)
         (assemble-make-procedure builder exp))
        ((eq? (op-name exp) 'cons)
         (assemble-cons-op builder exp))
        ((eq? (op-name exp) 'list)
         (assemble-list-op builder exp))
        ((eq? (op-name exp) 'compiled-procedure-entry)
         (assemble-compiled-entry builder exp))
        ((eq? (op-name exp) 'procedure-call)
         (assemble-apply builder exp))
        (else (error 'assemble-op "invalid assembly operation" exp))))

(define (followed-by-entry? asm)
  (and (not (null? (cdr asm)))
       (entry? (cadr asm))))

(define (assemble-close-entry builder asm)
  (assembly-append builder "return ")
  (assembly-append
    builder
    (if (null? asm)
      "false"
      (entry-name (car asm))))
  (assembly-append builder ";};"))

(define (assembly-entry? asm) (entry? (car asm)))

(define (assemble-entry builder asm)
  (assembly-append builder "var ")
  (assembly-append builder (entry-name (car asm)))
  (assembly-append builder "=function(){")
  (cond ((followed-by-entry? asm)
         (assemble-close-entry builder (cdr asm))))
  (assemble builder (cdr asm)))

(define (assembly-goto? asm)
  (tagged-list? (car asm) 'goto))

(define (assemble-goto builder asm)
  (assembly-append builder "return ")
  (assembly-append
    builder
    (cond ((register? (cadar asm))
           (register-name (cadar asm)))
          ((label? (cadar asm))
           (label-name (cadar asm)))
          (else
            (error 'assemble-goto "invalid assembly goto" asm))))
  (assembly-append builder ";};")
  (assemble builder (cdr asm)))

(define (assembly-save? asm)
  (tagged-list? (car asm) 'save))

(define (assemble-save builder asm)
  (assembly-append builder "save(")
  (assembly-append builder (register-name (car asm)))
  (assembly-append builder ");")
  (cond ((followed-by-entry? asm)
         (assemble-close-entry builder (cdr asm))))
  (assemble builder (cdr asm)))

(define (assembly-restore? asm)
  (tagged-list? (car asm) 'restore))

(define (assemble-restore builder asm)
  (assembly-append builder (register-name (car asm)))
  (assembly-append builder "= restore();")
  (cond ((followed-by-entry? asm)
         (assemble-close-entry builder (cdr asm))))
  (assemble builder (cdr asm)))

(define (assembly-assign? asm) (tagged-list? (car asm) 'assign))

(define (assemble-assign builder asm)
  (assembly-append builder (register-name (car asm)))
  (assembly-append builder "=")
  (cond ((register? (caddar asm))
         (assembly-append builder (register-name (caddar asm))))
        ((label? (caddar asm))
         (assembly-append builder (label-name (caddar asm))))
        ((const? (caddar asm))
         (assemble-const builder (cadr (caddar asm))))
        ((assembly-op? (caddar asm))
         (assemble-op builder (cddar asm)))
        (else (error 'assemble-assign "invalid assembly assign" asm)))
  (assembly-append builder ";")
  (cond ((followed-by-entry? asm)
         (assemble-close-entry builder (cdr asm))))
  (assemble builder (cdr asm)))

(define (assembly-perform? asm) (tagged-list? (car asm) 'perform))

(define (assemble-perform builder asm)
  (assemble-op builder (cdar asm))
  (assembly-append builder ";")
  (assemble builder (cdr asm)))

(define (assembly-perform-continue? asm)
  (tagged-list? (car asm) 'perform-continue))

(define (assemble-perform-continue builder asm)
  (assembly-append builder "return ")
  (assemble-op builder (cdar asm))
  (assembly-append builder ";};")
  (assemble builder (cdr asm)))

(define (assembly-test-branch? asm)
  (and (eq? (caar asm) 'test)
       (eq? (caadr asm) 'branch)))

(define (assemble-test-branch builder asm)
  (assembly-append builder "if(")
  (cond ((register? (cadar asm))
         (assembly-append builder (register-name (cadar asm))))
        ((assembly-op? (cadar asm))
         (assemble-op builder (cdar asm)))
        (else (error 'assemble-test-branch "invalid test" asm)))
  (assembly-append builder "){return ")
  (assembly-append builder (label-name (cadadr asm)))
  (assembly-append builder ";};")
  (assemble-close-entry builder (cddr asm))
  (assemble builder (cddr asm)))

(define (assemble builder asm)
  (cond ((null? asm)
         (assemble-close-entry builder asm)
         builder)
        ((assembly-entry? asm)
         (assemble-entry builder asm))
        ((assembly-goto? asm)
         (assemble-goto builder asm))
        ((assembly-save? asm)
         (assemble-save builder asm))
        ((assembly-restore? asm)
         (assemble-restore builder asm))
        ((assembly-assign? asm)
         (assemble-assign builder asm))
        ((assembly-perform? asm)
         (assemble-perform builder asm))
        ((assembly-perform-continue? asm)
         (assemble-perform-continue builder asm))
        ((assembly-test-branch? asm)
         (assemble-test-branch builder asm))
        (else (error 'assemble "invalid assembly" asm))))

(define (compile-js exp)
  (assembly->string
    (assemble (make-assembly-builder)
              (optimize-asm
                (cons 'start (statements (compile exp 'val 'next)))))))

; (define exp
;   '((lambda ()
; 
;       (define (factorial n)
;         (if (= n 0)
;           1
;           (* n (factorial (- n 1)))))
; 
;       (define n 150)
; 
;       (out (factorial n))
;       (out (((lambda () factorial)) n))
; 
;       (apply * '())
;       (apply factorial '(1))
; 
;       (out (apply * '()))
;       (out (apply factorial '(1)))
; 
;       (define (list . x) x)
; 
;       (out (apply apply (list * '(1 2 3))))
;       (out (apply apply (list apply (list * '(1 2 3)))))
; 
;       (define (loop counter)
;         (if (= counter 0)
;           (out 'done)
;           (begin
;             (factorial n)
;             (loop (- counter 1)))))
;       (loop 100)
; 
;       'ok)))

; (print (cons 'start (statements (compile exp 'val 'next))))

(define exp
  '((lambda ()
      (define (assert condition message)
        (cond ((not condition) (out message))))

      (define (map f l)
        (if (null? l)
          '()
          (cons (f (car l)) (map f (cdr l)))))

      (define (append . lists)
        (cond ((null? lists) '())
              ((null? (car lists)) (apply append (cdr lists)))
              (else (cons (car (car lists))
                          (apply append (cons (cdr (car lists))
                                              (cdr lists)))))))

      (define (reverse l)
        (if (null? l)
          '()
          (append (reverse (cdr l)) (list (car l)))))

      (define (symbol-eq? left right)
        (and (symbol? left)
             (symbol? right)
             (eq? (symbol->string left)
                  (symbol->string right))))

      (define (tagged-list? exp tag)
        (and (pair? exp)
             (eq? (car exp) tag)))

      (define (quote? exp) (tagged-list? exp 'quote))

      (define (quote-text exp) (car (cdr exp)))

      (define (quote-eq? left right)
        (and (quote? left)
             (quote? right)
             (eq? (quote-text left)
                  (quote-text right))))

      (define (eq? . x)
        (cond ((or (null? x) (null? (cdr x))) true)
              (else
                (and (or (and (null? (car x)) (null? (car (cdr x))))
                         (symbol-eq? (car x) (car (cdr x)))
                         (quote-eq? (car x) (car (cdr x)))
                         (primitive-eq? (car x) (car (cdr x))))
                     (apply eq? (cdr x))))))

      (define (list . x) x)

      (define (not x) (eq? x false))

      (define test-struct (struct))
      (assert (struct? test-struct) "is struct")
      (assert (not (struct-defined? test-struct 'a-name))
              "not defined")
      (struct-define test-struct 'a-name 1)
      (assert (struct-defined? test-struct 'a-name)
              "defined")
      (assert (eq? (struct-lookup test-struct 'a-name) 1)
              "lookup value 1")
      (struct-set! test-struct 'a-name 2)
      (assert (eq? (struct-lookup test-struct 'a-name) 2)
              "lookup value 2")
      ; (define test-struct (struct (a-name 3) (another-name 4)))


      (assert (eq? (cond (1 2) (else 3)) 2) "cond clause")
      (assert (eq? (cond (false 2) (else 3)) 3) "cond clause")

      (assert (eq? (and 1 2 3) 3) "and true")
      (assert (eq? (and 1 false 3) false) "and false")

      (assert (eq? (or 1 2 3) 1) "or true")
      (assert (eq? (or false 1 2) 1) "or false")

      (assert (eq? (let ((a 1) (b 2)) (- a b)) -1) "let")

      (assert (not (symbol-eq? 'a 'b)) "symbol-eq? a b")
      (assert (not (symbol-eq? 'a '1)) "symbol-eq? a 1")
      (assert (symbol-eq? 'a 'a) "symbol-eq? a a")

      (assert (not (quote? 1)) "quote? 1")
      (assert (not (quote? (quote 1))) "quote? '1")
      (assert (quote? (quote (quote 1))) "quote? ''1")
      (assert (not (quote? 'a)) "quote? 'a")
      (assert (quote? (quote 'a)) "quote? ''a")
      (assert (quote? ''a) "quote? ''a")
      (assert (eq? (quote-text (quote 'a)) 'a) "quote text (quote a)")
      (assert (eq? (quote-text ''a) 'a) "quote text 'a")
      (assert (not (quote-eq? 1 2)) "quote-eq? 1 2")
      (assert (not (quote-eq? 'a 'a)) "quote-eq? a a")
      (assert (quote-eq? ''a ''a) "quote-eq? ''a ''a")
      (assert (not (quote-eq? 'a ''a)) "quote-eq? 'a ''a")
      (assert (not (quote-eq? 'a 'b)) "quote-eq? 'a 'b")
      (assert (not (quote-eq? ''a ''b)) "quote-eq? ''a ''b")
      (assert (quote-eq? '''a '''a) "quote-eq? '''a '''a")

      (assert (eq?) "eq?")
      (assert (eq? 1) "eq?")
      (assert (eq? 1 1) "eq? 1 1")
      (assert (eq? 1 1 1) "eq? 1 1 1")
      (assert (not (eq? 1 2)) "eq? 1 2")
      (assert (not (eq? 1 1 2)) "eq? 1 1 2")
      (assert (not (eq? 1 2 3)) "eq? 1 2 3")
      (assert (eq? "a" "a") "eq? \"a\" \"a\"")
      (assert (not (eq? "a" "b")) "eq? \"a\" \"b\"")
      (assert (eq? '() '()) "eq? '() '()")
      (assert (not (eq? '() 1)) "eq? '() 1")
      (assert (eq? 'a 'a) "eq? 'a 'a")
      (assert (not (eq? 'a 'b)) "eq? 'a 'b")
      (assert (eq? ''a ''a) "eq? ''a ''a")
      (assert (not (eq? ''a ''b)) "eq? ''a ''b")
      (define (f x) x)
      (assert (eq? f f) "eq? ref")
      (define g f)
      (assert (eq? f g) "eq? ref copy")
      (define (h x) x)
      (assert (not (eq? f h)) "no eq? ref")

      (define l (list 1 2 3))
      (assert (eq? (car l) 1) "list first")
      (assert (eq? (car (cdr l)) 2) "list second")
      (assert (eq? (car (cdr (cdr l))) 3) "list third")

      (define l (apply list '(1 2 3)))
      (assert (eq? (car l) 1) "list apply first")
      (assert (eq? (car (cdr l)) 2) "list apply second")
      (assert (eq? (car (cdr (cdr l))) 3) "list apply third")

      (define l (call list 1 2 3))
      (assert (eq? (car l) 1) "list call first")
      (assert (eq? (car (cdr l)) 2) "list call second")
      (assert (eq? (car (cdr (cdr l))) 3) "list call third")
      
      (define l (map (lambda (x) (* 2 x)) '(1 2 3)))
      (assert (eq? (car l) 2) "map first")
      (assert (eq? (car (cdr l)) 4) "map second")
      (assert (eq? (car (cdr (cdr l))) 6) "map third")

      (define (make-port object)
        (let ((port (struct)))
          (struct-define port 'port-object object)
          port))

      (define (port? port)
        (and (struct? port)
             (struct-defined? port 'port-object)))

      (define (input-port? port)
        (and (port? port)
             ((struct-lookup port 'port-object)
              'input-port?)))
      (define (output-port? port)
        (and (port? port)
             ((struct-lookup port 'port-object)
              'output-port?)))

      (define (close-port port)
        ((struct-lookup port 'port-object) 'close))

      (close-port (make-port (lambda (x) x)))

      (define (end-of-file) end-of-file)

      (define (eof-object? eof) (eq? eof end-of-file))

      (assert (not (port? 'not-a-port)) "not a port?")
      (assert (port? (make-port 'port-object)) "port?")
      (assert (input-port? (make-port (lambda (message)
                                        (eq? message 'input-port?))))
              "input-port?")
      (assert (output-port? (make-port (lambda (message)
                                         (eq? message 'output-port?))))
              "output-port?")
      (assert (eof-object? end-of-file) "end-of-file")

      (define (read-string port . args)
        (apply (struct-lookup port 'port-object)
               (cons 'read-string args)))

      (define (write-string port string)
        ((struct-lookup port 'port-object) 'write-string string))

      (define (read-char port)
        ((struct-lookup port 'port-object) 'read-char))

      (define (write-char port char)
        ((struct-lookup port 'port-object) 'write-char char))

      (define (read port)
        (read-string port)
        '())

      (define (write port)
        ((struct-lookup port 'port-object 'write-string "")))

      (define (open-string-port . strings)
        (let ((ref 0))
          (define (read-string . args)
            (define (read-string parts length)
              (cond ((and (null? strings) (null? parts)) end-of-file)
                    ((< length 0)
                     (let ((string
                             (apply
                               string-append
                               (cons (string-copy (car strings) ref)
                                     (cdr strings)))))
                       (set! strings '())
                       (set! ref 0)
                       string))
                    ((null? strings) (apply string-append (reverse parts)))
                    ((>= (+ ref length) (string-length (car strings)))
                     (let ((parts (cons (string-copy (car strings) ref)
                                        parts)))
                       (set! strings (cdr strings))
                       (set! ref 0)
                       (read-string parts
                                    (- length (string-length (car parts))))))
                    (else
                      (let ((string
                              (apply
                                string-append
                                (reverse
                                  (cons (string-copy (car strings) ref (+ ref length))
                                        parts)))))
                        (set! ref (+ ref length))
                        string))))
            (read-string '() (if (null? args) -1 (car args))))
          (define (write-string string)
            (set! strings (append strings (list string))))
          (define (close)
            (set! strings '())
            (set! ref 0))
          (make-port (lambda (message . args)
                       (cond ((eq? message 'port) true)
                             ((eq? message 'input-port?) true)
                             ((eq? message 'output-port?) true)
                             ((eq? message 'close) (close))
                             ((eq? message 'read-string) (apply read-string args))
                             ((eq? message 'write-string) (write-string (car args)))
                             ((eq? message 'read-char) (read-string 1))
                             ((eq? message 'write-char) (write-string (string-copy (car args) 0 1)))
                             (else (error "invalid message" message)))))))

      (define string-port (open-string-port))
      (assert (port? string-port) "port? string-port")
      (assert (input-port? string-port) "input-port? string-port")
      (assert (output-port? string-port) "output-port? string-port")
      (assert (eof-object? (read-string string-port)) "read eof")
      (write-string string-port "hello")
      (assert (eq? (read-string string-port) "hello") "read hello")
      (assert (eof-object? (read-string string-port)) "read eof again")

      (define string-port (open-string-port "some string"))
      (assert (eq? (read-string string-port) "some string") "read initial")
      (assert (eof-object? (read-string string-port)) "read eof after initial")

      (define string-port (open-string-port))
      (write-string string-port "some")
      (write-string string-port " string")
      (assert (eq? (read-string string-port) "some string") "read all")
      (assert (eof-object? (read-string string-port)) "after read all")

      (define string-port (open-string-port))
      (write-string string-port "some")
      (write-string string-port " string")
      (assert (eq? (read-string string-port 2) "so") "read 2")
      (assert (eq? (read-string string-port 4) "me s") "read 4")
      (assert (eq? (read-string string-port 9) "tring") "read to end")
      (assert (eof-object? (read-string string-port 3)) "read eof after parts")

      (define string-port (open-string-port))
      (write-string string-port "some")
      (write-string string-port " string")
      (write-string string-port " indeed")
      (assert (eq? (read-string string-port 15) "some string ind") "read across parts")
      (assert (eq? (read-string string-port 15) "eed") "read rest")
      (assert (eof-object? (read-string string-port 15)) "eof after reading across parts")

      (define a (+ 1 2))
      (define b (call/cc (lambda (return) (return (+ a 1)))))
      (define c (+ a b))
      (assert (eq? c 7) "simple call/cc")

      (define (search wanted? list)
        (call/cc
          (lambda (return)
            (define (iter list)
              (cond ((null? list) false)
                    ((wanted? (car list)) (return (car list)))
                    (else (iter (cdr list)))))
            (iter list))))
      (assert (eq? (search (lambda (x) (eq? x 3)) '(1 2 3 4)) 3) "call/cc in iteration")

      (define (treat-element element like-it)
        (cond ((eq? element 3) (like-it element))))
      (define (search list)
        (call/cc
          (lambda (return)
            (define (iter list)
              (cond ((null? list) false))
              (treat-element (car list) return)
              (iter (cdr list)))
            (iter list))))
      (assert (eq? (search '(1 2 3 4)) 3) "call/cc in iteration, call from outside")

      (define continued false)
      (define cont false)
      (define result (+ (call/cc
                          (lambda (return)
                            (set! cont return)
                            1))
                        1))
      (assert (eq? result
                   (if continued 16 2))
              "call/cc with saved reference")
      (cond ((not continued)
             (set! continued true)
             (cont 15)))

      'ok)))

(define js (compile-js exp))
(out js)
