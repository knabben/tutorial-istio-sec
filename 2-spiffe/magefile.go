//go:build mage
// +build mage

package main

import (
	"fmt"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"time"
)

type K8S mg.Namespace

const (
	SPEC_PATH = "./specs/install.yaml"
	SPEC_APPS = "./specs/app.yaml"
)

// Install installs resources into the cluster
func (K8S) Install() error {
	if err := sh.RunV("kind", "delete", "cluster"); err != nil {
		return err
	}
	if err := sh.RunV("kind", "create", "cluster"); err != nil {
		return err
	}

	if err := sh.RunV("kubectl", "apply", "-f", SPEC_PATH); err != nil {
		return err
	}

	for {
		// wait for server run, if pod does not exist keep trying.
		args := []string{"wait", "--for=condition=Ready", "-n", "spire", "pod", "-l", "app=spire-server", "--timeout", "300s"}
		if err := sh.RunV("kubectl", args...); err == nil {
			break
		}
		time.Sleep(time.Second)
	}

	// create the SPIRE node attestation
	if err := sh.RunV("kubectl",
		"exec", "-n", "spire",
		"spire-server-0", "--",
		"/opt/spire/bin/spire-server", "entry", "create",
		"-spiffeID", "spiffe://opssec.in/ns/spire/sa/spire-agent",
		"-selector", "k8s_sat:cluster:kind",
		"-selector", "k8s_sat:agent_ns:spire",
		"-selector", "k8s_sat:agent_sa:spire-agent",
		"-node"); err != nil {
		return err
	}

	// create client and server wl attestation
	for _, n := range []string{"client", "server"} {
		if err := sh.RunV("kubectl", "create", "serviceaccount", n); err != nil {
			return err
		}
		err := sh.RunV("kubectl",
			"exec", "-n", "spire",
			"spire-server-0", "--",
			"/opt/spire/bin/spire-server", "entry", "create",
			"-spiffeID", fmt.Sprintf("spiffe://opssec.in/ns/default/sa/%s", n),
			"-parentID", "spiffe://opssec.in/ns/spire/sa/spire-agent",
			"-selector", "k8s:ns:default",
			"-selector", "k8s:sa:"+n)
		if err != nil {
			return err
		}
	}

	// Install client/server deployment
	if err := sh.RunV("kubectl", "apply", "-f", SPEC_APPS); err != nil {
		return err
	}

	return nil
}

// Delete cleans up resources from cluster
func (K8S) Delete() error {
	if err := sh.RunV("kubectl", "delete", "-f", SPEC_PATH); err != nil {
		return err
	}

	if err := sh.RunV("kubectl", "apply", "-f", SPEC_APPS); err != nil {
		return err
	}

	if err := sh.RunV("kubectl", "delete", "serviceaccount", "client", "server"); err != nil {
		return err
	}
	return nil
}

// ContainerClient build the SPIFFE gRPC client and push to registry
func (K8S) ContainerClient() error {
	if err := sh.RunV("docker", "build", "-f", "code/client/Dockerfile", "-t", "knabben/spiffe-client", "code/client"); err != nil {
		return err
	}
	if err := sh.RunV("docker", "push", "knabben/spiffe-client"); err != nil {
		return err
	}
	return nil
}

// ContainerServer build the SPIFFE gRPC server and push to registry
func (K8S) ContainerServer() error {
	if err := sh.RunV("docker", "build", "-f", "code/server/Dockerfile", "-t", "knabben/spiffe-server", "code/server"); err != nil {
		return err
	}
	if err := sh.RunV("docker", "push", "knabben/spiffe-server"); err != nil {
		return err
	}
	return nil
}
