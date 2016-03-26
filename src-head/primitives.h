// requires first:
// value.h

value eofval;
value stdinval;
value stdoutval;
value stderrval;

value errorval(value args);
value iseqval(value args);
value sumval(value args);
value diffval(value args);
value bitorval(value args);
value openfileval(value args);
value seekfileval(value args);
value readfileval(value args);
value writefileval(value args);
value iseofval(value args);
value closefileval(value args);
value mkregexval(value args);
value regexmatch(value args);
value byteslenval(value args);
value stringlenval(value args);
value isutf8val(value args);
value copystrval(value args);
value stringappendval(value args);
value islessval(value args);
value islessval(value args);
value isgreaterval(value args);
value islessoreqval(value args);
value isgreateroreqval(value args);
value notval(value args);
value isnullvalp(value args);
value consval(value args);
value carvalp(value args);
value cdrvalp(value args);
value stringtonumsafe(value args);
value isnumvalp(value args);
value isstringvalp(value args);
value ispairvalp(value args);
value issymbolvalp(value args);
value isintvalp(value args);
value numbertostring(value args);
value stringtosymbol(value args);
value symboltostring(value args);
void initmodule_primitives();
void freemodule_primitives();
