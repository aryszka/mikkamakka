#include <stdlib.h>
#include <string.h>
#include "error.h"
#include "symbol.h"
#include "number.h"
#include "string.h"
#include "compound-types.h"
#include "stack.h"
#include "io.h"
#include "value.h"
#include "regex.h"
#include "primitives.h"

#include <stdio.h>

typedef number (*numop)(number, number);
typedef int (*typecheck)(value);
typedef int (*eqcheck)(void *, void *);
typedef int (*comparer)(int);

static int symeqcheck(void *s1, void *s2) {
	return symeq((symbol)s1, (symbol)s2);
}

static int numeqcheck(void *n1, void *n2) {
	return comparenum((number)n1, (number)n2) == 0;
}

static int stringeqcheck(void *s1, void *s2) {
	return comparestring((string)s1, (string)s2) == 0;
}

static value *parseargs(long len, typecheck *checks, value args) {
	value *parsed = malloc(len * sizeof(value));
	for (long i = 0; i < len; i++) {
		if (isnulltype(args)) {
			error(invalidnumberofargs, "invalid number of args");
			free(parsed);
			return 0;
		}

		*(parsed + i) = carval(args);
		if (!(**(checks + i))(*(parsed + i))) {
			error(invalidtype, "wrong type");
			free(parsed);
			return 0;
		}

		args = cdrval(args);
	}

	return parsed;
}

static void parsevarargs(typecheck check, value args, long *len, value **parsed, int *haserr) {
	*haserr = 0;
	*len = 0;
	*parsed = 0;
	value *parsedprev = 0;
	value curr;
	for (;;) {
		if (isnulltype(args)) {
			return;
		}

		curr = carval(args);
		if (!(*check)(curr)) {
			error(invalidtype, "wrong type");

			if (*len) {
				free(*parsed);
			}

			*haserr = 1;
			return;
		}

		parsedprev = *parsed;
		*parsed = malloc((*len + 1) * sizeof(value));
		for (long i = 0; i < *len; i++) {
			*(*parsed + i) = *(parsedprev + i);
		}
		free(parsedprev);

		*(*parsed + *len) = curr;
		(*len)++;
		args = cdrval(args);
	}
}

static value callnumop(number defval, numop op, value args) {
	if (isnulltype(args)) {
		return mknumval(defval);
	}

	if (!ispairtype(args)) {
		error(invalidtype, "pair expected");
		return 0;
	}

	value n = carval(args);
	if (!isnumtype(n)) {
		error(invalidtype, "primitives: number expected, 1");
		return 0;
	}

	value rr = callnumop(defval, op, cdrval(args));
	number nr = (*op)(numval(n), numval(rr));
	// freeval(rr);
	value vr = mknumval(nr);
	freenum(nr);
	return vr;
}

static value callnumopneg(number defval, numop op, numop inverse, value args) {
	if (isnulltype(args)) {
		return mknumval(defval);
	}

	if (!ispairtype(args)) {
		error(invalidtype, "pair expected");
		return 0;
	}

	value a = carval(args);
	if (!isnumtype(a)) {
		error(invalidtype, "primitives: number expected, 2");
		return 0;
	}

	number an = numval(a);
	number r;
	value rv;

	value d = cdrval(args);
	if (isnulltype(d)) {
		r = (*op)(defval, an);
		rv = mknumval(r);
		freenum(r);
		return rv;
	}

	value rrv = callnumop(defval, inverse, d);
	r = (*op)(an, numval(rrv));
	// freeval(rrv);
	rv = mknumval(r);
	freenum(r);
	return rv;
}

static int npeq(value v1, value v2, typecheck typecheck, eqcheck eqcheck) {
	if (!(*typecheck)(v2)) {
		return 1;
	}

	return !eqcheck(rawval(v1), rawval(v2));
}

static value comparison(value args, comparer comp) {
	if (isnulltype(args)) {
		return true;
	}

	value a = carval(args);
	int isnum = isnumtype(a);
	int isstring = isstringtype(a);

	if (!isnum && !isstring) {
		error(invalidtype, "number or string expected");
		return 0;
	}

	value d = cdrval(args);
	if (isnulltype(d)) {
		return true;
	}

	value ad = carval(d);
	if ((isnum && !isnumtype(ad)) ||
		(isstring && !isstringtype(ad))) {
		error(invalidtype, "values of same type expected");
		return 0;
	}

	int c;
	if (isnum) {
		c = comparenum(numval(a), numval(ad));
	} else {
		c = comparestring(valstring(a), valstring(ad));
	}

	if (!comp(c)) {
		return false;
	}

	return comparison(d, comp);
}

static char *copystring(long len, char *s) {
	char *c = malloc(len);
	for (long i = 0; i < len; i++) {
		*(c + i) = *(s + i);
	}

	return c;
}

value errorval(value args) {
	value *argsparsed = parseargs(1, (typecheck[]){&isstringtype}, args);
	if (!argsparsed) {
		return 0;
	}

	value msg = *argsparsed;
	free(argsparsed);

	error(usererror, valrawstring(msg));
	return mksymval(5, "error");
}

value iseqval(value args) {
	if (isnulltype(args)) {
		return true;
	}

	if (!ispairtype(args)) {
		error(invalidtype, "pair expected");
		return 0;
	}

	value a = carval(args);
	value d = cdrval(args);

	if (isnulltype(d)) {
		return true;
	}

	if (!ispairtype(d)) {
		error(invalidtype, "pair expected");
		return 0;
	}

	value ad = carval(d);

	if (a == ad) {
		return iseqval(d);
	}

	if (issymtype(a)) {
		if (npeq(a, ad, &issymtype, &symeqcheck)) {
			return false;
		}
	} else if (isnumtype(a)) {
		if (npeq(a, ad, &isnumtype, &numeqcheck)) {
			return false;
		}
	} else if (isstringtype(a)) {
		if (npeq(a, ad, &isstringtype, &stringeqcheck)) {
			return false;
		}
	} else if (a != ad) {
		return false;
	}

	return iseqval(d);
}

value sumval(value args) {
	number defval = mknumi(0, 1);
	value r = callnumop(defval, &sum, args);
	freenum(defval);
	return r;
}

value diffval(value args) {
	number defval = mknumi(0, 1);
	value r = callnumopneg(defval, &diff, &sum, args);
	freenum(defval);
	return r;
}

value bitorval(value args) {
	number defval = mknumi(0, 1);
	value r = callnumop(defval, &bitor, args);
	freenum(defval);
	return r;
}

int isfiletype(value v) {
	return gettypecheck(v) == &isfiletype;
}

value mkfileval(file f) {
	return mkval(f, 0, 0, &isfiletype);
}

value openfileval(value args) {
	value *argsparsed = parseargs(2, (typecheck[]){&isstringtype, &issmallinttype}, args);
	if (!argsparsed) {
		return 0;
	}

	value fn = *argsparsed;
	value mode = *(argsparsed + 1);
	free(argsparsed);

	file f;
	ioerror err;
	openfile(valrawstring(fn), valrawint(mode), &f, &err);

	if (err.code) {
		// todo: error handling
		error(err.code, "file error");
		return 0;
	}

	return mkfileval(f);
}

value seekfileval(value args) {
	value *argsparsed = parseargs(3, (typecheck[]){&isfiletype, &issmallinttype, &issmallinttype}, args);
	if (!argsparsed) {
		return 0;
	}

	value f = *argsparsed;
	value pos = *(argsparsed + 1);
	value mode = *(argsparsed + 2);
	free(argsparsed);

	ioerror err = seekfile(rawval(f), valrawint(pos), valrawint(mode));
	if (err.code) {
		// todo: error handling
		error(err.code, "file error");
		return 0;
	}

	return mksymval(2, "ok");
}

value readfileval(value args) {
	value *argsparsed = parseargs(2, (typecheck[]){&isfiletype, &issmallinttype}, args);
	if (!argsparsed) {
		return 0;
	}

	value f = *argsparsed;
	value len = *(argsparsed + 1);
	free(argsparsed);

	char *s = malloc(valrawint(len));
	long rlen;
	ioerror err;
	readfile(rawval(f), valrawint(len), s, &rlen, &err);

	if (err.code && err.code != eof) {
		// todo: error handling
		error(err.code, "file error");
		free(s);
		return 0;
	}

	if (err.code && err.code == eof && !rlen) {
		free(s);
		return eofval;
	}

	value r = mkstringvalc(rlen, s);
	free(s);

	return r;
}

value writefileval(value args) {
	// this should accept only strings
	value *argsparsed = parseargs(2, (typecheck[]){&isfiletype, &typecheckany}, args);
	if (!argsparsed) {
		return 0;
	}

	value f = *argsparsed;
	value v = *(argsparsed + 1);
	free(argsparsed);

	char *s = sprintraw(v);
	ioerror err = writefile(rawval(f), strlen(s), s);
	free(s);

	if (err.code) {
		error(err.code, "file error");
		return 0;
	}

	return mksymval(2, "ok");
}

value closefileval(value args) {
	value *argsparsed = parseargs(1, (typecheck[]){&isfiletype}, args);
	if (!argsparsed) {
		return 0;
	}

	value fv = *argsparsed;
	free(argsparsed);

	file f = rawval(fv);
	// freeval(fv);
	ioerror err = closefile(f);
	if (err.code) {
		// todo: error handling
		error(err.code, "file error");
		return 0;
	}

	return mksymval(2, "ok");
}

int isregextype(value v) {
	return gettypecheck(v) == &isregextype;
}

void freeregexraw(void *rx) {
	freeregex(rx);
}

value mkregexval(value args) {
	value *argsparsed = parseargs(2, (typecheck[]){&isstringtype, &isnumtype}, args);
	if (!argsparsed) {
		return 0;
	}

	value expression = *argsparsed;
	value flags = *(argsparsed + 1);
	free(argsparsed);

	char *rexp = valrawstring(expression);
	regex rx = mkregex(strlen(rexp), rexp, valrawint(flags));
	return mkval(rx, &freeregexraw, 0, &isregextype);
}

value regexmatch(value args) {
	value *argsparsed = parseargs(2, (typecheck[]){&isregextype, &isstringtype}, args);
	if (!argsparsed) {
		return 0;
	}

	value rxv = *argsparsed;
	regex rx = rawval(rxv);

	value sv = *(argsparsed + 1);
	char *s = valrawstring(sv);
	long l = strlen(s);

	free(argsparsed);

	match m = matchrx(rx, l, s);
	if (!m) {
		return null;
	}

	long len = matchlen(m);

	value r = null;
	for (long i = len - 1; i >= 0; i--) {
		submatch sm = smatch(m, i);
		if (sm.index < 0) {
			continue;
		}

		r = mkpairval(
			mkpairval(
				mknumvali(bytestochars(sm.index, s), 1),
				mkpairval(
					mknumvali(bytestochars(sm.len, s + sm.index), 1),
					null
				)
			),
			r
		);
	}

	return r;
}

value isutf8val(value args) {
	value *argsparsed = parseargs(1, (typecheck[]){&isstringtype}, args);
	if (!argsparsed) {
		return 0;
	}

	string s = valstring(*argsparsed);
	free(argsparsed);

	return isutf8(s) ? true : false;
}

value copystrval(value args) {
	// todo: accept negative number as last argument to copy to end
	value *argsparsed = parseargs(3, (typecheck[]){&isstringtype, &issmallinttype, &issmallinttype}, args);
	if (!argsparsed) {
		return 0;
	}

	string s = valstring(*argsparsed);
	long from = valrawint(*(argsparsed + 1));
	long len = valrawint(*(argsparsed + 2));
	free(argsparsed);

	string ss = substr(s, from, len);
	value sv = mkstringval(ss);
	freestring(ss);
	return sv;
}

value byteslenval(value args) {
	value *argsparsed = parseargs(1, (typecheck[]){&isstringtype}, args);
	if (!argsparsed) {
		return 0;
	}

	string s = valstring(*argsparsed);
	free(argsparsed);

	return mknumvali(byteslen(s), 1);
}

value stringappendval(value args) {
	int haserr;
	long len;
	value *parsed;
	parsevarargs(&isstringtype, args, &len, &parsed, &haserr);
	if (haserr) {
		return 0;
	}

	string *s = malloc(len * sizeof(string));
	for (long i = 0; i < len; i++) {
		*(s + i) = valstring(*(parsed + i));
	}
	free(parsed);

	string as = appendstr(len, s);
	free(s);
	value asv = mkstringval(as);
	freestring(as);
	return asv;
}

value stringlenval(value args) {
	value *argsparsed = parseargs(1, (typecheck[]){&isstringtype}, args);
	if (!argsparsed) {
		return 0;
	}

	string s = valstring(*argsparsed);
	free(argsparsed);

	return mknumvali(stringlen(s), 1);
}

int iseoftype(value v) {
	return gettypecheck(v) == &iseoftype;
}

value iseofval(value args) {
	value *argsparsed = parseargs(1, (typecheck[]){&typecheckany}, args);
	if (!argsparsed) {
		return 0;
	}

	value v = *argsparsed;
	free(argsparsed);

	if (iseoftype(v)) {
		return true;
	}

	return false;
}

int compless(int c) {
	return c < 0;
}

value islessval(value args) {
	return comparison(args, &compless);
}

int compgreater(int c) {
	return c > 0;
}

value isgreaterval(value args) {
	return comparison(args, &compgreater);
}

int complesseq(int c) {
	return c <= 0;
}

value islessoreqval(value args) {
	return comparison(args, &complesseq);
}

int compgreatereq(int c) {
	return c >= 0;
}

value isgreateroreqval(value args) {
	return comparison(args, &compgreatereq);
}

value notval(value args) {
	value *argsparsed = parseargs(1, (typecheck[]){&typecheckany}, args);
	if (!argsparsed) {
		return 0;
	}

	value v = *argsparsed;
	free(argsparsed);

	if (isfalsetype(v)) {
		return true;
	}

	return false;
}

value isnullvalp(value args) {
	value *argsparsed = parseargs(1, (typecheck[]){&typecheckany}, args);
	if (!argsparsed) {
		return 0;
	}

	value v = *argsparsed;
	free(argsparsed);

	if (isnulltype(v)) {
		return true;
	}

	return false;
}

value consval(value args) {
	value *argsparsed = parseargs(2, (typecheck[]){&typecheckany, &typecheckany}, args);
	if (!argsparsed) {
		return 0;
	}

	value a = *argsparsed;
	value b = *(argsparsed + 1);
	free(argsparsed);

	return mkpairval(a, b);
}

value carvalp(value args) {
	value *argsparsed = parseargs(1, (typecheck[]){&ispairtype}, args);
	if (!argsparsed) {
		return 0;
	}

	value p = *argsparsed;
	free(argsparsed);

	return carval(p);
}

value cdrvalp(value args) {
	value *argsparsed = parseargs(1, (typecheck[]){&ispairtype}, args);
	if (!argsparsed) {
		return 0;
	}

	value p = *argsparsed;
	free(argsparsed);

	return cdrval(p);
}

value stringtonumsafe(value args) {
	value *argsparsed = parseargs(1, (typecheck[]){&isstringtype}, args);
	if (!argsparsed) {
		return 0;
	}

	value s = *argsparsed;
	free(argsparsed);

	number n = mknumcsafe(byteslen(valstring(s)), valrawstring(s));
	if (n) {
		value v = mknumval(n);
		freenum(n);
		return v;
	}

	return false;
}

value isnumvalp(value args) {
    value *argsparsed = parseargs(1, (typecheck[]){&typecheckany}, args);
    if (!argsparsed) {
        return 0;
    }

    value v = *argsparsed;
    free(argsparsed);

    return isnumval(v);
}

value isstringvalp(value args) {
    value *argsparsed = parseargs(1, (typecheck[]){&typecheckany}, args);
    if (!argsparsed) {
        return 0;
    }

    value v = *argsparsed;
    free(argsparsed);

    return isstringval(v);
}

value ispairvalp(value args) {
    value *argsparsed = parseargs(1, (typecheck[]){&typecheckany}, args);
    if (!argsparsed) {
        return 0;
    }

    value v = *argsparsed;
    free(argsparsed);

    return ispairval(v);
}

value issymbolvalp(value args) {
    value *argsparsed = parseargs(1, (typecheck[]){&typecheckany}, args);
    if (!argsparsed) {
        return 0;
    }

    value v = *argsparsed;
    free(argsparsed);

    return issymval(v);
}

value isintvalp(value args) {
    value *argsparsed = parseargs(1, (typecheck[]){&typecheckany}, args);
    if (!argsparsed) {
        return 0;
    }

    value v = *argsparsed;
    free(argsparsed);

    return isintval(v);
}

value numbertostring(value args) {
    value *argsparsed = parseargs(1, (typecheck[]){&isnumtype}, args);
    if (!argsparsed) {
        return 0;
    }

    value v = *argsparsed;
    free(argsparsed);

    return sprint(v);
}

value stringtosymbol(value args) {
    value *argsparsed = parseargs(1, (typecheck[]){&isstringtype}, args);
    if (!argsparsed) {
        return 0;
    }

    value s = *argsparsed;
    free(argsparsed);

    char *raw = valrawstring(s);
    return mksymval(strlen(raw), raw);
}

value symboltostring(value args) {
    value *argsparsed = parseargs(1, (typecheck[]){&issymtype}, args);
    if (!argsparsed) {
        return 0;
    }

    value s = *argsparsed;
    free(argsparsed);

    char *raw = sprintraw(s);
    return mkstringvalc(strlen(raw), raw);
}
char *sprinteof(void *raw) {
	return copystring(5, "<eof>");
}

char *sprintstdin(void *raw) {
	char *s = "<stdin>";
	return copystring(strlen(s), s);
}

char *sprintstdout(void *raw) {
	char *s = "<stdout>";
	return copystring(strlen(s), s);
}

char *sprintstderr(void *raw) {
	char *s = "<stderr>";
	return copystring(strlen(s), s);
}

void initmodule_primitives() {
	eofval = mkval(0, &freenoop, &sprinteof, &iseoftype);
	stdinval = mkval(iostdin, &freenoop, &sprintstdin, &isfiletype);
	stdoutval = mkval(iostdout, &freenoop, &sprintstdout, &isfiletype);
	stderrval = mkval(iostderr, &freenoop, &sprintstderr, &isfiletype);
}

void freemodule_primitives() {
	// freeval(eofval);
}
