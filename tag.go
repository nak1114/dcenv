package main

import (
  "fmt"

  "github.com/urfave/cli"
)

var Command_tag = cli.Command{
  Name:      "tag",
  Aliases:   []string{"t"},
  Usage:     "Change image tag",
  ArgsUsage: "[options...] image:tag",
  //      Flags:   []cli.Flag{
  //      },

  Action: tag,
}

func tag(c *cli.Context) {
  if isV {
    fmt.Println("dcenv install ", c.Args())
  }
  if len(c.Args()) < 1 {
    fmt.Println("No image:tag.")
    cli.ShowSubcommandHelp(c)
    return
  }
  tName, _, tTag := ParseImageTag(c.Args()[0])
  cfg, fname := SearchConfig(tName, CheckImage)
  if cfg == nil {
    return
  }
  //(*cfg).Images[tName].Tag=tTag
  //Work around
  tc := (*cfg).Images[tName]
  tc.Tag = tTag
  (*cfg).Images[tName] = tc
  (*cfg).WriteToFile(fname)

}
