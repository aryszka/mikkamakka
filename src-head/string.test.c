#include <stdlib.h>
#include <string.h>
#include "sys.h"
#include "testing.h"
#include "string.h"

char *teststring = "fűzfánfütyülő";
int teststringlen = 13;

void test_init_free() {
	string s = mkstring(3, "123");
	freestring(s);
}

void test_bytes_len() {
	string s = mkstring(3, "123");
	assert(byteslen(s) == 3, "string len");
	freestring(s);
}

void test_bytes_length_from_arg() {
	string s = mkstring(3, "123123");
	assert(byteslen(s) == 3, "string len");
	freestring(s);
}

void test_string_length() {
	string s = mkstring(strlen(teststring), teststring);
	assert(stringlen(s) == teststringlen, "string length");
	freestring(s);
}

void test_string_length_damaged() {
	char *s = teststring;
	int len = strlen(s);
	int dpos = 3;
	char *sd = malloc(len + 1);
	for (int i = 0; i < len; i++) {
		*(sd + i) = i == dpos ? -1 : *(s + i);
	}
	*(sd + len) = 0;

	string st = mkstring(len, sd);
	assert(byteslen(st) == len, "string: damaged, bytes length");
	assert(stringlen(st) == 2, "string: damaged, string length");

	free(sd);
	freestring(st);
}

void test_raw() {
	string s = mkstring(3, "123");
	assert(!strcoll(rawstring(s), "123"), "string raw");
	freestring(s);
}

void test_raw_len_from_arg() {
	string s = mkstring(3, "123123");
	assert(!strcoll(rawstring(s), "123"), "string raw from arg");
	freestring(s);
}

void test_returns_null_terminated() {
	string s = mkstring(3, (char[]){'1', '2', '3'});
	assert(byteslen(s) == 3, "string len from non-null terminated");
	char *sout = rawstring(s);
	assert(strlen(sout) == 3, "string strlen from non-null terminated");
	assert(!strcoll(sout, "123"), "string compare from non-null terminated");
	freestring(s);
}

void test_isutf8_positive() {
	int len = strlen(teststring);
	string s = mkstring(len, teststring);
	assert(isutf8(s), "check is unicode, positive");
	freestring(s);
}

void test_isutf8_negative() {
	char *raw = (char[]){-1, 0};
	int len = strlen(raw);
	string s = mkstring(len, raw);
	assert(!isutf8(s), "check is unicode, negative");
	freestring(s);
}

void test_sprint() {
	char *raw = "some string";
	int len = strlen(raw);
	string s = mkstring(len, raw);
	char *ss = sprintstring(s);
	assert(!strcoll(ss, rawstring(s)), "string: sprint");
	free(ss);
	freestring(s);
}

void test_comparestring_less() {
	string a = mkstring(1, "a");
	string b = mkstring(1, "b");
	assert(comparestring(a, b) < 0, "string: compare, less");
	freestring(a);
	freestring(b);
}

void test_comparestring_eq() {
	string a = mkstring(1, "a");
	string aa = mkstring(1, "a");
	assert(comparestring(a, aa) == 0, "string: compare, eq");
	freestring(a);
	freestring(aa);
}

void test_comparestring_greater() {
	string a = mkstring(1, "a");
	string b = mkstring(1, "b");
	assert(comparestring(b, a) > 0, "string: compare, less");
	freestring(a);
	freestring(b);
}

void test_substring_lower_boundary() {
	char *s = "some string";
	long len = strlen(s);
	string st = mkstring(len, s);
	string sub = substr(st, 0, -1);
	assert(byteslen(sub) == 0, "string: substring, lower boundary");
	freestring(st);
	freestring(sub);
}

void test_substring_higher_boundary() {
	char *s = "some string";
	long len = strlen(s);
	string st = mkstring(len, s);
	string sub = substr(st, len, 1);
	assert(byteslen(sub) == 0, "string: substring, lower boundary");
	freestring(st);
	freestring(sub);
}

void test_substring() {
	char *s = teststring;
	long len = strlen(s);
	string st = mkstring(len, s);
	string sub = substr(st, 1, 2);
	assert(!strcoll(rawstring(sub), "űz"), "string: substring");
	freestring(st);
	freestring(sub);
}

void test_substring_utf8() {
	string s = mkstring(strlen(teststring), teststring);
	string sub = substr(s, 0, stringlen(s));
	assert(!strcoll(rawstring(sub), rawstring(sub)), "string: substring, utf8");
	freestring(sub);
	freestring(s);
}

void test_append_zero() {
	string s = appendstr(0, 0);
	assert(byteslen(s) == 0, "string: append, zero");
	freestring(s);
}

void test_append_one() {
	char *raw = "some string";
	long len = strlen(raw);
	string s = mkstring(len, raw);
	string a = appendstr(1, &s);
	assert(comparestring(a, s) == 0, "string: append, one");
	freestring(s);
	freestring(a);
}

void test_append_multiple() {
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

	string a = appendstr(3, ss);
	assert(!strcoll(rawstring(a), "some strings to append"), "string: append, multiple");

	free(ss);
	freestring(a);
	freestring(s1);
	freestring(s2);
	freestring(s3);
}

void test_clonestring() {
	char *raw = "some string";
	long len = strlen(raw);
	string s = mkstring(len, raw);
	string cs = clonestring(s);
	assert(!strcoll(rawstring(cs), rawstring(s)), "string: clone");
	freestring(cs);
	freestring(s);
}

void test_bytestochars() {
	long c = bytestochars(strlen(teststring) - 1, teststring + 1);
	assert(c == 12, "string: bytestochars");
}

int main() {
	initsys();

	test_init_free();
	test_bytes_len();
	test_bytes_length_from_arg();
	test_string_length_damaged();
	test_string_length();
	test_raw();
	test_raw_len_from_arg();
	test_returns_null_terminated();
	test_isutf8_positive();
	test_isutf8_negative();
	test_sprint();
	test_comparestring_less();
	test_comparestring_eq();
	test_comparestring_greater();
	test_substring_lower_boundary();
	test_substring_higher_boundary();
	test_substring();
	test_substring_utf8();
	test_append_zero();
	test_append_one();
	test_append_multiple();
	test_clonestring();
	test_bytestochars();

	return 0;
}
