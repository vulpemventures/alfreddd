package main

import "fmt"

func makeScriptBuild(app string) []byte {
	content := fmt.Sprintf(`
#!/bin/bash

set -e

PARENT_PATH=$(dirname $(
	cd $(dirname $0)
	pwd -P
))

OS=$(eval "go env GOOS")
ARCH=$(eval "go env GOARCH")

pushd $PARENT_PATH
mkdir -p build
GO111MODULE=on go build -ldflags="-s -w" -o build/%s-$OS-$ARCH cmd/%s/main.go
popd	
  `, app, app)

	return []byte(content)
}
