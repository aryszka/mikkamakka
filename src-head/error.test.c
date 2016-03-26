#include <string.h>
#include "sys.h"
#include "testing.h"
#include "error.h"

void test_all_codes_have_string() {
	for (int code = invalidnumber; code <= unknownerror; code++) {
		if (code == usererror) {
			continue;
		}

		char *s = errorstring(code);
		assert(!!s, "error: string not null");
		assert(strcoll(s, ""), "error: string not empty");
	}
}

int main() {
	initsys();
	test_all_codes_have_string();
}
