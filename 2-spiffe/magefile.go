//go:build mage
// +build mage

package main

import (
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type K8S mg.Namespace

const SPEC_PATH = "./specs/install.yaml"

func (K8S) Install() error {
	if err := sh.RunV("kubectl", "apply", "-f", SPEC_PATH); err != nil {
		return err
	}
	return nil
}

func (K8S) Delete() error {
	if err := sh.RunV("kubectl", "delete", "-f", SPEC_PATH); err != nil {
		return err
	}
	return nil
}
