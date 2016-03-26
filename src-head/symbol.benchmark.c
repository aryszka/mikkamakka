#include <stdlib.h>
#include <string.h>
#include "sys.h"
#include "testing.h"
#include "symbol.h"

void benchmark_init_free_100() {
	int m = 100;

	char *name = "symbol";
	int len = strlen(name);

	symbol s;
	for (int i = 0; i < m; i++) {
		s = mksymbol(len, name);
		freesymbol(s);
	}
}

void benchmark_symbol_name_100() {
	int m = 100;

	char *name = "symbol";
	int len = strlen(name);
	symbol s = mksymbol(len, name);

	for (int i = 0; i < m; i++) {
		symbolname(s);
	}

	freesymbol(s);
}

void benchmark_sprint_100() {
	int m = 100;

	char *name = "symbol";
	int len = strlen(name);
	symbol s = mksymbol(len, name);

	for (int i = 0; i < m; i++) {
		char *ss = sprintsymbol(s);
		free(ss);
	}

	freesymbol(s);
}

void benchmark_symeq_false_100() {
	int m = 100;

	symbol s1 = mksymbol(3, "abc");
	symbol s2 = mksymbol(3, "Abc");

	for (int i = 0; i < m; i++) {
		symeq(s1, s2);
	}

	freesymbol(s1);
	freesymbol(s2);
}

void benchmark_symeq_true_100() {
	int m = 100;

	symbol s1 = mksymbol(3, "abc");
	symbol s2 = mksymbol(3, "abc");

	for (int i = 0; i < m; i++) {
		symeq(s1, s2);
	}

	freesymbol(s1);
	freesymbol(s2);
}

int main(int argc, char **argv) {
	initsys();
	int n = repeatcount(argc, argv);

	int err = 0;
	err += benchmark(n, &benchmark_init_free_100, "symbol: init and free, 100 times");
	err += benchmark(n, &benchmark_symbol_name_100, "symbol: symbol name, 100 times");
	err += benchmark(n, &benchmark_sprint_100, "symbol: sprint, 100 times");
	err += benchmark(n, &benchmark_symeq_false_100, "symbol: eq, false, 100 times");
	err += benchmark(n, &benchmark_symeq_true_100, "symbol: eq, true, 100 times");

	return err;
}
