#include <stdlib.h>
#include <string.h>
#include "sys.h"
#include "testing.h"
#include "error.h"
#include "sysio.h"
#include "sysio.mock.h"
#include "io.h"

void test_open_close() {
	file f;
	ioerror err;
	openfile(fnok, ioread, &f, &err);
	assert(!err.code, "io: open file");
	err = closefile(f);
	assert(!err.code, "io: close file");
}

void test_open_failed() {
	file f;
	ioerror err;
	openfile(fnopenfail, ioread, &f, &err);
	assert(!f, "io: open failed file null");
	assert(err.code == fileopenfailed, "io: open failed code");
	assert(err.syserror == mnopenfailed, "io: open failed errno");
}

void test_close_failed() {
	file f;
	ioerror err;
	openfile(fnclosefail, ioread, &f, &err);
	err = closefile(f);
	assert(err.code == fileclosefailed, "io: close failed code");
	assert(err.syserror = mnclosefailed, "io: close failed errno");
}

void test_seek() {
	file f;
	ioerror err;
	openfile(fnok, ioread, &f, &err);
	err = seekfile(f, 30, iostart);
	assert(!err.code, "io: seek file");
	closefile(f);
}

void test_seek_failed() {
	file f;
	ioerror err;
	openfile(fnseekfail, ioread, &f, &err);
	err = seekfile(f, 30, iostart);
	assert(err.code == fileseekfailed, "io: seek failed code");
	assert(err.syserror == mnseekfailed, "io: seek failed errno");
	closefile(f);
}

void test_read() {
	file f;
	ioerror err;
	openfile(fnok, ioread, &f, &err);
	seekfile(f, 15, iostart);
	char *c = malloc(30);
	long rlen;
	readfile(f, 30, c, &rlen, &err);
	assert(rlen == 30, "io: read length");
	assert(!err.code, "io: read no error");
	free(c);
	closefile(f);
}

void test_read_eof() {
	file f;
	ioerror err;
	openfile(fnok, ioread, &f, &err);
	seekfile(f, 15, iostart);
	char *c = malloc(testfilelength);
	long rlen;
	readfile(f, testfilelength, c, &rlen, &err);
	assert(rlen == testfilelength - 15, "io: read length");
	assert(err.code == eof, "io: read no error");
	free(c);
	closefile(f);
}

void test_read_failed() {
	file f;
	ioerror err;
	openfile(fnreadfail, ioread, &f, &err);
	seekfile(f, 15, iostart);
	char *c = malloc(30);
	long rlen;
	readfile(f, 30, c, &rlen, &err);
	assert(err.code == filereadfailed, "io: read error");
	assert(err.syserror == mnreadfailed, "io: read errno");
	free(c);
	closefile(f);
}

void test_write() {
	file f;
	ioerror err;
	openfile(fnok, iowrite, &f, &err);
	seekfile(f, 15, iostart);
	char *c = "some data";
	long len = strlen(c);
	err = writefile(f, len, c);
	assert(!err.code, "io: write no error");
	closefile(f);
}

void test_write_failed() {
	file f;
	ioerror err;
	openfile(fnwritefail, iowrite, &f, &err);
	seekfile(f, 15, iostart);
	char *c = "some data";
	long len = strlen(c);
	err = writefile(f, len, c);
	assert(err.code == filewritefailed, "io: write error");
	assert(err.syserror == mnwritefailed, "io: write errno");
	closefile(f);
}

int main() {
	initsys();
	initmodule_sysio();

	test_open_close();
	test_open_failed();
	test_close_failed();
	test_seek();
	test_seek_failed();
	test_read();
	test_read_eof();
	test_read_failed();
	test_write();

	freemodule_sysio();
	return 0;
}
