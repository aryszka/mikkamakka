include config.mk

all: bin/mm

libs           = -lgmp -lunistring -lpcre2-8
compile        = ${CC} -o $@ ${libs} $^
test-benchmark = ${compile} && ./$@ && echo $@ ok 1>&2 || { echo $@ failed; false; }
test           = mkdir -p test && ${test-benchmark}
benchmark      = mkdir -p benchmark && ${test-benchmark}

test: \
	test/error \
	test/symbol \
	test/number \
	test/string \
	test/regex \
	test/io \
	test/value \
	test/pair \
	test/procedure \
	test/environment \
	test/symbol-table \
	test/stack \
	test/register-machine \
	test/sprint-list \
	test/primitives \
	test/bytemap

benchmark: \
	benchmark/symbol \
	benchmark/number \
	benchmark/string \
	benchmark/regex \
	benchmark/io \
	benchmark/value \
	benchmark/pair \
	benchmark/procedure \
	benchmark/environment \
	benchmark/symbol-table \
	benchmark/stack \
	benchmark/register-machine \
	benchmark/sprint-list \
	benchmark/primitives \
	benchmark/bytemap

clean:
	rm -f src-head/*.o
	rm -rf test benchmark
	rm -rf bin
	rm -rf .tmp

primitive-type-headers = src-head/symbol.h src-head/number.h src-head/string.h

compound-type-headers = \
	src-head/compound-types.h src-head/procedure.h \
	src-head/environment.h src-head/pair.h src-head/value.h \
	src-head/stack.h src-head/bytemap.h

value-level-headers = ${primitive-type-headers} ${compound-type-headers} src-head/sprint-list.h src-head/registry.h 

test-headers = \
	src-head/sys.h src-head/testing.h src-head/error.h \
	src-head/error.mock.h src-head/sysio.h src-head/sysio.mock.h

test-headers-error = \
	src-head/sys.h src-head/testing.h src-head/error.h \
	src-head/sysio.h src-head/sysio.mock.h

src-head/sys.o:                    src-head/sys.h
src-head/testing.o:                src-head/sys.h src-head/testing.h
src-head/error.o:                  src-head/error.h
src-head/error.mock.o:             src-head/error.h src-head/error.mock.h
src-head/error.test.o:             src-head/error.h ${test-headers}
src-head/symbol.o:                 src-head/symbol.h
src-head/symbol.test.o:            src-head/symbol.h ${test-headers}
src-head/symbol.benchmark.o:       src-head/symbol.h ${test-headers}
src-head/number.o:                 src-head/sys.h src-head/error.h src-head/number.h
src-head/number.test.o:            src-head/number.h ${test-headers}
src-head/number.benchmark.o:       src-head/number.h ${test-headers}
src-head/string.o:                 src-head/string.h
src-head/string.test.o:            src-head/string.h ${test-headers}
src-head/string.benchmark.o:       src-head/string.h ${test-headers}
src-head/value.o:                  src-head/error.h ${value-level-headers}
src-head/value.test.o:             ${value-level-headers} ${test-headers}
src-head/value.benchmark.o:        ${value-level-headers} ${test-headers}
src-head/pair.o:                   ${compound-type-headers}
src-head/pair.test.o:              ${compound-type-headers} ${test-headers}
src-head/pair.benchmark.o:         ${compound-type-headers} ${test-headers}
src-head/environment.o:            src-head/error.h src-head/symbol.h ${compound-type-headers}
src-head/environment.test.o:       ${compound-type-headers} ${test-headers}
src-head/environment.benchmark.o:  ${compound-type-headers} ${test-headers}
src-head/symbol-table.o:           src-head/symbol.h src-head/bytemap.h
src-head/symbol-table.test.o:      src-head/symbol.h src-head/bytemap.h ${test-headers}
src-head/symbol-table.benchmark.o: src-head/symbol.h src-head/bytemap.h ${test-headers}
src-head/procedure.o:              src-head/error.h ${compound-type-headers}
src-head/procedure.test.o:         ${compound-type-headers} ${test-headers}
src-head/procedure.benchmark.o:    ${compound-type-headers} ${test-headers}
src-head/sprint-list.o:            ${value-level-headers}
src-head/sprint-list.test.o:       ${value-level-headers} ${test-headers}
src-head/sprint-list.benchmark.o:  ${value-level-headers} ${test-headers}
src-head/regex.o:                  src-head/error.h src-head/regex.h
src-head/regex.test.o:             src-head/regex.h ${test-headers}
src-head/regex.benchmark.o:        src-head/regex.h ${test-headers}
src-head/sysio.o:                  src-head/sysio.h
src-head/sysio.mock.o:             src-head/sysio.h src-head/sysio.mock.h
src-head/io.o:                     src-head/error.h src-head/sysio.h src-head/io.h
src-head/io.test.o:                src-head/io.h ${test-headers}
src-head/io.benchmark.o:           src-head/io.h ${test-headers}
src-head/stack.o:                  src-head/stack.h
src-head/stack.test.o:             src-head/stack.h ${test-headers}
src-head/stack.benchmark.o:        src-head/stack.h ${test-headers}
src-head/primitives.o:             src-head/primitives.h src-head/regex.h src-head/error.h ${value-level-headers}
src-head/primitives.test.o:        src-head/primitives.h src-head/regex.h src-head/error.h ${value-level-headers} ${test-headers}
src-head/primitives.benchmark.o:   src-head/primitives.h src-head/regex.h src-head/error.h ${value-level-headers} ${test-headers}
src-head/registry.o:               src-head/registry.h ${compound-type-headers}
src-head/bytemap.o:                src-head/bytemap.h
src-head/bytemap.test.o:           src-head/bytemap.h ${test-headers}
src-head/bytemap.benchmark.o:      src-head/bytemap.h ${test-headers}

register-machine.o:
	src-head/stack.h src-head/register-machine.h src-head/primitivies.h src-head/regex.h \
	${value-level-headers}
register-machine.test.o: \
	src-head/stack.h src-head/register-machine.h src-head/primitives.h src-head/regex.h \
	${value-level-headers} \
	${test-headers-error}
register-machine.benchmark.o: \
	src-head/stack.h src-head/register-machine.h src-head/primitives.h src-head/regex.h \
	${value-level-headers} \
	${test-headers-error}

primitive-type-objects = src-head/symbol.o src-head/number.o src-head/string.o
compound-type-objects = \
	src-head/pair.o src-head/procedure.o src-head/environment.o \
	src-head/symbol-table.o src-head/value.o src-head/registry.o \
	src-head/stack.o src-head/bytemap.o

value-level-objects-base = ${primitive-type-objects} ${compound-type-objects} src-head/sprint-list.o
io-objects               = src-head/sysio.o src-head/io.o
io-objects-test          = src-head/sysio.mock.o src-head/io.o
test-objects             = src-head/sys.o src-head/testing.o src-head/error.mock.o
value-level-objects      = ${value-level-objects-base} ${io-objects}
value-level-objects-test = ${value-level-objects-base} ${io-objects-test}

register-machine-objects = \
	src-head/error.o src-head/sys.o src-head/regex.o \
	src-head/register-machine.o src-head/primitives.o \
	${value-level-objects}

register-machine-test-objects = \
	src-head/io.o src-head/regex.o \
	src-head/stack.o src-head/register-machine.o src-head/primitives.o \
	${value-level-objects-test} \
	${test-objects}

test/error: src-head/sys.o src-head/testing.o src-head/error.o src-head/error.test.o
	${test}

test/symbol: src-head/symbol.o src-head/symbol.test.o ${test-objects}
	${test}

benchmark/symbol: src-head/symbol.o src-head/symbol.benchmark.o ${test-objects}
	${benchmark}

test/number: src-head/number.o src-head/number.test.o ${test-objects}
	${test}

benchmark/number: src-head/number.o src-head/number.benchmark.o ${test-objects}
	${benchmark}

test/string: src-head/string.o src-head/string.test.o ${test-objects}
	${test}

benchmark/string: src-head/string.o src-head/string.benchmark.o ${test-objects}
	${benchmark}

test/value: src-head/value.test.o ${value-level-objects-test} ${test-objects}
	${test}

benchmark/value: src-head/value.benchmark.o ${value-level-objects-test} ${test-objects}
	${benchmark}

test/pair: src-head/pair.test.o ${value-level-objects-test} ${test-objects}
	${test}

benchmark/pair: src-head/pair.benchmark.o ${value-level-objects-test} ${test-objects}
	${benchmark}

test/environment: src-head/environment.test.o ${value-level-objects-test} ${test-objects}
	${test}

benchmark/environment: src-head/environment.benchmark.o ${value-level-objects-test} ${test-objects}
	${benchmark}

test/symbol-table: src-head/symbol-table.o src-head/bytemap.o src-head/symbol-table.test.o ${test-objects}
	${test}

benchmark/symbol-table: src-head/symbol-table.o src-head/bytemap.o src-head/symbol-table.benchmark.o ${test-objects}
	${benchmark}

test/procedure: src-head/procedure.test.o ${value-level-objects-test} ${test-objects}
	${test}

benchmark/procedure: src-head/procedure.benchmark.o ${value-level-objects-test} ${test-objects}
	${benchmark}

test/sprint-list: src-head/sprint-list.test.o ${value-level-objects-test} ${test-objects}
	${test}

benchmark/sprint-list: src-head/sprint-list.benchmark.o ${value-level-objects-test} ${test-objects}
	${benchmark}

test/regex: src-head/regex.o src-head/regex.test.o ${test-objects}
	${test}

benchmark/regex: src-head/regex.o src-head/regex.benchmark.o ${test-objects}
	${benchmark}

test/io: src-head/io.o src-head/io.test.o ${test-objects} ${io-objects-test}
	${test}

benchmark/io: src-head/io.o src-head/io.benchmark.o ${test-objects} ${io-objects-test}
	${benchmark}

test/stack: src-head/stack.o src-head/stack.test.o ${test-objects}
	${test}

benchmark/stack: src-head/stack.o src-head/stack.benchmark.o ${test-objects}
	${benchmark}

test/register-machine: src-head/register-machine.test.o ${register-machine-test-objects}
	${test}

benchmark/register-machine: src-head/register-machine.benchmark.o ${register-machine-test-objects}
	${benchmark}

test/primitives: src-head/primitives.test.o src-head/primitives.o src-head/regex.o ${value-level-objects-test} ${test-objects}
	${test}

benchmark/primitives: src-head/primitives.benchmark.o src-head/primitives.o src-head/regex.o ${value-level-objects-test} ${test-objects}
	${benchmark}

test/bytemap: src-head/bytemap.o src-head/bytemap.test.o ${test-objects}
	${test}

benchmark/bytemap: src-head/bytemap.o src-head/bytemap.benchmark.o ${test-objects}
	${benchmark}

.tmp/mm-out.c: src-mm/mm-boot.scm
	mkdir -p .tmp && guile src-mm/mm-boot.scm > .tmp/mm-out.c

.tmp/mm-out.o: .tmp/mm-out.c
	${CC} -c -o .tmp/mm-out.o $^

bin/mm: .tmp/mm-out.o ${register-machine-objects}
	mkdir -p bin && ${CC} -o bin/mm ${libs} .tmp/mm-out.o ${register-machine-objects}

.tmp/output-test.c: bin/mm src-mm/mm.scm
	mkdir -p .tmp && bin/mm > .tmp/output-test.c

.tmp/output-test.o: .tmp/output-test.c
	${CC} -c -o .tmp/output-test.o .tmp/output-test.c

bin/output-test: .tmp/output-test.o ${register-machine-objects}
	mkdir -p bin && ${CC} -o bin/output-test ${libs} .tmp/output-test.o ${register-machine-objects}
