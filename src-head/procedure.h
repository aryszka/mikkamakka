// requires first:
// compound-types.h
// value.h

procedure mkprimitiveproc(primitive p);
procedure mkcompiledproc(value label, environment env);
int isprimitive(procedure p);
value applyprimitive(procedure p, value args);
value proclabel(procedure p);
environment procenv(procedure p);
void freeproc(procedure p);
