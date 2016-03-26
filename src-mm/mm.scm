((lambda ()
   (define (writeln v)
     (write-file stdout v)
     (write-file stdout "\n"))


   (define (fibonacci n)
     (cond ((== n 0) 1)
           ((== n 1) 1)
           (else (+ (fibonacci (- n 1))
                    (fibonacci (- n 2))))))


   (writeln (fibonacci 24))))
