#include <stdlib.h>
#include <string.h>
#include "sys.h"
#include "testing.h"
#include "error.h"
#include "sysio.h"
#include "sysio.mock.h"
#include "io.h"

void benchmark_open_close_100() {
	int m = 100;

	file f;
	ioerror err;

	for (int i = 0; i < m; i++) {
		openfile(fnok, ioread, &f, &err);
		closefile(f);
	}
}

void benchmark_open_fail_100() {
	int m = 100;

	file f;
	ioerror err;

	for (int i = 0; i < m; i++) {
		openfile(fnopenfail, ioread, &f, &err);
	}
}

void benchmark_seek_100() {
	int m = 100;

	file f;
	ioerror err;
	openfile(fnok, ioread, &f, &err);

	for (int i = 0; i < m; i++) {
		seekfile(f, 30, iostart);
	}

	closefile(f);
}

void benchmark_seek_fail_100() {
	int m = 100;

	file f;
	ioerror err;
	openfile(fnseekfail, ioread, &f, &err);

	for (int i = 0; i < m; i++) {
		seekfile(f, 30, iostart);
	}

	closefile(f);
}

void benchmark_read_100() {
	int m = 100;

	file f;
	ioerror err;
	openfile(fnok, ioread, &f, &err);

	char *c = malloc(30);
	long rlen;
	for (int i = 0; i < m; i++) {
		readfile(f, 30, c, &rlen, &err);
	}

	free(c);
	closefile(f);
}

void benchmark_read_eof_100() {
	int m = 100;

	file f;
	ioerror err;
	openfile(fnok, ioread, &f, &err);

	char *c = malloc(testfilelength);
	long rlen;
	for (int i = 0; i < m; i++) {
		seekfile(f, 15, iostart);
		readfile(f, testfilelength, c, &rlen, &err);
	}

	free(c);
	closefile(f);
}

void benchmark_read_fail_100() {
	int m = 100;

	file f;
	ioerror err;
	openfile(fnreadfail, ioread, &f, &err);

	char *c = malloc(30);
	long rlen;
	for (int i = 0; i < m; i++) {
		readfile(f, 30, c, &rlen, &err);
	}

	free(c);
	closefile(f);
}

void benchmark_write_100() {
	int m = 100;

	file f;
	ioerror err;
	openfile(fnok, iowrite, &f, &err);

	char *c = "some data";
	long len = strlen(c);
	for (int i = 0; i < m; i++) {
		seekfile(f, 15, iostart);
		writefile(f, len, c);
	}

	closefile(f);
}

void benchmark_write_fail_100() {
	int m = 100;

	file f;
	ioerror err;
	openfile(fnwritefail, iowrite, &f, &err);

	char *c = "some data";
	long len = strlen(c);
	for (int i = 0; i < m; i++) {
		seekfile(f, 15, iostart);
		writefile(f, len, c);
	}

	closefile(f);
}

int main(int argc, char **argv) {
	initsys();
	initmodule_sysio();

	int n = repeatcount(argc, argv);

	int err = 0;
	err += benchmark(n, &benchmark_open_close_100, "io: open and close file, 100 times");
	err += benchmark(n, &benchmark_open_fail_100, "io: open fail, 100 times");
	err += benchmark(n, &benchmark_seek_100, "io: seek file, 100 times");
	err += benchmark(n, &benchmark_seek_fail_100, "io: seek fail, 100 times");
	err += benchmark(n, &benchmark_read_100, "io: read file, 100 times");
	err += benchmark(n, &benchmark_read_eof_100, "io: read file eof, 100 times");
	err += benchmark(n, &benchmark_read_fail_100, "io: read file eof, 100 times");
	err += benchmark(n, &benchmark_write_100, "io: write file, 100 times");
	err += benchmark(n, &benchmark_write_fail_100, "io: write fail, 100 times");

	freemodule_sysio();
	return err;
}
