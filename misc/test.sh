#!/bin/bash
go test --coverprofile=./misc/tmp/cover.out
go tool cover --html=./misc/tmp/cover.out -o ./misc/tmp/cover.html
