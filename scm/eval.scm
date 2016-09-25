(import "scm/lang.scm")
(import "scm/code.scm")
(import "scm/read.scm")


(def definition-expression (string->error "definition in expression position"))


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
    (fn a (call eval-seq (extend-env env signature:names a) body))))


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
        (else not-implemented)))


(export eval-env eval-env)
