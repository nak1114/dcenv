package main

import (
  "fmt"
  "io/ioutil"
  "os"
  "path/filepath"
  "runtime"

  "github.com/urfave/cli"
  "gopkg.in/yaml.v2"
)

var Command_rehash = cli.Command{
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
    os.RemoveAll(sdir)
  }
  os.RemoveAll(tdir)
  if c.Bool("clean") {
    return
  }
  os.Mkdir(sdir, 0777)
  os.Mkdir(tdir, 0777)
  SearchConfig("", makeShimCommands)
}

func makeShimCommands(cmd string, fname string) *Config {
  makeShimCommand(fname)
  return nil
}
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
  buf, err := ioutil.ReadFile(fname)
  if err != nil {
    if isV {
      fmt.Println("  File can not read.:", fname)
    }
    return false
  }
  m := Config{}
  if err = yaml.Unmarshal([]byte(buf), &m); err != nil {
    if isV {
      fmt.Println("  File can not unmarshal.:", fname)
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
  for key, _ := range m.Commands {
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
