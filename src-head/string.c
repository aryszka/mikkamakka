#include <stdlib.h>
#include <string.h>
#include <unistr.h>
#include "string.h"

struct string {
	long blen;
	long len;
	char *raw;
};

static void copyrawstringto(char *to, long len, char *from) {
	for (long i = 0; i < len; i++) {
		*(to + i) = *(from + i);
	}
}

static char *copyrawstring(long len, char *raw) {
	char *c = malloc(len + 1);
	copyrawstringto(c, len, raw);
	*(c + len) = 0;
	return c;
}

static long utf8bytes(char *s, long len) {
	const uint8_t *c = u8_check((uint8_t *)s, len);
	if (!c) {
		return len;
	}

	return (char *)c - s;
}

static long utf8bytepos(char *s, long blen, long pos) {
	if (blen < pos) {
		return -1;
	}

	long valid = pos;
	for (long i = pos; i <= blen; i++) {
		valid = utf8bytes(s, i);
		long utf8len = u8_mbsnlen((uint8_t *)s, valid);
		if (utf8len == pos) {
			break;
		}
	}

	return valid;
}

string mkstring(long len, char *raw) {
	string s = malloc(sizeof(struct string));
	s->blen = len;
	s->raw = copyrawstring(len, raw);
	s->len = -1;
	return s;
}

string clonestring(string s) {
	return appendstr(1, &s);
}

long byteslen(string s) {
	return s->blen;
}

long bytestochars(long len, char *s) {
	long u8part = utf8bytes(s, len);
	// long u8part = 0;
	return u8_mbsnlen((uint8_t *)s, u8part);
}

long stringlen(string s) {
	if (s->len >= 0) {
		return s->len;
	}

	s->len = bytestochars(byteslen(s), rawstring(s));
	return s->len;
}

string substr(string s, long from, long len) {
	if (from < 0) {
		from = 0;
	}

	if (len < 0) {
		len = 0;
	}

	long slen = stringlen(s);
	if (from + len > slen) {
		len = slen - from;
	}

	long blen = byteslen(s);
	char *raw = rawstring(s);
	long bytefrom = utf8bytepos(raw, blen, from);
	long byteto = utf8bytepos(raw, blen, from + len);
	long bytelen = byteto - bytefrom;
	return mkstring(bytelen, raw + bytefrom);
}

string appendstr(long len, string *ss) {
	long totallen = 0;
	for (int i = 0; i < len; i++) {
		totallen += byteslen(*(ss + i));
	}

	char *raw = malloc(totallen + 1);

	long copiedlen = 0;
	for (int i = 0; i < len; i++) {
		string s = *(ss + i);
		long slen = byteslen(s);
		copyrawstringto(raw + copiedlen, slen, rawstring(s));
		copiedlen += slen;
	}

	*(raw + totallen) = 0;

	string a = mkstring(totallen, raw);
	free(raw);
	return a;
}

char *rawstring(string s) {
	return s->raw;
}

int comparestring(string s1, string s2) {
	return strcoll(rawstring(s1), rawstring(s2));
}

char *sprintstring(string s) {
	return copyrawstring(s->blen, rawstring(s));
}

int isutf8(string s) {
	const uint8_t *check = u8_check((uint8_t *)s->raw, s->blen);
	return check == 0;
}

void freestring(string s) {
	free(s->raw);
	free(s);
}
