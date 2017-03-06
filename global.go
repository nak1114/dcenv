package main

import (
	"fmt"
	"path/filepath"

	"github.com/kr/pretty"
	"github.com/urfave/cli"
)

var commandGlobal = cli.Command{
	Name:      "global",
	Aliases:   []string{"g"},
	Usage:     "Deply image config data to DCENV dir",
	ArgsUsage: "[options...] image[:tag] [commands...]",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "force, f",
			Usage: "Ignore confirm",
		},
		cli.IntFlag{
			Name:  "num, n",
			Value: 0,
			Usage: "Specify the number of the image config data",
		},
		cli.BoolFlag{
			Name:  "remove, r",
			Usage: "Remove this image config data from the global file",
		},
		cli.BoolFlag{
			Name:  "check, c",
			Usage: "Check and rewite the global file",
		},
		cli.StringSliceFlag{
			Name:  "environment, e",
			Usage: "Environment variable used in script",
		},
	},

	Action: global,
}

func global(c *cli.Context) {
	isForce := c.Bool("force")

	if isV {
		fmt.Println("dcenv local ", c.Args())
	}

	p := filepath.Join(envHome, "files")

	fname := filepath.Join(p, ".dcenv_"+envShell)
	cfg := GetConfig(fname)
	if !c.Bool("check") {
		if len(c.Args()) < 1 {
			fmt.Println("No image name.")
			cli.ShowSubcommandHelp(c)
			return
		}
		tName, tCommand, tTag := ParseImageTag(c.Args()[0])
		if !c.Bool("remove") {
			tc := SearchImageFromYard(tName, tCommand, tTag, c.Int("num"))
			if len(c.Args()) > 2 {
				tc.MakeCommands(c.Args()[1:], c.StringSlice("environment"))
			}
			cfg.AddImage(tc, tName, isForce)
		} else {
			cfg.DelImage(tName, isForce)
		}
	}
	cfg.WriteToFile(fname)
	fmt.Println("Complete!.")
	if isV {
		pretty.Printf("--- cofig %s:\n%# v\n\n", fname, cfg)
	}
}
