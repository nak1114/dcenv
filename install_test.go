package main

import (
	"testing"

	a "github.com/nak1114/goutil/assert"
)

//NewYardPackFromYard
func TestNewYardPackFromYard(t *testing.T) {
	a.Set(t, "NewYardPackFromYard")
	outb := ""
	envHome = "./misc/test/dcenv_home"
	envShell = "bash"
	tc := YardPack{}
	outexpect := ""

	//Success load yard pack
	tc = YardPack{}
	tc = NewYardPackFromYard("naktak/inst-test")
	a.Eq(len(tc), 5)

	//Load error (illigal format)
	tc = YardPack{}
	PseudoExit(func() {
		outb, _ = a.StubIO("", func() {
			tc = NewYardPackFromYard("naktak/inst-test-illigal")
		})
	}, 1, 1)
	outexpect = "yaml: unmarshal errors:\n  line 52: cannot unmarshal !!str `0.0.1` into main.Image\n"
	a.Eq(outb, outexpect)

	//Create new obj
	tc = YardPack{}
	tc = NewYardPackFromYard("naktak/inst-test-nofile")
	a.Eq(len(tc), 0)

}
