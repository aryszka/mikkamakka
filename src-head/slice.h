struct slice;

typedef void *sliceitem;
typedef struct slice *slice;

slice *append(slice s, sliceitem i);
slice *take(slice s, long offset, long length);
void freeslice(slice s);
