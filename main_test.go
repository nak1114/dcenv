package main

import (
	"os"
	"testing"
)

var flgExit = false

func pseudoExit(n int) {
	flgExit = true
}

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}
