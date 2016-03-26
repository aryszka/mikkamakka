#include <string.h>
#include "sys.h"
#include "testing.h"
#include "symbol-table.h"

void test_init_free() {
	symtable st = mksymtable();
	freesymtable(st);
}

void test_hassym_false() {
	symtable st = mksymtable();
	char *sname = "some-symbol";
	int slen = strlen(sname);
	assert(!hassym(st, slen, sname), "symtable: hassym, false");
	freesymtable(st);
}

void test_hassym() {
	symtable st = mksymtable();
	char *sname = "some-symbol";
	int slen = strlen(sname);
    long n = 42;
	setsym(st, slen, sname, &n);
	assert(hassym(st, slen, sname), "symtable: hassym");
	freesymtable(st);
}

void test_getsym() {
	symtable st = mksymtable();
	char *sname = "some-symbol";
	int slen = strlen(sname);
    long n = 42;
	setsym(st, slen, sname, &n);
    void *vback = getsym(st, slen, sname);
    long nback = *((long *)vback);
	assert(nback == n, "symtable: getsym");
	freesymtable(st);
}

void test_setsym_new() {
	symtable st = mksymtable();
	char *sname = "some-symbol";
	int slen = strlen(sname);
    long n = 42;
	setsym(st, slen, sname, &n);
	assert(*(long *)getsym(st, slen, sname) == n, "symtable: setsym, new");
	freesymtable(st);
}

void test_setsym_existing() {
	symtable st = mksymtable();
	char *sname = "some-symbol";
	int slen = strlen(sname);
    long nn = 36;
    long nm = 42;
	setsym(st, slen, sname, &nn);
	setsym(st, slen, sname, &nm);
	assert(*(long *)getsym(st, slen, sname) == nm, "symtable: setsym, existing");
	freesymtable(st);
}

int main() {
	initsys();

	test_init_free();
	test_hassym_false();
	test_hassym();
	test_getsym();
	test_setsym_new();
	test_setsym_existing();

	return 0;
}
