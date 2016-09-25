(import "scm/lang.scm")
(import "scm/loop.scm")


(def (open-file name)
  (let (f (fopen name))
	(if (error? f) (fatal f) f)))


(read-compile-write
  (open-file "scm/mm.scm")
  (stdout))
