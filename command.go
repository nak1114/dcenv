package main

import (
  "fmt"
  "sort"

  "github.com/urfave/cli"
)

var Command_command = cli.Command{
  Name:      "command",
  Aliases:   []string{"c"},
  Usage:     "Display command detail",
  ArgsUsage: "[options...] command...",
  Flags: []cli.Flag{
    cli.BoolFlag{
      Name:  "list, l",
      Usage: "show a list of all command",
    },
    cli.BoolFlag{
      Name:  "display, d",
      Usage: "Show a detail fo the command",
    },
    // cli.BoolFlag{
    //   Name: "json, j",
    //   Usage: "display by json style.",
    // },
  },

  Action: command,
}

type Strp [][3]string

func (p Strp) Len() int {
  return len(p)
}
func (p Strp) Swap(i, j int) {
  p[i], p[j] = p[j], p[i]
}
func (p Strp) Less(i, j int) bool {
  return p[i][2] < p[j][2]
}

func command(c *cli.Context) {
  if isV {
    fmt.Println("dcenv command ", c.Args())
  }

  if c.Bool("list") {
    listCommand()
    return
  }

  if len(c.Args()) < 1 {
    fmt.Println("No command.")
    cli.ShowSubcommandHelp(c)
    return
  }

  //isJson := c.Bool("json")
  for _, cmd := range c.Args() {
    cfg, fname := SearchConfig(cmd, CheckCommand)
    fmt.Printf("---[ %s ]------------------------------\n", cmd)
    if cfg == nil {
      fmt.Println("    Command not found.")
      continue
    }
    cval := (*cfg).Commands[cmd]
    cnt := (*cfg).Images[cval]

    fmt.Println("Filename:", fname)
    fmt.Println("Image   :", cval)
    fmt.Println("Tag     :", cnt.Tag)
    fmt.Println("Fake    :", cnt.Fake)
    fmt.Printf("Envs    :%v\n", cnt.Commands[cmd])
    fmt.Printf("Script  :\n%s\n", cnt.Script)
    fmt.Println("Execute :")
    ShowExecFile(cfg, cmd, fname)
  }

}

func listCommand() {
  InitCheckCmds()
  SearchConfig("", CheckCmds)
  lk := 7
  lv := 5
  arys := Strp{}
  for key, val := range CheckCmdsMap {
    if len(key) > lk {
      lk = len(key)
    }
    if len(val[0]) > lv {
      lv = len(val[0])
    }
    arys = append(arys, [3]string{val[0], val[1], key})
  }
  sort.Sort(arys)
  fmt.Printf("%-*s | %-*s (%s)\n", lk, "COMMAND", lv, "IMAGE", "FILE")
  for _, val := range arys {
    fmt.Printf("%-*s | %-*s (%s)\n", lk, val[2], lv, val[0], val[1])
  }
}
