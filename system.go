package main

import (
  "fmt"
  "os"

  "github.com/urfave/cli"
)

var Command_system = cli.Command{
  Name:      "system",
  Aliases:   []string{"s"},
  Usage:     "Execute command without shim dir",
  ArgsUsage: "<command> [args...]",
  //      Flags:   []cli.Flag{
  //      },
  Action: system,
}

func system(c *cli.Context) {

  if isV {
    fmt.Println("dcenv system ", c.Args())
  }
  if len(c.Args()) < 1 {
    fmt.Println("No command.")
    cli.ShowSubcommandHelp(c)
    return
  }

  MakeArgsFile(len(os.Args) - len(c.Args()) + 1)
  MakeExecFileSystem(c.Args()[0])

}
