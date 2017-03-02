package main

import (
  "fmt"
  "io/ioutil"
  "os"
  "runtime"
  "strings"
  "text/template"

  "gopkg.in/yaml.v2"
  "path/filepath"
)

func CheckCommand(cmd string, fname string) *Config {
  if isV {
    fmt.Println("search:", fname)
  }
  if _, err := os.Stat(fname); err != nil {
    return nil
  }
  if isV {
    fmt.Println("  Found the config file.:", fname)
  }
  buf, err := ioutil.ReadFile(fname)
  if err != nil {
    if isV {
      fmt.Println("  File can not read.:", fname)
    }
    return nil
  }

  m := Config{}
  if err = yaml.Unmarshal([]byte(buf), &m); err != nil {
    if isV {
      fmt.Println("  File can not unmarshal.:", fname)
    }
    return nil
  }
  cval, ok := m.Commands[cmd]
  if !ok {
    if isV {
      fmt.Println("  Command not found.:", fname)
    }
    return nil
  }
  if isV {
    fmt.Println("found.:", cval)
  }
  return &m
}

func CheckImage(cmd string, fname string) *Config {
  if isV {
    fmt.Println("search:", fname)
  }
  if _, err := os.Stat(fname); err != nil {
    return nil
  }
  if isV {
    fmt.Println("  Found the config file.:", fname)
  }
  buf, err := ioutil.ReadFile(fname)
  if err != nil {
    if isV {
      fmt.Println("  File can not read.:", fname)
    }
    return nil
  }

  m := Config{}
  if err = yaml.Unmarshal([]byte(buf), &m); err != nil {
    if isV {
      fmt.Println("  File can not unmarshal.:", fname)
    }
    return nil
  }
  cval, ok := m.Images[cmd]
  if !ok {
    if isV {
      fmt.Println("  Image not found.:", fname)
    }
    return nil
  }
  if isV {
    fmt.Println("found.:", cval)
  }
  return &m
}

var CheckCmdsMap = make(map[string][]string)

func InitCheckCmds() {
  CheckCmdsMap = make(map[string][]string)
}

func CheckCmds(cmd string, fname string) *Config {
  if isV {
    fmt.Println("search:", fname)
  }
  if _, err := os.Stat(fname); err != nil {
    return nil
  }
  if isV {
    fmt.Println("  Found the config file.:", fname)
  }
  buf, err := ioutil.ReadFile(fname)
  if err != nil {
    if isV {
      fmt.Println("  File can not read.:", fname)
    }
    return nil
  }

  m := Config{}
  if err = yaml.Unmarshal([]byte(buf), &m); err != nil {
    if isV {
      fmt.Println("  File can not unmarshal.:", fname)
    }
    return nil
  }
  for key, val := range m.Commands {
    _, ok := CheckCmdsMap[key]
    if !ok {
      CheckCmdsMap[key] = []string{val, fname}
    }
  }

  return nil
}

func SearchConfig(cmd string, f func(string, string) *Config) (*Config, string) {
  fname := `.dcenv_` + envShell

  p := os.Getenv("DCENV_DIR")
  if p != "" {
    t := strings.Split(p, string(os.PathSeparator))
    for i := len(t); i > 0; i -= 1 {
      s := strings.Join(t[0:i], string(os.PathSeparator)) + string(os.PathSeparator) + fname
      tc := f(cmd, s)
      if tc != nil {
        return tc, s
      }
    }
  }

  p, err := os.Getwd()
  if err != nil {
    fmt.Println(err)
    exit(1)
    return nil, p
  }
  t := strings.Split(p, string(os.PathSeparator))
  for i := len(t); i > 0; i -= 1 {
    s := strings.Join(t[0:i], string(os.PathSeparator)) + string(os.PathSeparator) + fname
    tc := f(cmd, s)
    if tc != nil {
      return tc, s
    }
  }

  s := filepath.Join(envHome, "files", fname)
  tc := f(cmd, s)
  if tc != nil {
    return tc, s
  }

  return nil, s
}

type scriptArg struct {
  Tag    string
  Image  string
  Cmd    string
  CfgDir string
  Envs   map[string]string
}

func MakeArgsFile(sft int) {
  s := ""
  for _, t := range os.Args[sft:] {
    if strings.Contains(t, " ") {
      s += ` "` + t + `"`
    } else {
      s += " " + t
    }
  }
  //  if len(os.Args[sft:]) > 0 {
  //    s = `"`+strings.Join(os.Args[sft:],`" "`)+`"`
  //  }
  fname := filepath.Join(envHome, "files", "__args__")
  err := ioutil.WriteFile(fname, []byte(s), 0666)
  if err != nil {
    fmt.Println(err)
    exit(1)
    return
  }
}
func MakeExecFile(m *Config, cmd string, fname string) bool {
  //make script
  cval := (*m).Commands[cmd]
  cnt := (*m).Images[cval]

  sarg := scriptArg{
    Tag:    cnt.Tag,
    Image:  cval,
    Cmd:    cmd,
    CfgDir: filepath.Dir(fname),
    Envs:   cnt.Commands[cmd],
  }

  tf := template.FuncMap{
  //    "tocmd": func(s string) string { sl:=strings.Split(s,"/");return sl[len(sl)-1] },
  //    "todir": func(s string) string { sl:=strings.Split(s,"/");return strings.Join(sl[:len(sl)-1],"/") },
  }
  hname := `header_` + envShell
  buf, err := ioutil.ReadFile(filepath.Join(envHome, "files", hname))
  if err != nil {
    fmt.Println(err)
    exit(1)
    return false
  }
  tpl := template.Must(template.New("script").Funcs(tf).Parse(string(buf) + cnt.Script))
  uname := filepath.Join(envHome, "tmp", envCommand)
  //fp=os.Stdout
  fp, err := os.OpenFile(uname, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
  if err != nil {
    fmt.Println(err)
    exit(1)
    return false
  }
  if err := tpl.Execute(fp, sarg); err != nil {
    fmt.Println(err, hname)
    exit(1)
    return false
  }

  if err := fp.Close(); err != nil {
    fmt.Println(err)
    exit(1)
    return false
  }
  if isV {
    fmt.Println("write file:", uname)
  }
  return true
}

func MakeExecFileSystem(cmd string) {
  s := filepath.Join(envHome, "shims")
  plist := filepath.SplitList(os.Getenv("PATH"))
  for i, p := range plist {
    if p == s {
      plist = append(plist[0:i], plist[i+1:]...)
      break
    }
  }
  ret := strings.Join(plist, string(os.PathListSeparator))

  if runtime.GOOS == `windows` {
    ioutil.WriteFile(filepath.Join(envHome, "tmp", envCommand),
      []byte("@setlocal\n@set PATH="+ret+"\n\""+cmd+"\" %*\n"),
      0777)
  } else {
    ioutil.WriteFile(filepath.Join(envHome, "tmp", envCommand),
      []byte("#!/bin/bash\nexport PATH="+ret+"\n\""+cmd+"\" \"$@\"\n"),
      0777)
  }
  return
}

func ShowExecFile(m *Config, cmd string, fname string) bool {
  //make script
  cval := (*m).Commands[cmd]
  cnt := (*m).Images[cval]

  sarg := scriptArg{
    Tag:    cnt.Tag,
    Image:  cval,
    Cmd:    cmd,
    CfgDir: filepath.Dir(fname),
    Envs:   cnt.Commands[cmd],
  }

  tf := template.FuncMap{
  //    "tocmd": func(s string) string { sl:=strings.Split(s,"/");return sl[len(sl)-1] },
  //    "todir": func(s string) string { sl:=strings.Split(s,"/");return strings.Join(sl[:len(sl)-1],"/") },
  }
  hname := `header_` + envShell
  buf, err := ioutil.ReadFile(filepath.Join(envHome, "files", hname))
  if err != nil {
    fmt.Println(err)
    exit(1)
    return false
  }
  tpl := template.Must(template.New("script").Funcs(tf).Parse(string(buf) + cnt.Script))
  if err := tpl.Execute(os.Stdout, sarg); err != nil {
    fmt.Println(err, hname)
    exit(1)
    return false
  }
  return true
}
