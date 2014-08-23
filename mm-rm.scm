(define (label? exp) (symbol? exp))

(define (make-instruction-sequence needs modifies statements)
  (list needs modifies statements))

(define (registers-needed s) (if (label? s) '() (car s)))

(define (registers-modified s) (if (label? s) '() (cadr s)))

(define (needs-register? seq reg)
  (memq reg (registers-needed seq)))

(define (modifies-register? seq reg)
  (memq reg (registers-modified seq)))

(define (statements s)
  (if (label? s) (list s) (caddr s)))

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
                                   (statements seq1)
                                   (list (list 'restore (car regs)))))
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
    (cats (symbol-name name)
          (number->string (new-label-number)))))

(define all-regs '(env proc val args cont))

(define (compile-self-evaluating exp target linkage)
  (end-with-linkage
    linkage
    (make-instruction-sequence
      '()
      (list target)
      (list (list 'assign target (list 'const exp))))))

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
  (let ((after-if (make-label 'after-if))
        (true-branch (make-label 'true-branch))
        (false-branch (make-label 'false-branch)))
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
  (let ((after-lambda (make-label 'after-lambda))
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

(define (compile-proc-appl target linkage)
  (cond ((and (eq? target 'val) (not (eq? linkage 'return)))
         (make-instruction-sequence
           '(proc) all-regs
           (list (list 'assign
                       'cont
                       (list 'label linkage))
                 '(assign val
                          (op compiled-procedure-entry)
                          (reg proc))
                 '(goto (reg val)))))
        ((and (not (eq? target 'val))
              (not (eq? linkage 'return)))
         (let ((proc-return (make-label 'proc-return)))
           (make-instruction-sequence
             '(proc) all-regs
             (list (list 'assign
                         'cont
                         (list 'label proc-return))
                   '(assign val
                            (op compiled-procedure-entry)
                            (reg proc))
                   '(goto (reg val))
                   proc-return
                   (list target '(reg val))
                   (list 'goto (list 'label linkage))))))
        ((and (eq? target 'val) (eq? linkage 'return))
         (make-instruction-sequence
           '(proc cont)
           all-regs
           '((assign val
                     (op compiled-procedure-entry)
                     (reg proc))
             (goto (reg val)))))
        (else
          (error 'compile-proc-appl
                 "return linkage, target not val"
                 target))))

(define (compile-procedure-call target linkage)
  (let ((after-call (make-label 'after-call))
        (primitive-branch (make-label 'primitive-branch))
        (compiled-branch (make-label 'compiled-branch)))
    (let ((compiled-linkage
            (if (eq? linkage 'next) after-call linkage)))
      (append-instruction-sequences
        (make-instruction-sequence
          '(proc) '()
          (list '(test (op primitive-procedure?) (reg proc))
                (list 'branch (list 'label primitive-branch))))
        (parallel-instruction-sequences
          (append-instruction-sequences
            compiled-branch
            (compile-proc-appl target compiled-linkage))
          (append-instruction-sequences
            primitive-branch
            (end-with-linkage
              linkage
              (make-instruction-sequence
                '(proc args) (list target)
                (list (list 'assign
                            target
                            '(op apply-primitive-procedure)
                            '(reg proc)
                            '(reg args)))))))
        after-call))))

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
        ((quoted? exp) (compile-quoted exp target linkage))
        ((variable? exp)
         (compile-variable exp target linkage))
        ((assignment? exp)
         (compile-assignment exp target linkage))
        ((definition? exp)
         (compile-definition exp target linkage))
        ((if? exp) (compile-if exp target linkage))
        ((lambda? exp) (compile-lambda exp target linkage))
        ((begin? exp)
         (compile-sequence (begin-actions exp)
                           target
                           linkage))
        ((application? exp)
         (compile-application exp target linkage))
        (else (error 'compile-asm "invalid expression type" exp))))

(define (optimize-asm asm) asm)

(define (assemble asm) asm)

(define (compile-js exp)
  (assemble (optimize-asm (compile exp 'val 'next))))

(define exp '(define (a x)
               (define (b y) (* 2 y))
               (b (lambda () 1))
               (if 1 2 3)
               (set! x 3)
               (begin (a 1) (b "some") (a 'some))))
(define js (compile-js exp))
(print js)
