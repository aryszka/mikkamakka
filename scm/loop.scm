(import "scm/lang.scm")
(import "scm/read.scm")
(import "scm/eval.scm")
(import "scm/print.scm")
(import "scm/compile.scm")


(def (read-one state)
  (let (r (read state:reader))
    (cond ((= r:state eof) (assign state {done true}))
          ((error? r:state) (fatal r:state))
          (else (assign state {reader r expression r:state})))))


(def (eval-one state)
  (let (v (eval-env state:environment state:expression))
    (cond ((error? v) (fatal v))
          (else (assign state {value v})))))


(def (print-out state print-func exp)
  (let (p (print-func state:printer exp))
    (cond ((error? p:state) (fatal p:state))
          (else (assign state {printer p})))))


(def (compile-one state)
  (let (c (compile state:compiler state:expression))
    (cond (c:error (fatal c:error))
          (else (assign state {compiler c content c:output})))))


(def (write-out output content)
  (let (output (fwrite output content))
    (cond ((error? (fstate output)) (fatal (fstate output)))
          (else output))))


(def (loop state . f)
  (def (if-not-done f)
    (fn (state)
        (if state:done state (f state))))

  (def (loop state)
    ((apply
	   compose
	   (map if-not-done (append f (list loop))))
     state))

  (loop state))


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


(export read-eval-print    read-eval-print
		read-compile-write read-compile-write)
