package main

import (
	"os"
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

// Load/Save/Delete yaml
func TestYaml(t *testing.T) {
	a.Set(t, "TestYaml")
	ret := make(map[string]string)
	tmp := map[string]string{
		"hoge": "fuga",
		"foo":  "bar",
	}
	p := TestP(0)
	// Success save yaml
	e := SaveYaml(tmp, "./misc/tmp/test.yml")
	a.Eq(e, error(nil))
	// Success load yaml
	e = LoadYaml(&ret, "./misc/tmp/test.yml")
	a.Eq(e, error(nil))
	// Success delete yaml
	e = DeleteYaml("./misc/tmp/test.yml")
	a.Eq(e, error(nil))
	a.Eq(ret["hoge"], tmp["hoge"])

	// Error load yaml
	e = LoadYaml(&ret, "./misc/tmp/test.yml")
	a.Neq(e, error(nil))
	// Error save yaml
	e = SaveYaml(p, "./misc/tmp/test.yml")
	a.Neq(e, error(nil))

}

func TestSearchImageFromYard(t *testing.T) {
	a.Set(t, "SearchImageFromYard")
	outb := ""
	envHome = "./misc/test/dcenv_home"
	envShell = "bash"
	tc := Image{}

	// Error : file not found
	PseudoExit(func() {
		outb, _ = a.StubIO("", func() {
			tc = SearchImageFromYard("naktak/inst-test-null", "inst-test-null", "", 0)
		})
	}, 1, 1)
	outexpect := `[naktak/inst-test-null] config data were not found in the yard.
Please command 'dcenv yard --list naktak/inst-test-null'
`
	a.Eq(tc.Tag, "")
	a.Eq(outb, outexpect)

	// Error : image not found
	tc = Image{}
	PseudoExit(func() {
		outb, _ = a.StubIO("", func() {
			tc = SearchImageFromYard("naktak/inst-test", "inst-test-null", "", 20)
		})
	}, 1, 1)
	outexpect = `5 [naktak/inst-test] config data were found in the yard.
Too large number you choice.:(20)
Please command 'dcenv yard -d naktak/inst-test'
`
	a.Eq(tc.Tag, "")
	a.Eq(outb, outexpect)

	// Error : shell not found
	tc = Image{}
	envShell = "csh"
	PseudoExit(func() {
		outb, _ = a.StubIO("", func() {
			tc = SearchImageFromYard("naktak/inst-test", "inst-test", "1.1.1", 1)
		})
	}, 1, 1)
	a.Eq(tc.Tag, "")
	a.Eq(outb, `Not found command for this shell.:csh
Please command 'dcenv yard -d naktak/inst-test'
`)
	envShell = "bash"

	// Success : image found.
	// slip image set
	// add default tag
	tc = Image{}
	tc = SearchImageFromYard("naktak/inst-test", "inst-test", "", 0)
	a.Eq(tc.Tag, "latest")
	a.Eq(tc.Script, `echo "hogeid2 {{.CfgDir}} {{.Image}}:{{.Tag}} {{.Cmd}}"`+"\n")

	// Success : image found.
	// add a tag by user
	tc = Image{}
	tc = SearchImageFromYard("naktak/inst-test", "inst-test", "1.1.1", 3)
	a.Eq(tc.Tag, "1.1.1")
	a.Eq(tc.Script, `echo "hogeid3 {{.CfgDir}} {{.Image}}:{{.Tag}} {{.Cmd}}"`+"\n")

	// Success : image found.
	// Verbose modes
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

//LoadConfig
func TestLoadConfig(t *testing.T) {
	a.Set(t, "LoadConfig")
	cfg := &Config{}
	outb := ""

	// Verbose
	isV = true
	// success Load config
	outb, _ = a.StubIO("", func() {
		cfg = LoadConfig("./misc/test/inst/dcenv_LoadConfig.yml")
	})
	a.Neq(cfg, (*Config)(nil))
	a.Eq(outb, ``)

	// Load error(Yaml format error)
	outb, _ = a.StubIO("", func() {
		cfg = LoadConfig("./misc/test/inst/dcenv_LoadConfig_fail.yml")
	})
	a.Eq(cfg, (*Config)(nil))
	a.Eq(outb, `  File can not unmarshal.: ./misc/test/inst/dcenv_LoadConfig_fail.yml
yaml: unmarshal errors:
  line 8: cannot unmarshal !!bool `+"`false`"+` into main.Image
`)

	// Load error(file not found)
	outb, _ = a.StubIO("", func() {
		cfg = LoadConfig("./misc/test/dcenv_home/files/dcenv_LoadConfig_nofile")
	})
	a.Eq(cfg, (*Config)(nil))
	a.EqRegexp(outb, `^  File can not unmarshal\.\: \./misc/test/dcenv_home/files/dcenv_LoadConfig_nofile
open \./misc/test/dcenv_home/files/dcenv_LoadConfig_nofile\:`)

	// no Verbose
	isV = false
	// success Load config
	outb, _ = a.StubIO("", func() {
		cfg = LoadConfig("./misc/test/inst/dcenv_LoadConfig.yml")
	})
	a.Neq(cfg, (*Config)(nil))
	a.Eq(outb, ``)

	// Load error(Yaml format error)
	outb, _ = a.StubIO("", func() {
		cfg = LoadConfig("./misc/test/inst/dcenv_LoadConfig_fail.yml")
	})
	a.Eq(cfg, (*Config)(nil))
	a.Eq(outb, ``)

	// Load error(file not found)
	outb, _ = a.StubIO("", func() {
		cfg = LoadConfig("./misc/test/dcenv_home/files/dcenv_LoadConfig_nofile")
	})
	a.Eq(cfg, (*Config)(nil))
	a.Eq(outb, ``)

}

// NewConfig
func TestNewConfig(t *testing.T) {
	a.Set(t, "NewConfig")
	outb := ""
	cfg := Config{}

	// Verbose
	isV = true
	// Config load from file
	outb, _ = a.StubIO("", func() {
		cfg = NewConfig("./misc/test/inst/dcenv_NewConfig.yml")
	})
	a.Eq(outb, `Found the config file.: ./misc/test/inst/dcenv_NewConfig.yml
Config file.: ./misc/test/inst/dcenv_NewConfig.yml
`)

	// Config load error
	PseudoExit(func() {
		outb, _ = a.StubIO("", func() {
			cfg = NewConfig("./misc/test/inst/dcenv_NewConfig_fail.yml")
		})
	}, 1, 1)
	a.Eq(outb, `Found the config file.: ./misc/test/inst/dcenv_NewConfig_fail.yml
yaml: unmarshal errors:
  line 8: cannot unmarshal !!bool `+"`false`"+` into main.Image
`)

	// Config create new obj
	outb, _ = a.StubIO("", func() {
		cfg = NewConfig("./misc/test/inst/dcenv_NewConfig_nofile.yml")
	})
	a.Eq(outb, `Not found the config file.Create it.: ./misc/test/inst/dcenv_NewConfig_nofile.yml
Config file.: ./misc/test/inst/dcenv_NewConfig_nofile.yml
`)

	// no Verbose
	isV = false
	// Config load from file
	outb, _ = a.StubIO("", func() {
		cfg = NewConfig("./misc/test/inst/dcenv_NewConfig.yml")
	})
	a.Eq(outb, `Config file.: ./misc/test/inst/dcenv_NewConfig.yml
`)

	// Config load error
	PseudoExit(func() {
		outb, _ = a.StubIO("", func() {
			cfg = NewConfig("./misc/test/inst/dcenv_NewConfig_fail.yml")
		})
	}, 1, 1)
	a.Eq(outb, `yaml: unmarshal errors:
  line 8: cannot unmarshal !!bool `+"`false`"+` into main.Image
`)

	// Config create new obj
	outb, _ = a.StubIO("", func() {
		cfg = NewConfig("./misc/test/inst/dcenv_NewConfig_nofile.yml")
	})
	a.Eq(outb, `Config file.: ./misc/test/inst/dcenv_NewConfig_nofile.yml
`)
}

func TestAddImage(t *testing.T) {
	a.Set(t, "AddImage")
	envHome = "./misc/test/dcenv_home"
	img := SearchImageFromYard("naktak/config-test", "config-test", "", 0)
	cfg := Config{}
	outb := ""

	//Forced
	a.Sub(1, "Forced")
	isF := true
	//Success image name
	a.Sub(2, "Success image name")
	cfg = NewConfig("./misc/test/inst/dcenv_ConfigAddImage.yml")
	cfg.AddImage(img, "naktak/config-test2", isF)
	a.Eq(&cfg, LoadConfig("./misc/test/inst/dcenv_ConfigAddImage_exp_success.yml"))

	//Conflict image name
	a.Sub(2, "Conflict image name")
	cfg = NewConfig("./misc/test/inst/dcenv_ConfigAddImage.yml")
	cfg.AddImage(img, "naktak/config-test", isF)
	a.Eq(&cfg, LoadConfig("./misc/test/inst/dcenv_ConfigAddImage_exp_conflict.yml"))

	//Unforced
	a.Sub(1, "Unforced")
	isF = false
	//Success image name
	a.Sub(2, "Success image name")
	cfg = NewConfig("./misc/test/inst/dcenv_ConfigAddImage.yml")
	PseudoExit(func() {
		outb, _ = a.StubIO("", func() {
			cfg.AddImage(img, "naktak/config-test2", isF)
		})
	}, 0, 0)
	a.Eq(outb, ``)
	a.Eq(&cfg, LoadConfig("./misc/test/inst/dcenv_ConfigAddImage_exp_success.yml"))

	//Conflict image name(Y)
	a.Sub(2, "Conflict image name(Y)")
	cfg = NewConfig("./misc/test/inst/dcenv_ConfigAddImage.yml")
	PseudoExit(func() {
		PseudoPrompterYN(true, func() {
			outb, _ = a.StubIO("Y\n", func() {
				cfg.AddImage(img, "naktak/config-test", isF)
			})
		})
	}, 0, 0)
	a.EqRegexp(outb, `config add script`)
	a.Eq(&cfg, LoadConfig("./misc/test/inst/dcenv_ConfigAddImage_exp_conflict.yml"))

	//Conflict image name(N)
	a.Sub(2, "Conflict image name(N)")
	cfg = NewConfig("./misc/test/inst/dcenv_ConfigAddImage.yml")
	PseudoExit(func() {
		PseudoPrompterYN(false, func() {
			outb, _ = a.StubIO("", func() {
				cfg.AddImage(img, "naktak/config-test", isF)
			})
		})
	}, 1, 1)
	a.EqRegexp(outb, `config add script`)
	a.Eq(&cfg, LoadConfig("./misc/test/inst/dcenv_ConfigAddImage.yml"))
}

func TestDelImage(t *testing.T) {
	a.Set(t, "DelImage")
	cfg := Config{}

	isF := true
	//Success image name
	a.Sub(1, "Success image name")
	cfg = NewConfig("./misc/test/inst/dcenv_ConfigDelImage.yml")
	cfg.DelImage("golang", isF)
	a.Eq(&cfg, LoadConfig("./misc/test/inst/dcenv_ConfigDelImage_success.yml"))
	//image not found
	a.Sub(1, "image not found")
	cfg = NewConfig("./misc/test/inst/dcenv_ConfigDelImage.yml")
	PseudoExit(func() {
		cfg.DelImage("golanger", isF)
	}, 1, 1)
	a.Eq(&cfg, LoadConfig("./misc/test/inst/dcenv_ConfigDelImage.yml"))
}

func TestConfigWriteToFile(t *testing.T) {
	a.Set(t, "ConfigWriteToFile")
	isV = false
	cfg := Config{}
	//Success
	cfg = NewConfig("./misc/test/inst/dcenv_ConfigWriteToFile.yml")
	cfg.WriteToFile("./misc/tmp/dcenv_ConfigWriteToFile_success")
	a.Eq(
		LoadConfig("./misc/tmp/dcenv_ConfigWriteToFile_success"),
		LoadConfig("./misc/test/inst/dcenv_ConfigWriteToFile_exp_success.yml"))

	//Delete yaml
	cfg = NewConfig("./misc/test/inst/dcenv_ConfigWriteToFile_noimage.yml")
	cfg.WriteToFile("./misc/tmp/dcenv_ConfigWriteToFile_success")
	a.Eq(
		isExistFile("./misc/tmp/dcenv_ConfigWriteToFile_success"),
		false)

	//Fail Delete yaml
	cfg = NewConfig("./misc/test/inst/dcenv_ConfigWriteToFile_noimage.yml")
	PseudoExit(func() {
		cfg.WriteToFile("./misc/tmp/nodir/dcenv_ConfigWriteToFile")
	}, 1, 1)

	//Command conflict
	cfg = NewConfig("./misc/test/inst/dcenv_ConfigWriteToFile_conflict.yml")
	PseudoExit(func() {
		cfg.WriteToFile("./misc/tmp/dcenv_ConfigWriteToFile_success")
	}, 1, 1)
	a.Eq(
		isExistFile("./misc/tmp/dcenv_ConfigWriteToFile_success"),
		false)

	//Fail save yaml
	cfg = NewConfig("./misc/test/inst/dcenv_ConfigWriteToFile.yml")
	PseudoExit(func() {
		cfg.WriteToFile("./misc/tmp/nodir/dcenv_ConfigWriteToFile_success")
	}, 1, 1)

	//Verbose mode(success)
	isV = true
	cfg = NewConfig("./misc/test/inst/dcenv_ConfigWriteToFile.yml")
	cfg.WriteToFile("./misc/tmp/dcenv_ConfigWriteToFile_success")
	a.Eq(
		LoadConfig("./misc/tmp/dcenv_ConfigWriteToFile_success"),
		LoadConfig("./misc/test/inst/dcenv_ConfigWriteToFile_exp_success.yml"))

	//Verbose mode(delete)
	isV = true
	cfg = NewConfig("./misc/test/inst/dcenv_ConfigWriteToFile_noimage.yml")
	cfg.WriteToFile("./misc/tmp/dcenv_ConfigWriteToFile_success")
	a.Eq(
		isExistFile("./misc/tmp/dcenv_ConfigWriteToFile_success"),
		false)
}

func isExistFile(fname string) bool {
	_, err := os.Stat(fname)
	return err == nil
}
