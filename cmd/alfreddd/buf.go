package main

import (
	"fmt"
	"strings"
)

func makeBufYaml(project string) []byte {
	content := fmt.Sprintf(`
version: v1
name: buf.build/%s
deps:
  - buf.build/googleapis/googleapis
breaking:
	use:
		- FILE
	ignore_unstable_packages: true
lint:
	use:
		- DEFAULT
	`, project)
	content = strings.Replace(content, "\t", "  ", -1)
	return []byte(strings.Trim(content, "\n"))
}

func makeBufWorkYaml() []byte {
	content := `
version: v1
directories:
	- api-spec/protobuf	
	`
	content = strings.Replace(content, "\t", "  ", -1)
	return []byte(strings.Trim(content, "\n"))
}

func makeBufGenYaml(module string) []byte {
	content := fmt.Sprintf(`
version: v1
managed:
	enabled: true
	go_package_prefix:
		default: %s/api-spec/protobuf/gen
plugins:
	# Golang
	- remote: buf.build/protocolbuffers/plugins/go
		out: api-spec/protobuf/gen/go
		opt: paths=source_relative
	- remote: buf.build/grpc/plugins/go
		out: api-spec/protobuf/gen/go
		opt: paths=source_relative,require_unimplemented_servers=false
	- remote: buf.build/grpc-ecosystem/plugins/grpc-gateway
    out: api-spec/protobuf/gen
    opt: paths=source_relative
	#OpenApi
	- remote: buf.build/grpc-ecosystem/plugins/openapiv2
		out: api-spec/openapi/swagger
		`, module)

	content = strings.Replace(content, "\t", "  ", -1)

	return []byte(strings.Trim(content, "\n"))
}
