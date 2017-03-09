package main

import (
	"runtime"
	"testing"

	"fmt"

	a "github.com/nak1114/goutil/assert"
)

func TestParseImageTag(t *testing.T) {
	a.Set(t, "TestParseImageTag")

	cnt, com, tag := ParseImageTag("local:3000/hoge-fuga/app:1.0.0-alpine")
	a.Eq(cnt, "local:3000/hoge-fuga/app")
	a.Eq(com, "app")
	a.Eq(tag, "1.0.0-alpine")

	cnt, com, tag = ParseImageTag("app:1.0.0-alpine")
	a.Eq(cnt, "app")
	a.Eq(com, "app")
	a.Eq(tag, "1.0.0-alpine")

	cnt, com, tag = ParseImageTag("app")
	a.Eq(cnt, "app")
	a.Eq(com, "app")
	a.Eq(tag, "")

	cnt, com, tag = ParseImageTag("local:3000/hoge-fuga/app")
	a.Eq(cnt, "local:3000/hoge-fuga/app")
	a.Eq(com, "app")
	a.Eq(tag, "")

	cnt, com, tag = ParseImageTag("hoge-fuga/app")
	a.Eq(cnt, "hoge-fuga/app")
	a.Eq(com, "app")
	a.Eq(tag, "")

	cnt, com, tag = ParseImageTag("hoge-fuga/")
	a.Eq(cnt, "hoge-fuga/")
	a.Eq(com, "")
	a.Eq(tag, "")

	cnt, com, tag = ParseImageTag(":1.0.0-alpine")
	a.Eq(cnt, "")
	a.Eq(com, "")
	a.Eq(tag, "1.0.0-alpine")

	cnt, com, tag = ParseImageTag("local:3000/hoge-fuga/:1.0.0-alpine")
	a.Eq(cnt, "local:3000/hoge-fuga/")
	a.Eq(com, "")
	a.Eq(tag, "1.0.0-alpine")

}

func TestCheckCmdName(t *testing.T) {
	a.Set(t, "TestCheckCmdName")
	isLinux := true
	if runtime.GOOS == `windows` {
		isLinux = false
	}

	a.Eq(CheckCmdName("aaa"), true)
	a.Eq(CheckCmdName("aaa.bat"), true)
	a.Eq(CheckCmdName("a/aa"), false)
	a.Eq(CheckCmdName(`a\aa`), isLinux)
	a.Eq(CheckCmdName(`a<aa`), isLinux)
	a.Eq(CheckCmdName(`a>aa`), isLinux)
	a.Eq(CheckCmdName(`|aaa`), isLinux)
	a.Eq(CheckCmdName(`aa|a`), isLinux)
	a.Eq(CheckCmdName(`aa*a`), isLinux)
	a.Eq(CheckCmdName(`aa:a`), isLinux)
	a.Eq(CheckCmdName(`aa"a`), isLinux)
}

func TestMakeCommands(t *testing.T) {
	a.Set(t, "TestMakeCommands")
	m := Image{}
	m.MakeCommands([]string{"test", "hoge"}, []string{"aaa=bbb", "ccc=ddd"})
	exp := Cmder{
		"test": map[string]string{
			"aaa": "bbb",
			"ccc": "ddd",
		},
		"hoge": map[string]string{
			"aaa": "bbb",
			"ccc": "ddd",
		},
	}
	a.Eq(m.Commands["test"]["aaa"], exp["test"]["aaa"])
	a.Eq(m.Commands["test"]["ccc"], exp["test"]["ccc"])
	a.Eq(m.Commands["hoge"]["aaa"], exp["hoge"]["aaa"])
	a.Eq(m.Commands["hoge"]["ccc"], exp["hoge"]["ccc"])
	a.Eq(len(m.Commands), 2)

	m.MakeCommands([]string{}, []string{"aaa=bbb", "ccc=ddd"})
	a.Eq(len(m.Commands), 0)

	m.Commands = Cmder{
		"te/st": map[string]string{},
		"hoge":  map[string]string{},
	}

	a.Eq(len(m.Commands), 2)
	m.ValidCommands()
	a.Eq(len(m.Commands), 1)

}

type TestP int

func (p TestP) MarshalYAML() (interface{}, error) {
	return nil, fmt.Errorf("test")
}

func TestYaml(t *testing.T) {
	a.Set(t, "TestYaml")
	ret := make(map[string]string)
	tmp := map[string]string{
		"hoge": "fuga",
		"foo":  "bar",
	}
	p := TestP(0)
	e := SaveYaml(tmp, "./misc/tmp/test.yml")
	a.Eq(e, error(nil))
	e = LoadYaml(&ret, "./misc/tmp/test.yml")
	a.Eq(e, error(nil))
	e = DeleteYaml("./misc/tmp/test.yml")
	a.Eq(e, error(nil))
	a.Eq(ret["hoge"], tmp["hoge"])
	e = LoadYaml(&ret, "./misc/tmp/test.yml")
	a.Neq(e, error(nil))
	e = SaveYaml(p, "./misc/tmp/test.yml")
	a.Neq(e, error(nil))

}

func TestSearchImageFromYard(t *testing.T) {
	a.Set(t, "SearchImageFromYard")
	orgExit := exit
	exit = pseudoExit
	flgExit = false
	envHome = "./misc/test/dcenv_home"
	envShell = "bash"
	tc := Image{}

	outb, _ := a.StubIO("", func() {
		tc = SearchImageFromYard("naktak/inst-test-null", "inst-test-null", "", 0)
	})
	outexpect := `[naktak/inst-test-null] config data were not found in the yard.
Please command 'dcenv yard --list naktak/inst-test-null'
`
	a.Eq(flgExit, true)
	a.Eq(tc.Tag, "")
	a.Eq(outb, outexpect)

	flgExit = false
	tc = Image{}
	outb, _ = a.StubIO("", func() {
		tc = SearchImageFromYard("naktak/inst-test", "inst-test-null", "", 20)
	})
	outexpect = `5 [naktak/inst-test] config data were found in the yard.
Too large number you choice.:(20)
Please command 'dcenv yard -d naktak/inst-test'
`
	a.Eq(flgExit, true)
	a.Eq(tc.Tag, "")
	a.Eq(outb, outexpect)

	exit = orgExit

	tc = Image{}
	tc = SearchImageFromYard("naktak/inst-test", "inst-test", "", 0)
	a.Eq(tc.Tag, "latest")
	a.Eq(tc.Script, `echo "hogeid2 {{.CfgDir}} {{.Image}}:{{.Tag}} {{.Cmd}}"`+"\n")

	tc = Image{}
	tc = SearchImageFromYard("naktak/inst-test", "inst-test", "1.1.1", 3)
	a.Eq(tc.Tag, "1.1.1")
	a.Eq(tc.Script, `echo "hogeid3 {{.CfgDir}} {{.Image}}:{{.Tag}} {{.Cmd}}"`+"\n")

	isV = true
	tc = Image{}
	outb, _ = a.StubIO("", func() {
		tc = SearchImageFromYard("naktak/inst-test", "inst-test", "", 4)
	})
	outexpect = `5 [naktak/inst-test] config data were found in the yard.
Invalid command. Removed: te/st2
`
	a.Eq(tc.Tag, "0.0.1")
	a.Eq(tc.Script, `echo "hogeid4 {{.CfgDir}} {{.Image}}:{{.Tag}} {{.Cmd}}"`+"\n")
	a.Eq(outb, outexpect)

}

//NewYardPackFromYard
func TestNewYardPackFromYard(t *testing.T) {
	a.Set(t, "NewYardPackFromYard")
	orgExit := exit
	exit = pseudoExit
	flgExit = false
	envHome = "./misc/test/dcenv_home"
	envShell = "bash"
	tc := YardPack{}
	outexpect := ""

	tc = YardPack{}
	outb, _ := a.StubIO("", func() {
		tc = NewYardPackFromYard("naktak/inst-test-illigal")
	})
	outexpect = "yaml: unmarshal errors:\n  line 52: cannot unmarshal !!str `0.0.1` into main.Image\n"
	a.Eq(flgExit, true)
	a.Eq(outb, outexpect)

	exit = orgExit

	tc = YardPack{}
	tc = NewYardPackFromYard("naktak/inst-test-nofile")
	a.Eq(len(tc), 0)

	tc = YardPack{}
	tc = NewYardPackFromYard("naktak/inst-test")
	a.Eq(len(tc), 5)

}

//LoadConfig
func TestLoadConfig(t *testing.T) {
	a.Set(t, "LoadConfig")
	cfg := &Config{}
	outb := ""
	isV = true
	outb, _ = a.StubIO("", func() {
		cfg = LoadConfig("./misc/test/dcenv_home/files/dcenv_LoadConfig_fail")
	})
	a.Eq(cfg, (*Config)(nil))
	a.Eq(outb, `  File can not unmarshal.: ./misc/test/dcenv_home/files/dcenv_LoadConfig_fail
yaml: unmarshal errors:
  line 8: cannot unmarshal !!bool `+"`false`"+` into main.Image
`)

	outb, _ = a.StubIO("", func() {
		cfg = LoadConfig("./misc/test/dcenv_home/files/dcenv_LoadConfig_nofile")
	})
	a.Eq(cfg, (*Config)(nil))
	a.Eq(outb, `  File can not unmarshal.: ./misc/test/dcenv_home/files/dcenv_LoadConfig_nofile
open ./misc/test/dcenv_home/files/dcenv_LoadConfig_nofile: The system cannot find the file specified.
`)

	outb, _ = a.StubIO("", func() {
		cfg = LoadConfig("./misc/test/dcenv_home/files/dcenv_LoadConfig")
	})
	a.Neq(cfg, (*Config)(nil))
	a.Eq(outb, ``)
	isV = false

	outb, _ = a.StubIO("", func() {
		cfg = LoadConfig("./misc/test/dcenv_home/files/dcenv_LoadConfig_fail")
	})
	a.Eq(cfg, (*Config)(nil))
	a.Eq(outb, ``)

}
