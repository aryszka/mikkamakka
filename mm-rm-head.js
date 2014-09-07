(function () {
    var identity = function (x) { return x; };

    var error = function (msg) {
        throw new Error(msg);
    };

    var func = function (argLength, allowMore, customCheck, f) {
        return function () {
            var args = Array.prototype.slice.call(arguments);

            if (allowMore && args.length < argLength ||
                !allowMore && args.length !== argLength) {
                return error("invalid arity");
            }

            if (customCheck && !customCheck.apply(undefined, args)) {
                return error("argument error");
            }

            return f.apply(undefined, args);
        };
    };

    var isFunction = function (f) {
        return typeof f === "function";
    };

    var primitiveEq = function (left, right) {
        return left === right;
    };

    // boolean
    var isBoolean = function (value) {
        return value === true || value === false;
    };

    var isFalse = function (value) {
        return value === false;
    };

    // number
    var isNumber = function (number) {
        return typeof number === "number";
    };

    // symbol
    var isSymbol = function (symbol) {
        return symbol instanceof Array &&
            symbol.length === 1 &&
            typeof symbol[0] === "string";
    };

    var stringToSymbol = function (string) {
        return [string];
    };

    var symbolToString = function (symbol) {
        return symbol[0];
    };

    // null
    var isNull = function (value) {
        return value instanceof Array && value.length === 0;
    };

    // pair
    var isPair = function (pair) {
        return pair instanceof Array && pair.length === 2;
    };

    var cons = function (car, cdr) {
        return [car, cdr];
    };

    var car = function (pair) {
        return pair[0];
    };

    var cdr = function (pair) {
        return pair[1];
    };

    // list
    var list = function () {
        var l = [];
        for (var i = 0; i < arguments.length; i++) {
            l = cons(arguments[i], l);
        }
        return l;
    };

    // primitiveProcedure
    var isPrimitiveProcedure = function (p) {
        return !!(p && typeof p.primitive === "function");
    };

    var applyPrimitive = function (p, args) {
        return p.primitive(args);
    };

    // compiled procedure
    var isCompiledProcedure = function (p) {
        return !!(p && p.entry && p.env);
    };

    var makeProcedure = function (entry, env) {
        return {entry: entry, env: env};
    };

    var compiledEntry = function (p) {
        return p.entry;
    };

    var compiledProcedureEnv = function (p) {
        return p.env;
    };

    // call
    var call = function (regs, cont) {
        if (isPrimitiveProcedure(regs.proc)) {
            regs.val = applyPrimitive(regs.proc, regs.args);
            return cont || regs.cont;
        };

        if (cont) {
            regs.cont = cont;
        }
        return compiledEntry(regs.proc);
    };

    // struct
    var makeStruct = function () {
        var s = {struct: {}};
        for (var i = 0; i < arguments.length; i++) {
            var def = arguments[i];
            if (!isPair(def) ||
                !isSymbol(car(def)) ||
                !isPair(cdr(def))) {
                return error("invalid struct definition");
            }
            s.struct[symbolToString(car(def))] = car(cdr(def));
        }
        return s;
    };

    var isStruct = function (s) {
        return !!(s && s.struct);
    };

    var structDefined = function (s, name) {
        return symbolToString(name) in s.struct;
    };

    var structDefine = function (s, name, value) {
        s.struct[symbolToString(name)] = value;
        return value;
    };

    var structLookup = function (s, name) {
        if (!structDefined(s, name)) {
            return error("undefined struct field");
        }
        return s.struct[symbolToString(name)];
    };

    var structSet = function (s, name, value) {
        if (!structDefined(s, name)) {
            return error("undefined struct field");
        }
        s.struct[symbolToString(name)] = value;
        return value;
    };

    var structCheck = function (s, name) {
        return isStruct(s) && isSymbol(name);
    };

    // env
    var makeEnv = function (parent) {
        return makeStruct(
            [stringToSymbol("env"), [true, []]],
            [stringToSymbol("parent"), [parent, []]]);
    };

    var isEnv = function (env) {
        return isStruct(env) &&
            structDefined(env, stringToSymbol("env")),
            structDefined(env, stringToSymbol("parent"));
    };

    var withBoundVar = function (env, name, mutate) {
        for (; env; env = structLookup(env, stringToSymbol("parent"))) {
            if (structDefined(env, name)) {
                structSet(env, name, mutate(structLookup(env, name)));
                return structLookup(env, name);
            }
        }
        return error("unbound variable");
    };

    var lookupVar = function (env, name) {
        return withBoundVar(env, name, identity);
    };

    var setVar = function (env, name, value) {
        return withBoundVar(env, name, function () {
            return value;
        });
    };

    var defineVar = function (env, name, value) {
        structDefine(env, name, value);
    };

    var extendEnv = function (env, names, values) {
        env = makeEnv(env);
        for (;;) {
            if (isNull(names) && !isNull(values) ||
                !isNull(names) && !isSymbol(names) && isNull(values)) {
                return error("invalid arity");
            }

            if (isNull(names) && isNull(values)) {
                return env;
            }

            if (isSymbol(names)) {
                defineVar(env, names, values);
                return env;
            }

            if (!isPair(names) ||
                !isPair(values)) {
                return error("argument error");
            }

            defineVar(env, car(names), car(values));
            names = cdr(names);
            values = cdr(values);
        }
    };

    // import
    var importPrimitive = function (p) {
        return {primitive: p};
    };

    var argsToArray = function (args) {
        var array = [];
        for (;;) {
            if (isNull(args)) {
                return array;
            }

            if (!isPair(args)) {
                return error("argument error");
            }

            array.push(car(args));
            args = cdr(args);
        }
    };

    var importFunction = function (f, ctx, value) {
        return importPrimitive(function (args) {
            var result = f.apply(ctx || this, argsToArray(args));
            return value === undefined ? result : value;
        });
    };

    // primitive procedures
    var numberEqual = function (args) {
        var first = false;
        var multiple = false;
        for (;;) {
            if (isNull(args)) {
                if (isFalse(first) || !multiple) {
                    return error("arity error");
                }
                return true;
            }

            if (!isPair(args) || !isNumber(car(args))) {
                return error("argument error");
            }

            if (isFalse(first)) {
                first = car(args);
            } else {
                multiple = true;
                if (car(args) !== first) {
                    return false;
                }
            }

            args = cdr(args);
        }
    };

    var multiply = function (args) {
        var product = 1;
        for (;;) {
            if (isNull(args)) {
                return product;
            }

            if (!isPair(args) || !isNumber(car(args))) {
                return error("argument error");
            }

            product *= car(args);
            args = cdr(args);
        }
    };

    var subtract = function (args) {
        var diff = 0;
        var first = false;
        var multiple = false;
        for (;;) {
            if (isNull(args)) {
                if (!first) {
                    return error("arity error");
                }
                
                if (!multiple) {
                    return 0 - diff;
                }

                return diff;
            }

            if (!isPair(args) || !isNumber(car(args))) {
                return error("argument error");
            }

            if (!first) {
                diff = car(args);
                first = true;
            } else {
                multiple = true;
                diff -= car(args);
            }

            args = cdr(args);
        }
    };

    // ops
    var envCheck = function (env, name) {
        return isEnv(env) && isSymbol(name);
    };

    var ops = {
        lookupVar: func(2, false, envCheck, lookupVar),
        setVar: func(3, false, envCheck, setVar),
        defineVar: func(3, false, envCheck, defineVar),
        extendEnv: func(3, false, isEnv, extendEnv),
        isFalse: func(1, false, false, isFalse),
        makeProcedure: func(2, false, function (entry, env) {
            return isFunction(entry) && isEnv(env);
        }, makeProcedure),
        compiledEntry: func(1, false, isCompiledProcedure, compiledEntry),
        compiledProcedureEnv: func(1, false, isCompiledProcedure, compiledProcedureEnv),
        cons: func(2, false, false, cons),
        list: func(0, true, false, list),
        call: func(2, false, false, call)
    };

    // registers
    var regs = {
        env: makeEnv(false),
        proc: false,
        val: false,
        args: [],
        next: false,
        cont: false
    };

    var stack = [];

    var save = function (reg) {
        return stack.push(reg);
    };

    var restore = function () {
        return stack.pop();
    };

    var apply = function () {
        if (isNull(regs.args) ||
            isNull(cdr(regs.args)) ||
            !isNull(cdr(cdr(regs.args)))) {
            return error("invalid arity");
        }

        regs.proc = car(regs.args);
        regs.args = car(cdr(regs.args));
        return ops.call(regs, regs.cont);
    };

    // primitive definitions
    defineVar(regs.env, stringToSymbol("true"), true);

    defineVar(regs.env, stringToSymbol("false"), false);

    defineVar(regs.env, stringToSymbol("apply"),
        makeProcedure(apply, regs.env));

    defineVar(regs.env, stringToSymbol("="),
        importPrimitive(numberEqual));

    defineVar(regs.env, stringToSymbol("*"),
        importPrimitive(multiply));

    defineVar(regs.env, stringToSymbol("-"),
        importPrimitive(subtract));

    defineVar(regs.env, stringToSymbol("out"),
        importFunction(console.log, console, stringToSymbol("ok")));

    defineVar(regs.env, stringToSymbol("struct"),
        importFunction(func(0, true, false, makeStruct)));

    defineVar(regs.env, stringToSymbol("struct?"),
        importFunction(func(1, false, false, isStruct)));

    defineVar(regs.env, stringToSymbol("struct-defined?"),
        importFunction(func(2, false, structCheck, structDefined)));

    defineVar(regs.env, stringToSymbol("struct-define"),
        importFunction(func(3, false, structCheck, structDefine)));

    defineVar(regs.env, stringToSymbol("struct-lookup"),
        importFunction(func(2, false, structCheck, structLookup)));

    defineVar(regs.env, stringToSymbol("struct-set!"),
        importFunction(func(3, false, structCheck, structSet)));

    defineVar(regs.env, stringToSymbol("null?"),
        importFunction(func(1, false, false, isNull)));

    defineVar(regs.env, stringToSymbol("symbol?"),
        importFunction(func(1, false, false, isSymbol)));

    defineVar(regs.env, stringToSymbol("symbol->string"),
        importFunction(func(1, false, isSymbol, symbolToString)));

    defineVar(regs.env, stringToSymbol("peq?"),
        importFunction(func(2, false, false, primitiveEq)));

    defineVar(regs.env, stringToSymbol("cons"),
        importFunction(func(2, false, false, cons)));

    defineVar(regs.env, stringToSymbol("car"),
        importFunction(func(1, false, isPair, car)));

    defineVar(regs.env, stringToSymbol("cdr"),
        importFunction(func(1, false, isPair, cdr)));

    // program

    // control
    regs.next = start;
    var control = function () {
        for (; regs.next; regs.next = regs.next());
    };
    control();
})();
