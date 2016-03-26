#include <string.h>
#include "sys.h"
#include "testing.h"
#include "error.h"
#include "error.mock.h"
#include "number.h"
#include "string.h"
#include "compound-types.h"
#include "sysio.h"
#include "io.h"
#include "value.h"
#include "primitives.h"
#include "environment.h"

#include <stdio.h>

void test_sumnumbers_fail_not_pair() {
	value v = mkstringvalc(1, "s");
	clearerrors();
	sumval(v);
	assert(poperror() == invalidtype, "primitives: sum, not pair");
	freeval(v);
}

void test_sumnumbers_fail_not_number() {
	value v = mkstringvalc(1, "s");
	value p = mkpairval(v, null);
	clearerrors();
	sumval(p);
	assert(poperror() == invalidtype, "primitives: sum, not number");
	freeval(p);
	freeval(v);
}

void test_sumnumbers_null() {
	value s = sumval(null);
	assert(valrawint(s) == 0, "primitives: sum, null");
	freeval(s);
}

void test_sumnumbers() {
	value n1 = mknumvali(1, 1);
	value n2 = mknumvali(2, 1);
	value n3 = mknumvali(3, 1);
	value p1 = mkpairval(n3, null);
	value p2 = mkpairval(n2, p1);
	value p3 = mkpairval(n1, p2);
	value s = sumval(p3);
	assert(valrawint(s) == 6, "primitives: sum numbers");
	freeval(s);
	freeval(p3);
	freeval(p2);
	freeval(p1);
	freeval(n1);
	freeval(n2);
	freeval(n3);
}

void test_diff_error() {
	value s = mkstringvalc(11, "some string");
	value p = mkpairval(s, null);
	clearerrors();
	diffval(p);
	assert(poperror() == invalidtype, "primitives: diff, error");
	freeval(p);
	freeval(s);
}

void test_diff_zero() {
	value d = diffval(null);
	assert(valrawint(d) == 0, "primitives: diff, zero");
	freeval(d);
}

void test_diff_one() {
	value n = mknumvali(42, 1);
	value p = mkpairval(n, null);
	value d = diffval(p);
	assert(valrawint(d) == -42, "primitives: diff, one");
	freeval(d);
	freeval(p);
	freeval(n);
}

void test_diff_all() {
	value n0 = mknumvali(1, 1);
	value n1 = mknumvali(2, 1);
	value n2 = mknumvali(3, 1);
	value p0 = mkpairval(n2, null);
	value p1 = mkpairval(n1, p0);
	value p2 = mkpairval(n0, p1);
	value d = diffval(p2);
	assert(valrawint(d) == -4, "primitives: diff, all");
	freeval(d);
	freeval(p2);
	freeval(p1);
	freeval(p0);
	freeval(n0);
	freeval(n1);
	freeval(n2);
}

void test_bitor_fail_not_number() {
	value s = mkstringvalc(11, "some string");
	clearerrors();
	bitorval(s);
	assert(poperror() == invalidtype, "primitives: bitor fail, not a number");
	freeval(s);
}

void test_bitor_fail_not_integer() {
	value n = mknumvalc(4, "3.14");
	clearerrors();
	bitorval(n);
	assert(poperror() == invalidtype, "primitives: bitor fail, not integer");
	freeval(n);
}

void test_bitor_zero_args() {
	value n = bitorval(null);
	assert(valrawint(n) == 0, "primitives: bitor, zero args");
	freeval(n);
}

void test_bitor_one_arg() {
	value n = mknumvali(42, 1);
	value p = mkpairval(n, null);
	value ior = bitorval(p);
	assert(valrawint(ior) == valrawint(n), "primitives: bitor, one arg");
	freeval(n);
	freeval(p);
	freeval(ior);
}

void test_bitor_multiple_args() {
	value n1 = mknumvali(1, 1);
	value n2 = mknumvali(2, 1);
	value n3 = mknumvali(4, 1);
	value p1 = mkpairval(n3, null);
	value p2 = mkpairval(n2, p1);
	value p3 = mkpairval(n1, p2);
	value ior = bitorval(p3);
	assert(valrawint(ior) == 7, "primitives: bitor, multiple args");
	freeval(ior);
	freeval(p3);
	freeval(p2);
	freeval(p1);
	freeval(n1);
	freeval(n2);
	freeval(n3);
}

void test_openfile_wrong_arg() {
	value v = mknumvali(1, 1);
	value p = mkpairval(v, null);
	clearerrors();
	openfileval(p);
	assert(poperror() == invalidtype, "primitives: open file, wrong arg type");
	freeval(p);
	freeval(v);
}

void test_openfile_wrong_number_of_args() {
	value fn = mkstringvalc(11, "some string");
	value p = mkpairval(fn, null);
	clearerrors();
	openfileval(p);
	assert(poperror() == invalidnumberofargs, "primitives: open file, wrong number of args");
	freeval(p);
	freeval(fn);
}

void test_openfile() {
	value fn = mkstringvalc(11, "some string");
	value mode = mknumvali(ioread, 1);
	value p0 = mkpairval(mode, null);
	value p1 = mkpairval(fn, p0);
	value f = openfileval(p1);
	freeval(p1);
	freeval(fn);
	closefile(rawval(f));
	freeval(f);
}

void test_closefile() {
	value fn = mkstringvalc(11, "some string");
	value mode = mknumvali(ioread, 1);
	value p0 = mkpairval(mode, null);
	value p1 = mkpairval(fn, p0);
	value f = openfileval(p1);
	value p = mkpairval(f, null);
	value ok = closefileval(p);
	assert(!strcoll(valsymbolnameraw(ok), "ok"), "primitives: close file");
	freeval(ok);
	freeval(p1);
	freeval(p0);
	freeval(fn);
	freeval(mode);
	freeval(p);
}

void test_seekfile() {
	value fn = mkstringvalc(11, "some string");
	value mode = mknumvali(ioread, 1);
	value p0 = mkpairval(mode, null);
	value p1 = mkpairval(fn, p0);
	value f = openfileval(p1);

	value seekmode = mknumvali(iostart, 1);
	value pos = mknumvali(42, 1);
	value p2 = mkpairval(seekmode, null);
	value p3 = mkpairval(pos, p2);
	value p4 = mkpairval(f, p3);
	value ok = seekfileval(p4);
	assert(!strcoll(valsymbolnameraw(ok), "ok"), "primitives: close file");

	value p5 = mkpairval(f, null);
	closefileval(p5);

	freeval(p5);
	freeval(ok);
	freeval(p1);
	freeval(p0);
	freeval(fn);
	freeval(mode);
	freeval(p4);
	freeval(p3);
	freeval(p2);
	freeval(pos);
	freeval(seekmode);
}

void test_readfile() {
	value fn = mkstringvalc(11, "some string");
	value mode = mknumvali(ioread, 1);
	value p0 = mkpairval(mode, null);
	value p1 = mkpairval(fn, p0);
	value f = openfileval(p1);

	value len = mknumvali(42, 1);
	value p2 = mkpairval(len, null);
	value p3 = mkpairval(f, p2);
	value s = readfileval(p3);

	value p = mkpairval(f, null);
	value ok = closefileval(p);

	freeval(ok);
	freeval(p1);
	freeval(p0);
	freeval(fn);
	freeval(mode);
	freeval(p3);
	freeval(p2);
	freeval(s);
	freeval(len);
}

void test_writefile() {
	value fn = mkstringvalc(11, "some string");
	value mode = mknumvali(iowrite, 1);
	value p0 = mkpairval(mode, null);
	value p1 = mkpairval(fn, p0);
	value f = openfileval(p1);

	value s = mkstringvalc(17, "some other string");
	value p2 = mkpairval(s, null);
	value p3 = mkpairval(f, p2);
	value okw = writefileval(p3);

	value p = mkpairval(f, null);
	value ok = closefileval(p);

	freeval(ok);
	freeval(okw);
	freeval(p1);
	freeval(p0);
	freeval(fn);
	freeval(mode);
	freeval(p3);
	freeval(p2);
	freeval(s);
	freeval(p);
}

void test_init_free_regex() {
	char *exps = "some expression";
	int len = strlen(exps);
	value exp = mkstringvalc(len, exps);
	value flags = mknumvali(0, 1);
	value p0 = mkpairval(flags, null);
	value p1 = mkpairval(exp, p0);
	value rx = mkregexval(p1);
	freeval(rx);
	freeval(p1);
	freeval(p0);
	freeval(flags);
	freeval(exp);
}

void test_regex_match() {
	char *exps = "\\([^(]*\\)";
	int len = strlen(exps);
	value exp = mkstringvalc(len, exps);
	value flags = mknumvali(0, 1);
	value p0 = mkpairval(flags, null);
	value p1 = mkpairval(exp, p0);
	value rx = mkregexval(p1);

	char *s = "((some list) (of lists))";
	int slen = strlen(s);
	value sval = mkstringvalc(slen, s);
	value p2 = mkpairval(sval, null);
	value p3 = mkpairval(rx, p2);
	value m = regexmatch(p3);
	assert(ispairtype(m), "primitives: regex match");
	assert(ispairtype(carval(m)), "primitives: regex match, first");
	assert(isnumtype(carval(carval(m))), "primitives: regex match, first, index");
	assert(valrawint(carval(carval(m))) == 1, "primitives: regex match, first, index value");
	assert(isnumtype(carval(cdrval(carval(m)))), "primitives: regex match, first, length");
	assert(valrawint(carval(cdrval(carval(m)))) == 11, "primitives: regex match, first, length value");

	freeval(p3);
	freeval(p2);
	freeval(p1);
	freeval(p0);
	freeval(flags);
	freeval(exp);
	freeval(sval);
	freeval(rx);
	freeval(carval(carval(m)));
	freeval(carval(cdrval(carval(m))));
	freeval(carval(m));
	freeval(m);
}

void test_regex_nomatch() {
	value exp = mkstringvalc(1, "a");
	value flags = mknumvali(0, 1);
	value p0 = mkpairval(flags, null);
	value p1 = mkpairval(exp, p0);
	value rx = mkregexval(p1);
	value s = mkstringvalc(1, "b");
	value p2 = mkpairval(s, null);
	value p3 = mkpairval(rx, p2);
	assert(regexmatch(p3) == null, "primitives: regex, nomatch");
	freeval(p3);
	freeval(p2);
	freeval(s);
	freeval(rx);
	freeval(p1);
	freeval(p0);
	freeval(flags);
	freeval(exp);
}

void test_regex_match_utf8() {
	value exp = mkstringvalc(2, "nf");
	value f = mknumvali(0, 1);
	value p0 = mkpairval(f, null);
	value p1 = mkpairval(exp, p0);
	value rx = mkregexval(p1);

	char *s = "fűzfánfütyülő";
	long len = strlen(s);
	value sv = mkstringvalc(len, s);
	value p2 = mkpairval(sv, null);
	value p3 = mkpairval(rx, p2);

	value m = regexmatch(p3);
	assert(valrawint(carval(carval(m))) == 5, "primitives: regex, match, utf8");

	freeval(m);
	freeval(p3);
	freeval(p2);
	freeval(sv);
	freeval(rx);
	freeval(p1);
	freeval(p0);
	freeval(f);
	freeval(exp);
}

void test_regex_capture_groups() {
	value exp = mkstringvalc(9, "a(b(c)d)e");
	value f = mknumvali(0, 1);
	value p0 = mkpairval(f, null);
	value p1 = mkpairval(exp, p0);
	value rx = mkregexval(p1);

	value s = mkstringvalc(5, "abcde");
	value p2 = mkpairval(s, null);
	value p3 = mkpairval(rx, p2);

	value m = regexmatch(p3);
	assert(!isnulltype(m), "primitives: regex, capture groups, not null 0");
	assert(!isnulltype(cdrval(m)), "primitives: regex, capture groups, not null 1");
	assert(!isnulltype(cdrval(cdrval(m))), "primitives: regex, capture groups, not null 2");
	assert(isnulltype(cdrval(cdrval(cdrval(m)))), "primitives: regex, capture groups, length");
	assert(valrawint(carval(carval(m))) == 0, "primitives: regex, capture groups, first index");
	assert(valrawint(carval(cdrval(carval(m)))) == 5, "primitives: regex, capture groups, first len");
	assert(valrawint(carval(carval(cdrval(m)))) == 1, "primitives: regex, capture groups, second index");
	assert(valrawint(carval(cdrval(carval(cdrval(m))))) == 3, "primitives: regex, capture groups, second len");
	assert(valrawint(carval(carval(cdrval(cdrval(m))))) == 2, "primitives: regex, capture groups, third index");
	assert(valrawint(carval(cdrval(carval(cdrval(cdrval(m)))))) == 1, "primitives: regex, capture groups, third len");

	freeval(m);
	freeval(p3);
	freeval(p2);
	freeval(s);
	freeval(rx);
	freeval(p1);
	freeval(p0);
	freeval(f);
	freeval(exp);
}

void test_error() {
	char *msg = "some message";
	int len = strlen(msg);
	value msgval = mkstringvalc(len, msg);
	value args = mkpairval(msgval, null);
	clearerrors();
	value r = errorval(args);
	assert(poperror() == usererror, "primitives: error");
	freeval(r);
	freeval(args);
	freeval(msgval);
}

void test_eq_symbol_false() {
	value s1 = mksymval(3, "abc");
	value s2 = mksymval(3, "Abc");
	value args = mkpairval(s2, null);
	args = mkpairval(s1, args);
	assert(iseqval(args) == false, "primitives: eq, symbol, false");
	freeval(cdrval(args));
	freeval(args);
	freeval(s1);
	freeval(s2);
}

void test_eq_symbol_true() {
	value s1 = mksymval(3, "abc");
	value s2 = mksymval(3, "abc");
	value args = mkpairval(s2, null);
	args = mkpairval(s1, args);
	assert(iseqval(args) != false, "primitives: eq, symbol, true");
	freeval(cdrval(args));
	freeval(args);
	freeval(s1);
	freeval(s2);
}

void test_eq_number_false() {
	value n1 = mknumvali(42, 1);
	value n2 = mknumvali(3, 1);
	value args = mkpairval(n2, null);
	args = mkpairval(n1, args);
	assert(iseqval(args) == false, "primitives: eq, number, false");
	freeval(cdrval(args));
	freeval(args);
	freeval(n1);
	freeval(n2);
}

void test_eq_number_true() {
	value n1 = mknumvali(3, 1);
	value n2 = mknumvali(3, 1);
	value args = mkpairval(n2, null);
	args = mkpairval(n1, args);
	assert(iseqval(args) != false, "primitives: eq, number, true");
	freeval(cdrval(args));
	freeval(args);
	freeval(n1);
	freeval(n2);
}

void test_eq_string_false() {
	value s1 = mkstringvalc(3, "abc");
	value s2 = mkstringvalc(3, "Abc");
	value args = mkpairval(s2, null);
	args = mkpairval(s1, args);
	assert(iseqval(args) == false, "primitives: eq, string, false");
	freeval(cdrval(args));
	freeval(args);
	freeval(s1);
	freeval(s2);
}

void test_eq_string_true() {
	value s1 = mkstringvalc(3, "abc");
	value s2 = mkstringvalc(3, "abc");
	value args = mkpairval(s2, null);
	args = mkpairval(s1, args);
	assert(iseqval(args) != false, "primitives: eq, string, true");
	freeval(cdrval(args));
	freeval(args);
	freeval(s1);
	freeval(s2);
}

void test_eq_refs_false() {
	value entry = mknumvali(1, 1);
	environment env = mkenvironment(0);
	value v1 = mkcompiledprocval(entry, env);
	value v2 = mkcompiledprocval(entry, env);
	value args = mkpairval(v2, null);
	args = mkpairval(v1, args);
	args = mkpairval(v1, args);
	assert(iseqval(args) == false, "primitives: eq, refs, false");
	freeval(cdrval(cdrval(args)));
	freeval(cdrval(args));
	freeval(args);
	freeval(entry);
	freenv(env);
	freeval(v1);
	freeval(v2);
}

void test_eq_refs_true() {
	value entry = mknumvali(1, 1);
	environment env = mkenvironment(0);
	value v1 = mkcompiledprocval(entry, env);
	value args = mkpairval(v1, null);
	args = mkpairval(v1, args);
	args = mkpairval(v1, args);
	assert(iseqval(args) != false, "primitives: eq, refs, true");
	freeval(cdrval(cdrval(args)));
	freeval(cdrval(args));
	freeval(args);
	freeval(entry);
	freenv(env);
	freeval(v1);
}

void test_eq_refs_one() {
	value v = mknumvali(1, 1);
	value args = mkpairval(v, null);
	assert(iseqval(args) != false, "primitives: eq, refs, true");
	freeval(args);
	freeval(v);
}

void test_eq_refs_zero() {
	assert(iseqval(null) != false, "primitives: eq, refs, true");
}

void test_isutf8_false() {
	char *raw = (char []){255, 1, 0};
	value s = mkstringvalc(2, raw);
	value p = mkpairval(s, null);
	assert(isutf8val(p) == false, "primitives: isutf8, false");
	freeval(p);
	freeval(s);
}

void test_isutf8_true() {
	value s = mkstringvalc(11, "some string");
	value p = mkpairval(s, null);
	assert(isutf8val(p) != false, "primitives: isutf8, true");
	freeval(p);
	freeval(s);
}

void test_copystrval_error() {
	value s = mknumvali(42, 1);
	value from = mknumvali(3, 1);
	value len = mknumvali(3, 1);
	value p0 = mkpairval(len, null);
	value p1 = mkpairval(from, p0);
	value p2 = mkpairval(s, p1);
	clearerrors();
	copystrval(p2);
	assert(poperror() == invalidtype, "primitives: substring, error");
	freeval(p0);
	freeval(p1);
	freeval(p2);
	freeval(len);
	freeval(from);
	freeval(s);
}

void test_copystrval() {
	value s = mkstringvalc(11, "some string");
	value from = mknumvali(3, 1);
	value len = mknumvali(3, 1);
	value p0 = mkpairval(len, null);
	value p1 = mkpairval(from, p0);
	value p2 = mkpairval(s, p1);
	value ss = copystrval(p2);
	assert(!strcoll(valrawstring(ss), "e s"), "primitives: substring");
	freeval(p0);
	freeval(p1);
	freeval(p2);
	freeval(len);
	freeval(from);
	freeval(s);
	freeval(ss);
}

void test_byteslenval_error() {
	value v = mknumvali(42, 1);
	value p = mkpairval(v, null);
	clearerrors();
	byteslenval(p);
	assert(poperror() == invalidtype, "primitives: byteslenval, error");
	freeval(p);
	freeval(v);
}

void test_byteslenval() {
	value s = mkstringvalc(11, "some string");
	value p = mkpairval(s, null);
	value l = byteslenval(p);
	assert(valrawint(l) == 11, "primitives: byteslenval");
	freeval(p);
	freeval(s);
	freeval(l);
}

void test_stringappend_error() {
	value v = mknumvali(42, 1);
	value p = mkpairval(v, null);
	clearerrors();
	stringappendval(p);
	assert(poperror() == invalidtype, "primitives: stringappend, error");
	freeval(p);
	freeval(v);
}

void test_stringappend() {
	value s1 = mkstringvalc(5, "some ");
	value s2 = mkstringvalc(6, "string");
	value p0 = mkpairval(s2, null);
	value p1 = mkpairval(s1, p0);
	value sa = stringappendval(p1);
	assert(!strcoll(valrawstring(sa), "some string"), "primitives: stringappend");
	freeval(sa);
	freeval(p1);
	freeval(p0);
	freeval(s1);
	freeval(s2);
}

void test_stringlenval_error() {
	value v = mknumvali(42, 1);
	value p = mkpairval(v, null);
	clearerrors();
	stringlenval(p);
	assert(poperror() == invalidtype, "primitives: stringlenval, error");
	freeval(p);
	freeval(v);
}

void test_stringlenval() {
	value s = mkstringvalc(11, "some string");
	value p = mkpairval(s, null);
	value l = stringlenval(p);
	assert(valrawint(l) == 11, "primitives: stringlenval");
	freeval(l);
	freeval(p);
	freeval(s);
}

void test_iseofval_false() {
	value v = mknumvali(42, 1);
	value p = mkpairval(v, null);
	assert(iseofval(p) == false, "primitives: iseof, false");
	freeval(v);
	freeval(p);
}

void test_iseofval_true() {
	value p = mkpairval(eofval, null);
	assert(iseofval(p) == true, "primitives: iseof, true");
	freeval(p);
}

void test_lessval_error() {
	value p = mkpairval(null, null);
	clearerrors();
	islessval(p);
	assert(poperror() == invalidtype, "primitives: islessval, error");
	freeval(p);
}

void test_lessval_false() {
	value n1 = mknumvali(1, 1);
	value n2 = mknumvali(2, 1);
	value n3 = mknumvali(0, 1);
	value p0 = mkpairval(n3, null);
	value p1 = mkpairval(n2, p0);
	value p2 = mkpairval(n1, p1);
	assert(islessval(p2) == false, "primitives: islessval, false");
	freeval(p2);
	freeval(p1);
	freeval(p0);
	freeval(n1);
	freeval(n2);
	freeval(n3);
}

void test_lessval_true() {
	value n1 = mknumvali(1, 1);
	value n2 = mknumvali(2, 1);
	value n3 = mknumvali(3, 1);
	value p0 = mkpairval(n3, null);
	value p1 = mkpairval(n2, p0);
	value p2 = mkpairval(n1, p1);
	assert(islessval(p2) != false, "primitives: islessval, false");
	freeval(p2);
	freeval(p1);
	freeval(p0);
	freeval(n1);
	freeval(n2);
	freeval(n3);
}

void test_lessval_string() {
	value s1 = mkstringvalc(1, "a");
	value s2 = mkstringvalc(1, "b");
	value p0 = mkpairval(s2, null);
	value p1 = mkpairval(s1, p0);
	assert(islessval(p1) != false, "primitives: islessval, string");
	freeval(p1);
	freeval(p0);
	freeval(s1);
	freeval(s2);
}

void test_greaterval_false() {
	value n1 = mknumvali(3, 1);
	value n2 = mknumvali(1, 1);
	value n3 = mknumvali(2, 1);
	value p0 = mkpairval(n3, null);
	value p1 = mkpairval(n2, p0);
	value p2 = mkpairval(n1, p1);
	assert(isgreaterval(p2) == false, "primitives: isgreaterval, false");
	freeval(p2);
	freeval(p1);
	freeval(p0);
	freeval(n1);
	freeval(n2);
	freeval(n3);
}

void test_greaterval_true() {
	value n1 = mknumvali(3, 1);
	value n2 = mknumvali(2, 1);
	value n3 = mknumvali(1, 1);
	value p0 = mkpairval(n3, null);
	value p1 = mkpairval(n2, p0);
	value p2 = mkpairval(n1, p1);
	assert(isgreaterval(p2) != false, "primitives: isgreaterval, false");
	freeval(p2);
	freeval(p1);
	freeval(p0);
	freeval(n1);
	freeval(n2);
	freeval(n3);
}

void test_lessoreqval_false() {
	value n1 = mknumvali(1, 1);
	value n2 = mknumvali(1, 1);
	value n3 = mknumvali(0, 1);
	value p0 = mkpairval(n3, null);
	value p1 = mkpairval(n2, p0);
	value p2 = mkpairval(n1, p1);
	assert(islessoreqval(p2) == false, "primitives: islessoreqval, false");
	freeval(p2);
	freeval(p1);
	freeval(p0);
	freeval(n1);
	freeval(n2);
	freeval(n3);
}

void test_lessoreqval_true() {
	value n1 = mknumvali(1, 1);
	value n2 = mknumvali(1, 1);
	value n3 = mknumvali(2, 1);
	value p0 = mkpairval(n3, null);
	value p1 = mkpairval(n2, p0);
	value p2 = mkpairval(n1, p1);
	assert(islessoreqval(p2) != false, "primitives: islessoreqval, false");
	freeval(p2);
	freeval(p1);
	freeval(p0);
	freeval(n1);
	freeval(n2);
	freeval(n3);
}

void test_greateroreqval_false() {
	value n1 = mknumvali(1, 1);
	value n2 = mknumvali(1, 1);
	value n3 = mknumvali(2, 1);
	value p0 = mkpairval(n3, null);
	value p1 = mkpairval(n2, p0);
	value p2 = mkpairval(n1, p1);
	assert(isgreateroreqval(p2) == false, "primitives: isgreateroreqval, false");
	freeval(p2);
	freeval(p1);
	freeval(p0);
	freeval(n1);
	freeval(n2);
	freeval(n3);
}

void test_greateroreqval_true() {
	value n1 = mknumvali(1, 1);
	value n2 = mknumvali(1, 1);
	value n3 = mknumvali(0, 1);
	value p0 = mkpairval(n3, null);
	value p1 = mkpairval(n2, p0);
	value p2 = mkpairval(n1, p1);
	assert(isgreateroreqval(p2) != false, "primitives: isgreateroreqval, false");
	freeval(p2);
	freeval(p1);
	freeval(p0);
	freeval(n1);
	freeval(n2);
	freeval(n3);
}

void test_notval_false() {
	value v = mknumvali(42, 1);
	value p = mkpairval(v, null);
	assert(notval(p) == false, "primitives: not, false");
	freeval(v);
	freeval(p);
}

void test_notval_true() {
	value p = mkpairval(false, null);
	assert(notval(p) != false, "primitives: not, true");
	freeval(p);
}

void test_isnullvalp_false() {
	value v = mknumvali(42, 1);
	value p = mkpairval(v, null);
	assert(isnullvalp(p) == false, "primitives: isnullvalp, false");
	freeval(p);
	freeval(v);
}

void test_isnullvalp_true() {
	value p = mkpairval(null, null);
	assert(isnullvalp(p) != false, "primitives: isnullvalp, true");
	freeval(p);
}

void test_consval_error() {
	value a = mknumvali(42, 1);
	value p = mkpairval(a, null);
	clearerrors();
	consval(p);
	assert(poperror() == invalidnumberofargs, "primitives: consval, error");
	freeval(p);
	freeval(a);
}

void test_consval() {
	value a = mknumvali(42, 1);
	value b = mknumvali(36, 1);
	value p0 = mkpairval(b, null);
	value p1 = mkpairval(a, p0);
	value p = consval(p1);
	assert(carval(p) == a, "primitives: consval, car");
	assert(cdrval(p) == b, "primitives: consval, cdr");
	freeval(p);
	freeval(p1);
	freeval(p0);
	freeval(b);
	freeval(a);
}

void test_carvalp_error() {
	value v = mknumvali(42, 1);
	value args = mkpairval(v, null);
	clearerrors();
	carvalp(args);
	assert(poperror() == invalidtype, "primitives: carvalp, error");
	freeval(args);
	freeval(v);
}

void test_carvalp() {
	value v1 = mknumvali(42, 1);
	value v2 = mknumvali(84, 1);
	value p = mkpairval(v1, v2);
	value args = mkpairval(p, null);
	assert(carvalp(args) == v1, "primitives: carvalp");
	freeval(args);
	freeval(p);
	freeval(v1);
	freeval(v2);
}

void test_cdrvalp_error() {
	value v = mknumvali(42, 1);
	value args = mkpairval(v, null);
	clearerrors();
	cdrvalp(args);
	assert(poperror() == invalidtype, "primitives: carvalp, error");
	freeval(args);
	freeval(v);
}

void test_cdrvalp() {
	value v1 = mknumvali(42, 1);
	value v2 = mknumvali(84, 1);
	value p = mkpairval(v1, v2);
	value args = mkpairval(p, null);
	assert(cdrvalp(args) == v2, "primitives: carvalp");
	freeval(args);
	freeval(p);
	freeval(v1);
	freeval(v2);
}

void test_stringtonumsafe_fail() {
	value s = mkstringvalc(1, "a");
	value p = mkpairval(s, null);
	value n = stringtonumsafe(p);
	assert(n == false, "primitives: stringtonumsafe, fail");
	freeval(p);
	freeval(s);
}

void test_stringtonumsafe() {
	value s = mkstringvalc(1, "1");
	value p = mkpairval(s, null);
	value n = stringtonumsafe(p);
	assert(valrawint(n) == 1, "primitives: stringtonumsafe");
	freeval(n);
	freeval(p);
	freeval(s);
}

int main() {
	initsys();
	initmodule_errormock();
	initmodule_number();
	initmodule_value();
	initmodule_sysio();
	initmodule_primitives();

	test_sumnumbers_fail_not_pair();
	test_sumnumbers_fail_not_number();
	test_sumnumbers_null();
	test_sumnumbers();
	test_diff_error();
	test_diff_zero();
	test_diff_one();
	test_diff_all();
	test_bitor_fail_not_number();
	test_bitor_fail_not_integer();
	test_bitor_zero_args();
	test_bitor_one_arg();
	test_bitor_multiple_args();
	test_openfile_wrong_arg();
	test_openfile_wrong_number_of_args();
	test_openfile();
	test_closefile();
	test_seekfile();
	test_readfile();
	test_writefile();
	test_init_free_regex();
	test_regex_match();
	test_regex_nomatch();
	test_regex_match_utf8();
	test_regex_capture_groups();
	test_error();
	test_eq_symbol_false();
	test_eq_symbol_true();
	test_eq_number_false();
	test_eq_number_true();
	test_eq_string_false();
	test_eq_string_true();
	test_eq_refs_false();
	test_eq_refs_true();
	test_eq_refs_one();
	test_eq_refs_zero();
	test_isutf8_false();
	test_isutf8_true();
	test_copystrval_error();
	test_copystrval();
	test_byteslenval_error();
	test_byteslenval();
	test_stringappend_error();
	test_stringappend();
	test_stringlenval_error();
	test_stringlenval();
	test_iseofval_false();
	test_iseofval_true();
	test_lessval_error();
	test_lessval_false();
	test_lessval_true();
	test_lessval_string();
	test_greaterval_false();
	test_greaterval_true();
	test_lessoreqval_false();
	test_lessoreqval_true();
	test_greateroreqval_false();
	test_greateroreqval_true();
	test_notval_false();
	test_notval_true();
	test_isnullvalp_false();
	test_isnullvalp_true();
	test_consval_error();
	test_consval();
	test_carvalp_error();
	test_carvalp();
	test_cdrvalp_error();
	test_cdrvalp();
	test_stringtonumsafe_fail();
	test_stringtonumsafe();

	freemodule_primitives();
	freemodule_sysio();
	freemodule_value();
	freemodule_number();
	freemodule_errormock();
	return 0;
}
