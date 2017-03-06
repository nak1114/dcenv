package main

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/urfave/cli"
)

var Version = "0.0.1"
var isV = false
var isD = false
var envShell = "bash"
var envHome = ""
var envCommand = ""

func main() {
	app := cli.NewApp()
	//app.Name = "dcenv"
	app.Usage = "Simple Docker Environment Management"
	app.Version = Version

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "dryrun, d",
			Usage:       "Don't create/write any file and io.(not impl yet)",
			Destination: &isD,
		},
		cli.BoolFlag{
			Name:        "verbose, V",
			Usage:       "print verbose mode.",
			Destination: &isV,
		},
	}

	app.Commands = []cli.Command{
		commandGlobal,
		commandExec,
		commandRehash,
		commandLocal,
		commandSystem,
		commandCommand,
		commandInstall,
		commandUninstall,
		commandTag,
		commandYard,
		commandPush,
		commandLogout,
	}

	envCommand = os.Getenv("DCENV_COMMAND")
	if runtime.GOOS == `windows` {
		envCommand = os.Getenv("DCENV_COMMAND") + `.bat`
	}
	envShell = os.Getenv("DCENV_SHELL")
	if envShell == "" {
		if runtime.GOOS == `windows` {
			envShell = `windows`
		} else {
			envShell = `bash`
		}
	}
	envHome = os.Getenv("DCENV_HOME")
	if envHome == "" {
		str, err := os.Executable()
		if err != nil {
			exit(1)
			return
		}
		envHome = strings.TrimSuffix(filepath.Dir(str), string(os.PathSeparator)+"files")
	}

	app.Run(os.Args)
}
