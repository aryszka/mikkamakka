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

void benchmark_null_100() {
	int m = 100;
	for (int i = 0; i < m; i++) {
		char *s = sprintlistraw(null);
		free(s);
	}
}

void benchmark_pair_25() {
	int m = 25;

	value n0 = mknumvali(1, 1);
	value n1 = mknumvali(2, 1);
	value p = mkpairval(n0, n1);

	for (int i = 0; i < m; i++) {
		char *s = sprintlistraw(p);
		free(s);
	}

	freeval(p);
	freeval(n0);
	freeval(n1);
}

void benchmark_list_16() {
	int m = 16;

	value n0 = mknumvali(1, 1);
	value n1 = mknumvali(2, 1);
	value n2 = mknumvali(3, 1);
	value p0 = mkpairval(n2, null);
	value p1 = mkpairval(n1, p0);
	value p2 = mkpairval(n0, p1);

	for (int i = 0; i < m; i++) {
		char *s = sprintlistraw(p2);
		free(s);
	}

	freeval(p2);
	freeval(p1);
	freeval(p0);
	freeval(n0);
	freeval(n1);
	freeval(n2);
}

void benchmark_embedded_list_16() {
	int m = 16;

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

	for (int i = 0; i < m; i++) {
		char *s = sprintlistraw(p7);
		free(s);
	}

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

void benchmark_pyramid_16() {
	int m = 16;

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

	for (int i = 0; i < m; i++) {
		char *s = sprintlistraw(p6);
		free(s);
	}

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

void benchmark_irregular_list_16() {
	int m = 16;

	value n0 = mknumvali(1, 1);
	value n1 = mknumvali(2, 1);
	value n2 = mknumvali(3, 1);
	value p0 = mkpairval(n1, n2);
	value p1 = mkpairval(n0, p0);

	for (int i = 0; i < m; i++) {
		char *s = sprintlistraw(p1);
		free(s);
	}

	freeval(p1);
	freeval(p0);
	freeval(n0);
	freeval(n1);
	freeval(n2);
}

int main(int argc, char **argv) {
	initsys();
	initmodule_number();
	initmodule_value();
	int n = repeatcount(argc, argv);

	int err = 0;
	err += benchmark(n, &benchmark_null_100, "sprint list: null, 100 times");
	err += benchmark(n, &benchmark_pair_25, "sprint list: pair, 25 times");
	err += benchmark(n, &benchmark_list_16, "sprint list: list, 16 times");
	err += benchmark(n, &benchmark_embedded_list_16, "sprint list: embedded list, 16 times");
	err += benchmark(n, &benchmark_pyramid_16, "sprint list: pyramid, 16 times");
	err += benchmark(n, &benchmark_irregular_list_16, "sprint list: irregular list, 16 times");

	freemodule_value();
	freemodule_number();

	return err;
}
