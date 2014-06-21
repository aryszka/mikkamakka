// todo:
// fix the situation with the primitive procedures and higher order (try evaluation +)
// create quasiquotation
// arbitrary number of variables
// debug utilities: stack
// sprint escaping: compound procedures (lookup in env)
// introduce slicing
// introduce macros
// rename sbuilder, find similar lisp name
// rewrite the string functions exploiting sslice and builder
// primitive error checks
// delete functionality
// tail recursion optimization:
// - tail context definition in r6rs
// - benchmarks and profiling
// missing syntax:
// - let, let*, letrec
// cleanup and refactor
// - rename exp to str where obvious
// - handling of eof in reading, and just scheme streams in general
// - fix error not pass undefined to sprint
// - subs not failing on out of range
// - eval environment
// - higher order callback situation has to be cleared up
// tracking and minifying identical symbols
// replace checking constructor with instanceof? matter of optimization
// reenact tests
// cleanup the case of the global environment
// generate ca/dr shortcuts
// introduce ports
// modules, tests
// make exit exit only from the current run-time
// vim repl: http://www.vim.org/scripts/script.php?script_id=4336
// check validation. e.g. don't allow lambda parameters other then symbols
// rewrite compilation with quasiquotation
// eq to variadic, check what else
// parse-number can return non mikkamakka types (NaN)
// don't override members in the global environment

// lang
var out = function (str) {
    process.stdout.write(str);
    return noPrint;
};
var error = function (where, what, arg) {
    throw new Error(String(where) + ":" + String(what) + ":" + sprintq(arg, true, false));
};
var isError = function (exp) {
    return exp instanceof Error;
};
var sprintError = function (error) {
    return error.message || "error";
};
var sprintStack = function (error) {
    return error.stack || "";
};
var noPrint = function () {
    return noPrint;
};
var print = function (exp) {
    if (exp !== noPrint) {
        out(sprint(exp));
    }
    return noPrint;
};
var display = function (exp) {
    out(sunescape(sprint(exp)));
    return noPrint;
};
var exit = function (val) {
    process.exit(val);
};

// fix error handling
var readFile = function (f, clb) {
    require("fs").readFile(f, function (err, data) {
        if (err) {
            error("readFile", "error while reading file", list(f, err));
        }
        apply(clb, list(data.toString()));
    });
};

// requires the hack with the primitive procedures for higher order functions
// try to eliminate by a language construct
var tryc = function (t, c) {
    try {
        return apply(t, list());
    } catch (error) {
        return apply(c, list(error));
    }
};

// requires the hack with the primitive procedures for higher order functions
// try to eliminate by a language construct
// temporary feature until tail recursion optimization
var forLoop = function (body) {
    for (;;) {
        var result = apply(body, list());
        if (isTrue(result)) {
            return result;
        }
    }
};

var argv = function () {
    var consArgs = function (args, list) {
        if (!args.length) {
            return list;
        }
        return consArgs(args.slice(0, args.length - 1),
            cons(args[args.length - 1], list));
    };
    return consArgs(process.argv, list());
};
var isEnv = function (env) {
    if (typeof env !== "object") {
        return false;
    }
    if (typeof env.bindings !== "object") {
        return false;
    }
    return true;
};
var mkenv = function () {
    return {bindings: {}};
};
var extendEnv = function (parent) {
    if (not(isEnv(parent))) {
        return error("extendEnv", "invalid environment", parent);
    }
    var env = mkenv();
    env.parent = parent;
    return env;
};
var lookupVar = function (variable, env) {
    if (!isEnv(env)) {
        return error("lookupVar", "unbound variable", variable);
    }
    var varName = symbolName(variable);
    if (varName in env.bindings) {
        return env.bindings[varName];
    }
    return lookupVar(variable, env.parent);
};
var setVar = function (variable, val, env) {
    if (!isEnv(env)) {
        return error("setVar", "unbound variable", variable);
    }
    var varName = symbolName(variable);
    if (varName in env.bindings) {
        env.bindings[varName] = val;
        return val;
    }
    return setVar(variable, val, env.parent);
};
var defineVar = function (variable, val, env) {
    if (!isEnv(env)) {
        return error("defineVar", "invalid environment", env);
    }
    env.bindings[symbolName(variable)] = val;
    return val;
};

// lists
var isNull = function (l) {
    return !!(l && l.constructor === Array && !l.length);
};
var cons = function (left, right) {
    return [left, right];
};
var isPair = function (exp) {
    return !!(exp && exp.constructor === Array && exp.length === 2);
};
var car = function (p) {
    if (not(isPair(p))) {
        error("car", "not a pair", p);
    }
    return p[0];
};
var cdr = function (p) {
    if (not(isPair(p))) {
        error("cdr", "not a pair", p);
    }
    return p[1];
};
// can move to scm after variadic arguments implemented
var list = function () {
    var args = Array.prototype.slice.call(arguments);
    if (!args.length) {
        return [];
    }
    return cons(args.shift(), list.apply(undefined, args));
};

// values

// numbers
var isNumber = function (exp) {
    return typeof exp === "number" && !Number.isNaN(exp);
};
var parseNumber = function (exp) {
    var num = parseFloat(exp);
    if (isNumber(num)) {
        return num;
    }
    return exp;
};
var add = function () {
    var args = Array.prototype.slice.call(arguments);
    if (!args.length) {
        return 0;
    }
    return args.shift() + add.apply(undefined, args);
};
var sub = function () {
    var args = Array.prototype.slice.call(arguments);
    if (!args.length) {
        return error("sub", "wrong arity", args.length);
    }
    if (args.length === 1) {
        if (not(isNumber(args[0]))) {
            return error("subs", "not a number", args[0]);
        }
        return 0 - args[0];
    }
    var sub = function (sum, args) {
        if (!args.length) {
            return sum;
        }
        return sub(sum - args.shift(), args);
    };
    return sub(args.shift(), args);
};
var mul = function () {
    var args = Array.prototype.slice.call(arguments);
    if (!args.length) {
        return 1;
    }
    return args.shift() * mul.apply(undefined, args);
};
var div = function () {
    var args = Array.prototype.slice.call(arguments);
    if (!args.length) {
        return error("div", "wrong arity", args.length);
    }
    if (args.length === 1) {
        return 1 / args.shift();
    }
    var div = function (d, args) {
        if (!args.length) {
            return d;
        }
        return div(d / args.shift(), args);
    };
    return div(1, args);
};
var mod = function (dnd, dvs) {
    if (not(isNumber(dnd))) {
        return error("mod", "not a number", dnd);
    }
    if (not(isNumber(dvs))) {
        return error("mod", "not a number", dvs);
    }
    return dnd % dvs;
};
var gt = function (left, right) {
    return left > right;
};
var lt = function (left, right) {
    return not(gte(left, right));
};
var gte = function (left, right) {
    return gt(left, right) || eq(left, right);
};
var lte = function (left, right) {
    return not(gt(left, right));
};

// strings
var isString = function (exp) {
    return typeof exp === "string";
};
var slen = function (s) {
    if (not(isString(s))) {
        error("slen", "not a string", s);
    }
    return s.length;
};
var sidx = function (s, expression) {
    if (not(isString(s))) {
        error("sidx", "not a string", s);
    }
    if (not(isString(expression))) {
        error("sidx", "not a string", expression);
    }
    var m = s.match(new RegExp(expression), "");
    if (!m) {
        return -1;
    }
    return m.index;
};
var subs = function (s, offset, count) {
    if (not(isString(s))) {
        error("subs", "not a string", s);
    }
    if (not(isNumber(offset))) {
        error("subs", "not a number", offset);
    }
    if (not(isNumber(count))) {
        error("subs", "not a number", count);
    }
    if (gt(0, count)) {
        return s.substr(offset);
    }
    return s.substr(offset, count);
};
var sreplace = function (s, expression, replace) {
    if (not(isString(s))) {
        error("sreplace", "not a string", s);
    }
    if (not(isString(expression))) {
        error("sreplace", "not a string", expression);
    }
    if (not(isString(replace))) {
        error("sreplace", "not a string", replace);
    }
    return s.replace(new RegExp(expression, "g"), replace);
};
var cats = function () {
    var args = Array.prototype.slice.call(arguments);
    if (!args.length) {
        return "";
    }
    var first = args.shift();
    if (not(isString(first))) {
        return error("cats", "not a string", first);
    }
    return first + cats.apply(undefined, args);
};
var sescapeChar = function (c) {
    switch (c) {
    case "\b": return "\\b";
    case "\t": return "\\t";
    case "\n": return "\\n";
    case "\v": return "\\v";
    case "\f": return "\\f";
    case "\r": return "\\r";
    case "\"": return "\\\"";
    case "\\": return "\\\\";
    default: return c;
    }
};
var isEscapeChar = function (c) {
    return sescapeChar(c) !== c;
};
var sescapedChar = function (charSymbol) {
    switch (charSymbol) {
    case "b": return "\b";
    case "t": return "\t";
    case "n": return "\n";
    case "v": return "\v";
    case "f": return "\f";
    case "r": return "\r";
    default: return charSymbol;
    }
};

// possible to migrate after string ports/builders are implemented
var sescape = function (s) {
    var sescapes = function (ss, ps) {
        if (!ss.length) {
            return ps;
        }
        ps[ps.length] = sescapeChar(ss.shift());
        return sescapes(ss, ps);
    };
    return sescapes(s.split(""), []).join("");
};

var sescapeSymbol = function (s) {
    var findEscapeChar = function (ss) {
        if (!ss.length) {
            return symbolName(s);
        }
        if (isEscapeChar(ss.shift())) {
            return "|" + symbolName(s) + "|";
        }
        return findEscapeChar(ss);
    };
    return findEscapeChar(symbolName(s).split(""));
};
var sunescape = function (s) {
    var sunescapes = function (ss, ps, escaped) {
        if (!ss.length) {
            if (escaped) {
                return error("sunescape", "invalid escape sequence", s);
            }
            return ps;
        }
        var c = ss.shift();
        if (escaped) {
            ps[ps.length] = sescapedChar(c);
            return sunescapes(ss, ps, false);
        }
        if (c === "\\") {
            return sunescapes(ss, ps, true);
        }
        ps[ps.length] = c;
        return sunescapes(ss, ps, false);
    };
    return sunescapes(s.split(""), [], false).join("");
};
var sprintq = function (exp, quoted, inList) {
    if (typeof exp === "undefined") {
        error("sprint", "unknown type", String(exp));
    }
    if (isNumber(exp)) {
        return String(exp);
    }
    if (isString(exp)) {
        return "\"" + sescape(exp) + "\"";
    }
    if (isSymbol(exp)) {
        if (quoted) {
            return sescapeSymbol(exp);
        }
        return "'" + sescapeSymbol(exp);
    }
    if (isNull(exp)) {
        if (quoted) {
            return "()";
        }
        return "'()";
    }
    if (isQuote(exp)) {
        var text = "???";
        try {
            text = quoteText(exp);
        } catch (error) {}
        return "'" + sprintq(text, true, false);
    }
    if (isPair(exp)) {
        var product = [];
        if (!quoted) {
            product[product.length] = "'(";
        } else if (!inList) {
            product[product.length] = "(";
        }
        if (inList) {
            product[product.length] = " ";
        }
        product[product.length] = sprintq(car(exp), true, false);
        if (isNull(cdr(exp))) {
            product[product.length] = ")";
        } else if (isPair(cdr(exp))) {
            product[product.length] = sprintq(cdr(exp), true, true);
        } else {
            product[product.length] = " . ";
            product[product.length] = sprintq(cdr(exp), true, false);
            product[product.length] = ")";
        }
        return product.join("");
    }
    return sescape(String(exp));
};
var sprint = function (exp) {
    return sprintq(exp, false, false);
};

// booleans
var isTrue = function (exp) {
    return exp !== false;
};
var not = function (exp) {
    if (exp === false) {
        return true;
    }
    return false;
};

// symbols
var symbol = function (name) {
    return [name];
};
var isSymbol = function (exp) {
    return !!(exp && exp.constructor === Array && exp.length === 1 && typeof exp[0] === "string");
};
var symbolName = function (symbol) {
    return symbol[0];
};
var symbolEq = function (left, right) {
    if (!isSymbol(left)) {
        return false;
    }
    if (!isSymbol(right)) {
        return false;
    }
    return symbolName(left) === symbolName(right);
};

// primitives
var primitiveEq = function (left, right) {
    return left === right;
};
var applyInJs = function (f, args) {
    var toJsArray = function (args) {
        if (isNull(args)) {
            return [];
        }
        return [car(args)].concat(toJsArray(cdr(args)));
    };
    return f.apply(undefined, toJsArray(args));
};

// -- lang

// lists
var caar = function (l) { return car(car(l)); };
var cadr = function (l) { return car(cdr(l)); };
var cdar = function (l) { return cdr(car(l)); };
var cddr = function (l) { return cdr(cdr(l)); };
var caaar = function (l) { return car(car(car(l))); };
var caadr = function (l) { return car(car(cdr(l))); };
var cadar = function (l) { return car(cdr(car(l))); };
var cdaar = function (l) { return cdr(car(car(l))); };
var caddr = function (l) { return car(cdr(cdr(l))); };
var cdadr = function (l) { return cdr(car(cdr(l))); };
var cddar = function (l) { return cdr(cdr(car(l))); };
var cdddr = function (l) { return cdr(cdr(cdr(l))); };
var caaaar = function (l) { return car(car(car(car(l)))); };
var caaadr = function (l) { return car(car(car(cdr(l)))); };
var caadar = function (l) { return car(car(cdr(car(l)))); };
var cadaar = function (l) { return car(cdr(car(car(l)))); };
var cdaaar = function (l) { return cdr(car(car(car(l)))); };
var caaddr = function (l) { return car(car(cdr(cdr(l)))); };
var cadadr = function (l) { return car(cdr(car(cdr(l)))); };
var cdaadr = function (l) { return cdr(car(car(cdr(l)))); };
var caddar = function (l) { return car(cdr(cdr(car(l)))); };
var cdadar = function (l) { return cdr(car(cdr(car(l)))); };
var cddaar = function (l) { return cdr(cdr(car(car(l)))); };
var cadddr = function (l) { return car(cdr(cdr(cdr(l)))); };
var cdaddr = function (l) { return cdr(car(cdr(cdr(l)))); };
var cddadr = function (l) { return cdr(cdr(car(cdr(l)))); };
var cdddar = function (l) { return cdr(cdr(cdr(car(l)))); };
var cddddr = function (l) { return cdr(cdr(cdr(cdr(l)))); };
var len = function (l) {
    if (isNull(l)) {
        return 0;
    }
    return add(1, len(cdr(l)));
};
var isList = function (exp) {
    if (isPair(exp)) {
        return isList(cdr(exp));
    }
    if (isNull(exp)) {
        return true;
    }
    return false;
};
var map = function (proc, l) {
    if (isNull(l)) {
        return list();
    }
    return cons(proc(car(l)), map(proc, cdr(l)));
};
var last = function (l) {
    if (isNull(cdr(l))) {
        return car(l);
    }
    return last(cdr(l));
};
var dropLast = function (l) {
    if (isNull(cdr(l))) {
        return list();
    }
    return cons(car(l), dropLast(cdr(l)));
};
var append = function (left, right) {
    if (isNull(left)) {
        return right;
    }
    return cons(car(left), append(cdr(left), right));
};
var reverse = function (l) {
    if (isNull(l)) {
        return list();
    }
    return append(reverse(cdr(l)), list(car(l)));
};
var deepReverse = function (l) {
    return map(function (i) {
        if (isPair(i)) {
            return deepReverse(i);
        }
        return i;
    }, reverse(l));
};

// eq
var eq = function (left, right) {
    if (isNull(left) && isNull(right)) {
        return true;
    }
    if (symbolEq(left, right)) {
        return true;
    }
    if (isQuote(left) && isQuote(right)) {
        return eq(len(left), 1) && eq(len(right), 1) ||
            eq(quoteText(left), quoteText(right));
    }
    return primitiveEq(left, right);
};
var isTaggedList = function (exp, tag) {
    if (not(isPair(exp))) {
        return false;
    }
    return eq(car(exp), tag);
};
var isQuote = function (exp) {
    if (not(isTaggedList(exp, symbol("quote")))) {
        return false;
    }
    return true;
};
var quote = function (exp) {
    return list(symbol("quote"), exp);
};
var quoteText = function (exp) {
    return cadr(exp);
};
var validateQuote = function (exp) {
    if (not(isTaggedList(exp, symbol("quote"))) ||
        not(eq(len(exp), 2))) {
        return error("validateQuote", "invalid quotation", exp);
    }
    return noPrint;
};

// environment
var extendEnvironment = function (vars, vals, baseEnv) {
    if (eq(len(vars), len(vals))) {
        var env = extendEnv(baseEnv);
        var defineVars = function (vars, vals) {
            if (isNull(vars)) {
                return env;
            }
            defineVar(car(vars), car(vals), env);
            return defineVars(cdr(vars), cdr(vals));
        };
        return defineVars(vars, vals);
    }
    if (gt(len(vars), len(vals))) {
        return error("extendEnvironment", "too few arguments supplied", list(vars, vals));
    }
    return error("extendEnvironment", "too many arguments suppplied", list(vars, vals));
};
var lookupVariableValue = function (variable, env) {
    return lookupVar(variable, env);
};

// primitive values
var isSelfEvaluating = function (exp) {
    return exp === true ||
        exp === false ||
        isNumber(exp) ||
        isString(exp);
};

// assignment
var isVariable = function (exp) {
    return isSymbol(exp);
};
var isAssignment = function (exp) {
    return isTaggedList(exp, symbol("set!"));
};
var validateAssignment = function (exp) {
    if (not(eq(len(exp), 3))) {
        error("validateAssignment", "invalid arity", exp);
    }
    if (not(isSymbol(assignmentVariable(exp)))) {
        error("validateAssignment", "invalid variable name", exp);
    }
    return noPrint;
};
var assignmentVariable = function (exp) {
    return cadr(exp);
};
var assignmentValue = function (exp) {
    return caddr(exp);
};
var setVariableValue = function (variable, val, env) {
    return setVar(variable, val, env);
};
var evalAssignment = function (exp, env) {
    return setVariableValue(assignmentVariable(exp), eval(assignmentValue(exp), env), env);
};

// lambda
var isLambda = function (exp) {
    return isTaggedList(exp, symbol("lambda"));
};
var validateLambda = function (exp) {
    if (gt(3, len(exp))) {
        error("valdiateLambda", "invalid arity", exp);
    }
    if (not(isList(cadr(exp)))) {
        error("validateLambda", "invalid argument list", exp);
    }
    return noPrint;
};
var lambdaParameters = function (exp) {
    return cadr(exp);
};
var lambdaBody = function (exp) {
    return cddr(exp);
};
var makeLambda = function (parameters, body) {
    return cons(symbol("lambda"), cons(parameters, body));
};
var makeProcedure = function (parameters, body, env) {
    return list(symbol("procedure"), parameters, body, env);
};

// definition
var evalDefinition = function (exp, env) {
    return defineVariable(definitionVariable(exp), eval(definitionValue(exp), env), env);
};
var isDefinition = function (exp) {
    return isTaggedList(exp, symbol("define"));
};
var validateDefinition = function (exp) {
    if (not(isSymbol(cadr(exp)) && eq(len(exp), 3) ||
        isList(cadr(exp)) && gt(len(exp), 2))) {
        error("validateDefinition", "invalid format", exp);
    }
    if (not(isSymbol(definitionVariable(exp)))) {
        error("validateDefinition", "invalid variable name", exp);
    }
    return noPrint;
};
var definitionVariable = function (exp) {
    if (isSymbol(cadr(exp))) {
        return cadr(exp);
    }
    return caadr(exp);
};
var definitionValue = function (exp) {
    if (isSymbol(cadr(exp))) {
        return caddr(exp);
    }
    return makeLambda(cdadr(exp), cddr(exp));
};
var defineVariable = function (variable, val, env) {
    return defineVar(variable, val, env);
};

// if
var isIf = function (exp) {
    return isTaggedList(exp, symbol("if"));
};
var validateIf = function (exp) {
    if (not(eq(len(exp), 4))) {
        error("validateIf", "invalid arity", exp);
    }
    return noPrint;
};
var ifPredicate = function (exp) {
    return cadr(exp);
};
var ifConsequent = function (exp) {
    return caddr(exp);
};
var ifAlternative = function (exp) {
    return cadddr(exp);
};
var makeIf = function (predicate, consequent, alternative) {
    return list(symbol("if"), predicate, consequent, alternative);
};
var evalIf = function (exp, env) {
    if (isTrue(eval(ifPredicate(exp), env))) {
        return eval(ifConsequent(exp), env);
    }
    return eval(ifAlternative(exp), env);
};

// begin
var isBegin = function (exp) {
    return isTaggedList(exp, symbol("begin"));
};
var validateBegin = function (exp) {
    if (gt(2, len(exp))) {
        return error("validateBegin", "invalid arity", exp);
    }
    return noPrint;
};
var beginActions = function (exp) {
    return cdr(exp);
};
var makeBegin = function (seq) {
    return cons(symbol("begin"), seq);
};

// sequence
var isLastExp = function (exp) {
    return isNull(cdr(exp));
};
var firstExp = function (exps) {
    return car(exps);
};
var restExps = function (exps) {
    return cdr(exps);
};
var evalSequence = function (exps, env) {
    if (isNull(exps)) {
        error("evalSequence", "unspecified sequence value", exps);
    }
    if (isLastExp(exps)) {
        return eval(firstExp(exps), env);
    }
    eval(firstExp(exps), env);
    return evalSequence(restExps(exps), env);
};
var sequenceToExp = function (seq) {
    if (isNull(seq)) {
        return seq;
    }
    if (isLastExp(seq)) {
        return firstExp(seq);
    }
    return makeBegin(seq);
};

// cond
var isCond = function (exp) {
    return isTaggedList(exp, symbol("cond"));
};
var condClauses = function (exp) {
    return cdr(exp);
};
var condPredicate = function (clause) {
    return car(clause);
};
var isCondElseClause = function (clause) {
    return eq(condPredicate(clause), symbol("else"));
};
var condActions = function (clause) {
    return cdr(clause);
};
var expandClauses = function (clauses) {
    if (isNull(clauses)) {
        return symbol("true");
    }
    var first = car(clauses);
    if (not(isPair(first))) {
        error("expandClauses", "invalid syntax", first);
    }
    var theRest = cdr(clauses);
    if (isCondElseClause(first)) {
        if (isNull(theRest)) {
            return sequenceToExp(condActions(first));
        }
        error("expandClauses", "else clause isn't last", clauses);
    }
    return makeIf(condPredicate(first), sequenceToExp(condActions(first)), expandClauses(theRest));
};
var condToIf = function (exp) {
    return expandClauses(condClauses(exp));
};

// and
var isAnd = function (exp) {
    return isTaggedList(exp, symbol("and"));
};
var andExpressions = function (exp) {
    return cdr(exp);
};
var expandAnd = function (expressions) {
    if (isNull(expressions)) {
        return symbol("true");
    }
    var first = car(expressions);
    var theRest = cdr(expressions);
    if (isNull(theRest)) {
        return first;
    }
    return makeIf(first, expandAnd(theRest), first);
};
var andToIf = function (exp) {
    return expandAnd(andExpressions(exp));
};

// or
var isOr = function (exp) {
    return isTaggedList(exp, symbol("or"));
};
var orExpressions = function (exp) {
    return cdr(exp);
};
var expandOr = function (expressions) {
    if (isNull(expressions)) {
        return symbol("false");
    }
    var first = car(expressions);
    var theRest = cdr(expressions);
    if (isNull(theRest)) {
        return first;
    }
    return makeIf(first, first, expandOr(theRest));
};
var orToIf = function (exp) {
    return expandOr(orExpressions(exp));
};

// let
var isLet = function (exp) {
    return isTaggedList(exp, symbol("let"));
};
var validateLet = function (exp) {
    if (lt(len(exp), 3)) {
        return error("validateLet", "invalid arity", exp);
    }
    if (not(isPair(cadr(exp)))) {
        return error("validateLet", "invalid syntax", exp);
    }
    return noPrint;
};
var letDefs = function (exp) {
    return cadr(exp);
};
var letVariables = function (defs) {
    return map(car, defs);
};
var letValues = function (defs) {
    return map(cadr, defs);
};
var letBody = function (exp) {
    return cddr(exp);
};
var letToApplication = function (exp) {
    return cons(makeLambda(letVariables(letDefs(exp)), letBody(exp)), letValues(letDefs(exp)));
};

// application
var isApplication = function (exp) {
    if (not(isPair(exp))) {
        return false;
    }
    if (gt(1, len(exp))) {
        return false;
    }
    return true;
};
var operator = function (exp) {
    return car(exp);
};
var operands = function (exp) {
    return cdr(exp);
};
var firstOperand = function (ops) {
    return car(ops);
};
var restOperands = function (ops) {
    return cdr(ops);
};
var isCompoundProcedure = function (p) {
    return isTaggedList(p, symbol("procedure"));
};
var procedureParameters = function (p) {
    return cadr(p);
};
var procedureEnvironment = function (p) {
    return cadddr(p);
};
var procedureBody = function (p) {
    return caddr(p);
};
var hasNoOperands = function (ops) {
    return isNull(ops);
};
var listOfValues = function (exps, env) {
    if (hasNoOperands(exps)) {
        return list();
    }
    return cons(eval(firstOperand(exps), env), listOfValues(restOperands(exps), env));
};

// primitive procedures
var isPrimitiveProcedure = function (proc) {
    if (not(isTaggedList(proc, symbol("primitive")))) {
        return false;
    }
    return true;
};
var primitiveProcedureNames = function () {
    return map(car, primitiveProcedures());
};
var primitiveProcedureObjects = function () {
    return map(function (proc) { return list(symbol("primitive"), cadr(proc)); }, primitiveProcedures());
};
var primitiveImplementation = function (proc) {
    // hack for the higher order primitive functions
    if (isPrimitiveProcedure(proc)) {
        return primitiveImplementation(cadr(proc));
    }
    return proc;
};
var applyPrimitiveProcedure = function (proc, args) {
    return applyInJs(primitiveImplementation(proc), args);
};

// eval/apply
var apply = function (procedure, arguments) {
    if (isPrimitiveProcedure(procedure)) {
        return applyPrimitiveProcedure(procedure, arguments);
    }
    if (isCompoundProcedure(procedure)) {
        return evalSequence(
            procedureBody(procedure),
            extendEnvironment(
                procedureParameters(procedure),
                arguments,
                procedureEnvironment(procedure)));
    }
    error("apply", "unknown procedure type", procedure);
};
var eval = function (exp, env) {
    if (isSelfEvaluating(exp)) {
        return exp;
    }
    if (isVariable(exp)) {
        return lookupVariableValue(exp, env);
    }
    if (isQuote(exp)) {
        validateQuote(exp);
        return quoteText(exp);
    }
    if (isAssignment(exp)) {
        validateAssignment(exp);
        return evalAssignment(exp, env);
    }
    if (isDefinition(exp)) {
        validateDefinition(exp);
        return evalDefinition(exp, env);
    }
    if (isIf(exp)) {
        validateIf(exp);
        return evalIf(exp, env);
    }
    if (isLambda(exp)) {
        validateLambda(exp);
        return makeProcedure(lambdaParameters(exp), lambdaBody(exp), env);
    }
    if (isBegin(exp)) {
        validateBegin(exp);
        return evalSequence(beginActions(exp), env);
    }
    if (isCond(exp)) {
        return eval(condToIf(exp), env);
    }
    if (isAnd(exp)) {
        return eval(andToIf(exp), env);
    }
    if (isOr(exp)) {
        return eval(orToIf(exp), env);
    }
    if (isLet(exp)) {
        validateLet(exp);
        return eval(letToApplication(exp), env);
    }
    if (isApplication(exp)) {
        return apply(eval(operator(exp), env),
            listOfValues(operands(exp), env));
    }
    error("eval", "unknown expression type", exp);
};
var evalShare = function (exp) {
    return eval(exp, setupEnvironment(mkenv()));
};

// setup
var setupEnvironment = function (empty) {
    var initialEnv = extendEnvironment(
        primitiveProcedureNames(),
        primitiveProcedureObjects(),
        empty);
    defineVariable(symbol("true"), true, initialEnv);
    defineVariable(symbol("false"), false, initialEnv);
    defineVariable(symbol("no-print"), noPrint, initialEnv);
    return initialEnv;
};
var createEnvironment = function () {
    return setupEnvironment(mkenv());
};
var primitiveProcedures = function () {
    return list(
        list(symbol("eval"), evalShare),
        list(symbol("apply"), apply),
        list(symbol("try"), tryc),
        list(symbol("for"), forLoop),
        list(symbol("read-file"), readFile),
        list(symbol("symbol?"), isSymbol),
        list(symbol("symbol-name"), symbolName),
        list(symbol("string"), String),
        list(symbol("string->symbol"), symbol),
        list(symbol("primitive-eq?"), primitiveEq),
        list(symbol("make-env"), mkenv),
        list(symbol("define-var"), defineVar),
        list(symbol("lookup-var"), lookupVar),
        list(symbol("set-var"), setVar),
        list(symbol("extend-env"), extendEnv),
        list(symbol("car"), car),
        list(symbol("cdr"), cdr),
        list(symbol("cons"), cons),
        list(symbol("pair?"), isPair),
        list(symbol("null?"), isNull),
        list(symbol("list"), list),
        list(symbol("eq?"), eq),
        list(symbol("number?"), isNumber),
        list(symbol("parse-number"), parseNumber),
        list(symbol("+"), add),
        list(symbol("-"), sub),
        list(symbol("*"), mul),
        list(symbol("/"), div),
        list(symbol("%"), mod),
        list(symbol(">"), gt),
        list(symbol("<"), lt),
        list(symbol(">="), gte),
        list(symbol("<="), lte),
        list(symbol("string?"), isString),
        list(symbol("slen"), slen),
        list(symbol("sidx"), sidx),
        list(symbol("subs"), subs),
        list(symbol("sreplace"), sreplace),
        list(symbol("cats"), cats),
        list(symbol("sescape"), sescape),
        list(symbol("sunescape"), sunescape),
        list(symbol("sprint"), sprint),
        list(symbol("sescaped-char?"), sescapedChar),
        list(symbol("not"), not),
        list(symbol("error"), error),
        list(symbol("out"), out),
        list(symbol("error?"), isError),
        list(symbol("sprint-error"), sprintError),
        list(symbol("sprint-stack"), sprintStack),
        list(symbol("print"), print),
        list(symbol("display"), display),
        list(symbol("exit"), exit),
        list(symbol("proc-argv"), argv),
        list(symbol("sread"), sreadFull),
        list(symbol("caar"), caar),
        list(symbol("cadr"), cadr),
        list(symbol("cdar"), cdar),
        list(symbol("cddr"), cddr),
        list(symbol("caaar"), caaar),
        list(symbol("caadr"), caadr),
        list(symbol("cadar"), cadar),
        list(symbol("cdaar"), cdaar),
        list(symbol("caddr"), caddr),
        list(symbol("cdadr"), cdadr),
        list(symbol("cddar"), cddar),
        list(symbol("cdddr"), cdddr),
        list(symbol("caaaar"), caaaar),
        list(symbol("caaadr"), caaadr),
        list(symbol("caadar"), caadar),
        list(symbol("cadaar"), cadaar),
        list(symbol("cdaaar"), cdaaar),
        list(symbol("caaddr"), caaddr),
        list(symbol("cadadr"), cadadr),
        list(symbol("cdaadr"), cdaadr),
        list(symbol("caddar"), caddar),
        list(symbol("cdadar"), cdadar),
        list(symbol("cddaar"), cddaar),
        list(symbol("cadddr"), cadddr),
        list(symbol("cdaddr"), cdaddr),
        list(symbol("cddadr"), cddadr),
        list(symbol("cdddar"), cdddar),
        list(symbol("cddddr"), cddddr));
};
var mikkamakka = createEnvironment();
defineVar(symbol("mikkamakka"), mikkamakka, mikkamakka);

// reader
var mkstate = function () {
    var top = list();
    return list(list(), top, top, list(symbol("list")));
}
var hasExpression = function (state) {
    return not(isNull(car(state)));
};
var hasUncompleteExpression = function (state) {
    return not(isNull(caddr(state)));
};
var nextExpression = function (state) {
    return caar(state);
};
var dropExpression = function (state) {
    return list(cdar(state), cadr(state), caddr(state), cadddr(state));
};
var closeState = function (state) {
    return sread(state, symbol("eof"));
};
var sread = function (state, exp) {
    // state
    var replaceLevel = function (state, level) {
        var findAndReplace = function (target, current) {
            if (eq(target, current)) {
                return list(level, level);
            }
            var next = findAndReplace(target, car(current));
            return list(car(next), cons(cadr(next), cdr(current)));
        };
        var next = findAndReplace(cadr(state), caddr(state));
        return list(car(state), car(next), cadr(next), cadddr(state));
    };
    var addToLevel = function (state, tag) {
        return replaceLevel(state, cons(tag, cadr(state)));
    };
    var isExpressionToComplete = function (state) {
        return eq(cadr(state), caddr(state));
    };
    var completeExpression = function (state) {
        var top = list();
        var exp = tag(state);
        return list(append(car(state), deepReverse(list(exp))), top, top, cadddr(state));
    };
    var isTag = function (state) {
        return not(isNull(cadr(state))) &&
            not(isPair(caadr(state)));
    };
    var tag = function (state) {
        return caadr(state);
    };
    var startTag = function (state) {
        return addToLevel(state, "");
    };
    var replaceTag = function (state, tag) {
        return replaceLevel(state, cons(tag, cdadr(state)));
    };
    var appendToTag = function (state, char) {
        return replaceTag(state, cats(tag(state), char));
    };
    var applyTagType = function (tag) {
        if (isPair(tag)) {
            return tag;
        }
        var num = parseNumber(tag);
        if (isNumber(num)) {
            return num;
        }
        if (gt(slen(tag), 1) &&
            eq(subs(tag, 0, 1), "\"") &&
            eq(subs(tag, sub(slen(tag), 1), 1), "\"")) {
            return subs(tag, 1, sub(slen(tag), 2));
        }
        return symbol(tag);
    };
    var finishTag = function (state) {
        var state = replaceTag(state, applyTagType(tag(state)));
        if (isExpressionToComplete(state)) {
            return completeExpression(state);
        }
        return state;
    };
    var pushLevel = function (state) {
        var next = addToLevel(state, list());
        return list(car(state), caadr(next), caddr(next), cadddr(state));
    };
    var popLevel = function (state) {
        var findParent = function (target, current) {
            if (eq(target, current)) {
                return error("sread", "unexpected closing paren", exp);
            }
            if (eq(target, car(current))) {
                return current;
            }
            return findParent(target, car(current));
        };
        var newState = list(car(state),
            findParent(cadr(state), caddr(state)),
            caddr(state),
            cadddr(state));
        if (isExpressionToComplete(newState)) {
            return completeExpression(newState);
        }
        return newState;
    };
    var readState = function (state) {
        return car(cadddr(state));
    };
    var pushReadState = function (state, nextState) {
        return list(car(state), cadr(state), caddr(state), cons(nextState, cadddr(state)));
    };
    var popReadState = function (state) {
        if (isNull(cdr(cadddr(state)))) {
            return error("sread", "unexpected closing paren", exp);
        }
        return list(car(state), cadr(state), caddr(state), cdr(cadddr(state)));
    };
    var finishRead = function (state) {
        if (isInTag(state)) {
            state = popReadState(finishTag(state));
        }
        state = popQuote(state);
        if (not(isExpressionToComplete(state))) {
            return error("sread", "invalid input", exp);
        }
        return state;
    };

    // read state
    var isInReadState = function (state, stateFlag) {
        return eq(readState(state), stateFlag);
    };
    var isInComment = function (state) {
        return isInReadState(state, symbol("comment"));
    };
    var pushComment = function (state) {
        return pushReadState(state, symbol("comment"));
    };
    var isInRangeEscape = function (state) {
        return isInReadState(state, symbol("range-escape"));
    };
    var pushRangeEscape = function (state) {
        return pushReadState(state, symbol("range-escape"));
    };
    var isInSingleEscape = function (state) {
        return isInReadState(state, symbol("single-escape"));
    };
    var pushSingleEscape = function (state) {
        return pushReadState(state, symbol("single-escape"));
    };
    var isInString = function (state) {
        return isInReadState(state, symbol("string"));
    };
    var pushString = function (state) {
        return pushReadState(state, symbol("string"));
    };
    var isInTag = function (state) {
        return isInReadState(state, symbol("tag"));
    };
    var pushTag = function (state) {
        return pushReadState(state, symbol("tag"));
    };
    var isInQuote = function (state) {
        return isInReadState(state, symbol("quote"));
    };
    var pushQuote = function (state) {
        return pushReadState(state, symbol("quote"));
    };
    var popQuote = function (state) {
        if (not(isInQuote(state))) {
            return state;
        }
        return popQuote(popLevel(popReadState(state)));
    };
    var isInList = function (state) {
        return isInReadState(state, symbol("list"));
    };
    var pushList = function (state) {
        return pushReadState(state, symbol("list"));
    };

    // read char
    var isStartComment = function (char) {
        return eq(char, ";");
    };
    var isNewLine = function (char) {
        return eq(char, "\n");
    };
    var isWhiteSpace = function (char) {
        return eq(sidx(char, "\\s"), 0);
    };
    var isSingleEscape = function (char) {
        return eq(char, "\\");
    };
    var isRangeEscape = function (char) {
        return eq(char, "|");
    };
    var isStringDelimiter = function (char) {
        return eq(char, "\"");
    };
    var isQuoteChar = function (char) {
        return eq(char, "'");
    };
    var isListOpen = function (char) {
        return eq(char, "(");
    };
    var isListClose = function (char) {
        return eq(char, ")");
    };
    var readChar = function (state, char) {
        if (isInComment(state)) {
            if (isNewLine(char)) {
                return popReadState(state);
            }
            return state;
        }

        if (isInSingleEscape(state)) {
            return appendToTag(popReadState(state), sescapedChar(char));
        }
        
        if (isInRangeEscape(state)) {
            if (isSingleEscape(char)) {
                return pushSingleEscape(state);
            }
            if (isRangeEscape(char)) {
                return popReadState(state);
            }
            return appendToTag(state, char);
        }

        if (isInString(state)) {
            if (isSingleEscape(char)) {
                return pushSingleEscape(state);
            }
            if (isStringDelimiter(char)) {
                state = finishTag(appendToTag(popReadState(state), char));
                if (isInQuote(state)) {
                    return popQuote(state);
                }
                return state;
            }
            return appendToTag(state, char);
        }

        if (isInTag(state)) {
            if (isStartComment(char)) {
                state = popReadState(finishTag(state));
                if (isInQuote(state)) {
                    state = popQuote(state);
                }
                return pushComment(state);
            }
            if (isWhiteSpace(char)) {
                state = popReadState(finishTag(state));
                if (isInQuote(state)) {
                    return popQuote(state);
                }
                return state;
            }
            if (isSingleEscape(char)) {
                return pushSingleEscape(state);
            }
            if (isRangeEscape(char)) {
                return pushRangeEscape(state);
            }
            if (isStringDelimiter(char)) {
                state = popReadState(finishTag(state));
                if (isInQuote(state)) {
                    state = popQuote(state);
                }
                return pushString(appendToTag(startTag(state), char));
            }
            if (isQuoteChar(char)) {
                state = popReadState(finishTag(state));
                if (isInQuote(state)) {
                    state = popQuote(state);
                }
                return pushQuote(
                    finishTag(replaceTag(startTag(pushLevel(state)), "quote")));
            }
            if (isListOpen(char)) {
                state = popReadState(finishTag(state));
                if (isInQuote(state)) {
                    state = popQuote(state);
                }
                return pushList(pushLevel(state));
            }
            if (isListClose(char)) {
                state = popReadState(finishTag(state));
                if (isInQuote(state)) {
                    state = popQuote(state);
                }
                state = popLevel(popReadState(state));
                if (isInQuote(state)) {
                    return popQuote(state);
                }
                return state;
            }
            return appendToTag(state, char);
        }

        if (isInQuote(state)) {
            if (isStartComment(char)) {
                return pushComment(state);
            }
            if (isWhiteSpace(char)) {
                return state;
            }
            if (isSingleEscape(char)) {
                return pushSingleEscape(pushTag(startTag(state)));
            }
            if (isRangeEscape(char)) {
                return pushRangeEscape(pushTag(startTag(state)));
            }
            if (isStringDelimiter(char)) {
                return pushString(appendToTag(startTag(state), char));
            }
            if (isQuoteChar(char)) {
                return pushQuote(
                    finishTag(
                    replaceTag(startTag(pushLevel(state)), "quote")));
            }
            if (isListOpen(char)) {
                return pushList(pushLevel(state));
            }
            if (isListClose(char)) {
                return error("sread", "invalid quotation", exp);

            }
            return pushTag(appendToTag(startTag(state), char));
        }

        if (isInList(state)) {
            if (isStartComment(char)) {
                return pushComment(state);
            }
            if (isWhiteSpace(char)) {
                return state;
            }
            if (isSingleEscape(char)) {
                return pushSingleEscape(pushTag(startTag(state)));
            }
            if (isRangeEscape(char)) {
                return pushRangeEscape(pushTag(startTag(state)));
            }
            if (isStringDelimiter(char)) {
                return pushString(appendToTag(startTag(state), char));
            }
            if (isQuoteChar(char)) {
                return pushQuote(
                    finishTag(
                    replaceTag(startTag(pushLevel(state)), "quote")));
            }
            if (isListOpen(char)) {
                return pushList(pushLevel(state));
            }
            if (isListClose(char)) {
                return popQuote(popReadState(popLevel(state)));
            }
            return pushTag(appendToTag(startTag(state), char));
        }

        error("sread", "unrecognized reader state", readState(state));
    };

    // tail recursion
    for (;;) {
        if (not(isString(exp)) && not(eq(exp, symbol("eof")))) {
            error("sread", "invalid expression", exp);
        }
        if (eq(exp, symbol("eof"))) {
            return finishRead(state);
        }
        if (eq(slen(exp), 0)) {
            return state;
        }
        state = readChar(state, subs(exp, 0, 1))
        exp = subs(exp, 1, -1);
    }
};
var sreadFull = function (str) {
    var state = mkstate();
    state = sread(state, str);
    state = closeState(state);
    if (hasExpression(state)) {
        return nextExpression(state);
    }
    if (hasUncompleteExpression(state)) {
        return error("sreadFull", "invalid expression", str);
    }
    return symbol("nothing");
};

// commands
var isCommand = function (arg) {
    switch (arg) {
    case "repl":
    case "run":
    case "test":
        return true;
    }
    return false;
};
var getCommand = function () {
    if (process.argv.length <= 2) {
        return "repl";
    }
    if (isCommand(process.argv[2])) {
        return process.argv[2];
    }
    return "run";
};
var getArg = function () {
    var checkArg = function (argIndex) {
        if (process.argv.length <= argIndex) {
            return error("getArg", "missing argument", getCommand());
        }
    }
    var argIndex = 2;
    checkArg(argIndex);
    if (isCommand(process.argv[argIndex])) {
        argIndex = 3;
    }
    checkArg(argIndex);
    return process.argv[argIndex];
};
var repl = function () {
    var theGlobalEnvironment = mikkamakka;
    var readState = mkstate();
    var readStdin = function (clb) {
        var rl = require("readline").createInterface(process.stdin, process.stdout);
        rl.on("line", clb);
        return rl;
    };
    var processLine = function (l) {
        var fetchAndEval = function (state) {
            if (not(hasExpression(state))) {
                return state;
            };
            var exp = nextExpression(state);
            try {
                var out = eval(exp, theGlobalEnvironment);
                print(out);
            } catch (error) {
                display(error);
            }
            state = dropExpression(state);
            return fetchAndEval(state);
        };
        readState = sread(readState, cats(l, "\n"));
        readState = fetchAndEval(readState);
        if (hasUncompleteExpression(readState)) {
            rl.setPrompt(". ");
        } else {
            rl.setPrompt("> ");
        }
        rl.prompt();
    };
    var rl = readStdin(processLine);
    rl.setPrompt("> ");
    rl.prompt();
};
var run = function () {
    var readScm = function (f, clb) {
        require("fs").readFile(f, function (err, scm) {
            if (err) {
                error("run", "error while reading a file", err);
            }
            clb(scm.toString());
        });
    };
    var readAndEval = function (scm) {
        var theGlobalEnvironment = mikkamakka;
        var fetchAndEval = function (state) {
            if (not(hasExpression(state))) {
                return state;
            }
            var exp = nextExpression(state);
            eval(exp, theGlobalEnvironment);
            state = dropExpression(state);
            return fetchAndEval(state);
        };
        var state = mkstate();
        state = sread(state, scm);
        state = closeState(state);
        fetchAndEval(state);
    };
    readScm(getArg(), readAndEval);
};

switch (getCommand()) {
case "run":
    run();
    break;
case "repl":
    repl();
    break;
}

// test
var runTests = function () {
    var lastTestName = false;
    var stringContains = function (o, p) {
        return String(o).indexOf(p) >= 0;
    };
    var listEquals = function (left, right) {
        if (left.constructor !== Array ||
            right.constructor !== Array) {
            return left === right;
        }
        if (left.length !== right.length) {
            return false;
        }
        if (!left.length) {
            return true;
        }
        if (left.length === 1) {
            return listEquals(left[0], right[0]);
        }
        if (left.length !== 2) {
            return false;
        }
        if (!listEquals(left[0], right[0])) {
            return false;
        }
        return listEquals(left[1], right[1]);
    };
    var assert = function (exp, msg) {
        if (exp) {
            return exp;
        }
        console.log(
            (lastTestName && (lastTestName + ": ") || "") +
            (msg && ("failed: " + msg) || "failed"));
        return process.exit(-1);
    };
    var fail = function (msg, f) {
        try {
            f();
            assert(false, msg);
        } catch (error) {
            if (!stringContains(error, msg)) {
                console.log(error.stack);
                assert(false, msg);
            }
        }
    };
    var test = function () {
        var name = arguments.length > 1 && arguments[0] || null;
        var theGlobalEnvironment = mikkamakka;
        lastTestName = name;
        try {
            (name && arguments[1] || arguments[0])(function (sexp) {
                if (not(isString(sexp))) {
                    return error("mikkamakka", "invalid expression", "??");
                }
                var state = mkstate();
                state = sread(state, sexp);
                state = closeState(state);
                if (not(hasExpression(state))) {
                    return error("mikkamakka", "no expression to evaluate", sexp);
                }
                var exp = nextExpression(state);
                state = dropExpression(state);
                if (hasExpression(state)) {
                    return error("mikkamakka", "currently not supporting more than one expression", sexp);
                }
                return eval(exp, theGlobalEnvironment);
            });
        } catch (error) {
            console.log((name && (name + ": ") || "") + "error during test:");
            console.log(error.stack);
            process.exit(-1);
        }
    };

    test("test no input", function (mm) {
        fail("invalid expression", function () {
            mm();
        });

        fail("invalid expression", function () {
            mm(1);
        });

        fail("no expression", function () {
            mm("");
        });

        fail("no expression", function () {
            mm(" ");
        });

        fail("no expression", function () {
            mm("\n");
        });
    });

    test("false, true", function (mm) {
        assert(mm("false") === false);
        assert(mm("true") === true);
    });

    test("self evaluation", function (mm) {
        assert(mm("1") === 1);
        assert(mm("\"some string\"") === "some string");
        fail("no expression", function () {
            mm("\"some unclosed string");
        });
    });

    test("variable", function (mm) {
        fail("unbound variable", function () {
            mm("a");
        });
        mm("(define a 1)");
        assert(mm("a") === 1, "bound variable");
    });

    test("quote", function (mm) {
        fail("invalid quotation", function () {
            mm("'");
        });
        fail("invalid quotation", function () {
            mm("(quote)");
        });
        assert(mm("'1") === 1);
        assert(mm("' 1") === 1);
        assert(mm("(quote 1)") === 1);
        assert(mm("'\"string\"") === "string");
        assert(mm("(quote \"string\")") === "string");
        assert(listEquals(mm("'a"), ["a"]));
        assert(listEquals(mm("(quote a)"), ["a"]));
        assert(listEquals(mm("'()"), []));
        assert(listEquals(mm("(quote ())"), []));
        assert(listEquals(mm("'(a b c)"), [["a"], [["b"], [["c"], []]]]));
        assert(listEquals(mm("(quote (a b c))"), [["a"], [["b"], [["c"], []]]]));
        assert(listEquals(mm("'('a 'b)"), [[["quote"], [["a"], []]], [[["quote"], [["b"], []]], []]]));
        assert(listEquals(mm("''"), [["quote"], []]));
        assert(listEquals(mm("'(quote)"), [["quote"], []]));
        fail("invalid quotation", function () {
            mm("(quote ')");
        });
        assert(listEquals(mm("(quote (quote))"), [["quote"], []]));
    });

    test("assignment", function (mm) {
        mm("(define a 1)");
        mm("(set! a 2)");
        assert(mm("a") === 2);
    });

    test("define", function (mm) {
        mm("(define a 1)");
        assert(mm("a") === 1, "variable, check value");
        mm("(define (a) 1)");
        assert(mm("(a)") === 1, "procedure, check return value");
        (function () {
            mm("(define (a) (begin))");
            assert(true, "empty begin is legal here");
        })();
        fail("invalid arity", function () {
            mm("(a)");
        });
    });

    test("if", function (mm) {
        assert(mm("(if false 1 0)") === 0, "if, false");
        assert(mm("(if true 1 0)") === 1, "if, true");
        assert(mm("(if \"anything else\" 1 0)") === 1, "if, anything");
    });

    test("lambda", function (mm) {
        mm("(define a (lambda () 1))");
        assert(mm("(a)") === 1);
    });

    test("begin", function (mm) {
        assert(mm("(begin (define a 1) 2)") === 2, "return value");
        assert(mm("a") === 1, "side effect");
        fail("arity", function () {
            mm("(begin)");
        });
    });

    test("cond", function (mm) {
        assert(mm("(cond (true 0) (false 1) (else 2))") === 0, "first");
        assert(mm("(cond (false 0) (true 1) (else 2))") === 1, "second");
        assert(mm("(cond (false 0) (false 1) (else 2))") === 2, "else");
        assert(mm("(cond)"), "no clauses");
        assert(mm("(cond (false))"), "implicit else clause");
        fail("invalid syntax", function () {
            mm("(cond 1)");
        });
    });

    test("application", function (mm) {
        assert(mm("((lambda () 1))") === 1, "lambda");
        mm("(define (a) 1)");
        assert(mm("(a)") === 1, "procedure from shortcut");
        mm("(define a (lambda () 2))");
        assert(mm("(a)") === 2, "procedure from variable");
        mm("(define (a x) (if x \"it's true\" \"it's false\"))");
        assert(mm("(a true)") === "it's true", "with args, true");
        assert(mm("(a false)") === "it's false", "with args, false");
        assert(mm("(a \"just anything\")") === "it's true", "with args, anything");
        assert(listEquals(mm("(cons 1 2)"), [1, 2]), "cons");
        mm("(define (a) (begin))");
        fail("arity", function () {
            mm("(a)");
        });
    });

    test("escape", function (mm) {
        assert(sescape("\b\t\v\f\n\r\\\"") === "\\b\\t\\v\\f\\n\\r\\\\\\\"");
        assert(mm("\"\b\t\v\f\n\r\\\\\\\"\"") === "\b\t\v\f\n\r\\\"");
        assert(sunescape("\\b\\t\\v\\f\\n\\r\\\\\\\"") === "\b\t\v\f\n\r\\\"");
        assert(mm("\"\\b\\t\\v\\f\\n\\r\\\\\\\"\"") === "\b\t\v\f\n\r\\\"");
    });

    test("sprint", function (mm) {
        assert(sprint(mm("1")) === "1");
        assert(sprint(mm("\"some string\"")) === "\"some string\"");
        assert(sprint(mm("\"some string with \\\"apostrophs\\\" in it\"")) ===
            "\"some string with \\\"apostrophs\\\" in it\"");
        assert(sprint(mm("'a")) === "'a");
        assert(sprint(mm("'(a b c)")) === "'(a b c)");
        assert(sprint(mm("'(a b '(c 'd))")) === "'(a b '(c 'd))");
        assert(sprint(mm("'(a b '(c 'd ''d))")) === "'(a b '(c 'd ''d))");
    });

    test("numbers", function (mm) {
        assert(mm("(/ 3 5)"), 3 / 5);
    });

    test("and", function (mm) {
        assert(mm("(and)") === true);
        assert(mm("(and false)") === false);
        assert(mm("(and true)") === true);
        assert(mm("(and \"something\")") === "something");
        assert(mm("(and 1 2)") === 2);
        assert(mm("(and 1 2 3)") === 3);
        assert(mm("(and 1 2 false)") === false);
        assert(mm("(and 1 false 2)") === false);
        assert(mm("(and false 1 2)") === false);
    });

    test("or", function (mm) {
        assert(mm("(or)") === false);
        assert(mm("(or false)") === false);
        assert(mm("(or true)") === true);
        assert(mm("(or \"something\")") === "something");
        assert(mm("(or 1)") === 1);
        assert(mm("(or 1 2)") === 1);
        assert(mm("(or 1 2 3)") === 1);
        assert(mm("(or 1 false)") === 1);
        assert(mm("(or false 1)") === 1);
        assert(mm("(or false false 1)") === 1);
    });

    test("let", function (mm) {
        assert(mm("(let ((a 1) (b 2)) b a)") === 1);
        assert(mm("(let ((a 1) (b 2)) (+ a b))") === 3);
        fail("invalid arity", function () {
            mm("(let)");
        });
        fail("invalid syntax", function () {
            mm("(let 1 2 3)");
        });
        fail("not a pair", function () {
            mm("(let (1) 2 3)");
        });
    });

    test("try", function (mm) {
        assert(eq(mm("(try (lambda () 'ok) (lambda (error) 'notok))"),
            symbol("ok")));
        assert(eq(mm("(try (lambda () (error \"test\" \"error\" 0)) (lambda (error) 'notok))"),
            symbol("notok")));
    });
};

switch (getCommand()) {
case "test":
    runTests();
    break;
}
