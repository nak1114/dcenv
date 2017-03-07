package main

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/Songmu/prompter"
	baas "github.com/nak1114/dcenv/kii"
	"github.com/urfave/cli"
)

const (
	KiiAppName = "efo16zkedmd7"
	KiiAppKey  = "56b6b2e56d32447eb636f31c6e33dead"
	KiiSite    = "api-jp.kii.com"
	KiiBuckets = "dcenv"
)

var commandPush = cli.Command{
	Name:      "push",
	Aliases:   []string{"p"},
	Usage:     "Push a image config data to the resistry",
	ArgsUsage: "[options...] image_name",
	Flags: []cli.Flag{
		cli.IntFlag{
			Name:  "num, n",
			Value: 0,
			Usage: "Specify the number of the image config data",
		},
		cli.BoolFlag{
			Name:  "delete, d",
			Usage: "Delete a image config data from the registry",
		},
	},

	Action: push,
}

func push(c *cli.Context) {
	//  a := baas.NewApp("efo16zkedmd7","56b6b2e56d32447eb636f31c6e33dead","api-jp.kii.com","dcenv")
	a := baas.NewApp(KiiAppName, KiiAppKey, KiiSite, KiiBuckets)
	u, err := a.Relogin(fnameAccount())
	if err != nil {
		fmt.Println(err)
		baas.Logout(fnameAccount())
		return
	}
	u.WriteToFile(fnameAccount())

	if len(c.Args()) < 1 {
		fmt.Println("No filename.")
		cli.ShowSubcommandHelp(c)
		return
	}
	numPush := c.Int("num")
	dname, _, _ := ParseImageTag(c.Args()[0])
	if c.Bool("delete") {
		DeleteFromResistry(u, dname, numPush)
		return
	}
	PushToResistry(u, dname, numPush)
}
func fnameAccount() string {
	return filepath.Join(envHome, "files", ".kii")
}

func DeleteFromResistry(u baas.User, dname string, numPush int) {
	yp := NewYardPackFromYard(dname)
	if len(yp) <= numPush {
		fmt.Printf("Number(%d) you set is too large.Max %d.", numPush, len(yp)-1)
		fmt.Printf("Type command `dcenv yard -d %s` for more detail.\n", dname)
		return
	}
	y := &(yp[numPush])
	if y.ID != "" {
		oy := Yard{}
		if err := u.ExistObj(y.ID, &oy); err != nil {
			y.Disp("yard")
			fmt.Println(err)
			return
		}
		oy.Disp("resistry")
		if ret := prompter.YN("Delete?", false); ret {
			if err := u.DeleteObj(y.ID); err != nil {
				fmt.Println(err)
				exit(1)
				return
			}
			y.ID = ""
			y.Owner = ""
			y.Pri = 0
			yp.WriteToYard(dname)
			fmt.Println("Deleted!")
			return
		}
		fmt.Println("Revoked!")
		return
	}
	fmt.Println("Not found in the resistry.:", dname, numPush)
}

func PushToResistry(u baas.User, dname string, numPush int) {
	yp := NewYardPackFromYard(dname)
	if len(yp) <= numPush {
		fmt.Printf("Number(%d) you set is too large.Max %d.", numPush, len(yp)-1)
		fmt.Printf("Type command `dcenv yard -d %s` for more detail.\n", dname)
		return
	}
	y := &(yp[numPush])
	yb, err := json.Marshal(y)
	if err != nil {
		fmt.Println(err)
		exit(1)
		return
	}

	if y.ID != "" {
		oy := Yard{}
		if err := u.ExistObj(y.ID, &oy); err != nil {
			y.Disp("yard")
			fmt.Println(err)
			return
		}
		oy.Disp("resistry")
		y.Disp("yard")
		if ret := prompter.YN("Overwrite?", false); ret {
			if _, err := u.UpdateObj(y.ID, string(yb)); err != nil {
				fmt.Println(err)
				exit(1)
				return
			}
			fmt.Println("Overwited!")
			return
		}
		fmt.Println("Revoked!")
		return
	}
	y.ID = ""
	y.Owner = ""
	y.Pri = 0
	y.Disp("yard")
	if ret := prompter.YN("Create new config data in the resistry?", false); !ret {
		fmt.Println("Revoked!")
		return
	}
	ret, err := u.CreateObj(string(yb))
	if err != nil {
		fmt.Println(err)
		exit(1)
		return
	}
	if err := u.ACLObj(ret.ID); err != nil {
		fmt.Println(err)
		exit(1)
		return
	}
	y.ID = ret.ID
	y.Owner = u.ID
	yp.WriteToYard(dname)
	fmt.Println("Created!")
	return
}
