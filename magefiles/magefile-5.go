package main

import (
	"github.com/knabben/tutorial-istio-sec/magefiles/pkg/apps"
	"github.com/knabben/tutorial-istio-sec/magefiles/pkg/helm"
	"github.com/knabben/tutorial-istio-sec/magefiles/pkg/kind"
	"github.com/knabben/tutorial-istio-sec/magefiles/pkg/sigstore"
	"github.com/magefile/mage/mg"
)

const (
	NAMESPACE5    = "default"
	CLUSTER5_NAME = "sigstore"
	SPECS5_FOLDER = "5-supply-chain/specs"
)

type SM5 mg.Namespace

// Install installs kind w/ metalLB in the cluster
func (SM5) Install() error {
	return kind.InstallKind(CLUSTER5_NAME, SPECS5_FOLDER, false)
}

// Delete cleans up kind from cluster
func (SM5) Delete() error {
	return kind.DeleteKind(CLUSTER5_NAME)
}

// InstallPC installs the policy controller in the cluster
func (SM5) InstallPC() error {
	return helm.InstallPC(NAMESPACE5)
}

// SignAndVerify sign a container with OIDC and verify it in the sequence
func (SM5) SignAndVerify() error {
	email := "aknabben@vmware.com"
	container := "ttl.sh/knabben/netshoot:3h"
	sigstore.Sign(container)
	return sigstore.Verify(container, email, "https://github.com/login/oauth")
}

// SyftSBOM generates a new SBOM based on container and attest it
func (SM5) SyftSBOM() error {
	container := "ttl.sh/knabben/netshoot:4h"
	sigstore.SBOM(container)
	return sigstore.Attest(container)
}

// Policy applies the policy in the cluster
func (SM5) Policy() error {
	return apps.ApplyPolicies(SPECS5_FOLDER, NAMESPACE5)
}
