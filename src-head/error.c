#include <stdlib.h>
#include <stdio.h>
#include "error.h"

char *snone = "";
char *sinvalidnumber = "invalid number";
char *sinvalidregex = "invalid regex";
char *sinvalidtype = "invalid type";
char *sfileopenfailed = "file open failed";
char *sfileclosefailed = "file close failed";
char *sfileseekfailed = "file seek failed";
char *sfilereadfailed = "file read failed";
char *sfilewritefailed = "file write failed";
char *ssymbolnotfound = "symbol not found";
char *sinvalidnumberofargs = "invalid number of arguments";
char *snumbernotint = "number not an integer";
char *snumbertoobig = "number too big";
char *seof = "eof";
char *susererror = "";
char *sunknownerror = "unknown error";

char *errorstring(errorcode code) {
	switch (code) {
	case none:
		return snone;
	case invalidnumber:
		return sinvalidnumber;
	case invalidregex:
		return sinvalidregex;
	case invalidtype:
		return sinvalidtype;
	case fileopenfailed:
		return sfileopenfailed;
	case fileclosefailed:
		return sfileclosefailed;
	case fileseekfailed:
		return sfileseekfailed;
	case filereadfailed:
		return sfilereadfailed;
	case filewritefailed:
		return sfilewritefailed;
	case symbolnotfound:
		return ssymbolnotfound;
	case invalidnumberofargs:
		return sinvalidnumberofargs;
	case numbernotint:
		return snumbernotint;
	case numbertoobig:
		return snumbertoobig;
	case eof:
		return seof;
	case usererror:
		return susererror;
	case unknownerror:
		return sunknownerror;
	}
}

void error(errorcode code, char *info) {
	if (!code) {
		return;
	}

	if (code == usererror) {
		fprintf(stderr, "error %d: %s\n", code, info);
	} else {
		fprintf(stderr, "error %d: %s (%s)\n", code, errorstring(code), info);
	}

	exit(code);
}
