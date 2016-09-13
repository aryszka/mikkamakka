all: test

mmfc:
	go install ./cmd/mmfc

obj:
	mkdir -p obj

eval-compile: obj
	time mmfc scm/mm.scm < scm/mm.scm > obj/mm.go

compile:
	time go run obj/mm.go scm/mm.scm > obj/mm-out.go

test: mmfc obj eval-compile compile
	diff obj/mm{,-out}.go
