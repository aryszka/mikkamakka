.PHONY: bootstrap

all: build

obj/mmfc:
	mkdir -p obj/mmfc

obj/mmfc/mm.go: obj/mmfc scm/*.scm
	time mmfc compile < scm/mm.scm > obj/mmfc/mm.go

obj/mmfc/mm-recompile.go: obj/mmfc scm/*.scm
	time go run obj/mmfc/mm.go compile < scm/mm.scm > obj/mmfc/mm-recompile.go

build: obj/mmfc/mm.go
	time go run obj/mmfc/mm.go compile < scm/mm.scm > obj/mm-check.go
	diff obj/mmfc/mm.go obj/mm-check.go

rebuild: obj/mmfc/mm-recompile.go
	time go run obj/mmfc/mm-recompile.go compile < scm/mm.scm > obj/mm-check.go
	diff obj/mmfc/mm-recompile.go obj/mm-check.go

install: build
	go install ./obj/mmfc

reinstall: rebuild
	mv obj/mmfc/mm{-recompile,}.go
	go install ./obj/mmfc

bootstrap/obj:
	mkdir -p bootstrap/obj

bootstrap: obj/mmfc
	mkdir -p bootstrap/obj
	go run bootstrap/bootstrap.go > bootstrap/obj/mm.go
	go run bootstrap/obj/mm.go compile < scm/mm.scm > bootstrap/obj/mm-check.go
	diff bootstrap/obj/mm{,-check}.go
	mv bootstrap/obj/mm.go obj/mmfc/mm.go
	go install ./obj/mmfc

gen-bootstrap:
	mmfc compile < bootstrap/bootstrap.scm > bootstrap/bootstrap.go

release: build gen-bootstrap

head:
	go install

test: test-compile

test-repl: obj/mmfc/mm.go
	go run obj/mmfc/mm.go

test-compile: obj/mmfc/mm.go
	go run obj/mmfc/mm.go compile < test/simple.scm > obj/simple.go
	go run obj/simple.go

test-run: obj/mmfc/mm.go
	go run obj/mmfc/mm.go < test/simple.scm

test-compile-modules: obj/mmfc/mm.go
	go run obj/mmfc/mm.go compile < test/module-4.scm > obj/module-4.go
	go run obj/module-4.go

test-run-modules: obj/mmfc/mm.go
	go run obj/mmfc/mm.go < test/module-4.scm

clean:
	rm -rf obj
