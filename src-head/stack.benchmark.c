#include "sys.h"
#include "testing.h"
#include "compound-types.h"
#include "stack.h"

void benchmark_init_free_100() {
	int m = 100;
	for (int i = 0; i < m; i++) {
		stack s = mkstack();
		freestack(s);
	}
}

void benchmark_push_100() {
	int m = 100;
	stack s = mkstack();

	for (int i = 0; i < m; i++) {
		push(s, (void *)1);
	}

	freestack(s);
}

void benchmark_push_pop_100() {
	int m = 100;
	stack s = mkstack();

	for (int i = 0; i < m; i++) {
		push(s, (void *)1);
		pop(s);
	}

	freestack(s);
}

void benchmark_peek_100() {
	int m = 100;

	stack s = mkstack();
	push(s, (void *)1);

	for (int i = 0; i < m; i++) {
		peek(s);
	}

	freestack(s);
}

int main(int argc, char **argv) {
	initsys();
	int n = repeatcount(argc, argv);

	int err = 0;
	err += benchmark(n, &benchmark_init_free_100, "stack: init and free, 100 times");
	err += benchmark(n, &benchmark_push_100, "stack: push, 100 times");
	err += benchmark(n, &benchmark_push_pop_100, "stack: push and pop, 100 times");
	err += benchmark(n, &benchmark_peek_100, "stack: peek, 100 times");
	return err;
}
