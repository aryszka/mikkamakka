struct symbol;
typedef struct symbol *symbol;

symbol mksymbol(size_t len, char *name);
char *symbolname(symbol s);
size_t symbollen(symbol s);
int symeq(symbol s1, symbol s2);
char *sprintsymbol(symbol s);
void freesymbol(symbol s);
