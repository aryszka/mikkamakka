// require first:
// compound-types.h

stack mkstack();
void push(stack, void *);
void *pop(stack);
void *peek(stack);
valuenode allstackvals(stack, valuenode);
envnode allstackenvs(stack, envnode);
void freestack(stack);
