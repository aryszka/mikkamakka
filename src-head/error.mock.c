#include <stdlib.h>
#include "error.h"
#include "error.mock.h"

struct errornode {
	struct errornode *prev;
	errorcode code;
};

struct errornode *lasterror;

void error(errorcode code, char *info) {
	if (!code) {
		return;
	}

	struct errornode *node = malloc(sizeof(struct errornode));
	node->prev = lasterror;
	node->code = code;
	lasterror = node;
}

char *errorstring(errorcode code) {
	return "error";
}

errorcode poperror() {
	if (!lasterror) {
		return 0;
	}

	struct errornode *last = lasterror;
	errorcode code = last->code;
	lasterror = last->prev;
	free(last);
	return code;
}

void clearerrors() {
	for (;poperror(););
}

void initmodule_errormock() {
	lasterror = 0;
}

void freemodule_errormock() {
	clearerrors();
}
