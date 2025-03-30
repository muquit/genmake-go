##----------------------------------------------------------------------------
# Created with genmake-go v1.0.1 on Sun, 30 Mar 2025 18:01:58 EDT
# genmake-go https://github.com/muquit/genmake-go
# License is MIT
##----------------------------------------------------------------------------
rm=/bin/rm -f
CC= cc
DEFS=  
PROGNAME= hello
INCLUDES=  -I.
LIBS=


DEFINES= $(INCLUDES) $(DEFS) -DSYS_UNIX=1
CFLAGS= -g $(DEFINES)

SRCS = main.c

OBJS = main.o

.c.o:
	$(rm) $@
	$(CC) $(CFLAGS) -c $*.c

all: $(PROGNAME)

$(PROGNAME) : $(OBJS)
	$(CC) $(CFLAGS) -o $(PROGNAME) $(OBJS) $(LIBS)

clean:
	$(rm) $(OBJS) $(PROGNAME) core *~
