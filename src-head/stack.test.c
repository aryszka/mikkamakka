#include "sys.h"
#include "testing.h"
#include "compound-types.h"
#include "stack.h"

void test_init_free() {
	stack s = mkstack();
	freestack(s);
}

void test_push() {
	stack s = mkstack();
	push(s, (void *)1);
	assert((int)peek(s) == 1, "stack: pushed to stack");
	freestack(s);
}

void test_pop_empty() {
	stack s = mkstack();
	assert(pop(s) == 0, "stack: pop empty stack");
	freestack(s);
}

void test_pop() {
	stack s = mkstack();
	push(s, (void *)1);
	assert((int)pop(s) == 1, "postack: p stack");
	assert(peek(s) == 0, "stack: after pop");
	freestack(s);
}

void test_peek_empty() {
	stack s = mkstack();
	assert(peek(s) == 0, "stack: peek empty stack");
	freestack(s);
}

void test_peek() {
	stack s = mkstack();
	push(s, (void *)1);
	assert((int)peek(s) == 1, "stack: peek empty stack");
	assert((int)peek(s) == 1, "stack: peek does not change");
	freestack(s);
}

int main() {
	initsys();

	test_init_free();
	test_push();
	test_pop_empty();
	test_pop();
	test_peek_empty();
	test_peek();

	return 0;
}
