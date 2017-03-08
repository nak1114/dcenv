package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kr/pretty"
	"github.com/urfave/cli"
)

var commandLocal = cli.Command{
	Name:      "local",
	Aliases:   []string{"l"},
	Usage:     "Deply image config data to current dir",
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
			Usage: "Remove this image config data from a current config file",
		},
		cli.BoolFlag{
			Name:  "check, c",
			Usage: "Check and rewite a current config file",
		},
		cli.StringSliceFlag{
			Name:  "environment, e",
			Usage: "Environment variable used in script",
		},
	},

	Action: local,
}

func local(c *cli.Context) {
	isForce := c.Bool("force")

	if isV {
		fmt.Println("dcenv local ", c.Args())
	}

	p, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return
	}

	fname := filepath.Join(p, ".dcenv_"+envShell)
	cfg := NewConfig(fname)
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
