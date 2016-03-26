#include "sys.h"
#include "testing.h"
#include "error.h"
#include "error.mock.h"
#include "regex.h"

void benchmark_init_free_100() {
	int m = 100;
	for (int i = 0; i < m; i++) {
		regex rx = mkregex(3, "123", 0);
		freeregex(rx);
	}
}

void benchmark_error_compiling_100() {
	int m = 100;
	for (int i = 0; i < m; i++) {
		mkregex(2, "][", 0);
	}

	clearerrors();
}

void benchmark_init_with_flags_100() {
	int m = 100;
	rxflags flags = rxignorecase | rxmultiline;

	for (int i = 0; i < m; i++) {
		regex rx = mkregex(3, "123", flags);
		freeregex(rx);
	}
}

void benchmark_match_simple_string_100() {
	int m = 100;
	regex rx = mkregex(3, "123", 0);

	for (int i = 0; i < m; i++) {
		match m = matchrx(rx, 6, "012345");
		freematch(m);
	}

	freeregex(rx);
}

void benchmark_match_with_capture_groups_100() {
	int m = 100;
	regex rx = mkregex(10, "12(3(4)5)6", 0);

	for (int i = 0; i < m; i++) {
		match m = matchrx(rx, 10, "0123456789");
		freematch(m);
	}

	freeregex(rx);
}

void benchmark_match_len_100() {
	int m = 100;

	regex rx = mkregex(3, "123", 0);
	match match = matchrx(rx, 6, "012345");

	for (int i = 0; i < m; i++) {
		matchlen(match);
	}

	freematch(match);
	freeregex(rx);
}

void benchmark_smatch_100() {
	int m = 100;

	regex rx = mkregex(3, "123", 0);
	match match = matchrx(rx, 6, "012345");

	for (int i = 0; i < m; i++) {
		smatch(match, 0);
	}

	freematch(match);
	freeregex(rx);
}

void benchmark_match_fail_100() {
	int m = 100;
	regex rx = mkregex(3, "321", 0);

	for (int i = 0; i < m; i++) {
		matchrx(rx, 6, "012345");
	}

	freeregex(rx);
}

int main(int argc, char **argv) {
	initsys();
	initmodule_errormock();

	int n = repeatcount(argc, argv);

	int err = 0;
	err += benchmark(n, &benchmark_init_free_100, "regex: init and free, 100 times");
	err += benchmark(n, &benchmark_error_compiling_100, "regex: error compiling, 100 times");
	err += benchmark(n, &benchmark_init_with_flags_100, "regex: init with flags, 100 times");
	err += benchmark(n, &benchmark_match_simple_string_100, "regex: match simple string, 100 times");
	err += benchmark(n, &benchmark_match_with_capture_groups_100, "regex match with capture groups, 100 times");
	err += benchmark(n, &benchmark_match_len_100, "regex: match len, 100 times");
	err += benchmark(n, &benchmark_smatch_100, "regex: submatch, 100 times");
	err += benchmark(n, &benchmark_match_fail_100, "regex: match fail, 100 times");

	freemodule_errormock();
	return err;
}
