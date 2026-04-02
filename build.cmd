@echo off

go test -v -count=1 ./...

golangci-lint run

revive --formatter stylish ./...
