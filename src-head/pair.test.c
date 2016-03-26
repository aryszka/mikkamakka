#include "sys.h"
#include "error.h"
#include "testing.h"
#include "number.h"
#include "string.h"
#include "compound-types.h"
#include "io.h"
#include "value.h"
#include "pair.h"

void test_pair() {
	value n1 = mknumvali(1, 1);
	value n2 = mknumvali(2, 1);
	pair p = mkpair(n1, n2);
	assert(car(p) == n1, "pair: car");
	assert(cdr(p) == n2, "pair: cdr");
	freepair(p);
	freeval(n1);
	freeval(n2);
}

int main() {
	initsys();
	initmodule_value();

	test_pair();

	freemodule_value();
	return 0;
}
