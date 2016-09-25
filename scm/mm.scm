(import "scm/lang.scm")
(import "scm/read.scm")
(import "scm/eval.scm")
(import "scm/print.scm")
(import "scm/compile.scm")


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


(def (read-one loop-state)
  (let (r (read loop-state:reader))
    (cond ((= r:state eof) (assign loop-state {done true}))
          ((error? r:state) (fatal r:state))
          (else (assign loop-state {reader r expression r:state})))))


(def (eval-one loop-state)
  (let (v (eval-env loop-state:environment loop-state:expression))
    (cond ((error? v) (fatal v))
          (else (assign loop-state {value v})))))


(def (print-out loop-state print-func exp)
  (let (p (print-func loop-state:printer exp))
    (cond ((error? p:state) (fatal p:state))
          (else (assign loop-state {printer p})))))


(def (compile-one loop-state)
  (let (c (compile loop-state:compiler loop-state:expression))
    (cond (c:error (fatal c:error))
          (else (assign loop-state {compiler c content c:output})))))


(def (write-out output content)
  (let (output (fwrite output content))
    (cond ((error? (fstate output)) (fatal (fstate output)))
          (else output))))


(def (loop loop-state . f)
  (def (if-not-done f)
    (fn (loop-state)
        (if loop-state:done loop-state (f loop-state))))

  (def (loop loop-state)
    ((apply compose (map if-not-done (append f (list loop))))
     loop-state))

  (loop loop-state))


(def (read-eval-print input output)
  (def (print-one repl) (print-out repl print repl:value))
  (def (print-line repl) (print-out repl print-raw "\n"))

  (loop
    {reader      (reader input)
     environment (current-env)
     printer     (printer output)
     done        false}

    read-one
    eval-one
    print-one
    print-line))


(def (read-compile-write input output)
  ((compose
     (partr loop read-one compile-one)
     (partr call 'compiler)
     close-compiler
     (partr call 'output)
     (part write-out output))

   {reader   (reader input)
    output   output
    compiler (compiler)
    done     false}))


(def (input argv)
  (cond ((= (len argv) 0) (stdin))
        (else
          (let (f (fopen (car argv)))
            (cond ((error? (fstate f))
                   (fatal (fstate f)))
                  (else f))))))


(def (command argv)
  (cond ((= (len argv) 0)
         {command read-eval-print
          input   (input argv)})

        (else
         (cond ((= (car argv) "compile")
                {command read-compile-write
                 input   (input (cdr argv))})

               ((= (car argv) "run")
                {command read-eval-print
                 input   (input (cdr argv))})

               (else
                 (fatal (string->error "invalid command")))))))


(def cmd (command (cdr (argv))))
(cmd:command cmd:input (stdout))
