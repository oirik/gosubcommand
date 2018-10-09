package gosubcommand

import (
	"flag"
	"fmt"
	"os"
	"path"
	"sort"
	"text/tabwriter"

	"github.com/pkg/errors"
)

// AppName is used to output root command name in help messages.
var AppName = path.Base(os.Args[0])

// Version is used to show version text in version command.
var Version = ""

// Summary is used to output root command summary in help messages.
var Summary = AppName

// ExitCode represents processes' return code which would be passed to os.Exit() method.
type ExitCode int

const (
	// ExitCodeSuccess represents the process successed.
	ExitCodeSuccess ExitCode = iota
	// ExitCodeError represents the process failed.
	ExitCodeError
	// ExitCodeUsageError represents wrong command args.
	ExitCodeUsageError
)

// Command represents subcommand which is implemented fro each subcommand.
type Command interface {
	Summary() string
	SetFlag(*flag.FlagSet)
	Execute(*flag.FlagSet) ExitCode
}

// Register regists the subcommand to gosubcommand.
func Register(name string, command Command) {
	commands[name] = command
}

// RegisterFunc is easy way of Register func.
func RegisterFunc(name string, summaryHandler func() string, setFlagHandler func(*flag.FlagSet), executeHandler func(*flag.FlagSet) ExitCode) {
	Register(name, &commandImp{summary: summaryHandler, setFlag: setFlagHandler, execute: executeHandler})
}

// Execute should be called afer all necessary Registers are called.
func Execute() ExitCode {

	Register("help", &helpCommand{})

	if Version != "" {
		Register("version", &versionCommand{})
	}

	flag.CommandLine.Usage = func() { printUsage() }
	flag.Parse()

	if flag.NArg() == 0 {
		printUsage()
		return ExitCodeUsageError
	}

	name, command, err := getCommand(flag.Arg(0))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return ExitCodeUsageError
	}
	fs := flag.NewFlagSet(name, flag.ExitOnError)
	fs.Usage = func() {
		printCommandUsage(name, fs)
	}
	command.SetFlag(fs)
	fs.Parse(flag.Args()[1:])
	return command.Execute(fs)
}

var commands = map[string]Command{}

func getCommand(commandName string) (string, Command, error) {
	for name, command := range commands {
		if commandName == name {
			return name, command, nil
		}
	}
	return "", nil, errors.Errorf("%s %s: unknown command\nRun '%s help' for usage", AppName, flag.Arg(0), AppName)
}

func printUsage() {
	w := flag.CommandLine.Output()
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)

	fmt.Fprintln(w, Summary)
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Usage:")
	fmt.Fprintln(w)
	fmt.Fprintf(tw, "\t%s <command> [arguments]\n", AppName)
	tw.Flush()
	fmt.Fprintln(w)
	fmt.Fprintln(w, "The commands are:")
	fmt.Fprintln(w)
	sortedKeys := []string{}
	for key := range commands {
		sortedKeys = append(sortedKeys, key)
	}
	sort.Strings(sortedKeys)
	for _, name := range sortedKeys {
		fmt.Fprintf(tw, "\t%s\t%s\n", name, commands[name].Summary())
	}
	tw.Flush()
	fmt.Fprintln(w)
	fmt.Fprintf(w, "Use \"%s help <command>\" for more information about a command.\n", AppName)
	fmt.Fprintln(w)
}

func printCommandUsage(name string, fs *flag.FlagSet) {
	w := fs.Output()
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)

	fmt.Fprintln(w)
	fmt.Fprintln(w, "Usage:")
	fmt.Fprintln(w)
	fmt.Fprintf(tw, "\t%s %s [arguments]\n", AppName, name)
	tw.Flush()
	fmt.Fprintln(w)
	fmt.Fprintln(w, "The flags are:")
	fmt.Fprintln(w)
	fs.PrintDefaults()
	fmt.Fprintln(w)
}

type commandImp struct {
	summary func() string
	setFlag func(*flag.FlagSet)
	execute func(*flag.FlagSet) ExitCode
}

func (imp *commandImp) Summary() string                   { return imp.summary() }
func (imp *commandImp) SetFlag(fs *flag.FlagSet)          { imp.setFlag(fs) }
func (imp *commandImp) Execute(fs *flag.FlagSet) ExitCode { return imp.execute(fs) }

type helpCommand struct {
}

func (help *helpCommand) Summary() string          { return "Show help information" }
func (help *helpCommand) SetFlag(fs *flag.FlagSet) {}
func (help *helpCommand) Execute(fs *flag.FlagSet) ExitCode {
	if fs.NArg() == 0 {
		printUsage()
		return ExitCodeSuccess
	}

	name, command, err := getCommand(fs.Arg(0))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return ExitCodeUsageError
	}
	tmpfs := flag.NewFlagSet(name, flag.ExitOnError)
	tmpfs.Usage = func() {
		printCommandUsage(name, tmpfs)
	}
	command.SetFlag(tmpfs)
	tmpfs.Usage()
	return ExitCodeSuccess
}

type versionCommand struct {
}

func (version *versionCommand) Summary() string          { return "Show version information" }
func (version *versionCommand) SetFlag(fs *flag.FlagSet) {}
func (version *versionCommand) Execute(fs *flag.FlagSet) ExitCode {
	fmt.Fprintln(fs.Output(), Version)
	return ExitCodeSuccess
}
