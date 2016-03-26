#include <stdlib.h>
#include "error.h"
#include "number.h"
#include "string.h"
#include "compound-types.h"
#include "stack.h"
#include "io.h"
#include "value.h"
#include "pair.h"

struct pair {
	value a;
	value d;
};

pair mkpair(value a, value d) {
	pair p = malloc(sizeof(struct pair));
	p->a = a;
	p->d = d;
	return p;
}

value car(pair p) {
	return p->a;
}

value cdr(pair p) {
	return p->d;
}

void freepair(pair p) {
	free(p);
}
