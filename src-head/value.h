// requires first:
// number.h
// string.h
// compound-types.h
// io.h

typedef void (*freeraw)(void *raw);
typedef char *(*sprintrawfunc)(void *raw);
typedef int (*istype)(value v);

value false;
value true;
value null;

void freenoop(void *);
int typecheckany(value v);
char *sprintrawdefault(void *);
value mksymval(int, char *);
value mknumval(number);
value mknumvalc(int, char *);
value mknumvali(long, long);
value mkstringval(string);
value mkstringvalc(int, char *);
value mkpairval(value, value);
value mkprimitiveprocval(primitive);
value mkcompiledprocval(value, environment);
value mkval(void *, freeraw, sprintrawfunc, istype);
void *rawval(value v);
istype gettypecheck(value);
value istypeval(value, istype);
int issymtype(value);
int isnumtype(value);
int isinttype(value);
int issmallinttype(value);
int isstringtype(value);
int isfalsetype(value);
int isnulltype(value);
int ispairtype(value);
int isproctype(value);
int isprimitiveproctype(value);
value issymval(value);
value isnumval(value);
value isintval(value);
value issmallintval(value);
value isstringval(value);
value isfalseval(value);
value isnullval(value);
value ispairval(value);
value isprocval(value);
value isprimitiveprocval(value);
char *valsymbolnameraw(value);
int valsymbollen(value);
number numval(value);
long valrawint(value);
string valstring(value);
char *valrawstring(value);
value carval(value);
value cdrval(value);
value valapplyprimitive(value, value);
value valproclabel(value);
value valapply(value);
environment valprocenv(value);
registryflag getregistryflag(value);
void setregistryflag(value, registryflag);
value sprint(value);
char *sprintraw(value);
void freeval(value);

void initmodule_value();
void freemodule_value();
