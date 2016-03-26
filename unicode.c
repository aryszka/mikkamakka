#include <stdlib.h>
#include <locale.h>
#include <stdio.h>
#include <string.h>
#include <unistr.h>

int utf8bytes(char *s) {
    int len = strlen(s);
    const uint8_t *c = u8_check((uint8_t *)s, len);
    if (!c) {
        return len;
    }

    return (char *)c - s;
}

int main() {
	setlocale(LC_ALL, "");

    printf("%d - %ld\n", utf8bytes("abc"), u8_mbsnlen((uint8_t *)"abc", 3));

    char *s = "fűzfánfütyülő";
    int len = strlen(s);
    printf("%d - %d - %ld\n", len, utf8bytes(s), u8_mbsnlen((uint8_t *)s, (size_t)len));

    char *base = "fűzfánfütyülő";
    int baselen = strlen(base);
    int brokenlen = baselen + 1;
    char *broken = malloc(brokenlen + 1);
    for (int i = 0; i < baselen; i++) {
        *(broken + i) = *(base + i);
    }
    *(broken + brokenlen - 1) = 255;
    *(broken + brokenlen) = 0;
	const uint8_t *check = u8_check((uint8_t *)broken, brokenlen);
    if (check == 0) {
        fprintf(stderr, "failed to create non-utf8 string\n");
        return 1;
    }
    printf("%d - %d - %d - %ld\n", baselen, brokenlen, utf8bytes(broken), u8_mbsnlen((uint8_t *)broken, (size_t)brokenlen));

    base = "fűzfánfütyülő";
    baselen = strlen(base);
    broken = malloc(baselen + 1);
    for (int i = 0; i < baselen; i++) {
        if (i == 6) {
            *(broken + i) = 255;
        } else {
            *(broken + i) = *(base + i);
        }
    }
	check = u8_check((uint8_t *)broken, baselen);
    if (check == 0) {
        fprintf(stderr, "failed to create non-utf8 string\n");
        return 1;
    }
    printf("%d - %d - %ld\n", baselen, utf8bytes(broken), u8_mbsnlen((uint8_t *)broken, (size_t)baselen));

    printf("u8len: %ld\n", u8_strlen((uint8_t *)base));

    return 0;
}
