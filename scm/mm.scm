(def tnone 0)
(def undefined (fn () 'ok))

(def (reader input)
  { input         input 
    token-type    tnone
    value         undefined
    current-token ""
    in-list?      false
    close-list?   false
    list-items    nil
    list-cons     false
    cons-items    0
  })


(def (read reader) reader)


; (test
;   (test "as an initial test, returns the reader"
;     (let (in-initial (fopen "buffer:" file-mode-none)
;           in-ready   (fwrite in "test data")
;           r          (reader in))
;       (= (read r) r))))
