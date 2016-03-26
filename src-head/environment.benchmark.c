#include <stdlib.h>
#include <string.h>
#include "sys.h"
#include "testing.h"
#include "error.h"
#include "error.mock.h"
#include "number.h"
#include "string.h"
#include "compound-types.h"
#include "io.h"
#include "value.h"
#include "environment.h"

void benchmark_init_free_100() {
	int m = 100;
	for (int i = 0; i < m; i++) {
		environment env = mkenvironment(0);
		freenv(env);
	}
}

void benchmark_init_free_with_parent_100() {
	int m = 100;
	environment parent = mkenvironment(0);

	for (int i = 0; i < m; i++) {
		environment env = mkenvironment(parent);
		freenv(env);
	}

	freenv(parent);
}

void benchmark_defvar_fail_100() {
	int m = 100;

	environment env = mkenvironment(0);
	value notsym = mknumvali(1, 1);
	value val = mknumvali(2, 1);

	for (int i = 0; i < m; i++) {
		defvar(env, notsym, val);
	}

	freenv(env);
	freeval(notsym);
	freeval(val);
	clearerrors();
}

void benchmark_defvar_100() {
	int m = 100;

	environment env = mkenvironment(0);
	char *sname = "some-symbol";
	size_t slen = strlen(sname);
	value sym = mksymval(slen, sname);
	value val = mknumvali(1, 1);

	for (int i = 0; i < m; i++) {
		defvar(env, sym, val);
	}

	freenv(env);
	freeval(sym);
	freeval(val);
}

void benchmark_hasvar_fail_100() {
	int m = 100;

	environment env = mkenvironment(0);
	value notsym = mknumvali(1, 1);
	
	for (int i = 0; i < m; i++) {
		hasvar(env, notsym);
	}

	freenv(env);
	freeval(notsym);
	clearerrors();
}

void benchmark_hasvar_false_100() {
	int m = 100;

	environment env = mkenvironment(0);
	char *sname = "some-symbol";
	size_t slen = strlen(sname);
	value sym = mksymval(slen, sname);
	
	for (int i = 0; i < m; i++) {
		hasvar(env, sym);
	}

	freenv(env);
	freeval(sym);
	clearerrors();
}

void benchmark_hasvar_100() {
	int m = 100;

	environment env = mkenvironment(0);
	char *sname = "some-symbol";
	size_t slen = strlen(sname);
	value sym = mksymval(slen, sname);
	value val = mknumvali(1, 1);
	defvar(env, sym, val);
	
	for (int i = 0; i < m; i++) {
		hasvar(env, sym);
	}

	freenv(env);
	freeval(sym);
	freeval(val);
	clearerrors();
}

void benchmark_getvar_fail_not_symbol_100() {
	int m = 100;

	environment env = mkenvironment(0);
	value notsym = mknumvali(1, 1);

	for (int i = 0; i < m; i++) {
		getvar(env, notsym);
	}

	freenv(env);
	freeval(notsym);
	clearerrors();
}

void benchmark_getvar_fail_not_exists_100() {
	int m = 100;

	environment env = mkenvironment(0);
	char *sname = "some-symbol";
	size_t slen = strlen(sname);
	value sym = mksymval(slen, sname);

	for (int i = 0; i < m; i++) {
		getvar(env, sym);
	}

	freenv(env);
	freeval(sym);
	clearerrors();
}

void benchmark_getvar_100() {
	int m = 100;

	environment env = mkenvironment(0);
	char *sname = "some-symbol";
	size_t slen = strlen(sname);
	value sym = mksymval(slen, sname);
	value val = mknumvali(1, 1);
	defvar(env, sym, val);

	for (int i = 0; i < m; i++) {
		getvar(env, sym);
	}

	freenv(env);
	freeval(sym);
	freeval(val);
}

void benchmark_setvar_fail_not_symbol_100() {
	int m = 100;

	environment env = mkenvironment(0);
	value notsym = mknumvali(1, 1);
	value val = mknumvali(2, 1);

	for (int i = 0; i < m; i++) {
		setvar(env, notsym, val);
	}

	freenv(env);
	freeval(notsym);
	freeval(val);
	clearerrors();
}

void benchmark_setvar_fail_not_exists_100() {
	int m = 100;

	environment env = mkenvironment(0);
	char *sname = "some-symbol";
	size_t slen = strlen(sname);
	value sym = mksymval(slen, sname);
	value val = mknumvali(2, 1);

	for (int i = 0; i < m; i++) {
		setvar(env, sym, val);
	}

	freenv(env);
	freeval(sym);
	freeval(val);
	clearerrors();
}

void benchmark_setvar_100() {
	int m = 100;

	environment env = mkenvironment(0);
	char *sname = "some-symbol";
	size_t slen = strlen(sname);
	value sym = mksymval(slen, sname);
	value val = mknumvali(2, 1);
	defvar(env, sym, val);

	for (int i = 0; i < m; i++) {
		setvar(env, sym, val);
	}

	freenv(env);
	freeval(sym);
	freeval(val);
}

void benchmark_getvar_lookup_100() {
	int m = 100;
	int n = 24;

	environment *env = malloc(n * sizeof(environment));
	char *sname = "some-symbol";
	size_t slen = strlen(sname);
	value sym = mksymval(slen, sname);
	value val = mknumvali(1, 1);
	*(env + n - 1) = mkenvironment(0);
	defvar(*(env + n - 1), sym, val);
	for (int i = n - 1; i > 0; i--) {
		*(env + i - 1) = mkenvironment(*(env + i));
	}

	for (int i = 0; i < m; i++) {
		getvar(*env, sym);
	}

	for (int i = 0; i < n; i++) {
		freenv(*(env + i));
	}

	free(env);
	freeval(sym);
	freeval(val);
}

void benchmark_extenv_100() {
	int m = 100;

	char *sname1 = "symbol1";
	size_t slen1 = strlen(sname1);
	value sym1 = mksymval(slen1, sname1);
	value val1 = mknumvali(1, 1);
	char *sname2 = "symbol2";
	size_t slen2 = strlen(sname2);
	value sym2 = mksymval(slen2, sname2);
	value val2 = mknumvali(2, 1);
	char *sname3 = "symbol3";
	size_t slen3 = strlen(sname3);
	value sym3 = mksymval(slen3, sname3);
	value val3 = mknumvali(3, 1);
	environment parent = mkenvironment(0);
	defvar(parent, sym1, val1);
	value symp1 = mkpairval(sym2, null);
	value symp2 = mkpairval(sym3, symp1);
	value valp1 = mkpairval(val2, null);
	value valp2 = mkpairval(val3, valp1);

	for (int i = 0; i < m; i++) {
		environment env = extenv(parent, symp2, valp2);
		freenv(env);
	}

	freenv(parent);
	freeval(symp2);
	freeval(symp1);
	freeval(valp2);
	freeval(valp1);
	freeval(sym1);
	freeval(val1);
	freeval(sym2);
	freeval(val2);
	freeval(sym3);
	freeval(val3);
}

void benchmark_extend_env_variadic_100() {
	int m = 100;

	char *sname = "symbol";
	size_t slen = strlen(sname);
	value sym = mksymval(slen, sname);
	value val1 = mknumvali(1, 1);
	value val2 = mknumvali(2, 1);
	value val3 = mknumvali(3, 1);
	environment parent = mkenvironment(0);
	value p0 = mkpairval(val3, null);
	value p1 = mkpairval(val2, p0);
	value p2 = mkpairval(val1, p1);

	for (int i = 0; i < m; i++) {
		environment env = extenv(parent, sym, p2);
		freenv(env);
	}

	freeval(p2);
	freeval(p1);
	freeval(p0);
	freenv(parent);
	freeval(val3);
	freeval(val2);
	freeval(val1);
	freeval(sym);
}

int main(int argc, char **argv) {
	initsys();
	initmodule_errormock();
	initmodule_value();
	int n = repeatcount(argc, argv);

	int err = 0;
	err += benchmark(n, &benchmark_init_free_100, "environment: init and free, 100 times");
	err += benchmark(n, &benchmark_init_free_with_parent_100, "environment: init and free with parent, 100 times");
	err += benchmark(n, &benchmark_defvar_fail_100, "environment: defvar, fail, 100 times");
	err += benchmark(n, &benchmark_defvar_100, "environment: defvar, 100 times");
	err += benchmark(n, &benchmark_hasvar_fail_100, "environment: hasvar, fail, 100 times");
	err += benchmark(n, &benchmark_hasvar_false_100, "environment: hasvar, false, 100 times");
	err += benchmark(n, &benchmark_getvar_fail_not_symbol_100, "environment: getvar, fail, not a symbol, 100 times");
	err += benchmark(n, &benchmark_getvar_fail_not_exists_100, "environment: getvar, fail, not defined, 100 times");
	err += benchmark(n, &benchmark_getvar_100, "environment: getvar, 100 times");
	err += benchmark(n, &benchmark_setvar_fail_not_symbol_100, "environment: setvar, fail, not symbol, 100 times");
	err += benchmark(n, &benchmark_setvar_fail_not_exists_100, "environment: setvar, fail, not defined, 100 times");
	err += benchmark(n, &benchmark_setvar_100, "environment: setvar, 100 times");
	err += benchmark(n, &benchmark_getvar_lookup_100, "environment: lookup, 100 times");
	err += benchmark(n, &benchmark_extenv_100, "environment: extend, 100 times");
	err += benchmark(n, &benchmark_extend_env_variadic_100, "environment: extend, variadic, 100 times");

	freemodule_value();
	freemodule_errormock();
	return err;
}
