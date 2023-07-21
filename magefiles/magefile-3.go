package main

import (
	"github.com/knabben/tutorial-istio-sec/magefiles/pkg/apps"
	"github.com/knabben/tutorial-istio-sec/magefiles/pkg/istio"
	"github.com/knabben/tutorial-istio-sec/magefiles/pkg/kind"
	"github.com/magefile/mage/mg"
)

const (
	NAMESPACE3    = "default"
	CLUSTER3_NAME = "ambient"
	SPECS3_FOLDER = "3-istio-gw/specs"
)

type SM3 mg.Namespace

// Install installs kind and metallb into the cluster
func (SM3) Install() error {
	return kind.InstallKind(CLUSTER3_NAME, SPECS3_FOLDER, true)
}

// InstallIstio install ambient
func (SM3) InstallIstio() error {
	return istio.InstallIstio(SPECS3_FOLDER, NAMESPACE3, true, true)
}

// Deploy creates the pre-defined application topology
func (SM3) Deploy() error {
	return apps.DeployApplication(SPECS3_FOLDER, NAMESPACE3, true, true)
}

// Policies create a VirtualService and define application Authorization files
func (SM3) Policies() error {
	return apps.ApplyPolicies(SPECS3_FOLDER, NAMESPACE3)
}

// Delete cleans up kind from cluster
func (SM3) Delete() error {
	return kind.DeleteKind(CLUSTER3_NAME)
}

// DeleteIstio cleans up resources from cluster
func (SM3) DeleteIstio() error {
	return istio.DeleteIstio(SPECS3_FOLDER, true)
}
