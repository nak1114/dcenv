package main

import (
  "fmt"
  "github.com/urfave/cli"
  //"github.com/Songmu/prompter"
  baas "github.com/nak1114/dcenv/kii"
)

var Command_logout = cli.Command{
  Name:      "logout",
  Usage:     "Log out from the registry",
  ArgsUsage: "[options...]",
  //  Flags:   []cli.Flag{
  //  },

  Action: logout,
}

func logout(c *cli.Context) {
  baas.Logout(fnameAccount())
  fmt.Println("Logout!")
}
