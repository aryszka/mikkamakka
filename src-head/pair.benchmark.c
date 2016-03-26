#include "sys.h"
#include "error.h"
#include "testing.h"
#include "number.h"
#include "string.h"
#include "compound-types.h"
#include "io.h"
#include "value.h"
#include "pair.h"

void benchmark_init_free_100() {
	int m = 100;

	pair p;
	value n1 = mknumvali(1, 1);
	value n2 = mknumvali(2, 1);
	for (int i = 0; i < m; i++) {
		p = mkpair(n1, n2);
		freepair(p);
	}

	freeval(n1);
	freeval(n2);
}

void benchmark_car_100() {
	int m = 100;

	value n1 = mknumvali(1, 1);
	value n2 = mknumvali(2, 1);
	pair p = mkpair(n1, n2);
	for (int i = 0; i < m; i++) {
		car(p);
	}

	freepair(p);
	freeval(n1);
	freeval(n2);
}

void benchmark_cdr_100() {
	int m = 100;

	value n1 = mknumvali(1, 1);
	value n2 = mknumvali(2, 1);
	pair p = mkpair(n1, n2);
	for (int i = 0; i < m; i++) {
		cdr(p);
	}

	freepair(p);
	freeval(n1);
	freeval(n2);
}

int main(int argc, char **argv) {
	initsys();
	initmodule_value();
	int n = repeatcount(argc, argv);

	int err = 0;
	err += benchmark(n, &benchmark_init_free_100, "pair: init free 100 times");
	err += benchmark(n, &benchmark_car_100, "pair: car 100 times");
	err += benchmark(n, &benchmark_cdr_100, "pair: cdr 100 times");

	freemodule_value();
	return err;
}
