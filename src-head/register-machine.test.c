#include <stdlib.h>
#include <string.h>
#include "sys.h"
#include "testing.h"
#include "error.h"
#include "error.mock.h"
#include "sysio.h"
#include "sysio.mock.h"
#include "compound-types.h"
#include "stack.h"
#include "number.h"
#include "string.h"
#include "io.h"
#include "value.h"
#include "sprint-list.h"
#include "environment.h"
#include "register-machine.h"
#include "primitives.h"

// ((lambda () 42))
//
// ->
//
// (make-compiled-procedure proc 1)
// (goto (label 2))
// 1
// (init-proc-env ())
// (initreg val (const 42))
// (goto (reg continue))
// 2
// (initargs)
// (branchproc 3)
// 4
// (goto proclabel)
// 3
// (apply-primitive-procedure val)
// (goto (reg continue))
// 5
// 
void test_number_small() {
	regmachine rm = mkregmachine();
	for (;;) {
		long labelval = valrawint(rm->label);
		switch (labelval) {
		case 0:
			mkcompiledprocreg(rm, (void *)&rm->proc, 1);
			gotolabel(rm, 2); break;
		case 1:
			initprocenv(rm, null);
			initreg((void *)&rm->val, numvali(42));
			gotoreg(rm, rm->cont); break;
		case 2:
			initargs(rm);
			if (branchproc(rm, 3)) { break; }
		case 4:
			gotoproc(rm); break;
		case 3:
			applyprimitivereg(rm, (void *)&rm->val);
			gotoreg(rm, rm->cont); break;
		case 5:
		default:
			assert(valrawint(rm->val) == 42, "reg-machine: expression, number");
			return;
		}
	}

	freeregmachine(rm);
}

// ((lambda () 2.25))
//
// ->
//
// (make-compiled-procedure proc 1)
// (goto (label 2))
// 1
// (init-proc-env ())
// (initreg val (const 2.25))
// (goto (reg continue))
// 2
// (initargs)
// (branchproc 3)
// 4
// (goto proclabel)
// 3
// (apply-primitive-procedure val)
// (goto (reg continue))
// 5
// 
void test_number_rational() {
	regmachine rm = mkregmachine();
	char *s;
	for (;;) {
		long labelval = valrawint(rm->label);
		switch (labelval) {
		case 0:
		mkcompiledprocreg(rm, (void *)&rm->proc, 1);
			gotolabel(rm, 2); break;
		case 1:
			initprocenv(rm, null);
			initreg((void *)&rm->val, numvalc("2.25"));
			gotoreg(rm, rm->cont); break;
		case 2:
			initargs(rm);
			if (branchproc(rm, 3)) { break; }
		case 4:
			gotoproc(rm); break;
		case 3:
			applyprimitivereg(rm, (void *)&rm->val);
			gotoreg(rm, rm->cont); break;
		case 5:
		default:
			s = sprintraw(rm->val);
			assert(!strcoll(s, "2.25"), "reg-machine: rational number");
			free(s);
			return;
		}
	}

	freeregmachine(rm);
}

// ((lambda () "some string"))
//
// ->
//
// (make-compiled-procedure proc 1)
// (goto (label 2))
// 1
// (init-proc-env ())
// (initreg val (const "some string"))
// (goto (reg continue))
// 2
// (initargs)
// (branchproc 3)
// 4
// (goto proclabel)
// 3
// (apply-primitive-procedure val)
// (goto (reg continue))
// 5
//
void test_string() {
	regmachine rm = mkregmachine();
	char *s;
	for (;;) {
		long labelval = valrawint(rm->label);
		switch (labelval) {
		case 0:
			mkcompiledprocreg(rm, (void *)&rm->proc, 1);
			gotolabel(rm, 2); break;
		case 1:
			initprocenv(rm, null);
			initreg((void *)&rm->val, stringval("some string"));
			gotoreg(rm, rm->cont); break;
		case 2:
			initargs(rm);
			if (branchproc(rm, 3)) { break; }
		case 4:
			gotoproc(rm); break;
		case 3:
			applyprimitivereg(rm, (void *)&rm->val);
			gotoreg(rm, rm->cont); break;
		case 5:
		default:
			s = sprintraw(rm->val);
			assert(!strcoll(s, "some string"), "reg-machine: string");
			free(s);
			return;
		}
	}

	freeregmachine(rm);
}

// ((lambda () null))
//
// ->
//
// (make-compiled-procedure proc 1)
// (goto (label 2))
// 1
// (init-proc-env ())
// (initreg val (const ()))
// (goto (reg continue))
// 2
// (initargs)
// (branchproc 3)
// 4
// (goto proclabel)
// 3
// (apply-primitive-procedure val)
// (goto (reg continue))
// 5
//
void test_null() {
	regmachine rm = mkregmachine();
	char *s;
	for (;;) {
		long labelval = valrawint(rm->label);
		switch (labelval) {
		case 0:
			mkcompiledprocreg(rm, (void *)&rm->proc, 1);
			gotolabel(rm, 2); break;
		case 1:
			initprocenv(rm, null);
			initreg((void *)&rm->val, null);
			gotoreg(rm, rm->cont); break;
		case 2:
			initargs(rm);
			if (branchproc(rm, 3)) { break; }
		case 4:
			gotoproc(rm); break;
		case 3:
			applyprimitivereg(rm, (void *)&rm->val);
			gotoreg(rm, rm->cont); break;
		case 5:
		default:
			assert(rm->val == null, "reg-machine: null");
			return;
		}
	}

	freeregmachine(rm);
}

// ((lambda () 'a))
//
// ->
//
// (make-compiled-procedure proc 1)
// (goto (label 2))
// 1
// (init-proc-env ())
// (initreg val (const a))
// (goto (reg continue))
// 2
// (initargs)
// (branchproc 3)
// 4
// (goto proclabel)
// 3
// (apply-primitive-procedure val)
// (goto (reg continue))
// 5
//
void test_symbol() {
	regmachine rm = mkregmachine();
	char *s;
	for (;;) {
		long labelval = valrawint(rm->label);
		char *s;
		switch (labelval) {
		case 0:
			mkcompiledprocreg(rm, (void *)&rm->proc, 1);
			gotolabel(rm, 2); break;
		case 1:
			initprocenv(rm, null);
			initreg((void *)&rm->val, symval("a"));
			gotoreg(rm, rm->cont); break;
		case 2:
			initargs(rm);
			if (branchproc(rm, 3)) { break; }
		case 4:
			gotoproc(rm); break;
		case 3:
			applyprimitivereg(rm, (void *)&rm->val);
			gotoreg(rm, rm->cont); break;
		case 5:
		default:
			s = sprintraw(rm->val);
			assert(!strcoll(s, "a"), "reg-machine: symbol");
			free(s);
			return;
		}
	}

	freeregmachine(rm);
}

// ((lambda ()
//	'(a . b)))
//
// ->
//
// (make-compiled-procedure proc 1)
// (goto (label 2))
// 1
// (init-proc-env ())
// (initreg val (const (a . b)))
// (goto (reg continue))
// 2
// (initargs)
// (branchproc 3)
// 4
// (goto proclabel)
// 3
// (apply-primitive-procedure val)
// (goto (reg continue))
// 5
// 
void test_pair() {
	regmachine rm = mkregmachine();
	char *s;
	for (;;) {
		long labelval = valrawint(rm->label);
		switch (labelval) {
		case 0:
			mkcompiledprocreg(rm, (void *)&rm->proc, 1);
			gotolabel(rm, 2); break;
		case 1:
			initprocenv(rm, null);
			initreg((void *)&rm->val, pairval(symval("a"), symval("b")));
			gotoreg(rm, rm->cont); break;
		case 2:
			initargs(rm);
			if (branchproc(rm, 3)) { break; }
		case 4:
			gotoproc(rm); break;
		case 3:
			applyprimitivereg(rm, (void *)&rm->val);
			gotoreg(rm, rm->cont); break;
		case 5:
		default:
			s = sprintraw(rm->val);
			assert(!strcoll(s, "(a . b)"), "reg-machine: pair");
			free(s);
			return;
		}
	}

	freeregmachine(rm);
}

// '((lambda ()
//	 (define a 1)
//	 a))
//
// ->
//
// (make-compiled-procedure proc 1)
// (goto (label 2))
// 1
// (init-proc-env ())
// (initreg val (const 1))
// (define-variable a)
// (get-variable val a)
// (goto (reg continue))
// 2
// (initargs)
// (branchproc 3)
// 4
// (goto proclabel)
// 3
// (apply-primitive-procedure val)
// (goto (reg continue))
// 5
// 
void test_defvar_getvar() {
	regmachine rm = mkregmachine();
	for (;;) {
		long labelval = valrawint(rm->label);
		switch (labelval) {
		case 0:
			mkcompiledprocreg(rm, (void *)&rm->proc, 1);
			gotolabel(rm, 2); break;
		case 1:
			initprocenv(rm, null);
			initreg((void *)&rm->val, numvali(1));
			defenvvar(rm, "a");
			getenvvar(rm, (void *)&rm->val, "a");
			gotoreg(rm, rm->cont); break;
		case 2:
			initargs(rm);
			if (branchproc(rm, 3)) { break; }
		case 4:
			gotoproc(rm); break;
		case 3:
			applyprimitivereg(rm, (void *)&rm->val);
			gotoreg(rm, rm->cont); break;
		case 5:
		default:
			assert(valrawint(rm->val) == 1, "reg-machine: defvar, getvar");
			return;
		}
	}

	freeregmachine(rm);
}

// '((lambda ()
//	 (define a 1)
//	 (set! a 2)
//	 a))
//
// ->
//
// (make-compiled-procedure proc 1)
// (goto (label 2))
// 1
// (init-proc-env ())
// (initreg val (const 1))
// (define-variable a)
// (initreg val (const 2))
// (set-variable-value val a)
// (get-variable val a)
// (goto (reg continue))
// 2
// (initargs)
// (branchproc 3)
// 4
// (goto proclabel)
// 3
// (apply-primitive-procedure val)
// (goto (reg continue))
// 5
// 
void test_defvar_setvar_getvar() {
	regmachine rm = mkregmachine();
	for (;;) {
		long labelval = valrawint(rm->label);
		switch (labelval) {
		case 0:
			mkcompiledprocreg(rm, (void *)&rm->proc, 1);
			gotolabel(rm, 2); break;
		case 1:
			initprocenv(rm, null);
			initreg((void *)&rm->val, numvali(1));
			defenvvar(rm, "a");
			initreg((void *)&rm->val, numvali(2));
			setenvvar(rm, (void *)&rm->val, "a");
			getenvvar(rm, (void *)&rm->val, "a");
			gotoreg(rm, rm->cont); break;
		case 2:
			initargs(rm);
			if (branchproc(rm, 3)) { break; }
		case 4:
			gotoproc(rm); break;
		case 3:
			applyprimitivereg(rm, (void *)&rm->val);
			gotoreg(rm, rm->cont); break;
		case 5:
		default:
			assert(valrawint(rm->val) == 2, "reg-machine: defvar, setvar, getvar");
			return;
		}
	}

	freeregmachine(rm);
}

// (lambda ()
//   (if 1 2 3))
//
// ->
//
// (make-compiled-procedure proc 1)
// (goto (label 2))
// 1
// (init-proc-env ())
// (initreg val (const 1))
// (branchval 4)
// 3
// (initreg val (const 2))
// (goto (reg continue))
// 4
// (initreg val (const 3))
// (goto (reg continue))
// 5
// 2
// (initargs)
// (branchproc 6)
// 7
// (goto proclabel)
// 6
// (apply-primitive-procedure val)
// (goto (reg continue))
// 8
// 
void test_if_true() {
	regmachine rm = mkregmachine();
	for (;;) {
		long labelval = valrawint(rm->label);
		switch (labelval) {
		case 0:
			mkcompiledprocreg(rm, (void *)&rm->proc, 1);
			gotolabel(rm, 2); break;
		case 1:
			initprocenv(rm, null);
			initreg((void *)&rm->val, numvali(1));
			if (branchval(rm, 4)) { break; }
		case 3:
			initreg((void *)&rm->val, numvali(2));
			gotoreg(rm, rm->cont); break;
		case 4:
			initreg((void *)&rm->val, numvali(3));
			gotoreg(rm, rm->cont); break;
		case 5:
		case 2:
			initargs(rm);
			if (branchproc(rm, 6)) { break; }
		case 7:
			gotoproc(rm); break;
		case 6:
			applyprimitivereg(rm, (void *)&rm->val);
			gotoreg(rm, rm->cont); break;
		case 8:
		default:
			assert(valrawint(rm->val) == 2, "reg-machine: if, true");
			return;
		}
	}

	freeregmachine(rm);
}

// ((lambda ()
//	(if false 2 3)))
//
// ->
//
void test_if_false() {
	regmachine rm = mkregmachine();
	for (;;) {
		long labelval = valrawint(rm->label);
		switch (labelval) {
		case 0:
		mkcompiledprocreg(rm, (void *)&rm->proc, 1);
			gotolabel(rm, 2); break;
		case 1:
			initprocenv(rm, null);
			getenvvar(rm, (void *)&rm->val, "false");
			if (branchval(rm, 4)) { break; }
		case 3:
			initreg((void *)&rm->val, numvali(2));
			gotoreg(rm, rm->cont); break;
		case 4:
			initreg((void *)&rm->val, numvali(3));
			gotoreg(rm, rm->cont); break;
		case 5:
		case 2:
			initargs(rm);
			if (branchproc(rm, 6)) { break; }
		case 7:
			gotoproc(rm); break;
		case 6:
			applyprimitivereg(rm, (void *)&rm->val);
			gotoreg(rm, rm->cont); break;
		case 8:
		default:
			assert(valrawint(rm->val) == 3, "reg-machine: if, false");
			return;
		}
	}

	freeregmachine(rm);
}

// ((lambda ()
//	(begin 1 2 3)))
//
// ->
//
// (make-compiled-procedure proc 1)
// (goto (label 2))
// 1
// (init-proc-env ())
// (initreg val (const 1))
// (initreg val (const 2))
// (initreg val (const 3))
// (goto (reg continue))
// 2
// (initargs)
// (branchproc 3)
// 4
// (goto proclabel)
// 3
// (apply-primitive-procedure val)
// (goto (reg continue))
// 5
// 
void test_begin() {
	regmachine rm = mkregmachine();
	for (;;) {
		long labelval = valrawint(rm->label);
		switch (labelval) {
		case 0:
			mkcompiledprocreg(rm, (void *)&rm->proc, 1);
			gotolabel(rm, 2); break;
		case 1:
			initprocenv(rm, null);
			initreg((void *)&rm->val, numvali(1));
			initreg((void *)&rm->val, numvali(2));
			initreg((void *)&rm->val, numvali(3));
			gotoreg(rm, rm->cont); break;
		case 2:
			initargs(rm);
			if (branchproc(rm, 3)) { break; }
		case 4:
			gotoproc(rm); break;
		case 3:
			applyprimitivereg(rm, (void *)&rm->val);
			gotoreg(rm, rm->cont); break;
		case 5:
		default:
			assert(valrawint(rm->val) == 3, "reg-machine: begin");
			return;
		}
	}

	freeregmachine(rm);
}

// ((lambda ()
//	(cond (false 2) (3 4) (else 5))))
//
// ->
//
// (make-compiled-procedure proc 1)
// (goto (label 2))
// 1
// (init-proc-env ())
// (get-variable val false)
// (branchval 4)
// 3
// (initreg val (const 2))
// (goto (reg continue))
// 4
// (initreg val (const 3))
// (branchval 7)
// 6
// (initreg val (const 4))
// (goto (reg continue))
// 7
// (initreg val (const 5))
// (goto (reg continue))
// 8
// 5
// 2
// (initargs)
// (branchproc 9)
// 10
// (goto proclabel)
// 9
// (apply-primitive-procedure val)
// (goto (reg continue))
// 11
// 
void test_cond_false() {
	regmachine rm = mkregmachine();
	for (;;) {
		long labelval = valrawint(rm->label);
		switch (labelval) {
		case 0:
			mkcompiledprocreg(rm, (void *)&rm->proc, 1);
			gotolabel(rm, 2); break;
		case 1:
			initprocenv(rm, null);
			getenvvar(rm, (void *)&rm->val, "false");
			if (branchval(rm, 4)) { break; }
		case 3:
			initreg((void *)&rm->val, numvali(2));
			gotoreg(rm, rm->cont); break;
		case 4:
			initreg((void *)&rm->val, numvali(3));
			if (branchval(rm, 7)) { break; }
		case 6:
			initreg((void *)&rm->val, numvali(4));
			gotoreg(rm, rm->cont); break;
		case 7:
			initreg((void *)&rm->val, numvali(5));
			gotoreg(rm, rm->cont); break;
		case 8:
		case 5:
		case 2:
			initargs(rm);
			if (branchproc(rm, 9)) { break; }
		case 10:
			gotoproc(rm); break;
		case 9:
			applyprimitivereg(rm, (void *)&rm->val);
			gotoreg(rm, rm->cont); break;
		case 11:
		default:
			assert(valrawint(rm->val) == 4, "reg-machine: cond, false");
			return;
		}
	}

	freeregmachine(rm);
}

// ((lambda ()
//	(cond (1 2) (3 4) (else 5))))
//
// ->
//
// (make-compiled-procedure proc 1)
// (goto (label 2))
// 1
// (init-proc-env ())
// (initreg val (const 1))
// (branchval 4)
// 3
// (initreg val (const 2))
// (goto (reg continue))
// 4
// (initreg val (const 3))
// (branchval 7)
// 6
// (initreg val (const 4))
// (goto (reg continue))
// 7
// (initreg val (const 5))
// (goto (reg continue))
// 8
// 5
// 2
// (initargs)
// (branchproc 9)
// 10
// (goto proclabel)
// 9
// (apply-primitive-procedure val)
// (goto (reg continue))
// 11
// 
void test_cond_true() {
	regmachine rm = mkregmachine();
	for (;;) {
		long labelval = valrawint(rm->label);
		switch (labelval) {
		case 0:
			mkcompiledprocreg(rm, (void *)&rm->proc, 1);
			gotolabel(rm, 2); break;
		case 1:
			initprocenv(rm, null);
			initreg((void *)&rm->val, numvali(1));
			if (branchval(rm, 4)) { break; }
		case 3:
			initreg((void *)&rm->val, numvali(2));
			gotoreg(rm, rm->cont); break;
		case 4:
			initreg((void *)&rm->val, numvali(3));
			if (branchval(rm, 7)) { break; }
		case 6:
			initreg((void *)&rm->val, numvali(4));
			gotoreg(rm, rm->cont); break;
		case 7:
			initreg((void *)&rm->val, numvali(5));
			gotoreg(rm, rm->cont); break;
		case 8:
		case 5:
		case 2:
			initargs(rm);
			if (branchproc(rm, 9)) { break; }
		case 10:
			gotoproc(rm); break;
		case 9:
			applyprimitivereg(rm, (void *)&rm->val);
			gotoreg(rm, rm->cont); break;
		case 11:
		default:
			assert(valrawint(rm->val) == 2, "reg-machine: cond, true");
			return;
		}
	}

	freeregmachine(rm);
}

// ((lambda ()
//	(cond (false 2) (false 4) (else 5))))
//
// ->
//
// (make-compiled-procedure proc 1)
// (goto (label 2))
// 1
// (init-proc-env ())
// (get-variable val false)
// (branchval 4)
// 3
// (initreg val (const 2))
// (goto (reg continue))
// 4
// (get-variable val false)
// (branchval 7)
// 6
// (initreg val (const 4))
// (goto (reg continue))
// 7
// (initreg val (const 5))
// (goto (reg continue))
// 8
// 5
// 2
// (initargs)
// (branchproc 9)
// 10
// (goto proclabel)
// 9
// (apply-primitive-procedure val)
// (goto (reg continue))
// 11
// 
void test_cond_else() {
	regmachine rm = mkregmachine();
	for (;;) {
		long labelval = valrawint(rm->label);
		switch (labelval) {
		case 0:
			mkcompiledprocreg(rm, (void *)&rm->proc, 1);
			gotolabel(rm, 2); break;
		case 1:
			initprocenv(rm, null);
			getenvvar(rm, (void *)&rm->val, "false");
			if (branchval(rm, 4)) { break; }
		case 3:
			initreg((void *)&rm->val, numvali(2));
			gotoreg(rm, rm->cont); break;
		case 4:
			getenvvar(rm, (void *)&rm->val, "false");
			if (branchval(rm, 7)) { break; }
		case 6:
			initreg((void *)&rm->val, numvali(4));
			gotoreg(rm, rm->cont); break;
		case 7:
			initreg((void *)&rm->val, numvali(5));
			gotoreg(rm, rm->cont); break;
		case 8:
		case 5:
		case 2:
			initargs(rm);
			if (branchproc(rm, 9)) { break; }
		case 10:
			gotoproc(rm); break;
		case 9:
			applyprimitivereg(rm, (void *)&rm->val);
			gotoreg(rm, rm->cont); break;
		case 11:
		default:
			assert(valrawint(rm->val) == 5, "reg-machine: cond, else");
			return;
		}
	}

	freeregmachine(rm);
}

// ((lamda ()
//	((lambda (x) x) 1)))
//
// ->
//
// (make-compiled-procedure proc 1)
// (goto (label 2))
// 1
// (init-proc-env ())
// (make-compiled-procedure proc 3)
// (goto (label 4))
// 3
// (init-proc-env (x))
// (get-variable val x)
// (goto (reg continue))
// 4
// (initreg val (const 1))
// (initargs)
// (addarg)
// (branchproc 5)
// 6
// (goto proclabel)
// 5
// (apply-primitive-procedure val)
// (goto (reg continue))
// 7
// 2
// (initargs)
// (branchproc 8)
// 9
// (goto proclabel)
// 8
// (apply-primitive-procedure val)
// (goto (reg continue))
// 10
// 
void test_call_lambda() {
	regmachine rm = mkregmachine();
	for (;;) {
		long labelval = valrawint(rm->label);
		switch (labelval) {
		case 0:
			mkcompiledprocreg(rm, (void *)&rm->proc, 1);
			gotolabel(rm, 2); break;
		case 1:
			initprocenv(rm, null);
			mkcompiledprocreg(rm, (void *)&rm->proc, 3);
			gotolabel(rm, 4); break;
		case 3:
			initprocenv(rm, pairval(symval("x"), null));
			getenvvar(rm, (void *)&rm->val, "x");
			gotoreg(rm, rm->cont); break;
		case 4:
			initreg((void *)&rm->val, numvali(1));
			initargs(rm);
			addarg(rm);
			if (branchproc(rm, 5)) { break; }
		case 6:
			gotoproc(rm); break;
		case 5:
			applyprimitivereg(rm, (void *)&rm->val);
			gotoreg(rm, rm->cont); break;
		case 7:
		case 2:
			initargs(rm);
			if (branchproc(rm, 8)) { break; }
		case 9:
			gotoproc(rm); break;
		case 8:
			applyprimitivereg(rm, (void *)&rm->val);
			gotoreg(rm, rm->cont); break;
		case 10:
		default:
			assert(valrawint(rm->val) == 1, "reg-machine: call lambda");
			return;
		}
	}

	freeregmachine(rm);
}

// ((lambda ()
//	(apply (lambda (x) x) '(1))))
//
// ->
//
// (make-compiled-procedure proc 1)
// (goto (label 2))
// 1
// (init-proc-env ())
// (get-variable proc apply)
// (initreg val (const (1)))
// (initargs)
// (addarg)
// (make-compiled-procedure val 3)
// (goto (label 4))
// 3
// (init-proc-env (x))
// (get-variable val x)
// (goto (reg continue))
// 4
// (addarg)
// (branchproc 5)
// 6
// (goto proclabel)
// 5
// (apply-primitive-procedure val)
// (goto (reg continue))
// 7
// 2
// (initargs)
// (branchproc 8)
// 9
// (goto proclabel)
// 8
// (apply-primitive-procedure val)
// (goto (reg continue))
// 10
// 
void test_apply_lambda() {
	regmachine rm = mkregmachine();
	for (;;) {
		long labelval = valrawint(rm->label);
		switch (labelval) {
		case 0:
			mkcompiledprocreg(rm, (void *)&rm->proc, 1);
			gotolabel(rm, 2); break;
		case 1:
			initprocenv(rm, null);
			getenvvar(rm, (void *)&rm->proc, "apply");
			initreg((void *)&rm->val, pairval(numvali(1), null));
			initargs(rm);
			addarg(rm);
			mkcompiledprocreg(rm, (void *)&rm->val, 3);
			gotolabel(rm, 4); break;
		case 3:
			initprocenv(rm, pairval(symval("x"), null));
			getenvvar(rm, (void *)&rm->val, "x");
			gotoreg(rm, rm->cont); break;
		case 4:
			addarg(rm);
			if (branchproc(rm, 5)) { break; }
		case 6:
			gotoproc(rm); break;
		case 5:
			applyprimitivereg(rm, (void *)&rm->val);
			gotoreg(rm, rm->cont); break;
		case 7:
		case 2:
			initargs(rm);
			if (branchproc(rm, 8)) { break; }
		case 9:
			gotoproc(rm); break;
		case 8:
			applyprimitivereg(rm, (void *)&rm->val);
			gotoreg(rm, rm->cont); break;
		case 10:
		default:
			assert(valrawint(rm->val) == 1, "reg-machine: apply lambda");
			return;
		}
	}

	freeregmachine(rm);
}

// ((lambda ()
//	(define (f x) x)
//	(f 1)))
//
// ->
//
// (make-compiled-procedure proc 1)
// (goto (label 2))
// 1
// (init-proc-env ())
// (make-compiled-procedure val 3)
// (goto (label 4))
// 3
// (init-proc-env (x))
// (get-variable val x)
// (goto (reg continue))
// 4
// (define-variable f)
// (get-variable proc f)
// (initreg val (const 1))
// (initargs)
// (addarg)
// (branchproc 5)
// 6
// (goto proclabel)
// 5
// (apply-primitive-procedure val)
// (goto (reg continue))
// 7
// 2
// (initargs)
// (branchproc 8)
// 9
// (goto proclabel)
// 8
// (apply-primitive-procedure val)
// (goto (reg continue))
// 10
// 
void test_call_func_var() {
	regmachine rm = mkregmachine();
	for (;;) {
		long labelval = valrawint(rm->label);
		switch (labelval) {
		case 0:
			mkcompiledprocreg(rm, (void *)&rm->proc, 1);
			gotolabel(rm, 2); break;
		case 1:
			initprocenv(rm, null);
			mkcompiledprocreg(rm, (void *)&rm->val, 3);
			gotolabel(rm, 4); break;
		case 3:
			initprocenv(rm, pairval(symval("x"), null));
			getenvvar(rm, (void *)&rm->val, "x");
			gotoreg(rm, rm->cont); break;
		case 4:
			defenvvar(rm, "f");
			getenvvar(rm, (void *)&rm->proc, "f");
			initreg((void *)&rm->val, numvali(1));
			initargs(rm);
			addarg(rm);
			if (branchproc(rm, 5)) { break; }
		case 6:
			gotoproc(rm); break;
		case 5:
			applyprimitivereg(rm, (void *)&rm->val);
			gotoreg(rm, rm->cont); break;
		case 7:
		case 2:
			initargs(rm);
			if (branchproc(rm, 8)) { break; }
		case 9:
			gotoproc(rm); break;
		case 8:
			applyprimitivereg(rm, (void *)&rm->val);
			gotoreg(rm, rm->cont); break;
		case 10:
		default:
			assert(valrawint(rm->val) == 1, "reg-machine: call func var");
			return;
		}
	}

	freeregmachine(rm);
}

// ((lambda ()
//	(define (f x) x)
//	(apply f '(1))))
//
// ->
//
// (make-compiled-procedure proc 1)
// (goto (label 2))
// 1
// (init-proc-env ())
// (make-compiled-procedure val 3)
// (goto (label 4))
// 3
// (init-proc-env (x))
// (get-variable val x)
// (goto (reg continue))
// 4
// (define-variable f)
// (get-variable proc apply)
// (initreg val (const (1)))
// (initargs)
// (addarg)
// (get-variable val f)
// (addarg)
// (branchproc 5)
// 6
// (goto proclabel)
// 5
// (apply-primitive-procedure val)
// (goto (reg continue))
// 7
// 2
// (initargs)
// (branchproc 8)
// 9
// (goto proclabel)
// 8
// (apply-primitive-procedure val)
// (goto (reg continue))
// 10
// 
void test_apply_func_var() {
	regmachine rm = mkregmachine();
	for (;;) {
		long labelval = valrawint(rm->label);
		switch (labelval) {
		case 0:
			mkcompiledprocreg(rm, (void *)&rm->proc, 1);
			gotolabel(rm, 2); break;
		case 1:
			initprocenv(rm, null);
			mkcompiledprocreg(rm, (void *)&rm->val, 3);
			gotolabel(rm, 4); break;
		case 3:
			initprocenv(rm, pairval(symval("x"), null));
			getenvvar(rm, (void *)&rm->val, "x");
			gotoreg(rm, rm->cont); break;
		case 4:
			defenvvar(rm, "f");
			getenvvar(rm, (void *)&rm->proc, "apply");
			initreg((void *)&rm->val, pairval(numvali(1), null));
			initargs(rm);
			addarg(rm);
			getenvvar(rm, (void *)&rm->val, "f");
			addarg(rm);
			if (branchproc(rm, 5)) { break; }
		case 6:
			gotoproc(rm); break;
		case 5:
			applyprimitivereg(rm, (void *)&rm->val);
			gotoreg(rm, rm->cont); break;
		case 7:
		case 2:
			initargs(rm);
			if (branchproc(rm, 8)) { break; }
		case 9:
			gotoproc(rm); break;
		case 8:
			applyprimitivereg(rm, (void *)&rm->val);
			gotoreg(rm, rm->cont); break;
		case 10:
		default:
			assert(valrawint(rm->val) == 1, "reg-machine: apply func var");
			return;
		}
	}

	freeregmachine(rm);
}

// ((lambda ()
//	(+ 1 2 3)))
//
// ->
//
// (make-compiled-procedure proc 1)
// (goto (label 2))
// 1
// (init-proc-env ())
// (get-variable proc +)
// (initreg val (const 3))
// (initargs)
// (addarg)
// (initreg val (const 2))
// (addarg)
// (initreg val (const 1))
// (addarg)
// (branchproc 3)
// 4
// (goto proclabel)
// 3
// (apply-primitive-procedure val)
// (goto (reg continue))
// 5
// 2
// (initargs)
// (branchproc 6)
// 7
// (goto proclabel)
// 6
// (apply-primitive-procedure val)
// (goto (reg continue))
// 8
// 
void test_exp_sum() {
	regmachine rm = mkregmachine();
	for (;;) {
		long labelval = valrawint(rm->label);
		switch (labelval) {
		case 0:
			mkcompiledprocreg(rm, (void *)&rm->proc, 1);
			gotolabel(rm, 2); break;
		case 1:
			initprocenv(rm, null);
			getenvvar(rm, (void *)&rm->proc, "+");
			initreg((void *)&rm->val, numvali(3));
			initargs(rm);
			addarg(rm);
			initreg((void *)&rm->val, numvali(2));
			addarg(rm);
			initreg((void *)&rm->val, numvali(1));
			addarg(rm);
			if (branchproc(rm, 3)) { break; }
		case 4:
			gotoproc(rm); break;
		case 3:
			applyprimitivereg(rm, (void *)&rm->val);
			gotoreg(rm, rm->cont); break;
		case 5:
		case 2:
			initargs(rm);
			if (branchproc(rm, 6)) { break; }
		case 7:
			gotoproc(rm); break;
		case 6:
			applyprimitivereg(rm, (void *)&rm->val);
			gotoreg(rm, rm->cont); break;
		case 8:
		default:
			assert(valrawint(rm->val) == 6, "reg-machine: exp sum");
			return;
		}
	}

	freeregmachine(rm);
}

// ((lambda ()
//	(define (sum . n)
//	  (apply + n))
//	(sum 1 2 3)))
// 
// ->
//
// (make-compiled-procedure proc 1)
// (goto (label 2))
// 1
// (init-proc-env ())
// (make-compiled-procedure val 3)
// (goto (label 4))
// 3
// (init-proc-env n)
// (get-variable proc apply)
// (get-variable val n)
// (initargs)
// (addarg)
// (get-variable val +)
// (addarg)
// (branchproc 5)
// 6
// (goto proclabel)
// 5
// (apply-primitive-procedure val)
// (goto (reg continue))
// 7
// 4
// (define-variable sum)
// (get-variable proc sum)
// (initreg val (const 3))
// (initargs)
// (addarg)
// (initreg val (const 2))
// (addarg)
// (initreg val (const 1))
// (addarg)
// (branchproc 8)
// 9
// (goto proclabel)
// 8
// (apply-primitive-procedure val)
// (goto (reg continue))
// 10
// 2
// (initargs)
// (branchproc 11)
// 12
// (goto proclabel)
// 11
// (apply-primitive-procedure val)
// (goto (reg continue))
// 13
void test_sum() {
	regmachine rm = mkregmachine();
	for (;;) {
		long labelval = valrawint(rm->label);
		switch (labelval) {
		case 0:
			mkcompiledprocreg(rm, (void *)&rm->proc, 1);
			gotolabel(rm, 2); break;
		case 1:
			initprocenv(rm, null);
			mkcompiledprocreg(rm, (void *)&rm->val, 3);
			gotolabel(rm, 4); break;
		case 3:
			initprocenv(rm, symval("n"));
			getenvvar(rm, (void *)&rm->proc, "apply");
			getenvvar(rm, (void *)&rm->val, "n");
			initargs(rm);
			addarg(rm);
			getenvvar(rm, (void *)&rm->val, "+");
			addarg(rm);
			if (branchproc(rm, 5)) { break; }
		case 6:
			gotoproc(rm); break;
		case 5:
			applyprimitivereg(rm, (void *)&rm->val);
			gotoreg(rm, rm->cont); break;
		case 7:
		case 4:
			defenvvar(rm, "sum");
			getenvvar(rm, (void *)&rm->proc, "sum");
			initreg((void *)&rm->val, numvali(3));
			initargs(rm);
			addarg(rm);
			initreg((void *)&rm->val, numvali(2));
			addarg(rm);
			initreg((void *)&rm->val, numvali(1));
			addarg(rm);
			if (branchproc(rm, 8)) { break; }
		case 9:
			gotoproc(rm); break;
		case 8:
			applyprimitivereg(rm, (void *)&rm->val);
			gotoreg(rm, rm->cont); break;
		case 10:
		case 2:
			initargs(rm);
			if (branchproc(rm, 11)) { break; }
		case 12:
			gotoproc(rm); break;
		case 11:
			applyprimitivereg(rm, (void *)&rm->val);
			gotoreg(rm, rm->cont); break;
		case 13:
		default:
			assert(valrawint(rm->val) == 6, "reg-machine: exp sum");
			return;
		}
	}

	freeregmachine(rm);
}

// ((lambda ()
//	(define f (file-open "some-file" file-mode-read))
//	(seek-file f 2 file-seek-mode-start)
//	(define s (read-file f 3))
//	(close-file s)
//	s))
//
// ->
//
// (make-compiled-procedure proc 1)
// (goto (label 2))
// 1
// (init-proc-env ())
// (save continue)
// (save env)
// (get-variable proc open-file)
// (get-variable val file-mode-read)
// (initargs)
// (addarg)
// (initreg val (const "some-file"))
// (addarg)
// (branchproc 3)
// 4
// (initreg continue (label 5))
// (takeproclabel)
// (goto (reg val))
// 3
// (apply-primitive-procedure val)
// 5
// (restore env)
// (define-variable f)
// (restore continue)
// (save continue)
// (save env)
// (get-variable proc seek-file)
// (get-variable val file-seek-mode-start)
// (initargs)
// (addarg)
// (initreg val (const 2))
// (addarg)
// (get-variable val f)
// (addarg)
// (branchproc 6)
// 7
// (initreg continue (label 8))
// (takeproclabel)
// (goto (reg val))
// 6
// (apply-primitive-procedure val)
// 8
// (restore env)
// (restore continue)
// (save continue)
// (save env)
// (get-variable proc read-file)
// (initreg val (const 3))
// (initargs)
// (addarg)
// (get-variable val f)
// (addarg)
// (branchproc 9)
// 10
// (initreg continue (label 11))
// (takeproclabel)
// (goto (reg val))
// 9
// (apply-primitive-procedure val)
// 11
// (restore env)
// (define-variable s)
// (restore continue)
// (save continue)
// (save env)
// (get-variable proc close-file)
// (get-variable val f)
// (initargs)
// (addarg)
// (branchproc 12)
// 13
// (initreg continue (label 14))
// (takeproclabel)
// (goto (reg val))
// 12
// (apply-primitive-procedure val)
// 14
// (restore env)
// (restore continue)
// (get-variable val s)
// (goto (reg continue))
// 2
// (initargs)
// (branchproc 15)
// 16
// (goto proclabel)
// 15
// (apply-primitive-procedure val)
// (goto (reg continue))
// 17
// 
void test_file_read() {
	initfilecontent("some-file", "Hello World!");
	regmachine rm = mkregmachine();
	for (;;) {
		long labelval = valrawint(rm->label);
		char *s;
		switch (labelval) {
		case 0:
			mkcompiledprocreg(rm, (void *)&rm->proc, 1);
			gotolabel(rm, 2); break;
		case 1:
			initprocenv(rm, null);
			savereg(rm, rm->cont);
			savereg(rm, rm->env);
			getenvvar(rm, (void *)&rm->proc, "open-file");
			getenvvar(rm, (void *)&rm->val, "file-mode-read");
			initargs(rm);
			addarg(rm);
			initreg((void *)&rm->val, stringval("some-file"));
			addarg(rm);
			if (branchproc(rm, 3)) { break; }
		case 4:
			initreg((void *)&rm->cont, numvali(5));
			takeproclabel(rm);
			gotoreg(rm, rm->val); break;
		case 3:
			applyprimitivereg(rm, (void *)&rm->val);
		case 5:
			restorereg(rm, (void *)&rm->env);
			defenvvar(rm, "f");
			restorereg(rm, (void *)&rm->cont);
			savereg(rm, rm->cont);
			savereg(rm, rm->env);
			getenvvar(rm, (void *)&rm->proc, "seek-file");
			getenvvar(rm, (void *)&rm->val, "file-seek-mode-start");
			initargs(rm);
			addarg(rm);
			initreg((void *)&rm->val, numvali(2));
			addarg(rm);
			getenvvar(rm, (void *)&rm->val, "f");
			addarg(rm);
			if (branchproc(rm, 6)) { break; }
		case 7:
			initreg((void *)&rm->cont, numvali(8));
			takeproclabel(rm);
			gotoreg(rm, rm->val); break;
		case 6:
			applyprimitivereg(rm, (void *)&rm->val);
		case 8:
			restorereg(rm, (void *)&rm->env);
			restorereg(rm, (void *)&rm->cont);
			savereg(rm, rm->cont);
			savereg(rm, rm->env);
			getenvvar(rm, (void *)&rm->proc, "read-file");
			initreg((void *)&rm->val, numvali(3));
			initargs(rm);
			addarg(rm);
			getenvvar(rm, (void *)&rm->val, "f");
			addarg(rm);
			if (branchproc(rm, 9)) { break; }
		case 10:
			initreg((void *)&rm->cont, numvali(11));
			takeproclabel(rm);
			gotoreg(rm, rm->val); break;
		case 9:
			applyprimitivereg(rm, (void *)&rm->val);
		case 11:
			restorereg(rm, (void *)&rm->env);
			defenvvar(rm, "s");
			restorereg(rm, (void *)&rm->cont);
			savereg(rm, rm->cont);
			savereg(rm, rm->env);
			getenvvar(rm, (void *)&rm->proc, "close-file");
			getenvvar(rm, (void *)&rm->val, "f");
			initargs(rm);
			addarg(rm);
			if (branchproc(rm, 12)) { break; }
		case 13:
			initreg((void *)&rm->cont, numvali(14));
			takeproclabel(rm);
			gotoreg(rm, rm->val); break;
		case 12:
			applyprimitivereg(rm, (void *)&rm->val);
		case 14:
			restorereg(rm, (void *)&rm->env);
			restorereg(rm, (void *)&rm->cont);
			getenvvar(rm, (void *)&rm->val, "s");
			gotoreg(rm, rm->cont); break;
		case 2:
			initargs(rm);
			if (branchproc(rm, 15)) { break; }
		case 16:
			gotoproc(rm); break;
		case 15:
			applyprimitivereg(rm, (void *)&rm->val);
			gotoreg(rm, rm->cont); break;
		case 17:
		default:
			s = valrawstring(rm->val);
			assert(!strcoll(s, "llo"), "reg-machine: read file");
			free(s);
			return;
		}
	}

	freeregmachine(rm);
}

// (lambda ()
//   ; ideally, this will become simpler with block escaped symbols
//   ; double escaping needs to be fixed after disattached from guile
//   (define rx (make-regex "\\\\([^(]*\\\\)" 0)) ")"
//   (regex-match rx "((some list) (of lists))"))
//
// ->
//
// (make-compiled-procedure proc 1)
// (goto (label 2))
// 1
// (init-proc-env ())
// (save continue)
// (save env)
// (get-variable proc make-regex)
// (initreg val (const 0))
// (initargs)
// (addarg)
// (initreg val (const "\\\\([^(]*\\\\)"))
// (addarg)
// (branchproc 3)
// 4
// (initreg continue (label 5))
// (takeproclabel)
// (goto (reg val))
// 3
// (apply-primitive-procedure val)
// 5
// (restore env)
// (define-variable rx)
// (restore continue)
// (initreg val (const ")"))
// (get-variable proc regex-match)
// (initreg val (const "((some list) (of lists))"))
// (initargs)
// (addarg)
// (get-variable val rx)
// (addarg)
// (branchproc 6)
// 7
// (goto proclabel)
// 6
// (apply-primitive-procedure val)
// (goto (reg continue))
// 8
// 2
// (initargs)
// (branchproc 9)
// 10
// (goto proclabel)
// 9
// (apply-primitive-procedure val)
// (goto (reg continue))
// 11
// 
void test_regex() {
	regmachine rm = mkregmachine();
	for (;;) {
		long labelval = valrawint(rm->label);
		char *s;
		switch (labelval) {
		case 0:
			mkcompiledprocreg(rm, (void *)&rm->proc, 1);
			gotolabel(rm, 2); break;
		case 1:
			initprocenv(rm, null);
			savereg(rm, rm->cont);
			savereg(rm, rm->env);
			getenvvar(rm, (void *)&rm->proc, "make-regex");
			initreg((void *)&rm->val, numvali(0));
			initargs(rm);
			addarg(rm);
			initreg((void *)&rm->val, stringval("\\([^(]*\\)"));
			addarg(rm);
			if (branchproc(rm, 3)) { break; }
		case 4:
			initreg((void *)&rm->cont, numvali(5));
			takeproclabel(rm);
			gotoreg(rm, rm->val); break;
		case 3:
			applyprimitivereg(rm, (void *)&rm->val);
		case 5:
			restorereg(rm, (void *)&rm->env);
			defenvvar(rm, "rx");
			restorereg(rm, (void *)&rm->cont);
			initreg((void *)&rm->val, stringval(")"));
			getenvvar(rm, (void *)&rm->proc, "regex-match");
			initreg((void *)&rm->val, stringval("((some list) (of lists))"));
			initargs(rm);
			addarg(rm);
			getenvvar(rm, (void *)&rm->val, "rx");
			addarg(rm);
			if (branchproc(rm, 6)) { break; }
		case 7:
			gotoproc(rm); break;
		case 6:
			applyprimitivereg(rm, (void *)&rm->val);
			gotoreg(rm, rm->cont); break;
		case 8:
		case 2:
			initargs(rm);
			if (branchproc(rm, 9)) { break; }
		case 10:
			gotoproc(rm); break;
		case 9:
			applyprimitivereg(rm, (void *)&rm->val);
			gotoreg(rm, rm->cont); break;
		case 11:
		default:
			s = sprintraw(rm->val);
			assert(!strcoll(s, "((1 11))"), "reg-machine: regex");
			return;
		}
	}

	freeregmachine(rm);
}

// (define code
//   '((lambda ()
//       (define (make-object)
//         (define (internal-method)
//           'object-ok)
//         (lambda () (internal-method)))
//       (define object (make-object))
//       (object))))
//
// ->
//
// (make-compiled-procedure proc 1)
// (goto (label 2))
// 1
// (init-proc-env ())
// (make-compiled-procedure val 3)
// (goto (label 4))
// 3
// (init-proc-env ())
// (make-compiled-procedure val 5)
// (goto (label 6))
// 5
// (init-proc-env ())
// (initreg val (const object-ok))
// (goto (reg continue))
// 6
// (define-variable internal-method)
// (make-compiled-procedure val 7)
// (goto (reg continue))
// 7
// (init-proc-env ())
// (get-variable proc internal-method)
// (initargs)
// (branchproc 9)
// 10
// (goto proclabel)
// 9
// (apply-primitive-procedure val)
// (goto (reg continue))
// 11
// 8
// 4
// (define-variable make-object)
// (save continue)
// (save env)
// (get-variable proc make-object)
// (initargs)
// (branchproc 12)
// 13
// (initreg continue (label 14))
// (takeproclabel)
// (goto (reg val))
// 12
// (apply-primitive-procedure val)
// 14
// (restore env)
// (define-variable object)
// (restore continue)
// (get-variable proc object)
// (initargs)
// (branchproc 15)
// 16
// (goto proclabel)
// 15
// (apply-primitive-procedure val)
// (goto (reg continue))
// 17
// 2
// (initargs)
// (branchproc 18)
// 19
// (goto proclabel)
// 18
// (apply-primitive-procedure val)
// (goto (reg continue))
// 20
// 
void test_scope() {
	regmachine rm = mkregmachine();
	char *result;
	for (;;) {
		long labelval = valrawint(rm->label);
		switch (labelval) {
		case 0:
			mkcompiledprocreg(rm, (void *)&rm->proc, 1);
			gotolabel(rm, 2); break;
		case 1:
			initprocenv(rm, null);
			mkcompiledprocreg(rm, (void *)&rm->val, 3);
			gotolabel(rm, 4); break;
		case 3:
			initprocenv(rm, null);
			mkcompiledprocreg(rm, (void *)&rm->val, 5);
			gotolabel(rm, 6); break;
		case 5:
			initprocenv(rm, null);
			initreg((void *)&rm->val, symval("object-ok"));
			gotoreg(rm, rm->cont); break;
		case 6:
			defenvvar(rm, "internal-method");
			mkcompiledprocreg(rm, (void *)&rm->val, 7);
			gotoreg(rm, rm->cont); break;
		case 7:
			initprocenv(rm, null);
			getenvvar(rm, (void *)&rm->proc, "internal-method");
			initargs(rm);
			if (branchproc(rm, 9)) { break; }
		case 10:
			takeproclabel(rm);
			gotoreg(rm, rm->val); break;
		case 9:
			applyprimitivereg(rm, (void *)&rm->val);
			gotoreg(rm, rm->cont); break;
		case 11:
		case 8:
		case 4:
			defenvvar(rm, "make-object");
			savereg(rm, rm->cont);
			savereg(rm, rm->env);
			getenvvar(rm, (void *)&rm->proc, "make-object");
			initargs(rm);
			if (branchproc(rm, 12)) { break; }
		case 13:
			initreg((void *)&rm->cont, numvali(14));
			takeproclabel(rm);
			gotoreg(rm, rm->val); break;
		case 12:
			applyprimitivereg(rm, (void *)&rm->val);
		case 14:
			restorereg(rm, (void *)&rm->env);
			defenvvar(rm, "object");
			restorereg(rm, (void *)&rm->cont);
			getenvvar(rm, (void *)&rm->proc, "object");
			initargs(rm);
			if (branchproc(rm, 15)) { break; }
		case 16:
			takeproclabel(rm);
			gotoreg(rm, rm->val); break;
		case 15:
			applyprimitivereg(rm, (void *)&rm->val);
			gotoreg(rm, rm->cont); break;
		case 17:
		case 2:
			initargs(rm);
			if (branchproc(rm, 18)) { break; }
		case 19:
			takeproclabel(rm);
			gotoreg(rm, rm->val); break;
		case 18:
			applyprimitivereg(rm, (void *)&rm->val);
			gotoreg(rm, rm->cont); break;
		case 20:
		default:
			result = sprintraw(rm->val);
			assert(!strcoll(result, "object-ok"), "reg-machine: scope");
			return;
		}
	}

	freeregmachine(rm);
}

// ((lambda () (or)))
//
// ->
//
// (make-compiled-procedure proc 1)
// (goto (label 2))
// 1
// (init-proc-env ())
// (get-variable val false)
// (goto (reg continue))
// 2
// (initargs)
// (branchproc 3)
// 4
// (takeproclabel)
// (goto (reg val))
// 3
// (apply-primitive-procedure val)
// (goto (reg continue))
// 5
// 
void test_or_empty() {
	regmachine rm = mkregmachine();
	for (;;) {
		long labelval = valrawint(rm->label);
		switch (labelval) {
		case 0:
			mkcompiledprocreg(rm, (void *)&rm->proc, 1);
			gotolabel(rm, 2); break;
		case 1:
			initprocenv(rm, null);
			getenvvar(rm, (void *)&rm->val, "false");
			gotoreg(rm, rm->cont); break;
		case 2:
			initargs(rm);
			if (branchproc(rm, 3)) { break; }
		case 4:
			takeproclabel(rm);
			gotoreg(rm, rm->val); break;
		case 3:
			applyprimitivereg(rm, (void *)&rm->val);
			gotoreg(rm, rm->cont); break;
		case 5:
		default:
			assert(rm->val == false, "reg-machine: or, empty");
			return;
		}
	}

	freeregmachine(rm);
}

// ((lambda () (or false)))
//
// ->
//
// (make-compiled-procedure proc 1)
// (goto (label 2))
// 1
// (init-proc-env ())
// (get-variable val false)
// (branchval 4)
// 3
// (get-variable val false)
// (goto (reg continue))
// 4
// (get-variable val false)
// (goto (reg continue))
// 5
// 2
// (initargs)
// (branchproc 6)
// 7
// (takeproclabel)
// (goto (reg val))
// 6
// (apply-primitive-procedure val)
// (goto (reg continue))
// 8
// 
void test_or_false() {
	regmachine rm = mkregmachine();
	for (;;) {
		long labelval = valrawint(rm->label);
		switch (labelval) {
		case 0:
			mkcompiledprocreg(rm, (void *)&rm->proc, 1);
			gotolabel(rm, 2); break;
		case 1:
			initprocenv(rm, null);
			getenvvar(rm, (void *)&rm->val, "false");
			if (branchval(rm, 4)) { break; }
		case 3:
			getenvvar(rm, (void *)&rm->val, "false");
			gotoreg(rm, rm->cont); break;
		case 4:
			getenvvar(rm, (void *)&rm->val, "false");
			gotoreg(rm, rm->cont); break;
		case 5:
		case 2:
			initargs(rm);
			if (branchproc(rm, 6)) { break; }
		case 7:
			takeproclabel(rm);
			gotoreg(rm, rm->val); break;
		case 6:
			applyprimitivereg(rm, (void *)&rm->val);
			gotoreg(rm, rm->cont); break;
		case 8:
		default:
			assert(rm->val == false, "reg-machine: or, false");
			return;
		}
	}

	freeregmachine(rm);
}

// ((lambda () (or false 1 2)))
//
// ->
//
// (make-compiled-procedure proc 1)
// (goto (label 2))
// 1
// (init-proc-env ())
// (get-variable val false)
// (branchval 4)
// 3
// (get-variable val false)
// (goto (reg continue))
// 4
// (initreg val (const 1))
// (branchval 7)
// 6
// (initreg val (const 1))
// (goto (reg continue))
// 7
// (initreg val (const 2))
// (branchval 10)
// 9
// (initreg val (const 2))
// (goto (reg continue))
// 10
// (get-variable val false)
// (goto (reg continue))
// 11
// 8
// 5
// 2
// (initargs)
// (branchproc 12)
// 13
// (takeproclabel)
// (goto (reg val))
// 12
// (apply-primitive-procedure val)
// (goto (reg continue))
// 14
// 
void test_or_true () {
	regmachine rm = mkregmachine();
	for (;;) {
		long labelval = valrawint(rm->label);
		switch (labelval) {
		case 0:
			mkcompiledprocreg(rm, (void *)&rm->proc, 1);
			gotolabel(rm, 2); break;
		case 1:
			initprocenv(rm, null);
			getenvvar(rm, (void *)&rm->val, "false");
			if (branchval(rm, 4)) { break; }
		case 3:
			getenvvar(rm, (void *)&rm->val, "false");
			gotoreg(rm, rm->cont); break;
		case 4:
			initreg((void *)&rm->val, numvali(1));
			if (branchval(rm, 7)) { break; }
		case 6:
			initreg((void *)&rm->val, numvali(1));
			gotoreg(rm, rm->cont); break;
		case 7:
			initreg((void *)&rm->val, numvali(2));
			if (branchval(rm, 10)) { break; }
		case 9:
			initreg((void *)&rm->val, numvali(2));
			gotoreg(rm, rm->cont); break;
		case 10:
			getenvvar(rm, (void *)&rm->val, "false");
			gotoreg(rm, rm->cont); break;
		case 11:
		case 8:
		case 5:
		case 2:
			initargs(rm);
			if (branchproc(rm, 12)) { break; }
		case 13:
			takeproclabel(rm);
			gotoreg(rm, rm->val); break;
		case 12:
			applyprimitivereg(rm, (void *)&rm->val);
			gotoreg(rm, rm->cont); break;
		case 14:
		default:
			assert(valrawint(rm->val) == 1, "reg-machine: or, true");
			return;
		}
	}

	freeregmachine(rm);
}

// ((lambda () (and)))
//
// ->
//
// (make-compiled-procedure proc 1)
// (goto (label 2))
// 1
// (init-proc-env ())
// (get-variable val true)
// (goto (reg continue))
// 2
// (initargs)
// (branchproc 3)
// 4
// (takeproclabel)
// (goto (reg val))
// 3
// (apply-primitive-procedure val)
// (goto (reg continue))
// 5
// 
void test_and_empty() {
	regmachine rm = mkregmachine();
	for (;;) {
		long labelval = valrawint(rm->label);
		switch (labelval) {
		case 0:
			mkcompiledprocreg(rm, (void *)&rm->proc, 1);
			gotolabel(rm, 2); break;
		case 1:
			initprocenv(rm, null);
			getenvvar(rm, (void *)&rm->val, "true");
			gotoreg(rm, rm->cont); break;
		case 2:
			initargs(rm);
			if (branchproc(rm, 3)) { break; }
		case 4:
			takeproclabel(rm);
			gotoreg(rm, rm->val); break;
		case 3:
			applyprimitivereg(rm, (void *)&rm->val);
			gotoreg(rm, rm->cont); break;
		case 5:
		default:
			assert(rm->val == true, "reg-machine: and, empty");
			return;
		}
	}

	freeregmachine(rm);
}

// ((lambda () (and 1 2 false)))
//
// ->
//
// (make-compiled-procedure proc 1)
// (goto (label 2))
// 1
// (init-proc-env ())
// (initreg val (const 1))
// (branchval 4)
// 3
// (initreg val (const 2))
// (branchval 7)
// 6
// (get-variable val false)
// (branchval 10)
// 9
// (get-variable val true)
// (goto (reg continue))
// 10
// (get-variable val false)
// (goto (reg continue))
// 11
// 7
// (get-variable val false)
// (goto (reg continue))
// 8
// 4
// (get-variable val false)
// (goto (reg continue))
// 5
// 2
// (initargs)
// (branchproc 12)
// 13
// (takeproclabel)
// (goto (reg val))
// 12
// (apply-primitive-procedure val)
// (goto (reg continue))
// 14
// 
void test_and_false() {
	regmachine rm = mkregmachine();
	for (;;) {
		long labelval = valrawint(rm->label);
		switch (labelval) {
		case 0:
			mkcompiledprocreg(rm, (void *)&rm->proc, 1);
			gotolabel(rm, 2); break;
		case 1:
			initprocenv(rm, null);
			initreg((void *)&rm->val, numvali(1));
			if (branchval(rm, 4)) { break; }
		case 3:
			initreg((void *)&rm->val, numvali(2));
			if (branchval(rm, 7)) { break; }
		case 6:
			getenvvar(rm, (void *)&rm->val, "false");
			if (branchval(rm, 10)) { break; }
		case 9:
			getenvvar(rm, (void *)&rm->val, "true");
			gotoreg(rm, rm->cont); break;
		case 10:
			getenvvar(rm, (void *)&rm->val, "false");
			gotoreg(rm, rm->cont); break;
		case 11:
		case 7:
			getenvvar(rm, (void *)&rm->val, "false");
			gotoreg(rm, rm->cont); break;
		case 8:
		case 4:
			getenvvar(rm, (void *)&rm->val, "false");
			gotoreg(rm, rm->cont); break;
		case 5:
		case 2:
			initargs(rm);
			if (branchproc(rm, 12)) { break; }
		case 13:
			takeproclabel(rm);
			gotoreg(rm, rm->val); break;
		case 12:
			applyprimitivereg(rm, (void *)&rm->val);
			gotoreg(rm, rm->cont); break;
		case 14:
		default:
			assert(rm->val == false, "reg-machine: and, false");
			return;
		}
	}

	freeregmachine(rm);
}

// ((lambda () (and 1 2 3)))
//
// ->
//
// (make-compiled-procedure proc 1)
// (goto (label 2))
// 1
// (init-proc-env ())
// (initreg val (const 1))
// (branchval 4)
// 3
// (initreg val (const 2))
// (branchval 7)
// 6
// (initreg val (const 3))
// (branchval 10)
// 9
// (get-variable val true)
// (goto (reg continue))
// 10
// (get-variable val false)
// (goto (reg continue))
// 11
// 7
// (get-variable val false)
// (goto (reg continue))
// 8
// 4
// (get-variable val false)
// (goto (reg continue))
// 5
// 2
// (initargs)
// (branchproc 12)
// 13
// (takeproclabel)
// (goto (reg val))
// 12
// (apply-primitive-procedure val)
// (goto (reg continue))
// 14
// 
void test_and_true() {
	regmachine rm = mkregmachine();
	for (;;) {
		long labelval = valrawint(rm->label);
		switch (labelval) {
		case 0:
			mkcompiledprocreg(rm, (void *)&rm->proc, 1);
			gotolabel(rm, 2); break;
		case 1:
			initprocenv(rm, null);
			initreg((void *)&rm->val, numvali(1));
			if (branchval(rm, 4)) { break; }
		case 3:
			initreg((void *)&rm->val, numvali(2));
			if (branchval(rm, 7)) { break; }
		case 6:
			initreg((void *)&rm->val, numvali(3));
			if (branchval(rm, 10)) { break; }
		case 9:
			getenvvar(rm, (void *)&rm->val, "true");
			gotoreg(rm, rm->cont); break;
		case 10:
			getenvvar(rm, (void *)&rm->val, "false");
			gotoreg(rm, rm->cont); break;
		case 11:
		case 7:
			getenvvar(rm, (void *)&rm->val, "false");
			gotoreg(rm, rm->cont); break;
		case 8:
		case 4:
			getenvvar(rm, (void *)&rm->val, "false");
			gotoreg(rm, rm->cont); break;
		case 5:
		case 2:
			initargs(rm);
			if (branchproc(rm, 12)) { break; }
		case 13:
			takeproclabel(rm);
			gotoreg(rm, rm->val); break;
		case 12:
			applyprimitivereg(rm, (void *)&rm->val);
			gotoreg(rm, rm->cont); break;
		case 14:
		default:
			assert(valrawint(rm->val) == 3, "reg-machine: and, not false");
			return;
		}
	}

	freeregmachine(rm);
}

int main(int argc, char **argv) {
	initsys();
	initmodule_sysio();
	initmodule_number();
	initmodule_value();
	initmodule_errormock();
	initmodule_primitives();

	test_number_small();
	test_number_rational();
	// test_number_big();
	test_string();
	test_null();
	test_symbol();
	test_pair();
	test_defvar_getvar();
	test_defvar_setvar_getvar();
	test_if_true();
	test_if_false();
	test_begin();
	test_cond_false();
	test_cond_true();
	test_cond_else();
	test_call_lambda();
	test_apply_lambda();
	test_call_func_var();
	test_apply_func_var();
	test_exp_sum();
	test_sum();
	test_file_read();
	test_regex();
	test_scope();
	test_or_empty();
	test_or_false();
	test_or_true();
	test_or_empty();
	test_or_false();
	test_or_true();

	freemodule_primitives();
	freemodule_errormock();
	freemodule_value();
	freemodule_number();
	freemodule_sysio();
	return 0;
}
