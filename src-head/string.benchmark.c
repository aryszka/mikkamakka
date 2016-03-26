#include <stdlib.h>
#include <string.h>
#include "sys.h"
#include "testing.h"
#include "string.h"

char *teststring = "fűzfánfütyülő";

void benchmark_init_free_100() {
	int m = 100;
	for (int i = 0; i < m; i++) {
		string s = mkstring(3, "123");
		freestring(s);
	}
}

void benchmark_bytes_len_100() {
	int m = 100;
	string s = mkstring(3, "123");

	for (int i = 0; i < m; i++) {
		byteslen(s);
	}

	freestring(s);
}

void benchmark_string_length_damaged_100() {
	int m = 100;

	char *s = "fűzfánfütyülő";
	int len = strlen(s);
	int dpos = 3;
	char *sd = malloc(len + 1);
	for (int i = 0; i < len; i++) {
		*(sd + i) = i == dpos ? -1 : *(s + i);
	}
	*(sd + len) = 0;

	string st = mkstring(len, sd);

	for (int i = 0; i < m; i++) {
		stringlen(st);
	}

	free(sd);
	freestring(st);
}

void benchmark_string_length_100() {
	int m = 100;
	string s = mkstring(3, "123");

	for (int i = 0; i < m; i++) {
		stringlen(s);
	}

	freestring(s);
}

void benchmark_raw_100() {
	int m = 100;
	string s = mkstring(3, "123");

	for (int i = 0; i < m; i++) {
		rawstring(s);
	}

	freestring(s);
}

void benchmark_isutf8_100() {
	int m = 100;
	string s = mkstring(3, "123");

	for (int i = 0; i < m; i++) {
		isutf8(s);
	}

	freestring(s);
}

void benchmark_sprint_100() {
	int m = 100;

	char *raw = "some string";
	int len = strlen(raw);
	string s = mkstring(len, raw);

	for (int i = 0; i < m; i++) {
		char *ss = sprintstring(s);
		free(ss);
	}

	freestring(s);
}

void benchmark_comparestring_less_100() {
	int m = 100;

	string a = mkstring(1, "a");
	string b = mkstring(1, "b");

	for (int i = 0; i < m; i++) {
		comparestring(a, b);
	}

	freestring(a);
	freestring(b);
}

void benchmark_comparestring_eq_100() {
	int m = 100;

	string a = mkstring(1, "a");
	string aa = mkstring(1, "a");

	for (int i = 0; i < m; i++) {
		comparestring(a, aa);
	}

	freestring(a);
	freestring(aa);
}

void benchmark_comparestring_greater_100() {
	int m = 100;

	string a = mkstring(1, "a");
	string b = mkstring(1, "b");

	for (int i = 0; i < m; i++) {
		comparestring(b, a);
	}

	freestring(a);
	freestring(b);
}

void benchmark_substring_lower_boundary_100() {
	int m = 100;

	char *s = "some string";
	long len = strlen(s);
	string st = mkstring(len, s);

	for (int i = 0; i < m; i++) {
		string sub = substr(st, 0, -1);
		freestring(sub);
	}

	freestring(st);
}

void benchmark_substring_higher_boundary_100() {
	int m = 100;

	char *s = "some string";
	long len = strlen(s);
	string st = mkstring(len, s);

	for (int i = 0; i < m; i++) {
		string sub = substr(st, len, 1);
		freestring(sub);
	}

	freestring(st);
}

void benchmark_substring_100() {
	int m = 100;

	char *s = teststring;
	long len = strlen(s);
	string st = mkstring(len, s);

	for (int i = 0; i < m; i++) {
		string sub = substr(st, 1, 2);
		freestring(sub);
	}

	freestring(st);
}

void benchmark_append_zero_100() {
	int m = 100;
	for (int i = 0; i < m; i++) {
		string s = appendstr(0, 0);
		freestring(s);
	}
}

void benchmark_append_one_100() {
	int m = 100;

	char *raw = "some string";
	long len = strlen(raw);
	string s = mkstring(len, raw);

	for (int i = 0; i < m; i++) {
		string a = appendstr(1, &s);
		freestring(a);
	}

	freestring(s);
}

void benchmark_append_multiple_100() {
	int m = 100;

	char *raw1 = "some ";
	long len1 = strlen(raw1);
	string s1 = mkstring(len1, raw1);

	char *raw2 = "strings ";
	long len2 = strlen(raw2);
	string s2 = mkstring(len2, raw2);

	char *raw3 = "to append";
	long len3 = strlen(raw3);
	string s3 = mkstring(len3, raw3);

	string *ss = malloc(3 * sizeof(string));
	*ss = s1;
	*(ss + 1) = s2;
	*(ss + 2) = s3;

	for (int i = 0; i < m; i++) {
		string a = appendstr(3, ss);
		freestring(a);
	}

	free(ss);
	freestring(s1);
	freestring(s2);
	freestring(s3);
}

void benchmark_clonestring_100() {
	int m = 100;

	char *raw = "some string";
	long len = strlen(raw);
	string s = mkstring(len, raw);

	for (int i = 0; i < m; i++) {
		string cs = clonestring(s);
		freestring(cs);
	}

	freestring(s);
}

void benchmark_bytestochars_100() {
	int m = 100;
	for (int i = 0; i < m; i++) {
		bytestochars(strlen(teststring) - 1, teststring + 1);
	}
}

int main(int argc, char **argv) {
	initsys();
	int n = repeatcount(argc, argv);

	int err = 0;
	err += benchmark(n, &benchmark_init_free_100, "string: init, free string, 100 times");
	err += benchmark(n, &benchmark_bytes_len_100, "string: bytes length, 100 times");
	err += benchmark(n, &benchmark_string_length_damaged_100, "string: string length, damaged, 100 times");
	err += benchmark(n, &benchmark_string_length_100, "string: string length, 100 times");
	err += benchmark(n, &benchmark_raw_100, "string: raw, 100 times");
	err += benchmark(n, &benchmark_isutf8_100, "string: is utf8, 100 times");
	err += benchmark(n, &benchmark_sprint_100, "string: sprint, 100 times");
	err += benchmark(n, &benchmark_comparestring_less_100, "string: compare, less, 100 times");
	err += benchmark(n, &benchmark_comparestring_eq_100, "string: compare, eq, 100 times");
	err += benchmark(n, &benchmark_comparestring_greater_100, "string: compare, greater, 100 times");
	err += benchmark(n, &benchmark_substring_lower_boundary_100, "string: substring, lower bound, 100 times");
	err += benchmark(n, &benchmark_substring_higher_boundary_100, "string: substring, higher bound, 100 times");
	err += benchmark(n, &benchmark_substring_100, "string: substring, 100 times");
	err += benchmark(n, &benchmark_append_zero_100, "string: append, zero items, 100 times");
	err += benchmark(n, &benchmark_append_one_100, "string: append, one item, 100 times");
	err += benchmark(n, &benchmark_append_multiple_100, "string: append, multiple items, 100 times");
	err += benchmark(n, &benchmark_clonestring_100, "string: clone, 100 times");
	err += benchmark(n, &benchmark_bytestochars_100, "string: bytestochars, 100 times");

	return err;
}
