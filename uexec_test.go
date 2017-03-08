package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	a "github.com/nak1114/goutil/assert"
	"github.com/nak1114/goutil/cp"
)

func TestShowExecFile(t *testing.T) {
	envHome = "./misc/tmp/item"
	envShell = "bash"
	envCommand = "cmd.bat"
	cmd := "cmd"
	fname := "./misc/test/uexec/MakeExecFile_cfg.txt"
	envShell = "bash"
	if runtime.GOOS == `windows` {
		envShell = "windows"
	}
	buf, _ := ioutil.ReadFile("./misc/test/uexec/MakeExecFile_" + envShell + ".txt")

	a.Set(t, "ShowExecFile")
	os.RemoveAll(envHome)
	cp.Dir("./misc/item", envHome)

	m := NewConfig(fname)
	bufr, _ := a.StubIO("", func() {
		ShowExecFile(&m, cmd, fname)
	})
	a.Eq(bufr, string(buf))

	os.RemoveAll(envHome)
}

func TestMakeExecFile(t *testing.T) {
	envHome = "./misc/tmp/item"
	envShell = "bash"
	envCommand = "cmd.bat"
	cmd := "cmd"
	fname := "./misc/test/uexec/MakeExecFile_cfg.txt"
	envShell = "bash"
	if runtime.GOOS == `windows` {
		envShell = "windows"
	}
	ofile := filepath.Join(envHome, "tmp", envCommand)
	buf, _ := ioutil.ReadFile("./misc/test/uexec/MakeExecFile_" + envShell + ".txt")

	a.Set(t, "MakeExecFile")
	os.RemoveAll(envHome)
	cp.Dir("./misc/item", envHome)

	m := NewConfig(fname)
	MakeExecFile(&m, cmd, fname)
	a.FileContent(ofile, string(buf))

	os.RemoveAll(envHome)
}

func TestMakeShimsFile(t *testing.T) {
	envHome = "./misc/tmp/item"
	envShell = "bash"
	envExt = ""
	cmd := "hoge"
	ofile := filepath.Join(envHome, "shims", cmd+envExt)
	buf, _ := ioutil.ReadFile("./misc/test/uexec/MakeShimsFile.txt")

	a.Set(t, "MakeShimFile")
	os.RemoveAll(envHome)
	cp.Dir("./misc/item", envHome)

	MakeShimsFile(cmd)
	a.FileContent(ofile, string(buf))

	os.RemoveAll(envHome)

}

func TestMakeExecFileSystem(t *testing.T) {
	envHome = "./misc/tmp/item"
	envCommand = "fuga.bat"
	cmd := "fuga"
	ofile := filepath.Join(envHome, "tmp", envCommand)
	plist := []string{
		"/bin",
		filepath.Join(envHome, "shims"),
		"/sbin",
	}
	os.Setenv("PATH", strings.Join(plist, string(os.PathListSeparator)))
	envShell = "bash"
	if runtime.GOOS == `windows` {
		envShell = "windows"
	}
	buf, _ := ioutil.ReadFile("./misc/test/uexec/MakeExecFileSystem_" + envShell + ".txt")

	a.Set(t, "MakeExecFileSystem")
	os.RemoveAll(envHome)
	cp.Dir("./misc/item", envHome)

	MakeExecFileSystem(cmd)
	a.FileContent(ofile, string(buf))

	os.RemoveAll(envHome)

}

func TestMakeArgsFile(t *testing.T) {
	envHome = "./misc/tmp/item"
	ofile := filepath.Join(envHome, "files", "__args__")

	a.Set(t, "MakeArgsFile")
	os.RemoveAll(envHome)
	cp.Dir("./misc/item", envHome)

	os.Args = []string{"dcenv", "install", "--list", "hoge fuga"}
	MakeArgsFile(2)
	a.FileContent(ofile, ` --list "hoge fuga"`)

	os.RemoveAll(envHome)
}
