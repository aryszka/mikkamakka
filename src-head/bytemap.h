struct bytemap;
typedef struct bytemap *bytemap;

bytemap mkbytemap();
void *bmget(bytemap, size_t, char *);

// setting 0 undefined
void bmset(bytemap, size_t, char *, void *);

void bmdelete(bytemap, size_t, char *);
void freebytemap(bytemap);
