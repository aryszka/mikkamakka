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
#include "procedure.h"

value testprimitive(value args) {
	return carval(args);
}

void test_init_free_primitive() {
	procedure p = mkprimitiveproc(&testprimitive);
	freeproc(p);
}

void test_init_free_compiled() {
	value label = mknumvali(1, 1);
	environment env = mkenvironment(0);
	procedure p = mkcompiledproc(label, env);
	freeproc(p);
	freenv(env);
	freeval(label);
}

void test_is_primitive_false() {
	value label = mknumvali(1, 1);
	environment env = mkenvironment(0);
	procedure p = mkcompiledproc(label, env);
	assert(!isprimitive(p), "proc: is primtivie false");
	freeproc(p);
	freeval(label);
	freenv(env);
}

void test_is_primitive_true() {
	procedure p = mkprimitiveproc(testprimitive);
	assert(isprimitive(p), "proc: is primtivie true");
	freeproc(p);
}

void test_apply_primitive_fail() {
	value label = mknumvali(1, 1);
	environment env = mkenvironment(0);
	procedure p = mkcompiledproc(label, env);
	value n = mknumvali(1, 1);
	value args = mkpairval(n, null);
	clearerrors();
	applyprimitive(p, args);
	assert(poperror() == invalidtype, "proc: apply primitive, fail");
	freeval(args);
	freeval(n);
	freeproc(p);
	freeval(label);
	freenv(env);
}

void test_apply_primitive() {
	procedure p = mkprimitiveproc(testprimitive);
	value n = mknumvali(1, 1);
	value args = mkpairval(n, null);
	value r = applyprimitive(p, args);
	assert(r == n, "proc: apply primitive");
	freeval(args);
	freeval(n);
	freeproc(p);
}

void test_proclabel_fail() {
	procedure p = mkprimitiveproc(testprimitive);
	clearerrors();
	proclabel(p);
	assert(poperror() == invalidtype, "proc: label, fail");
	freeproc(p);
}

void test_proclabel() {
	value label = mknumvali(1, 1);
	environment env = mkenvironment(0);
	procedure p = mkcompiledproc(label, env);
	assert(proclabel(p) == label, "proc: label");
	freeproc(p);
	freeval(label);
	freenv(env);
}

void test_procenv_fail() {
	procedure p = mkprimitiveproc(testprimitive);
	clearerrors();
	procenv(p);
	assert(poperror() == invalidtype, "proc: label, fail");
	freeproc(p);
}

void test_procenv() {
	value label = mknumvali(1, 1);
	environment env = mkenvironment(0);
	procedure p = mkcompiledproc(label, env);
	assert(procenv(p) == env, "proc: label");
	freeproc(p);
	freeval(label);
	freenv(env);
}

int main() {
	initsys();
	initmodule_value();
	initmodule_errormock();

	test_init_free_primitive();
	test_init_free_compiled();
	test_is_primitive_false();
	test_is_primitive_true();
	test_apply_primitive_fail();
	test_apply_primitive();
	test_proclabel_fail();
	test_procenv();
	test_procenv_fail();

	freemodule_errormock();
	freemodule_value();
	return 0;
}
