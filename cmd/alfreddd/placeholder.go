package main

import (
	"fmt"
	"strings"
)

func makeProtoPlaceholder(service string) []byte {
	content := fmt.Sprintf(`
syntax = "proto3";

package %s.v1;

import "google/api/annotations.proto";

// TODO: Edit this proto to something more meaningful for your application.
service Service {
	rpc GetVersion(GetVersionRequest) returns (GetVersionResponse) {
		option (google.api.http) = {
      post: "/v1/hello"
			body: "*"
    };
	}
}

message GetVersionRequest {
	string name = 1;
}
message GetVersionResponse {
	string message = 1;
}
	`, service)

	return []byte(strings.Trim(content, "\n"))
}

func makeMainPlaceholder(module string) []byte {
	content := fmt.Sprintf(`
package main

import (
	"os"
	"os/signal"
	"syscall"
	
	log "github.com/sirupsen/logrus"
	service_interface "%s/internal/interface"
)

//nolint:all
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// TODO: Edit this file to something more meaningful for your application.
func main() {
	svc, err := service_interface.NewService()
	if err != nil {
		log.Fatal(err)
	}

	log.RegisterExitHandler(svc.Stop)
	
	log.Info("starting service...")
	if err := svc.Start(); err != nil {
		log.Fatal(err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)
	<-sigChan

	log.Info("shutting down service...")
	log.Exit(0)
}
	`, module)

	return []byte(strings.Trim(content, "\n"))
}

func makeInterfacePlaceholder(module string) []byte {
	content := fmt.Sprintf(`
package service_interface

import (
	grpc_interface "%s/internal/interface/grpc"
)

// TODO: Edit this file to something more meaningful for your application.
type Service interface {
	Start() error
	Stop()
}

func NewService() (Service, error) {
	return grpc_interface.NewService()
}
	`, module)

	return []byte(strings.Trim(content, "\n"))
}

func makeServicePlaceholder() []byte {
	content := `
package grpc_interface

import (
	log "github.com/sirupsen/logrus"
)

// TODO: Edit this file to something more meaningful for your application.
type service struct {}

func NewService() (*service, error) {
	return &service{}, nil
}

func (s *service) Start() error {
	log.Debug("service is listening")
	return nil
}

func (s *service) Stop() {
	log.Debug("service stopped")
}
	`

	return []byte(strings.Trim(content, "\n"))
}
