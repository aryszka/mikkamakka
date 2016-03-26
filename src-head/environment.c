#include <stdlib.h>
#include "error.h"
#include "number.h"
#include "string.h"
#include "compound-types.h"
#include "stack.h"
#include "io.h"
#include "value.h"
#include "procedure.h"
#include "symbol-table.h"
#include "environment.h"
#include "registry.h"

#include <stdio.h>

struct environment {
    registryflag rflag;
	environment parent;
	symtable symtable;
};

static int defall(environment env, value symbols, value values) {
	if (!issymtype(symbols) &&
		((isnulltype(symbols) && !isnulltype(values)) ||
		(!isnulltype(symbols) && isnulltype(values)))) {
		error(invalidnumberofargs, "number of symbols and values don't match");
		return 1;
	}

	if (isnulltype(symbols)) {
		return 0;
	}

	if (issymtype(symbols)) {
		defvar(env, symbols, values);
		return 0;
	}

	defvar(env, carval(symbols), carval(values));
	return defall(env, cdrval(symbols), cdrval(values));
}

environment mkenvironment(environment parent) {
	environment env = malloc(sizeof(struct environment));
	env->parent = parent;
	env->symtable = mksymtable();
    env->rflag = unregistered;
    registerenv(env);
	return env;
}

void defvar(environment env, value sym, value val) {
    if (!issymtype(sym)) {
        error(invalidtype, "symbol expected");
        return;
    }

    printf("%u %s\n", valsymbollen(sym), valsymbolnameraw(sym));
	setsym(env->symtable, valsymbollen(sym), valsymbolnameraw(sym), val);
}

int hasvar(environment env, value sym) {
	int r = (
		hassym(env->symtable, valsymbollen(sym), valsymbolnameraw(sym)) ||
		(env->parent && hasvar(env->parent, sym))
	);
    return r;
}

value getvar(environment env, value sym) {
	if (hassym(env->symtable, valsymbollen(sym), valsymbolnameraw(sym))) {
		value r = getsym(env->symtable, valsymbollen(sym), valsymbolnameraw(sym));
        return r;
	}

	if (env->parent) {
		value r = getvar(env->parent, sym);
        return r;
	}

	error(symbolnotfound, valsymbolnameraw(sym));
	return 0;
}

void setvar(environment env, value sym, value val) {
	if (hassym(env->symtable, valsymbollen(sym), valsymbolnameraw(sym))) {
		setsym(env->symtable, valsymbollen(sym), valsymbolnameraw(sym), val);
		return;
	}

	if (env->parent) {
		setvar(env->parent, sym, val);
		return;
	}

	error(symbolnotfound, valsymbolnameraw(sym));
	return;
}

environment extenv(environment env, value symbols, value values) {
	environment extended = mkenvironment(env);
	if (defall(extended, symbols, values)) {
		// freenv(extended);
		return 0;
	}

	return extended;
}

registryflag getenvregistryflag(environment env) {
    registryflag r = env->rflag;
    return r;
}

void setenvregistryflag(environment env, registryflag f) {
    env->rflag = f;
}

// valuenode allenvvalues(environment env, valuenode prev) {
//     if (env->parent) {
//         prev = allenvvalues(env->parent, prev);
//     }
// 
//     valuenode r = allvalues(env->symtable, prev);
//     return r;
// }
// 
// envnode allenvenvs(environment env, envnode prev) {
//     if (env->parent) {
//         prev = allenvenvs(env->parent, prev);
//         envnode p = malloc(sizeof(struct envnode));
//         p->env = env->parent;
//         p->prev = prev;
//         prev = p;
//     }
// 
//     return allenvs(env->symtable, prev);
// }

void freenv(environment env) {
    // printf("really freeing env: ");
    // printf("%ld\n", env);
	freesymtable(env->symtable);
	free(env);
}
