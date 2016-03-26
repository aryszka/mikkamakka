// requires first:
// compound-types.h
// value.h

environment mkenvironment(environment parent);
void defvar(environment env, value sym, value val);
int hasvar(environment env, value sym);
value getvar(environment env, value sym);
void setvar(environment env, value sym, value val);
environment extenv(environment env, value syms, value vals);
registryflag getenvregistryflag(environment);
void setenvregistryflag(environment, registryflag);
valuenode allenvvalues(environment, valuenode);
envnode allenvenvs(environment, envnode);
void freenv(environment env);
