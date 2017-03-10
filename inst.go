package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"

	"github.com/kr/pretty"

	"gopkg.in/yaml.v2"
)

// Cmder are command and command's environments
type Cmder map[string]map[string]string

// Image is a docker image config.
type Image struct {
	Tag      string `json:",omitempty"`
	Fake     bool   `json:",omitempty"`
	Script   string
	Commands Cmder
}

// ImagePack are some docker image config.
type ImagePack map[string]Image

//Config is a config file format.
type Config struct {
	Commands map[string]string
	Images   ImagePack
}

// ParseImageTag is split imagetag name and tag name.
// s input image name&Tag
// return fullname,imagename,tag
func ParseImageTag(s string) (string, string, string) {
	c := strings.LastIndex(s, "/")
	t := strings.LastIndex(s, ":")
	if c < t {
		return s[:t], s[c+1 : t], s[t+1:]
	}
	return s, s[c+1:], ""
}

// CheckCmdName lints a command name.
// s command name
// return command name can use filename?
func CheckCmdName(s string) bool {
	if runtime.GOOS == `windows` {
		return strings.IndexAny(s, `\/:*?"<>|`) < 0
	}
	return strings.IndexAny(s, `/`) < 0
}

// MakeCommands makes some original commands & envs.
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

// ValidCommands lints all command name and remove it.
func (m *Image) ValidCommands() {
	//normarize command
	for cmd := range m.Commands {
		if !CheckCmdName(cmd) {
			fmt.Println("Invalid command. Removed:", cmd)
			delete(m.Commands, cmd)
		}
	}
	return
}

// SearchImageFromYard find some 'dname' image in your yard.and select a 'cnum' config set.at last, normalize all commands and tag.
func SearchImageFromYard(dname string, tCommand string, tTag string, cnum int) (tc Image) {
	yp := NewYardPackFromYard(dname)
	if len(yp) == 0 {
		fmt.Printf("[%s] config data were not found in the yard.\n", dname)
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
	fmt.Printf("Not found command for this shell.:%s\n", envShell)
	fmt.Printf("Please command 'dcenv yard -d %s'\n", dname)
	exit(1)
	return
}

// LoadConfig gets a 'Config' from the 'fname' file.
func LoadConfig(fname string) (m *Config) {
	m = &Config{}
	if err := LoadYaml(m, fname); err != nil {
		if isV {
			fmt.Println("  File can not unmarshal.:", fname)
			fmt.Println(err)
		}
		return nil
	}
	return
}

// NewConfig finds and creates a 'Config' from the 'fname' file.
func NewConfig(fname string) (m Config) {
	//read config file
	if _, err := os.Stat(fname); err == nil {
		if isV {
			fmt.Println("Found the config file.:", fname)
		}
		if err := LoadYaml(&m, fname); err != nil {
			fmt.Println(err)
			exit(1)
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

// AddImage add a image config data to itself.
func (m Config) AddImage(c Image, name string, isForce bool) {
	//check exists image
	if val, ok := m.Images[name]; ok {
		if !isForce {
			pretty.Printf("--- cur %s:\n%# v\n\n", name, val)
			pretty.Printf("--- new %s:\n%# v\n\n", name, c)
			fmt.Printf("\nFound same image[%s]\n\n", name)
			if ret := prompterYN("Overwrite?", false); !ret { //askForConfirmation("Overwrite?"); !ret {
				exit(1)
				return
			}
		}
	}
	//insert image
	m.Images[name] = c
	return
}

// DelImage delete a image config data from itself.
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

// WriteToFile write itself to file with renewal 'commands' field.
func (m *Config) WriteToFile(fname string) {
	if len((*m).Images) == 0 {
		if err := DeleteYaml(fname); err != nil {
			fmt.Println(err)
			exit(1)
		}
		if isV {
			fmt.Println("Delete config file. :", fname)
		}
		return
	}
	//renewal commands
	cc := make(map[string]string)
	for key, val := range (*m).Images {
		for com := range val.Commands {
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
	if err := SaveYaml(m, fname); err != nil {
		fmt.Println(err)
		exit(1)
		return
	}
	if isV {
		fmt.Println("Wrote config file. :", fname)
	}
}

// LoadYaml loads a yaml file to object.
func LoadYaml(v interface{}, fname string) (err error) {
	buf, err := ioutil.ReadFile(fname)
	if err == nil {
		err = yaml.Unmarshal([]byte(buf), v)
	}
	return
}

// SaveYaml saves a yaml file from object.
func SaveYaml(v interface{}, fname string) (err error) {
	buf, err := yaml.Marshal(v)
	if err == nil {
		err = ioutil.WriteFile(fname, buf, 0666)
	}
	return
}

// DeleteYaml delete a file.
func DeleteYaml(fname string) (err error) {
	err = os.Remove(fname)
	return
}
