var env = (function () {
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

    var filter = function (pred, l) {
        if (isNull(l)) {
            return list();
        }
        return pred(car(l)) ?
            cons(car(l), filter(pred, cdr(l))) :
            filter(pred, cdr(l));
    };

    var isTable = function (t) {
        return !!(t && t.isTable);
    };

    var makeNameTable = function () {
        var names = [];
        var table = {};
        var t = function (modNames, op) {
            names = modNames(names, table);
            return op(names, table);
        };
        t.isTable = true;
        return t;
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
            return filter(function (i) { return i[0] !== name[0]; }, names);
        }, function (_, table) {
            var has = name[0] in table;
            delete table[name[0]];
            return has;
        });
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

    var shared = makeNameTable();
    var share = function (name, member) {
        return tableDefine(shared, [name], member);
    };
    share("compiled-procedure?", isCompiledProcedure);
    share("capply", capply);
    share("symbol?", isSymbol);
    share("symbol-name", symbolName);
    share("string->symbol", stringToSymbol);
    share("identity", identity);
    share("null?", isNull);
    share("cons", cons);
    share("pair?", isPair);
    share("car", car);
    share("cdr", cdr);
    share("list", list);
    share("peq?", peq);
    share("try", tryc);
    share("for", forLoop);
    share("read-file", readFile);
    share("read-line", readLine);
    share("prompt", prompt);
    share("set-prompt", setPrompt);
    share("number?", isNumber);
    share("parse-number", parseNumber);
    share("+", add);
    share("-", sub);
    share("*", mul);
    share("/", div);
    share("%", mod);
    share(">", gt);
    share("<", lt);
    share(">=", gte);
    share("<=", lte);
    share("string?", isString);
    share("number->string", String);
    share("slen", slen);
    share("sidx", sidx);
    share("char-at", charAt);
    share("subs", subs);
    share("sreplace", sreplace);
    share("make-string-builder", mkStringBuilder);
    share("string-builder?", isStringBuilder);
    share("sbempty?", sbempty);
    share("sbappend", sbappend);
    share("builder->string", builderToString);
    share("cats", cats);
    share("out", out);
    share("clog", log);
    share("now", now);
    share("time?", isTime);
    share("number->time", numberToTime);
    share("time->number", timeToNumber);
    share("time->string", timeToString);
    share("filter", filter);
    share("table?", isTable);
    share("make-name-table", makeNameTable);
    share("table-has-name?", tableHasName);
    share("table-names", tableNames);
    share("table-lookup", tableLookup);
    share("table-define", tableDefine);
    share("table-set!", tableSet);
    share("table-delete!", tableDelete);
    share("cerror", cerror);
    share("error?", isError);
    share("sprint-error", sprintError);
    share("sprint-stack", sprintStack);
    share("exit", exit);
    share("proc-argv", argv);
    share("make-regexp", makeRegexp);
    share("true", true);
    share("false", false);
    share("no-print", noPrint);
    share("vector?", isVector);
    share("vector", vector);
    share("vlen", vlen);
    share("vref", vref);
    share("list->vector", listToVector);
    share("vector->list", vectorToList);

    // patch
    var extendEnv = function (env, shared) {
        var current = {
            parent: env,
            shared: shared || makeNameTable()
        };
        return function () {
            if (arguments.length === 1) {
                if (tableHasName(current.shared, [arguments[0]])) {
                    return tableLookup(current.shared, [arguments[0]]);
                }
                if (!current.parent) {
                    return cerror("extend-env", "unbound variable", arguments[0]);
                }
                return current.parent(arguments[0]);
            }
            if (arguments[1]) {
                return tableDefine(current.shared, [arguments[0]], arguments[2]);
            }
            if (tableHasName(current.shared, [arguments[0]])) {
                return tableSet(current.shared, [arguments[0]], arguments[2]);
            }
            if (!current.parent) {
                return cerror("extend-env", "unbound variable", arguments[0]);
            }
            return current.parent(arguments[0], false, arguments[2]);
        };
    };
    var isDefined = function (env, name) {
        try {
            env(name);
            return true;
        } catch (_) {
            return false;
        }
    };
    var head = extendEnv(null, shared);
    var mikkamakka = extendEnv(head);
    share("extend-env", extendEnv);
    share("defined?", isDefined);
    share("head", head);
    share("mikkamakka", mikkamakka);
    return mikkamakka;
})();