package main

import (
  "runtime"
  //"regexp"
  //"bufio"
  "strings"
  //"strconv"
  "os"
  //"path/filepath"
  //"sort"
  "fmt"
  "io/ioutil"

  "github.com/Songmu/prompter"
  "github.com/kr/pretty"
  "gopkg.in/yaml.v2"
)

var exit = os.Exit

type Cmder map[string]map[string]string

type Image struct {
  Tag      string `json:",omitempty"`
  Fake     bool   `json:",omitempty"`
  Script   string
  Commands Cmder
}
type ImagePack map[string]Image

type Config struct {
  Commands map[string]string
  Images   ImagePack
}

func ParseImageTag(s string) (string, string, string) {
  c := strings.LastIndex(s, "/")
  t := strings.LastIndex(s, ":")
  if c < t {
    return s[:t], s[c+1 : t], s[t+1:]
  } else {
    return s, s[c+1:], ""
  }
}

func CheckCmdName(s string) bool {
  if runtime.GOOS == `windows` {
    return strings.IndexAny(s, `\/:*?"<>|`) < 0
  } else {
    return strings.IndexAny(s, `/`) < 0
  }
}

func (m *Image) MakeCommands(cmds []string, envs []string) {
  c := make(Cmder)
  e := make(map[string]string)
  for _, env := range envs {
    t := strings.SplitN(env, `=`, 2)
    if len(t) == 2 {
      e[t[0]] = t[1]
    }
  }
  for _, cmd := range cmds {
    if CheckCmdName(cmd) {
      c[cmd] = e
    }
  }
  m.Commands = c
  return
}
func (tc *Image) ValidCommands() {
  //normarize command
  for cmd, _ := range tc.Commands {
    if !CheckCmdName(cmd) {
      fmt.Println("Invalid command. Removed:", cmd)
      delete(tc.Commands, cmd)
    }
  }
  return
}

func SearchImageFromYard(dname string, tCommand string, tTag string, cnum int) (tc Image) {
  yp := NewYardPackFromYard(dname)
  if len(yp) == 0 {
    fmt.Printf("[%s] config data were not found in the yard.\n", len(yp), dname)
    fmt.Printf("Please command 'dcenv yard --list %s'\n", dname)
    exit(1)
    return
  }
  if isV {
    fmt.Printf("%d [%s] config data were found in the yard.\n", len(yp), dname)
  }
  if len(yp) <= cnum {
    fmt.Printf("%d [%s] config data were found in the yard.\n", len(yp), dname)
    fmt.Printf("Too large number you choice.:(%d)\n", cnum)
    fmt.Printf("Please command 'dcenv yard -d %s'\n", dname)
    exit(1)
    return
  }
  for _, val := range yp[cnum:] {
    ok := true
    tc, ok = val.Config[envShell]
    if ok {
      //normarize image
      if tc.Script == "" {
        continue
      }

      //Commands
      if len(tc.Commands) == 0 {
        tc.MakeCommands([]string{tCommand}, []string{})
      } else {
        tc.ValidCommands()
      }
      //Tag
      if tTag != "" {
        tc.Tag = tTag
      } else if tc.Tag == "" {
        tc.Tag = "latest"
      }
      return
    }
  }
  return
}

func GetConfig(fname string) (m Config) {
  //m := Config{}
  //read config file
  if _, err := os.Stat(fname); err == nil {
    if isV {
      fmt.Println("Found the config file.:", fname)
    }
    buf, err := ioutil.ReadFile(fname)
    if err != nil {
      fmt.Println("File can not read.:", fname)
      os.Exit(1)
      return
    }
    if err = yaml.Unmarshal([]byte(buf), &m); err != nil {
      fmt.Println(err)
      os.Exit(1)
      return
    }
  } else {
    if isV {
      fmt.Println("Not found the config file.Create it.:", fname)
    }
    m.Images = make(ImagePack)
  }
  fmt.Println("Config file.:", fname)
  return
}

func (m Config) AddImage(c Image, name string, isForce bool) {
  //check exists image
  if val, ok := m.Images[name]; ok {
    if !isForce {
      pretty.Printf("--- cur %s:\n%# v\n\n", name, val)
      pretty.Printf("--- new %s:\n%# v\n\n", name, c)
      fmt.Printf("\nFound same image[%s]\n\n", name)
      if ret := prompter.YN("Overwrite?", true); !ret { //askForConfirmation("Overwrite?"); !ret{
        exit(1)
        return
      }
    }
  }
  //insert image
  m.Images[name] = c
  return
}

func (m Config) DelImage(name string, isForce bool) {
  //check exists image
  if _, ok := m.Images[name]; ok {
    //insert image
    delete(m.Images, name)
  } else {
    fmt.Println("Image not found.:", name)
    exit(1)
    return
  }
  return
}

func (m *Config) WriteToFile(fname string) {
  if len((*m).Images) == 0 {
    if err := os.Remove(fname); err != nil {
      exit(1)
      fmt.Println(err)
    }
    if isV {
      fmt.Println("Delete config file. :", fname)
    }
    return
  }
  //renewal commands
  cc := make(map[string]string)
  for key, val := range (*m).Images {
    for com, _ := range val.Commands {
      if cnt, ok := cc[com]; ok {
        fmt.Println("Confrict command name.:", com)
        fmt.Printf("  image [%s] and [%s]\n", cnt, key)
        fmt.Printf("  in file:%s\n", fname)
        exit(1)
        return
      }
      cc[com] = key
    }
  }
  (*m).Commands = cc

  //write config file
  buf, err := yaml.Marshal(m)
  if err != nil {
    fmt.Println(err)
    exit(1)
    return
  }

  err = ioutil.WriteFile(fname, buf, 0666)
  if err != nil {
    fmt.Println(err)
    exit(1)
    return
  }
  if isV {
    fmt.Println("Wrote config file. :", fname)
  }
}

func (ip *ImagePack) WriteToFile(fname string) {
  //write config file
  buf, err := yaml.Marshal(ip)
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
    return
  }
  err = ioutil.WriteFile(fname, buf, 0666)
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
    return
  }
}
func NewImagePackFromFile(fname string) (m ImagePack) {
  buf, err := ioutil.ReadFile(fname)
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
    return
  }
  if err = yaml.Unmarshal([]byte(buf), &m); err != nil {
    fmt.Println(err)
    os.Exit(1)
    return
  }
  return
}

/*
// askForConfirmation asks the user for confirmation. A user must type in "yes" or "no" and
// then press enter. It has fuzzy matching, so "y", "Y", "yes", "YES", and "Yes" all count as
// confirmations. If the input is not recognized, it will ask again. The function does not return
// until it gets a valid response from the user.
func askForConfirmation(s string) bool {
  reader := bufio.NewReader(os.Stdin)

  for {
    fmt.Printf("%s [y/n]: ", s)

    response, err := reader.ReadString('\n')
    if err != nil {
      fmt.Println(err)
      os.Exit(1)
    }

    response = strings.ToLower(strings.TrimSpace(response))

    if response == "y" || response == "yes" {
      return true
    } else if response == "n" || response == "no" {
      return false
    }
  }
}
*/
