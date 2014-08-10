(define (assoc v l)
  (cond ((null? l) false)
        ((eq? v (caar l)) (car l))
        (else (assoc v (cdr l)))))

(define (for-each p l)
  (cond ((null? l) 'ok)
        (else
          (p (car l))
          (for-each p (cdr l)))))

(define (make-register name)
  (let ((contents 'just-defined))
    (define (dispatch message)
      (cond ((eq? message 'get) contents)
            ((eq? message 'set)
             (lambda (value) (set! contents value)))
            (else (error 'make-register "unknown request" message))))
    dispatch))

(define (get-contents register) (register 'get))

(define (set-contents! register value) ((register 'set) value))

(define (make-stack)
  (let ((s '()))
    (define (push x)
      (set! s (cons x s)))
    (define (pop)
      (if (null? s)
        (error 'pop "empty stack" s)
        (let ((top (car s)))
          (set! s (cdr s))
          top)))
    (define (initialize)
      (set! s '())
      'done)
    (define (dispatch message)
      (cond ((eq? message 'push) push)
            ((eq? message 'pop) (pop))
            ((eq? message 'initialize) (initialize))
            (else (error 'make-stack "unknown request" message))))
    dispatch))

(define (pop stack) (stack 'pop))

(define (push stack value) ((stack 'push) value))

(define (make-new-machine)
  (let ((pc (make-register 'pc))
        (flag (make-register 'flag))
        (stack (make-stack))
        (the-instruction-sequence '()))
    (let ((the-ops
            (list (list 'initialize-stack
                        (lambda () (stack 'initialize)))))
          (register-table
            (list (list 'pc pc) (list 'flag flag))))
      (define (allocate-register name)
        (if (assoc name register-table)
          (error 'allocate-register "register redefined" name)
          (set! register-table
            (cons (list name (make-register name))
                  register-table)))
        'register-allocated)
      (define (lookup-register name)
        (let ((val (assoc name register-table)))
          (if val
            (cadr val)
            (error 'lookup-register "unknown register" name))))
      (define (execute)
        (let ((insts (get-contents pc)))
          (if (null? insts)
            'done
            (begin
              ((instruction-execution-proc (car insts)))
              (execute)))))
      (define (dispatch message)
        (cond ((eq? message 'start)
               (set-contents! pc the-instruction-sequence)
               (execute))
              ((eq? message 'install-instruction-sequence)
               (lambda (seq) (set! the-instruction-sequence seq)))
              ((eq? message 'allocate-register) allocate-register)
              ((eq? message 'get-register) lookup-register)
              ((eq? message 'install-operations)
               (lambda (ops) (set! the-ops (append the-ops ops))))
              ((eq? message 'stack) stack)
              ((eq? message 'operations) the-ops)
              (else (error 'make-new-machine "unknown request" message))))
      dispatch)))

(define (start machine) (machine 'start))

(define (get-register machine reg-name)
  ((machine 'get-register) reg-name))

(define (get-register-contents machine register-name)
  (get-contents (get-register machine register-name)))

(define (set-register-contents machine register-name value)
  (set-contents! (get-register machine register-name) value)
  'done)

(define (instruction-text inst)
  (table-lookup inst 'text))

(define (instruction-execution-proc inst)
  (table-lookup inst 'proc))

(define (set-instruction-execution-proc! inst proc)
  (table-define inst 'proc proc))

(define (assign-reg-name inst) (cadr inst))

(define (assign-value-exp inst) (cddr inst))

(define (advance-pc pc)
  (set-contents! pc (cdr (get-contents pc))))

(define (register-exp? exp) (tagged-list? exp 'reg))

(define (register-exp-reg exp) (cadr exp))

(define (constant-exp? exp) (tagged-list? exp 'const))

(define (constant-exp-value exp) (cadr exp))

(define (label-exp? exp) (tagged-list? exp 'label))

(define (label-exp-label exp) (cadr exp))

(define (make-primitive-exp exp machine labels)
  (cond ((constant-exp? exp)
         (let ((c (constant-exp-value exp)))
           (lambda () c)))
        ((label-exp? exp)
         (let ((insts
                 (lookup-label labels
                               (label-exp-label exp))))
           (lambda () insts)))
        ((register-exp? exp)
         (let ((r (get-register machine
                                (register-exp-reg exp))))
           (lambda () (get-contents r))))
        (else (error 'make-primitive-exp "bad expression" exp))))

(define (operation-exp? exp)
  (and (pair? exp) (tagged-list? (car exp) 'op)))

(define (operation-exp-op operation-exp)
  (cadar operation-exp))

(define (operation-exp-operands operation-exp)
  (cdr operation-exp))

(define (lookup-prim symbol operations)
  (let ((val (assoc symbol operations)))
    (if val
      (cadr val)
      (error 'lookup-prim "unknown operation" symbol))))

(define (make-operation-exp exp machine labels operations)
  (let ((op (lookup-prim (operation-exp-op exp) operations))
        (aprocs
          (map (lambda (e)
                 (make-primitive-exp e machine labels))
               (operation-exp-operands exp))))
    (lambda ()
      (apply op (map (lambda (p) (p)) aprocs)))))

(define (make-assign inst machine labels ops pc)
  (let ((target
          (get-register machine (assign-reg-name inst)))
        (value-exp (assign-value-exp inst)))
    (let ((value-proc
            (if (operation-exp? value-exp)
              (make-operation-exp
                value-exp machine labels ops)
              (make-primitive-exp
                (car value-exp) machine labels))))
      (lambda ()
        (set-contents! target (value-proc))
        (advance-pc pc)))))

(define (test-condition test-instruction)
  (cdr test-instruction))

(define (make-test inst machine labels operations flag pc)
  (let ((condition (test-condition inst)))
    (if (operation-exp? condition)
      (let ((condition-proc
              (make-operation-exp
                condition machine labels operations)))
        (lambda ()
          (set-contents! flag (condition-proc))
          (advance-pc pc)))
      (error 'make-test "bad test" inst))))

(define (branch-dest branch-instruction)
  (cadr branch-instruction))

(define (make-branch inst machine labels flag pc)
  (let ((dest (branch-dest inst)))
    (if (label-exp? dest)
      (let ((insts
              (lookup-label labels (label-exp-label dest))))
        (lambda ()
          (if (get-contents flag)
            (set-contents! pc insts)
            (advance-pc pc))))
      (error 'make-branch "bad branch" inst))))

(define (goto-dest goto-instruction)
  (cadr goto-instruction))

(define (make-goto inst machine labels pc)
  (let ((dest (goto-dest inst)))
    (cond ((label-exp? dest)
           (let ((insts
                   (lookup-label labels
                                 (label-exp-label dest))))
             (lambda () (set-contents! pc insts))))
          ((register-exp? dest)
           (let ((reg
                   (get-register machine
                                 (register-exp-reg dest))))
             (lambda ()
               (set-contents! pc (get-contents reg)))))
          (else (error 'make-goto "bad goto" inst)))))

(define (stack-inst-reg-name stack-instruction)
  (cadr stack-instruction))

(define (make-save inst machine stack pc)
  (let ((reg (get-register machine
                           (stack-inst-reg-name inst))))
    (lambda ()
      (push stack (get-contents reg))
      (advance-pc pc))))

(define (make-restore inst machine stack pc)
  (let ((reg (get-register machine
                           (stack-inst-reg-name inst))))
    (lambda ()
      (set-contents! reg (pop stack))
      (advance-pc pc))))

(define (perform-action inst) (cdr inst))

(define (make-perform inst machine labels operations pc)
  (let ((action (perform-action inst)))
    (if (operation-exp? action)
      (let ((action-proc
              (make-operation-exp
                action machine labels operations)))
        (lambda ()
          (action-proc)
          (advance-pc pc)))
      (error 'make-perform "bad perform" inst))))

(define (make-execution-procedure inst labels machine
                                  pc flag stack ops)
  (cond ((eq? (car inst) 'assign)
         (make-assign inst machine labels ops pc))
        ((eq? (car inst) 'test)
         (make-test inst machine labels ops flag pc))
        ((eq? (car inst) 'branch)
         (make-branch inst machine labels flag pc))
        ((eq? (car inst) 'goto)
         (make-goto inst machine labels pc))
        ((eq? (car inst) 'save)
         (make-save inst machine stack pc))
        ((eq? (car inst) 'restore)
         (make-restore inst machine stack pc))
        ((eq? (car inst) 'perform)
         (make-perform inst machine labels ops pc))
        (else (error 'make-execution-procedure
                     "invalid instruction"
                     (car inst)))))

(define (update-insts! insts labels machine)
  (let ((pc (get-register machine 'pc))
        (flag (get-register machine 'flag))
        (stack (machine 'stack))
        (ops (machine 'operations)))
    (for-each
      (lambda (inst)
        (set-instruction-execution-proc!
          inst
          (make-execution-procedure
            (instruction-text inst) labels machine
            pc flag stack ops)))
      insts)))

(define (make-instruction text)
  (let ((inst (make-name-table)))
    (table-define inst 'text text)
    inst))

(define (make-label-entry name insts) (cons name insts))

(define (extract-labels text receive)
  (if (null? text)
    (receive '() '())
    (extract-labels
      (cdr text)
      (lambda (insts labels)
        (let ((next-inst (car text)))
          (if (symbol? next-inst)
            (receive insts
                     (cons (make-label-entry next-inst insts)
                           labels))
            (receive (cons (make-instruction next-inst)
                           insts)
                     labels)))))))

(define (assemble controller-text machine)
  (extract-labels
    controller-text
    (lambda (insts labels)
      (update-insts! insts labels machine)
      insts)))

(define (lookup-label labels name)
  (let ((val (assoc name labels)))
    (if val
      (cdr val)
      (error 'lookup-label "undefined label" name))))

(define (make-machine register-names ops controller-text)
  (let ((machine (make-new-machine)))
    (for-each (lambda (name)
                ((machine 'allocate-register) name))
              register-names)
    ((machine 'install-operations) ops)
    ((machine 'install-instruction-sequence)
     (assemble controller-text machine))
    machine))

; test

; ((make-machine
;   '(a)
;   (list (list 'log log))
;   '(start
;      (goto (label here))
;      here
;      (assign a (const 3))
;      (goto (label there))
;      there
;      (perform (op log) (reg a))))
;  'start)

; ((make-machine
;    '(a b t)
;    (list (list 'log log)
;          (list '% %)
;          (list 'eq? eq?))
;    '((assign a (const 18))
;      (assign b (const 12))
;      test-b
;      (test (op eq?) (reg b) (const 0))
;      (branch (label gcd-done))
;      (assign t (op %) (reg a) (reg b))
;      (assign a (reg b))
;      (assign b (reg t))
;      (goto (label test-b))
;      gcd-done
;      (perform (op log) (reg a))))
;  'start)

(define gc
  '(begin-garbage-collection
     (assign free (const 0))
     (assign scan (const 0))
     (assign old (reg root))
     (assign relocate-continue (label reassign-root))
     (goto (label relocate-old-result-in-new))

     reassign-root
     (assign root (reg new))
     (goto (label gc-loop))

     gc-loop
     (test (op =) (reg scan) (reg free))
     (branch (label gc-flip))
     (assign old (op vector-ref) (reg new-cars) (reg scan))
     (assign relocate-continue (label update-car))
     (goto (label relocate-old-result-in-new))

     update-car
     (perform
       (op vector-set!) (reg new-cars) (reg scan) (reg new))
     (assign old (op vector-ref) (reg new-cdrs) (reg scan))
     (assign relocate-continue (label update-cdr))
     (goto (label relocate-old-result-in-new))

     update-cdr
     (perform
       (op vector-set!) (reg new-cdrs) (reg scan) (reg new))
     (assign scan (op +) (reg scan) (const 1))
     (goto (label gc-loop))

     relocate-old-result-in-new
     (test (op pair?) (reg old))
     (branch (label pair))
     (assign new (reg old))
     (goto (reg relocate-continue))

     pair
     (assign oldcr (op vector-ref) (reg the-cars) (reg old))
     (test (op =) (reg oldcr) (const broken-heart))
     (branch (label already-moved))
     (assign new (reg free))
     (assign free (op +) (reg free) (const 1))
     (perform (op vector-set!)
              (reg new-cars) (reg new) (reg oldcr))
     (assign oldcr (op vector-ref) (reg the-cdrs) (reg old))
     (perform (op vector-set!)
              (reg new-cdrs) (reg new) (reg oldcr))
     (perform (op vector-set!)
              (reg the-cars) (reg old) (const broken-heart))
     (perform
       (op vector-set!) (reg the-cdrs) (reg old) (reg new))
     (goto (reg relocate-continue))

     already-moved
     (assign new (op vector-ref) (reg the-cdrs) (reg old))
     (goto (reg relocate-continue))

     gc-flip
     (assign temp (reg the-cdrs))
     (assign the-cdrs (reg new-cdrs))
     (assign new-cdrs (reg temp))
     (assign temp (reg the-cars))
     (assign the-cars (reg new-cars))
     (assign new-cars (reg temp))))

(define (make-mvector size)
  (struct (list 'size size)))

(define (mvector-ref v ref)
  (if (>= ref (table-lookup v 'size))
    (error 'mvector-ref "ref out of bounds" ref)
    (if (table-has-name? v (string->symbol (number->string ref)))
      (table-lookup v (string->symbol (number->string ref)))
      0)))

(define (mvector-set! v ref value)
  (if (>= ref (table-lookup v 'size))
    (error 'mvector-set! "ref out of bounds" ref)
    (table-define v (string-symbol (number->string ref)) value)))

; ((make-machine
;     '(free scan old new root
;            oldcr temp
;            relocate-continue
;            new-cars new-cdrs
;            the-cars the-cdrs)
;     (list (list 'vector-ref mvector-ref)
;           (list 'vector-set! mvector-set!)
;           (list '= eq?)
;           (list '+ +)
;           (list 'pair? pair?))
;     gc)
;  'start)

(define ev
  '((goto (label read-eval-print-loop))
    eval-dispatch
    (test (op self-evaluating?) (reg exp))
    (branch (label ev-self-eval))
    (test (op variable?) (reg exp))
    (branch (label ev-variable))
    (test (op quoted?) (reg exp))
    (branch (label ev-quoted))
    (test (op assignment?) (reg exp))
    (branch (label ev-assignment))
    (test (op definition?) (reg exp))
    (branch (label ev-definition))
    (test (op if?) (reg exp))
    (branch (label ev-if))
    (test (op lambda?) (reg exp))
    (branch (label ev-lambda))
    (test (op begin?) (reg exp))
    (branch (label ev-begin))
    (test (op application?) (reg exp))
    (branch (label ev-application))
    (goto (label unknown-expression-type))

    ev-self-eval
    (assign val (reg exp))
    (goto (reg cont))

    ev-variable
    (assign val (op lookup-variable-value) (reg exp) (reg env))
    (goto (reg cont))

    ev-quoted
    (assign val (op text-of-quotation) (reg exp))
    (goto (reg cont))

    ev-lambda
    (assign unev (op lambda-parameters) (reg exp))
    (assign exp (op lambda-body) (reg exp))
    (assign val
            (op make-procedure)
            (reg unev) (reg exp) (reg env))
    (goto (reg cont))

    ev-application
    (save cont)
    (save env)
    (assign unev (op operands) (reg exp))
    (save unev)
    (assign exp (op operator) (reg exp))
    (assign cont (label ev-apply-did-operator))
    (goto (label eval-dispatch))

    ev-apply-did-operator
    (restore unev)
    (restore env)
    (assign args (op empty-arglist))
    (assign proc (reg val))
    (test (op no-operands?) (reg unev))
    (branch (label apply-dispatch))
    (save proc)

    ev-appl-operand-loop
    (save args)
    (assign exp (op first-operand) (reg unev))
    (test (op last-operand?) (reg unev))
    (branch (label ev-appl-last-arg))
    (save env)
    (save unev)
    (assign cont (label ev-appl-accumulate-arg))
    (goto (label eval-dispatch))

    ev-appl-accumulate-arg
    (restore unev)
    (restore env)
    (restore args)
    (assign args (op adjoin-arg) (reg val) (reg args))
    (assign unev (op rest-operands) (reg unev))
    (goto (label ev-appl-operand-loop))

    ev-appl-last-arg
    (assign cont (label ev-appl-accum-last-arg))
    (goto (label eval-dispatch))

    ev-appl-accum-last-arg
    (restore args)
    (assign args (op adjoin-arg) (reg val) (reg args))
    (restore proc)
    (goto (label apply-dispatch))

    apply-dispatch
    (test (op primitive-procedure?) (reg proc))
    (branch (label primitive-apply))
    (test (op compound-procedure?) (reg proc))
    (branch (label compound-apply))
    (goto (label unknown-procedure-type))

    primitive-apply
    (assign val
            (op apply-primitive-procedure)
            (reg proc)
            (reg args))
    (restore cont)
    (goto (reg cont))

    compound-apply
    (assign unev (op procedure-parameters) (reg proc))
    (assign env (op procedure-environment) (reg proc))
    (assign env
            (op extend-environment)
            (reg unev) (reg args) (reg env))
    (assign unev (op procedure-body) (reg proc))
    (goto (label ev-sequence))

    ev-begin
    (assign unev (op begin-actions) (reg exp))
    (save cont)
    (goto (label ev-sequence))

    ev-sequence
    (assign exp (op first-exp) (reg unev))
    (test (op last-exp?) (reg unev))
    (branch (label ev-sequence-last-exp))
    (save unev)
    (save env)
    (assign cont (label ev-sequence-continue))
    (goto (label eval-dispatch))

    ev-sequence-continue
    (restore env)
    (restore unev)
    (assign unev (op rest-exps) (reg unev))
    (goto (label ev-sequence))

    ev-sequence-last-exp
    (restore cont)
    (goto (label eval-dispatch))

    ev-if
    (save exp)
    (save env)
    (save cont)
    (assign cont (label ev-if-decide))
    (assign exp (op if-predicate) (reg exp))
    (goto (label eval-dispatch))

    ev-if-decide
    (restore cont)
    (restore env)
    (restore exp)
    (test (op true?) (reg val))
    (branch (label ev-if-consequent))

    ev-if-alternative
    (assign exp (op if-alternative) (reg exp))
    (goto (label eval-dispatch))

    ev-if-consequent
    (assign exp (op if-consequent) (reg exp))
    (goto (label eval-dispatch))

    ev-assignment
    (assign unev (op assignment-variable) (reg exp))
    (save unev)
    (assign exp (op assignment-value) (reg exp))
    (save env)
    (save cont)
    (assign cont (label ev-assignment-1))
    (goto (label eval-dispatch))

    ev-assignment-1
    (restore cont)
    (restore env)
    (restore unev)
    (perform
      (op set-variable-value!) (reg unev) (reg val) (reg env))
    (assign val (const ok))
    (goto (reg cont))

    ev-definition
    (assign unev (op definition-variable) (reg exp))
    (save unev)
    (assign exp (op definition-value) (reg exp))
    (save env)
    (save cont)
    (assign cont (label ev-definition-1))
    (goto (label eval-dispatch))

    ev-definition-1
    (restore cont)
    (restore env)
    (restore unev)
    (perform
      (op define-variable!) (reg unev) (reg val) (reg env))
    (assign val (const ok))
    (goto (reg cont))

    read-eval-print-loop
    (perform (op initialize-stack))
    (assign exp (op read))
    (assign env (op get-global-environment))
    (assign cont (label print-result))
    (goto (label eval-dispatch))

    print-result
    (perform (op user-print) (reg val))
    (goto (label exit))

    unknown-expression-type
    (assign val (const unknown-expression-type-error))
    (goto (label signal-error))

    unknown-procedure-type
    (restore cont)
    (assign val (const unknown-procedure-type-error))
    (goto (label signal-error))

    signal-error
    (perform (op user-print) (reg val))
    (goto (label exit))

    exit))

(define (quoted? exp) (tagged-list? exp 'quote))

(define (variable? exp) (symbol? exp))

(define (assignment? exp) (tagged-list? exp 'set!))

(define (definition? exp) (tagged-list? exp 'define))

(define (if? exp) (tagged-list? exp 'if))

(define (lambda? exp) (tagged-list? exp 'lambda))

(define (begin? exp) (tagged-list? exp 'begin))

(define (application? exp) (pair? exp))

(define (last-exp? exps) (null? (cdr exps)))

(define (first-exp exps) (car exps))

(define (rest-exps exps) (cdr exps))

(define (false? exp) (eq? exp false))

(define (define-variable! variable value env) (env variable true value))

(define (text-of-quotation exp) (cadr exp))

(define ev-ops (list
                 (list 'adjoin-arg (lambda (arg arg-list) (append arg-list (list arg))))

                 (list 'application? application?)

                 (list 'apply-primitive-procedure (lambda (proc args) (apply proc args)))

                 (list 'assignment-value (lambda (exp) (caddr exp)))

                 (list 'assignment-variable (lambda (exp) (cadr exp)))

                 (list 'assignment? assignment?)

                 (list 'begin-actions (lambda (exp) (cdr exp)))

                 (list 'begin? begin?)

                 (list 'compound-procedure? (lambda (exp) (tagged-list? exp 'compound)))

                 (list 'define-variable! define-variable!)

                 (list 'definition-value (lambda (exp)
                                           (if (pair? (cadr exp))
                                             (cons 'lambda (cons (cdadr exp) (cddr exp)))
                                             (caddr exp))))

                 (list 'definition-variable (lambda (exp)
                                              (if (pair? (cadr exp))
                                                (caadr exp)
                                                (cadr exp))))

                 (list 'definition? definition?)

                 (list 'empty-arglist (lambda _ _))

                 (list 'extend-environment (lambda (param-list arg-list env)
                                             (define (iter param-list arg-list env)
                                               (cond ((null? param-list) env)
                                                     (else (env (car param-list) true (car arg-list))
                                                           (iter (cdr param-list) (cdr arg-list) env))))
                                             (let ((env (extend-env env false)))
                                               (iter param-list arg-list env))))

                 (list 'first-exp first-exp)

                 (list 'first-operand (lambda (ops) (car ops)))

                 (list 'get-global-environment (lambda () mikkamakka))

                 (list 'if-alternative (lambda (exp) (cadddr exp)))

                 (list 'if-consequent (lambda (exp) (caddr exp)))

                 (list 'if-predicate (lambda (exp) (cadr exp)))

                 (list 'if? if?)

                 (list 'initialize-stack (lambda x x))

                 (list 'lambda-body (lambda (exp)
                                      (cddr exp)))

                 (list 'lambda-parameters (lambda (exp) (cadr exp)))

                 (list 'lambda? lambda?)

                 (list 'last-exp? last-exp?)

                 (list 'last-operand? (lambda (exps) (null? (cdr exps))))

                 (list 'lookup-variable-value (lambda (var env) (env var)))

                 (list 'make-procedure (lambda (parameters body env)
                                         (list 'compound parameters body env)))

                 (list 'no-operands? (lambda (args) (null? args)))

                 (list 'operands (lambda (exp) (cdr exp)))

                 (list 'operator (lambda (exp) (car exp)))

                 (list 'primitive-procedure? (lambda (exp) (primitive-procedure? exp)))

                 (list 'procedure-body (lambda (exp)
                                         (caddr exp)))

                 (list 'procedure-environment (lambda (exp) (cadddr exp)))

                 (list 'procedure-parameters (lambda (exp) (cadr exp)))

                 (list 'quoted? quoted?)

                 (list 'rest-exps rest-exps)

                 (list 'rest-operands (lambda (exps) (cdr exps)))

                 (list 'self-evaluating? self-evaluating?)

                 (list 'set-variable-value! (lambda (var value env)
                                              (env var false value)))

                 (list 'text-of-quotation text-of-quotation)

                 (list 'true? (lambda (exp) (not (eq? exp false))))

                 (list 'user-print (lambda (exp) (print exp)))

                 (list 'variable? variable?)

                 (list 'read (lambda () '((lambda ()
                                            (define (gcd a b)
                                              (if (eq? 0 b)
                                                a
                                                (gcd b (% a b))))
                                            (gcd 36 90)))))))

; ((make-machine
;    '(exp env val proc args cont unev)
;    ev-ops
;    ev) 'start)

(define (make-instruction-sequence needs modifies statements)
  (list needs modifies statements))

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

; only return needs preserving?
(define (end-with-linkage linkage instruction-sequence)
  (preserving '(cont)
              instruction-sequence
              (compile-linkage linkage)))

(define (compile-self-evaluating exp target linkage)
  (end-with-linkage
    linkage
    (make-instruction-sequence
      '()
      (list target)
      (list (list 'assign target (list 'const exp))))))

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

(define (compile-assignment exp target linkage)
  (let ((var (assignment-variable exp))
        (get-value-code
          (compile (assignment-value exp) 'val 'next)))
    (end-with-linkage
      linkage
      (preserving '(env)
                  get-value-code
                  (make-instruction-sequence
                    '(env val)
                    (list target)
                    (list (list 'perform
                                '(op set-variable-value!)
                                (list 'const var)
                                '(reg val)
                                '(reg env))
                          (list 'assign
                                target
                                '(const ok))))))))

(define (compile-definition exp target linkage)
  (let ((var (definition-variable exp))
        (get-value-code
          (compile (definition-value exp) 'val 'next)))
    (end-with-linkage
      linkage
      (preserving '(env)
                  get-value-code
                  (make-instruction-sequence
                    '(env val)
                    (list target)
                    (list (list 'perform
                                '(op define-variable!)
                                (list 'const var)
                                '(reg val)
                                '(reg env))
                          (list 'assign
                                target
                                '(const ok))))))))

(define label-counter 0)

(define (new-label-number)
  (set! label-counter (+ label-counter 1))
  label-counter)

(define (make-label name)
  (string->symbol
    (cats (symbol-name name)
          (number->string (new-label-number)))))

(define (registers-needed s)
  (if (symbol? s) '() (car s)))

(define (registers-modified s)
  (if (symbol? s) '() (cadr s)))

(define (statements s)
  (if (symbol? s) (list s) (caddr s)))

(define (needs-register? seq reg)
  (memq reg (registers-needed seq)))

(define (modifies-register? seq reg)
  (memq reg (registers-modified seq)))

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

(define (preserving regs seq1 seq2)
  (if (null? regs)
    (append-instruction-sequences seq1 seq2)
    (let ((first-reg (car regs)))
      (if (and (needs-register? seq2 first-reg)
               (modifies-register? seq1 first-reg))
        (preserving (cdr regs)
                    (make-instruction-sequence
                      (list-union (list first-reg)
                                  (registers-needed seq1))
                      (list-difference (registers-modified seq1)
                                       (list first-reg))
                      (append (list (list 'save first-reg))
                              (statements seq1)
                              (list (list 'restore first-reg))))
                    seq2)
        (preserving (cdr regs) seq1 seq2)))))

(define (tack-on-instruction-sequence seq body-seq)
  (make-instruction-sequence
    (registers-needed seq)
    (registers-modified seq)
    (append (statements seq) (statements body-seq))))

(define (parallel-instruction-sequences seq1 seq2)
  (make-instruction-sequence
    (list-union (registers-needed seq1)
                (registers-needed seq2))
    (list-union (registers-modified seq1)
                (registers-modified seq2))
    (append (statements seq1) (statements seq2))))

(define (compile-if exp target linkage)
  (let ((t-branch (make-label 'true-branch))
        (f-branch (make-label 'false-branch))
        (after-if (make-label 'after-if)))
    (let ((consequent-linkage
            (if (eq? linkage 'next) after-if linkage)))
      (let ((p-code (compile (if-predicate exp) 'val 'next))
            (c-code (compile (if-consequent exp) target consequent-linkage))
            (a-code (compile (if-alternative exp) target linkage)))
        (preserving '(env cont)
                    p-code
                    (append-instruction-sequences
                      (make-instruction-sequence
                        '(val) '()
                        (list '(test (op false?) (reg val))
                              (list 'branch (list 'label f-branch))))
                      (parallel-instruction-sequences
                        (append-instruction-sequences t-branch c-code)
                        (append-instruction-sequences f-branch a-code))
                      after-if))))))

(define (compile-sequence seq target linkage)
  (if (last-exp? seq)
    (compile (first-exp seq) target linkage)
    (preserving '(env cont)
                (compile (first-exp seq) target 'next)
                (compile-sequence (rest-exps seq) target linkage))))

(define (make-compiled-procedure entry env)
  (list 'compiled-procedure entry env))

(define primitive-procedure? compiled-procedure?)

(define (compiled-procedure? proc)
  (tagged-list? proc 'compiled-procedure))

(define (compiled-procedure-entry c-proc) (cadr c-proc))

(define (compiled-procedure-env c-proc) (caddr c-proc))

(define (compile-lambda-body exp proc-entry)
  (let ((formals (lambda-parameters exp)))
    (append-instruction-sequences
      (make-instruction-sequence
        '(env proc args) '(env)
        (list proc-entry
              '(assign env (op compiled-procedure-env) (reg proc))
              (list 'assign
                    'env
                    '(op extend-environment)
                    (list 'const formals)
                    '(reg args)
                    '(reg env))))
      (compile-sequence (lambda-body exp) 'val 'return))))

(define (compile-lambda exp target linkage)
  (let ((proc-entry (make-label 'entry))
        (after-lambda (make-label 'after-lambda)))
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
                          (list 'label proc-entry)
                          '(reg env)))))
          (compile-lambda-body exp proc-entry))
        after-lambda))))

(define (code-to-get-rest-args operand-codes)
  (let ((code-for-next-arg
          (preserving
            '(args)
            (car operand-codes)
            (make-instruction-sequence
              '(val args) '(args)
              '((assign args
                        (op cons)
                        (reg val)
                        (reg args)))))))
    (if (null? (cdr operand-codes))
      code-for-next-arg
      (preserving
        '(env)
        code-for-next-arg
        (code-to-get-rest-args (cdr operand-codes))))))

(define (construct-arglist operand-codes)
  (let ((operand-codes (reverse operand-codes)))
    (if (null? operand-codes)
      (make-instruction-sequence
        '() '(args)
        '((assign args (const ()))))
      (let ((code-to-get-last-arg
              (append-instruction-sequences
                (car operand-codes)
                (make-instruction-sequence
                  '(val) '(args)
                  '((assign args (op list) (reg val)))))))
        (if (null? (cdr operand-codes))
          code-to-get-last-arg
          (preserving
            '(env)
            code-to-get-last-arg
            (code-to-get-rest-args
              (cdr operand-codes))))))))

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
        ((and (not (eq? target 'val)) (eq? linkage 'return))
         (error 'compile-proc-appl
                "return linkage, target not val"
                target))))

(define (compile-procedure-call target linkage)
  (let ((primitive-branch (make-label 'primitive-branch))
        (compiled-branch (make-label 'compiled-branch))
        (after-call (make-label 'after-call)))
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

(define all-regs '(env proc val args cont))

(define (compile-application exp target linkage)
  (let ((proc-code (compile (operator exp) 'proc 'next))
        (operand-codes
          (map (lambda (operand) (compile operand 'val 'next))
               (operands exp))))
    (preserving
      '(env cont)
      proc-code
      (preserving
        '(proc cont)
        (construct-arglist operand-codes)
        (compile-procedure-call target linkage)))))

(define (self-evaluating? exp)
  (or (number? exp)
      (string? exp)
      (eq? exp true)
      (eq? exp false)
      (eq? exp no-print)))

(define (js-import-code? exp)
  (tagged-list? exp 'js-import-code))

(define (compile-js-import-code exp target linkage)
  (make-instruction-sequence
    '(env)
    '()
    (list (list 'js-import-code
                (list 'const (cadr exp))
                (list 'const (caddr exp))))))

(define (compile exp target linkage)
  (cond ((self-evaluating? exp)
         (compile-self-evaluating exp target linkage))
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
        ((js-import-code? exp)
         (compile-js-import-code exp target linkage))
        ((application? exp)
         (compile-application exp target linkage))
        (else
          (error 'compile "Unknown expression type" exp))))

; (print (compile
;          '(define (factorial n)
;             (if (eq? n 1)
;               1
;               (* (factorial (- n 1)) n)))
;          'val
;          'next))

(define compile-ops
  (append ev-ops
          (list (list 'make-compiled-procedure make-compiled-procedure)
                (list 'compiled-procedure-env compiled-procedure-env)
                (list 'compiled-procedure-entry compiled-procedure-entry)
                (list 'list list)
                (list 'cons cons)
                (list 'false? false?))))

(define js-head "(function () {
    \"use strict\";

    // ops
    var error = function (msg) {
        throw msg;
    };

    var func = function (argLength, allowMore, customCheck, f) {
        return function () {
            var args = Array.prototype.slice.call(arguments);

            if (allowMore && args.length < argLength ||
                !allowMore && args.length !== argLength) {
                return error(\"invalid arity\");
            }

            if (customCheck && !customCheck.apply(undefined, args)) {
                return error(\"argument error\");
            }

            return f.apply(undefined, args);
        };
    };

    var isFalse = func(1, false, false, function (val) {
        return val === false;
    });

    var isNull = func(1, false, false, function (val) {
        return val instanceof Array && !val.length;
    });

    var cons = func(2, false, false, function (left, right) {
        return [left, right];
    });

    var car = func(1, false, function (p) {
        return isPair(p);
    }, function (p) {
        return p[0];
    });

    var cdr = func(1, false, function (p) {
        return isPair(p);
    }, function (p) {
        return p[1];
    });

    var isPair = func(1, false, false, function (p) {
        return p instanceof Array &&
            p.length === 2;
    });

    var consAll = func(1, true, false, function () {
        var l = arguments[arguments.length - 1];
        for (var i = arguments.length - 2; i >= 0; i--) {
            l = [arguments[i], l];
        }
        return l;
    });

    var list = function () {
        var args = Array.prototype.slice.call(arguments);
        args.push([]);
        return consAll.apply(undefined, args);
    };

    var isSymbol = func(1, false, false, function (symbol) {
        return symbol instanceof Array && symbol.length === 1;
    });

    var symbolEq = func(2, false, function (left, right) {
        return isSymbol(left) && isSymbol(right);
    }, function (left, right) {
        return symbolToString(left) === symbolToString(right);
    });

    var symbolToString = func(1, false, function (symbol) {
        return isSymbol(symbol);
    }, function (symbol) {
        return symbol[0];
    });

    var stringToSymbol = func(1, false, function (string) {
        return isString(string);
    }, function (string) {
        return [string];
    });

    var isString = func(1, false, false, function (string) {
        return typeof string === \"string\";
    });

    var isNumber = func(1, false, false, function (number) {
        return typeof number === \"number\" && !Number.isNaN(number);
    });

    var neq = func(2, true, function () {
        for (var i = 0; i < arguments.length; i++) {
            if (!isNumber(arguments[i])) {
                return false;
            }
        }
        return true;
    }, function () {
        for (var i = 0; i < arguments.length - 1; i++) {
            if (arguments[i] !== arguments[i + 1]) {
                return false;
            }
        }
        return true;
    });

    var identity = func(1, false, false, function (x) { return x; });

    var listToPlist = func(1, false, false, function (l) {
        var plist = [];
        for (;!(isNull(l));) {
            if (!(isPair(l))) {
                return error(\"invalid argument\");
            }
            plist[plist.length] = car(l);
            l = cdr(l);
        }
        return plist;
    });

    var isCompoundProcedure = func(1, false, false, function (proc) {
        return false;
    });

    var isCompiledProcedure = func(1, false, false, function (proc) {
        return isPair(proc) &&
            symbolEq(car(proc), stringToSymbol(\"compiled\"));
    });

    var isPrimitiveProcedure = func(1, false, false, function (proc) {
        return proc instanceof Function;
    });

    var makeCompiledProcedure = func(2, false, false, function (entry, env) {
        return list(stringToSymbol(\"compiled\"), entry, env);
    });

    var compiledProcedureEnv = func(1, false, false, function (proc) {
        return car(cdr(cdr(proc)));
    });

    var compiledProcedureEntry = func(1, false, false, function (proc) {
        return car(cdr(proc));
    });

    var applyPrimitive = func(2, false, function (proc, _) {
        return isPrimitiveProcedure(proc);
    }, function (proc, args) {
        return proc.apply(undefined, listToPlist(args));
    });

    var isEnv = func(1, false, false, function (env) {
        return env && typeof env.current === \"object\";
    });

    var extendEnv = func(3, false, function (_, __, env) {
        return isEnv(env) || isFalse(env);
    }, function (names, values, env) {
        var ext = {parent: env, current: {}};

        for (; isPair(names);) {
            if (!isPair(values)) {
                throw \"not enough values\";
            }

            defineVar(car(names), car(values), ext);

            names = cdr(names);
            values = cdr(values);
        }

        if (!isNull(names)) {
            defineVar(names, values, ext);
            values = list();
        }

        if (!isNull(values)) {
            throw \"not enough names\";
        }

        return ext;
    });

    var defineVar = func(3, false, function (name, _, env) {
        return isEnv(env);
    }, function (name, val, env) {
        env.current[symbolToString(name)] = val;
        return val;
    });

    var findVar = func(3, false, function (_, env, f) {
        return isEnv(env) && typeof f === \"function\";
    }, function (name, env, f) {
        for (; env;) {
            if (symbolToString(name) in env.current) {
                return f(env);
            }
            env = env.parent;
        }
        throw \"unbound variable\";
    });

    var lookupVar = func(2, false, function (_, env) {
        return isEnv(env);
    }, function (name, env) {
        return findVar(name, env, function (env) {
            return env.current[symbolToString(name)];
        });
    });

    var setVar = func(3, false, function(_, __, env) {
        return isEnv(env);
    }, function (name, val, env) {
        return findVar(name, env, function (env) {
            env.current[symbolToString(name)] = val;
            return val;
        });
    });

    var mkNumOp = func(3, false, false, function (initial, single, reduce) {
        return func(0, true, function () {
            for (var i = 0; i < arguments.length; i++) {
                if (!isNumber(arguments[i])) {
                    return false;
                }
            }
            return true;
        }, function () {
            var args = Array.prototype.slice.call(arguments);
            if (!args.length) {
                return initial;
            }
            if (args.length === 1) {
                return single(args[0]);
            }
            return args.slice(1).reduce(function (previous, current) {
                return reduce(previous, current);
            }, args[0]);
        });
    });

    var checkImportType = function (val) {
        switch (true) {
        case typeof val === \"number\":
        case typeof val === \"string\":
            break;
        case val instanceof Array:
            for (var i = 0; i < val.length; i++) {
                checkImportType(val[i]);
            }
            break;
        default:
            return error(\"invalid argument type for import\");
        }
    };

    var importFunction = function (module, key) {
        return function () {
            var args = Array.prototype.slice.call(arguments);
            for (var i = 0; i < args.length; i++) {
                checkImportType(args[i]);
            }
            var res = module[key].apply(module, args);
            if (typeof res === \"undefined\") {
                res = stringToSymbol(\"unknown value\");
            } else {
                checkImportType(res);
            }
            return res;
        };
    };

    var importModule = function (name, module) {
        for (var key in module) {
            switch (typeof module[key]) {
            case \"number\":
            case \"string\":
                defineVar(stringToSymbol(key), module[key], env);
            case \"function\":
                defineVar(stringToSymbol(key), importFunction(module, key), env);
                break;
            default:
                return error(\"invalid import type\");
            }
        }
    };

    // registers
    var env = extendEnv(list(), list(), false);
    var val = stringToSymbol(\"initial\");
    var proc = function () { return proc };
    var args = list();
    var next = false;
    var cont = false;

    // stack
    var stack = list();

    var save = func(1, false, false, function (register) {
        stack = cons(register, stack);
        return register;
    });

    var restore = func(0, false, false, function () {
        var register = car(stack);
        stack = cdr(stack);
        return register;
    });

    // primitive procedures
    defineVar(stringToSymbol(\"*\"),
        mkNumOp(1, identity, function (x, y) { return x * y; }),
        env);

    defineVar(stringToSymbol(\"-\"),
        mkNumOp(0, function (x) {
            return 0 - x;
        }, function (x, y) {
            return x - y;
        }),
        env);

    defineVar(stringToSymbol(\"=\"), neq, env);

    // program
    next = function () {
")

(define js-tail "        return false;
    };

    // control
    for (; next; next = next());
})();")

(define (tab builder)
  (sbappend builder "    "))

(define (js-name symbol)
  (cats "_" (sreplace (symbol-name symbol)
                      "\\W"
                      (lambda (s _ __)
                        (cats "_"
                              (sprint (char-code-at s 0))
                              "_")))))

(define (assembly-const? expression)
  (tagged-list? (car expression) 'const))

(define (assemble-const-pair expression)
  (sbappend
    (sbappend
      (sbappend
        (sbappend
          (sbappend (make-string-builder) "cons(")
          (assemble-const (list (list 'const (car expression)))))
        ", ")
      (assemble-const (list (list 'const (cdr expression)))))
    ")"))

(define (assemble-const expression)
  (cond ((self-evaluating? (cadar expression)) (cadar expression))
        ((null? (cadar expression)) "list()")
        ((pair? (cadar expression))
         (assemble-const-pair (cadar expression)))
        ((symbol? (cadar expression))
         (sbappend
           (sbappend
             (sbappend
               (make-string-builder) "stringToSymbol(\"")
             (symbol-name (cadar expression)))
           "\")"))
        (else
          (error 'assemble-const
                 "invalid const expression"
                 expression))))

(define (assembly-reg? expression)
  (tagged-list? (car expression) 'reg))

(define (assembly-label? expression)
  (tagged-list? (car expression) 'label))

(define (assembly-entry? expression)
  (symbol? (car expression)))

(define (assembly-op? expression)
  (tagged-list? (car expression) 'op))

(define (op-name symbol)
  (cond ((eq? symbol 'lookup-variable-value) "lookupVar")
        ((eq? symbol 'list) "list")
        ((eq? symbol 'primitive-procedure?) "isPrimitiveProcedure")
        ((eq? symbol 'compiled-procedure-entry) "compiledProcedureEntry")
        ((eq? symbol 'apply-primitive-procedure) "applyPrimitive")
        ((eq? symbol 'define-variable!) "defineVar")
        ((eq? symbol 'false?) "isFalse")
        ((eq? symbol 'set-variable-value!) "setVar")
        ((eq? symbol 'make-compiled-procedure) "makeCompiledProcedure")
        ((eq? symbol 'compiled-procedure-env) "compiledProcedureEnv")
        ((eq? symbol 'extend-environment) "extendEnv")
        ((eq? symbol 'cons) "cons")
        (else
          (error 'op-name
                 "invalid operation"
                 symbol))))

(define (op-arg-list args)
  (define (op-arg-list args builder first?)
    (cond ((null? args) builder)
          (else
            (op-arg-list
              (cdr args)
              (sbappend
                (sbappend builder (if first? "" ", "))
                (assemble-expression args))
              false))))
  (op-arg-list args (make-string-builder) true))

(define (assemble-op expression)
  (sbappend
    (sbappend
      (sbappend
        (sbappend (make-string-builder)
                  (op-name (cadar expression)))
        "(")
      (op-arg-list (cdr expression)))
    ")"))

(define (assemble-expression expression)
  (cond ((assembly-const? expression)
         (assemble-const expression))
        ((assembly-reg? expression)
         (symbol-name (cadar expression)))
        ((assembly-label? expression)
         (js-name (cadar expression)))
        ((assembly-entry? expression)
         (js-name (car expression)))
        ((assembly-op? expression)
         (assemble-op expression))
        (else
          (error 'assemble-expression
                 "invalid assembly expression"
                 expression))))

(define (assembly-assign? program)
  (tagged-list? (car program) 'assign))

(define (assemble-assign program)
  (list
    (cdr program)
    (sbappend
      (sbappend
        (sbappend
          (sbappend
            (tab (tab (make-string-builder)))
            (symbol-name (cadar program)))
          " = ")
        (assemble-expression (cddar program)))
      ";\n")))

(define (assembly-test? program)
  (tagged-list? (car program) 'test))

(define (assemble-test program)
  (list
    (cddr program)
    (sbappend
      (sbappend
        (sbappend
          (sbappend
            (sbappend
              (sbappend
                (sbappend
                  (sbappend
                    (sbappend
                      (sbappend (tab (tab (make-string-builder))) "if (")
                      (assemble-expression (cdar program)))
                    ") {\n")
                  (sbappend (tab (tab (tab (make-string-builder))))
                            "return "))
                (assemble-expression (cdadr program)))
              ";\n")
            (sbappend (tab (tab (make-string-builder))) "}\n"))
          (sbappend (tab (tab (make-string-builder))) "return "))
        (assemble-expression (cddr program)))
      ";\n")))

(define (assemble-entry program)
  (list (cdr program)
        (sbappend
          (sbappend
            (sbappend
              (sbappend
                (sbappend
                  (sbappend
                    (sbappend (tab (tab (make-string-builder)))
                              "return ")
                    (assemble-expression program))
                  ";\n")
                (sbappend (tab (make-string-builder)) "};\n\n"))
              (sbappend (tab (make-string-builder)) "var "))
            (assemble-expression program))
          " = function () {\n")))

(define (assembly-goto? program)
  (tagged-list? (car program) 'goto))

(define (assemble-goto program)
  (list (cdr program)
        (sbappend
          (sbappend
            (sbappend (tab (tab (make-string-builder)))
                      "return ")
            (assemble-expression (cdar program)))
          ";\n")))

(define (assembly-save? program)
  (tagged-list? (car program) 'save))

(define (assemble-save program)
  (list
    (cdr program)
    (sbappend
      (sbappend
        (sbappend (tab (tab (make-string-builder))) "save(")
        (symbol-name (cadar program)))
      ");\n")))

(define (assembly-restore? program)
  (tagged-list? (car program) 'restore))

(define (assemble-restore program)
  (list
    (cdr program)
    (sbappend
      (sbappend (tab (tab (make-string-builder)))
                (symbol-name (cadar program)))
      " = restore();\n")))

(define (assembly-perform? program)
  (tagged-list? (car program) 'perform))

(define (assemble-perform program)
  (list
    (cdr program)
    (sbappend
      (sbappend (tab (tab (make-string-builder)))
                (assemble-expression (cdar program)))
      ";\n")))

(define (assembly-js-import-code? program)
  (tagged-list? (car program) 'js-import-code))

(define (assemble-js-import-code program)
  (list
    (cdr program)
    (sbappend
      (sbappend
        (sbappend
          (sbappend
            (sbappend (make-string-builder)
                      "importModule(stringToSymbol(\"")
            (symbol-name (cadr (cadar program))))
          "\"), (function (exports) {\n")
        (cadr (caddar program)))
      "\nreturn exports;\n})({}));\n")))

(define (assemble-body program builder)
  (if (null? program)
    builder
    (let ((step
            (cond ((assembly-assign? program)
                   (assemble-assign program))
                  ((assembly-test? program)
                   (assemble-test program))
                  ((assembly-entry? program)
                   (assemble-entry program))
                  ((assembly-goto? program)
                   (assemble-goto program))
                  ((assembly-save? program)
                   (assemble-save program))
                  ((assembly-restore? program)
                   (assemble-restore program))
                  ((assembly-perform? program)
                   (assemble-perform program))
                  ((assembly-js-import-code? program)
                   (assemble-js-import-code program))
                  (else
                    (error 'assemble-body
                           "invalid instruction"
                           program)))))
      (assemble-body (car step)
                     (sbappend builder (cadr step))))))

(define (assemble program)
  (builder->string
    (sbappend
      (sbappend
        (sbappend (make-string-builder)
                  js-head)
        (assemble-body program (make-string-builder)))
      js-tail)))

(let ((program (compile-sequence
                 ; '((define (factorial n)
                 ;     (if (eq? n 1)
                 ;       1
                 ;       (* (factorial (- n 1)) n)))
                 ;   (print (factorial 3)))
                 '((js-import-code / "
        exports.log = function () {
            console.log.apply(console, Array.prototype.slice.call(arguments));
        };
                                   ")
                   1
                   (log 1)
                   (log 2)
                   (define a 3)
                   (log a)
                   (if 4 5 6)
                   (log (if 7 8 9))
                   '(some)
                   (set! a 10)
                   (log a)
                   (lambda (x) x)
                   (lambda x x)
                   ((lambda (x) x) 11)
                   (log ((lambda x x) 12 13 14))
                   (define (factorial n)
                     (if (= n 0)
                       1
                       (* n (factorial (- n 1)))))
                   (log (factorial 120)))
                   ; (js-export name definition)
                   ; (js-import-file name path)
                   ; (js-import-code name code))
                 'val
                 'next)))
  ; (start (make-machine
  ;    (list-union '(env)
  ;                (list-union (registers-needed program)
  ;                            (registers-modified program)))
  ;    compile-ops
  ;    (append '((assign env (op get-global-environment)))
  ;            (caddr program)))))
  (let ((js (assemble (caddr program))))
    (out js)))

; do lexical addressing of variables
