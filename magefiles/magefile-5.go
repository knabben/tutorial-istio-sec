package main

import (
	"github.com/knabben/tutorial-istio-sec/magefiles/pkg/helm"
	"github.com/knabben/tutorial-istio-sec/magefiles/pkg/kind"
	"github.com/magefile/mage/mg"
)

const (
	NAMESPACE5    = "default"
	CLUSTER5_NAME = "sigstore"
	SPECS5_FOLDER = "5-supply-chain/specs"
)

type SM5 mg.Namespace

// Install installs kind and metallb into the cluster
func (SM5) Install() error {
	return kind.InstallKind(CLUSTER4_NAME, SPECS4_FOLDER, false)
}

// Delete cleans up kind from cluster
func (SM5) Delete() error {
	return kind.DeleteKind(CLUSTER4_NAME)
}

// InstallPC
func (SM5) InstallPC() error {
	return helm.InstallPC(CLUSTER5_NAME)
}

// SyftSBOM
func (SM5) SyftSBOM() error {
	return nil
}

// AttestSBOM
func (SM5) AttestSBOM() error {
	return nil
}

// CreatePolicy
func (SM5) CreatePolicy() error {
	return nil
}
