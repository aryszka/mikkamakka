var lang = (function () {
    var isSymbol = function (s) {
        return !!s &&
            s instanceof Array &&
            s.length === 1 &&
            typeof s[0] === "string";
    };

    var symbolName = function (symbol) {
        return symbol[0];
    };

    var stringToSymbol = function (name) {
        return [name];
    };

    var peq = function (left, right) {
        return left === right;
    };

    var filter = function (l, pred) {
        if (isNull(l)) {
            return list();
        }
        return pred(car(l)) ?
            cons(car(l), filter(l, cdr(l))) :
            filter(l, cdr(l));
    };

    var makeNameTable = function () {
        var names = [];
        var table = {};
        return function (modNames, op) {
            names = modNames(names, table);
            return op(names, table);
        };
    };

    var tableNoop = function (names, _) {
        return names;
    };

    var tableNames = function (table) {
        return table(tableNoop, tableNoop);
    };

    var tableHasName = function (table, name) {
        return table(tableNoop, function (_, table) {
            return symbolName(name) in table;
        });
    };

    var tableLookup = function (table, name) {
        return table(tableNoop, function (_, table) {
            name = symbolName(name);
            var val = table[name];
            if (val === undefined) {
                return cerror("tableLookup", "name is not defined", name);
            }
            return val;
        });
    };

    var tableDefine = function (table, name, val) {
        var sname = symbolName(name);
        return table(function (names, table) {
            return sname in table ? names : [name, names];
        }, function (_, table) {
            table[sname] = val;
            return val;
        });
    };

    var tableSet = function (table, name, val) {
        return table(tableNoop, function (_, table) {
            name = symbolName(name);
            if (!(name in table)) {
                return cerror("tableSet", "name is not defined", name);
            }
            table[name] = val;
            return val;
        });
    };

    var tableDelete = function (table, name) {
        return table(function (names, _) {
            return filter(names, function (i) { return i[0] === name[0]; });
        }, function (_, table) {
            var has = name[0] in table;
            delete table[name[0]];
            return has;
        });
    };

    var extendEnv = function (env) {
        var frame = {};
        return function (name, def, value, check) {
            if (isSymbol(name)) {
                name = symbolName(name);
            }
            if (check) {
                return name in frame;
            }
            if (def === true) {
                frame[name] = value;
                return value;
            }
            if (def === false) {
                if (name in frame) {
                    frame[name] = value;
                    return value;
                }
                return env(name, def, value, check);
            }
            if (name in frame) {
                return frame[name];
            }
            return env(name, def, value, check);
        };
    };

    var isDefined = function (environment, name) {
        if (arguments.length < 2) {
            name = environment;
            environment = env;
        }
        return environment(name, false, false, true);
    };

    var isNull = function (l) {
        return !!l &&
            l instanceof Array &&
            l.length === 0;
    };

    var cons = function (left, right) {
        return [left, right];
    };

    var isPair = function (p) {
        return !!p &&
            p instanceof Array &&
            p.length === 2;
    };

    var car = function (p) {
        return p[0];
    };

    var cdr = function (p) {
        return p[1];
    };

    var list = function () {
        var l = [];
        for (var i = arguments.length - 1; i >= 0; i--) {
            l = cons(arguments[i], l);
        }
        return l;
    };

    var isCompiledProcedure = function (p) {
        return !!p && typeof p === "function";
    };

    var capply = function (p, args) {
        var jsArgs = [];
        for (;;) {
            if (isNull(args)) {
                break;
            }
            jsArgs[jsArgs.length] = car(args);
            args = cdr(args);
        }
        return p.apply(undefined, jsArgs);
    };

    var tryc = function (t, c) {
        try {
            return env("apply")(t, list());
        } catch (error) {
            return env("apply")(c, list(error));
        }
    };

    var forLoop = function (body) {
        for (;;) {
            var result = env("apply")(body, list());
            if (false !== result) {
                return result;
            }
        }
    };

    var readFile = function (f, clb) {
        require("fs").readFile(f, function (err, data) {
            if (err) {
                return cerror("readFile", "error while reading file", err);
            }
            return env("apply")(clb, list(data.toString()));
        });
        return noPrint;
    };

    var readLine = function (clb, prompt) {
        var rl = require("readline").createInterface(process.stdin, process.stdout);
        var frl = function (f) {
            return f(rl);
        };
        rl.on("line", function (l) {
            env("apply")(clb, list(frl, l));
        });
        return frl;
    };

    var prompt = function (rl) {
        return rl(function (rrl) {
            rrl.prompt();
            return noPrint;
        });
    };

    var setPrompt = function (rl, p) {
        return rl(function (rrl) {
            rrl.setPrompt(p);
            return noPrint;
        });
    };

    var isNumber = function (n) {
        return typeof n === "number" && !Number.isNaN(n);
    };

    var parseNumber = function (s) {
        var num = parseFloat(s);
        if (isNumber(num)) {
            return num;
        }
        return s;
    };

    var identity = function (x) { return x; };

    var mkNumOp = function (initial, single, reduce) {
        return function () {
            var args = Array.prototype.slice.call(arguments);
            if (!args.length) {
                return initial;
            }
            if (args.length === 1) {
                return single(args[0]);
            }
            return args.slice(1).reduce(function (previous, current) {
                return reduce(previous, current);
            }, args[0]);
        };
    };

    var add = mkNumOp(0, identity, function (x, y) { return x + y; });

    var sub = mkNumOp(0, function (x) {
        return 0 - x;
    }, function (x, y) {
        return x - y;
    });

    var mul = mkNumOp(1, identity, function (x, y) { return x * y; });

    var div = mkNumOp(1, function (x) {
        return 1 / x;
    }, function (x, y) {
        return x / y;
    });

    var mod = function (dnd, dvs) {
        return dnd % dvs;
    };

    var gt = function (left, right) {
        return left > right;
    };

    var lt = function (left, right) {
        return left < right;
    };

    var gte = function (left, right) {
        return left >= right;
    };

    var lte = function (left, right) {
        return left <= right;
    };

    var isString = function (s) {
        return typeof s === "string";
    };

    var slen = function (s) {
        return s.length;
    };

    var sidx = function (s, expression) {
        var m = s.match(new RegExp(expression), "");
        if (!m) {
            return -1;
        }
        return m.index;
    };

    var charAt = function (s, i) {
        return s[i];
    };

    var subs = function (s, offset, count) {
        if (count < 0) {
            return s.substr(offset);
        }
        return s.substr(offset, count);
    };

    var sreplace = function (s, expression, replace) {
        return s.replace(new RegExp(expression, "g"), replace);
    };

    var mkStringBuilder = function () {
        return {sb: []};
    };

    var isStringBuilder = function (b) {
        return !!b && b.sb;
    };

    var sbempty = function (b) {
        for (var i = 0; i < b.sb.length; i++) {
            if (b.sb[i].length) {
                return false;
            }
        }
        return true;
    };

    var sbappend = function (b, s) {
        if (isStringBuilder(s)) {
            return {
                sb: b.sb.concat(s.sb)
            };
        }
        return {
            sb: b.sb.concat(s)
        }
    };

    var builderToString = function (b) {
        return b.sb.join("");
    };

    // no-head, once variadic arguments
    var cats = function () {
        var b = mkStringBuilder();
        for (var i = 0; i < arguments.length; i++) {
            b = sbappend(b, arguments[i]);
        }
        return builderToString(b);
    };

    var noPrint = function () { return noPrint; };

    var out = process && process.stdout && function (str) {
        if (str !== noPrint) {
            process.stdout.write(String(str));
        }
        return noPrint;
    };

    var log = process && process.stderr && function (s) {
        if (s !== noPrint) {
            process.stderr.write(String(s));
        }
        process.stderr.write("\n");
        return noPrint;
    } || function (s) {
        if (s !== noPrint) {
            console.log(s);
        } else {
            console.log();
        }
        return noPrint;
    };

    var cerror = function (where, what, arg) {
        throw new Error(String(where) + ":" + String(what) + ":" + String(arg));
        return noPrint;
    };

    var isError = function (error) {
        return error instanceof Error;
    };

    var sprintError = function (error) {
        return error && error.message || "error";
    };

    var sprintStack = function (error) {
        return error && error.stack || "";
    };

    var now = function () {
        return {t: new Date()};
    };

    var isTime = function (t) {
        return !!t && t.t instanceof Date;
    };

    var numberToTime = function (n) {
        return new Date(n);
    };

    var timeToNumber = function (t) {
        return t.t.valueOf();
    };

    var timeToString = function (t) {
        return t.t.toString();
    };

    var exit = function (val) {
        process.exit(val);
    };

    var argv = function () {
        return list.apply(undefined, process.argv);
    };

    var makeRegexp = function (expression, flags) {
        var regexp = new RegExp(expression, flags);
        return function (text) {
            return vector.apply(undefined, text.match(regexp) || []);
        };
    };

    var isVector = function (exp) {
        return !!exp &&
            exp.vector instanceof Array;
    };

    var vector = function () {
        return {vector: Array.prototype.slice.call(arguments)};
    };

    var vlen = function (v) {
        return v.vector.length;
    };

    var vref = function (v, i) {
        return v.vector[i];
    };

    var listToVector = function (list, reverse) {
        var v = vector();
        for (;;) {
            if (isNull(list)) {
                return v;
            }
            v.vector[reverse ? "unshift" : "push"](car(list));
            list = cdr(list);
        }
    };

    var vectorToList = function (v, reverse) {
        var l = list();
        var len = vlen(v);
        for (var i = 0; i < len; i++) {
            l = cons(vref(reverse ? len - i - 1 : i), l);
        }
        return l;
    };

    var lang = makeNameTable();
    tableDefine(lang, "compiled-procedure?", isCompiledProcedure);
    tableDefine(lang, "capply", capply);
    tableDefine(lang, "symbol?", isSymbol);
    tableDefine(lang, "symbol-name", symbolName);
    tableDefine(lang, "string->symbol", stringToSymbol);
    tableDefine(lang, "identity", identity);
    tableDefine(lang, "null?", isNull);
    tableDefine(lang, "cons", cons);
    tableDefine(lang, "pair?", isPair);
    tableDefine(lang, "car", car);
    tableDefine(lang, "cdr", cdr);
    tableDefine(lang, "list", list);
    tableDefine(lang, "peq?", peq);
    tableDefine(lang, "try", tryc);
    tableDefine(lang, "for", forLoop);
    tableDefine(lang, "read-file", readFile);
    tableDefine(lang, "read-line", readLine);
    tableDefine(lang, "prompt", prompt);
    tableDefine(lang, "set-prompt", setPrompt);
    tableDefine(lang, "number?", isNumber);
    tableDefine(lang, "parse-number", parseNumber);
    tableDefine(lang, "+", add);
    tableDefine(lang, "-", sub);
    tableDefine(lang, "*", mul);
    tableDefine(lang, "/", div);
    tableDefine(lang, "%", mod);
    tableDefine(lang, ">", gt);
    tableDefine(lang, "<", lt);
    tableDefine(lang, ">=", gte);
    tableDefine(lang, "<=", lte);
    tableDefine(lang, "string?", isString);
    tableDefine(lang, "number->string", String);
    tableDefine(lang, "slen", slen);
    tableDefine(lang, "sidx", sidx);
    tableDefine(lang, "char-at", charAt);
    tableDefine(lang, "subs", subs);
    tableDefine(lang, "sreplace", sreplace);
    tableDefine(lang, "make-string-builder", mkStringBuilder);
    tableDefine(lang, "string-builder?", isStringBuilder);
    tableDefine(lang, "sbempty?", sbempty);
    tableDefine(lang, "sbappend", sbappend);
    tableDefine(lang, "builder->string", builderToString);
    tableDefine(lang, "cats", cats);
    tableDefine(lang, "out", out);
    tableDefine(lang, "clog", log);
    tableDefine(lang, "now", now);
    tableDefine(lang, "time?", isTime);
    tableDefine(lang, "number->time", numberToTime);
    tableDefine(lang, "time->number", timeToNumber);
    tableDefine(lang, "time->string", timeToString);
    tableDefine(lang, "filter", filter);
    tableDefine(lang, "make-name-table", makeNameTable);
    tableDefine(lang, "table-has-name?", tableHasName);
    tableDefine(lang, "table-names", tableNames);
    tableDefine(lang, "table-lookup", tableLookup);
    tableDefine(lang, "table-define", tableDefine);
    tableDefine(lang, "table-set!", tableSet);
    tableDefine(lang, "table-delete!", tableDelete);
    tableDefine(lang, "cerror", cerror);
    tableDefine(lang, "error?", isError);
    tableDefine(lang, "sprint-error", sprintError);
    tableDefine(lang, "sprint-stack", sprintStack);
    tableDefine(lang, "exit", exit);
    tableDefine(lang, "proc-argv", argv);
    tableDefine(lang, "make-regexp", makeRegexp);
    tableDefine(lang, "true");
    tableDefine(lang, "false", false);
    tableDefine(lang, "no-print", noPrint);
    tableDefine(lang, "vector?", isVector);
    tableDefine(lang, "vector", vector);
    tableDefine(lang, "vlen", vlen);
    tableDefine(lang, "vref", vref);
    tableDefine(lang, "list->vector", listToVector);
    tableDefine(lang, "vector->list", vectorToList);
    tableDefine(lang, "lang", lang);
    return lang;
})();
