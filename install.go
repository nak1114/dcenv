package main

import (
  "sort"
  //"bytes"
  //"strings"
  "encoding/json"
  "fmt"
  "io/ioutil"
  "os"
  "path/filepath"
  //"net/http"
  "net/url"

  "github.com/urfave/cli"
  "gopkg.in/yaml.v2"
  //"github.com/kr/pretty"

  baas "github.com/nak1114/dcenv/kii"
)

type Yard struct {
  Id     string    `json:"_id,omitempty"`
  Owner  string    `json:"_owner,omitempty"`
  Image  string    `json:"image,omitempty"`
  Brief  string    `json:"brief,omitempty"`
  Desc   string    `json:"desc,omitempty"`
  Pri    int32     `json:"pri,omitempty"`
  Config ImagePack `json:"config,omitempty"`
}
type YardPack []Yard

var Command_install = cli.Command{
  Name:      "install",
  Aliases:   []string{"i"},
  Usage:     "Install a config file to your yard",
  ArgsUsage: "[options...] image_name",
  Flags: []cli.Flag{
    cli.StringFlag{
      Name:  "export, e",
      Usage: "Export ajson file from the yard to `FILE`",
    },
    cli.StringFlag{
      Name:  "import, i",
      Usage: "Import a json file from `FILE` to the yard",
    },
    cli.BoolFlag{
      Name:  "list, l",
      Usage: "Show a list of all image data in the dcenv registry",
    },
  },

  Action: install,
}

func install(c *cli.Context) {
  //isF := c.Bool("force")

  if isV {
    fmt.Println("dcenv install ", c.Args())
  }
  if c.Bool("list") {
    //ListYardAll()
    ListRepoAll()
    return
  }
  importf := c.String("import")
  if importf != "" {
    yp := NewYardPackFromFile(importf)
    if isV {
      fmt.Println("Read import file.:", importf)
    }
    for i := len(yp); i > 0; i -= 1 {
      dname := yp[i-1].Image
      yd := NewYardPackFromYard(dname)
      yd.AddItem(yp[i-1])
      yd.WriteToYard(dname)
      if isV {
        fmt.Println("Add image to yard.:", dname)
        fmt.Println("  Total images in yard.:", len(yd))
      }
    }
    fmt.Println("Import!")
    return
  }
  if len(c.Args()) < 1 {
    fmt.Println("No image name.")
    cli.ShowSubcommandHelp(c)
    return
  }
  dname, _, _ := ParseImageTag(c.Args()[0])
  export := c.String("export")
  if export != "" {
    yp := NewYardPackFromYard(dname)
    if len(yp) == 0 {
      fmt.Println("Image not found in yard.")
      exit(1)
      return
    }
    yp.WriteToFile(export)
    fmt.Println("Export!")
    return
  }

  yp := NewYardPackFromKii(dname)
  yd := NewYardPackFromYard(dname)
  for _, val := range yp {
    yd.AddItem(val)
  }
  yd.WriteToYard(dname)
  fmt.Println("Add to your yard.")
  return

}

func NewYardPackFromKii(dname string) (yp YardPack) {
  a := baas.NewApp(KiiAppName, KiiAppKey, KiiSite, KiiBuckets)
  q, err := a.Query(fmt.Sprintf(`{"clause": {"type": "eq", "field": "image", "value": "%s"}}`, dname))
  if err != nil {
    fmt.Println(err)
    exit(1)
    return
  }
  for !q.EoQ {
    t := Yard{}
    if err := q.Next(&t); err != nil {
      fmt.Println(err)
      exit(1)
      return
    }
    t.Pri = 0
    yp = append(yp, t)
    //pretty.Printf("--- cur:%s\n%# v\n\n",t.Image,t)
  }
  return
}

func NewYardPackFromFile(fname string) (yp YardPack) {
  buf, err := ioutil.ReadFile(fname)
  if err != nil {
    fmt.Println(err)
    exit(1)
    return
  }
  if err := json.Unmarshal(buf, &yp); err != nil {
    fmt.Println(err)
    exit(1)
    return
  }
  return
}
func NewYardPackFromYard(dname string) (yp YardPack) {
  fname := filepath.Join(envHome, "image_yard", url.QueryEscape(dname)+".yml")
  if _, err := os.Stat(fname); err != nil {
    return
  }
  buf, err := ioutil.ReadFile(fname)
  if err != nil {
    fmt.Println(err)
    exit(1)
    return
  }
  if err = yaml.Unmarshal([]byte(buf), &yp); err != nil {
    fmt.Println(err)
    exit(1)
    return
  }
  sort.SliceStable(yp, func(i, j int) bool {
    return yp[i].Id < yp[j].Id
  })
  return
}
func (yp YardPack) WriteToFile(fname string) {
  b, err := json.MarshalIndent(yp, "", "  ")
  if err != nil {
    fmt.Println("error:", err)
    exit(1)
    return
  }
  err = ioutil.WriteFile(fname, b, 0666)
  if err != nil {
    fmt.Println(err)
    exit(1)
    return
  }
}
func (yp *YardPack) AddItem(y Yard) {
  if y.Id == "" {
    (*yp) = append((*yp), y)
    sort.SliceStable((*yp), func(i, j int) bool {
      return (*yp)[i].Id < (*yp)[j].Id
    })
    return
  }
  i := sort.Search(len((*yp)), func(i int) bool { return (*yp)[i].Id >= y.Id })
  if i >= len((*yp)) {
    (*yp) = append((*yp), y)
    return
  }
  if (*yp)[i].Id == y.Id {
    (*yp)[i] = y
    return
  }
  (*yp) = append((*yp), y)
  sort.SliceStable(yp, func(i, j int) bool {
    return (*yp)[i].Id < (*yp)[j].Id
  })
  return
}

func (yp YardPack) WriteToYard(dname string)(fname string) {
  fname = filepath.Join(envHome, "image_yard", url.QueryEscape(dname)+".yml")
  b, err := yaml.Marshal(yp)
  if err != nil {
    fmt.Println(err)
    exit(1)
    return
  }
  err = ioutil.WriteFile(fname, b, 0666)
  if err != nil {
    fmt.Println(err)
    exit(1)
    return
  }
  return
}

func ListRepoAll() {
  a := baas.NewApp(KiiAppName, KiiAppKey, KiiSite, KiiBuckets)
  q, err := a.Query(`{"clause": 
                      {"type": "hasField", "field": "image", "fieldType": "STRING"
                      },
                      "orderBy": "image", "descending": false
                   }`)
  if err != nil {
    fmt.Println(err)
    exit(1)
    return
  }
  ret := ""
  for !q.EoQ {
    t := Yard{}
    if err := q.Next(&t); err != nil {
      fmt.Println(err)
      exit(1)
      return
    }
    if ret != t.Image {
      ret = t.Image
      fmt.Println(ret)
    }
    //pretty.Printf("--- cur:%s\n%# v\n\n",t.Image,t)
  }
}
