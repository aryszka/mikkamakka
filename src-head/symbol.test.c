#include <stdlib.h>
#include <string.h>
#include "sys.h"
#include "testing.h"
#include "symbol.h"

void test_init_free() {
	char *name = "symbol";
	symbol s = mksymbol(strlen(name), name);
	freesymbol(s);
}

void test_symbol_name() {
	char *name = "symbol";
	symbol s = mksymbol(strlen(name), name);
	assert(!strcoll(symbolname(s), name), "symbol: name");
	freesymbol(s);
}

void test_sprint() {
	char *name = "symbol";
	int len = strlen(name);
	symbol s = mksymbol(len, name);
	char *ss = sprintsymbol(s);
	assert(!strcoll(ss, name), "symbol: sprint");
	free(ss);
	freesymbol(s);
}

void test_symeq_false() {
	symbol s1 = mksymbol(3, "abc");
	symbol s2 = mksymbol(3, "Abc");
	assert(!symeq(s1, s2), "symbol: eq, false");
	freesymbol(s1);
	freesymbol(s2);
}

void test_symeq_true() {
	symbol s1 = mksymbol(3, "abc");
	symbol s2 = mksymbol(3, "abc");
	assert(symeq(s1, s2), "symbol: eq, true");
	freesymbol(s1);
	freesymbol(s2);
}

int main() {
	initsys();

	test_init_free();
	test_symbol_name();
	test_sprint();
	test_symeq_false();
	test_symeq_true();

	return 0;
}
