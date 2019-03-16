#include <rpc/rpc.h>

#include <stdio.h>

 

#include "trans.h"/* Client-side stub interface routines written by programmer */

extern CLIENT * handle;

static char **ret;

 

char * readfile(char * name)

{

   char ** arg;

   arg = &name;

   ret = readfile_1(arg, handle);

 

   return ret==NULL ? NULL : *ret;

}