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
	return null;
}

void benchmark_init_free_primitive_100() {
	int m = 100;
	for (int i = 0; i < m; i++) {
		procedure p = mkprimitiveproc(&testprimitive);
		freeproc(p);
	}
}

void benchmark_init_free_compiled_100() {
	int m = 100;

	value label = mknumvali(1, 1);
	environment env = mkenvironment(0);

	for (int i = 0; i < m; i++) {
		procedure p = mkcompiledproc(label, env);
		freeproc(p);
	}

	freeval(label);
	freenv(env);
}

void benchmark_is_primitive_false_100() {
	int m = 100;

	value label = mknumvali(1, 1);
	environment env = mkenvironment(0);
	procedure p = mkcompiledproc(label, env);

	for (int i = 0; i < m; i++) {
		isprimitive(p);
	}

	freeproc(p);
	freeval(label);
	freenv(env);
}

void benchmark_is_primitive_true_100() {
	int m = 100;
	procedure p = mkprimitiveproc(testprimitive);

	for (int i = 0; i < m; i++) {
		isprimitive(p);
	}

	freeproc(p);
}

void benchmark_apply_primitive_fail_100() {
	int m = 100;

	value label = mknumvali(1, 1);
	environment env = mkenvironment(0);
	procedure p = mkcompiledproc(label, env);
	value n = mknumvali(1, 1);
	value args = mkpairval(n, null);

	for (int i = 100; i < m; i++) {
		applyprimitive(p, args);
	}

	freeval(args);
	freeval(n);
	freeproc(p);
	freenv(env);
	freeval(label);
	clearerrors();
}

void benchmark_apply_primitive_100() {
	int m = 100;

	procedure p = mkprimitiveproc(testprimitive);
	value n = mknumvali(1, 1);
	value args = mkpairval(n, null);

	for (int i = 100; i < m; i++) {
		applyprimitive(p, args);
	}

	freeval(args);
	freeval(n);
	freeproc(p);
}

void benchmark_proclabel_fail_100() {
	int m = 100;
	procedure p = mkprimitiveproc(testprimitive);

	for (int i = 0; i < m; i++) {
		proclabel(p);
	}

	freeproc(p);
	clearerrors();
}

void benchmark_proclabel_100() {
	int m = 100;

	value label = mknumvali(1, 1);
	environment env = mkenvironment(0);
	procedure p = mkcompiledproc(label, env);

	for (int i = 0; i < m; i++) {
		proclabel(p);
	}

	freeproc(p);
	freeval(label);
	freenv(env);
}

void benchmark_procenv_fail_100() {
	int m = 100;
	procedure p = mkprimitiveproc(testprimitive);

	for (int i = 0; i < m; i++) {
		procenv(p);
	}

	freeproc(p);
	clearerrors();
}

void benchmark_procenv_100() {
	int m = 100;

	value label = mknumvali(1, 1);
	environment env = mkenvironment(0);
	procedure p = mkcompiledproc(label, env);

	for (int i = 0; i < m; i++) {
		procenv(p);
	}

	freeproc(p);
	freeval(label);
	freenv(env);
}

int main(int argc, char **argv) {
	initsys();
	initmodule_errormock();
	initmodule_value();

	int n = repeatcount(argc, argv);

	int err = 0;
	err += benchmark(n, benchmark_init_free_primitive_100, "procedure: init free primitive, 100 times");
	err += benchmark(n, benchmark_init_free_compiled_100, "procedure: init free compiled, 100 times");
	err += benchmark(n, benchmark_is_primitive_false_100, "procedure: is primitive false, 100 times");
	err += benchmark(n, benchmark_is_primitive_true_100, "procedure: is primitive true, 100 times");
	err += benchmark(n, benchmark_apply_primitive_fail_100, "procedure: apply primitive, fail, 100 times");
	err += benchmark(n, benchmark_apply_primitive_100, "procedure: apply primitive, 100 times");
	err += benchmark(n, benchmark_proclabel_fail_100, "procedure: proclabel, fail, 100 times");
	err += benchmark(n, benchmark_proclabel_100, "procedure: proclabel, 100 times");
	err += benchmark(n, benchmark_procenv_fail_100, "procedure: procenv, fail, 100 times");
	err += benchmark(n, benchmark_procenv_100, "procedure: procenv, 100 times");

	freemodule_value();
	freemodule_errormock();
	return err;
}
