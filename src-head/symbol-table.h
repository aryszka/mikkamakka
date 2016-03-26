struct symtable;
typedef struct symtable *symtable;

symtable mksymtable();
int hassym(symtable, size_t, char *);
void *getsym(symtable, size_t, char *);
void setsym(symtable, size_t, char *, void *);
void freesymtable(symtable);
