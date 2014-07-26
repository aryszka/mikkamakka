var env = (function () {
    var noPrint = function () { return noPrint; };

    var cerror = function (where, what, arg) {
        throw new Error(String(where) + ":" + String(what) + ":" + String(arg));
        return noPrint;
    };

    var argCheck = function (argObj, length, variadic) {
        if (argObj.length < length || !variadic && argObj.length !== length) {
            return cerror("procedure", "invalid number of arguments", argObj.length);
        }
        return noPrint;
    };

    var isSymbol = function (s) {
        argCheck(arguments, 1, false);
        return !!s &&
            s instanceof Array &&
            s.length === 1 &&
            typeof s[0] === "string";
    };

    var symbolName = function (symbol) {
        argCheck(arguments, 1, false);
        return symbol[0];
    };

    var stringToSymbol = function (name) {
        argCheck(arguments, 1, false);
        return [name];
    };

    var peq = function (left, right) {
        argCheck(arguments, 2, false);
        return left === right;
    };

    var filter = function (pred, l) {
        argCheck(arguments, 2, false);
        if (isNull(l)) {
            return list();
        }
        return pred(car(l)) ?
            cons(car(l), filter(pred, cdr(l))) :
            filter(pred, cdr(l));
    };

    var isTable = function (t) {
        argCheck(arguments, 1, false);
        return !!(t && t.table);
    };

    var makeNameTable = function () {
        argCheck(arguments, 0, false);
        return {table: {}};
    };

    var tableNames = function (table) {
        argCheck(arguments, 1, false);
        var names = [];
        for (var sname in table.table) {
            names = [[sname], names];
        }
        return names;
    };

    var tableHasName = function (table, name) {
        argCheck(arguments, 2, false);
        return symbolName(name) in table.table;
    };

    var tableLookup = function (table, name) {
        argCheck(arguments, 2, false);
        var val = table.table[symbolName(name)];
        if (typeof val === "undefined") {
            return cerror("tableLookup", "name is not defined", name);
        }
        return val;
    };

    var tableDefine = function (table, name, val) {
        argCheck(arguments, 3, false);
        table.table[symbolName(name)] = val;
        return val;
    };

    var tableSet = function (table, name, val) {
        argCheck(arguments, 3, false);
        if (!tableHasName(table, name)) {
            return cerror("tableSet", "name is not defined", name);
        }
        table.table[symbolName(name)] = val;
        return val;
    };

    var tableDelete = function (table, name) {
        argCheck(arguments, 2, false);
        var has = tableHasName(table, name);
        if (has) {
            delete table.table[name];
        }
        return has;
    };

    var isNull = function (l) {
        argCheck(arguments, 1, false);
        return !!l &&
            l instanceof Array &&
            l.length === 0;
    };

    var cons = function (left, right) {
        argCheck(arguments, 2, false);
        return [left, right];
    };

    var isPair = function (p) {
        argCheck(arguments, 1, false);
        return !!p &&
            p instanceof Array &&
            p.length === 2;
    };

    var car = function (p) {
        argCheck(arguments, 1, false);
        return p[0];
    };

    var cdr = function (p) {
        argCheck(arguments, 1, false);
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
        argCheck(arguments, 1, false);
        return !!p && (typeof p.main === "function" || typeof p === "function");
    };

    var getMain = function (f) {
        argCheck(arguments, 1, false);
        return f.main || f;
    };

    var capply = function (p, args) {
        argCheck(arguments, 2, false);
        var jsArgs = [];
        for (;;) {
            if (isNull(args)) {
                break;
            }
            jsArgs[jsArgs.length] = car(args);
            args = cdr(args);
        }
        return getMain(p).apply(undefined, jsArgs);
    };

    var tryc = function (t, c) {
        argCheck(arguments, 2, false);
        try {
            return getMain(env("apply"))(t, list());
        } catch (error) {
            return getMain(env("apply"))(c, list(error));
        }
    };

    var readFile = function (f, clb) {
        argCheck(arguments, 2, false);
        require("fs").readFile(f, function (err, data) {
            if (err) {
                return cerror("readFile", "error while reading file", err);
            }
            return getMain(env("apply"))(clb, list(data.toString()));
        });
        return noPrint;
    };

    var readLine = function (clb, prompt) {
        argCheck(arguments, 2, false);
        var rl = require("readline").createInterface(process.stdin, process.stdout);
        var frl = function (f) {
            return f(rl);
        };
        rl.on("line", function (l) {
            getMain(env("apply"))(clb, list(frl, l));
        });
        return frl;
    };

    var prompt = function (rl) {
        argCheck(arguments, 1, false);
        return rl(function (rrl) {
            rrl.prompt();
            return noPrint;
        });
    };

    var setPrompt = function (rl, p) {
        argCheck(arguments, 2, false);
        return rl(function (rrl) {
            rrl.setPrompt(p);
            return noPrint;
        });
    };

    var isNumber = function (n) {
        argCheck(arguments, 1, false);
        return typeof n === "number" && !Number.isNaN(n);
    };

    var parseNumber = function (s) {
        argCheck(arguments, 1, false);
        var num = parseFloat(s);
        if (isNumber(num)) {
            return num;
        }
        return s;
    };

    var identity = function (x) {
        argCheck(arguments, 1, false);
        return x;
    };

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
        argCheck(arguments, 2, false);
        return dnd % dvs;
    };

    var gt = function (left, right) {
        argCheck(arguments, 2, false);
        return left > right;
    };

    var lt = function (left, right) {
        argCheck(arguments, 2, false);
        return left < right;
    };

    var gte = function (left, right) {
        argCheck(arguments, 2, false);
        return left >= right;
    };

    var lte = function (left, right) {
        argCheck(arguments, 2, false);
        return left <= right;
    };

    var isString = function (s) {
        argCheck(arguments, 1, false);
        return typeof s === "string";
    };

    var slen = function (s) {
        argCheck(arguments, 1, false);
        return s.length;
    };

    var sidx = function (s, expression) {
        argCheck(arguments, 2, false);
        var m = s.match(new RegExp(expression), "");
        if (!m) {
            return -1;
        }
        return m.index;
    };

    var charAt = function (s, i) {
        argCheck(arguments, 2, false);
        return s[i];
    };

    var charCodeAt = function (s, i) {
        argCheck(arguments, 2, false);
        return s.charCodeAt(s, i);
    };

    var subs = function (s, offset, count) {
        argCheck(arguments, 3, false);
        if (count < 0) {
            return s.substr(offset);
        }
        return s.substr(offset, count);
    };

    var sreplace = function(s, expression, replace) {
        argCheck(arguments, 3, false);
        return s.replace(new RegExp(expression, "g"), function (s) {
            return env("apply")(replace, list(s, false, false));
        });
    };

    var mkStringBuilder = function () {
        argCheck(arguments, 0, false);
        return {sb: []};
    };

    var isStringBuilder = function (b) {
        argCheck(arguments, 1, false);
        return !!b && b.sb;
    };

    var sbempty = function (b) {
        argCheck(arguments, 1, false);
        for (var i = 0; i < b.sb.length; i++) {
            if (b.sb[i].length) {
                return false;
            }
        }
        return true;
    };

    var sbappend = function (b, s) {
        argCheck(arguments, 2, false);
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
        argCheck(arguments, 1, false);
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

    var out = process && process.stdout && function (str) {
        argCheck(arguments, 1, false);
        if (str !== noPrint) {
            process.stdout.write(String(str));
        }
        return noPrint;
    } || console && console.log || noPrint;

    var log = process && process.stderr && function (s) {
        argCheck(arguments, 1, false);
        if (s !== noPrint) {
            process.stderr.write(String(s));
        }
        process.stderr.write("\n");
        return noPrint;
    } || function (s) {
        argCheck(arguments, 1, false);
        if (s !== noPrint) {
            console.log(s);
        } else {
            console.log();
        }
        return noPrint;
    };

    var isError = function (error) {
        argCheck(arguments, 1, false);
        return error instanceof Error;
    };

    var sprintError = function (error) {
        argCheck(arguments, 1, false);
        return error && error.message || "error";
    };

    var sprintStack = function (error) {
        argCheck(arguments, 1, false);
        return error && error.stack || "";
    };

    var now = function () {
        argCheck(arguments, 0, false);
        return {t: new Date()};
    };

    var isTime = function (t) {
        argCheck(arguments, 1, false);
        return !!t && t.t instanceof Date;
    };

    var numberToTime = function (n) {
        argCheck(arguments, 1, false);
        return new Date(n);
    };

    var timeToNumber = function (t) {
        argCheck(arguments, 1, false);
        return t.t.valueOf();
    };

    var timeToString = function (t) {
        argCheck(arguments, 1, false);
        return t.t.toString();
    };

    var exit = function (val) {
        argCheck(arguments, 1, false);
        process.exit(val);
    };

    var argv = function () {
        argCheck(arguments, 0, false);
        return list.apply(undefined, process.argv);
    };

    var makeRegexp = function (expression, flags) {
        argCheck(arguments, 2, false);
        var regexp = new RegExp(expression, flags);
        var wrap = function (text) {
            return vector.apply(undefined, text.match(regexp) || []);
        };
        wrap.main = wrap;
        wrap.body = wrap;
        return wrap;
    };

    var isVector = function (exp) {
        argCheck(arguments, 1, false);
        return !!exp &&
            exp.vector instanceof Array;
    };

    var vector = function () {
        return {vector: Array.prototype.slice.call(arguments)};
    };

    var vlen = function (v) {
        argCheck(arguments, 1, false);
        return v.vector.length;
    };

    var vref = function (v, i) {
        argCheck(arguments, 2, false);
        return v.vector[i];
    };

    var listToVector = function (list, reverse) {
        argCheck(arguments, 2, false);
        var v = vector();
        for (;;) {
            if (isNull(list)) {
                return v;
            }
            v.vector[reverse ? "unshift" : "push"](car(list));
            list = cdr(list);
        }
    };

    var vectorToList = function (v) {
        argCheck(arguments, 1, false);
        var l = list();
        var len = vlen(v);
        for (var i = len - 1; i >= 0; i--) {
            l = cons(vref(v, i), l);
        }
        return l;
    };

    var struct = function () {
        var s = makeNameTable();
        for (var i = 0; i < arguments.length; i++) {
            var l = arguments[i];
            if (l.length < 2 ||
                l[1].length < 2 ||
                l[1][1].length !== 0) {
                return cerror("struct", "invalid struct member", l);
            }
            tableDefine(s, l[0], l[1][0]);
        }
        return s;
    };

    var mktail = function (f, args) {
        argCheck(arguments, 2, false);
        var tail = function () {
            return f.apply(undefined, args);
        };
        tail.isTail = true;
        return tail;
    };

    var tailCall = function (f) {
        argCheck(arguments, 1, false);
        var fi = f;
        for (;;) {
            fi = fi.apply();
            if (!fi.isTail) {
                return fi;
            }
        }
    };

    var shared = makeNameTable();
    var share = function (name, member) {
        if (typeof member === "function") {
            member.main = member;
            member.body = member;
        }
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
    share("char-code-at", charCodeAt);
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
    share("struct", struct);
    share("struct?", isTable);
    share("mktail", mktail);
    share("tail-call", tailCall);

    // patch
    var extendEnv = function (env, shared) {
        argCheck(arguments, 2, false);
        var current = {
            parent: env,
            shared: shared || makeNameTable()
        };
        var wrap = function () {
            argCheck(arguments, 1, true);
            if (arguments.length === 4 && arguments[3]) {
                return tableNames(current.shared);
            }
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
        wrap.main = wrap;
        wrap.body = wrap;
        return wrap;
    };

    var isDefined = function (env, name) {
        argCheck(arguments, 2, false);
        try {
            env(name);
            return true;
        } catch (_) {
            return false;
        }
    };

    var head = extendEnv(null, shared);
    var mikkamakka = extendEnv(head, null);
    share("extend-env", extendEnv);
    share("defined?", isDefined);
    share("head", head);
    share("mikkamakka", mikkamakka);
    return mikkamakka;
})();
