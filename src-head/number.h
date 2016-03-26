struct number;
typedef struct number *number;

number mknumc(size_t, char *);
number mknumi(long, long);
number mknumcsafe(size_t, char *);
number clonenum(number);
number sum(number, number);
number diff(number, number);
number bitor(number, number);
long rawint(number);
int isint(number);
int issmallint(number);
int comparenum(number, number);
char *sprintnum(number);
void freenum(number);
void initmodule_number();
void freemodule_number();
