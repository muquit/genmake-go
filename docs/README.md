## Introduction

`genmake-go` is a simple cross platform tool to generate _simple_ Makefiles for C programs. It is a golang port of my
perl script [genmake](https://muquit.com/muquit/software/genmake/genmake.html).

# Synopsis

```
genmake-go v1.0.1 genmake-go https://github.com/muquit/genmake-go
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

Usage of ./genmake-go:
  -app string
    	generate Makefile for an application
  -debug
    	enable debug output
  -dll string
    	generate Makefile for a shared library in Unix and DLL in Windows
  -help
    	show help
  -lib string
    	generate Makefile for a static library
  -unix
    	generate Makefile for Unix
  -version
    	show version info
  -win
    	generate Makefile for MS Windows
```

# Version
The current version is 1.0.1

Please look at [ChangeLog](ChangeLog.md) for what has changed in the current version.

# Quick Start

Install [go](https://go.dev/) first

```bash
go install github.com/muquit/genmake-go@latest
genmake-go -version
```

# Download

Download pre-compiled binaries from
[Releases](https://github.com/muquit/genmake-go/releases) page

# Building from source

Install [go](https://go.dev/) first

```bash
git clone https://github.com/muquit/genmake-go
cd genmake-go
go build .
```

Look at Makefile for more info

# License

MIT License - See LICENSE.txt file for details.

# Authors

Developed with Claude AI 3.7 Sonnet, working under my guidance and instructions - translated from
the perl script [genmake](https://muquit.com/muquit/software/genmake/genmake.html) to go
