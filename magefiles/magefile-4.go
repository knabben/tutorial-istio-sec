package main

import (
	"github.com/knabben/tutorial-istio-sec/magefiles/pkg/apps"
	"github.com/knabben/tutorial-istio-sec/magefiles/pkg/istio"
	"github.com/knabben/tutorial-istio-sec/magefiles/pkg/kind"
	"github.com/magefile/mage/mg"
)

const (
	NAMESPACE4    = "default"
	CLUSTER4_NAME = "ambient"
	SPECS4_FOLDER = "4-istio-s2s/specs"
)

type SM4 mg.Namespace

var handleCM = false

// Install installs kind and metallb into the cluster
func (SM4) Install() error {
	return kind.InstallKind(CLUSTER4_NAME, SPECS4_FOLDER, false)
}

// Delete cleans up kind from cluster
func (SM4) Delete() error {
	return kind.DeleteKind(CLUSTER4_NAME)
}

// InstallIstio install ambient
func (SM4) InstallIstio() error {
	return istio.InstallIstio(SPECS4_FOLDER, NAMESPACE4, true, true)
}

// DeleteIstio cleans up resources from cluster
func (SM4) DeleteIstio() error {
	return istio.DeleteIstio(SPECS4_FOLDER, handleCM)
}

// Deploy creates the pre-defined application topology
func (SM4) Deploy() error {
	return apps.DeployApplication(SPECS4_FOLDER, NAMESPACE4, handleCM, false, "grpc")
}

// CheckApp check the application information
func (SM4) CheckApp() error {
	return nil
	//return checkers(SPECS4_FOLDER, NAMESPACE4, handleCM, false)
}
