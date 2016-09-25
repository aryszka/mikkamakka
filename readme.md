# Mikkamakka Functional Calculator

MMFC is an extendable and programmable software calculator.

### Work In Progress

The current pre-release version is only capable to compile or interpret itself. Actual calculator functionality
is going to be added as the programming platform reaches a more stable state.

### Bootstrapping

##### Prerequisits

- the Go platform [https://golang.org](https://golang.org)
- make (preferrably GNU Make)

##### Generate and install initial version from source:

```
make bootstrap
```

### Running MMFC

MMFC can compile or run a custom dialect of Scheme LISP. Compilation means only transpilation to Go and running
or installing the output requires the Go build environment. Interpreting MMFC code is possible without the Go
environment:

```
mmfc my-code.scm
```

or:

```
mmfc run my-coce.scm
```

The interactive REPL mode can be started by simply:

```
mmfc
```

...and, e.g, one can run:

```
(+ 19 23)
```

(outputs: 42).
