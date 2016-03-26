typedef enum {
	sstart,
	scurr,
	send
} sseekmode;

struct sfile;
typedef struct sfile *sfile;

int serrno;
sfile sstdin;
sfile sstdout;
sfile sstderr;

sfile sfopen(char *n, char *mode);
int sfclose(sfile sf);
int sfseek(sfile sf, int offset, sseekmode mode);
int sfread(char *c, int len, sfile sf);
int sfeof(sfile sf);
int sfwrite(char *c, int len, sfile sf);

void initmodule_sysio();
void freemodule_sysio();
