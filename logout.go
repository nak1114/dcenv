package main

import (
	"fmt"

	baas "github.com/nak1114/dcenv/kii"
	"github.com/urfave/cli"
)

var commandLogout = cli.Command{
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
