#include <stdlib.h>
#include <string.h>
#include "error.h"
#include "symbol.h"
#include "number.h"
#include "string.h"
#include "compound-types.h"
#include "stack.h"
#include "procedure.h"
#include "environment.h"
#include "io.h"
#include "value.h"
#include "pair.h"
#include "sprint-list.h"
#include "registry.h"

#include <stdio.h>

typedef enum {
    vundefined,
	vsymbol,
	vnumber,
	vstring,
	vfalse,
	vtrue,
	vnull,
	vpair,
	vproc,
	vgeneric
} vtype;

struct value {
	vtype type;
	void *raw;
	freeraw freeraw;
	sprintrawfunc sprintraw;
	istype typecheck;
    registryflag rflag;
};

static value mkvalt(vtype t, void *raw, freeraw freeraw, sprintrawfunc sprintrawf, istype typecheck) {
	value v = malloc(sizeof(struct value));
	v->type = t;
	v->raw = raw;
	v->freeraw = freeraw;
	v->sprintraw = sprintrawf;
	v->typecheck = typecheck;
    v->rflag = unregistered;
    registerval(v);
	return v;
}

static char *copystring(long len, char *s) {
	char *c = malloc(len);
	for (long i = 0; i < len; i++) {
		*(c + i) = *(s + i);
	}

	return c;
}

void freenoop(void *raw) {}

char *sprintrawdefault(void *raw) {
	error(invalidtype, "not implemented, 2");
	return 0;
}

int typecheckany(value v) {
	return 1;
}

value mkval(void *raw, freeraw freeraw, sprintrawfunc sprintrawf, istype typecheck) {
	return mkvalt(vgeneric, raw, freeraw, sprintrawf, typecheck);
}

value mksymval(int len, char *c) {
	return mkvalt(vsymbol, mksymbol(len, c), (freeraw)&freesymbol, (sprintrawfunc)&sprintsymbol, &issymtype);
}

value mknumval(number n) {
	return mkvalt(vnumber, clonenum(n), (freeraw)&freenum, (sprintrawfunc)&sprintnum, &isinttype);
}

value mknumvalc(int len, char *c) {
	return mkvalt(vnumber, mknumc(len, c), (freeraw)&freenum, (sprintrawfunc)&sprintnum, &isinttype);
}

value mknumvali(long nom, long den) {
	value n = mkvalt(vnumber, mknumi(nom, den), (freeraw)&freenum, (sprintrawfunc)&sprintnum, &isinttype);
    return n;
}

value mkstringval(string s) {
	return mkvalt(vstring, clonestring(s), (freeraw)&freestring, (sprintrawfunc)&sprintstring, &isstringtype);
}

value mkstringvalc(int len, char *c) {
	return mkvalt(vstring, mkstring(len, c), (freeraw)&freestring, (sprintrawfunc)&sprintstring, &isstringtype);
}

value mkpairval(value car, value cdr) {
	return mkvalt(vpair, mkpair(car, cdr), (freeraw)&freepair, &sprintrawdefault, &ispairtype);
}

char *sprintprimitiveprocedure(void *raw) {
	char *s = "<primitive-procedure>";
	return copystring(strlen(s), s);
}

value mkprimitiveprocval(primitive p) {
	return mkvalt(vproc, mkprimitiveproc(p), (freeraw)&freeproc, &sprintrawdefault, &isproctype);
}

char *sprintcompiledprocedure(void *raw) {
	char *s = "<compiled-procedure>";
	return copystring(strlen(s), s);
}

value mkcompiledprocval(value label, environment env) {
	return mkvalt(vproc, mkcompiledproc(label, env), (freeraw)&freeproc, &sprintcompiledprocedure, &isproctype);
}

void *rawval(value v) {
	return v->raw;
}

istype gettypecheck(value v) {
	return v->typecheck;
}

value istypeval(value v, istype istype) {
	if ((*istype)(v)) {
		return true;
	}

	return false;
}

int issymtype(value v)           { return v->type == vsymbol; }
int isnumtype(value v)           { return v->type == vnumber; }
int isinttype(value v)           { return v->type == vnumber && isint(v->raw); }
int issmallinttype(value v)      { return v->type == vnumber && issmallint(v->raw); }
int isstringtype(value v)        { return v->type == vstring; }
int isfalsetype(value v)         { return v->type == vfalse; }
int isnulltype(value v)          { return v->type == vnull; }
int ispairtype(value v)          { return v->type == vpair; }
int isproctype(value v)          { return v->type == vproc; }
int isprimitiveproctype(value v) { return v->type == vproc && isprimitive(v->raw); }

value issymval(value v)           { return istypeval(v, &issymtype); }
value isnumval(value v)           { return istypeval(v, &isnumtype); }
value isintval(value v)           { return istypeval(v, &isinttype); }
value issmallintval(value v)      { return istypeval(v, &issmallinttype); }
value isstringval(value v)        { return istypeval(v, &isstringtype); }
value isfalseval(value v)         { return istypeval(v, &isfalsetype); }
value isnullval(value v)          { return istypeval(v, &isnulltype); }
value ispairval(value v)          { return istypeval(v, &ispairtype); }
value isprocval(value v)          { return istypeval(v, &isproctype); }
value isprimitiveprocval(value v) { return istypeval(v, &isprimitiveproctype); }

char *valsymbolnameraw(value s) {
	if (!issymtype(s)) {
		error(invalidtype, "value: symbol expected");
		return 0;
	}

	return symbolname(s->raw);
}

int valsymbollen(value s) {
	if (!issymtype(s)) {
		error(invalidtype, "value: symbol expected");
		return 0;
	}

	return symbollen(s->raw);
}

number numval(value n) {
	if (!isnumtype(n)) {
		error(invalidtype, "value: number expected, 1");
		return 0;
	}

	return n->raw;
}

long valrawint(value v) {
	if (!isnumtype(v)) {
		error(invalidtype, "value: number expected, 2");
		return 0;
	}

	return rawint(v->raw);
}

string valstring(value v) {
	if (!isstringtype(v)) {
		error(invalidtype, "string expected");
		return 0;
	}

	return v->raw;
}

char *valrawstring(value v) {
	if (!isstringtype(v)) {
		error(invalidtype, "string expected");
		return 0;
	}

	return rawstring(v->raw);
}

value carval(value p) {
	if (!ispairtype(p)) {
		error(invalidtype, "not a pair");
		return 0;
	}

	return car(p->raw);
}

value cdrval(value p) {
	if (!ispairtype(p)) {
		error(invalidtype, "not a pair");
		return 0;
	}

	return cdr(p->raw);
}

value valapplyprimitive(value p, value args) {
	if (!isproctype(p)) {
		error(invalidtype, "not a procedure");
		return 0;
	}

	return applyprimitive(p->raw, args);
}

value valproclabel(value p) {
	if (!isproctype(p)) {
		error(invalidtype, "not a procedure");
		return 0;
	}

	return proclabel(p->raw);
}

value valapply(value args) {
	return 0;
}

environment valprocenv(value p) {
	if (!isproctype(p)) {
		error(invalidtype, "not a procedure");
		return 0;
	}

	return procenv(p->raw);
}

registryflag getregistryflag(value v) {
    return v->rflag;
}

void setregistryflag(value v, registryflag f) {
    v->rflag = f;
}

char *sprintraw(value v) {
	switch (v->type) {
	case vfalse:
		return copystring(5, "false");
	case vtrue:
		return copystring(4, "true");
	case vnull:
		return copystring(2, "()");
	case vpair:
		return sprintlistraw(v);
	default:
		if (!v->sprintraw) {
			error(invalidtype, "not implemented, 1");
			return 0;
		}

		return (*v->sprintraw)(v->raw);
	}

}

value sprint(value v) {
    char *raw = sprintraw(v);
    value s = mkstringvalc(strlen(raw), raw);
    free(raw);
    return s;
}

void freeval(value v) {
	if (v->freeraw) {
		(*v->freeraw)(v->raw);
	}

	free(v);
}

void initmodule_value() {
	false = mkvalt(vfalse, 0, &freenoop, &sprintrawdefault, &isfalsetype);
	true = mkvalt(vtrue, 0, &freenoop, &sprintrawdefault, 0);
	null = mkvalt(vnull, 0, &freenoop, &sprintrawdefault, &isnulltype);
}

void freemodule_value() {
	freeval(false);
	freeval(true);
	freeval(null);
}
