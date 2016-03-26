#include <stdlib.h>
#include <string.h>
#include "sys.h"
#include "error.h"
#include "testing.h"
#include "number.h"
#include "string.h"
#include "compound-types.h"
#include "io.h"
#include "value.h"
#include "sprint-list.h"

void test_null() {
	char *s = sprintlistraw(null);
	assert(!strcoll(s, "()"), "sprint list: null");
	free(s);
}

void test_pair() {
	value n0 = mknumvali(1, 1);
	value n1 = mknumvali(2, 1);
	value p = mkpairval(n0, n1);
	char *s = sprintlistraw(p);
	assert(!strcoll(s, "(1 . 2)"), "sprint list: pair");
	free(s);
	freeval(p);
	freeval(n0);
	freeval(n1);
}

void test_list() {
	value n0 = mknumvali(1, 1);
	value n1 = mknumvali(2, 1);
	value n2 = mknumvali(3, 1);
	value p0 = mkpairval(n2, null);
	value p1 = mkpairval(n1, p0);
	value p2 = mkpairval(n0, p1);
	char *s = sprintlistraw(p2);
	assert(!strcoll(s, "(1 2 3)"), "sprint list: list");
	free(s);
	freeval(p2);
	freeval(p1);
	freeval(p0);
	freeval(n0);
	freeval(n1);
	freeval(n2);
}

void test_embedded_list() {
	value n0 = mknumvali(1, 1);
	value n1 = mknumvali(2, 1);
	value n2 = mknumvali(3, 1);
	value n3 = mknumvali(4, 1);
	value n4 = mknumvali(5, 1);
	value n5 = mknumvali(6, 1);
	value n6 = mknumvali(7, 1);
	value p0 = mkpairval(n6, null);
	value p1 = mkpairval(n5, p0);
	value p2 = mkpairval(n4, null);
	value p3 = mkpairval(n3, p2);
	value p4 = mkpairval(n2, p3);
	value p5 = mkpairval(p4, p1);
	value p6 = mkpairval(n1, p5);
	value p7 = mkpairval(n0, p6);
	char *s = sprintlistraw(p7);
	assert(!strcoll(s, "(1 2 (3 4 5) 6 7)"), "sprint list: embedded list");
	free(s);
	freeval(p7);
	freeval(p6);
	freeval(p5);
	freeval(p4);
	freeval(p3);
	freeval(p2);
	freeval(p1);
	freeval(p0);
	freeval(n0);
	freeval(n1);
	freeval(n2);
	freeval(n3);
	freeval(n4);
	freeval(n5);
	freeval(n6);
}

void test_embedded_list_as_first() {
	value n0 = mknumvali(1, 1);
	value n1 = mknumvali(2, 1);
	value n2 = mknumvali(3, 1);
	value p0 = mkpairval(n2, null);
	value p1 = mkpairval(n1, p0);
	value p2 = mkpairval(n0, null);
	value p3 = mkpairval(p2, p1);
	char *s = sprintlistraw(p3);
	assert(!strcoll(s, "((1) 2 3)"), "sprint list: embedded list as first");
	free(s);
	freeval(p2);
	freeval(p1);
	freeval(p0);
	freeval(n0);
	freeval(n1);
	freeval(n2);
}

void test_embedded_list_as_last() {
	value n0 = mknumvali(1, 1);
	value n1 = mknumvali(2, 1);
	value n2 = mknumvali(3, 1);
	value p0 = mkpairval(n2, null);
	value p1 = mkpairval(p0, null);
	value p2 = mkpairval(n1, p1);
	value p3 = mkpairval(n0, p2);
	char *s = sprintlistraw(p3);
	assert(!strcoll(s, "(1 2 (3))"), "sprint list: embedded list as last");
	free(s);
	freeval(p2);
	freeval(p1);
	freeval(p0);
	freeval(n0);
	freeval(n1);
	freeval(n2);
}

void test_pyramid() {
	value n0 = mknumvali(1, 1);
	value n1 = mknumvali(2, 1);
	value n2 = mknumvali(3, 1);
	value n3 = mknumvali(4, 1);
	value n4 = mknumvali(5, 1);
	value p0 = mkpairval(n4, null);
	value p1 = mkpairval(n3, null);
	value p2 = mkpairval(n2, null);
	value p3 = mkpairval(p2, p1);
	value p4 = mkpairval(n1, p3);
	value p5 = mkpairval(p4, p0);
	value p6 = mkpairval(n0, p5);
	char *s = sprintlistraw(p6);
	assert(!strcoll(s, "(1 (2 (3) 4) 5)"), "sprint list: pyramid");
	free(s);
	freeval(p6);
	freeval(p5);
	freeval(p4);
	freeval(p3);
	freeval(p2);
	freeval(p1);
	freeval(p0);
	freeval(n0);
	freeval(n1);
	freeval(n2);
	freeval(n3);
	freeval(n4);
}

void test_irregular_list() {
	value n0 = mknumvali(1, 1);
	value n1 = mknumvali(2, 1);
	value n2 = mknumvali(3, 1);
	value p0 = mkpairval(n1, n2);
	value p1 = mkpairval(n0, p0);
	char *s = sprintlistraw(p1);
	assert(!strcoll(s, "(1 2 . 3)"), "sprint list: irregular");
	free(s);
	freeval(p1);
	freeval(p0);
	freeval(n0);
	freeval(n1);
	freeval(n2);
}

int main() {
	initsys();
	initmodule_number();
	initmodule_value();

	test_null();
	test_pair();
	test_list();
	test_embedded_list();
	test_embedded_list_as_first();
	test_embedded_list_as_last();
	test_pyramid();
	test_irregular_list();

	freemodule_value();
	freemodule_number();
	return 0;
}
