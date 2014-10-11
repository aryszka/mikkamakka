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

    // noprint
    var noprint = function () { return noprint; };

    var isNoprint = function (val) {
        return val === noprint;
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

    // string
    var isString = function (string) {
        return typeof string === "string";
    };

    var stringLength = function (string) {
        return string.length;
    };

    var stringCopyCheck = function (string, from, to) {
        return (string === undefined || isString(string)) &&
            (from === undefined || isNumber(from)) &&
            (to === undefined || isNumber(to));
    };

    var stringCopy = function (string, from, to) {
        return string.substring(from, to);
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
    var callCheck = function (regs) {
        return isPrimitiveProcedure(regs.proc) ||
            isCompiledProcedure(regs.proc);
    };

    var callOp = function (regs, cont) {
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
        return error("unbound variable: " + name);
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

    var exportPair = function (p) {
        var list = [];
        for (;;) {
            if (isNull(p)) {
                return list;
            }

            if (!isPair(p)) {
                return error("argument error");
            }

            list.push(exportVal(car(p)));
            p = cdr(p);
        }
    };

    var exportCompiledProcedure = function (p) {
        var f = function () {
            var regsSave = regs;
            regs.args = importArray(Array.prototype.slice.call(arguments));
            regs.proc = p;
            regs.next = compiledEntry(p);
            regs.cont = false;
            control();
            var val = regs.val;
            regs = regsSave;
            regs.cont = false;
            return exportVal(val);
        };
        f.__exported = p;
        return f;
    };

    var exportPrimitiveProcedure = function (p) {
        return function () {
            return p.primitive.call(this, importArray(Array.prototype.slice.call(arguments)));
        };
    };

    var exportVal = function (val) {
        switch (true) {
        case isBoolean(val):
        case isNumber(val):
        case isString(val):
            return val;
        case isCompiledProcedure(val):
            return exportCompiledProcedure(val);
        case isPrimitiveProcedure(val):
            return exportPrimitiveProcedure(val);
        case isNull(val):
            return [];
        case isPair(val):
            return exportPair(val);
        case isNoprint(val):
            return undefined;
        default:
            return error("invalid type to export");
        }
    };

    var listToArray = function (list) {
        var plist = [];
        for (;;) {
            if (isNull(list)) {
                return plist;
            }

            if (!isPair(list)) {
                return error("argument error");
            }

            plist.push(car(list));
            list = cdr(list);
        }
    };

    var importFunction = function (f, ctx, value, convert) {
        return importPrimitive(function (args) {
            var result = f.apply(ctx || this, (convert ? exportPair : listToArray)(args));
            if (value !== undefined) {
                return value;
            }
            return convert ? importVal(result) : result;
        });
    };

    var isExportedProcedure = function (p) {
        return isFunction(p) && p.__exported;
    };

    var importExportedProcedure = function (p) {
        return p.__exported;
    };

    var importArray = function (val) {
        var l = list();
        for (var i = val.length - 1; i >= 0; i--) {
            l = cons(importVal(val[i]), l);
        }
        return l;
    };

    var importVal = function (val, module) {
        switch (true) {
        case isBoolean(val):
        case isNumber(val):
        case isString(val):
            return val;
        case isExportedProcedure(val):
            return importExportedProcedure(val);
        case isFunction(val):
            return importFunction(val, module, undefined, true);
        case val instanceof Array:
            return importArray(val);
        case val === undefined:
            return noprint;
        default:
            return error("invalid import type");
        }
    };

    var importModule = function (name, module) {
        for (var key in module) {
            defineVar(regs.env, stringToSymbol(key), importVal(module[key], module));
        }
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

    var add = function (args) {
        var sum = 0;
        for (;;) {
            if (isNull(args)) {
                return sum;
            }

            if (!isPair(args) || !isNumber(car(args))) {
                return error("argument error");
            }

            sum += car(args);
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

            if (isFalse(first)) {
                diff = car(args);
                first = true;
            } else {
                multiple = true;
                diff -= car(args);
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

    var makeCompareNumbers = function (compareTwo) {
        return function (args) {
            var last = false;
            for (;;) {
                if (isNull(args)) {
                    return true;
                }

                if (!isPair(args) || !isNumber(car(args))) {
                    return error("argument error");
                }

                if (isFalse(last) || compareTwo(last, car(args))) {
                    last = car(args);
                    args = cdr(args);
                    continue;
                }

                return false;
            }
        };
    };

    var less = makeCompareNumbers(function (left, right) {
        return left < right;
    });

    var greaterOrEquals = makeCompareNumbers(function (left, right) {
        return left >= right;
    });

    var stringAppend = function (args) {
        var strings = [];
        for (;;) {
            if (isNull(args)) {
                return strings.join("");
            }

            if (!isPair(args) || !isString(car(args))) {
                return error("argument error");
            }

            strings.push(car(args));
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
        call: func(2, false, callCheck, callOp)
    };

    // registers
    var regs = {
        env: makeEnv(false),
        proc: false,
        val: false,
        args: [],
        cont: false,
        next: false
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
        return ops.call(regs, false);
    };

    var call = function () {
        if (isNull(regs.args)) {
            return error("invalid arity");
        }

        regs.args = [car(regs.args), [cdr(regs.args), []]];
        return apply;
    };

    var callCc = function () {
        if (isNull(regs.args) ||
            !isNull(cdr(regs.args))) {
            return error("invalid arity");
        }

        var regsSave = {
            env: regs.env,
            proc: regs.proc,
            val: regs.val,
            args: regs.args,
            cont: regs.cont,
            next: regs.next
        };
        var stackSave = stack.slice();

        regs.proc = car(regs.args);
        regs.args = ops.list(ops.makeProcedure(function () {
            if (isNull(regs.args) ||
                !isNull(cdr(regs.args))) {
                return error("invalid arity");
            }

            regs.val = car(regs.args);

            regs.env = regsSave.env;
            regs.proc = regsSave.proc;
            regs.args = regsSave.args;
            regs.cont = regsSave.cont;
            regs.next = regsSave.next;

            stack = stackSave.slice();

            return regs.cont;
        }, regs.env));

        return ops.call(regs, false);
    };

    var breakExecution = function () {
        if (!isNull(regs.args)) {
            return error("invalid arity");
        }

        return false;
    };

    // primitive definitions
    defineVar(regs.env, stringToSymbol("noprint"), noprint);

    defineVar(regs.env, stringToSymbol("true"), true);

    defineVar(regs.env, stringToSymbol("false"), false);

    defineVar(regs.env, stringToSymbol("apply"),
        makeProcedure(apply, regs.env));

    defineVar(regs.env, stringToSymbol("call/cc"),
        makeProcedure(callCc, regs.env));

    defineVar(regs.env, stringToSymbol("break-execution"),
        makeProcedure(breakExecution, regs.env));

    defineVar(regs.env, stringToSymbol("call"),
        makeProcedure(call, regs.env));

    defineVar(regs.env, stringToSymbol("="),
        importPrimitive(numberEqual));

    defineVar(regs.env, stringToSymbol("+"),
        importPrimitive(add));

    defineVar(regs.env, stringToSymbol("-"),
        importPrimitive(subtract));

    defineVar(regs.env, stringToSymbol("*"),
        importPrimitive(multiply));

    defineVar(regs.env, stringToSymbol("<"),
        importPrimitive(less));

    defineVar(regs.env, stringToSymbol(">="),
        importPrimitive(greaterOrEquals));

    defineVar(regs.env, stringToSymbol("out"),
        importFunction(console.log, console, noprint));

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

    defineVar(regs.env, stringToSymbol("primitive-eq?"),
        importFunction(func(2, false, false, primitiveEq)));

    defineVar(regs.env, stringToSymbol("pair?"),
        importFunction(func(1, false, false, isPair)));

    defineVar(regs.env, stringToSymbol("cons"),
        importFunction(func(2, false, false, cons)));

    defineVar(regs.env, stringToSymbol("car"),
        importFunction(func(1, false, isPair, car)));

    defineVar(regs.env, stringToSymbol("cdr"),
        importFunction(func(1, false, isPair, cdr)));

    defineVar(regs.env, stringToSymbol("string?"),
        importFunction(func(1, false, false, isString)));

    defineVar(regs.env, stringToSymbol("string-length"),
        importFunction(func(1, false, isString, stringLength)));

    defineVar(regs.env, stringToSymbol("string-copy"),
        importFunction(func(0, true, stringCopyCheck, stringCopy)));

    defineVar(regs.env, stringToSymbol("string-append"),
        importPrimitive(stringAppend));

    // program

    // control
    regs.next = start;
    var control = function () {
        for (; regs.next; regs.next = regs.next());
    };
    control();
})();