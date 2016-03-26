#include <string.h>
#include "sys.h"
#include "testing.h"
#include "error.h"
#include "error.mock.h"
#include "number.h"
#include "string.h"
#include "compound-types.h"
#include "sysio.h"
#include "io.h"
#include "value.h"
#include "primitives.h"
#include "environment.h"

void benchmark_sumnumbers_fail_not_pair_100() {
	int m = 100;
	value v = mkstringvalc(1, "s");

	for (int i = 0; i < m; i++) {
		sumval(v);
	}

	freeval(v);
	clearerrors();
}

void benchmark_sumnumbers_fail_not_number_100() {
	int m = 100;

	value v = mkstringvalc(1, "s");
	value p = mkpairval(v, null);

	for (int i = 0; i < m; i++) {
		sumval(p);
	}

	freeval(p);
	freeval(v);
	clearerrors();
}

void benchmark_sumnumbers_null_100() {
	int m = 100;
	for (int i = 0; i < m; i++) {
		value s = sumval(null);
		freeval(s);
	}
}

void benchmark_sumnumbers_100() {
	int m = 100;

	value n1 = mknumvali(1, 1);
	value n2 = mknumvali(2, 1);
	value n3 = mknumvali(3, 1);
	value p1 = mkpairval(n3, null);
	value p2 = mkpairval(n2, p1);
	value p3 = mkpairval(n1, p2);

	for (int i = 0; i < m; i++) {
		value s = sumval(p3);
		freeval(s);
	}

	freeval(p3);
	freeval(p2);
	freeval(p1);
	freeval(n1);
	freeval(n2);
	freeval(n3);
	clearerrors();
}

void benchmark_diff_error_100() {
	int m = 100;

	value s = mkstringvalc(11, "some string");
	value p = mkpairval(s, null);

	for (int i = 0; i < m; i++) {
		diffval(p);
	}

	freeval(p);
	freeval(s);
	clearerrors();
}

void benchmark_diff_zero_100() {
	int m = 100;
	for (int i = 0; i < m; i++) {
		value d = diffval(null);
		freeval(d);
	}
}

void benchmark_diff_one_100() {
	int m = 100;

	value n = mknumvali(42, 1);
	value p = mkpairval(n, null);

	for (int i = 0; i < m; i++) {
		value d = diffval(p);
		freeval(d);
	}

	freeval(p);
	freeval(n);
}

void benchmark_diff_all_100() {
	int m = 100;

	value n0 = mknumvali(1, 1);
	value n1 = mknumvali(2, 1);
	value n2 = mknumvali(3, 1);
	value p0 = mkpairval(n2, null);
	value p1 = mkpairval(n1, p0);
	value p2 = mkpairval(n0, p1);

	for (int i = 0; i < m; i++) {
		value d = diffval(p2);
		freeval(d);
	}

	freeval(p2);
	freeval(p1);
	freeval(p0);
	freeval(n0);
	freeval(n1);
	freeval(n2);
}

void benchmark_bitor_fail_not_number_100() {
	int m = 100;
	value s = mkstringvalc(11, "some string");

	for (int i = 0; i < m; i++) {
		bitorval(s);
	}

	freeval(s);
	clearerrors();
}

void benchmark_bitor_fail_not_integer_100() {
	int m = 100;
	value n = mknumvalc(4, "3.14");

	for (int i = 0; i < m; i++) {
		bitorval(n);
	}

	freeval(n);
	clearerrors();
}

void benchmark_bitor_zero_args_100() {
	int m = 100;
	for (int i = 0; i < m; i++) {
		value n = bitorval(null);
		freeval(n);
	}
}

void benchmark_bitor_one_arg_100() {
	int m = 100;

	value n = mknumvali(42, 1);
	value p = mkpairval(n, null);

	for (int i = 0; i < m; i++) {
		value ior = bitorval(p);
		freeval(ior);
	}

	freeval(n);
	freeval(p);
}

void benchmark_bitor_multiple_args_100() {
	int m = 100;

	value n1 = mknumvali(1, 1);
	value n2 = mknumvali(2, 1);
	value n3 = mknumvali(4, 1);
	value p1 = mkpairval(n3, null);
	value p2 = mkpairval(n2, p1);
	value p3 = mkpairval(n1, p2);

	for (int i = 0; i < m; i++) {
		value ior = bitorval(p3);
		freeval(ior);
	}

	freeval(p3);
	freeval(p2);
	freeval(p1);
	freeval(n1);
	freeval(n2);
	freeval(n3);
}

void benchmark_openfile_wrong_arg_100() {
	int m = 100;

	value v = mknumvali(1, 1);
	value p = mkpairval(v, null);

	for (int i = 0; i < m; i++) {
		openfileval(p);
	}

	freeval(p);
	freeval(v);
	clearerrors();
}

void benchmark_openfile_wrong_number_of_args_100() {
	int m = 100;

	value fn = mkstringvalc(11, "some string");
	value p = mkpairval(fn, null);

	for (int i = 0; i < m; i++) {
		openfileval(p);
	}

	freeval(p);
	freeval(fn);
	clearerrors();
}

void benchmark_openfile_100() {
	int m = 100;

	value fn = mkstringvalc(11, "some string");
	value mode = mknumvali(ioread, 1);
	value p0 = mkpairval(mode, null);
	value p1 = mkpairval(fn, p0);

	for (int i = 0; i < m; i++) {
		value f = openfileval(p1);
		closefile(rawval(f));
		freeval(f);
	}

	freeval(p1);
	freeval(p0);
	freeval(fn);
	freeval(mode);
}

void benchmark_closefile_100() {
	int m = 100;

	value fn = mkstringvalc(11, "some string");
	value mode = mknumvali(ioread, 1);
	value p0 = mkpairval(mode, null);
	value p1 = mkpairval(fn, p0);

	for (int i = 0; i < m; i++) {
		value f = openfileval(p1);
		value p = mkpairval(f, null);
		value ok = closefileval(p);
		freeval(p);
		freeval(ok);
	}

	freeval(p1);
	freeval(p0);
	freeval(fn);
	freeval(mode);
}

void benchmark_seekfile_100() {
	int m = 100;

	value fn = mkstringvalc(11, "some string");
	value mode = mknumvali(ioread, 1);
	value p0 = mkpairval(mode, null);
	value p1 = mkpairval(fn, p0);
	value f = openfileval(p1);

	value seekmode = mknumvali(iostart, 1);
	value pos = mknumvali(42, 1);
	value p2 = mkpairval(seekmode, null);
	value p3 = mkpairval(pos, p2);
	value p4 = mkpairval(f, p3);

	for (int i = 0; i < m; i++) {
		value ok = seekfileval(p4);
		freeval(ok);
	}

	value p5 = mkpairval(f, null);
	value ok = closefileval(p5);
	freeval(ok);

	freeval(p5);
	freeval(p1);
	freeval(p0);
	freeval(fn);
	freeval(mode);
	freeval(p4);
	freeval(p3);
	freeval(p2);
	freeval(pos);
	freeval(seekmode);
}

void benchmark_readfile_100() {
	int m = 100;

	value fn = mkstringvalc(11, "some string");
	value mode = mknumvali(ioread, 1);
	value p0 = mkpairval(mode, null);
	value p1 = mkpairval(fn, p0);
	value f = openfileval(p1);

	value len = mknumvali(42, 1);
	value p2 = mkpairval(len, null);
	value p3 = mkpairval(f, p2);

	for (int i = 0; i < m; i++) {
		value s = readfileval(p3);
		if (s != eofval) {
			freeval(s);
		}
	}

	value p = mkpairval(f, null);
	value ok = closefileval(p);

	freeval(ok);
	freeval(p1);
	freeval(p0);
	freeval(fn);
	freeval(mode);
	freeval(p3);
	freeval(p2);
	freeval(len);
	freeval(p);
}

void benchmark_writefile_100() {
	int m = 100;

	value fn = mkstringvalc(11, "some string");
	value mode = mknumvali(iowrite, 1);
	value p0 = mkpairval(mode, null);
	value p1 = mkpairval(fn, p0);
	value f = openfileval(p1);

	value s = mkstringvalc(17, "some other string");
	value p2 = mkpairval(s, null);
	value p3 = mkpairval(f, p2);

	for (int i = 0; i < m; i++) {
		value okw = writefileval(p3);
		freeval(okw);
	}

	value p = mkpairval(f, null);
	value ok = closefileval(p);

	freeval(ok);
	freeval(p1);
	freeval(p0);
	freeval(fn);
	freeval(mode);
	freeval(p3);
	freeval(p2);
	freeval(s);
	freeval(p);
}

void benchmark_init_free_regex_100() {
	int m = 100;

	char *exps = "some expression";
	int len = strlen(exps);
	value exp = mkstringvalc(len, exps);
	value flags = mknumvali(0, 1);
	value p0 = mkpairval(flags, null);
	value p1 = mkpairval(exp, p0);

	for (int i = 0; i < m; i++) {
		value rx = mkregexval(p1);
		freeval(rx);
	}

	freeval(p1);
	freeval(p0);
	freeval(flags);
	freeval(exp);
}

void benchmark_regex_match_100() {
	int m = 100;

	char *exps = "\\([^(]*\\)";
	int len = strlen(exps);
	value exp = mkstringvalc(len, exps);
	value flags = mknumvali(0, 1);
	value p0 = mkpairval(flags, null);
	value p1 = mkpairval(exp, p0);
	value rx = mkregexval(p1);

	char *s = "((some list) (of lists))";
	int slen = strlen(s);
	value sval = mkstringvalc(slen, s);
	value p2 = mkpairval(sval, null);
	value p3 = mkpairval(rx, p2);

	for (int i = 0; i < m; i++) {
		// value m = regexmatch(p3);
		// freeval(carval(carval(m)));
		// freeval(carval(cdrval(carval(m))));
		// freeval(carval(m));
		// freeval(m);
	}

	freeval(p3);
	freeval(p2);
	freeval(p1);
	freeval(p0);
	freeval(flags);
	freeval(exp);
	freeval(sval);
	freeval(rx);
}

void benchmark_regex_nomatch_100() {
	int m = 100;

	value exp = mkstringvalc(1, "a");
	value flags = mknumvali(0, 1);
	value p0 = mkpairval(flags, null);
	value p1 = mkpairval(exp, p0);
	value rx = mkregexval(p1);
	value s = mkstringvalc(1, "b");
	value p2 = mkpairval(s, null);
	value p3 = mkpairval(rx, p2);

	for (int i = 0; i < m; i++) {
		regexmatch(p3);
	}

	freeval(p3);
	freeval(p2);
	freeval(s);
	freeval(rx);
	freeval(p1);
	freeval(p0);
	freeval(flags);
	freeval(exp);
}

void benchmark_error_100() {
	int m = 100;

	char *msg = "some message";
	int len = strlen(msg);
	value msgval = mkstringvalc(len, msg);
	value args = mkpairval(msgval, null);

	for (int i = 0; i < m; i++) {
		value r = errorval(args);
		freeval(r);
	}

	freeval(args);
	freeval(msgval);
	clearerrors();
}

void benchmark_eq_symbol_false_100() {
	int m = 100;

	value s1 = mksymval(3, "abc");
	value s2 = mksymval(3, "Abc");
	value args = mkpairval(s2, null);
	args = mkpairval(s1, args);

	for (int i = 0; i < m; i++) {
		iseqval(args);
	}

	freeval(cdrval(args));
	freeval(args);
	freeval(s1);
	freeval(s2);
}

void benchmark_eq_symbol_true_100() {
	int m = 100;

	value s1 = mksymval(3, "abc");
	value s2 = mksymval(3, "abc");
	value args = mkpairval(s2, null);
	args = mkpairval(s1, args);

	for (int i = 0; i < m; i++) {
		iseqval(args);
	}

	freeval(cdrval(args));
	freeval(args);
	freeval(s1);
	freeval(s2);
}

void benchmark_eq_number_false_100() {
	int m = 100;

	value n1 = mknumvali(42, 1);
	value n2 = mknumvali(3, 1);
	value args = mkpairval(n2, null);
	args = mkpairval(n1, args);

	for (int i = 0; i < m; i++) {
		iseqval(args);
	}

	freeval(cdrval(args));
	freeval(args);
	freeval(n1);
	freeval(n2);
}

void benchmark_eq_number_true_100() {
	int m = 100;

	value n1 = mknumvali(3, 1);
	value n2 = mknumvali(3, 1);
	value args = mkpairval(n2, null);
	args = mkpairval(n1, args);

	for (int i = 0; i < m; i++) {
		iseqval(args);
	}

	freeval(cdrval(args));
	freeval(args);
	freeval(n1);
	freeval(n2);
}

void benchmark_eq_string_false_100() {
	int m = 100;

	value s1 = mkstringvalc(3, "abc");
	value s2 = mkstringvalc(3, "Abc");
	value args = mkpairval(s2, null);
	args = mkpairval(s1, args);

	for (int i = 0; i < m; i++) {
		iseqval(args);
	}

	freeval(cdrval(args));
	freeval(args);
	freeval(s1);
	freeval(s2);
}

void benchmark_eq_string_true_100() {
	int m = 100;

	value s1 = mkstringvalc(3, "abc");
	value s2 = mkstringvalc(3, "abc");
	value args = mkpairval(s2, null);
	args = mkpairval(s1, args);

	for (int i = 0; i < m; i++) {
		iseqval(args);
	}

	freeval(cdrval(args));
	freeval(args);
	freeval(s1);
	freeval(s2);
}

void benchmark_eq_refs_false_100() {
	int m = 100;

	value entry = mknumvali(1, 1);
	environment env = mkenvironment(0);
	value v1 = mkcompiledprocval(entry, env);
	value v2 = mkcompiledprocval(entry, env);
	value args = mkpairval(v2, null);
	args = mkpairval(v1, args);
	args = mkpairval(v1, args);

	for (int i = 0; i < m; i++) {
		iseqval(args);
	}

	freeval(cdrval(cdrval(args)));
	freeval(cdrval(args));
	freeval(args);
	freeval(entry);
	freenv(env);
	freeval(v1);
	freeval(v2);
}

void benchmark_eq_refs_true_100() {
	int m = 100;

	value entry = mknumvali(1, 1);
	environment env = mkenvironment(0);
	value v1 = mkcompiledprocval(entry, env);
	value args = mkpairval(v1, null);
	args = mkpairval(v1, args);
	args = mkpairval(v1, args);

	for (int i = 0; i < m; i++) {
		iseqval(args);
	}

	freeval(cdrval(cdrval(args)));
	freeval(cdrval(args));
	freeval(args);
	freeval(entry);
	freenv(env);
	freeval(v1);
}

void benchmark_eq_refs_one_100() {
	int m = 100;

	value v = mknumvali(1, 1);
	value args = mkpairval(v, null);

	for (int i = 0; i < m; i++) {
		iseqval(args);
	}

	freeval(args);
	freeval(v);
}

void benchmark_eq_refs_zero_100() {
	int m = 100;
	for (int i = 0; i < m; i++) {
		iseqval(null);
	}
}

void benchmark_isutf8_false_100() {
	int m = 100;
	char *raw = (char []){255, 1, 0};
	value s = mkstringvalc(2, raw);
	value p = mkpairval(s, null);

	for (int i = 0; i < m; i++) {
		isutf8val(p);
	}

	freeval(p);
	freeval(s);
}

void benchmark_isutf8_true_100() {
	int m = 100;
	value s = mkstringvalc(11, "some string");
	value p = mkpairval(s, null);

	for (int i = 0; i < m; i++) {
		isutf8val(p);
	}

	freeval(p);
	freeval(s);
}

void benchmark_copystrval_error_100() {
	int m = 100;

	value s = mknumvali(42, 1);
	value from = mknumvali(3, 1);
	value len = mknumvali(3, 1);
	value p0 = mkpairval(len, null);
	value p1 = mkpairval(from, p0);
	value p2 = mkpairval(s, p1);

	for (int i = 0; i < m; i++) {
		copystrval(p2);
	}

	freeval(p0);
	freeval(p1);
	freeval(p2);
	freeval(len);
	freeval(from);
	freeval(s);
	clearerrors();
}

void benchmark_copystrval_100() {
	int m = 100;

	value s = mkstringvalc(11, "some string");
	value from = mknumvali(3, 1);
	value len = mknumvali(3, 1);
	value p0 = mkpairval(len, null);
	value p1 = mkpairval(from, p0);
	value p2 = mkpairval(s, p1);

	for (int i = 0; i < m; i++) {
		value ss = copystrval(p2);
		freeval(ss);
	}

	freeval(p0);
	freeval(p1);
	freeval(p2);
	freeval(len);
	freeval(from);
	freeval(s);
}

void benchmark_byteslenval_error_100() {
	int m = 100;

	value v = mknumvali(42, 1);
	value p = mkpairval(v, null);

	for (int i = 0; i < m; i++) {
		byteslenval(p);
	}

	freeval(p);
	freeval(v);
	clearerrors();
}

void benchmark_byteslenval_100() {
	int m = 100;

	value s = mkstringvalc(11, "some string");
	value p = mkpairval(s, null);

	for (int i = 0; i < m; i++) {
		value l = byteslenval(p);
		freeval(l);
	}

	freeval(p);
	freeval(s);
}

void benchmark_stringappend_error_100() {
	int m = 100;

	value v = mknumvali(42, 1);
	value p = mkpairval(v, null);

	for (int i = 0; i < m; i++) {
		stringappendval(p);
	}

	freeval(p);
	freeval(v);
	clearerrors();
}

void benchmark_stringappend_100() {
	int m = 100;

	value s1 = mkstringvalc(5, "some ");
	value s2 = mkstringvalc(6, "string");
	value p0 = mkpairval(s1, null);
	value p1 = mkpairval(s2, p0);

	for (int i = 0; i < m; i++) {
		value sa = stringappendval(p1);
		freeval(sa);
	}

	freeval(p1);
	freeval(p0);
	freeval(s1);
	freeval(s2);
}

void benchmark_stringlenval_error_100() {
	int m = 100;

	value v = mknumvali(42, 1);
	value p = mkpairval(v, null);

	for (int i = 0; i < m; i++) {
		stringlenval(p);
	}

	freeval(p);
	freeval(v);
	clearerrors();
}

void benchmark_stringlenval_100() {
	int m = 100;

	value s = mkstringvalc(11, "some string");
	value p = mkpairval(s, null);

	for (int i = 0; i < m; i++) {
		value l = stringlenval(p);
		freeval(l);
	}

	freeval(p);
	freeval(s);
}

void benchmark_iseofval_false_100() {
	int m = 100;

	value v = mknumvali(42, 1);
	value p = mkpairval(v, null);

	for (int i = 0; i < m; i++) {
		iseofval(p);
	}

	freeval(v);
	freeval(p);
}

void benchmark_iseofval_true_100() {
	int m = 100;
	value p = mkpairval(eofval, null);

	for (int i = 0; i < m; i++) {
		iseofval(p);
	}

	freeval(p);
}

void benchmark_lessval_error_100() {
	int m = 100;
	value p = mkpairval(null, null);

	for (int i = 0; i < m; i++) {
		islessval(p);
	}

	freeval(p);
	clearerrors();
}

void benchmark_lessval_false_100() {
	int m = 100;

	value n1 = mknumvali(1, 1);
	value n2 = mknumvali(2, 1);
	value n3 = mknumvali(0, 1);
	value p0 = mkpairval(n3, null);
	value p1 = mkpairval(n2, p0);
	value p2 = mkpairval(n1, p1);

	for (int i = 0; i < m; i++) {
		islessval(p2);
	}

	freeval(p2);
	freeval(p1);
	freeval(p0);
	freeval(n1);
	freeval(n2);
	freeval(n3);
}

void benchmark_lessval_true_100() {
	int m = 100;

	value n1 = mknumvali(1, 1);
	value n2 = mknumvali(2, 1);
	value n3 = mknumvali(3, 1);
	value p0 = mkpairval(n3, null);
	value p1 = mkpairval(n2, p0);
	value p2 = mkpairval(n1, p1);

	for (int i = 0; i < m; i++) {
		islessval(p2);
	}

	freeval(p2);
	freeval(p1);
	freeval(p0);
	freeval(n1);
	freeval(n2);
	freeval(n3);
}

void benchmark_lessval_string_100() {
	int m = 100;

	value s1 = mkstringvalc(1, "a");
	value s2 = mkstringvalc(1, "b");
	value p0 = mkpairval(s2, null);
	value p1 = mkpairval(s1, p0);

	for (int i = 0; i < m; i++) {
		islessval(p1);
	}

	freeval(p1);
	freeval(p0);
	freeval(s1);
	freeval(s2);
}

void benchmark_notval_false_100() {
	int m = 100;

	value v = mknumvali(42, 1);
	value p = mkpairval(v, null);

	for (int i = 0; i < m; i++) {
		notval(p);
	}

	freeval(v);
	freeval(p);
}

void benchmark_notval_true_100() {
	int m = 100;
	value p = mkpairval(false, null);

	for (int i = 0; i < m; i++) {
		notval(p);
	}

	freeval(p);
}

void benchmark_isnullvalp_false_100() {
	int m = 100;

	value v = mknumvali(42, 1);
	value p = mkpairval(v, null);

	for (int i = 0; i < m; i++) {
		isnullvalp(p);
	}

	freeval(p);
	freeval(v);
}

void benchmark_isnullvalp_true_100() {
	int m = 100;
	value p = mkpairval(null, null);

	for (int i = 0; i < m; i++) {
		isnullvalp(p);
	}

	freeval(p);
}

void benchmark_consval_error_100() {
	int m = 100;

	value a = mknumvali(42, 1);
	value p = mkpairval(a, null);

	for (int i = 0; i < m; i++) {
		consval(p);
	}

	freeval(p);
	freeval(a);
	clearerrors();
}

void benchmark_consval_100() {
	int m = 100;

	value a = mknumvali(42, 1);
	value b = mknumvali(36, 1);
	value p0 = mkpairval(b, null);
	value p1 = mkpairval(a, p0);

	for (int i = 0; i < m; i++) {
		value p = consval(p1);
		freeval(p);
	}

	freeval(p1);
	freeval(p0);
	freeval(b);
	freeval(a);
}

void benchmark_carvalp_error_100() {
	int m = 100;

	value v = mknumvali(42, 1);
	value args = mkpairval(v, null);

	for (int i = 0; i < m; i++) {
		carvalp(args);
	}

	freeval(args);
	freeval(v);
	clearerrors();
}

void benchmark_carvalp_100() {
	int m = 100;

	value v1 = mknumvali(42, 1);
	value v2 = mknumvali(84, 1);
	value p = mkpairval(v1, v2);
	value args = mkpairval(p, null);

	for (int i = 0; i < m; i++) {
		carvalp(args);
	}

	freeval(args);
	freeval(p);
	freeval(v1);
	freeval(v2);
}

void benchmark_cdrvalp_error_100() {
	int m = 100;

	value v = mknumvali(42, 1);
	value args = mkpairval(v, null);

	for (int i = 0; i < m; i++) {
		cdrvalp(args);
	}

	freeval(args);
	freeval(v);
	clearerrors();
}

void benchmark_cdrvalp_100() {
	int m = 100;

	value v1 = mknumvali(42, 1);
	value v2 = mknumvali(84, 1);
	value p = mkpairval(v1, v2);
	value args = mkpairval(p, null);

	for (int i = 0; i < m; i++) {
		cdrvalp(args);
	}

	freeval(args);
	freeval(p);
	freeval(v1);
	freeval(v2);
}

void benchmark_stringtonumsafe_fail_100() {
	int m = 100;

	value s = mkstringvalc(1, "a");
	value p = mkpairval(s, null);

	for (int i = 0; i < m; i++) {
		stringtonumsafe(p);
	}

	freeval(p);
	freeval(s);
}

void benchmark_stringtonumsafe_100() {
	int m = 100;

	value s = mkstringvalc(1, "1");
	value p = mkpairval(s, null);

	for (int i = 0; i < m; i++) {
		value n = stringtonumsafe(p);
		freeval(n);
	}

	freeval(p);
	freeval(s);
}

int main(int argc, char **argv) {
	initsys();
	initmodule_errormock();
	initmodule_number();
	initmodule_value();
	initmodule_sysio();
	initmodule_primitives();
	int n = repeatcount(argc, argv);

	int err = 0;
	// err += benchmark(n, &benchmark_sumnumbers_fail_not_pair_100, "primitives: sumval, not pair, 100 times");
	// err += benchmark(n, &benchmark_sumnumbers_fail_not_number_100, "primitives: sumval, not number, 100 times");
	// err += benchmark(n, &benchmark_sumnumbers_null_100, "primitives: sumval, null, 100 times");
	err += benchmark(n, &benchmark_sumnumbers_100, "primitives: sumval, 100 times");
	// err += benchmark(n, &benchmark_diff_error_100, "primitives: diffval, error, 100 times");
	// err += benchmark(n, &benchmark_diff_zero_100, "primitives: diffval, zero, 100 times");
	// err += benchmark(n, &benchmark_diff_one_100, "primitives: diffval, one, 100 times");
	// err += benchmark(n, &benchmark_diff_all_100, "primitives: diffval, all, 100 times");
	// err += benchmark(n, &benchmark_bitor_fail_not_number_100, "primitives: bitor, fail, not a number, 100 times");
	// err += benchmark(n, &benchmark_bitor_fail_not_integer_100, "primitives: bitor, fail not integer, 100 times");
	// err += benchmark(n, &benchmark_bitor_zero_args_100, "primitives: bitor, zero args, 100 times");
	// err += benchmark(n, &benchmark_bitor_one_arg_100, "primitives: bitor, one arg, 100 times");
	// err += benchmark(n, &benchmark_bitor_multiple_args_100, "primitives: bitor, multiple args, 100 times");
	// err += benchmark(n, &benchmark_openfile_wrong_arg_100, "primitives: open file, wrong arg, 100 times");
	// err += benchmark(n, &benchmark_openfile_wrong_number_of_args_100, "primitives: open file, wrong number of arg, 100 times");
	// err += benchmark(n, &benchmark_openfile_100, "primitives: open file, 100 times");
	// err += benchmark(n, &benchmark_closefile_100, "primitives: close file, 100 times");
	// err += benchmark(n, &benchmark_seekfile_100, "primitives: seek file, 100 times");
	// err += benchmark(n, &benchmark_readfile_100, "primitives: read file, 100 times");
	// err += benchmark(n, &benchmark_writefile_100, "primitives: writefile, 100 times");
	// err += benchmark(n, &benchmark_init_free_regex_100, "primitives: init free regex, 100 times");
	// err += benchmark(n, &benchmark_regex_match_100, "primitives: match regex, 100 times");
	// err += benchmark(n, &benchmark_regex_nomatch_100, "primitives: regex, nomatch, 100 times");
	// err += benchmark(n, &benchmark_error_100, "primitives: error, 100 times");
	// err += benchmark(n, &benchmark_eq_symbol_false_100, "primitives: eq, symbol, false, 100 times");
	// err += benchmark(n, &benchmark_eq_symbol_true_100, "primitives: eq, symbol, true, 100 times");
	// err += benchmark(n, &benchmark_eq_number_false_100, "primitives: eq, number, false, 100 times");
	// err += benchmark(n, &benchmark_eq_number_true_100, "primitives: eq, number, true, 100 times");
	// err += benchmark(n, &benchmark_eq_string_false_100, "primitives: eq, string, false, 100 times");
	// err += benchmark(n, &benchmark_eq_string_true_100, "primitives: eq, string, true, 100 times");
	// err += benchmark(n, &benchmark_eq_refs_false_100, "primitives: eq, refs, false, 100 times");
	// err += benchmark(n, &benchmark_eq_refs_true_100, "primitives: eq, refs, true, 100 times");
	// err += benchmark(n, &benchmark_eq_refs_one_100, "primitives: eq, one, 100 times");
	// err += benchmark(n, &benchmark_eq_refs_zero_100, "primitives: eq, zero, 100 times");
	// err += benchmark(n, &benchmark_isutf8_false_100, "primitives: isutf8, false, 100 times");
	// err += benchmark(n, &benchmark_isutf8_true_100, "primitives: isutf8, true, 100 times");
	// err += benchmark(n, &benchmark_copystrval_error_100, "primitives: substr, error, 100 times");
	// err += benchmark(n, &benchmark_copystrval_100, "primitives: substr, 100 times");
	// err += benchmark(n, &benchmark_byteslenval_error_100, "primitives: byteslenval, error, 100 times");
	// err += benchmark(n, &benchmark_byteslenval_100, "primitives: byteslenval, 100 times");
	// err += benchmark(n, &benchmark_stringappend_error_100, "primitives: stringappend, error, 100 times");
	// err += benchmark(n, &benchmark_stringappend_100, "primitives: stringappend, 100 times");
	// err += benchmark(n, &benchmark_stringlenval_error_100, "primitives: stringlenval, error, 100 times");
	// err += benchmark(n, &benchmark_stringlenval_error_100, "primitives: stringlenval, 100 times");
	// err += benchmark(n, &benchmark_iseofval_false_100, "primitives: iseofval, false, 100 times");
	// err += benchmark(n, &benchmark_iseofval_true_100, "primitives: iseofval, true, 100 times");
	// err += benchmark(n, &benchmark_lessval_error_100, "primitives: lessval, error, 100 times");
	// err += benchmark(n, &benchmark_lessval_false_100, "primitives: lessval, false, 100 times");
	// err += benchmark(n, &benchmark_lessval_true_100, "primitives: lessval, true, 100 times");
	// err += benchmark(n, &benchmark_lessval_string_100, "primitives: lessval, string, 100 times");
	// err += benchmark(n, &benchmark_notval_false_100, "primitives: notval, false, 100 times");
	// err += benchmark(n, &benchmark_notval_true_100, "primitives: notval, true, 100 times");
	// err += benchmark(n, &benchmark_regex_nomatch_100, "primitives: regex, nomatch, 100 times");
	// err += benchmark(n, &benchmark_isnullvalp_false_100, "primitives: isnullvalp, false, 100 times");
	// err += benchmark(n, &benchmark_isnullvalp_true_100, "primitives: isnullvalp, true, 100 times");
	// err += benchmark(n, &benchmark_consval_error_100, "primitives: consval, error, 100 times");
	// err += benchmark(n, &benchmark_consval_100, "primitives: consval, 100 times");
	// err += benchmark(n, &benchmark_carvalp_error_100, "primitives: carvalp, error, 100 times");
	// err += benchmark(n, &benchmark_carvalp_100, "primitives: carvalp, error, 100 times");
	// err += benchmark(n, &benchmark_cdrvalp_error_100, "primitives: carvalp, error, 100 times");
	// err += benchmark(n, &benchmark_cdrvalp_100, "primitives: carvalp, error, 100 times");
	// err += benchmark(n, &benchmark_stringtonumsafe_fail_100, "primitives: stringtonumsafe, fail, 100 times");
	// err += benchmark(n, &benchmark_stringtonumsafe_100, "primitives: stringtonumsafe, 100 times");

	freemodule_primitives();
	freemodule_sysio();
	freemodule_value();
	freemodule_number();
	freemodule_errormock();
	return err;
}
