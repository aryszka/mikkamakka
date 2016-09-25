(import "scm/lang.scm")
(import "scm/read.scm")
(import "scm/compile.scm")
(import "scm/eval.scm")
(import "scm/print.scm")


(def (read-eval-print)
  (def (loop env r p)
    (let (r (read r))
      (cond
        ((= r:state eof) 'ok)
        ((error? r:state) (fatal r:state))
        (else
          (let (v (eval-env env r:state))
            (cond
              ((error? v) (fatal v))
              (else
                (let (p (print p v))
                  (cond
                    ((error? p:state) (fatal p:state))
                    (else (loop env r (print-raw p "\n"))))))))))))
  (loop (current-env) (reader (stdin)) (printer (stdout))))


(def (read-compile-write)
  (def (compile-top c exp)
    (compiler-compose
      c
      compile exp
      compiler-append ";"))
  (def (loop r c)
    (let (r (read r))
      (cond ((error? r:state)
             {input r:input
              output c:output
              read-error r:state
              compile-error c:error})
            (else (loop r (compile-top c r:state))))))
  (let (r  (loop (reader (stdin)) (compiler)))
    (cond ((error? r:compile-error) (fatal r:compile-error))
          ((= r:read-error eof)
           (fclose r:input)
           (fwrite (stdout) compiled-head)
           (fwrite (stdout) r:output)
           (fwrite (stdout) compiled-tail)
           'ok)
          (else (fatal r:read-error)))))


; (read-compile-write)
(read-eval-print)
