typedef enum {
	none,
	invalidnumber,
	invalidregex,
	invalidtype,
	fileopenfailed,
	fileclosefailed,
	fileseekfailed,
	filereadfailed,
	filewritefailed,
	symbolnotfound,
	invalidnumberofargs,
	numbernotint,
	numbertoobig,
	eof,
	usererror,
	unknownerror
} errorcode;

char *snone;
char *sinvalidnumber;
char *sinvalidregex;
char *sinvalidtype;
char *sfileopenfailed;
char *sfileclosefailed;
char *sfileseekfailed;
char *sfilereadfailed;
char *sfilewritefailed;
char *ssymbolnotfound;
char *sinvalidnumberofargs;
char *snumbernotint;
char *snumbertoobig;
char *seof;
char *susererror;
char *sunknownerror;

char *errorstring(errorcode code);
void error(errorcode code, char *info);
