#include <stdlib.h>
#include <gmp.h>
#include <limits.h>
#include <string.h>
#include "sys.h"
#include "error.h"
#include "number.h"

struct number {
	mpq_t val;
};

static number minint;
static number maxint;

static number mknum() {
	number n = malloc(sizeof(struct number));
	mpq_init(n->val);
	return n;
}

static int setmpqc(mpq_t q, size_t len, char *c) {
	if (len == 0) {
		return 1;
	}

	char dot = getdecchar();
	char *cc = malloc(len + 1);

	size_t pos = 0;
	size_t shift = 0;
	size_t dotpos = len;
	char curr;
	for (;;) {
		curr = *(c + pos + shift);

		if (curr == dot) {
			if (shift) {
				free(cc);
				return 1;
			}

			dotpos = pos;
			shift = 1;
			curr = *(c + pos + shift);
		}

		*(cc + pos) = curr;
		if (curr == 0) {
			break;
		}

		pos++;
	}

	if (len == 1 && dotpos == 0) {
		free(cc);
		return 1;
	}

	mpz_t z;
	mpz_init(z);
	int check = mpz_set_str(z, cc, 10);
	free(cc);
	if (check) {
		mpz_clear(z);
		return 1;
	}

	mpq_set_num(q, z);

	if (shift) {
		size_t denlen = len - dotpos;
		char *sd = malloc(denlen + 1);

		*sd = '1';
		for (pos = 1; pos < denlen; pos++) {
			*(sd + pos) = '0';
		}
		*(sd + denlen) = 0;

		check = mpz_set_str(z, sd, 10);
		free(sd);
		if (check) {
			mpz_clear(z);
			return 1;
		}

		mpq_set_den(q, z);
	}

	mpz_clear(z);
	mpq_canonicalize(q);

	return 0;
}

static void mpqseti(mpq_t q, long n, long d) {
	mpq_set_si(q, n, d);
	mpq_canonicalize(q);
}

static int less(number n1, number n2) {
	return comparenum(n1, n2) < 0;
}

static int greater(number n1, number n2) {
	return comparenum(n1, n2) > 0;
}

static int checkint(number n, int checksmall) {
	mpz_t num;
	mpz_t den;

	mpz_init(den);
	mpq_get_den(den, n->val);
	long di = mpz_get_si(den);
	mpz_clear(den);

	if (di != 1) {
		return 0;
	}

	if (!checksmall) {
		return 1;
	}

	return greater(n, minint) && less(n, maxint);
}

number mknumcsafe(size_t len, char *c) {
	number n = mknum();
	if (!setmpqc(n->val, len, c)) {
		return n;
	}

	mpq_clear(n->val);
	free(n);
	return 0;
}

number mknumc(size_t len, char *c) {
	number n = mknumcsafe(len, c);
	if (!n) {
		error(invalidnumber, "invalid number literal");
	}

	return n;
}

number mknumi(long nom, long den) {
	number n = mknum();
	if (den < 0) {
		den = 0 - den;
		nom = 0 - nom;
	}

	mpqseti(n->val, nom, den);
	return n;
}

number clonenum(number n) {
	if (issmallint(n)) {
		return mknumi(rawint(n), 1);
	}

	char *s = sprintnum(n);
	number cn = mknumc(strlen(s), s);
	free(s);
	return cn;
}

number sum(number n1, number n2) {
	number s = mknum();
	mpq_add(s->val, n1->val, n2->val);
	return s;
}

number diff(number n1, number n2) {
	number s = mknum();
	mpq_sub(s->val, n1->val, n2->val);
	return s;
}

number bitor(number n1, number n2) {
	if (!isint(n1) || !isint(n2)) {
		error(invalidnumber, "integer expected");
		return 0;
	}

	mpz_t ni1;
	mpz_init(ni1);
	mpq_get_num(ni1, n1->val);

	mpz_t ni2;
	mpz_init(ni2);
	mpq_get_num(ni2, n2->val);

	mpz_t ior;
	mpz_init(ior);
	mpz_ior(ior, ni1, ni2);

	number r = mknum();
	mpq_set_num(r->val, ior);

	mpz_clear(ni1);
	mpz_clear(ni2);
	mpz_clear(ior);

	return r;
}

long rawint(number n) {
	if (!isint(n)) {
		error(numbernotint, "integer expected");
		return 0;
	}

	char *toobigmsg = "number (abs) too big for integer conversion";

	if (less(n, minint)) {
		error(numbertoobig, toobigmsg);
		return 0;
	}

	if (greater(n, maxint)) {
		error(numbertoobig, toobigmsg);
		return 0;
	}

	mpz_t num;
	mpz_init(num);
	mpq_get_num(num, n->val);
	long ri = mpz_get_si(num);
	mpz_clear(num);
	return ri;
}

int isint(number n) {
	return checkint(n, 0);
}

int issmallint(number n) {
	return checkint(n, 1);
}

int comparenum(number n1, number n2) {
	return mpq_cmp(n1->val, n2->val);
}

char *sprintnum(number n) {
	mpf_t f;
	mpf_init(f);
	mpf_set_prec(f, 2048);
	mpf_set_q(f, n->val);
	char *fmt = "%.64Fg";
	char *s;
	size_t len = gmp_snprintf(s, 0, fmt, f) + 1;
	s = malloc(len);
	gmp_snprintf(s, len, fmt, f);
	mpf_clear(f);
	return s;
}

void freenum(number n) {
	mpq_clear(n->val);
	free(n);
}

void initmodule_number() {
	minint = mknumi(LONG_MIN, 1);
	maxint = mknumi(LONG_MAX, 1);
}

void freemodule_number() {
	freenum(minint);
	freenum(maxint);
}
