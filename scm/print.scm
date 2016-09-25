(import "scm/lang.scm")
(import "scm/code.scm")


(def (printer output)
  {output output state 'ok})


(def (print-raw p s)
  (if (error? p:state) p
    (let (output (fwrite p:output s))
      (assign p {output output state (fstate output)}))))


(def (print-quote-sign p) (print-raw p "'"))


(def (print-symbol p v quoted?)
  (print-raw
    (if quoted? p (print-quote-sign p))
    (symbol->string v)))


(def (print-quote p v)
  (printq (print-quote-sign p) (car (cdr v)) false))


(def (print-pair p v quoted?)
  (def (print-space p v)
    (if (nil? (cdr v)) p
      (print-raw p " ")))
  (def (print-items p v)
    (cond ((nil? v) (print-raw p ")"))
          ((not (pair? v)) (print-raw (printq (print-raw p ". ") v true) ")"))
          (else (print-items
                  (print-space
                    (printq p (car v) true)
                    v)
                  (cdr v)))))
  (let (p (print-raw
            (if quoted? p (print-quote-sign p)) "("))
    (print-items p v)))


(def (print-vector p v)
  (def (print-space p i)
    (if (>= i (- (len v) 1)) p
      (print-raw p " ")))
  (def (print-items p i)
    (cond ((= i (len v)) (print-raw p "]"))
          (else (print-items
                  (print-space
                    (printq p (vector-ref v i) true)
                    i)
                  (inc i)))))
  (let (p (print-raw p "["))
    (print-items p 0)))


(def (print-struct p s)
  (def (print-space p n)
    (if (nil? (cdr n)) p
      (print-raw p " ")))
  (def (print-items p n)
    (cond ((nil? n) (print-raw p "}"))
          (else (print-items
                  (print-space
                    (printq
                      (print-raw (printq p (car n) true) " ")
                      (field s (car n))
                      true)
                    n)
                  (cdr n)))))
  (let (p (print-raw p "{"))
    (print-items p (struct-names s))))


(def (printq p v quoted?)
  (cond ((number? v) (print-raw p (number->string v)))
        ((string? v) (print-raw p v))
        ((bool? v) (print-raw p (bool->string v)))
        ((symbol? v) (print-symbol p v quoted?))
        ((sys? v) (print-raw p (sys->string v)))
        ((error? v) (print-raw p (error->string v)))
        ((quote? v) (print-quote p v))
        ((or (pair? v) (nil? v)) (print-pair p v quoted?))
        ((vector? v) (print-vector p v))
        ((struct? v) (print-struct p v))
        ((env? v) (print-raw p (env->string v)))
        ((function? v) (print-raw p (function->string v)))
        (else not-implemented)))


(def (print p v) (printq p v false))


(export printer printer
        print print
        print-raw print-raw)
