;; Define compatibility functions to match ISLisp's standard input/output access
(defun standard-input () *standard-input*)
(defun standard-output () t)
(load "test.lsp")
