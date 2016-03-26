#include "sys.h"
#include "testing.h"
#include "error.h"
#include "error.mock.h"
#include "regex.h"

#include <stdio.h>
#include <string.h>

void test_init_free() {
	regex rx = mkregex(3, "123", 0);
	freeregex(rx);
}

void test_error_compiling() {
	clearerrors();
	regex rx = mkregex(2, "][", 0);
	assert(rx == 0, "regex null on error");
	assert(poperror() == invalidregex, "regex compile error code");
}

void test_init_with_flags() {
	regex rx = mkregex(3, "123", rxignorecase | rxmultiline);
	freeregex(rx);
}

void test_simple_match() {
	regex rx = mkregex(3, "123", 0);

	match m = matchrx(rx, 6, "012345");
	assert(!!m, "regex matches simple string");
	assert(matchlen(m) == 1, "regex simple match length");

	submatch sm = smatch(m, 0);
	assert(sm.index == 1, "regex simple match index");
	assert(sm.len == 3, "regex simple match len");

	freematch(m);
	freeregex(rx);
}

void test_match_with_capture() {
	regex rx = mkregex(10, "12(3(4)5)6", 0);

	match m = matchrx(rx, 10, "0123456789");
	assert(!!m, "regex captures groups");
	assert(matchlen(m) == 3, "regex capture group count");

	submatch sm;

	sm = smatch(m, 0);
	assert(sm.index == 1, "main match index");
	assert(sm.len == 6, "main match length");

	sm = smatch(m, 1);
	assert(sm.index == 3, "capture group index");
	assert(sm.len == 3, "capture group length");

	sm = smatch(m, 2);
	assert(sm.index == 4, "capture group index");
	assert(sm.len == 1, "capture group length");

	freematch(m);
	freeregex(rx);
}

void test_match_fail() {
	regex rx = mkregex(3, "321", 0);

	match m = matchrx(rx, 6, "012345");
	assert(!m, "regex simple match fail");

	freeregex(rx);
}

int main() {
	initsys();
	initmodule_errormock();

	test_init_free();
	test_error_compiling();
	test_init_with_flags();
	test_simple_match();
	test_match_with_capture();
	test_match_fail();

	freemodule_errormock();
	return 0;
}
