(defun equalp (x y)
  (equal x y))

(with-open-input-file
  (fd "sample.log")
  (with-standard-input
    fd
    (load "test.lsp")))
