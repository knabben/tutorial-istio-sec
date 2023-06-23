//go:build mage
// +build mage

package main

import (
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type K8S mg.Namespace

const (
	KIND_YAML = "./specs/kind.yaml"
)

// Install installs resources into the cluster
func (K8S) Install() error {
	if err := sh.RunV("kind", "delete", "cluster", "--name", "ambient"); err != nil {
		return err
	}
	if err := sh.RunV("kind", "create", "cluster", "--config", KIND_YAML); err != nil {
		return err
	}
	return nil
}

// Delete cleans up resources from cluster
func (K8S) Delete() error {
	if err := sh.RunV("kind", "delete", "cluster", "--name", "ambient"); err != nil {
		return err
	}
	return nil
}
