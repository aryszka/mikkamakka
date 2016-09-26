.PHONY: bootstrap

all: build

obj/mmfc:
	mkdir -p obj/mmfc

obj/mmfc/mm.go: obj/mmfc scm/*.scm
	time mmfc compile < scm/mm.scm > obj/mmfc/mm.go
	time go run obj/mmfc/mm.go compile < scm/mm.scm > obj/mm-check.go
	diff obj/mmfc/mm.go obj/mm-check.go

build: obj/mmfc/mm.go

install: build
	go install ./obj/mmfc

bootstrap/obj:
	mkdir -p bootstrap/obj

bootstrap: obj/mmfc
	mkdir -p bootstrap/obj
	go run bootstrap/bootstrap.go > bootstrap/obj/mm.go
	go run bootstrap/obj/mm.go compile scm/mm.scm > bootstrap/obj/mm-check.go
	diff bootstrap/obj/mm{,-check}.go
	mv bootstrap/obj/mm.go obj/mmfc/mm.go
	go install ./obj/mmfc

gen-bootstrap:
	mmfc compile bootstrap/bootstrap.scm > bootstrap/bootstrap.go

release: build gen-bootstrap
