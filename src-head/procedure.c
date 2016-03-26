#include <stdlib.h>
#include <stdio.h>
#include "error.h"
#include "number.h"
#include "string.h"
#include "compound-types.h"
#include "stack.h"
#include "io.h"
#include "value.h"
#include "procedure.h"

typedef enum {
	pprimitive,
	pcompiled
} proctype;

struct procedure {
	proctype type;
	primitive primitive;
	value label;
	environment env;
};

procedure mkprocedure(proctype type) {
	procedure p = malloc(sizeof(struct procedure));
	p->type = type;
	return p;
}

procedure mkprimitiveproc(primitive p) {
	procedure proc = mkprocedure(pprimitive);
	proc->primitive = p;
	return proc;
}

procedure mkcompiledproc(value label, environment env) {
	procedure proc = mkprocedure(pcompiled);
	proc->label = label;
	proc->env = env;
	return proc;
}

int isprimitive(procedure p) {
	return p->type == pprimitive;
}

value applyprimitive(procedure p, value args) {
	if (!isprimitive(p)) {
		error(invalidtype, "primitive procedure expected");
		return 0;
	}

	return (*(p->primitive))(args);
}

value proclabel(procedure p) {
	if (isprimitive(p)) {
		error(invalidtype, "compiled procedure expected");
		return 0;
	}

	return p->label;
}

environment procenv(procedure p) {
	if (isprimitive(p)) {
		error(invalidtype, "compiled procedure expected");
		return 0;
	}

	return p->env;
}

void freeproc(procedure p) {
	free(p);
}
