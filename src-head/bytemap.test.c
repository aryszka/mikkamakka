#include <stdlib.h>
#include "sys.h"
#include "testing.h"
#include "bytemap.h"

void empty() {
    bytemap m = mkbytemap();
    void *v = bmget(m, 1, "x");
    assert(v == 0, "empty map returns null");
    freebytemap(m);
}

void setget() {
    bytemap m = mkbytemap();
    long n = 42;
    void *v = (void *)n;
    bmset(m, 1, "x", v);
    void *vback = bmget(m, 1, "x");
    long nback = (long)vback;
    assert(vback == v && nback == n, "get returns set value");
    freebytemap(m);
}

void setgetlarge() {
    bytemap m = mkbytemap();
    long n = 42;
    void *v = (void *)n;
    bmset(m, 3, "xyz", v);
    void *vback = bmget(m, 3, "xyz");
    long nback = (long)vback;
    assert(vback == v && nback == n, "get returns set value, large");
    freebytemap(m);
}

void nonexisting() {
    bytemap m = mkbytemap();
    long n = 42;
    void *v = (void *)n;
    bmset(m, 1, "xyz", v);
    void *vback = bmget(m, 3, "xxz");
    assert(vback == 0, "get returns 0 on non-existing");
    freebytemap(m);
}

void deleteempty() {
    bytemap m = mkbytemap();
    bmdelete(m, 1, "x");
    freebytemap(m);
}

void setdelete() {
    bytemap m = mkbytemap();
    long n = 42;
    void *v = (void *)n;
    bmset(m, 1, "x", v);
    bmdelete(m, 1, "x");
    void *vback = bmget(m, 1, "x");
    assert(vback == 0, "delete value");
    freebytemap(m);
}

void setdeletelarge() {
    bytemap m = mkbytemap();
    long n = 42;
    void *v = (void *)n;
    bmset(m, 3, "xyz", v);
    bmdelete(m, 3, "xyz");
    void *vback = bmget(m, 3, "xyz");
    assert(vback == 0, "delete value, large");
    freebytemap(m);
}

void deletehigher() {
    bytemap m = mkbytemap();
    long nn = 36;
    long nm = 42;
    void *vn = (void *)nn;
    void *vm = (void *)nm;
    bmset(m, 1, "z", vn);
    bmset(m, 3, "xyz", vm);
    bmdelete(m, 1, "z");
    void *vnback = bmget(m, 1, "z");
    void *vmback = bmget(m, 3, "xyz");
    long nmback = (long)vmback;
    assert(vnback == 0, "delete value, deleted");
    assert(vmback == vm && nmback == nm, "delete value, keep");
    freebytemap(m);
}

int main() {
	initsys();

    empty();
    setget();
    setgetlarge();
    setdelete();
    setdeletelarge();
    nonexisting();
    deleteempty();

	return 0;
}
