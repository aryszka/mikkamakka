(import "scm/lang.scm")
(import "scm/print.scm")


(def (trace message . values)
  (def (trace out values)
    (let (out (print out (car values)))
      (cond ((nil? (cdr values))
             (fwrite out:output "\n")
             (car values))
            (else
              (trace
                (assign out {output (fwrite out:output " ")})
                (cdr values))))))
  (trace (printer (stderr)) (cons message values)))


(export trace trace)
