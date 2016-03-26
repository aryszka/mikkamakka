#include <stdlib.h>
#include "bytemap.h"
#include "symbol-table.h"

struct symtable {
	bytemap map;
};

symtable mksymtable() {
	symtable st = malloc(sizeof(struct symtable));
    st->map = mkbytemap();
	return st;
}

void *getsym(symtable st, size_t len, char *b) {
    return bmget(st->map, len, b);
}

int hassym(symtable st, size_t len, char *b) {
    return !!getsym(st, len, b);
}

void setsym(symtable st, size_t len, char *b, void *v) {
    bmset(st->map, len, b, v);
}

void freesymtable(symtable st) {
    freebytemap(st->map);
	free(st);
}
