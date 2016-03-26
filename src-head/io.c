#include <stdlib.h>
#include "sysio.h"
#include "error.h"
#include "io.h"

struct file {
	sfile sf;
};

static char *getsysfilemode(iofilemode mode) {
	switch (mode) {
	case ioread:
		return "r";
	case ioreadwrite:
		return "r+";
	case iowrite:
		return "w";
	case ioreadwritecreate:
		return "w+";
	case ioappend:
		return "a";
	case ioreadappend:
		return "a+";
	}
}

static sseekmode getsysseekmode(ioseekmode mode) {
	switch (mode) {
	case iostart:
		return sstart;
	case iocurr:
		return scurr;
	case ioend:
		return send;
	}
}

static file mkfile(sfile sf) {
	file rf = malloc(sizeof(struct file));
	rf->sf = sf;
	return rf;
}

static void freefile(file f) {
	free(f);
}

void openfile(char *n, iofilemode mode, file *f, ioerror *err) {
	ioerror rerr;
	sfile sf = sfopen(n, getsysfilemode(mode));
	if (!sf) {
		*f = 0;

		rerr.code = fileopenfailed;
		rerr.syserror = serrno;
		*err = rerr;

		return;
	}

	*f = mkfile(sf);

	rerr.code = 0;
	*err = rerr;
}

ioerror closefile(file f) {
	ioerror err;
	if (sfclose(f->sf)) {
		err.code = fileclosefailed;
		err.syserror = serrno;
		return err;
	}

	err.code = 0;
	freefile(f);
	return err;
}

ioerror seekfile(file f, long offset, ioseekmode mode) {
	ioerror err;
	if (sfseek(f->sf, offset, getsysseekmode(mode))) {
		err.code = fileseekfailed;
		err.syserror = serrno;
		return err;
	}

	err.code = 0;
	return err;
}

void readfile(file f, long len, char *c, long *rlen, ioerror *err) {
	ioerror rerr;
	long readlen = sfread(c, len, f->sf);
	if (readlen == len) {
		rerr.code = 0;
		*err = rerr;
		*rlen = readlen;
		return;
	}

	if (sfeof(f->sf)) {
		rerr.code = eof;
		*err = rerr;
		*rlen = readlen;
		return;
	}

	rerr.code = filereadfailed;
	rerr.syserror = serrno;
	*err = rerr;
	*rlen = readlen;
}

ioerror writefile(file f, long len, char *c) {
	ioerror err;
	long wlen = sfwrite(c, len, f->sf);
	if (wlen == len) {
		err.code = 0;
		return err;
	}

	err.code = filewritefailed;
	err.syserror = serrno;
	return err;
}

void initmodule_io() {
	iostdin = mkfile(sstdin);
	iostdout = mkfile(sstdout);
	iostderr = mkfile(sstderr);
}

void freemodule_io() {
	freefile(iostdin);
	freefile(iostdout);
	freefile(iostderr);
}
