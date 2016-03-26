#include <locale.h>
#include <stdio.h>
#include <string.h>

char getdecchar() {
	struct lconv *l = localeconv();
	if (strlen(l->decimal_point)) {
		return *(l->decimal_point);
	}

	return '.';
}

void initsys() {
	setlocale(LC_ALL, "");
}
