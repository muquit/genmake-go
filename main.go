package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

// Global variables
const (
	VERSION     = "1.0.1"
	AUTHOR      = "Muhammad A Muquit, muquit@muquit.com"
	URL         = "genmake-go https://github.com/muquit/genmake-go"
	LICENSE     = "License is MIT"
	GEN_APP     = 1
	GEN_LIB     = 2
	GEN_DLL     = 3
)

// OS specific compilation commands
var osCommands = map[string][]string{
	"SunOS": {"cc", "cc -KPIC", "cc", "ld -G"},
	"Linux": {"gcc", "gcc -fPIC", "gcc", "ld -shared"},
}

// Templates
const unixAppTemplate = `##----------------------------------------------------------------------------
# -=VERSION=-
# -=URL=-
# -=LICENSE=-
##----------------------------------------------------------------------------
rm=/bin/rm -f
CC= -=CC=-
DEFS=  
PROGNAME= -=PROG=-
INCLUDES=  -I.
LIBS=


DEFINES= $(INCLUDES) $(DEFS) -DSYS_UNIX=1
CFLAGS= -g $(DEFINES)

SRCS = -=SRCS=-

OBJS = -=OBJS=-

.c.o:
	$(rm) $@
	$(CC) $(CFLAGS) -c $*.c

all: $(PROGNAME)

$(PROGNAME) : $(OBJS)
	$(CC) $(CFLAGS) -o $(PROGNAME) $(OBJS) $(LIBS)

clean:
	$(rm) $(OBJS) $(PROGNAME) core *~`

const unixLibTemplate = `##----------------------------------------------------------------------------
# -=VERSION=-
# -=URL=-
# -=LICENSE=-
##----------------------------------------------------------------------------
rm=/bin/rm -f
CC= -=CC=-
LINK= -=LINK=-
DEFS=  
AR= ar cq
RANLIB= -=ranlib=-
LIBNAME= -=LIBN=-

INCLUDES=  -I. -I..

DEFINES= $(INCLUDES) $(DEFS) -DSYS_UNIX=1
CFLAGS= -g $(DEFINES)

SRCS = -=SRCS=-

OBJS = -=OBJS=-

.c.o:
	$(rm) -f $@
	$(CC) $(CFLAGS) -c $*.c

all: $(LIBNAME)

$(LIBNAME) : $(OBJS)
	$(rm) $@
	$(AR) $@ $(OBJS)
	$(RANLIB) $@

clean:
	$(rm) $(OBJS) $(LIBNAME) core *~`

const winAppTemplate = `##----------------------------------------------------------------------------
# -=VERSION=-
# -=URL=-
# -=LICENSE=-
##----------------------------------------------------------------------------

CC= cl
#DEFS=  -nologo -DSTRICT -G3 -Ow -W3 -Zp -Tp
DEFS=  -nologo -G3
PROGNAME= -=PROGN=-
LINKER=link -nologo

INCLUDES=  -I. 

# don't define -DSYS_WIN32.. win2k complains
DEFINES= $(INCLUDES) $(DEFS) -DWINNT=1 

CFLAGS= $(DEFINES)
GUIFLAGS=user32.lib gdi32.lib winmm.lib comdlg32.lib comctl32.lib
WINSOCK_LIB=wsock32.lib
LIBS=$(WINSOCK_LIB) $(GUIFLAGS)
RC=rc
RCVARS=-r -DWIN32

SRCS = -=SRCS=-

OBJS = -=OBJS=-

.c.obj:
	$(CC) $(CFLAGS) -c $< -Fo$@

all: $(PROGNAME)

$(PROGNAME) : $(OBJS)
	$(LINKER) $(OBJS) /OUT:$(PROGNAME) $(LIBS)

clean:
	del $(OBJS) $(PROGNAME) core`

const winLibTemplate = `##----------------------------------------------------------------------------
# -=VERSION=-
# -=URL=-
# -=LICENSE=-
##----------------------------------------------------------------------------

CC= cl
DEFS=  -DWINNT=1

INCLUDES=  -I. -I..
LIBRARY= -=LIBN=-

# replace -O with -g in order to debug

DEFINES= $(INCLUDES) $(DEFS) 
#CFLAGS=  $(cvars) $(cdebug) -nologo -G4 $(DEFINES)


SRCS = -=SRCS=-
OBJS = -=OBJS=-

.c.obj:
	$(CC) $(CFLAGS) $(DEFS) -c $< -Fo$@

all: $(LIBRARY)

$(LIBRARY): $(OBJS)
	link -=WHAT=- /OUT:$(LIBRARY) $(OBJS) 

clean:
	del $(OBJS) $(LIBRARY) *.bak`

// Debug flag
var debug bool

// Command line options
var (
	unix    bool
	windows bool
	app     string
	lib     string
	dll     string
	version bool
	help    bool
)

func main() {
	// Parse command line flags
	flag.BoolVar(&unix, "unix", false, "generate Makefile for Unix")
	flag.BoolVar(&windows, "win", false, "generate Makefile for MS Windows")
	flag.StringVar(&app, "app", "", "generate Makefile for an application")
	flag.StringVar(&lib, "lib", "", "generate Makefile for a static library")
	flag.StringVar(&dll, "dll", "", "generate Makefile for a shared library in Unix and DLL in Windows")
	flag.BoolVar(&version, "version", false, "show version info")
	flag.BoolVar(&debug, "debug", false, "enable debug output")
	flag.BoolVar(&help, "help", false, "show help")

	flag.Parse()

	if version {
		fmt.Printf("genmake-go v%s\n", VERSION)
		os.Exit(0)
	}

	if help {
		showUsage()
	}

	if unix && windows {
		printError("--unix and --win are mutually exclusive\n")
		os.Exit(1)
	}

	printDebug(fmt.Sprintf("unix: %v, win: %v\n", unix, windows))

	// Auto-detect platform if not specified
	if !unix && !windows {
		os := detectPlatform()
		printDebug(fmt.Sprintf("os: %s\n", os))
		if os != "win" {
			unix = true
		} else {
			windows = true
		}
	}

	if app == "" && lib == "" && dll == "" {
		showUsage()
	}

	if lib != "" && dll != "" {
		printError("--lib and --dll are mutually exclusive\n")
		os.Exit(1)
	}

	if app != "" && lib != "" {
		printError("--app and --lib are mutually exclusive\n")
		os.Exit(1)
	}

	// Get source files from remaining arguments
	args := flag.Args()
	if len(args) == 0 {
		showUsage()
	}

	var doWhat int
	if app != "" {
		doWhat = GEN_APP
	} else if lib != "" {
		doWhat = GEN_LIB
	} else if dll != "" {
		doWhat = GEN_DLL
	}

	if unix {
		genUnixMakefile(doWhat, args)
	} else {
		genWindowsMakefile(doWhat, args)
	}
}

func showUsage() {
	fmt.Printf(`genmake-go v%s %s
A program to generate nice/simple Makefiles for Linux/Unix and MS Windows

Usage: genmake-go [options]
Where the options include:
    --unix      generate Makefile for Unix 
    --win       generate Makefile for MS Windows
    --app=name  generate Makefile for an application
    --lib=name  generate Makefile for a static library
    --dll=name  generate Makefile for a shared library in Unix and DLL in Windows
    --version   show version info

Example:
    genmake-go --unix --app=myapp *.c > Makefile
    genmake-go --win --app=myapp.exe main.c bar.c > Makefile.win
    genmake-go --unix --lib=libmyapp.a main.c bar.c > Makefile
    genmake-go --win --lib=myapp.lib main.c bar.c > Makefile.win
    genmake-go --unix --dll=libmyapp.so main.c bar.c > Makefile
    genmake-go --win --dll=myapp.dll main.c bar.c > Makefile.win

If no --unix or --win flag is specified, OS type will be guessed

Edit the generated Makefile if needed.
`, VERSION, URL)
	os.Exit(1)
}

func printError(msg string) {
	fmt.Fprintf(os.Stderr, "%s", msg)
}

func printDebug(msg string) {
	if debug {
		fmt.Fprintf(os.Stderr, "(debug) %s", msg)
	}
}

func detectPlatform() string {
	output, err := exec.Command("uname").Output()
	if err != nil {
		return "win"
	}

	os := strings.TrimSpace(string(output))
	printDebug(fmt.Sprintf("os: %s\n", os))

	if strings.HasPrefix(strings.ToLower(os), "cygwin") {
		return "win"
	}

	return os
}

func runCommand(cmd string) string {
	parts := strings.Split(cmd, " ")
	output, err := exec.Command(parts[0], parts[1:]...).Output()
	if err != nil {
		printError(fmt.Sprintf("Could not run: %s\n", cmd))
		return ""
	}
	return strings.TrimSpace(string(output))
}

func hasRanlib() string {
	return runCommand("which ranlib")
}

func genUnixMakefile(doWhat int, args []string) {
	var template string
	srcs := strings.Join(args, " ")
	
	// Convert source files to object files
	re := regexp.MustCompile(`\.[a-zA-Z]+`)
	objs := re.ReplaceAllString(srcs, ".o")
	
	date := time.Now().Format(time.RFC1123)
	created := fmt.Sprintf("Created with genmake-go v%s on %s", VERSION, date)
	
	switch doWhat {
	case GEN_APP:
		template = unixAppTemplate
		template = strings.ReplaceAll(template, "-=PROG=-", app)
		template = strings.ReplaceAll(template, "-=VERSION=-", created)
		template = strings.ReplaceAll(template, "-=URL=-", URL)
		template = strings.ReplaceAll(template, "-=LICENSE=-", LICENSE)
		template = strings.ReplaceAll(template, "-=CC=-", "cc")
		template = strings.ReplaceAll(template, "-=SRCS=-", srcs)
		template = strings.ReplaceAll(template, "-=OBJS=-", objs)
		
	case GEN_LIB:
		ranlib := hasRanlib()
		template = unixLibTemplate
		template = strings.ReplaceAll(template, "-=LIBN=-", lib)
		template = strings.ReplaceAll(template, "-=VERSION=-", created)
		template = strings.ReplaceAll(template, "-=URL=-", URL)
		template = strings.ReplaceAll(template, "-=LICENSE=-", LICENSE)
		template = strings.ReplaceAll(template, "-=CC=-", "cc")
		template = strings.ReplaceAll(template, "-=SRCS=-", srcs)
		template = strings.ReplaceAll(template, "-=OBJS=-", objs)
		template = strings.ReplaceAll(template, "-=ranlib=-", ranlib)
		
	case GEN_DLL:
		ranlib := hasRanlib()
		os := detectPlatform()
		printDebug(fmt.Sprintf("OS: %s\n", os))
		
		cc := "cc"
		linker := "cc"
		
		if cmds, ok := osCommands[os]; ok {
			cc = cmds[1]
			linker = cmds[3]
		}
		
		template = unixLibTemplate
		template = strings.ReplaceAll(template, "-=LIBN=-", dll)
		template = strings.ReplaceAll(template, "-=VERSION=-", created)
		template = strings.ReplaceAll(template, "-=URL=-", URL)
		template = strings.ReplaceAll(template, "-=LICENSE=-", LICENSE)
		template = strings.ReplaceAll(template, "-=CC=-", cc)
		template = strings.ReplaceAll(template, "-=LINK=-", linker)
		template = strings.ReplaceAll(template, "-=SRCS=-", srcs)
		template = strings.ReplaceAll(template, "-=OBJS=-", objs)
		template = strings.ReplaceAll(template, "-=ranlib=-", ranlib)
		template = strings.ReplaceAll(template, "	$(AR) $@ $(OBJS)", "	$(LINK) -o $(LIBNAME) $(LIBS)")
		template = strings.ReplaceAll(template, "	$(RANLIB) $@", "")
	}
	
	fmt.Println(template)
}

func genWindowsMakefile(doWhat int, args []string) {
	var template string
	srcs := strings.Join(args, " ")
	
	// Convert source files to object files for Windows
	re := regexp.MustCompile(`\.[a-zA-Z]+`)
	objs := re.ReplaceAllString(srcs, ".obj")
	
	date := time.Now().Format(time.RFC1123)
	created := fmt.Sprintf("Created with genmake-go v%s on %s", VERSION, date)
	
	switch doWhat {
	case GEN_APP:
		template = winAppTemplate
		template = strings.ReplaceAll(template, "-=PROGN=-", app)
		template = strings.ReplaceAll(template, "-=VERSION=-", created)
		template = strings.ReplaceAll(template, "-=URL=-", URL)
		template = strings.ReplaceAll(template, "-=LICENSE=-", LICENSE)
		template = strings.ReplaceAll(template, "-=CC=-", "cc")
		template = strings.ReplaceAll(template, "-=SRCS=-", srcs)
		template = strings.ReplaceAll(template, "-=OBJS=-", objs)
		
	case GEN_LIB:
		template = winLibTemplate
		template = strings.ReplaceAll(template, "-=LIBN=-", lib)
		template = strings.ReplaceAll(template, "-=VERSION=-", created)
		template = strings.ReplaceAll(template, "-=URL=-", URL)
		template = strings.ReplaceAll(template, "-=LICENSE=-", LICENSE)
		template = strings.ReplaceAll(template, "-=CC=-", "cc")
		template = strings.ReplaceAll(template, "-=SRCS=-", srcs)
		template = strings.ReplaceAll(template, "-=OBJS=-", objs)
		template = strings.ReplaceAll(template, "-=WHAT=-", "/lib")
		
	case GEN_DLL:
		template = winLibTemplate
		template = strings.ReplaceAll(template, "-=LIBN=-", dll)
		template = strings.ReplaceAll(template, "-=VERSION=-", created)
		template = strings.ReplaceAll(template, "-=URL=-", URL)
		template = strings.ReplaceAll(template, "-=LICENSE=-", LICENSE)
		template = strings.ReplaceAll(template, "-=CC=-", "cc")
		template = strings.ReplaceAll(template, "-=SRCS=-", srcs)
		template = strings.ReplaceAll(template, "-=OBJS=-", objs)
		template = strings.ReplaceAll(template, "-=WHAT=-", "/dll")
	}
	
	fmt.Println(template)
}
