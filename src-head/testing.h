typedef struct {
	int utime;
	int stime;
	long rss;
} resources;

typedef void (*testproc)();

char *locnumberstr(char *sin);
void assert(int test, char *message);
void times(int n, testproc p);
resources rusage();
resources diffrusage(resources r1, resources r2);
int benchmark(int n, testproc p, char *testname);
int repeatcount(int argc, char **argv);
