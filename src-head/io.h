// requires first:
// error.h

typedef enum {
	ioread,
	ioreadwrite,
	iowrite,
	ioreadwritecreate,
	ioappend,
	ioreadappend
} iofilemode;

typedef enum {
	iostart,
	iocurr,
	ioend
} ioseekmode;

typedef struct {
	errorcode code;
	int syserror;
} ioerror;

struct file;
typedef struct file *file;

file iostdin;
file iostdout;
file iostderr;

void openfile(char *n, iofilemode mode, file *f, ioerror *err);
ioerror seekfile(file f, long pos, ioseekmode);
void readfile(file f, long len, char *s, long *rlen, ioerror *err);
ioerror writefile(file f, long len, char *s);
ioerror closefile(file f);
void initmodule_io();
void freemodule_io();
