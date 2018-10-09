# gosubcommand

[![Build Status](https://travis-ci.org/oirik/gosubcommand.svg?branch=master)](https://travis-ci.org/oirik/gosubcommand)
[![GoDoc](https://godoc.org/github.com/oirik/gosubcommand?status.svg)](https://godoc.org/github.com/oirik/gosubcommand)
[![apache license](https://img.shields.io/badge/license-Apache-blue.svg)](LICENSE)

gosubcommand is a Go package that enable to develop command-line application with `subcommand` like `go` more easily.

This is inspired by [google/subcommand](https://github.com/google/subcommands), and normally you should use google's one. 

This could be suitable if you were looking for much simpler one.

## Usage

Call `Register` and `Execute` like this:

```go
func main() {
	gosubcommand.Register("init", &initCommand{})
	gosubcommand.Register("status", &statusCommand{})
	os.Exit(int(gosubcommand.Execute()))
}
```

Each subcommand should be implemented like this: 

```go
type initCommand struct {
    flagX bool
}

func (init *initCommand) Summary() string {
    return "Init something"
}

func (init *initCommand) SetFlag(fs *flag.FlagSet) {
    fs.BoolVar(&init.flagX, "x", false, "x information")
}

func (init *initCommand) Execute(fs *flag.FlagSet) gosubcommand.ExitCode {
    // Do something
    return gosubcommand.ExitCodeSuccess
}
```

And your application will be implemented with subcommand like this:

```sh
$ yourapp help
yourapp's Summary.

Usage:

  yourapp <command> [arguments]

The commands are:

  init      Init something
  help      Show help information
  version   Show version information

Use "yourapp help <command>" for more information about a command.
```
