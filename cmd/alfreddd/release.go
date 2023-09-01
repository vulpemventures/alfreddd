package main

import (
	"fmt"
	"strings"
)

func makeGoreleaserYaml(project, app string) []byte {
	org := strings.Split(project, "/")[0]
	image := fmt.Sprintf("ghcr.io/%s/%s", org, app)
	content := fmt.Sprintf(`
builds:
  - id: "%s"
    main: ./cmd/%s
    ldflags:
      - -s -X 'main.version={{.Version}}' -X 'main.commit={{.Commit}}' -X 'main.date={{.Date}}'
    goos:
      - linux
      - darwin
    goarch:
      - amd64
      - arm64
    binary: %s

## flag the semver v**.**.**-<tag>.* as pre-release on Github
release:
  prerelease: auto

checksum:
  name_template: "checksums.txt"

snapshot:
  name_template: "{{ .Tag }}-next"

changelog:
  use: github-native

archives:
  - id: %s
    format: binary
    builds:
      - %s
    name_template: "%s-v{{ .Version }}-{{ .Os }}-{{ .Arch }}"

dockers:
  ###########################
  # tag latest & prerelease #
  ###########################

  #amd64
  - image_templates:
      - "%s:{{ .Tag }}-amd64"
        # push always either release or prerelease with a docker tag with the semver only
    skip_push: "false"
    use: buildx
    ids:
      - %s
    dockerfile: goreleaser.Dockerfile
    # GOOS of the built binaries/packages that should be used.
    goos: linux
    # GOARCH of the built binaries/packages that should be used.
    goarch: amd64
    # Template of the docker build flags.
    build_flag_templates:
      - "--platform=linux/amd64"
      - "--pull"
      - "--label=org.opencontainers.image.created={{.Date}}"
      - "--label=org.opencontainers.image.title=%s"
      - "--label=org.opencontainers.image.revision={{.FullCommit}}"
      - "--label=org.opencontainers.image.version={{.Version}}"
      - "--build-arg=VERSION={{.Version}}"
      - "--build-arg=COMMIT={{.Commit}}"
      - "--build-arg=DATE={{.Date}}"
  - image_templates:
      - "%s:{{ .Tag }}-arm64v8"
        # push always either release or prerelease with a docker tag with the semver only
    skip_push: "false"
    use: buildx
    ids:
      - %s
    dockerfile: goreleaser.Dockerfile
    # GOOS of the built binaries/packages that should be used.
    goos: linux
    # GOARCH of the built binaries/packages that should be used.
    goarch: arm64
    # Template of the docker build flags.
    build_flag_templates:
      - "--platform=linux/arm64/v8"
      - "--pull"
      - "--label=org.opencontainers.image.created={{.Date}}"
      - "--label=org.opencontainers.image.title=%s"
      - "--label=org.opencontainers.image.revision={{.FullCommit}}"
      - "--label=org.opencontainers.image.version={{.Version}}"
      - "--build-arg=VERSION={{.Version}}"
      - "--build-arg=COMMIT={{.Commit}}"
      - "--build-arg=DATE={{.Date}}"

docker_manifests:
  - name_template: %s:{{ .Tag }}
    image_templates:
      - %s:{{ .Tag }}-amd64
      - %s:{{ .Tag }}-arm64v8
    skip_push: "false"

  - name_template: %s:latest
    image_templates:
      - %s:{{ .Tag }}-amd64
      - %s:{{ .Tag }}-arm64v8
    skip_push: auto
	`, app, app, app, app, app, app, image, app, app, image, app, app, image, image, image, image, image, image)
	content = strings.Replace(content, "\t", "  ", -1)

	return []byte(strings.Trim(content, "\n"))
}
