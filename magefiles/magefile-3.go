package main

import (
	"github.com/knabben/tutorial-istio-sec/magefiles/pkg/istio"
	"github.com/knabben/tutorial-istio-sec/magefiles/pkg/kind"
	"github.com/magefile/mage/mg"
)

const (
	namespace    = "default"
	CLUSTER_NAME = "ambient"
	specsFolder  = "3-istio-gw/specs"
)

type SM3 mg.Namespace

// Install installs kind and metallb into the cluster
func (SM3) Install() error {
	return kind.InstallKind(CLUSTER_NAME, specsFolder, true)
}

// InstallIstio install ambient
func (SM3) InstallIstio() error {
	return istio.InstallIstio(specsFolder, namespace)
}

/*
// Deploy creates the pre-defined topology for tests
func (SM3) Deploy() error {
	if err := DeployApplication(namespace); err != nil {
		return err
	}

	return nil
}

// Policies create a VirtualService and define application Authorization files
func (SM3) Policies() error {
	if err := ApplyPolicies(namespace); err != nil {
		return err
	}

	return nil
}
*/

// Delete cleans up kind from cluster
func (SM3) Delete() error {
	return kind.DeleteKind(CLUSTER_NAME)
}

// DeleteIstio cleans up resources from cluster
func (SM3) DeleteIstio() error {
	return istio.DeleteIstio(specsFolder)
}
