(import "scm/lang.scm")
(import "scm/loop.scm")


; TODO: closing opened files


(def (input argv)
  (if (= (len argv) 0)
	(stdin)
	(let (f (fopen (car argv)))
	  (if (error? f) (fatal f) f))))


(def (command argv)
  (cond ((= (len argv) 0)
         {command read-eval-print
          input   (input argv)})

		((= (car argv) "compile")
		 {command read-compile-write
		 input   (input (cdr argv))})

		((= (car argv) "run")
		 {command read-eval-print
		 input   (input (cdr argv))})

		(else
		  (fatal (string->error "invalid command")))))


(def cmd (command (cdr (argv))))
(cmd:command cmd:input (stdout))
