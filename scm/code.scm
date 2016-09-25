(import "scm/lang.scm")


(def definition-expression (string->error "definition in expression position"))
(def invalid-expression (string->error "invalid expression"))
(def invalid-statement (string->error "invalid expression in statement position"))
(def invalid-cond (string->error "invalid cond expression"))
(def circular-import (string->error "circular-import"))
(def not-implemented (string->error "not implemented"))
(def invalid-literal (string->error "invalid literal"))
(def invalid-value-list (string->error "invalid value list"))
(def invalid-import (string->error "invalid import expression"))


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


(def (and->if v)
  (cond ((nil? v) true)
        ((nil? (cdr v)) (car v))
        (else (list 'if (car v) (and->if (cdr v)) false))))


(def (or->if v)
  (cond ((nil? v) false)
        ((nil? (cdr v)) (car v))
        (else (list 'if (car v) (car v) (or->if (cdr v))))))


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


(def (fn-signature v)
  (cond ((nil? v) {count 0 var? false names '()})
        ((symbol? v) {count 0 var? true names v})
        ((pair? v) (let (signature (fn-signature (cdr v)))
                     (assign signature {count (inc signature:count)
                                         names (cons (car v) signature:names)})))
        (else invalid-fn)))


(def (import-def exp)
  (cond ((= (len exp) 2) {import-name false module-name (car (cdr exp))})
        ((= (len exp) 3) {import-name (car (cdr exp)) module-name (car (cdr (cdr exp)))})
        (else invalid-import)))


(export quote? quote?
        def? def?
        vector-form? vector-form?
        struct-form? struct-form?
        if? if?
        and? and?
        or? or?
        fn? fn?
        begin? begin?
        cond? cond?
        let? let?
        test? test?
        export? export?
        import? import?
        application? application?
        current-env? current-env?
        make-fn make-fn
        value-def? value-def?
        valid-def? valid-def?
        valid-value-def? valid-value-def?
        valid-function-def? valid-function-def?
        def-name def-name
        def-value def-value
        and->if and->if
        or->if or->if
        cond->if cond->if
        let-body let-body
        fn-signature fn-signature
        import-def import-def
		invalid-expression invalid-expression)
