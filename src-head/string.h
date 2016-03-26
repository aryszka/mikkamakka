struct string;
typedef struct string *string;

string mkstring(long, char *);
string clonestring(string);
long byteslen(string);
long stringlen(string);
int isutf8(string);
long bytestochars(long, char *);
string substr(string, long, long);
string byteslice(string, long, long);
string appendstr(long, string *);
int comparestring(string, string);
char *rawstring(string);
char *sprintstring(string);
void freestring(string);
