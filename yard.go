package main

import (
  "strconv"
  "strings"
  "fmt"
  "path/filepath"
  "net/url"

  "github.com/Songmu/prompter"
  pty "github.com/kr/pretty"
  "github.com/urfave/cli"
)

var Command_yard = cli.Command{
  Name:      "yard",
  Aliases:   []string{"y"},
  Usage:     "Control your yard",
  ArgsUsage: "[options...] image_name",
  Flags: []cli.Flag{
    cli.BoolFlag{
      Name:  "display, d",
      Usage: "Show a detail fo the image config data",
    },
    cli.BoolFlag{
      Name:  "list, l",
      Usage: "Show the list of all image config data in your yard",
    },
    cli.BoolFlag{
      Name:  "create, c",
      Usage: "Create new image config data in your yard",
    },
  },

  Action: yard,
}

func yard(c *cli.Context) {
  if isV {
    fmt.Println("dcenv yard ", c.Args())
  }
  if c.Bool("list") {
    ListYardAll()
    return
  }
  if len(c.Args()) < 1 {
    fmt.Println("No image name.")
    cli.ShowSubcommandHelp(c)
    return
  }
  dname, _, _ := ParseImageTag(c.Args()[0])
  if c.Bool("create") {
    CreateNewYard(dname)
    return
  }
  ListYardDetail(dname)
  return

}

func ListYardAll() {
  //  fname:=filepath.Join(envHome ,"image_yard",url.QueryEscape(dname)+".yml")
  if filenames, err := filepath.Glob(filepath.Join(envHome, "image_yard", "*.yml")); err == nil {
    for _, val := range filenames {
      str, err := url.QueryUnescape(strings.TrimSuffix(filepath.Base(val), ".yml"))
      if err != nil {
        fmt.Println(err)
      } else {
        fmt.Println(str)
      }
    }
  }
}

func (yd *Yard) Disp(i string) {
  fmt.Printf("---[ %s ]---------\n", i)
  fmt.Println("  Id      :", yd.Id)
  fmt.Println("  Owner   :", yd.Owner)
  //  fmt.Println("   Createat:",yd.CreateAt)
  fmt.Println("  Image   :", yd.Image)
  fmt.Println("  Brief   :", yd.Brief)
  fmt.Println("  Desc    :", yd.Desc)
  fmt.Println("  Pri     :", yd.Pri)
  for key, cnt := range yd.Config {
    fmt.Println("  Config  :", key)
    fmt.Println("    Tag     :", cnt.Tag)
    fmt.Println("    Fake    :", cnt.Fake)
    pty.Printf("    Commands: %v\n", cnt.Commands)
    fmt.Printf("    Script  :\n%s\n", cnt.Script)
  }
  return
}
func ListYardDetail(dname string) {
  yp := NewYardPackFromYard(dname)
  for i, yd := range yp {
    yd.Disp(strconv.Itoa(i))
  }
}
func CreateNewYard(dname string) {
  yp := NewYardPackFromYard(dname)
  if len(yp) > 0 {
    fmt.Printf("*** A file already exists. ***\n")
    ListYardDetail(dname)
    if ret := prompter.YN("Erase it?", false); !ret {
      exit(1)
      return
    }
  }
  nyp := YardPack{}
  y := Yard{
    Image: dname,            //    string    `json:"image,omitempty"`
    Brief: "Brief " + dname, //    string    `json:"brief,omitempty"`
  }
  y.Config = make(ImagePack)
  y.Config["bash"] = Image{}
  y.Config["windows"] = Image{}
  nyp.AddItem(y)
  f:=nyp.WriteToYard(dname)
  fmt.Println("Edit:", f)
  return
}
