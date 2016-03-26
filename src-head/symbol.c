#include <stdlib.h>
#include <string.h>
#include "symbol.h"

struct symbol {
	char *name;
    size_t len;
};

static char *copyrawstring(size_t len, char *s) {
	char *c = malloc(len + 1);
	for (int i = 0; i < len; i++) {
		*(c + i) = *(s + i);
	}
	*(c + len) = 0;
	return c;
}

symbol mksymbol(size_t len, char *name) {
	char *cn = copyrawstring(len, name);
	symbol s = malloc(sizeof(struct symbol));
	s->name = cn;
	return s;
}

char *symbolname(symbol s) {
	return s->name;
}

size_t symbollen(symbol s) {
    return s->len;
}

int symeq(symbol s1, symbol s2) {
	return !strcoll(symbolname(s1), symbolname(s2));
}

char *sprintsymbol(symbol s) {
	char *raw = symbolname(s);
	return copyrawstring(strlen(raw), raw);
}

void freesymbol(symbol s) {
	free(s->name);
	free(s);
}
