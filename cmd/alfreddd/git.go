package main

import (
	"fmt"
	"strings"
)

func makeGitIngnore() []byte {
	return []byte(strings.Trim(`
# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary, built with 'go test -c'
*.test

# Output of the go coverage tool, specifically when used with LiteIDE
*.out

/build/
/dist/

DS_Store
._.DS_Store
**/.DS_Store
**/._.DS_Store
  `, "\n"))
}

func makeLicense(project string) []byte {
	org := strings.Split(project, "/")[0]

	content := fmt.Sprintf(`
MIT License

Copyright (c) 2023 %s

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.	
	`, org)

	return []byte(strings.Trim(content, "\n"))
}

func makeReadme(project string) []byte {
	projectName := strings.Split(project, "/")[1]
	return []byte(fmt.Sprintf(`# %s`, projectName))
}

func makeUnitTestActionYaml() []byte {
	content := `
name: ci_unit

on:
  push:
    branches: [master]
  pull_request:
    branches: [master]

jobs:
  test:
    name: unit tests
    runs-on: ubuntu-latest
    steps:
      - uses: actions/setup-go@v3
        with:
          go-version: ">1.17.2"
      - uses: actions/checkout@v3
      - name: check linting
        uses: golangci/golangci-lint-action@v3
        with:
          version: v1.54
      - name: check code integrity
        uses: securego/gosec@master
        with:
          args: '-severity high -quiet ./...'
      - uses: bufbuild/buf-setup-action@v1.3.1
      - name: check proto linting
        run: make proto-lint
      - run: go get -v -t -d ./...
      - name: unit testing
        run: make test
	`

	content = strings.Replace(content, "\t", "  ", -1)
	return []byte(strings.Trim(content, "\n"))
}

func makeIntegrationTestActionYaml() []byte {
	content := `
name: ci_integration

on:
  push:
    branches: [master]

jobs:
  test:
    name: integration tests
    runs-on: ubuntu-latest
    steps:
      - uses: actions/setup-go@v3
        with:
          go-version: ">1.17.2"
      - uses: actions/checkout@v3
      - run: go get -v -t -d ./...
      - name: integration testing
        run: make integrationtest
	`

	content = strings.Replace(content, "\t", "  ", -1)
	return []byte(strings.Trim(content, "\n"))
}

func makeReleaseActionYaml() []byte {
	content := `
name: release

on:
  workflow_dispatch:
  push:
    tags:
      - "*"

jobs:
  goreleaser:
    runs-on: ubuntu-20.04
    env:
      DOCKER_CLI_EXPERIMENTAL: "enabled"

    steps:
      - uses: actions/checkout@v3
        with:
          fetch-depth: 0
      - uses: actions/setup-go@v3
        with:
          go-version: ">1.17.2"
      - uses: actions/cache@v1
        with:
          path: ~/go/pkg/mod
          key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
          restore-keys: |
            ${{ runner.os }}-go-
      - uses: docker/login-action@v1
        with:
          registry: ghcr.io
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}
      - name: release artifacts
        uses: goreleaser/goreleaser-action@v2
        with:
          version: latest
          args: release --rm-dist --debug
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
          GPG_FINGERPRINT: ${{ steps.import_gpg.outputs.fingerprint }}
      - uses: bufbuild/buf-setup-action@v1.3.1
			- name: release protos
        uses: bufbuild/buf-push-action@v1
        with:
          input: api-spec/protobuf
          buf_token: ${{ secrets.BUF_TOKEN }}
	`

	content = strings.Replace(content, "\t", "  ", -1)
	return []byte(strings.Trim(content, "\n"))
}
