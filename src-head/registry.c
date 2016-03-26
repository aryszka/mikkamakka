#include <stdlib.h>
#include <stdio.h>
#include "compound-types.h"
#include "stack.h"
#include "registry.h"
#include "number.h"
#include "string.h"
#include "error.h"
#include "io.h"
#include "value.h"
#include "environment.h"
#include "stack.h"

registryflag flag;
regmachine regs = 0;
valuenode nodes;
envnode envnodes;
long size;
long highercollectsize;
long lowercollectsize;

static void freeall(valuenode current) {
    if (!current) {
        return;
    }

    freeall(current->prev);
    free(current);
}

static void freeallenv(envnode current) {
    if (!current) {
        return;
    }

    freeallenv(current->prev);
    free(current);
}

static void markval(value v);
static void markvaluenodes(valuenode nodes);
static void freevaluenodes(valuenode nodes);
static void markenvnodes(envnode nodes);
static void freenvnodes(envnode nodes);
static void markenv(environment env);
static void markstack(stack s);
static void markenvstack(stack s);

static void markval(value v) {
    registryflag vf = getregistryflag(v);

    if (vf == flag) {
        return;
    }

    if (vf != unregistered) {
        setregistryflag(v, flag);
    }

    if (ispairtype(v)) {
        markval(carval(v));
        markval(cdrval(v));
    }

    if (isproctype(v) && !isprimitiveproctype(v)) {
        markenv(valprocenv(v));
        markval(valproclabel(v));
    }
}

static void markvaluenodes(valuenode nodes) {
    if (!nodes) {
        return;
    }

    markvaluenodes(nodes->prev);
    markval(nodes->v);
}

static void freevaluenodes(valuenode nodes) {
    if (!nodes) {
        return;
    }

    freevaluenodes(nodes->prev);
    free(nodes);
}

static void markenvnodes(envnode nodes) {
    if (!nodes) {
        return;
    }

    markenvnodes(nodes->prev);
    markenv(nodes->env);
}

static void freenvnodes(envnode nodes) {
    if (!nodes) {
        return;
    }

    freenvnodes(nodes->prev);
    free(nodes);
}

static void markenv(environment env) {
    registryflag ef = getenvregistryflag(env);

    if (ef == flag) {
        return;
    }

    if (ef != unregistered) {
        setenvregistryflag(env, flag);
    }

    // valuenode values = allenvvalues(env, 0);
    // markvaluenodes(values);
    // freevaluenodes(values);

    // envnode envs = allenvenvs(env, 0);
    // markenvnodes(envs);
    // freenvnodes(envs);
}

static void markstack(stack s) {
    valuenode values = allstackvals(s, 0);
    markvaluenodes(values);
    freevaluenodes(values);
}

static void markenvstack(stack s) {
    envnode envs = allstackenvs(s, 0);
    markenvnodes(envs);
    freenvnodes(envs);
}

static valuenode collectvalues(valuenode nodes) {
    if (!nodes) {
        return 0;
    }

    valuenode prev = collectvalues(nodes->prev);

    if (getregistryflag(nodes->v) == flag) {
        nodes->prev = prev;
        return nodes;
    }

    // printf("freeing value: %ld\n", nodes->v);
    freeval(nodes->v);
    free(nodes);
    size--;
    return prev;
}

static envnode collectenvs(envnode nodes) {
    if (!nodes) {
        return 0;
    }

    envnode prev = collectenvs(nodes->prev);

    if (getenvregistryflag(nodes->env) == flag) {
        nodes->prev = prev;
        return nodes;
    }

    // printf("freeing env: ");
    // printf("%ld\n", nodes->env);
    freenv(nodes->env);
    free(nodes);
    size--;
    // printf("free done\n");
    return prev;
}

void registerval(value v) {
    if (!regs) {
        return;
    }

    // printf("registering value: %ld\n", v);
    valuenode n = malloc(sizeof(struct valuenode));
    n->v = v;
    n->prev = nodes;
    nodes = n;
    size++;
    setregistryflag(v, flag);
}

void registerenv(environment env) {
    if (!regs) {
        return;
    }

    // printf("registering env: %ld\n", env);
    envnode n = malloc(sizeof(struct envnode));
    n->env = env;
    n->prev = envnodes;
    envnodes = n;
    size++;
    setenvregistryflag(env, flag);
}

void collect() {
    if (size >= lowercollectsize && size < highercollectsize) {
        return;
    }

    flag = flag == left ? right : left;

    markval(regs->label);
    markval(regs->cont);
    markval(regs->proc);
    markval(regs->args);
    markval(regs->val);
    markenv(regs->env);
    markstack(regs->stack);
    markenvstack(regs->envstack);

    nodes = collectvalues(nodes);
    envnodes = collectenvs(envnodes);

    // todo: make this adjustment logic more performing

    if (size >= highercollectsize) {
        lowercollectsize = highercollectsize;
        highercollectsize = highercollectsize * 2;
    }

    if (size < lowercollectsize) {
        highercollectsize = lowercollectsize;
        lowercollectsize = lowercollectsize / 2;
    }

    // printf("collected: %ld %ld %ld\n", size, highercollectsize, lowercollectsize);
}

void setregistermachine(regmachine rm) {
    regs = rm;
}

void initmodule_registry() {
    flag = left;
    nodes = 0;
    envnodes = 0;
    size = 0;
    lowercollectsize = 0;

    // todo: raise this
    highercollectsize = 1;
}

void freemodule_registry() {
    freeall(nodes);
    freeallenv(envnodes);
}
