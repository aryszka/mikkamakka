#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include "sysio.h"
#include "sysio.mock.h"

struct sfile {
	char *n;
	char *mode;
	int pos;
	int eof;
};

struct testcontent;
typedef struct testcontent *testcontent;

struct testcontent {
	char *name;
	int len;
	char *content;
	testcontent prev;
};

testcontent testfiles;

static testcontent findfile(testcontent f, char *name) {
	if (!f) {
		return 0;
	}

	if (!strcoll(f->name, name)) {
		return f;
	}

	return findfile(f->prev, name);
}

sfile sfopen(char *n, char *mode) {
	if (!strcoll(n, fnopenfail)) {
		serrno = mnopenfailed;
		return 0;
	}

	sfile sf = malloc(sizeof(struct sfile));
	sf->n = n;
	sf->mode = mode;
	sf->pos = 0;
	return sf;
}

int sfclose(sfile sf) {
	if (!strcoll(sf->n, fnclosefail)) {
		serrno = mnclosefailed;
		return EOF;
	}

	free(sf);
	return 0;
}

int sfseek(sfile sf, int offset, sseekmode mode) {
	if (!strcoll(sf->n, fnseekfail)) {
		serrno = mnseekfailed;
		return -1;
	}

	int len = testfilelength;
	testcontent f = findfile(testfiles, sf->n);
	if (f) {
		len = f->len;
	}

	sf->eof = 0;
	sf->pos = offset;
	if (sf->pos > len) {
		sf->pos = len;
	}

	return 0;
}

int sfread(char *c, int len, sfile sf) {
	if (!strcoll(sf->n, fnreadfail)) {
		serrno = mnreadfailed;
		return len / 2;
	}

	int tlen = testfilelength;
	testcontent f = findfile(testfiles, sf->n);
	if (f) {
		tlen = f->len;
	}

	sf->eof = 0;
	int rlen = len;
	if (sf->pos + len > tlen) {
		rlen = tlen - sf->pos;
		sf->eof = 1;
	}

	if (f) {
		for (int i = 0; i < rlen; i++) {
			*(c + i) = *(f->content + sf->pos + i);
		}
	}

	sf->pos += rlen;

	return rlen;
}

int sfeof(sfile sf) {
	return sf->eof;
}

int sfwrite(char *c, int len, sfile sf) {
	if (!strcoll(sf->n, fnwritefail)) {
		serrno = mnwritefailed;
		return len / 2;
	}

	sf->eof = 0;
	return len;
}

void initfilecontent(char *name, char *content) {
	testcontent f = findfile(testfiles, name);
	if (f) {
		f->content = content;
		return;
	}

	f = malloc(sizeof(struct testcontent));
	f->name = name;
	f->len = strlen(content);
	f->content = content;
	f->prev = testfiles;
	testfiles = f;
}

void cleartestcontent() {
	if (!testfiles) {
		return;
	}

	testcontent f = testfiles;
	testfiles = f->prev;
	free(f);
	cleartestcontent();
}

void initmodule_sysio() {
	fnok = "f0";
	fnopenfail = "f1";
	fnclosefail = "f2";
	fnseekfail = "f3";
	fnreadfail = "f4";
	fnwritefail = "f5";

	testfilelength = 120;
	testfiles = 0;

	serrno = 0;
}

void freemodule_sysio() {
	cleartestcontent();
}
