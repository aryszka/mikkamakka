#include <stdlib.h>
#include <string.h>
#include <stdio.h>
#include "error.h"
#include "number.h"
#include "string.h"
#include "compound-types.h"
#include "stack.h"
#include "io.h"
#include "value.h"
#include "environment.h"
#include "register-machine.h"
#include "primitives.h"
#include "regex.h"
#include "registry.h"

struct rminternals {
	value apply;
};

static void defvarval(environment env, char *name, value v) {
	defvar(env, symval(name), v);
}

static void defvarproc(environment env, char *name, primitive p) {
	defvarval(env, name, mkprimitiveprocval(p));
}

static void defvarnum(environment env, char *name, long n) {
	defvarval(env, name, numvali(n));
}

static environment initialenv(value apply) {
	environment env = mkenvironment(0);
	defvarproc(env, "+", &sumval);
	defvarproc(env, "-", &diffval);
	defvarval(env, "apply", apply);
	defvarval(env, "false", false);
	defvarval(env, "true", true);
	defvarnum(env, "file-mode-read", ioread);
	defvarnum(env, "file-mode-read-write", ioreadwrite);
	defvarnum(env, "file-mode-write", iowrite);
	defvarnum(env, "file-mode-read-write-create", ioreadwritecreate);
	defvarnum(env, "file-mode-append", ioappend);
	defvarnum(env, "file-mode-read-append", ioreadappend);
	defvarnum(env, "file-seek-mode-start", iostart);
	defvarnum(env, "file-seek-mode-current", iocurr);
	defvarnum(env, "file-seek-mode-end", ioend);
	defvarproc(env, "open-file", &openfileval);
	defvarproc(env, "seek-file", &seekfileval);
	defvarproc(env, "read-file", &readfileval);
	defvarproc(env, "write-file", &writefileval);
	defvarproc(env, "close-file", &closefileval);
	defvarnum(env, "rx-ignore-case", rxignorecase);
	defvarnum(env, "rx-multiline", rxmultiline);
	defvarproc(env, "make-regex", &mkregexval);
	defvarproc(env, "regex-match", &regexmatch);
	defvarproc(env, "error", &errorval);
	defvarproc(env, "utf8-string?", &isutf8val);
	defvarproc(env, "string-copy", &copystrval);
	defvarproc(env, "bytes-length", &byteslenval);
	defvarproc(env, "eof?", &iseofval);
	defvarval(env, "eof", eofval);
	defvarproc(env, "string-append", &stringappendval);
	defvarproc(env, "string-length", &stringlenval);
	defvarproc(env, "<", &islessval);
	defvarproc(env, ">", &isgreaterval);
	defvarproc(env, "<=", &islessoreqval);
	defvarproc(env, ">=", &isgreateroreqval);
	defvarproc(env, "eq?", &iseqval);
	defvarproc(env, "==", &iseqval);
	defvarproc(env, "not", &notval);
	defvarproc(env, "null?", &isnullvalp);
	defvarproc(env, "cons", &consval);
	defvarproc(env, "car", &carvalp);
	defvarproc(env, "cdr", &cdrvalp);
	defvarval(env, "stdin", stdinval);
	defvarval(env, "stdout", stdoutval);
	defvarval(env, "stderr", stderrval);
	defvarproc(env, "string->number-safe", &stringtonumsafe);
    defvarproc(env, "number?", &isnumvalp);
    defvarproc(env, "string?", &isstringvalp);
    defvarproc(env, "pair?", &ispairvalp);
    defvarproc(env, "symbol?", &issymbolvalp);
    defvarproc(env, "integer?", &isintvalp);
    defvarproc(env, "number->string", &numbertostring);
    defvarproc(env, "string->symbol", &stringtosymbol);
    defvarproc(env, "symbol->string", &symboltostring);
	return env;
}

regmachine mkregmachine() {
	value apply = mkprimitiveprocval(&valapply);
	environment env = initialenv(apply);
	regmachine rm = malloc(sizeof(struct regmachine));
	rm->internals = malloc(sizeof(struct rminternals));
	rm->internals->apply = apply;
	rm->flag = 0;
	rm->label = numvali(0);
	rm->cont = numvali(-1);
	rm->proc = 0;
	rm->args = null;
	rm->val = symval("initial");
	rm->env = env;
	rm->stack = mkstack();
	rm->envstack = mkstack();
    setregistermachine(rm);
	return rm;
}

value symval(char *name) {
	return mksymval(strlen(name), name);
}

value numvali(long n) {
	value nv = mknumvali(n, 1);
    return nv;
}

value numvalc(char *c) {
	return mknumvalc(strlen(c), c);
}

value stringval(char *c) {
	return mkstringvalc(strlen(c), c);
}

value pairval(value car, value cdr) {
	return mkpairval(car, cdr);
}

void saveenv(regmachine rm) {
	push(rm->envstack, rm->env);
}

void restoreenv(regmachine rm) {
	rm->env = pop(rm->envstack);
}

void savereg(regmachine rm, void *reg) {
	push(rm->stack, reg);
}

void restorereg(regmachine rm, void **reg) {
	*reg = pop(rm->stack);
}

void gotoreg(regmachine rm, value reg) {
	rm->label = reg;
    collect();
}

void gotolabel(regmachine rm, long label) {
	rm->label = numvali(label);
    collect();
}

void gotoproc(regmachine rm) {
	rm->label = valproclabel(rm->proc);
    collect();
}

void takeproclabel(regmachine rm) {
	rm->val = valproclabel(rm->proc);
}

void initreg(void **reg, value val) {
	*reg = val;
}

void getenvvar(regmachine rm, void **reg, char *name) {
	*reg = getvar(rm->env, symval(name));
}

void setenvvar(regmachine rm, void **reg, char *name) {
	setvar(rm->env, symval(name), rm->val);
	*reg = symval("ok");
}

void defenvvar(regmachine rm, char *name) {
	defvar(rm->env, symval(name), rm->val);
}

int branchval(regmachine rm, long label) {
	if (isfalsetype(rm->val)) {
		rm->label = numvali(label);
		return 1;
	}

	return 0;
}

int branchproc(regmachine rm, long label) {
	if (rm->proc == rm->internals->apply) {
		rm->proc = carval(rm->args);
		rm->args = carval(cdrval(rm->args));
		return branchproc(rm, label);
	}

	if (isprimitiveproctype(rm->proc)) {
		rm->label = numvali(label);
		return 1;
	}

	return 0;
}

void initprocenv(regmachine rm, value names) {
	rm->env = extenv(valprocenv(rm->proc), names, rm->args);
}

void mkcompiledprocreg(regmachine rm, void **reg, long label) {
    value l = numvali(label);
	*reg = mkcompiledprocval(l, rm->env);
}

void initargs(regmachine rm) {
	rm->args = null;
}

void addarg(regmachine rm) {
	rm->args = pairval(rm->val, rm->args);
}

void applyprimitivereg(regmachine rm, void **reg) {
	*reg = valapplyprimitive(rm->proc, rm->args);
}

void freeregmachine(regmachine rm) {
	free(rm->internals);
	free(rm);
}
