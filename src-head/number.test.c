#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <locale.h>
#include "sys.h"
#include "testing.h"
#include "error.h"
#include "error.mock.h"
#include "number.h"

void test_invalid_length() {
	clearerrors();
	mknumc(0, "");
	assert(poperror() == invalidnumber, "error on 0 length");
}

void test_invalid_multidots() {
	char *sin = locnumberstr("1.2.3");
	size_t len = strlen(sin);
	clearerrors();
	number n = mknumc(len, sin);
	assert(n == 0, "number null on error");
	assert(poperror() == invalidnumber, "error on 0 length");
	free(sin);
}

void test_invalid_multidots_next() {
	char *sin = locnumberstr("1..3");
	size_t len = strlen(sin);
	clearerrors();
	number n = mknumc(len, sin);
	assert(n == 0, "number null on error");
	assert(poperror() == invalidnumber, "error on 0 length");
	free(sin);
}

void test_invalid_dotonly() {
	char *sin = locnumberstr(".");
	size_t len = strlen(sin);
	clearerrors();
	number n = mknumc(len, sin);
	assert(n == 0, "number null on error");
	assert(poperror() == invalidnumber, "error on 0 length");
	free(sin);
}

void test_invalid_signonly() {
	char *sin = locnumberstr("-");
	size_t len = strlen(sin);
	clearerrors();
	number n = mknumc(len, sin);
	assert(n == 0, "number null on error");
	assert(poperror() == invalidnumber, "sign only");
	free(sin);
}

void test_usescurrentlocale() {
	char *sinfmt = "123%c21";
	size_t len = strlen(sinfmt) - 1;
	char *sin = malloc(len + 1);
	sprintf(sin, sinfmt, getdecchar());
	number n = mknumc(len, sin);
	char *sout = sprintnum(n);
	assert(!strcoll(sout, sin), "uses current locale");
	free(sin);
	free(sout);
	freenum(n);
}

void test_int() {
	char *sin = locnumberstr("321");
	size_t len = strlen(sin);
	number n = mknumc(len, sin);
	char *sout = sprintnum(n);
	assert(!strcoll(sout, sin), "int");
	free(sin);
	free(sout);
	freenum(n);
}

void test_int_sign() {
	char *sin = locnumberstr("-321");
	size_t len = strlen(sin);
	number n = mknumc(len, sin);
	char *sout = sprintnum(n);
	assert(!strcoll(sout, sin), "int");
	free(sin);
	free(sout);
	freenum(n);
}

void test_int_enddot() {
	char *sin = locnumberstr("321.");
	size_t len = strlen(sin);
	number n = mknumc(len, sin);
	char *sout = sprintnum(n);
	assert(!strcoll(sout, "321"), "int with closing dot");
	free(sin);
	free(sout);
	freenum(n);
}

void test_int_enddot_sign() {
	char *sin = locnumberstr("-321.");
	size_t len = strlen(sin);
	number n = mknumc(len, sin);
	char *sout = sprintnum(n);
	assert(!strcoll(sout, "-321"), "int with closing dot");
	free(sin);
	free(sout);
	freenum(n);
}

void test_rational() {
	char *sin = locnumberstr("321.23");
	size_t len = strlen(sin);
	number n = mknumc(len, sin);
	char *sout = sprintnum(n);
	assert(!strcoll(sout, sin), "rational number");
	free(sin);
	free(sout);
	freenum(n);
}

void test_rational_sign() {
	char *sin = locnumberstr("-321.23");
	size_t len = strlen(sin);
	number n = mknumc(len, sin);
	char *sout = sprintnum(n);
	assert(!strcoll(sout, sin), "rational number");
	free(sin);
	free(sout);
	freenum(n);
}

void test_rational_startdot() {
	char *sin = locnumberstr(".23");
	size_t len = strlen(sin);
	number n = mknumc(len, sin);
	char *sout = sprintnum(n);
	assert(!strcoll(sout, "0.23"), "rational number start dot");
	free(sin);
	free(sout);
	freenum(n);
}

void test_rational_startdot_sign() {
	char *sin = locnumberstr("-.23");
	size_t len = strlen(sin);
	number n = mknumc(len, sin);
	char *sout = sprintnum(n);
	assert(!strcoll(sout, "-0.23"), "rational number start dot");
	free(sin);
	free(sout);
	freenum(n);
}

void test_small_int() {
	number n = mknumi(42, 35);
	char *sout = sprintnum(n);
	assert(!strcoll(sout, "1.2"), "short int");
	free(sout);
	freenum(n);
}

void test_small_int_sign() {
	number n = mknumi(42, -35);
	char *sout = sprintnum(n);
	assert(!strcoll(sout, "-1.2"), "short int sign");
	free(sout);
	freenum(n);
}

void test_mknumcsafe_fail() {
	number n = mknumcsafe(1, "a");
	assert(!n, "number: mknumcsafe, fail");
}

void test_mknumcsafe() {
	number n = mknumcsafe(1, "1");
	assert(rawint(n) == 1, "number: mknumcsafe, fail");
	freenum(n);
}

void test_rawint_not_int() {
	number n = mknumi(1, 2);
	clearerrors();
	rawint(n);
	assert(poperror() == numbernotint, "raw int from not int");
	freenum(n);
}

void test_rawint_too_big_positive() {
	char *c = "9847947694876985769576985769857698576";
	size_t len = strlen(c);
	number n = mknumc(len, c);
	clearerrors();
	rawint(n);
	assert(poperror() == numbertoobig, "raw int from too big");
	freenum(n);
}

void test_rawint_too_big_negative() {
	char *c = "-9847947694876985769576985769857698576";
	size_t len = strlen(c);
	number n = mknumc(len, c);
	clearerrors();
	rawint(n);
	assert(poperror() == numbertoobig, "raw int from too big");
	freenum(n);
}

void test_rawint() {
	number n = mknumi(42, 1);
	assert(rawint(n) == 42, "raw int");
	freenum(n);
}

void test_sum_small() {
	number n1 = mknumi(2, 1);
	number n2 = mknumi(3, 1);
	number s = sum(n1, n2);
	assert(rawint(s) == 5, "number: sum small");
	freenum(n1);
	freenum(n2);
	freenum(s);
}

void test_sum() {
	number n1 = mknumc(18, "123456789012345678");
	number n2 = mknumc(18, "876543210987654321");
	number s = sum(n1, n2);
	char *ss = sprintnum(s);
	assert(!strcoll(ss, "999999999999999999"), "number: sum");
	freenum(n1);
	freenum(n2);
	freenum(s);
	free(ss);
}

void test_diff_small() {
	number n1 = mknumi(2, 1);
	number n2 = mknumi(3, 1);
	number s = diff(n1, n2);
	assert(rawint(s) == -1, "number: sum small");
	freenum(n1);
	freenum(n2);
	freenum(s);
}

void test_diff() {
	number n1 = mknumc(18, "999999999999999999");
	number n2 = mknumc(18, "876543210987654321");
	number s = diff(n1, n2);
	char *ss = sprintnum(s);
	assert(!strcoll(ss, "123456789012345678"), "number: sum");
	freenum(n1);
	freenum(n2);
	freenum(s);
	free(ss);
}

void test_isint_small_false() {
	number n = mknumi(3, 42);
	assert(!isint(n), "number: isint, small, false");
	freenum(n);
}

void test_isint_small_true() {
	number n = mknumi(42, 1);
	assert(isint(n), "number: isint, small, true");
	freenum(n);
}

void test_isint_big_false() {
	char *s = "9874596769769844576984.55769479687";
	size_t len = strlen(s);
	number n = mknumc(len, s);
	assert(!isint(n), "number: isint, small, false");
	freenum(n);
}

void test_isint_big_true() {
	char *s = "987459676976984457698455769479687";
	size_t len = strlen(s);
	number n = mknumc(len, s);
	assert(isint(n), "number: isint, small, false");
	freenum(n);
}

void test_issmallint_small_notint() {
	number n = mknumi(3, 42);
	assert(!issmallint(n), "number: issmallint, small, not int");
	freenum(n);
}

void test_issmallint_small_int() {
	number n = mknumi(42, 1);
	assert(issmallint(n), "number: issmallint, small, int");
	freenum(n);
}

void test_issmallint_big_notint() {
	char *c = "9874598765967947695876.4586795679";
	size_t len = strlen(c);
	number n = mknumc(len, c);
	assert(!issmallint(n), "number: issmallint, big, not int");
	freenum(n);
}

void test_issmallint_big_int() {
	char *c = "98745987659679476958764586795679";
	size_t len = strlen(c);
	number n = mknumc(len, c);
	assert(!issmallint(n), "number: issmallint, bit, not int");
	freenum(n);
}

void test_clonenum_small() {
	number n = mknumi(42, 1);
	number cn = clonenum(n);
	assert(rawint(cn) == 42, "number: clonenum, small");
	freenum(cn);
	freenum(n);
}

void test_clonenum_big() {
	char *s = "94876984769747684764878576";
	size_t len = strlen(s);
	number n = mknumc(len, s);
	number cn = clonenum(n);
	char *ss = sprintnum(cn);
	assert(!strcoll(ss, s), "number: clonenum, big");
	free(ss);
	freenum(cn);
	freenum(n);
}

void test_bitor_fail() {
	number n1 = mknumi(42, 1);
	number n2 = mknumi(3, 42);
	clearerrors();
	bitor(n1, n2);
	assert(poperror() == invalidnumber, "number: bitor, fail");
	freenum(n1);
	freenum(n2);
}

void test_bitor() {
	number n1 = mknumi(8, 1);
	number n2 = mknumi(1, 1);
	number ior = bitor(n1, n2);
	assert(rawint(ior) == 9, "number: bitor");
	freenum(n1);
	freenum(n2);
	freenum(ior);
}

void test_comparenum_less() {
	number n1 = mknumi(42, 2);
	number n2 = mknumi(42, 1);
	assert(comparenum(n1, n2) < 0, "number: compare, less");
	freenum(n1);
	freenum(n2);
}

void test_comparenum_eq() {
	number n1 = mknumi(42, 1);
	number n2 = mknumi(42, 1);
	assert(comparenum(n1, n2) == 0, "number: compare, less");
	freenum(n1);
	freenum(n2);
}

void test_comparenum_greater() {
	number n1 = mknumi(42, 1);
	number n2 = mknumi(42, 2);
	assert(comparenum(n1, n2) > 0, "number: compare, less");
	freenum(n1);
	freenum(n2);
}

int main() {
	initsys();
	initmodule_errormock();
	initmodule_number();

	test_invalid_length();
	test_invalid_multidots();
	test_invalid_multidots_next();
	test_invalid_dotonly();
	test_invalid_signonly();
	test_int();
	test_int_sign();
	test_int_enddot();
	test_int_enddot_sign();
	test_rational();
	test_rational_sign();
	test_rational_startdot();
	test_rational_startdot_sign();
	test_usescurrentlocale();
	test_small_int();
	test_small_int_sign();
	test_mknumcsafe_fail();
	test_mknumcsafe();
	test_rawint_not_int();
	test_rawint_too_big_positive();
	test_rawint_too_big_negative();
	test_rawint();
	test_sum_small();
	test_sum();
	test_diff_small();
	test_diff();
	test_isint_small_false();
	test_isint_small_true();
	test_isint_big_false();
	test_isint_big_true();
	test_issmallint_small_notint();
	test_issmallint_small_int();
	test_issmallint_big_notint();
	test_issmallint_big_int();
	test_clonenum_small();
	test_clonenum_big();
	test_bitor_fail();
	test_bitor();
	test_comparenum_less();
	test_comparenum_eq();
	test_comparenum_greater();

	freemodule_number();
	freemodule_errormock();
	return 0;
}
