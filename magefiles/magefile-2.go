package main

import (
	"github.com/knabben/tutorial-istio-sec/magefiles/pkg/kind"
	"github.com/knabben/tutorial-istio-sec/magefiles/pkg/spire"
	"github.com/magefile/mage/mg"
)

const (
	CLUSTER_NAME_2 = "spire"
	SPEC2_PATH     = "2-spire/specs/"
)

type SM2 mg.Namespace

// Install installs resources into the cluster
func (SM2) Install() error {
	return kind.InstallKind(CLUSTER_NAME_2, SPEC2_PATH, false)
}

// Delete cleans up resources from cluster
func (SM2) Delete() error {
	return kind.DeleteKind(CLUSTER_NAME_2)
}

// InstallSpire install SPIRE server and application
func (SM2) InstallSpire() error {
	return spire.Bootstrap(SPEC2_PATH)
}

func (SM2) Deploy() error {
	if err := spire.Deploy(SPEC2_PATH); err != nil {
		return err
	}

	return spire.Check()
}
