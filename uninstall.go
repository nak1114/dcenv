package main

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/urfave/cli"
)

var commandUninstall = cli.Command{
	Name:      "uninstall",
	Aliases:   []string{"u"},
	Usage:     "Remove config file from your yard",
	ArgsUsage: "[options...] image_name",
	//      Flags:   []cli.Flag{
	//      },

	Action: uninstall,
}

func uninstall(c *cli.Context) {
	if isV {
		fmt.Println("dcenv uninstall ", c.Args())
	}
	if len(c.Args()) < 1 {
		fmt.Println("No filename.")
		cli.ShowSubcommandHelp(c)
		return
	}
	dname, _, _ := ParseImageTag(c.Args()[0])
	fname := filepath.Join(envHome, "image_yard", url.QueryEscape(dname)+".yml")
	if _, err := os.Stat(fname); err != nil {
		fmt.Println(err)
		return
	}
	os.Remove(fname)

}
