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


(def (id x) x)


(def (list . x) x)


(def (apply f a)
  (cond ((vector? f) (vector-ref f (car a)))
        ((struct? f) (field f (car a)))
        ((compiled-function? f) (apply-compiled f a))
        (else (fatal not-function))))
        ; (else (let (c (composite f))
        ;         (eval-seq (extend-env (car c) (car (cdr c)) a) (cdr (cdr c)))))))


(def (call f . a) (apply f a))


(def (fold f i l)
  (cond ((nil? l) i)
        (else (fold f (f (car l) i) (cdr l)))))


(def (foldr f i l)
  (cond ((nil? l) i)
        (else (f (car l) (foldr f i (cdr l))))))


(def (map f . l)
  (cond ((nil? (car l)) '())
        ((nil? (cdr l))
         (cons (f (car (car l)))
               (map f (cdr (car l)))))
        (else
          (cons (apply f (map car l))
                (apply map (cons f (map cdr l)))))))


(def (append . l)
  (cond ((nil? l) nil)
        ((nil? (cdr l)) (car l))
        (else (foldr cons
                     (apply append (cdr l))
                     (car l)))))


(def (part f . a) (fn b (apply f (append a b))))


(def (partr f . a) (fn b (apply f (append b a))))


(def reverse (part fold cons nil))


(def (reverse-irregular l)
  (cond ((or (nil? l) (nil? (cdr l))) irregular-cons)
        (else (fold cons (cons (car (cdr l)) (car l)) (cdr (cdr l))))))


(def (inc n) (+ n 1))


(def (dec n) (- n 1))


(def (>= . n)
  (cond ((nil? n) false)
        ((nil? (cdr n)) true)
        ((and (not (> (car n) (car (cdr n))))
              (not (= (car n) (car (cdr n)))))
         false)
        (else (apply >= (cdr n)))))


(def list-len (part fold (fn (_ c) (inc c)) 0))


(def (len v)
  (cond ((vector? v) (vector-len v))
        ((struct? v) (len (struct-names s)))
        (else (list-len v))))


(def (mem f l)
  (cond ((nil? l) false)
        ((f (car l)) l)
        (else (mem f (cdr l)))))


(def (memq v l) (mem (part = v) l))


(def (notf f) (fn (i) (not (f i))))


(def (every? f . l) (not (mem (notf f) l)))


(def (any? . v) (and (mem id v) true))


(def (take n l)
  (cond ((= n 0) nil)
        (else (cons (car l)
                    (take (- n 1) (cdr l))))))


(def (drop n l)
  (cond ((= n 0) l)
        (else (drop (dec n) (cdr l)))))


(def (pad f n) (fn a (apply f (drop n a))))


(def (padr f n) (fn a (apply f (take (- (len a) n) a))))


(def (flip f) (fn a (apply f (reverse a))))


(def (check-types v . types)
  (apply any? (map (fn (t?) (t? v)) types)))


(def (compose . f) (partr (part fold call) f))


(export trace trace
        id id
        list list
        apply apply
        call call
        fold fold
        foldr foldr
        map map
        append append
        part part
        partr partr
        reverse reverse
        reverse-irregular reverse-irregular
        inc inc
        dec dec
        >= >=
        len len
        mem mem
        memq memq
        notf notf
        every? every?
        any? any?
        take take
        drop drop
        pad pad
        padr padr
        flip flip
        check-types check-types
        compose compose)
