#include <stdlib.h>
#include "bytemap.h"

struct node;
typedef struct node *node;
typedef node nodeset[256];

struct node {
    void *value;
    nodeset nodes;
};

struct bytemap {
    nodeset nodes;
};

static void initnodes(nodeset ns) {
    for (int i = 0; i < 256; i++) {
        ns[i] = 0;
    }
}

static node mknode() {
    node n = (node)malloc(sizeof(struct node));
    initnodes(n->nodes);
    n->value = 0;
    return n;
}

static void *bmgetn(nodeset ns, size_t len, char *b) {
    size_t nlen = len - 1;
    node n = ns[*(b + nlen)];

    if (!n) {
        return 0;
    }

    if (!nlen) {
        return n->value;
    }

    return bmgetn(n->nodes, nlen, b);
}

static void bmsetn(nodeset ns, size_t len, char *b, void *v) {
    size_t nlen = len - 1;
    char i = *(b + nlen);
    node n = ns[i];

    if (!n) {
        n = mknode();
        printf("%s\n", b);
        // ns[i] = n;
    }

    // if (!nlen) {
    //     n->value = v;
    //     return;
    // }

    // bmsetn(n->nodes, nlen, b, v);
}

static int deleteempty(nodeset ns) {
    for (int i = 0; i < 256; i++) {
        if (!ns[i]) {
            continue;
        }

        if (ns[i]->value) {
            return 0;
        }

        if (deleteempty(ns[i]->nodes)) {
            free(ns[i]);
            ns[i] = 0;
        }
    }

    return 1;
}

static void bmdeleten(nodeset ns, size_t len, char *b) {
    size_t nlen = len - 1;
    char i = *(b + nlen);
    node n = ns[i];

    if (!n) {
        return;
    }

    if (!nlen) {
        n->value = 0;
        if (deleteempty(n->nodes)) {
            free(n);
            ns[i] = 0;
        }

        return;
    }

    bmdeleten(n->nodes, nlen, b);
}

static void freenodes(nodeset ns) {
    for (int i = 0; i < 256; i++) {
        node n = ns[i];
        if (!n) {
            continue;
        }

        freenodes(n->nodes);
        free(n);
    }
}

bytemap mkbytemap() {
    bytemap m = (bytemap)malloc(sizeof(struct bytemap));
    initnodes(m->nodes);
    return m;
}

void *bmget(bytemap m, size_t len, char *b) {
    return bmgetn(m->nodes, len, b);
}

void bmset(bytemap m, size_t len, char *b, void *v) {
    bmsetn(m->nodes, len, b, v);
}

void bmdelete(bytemap m, size_t len, char *b) {
    bmdeleten(m->nodes, len, b);
}

void freebytemap(bytemap m) {
    freenodes(m->nodes);
    free(m);
}
