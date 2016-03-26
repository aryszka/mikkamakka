// requires first:
// stack.h
// compound-types.h
// value.h
// environment.h

regmachine mkregmachine();
value symval(char *);
value numvali(long);
value numvalc(char *);
value stringval(char *);
value pairval(value, value);
void saveenv(regmachine);
void restoreenv(regmachine);
void savereg(regmachine, void *);
void restorereg(regmachine, void **);
void gotoreg(regmachine, value);
void gotolabel(regmachine, long);
void gotoproc(regmachine);
void takeproclabel(regmachine);
void initreg(void **reg, value);
void getenvvar(regmachine, void **, char *);
void setenvvar(regmachine, void **, char *);
void defenvvar(regmachine, char *);
int branchval(regmachine, long);
int branchproc(regmachine, long);
void initprocenv(regmachine, value);
void mkcompiledprocreg(regmachine, void **, long);
void initargs(regmachine);
void addarg(regmachine);
void applyprimitivereg(regmachine, void **);
void freeregmachine(regmachine);
