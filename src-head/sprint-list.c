#include <stdlib.h>
#include <string.h>
#include "error.h"
#include "number.h"
#include "string.h"
#include "compound-types.h"
#include "stack.h"
#include "io.h"
#include "value.h"
#include "sprint-list.h"

static char *copyrawstring(int len, char *raw) {
	char *c = malloc(len + 1);
	for (int i = 0; i < len; i++) {
		*(c + i) = *(raw + i);
	}
	*(c + len) = 0;
	return c;
}

static char *concatstring(int len, char **ss) {
	int *lens = malloc(len * sizeof(int));
	int slen = 0;
	for (int i = 0; i < len; i++) {
		*(lens + i) = strlen(*(ss + i));
		slen += *(lens + i);
	}

	char *s = malloc(slen + 1);
	int c = 0;
	for (int i = 0; i < len; i++)
	for (int j = 0; j < *(lens + i); j++) {
		*(s + c) = *(*(ss + i) + j);
		c++;
	}

	*(s + c) = 0;
	return s;
}

static char *sprintlist(value p, int inlist) {
	// todo: indeed
	if (isnulltype(p)) {
		if (inlist) {
			return copyrawstring(1, "");
		}

		return copyrawstring(2, "()");
	}

	value car = carval(p);
	value cdr = cdrval(p);

	char *cars = sprintraw(car);
	char *carlead = inlist ? " " : "(";

	char *cdrs;
	char *cdrlead;
	if (ispairtype(cdr) || isnulltype(cdr)) {
		cdrs = sprintlist(cdr, 1);
		cdrlead = "";
	} else {
		cdrs = sprintraw(cdr);
		cdrlead = " . ";
	}

	char *s;
	if (inlist) {
		s = concatstring(4, (char *[]){" ", cars, cdrlead, cdrs});
	} else {
		s = concatstring(5, (char *[]){"(", cars, cdrlead, cdrs, ")"});
	}

	free(cars);
	free(cdrs);

	return s;
}

char *sprintlistraw(value p) {
	return sprintlist(p, 0);
}
