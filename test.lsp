(defmacro test (source expect)
  (let ((result (gensym)))
    `(let ((,result ,source))
       (if (equalp ,result ,expect)
           (format (standard-output) "PASS: (test ~S ~S)~%"
                   (quote ,source)
                   ,expect)
           (format (standard-output) "FAIL: (test ~S ~S)~%  but ~S~%"
                   (quote ,source)
                   ,expect
                   ,result)
       ))))

(defun field (key m)
  (and
    m
    (consp m)
    (if (equal (car (car m)) key)
      (car (cdr (car m)))
      (field key (cdr m)))))

(let ((data (read (standard-input) nil nil)))
  (test (field 'struct data) 'Foo)
  (test (field 'bar data) "hogehoge")
  (test (field 'baz data) 0.1)
  (test (field 'qux data) #(1 2 3 4))
  (let ((m (field 'quux data)))
    (test (field "ahaha" m) 1)
    (test (field "ihihi" m) 2)
    (test (field "ufufu" m) 3))
  (test (field 'quuux data) "a\"\\
	b"))
