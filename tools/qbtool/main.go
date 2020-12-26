package main

import (
	"flag"
	"fmt"
	"git.moumentei.com/plat_go/golib/tools/qbtool/cmd/base"
	"git.moumentei.com/plat_go/golib/tools/qbtool/cmd/db"
	"git.moumentei.com/plat_go/golib/tools/qbtool/cmd/help"
	"git.moumentei.com/plat_go/golib/tools/qbtool/cmd/swag"
	"log"
	"os"
	"strings"
)

const version = "0.1.1"

func init() {
	base.Cmd.Commands = []*base.Command{
		help.Cmd,
		swag.Cmd,
		db.Cmd,
	}
	base.Usage = mainUsage
}
func mainUsage() {
	help.PrintUsage(os.Stderr, base.Cmd)
	os.Exit(2)
}

func main() {
	flag.Usage = base.Usage
	flag.Parse()
	log.SetFlags(0)
	args := flag.Args()
	if len(args) < 1 {
		base.Usage()
		return
	}

	base.CmdName = args[0] // for error message
	// Show help documentation.
	if args[0] == "help" {
		help.Help(os.Stdout, args[1:])
		return
	}
BigCmdLoop:
	for bigCmd := base.Cmd; ; {
		for _, cmd := range bigCmd.Commands {
			if cmd.Name() != args[0] {
				continue
			}
			if len(cmd.Commands) > 0 {
				bigCmd = cmd
				args = args[1:]
				if len(args) == 0 {
					help.PrintUsage(os.Stderr, bigCmd)
					base.SetExitStatus(2)
					base.Exit()
				}
				if args[0] == "help" {
					// Accept 'go mod help' and 'go mod help foo' for 'go help mod' and 'go help mod foo'.
					help.Help(os.Stdout, append(strings.Split(base.CmdName, " "), args[1:]...))
					return
				}
				base.CmdName += " " + args[0]
				continue BigCmdLoop
			}
			if !cmd.Runnable() {
				continue
			}
			cmd.Flag.Usage = func() { cmd.Usage() }
			if cmd.CustomFlags {
				args = args[0:]
			} else {
				cmd.Flag.Parse(args[1:])
				args = cmd.Flag.Args()
			}
			cmd.Run(cmd, args)
			base.Exit()
			return
		}
		helpArg := ""
		if i := strings.LastIndex(base.CmdName, " "); i >= 0 {
			helpArg = " " + base.CmdName[:i]
		}
		fmt.Fprintf(os.Stderr, "%s %s: unknown command\nRun '%s help%s' for usage.\n", base.ToolName, base.ToolName, base.CmdName, helpArg)
		base.SetExitStatus(2)
		base.Exit()
	}
}
