package main

import (
	"fmt"
	"os"
	"testing"
	//"os/exec"
	//"runtime"
)

var message = "initial"
var count = 0

func set(mes string) {
	message = mes
	count = 0
}
func eq(t *testing.T, actual interface{}, expected interface{}) {
	count += 1
	if actual != expected {
		fmt.Println(message, count)
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}

func TestMain(t *testing.T) {
	//exec.Command("ls", "-la").Run()
	os.Setenv("DCENV_HOME", "")
	os.Args = []string{"dcenv", "unknown"}
	main()
}
