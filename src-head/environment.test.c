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

void test_init_free() {
	environment env = mkenvironment(0);
	freenv(env);
}

void test_init_free_with_parent() {
	environment parent = mkenvironment(0);
	environment env = mkenvironment(parent);
	freenv(env);
	freenv(parent);
}

void test_defvar_fail() {
	environment env = mkenvironment(0);
	value notsym = mknumvali(1, 1);
	value val = mknumvali(2, 1);
	clearerrors();
	defvar(env, notsym, val);
	assert(poperror() == invalidtype, "env: defvar, fail");
	freenv(env);
	freeval(notsym);
	freeval(val);
}

void test_defvar_hasvar_getvar() {
	environment env = mkenvironment(0);
	char *symstr = "some-symbol";
	size_t symlen = strlen(symstr);
	value sym = mksymval(symlen, symstr);
	value val = mknumvali(1, 1);
	defvar(env, sym, val);
	// assert(hasvar(env, sym), "env: defvar, hasvar");
	// assert(getvar(env, sym) == val, "env: defvar, getvar");
	// freenv(env);
	// freeval(sym);
	// freeval(val);
}

void test_hasvar_fail_not_sym() {
	environment env = mkenvironment(0);
	value notsym = mknumvali(1, 1);
	clearerrors();
	hasvar(env, notsym);
	assert(poperror(), "env: getvar, not a symbol");
	freenv(env);
}

void test_getvar_fail_not_sym() {
	environment env = mkenvironment(0);
	value notsym = mknumvali(1, 1);
	clearerrors();
	getvar(env, notsym);
	assert(poperror(), "env: getvar, not a symbol");
	freenv(env);
}

void test_getvar_fail_not_exists() {
	environment env = mkenvironment(0);
	char *symstr = "some-symbol";
	size_t symlen = strlen(symstr);
	value sym = mksymval(symlen, symstr);
	clearerrors();
	getvar(env, sym);
	assert(poperror(), "env: getvar, not exists");
	freenv(env);
	freeval(sym);
}

void test_setvar_fail_not_sym() {
	environment env = mkenvironment(0);
	value notsym = mknumvali(1, 1);
	value val = mknumvali(1, 1);
	clearerrors();
	setvar(env, notsym, val);
	assert(poperror(), "env: setvar, not sym");
	freenv(env);
	freeval(notsym);
	freeval(val);
}

void test_setvar_fail_not_exists() {
	environment env = mkenvironment(0);
	char *symstr = "some-symbol";
	size_t symlen = strlen(symstr);
	value sym = mksymval(symlen, symstr);
	value val = mknumvali(1, 1);
	clearerrors();
	setvar(env, sym, val);
	assert(poperror(), "env: setvar, not exists");
	freenv(env);
	freeval(sym);
	freeval(val);
}

void test_setvar() {
	environment env = mkenvironment(0);
	char *symstr = "some-symbol";
	size_t symlen = strlen(symstr);
	value sym = mksymval(symlen, symstr);
	value val = mknumvali(1, 1);
	defvar(env, sym, val);
	value newval = mknumvali(2, 1);
	clearerrors();
	setvar(env, sym, newval);
	assert(getvar(env, sym) == newval, "env: setvar");
	freenv(env);
	freeval(sym);
	freeval(val);
	freeval(newval);
}

void test_hasvar_getvar_parent() {
	environment parent = mkenvironment(0);
	environment env = mkenvironment(parent);
	char *symstr = "some-symbol";
	size_t symlen = strlen(symstr);
	value sym = mksymval(symlen, symstr);
	value val = mknumvali(1, 1);
	defvar(parent, sym, val);
	assert(hasvar(env, sym), "env: has var in parent");
	assert(getvar(env, sym) == val, "env: get var from parent");
	freenv(env);
	freenv(parent);
	freeval(sym);
	freeval(val);
}

void test_setvar_in_parent() {
	environment parent = mkenvironment(0);
	environment env = mkenvironment(parent);
	char *symstr = "some-symbol";
	size_t symlen = strlen(symstr);
	value sym = mksymval(symlen, symstr);
	value val = mknumvali(1, 1);
	defvar(parent, sym, val);
	value newval = mknumvali(2, 1);
	setvar(env, sym, newval);
	assert(hasvar(parent, sym), "env: has var in parent");
	assert(getvar(parent, sym) == newval, "env: get var from parent");
	freenv(env);
	freenv(parent);
	freeval(sym);
	freeval(val);
	freeval(newval);
}

void test_override_var_of_parent() {
	environment parent = mkenvironment(0);
	environment env = mkenvironment(parent);
	char *symstr = "some-symbol";
	size_t symlen = strlen(symstr);
	value sym = mksymval(symlen, symstr);
	value val = mknumvali(1, 1);
	defvar(parent, sym, val);
	value newval = mknumvali(2, 1);
	defvar(env, sym, newval);
	assert(hasvar(parent, sym), "env: has var in parent");
	assert(getvar(parent, sym) == val, "env: get var from parent");
	assert(hasvar(env, sym), "env: has overridden val");
	assert(getvar(env, sym) == newval, "env: get overridden val");
	freenv(env);
	freenv(parent);
	freeval(sym);
	freeval(val);
	freeval(newval);
}

void test_extend_env() {
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
	environment env = extenv(parent, symp2, valp2);
	assert(getvar(env, sym1) == val1, "env: extenv, parent var");
	assert(getvar(env, sym2) == val2, "env: extenv, ext var");
	assert(getvar(env, sym3) == val3, "env: extenv, ext var");
	freenv(env);
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

void test_extend_env_variadic() {
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
	environment env = extenv(parent, sym, p2);
	value var = getvar(env, sym);
	assert(carval(var) == val1, "env: extenv, variadic, first");
	assert(carval(cdrval(var)) == val2, "env: extenv, variadic, second");
	assert(carval(cdrval(cdrval(var))) == val3, "env: extenv, variadic, third");
	assert(cdrval(cdrval(cdrval(var))) == null, "env: extenv, variadic, end");
	freenv(env);
	freeval(p2);
	freeval(p1);
	freeval(p0);
	freenv(parent);
	freeval(val3);
	freeval(val2);
	freeval(val1);
	freeval(sym);
}

int main() {
	initsys();
	initmodule_errormock();
	initmodule_value();

	test_init_free();
	test_init_free_with_parent();
	test_defvar_fail();
	test_defvar_hasvar_getvar();
	// test_hasvar_fail_not_sym();
	// test_getvar_fail_not_sym();
	// test_getvar_fail_not_exists();
	// test_setvar_fail_not_sym();
	// test_setvar_fail_not_exists();
	// test_setvar();
	// test_hasvar_getvar_parent();
	// test_setvar_in_parent();
	// test_override_var_of_parent();
	// test_extend_env();
	// test_extend_env_variadic();

	freemodule_value();
	freemodule_errormock();
	return 0;
}
