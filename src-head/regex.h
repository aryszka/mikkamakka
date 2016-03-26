typedef enum {
	rxignorecase = 1,
	rxmultiline  = 2
} rxflags;

struct regex;
typedef struct regex *regex;

struct match;
typedef struct match *match;

typedef struct {
	long index;
	long len;
} submatch;

regex mkregex(long len, char *raw, rxflags flags);
match matchrx(regex rx, long len, char *s);
int matchlen(match m);
submatch smatch(match m, int i);
void freematch(match m);
void freeregex(regex rx);
