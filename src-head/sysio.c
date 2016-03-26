#include <stdlib.h>
#include <stdio.h>
#include <errno.h>
#include "sysio.h"

struct sfile {
	FILE *f;
	int eof;
};

static int getwhence(sseekmode mode) {
	switch (mode) {
	case sstart:
		return SEEK_SET;
	case scurr:
		return SEEK_CUR;
	case send:
		return SEEK_END;
	}
}

static sfile mkfile(FILE *f) {
	sfile sf = malloc(sizeof(struct sfile));
	sf->f = f;
	sf->eof = 0;
	return sf;
}

static void freesfile(sfile sf) {
	free(sf);
}

sfile sfopen(char *n, char *mode) {
	FILE *f = fopen(n, mode);
	if (!f) {
		serrno = errno;
		return 0;
	}

	return mkfile(f);
}

int sfclose(sfile sf) {
	int ret = fclose(sf->f);
	if (ret) {
		serrno = errno;
	} else {
		freesfile(sf);
	}

	return ret;
}

int sfseek(sfile sf, int offset, sseekmode mode) {
	sf->eof = 0;
	int ret = fseek(sf->f, offset, getwhence(mode));
	if (ret) {
		serrno = errno;
	}

	return ret;
}

int sfread(char *c, int len, sfile sf) {
	sf->eof = 0;
	int ret = fread(c, 1, len, sf->f);
	if (ret == len) {
		return ret;
	}

	if (feof(sf->f)) {
		sf->eof = 1;
		return ret;
	}

	serrno = errno;
	return ret;
}

int sfeof(sfile sf) {
	return sf->eof;
}

int sfwrite(char *c, int len, sfile sf) {
	sf->eof = 0;
	int ret = fwrite(c, 1, len, sf->f);
	if (ret == len) {
		return ret;
	}

	serrno = errno;
	return ret;
}

void initmodule_sysio() {
	serrno = errno;
	sstdin = mkfile(stdin);
	sstdout = mkfile(stdout);
	sstderr = mkfile(stderr);
}

void freemodule_sysio() {
	freesfile(sstdin);
	freesfile(sstdout);
	freesfile(sstderr);
}
