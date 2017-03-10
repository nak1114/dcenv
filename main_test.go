package main

import (
	"os"
	"testing"

	a "github.com/nak1114/goutil/assert"
)

func PseudoExit(fn func(), expectedCode int, expectedCount int) (code int, count int) {
	code = 0
	count = 0
	orgObj := exit
	exit = func(n int) {
		if count == 0 {
			code = n
		}
		count++
	}
	fn()
	exit = orgObj
	a.Eq(code, expectedCode)
	a.Eq(count, expectedCount)
	return
}
func PseudoPrompterYN(pb bool, fn func()) {
	orgObj := prompterYN
	prompterYN = func(n string, b bool) bool {
		return pb
	}
	fn()
	prompterYN = orgObj
	return
}

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}
