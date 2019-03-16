#include <rpc/rpc.h>

#include <stdio.h>

#include "trans.h"

 

char * readfile(char *);

static char * retcode;

 

char ** readfile_1(char ** w, CLIENT *clnt)

{

   retcode = readfile(*(char**)w);

   return &retcode;

}

