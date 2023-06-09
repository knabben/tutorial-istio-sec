package tls

import (
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type GRPC mg.Namespace

// RunServer runs a TLS based gRPC server
func (GRPC) RunServer() {
	sh.Run("go", "run", "./grpc/server")
}

// RunClient runs a TLS based gRPC client
func (GRPC) RunClient() {
	sh.Run("go", "run", "./grpc/client")
}

// CompileProto generates the go protobuf code, execute with [-w grpc/proto]
func (GRPC) CompileProto() error {
	err := sh.Run("go", "install", "google.golang.org/protobuf/cmd/protoc-gen-go@latest")
	if err != nil {
		return err
	}

	err = sh.Run(
		"protoc", "--go_out", ".", "--go_opt", "paths=source_relative",
		"--go-grpc_out", ".", "--go-grpc_opt", "paths=source_relative", "bank.proto",
	)
	if err != nil {
		return err
	}

	return nil
}
