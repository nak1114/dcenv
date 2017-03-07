package main

import (
	"runtime"
	"testing"

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
