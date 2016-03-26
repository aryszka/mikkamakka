#include <stdio.h>

int main() {
    long n = 1;
    long *p = &n;
    char *c = (char *)p;
    int l = sizeof(n);
    printf("%d\n", l);
    for (int i = 0; i < l; i++) {
        printf("%d\n", *(c + i));
    }

    return 0;
}
