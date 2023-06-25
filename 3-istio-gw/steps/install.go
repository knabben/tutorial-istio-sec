package steps

import (
	"github.com/magefile/mage/sh"
	"path"
)

const (
	METALLB_URL = "https://raw.githubusercontent.com/metallb/metallb/v0.13.10/config/manifests/metallb-native.yaml"
	METALLB_CR  = "metallb_cr.yaml"
)

// InstallKind installs kind as a base K8s cluster for local development
func (s *ServiceMesh) InstallKind(name string, config string, withLB bool) error {
	// Delete old cluster if exists
	if err := kind("delete", "cluster", "--name", name); err != nil {
		return err
	}

	// Create a new cluster using predefined configuration file
	configPath := path.Join(SpecsFolder, config)
	if err := kind("create", "cluster", "--name", name, "--config", configPath); err != nil {
		return err
	}

	// Set sysctl parameters on each node
	for _, host := range []string{"ambient-worker2", "ambient-control-plane", "ambient-worker"} {
		args := []string{"exec", host, "sysctl", "-w"}
		for _, s := range []string{"fs.inotify.max_user_instances=1024", "fs.inotify.max_user_watches=1048576"} {
			if err := sh.RunV("docker", append(args, s)...); err != nil {
				return nil
			}
		}
	}

	if withLB {
		// Install MetalLB in the cluster
		for _, spec := range []string{METALLB_URL, path.Join(SpecsFolder, METALLB_CR)} {
			_ = kubectl("-n", "metallb-system", "wait", "--for=condition=Ready", "pod", "-l", "app=metallb", "--timeout", "300s")
			// Create a new cluster using predefined configuration file
			if err := kubectl("apply", "-f", spec); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *ServiceMesh) InstallIstio(p, namespace string) error {
	// Apply the Gateway API custom resources.
	args := [][]string{
		{"kustomize", "github.com/kubernetes-sigs/gateway-api/config/crd?ref=v0.6.2", "-o", "/tmp/kustomized"},
		{"apply", "-f", "/tmp/kustomized"},
	}
	for _, a := range args {
		if err := kubectl(a...); err != nil {
			return err
		}
	}

	// Install Istio with custom ambient
	argss := []string{"install", "-y", "--set", "values.global.proxy.logLevel=debug", "-f", path.Join(SpecsFolder, p)}
	if err := istioctl(argss...); err != nil {
		return err
	}

	// Apply otel addons
	if err := kubectl("apply", "-f", path.Join(SpecsFolder, "otel/")); err != nil {
		return err
	}

	// Enable ambient mode on default namespace and wait Kiali for completion
	args = [][]string{
		{
			"label", "namespace", namespace, "istio.io/dataplane-mode=ambient",
		},
		{
			"-n",
			"istio-system",
			"wait",
			"--for=condition=Ready",
			"pod",
			"-l",
			"app=kiali",
			"--timeout",
			"300s",
		},
	}
	for _, a := range args {
		if err := kubectl(a...); err != nil {
			return err
		}
	}

	// Start Kiali dashboard
	if err := istioctl("dashboard", "kiali"); err != nil {
		return err
	}

	return nil
}
