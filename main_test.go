package main

import (
	"fmt"
	"os"
	"testing"
)

var message = "initial"
var count = 0

func set(mes string) {
	message = mes
	count = 0
}
func eq(t *testing.T, actual interface{}, expected interface{}) {
	count++
	if actual != expected {
		fmt.Println(message, count)
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}
