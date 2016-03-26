#include <stdio.h>

int main() {
    int a = 1;
    int b = 2;
    printf("%ld %ld\n", &a - &b, &b - &a);
    printf("%ld %ld %ld %ld\n", sizeof(long), sizeof(long *), sizeof(int), sizeof(size_t));
}
