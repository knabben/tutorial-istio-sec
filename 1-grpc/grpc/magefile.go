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
