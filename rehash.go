package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/urfave/cli"
)

var commandRehash = cli.Command{
	Name:      "rehash",
	Aliases:   []string{"r"},
	Usage:     "Rehash shim dir",
	ArgsUsage: "[options...]",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "clean, c",
			Usage: "Only remove all comannd in shims",
		},
		cli.BoolFlag{
			Name:  "append, a",
			Usage: "Append a comannd to current shims",
		},
	},

	Action: rehash,
}

func rehash(c *cli.Context) {
	sdir := filepath.Join(envHome, "shims")
	tdir := filepath.Join(envHome, "tmp")
	if !c.Bool("append") {
		for i := 0; i < 100; i++ {
			_, e := os.Stat(sdir)
			if e != nil {
				for j := 0; j < 100; j++ {
					_, e := os.Stat(sdir)
					if e == nil {
						break
					}
					os.Mkdir(sdir, 0777)
				}
				break
			}
			os.RemoveAll(sdir)
		}
	}
	for i := 0; i < 100; i++ {
		_, e := os.Stat(tdir)
		if e != nil {
			for j := 0; j < 100; j++ {
				_, e := os.Stat(tdir)
				if e == nil {
					break
				}
				os.Mkdir(tdir, 0777)
			}
			break
		}
		os.RemoveAll(tdir)
	}
	if c.Bool("clean") {
		return
	}
	SearchConfig("", makeShimCommands)
}

func makeShimCommands(cmd string, fname string) *Config {
	if isV {
		fmt.Println("search:", fname)
	}
	if _, err := os.Stat(fname); err != nil {
		return nil
	}
	if isV {
		fmt.Println("File found.:", fname)
	}
	m := LoadConfig(fname)
	if m == nil {
		return nil
	}

	if m.Commands == nil {
		if isV {
			fmt.Println("  Illigal file format.:", fname)
		}
		return nil
	}
	for key := range m.Commands {
		MakeShimsFile(key)
		if isV {
			fmt.Println("Write command in shims.:", key)
		}
	}
	return nil
}

/*
func makeShimCommand(fname string) bool {
	if isV {
		fmt.Println("search:", fname)
	}
	if _, err := os.Stat(fname); err != nil {
		return false
	}
	if isV {
		fmt.Println("File found.:", fname)
	}
	m := Config{}
	if err := LoadYaml(&m, fname); err != nil {
		if isV {
			fmt.Println("  File can not load.:", fname)
		}
		return false
	}
	if m.Commands == nil {
		if isV {
			fmt.Println("  Illigal file format.:", fname)
		}
		return false
	}
	shimDir := filepath.Join(envHome, "shims")
	for key := range m.Commands {
		if runtime.GOOS == `windows` {
			ioutil.WriteFile(filepath.Join(shimDir, key+`.bat`),
				[]byte("@echo off\nset DCENV_COMMAND="+key+"\ndcenv exec \""+key+"\" %*\n"),
				0777)
		} else {
			ioutil.WriteFile(filepath.Join(shimDir, key),
				[]byte("#!/bin/bash\nexport DCENV_COMMAND="+key+"\ndcenv exec \""+key+"\" \"$@\"\n"),
				0777)
		}
		if isV {
			fmt.Println("Write command in shims.:", key)
		}
	}
	return true
}
*/
