package main

import (
	"fmt"
	"github.com/urfave/cli"
	"os"
)

var commandExec = cli.Command{
	Name:      "exec",
	Aliases:   []string{"e"},
	Usage:     "Execute command",
	ArgsUsage: "[options...] <command> [args...]",
	//  Flags:   []cli.Flag{
	//  },
	SkipFlagParsing: true,
	Action:          execute,
}

func execute(c *cli.Context) {
	if isV {
		fmt.Println("dcenv exec ", c.Args())
	}
	if len(c.Args()) < 1 {
		fmt.Println("No command.")
		cli.ShowSubcommandHelp(c)
		return
	}

	cmd := (c.Args()[0])
	MakeArgsFile(len(os.Args) - len(c.Args()) + 1)
	cfg, fname := SearchConfig(cmd, CheckCommand)
	if cfg == nil {
		MakeExecFileSystem(cmd)
	} else {
		MakeExecFile(cfg, cmd, fname)
	}
}
