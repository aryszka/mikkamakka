#include <stdlib.h>
#include <string.h>
#include "sys.h"
#include "sysio.h"
#include "testing.h"
#include "number.h"
#include "string.h"
#include "error.h"
#include "error.mock.h"
#include "compound-types.h"
#include "environment.h"
#include "io.h"
#include "value.h"

value testprimitive(value args) {
	return args;
}

void test_init_free_symbol() {
	char *name = "symbol";
	int len = strlen(name);
	value s = mksymval(len, name);
	freeval(s);
}

void test_init_free_number() {
	number n = mknumi(1, 1);
	value v = mknumval(n);
	freeval(v);
	freenum(n);
}

void test_init_free_number_char() {
	char *c = "12.3";
	int len = strlen(c);
	value n = mknumvalc(len, c);
	freeval(n);
}

void test_init_free_number_int() {
	value n = mknumvali(3, 4);
	freeval(n);
}

void test_init_free_string() {
	char *raw = "some string";
	int len = strlen(raw);
	string s = mkstring(len, raw);
	value sv = mkstringval(s);
	freeval(sv);
	freestring(s);
}

void test_init_free_string_char() {
	char *raw = "some string";
	int len = strlen(raw);
	value s = mkstringvalc(len, raw);
	freeval(s);
}

void test_issymtype_false() {
	value v = mknumvali(1, 1);
	assert(!issymtype(v), "value: type not symbol");
	freeval(v);
}

void test_issymtype() {
	char *name = "some-symbol";
	int nlen = strlen(name);
	value v = mksymval(nlen, name);
	assert(issymtype(v), "value: type symbol");
	freeval(v);
}

void test_issymval_false() {
	value v = mknumvali(1, 1);
	assert(issymval(v) == false, "value: not symbol");
	freeval(v);
}

void test_issymval() {
	char *name = "some-symbol";
	int nlen = strlen(name);
	value v = mksymval(nlen, name);
	assert(issymval(v) != false, "value: type symbol");
	freeval(v);
}

void test_isnumtype_false() {
	value v = mkstringvalc(1, "s");
	assert(!isnumtype(v), "value: type not number");
	freeval(v);
}

void test_isnumtype() {
	value v = mknumvali(1, 1);
	assert(isnumtype(v), "value: type number");
	freeval(v);
}

void test_isnumval_false() {
	value v = mkstringvalc(1, "s");
	assert(isnumval(v) == false, "value: not number");
	freeval(v);
}

void test_isnumval() {
	value v = mknumvali(1, 1);
	assert(isnumval(v) != false, "value: type number");
	freeval(v);
}

void test_isfalse_type_false() {
	value n = mknumvali(1, 1);
	assert(!isfalsetype(n), "value: type not false");
	freeval(n);
}

void test_isfalse_type() {
	assert(isfalsetype(false), "value: type false");
}

void test_isfalse_false() {
	value n = mknumvali(1, 1);
	assert(isfalseval(n) == false, "value: not false");
	freeval(n);
}

void test_isfalse() {
	assert(isfalseval(false) != false, "value: not false");
}

void test_pair() {
	value n1 = mknumvali(1, 1);
	value n2 = mknumvali(2, 1);
	value p = mkpairval(n1, n2);
	assert(carval(p) == n1, "value: car pair");
	assert(cdrval(p) == n2, "value: cdr pair");
	freeval(p);
	freeval(n1);
	freeval(n2);
}

void test_isnulltype_false() {
	value n = mknumvali(1, 1);
	assert(!isnulltype(n), "value: type not null");
	freeval(n);
}

void test_isnulltype() {
	assert(isnulltype(null), "value: type null");
}

void test_isnull_false() {
	value n = mknumvali(1, 1);
	assert(isnullval(n) == false, "value: is not null");
	freeval(n);
}

void test_isnull() {
	assert(isnullval(null) != false, "value: is null");
}

void test_ispairtype_false() {
	value n = mknumvali(1, 1);
	assert(!ispairtype(n), "value: type not pair");
	freeval(n);
}

void test_ispairtype() {
	value n0 = mknumvali(1, 1);
	value n1 = mknumvali(2, 1);
	value p = mkpairval(n0, n1);
	assert(ispairtype(p), "value: type pair");
	freeval(p);
	freeval(n0);
	freeval(n1);
}

void test_ispair_false() {
	value n = mknumvali(1, 1);
	assert(ispairval(n) == false, "value: is not pair");
	freeval(n);
}

void test_ispair() {
	value n0 = mknumvali(1, 1);
	value n1 = mknumvali(2, 1);
	value p = mkpairval(n0, n1);
	assert(ispairval(p) != false, "value: is pair");
	freeval(p);
	freeval(n0);
	freeval(n1);
}

void test_isproctype_false() {
	value n = mknumvali(1, 1);
	assert(!isproctype(n), "value: type not proc");
	freeval(n);
}

void test_isproctype() {
	value p = mkprimitiveprocval(testprimitive);
	assert(isproctype(p), "value: type proc");
	freeval(p);
}

void test_isprocval_false() {
	value n = mknumvali(1, 1);
	assert(isprocval(n) == false, "value: is not proc");
	freeval(n);
}

void test_isprocval() {
	value p = mkprimitiveprocval(testprimitive);
	assert(isprocval(p) != false, "value: is pair");
	freeval(p);
}

void test_car_not_pair() {
	value n = mknumvali(1, 1);

	clearerrors();
	carval(n);
	assert(poperror() == invalidtype, "value: car not pair");

	clearerrors();
	freeval(n);
}

void test_cdr_not_pair() {
	value n = mknumvali(1, 1);

	clearerrors();
	cdrval(n);
	assert(poperror() == invalidtype, "value: cdr not pair");

	clearerrors();
	freeval(n);
}

void test_init_free_primitive_proc() {
	value p = mkprimitiveprocval(&testprimitive);
	freeval(p);
}

void test_init_free_compiled_proc() {
	value label = mknumvali(1, 1);
	environment env = mkenvironment(0);
	value p = mkcompiledprocval(label, env);
	freeval(p);
	freeval(label);
	freenv(env);
}

void test_numval_fail() {
	value v = mkstringvalc(11, "some string");
	clearerrors();
	numval(v);
	assert(poperror() == invalidtype, "value: numval, fail");
	freeval(v);
}

void test_numval() {
	value n = mknumvali(1, 1);
	number nn = numval(n);
	assert(rawint(nn) == valrawint(n), "value: numval");
}

void test_rawint_fail_not_number() {
	value v = mkstringvalc(1, "s");
	clearerrors();
	valrawint(v);
	assert(poperror() == invalidtype, "value: rawint, not number");
	freeval(v);
}

void test_rawint_fail_not_integer() {
	value n = mknumvali(1, 2);
	clearerrors();
	valrawint(n);
	assert(poperror() == numbernotint, "value: rawint, not integer");
	freeval(n);
}

void test_rawint() {
	value n = mknumvali(42, 1);
	assert(valrawint(n) == 42, "value: rawint");
	freeval(n);
}

void test_applyprimitive_not_proc() {
	value n = mknumvali(1, 1);
	clearerrors();
	valapplyprimitive(n, null);
	assert(poperror() == invalidtype, "value: apply primitive, not proc");
	freeval(n);
}

void test_applyprimitive_not_primitive() {
	value label = mknumvali(1, 1);
	environment env = mkenvironment(0);
	value p = mkcompiledprocval(label, env);
	clearerrors();
	valapplyprimitive(p, null);
	assert(poperror() == invalidtype, "value: apply primitive, not primitive");
	freeval(p);
	freeval(label);
	freenv(env);
}

void test_applyprimitive() {
	value p = mkprimitiveprocval(testprimitive);
	assert(valapplyprimitive(p, null) == null, "value: apply primitive proc");
	freeval(p);
}

void test_proclabel_not_proc() {
	value n = mknumvali(1, 1);
	clearerrors();
	valproclabel(n);
	assert(poperror() == invalidtype, "value: proc label, not proc");
	freeval(n);
}

void test_proclabel_not_compiled() {
	value p = mkprimitiveprocval(testprimitive);
	clearerrors();
	valproclabel(p);
	assert(poperror() == invalidtype, "value: proc label, not compiled");
	freeval(p);
}

void test_proclabel() {
	value label = mknumvali(1, 1);
	environment env = mkenvironment(0);
	value p = mkcompiledprocval(label, env);
	assert(valproclabel(p) == label, "value: proc label");
	freeval(p);
	freeval(label);
	freenv(env);
}

void test_procenv_not_proc() {
	value n = mknumvali(1, 1);
	clearerrors();
	valprocenv(n);
	assert(poperror() == invalidtype, "value: proc env, not proc");
	freeval(n);
}

void test_procenv_not_compiled() {
	value p = mkprimitiveprocval(testprimitive);
	clearerrors();
	valprocenv(p);
	assert(poperror() == invalidtype, "value: proc env, not compiled");
	freeval(p);
}

void test_procenv() {
	value label = mknumvali(1, 1);
	environment env = mkenvironment(0);
	value p = mkcompiledprocval(label, env);
	assert(valprocenv(p) == env, "value: proc env");
	freeval(p);
	freeval(label);
	freenv(env);
}

void test_isprimitiveproctype_false() {
	value v = mknumvali(1, 1);
	assert(!isprimitiveproctype(v), "value: type not primitive proc");
	freeval(v);
}

void test_isprimitiveproctype() {
	value p = mkprimitiveprocval(testprimitive);
	assert(isprimitiveproctype(p), "value: type primitive proc");
	freeval(p);
}

void test_isprimitiveproc_false() {
	value v = mknumvali(1, 1);
	assert(isprimitiveprocval(v) == false, "value: not primitive proc");
	freeval(v);
}

void test_isprimitiveproc() {
	value p = mkprimitiveprocval(testprimitive);
	assert(isprimitiveprocval(p) != false, "value: primitive proc");
	freeval(p);
}

void test_symbolnameraw_not_symbol() {
	value v = mknumvali(1, 1);
	clearerrors();
	valsymbolnameraw(v);
	assert(poperror() == invalidtype, "value: symbolname, not symbol");
	freeval(v);
}

void test_symbolnameraw() {
	char *sname = "some-symbol";
	int slen = strlen(sname);
	value s = mksymval(slen, sname);
	assert(!strcoll(valsymbolnameraw(s), sname), "value: symbolname");
	freeval(s);
}

void test_sprintraw_number_small() {
	value n = mknumvali(9, 4);
	char *s = sprintraw(n);
	assert(!strcoll(s, "2.25"), "value: sprintraw, number, small");
	free(s);
	freeval(n);
}

void test_sprintraw_number_big() {
	char *c = "4798476948769769857028703847560387562984769.4897694769382765794876385760398476039865";
	int len = strlen(c);
	value n = mknumvalc(len, c);
	char *s = sprintraw(n);
	assert(!strcoll(s, c), "value: sprintraw, number, big");
	free(s);
	freeval(n);
}

// void test_init_free_file() {
// 	file f;
// 	ioerror err;
// 	openfile("test-file", ioread, &f, &err);
// 	value fv = mkfileval(f);
// 	freeval(fv);
// 	closefile(f);
// }

// void test_fileval_fail() {
// 	value v = mknumvali(1, 1);
// 	clearerrors();
// 	fileval(v);
// 	assert(poperror() == invalidtype, "value: fileval, wrong type");
// 	freeval(v);
// }
// 
// void test_fileval() {
// 	file f;
// 	ioerror err;
// 	openfile("test-file", ioread, &f, &err);
// 	value v = mkfileval(f);
// 	assert(fileval(v) == f, "value: fileval");
// 	freeval(v);
// 	closefile(f);
// }
// 
// void test_isfiletype_false() {
// 	value v = mknumvali(1, 1);
// 	assert(!isfiletype(v), "value: isfiletype, false");
// 	freeval(v);
// }
// 
// void test_isfiletype_true() {
// 	file f;
// 	ioerror err;
// 	openfile("test-file", ioread, &f, &err);
// 	value v = mkfileval(f);
// 	assert(isfiletype(v), "value: isfiletype, false");
// 	freeval(v);
// 	closefile(f);
// }
// 
// void test_isfileval_false() {
// 	value v = mknumvali(1, 1);
// 	assert(isfileval(v) == false, "value: isfiletype, false");
// 	freeval(v);
// }
// 
// void test_isfileval_true() {
// 	file f;
// 	ioerror err;
// 	openfile("test-file", ioread, &f, &err);
// 	value v = mkfileval(f);
// 	assert(isfileval(v) != false, "value: isfiletype, false");
// 	freeval(v);
// 	closefile(f);
// }

void test_valrawstring_fail() {
	value v = mknumvali(1, 1);
	clearerrors();
	valrawstring(v);
	assert(poperror() == invalidtype, "value: valrawstring, fail");
	freeval(v);
}

void test_valrawstring() {
	value v = mkstringvalc(11, "some string");
	char *s = valrawstring(v);
	assert(!strcoll(s, "some string"), "value: valrawstring");
	freeval(v);
}

void test_isstringtype_false() {
	value v = mknumvali(1, 1);
	assert(!isstringtype(v), "value: isstringtype, false");
	freeval(v);
}

void test_isstringtype_true() {
	value s = mkstringvalc(11, "some string");
	assert(isstringtype(s), "value: isstringtype, true");
	freeval(s);
}

void test_isstringval_false() {
	value v = mknumvali(1, 1);
	assert(isstringval(v) == false, "value: isstringval, false");
	freeval(v);
}

void test_isstringval_true() {
	value s = mkstringvalc(11, "some string");
	assert(isstringval(s) != false, "value: isstringval, true");
	freeval(s);
}

void test_isinttype_false() {
	value v = mkstringvalc(11, "some string");
	assert(!isinttype(v), "value: isinttype, false");
	freeval(v);
}

void test_isinttype_true() {
	value v = mknumvali(42, 1);
	assert(isinttype(v), "value: isinttype, true");
	freeval(v);
}

void test_isintval_false() {
	value v = mkstringvalc(11, "some string");
	assert(isintval(v) == false, "value: isintval, false");
	freeval(v);
}

void test_isintval_true() {
	value v = mknumvali(42, 1);
	assert(isintval(v) != false, "value: isintval, true");
	freeval(v);
}

void test_issmallinttype_false() {
	value v = mkstringvalc(11, "some string");
	assert(!issmallinttype(v), "value: isinttype, false");
	freeval(v);
}

void test_issmallinttype_true() {
	value v = mknumvali(42, 1);
	assert(issmallinttype(v), "value: isinttype, true");
	freeval(v);
}

void test_issmallintval_false() {
	value v = mkstringvalc(11, "some string");
	assert(issmallintval(v) == false, "value: isintval, false");
	freeval(v);
}

void test_issmallintval_true() {
	value v = mknumvali(42, 1);
	assert(issmallintval(v) != false, "value: isintval, true");
	freeval(v);
}

void test_valstring_error() {
	value v = mknumvali(42, 1);
	clearerrors();
	valstring(v);
	assert(poperror() == invalidtype, "value: valstring, error");
	freeval(v);
}

void test_valstring() {
	value v = mkstringvalc(11, "some string");
	string s = valstring(v);
	assert(!strcoll(rawstring(s), valrawstring(v)), "value: valstring");
	freeval(v);
}

int main() {
	initsys();
	initmodule_errormock();
	initmodule_number();
	initmodule_sysio();
	initmodule_value();

	test_init_free_symbol();
	test_init_free_number();
	test_init_free_number_char();
	test_init_free_number_int();
	test_init_free_string();
	test_issymtype_false();
	test_issymtype();
	test_issymval_false();
	test_issymval();
	test_isnumtype_false();
	test_isnumtype();
	test_isnumval_false();
	test_isnumval();
	test_isfalse_type_false();
	test_isfalse_type();
	test_isfalse_false();
	test_isfalse();
	test_pair();
	test_isnulltype_false();
	test_isnulltype();
	test_isnull_false();
	test_isnull();
	test_ispairtype_false();
	test_ispairtype();
	test_ispair_false();
	test_ispair();
	test_isproctype_false();
	test_isproctype();
	test_isprocval_false();
	test_isprocval();
	test_car_not_pair();
	test_cdr_not_pair();
	test_init_free_primitive_proc();
	test_init_free_compiled_proc();
	test_numval_fail();
	test_numval();
	test_rawint_fail_not_number();
	test_rawint_fail_not_integer();
	test_rawint();
	test_applyprimitive_not_proc();
	test_applyprimitive_not_primitive();
	test_applyprimitive();
	test_proclabel_not_proc();
	test_proclabel_not_compiled();
	test_proclabel();
	test_procenv_not_proc();
	test_procenv_not_compiled();
	test_procenv();
	test_isprimitiveproctype_false();
	test_isprimitiveproctype();
	test_isprimitiveproc_false();
	test_isprimitiveproc();
	test_symbolnameraw_not_symbol();
	test_symbolnameraw();
	test_sprintraw_number_small();
	test_valrawstring_fail();
	test_valrawstring();
	test_isstringtype_false();
	test_isstringtype_true();
	test_isstringval_false();
	test_isstringval_true();
	test_isinttype_false();
	test_isinttype_true();
	test_isintval_false();
	test_isintval_true();
	test_issmallinttype_false();
	test_issmallinttype_true();
	test_issmallintval_false();
	test_issmallintval_true();
	test_valstring_error();
	test_valstring();

	freemodule_value();
	freemodule_sysio();
	freemodule_number();
	freemodule_errormock();
}
