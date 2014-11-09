(function () {
    var identity = function (x) { return x; };

    var func = function (argLength, allowMore, customCheck, f) {
        return function () {
            var args = Array.prototype.slice.call(arguments);

            if (allowMore && args.length < argLength ||
                !allowMore && args.length !== argLength) {
                return error("invalid arity");
            }

            if (customCheck && !customCheck.apply(this, args)) {
                return error("argument error");
            }

            return f.apply(this, args);
        };
    };

    var isFunction = function (f) {
        return typeof f === "function";
    };

    var primitiveEq = function (left, right) {
        return left === right;
    };

    // error
    var error = function (msg) {
        throw new Error(msg);
    };

    var isError = function (object) {
        return object instanceof Error;
    };

    var errorToString = function (error) {
        return error.message;
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

    var numberToString = function (number) {
        return String(number);
    };

    var stringToNumber = function (string) {
        var num = Number(string);
        return Number.isNaN(num) ? false : num;
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
        return string ? string.substring(from, to) : "";
    };

    var stringIndex = function (string, expression) {
        var m = string.match(new RegExp(expression));
        if (!m) {
            return -1;
        }
        return m.index;
    };

    var stringIndexCheck = function (string, expression) {
        return isString(string) && isString(expression);
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
        return {entry: entry, env: env, id: entry};
    };

    var compiledEntry = function (p) {
        return p.entry;
    };

    var compiledProcedureEnv = function (p) {
        return p.env;
    };

    // vector
    var vector = function (args) {
        var v = {vector: [], slice: {from: 0, to: 0}};
        for (;;) {
            if (isNull(args)) {
                return v;
            }

            v.vector.push(car(args));
            v.slice.to++;
            args = cdr(args);
        }
    };

    var isVector = function (v) {
        return !!(v && v.vector);
    };

    var vectorRef = function (v, r) {
        return v.vector[v.slice.from + r];
    };

    var vectorRefCheck = function (v, r) {
        return isVector(v) &&
            isNumber(r) &&
            r >= 0 &&
            (v.slice.from + r) < v.slice.to;
    };

    var vectorLength = function (v) {
        return v.slice.to - v.slice.from;
    };

    var vectorSlice = function (v, from, to) {
        from = from || 0;
        to = to || vectorLength(v);
        return {
            vector: v.vector,
            slice: {
                from: v.slice.from + from,
                to: v.slice.from + to
            }
        };
    };

    var vectorSliceCheck = function (v, from, to) {
        return isVector(v) &&
            (arguments.length < 2 ||
            isNumber(from) &&
            from >= 0 &&
            from < (v.slice.to - v.slice.from)) &&
            (arguments.length < 3 ||
            isNumber(to) &&
            to >= from &&
            to < (v.slice.to - v.slice.from));
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

    var structNames = function (args) {
        if (isNull(args) || !isStruct(car(args))) {
            return error("argument error");
        }

        var names = [];
        for (var name in car(args).struct) {
            names = cons(stringToSymbol(name), names);
        }

        return names;
    };

    // env
    var makeEnv = function (parent) {
        return makeStruct(
            [stringToSymbol("env"), [true, []]],
            [stringToSymbol("parent"), [parent, []]]);
    };

    var isEnv = function (env) {
        return isStruct(env) &&
            structDefined(env, stringToSymbol("env")) &&
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

    var procIds = 0;

    var defineVar = function (env, name, value) {
        if ((isPrimitiveProcedure(value) ||
            isCompiledProcedure(value)) &&
            !value.id) {
            value.id = name + ":" + (procIds++);
        }
        structDefine(env, name, value);
    };

    var envIds = 0;

    var extendEnv = function (env, names, values) {
        env = makeEnv(env);
        env.id = (envIds++);
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

    // call
    var callCheck = function (regs) {
        return isPrimitiveProcedure(regs.proc) ||
            isCompiledProcedure(regs.proc);
    };

    var callOp = function (regs, cont) {
        for (;;) {
            if (isPrimitiveProcedure(regs.proc) &&
                regs.proc.primitive === apply) {
                regs.proc = car(regs.args);
                regs.args = car(cdr(regs.args));
            } else if (isPrimitiveProcedure(regs.proc) &&
                regs.proc.primitive === call) {
                regs.proc = car(regs.args);
                regs.args = cdr(regs.args);
            } else if (isPrimitiveProcedure(regs.proc) &&
                regs.proc.primitive === callCc) {
                var contSave = {
                    env: regs.env,
                    proc: regs.proc,
                    args: regs.args,
                    cont: cont || regs.cont,
                    stack: stack.slice()
                };

                regs.proc = car(regs.args);
                var originalProc = regs.proc;
                var continuation = importPrimitive(function (args) {
                    if (isNull(args) ||
                        !isNull(cdr(args))) {
                        return error("invalid arity");
                    }

                    regs.env = contSave.env;
                    regs.proc = contSave.proc;
                    regs.args = contSave.args;
                    regs.cont = contSave.cont;
                    stack = contSave.stack.slice();
                    return car(args);
                }, function () {
                    var regsSave = {
                        env: regs.env,
                        proc: regs.proc,
                        args: regs.args,
                        cont: regs.cont,
                        stack: stack.slice()
                    };
                    regs.args = importArray(Array.prototype.slice.call(arguments));
                    regs.proc = continuation;
                    callOp(regs, false);
                    control();
                    originalProc.continuedExternally = true;
                    var result = regs.val;
                    regs.env = regsSave.env;
                    regs.proc = regsSave.proc;
                    regs.args = regsSave.arsg;
                    regs.cont = regsSave.cont;
                    stack = regsSave.stack.slice();
                    return result;
                });
                continuation.cc = cont || regs.cont;
                regs.args = ops.list(continuation);
            } else if (isPrimitiveProcedure(regs.proc) &&
                regs.proc.primitive === breakExecution) {
                regs.proc.continuedExternally = true;
                break;
            } else {
                break;
            }
        }

        if (isPrimitiveProcedure(regs.proc)) {
            var proc = regs.proc;
            regs.val = applyPrimitive(regs.proc, regs.args);
            index = proc.continuedExternally && (program.length - 1) || proc.cc || cont || regs.cont;
            return;
        };

        if (cont) {
            regs.cont = cont;
        }

        index = compiledEntry(regs.proc);
    };

    // import/export
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
            index = compiledEntry(p);
            regs.cont = program.length - 1;
            control();
            var val = regs.val;
            regs = regsSave;
            regs.cont = program.length - 1;
            return exportVal(val);
        };
        f.__exported = p;
        return f;
    };

    var isExportedProcedure = function (p) {
        return isFunction(p) && p.__exported;
    };

    var importExportedProcedure = function (p) {
        return p.__exported;
    };

    var importPrimitive = function (p, external) {
        return {primitive: p, external: external};
    };

    var exportPrimitiveProcedure = function (p) {
        if (p.external) {
            return p.external;
        }

        var exported = function () {
            return p.primitive.call(this, importArray(Array.prototype.slice.call(arguments)));
        };
        exported.__exportedPrimitive = p;
        return exported;
    };

    var importFunction = function (f, ctx, value, convert) {
        if (f.__exportedPrimitive) {
            return f.__exportedPrimitive;
        }

        return importPrimitive(function (args) {
            var result = f.apply(ctx || this, (convert ? exportPair : listToArray)(args));
            if (value !== undefined) {
                return value;
            }
            return convert ? importVal(result) : result;
        }, f);
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

    var importArray = function (val) {
        var l = list();
        for (var i = val.length - 1; i >= 0; i--) {
            l = cons(importVal(val[i]), l);
        }
        return l;
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

    var jsImportCode = importModule;

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

    var bitwiseOr = function (args) {
        var result = 0;
        for (;;) {
            if (isNull(args)) {
                return result;
            }

            if (!isPair(args) || !isNumber(car(args))) {
                return error("argument error");
            }

            result |= car(args);
            args = cdr(args);
        }
    };

    var bitwiseAnd = function (args) {
        var result = -1;
        for (;;) {
            if (isNull(args)) {
                return result;
            }

            if (!isPair(args) || !isNumber(car(args))) {
                return error("argument error");
            }

            result &= car(args);
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

    var greater = makeCompareNumbers(function (left, right) {
        return left > right;
    });

    var greaterOrEquals = makeCompareNumbers(function (left, right) {
        return left >= right;
    });

    var floor = function (args) {
        if (isNull(args) ||
            !isNumber(car(args)) ||
            !isNull(cdr(args))) {
            return error("argument error");
        }
        return Math.floor(car(args));
    };

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
            return isNumber(entry) && isEnv(env);
        }, makeProcedure),
        compiledEntry: func(1, false, isCompiledProcedure, compiledEntry),
        compiledProcedureEnv: func(1, false, isCompiledProcedure, compiledProcedureEnv),
        cons: func(2, false, false, cons),
        list: func(0, true, false, list),
        call: func(2, false, callCheck, callOp)
    };

    // registers
    var regs = {
        flag: false,
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

    var apply = function () {};
    var call = function () {};

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

    defineVar(regs.env, stringToSymbol("error?"),
        importFunction(func(1, false, false, isError)));

    defineVar(regs.env, stringToSymbol("error"),
        importFunction(func(1, false, false, error)));

    defineVar(regs.env, stringToSymbol("error->string"),
        importFunction(func(1, false, isError, errorToString)));

    defineVar(regs.env, stringToSymbol("true"), true);

    defineVar(regs.env, stringToSymbol("false"), false);

    defineVar(regs.env, stringToSymbol("number?"),
        importFunction(func(1, false, false, isNumber)));

    defineVar(regs.env, stringToSymbol("number->string"),
        importFunction(func(1, false, isNumber, numberToString)));

    defineVar(regs.env, stringToSymbol("string->number"),
        importFunction(func(1, false, isString, stringToNumber)));

    defineVar(regs.env, stringToSymbol("compiled-procedure?"),
        importFunction(func(1, false, false, isCompiledProcedure)));

    defineVar(regs.env, stringToSymbol("primitive-procedure?"),
        importFunction(func(1, false, false, isPrimitiveProcedure)));

    defineVar(regs.env, stringToSymbol("apply"),
        importPrimitive(apply));

    defineVar(regs.env, stringToSymbol("call"),
        importPrimitive(call));

    defineVar(regs.env, stringToSymbol("call/cc"),
        importPrimitive(callCc));

    defineVar(regs.env, stringToSymbol("break-execution"),
        importPrimitive(breakExecution));

    defineVar(regs.env, stringToSymbol("="),
        importPrimitive(numberEqual));

    defineVar(regs.env, stringToSymbol("+"),
        importPrimitive(add));

    defineVar(regs.env, stringToSymbol("-"),
        importPrimitive(subtract));

    defineVar(regs.env, stringToSymbol("*"),
        importPrimitive(multiply));

    defineVar(regs.env, stringToSymbol("|"),
        importPrimitive(bitwiseOr));

    defineVar(regs.env, stringToSymbol("&"),
        importPrimitive(bitwiseAnd));

    defineVar(regs.env, stringToSymbol("<"),
        importPrimitive(less));

    defineVar(regs.env, stringToSymbol(">"),
        importPrimitive(greater));

    defineVar(regs.env, stringToSymbol(">="),
        importPrimitive(greaterOrEquals));

    defineVar(regs.env, stringToSymbol("floor"),
        importPrimitive(floor));

    defineVar(regs.env, stringToSymbol("out"),
        importFunction(console.log, console, noprint));

    defineVar(regs.env, stringToSymbol("vector"),
        importPrimitive(vector));

    defineVar(regs.env, stringToSymbol("vector?"),
        importFunction(func(1, false, false, isVector)));

    defineVar(regs.env, stringToSymbol("vector-ref"),
        importFunction(func(2, false, vectorRefCheck, vectorRef)));

    defineVar(regs.env, stringToSymbol("vector-length"),
        importFunction(func(1, false, isVector, vectorLength)));

    defineVar(regs.env, stringToSymbol("vector-slice"),
        importFunction(func(1, true, vectorSliceCheck, vectorSlice)));

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

    defineVar(regs.env, stringToSymbol("struct-names"),
        importPrimitive(structNames));

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

    defineVar(regs.env, stringToSymbol("string-index"),
        importFunction(func(2, false, stringIndexCheck, stringIndex)))

    defineVar(regs.env, stringToSymbol("string-append"),
        importPrimitive(stringAppend));

    defineVar(regs.env, stringToSymbol("string->symbol"),
        importFunction(func(1, false, isString, stringToSymbol)));

    var symbolNameEq = function (symbol, name) {
        return symbolToString(symbol) === name;
    }

    var makeInstructionQuery = function (name) {
        return function (inst) {
            return symbolNameEq(car(inst), name);
        };
    };

    var makeInstruction = function (query, ops) {
        return function () {
            if (!query(instruction)) {
                return false;
            }

            ops();
            return true;
        };
    };

    var isLabel = makeInstructionQuery("label");
    var isReg = makeInstructionQuery("reg");
    var isConst = makeInstructionQuery("const");
    var isOp = makeInstructionQuery("op");
    var isGoto = makeInstructionQuery("goto");
    var isSave = makeInstructionQuery("save");
    var isRestore = makeInstructionQuery("restore");
    var isAssign = makeInstructionQuery("assign");
    var isPerform = makeInstructionQuery("perform");
    var isPerformContinue = makeInstructionQuery("perform-continue");
    var isTest = makeInstructionQuery("test");
    var isBranch = makeInstructionQuery("branch");
    var isJsImportCode = makeInstructionQuery("js-import-code");

    var labelValue = function (label) {
        return car(cdr(label));
    };

    var constValue = function (c) {
        return car(cdr(c));
    };

    var regValue = function (reg) {
        return regs[symbolToString(car(cdr(reg)))];
    };

    var setRegValue = function (reg, value) {
        regs[symbolToString(car(cdr(reg)))] = value;
    };

    var gotoIndex = function () {
        var arg = car(cdr(instruction));
        switch (true) {
        case isLabel(arg):
            return labelValue(arg);
        case isReg(arg):
            return regValue(arg);
        default:
            return error("invalid goto instruction", arg);
        }
    };

    var instGoto = makeInstruction(isGoto, function () {
        index = gotoIndex();
    });

    var instSave = makeInstruction(isSave, function () {
        save(regValue(instruction));
    });

    var instRestore = makeInstruction(isRestore, function () {
        setRegValue(instruction, restore());
    });

    var opName = function (op) {
        switch (symbolToString(car(cdr(op)))) {
        case "lookup-variable-value":
            return "lookupVar";
        case "set-variable-value!":
            return "setVar";
        case "define-variable!":
            return "defineVar";
        case "false?":
            return "isFalse";
        case "compiled-procedure-env":
            return "compiledProcedureEnv";
        case "extend-environment":
            return "extendEnv";
        case "make-compiled-procedure":
            return "makeProcedure";
        case "cons":
            return "cons";
        case "list":
            return "list";
        case "compiled-procedure-entry":
            return "compiledEntry";
        case "procedure-call":
            return "call";
        default:
            return error("invalid operation", op);
        }
    };

    var getArgValue = function (arg) {
        switch (true) {
        case isReg(arg):
            return regValue(arg);
        case isLabel(arg):
            return labelValue(arg);
        case isConst(arg):
            return constValue(arg);
        case isSymbol(arg) && symbolNameEq(arg, "regs"):
            return regs;
        case isSymbol(arg) && symbolNameEq(arg, "false"):
            return false;
        default:
            return error("invalid op argument", arg);
        }
    };

    var opArgs = function (argList) {
        if (isNull(argList)) {
            return [];
        }

        var args = opArgs(cdr(argList));
        args.unshift(getArgValue(car(argList)));
        return args;
    };

    var instOpCall = function (opInst) {
        return ops[opName(car(opInst))].apply(this, opArgs(cdr(opInst)));
    };

    var assignValue = function () {
        var valueInst = car(cdr(cdr(instruction)));
        switch (true) {
        case isReg(valueInst):
            return regValue(valueInst);
        case isLabel(valueInst):
            return labelValue(valueInst);
        case isConst(valueInst):
            return constValue(valueInst);
        case isOp(valueInst):
            return instOpCall(cdr(cdr(instruction)));
        default:
            return error("invalid assignment", instruction);
        }
    };

    var instAssign = makeInstruction(isAssign, function () {
        setRegValue(instruction, assignValue());
    });

    var instPerform = makeInstruction(isPerform, function () {
        instOpCall(cdr(instruction));
    });

    var instPerformContinue = makeInstruction(isPerformContinue, function () {
        instOpCall(cdr(instruction));
        // index = regs.cont;
    });

    var testValue = function () {
        var arg = car(cdr(instruction));
        switch (true) {
        case isReg(arg):
            return regValue(arg);
        case isOp(arg):
            return instOpCall(cdr(instruction));
        default:
            return error("invalid test", instruction);
        }
    };

    var instTest = makeInstruction(isTest, function () {
        regs.flag = testValue();
    });

    var instBranch = makeInstruction(isBranch, function () {
        if (isFalse(regs.flag)) {
            return;
        }
        index = isFalse(regs.flag) ? index : labelValue(car(cdr(instruction)));
    });

    var instJsImportCode = makeInstruction(isJsImportCode, function () {
        jsImportCode(car(cdr(car(cdr(instruction)))), car(cdr(cdr(instruction))));
    });

    var isEnd = function () {
        return isSymbol(instruction) && symbolNameEq(instruction, "end");
    };

    var flushVal = function (val, flushed) {
        if ((isPair(val) ||
            typeof val === "object") &&
            flushed.indexOf(val) >= 0) {
            return "circular";
        }
        flushed.push(val);
        var s = "";
        if (isPair(val)) {
            s += "[";
            s += flushVal(car(val), flushed);
            s += ", ";
            s += flushVal(cdr(val), flushed);
            s += "]";
        } else if (isSymbol(val)) {
            s += val[0];
        } else if (isNull(val)) {
            s += "[]";
        } else if (isEnv(val)) {
            s += "env-" + val.id;
        } else if (isPrimitiveProcedure(val)) {
            s += "primitive-" + val.id;
        } else if (typeof val === "object") {
            s += "{";
            for (var key in val) {
                s += key;
                s += ": ";
                s += flushVal(val[key], flushed);
                s += ", ";
            }
            s += "}";
        } else {
            s += String(val);
        }
        return s;
    };

    var flushStack = function () {
        var s = "[";
        for (var i = 0; i < stack.length; i++) {
            s += flushVal(stack[i], []);
            s += ", ";
        }
        s += "]";
        return s;
    };

    var flushState = function (s) {
        s += "{env: ";
        s += regs.env.id;
        s += ", proc: ";
        s += regs.proc.id;
        s += ", val: ";
        s += flushVal(regs.val, []);
        s += ", args: ";
        s += flushVal(regs.args, []);
        s += ", cont: ";
        s += String(regs.cont);
        s += ", stack: ";
        s += flushStack();
        s += "}";

        console.error(s);
    };

    var program = [
    /* [program] */
    ];

    // control
    var index = 0;
    var instruction;
    var control = function () {
        for (;;) {
            try {
                instruction = program[index++];
                switch (true) {
                case instGoto(): break;
                case instSave(): break;
                case instRestore(): break;
                case instAssign(): break;
                case instPerform(): break;
                case instPerformContinue(): break;
                case instTest(): break;
                case instBranch(): break;
                case instJsImportCode(): break;
                case isEnd(): return;
                default:
                    error("invalid instruction", instruction);
                }
            } catch (error) {
                // flushState((index - 1) + ": " + flushVal(instruction, []) + " - ");
                throw error;
            } finally {
                // flushState((index - 1) + ": " + flushVal(instruction, []) + " - ");
            }
        }
    };
    control();
})();
