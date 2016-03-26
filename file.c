#include <stdio.h>
#include <errno.h>

int main() {
    FILE *f = fopen("file.c", "r");
    int err = fseek(f, -1200, SEEK_CUR);
    printf("%d %d\n", err, errno);
    return 0;
}
