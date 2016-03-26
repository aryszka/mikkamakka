struct stack;
typedef struct stack *stack;

struct value;
typedef struct value *value;

struct valuenode;
typedef struct valuenode *valuenode;

struct valuenode {
    value v;
    valuenode prev;
};

struct pair;
typedef struct pair *pair;

typedef value (*primitive)(value args);
struct procedure;
typedef struct procedure *procedure;

struct environment;
typedef struct environment *environment;

// todo: make the environment a value
struct envnode;
typedef struct envnode *envnode;

struct envnode {
    environment env;
    envnode prev;
};

struct symtable;
typedef struct symtable *symtable;

struct rminternals;
typedef struct rminternals *rminternals;

struct regmachine {
	rminternals internals;
	int flag;
	value label;
	value cont;
	value proc;
	value args;
	value val;
	environment env;
	stack stack;
    stack envstack;
};

typedef struct regmachine *regmachine;

typedef enum {
    unregistered,
    left,
    right
} registryflag;
