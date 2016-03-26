#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <errno.h>
#include <sys/time.h>
#include <sys/resource.h>
#include "sys.h"
#include "testing.h"

static int millis(struct timeval t) {
	return t.tv_sec * 1000 + t.tv_usec / 1000;
}

char *locnumberstr(char *sin) {
	int len = strlen(sin) + 1;
	char *nout = malloc(len);
	for (int i = 0; i < len; i++) {
		char curr = *(sin + i);
		if (curr == '.') {
			curr = getdecchar();
		}
		*(nout + i) = curr;
	}

	return nout;
}

void assert(int test, char *message) {
	if (!test) {
		fprintf(stderr, "failed: %s\n", message);
		exit(1);
	}
}

void times(int n, testproc p) {
	for (int i = 0; i < n; i++) {
		(*p)();
	}
}

resources rusage() {
	struct rusage r;
	if(getrusage(RUSAGE_SELF, &r)) {
		fprintf(stderr, "failed to retrieve resource usage: %s\n", strerror(errno));
		exit(errno);
	}

	resources res;
	res.utime = millis(r.ru_utime);
	res.stime = millis(r.ru_stime);
	res.rss = r.ru_maxrss;
	return res;
}

resources diffrusage(resources r1, resources r2) {
	resources d;
	d.utime = r1.utime - r2.utime;
	d.stime = r1.stime - r2.stime;
	d.rss = r1.rss - r2.rss;
	return d;
}

int benchmark(int n, testproc p, char *testname) {
	fprintf(
		stderr,
		"benchmark/%s * %d (utime, stime, rss): ",
		testname, n
	);

	resources start = rusage();
	times(n, p);
	resources end = rusage();

	resources diff = diffrusage(end, start);
	fprintf(
		stderr,
		"%d, %d, %ld\n",
		diff.utime, diff.stime, diff.rss
	);

	if (diff.rss > 0) {
		fprintf(stderr, ">>> rss increased, possible leak (%s)\n", testname);
		return 1;
	}

	return 0;
}

int repeatcount(int argc, char **argv) {
	if (argc > 1) {
		return atoi(*(argv + 1));
	}

	return 10000;
}
