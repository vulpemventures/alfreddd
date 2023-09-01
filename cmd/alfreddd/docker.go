package main

import (
	"fmt"
	"strings"
)

func makeDockerfile(ghProject, app string) []byte {
	// retrieve gh profile name and eventually use the first part if it's in the form `user-name`
	user := strings.Split(ghProject, "/")[0]
	user = strings.Split(user, "-")[0]

	content := strings.Trim(
		fmt.Sprintf(`
# first image used to build the sources
FROM golang:1.21 AS builder

ARG VERSION
ARG COMMIT
ARG DATE
ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

COPY . .

RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags="-X 'main.Version=${COMMIT}' -X 'main.Commit=${COMMIT}' -X 'main.Date=${COMMIT}'" -o bin/%s cmd/%s/main.go

# Second image, running the %s executable
FROM debian:buster-slim

# $USER name, and data $DIR to be used in the 'final' image
ARG USER=%s
ARG DIR=/home/%s

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates

COPY --from=builder /app/bin/* /usr/local/bin/

# NOTE: Default GID == UID == 1000
RUN adduser --disabled-password \
						--home "$DIR/" \
						--gecos "" \
						"$USER"
USER $USER

# Prevents 'VOLUME $DIR/.%s/' being created as owned by 'root'
RUN mkdir -p "$DIR/.%s/"

# Expose volume containing all '%s' data
VOLUME $DIR/.%s/

ENTRYPOINT [ "%s" ]
    `, app, app, app, user, user, app, app, app, app, app), "\n",
	)

	return []byte(content)
}

func makeDockerIngnore() []byte {
	return []byte(strings.Trim(`
*.md
build
dist
.git
.github
Dockerfile	
  `, "\n"))
}

func makeReleaseDockerfile(ghProject, app string) []byte {
	// retrieve gh profile name and eventually use the first part if it's in the form `user-name`
	user := strings.Split(ghProject, "/")[0]
	user = strings.Split(user, "-")[0]

	content := fmt.Sprintf(`
FROM debian:buster-slim

ARG TARGETPLATFORM

WORKDIR /app

COPY . .

RUN set -ex \
  && if [ "${TARGETPLATFORM}" = "linux/amd64" ]; then export TARGETPLATFORM=amd64; fi \
  && if [ "${TARGETPLATFORM}" = "linux/arm64" ]; then export TARGETPLATFORM=arm64; fi \
  && mv "%s-linux-$TARGETPLATFORM" /usr/local/bin/%s


# $USER name, and data $DIR to be used in the 'final' image
ARG USER=%s
ARG DIR=/home/%s

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates

# NOTE: Default GID == UID == 1000
RUN adduser --disabled-password \
            --home "$DIR/" \
            --gecos "" \
            "$USER"
USER $USER

# Prevents 'VOLUME $DIR/.%s/' being created as owned by 'root'
RUN mkdir -p "$DIR/.%s/"

# Expose volume containing all %s data
VOLUME $DIR/.%s/

ENTRYPOINT [ "%s" ]
	`, app, app, user, user, app, app, app, app, app)

	return []byte(strings.Trim(content, "\n"))
}
