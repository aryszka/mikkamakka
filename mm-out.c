
        #include <stdlib.h>
        #include <stdio.h>
        #include "../src-head/sys.h"
        #include "../src-head/error.h"
        #include "../src-head/sysio.h"
        #include "../src-head/number.h"
        #include "../src-head/string.h"
        #include "../src-head/compound-types.h"
        #include "../src-head/stack.h"
        #include "../src-head/io.h"
        #include "../src-head/value.h"
        #include "../src-head/register-machine.h"
        #include "../src-head/primitives.h"

        int main() {
            initsys();
            initmodule_sysio();
            initmodule_io();
            initmodule_number();
            initmodule_value();
            initmodule_primitives();

            regmachine rm = mkregmachine();
            for (;;) {
                long labelval = valrawint(rm->label);
                switch (labelval) {
                case 0:
                mkcompiledprocreg(rm, (void *)&rm->proc, 2);
gotolabel(rm, 1); break;
case 2:
initprocenv(rm, null);
getenvvar(rm, (void *)&rm->proc, "write-file");
getenvvar(rm, (void *)&rm->val, "stdout");
initargs(rm);
addarg(rm);
initreg((void *)&rm->val, stringval("Hello, world!"));
addarg(rm);
if (branchproc(rm, 5)) { break; }
case 4:
takeproclabel(rm);
gotoreg(rm, rm->val); break;
case 5:
applyprimitivereg(rm, (void *)&rm->val);
gotoreg(rm, rm->cont); break;
case 3:
case 1:
initargs(rm);
if (branchproc(rm, 8)) { break; }
case 7:
takeproclabel(rm);
gotoreg(rm, rm->val); break;
case 8:
applyprimitivereg(rm, (void *)&rm->val);
gotoreg(rm, rm->cont); break;
case 6:

                default:
                    // printf("%s\n", sprintraw(rm->val));
                    return 0;
                }
            }

            freeregmachine(rm);

            freemodule_primitives();
            freemodule_value();
            freemodule_number();
            freemodule_io();
            freemodule_sysio();

            return 0;
        }
        