#include <stdlib.h>
#include <string.h>
#include "sys.h"
#include "testing.h"
#include "error.h"
#include "error.mock.h"
#include "number.h"
#include "string.h"
#include "compound-types.h"
#include "environment.h"
#include "sysio.h"
#include "io.h"
#include "value.h"

value testprimitive(value args) {
	return null;
}

void benchmark_init_free_symbol_100() {
	int m = 100;

	char *name = "symbol";
	int len = strlen(name);

	for (int i = 0; i < m; i++) {
		value s = mksymval(len, name);
		freeval(s);
	}
}

void benchmark_init_free_number_100() {
	int m = 100;
	number n = mknumi(1, 1);

	for (int i = 0; i < m; i++) {
		value v = mknumval(n);
		freeval(v);
	}

	freenum(n);
}

void benchmark_init_free_number_char_100() {
	int m = 100;

	char *c = "12.3";
	int len = strlen(c);

	for (int i = 0; i < m; i++) {
		value n = mknumvalc(len, c);
		freeval(n);
	}
}

void benchmark_init_free_number_int_100() {
	int m = 100;


	for (int i = 0; i < m; i++) {
		value n = mknumvali(3, 4);
		freeval(n);
	}
}

void benchmark_init_free_string_100() {
	int m = 100;

	char *raw = "some string";
	int len = strlen(raw);

	for (int i = 0; i < m; i++) {
		value s = mkstringvalc(len, raw);
		freeval(s);
	}
}

void benchmark_issymtype_false_100() {
	int m = 100;
	value n = mknumvali(1, 1);

	for (int i = 0; i < m; i++) {
		issymtype(n);
	}

	freeval(n);
}

void benchmark_issymtype_100() {
	int m = 100;

	char *sname = "some-symbol";
	int slen = strlen(sname);
	value sym = mksymval(slen, sname);

	for (int i = 0; i < m; i++) {
		issymtype(false);
	}

	freeval(sym);
}

void benchmark_issymval_false_100() {
	int m = 100;
	value n = mknumvali(1, 1);

	for (int i = 0; i < m; i++) {
		isfalseval(n);
	}

	freeval(n);
}

void benchmark_issymval_100() {
	int m = 100;

	char *sname = "some-symbol";
	int slen = strlen(sname);
	value sym = mksymval(slen, sname);

	for (int i = 0; i < m; i++) {
		isfalseval(sym);
	}

	freeval(sym);
}

void benchmark_isnumtype_false_100() {
	int m = 100;
	value v = mkstringvalc(1, "s");

	for (int i = 0; i < m; i++) {
		isnumtype(v);
	}

	freeval(v);
}

void benchmark_isnumtype_100() {
	int m = 100;
	value v = mknumvali(1, 1);

	for (int i = 0; i < m; i++) {
		isnumtype(v);
	}

	freeval(v);
}

void benchmark_isnumval_false_100() {
	int m = 100;
	value v = mkstringvalc(1, "s");

	for (int i = 0; i < m; i++) {
		isnumval(v);
	}

	freeval(v);
}

void benchmark_isnumval_100() {
	int m = 100;
	value v = mknumvali(1, 1);

	for (int i = 0; i < m; i++) {
		isnumval(v);
	}

	freeval(v);
}

void benchmark_isfalsetype_false_100() {
	int m = 100;
	value n = mknumvali(1, 1);

	for (int i = 0; i < m; i++) {
		isfalsetype(n);
	}

	freeval(n);
}

void benchmark_isfalsetype_100() {
	int m = 100;
	for (int i = 0; i < m; i++) {
		isfalsetype(false);
	}
}

void benchmark_isfalse_false_100() {
	int m = 100;
	value n = mknumvali(1, 1);

	for (int i = 0; i < m; i++) {
		isfalseval(n);
	}

	freeval(n);
}

void benchmark_isfalse_100() {
	int m = 100;
	for (int i = 0; i < m; i++) {
		isfalseval(false);
	}
}

void benchmark_init_free_pair_100() {
	int m = 100;

	value n1 = mknumvali(1, 1);
	value n2 = mknumvali(2, 1);

	for (int i = 0; i < m; i++) {
		value p = mkpairval(n1, n2);
		freeval(p);
	}

	freeval(n1);
	freeval(n2);
}

void benchmark_isnulltype_false_100() {
	int m = 100;
	value n = mknumvali(1, 1);

	for (int i = 0; i < m; i++) {
		isnulltype(n);
	}

	freeval(n);
}

void benchmark_isnulltype_100() {
	int m = 100;
	for (int i = 0; i < m; i++) {
		isnulltype(null);
	}
}

void benchmark_isnull_false_100() {
	int m = 100;
	value n = mknumvali(1, 1);

	for (int i = 0; i < m; i++) {
		isnullval(n);
	}

	freeval(n);
}

void benchmark_isnull_100() {
	int m = 100;
	for (int i = 0; i < m; i++) {
		isnullval(null);
	}
}

void benchmark_ispairtype_false_100() {
	int m = 100;
	value n = mknumvali(1, 1);

	for (int i = 0; i < m; i++) {
		ispairtype(n);
	}

	freeval(n);
}

void benchmark_ispairtype_100() {
	int m = 100;

	value n0 = mknumvali(1, 1);
	value n1 = mknumvali(2, 1);
	value p = mkpairval(n0, n1);

	for (int i = 0; i < m; i++) {
		ispairtype(p);
	}

	freeval(p);
	freeval(n0);
	freeval(n1);
}

void benchmark_ispair_false_100() {
	int m = 100;
	value n = mknumvali(1, 1);

	for (int i = 0; i < m; i++) {
		ispairval(n);
	}

	freeval(n);
}

void benchmark_ispair_100() {
	int m = 100;

	value n0 = mknumvali(1, 1);
	value n1 = mknumvali(2, 1);
	value p = mkpairval(n0, n1);

	for (int i = 0; i < m; i++) {
		ispairval(p);
	}

	freeval(p);
	freeval(n0);
	freeval(n1);
}

void benchmark_pair_car_100() {
	int m = 100;

	value n1 = mknumvali(1, 1);
	value n2 = mknumvali(2, 1);
	value p = mkpairval(n1, n2);

	for (int i = 0; i < m; i++) {
		carval(p);
	}

	freeval(p);
	freeval(n1);
	freeval(n2);
}

void benchmark_pair_cdr_100() {
	int m = 100;

	value n1 = mknumvali(1, 1);
	value n2 = mknumvali(2, 1);
	value p = mkpairval(n1, n2);

	for (int i = 0; i < m; i++) {
		cdrval(p);
	}

	freeval(p);
	freeval(n1);
	freeval(n2);
}

void benchmark_init_free_primitive_proc_100() {
	int m = 100;
	for (int i = 0; i < m; i++) {
		value p = mkprimitiveprocval(&testprimitive);
		freeval(p);
	}
}

void benchmark_init_free_compiled_proc_100() {
	int m = 100;

	value label = mknumvali(1, 1);
	environment env = mkenvironment(0);

	for (int i = 0; i < m; i++) {
		value p = mkcompiledprocval(label, env);
		freeval(p);
	}

	freeval(label);
	freenv(env);
}

void benchmark_rawint_fail_not_number_100() {
	int m = 100;
	value v = mkstringvalc(1, "s");

	for (int i = 0; i < m; i++) {
		valrawint(v);
	}

	freeval(v);
	clearerrors();
}

void benchmark_rawint_fail_not_integer_100() {
	int m = 100;
	value n = mknumvali(1, 2);

	for (int i = 0; i < m; i++) {
		valrawint(n);
	}

	freeval(n);
	clearerrors();
}

void benchmark_rawint_100() {
	int m = 100;
	value n = mknumvali(42, 1);

	for (int i = 0; i < m; i++) {
		valrawint(n);
	}

	freeval(n);
}

void benchmakr_numval_100() {
	int m = 100;
	value n = mknumvali(42, 1);

	for (int i = 0; i < m; i++) {
		numval(n);
	}

	freeval(n);
}

void benchmark_applyprimitive_not_proc_100() {
	int m = 100;
	value n = mknumvali(1, 1);

	for (int i = 0; i < m; i++) {
		valapplyprimitive(n, null);
	}

	freeval(n);
	clearerrors();
}

void benchmark_applyprimitive_not_primitive_100() {
	int m = 100;

	value label = mknumvali(1, 1);
	environment env = mkenvironment(0);
	value p = mkcompiledprocval(label, env);

	for (int i = 0; i < m; i++) {
		valapplyprimitive(p, null);
	}

	freeval(p);
	freeval(label);
	freenv(env);
	clearerrors();
}

void benchmark_applyprimitive_100() {
	int m = 100;
	value p = mkprimitiveprocval(testprimitive);

	for (int i = 0; i < m; i++) {
		valapplyprimitive(p, null);
	}

	freeval(p);
}

void benchmark_proclabel_not_proc_100() {
	int m = 100;
	value n = mknumvali(1, 1);

	for (int i = 0; i < m; i++) {
		valproclabel(n);
	}

	freeval(n);
	clearerrors();
}

void benchmark_proclabel_not_compiled_100() {
	int m = 100;
	value p = mkprimitiveprocval(testprimitive);

	for (int i = 0; i < m; i++) {
		valproclabel(p);
	}

	freeval(p);
	clearerrors();
}

void benchmark_proclabel_100() {
	int m = 100;

	value label = mknumvali(1, 1);
	environment env = mkenvironment(0);
	value p = mkcompiledprocval(label, env);

	for (int i = 0; i < m; i++) {
		valproclabel(p);
	}

	freeval(p);
	freeval(label);
	freenv(env);
}

void benchmark_procenv_not_proc_100() {
	int m = 100;
	value n = mknumvali(1, 1);

	for (int i = 0; i < m; i++) {
		valprocenv(n);
	}

	freeval(n);
	clearerrors();
}

void benchmark_procenv_not_compiled_100() {
	int m = 100;
	value p = mkprimitiveprocval(testprimitive);

	for (int i = 0; i < m; i++) {
		valprocenv(p);
	}

	freeval(p);
	clearerrors();
}

void benchmark_procenv_100() {
	int m = 100;

	value label = mknumvali(1, 1);
	environment env = mkenvironment(0);
	value p = mkcompiledprocval(label, env);

	for (int i = 0; i < m; i++) {
		valprocenv(p);
	}

	freeval(p);
	freeval(label);
	freenv(env);
}

void benchmark_isprimitiveproctype_false_100() {
	int m = 100;
	value v = mknumvali(1, 1);

	for (int i = 0; i < m; i++) {
		isprimitiveproctype(v);
	}

	freeval(v);
}

void benchmark_isprimitiveproctype_100() {
	int m = 100;
	value p = mkprimitiveprocval(testprimitive);

	for (int i = 0; i < m; i++) {
		isprimitiveproctype(p);
	}

	freeval(p);
}

void benchmark_isprimitiveproc_false_100() {
	int m = 100;
	value v = mknumvali(1, 1);

	for (int i = 0; i < m; i++) {
		isprimitiveprocval(v);
	}

	freeval(v);
}

void benchmark_isprimitiveproc_100() {
	int m = 100;
	value p = mkprimitiveprocval(testprimitive);

	for (int i = 0; i < m; i++) {
		isprimitiveprocval(p);
	}

	freeval(p);
}

void benchmark_symbolnameraw_not_symbol_100() {
	int m = 100;
	value v = mknumvali(1, 1);

	for (int i = 0; i < m; i++) {
		valsymbolnameraw(v);
	}

	freeval(v);
	clearerrors();
}

void benchmark_symbolnameraw_100() {
	int m = 100;

	char *sname = "some-symbol";
	int slen = strlen(sname);
	value s = mksymval(slen, sname);

	for (int i = 0; i < m; i++) {
		valsymbolnameraw(s);
	}

	freeval(s);
}

void benchmark_sprintraw_number_small_100() {
	int m = 100;
	value n = mknumvali(9, 4);

	for (int i = 0; i < m; i++) {
		char *s = sprintraw(n);
		free(s);
	}

	freeval(n);
}

void benchmark_sprintraw_number_big_100() {
	int m = 100;

	char *c = "4798476948769769857028703847560387562984769.4897694769382765794876385760398476039865";
	int len = strlen(c);
	value n = mknumvalc(len, c);

	for (int i = 0; i < m; i++) {
		char *s = sprintraw(n);
		free(s);
	}

	freeval(n);
}

// void benchmark_init_free_file_100() {
// 	int m = 100;
// 
// 	file f;
// 	ioerror err;
// 	openfile("test-file", ioread, &f, &err);
// 
// 	for (int i = 0; i < m; i++) {
// 		value fv = mkfileval(f);
// 		freeval(fv);
// 	}
// 
// 	closefile(f);
// }
// 
// void benchmark_fileval_fail_100() {
// 	int m = 100;
// 	value v = mknumvali(1, 1);
// 
// 	for (int i = 0; i < m; i++) {
// 		fileval(v);
// 	}
// 
// 	freeval(v);
// 	clearerrors();
// }
// 
// void benchmark_fileval_100() {
// 	int m = 100;
// 
// 	file f;
// 	ioerror err;
// 	openfile("test-file", ioread, &f, &err);
// 	value v = mkfileval(f);
// 
// 	for (int i = 0; i < m; i++) {
// 		fileval(v);
// 	}
// 
// 	freeval(v);
// 	closefile(f);
// }
// 
// void benchmark_isfiletype_false_100() {
// 	int m = 100;
// 	value v = mknumvali(1, 1);
// 
// 	for (int i = 0; i < m; i++) {
// 		isfiletype(v);
// 	}
// 
// 	freeval(v);
// }
// 
// void benchmark_isfiletype_true_100() {
// 	int m = 100;
// 
// 	file f;
// 	ioerror err;
// 	openfile("test-file", ioread, &f, &err);
// 	value v = mkfileval(f);
// 
// 	for (int i = 0; i < m; i++) {
// 		isfiletype(v);
// 	}
// 
// 	freeval(v);
// 	closefile(f);
// }
// 
// void benchmark_isfileval_false_100() {
// 	int m = 100;
// 	value v = mknumvali(1, 1);
// 
// 	for (int i = 0; i < m; i++) {
// 		isfileval(v);
// 	}
// 
// 	freeval(v);
// }
// 
// void benchmark_isfileval_true_100() {
// 	int m = 100;
// 
// 	file f;
// 	ioerror err;
// 	openfile("test-file", ioread, &f, &err);
// 	value v = mkfileval(f);
// 
// 	for (int i = 0; i < m; i++) {
// 		isfileval(v);
// 	}
// 
// 	freeval(v);
// 	closefile(f);
// }

void benchmark_isinttype_false_100() {
	int m = 100;
	value v = mkstringvalc(11, "some string");

	for (int i = 0; i < m; i++) {
		isinttype(v);
	}

	freeval(v);
}

void benchmark_isinttype_true_100() {
	int m = 100;
	value v = mknumvali(42, 1);

	for (int i = 0; i < m; i++) {
		isinttype(v);
	}

	freeval(v);
}

void benchmark_isintval_false_100() {
	int m = 100;
	value v = mkstringvalc(11, "some string");

	for (int i = 0; i < m; i++) {
		isintval(v);
	}

	freeval(v);
}

void benchmark_isintval_true_100() {
	int m = 100;
	value v = mknumvali(42, 1);

	for (int i = 0; i < m; i++) {
		isintval(v);
	}

	freeval(v);
}

void benchmark_issmallinttype_false_100() {
	int m = 100;
	value v = mkstringvalc(11, "some string");

	for (int i = 0; i < m; i++) {
		issmallinttype(v);
	}

	freeval(v);
}

void benchmark_issmallinttype_true_100() {
	int m = 100;
	value v = mknumvali(42, 1);

	for (int i = 0; i < m; i++) {
		issmallinttype(v);
	}

	freeval(v);
}

void benchmark_issmallintval_false_100() {
	int m = 100;
	value v = mkstringvalc(11, "some string");

	for (int i = 0; i < m; i++) {
		issmallintval(v);
	}

	freeval(v);
}

void benchmark_issmallintval_true_100() {
	int m = 100;
	value v = mknumvali(42, 1);

	for (int i = 0; i < m; i++) {
		issmallintval(v);
	}

	freeval(v);
}

void benchmark_valstring_error_100() {
	int m = 100;
	value v = mknumvali(42, 1);

	for (int i = 0; i < m; i++) {
		valstring(v);
	}

	freeval(v);
	clearerrors();
}

void benchmark_valstring_100() {
	int m = 100;
	value v = mkstringvalc(11, "some string");

	for (int i = 0; i < m; i++) {
		valstring(v);
	}

	freeval(v);
}

int main(int argc, char **argv) {
	initsys();
	initmodule_errormock();
	initmodule_number();
	initmodule_value();
	initmodule_sysio();

	int n = repeatcount(argc, argv);

	int err = 0;
	err += benchmark(n, &benchmark_init_free_symbol_100, "value: init free symbol, 100 times");
	err += benchmark(n, &benchmark_init_free_number_100, "value: init free number, 100 times");
	err += benchmark(n, &benchmark_init_free_number_char_100, "value: init free number char, 100 times");
	err += benchmark(n, &benchmark_init_free_number_int_100, "value: init free number int, 100 times");
	err += benchmark(n, &benchmark_init_free_string_100, "value: init free string, 100 times");
	err += benchmark(n, &benchmark_issymtype_false_100, "value: type not false, 100 times");
	err += benchmark(n, &benchmark_issymtype_100, "value: type false, 100 times");
	err += benchmark(n, &benchmark_issymval_false_100, "value: is not false, 100 times");
	err += benchmark(n, &benchmark_issymval_100, "value: is false, 100 times");
	err += benchmark(n, &benchmark_isnumtype_false_100, "value: type not false, 100 times");
	err += benchmark(n, &benchmark_isnumtype_100, "value: type false, 100 times");
	err += benchmark(n, &benchmark_isnumval_false_100, "value: is not false, 100 times");
	err += benchmark(n, &benchmark_isnumval_100, "value: is false, 100 times");
	err += benchmark(n, &benchmark_isfalsetype_false_100, "value: type not false, 100 times");
	err += benchmark(n, &benchmark_isfalsetype_100, "value: type false, 100 times");
	err += benchmark(n, &benchmark_isfalse_false_100, "value: is not false, 100 times");
	err += benchmark(n, &benchmark_isfalse_100, "value: is false, 100 times");
	err += benchmark(n, &benchmark_init_free_pair_100, "value: init free pair, 100 times");
	err += benchmark(n, &benchmark_isnulltype_false_100, "value: type not null, 100 times");
	err += benchmark(n, &benchmark_isnulltype_100, "value: type null, 100 times");
	err += benchmark(n, &benchmark_isnull_false_100, "value: is not null, 100 times");
	err += benchmark(n, &benchmark_isnull_100, "value: is null, 100 times");
	err += benchmark(n, &benchmark_ispairtype_false_100, "value: type not pair, 100 times");
	err += benchmark(n, &benchmark_ispairtype_100, "value: type pair, 100 times");
	err += benchmark(n, &benchmark_ispair_false_100, "value: is not pair, 100 times");
	err += benchmark(n, &benchmark_ispair_100, "value: is pair, 100 times");
	err += benchmark(n, &benchmark_pair_car_100, "value: pair car, 100 times");
	err += benchmark(n, &benchmark_pair_cdr_100, "value: pair cdr, 100 times");
	err += benchmark(n, &benchmark_init_free_primitive_proc_100, "value: init free primitive proc, 100 times");
	err += benchmark(n, &benchmark_init_free_compiled_proc_100, "value: init free compiled proc, 100 times");
	err += benchmark(n, &benchmark_rawint_fail_not_number_100, "value: raw int, not number, 100 times");
	err += benchmark(n, &benchmark_rawint_fail_not_integer_100, "value: raw int, not integer, 100 times");
	err += benchmark(n, &benchmark_rawint_100, "value: raw int, 100 times");
	err += benchmark(n, &benchmakr_numval_100, "value: numval, 100 times");
	err += benchmark(n, &benchmark_applyprimitive_not_proc_100, "value: apply primitive, not proc, 100 times");
	err += benchmark(n, &benchmark_applyprimitive_not_primitive_100, "value: apply primitive, not primitive, 100 times");
	err += benchmark(n, &benchmark_applyprimitive_100, "value: apply primitive, 100 times");
	err += benchmark(n, &benchmark_proclabel_not_proc_100, "value: proc label, not proc, 100 times");
	err += benchmark(n, &benchmark_proclabel_not_compiled_100, "value: proc label, not compiled, 100 times");
	err += benchmark(n, &benchmark_proclabel_100, "value: proc label, 100 times");
	err += benchmark(n, &benchmark_procenv_not_proc_100, "value: proc env, not proc, 100 times");
	err += benchmark(n, &benchmark_procenv_not_compiled_100, "value: proc env, not compiled, 100 times");
	err += benchmark(n, &benchmark_procenv_100, "value: proc env, 100 times");
	err += benchmark(n, &benchmark_isprimitiveproctype_false_100, "value: proc not primitive type, 100 times");
	err += benchmark(n, &benchmark_isprimitiveproctype_100, "value: proc primitive type, 100 times");
	err += benchmark(n, &benchmark_isprimitiveproc_false_100, "value: proc not primitive, 100 times");
	err += benchmark(n, &benchmark_isprimitiveproc_100, "value: proc primitive, 100 times");
	err += benchmark(n, &benchmark_symbolnameraw_not_symbol_100, "value: symbol name, not symbol, 100 times");
	err += benchmark(n, &benchmark_symbolnameraw_100, "value: symbol name, 100 times");
	err += benchmark(n, &benchmark_sprintraw_number_small_100, "value: sprintraw, number, small, 100 times");
	err += benchmark(n, &benchmark_sprintraw_number_big_100, "value: sprintraw, number, big, 100 times");
	err += benchmark(n, &benchmark_isinttype_false_100, "value: isinttype, false, 100 times");
	err += benchmark(n, &benchmark_isinttype_true_100, "value: isinttype, true, 100 times");
	err += benchmark(n, &benchmark_isintval_false_100, "value: isintval, false, 100 times");
	err += benchmark(n, &benchmark_isintval_true_100, "value: isintval, true, 100 times");
	err += benchmark(n, &benchmark_issmallinttype_false_100, "value: issmallinttype, false, 100 times");
	err += benchmark(n, &benchmark_issmallinttype_true_100, "value: issmallinttype, true, 100 times");
	err += benchmark(n, &benchmark_issmallintval_false_100, "value: issmallintval, false, 100 times");
	err += benchmark(n, &benchmark_issmallintval_true_100, "value: issmallintval, true, 100 times");
	err += benchmark(n, &benchmark_valstring_error_100, "value: valstring, error, 100 times");
	err += benchmark(n, &benchmark_valstring_100, "value: valstring, 100 times");

	freemodule_sysio();
	freemodule_value();
	freemodule_number();
	freemodule_errormock();
	return err;
}
