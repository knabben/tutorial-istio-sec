package main

import (
	"github.com/knabben/tutorial-istio-sec/magefiles/pkg/kind"
	"github.com/knabben/tutorial-istio-sec/magefiles/pkg/spire"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

const (
	CLUSTER_NAME_2 = "grpc"
	SPEC_PATH      = "2-spiffe/specs/install.yaml"
	SPEC_APPS      = "2-spiffe/specs/app.yaml"
)

type SM2 mg.Namespace

// Install installs resources into the cluster
func (SM2) Install() error {
	return kind.InstallKind(CLUSTER_NAME_2, specsFolder, false)
}

// InstallSpire install SPIRE server and application
func (SM2) InstallSpire() error {
	return spire.InstallSpire(SPEC_PATH, SPEC_APPS)
}

// Delete cleans up resources from cluster
func (SM2) Delete() error {
	return kind.DeleteKind(CLUSTER_NAME_2)
}

// BuildClient builds the SPIFFE gRPC client and push to registry
func (SM2) BuildClient() error {
	if err := sh.RunV("docker", "build", "-f", "code/client/Dockerfile", "-t", "knabben/spiffe-client", "code/client"); err != nil {
		return err
	}
	if err := sh.RunV("docker", "push", "knabben/spiffe-client"); err != nil {
		return err
	}
	return nil
}

// BuildServer build the SPIFFE gRPC server and push to registry
func (SM2) BuildServer() error {
	if err := sh.RunV("docker", "build", "-f", "code/server/Dockerfile", "-t", "knabben/spiffe-server", "code/server"); err != nil {
		return err
	}
	if err := sh.RunV("docker", "push", "knabben/spiffe-server"); err != nil {
		return err
	}
	return nil
}
