enum mockerrno {
	mnopenfailed,
	mnclosefailed,
	mnseekfailed,
	mnreadfailed,
	mnwritefailed
};

char *fnok;
char *fnopenfail;
char *fnclosefail;
char *fnseekfail;
char *fnreadfail;
char *fnwritefail;

int testfilelength;

void initfilecontent(char *fn, char *content);
void cleartestcontent();
