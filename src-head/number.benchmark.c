#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include "testing.h"
#include "error.h"
#include "error.mock.h"
#include "sys.h"
#include "number.h"

char *bigintc = "32198437694876957694876938749677985769857695769769487498657947";
char *bigratc = "32198437694876957694876938749677985769.857695769769487498657947";

void benchmark_init_release_int_100() {
	int m = 100;

	char *sin = locnumberstr(bigintc);
	size_t len = strlen(sin);

	for (int i = 0; i < m; i++) {
		number n = mknumc(len, sin);
		freenum(n);
	}

	free(sin);
}

void benchmark_init_error_100() {
	int m = 100;

	char *sin = locnumberstr("123456890..0987654321");
	size_t len = strlen(sin);

	for (int i = 0; i < m; i++) {
		mknumc(len, sin);
	}

	free(sin);
	clearerrors();
}

void benchmark_sprint_int_50() {
	int m = 50;

	char *sin = locnumberstr(bigintc);
	size_t len = strlen(sin);
	number n = mknumc(len, sin);

	for (int i = 0; i < m; i++) {
		char *sout = sprintnum(n);
		free(sout);
	}

	freenum(n);
	free(sin);
}

void benchmark_init_release_small_int_100() {
	int m = 100;

	for (int i = 0; i < m; i++) {
		number n = mknumi(42, 1);
		freenum(n);
	}
}

void benchmark_sprint_small_int_50() {
	int m = 50;

	number n = mknumi(42, 1);

	for (int i = 0; i < m; i++) {
		char *sout = sprintnum(n);
		free(sout);
	}

	freenum(n);
}

void benchmark_init_release_rational_100() {
	int m = 100;

	char *sin = locnumberstr(bigratc);
	size_t len = strlen(sin);

	for (int i = 0; i < m; i++) {
		number n = mknumc(len, sin);
		freenum(n);
	}

	free(sin);
}

void benchmark_sprint_rational_50() {
	int m = 50;

	char *sin = locnumberstr(bigratc);
	size_t len = strlen(sin);
	number n = mknumc(len, sin);

	for (int i = 0; i < m; i++) {
		char *sout = sprintnum(n);
		free(sout);
	}

	free(sin);
	freenum(n);
}

void benchmark_init_release_small_rational_100() {
	int m = 100;

	char *sin = locnumberstr(bigratc);
	size_t len = strlen(sin);

	for (int i = 0; i < m; i++) {
		number n = mknumi(42, 35);
		freenum(n);
	}

	free(sin);
}

void benchmark_sprint_small_rational_50() {
	int m = 50;

	number n = mknumi(42, 35);

	for (int i = 0; i < m; i++) {
		char *sout = sprintnum(n);
		free(sout);
	}

	freenum(n);
}

void benchmark_rawint_not_int_100() {
	int m = 100;
	number n = mknumi(1, 2);

	for (int i = 0; i < m; i++) {
		rawint(n);
	}

	freenum(n);
	clearerrors();
}

void benchmark_rawint_too_big_positive_100() {
	int m = 100;

	char *c = "9847947694876985769576985769857698576";
	size_t len = strlen(c);
	number n = mknumc(len, c);

	for (int i = 0; i < m; i++) {
		rawint(n);
	}

	freenum(n);
	clearerrors();
}

void benchmark_rawint_too_big_negative_100() {
	int m = 100;

	char *c = "-9847947694876985769576985769857698576";
	size_t len = strlen(c);
	number n = mknumc(len, c);

	for (int i = 0; i < m; i++) {
		rawint(n);
	}

	freenum(n);
	clearerrors();
}

void benchmark_rawint_100() {
	int m = 100;
	number n = mknumi(42, 1);

	for (int i = 0; i < m; i++) {
		rawint(n);
	}

	freenum(n);
}

void benchmark_sum_small_100() {
	int m = 100;

	number n1 = mknumi(2, 1);
	number n2 = mknumi(3, 1);

	for (int i = 0; i < m; i++) {
		number s = sum(n1, n2);
		freenum(s);
	}

	freenum(n1);
	freenum(n2);
}

void benchmark_sum_100() {
	int m = 100;

	number n1 = mknumc(18, "123456789012345678");
	number n2 = mknumc(18, "876543210987654321");

	for (int i = 0; i < m; i++) {
		number s = sum(n1, n2);
		freenum(s);
	}

	freenum(n1);
	freenum(n2);
}

void benchmark_diff_small_100() {
	int m = 100;

	number n1 = mknumi(2, 1);
	number n2 = mknumi(3, 1);

	for (int i = 0; i < m; i++) {
		number s = diff(n1, n2);
		freenum(s);
	}

	freenum(n1);
	freenum(n2);
}

void benchmark_diff_100() {
	int m = 100;

	number n1 = mknumc(18, "999999999999999999");
	number n2 = mknumc(18, "876543210987654321");

	for (int i = 0; i < m; i++) {
		number s = diff(n1, n2);
		freenum(s);
	}

	freenum(n1);
	freenum(n2);
}

void benchmark_isint_small_false_100() {
	int m = 100;
	number n = mknumi(3, 42);

	for (int i = 0; i < m; i++) {
		isint(n);
	}

	freenum(n);
}

void benchmark_isint_small_true_100() {
	int m = 100;
	number n = mknumi(42, 1);

	for (int i = 0; i < m; i++) {
		isint(n);
	}

	freenum(n);
}

void benchmark_isint_big_false_100() {
	int m = 100;

	char *s = "9874596769769844576984.55769479687";
	size_t len = strlen(s);
	number n = mknumc(len, s);

	for (int i = 0; i < m; i++) {
		isint(n);
	}

	freenum(n);
}

void benchmark_isint_big_true_100() {
	int m = 100;

	char *s = "987459676976984457698455769479687";
	size_t len = strlen(s);
	number n = mknumc(len, s);

	for (int i = 0; i < m; i++) {
		isint(n);
	}

	freenum(n);
}

void benchmark_issmallint_small_notint_100() {
	int m = 100;
	number n = mknumi(3, 42);

	for (int i = 0; i < m; i++) {
		issmallint(n);
	}

	freenum(n);
}

void benchmark_issmallint_small_int_100() {
	int m = 100;
	number n = mknumi(42, 1);

	for (int i = 0; i < m; i++) {
		issmallint(n);
	}

	freenum(n);
}

void benchmark_issmallint_big_notint_100() {
	int m = 100;

	char *c = "9874598765967947695876.4586795679";
	size_t len = strlen(c);
	number n = mknumc(len, c);

	for (int i = 0; i < m; i++) {
		issmallint(n);
	}

	freenum(n);
}

void benchmark_issmallint_big_int_100() {
	int m = 100;

	char *c = "98745987659679476958764586795679";
	size_t len = strlen(c);
	number n = mknumc(len, c);

	for (int i = 0; i < m; i++) {
		issmallint(n);
	}

	freenum(n);
}

void benchmark_clonenum_small_100() {
	int m = 100;
	number n = mknumi(42, 1);

	for (int i = 0; i < m; i++) {
		number cn = clonenum(n);
		freenum(cn);
	}

	freenum(n);
}

void benchmark_clonenum_big_100() {
	int m = 100;

	char *s = "94876984769747684764878576";
	size_t len = strlen(s);
	number n = mknumc(len, s);

	for (int i = 0; i < m; i++) {
		number cn = clonenum(n);
		freenum(cn);
	}

	freenum(n);
}

void benchmark_bitor_fail_100() {
	int m = 100;

	number n1 = mknumi(42, 1);
	number n2 = mknumi(3, 42);

	for (int i = 0; i < m; i++) {
		bitor(n1, n2);
	}

	clearerrors();
	freenum(n1);
	freenum(n2);
}

void benchmark_bitor_100() {
	int m = 100;

	number n1 = mknumi(8, 1);
	number n2 = mknumi(1, 1);
	number ior = bitor(n1, n2);

	for (int i = 0; i < m; i++) {
		rawint(ior);
	}

	freenum(n1);
	freenum(n2);
	freenum(ior);
}

void benchmark_comparenum_less_100() {
	int m = 100;

	number n1 = mknumi(42, 2);
	number n2 = mknumi(42, 1);

	for (int i = 0; i < m; i++) {
		comparenum(n1, n2);
	}

	freenum(n1);
	freenum(n2);
}

void benchmark_comparenum_eq_100() {
	int m = 100;

	number n1 = mknumi(42, 1);
	number n2 = mknumi(42, 1);

	for (int i = 0; i < m; i++) {
		comparenum(n1, n2);
	}

	freenum(n1);
	freenum(n2);
}

void benchmark_comparenum_greater_100() {
	int m = 100;

	number n1 = mknumi(42, 1);
	number n2 = mknumi(42, 2);

	for (int i = 0; i < m; i++) {
		comparenum(n1, n2);
	}

	freenum(n1);
	freenum(n2);
}

void benchmark_mknumcsafe_fail_100() {
	int m = 100;
	for (int i = 0; i < m; i++) {
		mknumcsafe(1, "a");
	}
}

void benchmark_mknumcsafe_100() {
	int m = 100;
	for (int i = 0; i < m; i++) {
		number n = mknumcsafe(1, "1");
		freenum(n);
	}
}

int main(int argc, char **argv) {
	initsys();
	int n = repeatcount(argc, argv);
	initmodule_errormock();
	initmodule_number();

	int err = 0;
	err += benchmark(n, &benchmark_init_release_int_100, "number: init, release int, 100 times");
	err += benchmark(n, &benchmark_init_error_100, "number: init, error, 100 times");
	err += benchmark(n, &benchmark_sprint_int_50, "number: sprint int, 50 times");
	err += benchmark(n, &benchmark_init_release_small_int_100, "number: init, release small int, 100 times");
	err += benchmark(n, &benchmark_sprint_small_int_50, "number: sprint small int, 50 times");
	err += benchmark(n, &benchmark_init_release_rational_100, "number: init, release rational, 100 times");
	err += benchmark(n, &benchmark_sprint_rational_50, "number: sprint rational, 50 times");
	err += benchmark(n, &benchmark_init_release_small_rational_100, "number: init, release small rational, 100 times");
	err += benchmark(n, &benchmark_sprint_small_rational_50, "number: sprint small rational, 50 times");
	err += benchmark(n, &benchmark_rawint_not_int_100, "number: raw int, not int, 100 times");
	err += benchmark(n, &benchmark_rawint_too_big_positive_100, "number: raw int, too big+, 100 times");
	err += benchmark(n, &benchmark_rawint_too_big_negative_100, "number: raw int, too big-, 100 times");
	err += benchmark(n, &benchmark_rawint_100, "number: raw int, 100 times");
	err += benchmark(n, &benchmark_sum_small_100, "number: sum, small, 100 times");
	err += benchmark(n, &benchmark_sum_100, "number: sum, 100 times");
	err += benchmark(n, &benchmark_diff_small_100, "number: diff, small, 100 times");
	err += benchmark(n, &benchmark_diff_100, "number: diff, 100 times");
	err += benchmark(n, &benchmark_isint_small_false_100, "number: isint, small, false, 100 times");
	err += benchmark(n, &benchmark_isint_small_true_100, "number: isint, small, true, 100 times");
	err += benchmark(n, &benchmark_isint_big_false_100, "number: isint, big, false, 100 times");
	err += benchmark(n, &benchmark_isint_big_true_100, "number: isint, big, true, 100 times");
	err += benchmark(n, &benchmark_issmallint_small_notint_100, "number: issmallint, small, not int, 100 times");
	err += benchmark(n, &benchmark_issmallint_small_int_100, "number: issmallint, small, int, 100 times");
	err += benchmark(n, &benchmark_issmallint_big_notint_100, "number: issmallint, big, not int, 100 times");
	err += benchmark(n, &benchmark_issmallint_big_int_100, "number: issmallint, big, int, 100 times");
	err += benchmark(n, &benchmark_clonenum_small_100, "number: clonenum, small, 100 times");
	err += benchmark(n, &benchmark_clonenum_big_100, "number: clonenum, big, 100 times");
	err += benchmark(n, &benchmark_bitor_fail_100, "number: bitor, fail, 100 times");
	err += benchmark(n, &benchmark_bitor_100, "number: bitor, 100 times");
	err += benchmark(n, &benchmark_comparenum_less_100, "number: compare, less, 100 times");
	err += benchmark(n, &benchmark_comparenum_eq_100, "number: compare, eq, 100 times");
	err += benchmark(n, &benchmark_comparenum_greater_100, "number: compare, greater, 100 times");
	err += benchmark(n, &benchmark_mknumcsafe_fail_100, "number: mknumcsafe, fail, 100 times");
	err += benchmark(n, &benchmark_mknumcsafe_100, "number: mknumcsafe, 100 times");

	freemodule_number();
	freemodule_errormock();
	return err;
}
