#include <stdlib.h>
#include "compound-types.h"
#include "stack.h"
#include "number.h"
#include "string.h"
#include "value.h"
#include "environment.h"

struct stacknode;
typedef struct stacknode *stacknode;

struct stacknode {
	stacknode prev;
	void *val;
};

struct stack {
	stacknode nodes;
};

static void freenodes(stacknode nodes) {
	if (!nodes) {
		return;
	}

	freenodes(nodes->prev);
	free(nodes);
}

static valuenode getallstackvals(stacknode nodes, valuenode prev) {
    if (!nodes) {
        return prev;
    }

    valuenode n = malloc(sizeof(struct valuenode));
    n->v = (value)nodes->val;
    n->prev = getallstackvals(nodes->prev, prev);
    return n;
}

static envnode getallstackenvs(stacknode nodes, envnode prev) {
    if (!nodes) {
        return prev;
    }

    envnode n = malloc(sizeof(struct envnode));
    n->env = (environment)nodes->val;
    n->prev = getallstackenvs(nodes->prev, prev);
    return n;
}

stack mkstack() {
	stack s = malloc(sizeof(struct stack));
	s->nodes = 0;
	return s;
}

void push(stack s, void *v) {
	stacknode n = malloc(sizeof(struct stacknode));
	n->prev = s->nodes;
	n->val = v;
	s->nodes = n;
}

void *pop(stack s) {
	if (!s->nodes) {
		return 0;
	}

	stacknode n = s->nodes;
	void *v = n->val;
	s->nodes = n->prev;
	free(n);
	return v;
}

void *peek(stack s) {
	if (!s->nodes) {
		return 0;
	}

	return s->nodes->val;
}

valuenode allstackvals(stack s, valuenode prev) {
    return getallstackvals(s->nodes, prev);
}

envnode allstackenvs(stack s, envnode prev) {
    return getallstackenvs(s->nodes, prev);
}

void freestack(stack s) {
	freenodes(s->nodes);
	free(s);
}
