;; Display tabs as two spaces wide.
(setq tab-width 2)
(setq-default tab-width 2)

;; Use four spaces in python-mode.
(add-hook 'python-mode-hook
					(lambda ()
						(setq indent-tabs-mode nil)
						(setq tab-width 4)
						(setq python-indent 4)))

;; Use two tabs in js-mode, displayed with width two.
(add-hook 'js-mode-hook
					(lambda ()
						(setq indent-tabs-mode t)
						(setq tab-width 2)
						(setq js-indent-level 2)))

;; Use four spaces in java-mode as well.
(add-hook 'java-mode-hook
					(lambda ()
						(setq indent-tabs-mode nil)
						(setq tab-width 4)
						(setq c-basic-offset 4)
						(setq whitespace-line-column 100)))

;; Use two spaces in html-mode.
(add-hook 'html-mode-hook
					(lambda ()
						(setq indent-tabs-mode nil)
						(setq sh-basic-offset 2)
						(setq sh-indentation 2)))

;; Use two spaces in sh-mode.
(add-hook 'sh-mode-hook
					(lambda ()
						(setq indent-tabs-mode nil)
						(setq sh-basic-offset 2)
						(setq sh-indentation 2)))

(set-default-font "-adobe-courier-medium-r-normal--14-140-75-75-m-90-iso8859-1")

;; Sets the path for backup files generated automatically by emacs (represented
;; by the filename with a tilde appended to the end of it.)

;; (source: http://www.skrakes.com/?p=146)

(defvar backup-dir "~/.emacs.d/backups/")
(defvar autosave-dir "~/.emacs.d/autosaves/")

;; Create backup-directory and autosave-directory if they don't already exist

(make-directory backup-dir t)
(make-directory autosave-dir t)

(setq backup-directory-alist `(("." . ,backup-dir)))
(setq auto-save-file-name-trnsforms `(("." ,autosave-dir t)))
(setq backup-by-copying t)

(setq delete-old-versions t
			kept-new-versions 6
			kept-old-versions 2
			version-control t)

;; Use html-mode for .tmpl files.
(add-to-list 'auto-mode-alist '("\\.tmpl\\'" . html-mode))

;; We want go-mode, and goimports + gofmt hook.
(add-to-list 'load-path "~/.emacs.d/go-mode")
(require 'go-mode-autoloads)
(setq gofmt-command "goimports")
(add-hook 'before-save-hook #'gofmt-before-save)

(add-to-list 'load-path (concat (getenv "GOPATH") "/src/github.com/golang/lint/misc/emacs"))
(require 'golint)

;; We want protobuf-mode for .proto files.
;; (require 'protobuf-mode)

;; We want to automatically format .tf files on save.
(add-hook 'terraform-mode-hook #'terraform-format-on-save-mode)

;; ELPA packages; install interactively with M-x package-list-packages.
(require 'package)
(add-to-list 'package-archives
						 '("melpa" . "http://melpa.org/packages/") t)
(when (< emacs-major-version 24)
	;; For important compatibility libraries like cl-lib
	(add-to-list 'package-archives '("gnu" . "http://elpa.gnu.org/packages/")))
(package-initialize)

;; Add homebrew's bin directory to exec-path.
(setq exec-path (append exec-path '("/usr/local/homebrew/bin")))

(custom-set-variables
 ;; custom-set-variables was added by Custom.
 ;; If you edit it by hand, you could mess it up, so be careful.
 ;; Your init file should contain only one such instance.
 ;; If there is more than one, they won't work right.
 '(package-selected-packages
	 (quote
		(solidity-mode terraform-mode protobuf-mode smart-tabs-mode))))
(custom-set-faces
 ;; custom-set-faces was added by Custom.
 ;; If you edit it by hand, you could mess it up, so be careful.
 ;; Your init file should contain only one such instance.
 ;; If there is more than one, they won't work right.
 )

