#include <stdlib.h>
#include "error.h"
#include "regex.h"

#include <stdio.h>

#define PCRE2_CODE_UNIT_WIDTH 8
#include <pcre2.h>

struct regex {
	pcre2_code *code;
	pcre2_match_data *matchdata;
};

struct match {
	long len;
	submatch *smatches;
};

static uint32_t rxconvertoptions(rxflags flags) {
	uint32_t options = PCRE2_UTF;

	if (flags & rxignorecase) {
		options |= PCRE2_CASELESS;
	}

	if (flags & rxmultiline) {
		options |= PCRE2_MULTILINE;
	}

	return options;
}

static void rxcompileerror(int errorcode, int bufm) {
	int len = (bufm + 1) * 1024;
	char *msg = malloc(len);
	int msglen = pcre2_get_error_message(errorcode, (PCRE2_UCHAR *)msg, len);
	if (msglen < 0) {
		free(msg);
		rxcompileerror(errorcode, bufm + 1);
		return;
	}

	error(invalidregex, msg);
	free(msg);
}

static pcre2_code *rxcompile(long len, char *raw, rxflags flags) {
	int errorcode;
	PCRE2_SIZE erroroffset;
	pcre2_code *code = pcre2_compile(
		(PCRE2_SPTR)raw, len, rxconvertoptions(flags),
		&errorcode, &erroroffset, 0
	);

	if (!code) {
		rxcompileerror(errorcode, 0);
		return 0;
	}

	return code;
}

static void rxreadmatch(regex rx, match m, int mlen) {
	m->len = mlen;
	m->smatches = malloc(mlen * sizeof(submatch));

	// todo: handle umatched subpatterns
	PCRE2_SIZE *ovector = pcre2_get_ovector_pointer(rx->matchdata);
	PCRE2_SIZE sindex;
	PCRE2_SIZE slen;
	for (int i = 0; i < mlen; i++) {
		sindex = *(ovector + 2 * i);
		slen = *(ovector + 2 * i + 1) - sindex;

		submatch sm;
		sm.index = sindex;
		sm.len = slen;

		*(m->smatches + i) = sm;
	}
}

regex mkregex(long len, char *raw, rxflags flags) {
	pcre2_code *code = rxcompile(len, raw, flags);
	if (!code) {
		return 0;
	}

	regex rx = malloc(sizeof(struct regex));
	rx->code = code;
	rx->matchdata = pcre2_match_data_create_from_pattern(rx->code, 0);
	return rx;
}

match matchrx(regex rx, long len, char *s) {
	int mlen = pcre2_match(
		rx->code, (PCRE2_SPTR)s, len,
		0, 0, rx->matchdata, 0
	);

	if (mlen < 0) {
		return 0;
	}

	match m = malloc(sizeof(struct match));
	rxreadmatch(rx, m, mlen);

	return m;
}

int matchlen(match m) {
	return m->len;
}

submatch smatch(match m, int i) {
	return *(m->smatches + i);
}

void freematch(match m) {
	free(m->smatches);
	free(m);
}

void freeregex(regex rx) {
	pcre2_code_free(rx->code);
	pcre2_match_data_free(rx->matchdata);
	free(rx);
}
