#include <stdlib.h>
#include "sys.h"
#include "testing.h"
#include "symbol-table.h"

void benchmark_init_free_100() {
	int m = 100;
	for (int i = 0; i < m; i++) {
		symtable st = mksymtable();
		freesymtable(st);
	}
}

void benchmark_hassym_false_100() {
	int m = 100;

	symtable st = mksymtable();
	char *sname = "some-symbol";
	int slen = strlen(sname);

	for (int i = 0; i < m; i++) {
		hassym(st, slen, sname);
	}

	freesymtable(st);
}

// void benchmark_hassym_100() {
// 	int m = 100;
// 
// 	symtable st = mksymtable();
// 	char *sname = "some-symbol";
// 	int slen = strlen(sname);
// 	value sym = mksymval(slen, sname);
// 	value val = mknumvali(1, 1);
// 	setsym(st, sym, val);
// 
// 	for (int i = 0; i < m; i++) {
// 		hassym(st, sym);
// 	}
// 
// 	freesymtable(st);
// 	freeval(sym);
// 	freeval(val);
// 	clearerrors();
// }
// 
// void benchmark_getsym_fail_not_symbol_100() {
// 	int m = 100;
// 
// 	symtable st = mksymtable();
// 	value notsym = mknumvali(1, 1);
// 
// 	for (int i = 0; i < m; i++) {
// 		getsym(st, notsym);
// 	}
// 
// 	freesymtable(st);
// 	freeval(notsym);
// 	clearerrors();
// }
// 
// void benchmark_getsym_fail_not_exists_100() {
// 	int m = 100;
// 
// 	symtable st = mksymtable();
// 	char *sname = "some-symbol";
// 	int slen = strlen(sname);
// 	value sym = mksymval(slen, sname);
// 
// 	for (int i = 0; i < m; i++) {
// 		getsym(st, sym);
// 	}
// 
// 	freesymtable(st);
// 	freeval(sym);
// 	clearerrors();
// }
// 
// void benchmark_getsym_100() {
// 	int m = 100;
// 
// 	symtable st = mksymtable();
// 	char *sname = "some-symbol";
// 	int slen = strlen(sname);
// 	value sym = mksymval(slen, sname);
// 	value val = mknumvali(1, 1);
// 	setsym(st, sym, val);
// 
// 	for (int i = 0; i < m; i++) {
// 		getsym(st, sym);
// 	}
// 
// 	freesymtable(st);
// 	freeval(sym);
// 	freeval(val);
// 	clearerrors();
// }
// 
// void benchmark_setsym_fail_100() {
// 	int m = 100;
// 
// 	symtable st = mksymtable();
// 	value notsym = mknumvali(1, 1);
// 	value val = mknumvali(2, 1);
// 
// 	for (int i = 0; i < m; i++) {
// 		setsym(st, notsym, val);
// 	}
// 
// 	freesymtable(st);
// 	freeval(notsym);
// 	freeval(val);
// 	clearerrors();
// }
// 
// void benchmark_setsym_100() {
// 	int m = 100;
// 
// 	symtable st = mksymtable();
// 	char *sname = "some-symbol";
// 	int slen = strlen(sname);
// 	value sym = mksymval(slen, sname);
// 	value val = mknumvali(2, 1);
// 
// 	for (int i = 0; i < m; i++) {
// 		setsym(st, sym, val);
// 	}
// 
// 	freesymtable(st);
// 	freeval(sym);
// 	freeval(val);
// 	clearerrors();
// }

int main(int argc, char **argv) {
	initsys();
	int n = repeatcount(argc, argv);

	int err = 0;
	err += benchmark(n, &benchmark_init_free_100, "symtable: init and free, 100 times");
	err += benchmark(n, &benchmark_hassym_false_100, "symtable: hassym, false, 100 times");
	// err += benchmark(n, &benchmark_hassym_100, "symtable: hassym, 100 times");
	// err += benchmark(n, &benchmark_getsym_fail_not_symbol_100, "symtable: getsym, fail, not a symbol, 100 times");
	// err += benchmark(n, &benchmark_getsym_fail_not_exists_100, "symtable: getsym, fail, not exists, 100 times");
	// err += benchmark(n, &benchmark_getsym_100, "symtable: getsym, 100 times");
	// err += benchmark(n, &benchmark_setsym_fail_100, "symtable: setsym, fail, 100 times");
	// err += benchmark(n, &benchmark_setsym_100, "symtable: setsym, 100 times");

	return err;
}
