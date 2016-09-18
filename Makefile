all: release

release:
	time go run cmd/mm-next/mm.go scm/mm.scm > obj/mm-precompile.go
	time go run obj/mm-precompile.go scm/mm.scm > obj/mm-compile.go
	time go run obj/mm-compile.go scm/mm.scm > obj/mm-check.go
	diff obj/mm-{compile,check}.go
	mv obj/mm-compile.go cmd/mm-next/mm.go
